package report

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"

	capability "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
)

type authored struct {
	directive directive
	root      capability.Symbol
	document  string
}

type testDeclaration struct {
	id, pkg, name, variant, file string
	pos                          token.Pos
}

type caseAssociation struct {
	target    capability.Symbol
	marker    capability.Symbol
	directive directive
	test      testDeclaration
}

type annotations struct {
	operations map[capability.Symbol][]authored
	tests      []testDeclaration
	cases      []caseAssociation
}

func (p *program) fileImports(file *ast.File) (map[string]string, error) {
	imports := map[string]string{}
	for _, item := range file.Imports {
		path, err := strconv.Unquote(item.Path.Value)
		if err != nil {
			return nil, err
		}
		name := ""
		if item.Name != nil {
			name = item.Name.Name
		} else if pkg := p.packages[path]; pkg != nil {
			name = pkg.types.Name()
		} else {
			name = p.metadata[path].Name
		}
		// Production go list omits test-only standard dependencies. testing's
		// declared package name is part of the supported Go test signature.
		if name == "" && path == "testing" {
			name = "testing"
		}
		if name != "" && name != "_" && name != "." {
			imports[name] = path
		}
	}
	return imports, nil
}

func (p *program) resolveReference(reference, pkgPath string, imports map[string]string, external bool) (types.Object, error) {
	var symbol capability.Symbol
	var err error
	if strings.Contains(reference, "#") {
		symbol, err = capability.ParseSymbol(reference)
		if err != nil {
			return nil, err
		}
	} else {
		parts := strings.Split(reference, ".")
		symbol.Package = pkgPath
		if path, ok := imports[parts[0]]; ok {
			symbol.Package = path
			parts = parts[1:]
		} else if external {
			return nil, fmt.Errorf("external test references need an import qualifier or canonical identity: %s", reference)
		}
		switch len(parts) {
		case 1:
			symbol.Name = parts[0]
		case 2:
			symbol.Receiver, symbol.Name = parts[0], parts[1]
		default:
			return nil, fmt.Errorf("invalid symbol reference %q", reference)
		}
	}
	pkg := p.packages[symbol.Package]
	if pkg == nil {
		return nil, fmt.Errorf("symbol package is outside selected production source: %s", symbol.Package)
	}
	var obj types.Object
	if symbol.Receiver == "" {
		obj = pkg.types.Scope().Lookup(symbol.Name)
	} else {
		typ, ok := pkg.types.Scope().Lookup(symbol.Receiver).(*types.TypeName)
		if ok {
			obj = method(types.NewPointer(typ.Type()), symbol.Name)
		}
	}
	if obj == nil {
		return nil, fmt.Errorf("unresolved production symbol %s", symbol)
	}
	if fn, ok := obj.(*types.Func); ok {
		return fn.Origin(), nil
	}
	return obj, nil
}

func (p *program) readAnnotations() (*annotations, error) {
	result := &annotations{operations: map[capability.Symbol][]authored{}}
	paths := make([]string, 0, len(p.packages))
	for path := range p.packages {
		paths = append(paths, path)
	}
	slices.Sort(paths)
	for _, path := range paths {
		pkg := p.packages[path]
		for _, file := range pkg.files {
			imports, err := p.fileImports(file)
			if err != nil {
				return nil, err
			}
			attached := map[*ast.CommentGroup]bool{}
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				attached[fn.Doc] = true
				ds, err := rawDirectives(fn.Doc)
				if err != nil {
					return nil, p.at(fn.Pos(), err)
				}
				obj, _ := pkg.info.Defs[fn.Name].(*types.Func)
				for _, d := range ds {
					if d.name == "case" {
						return nil, p.at(fn.Pos(), fmt.Errorf("case requires a selected top-level Go test declaration"))
					}
					entry := authored{directive: d}
					if d.name == "grammar" {
						root, err := p.resolveReference(d.argument, path, imports, false)
						if err != nil {
							return nil, p.at(fn.Pos(), err)
						}
						switch root.(type) {
						case *types.Var, *types.Const, *types.TypeName:
						default:
							return nil, p.at(fn.Pos(), fmt.Errorf("grammar root must be a definition, not an operation"))
						}
						entry.root = canonicalSymbol(root)
					} else {
						doc, err := scopeDocument(p.root, d.argument)
						if err != nil {
							return nil, p.at(fn.Pos(), err)
						}
						entry.document = doc
					}
					key := canonicalSymbol(obj.Origin())
					for _, previous := range result.operations[key] {
						if previous.directive.name == d.name && previous.directive.profile == d.profile && previous.root == entry.root && previous.document == entry.document {
							return nil, p.at(fn.Pos(), fmt.Errorf("duplicate %s association", d.name))
						}
					}
					result.operations[key] = append(result.operations[key], entry)
				}
			}
			for _, group := range file.Comments {
				if attached[group] {
					continue
				}
				ds, err := rawDirectives(group)
				if err != nil {
					return nil, p.at(group.Pos(), err)
				}
				if len(ds) > 0 {
					return nil, p.at(group.Pos(), fmt.Errorf("capability directives require an operation declaration doc comment"))
				}
			}
		}
	}
	// Go's selected test files are a separate declaration index. They cannot add
	// candidates or replace objects in the production type graph.
	paths = paths[:0]
	for path, item := range p.metadata {
		if item.Module != nil && item.Module.Main && (path == p.module || strings.HasPrefix(path, p.module+"/")) && (len(item.TestGoFiles) > 0 || len(item.XTestGoFiles) > 0) {
			paths = append(paths, path)
		}
	}
	slices.Sort(paths)
	for _, path := range paths {
		item := p.metadata[path]
		for _, variant := range []struct {
			name     string
			files    []string
			external bool
		}{{"internal", item.TestGoFiles, false}, {"external", item.XTestGoFiles, true}} {
			for _, name := range variant.files {
				filename := filepath.Join(item.Dir, name)
				data, err := p.readSource(filename)
				if err != nil {
					return nil, err
				}
				file, err := parser.ParseFile(p.fset, filename, data, parser.ParseComments)
				if err != nil {
					return nil, err
				}
				imports, err := p.fileImports(file)
				if err != nil {
					return nil, err
				}
				attached := map[*ast.CommentGroup]bool{}
				for _, decl := range file.Decls {
					fn, ok := decl.(*ast.FuncDecl)
					if !ok {
						continue
					}
					attached[fn.Doc] = true
					ds, err := rawDirectives(fn.Doc)
					if err != nil {
						return nil, p.at(fn.Pos(), err)
					}
					valid := testFunction(fn, imports)
					if len(ds) > 0 && !valid {
						return nil, p.at(fn.Pos(), fmt.Errorf("case requires direct top-level func TestName(*testing.T) with no results"))
					}
					if !valid {
						continue
					}
					rel, _ := filepath.Rel(p.root, filename)
					test := testDeclaration{id: stableKey(path, variant.name, filepath.ToSlash(rel), fn.Name.Name), pkg: path, name: fn.Name.Name, variant: variant.name, file: filepath.ToSlash(rel), pos: fn.Name.Pos()}
					result.tests = append(result.tests, test)
					for _, d := range ds {
						if d.name != "case" {
							return nil, p.at(fn.Pos(), fmt.Errorf("test documentation accepts only case directives"))
						}
						obj, err := p.resolveReference(d.values["target"], path, imports, variant.external)
						if err != nil {
							return nil, p.at(fn.Pos(), err)
						}
						if _, ok := obj.(*types.Func); !ok {
							return nil, p.at(fn.Pos(), fmt.Errorf("case target must be a production operation"))
						}
						assoc := caseAssociation{target: canonicalSymbol(obj), directive: d, test: test}
						if value := d.values["language"]; value != "" {
							marker, err := p.resolveReference(value, path, imports, variant.external)
							if err != nil {
								return nil, p.at(fn.Pos(), err)
							}
								typeName, ok := marker.(*types.TypeName)
								if !ok {
								return nil, p.at(fn.Pos(), fmt.Errorf("language qualifier must name a production marker type"))
							}
								named, ok := types.Unalias(typeName.Type()).(*types.Named)
								if !ok {
									return nil, p.at(fn.Pos(), fmt.Errorf("language qualifier must denote a named marker, not a pointer or composite type"))
								}
								assoc.marker = canonicalSymbol(named.Obj())
						}
						for _, prev := range result.cases {
							if prev.target == assoc.target && prev.marker == assoc.marker && prev.directive.profile == d.profile && prev.directive.values["feature"] == d.values["feature"] && prev.directive.values["case"] == d.values["case"] && prev.test.id == test.id {
								return nil, p.at(fn.Pos(), fmt.Errorf("duplicate or conflicting case association"))
							}
						}
						result.cases = append(result.cases, assoc)
					}
				}
				for _, group := range file.Comments {
					if attached[group] {
						continue
					}
					ds, err := rawDirectives(group)
					if err != nil {
						return nil, p.at(group.Pos(), err)
					}
					if len(ds) > 0 {
						return nil, p.at(group.Pos(), fmt.Errorf("case must be attached to a top-level test declaration"))
					}
				}
			}
		}
	}
	return result, nil
}

func testFunction(fn *ast.FuncDecl, imports map[string]string) bool {
	if fn.Recv != nil || fn.Type.TypeParams != nil || fn.Type.Results != nil || !strings.HasPrefix(fn.Name.Name, "Test") {
		return false
	}
	if suffix := strings.TrimPrefix(fn.Name.Name, "Test"); suffix != "" {
		r := []rune(suffix)[0]
		if unicode.IsLower(r) {
			return false
		}
	}
	if fn.Type.Params == nil || len(fn.Type.Params.List) != 1 {
		return false
	}
	parameter := fn.Type.Params.List[0]
	if len(parameter.Names) > 1 {
		return false
	}
	unwrap := func(e ast.Expr) ast.Expr {
		for {
			p, ok := e.(*ast.ParenExpr)
			if !ok {
				return e
			}
			e = p.X
		}
	}
	ptr, ok := unwrap(parameter.Type).(*ast.StarExpr)
	if !ok {
		return false
	}
	selector, ok := unwrap(ptr.X).(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "T" {
		return false
	}
	id, ok := selector.X.(*ast.Ident)
	return ok && imports[id.Name] == "testing"
}

func (p *program) at(pos token.Pos, err error) error {
	return fmt.Errorf("%s: %w", p.fset.PositionFor(pos, false), err)
}
