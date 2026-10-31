package typescriptparse

import (
	check "github.com/Bitspark/bittype/check/typescript/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/internal/concrete"
	"github.com/Bitspark/bittype/internal/lex"
	"github.com/Bitspark/bittype/internal/parseutil"
	tr "github.com/Bitspark/bittype/internal/tree"
	lit "github.com/Bitspark/bittype/internal/typescript"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"strconv"
	"strings"
)

type parser struct{ *parseutil.Parser[s.Language] }

func decimalKey(s string) bool {
	if len(s) == 0 || (len(s) > 1 && s[0] == '0') {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	value, err := strconv.ParseUint(s, 10, 64)
	return err == nil && value <= 9007199254740991
}

var vocabulary = concrete.Vocabulary[s.Language]{Layout: s.ConcreteLayout, Token: s.ConcreteToken, Whitespace: s.ConcreteWhitespace, Comment: s.ConcreteComment}

func Parse(input text.Input) (cst.CST[s.Language], error) {
	base, err := parseutil.New[s.Language](input, lex.Options{Space: lit.Space, LineTerminator: lit.LineTerminator, RejectComment: lit.Directive, LeadingDotNumbers: true, Punctuation: []string{"...", "=>"}})
	if err != nil {
		return cst.CST[s.Language]{}, err
	}
	p := parser{base}
	var children []ast.Child
	declarations := false
	for !p.EOF() && p.Err == nil {
		if p.Peek() == "import" && !declarations {
			children = append(children, tr.Child("import", p.importDecl()))
			continue
		}
		declarations = true
		children = append(children, tr.Child("declaration", p.declaration()))
	}
	if err := p.End(); err != nil {
		return cst.CST[s.Language]{}, err
	}
	root := p.AddSpan(s.KindSourceFile, 0, len(input.Text), nil, children...)
	tree := p.Compact(root)
	if p.Err != nil {
		return cst.CST[s.Language]{}, p.Err
	}
	if err := check.Check(tree); err != nil {
		return cst.CST[s.Language]{}, err
	}
	return concrete.Parsed(tree, input, p.All, vocabulary)
}
func (p *parser) semi() {
	if p.Eat(";") || p.EOF() || p.Peek() == "}" {
		return
	}
	start := 0
	if p.Pos > 0 {
		start = p.Tokens[p.Pos-1].End
	}
	if !strings.ContainsAny(p.Input.Text[start:p.Token().Start], "\r\n\u2028\u2029") {
		p.Fail(text.SyntaxError, "expected semicolon or line break")
	}
}
func (p *parser) quoted() string {
	token := p.Take()
	value, err := lit.Unquote(token.Text)
	if err != nil {
		p.Fail(text.SyntaxError, "invalid string token")
		return ""
	}
	return value
}
func (p *parser) importDecl() ast.NodeID {
	start := p.Take().Start
	p.Want("type")
	var children []ast.Child
	if p.Eat("*") {
		begin := p.Tokens[p.Pos-1].Start
		p.Want("as")
		name := p.Ident()
		children = append(children, tr.Child("binding", p.Add(s.KindNamespaceImport, begin, tr.Attrs("name", name))))
	} else {
		p.Want("{")
		for p.Err == nil && !p.EOF() && p.Peek() != "}" {
			begin := p.Token().Start
			name := p.Ident()
			attrs := tr.Attrs("name", name)
			if p.Eat("as") {
				attrs = append(attrs, tr.Attrs("alias", p.Ident())...)
			}
			children = append(children, tr.Child("binding", p.Add(s.KindImportBinding, begin, attrs)))
			if !p.Eat(",") {
				break
			}
		}
		p.Want("}")
	}
	p.Want("from")
	module := p.quoted()
	p.semi()
	return p.Add(s.KindImport, start, tr.Attrs("module", module), children...)
}
func (p *parser) generics() []ast.Child {
	var out []ast.Child
	if !p.Eat("<") {
		return out
	}
	for p.Err == nil && !p.EOF() && p.Peek() != ">" {
		start := p.Token().Start
		name := p.Ident()
		var children []ast.Child
		if p.Eat("extends") {
			children = append(children, tr.Child("constraint", p.typ(0)))
		}
		if p.Eat("=") {
			children = append(children, tr.Child("default", p.typ(0)))
		}
		out = append(out, tr.Child("type-parameter", p.Add(s.KindTypeParameter, start, tr.Attrs("name", name), children...)))
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
func (p *parser) declaration() ast.NodeID {
	start := p.Token().Start
	exported := p.Eat("export")
	kind := p.Take().Text
	if kind != "type" && kind != "interface" {
		p.Fail(text.Unsupported, "expected a type alias or interface")
		return 0
	}
	name := p.Ident()
	attrs := tr.Attrs("name", name)
	if exported {
		attrs = append(attrs, tr.Attrs("exported", "true")...)
	}
	children := p.generics()
	if kind == "type" {
		p.Want("=")
		children = append(children, tr.Child("type", p.typ(0)))
		p.semi()
		return p.Add(s.KindTypeAlias, start, attrs, children...)
	}
	if p.Eat("extends") {
		for p.Err == nil {
			children = append(children, tr.Child("base", p.typ(3)))
			if !p.Eat(",") {
				break
			}
		}
	}
	children = append(children, p.members()...)
	p.Eat(";")
	return p.Add(s.KindInterface, start, attrs, children...)
}
func (p *parser) parameters() []ast.Child {
	p.Want("(")
	var out []ast.Child
	for p.Err == nil && !p.EOF() && p.Peek() != ")" {
		start := p.Token().Start
		rest := p.Eat("...")
		name := p.Ident()
		optional := p.Eat("?")
		p.Want(":")
		value := p.typ(0)
		attrs := tr.Attrs("name", name)
		if rest {
			attrs = append(attrs, tr.Attrs("rest", "true")...)
		}
		if optional {
			attrs = append(attrs, tr.Attrs("optional", "true")...)
		}
		out = append(out, tr.Child("parameter", p.Add(s.KindParameter, start, attrs, tr.Child("type", value))))
		if !p.Eat(",") {
			break
		}
	}
	p.Want(")")
	return out
}
func (p *parser) signature(start int, arrow bool) ast.NodeID {
	children := p.generics()
	children = append(children, p.parameters()...)
	if arrow {
		p.Want("=>")
	} else {
		p.Want(":")
	}
	children = append(children, tr.Child("result", p.typ(0)))
	return p.Add(s.KindFunctionType, start, nil, children...)
}
func (p *parser) members() []ast.Child {
	p.Want("{")
	var out []ast.Child
	for p.Err == nil && !p.EOF() && p.Peek() != "}" {
		start := p.Token().Start
		readonly := false
		if p.Peek() == "readonly" && p.At(1) != ":" && p.At(1) != "?" && p.At(1) != "(" {
			p.Take()
			readonly = true
		}
		var id ast.NodeID
		switch p.Peek() {
		case "(", "<":
			if readonly {
				p.Fail(text.Unsupported, "readonly call signature")
				break
			}
			id = p.Add(s.KindCallSignature, start, nil, tr.Child("signature", p.signature(start, false)))
		case "[":
			p.Take()
			begin := p.Token().Start
			name := p.Ident()
			p.Want(":")
			key := p.typ(0)
			parameter := p.Add(s.KindParameter, begin, tr.Attrs("name", name), tr.Child("type", key))
			p.Want("]")
			p.Want(":")
			value := p.typ(0)
			var attrs []ast.Attribute
			if readonly {
				attrs = tr.Attrs("readonly", "true")
			}
			id = p.Add(s.KindIndexSignature, start, attrs, tr.Child("parameter", parameter), tr.Child("type", value))
		default:
			var name string
			if p.Token().Class == "string" {
				name = p.quoted()
			} else if p.Token().Class == "number" {
				name = p.Take().Text
				// Only canonical decimal keys are admitted; other spellings need
				// JavaScript ToPropertyKey semantics before becoming property text.
				if !decimalKey(name) {
					p.Fail(text.Unsupported, "numeric property keys must be canonical decimal integers")
				}
			} else {
				name = p.Ident()
			}
			optional := p.Eat("?")
			attrs := tr.Attrs("name", name)
			if optional {
				attrs = append(attrs, tr.Attrs("optional", "true")...)
			}
			if p.Peek() == "(" || p.Peek() == "<" {
				if readonly {
					p.Fail(text.Unsupported, "readonly method")
					break
				}
				id = p.Add(s.KindMethod, start, attrs, tr.Child("signature", p.signature(p.Token().Start, false)))
			} else {
				if readonly {
					attrs = append(attrs, tr.Attrs("readonly", "true")...)
				}
				p.Want(":")
				id = p.Add(s.KindProperty, start, attrs, tr.Child("type", p.typ(0)))
			}
		}
		out = append(out, tr.Child("member", id))
		if !p.Eat(",") {
			p.semi()
		}
	}
	p.Want("}")
	return out
}
func (p *parser) functionAhead() bool {
	depth := 0
	for i := p.Pos; i < len(p.Tokens); i++ {
		switch p.Tokens[i].Text {
		case "(":
			depth++
		case ")":
			depth--
			if depth == 0 {
				return i+1 < len(p.Tokens) && p.Tokens[i+1].Text == "=>"
			}
		}
	}
	return false
}
func (p *parser) typ(min int) ast.NodeID {
	if !p.Enter() {
		return 0
	}
	defer p.Leave()
	start := p.Token().Start
	left := p.primary()
	for p.Err == nil {
		op := p.Peek()
		precedence := 0
		kind := s.KindUnionType
		if op == "|" {
			precedence = 1
		} else if op == "&" {
			precedence = 2
			kind = s.KindIntersectionType
		} else {
			break
		}
		if precedence < min {
			break
		}
		p.Take()
		right := p.typ(precedence + 1)
		children := []ast.Child{}
		if left != 0 && p.Tree.Nodes[left-1].Kind == kind {
			children = append(children, p.Tree.Nodes[left-1].Children...)
		} else {
			children = append(children, tr.Child("type", left))
		}
		children = append(children, tr.Child("type", right))
		left = p.Add(kind, start, nil, children...)
	}
	return left
}
func (p *parser) primary() ast.NodeID {
	if p.Err != nil {
		return 0
	}
	start := p.Token().Start
	var id ast.NodeID
	switch p.Peek() {
	case "keyof", "readonly", "unique":
		op := p.Take().Text
		child := p.typ(3)
		return p.Add(s.KindTypeOperator, start, tr.Attrs("operator", op), tr.Child("type", child))
	case "<":
		return p.signature(start, true)
	case "(":
		if p.functionAhead() {
			return p.signature(start, true)
		}
		p.Take()
		id = p.typ(0)
		p.Want(")")
		if id != 0 {
			p.Tree.Nodes[id-1].Origin.Spans = []ast.Span{{Source: p.Input.SourceID, Start: uint64(start), End: uint64(p.Tokens[p.Pos-1].End)}}
		}
	case "{":
		id = p.Add(s.KindObjectType, start, nil, p.members()...)
	case "[":
		p.Take()
		var elements []ast.Child
		for p.Err == nil && !p.EOF() && p.Peek() != "]" {
			begin := p.Token().Start
			rest := p.Eat("...")
			var attrs []ast.Attribute
			optional := false
			if p.Token().Class == "identifier" && (p.At(1) == ":" || (p.At(1) == "?" && p.At(2) == ":")) {
				attrs = tr.Attrs("name", p.Ident())
				optional = p.Eat("?")
				p.Want(":")
			}
			value := p.typ(0)
			if len(attrs) == 0 {
				optional = p.Eat("?")
			}
			if rest {
				attrs = append(attrs, tr.Attrs("rest", "true")...)
			}
			if optional {
				attrs = append(attrs, tr.Attrs("optional", "true")...)
			}
			elements = append(elements, tr.Child("element", p.Add(s.KindTupleElement, begin, attrs, tr.Child("type", value))))
			if !p.Eat(",") {
				break
			}
		}
		p.Want("]")
		id = p.Add(s.KindTupleType, start, nil, elements...)
	default:
		token := p.Token()
		sign := ""
		if token.Text == "-" {
			sign = p.Take().Text
			token = p.Token()
		}
		if token.Class == "string" || token.Class == "number" || token.Text == "true" || token.Text == "false" || token.Text == "null" || sign != "" {
			value := sign + p.Take().Text
			kind := "number"
			if token.Class == "string" {
				kind = "string"
			} else if strings.HasSuffix(value, "n") {
				kind = "bigint"
			} else if value == "true" || value == "false" {
				kind = "boolean"
			} else if value == "null" {
				kind = "null"
			}
			id = p.Add(s.KindLiteralType, start, tr.Attrs("literal-kind", kind, "text", value))
		} else if token.Class == "identifier" {
			name := p.Ident()
			kind := s.KindTypeReference
			switch name {
			case "any", "unknown", "never", "void", "undefined", "boolean", "number", "bigint", "string", "symbol", "object", "this":
				kind = s.KindKeywordType
			case "infer", "typeof", "new", "import":
				p.Fail(text.Unsupported, "unsupported type operator")
				return 0
			}
			id = p.Add(kind, start, tr.Attrs("name", name))
			for p.Eat(".") {
				id = p.Add(s.KindQualifiedType, start, tr.Attrs("name", p.Ident()), tr.Child("qualifier", id))
			}
			if p.Eat("<") {
				children := []ast.Child{tr.Child("constructor", id)}
				for p.Err == nil && !p.EOF() && p.Peek() != ">" {
					children = append(children, tr.Child("argument", p.typ(0)))
					if !p.Eat(",") {
						break
					}
				}
				p.Want(">")
				id = p.Add(s.KindTypeApplication, start, nil, children...)
			}
		} else {
			p.Fail(text.Unsupported, "unsupported type expression")
			return 0
		}
	}
	for p.Err == nil && p.Eat("[") {
		if p.Eat("]") {
			id = p.Add(s.KindArrayType, start, nil, tr.Child("element", id))
		} else {
			index := p.typ(0)
			p.Want("]")
			id = p.Add(s.KindIndexedAccessType, start, nil, tr.Child("object", id), tr.Child("index", index))
		}
	}
	return id
}

type Provider struct{ s.Language }

func (Provider) Parse(input text.Input) (cst.CST[s.Language], error) { return Parse(input) }

var _ text.ParserProvider[s.Language] = Provider{}
