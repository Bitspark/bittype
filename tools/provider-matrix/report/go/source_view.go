package report

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"

	capability "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
)

// canonicalSymbol accepts package declarations, named-receiver methods and Go
// builtins. Locals, fields and anonymous receivers cannot acquire a package
// identity just because their spelling happens to match one.
func canonicalSymbol(obj types.Object) capability.Symbol {
	if obj == nil {
		return capability.Symbol{}
	}
	if fn, ok := obj.(*types.Func); ok {
		if fn == nil {
			return capability.Symbol{}
		}
		fn = fn.Origin()
		if receiver := fn.Type().(*types.Signature).Recv(); receiver != nil {
			named := namedType(receiver.Type())
			if named == nil || named.Obj().Pkg() == nil {
				return capability.Symbol{}
			}
			return capability.Symbol{Package: named.Obj().Pkg().Path(), Receiver: named.Origin().Obj().Name(), Name: fn.Name()}
		}
		obj = fn
	}
	if obj.Pkg() == nil {
		if obj.Parent() == types.Universe {
			return capability.Symbol{Package: "builtin", Name: obj.Name()}
		}
		return capability.Symbol{}
	}
	if obj.Parent() != obj.Pkg().Scope() {
		return capability.Symbol{}
	}
	return capability.Symbol{Package: obj.Pkg().Path(), Name: obj.Name()}
}

func symbolPointer(obj types.Object) *capability.Symbol {
	symbol := canonicalSymbol(obj)
	if symbol.Package == "" || symbol.Name == "" {
		return nil
	}
	return &symbol
}

func sourceTypeFact(t types.Type) capability.TypeFact {
	if t == nil {
		return capability.TypeFact{}
	}
	result := capability.TypeFact{Text: fullType(t)}
	if named, ok := types.Unalias(t).(*types.Named); ok {
		result.Origin = symbolPointer(named.Origin().Obj())
		for i := 0; i < named.TypeArgs().Len(); i++ {
			result.Arguments = append(result.Arguments, sourceTypeFact(named.TypeArgs().At(i)))
		}
	}
	return result
}

type sourceView struct {
	fset         *token.FileSet
	files        map[capability.FileRef]capability.File
	paths        map[string]capability.FileRef
	declarations map[capability.Symbol]capability.Declaration
	expressions  map[capability.Span]expressionFacts
}

type expressionFacts struct {
	typ            types.Type
	symbol, callee types.Object
	value          constant.Value
}

var sourceSequence atomic.Uint64

// newSourceView snapshots strings and compiler facts. Extractors never receive
// AST nodes, go/types objects, maps or slices owned by the compiler/host.
func newSourceView(p *program) *sourceView {
	v := &sourceView{fset: p.fset, files: map[capability.FileRef]capability.File{}, paths: map[string]capability.FileRef{}, declarations: map[capability.Symbol]capability.Declaration{}, expressions: map[capability.Span]expressionFacts{}}
	sources := map[string]string{}
	for path, text := range p.sources {
		sources[filepath.Clean(path)] = text
	}
	for _, pkg := range p.packages {
		for path, text := range pkg.sources {
			sources[filepath.Clean(path)] = text
		}
	}
	paths := make([]string, 0, len(sources))
	for path := range sources {
		paths = append(paths, path)
	}
	slices.Sort(paths)
	viewID := sourceSequence.Add(1)
	for i, path := range paths {
		ref := capability.FileRef(fmt.Sprintf("source-%d/file-%d", viewID, i))
		v.paths[path] = ref
		v.files[ref] = capability.File{Ref: ref, Path: path, Text: sources[path]}
	}
	for _, pkg := range p.packages {
		for _, file := range pkg.files {
			for _, decl := range file.Decls {
				switch node := decl.(type) {
				case *ast.FuncDecl:
					v.addDeclaration(pkg.info.Defs[node.Name], node, nil, node.Type)
				case *ast.GenDecl:
					for _, spec := range node.Specs {
						switch spec := spec.(type) {
						case *ast.TypeSpec:
							v.addDeclaration(pkg.info.Defs[spec.Name], node, nil, spec.Type)
						case *ast.ValueSpec:
							for i, ident := range spec.Names {
								var value ast.Expr
								if len(spec.Names) == len(spec.Values) {
									value = spec.Values[i]
								}
								v.addDeclaration(pkg.info.Defs[ident], node, value, spec.Type)
							}
						}
					}
				}
			}
			ast.Inspect(file, func(node ast.Node) bool {
				expr, ok := node.(ast.Expr)
				if !ok {
					return true
				}
				span := v.span(expr)
				if span.File == "" {
					return true
				}
				facts := expressionFacts{typ: pkg.info.TypeOf(expr), symbol: expressionObject(pkg, expr), value: pkg.info.Types[expr].Value}
				if call, ok := expr.(*ast.CallExpr); ok {
					facts.callee = expressionObject(pkg, call.Fun)
				}
				v.expressions[span] = facts
				return true
			})
		}
	}
	return v
}

func expressionObject(pkg *sourcePackage, expr ast.Expr) types.Object {
	switch node := expr.(type) {
	case *ast.ParenExpr:
		return expressionObject(pkg, node.X)
	case *ast.IndexExpr:
		return expressionObject(pkg, node.X)
	case *ast.IndexListExpr:
		return expressionObject(pkg, node.X)
	case *ast.Ident:
		if obj := pkg.info.Uses[node]; obj != nil {
			return obj
		}
		return pkg.info.Defs[node]
	case *ast.SelectorExpr:
		if selection := pkg.info.Selections[node]; selection != nil {
			return selection.Obj()
		}
		return pkg.info.Uses[node.Sel]
	}
	return nil
}

func (v *sourceView) span(node ast.Node) capability.Span {
	if node == nil || !node.Pos().IsValid() || !node.End().IsValid() {
		return capability.Span{}
	}
	file := v.fset.File(node.Pos())
	if file == nil || v.fset.File(node.End()) != file {
		return capability.Span{}
	}
	ref := v.paths[filepath.Clean(file.Name())]
	if ref == "" {
		return capability.Span{}
	}
	return capability.Span{File: ref, Start: file.Offset(node.Pos()), End: file.Offset(node.End())}
}

func (v *sourceView) addDeclaration(obj types.Object, node ast.Node, value, typ ast.Expr) {
	symbol := canonicalSymbol(obj)
	if symbol.Package == "" || symbol.Name == "" {
		return
	}
	v.declarations[symbol] = capability.Declaration{Symbol: symbol, Span: v.span(node), Value: v.span(value), Type: v.span(typ)}
}

func (v *sourceView) File(ref capability.FileRef) (capability.File, error) {
	file, ok := v.files[ref]
	if !ok {
		return capability.File{}, fmt.Errorf("file does not belong to this source view")
	}
	return file, nil
}

func (v *sourceView) Declaration(symbol capability.Symbol) (capability.Declaration, error) {
	decl, ok := v.declarations[symbol]
	if !ok || decl.Span.File == "" {
		return capability.Declaration{}, fmt.Errorf("declaration source unavailable for %s", symbol)
	}
	return decl, nil
}

func (v *sourceView) Facts(span capability.Span) (capability.Facts, error) {
	if err := capability.ValidateSpan(v, span); err != nil {
		return capability.Facts{}, err
	}
	facts, ok := v.expressions[span]
	if !ok {
		return capability.Facts{}, fmt.Errorf("facts require an exact captured expression span")
	}
	result := capability.Facts{Symbol: symbolPointer(facts.symbol), Callee: symbolPointer(facts.callee), Type: sourceTypeFact(facts.typ)}
	if facts.value != nil {
		text := facts.value.ExactString()
		if facts.value.Kind() == constant.String {
			text = constant.StringVal(facts.value)
		}
		result.Constant = &capability.Constant{Kind: strings.ToLower(facts.value.Kind().String()), Value: text}
	}
	return result, nil
}

var _ capability.SourceView = (*sourceView)(nil)
