package cppparse

import (
	"github.com/Bitspark/bittype/trees/ast/go"
	check "github.com/Bitspark/bittype/check/cpp/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/internal/concrete"
	native "github.com/Bitspark/bittype/internal/cpp"
	"github.com/Bitspark/bittype/internal/lex"
	"github.com/Bitspark/bittype/internal/parseutil"
	tr "github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"strings"
)

type parser struct{ *parseutil.Parser[s.Language] }

var vocabulary = concrete.Vocabulary[s.Language]{Layout: s.ConcreteLayout, Token: s.ConcreteToken, Whitespace: s.ConcreteWhitespace, Comment: s.ConcreteComment}

func Parse(input text.Input) (cst.CST[s.Language], error) {
	if len(input.Text) > lex.MaxBytes {
		return cst.CST[s.Language]{}, &text.Diagnostic{Code: text.LimitExceeded, Message: "input exceeds 4 MiB"}
	}
	// C++ removes line continuations before recognizing comments. Reject that
	// preprocessing form instead of exposing a continued comment as declarations.
	for i := 0; i < len(input.Text); i++ {
		if input.Text[i] != '\\' {
			continue
		}
		j := i + 1
		for j < len(input.Text) && (input.Text[j] == ' ' || input.Text[j] == '\t') {
			j++
		}
		if j < len(input.Text) && (input.Text[j] == '\r' || input.Text[j] == '\n') {
			return cst.CST[s.Language]{}, &text.Diagnostic{Code: text.Unsupported, Message: "C++ line splicing is unsupported", Span: &ast.Span{Source: input.SourceID, Start: uint64(i), End: uint64(j + 1)}}
		}
	}
	base, err := parseutil.New[s.Language](input, lex.Options{Punctuation: []string{"...", "::", "&&", "->"}})
	if err != nil {
		return cst.CST[s.Language]{}, err
	}
	p := parser{base}
	var children []ast.Child
	declarations := false
	for !p.EOF() && p.Err == nil {
		if p.Peek() == "#" && !declarations {
			children = append(children, tr.Child("include", p.include()))
			continue
		}
		declarations = true
		children = append(children, tr.Child("declaration", p.declaration("")))
	}
	if err := p.End(); err != nil {
		return cst.CST[s.Language]{}, err
	}
	root := p.AddSpan(s.KindTranslationUnit, 0, len(input.Text), nil, children...)
	t := p.Compact(root)
	if p.Err != nil {
		return cst.CST[s.Language]{}, p.Err
	}
	if err := check.Check(t); err != nil {
		return cst.CST[s.Language]{}, err
	}
	return concrete.Parsed(t, input, p.All, vocabulary)
}
func (p *parser) include() ast.NodeID {
	start := p.Take().Start
	line := strings.LastIndex(p.Input.Text[:start], "\n") + 1
	if strings.TrimSpace(p.Input.Text[line:start]) != "" {
		p.Fail(text.SyntaxError, "preprocessor directive must start a line")
	}
	p.Want("include")
	system := p.Eat("<")
	path := ""
	if system {
		begin := p.Token().Start
		end := begin
		for !p.EOF() && p.Err == nil && p.Peek() != ">" {
			t := p.Take()
			if strings.ContainsAny(p.Input.Text[end:t.Start], "\r\n") {
				p.Fail(text.SyntaxError, "newline in header name")
			}
			end = t.End
		}
		path = p.Input.Text[begin:end]
		p.Want(">")
	} else {
		token := p.Take().Text
		if len(token) < 2 || token[0] != '"' || token[len(token)-1] != '"' {
			p.Fail(text.Unsupported, "include requires a literal header name")
		} else {
			path = token[1 : len(token)-1]
		}
	}
	if !p.EOF() && !strings.ContainsAny(p.Input.Text[p.Tokens[p.Pos-1].End:p.Token().Start], "\r\n") {
		p.Fail(text.Unsupported, "unsupported preprocessor directive tail")
	}
	if p.Pos > 0 && strings.ContainsAny(p.Input.Text[start:p.Tokens[p.Pos-1].End], "\r\n") {
		p.Fail(text.Unsupported, "multiline include directives are unsupported")
	}
	attrs := tr.Attrs("path", path)
	if system {
		attrs = append(attrs, tr.Attrs("system", "true")...)
	}
	return p.Add(s.KindInclude, start, attrs)
}
func (p *parser) templates() []ast.Child {
	var out []ast.Child
	if !p.Eat("template") {
		return out
	}
	p.Want("<")
	for p.Err == nil && !p.EOF() && p.Peek() != ">" {
		start := p.Token().Start
		key := p.Take().Text
		if key != "typename" && key != "class" {
			p.Fail(text.Unsupported, "only type template parameters are supported")
			break
		}
		name := p.Ident()
		var children []ast.Child
		if p.Eat("=") {
			children = append(children, tr.Child("default", p.typeID()))
		}
		out = append(out, tr.Child("template-parameter", p.Add(s.KindTypeParameter, start, tr.Attrs("name", name, "parameter-key", key), children...)))
		if !p.Eat(",") {
			break
		}
		if p.Peek() == ">" {
			p.Fail(text.SyntaxError, "trailing comma in template parameters")
		}
	}
	p.Want(">")
	if len(out) == 0 {
		p.Fail(text.Unsupported, "template specialization is outside this vocabulary")
	}
	return out
}
func (p *parser) declaration(access string) ast.NodeID {
	if !p.Enter() {
		return 0
	}
	defer p.Leave()
	start := p.Token().Start
	templates := p.templates()
	keyword := p.Take().Text
	attrs := []ast.Attribute{}
	if access != "" {
		attrs = tr.Attrs("access", access)
	}
	switch keyword {
	case "namespace":
		if len(templates) > 0 {
			p.Fail(text.SyntaxError, "namespace cannot be a template")
		}
		name := p.Ident()
		p.Want("{")
		var children []ast.Child
		for p.Err == nil && !p.EOF() && p.Peek() != "}" {
			children = append(children, tr.Child("declaration", p.declaration("")))
		}
		p.Want("}")
		return p.Add(s.KindNamespace, start, tr.Attrs("name", name), children...)
	case "using":
		name := p.Ident()
		attrs = append(attrs, tr.Attrs("name", name)...)
		p.Want("=")
		value := p.typeID()
		p.Want(";")
		return p.Add(s.KindTypeAlias, start, attrs, append(templates, tr.Child("type", value))...)
	case "struct", "class":
		name := p.Ident()
		attrs = append(attrs, tr.Attrs("name", name, "record-key", keyword)...)
		children := templates
		if p.Eat(";") {
			attrs = append(attrs, tr.Attrs("forward", "true")...)
			return p.Add(s.KindRecord, start, attrs, children...)
		}
		if p.Eat(":") {
			for p.Err == nil {
				begin := p.Token().Start
				var a []ast.Attribute
				seen := map[string]bool{}
				for p.Peek() == "virtual" || p.Peek() == "public" || p.Peek() == "protected" || p.Peek() == "private" {
					token := p.Take().Text
					key := "access"
					value := token
					if token == "virtual" {
						key = "virtual"
						value = "true"
					}
					if seen[key] {
						p.Fail(text.SyntaxError, "duplicate base modifier")
					}
					seen[key] = true
					a = append(a, tr.Attrs(key, value)...)
				}
				children = append(children, tr.Child("base", p.Add(s.KindBase, begin, a, tr.Child("type", p.named()))))
				if !p.Eat(",") {
					break
				}
			}
		}
		p.Want("{")
		memberAccess := ""
		for p.Err == nil && !p.EOF() && p.Peek() != "}" {
			if (p.Peek() == "public" || p.Peek() == "protected" || p.Peek() == "private") && p.At(1) == ":" {
				memberAccess = p.Take().Text
				p.Take()
				continue
			}
			var member ast.NodeID
			switch p.Peek() {
			case "template", "using", "struct", "class", "enum":
				member = p.declaration(memberAccess)
			default:
				member = p.member(name, memberAccess)
			}
			children = append(children, tr.Child("member", member))
		}
		p.Want("}")
		p.Want(";")
		return p.Add(s.KindRecord, start, attrs, children...)
	case "enum":
		if len(templates) > 0 {
			p.Fail(text.Unsupported, "enum templates are unsupported")
		}
		scoped := p.Eat("class")
		if !scoped {
			scoped = p.Eat("struct")
		}
		name := p.Ident()
		attrs = append(attrs, tr.Attrs("name", name)...)
		if scoped {
			attrs = append(attrs, tr.Attrs("scoped", "true")...)
		}
		var children []ast.Child
		if p.Eat(":") {
			children = append(children, tr.Child("underlying", p.typeID()))
		}
		if p.Eat(";") {
			attrs = append(attrs, tr.Attrs("forward", "true")...)
			return p.Add(s.KindEnum, start, attrs, children...)
		}
		p.Want("{")
		for p.Err == nil && !p.EOF() && p.Peek() != "}" {
			begin := p.Token().Start
			a := tr.Attrs("name", p.Ident())
			if p.Eat("=") {
				value := ""
				if p.Eat("-") {
					value = "-"
				} else {
					p.Eat("+")
				}
				value += p.Take().Text
				a = append(a, tr.Attrs("value", value)...)
			}
			children = append(children, tr.Child("enumerator", p.Add(s.KindEnumerator, begin, a)))
			if !p.Eat(",") {
				break
			}
		}
		p.Want("}")
		p.Want(";")
		return p.Add(s.KindEnum, start, attrs, children...)
	}
	p.Fail(text.Unsupported, "unsupported C++ declaration")
	return 0
}
func (p *parser) member(owner, access string) ast.NodeID {
	start := p.Token().Start
	virtual := p.Eat("virtual")
	attrs := []ast.Attribute{}
	if access != "" {
		attrs = tr.Attrs("access", access)
	}
	if virtual {
		attrs = append(attrs, tr.Attrs("virtual", "true")...)
	}
	if p.Eat("~") {
		name := p.Ident()
		if name != owner {
			p.Fail(text.SyntaxError, "destructor name differs from enclosing record")
		}
		p.Want("(")
		p.Want(")")
		if p.Eat("noexcept") {
			attrs = append(attrs, tr.Attrs("noexcept", "true")...)
		}
		if p.Eat("=") {
			p.Want("0")
			attrs = append(attrs, tr.Attrs("pure", "true")...)
		}
		p.Want(";")
		return p.Add(s.KindDestructor, start, attrs)
	}
	base := p.base()
	d := p.declarator(true)
	value := d.apply(base)
	if d.name == "" {
		p.Fail(text.SyntaxError, "record member requires a name")
		return 0
	}
	attrs = append(attrs, tr.Attrs("name", d.name)...)
	if p.Err != nil || value == 0 {
		return 0
	}
	if p.Tree.Nodes[value-1].Kind == s.KindFunctionType {
		function := p.Tree.Nodes[value-1]
		attrs = append(attrs, function.Attributes...)
		for _, q := range d.qualified {
			if q.id != value {
				p.Fail(text.Unsupported, "qualified function-pointer types are not represented")
				return 0
			}
			attrs = append(attrs, q.attrs...)
		}
		if p.Eat("=") {
			p.Want("0")
			attrs = append(attrs, tr.Attrs("pure", "true")...)
		}
		p.Want(";")
		return p.Add(s.KindMethod, start, attrs, function.Children...)
	}
	if virtual || len(d.qualified) > 0 {
		p.Fail(text.Unsupported, "unsupported qualified field type")
		return 0
	}
	p.Want(";")
	return p.Add(s.KindField, start, attrs, tr.Child("type", value))
}
func (p *parser) named() ast.NodeID {
	start := p.Token().Start
	typename := p.Eat("typename")
	global := p.Eat("::")
	var children []ast.Child
	for p.Err == nil {
		begin := p.Token().Start
		templ := p.Eat("template")
		name := p.Ident()
		attrs := tr.Attrs("name", name)
		if templ {
			attrs = append(attrs, tr.Attrs("template", "true")...)
		}
		var args []ast.Child
		if p.Eat("<") {
			for p.Err == nil && !p.EOF() && p.Peek() != ">" {
				args = append(args, tr.Child("argument", p.typeID()))
				if !p.Eat(",") {
					break
				}
				if p.Peek() == ">" {
					p.Fail(text.SyntaxError, "trailing comma in template arguments")
				}
			}
			p.Want(">")
			if len(args) == 0 {
				p.Fail(text.Unsupported, "empty template applications are not represented")
			}
		}
		children = append(children, tr.Child("segment", p.Add(s.KindNameSegment, begin, attrs, args...)))
		if p.Peek() != "::" || p.At(1) == "*" {
			break
		}
		p.Take()
	}
	var attrs []ast.Attribute
	if typename {
		attrs = tr.Attrs("typename", "true")
	}
	if global {
		attrs = append(attrs, tr.Attrs("global", "true")...)
	}
	return p.Add(s.KindNamedType, start, attrs, children...)
}
func (p *parser) cv() []ast.Attribute {
	var attrs []ast.Attribute
	seen := map[string]bool{}
	for p.Peek() == "const" || p.Peek() == "volatile" {
		name := p.Take().Text
		if seen[name] {
			p.Fail(text.SyntaxError, "duplicate cv qualifier")
		}
		seen[name] = true
		attrs = append(attrs, tr.Attrs(name, "true")...)
	}
	return attrs
}
func (p *parser) base() ast.NodeID {
	start := p.Token().Start
	qualifiers := p.cv()
	var id ast.NodeID
	if native.BuiltinWord(p.Peek()) {
		var words []string
		for native.BuiltinWord(p.Peek()) {
			words = append(words, p.Take().Text)
		}
		id = p.Add(s.KindBuiltinType, start, tr.Attrs("name", strings.Join(words, " ")))
	} else {
		switch p.Peek() {
		case "auto", "decltype", "struct", "class", "enum", "consteval", "constexpr", "static", "friend", "inline":
			p.Fail(text.Unsupported, "unsupported type specifier")
			return 0
		}
		id = p.named()
	}
	qualifiers = append(qualifiers, p.cv()...)
	if len(qualifiers) > 0 {
		id = p.Add(s.KindCVQualifiedType, start, qualifiers, tr.Child("element", id))
	}
	return id
}
func (p *parser) typeID() ast.NodeID {
	if !p.Enter() {
		return 0
	}
	defer p.Leave()
	base := p.base()
	d := p.declarator(false)
	id := d.apply(base)
	if len(d.qualified) > 0 {
		p.Fail(text.Unsupported, "qualified function types are not represented")
	}
	return id
}

type Provider struct{ s.Language }

func (Provider) Parse(input text.Input) (cst.CST[s.Language], error) { return Parse(input) }

var _ text.ParserProvider[s.Language] = Provider{}
