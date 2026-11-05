package main

import (
	"bytes"
	"encoding/json"
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
	"runtime"
	"strings"
)

// The Go command supplies the active build context and dependency export data.
// Repository packages are read from source so symbols retain their actual positions.
// No repository initializer, constructor or provider method is executed.
type packageInfo struct {
	ImportPath, Dir, Export string
	GoFiles, CgoFiles       []string
	Module                  *struct {
		Path, GoVersion string
		Main            bool
	}
	Error *struct{ Err string }
}

type sourcePackage struct {
	types *types.Package
	info  *types.Info
	files []*ast.File
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
}

type sourceImporter struct {
	program  *program
	metadata map[string]packageInfo
	external types.Importer
	loading  map[string]bool
}

func load(root string) (*program, error) {
	command := exec.Command("go", "list", "-deps", "-export", "-json", "./...")
	command.Dir = root
	var stderr bytes.Buffer
	command.Stderr = &stderr
	data, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("load Go package metadata: %w\n%s", err, stderr.String())
	}
	metadata := map[string]packageInfo{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	module := ""
	for {
		var item packageInfo
		if err := decoder.Decode(&item); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if item.Error != nil {
			return nil, fmt.Errorf("%s: %s", item.ImportPath, item.Error.Err)
		}
		metadata[item.ImportPath] = item
		if item.Module != nil && item.Module.Main {
			module = item.Module.Path
		}
	}
	if module == "" {
		return nil, fmt.Errorf("no main Go module at %s", root)
	}
	p := &program{root: root, module: module, fset: token.NewFileSet(), packages: map[string]*sourcePackage{}, functions: map[*types.Func]declaration{}}
	loader := &sourceImporter{program: p, metadata: metadata, loading: map[string]bool{}}
	loader.external = importer.ForCompiler(p.fset, "gc", func(path string) (io.ReadCloser, error) {
		item, ok := metadata[path]
		if !ok || item.Export == "" {
			return nil, fmt.Errorf("no export data for %s", path)
		}
		return os.Open(item.Export)
	})
	for path, item := range metadata {
		if strings.HasPrefix(path, module+"/") && len(item.GoFiles) > 0 && !strings.HasPrefix(path, module+"/tools/") {
			if _, err := loader.Import(path); err != nil {
				return nil, err
			}
		}
	}
	return p, nil
}

func (l *sourceImporter) Import(path string) (*types.Package, error) {
	p := l.program
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
	sp := &sourcePackage{info: &types.Info{
		Types: map[ast.Expr]types.TypeAndValue{}, Defs: map[*ast.Ident]types.Object{},
		Uses: map[*ast.Ident]types.Object{}, Selections: map[*ast.SelectorExpr]*types.Selection{},
		Instances: map[*ast.Ident]types.Instance{},
	}}
	for _, name := range item.GoFiles {
		file, err := parser.ParseFile(p.fset, filepath.Join(item.Dir, name), nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		sp.files = append(sp.files, file)
	}
	config := types.Config{Importer: l, Sizes: types.SizesFor(runtime.Compiler, runtime.GOARCH)}
	if item.Module != nil {
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

// Only transparent, single-return forwarding is followed. Dynamic dispatch and
// arbitrary call graphs do not establish a different implementation location.
func (p *program) entry(fn *types.Func) *types.Func {
	seen := map[*types.Func]bool{}
	for fn != nil {
		fn = fn.Origin()
		if seen[fn] {
			break
		}
		seen[fn] = true
		decl, ok := p.functions[fn]
		if !ok || decl.node.Body == nil || len(decl.node.Body.List) != 1 {
			break
		}
		ret, ok := decl.node.Body.List[0].(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			break
		}
		call, ok := ret.Results[0].(*ast.CallExpr)
		if !ok {
			break
		}
		next := p.called(decl.pkg, call.Fun)
		if next == nil || next.Pkg() == nil || !p.implementation(next.Pkg().Path()) {
			break
		}
		if _, ok := p.functions[next.Origin()]; !ok {
			break
		}
		fn = next
	}
	return fn
}
