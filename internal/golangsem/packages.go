package golangsem

import (
	"context"
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// capturedPackage is one checked package of captured trees.
type capturedPackage struct {
	pkg   *types.Package
	roots []ast.NodeRef
}

// checkedProgram is the checked package set. Supplied packages contribute the
// identities of the declarations captured code uses; their structure is never
// exported.
type checkedProgram struct {
	captured []capturedPackage // sorted by import path
	supplied []*types.Package  // sorted by import path
	imports  []importSite      // captured import declarations
}

func (p *checkedProgram) isSupplied(pkg *types.Package) bool {
	return pkg != nil && slices.Contains(p.supplied, pkg)
}

// closedImporter resolves only packages already checked from this request.
type closedImporter map[string]*types.Package

func (m closedImporter) Import(path string) (*types.Package, error) {
	if pkg := m[path]; pkg != nil {
		return pkg, nil
	}
	return nil, fmt.Errorf("package %q is not in the supplied environment", path)
}

// checkPackage type-checks one package beside the call and returns the
// context's error as soon as the context ends. go/types cannot be interrupted:
// an abandoned check keeps its goroutine until it ends by itself, and its
// result is dropped. A caller that is told the context ended must return
// without reading or writing anything it passed here, because the abandoned
// check may still read the files, the file set and the importer and write info.
func checkPackage(ctx context.Context, conf *types.Config, path string, fset *token.FileSet, files []*goast.File, info *types.Info) (*types.Package, error) {
	type outcome struct {
		pkg      *types.Package
		err      error
		panicked any
	}
	done := make(chan outcome, 1)
	go func() {
		var o outcome
		defer func() {
			// A panic of the checker belongs to the caller's goroutine, as it
			// did when the check ran there. Nobody receives it after an abandon.
			o.panicked = recover()
			done <- o
		}()
		o.pkg, o.err = conf.Check(path, fset, files, info)
	}()
	select {
	case o := <-done:
		if o.panicked != nil {
			panic(o.panicked)
		}
		return o.pkg, o.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// importSpec converts one captured import. Dot and blank imports stay outside the profile.
func (c *converter) importSpec(id ast.NodeID) (*goast.ImportSpec, error) {
	d := c.node(id)
	if d.Kind != s.KindImport || len(d.Children) != 0 {
		return nil, n.ErrUnsupported
	}
	path := attribute(d, s.AttributePath)
	if !n.ValidImportPath(path) {
		return nil, fmt.Errorf("%w: import path %q", n.ErrUnsupported, path)
	}
	spec := &goast.ImportSpec{Path: &goast.BasicLit{ValuePos: c.pos(id), Kind: token.STRING, Value: strconv.Quote(path)}}
	for _, a := range d.Attributes {
		if a.Name != s.AttributeAlias {
			continue
		}
		if a.Text == "." || a.Text == "_" || !token.IsIdentifier(a.Text) {
			return nil, fmt.Errorf("%w: import alias %q", n.ErrUnsupported, a.Text)
		}
		spec.Name = goast.NewIdent(a.Text)
		spec.Name.NamePos = c.pos(id)
	}
	c.sites = append(c.sites, importSite{At: ast.NodeRef{Tree: c.tree.ID, Node: id}, Path: path})
	return spec, nil
}

// environmentState is the validated environment: every tree's file path and
// every supplied package, checked before any captured code is judged.
type environmentState struct {
	module   string
	paths    map[ast.TreeID]string
	supplied map[string]*types.Package
	order    []*types.Package // sorted by import path
}

// prepareEnvironment validates the Files layout and checks the supplied
// packages in their own import order. Every failure here is an invalid request:
// it concerns the caller's environment, not the captured code.
func prepareEnvironment(ctx context.Context, fset *token.FileSet, environment n.ImportEnvironment, inputs []s.AST, captured *files.Context) (*environmentState, error) {
	if captured == nil {
		return nil, semantic.ErrInvalidRequest
	}
	filePaths := map[files.ID]files.Path{}
	for _, entry := range captured.Files.Entries {
		filePaths[entry.ID] = entry.Path
	}
	state := &environmentState{module: environment.Module, paths: map[ast.TreeID]string{}, supplied: map[string]*types.Package{}}
	for _, link := range captured.Links {
		state.paths[link.Tree] = string(filePaths[link.File])
	}
	for _, tree := range inputs {
		if p, ok := state.paths[tree.ID]; !ok || !strings.HasSuffix(p, ".go") {
			return nil, semantic.ErrInvalidRequest
		}
	}
	parsed := map[string][]*goast.File{}
	edges := map[string][]string{}
	var nodes []string
	for _, p := range environment.Packages {
		nodes = append(nodes, p.ImportPath)
		for _, f := range p.Files {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			file, err := parser.ParseFile(fset, p.ImportPath+"/"+f.Name, f.Text, parser.SkipObjectResolution)
			if err != nil {
				return nil, semantic.ErrInvalidRequest
			}
			parsed[p.ImportPath] = append(parsed[p.ImportPath], file)
			for _, spec := range file.Imports {
				dep, err := strconv.Unquote(spec.Path.Value)
				// Supplied packages close over the supplied environment only.
				if err != nil || !slices.ContainsFunc(environment.Packages, func(q n.Package) bool { return q.ImportPath == dep }) || !internalAllowed(p.ImportPath, dep) {
					return nil, semantic.ErrInvalidRequest
				}
				edges[p.ImportPath] = append(edges[p.ImportPath], dep)
			}
		}
	}
	order, cycle := importOrder(nodes, edges)
	if cycle != nil {
		return nil, semantic.ErrInvalidRequest
	}
	importer := closedImporter{}
	conf := types.Config{GoVersion: n.LanguageVersion, Sizes: types.SizesFor("gc", "amd64"), Importer: importer}
	for _, p := range order {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		for _, dep := range edges[p] {
			if importer[dep].Name() == "main" {
				return nil, semantic.ErrInvalidRequest
			}
		}
		pkg, err := checkPackage(ctx, &conf, p, fset, parsed[p], nil)
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		if err != nil {
			return nil, semantic.ErrInvalidRequest
		}
		importer[p] = pkg
		state.supplied[p] = pkg
		state.order = append(state.order, pkg)
	}
	sort.Slice(state.order, func(i, j int) bool { return state.order[i].Path() < state.order[j].Path() })
	return state, nil
}

type packageGroup struct {
	name  string
	files []*goast.File
	roots []ast.NodeRef
	paths []string
}

// checkLayout refuses, from file paths alone and in file-path order, captured
// files that are not part of a package under this profile. It runs before any
// tree is converted, so a layout refusal never waits behind a syntax limit.
func checkLayout(env *environmentState, inputs []s.AST) error {
	for _, tree := range inputs {
		filePath := env.paths[tree.ID]
		root := ast.NodeRef{Tree: tree.ID, Node: tree.Root}
		if ignoredByGo(filePath) {
			return &Failure{Err: fmt.Errorf("%w: %s is a test, ignored or vendored file", n.ErrUnsupported, filePath), At: root}
		}
		if !n.ValidImportPath(n.CapturedImportPath(env.module, filePath)) {
			return &Failure{Err: fmt.Errorf("%w: package directory of %s", n.ErrUnsupported, filePath), At: root}
		}
	}
	return nil
}

// checkCaptured judges the captured code against the validated environment.
// Unsupported layouts and imports are found before any type check, so a
// rejection never hides an unsupported construct.
func checkCaptured(ctx context.Context, fset *token.FileSet, env *environmentState, inputs []s.AST, goFiles []*goast.File, sites []importSite, info *types.Info, positions map[token.Pos]ast.NodeRef) (*checkedProgram, error) {
	groups := map[string]*packageGroup{}
	treePackage := map[ast.TreeID]string{}
	for i, tree := range inputs {
		filePath := env.paths[tree.ID]
		root := ast.NodeRef{Tree: tree.ID, Node: tree.Root}
		importPath := n.CapturedImportPath(env.module, filePath)
		g := groups[importPath]
		if g == nil {
			g = &packageGroup{name: goFiles[i].Name.Name}
			groups[importPath] = g
		}
		g.files = append(g.files, goFiles[i])
		g.roots = append(g.roots, root)
		g.paths = append(g.paths, filePath)
		treePackage[tree.ID] = importPath
	}
	for _, site := range sites {
		if groups[site.Path] == nil && env.supplied[site.Path] == nil {
			return nil, &Failure{Err: fmt.Errorf("%w: import %q is outside the supplied environment", n.ErrUnsupported, site.Path), At: site.At}
		}
	}
	// The go tool's rejections before type checking, in import-path and
	// file-path order: case-insensitive collisions and disagreeing package clauses.
	var groupPaths []string
	for p := range groups {
		groupPaths = append(groupPaths, p)
	}
	sort.Strings(groupPaths)
	folded := map[string]string{}
	for _, p := range groupPaths {
		g := groups[p]
		if other, ok := folded[strings.ToLower(p)]; ok {
			return nil, &Failure{Err: fmt.Errorf("case-insensitive import collision: %q and %q", other, p), At: g.roots[0], Invalid: true}
		}
		folded[strings.ToLower(p)] = p
		names := map[string]bool{}
		for i := range g.files {
			base := strings.ToLower(path.Base(g.paths[i]))
			if names[base] {
				return nil, &Failure{Err: fmt.Errorf("case-insensitive file name collision: %s", g.paths[i]), At: g.roots[i], Invalid: true}
			}
			names[base] = true
		}
	}
	for _, p := range groupPaths {
		g := groups[p]
		for i, f := range g.files {
			if f.Name.Name != g.name {
				return nil, &Failure{Err: fmt.Errorf("package %s; expected package %s", f.Name.Name, g.name), At: g.roots[i], Invalid: true}
			}
		}
	}
	edges := map[string][]string{}
	via := map[[2]string]ast.NodeRef{}
	for _, site := range sites {
		from := treePackage[site.At.Tree]
		if !internalAllowed(from, site.Path) {
			return nil, &Failure{Err: fmt.Errorf("use of internal package %q not allowed from %q", site.Path, from), At: site.At, Invalid: true}
		}
		name := ""
		if g := groups[site.Path]; g != nil {
			name = g.name
			edges[from] = append(edges[from], site.Path)
			if _, ok := via[[2]string{from, site.Path}]; !ok {
				via[[2]string{from, site.Path}] = site.At
			}
		} else {
			name = env.supplied[site.Path].Name()
		}
		if name == "main" {
			return nil, &Failure{Err: fmt.Errorf("import %q is a program, not an importable package", site.Path), At: site.At, Invalid: true}
		}
	}
	var nodes []string
	for p := range groups {
		nodes = append(nodes, p)
	}
	sort.Strings(nodes)
	order, cycle := importOrder(nodes, edges)
	if cycle != nil {
		for i := 0; i+1 < len(cycle); i++ {
			if at, ok := via[[2]string{cycle[i], cycle[i+1]}]; ok {
				return nil, &Failure{Err: fmt.Errorf("import cycle not allowed: %s", strings.Join(cycle, " -> ")), At: at, Invalid: true}
			}
		}
		return nil, semantic.ErrInvalidResult
	}
	importer := closedImporter{}
	for p, pkg := range env.supplied {
		importer[p] = pkg
	}
	conf := types.Config{GoVersion: n.LanguageVersion, Sizes: types.SizesFor("gc", "amd64"), Importer: importer}
	program := &checkedProgram{imports: sites, supplied: env.order}
	for _, p := range order {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		g := groups[p]
		pkg, err := checkPackage(ctx, &conf, p, fset, g.files, info)
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		if err != nil {
			return nil, typeFailure(err, positions)
		}
		importer[p] = pkg
		program.captured = append(program.captured, capturedPackage{pkg: pkg, roots: g.roots})
	}
	sort.Slice(program.captured, func(i, j int) bool { return program.captured[i].pkg.Path() < program.captured[j].pkg.Path() })
	return program, nil
}

// ignoredByGo reports the files outside a package under this profile: the go
// tool's test files, files and directories starting with "_" or ".", and
// testdata, plus vendor directories, whose import paths Go resolves differently.
func ignoredByGo(filePath string) bool {
	base := path.Base(filePath)
	if strings.HasSuffix(base, "_test.go") || strings.HasPrefix(base, "_") || strings.HasPrefix(base, ".") {
		return true
	}
	if dir := path.Dir(filePath); dir != "." {
		for _, element := range strings.Split(dir, "/") {
			if element == "testdata" || element == "vendor" || strings.HasPrefix(element, "_") || strings.HasPrefix(element, ".") {
				return true
			}
		}
	}
	return false
}

// internalAllowed applies Go's internal-package rule: a path with an "internal"
// element may be imported only from within the tree rooted at its parent.
func internalAllowed(importer, imported string) bool {
	elements := strings.Split(imported, "/")
	for i, element := range elements {
		if element != "internal" {
			continue
		}
		parent := strings.Join(elements[:i], "/")
		if importer != parent && !strings.HasPrefix(importer, parent+"/") {
			return false
		}
	}
	return true
}

// importOrder returns dependencies before their importers, or one import cycle
// as a path whose first and last element are the same package.
func importOrder(nodes []string, edges map[string][]string) ([]string, []string) {
	state := map[string]int{}
	var order, stack []string
	var cycle []string
	var visit func(string) bool
	visit = func(node string) bool {
		switch state[node] {
		case 2:
			return true
		case 1:
			i := slices.Index(stack, node)
			cycle = append(slices.Clone(stack[i:]), node)
			return false
		}
		state[node] = 1
		stack = append(stack, node)
		next := slices.Clone(edges[node])
		sort.Strings(next)
		for _, dep := range slices.Compact(next) {
			if !visit(dep) {
				return false
			}
		}
		stack = stack[:len(stack)-1]
		state[node] = 2
		order = append(order, node)
		return true
	}
	sorted := slices.Clone(nodes)
	sort.Strings(sorted)
	for _, node := range sorted {
		if !visit(node) {
			return nil, cycle
		}
	}
	return order, nil
}
