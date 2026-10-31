package golangparse

import (
	"github.com/Bitspark/bittype/ast/go"
	check "github.com/Bitspark/bittype/check/golang/go"
	"github.com/Bitspark/bittype/cst/go"
	"github.com/Bitspark/bittype/internal/concrete"
	"github.com/Bitspark/bittype/internal/lex"
	"github.com/Bitspark/bittype/internal/parseutil"
	tr "github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/text/go"
	stdast "go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"strconv"
)

var vocabulary = concrete.Vocabulary[s.Language]{Layout: s.ConcreteLayout, Token: s.ConcreteToken, Whitespace: s.ConcreteWhitespace, Comment: s.ConcreteComment}

func Parse(input text.Input) (cst.CST[s.Language], error) {
	p, err := parseutil.New[s.Language](input, lex.Options{BacktickStrings: true, Punctuation: []string{"...", "<-"}, Directives: []string{"//go:", "// +build", "//line "}})
	if err != nil {
		return cst.CST[s.Language]{}, err
	}
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, string(input.SourceID), input.Text, parser.AllErrors|parser.SkipObjectResolution|parser.ParseComments)
	if err != nil {
		d := &text.Diagnostic{Code: text.SyntaxError, Message: err.Error()}
		if errors, ok := err.(scanner.ErrorList); ok && len(errors) > 0 {
			pos := errors[0].Pos.Offset
			d.Span = &ast.Span{Source: input.SourceID, Start: uint64(pos), End: uint64(pos)}
		}
		return cst.CST[s.Language]{}, d
	}
	c := converter{Parser: p, positions: fs}
	var imports, decls []ast.Child
	for _, decl := range file.Decls {
		group, ok := decl.(*stdast.GenDecl)
		if !ok || (group.Tok != token.IMPORT && group.Tok != token.TYPE) {
			c.fail(decl, text.Unsupported, "only imports and type declarations are supported")
			break
		}
		for _, spec := range group.Specs {
			switch spec := spec.(type) {
			case *stdast.ImportSpec:
				path, e := strconv.Unquote(spec.Path.Value)
				if e != nil {
					c.fail(spec, text.SyntaxError, "invalid import string")
					continue
				}
				attrs := tr.Attrs("path", path)
				if spec.Name != nil {
					attrs = append(attrs, tr.Attrs("alias", spec.Name.Name)...)
				}
				imports = append(imports, tr.Child("import", c.node(s.KindImport, spec, attrs)))
			case *stdast.TypeSpec:
				children := tr.Children("type-parameter", c.typeParams(spec.TypeParams))
				children = append(children, tr.Child("type", c.typ(spec.Type)))
				kind := s.KindTypeDeclaration
				if spec.Assign.IsValid() {
					kind = s.KindTypeAlias
				}
				decls = append(decls, tr.Child("declaration", c.node(kind, spec, tr.Attrs("name", spec.Name.Name), children...)))
			default:
				c.fail(spec, text.Unsupported, "unsupported declaration")
			}
		}
	}
	if c.Err != nil {
		return cst.CST[s.Language]{}, c.Err
	}
	root := p.AddSpan(s.KindFile, 0, len(input.Text), tr.Attrs("package", file.Name.Name), append(imports, decls...)...)
	t := p.Compact(root)
	if err := check.Check(t); err != nil {
		return cst.CST[s.Language]{}, err
	}
	return concrete.Parsed(t, input, p.All, vocabulary)
}

type converter struct {
	*parseutil.Parser[s.Language]
	positions *token.FileSet
}

func (c *converter) node(kind s.Kind, n stdast.Node, attrs []ast.Attribute, children ...ast.Child) ast.NodeID {
	return c.AddSpan(kind, c.positions.Position(n.Pos()).Offset, c.positions.Position(n.End()).Offset, attrs, children...)
}
func (c *converter) fail(n stdast.Node, code text.Code, message string) {
	if c.Err != nil {
		return
	}
	start, end := c.positions.Position(n.Pos()).Offset, c.positions.Position(n.End()).Offset
	c.Err = &text.Diagnostic{Code: code, Message: message, Span: &ast.Span{Source: c.Input.SourceID, Start: uint64(start), End: uint64(end)}}
}
func (c *converter) typeParams(list *stdast.FieldList) []ast.NodeID {
	var out []ast.NodeID
	if list == nil {
		return out
	}
	for _, field := range list.List {
		for _, name := range field.Names {
			out = append(out, c.node(s.KindTypeParameter, field, tr.Attrs("name", name.Name), tr.Child("constraint", c.typ(field.Type))))
		}
	}
	return out
}
func (c *converter) params(list *stdast.FieldList) []ast.NodeID {
	var out []ast.NodeID
	if list == nil {
		return out
	}
	for _, field := range list.List {
		if field.Tag != nil {
			c.fail(field, text.Unsupported, "parameter tags are unsupported")
			return out
		}
		t := field.Type
		variadic := false
		if ellipsis, ok := t.(*stdast.Ellipsis); ok {
			t = ellipsis.Elt
			variadic = true
		}
		names := field.Names
		if len(names) == 0 {
			names = []*stdast.Ident{nil}
		}
		for _, name := range names {
			var attrs []ast.Attribute
			if name != nil {
				attrs = tr.Attrs("name", name.Name)
			}
			if variadic {
				attrs = append(attrs, tr.Attrs("variadic", "true")...)
			}
			out = append(out, c.node(s.KindParameter, field, attrs, tr.Child("type", c.typ(t))))
		}
	}
	return out
}
func (c *converter) typ(expr stdast.Expr) ast.NodeID {
	if expr == nil || !c.Enter() {
		return 0
	}
	defer c.Leave()
	node := func(kind s.Kind, attrs []ast.Attribute, children ...ast.Child) ast.NodeID {
		return c.node(kind, expr, attrs, children...)
	}
	switch e := expr.(type) {
	case *stdast.Ident:
		return node(s.KindNamedType, tr.Attrs("name", e.Name))
	case *stdast.SelectorExpr:
		q, ok := e.X.(*stdast.Ident)
		if !ok {
			c.fail(e, text.Unsupported, "qualified types require a package identifier")
			return 0
		}
		return node(s.KindQualifiedType, tr.Attrs("qualifier", q.Name, "name", e.Sel.Name))
	case *stdast.ParenExpr:
		id := c.typ(e.X)
		if id != 0 {
			c.Tree.Nodes[id-1].Origin.Spans = []ast.Span{{Source: c.Input.SourceID, Start: uint64(c.positions.Position(e.Pos()).Offset), End: uint64(c.positions.Position(e.End()).Offset)}}
		}
		return id
	case *stdast.IndexExpr:
		return node(s.KindTypeApplication, nil, tr.Child("constructor", c.typ(e.X)), tr.Child("argument", c.typ(e.Index)))
	case *stdast.IndexListExpr:
		children := []ast.Child{tr.Child("constructor", c.typ(e.X))}
		for _, arg := range e.Indices {
			children = append(children, tr.Child("argument", c.typ(arg)))
		}
		return node(s.KindTypeApplication, nil, children...)
	case *stdast.StarExpr:
		return node(s.KindPointerType, nil, tr.Child("element", c.typ(e.X)))
	case *stdast.ArrayType:
		if e.Len == nil {
			return node(s.KindSliceType, nil, tr.Child("element", c.typ(e.Elt)))
		}
		length, ok := e.Len.(*stdast.BasicLit)
		if !ok || length.Kind != token.INT {
			c.fail(e, text.Unsupported, "array lengths must be decimal integer literals")
			return 0
		}
		return node(s.KindArrayType, tr.Attrs("length", length.Value), tr.Child("element", c.typ(e.Elt)))
	case *stdast.MapType:
		return node(s.KindMapType, nil, tr.Child("key", c.typ(e.Key)), tr.Child("element", c.typ(e.Value)))
	case *stdast.ChanType:
		direction := "both"
		if e.Dir == stdast.SEND {
			direction = "send"
		}
		if e.Dir == stdast.RECV {
			direction = "receive"
		}
		return node(s.KindChannelType, tr.Attrs("direction", direction), tr.Child("element", c.typ(e.Value)))
	case *stdast.FuncType:
		if e.TypeParams != nil {
			c.fail(e, text.Unsupported, "generic function types are unsupported")
			return 0
		}
		children := tr.Children("parameter", c.params(e.Params))
		children = append(children, tr.Children("result", c.params(e.Results))...)
		return node(s.KindFunctionType, nil, children...)
	case *stdast.StructType:
		var children []ast.Child
		for _, field := range e.Fields.List {
			var attrs []ast.Attribute
			if field.Tag != nil {
				tag, err := strconv.Unquote(field.Tag.Value)
				if err != nil {
					c.fail(field, text.SyntaxError, "invalid struct tag")
					return 0
				}
				attrs = tr.Attrs("tag", tag)
			}
			if len(field.Names) == 0 {
				children = append(children, tr.Child("field", c.node(s.KindEmbeddedField, field, attrs, tr.Child("type", c.typ(field.Type)))))
				continue
			}
			for _, name := range field.Names {
				a := append(tr.Attrs("name", name.Name), attrs...)
				children = append(children, tr.Child("field", c.node(s.KindField, field, a, tr.Child("type", c.typ(field.Type)))))
			}
		}
		return node(s.KindStructType, nil, children...)
	case *stdast.InterfaceType:
		var children []ast.Child
		for _, field := range e.Methods.List {
			if len(field.Names) == 0 {
				children = append(children, tr.Child("member", c.node(s.KindEmbeddedType, field, nil, tr.Child("type", c.typ(field.Type)))))
				continue
			}
			if len(field.Names) != 1 {
				c.fail(field, text.SyntaxError, "method requires one name")
				return 0
			}
			children = append(children, tr.Child("member", c.node(s.KindMethod, field, tr.Attrs("name", field.Names[0].Name), tr.Child("signature", c.typ(field.Type)))))
		}
		return node(s.KindInterfaceType, nil, children...)
	case *stdast.UnaryExpr:
		if e.Op == token.TILDE {
			return node(s.KindTypeApproximation, nil, tr.Child("type", c.typ(e.X)))
		}
	case *stdast.BinaryExpr:
		if e.Op == token.OR {
			var terms []ast.Child
			var flatten func(stdast.Expr)
			flatten = func(x stdast.Expr) {
				if b, ok := x.(*stdast.BinaryExpr); ok && b.Op == token.OR {
					flatten(b.X)
					flatten(b.Y)
				} else {
					terms = append(terms, tr.Child("term", c.typ(x)))
				}
			}
			flatten(e)
			return node(s.KindTypeUnion, nil, terms...)
		}
	}
	c.fail(expr, text.Unsupported, "unsupported Go type syntax")
	return 0
}

type Provider struct{ s.Language }

func (Provider) Parse(input text.Input) (cst.CST[s.Language], error) { return Parse(input) }

var _ text.ParserProvider[s.Language] = Provider{}
