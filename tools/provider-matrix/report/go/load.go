package report

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/Bitspark/bittype/internal/processoutput"
)

// The Go command supplies the active build context and dependency export data.
// Repository packages are read from source so symbols retain their actual positions.
// No repository initializer, constructor or provider method is executed.
type packageInfo struct {
	ImportPath, Name, Dir, Export                        string
	GoFiles, CgoFiles, Imports                           []string
	TestGoFiles, XTestGoFiles, TestImports, XTestImports []string
	Module                                               *struct {
		Path, Dir, GoVersion string
		Main                 bool
	}
	Error *struct{ Err string }
}

type sourcePackage struct {
	types   *types.Package
	info    *types.Info
	files   []*ast.File
	sources map[string]string
}

type declaration struct {
	pkg  *sourcePackage
	node *ast.FuncDecl
}

type program struct {
	root, module string
	fset         *token.FileSet
	packages     map[string]*sourcePackage
	functions    map[*types.Func]declaration
	dimensions   map[*types.TypeName]dimensionMetadata
	metadata     map[string]packageInfo
	ctx          context.Context
	limits       LoadLimits
	sources      map[string]string
	sourceBytes  int
	build        goEnvironment
	resolvedRoot string
}

type goEnvironment struct {
	GOOS, GOARCH, GOFLAGS, GOVERSION string
}

type sourceImporter struct {
	program  *program
	metadata map[string]packageInfo
	external types.Importer
	loading  map[string]bool
}

// LoadLimits bounds package metadata and captured repository source. Zero fields
// use the defaults. The caller's context bounds the Go command; in-process Go
// type checking is cooperative between packages, not forcibly preempted.
type LoadLimits struct {
	MaxPackages, MaxFiles, MaxSourceBytes, MaxMetadataBytes int
}

func DefaultLoadLimits() LoadLimits {
	return LoadLimits{MaxPackages: 4096, MaxFiles: 32768, MaxSourceBytes: 128 << 20, MaxMetadataBytes: 64 << 20}
}

func (l LoadLimits) effective() (LoadLimits, error) {
	d := DefaultLoadLimits()
	for _, pair := range [][2]*int{{&l.MaxPackages, &d.MaxPackages}, {&l.MaxFiles, &d.MaxFiles}, {&l.MaxSourceBytes, &d.MaxSourceBytes}, {&l.MaxMetadataBytes, &d.MaxMetadataBytes}} {
		if *pair[0] < 0 {
			return LoadLimits{}, fmt.Errorf("load limits must not be negative")
		}
		if *pair[0] == 0 {
			*pair[0] = *pair[1]
		}
	}
	return l, nil
}

// Only build-selection flags supported by this source loader are accepted.
// Overlays and alternate package operands would invalidate the captured source
// view; -test would introduce a different compiler identity world.
func validateBuildFlags(flags []string) error {
	for i := 0; i < len(flags); i++ {
		flag := flags[i]
		switch flag {
		case "-race", "-msan", "-asan", "-trimpath":
			continue
		case "-tags", "-mod":
			if i+1 == len(flags) || strings.HasPrefix(flags[i+1], "-") || flags[i+1] == "" {
				return fmt.Errorf("missing value for build flag %s", flag)
			}
			i++
			continue
		}
		if strings.HasPrefix(flag, "-tags=") && len(flag) > len("-tags=") || strings.HasPrefix(flag, "-mod=") && len(flag) > len("-mod=") {
			continue
		}
		return fmt.Errorf("unsupported source-discovery build flag %q", flag)
	}
	return nil
}

func load(ctx context.Context, root string, limits LoadLimits, buildFlags []string) (*program, error) {
	if ctx == nil {
		return nil, fmt.Errorf("load context is required")
	}
	var err error
	limits, err = limits.effective()
	if err != nil {
		return nil, err
	}
	if err := validateBuildFlags(buildFlags); err != nil {
		return nil, err
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	commandCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	// One bound for the output streams and the error streams of both commands.
	budget := processoutput.NewBound(limits.MaxMetadataBytes, cancel)
	runGo := func(args ...string) ([]byte, error) {
		command := exec.CommandContext(commandCtx, "go", args...)
		command.Dir = root
		command.WaitDelay = 2 * time.Second
		stdout, stderr := budget.Writer(), budget.Writer()
		command.Stdout, command.Stderr = stdout, stderr
		err := command.Run()
		if budget.Exceeded() {
			return nil, fmt.Errorf("Go command output exceeds %d bytes", limits.MaxMetadataBytes)
		}
		if err != nil {
			return nil, fmt.Errorf("Go command: %w\n%s", errors.Join(err, ctx.Err()), stderr.Bytes())
		}
		return stdout.Bytes(), nil
	}
	environment, err := runGo("env", "-json", "GOOS", "GOARCH", "GOFLAGS", "GOVERSION")
	if err != nil {
		return nil, fmt.Errorf("load Go build context: %w", err)
	}
	var build goEnvironment
	if err := json.Unmarshal(environment, &build); err != nil {
		return nil, fmt.Errorf("decode Go build context: %w", err)
	}
	if strings.ContainsAny(build.GOFLAGS, "\"'") {
		return nil, fmt.Errorf("quoted GOFLAGS are unsupported by source discovery; use explicit build flags")
	}
	if err := validateBuildFlags(strings.Fields(build.GOFLAGS)); err != nil {
		return nil, fmt.Errorf("ambient GOFLAGS: %w", err)
	}
	if build.GOOS == "" || types.SizesFor("gc", build.GOARCH) == nil {
		return nil, fmt.Errorf("unsupported effective Go target %s/%s", build.GOOS, build.GOARCH)
	}
	args := append([]string{"list", "-deps", "-export", "-json"}, buildFlags...)
	args = append(args, "./...")
	data, err := runGo(args...)
	if err != nil {
		return nil, fmt.Errorf("load Go package metadata: %w", err)
	}
	metadata := map[string]packageInfo{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	module := ""
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		var item packageInfo
		if err := decoder.Decode(&item); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if item.Error != nil {
			return nil, fmt.Errorf("%s: %s", item.ImportPath, item.Error.Err)
		}
		if len(metadata) >= limits.MaxPackages {
			return nil, fmt.Errorf("Go package count exceeds %d", limits.MaxPackages)
		}
		metadata[item.ImportPath] = item
		if item.Module != nil && item.Module.Main && filepath.Clean(item.Module.Dir) == root {
			module = item.Module.Path
		}
	}
	if module == "" {
		return nil, fmt.Errorf("no main Go module at %s", root)
	}
	p := &program{root: root, module: module, fset: token.NewFileSet(), packages: map[string]*sourcePackage{}, functions: map[*types.Func]declaration{}, metadata: metadata, ctx: ctx, limits: limits, sources: map[string]string{}, build: build}
	loader := &sourceImporter{program: p, metadata: metadata, loading: map[string]bool{}}
	loader.external = importer.ForCompiler(p.fset, "gc", func(path string) (io.ReadCloser, error) {
		item, ok := metadata[path]
		if !ok || item.Export == "" {
			return nil, fmt.Errorf("no export data for %s", path)
		}
		return os.Open(item.Export)
	})
	paths := make([]string, 0, len(metadata))
	for path := range metadata {
		paths = append(paths, path)
	}
	slices.Sort(paths)
	for _, path := range paths {
		item := metadata[path]
		if strings.HasPrefix(path, module+"/") && len(item.GoFiles) > 0 && !strings.HasPrefix(path, module+"/tools/") {
			if _, err := loader.Import(path); err != nil {
				return nil, err
			}
		}
	}
	return p, nil
}

func (p *program) readSource(path string) (string, error) {
	if p.ctx != nil {
		if err := p.ctx.Err(); err != nil {
			return "", err
		}
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(p.root, path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("repository source is outside root: %s", path)
	}
	if source, ok := p.sources[path]; ok {
		return source, nil
	}
	if len(p.sources) >= p.limits.MaxFiles {
		return "", fmt.Errorf("source file count exceeds %d", p.limits.MaxFiles)
	}
	if p.resolvedRoot == "" {
		p.resolvedRoot, err = filepath.EvalSymlinks(p.root)
		if err != nil {
			return "", err
		}
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	if !contained(p.resolvedRoot, resolved) {
		return "", fmt.Errorf("repository source symlink resolves outside root: %s", path)
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	remaining := p.limits.MaxSourceBytes - p.sourceBytes
	data, err := io.ReadAll(io.LimitReader(file, int64(remaining)+1))
	if err != nil {
		return "", err
	}
	if len(data) > remaining {
		return "", fmt.Errorf("source bytes exceed %d", p.limits.MaxSourceBytes)
	}
	p.sourceBytes += len(data)
	source := string(data)
	p.sources[path] = source
	return source, nil
}

func (l *sourceImporter) Import(path string) (*types.Package, error) {
	p := l.program
	if err := p.ctx.Err(); err != nil {
		return nil, err
	}
	if loaded := p.packages[path]; loaded != nil {
		return loaded.types, nil
	}
	if !strings.HasPrefix(path, p.module+"/") {
		return l.external.Import(path)
	}
	if l.loading[path] {
		return nil, fmt.Errorf("import cycle at %s", path)
	}
	l.loading[path] = true
	defer delete(l.loading, path)
	item, ok := l.metadata[path]
	if !ok {
		return nil, fmt.Errorf("missing package %s", path)
	}
	if len(item.CgoFiles) != 0 {
		return nil, fmt.Errorf("source discovery does not support cgo package %s", path)
	}
	sp := &sourcePackage{sources: map[string]string{}, info: &types.Info{
		Types: map[ast.Expr]types.TypeAndValue{}, Defs: map[*ast.Ident]types.Object{},
		Uses: map[*ast.Ident]types.Object{}, Selections: map[*ast.SelectorExpr]*types.Selection{},
		Instances: map[*ast.Ident]types.Instance{},
	}}
	for _, name := range item.GoFiles {
		filename := filepath.Join(item.Dir, name)
		source, err := p.readSource(filename)
		if err != nil {
			return nil, err
		}
		file, err := parser.ParseFile(p.fset, filename, source, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		sp.files = append(sp.files, file)
		sp.sources[filename] = source
	}
	config := types.Config{Importer: l, Sizes: types.SizesFor("gc", p.build.GOARCH)}
	if item.Module != nil && item.Module.GoVersion != "" {
		config.GoVersion = "go" + item.Module.GoVersion
	}
	var err error
	sp.types, err = config.Check(path, p.fset, sp.files, sp.info)
	if err != nil {
		return nil, fmt.Errorf("read Go types: %w", err)
	}
	p.packages[path] = sp
	if err := p.readDimensions(sp); err != nil {
		return nil, err
	}
	for _, file := range sp.files {
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if obj, ok := sp.info.Defs[fn.Name].(*types.Func); ok {
					p.functions[obj.Origin()] = declaration{sp, fn}
				}
			}
		}
	}
	return sp.types, nil
}

func (p *program) implementation(path string) bool {
	rel := strings.TrimPrefix(path, p.module+"/")
	if rel == path {
		return false
	}
	// Architectural infrastructure is not a language implementation, even when an
	// adapter's method set has the same shape as the provider that it wraps.
	for _, dir := range []string{"contracts", "trees", "engine", "internal", "integration", "tools"} {
		if strings.HasPrefix(rel, dir+"/") {
			return false
		}
	}
	return true
}

func (p *program) named(path, name string) (*types.Named, error) {
	pkg := p.packages[p.module+"/"+path]
	if pkg == nil {
		return nil, fmt.Errorf("missing package %s", path)
	}
	obj := pkg.types.Scope().Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("missing %s.%s", path, name)
	}
	named, ok := types.Unalias(obj.Type()).(*types.Named)
	if !ok {
		return nil, fmt.Errorf("%s.%s is not a named type", path, name)
	}
	return named, nil
}

func method(t types.Type, name string) *types.Func {
	selection := types.NewMethodSet(t).Lookup(nil, name)
	if selection == nil {
		return nil
	}
	fn, _ := selection.Obj().(*types.Func)
	return fn
}

func (p *program) called(pkg *sourcePackage, expr ast.Expr) *types.Func {
	switch e := expr.(type) {
	case *ast.ParenExpr:
		return p.called(pkg, e.X)
	case *ast.IndexExpr:
		return p.called(pkg, e.X)
	case *ast.IndexListExpr:
		return p.called(pkg, e.X)
	case *ast.Ident:
		fn, _ := pkg.info.Uses[e].(*types.Func)
		return fn
	case *ast.SelectorExpr:
		if selection := pkg.info.Selections[e]; selection != nil {
			fn, _ := selection.Obj().(*types.Func)
			return fn
		}
		fn, _ := pkg.info.Uses[e.Sel].(*types.Func)
		return fn
	}
	return nil
}

// A single-return static call identifies a related implementation entry, not a
// transparent wrapper. Arguments may be rewritten and results may have behavior
// not established by this source walk; observations never transfer along it.
func (p *program) entryChain(fn *types.Func) (*types.Func, []callLink, string) {
	seen := map[*types.Func]bool{}
	var chain []callLink
	for fn != nil {
		fn = fn.Origin()
		if seen[fn] {
			return fn, chain, "cycle"
		}
		seen[fn] = true
		decl, ok := p.functions[fn]
		if !ok || decl.node.Body == nil || len(decl.node.Body.List) != 1 {
			return fn, chain, "not a single-return static call"
		}
		ret, ok := decl.node.Body.List[0].(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			return fn, chain, "not a single-return static call"
		}
		call, ok := ret.Results[0].(*ast.CallExpr)
		if !ok {
			return fn, chain, "not a single-return static call"
		}
		next := p.called(decl.pkg, call.Fun)
		if next == nil || next.Pkg() == nil || !p.implementation(next.Pkg().Path()) {
			return fn, chain, "callee is dynamic or outside implementation source"
		}
		if _, ok := p.functions[next.Origin()]; !ok {
			return fn, chain, "callee source unavailable"
		}
		chain = append(chain, callLink{caller: fn, callee: next.Origin(), pos: call.Pos(), receiver: callReceiver(decl.pkg, call.Fun), arguments: callArguments(decl.pkg, call.Fun)})
		fn = next
	}
	return fn, chain, "no callable entry"
}

func callReceiver(pkg *sourcePackage, expr ast.Expr) types.Type {
	switch expr := expr.(type) {
	case *ast.ParenExpr:
		return callReceiver(pkg, expr.X)
	case *ast.IndexExpr:
		return callReceiver(pkg, expr.X)
	case *ast.IndexListExpr:
		return callReceiver(pkg, expr.X)
	case *ast.SelectorExpr:
		if selection := pkg.info.Selections[expr]; selection != nil {
			return selection.Recv()
		}
	}
	return nil
}

func callArguments(pkg *sourcePackage, expr ast.Expr) []types.Type {
	var ident *ast.Ident
	switch expr := expr.(type) {
	case *ast.ParenExpr:
		return callArguments(pkg, expr.X)
	case *ast.IndexExpr:
		return callArguments(pkg, expr.X)
	case *ast.IndexListExpr:
		return callArguments(pkg, expr.X)
	case *ast.Ident:
		ident = expr
	case *ast.SelectorExpr:
		ident = expr.Sel
	}
	instance := pkg.info.Instances[ident]
	var result []types.Type
	for i := 0; i < instance.TypeArgs.Len(); i++ {
		result = append(result, instance.TypeArgs.At(i))
	}
	return result
}
