package parse

import (
	"context"
	"errors"
	check "github.com/Bitspark/bittype/check/syntax/rust/go"
	parsecontract "github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/internal/concrete"
	"github.com/Bitspark/bittype/internal/lex"
	"github.com/Bitspark/bittype/internal/parseutil"
	native "github.com/Bitspark/bittype/internal/rust"
	tr "github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"strings"
)

type parser struct{ *parseutil.Parser[s.Language] }

var vocabulary = concrete.Vocabulary[s.Language]{Layout: s.ConcreteLayout, Token: s.ConcreteToken, Whitespace: s.ConcreteWhitespace, Comment: s.ConcreteComment}

// bittype:scope docs/syntax.md#rust
func Parse(ctx context.Context, input parsecontract.Input) (cst.CST[s.Language], error) {
	if ctx == nil {
		return cst.CST[s.Language]{}, errors.New("nil context")
	}
	if err := ctx.Err(); err != nil {
		return cst.CST[s.Language]{}, err
	}
	base, err := parseutil.New[s.Language](input, lex.Options{Space: native.Space, NestedComments: true, Lifetimes: true, RawIdentifiers: true, Punctuation: []string{"::", "->", "..."}})
	if err != nil {
		return cst.CST[s.Language]{}, err
	}
	p := parser{base}
	children := p.items("")
	if err := p.End(); err != nil {
		return cst.CST[s.Language]{}, err
	}
	root := p.AddSpan(s.KindFile, 0, len(input.Text), nil, children...)
	t := p.Compact(root)
	if p.Err != nil {
		return cst.CST[s.Language]{}, p.Err
	}
	if err := check.Check(t); err != nil {
		return cst.CST[s.Language]{}, err
	}
	return concrete.Parsed(t, input, p.All, vocabulary)
}
func (p *parser) visibility() string {
	if !p.Eat("pub") {
		return ""
	}
	if !p.Eat("(") {
		return "public"
	}
	p.Want("crate")
	p.Want(")")
	return "crate"
}
func (p *parser) items(end string) []ast.Child {
	var out []ast.Child
	for !p.EOF() && p.Err == nil && p.Peek() != end {
		out = append(out, tr.Child("item", p.item()))
	}
	return out
}
func (p *parser) lifetime() ast.NodeID {
	start := p.Token().Start
	if p.Token().Class != "lifetime" {
		p.Fail(text.SyntaxError, "expected lifetime")
		return 0
	}
	name := strings.TrimPrefix(p.Take().Text, "'")
	return p.Add(s.KindLifetimeReference, start, tr.Attrs("name", name))
}
func (p *parser) bounds() []ast.Child {
	var out []ast.Child
	for p.Err == nil {
		start := p.Token().Start
		var id ast.NodeID
		if p.Eat("?") {
			id = p.Add(s.KindRelaxedBound, start, nil, tr.Child("path", p.path()))
		} else if p.Token().Class == "lifetime" {
			id = p.lifetime()
		} else {
			id = p.path()
		}
		out = append(out, tr.Child("bound", id))
		if !p.Eat("+") {
			break
		}
	}
	return out
}
func (p *parser) generics() []ast.Child {
	var out []ast.Child
	if !p.Eat("<") {
		return out
	}
	for p.Err == nil && !p.EOF() && p.Peek() != ">" {
		start := p.Token().Start
		kind := s.KindTypeParameter
		name := ""
		if p.Token().Class == "lifetime" {
			kind = s.KindLifetimeParameter
			name = strings.TrimPrefix(p.Take().Text, "'")
		} else {
			name = p.Ident()
		}
		var children []ast.Child
		if p.Eat(":") {
			children = append(children, p.bounds()...)
		}
		if p.Eat("=") {
			children = append(children, tr.Child("default", p.typ()))
		}
		out = append(out, tr.Child("generic-parameter", p.Add(kind, start, tr.Attrs("name", name), children...)))
		if !p.Eat(",") {
			break
		}
	}
	p.Want(">")
	if len(out) == 0 {
		p.Fail(text.SyntaxError, "empty generic parameter list")
	}
	return out
}
func (p *parser) where() []ast.Child {
	var out []ast.Child
	if !p.Eat("where") {
		return out
	}
	for p.Err == nil && !p.EOF() {
		start := p.Token().Start
		var subject ast.NodeID
		if p.Token().Class == "lifetime" {
			subject = p.lifetime()
		} else {
			subject = p.typ()
		}
		p.Want(":")
		children := []ast.Child{tr.Child("subject", subject)}
		children = append(children, p.bounds()...)
		out = append(out, tr.Child("where", p.Add(s.KindWherePredicate, start, nil, children...)))
		if !p.Eat(",") {
			break
		}
		if p.Peek() == "{" || p.Peek() == ";" || p.Peek() == "=" {
			break
		}
	}
	return out
}
func (p *parser) item() ast.NodeID {
	if !p.Enter() {
		return 0
	}
	defer p.Leave()
	start := p.Token().Start
	vis := p.visibility()
	keyword := p.Take().Text
	if keyword == "impl" {
		if vis != "" {
			p.Fail(text.SyntaxError, "trait impls have no visibility qualifier")
			return 0
		}
		return p.traitImpl(start)
	}
	if keyword == "use" {
		if vis != "" {
			p.Fail(text.Unsupported, "visible use declarations are outside this vocabulary")
			return 0
		}
		path := p.path()
		var attrs []ast.Attribute
		if p.Eat("as") {
			attrs = tr.Attrs("alias", p.Ident())
		}
		p.Want(";")
		return p.Add(s.KindUse, start, attrs, tr.Child("path", path))
	}
	if keyword != "mod" && keyword != "type" && keyword != "struct" && keyword != "enum" && keyword != "trait" {
		p.Fail(text.Unsupported, "unsupported Rust item")
		return 0
	}
	name := p.Ident()
	attrs := tr.Attrs("name", name)
	if vis != "" {
		attrs = append(attrs, tr.Attrs("visibility", vis)...)
	}
	if keyword == "mod" {
		p.Want("{")
		children := p.items("}")
		p.Want("}")
		return p.Add(s.KindModule, start, attrs, children...)
	}
	children := p.generics()
	if keyword == "trait" && p.Eat(":") {
		children = append(children, p.bounds()...)
	}
	if keyword == "struct" && p.Peek() == "(" {
		fields := p.fields(false, false)
		children = append(children, p.where()...)
		children = append(children, fields...)
		p.Want(";")
		return p.Add(s.KindTupleStruct, start, attrs, children...)
	}
	children = append(children, p.where()...)
	switch keyword {
	case "type":
		p.Want("=")
		children = append(children, tr.Child("type", p.typ()))
		p.Want(";")
		return p.Add(s.KindTypeAlias, start, attrs, children...)
	case "struct":
		if p.Peek() == ";" {
			p.Fail(text.Unsupported, "unit structs require a distinct AST form")
			return 0
		}
		children = append(children, p.fields(true, false)...)
		return p.Add(s.KindStruct, start, attrs, children...)
	case "enum":
		p.Want("{")
		for p.Err == nil && !p.EOF() && p.Peek() != "}" {
			begin := p.Token().Start
			variant := p.Ident()
			kind := s.KindUnitVariant
			var fields []ast.Child
			if p.Peek() == "(" {
				kind = s.KindTupleVariant
				fields = p.fields(false, true)
			} else if p.Peek() == "{" {
				kind = s.KindRecordVariant
				fields = p.fields(true, true)
			}
			children = append(children, tr.Child("variant", p.Add(kind, begin, tr.Attrs("name", variant), fields...)))
			if !p.Eat(",") {
				break
			}
		}
		p.Want("}")
		return p.Add(s.KindEnum, start, attrs, children...)
	case "trait":
		p.Want("{")
		for p.Err == nil && !p.EOF() && p.Peek() != "}" {
			begin := p.Token().Start
			if p.Eat("type") {
				member := p.Ident()
				declaration := p.generics()
				if p.Eat(":") {
					declaration = append(declaration, p.bounds()...)
				}
				declaration = append(declaration, p.where()...)
				if p.Peek() == "=" {
					p.Fail(text.Unsupported, "associated-type defaults are unsupported")
				}
				p.Want(";")
				children = append(children, tr.Child("member", p.Add(s.KindAssociatedType, begin, tr.Attrs("name", member), declaration...)))
			} else {
				children = append(children, tr.Child("member", p.method()))
			}
		}
		p.Want("}")
		return p.Add(s.KindTrait, start, attrs, children...)
	}
	return 0
}
func (p *parser) traitImpl(start int) ast.NodeID {
	children := p.generics()
	if p.Peek() == "!" {
		p.Fail(text.Unsupported, "negative trait impls are unsupported")
		return 0
	}
	trait := p.path()
	if p.Err == nil && p.Peek() != "for" {
		p.Fail(text.Unsupported, "only trait impls are supported; expected for")
		return 0
	}
	p.Want("for")
	children = append(children, tr.Child("trait", trait), tr.Child("subject", p.typ()))
	children = append(children, p.where()...)
	p.Want("{")
	for p.Err == nil && !p.EOF() && p.Peek() != "}" {
		begin := p.Token().Start
		if !p.Eat("type") {
			p.Fail(text.Unsupported, "type-only trait impls admit associated type definitions only")
			return 0
		}
		name := p.Ident()
		definition := p.generics()
		if p.Peek() == ":" {
			p.Fail(text.SyntaxError, "associated type definitions cannot declare result bounds")
			return 0
		}
		where := p.where()
		p.Want("=")
		value := p.typ()
		if p.Peek() == "where" {
			if len(where) != 0 {
				p.Fail(text.SyntaxError, "associated type definition has two where clauses")
				return 0
			}
			where = p.where()
		}
		definition = append(definition, where...)
		definition = append(definition, tr.Child("type", value))
		p.Want(";")
		children = append(children, tr.Child("member", p.Add(s.KindTypeAlias, begin, tr.Attrs("name", name), definition...)))
	}
	p.Want("}")
	return p.Add(s.KindTraitImpl, start, nil, children...)
}
func (p *parser) fields(named, variant bool) []ast.Child {
	open, close := "(", ")"
	kind := s.KindTupleField
	if named {
		open, close = "{", "}"
		kind = s.KindField
	}
	p.Want(open)
	var out []ast.Child
	for p.Err == nil && !p.EOF() && p.Peek() != close {
		start := p.Token().Start
		vis := p.visibility()
		var attrs []ast.Attribute
		if vis != "" {
			attrs = tr.Attrs("visibility", vis)
			if variant {
				p.Fail(text.SyntaxError, "variant fields inherit visibility")
			}
		}
		if named {
			attrs = append(attrs, tr.Attrs("name", p.Ident())...)
			p.Want(":")
		}
		value := p.typ()
		out = append(out, tr.Child("field", p.Add(kind, start, attrs, tr.Child("type", value))))
		if !p.Eat(",") {
			break
		}
	}
	p.Want(close)
	return out
}
func (p *parser) method() ast.NodeID {
	start := p.Token().Start
	p.Want("fn")
	name := p.Ident()
	generics := p.generics()
	p.Want("(")
	var receiver ast.NodeID
	var params []ast.Child
	if p.Peek() == "self" || p.Peek() == "&" {
		begin := p.Token().Start
		form := "value"
		var children []ast.Child
		if p.Eat("&") {
			form = "shared"
			if p.Token().Class == "lifetime" {
				children = append(children, tr.Child("lifetime", p.lifetime()))
			}
			if p.Eat("mut") {
				form = "mutable"
			}
		}
		p.Want("self")
		receiver = p.Add(s.KindReceiver, begin, tr.Attrs("receiver-form", form), children...)
		if p.Peek() != ")" {
			p.Want(",")
		}
	}
	for p.Err == nil && !p.EOF() && p.Peek() != ")" {
		begin := p.Token().Start
		name := p.Ident()
		p.Want(":")
		value := p.typ()
		params = append(params, tr.Child("parameter", p.Add(s.KindParameter, begin, tr.Attrs("name", name), tr.Child("type", value))))
		if !p.Eat(",") {
			break
		}
	}
	p.Want(")")
	var result ast.NodeID
	if p.Eat("->") {
		result = p.typ()
	}
	children := append(generics, p.where()...)
	if receiver != 0 {
		children = append(children, tr.Child("receiver", receiver))
	}
	children = append(children, params...)
	if result != 0 {
		children = append(children, tr.Child("result", result))
	}
	p.Want(";")
	return p.Add(s.KindMethodSignature, start, tr.Attrs("name", name), children...)
}
func (p *parser) path() ast.NodeID {
	start := p.Token().Start
	absolute := p.Eat("::")
	children := p.pathSegments()
	var attrs []ast.Attribute
	if absolute {
		attrs = tr.Attrs("absolute", "true")
	}
	return p.Add(s.KindPathType, start, attrs, children...)
}
func (p *parser) pathSegments() []ast.Child {
	var children []ast.Child
	for p.Err == nil {
		begin := p.Token().Start
		name := p.Ident()
		if p.Peek() == "::" && p.At(1) == "<" {
			p.Take() // Optional turbofish separator in a type path.
		}
		args := p.arguments()
		children = append(children, tr.Child("segment", p.Add(s.KindPathSegment, begin, tr.Attrs("name", name), args...)))
		if !p.Eat("::") {
			break
		}
	}
	return children
}
func (p *parser) arguments() []ast.Child {
	if !p.Eat("<") {
		return nil
	}
	var args []ast.Child
	for p.Err == nil && !p.EOF() && p.Peek() != ">" {
		start := p.Token().Start
		bindingName := p.Token().Class == "identifier" && native.Declaration(p.Peek())
		var arg ast.NodeID
		if p.Token().Class == "lifetime" {
			arg = p.lifetime()
		} else {
			arg = p.typ()
		}
		if p.Err != nil {
			break
		}
		if p.Eat("=") {
			if !bindingName {
				p.Fail(text.SyntaxError, "associated-type binding requires one unqualified name")
				break
			}
			// A binding's left side is exactly one TypePathSegment, which
			// retains its own arguments and origin. Compact drops the wrapper.
			n := p.Tree.Nodes[arg-1]
			var name ast.NodeID
			v := tr.View[s.Language]{Tree: p.Tree}
			if n.Kind == s.KindPrimitiveType {
				// Primitive spellings are also legal associated identifiers.
				p.Tree.Nodes[arg-1].Kind = s.KindPathSegment
				name = arg
			} else if n.Kind == s.KindPathType && !v.Flag(arg, "absolute") && len(n.Children) == 1 {
				name = n.Children[0].Node
			} else {
				p.Fail(text.SyntaxError, "associated-type binding requires one unqualified name")
				break
			}
			value := p.typ()
			arg = p.Add(s.KindAssociatedTypeBinding, start, nil, tr.Child("path", name), tr.Child("type", value))
		} else if p.Peek() == ":" {
			p.Fail(text.Unsupported, "associated-type bounds are unsupported; expected an equality binding")
		}
		args = append(args, tr.Child("argument", arg))
		if !p.Eat(",") {
			break
		}
	}
	p.Want(">")
	if len(args) == 0 {
		p.Fail(text.SyntaxError, "empty generic arguments")
	}
	return args
}
func (p *parser) qualifiedPath() ast.NodeID {
	start := p.Token().Start
	p.Want("<")
	children := []ast.Child{tr.Child("subject", p.typ())}
	if p.Eat("as") {
		children = append(children, tr.Child("trait", p.path()))
	}
	p.Want(">")
	p.Want("::")
	children = append(children, p.pathSegments()...)
	return p.Add(s.KindQualifiedPathType, start, nil, children...)
}
func (p *parser) typ() ast.NodeID {
	if !p.Enter() {
		return 0
	}
	defer p.Leave()
	start := p.Token().Start
	switch p.Peek() {
	case "<":
		return p.qualifiedPath()
	case "&":
		p.Take()
		var children []ast.Child
		if p.Token().Class == "lifetime" {
			children = append(children, tr.Child("lifetime", p.lifetime()))
		}
		var attrs []ast.Attribute
		if p.Eat("mut") {
			attrs = tr.Attrs("mutable", "true")
		}
		children = append([]ast.Child{tr.Child("element", p.typ())}, children...)
		return p.Add(s.KindReferenceType, start, attrs, children...)
	case "*":
		p.Take()
		var attrs []ast.Attribute
		if p.Eat("mut") {
			attrs = tr.Attrs("mutable", "true")
		} else {
			p.Want("const")
		}
		return p.Add(s.KindRawPointerType, start, attrs, tr.Child("element", p.typ()))
	case "[":
		p.Take()
		element := p.typ()
		if p.Eat(";") {
			length := p.Take().Text
			p.Want("]")
			return p.Add(s.KindArrayType, start, tr.Attrs("length", length), tr.Child("element", element))
		}
		p.Want("]")
		return p.Add(s.KindSliceType, start, nil, tr.Child("element", element))
	case "(":
		p.Take()
		var children []ast.Child
		comma := false
		for p.Err == nil && !p.EOF() && p.Peek() != ")" {
			children = append(children, tr.Child("element", p.typ()))
			if !p.Eat(",") {
				break
			}
			comma = true
		}
		p.Want(")")
		if len(children) == 1 && !comma {
			id := children[0].Node
			if p.Err != nil || id == 0 {
				return 0
			}
			p.Tree.Nodes[id-1].Origin.Spans = []ast.Span{{Source: p.Input.SourceID, Start: uint64(start), End: uint64(p.Tokens[p.Pos-1].End)}}
			return id
		}
		return p.Add(s.KindTupleType, start, nil, children...)
	case "!":
		p.Take()
		return p.Add(s.KindNeverType, start, nil)
	case "dyn":
		p.Take()
		return p.Add(s.KindTraitObjectType, start, nil, p.bounds()...)
	case "fn", "unsafe", "extern":
		var attrs []ast.Attribute
		if p.Eat("unsafe") {
			attrs = tr.Attrs("unsafe", "true")
		}
		if p.Eat("extern") {
			abi := "C"
			if p.Token().Class == "string" {
				value, err := native.UnquoteABI(p.Take().Text)
				if err != nil {
					p.Fail(text.SyntaxError, "invalid ABI string")
				}
				abi = value
			}
			attrs = append(attrs, tr.Attrs("abi", abi)...)
		}
		p.Want("fn")
		p.Want("(")
		var children []ast.Child
		for p.Err == nil && !p.EOF() && p.Peek() != ")" {
			begin := p.Token().Start
			var a []ast.Attribute
			if p.Token().Class == "identifier" && p.At(1) == ":" {
				a = tr.Attrs("name", p.Ident())
				p.Want(":")
			}
			children = append(children, tr.Child("parameter", p.Add(s.KindParameter, begin, a, tr.Child("type", p.typ()))))
			if !p.Eat(",") {
				break
			}
		}
		p.Want(")")
		if p.Eat("->") {
			children = append(children, tr.Child("result", p.typ()))
		}
		return p.Add(s.KindFunctionType, start, attrs, children...)
	case "impl", "for":
		p.Fail(text.Unsupported, "unsupported Rust type binder")
		return 0
	}
	name := p.Peek()
	switch name {
	case "bool", "char", "str", "f32", "f64", "i8", "i16", "i32", "i64", "i128", "isize", "u8", "u16", "u32", "u64", "u128", "usize":
		if p.At(1) == "::" || p.At(1) == "<" {
			return p.path()
		}
		p.Take()
		return p.Add(s.KindPrimitiveType, start, tr.Attrs("name", name))
	}
	return p.path()
}

type Provider struct{ s.Language }

func (Provider) Parse(ctx context.Context, input parsecontract.Input) (cst.CST[s.Language], error) {
	return Parse(ctx, input)
}

var _ parsecontract.Provider[parsecontract.Input, cst.CST[s.Language]] = Provider{}
