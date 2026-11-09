package golangsem

import (
	"context"
	"fmt"
	goast "go/ast"
	"go/parser"
	"go/token"
	"go/types"
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
// identities of the declarations captured code uses; their internals are never
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

type packageGroup struct {
	files []*goast.File
	roots []ast.NodeRef
}

// checkPackages checks captured packages and the supplied environment in import
// order. Captured type errors and import cycles through captured packages reject
// validity; imports outside the environment are unsupported coverage; a broken
// supplied package is an invalid request.
func checkPackages(ctx context.Context, fset *token.FileSet, environment n.ImportEnvironment, inputs []s.AST, goFiles []*goast.File, sites []importSite, captured *files.Context, info *types.Info, positions map[token.Pos]ast.NodeRef) (*checkedProgram, error) {
	if captured == nil {
		return nil, semantic.ErrInvalidRequest
	}
	filePaths := map[files.ID]files.Path{}
	for _, entry := range captured.Files.Entries {
		filePaths[entry.ID] = entry.Path
	}
	treePaths := map[ast.TreeID]string{}
	for _, link := range captured.Links {
		treePaths[link.Tree] = string(filePaths[link.File])
	}
	groups := map[string]*packageGroup{}
	treePackage := map[ast.TreeID]string{}
	for i, tree := range inputs {
		filePath, ok := treePaths[tree.ID]
		if !ok || !strings.HasSuffix(filePath, ".go") {
			return nil, semantic.ErrInvalidRequest
		}
		root := ast.NodeRef{Tree: tree.ID, Node: tree.Root}
		importPath := n.CapturedImportPath(environment.Module, filePath)
		if !n.ValidImportPath(importPath) {
			return nil, &Failure{Err: fmt.Errorf("%w: package directory of %s", n.ErrUnsupported, filePath), At: root}
		}
		g := groups[importPath]
		if g == nil {
			g = &packageGroup{}
			groups[importPath] = g
		}
		g.files = append(g.files, goFiles[i])
		g.roots = append(g.roots, root)
		treePackage[tree.ID] = importPath
	}
	supplied := map[string]n.Package{}
	for _, p := range environment.Packages {
		supplied[p.ImportPath] = p
	}
	parsed := map[string][]*goast.File{}
	edges := map[string][]string{}
	for _, p := range environment.Packages {
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
				path, err := strconv.Unquote(spec.Path.Value)
				if err != nil || supplied[path].ImportPath == "" {
					// Supplied packages close over the supplied environment only.
					return nil, semantic.ErrInvalidRequest
				}
				edges[p.ImportPath] = append(edges[p.ImportPath], path)
			}
		}
	}
	via := map[[2]string]ast.NodeRef{}
	for _, site := range sites {
		from := treePackage[site.At.Tree]
		if groups[site.Path] == nil && supplied[site.Path].ImportPath == "" {
			return nil, &Failure{Err: fmt.Errorf("%w: import %q is outside the supplied environment", n.ErrUnsupported, site.Path), At: site.At}
		}
		edges[from] = append(edges[from], site.Path)
		if _, ok := via[[2]string{from, site.Path}]; !ok {
			via[[2]string{from, site.Path}] = site.At
		}
	}
	var nodes []string
	for path := range groups {
		nodes = append(nodes, path)
	}
	for path := range supplied {
		nodes = append(nodes, path)
	}
	sort.Strings(nodes)
	order, cycle := importOrder(nodes, edges)
	if cycle != nil {
		for i := 0; i+1 < len(cycle); i++ {
			if at, ok := via[[2]string{cycle[i], cycle[i+1]}]; ok {
				return nil, &Failure{Err: fmt.Errorf("import cycle not allowed: %s", strings.Join(cycle, " -> ")), At: at, Invalid: true}
			}
		}
		return nil, semantic.ErrInvalidRequest
	}
	importer := closedImporter{}
	conf := types.Config{GoVersion: n.LanguageVersion, Sizes: types.SizesFor("gc", "amd64"), Importer: importer}
	program := &checkedProgram{imports: sites}
	for _, path := range order {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if _, ok := supplied[path]; ok {
			pkg, err := conf.Check(path, fset, parsed[path], nil)
			if err != nil {
				return nil, semantic.ErrInvalidRequest
			}
			importer[path] = pkg
			program.supplied = append(program.supplied, pkg)
			continue
		}
		g := groups[path]
		pkg, err := conf.Check(path, fset, g.files, info)
		if err != nil {
			return nil, typeFailure(err, positions)
		}
		importer[path] = pkg
		program.captured = append(program.captured, capturedPackage{pkg: pkg, roots: g.roots})
	}
	sort.Slice(program.captured, func(i, j int) bool { return program.captured[i].pkg.Path() < program.captured[j].pkg.Path() })
	sort.Slice(program.supplied, func(i, j int) bool { return program.supplied[i].Path() < program.supplied[j].Path() })
	return program, nil
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
	for _, node := range nodes {
		if !visit(node) {
			return nil, cycle
		}
	}
	return order, nil
}
