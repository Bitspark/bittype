package parseutil

import (
	"fmt"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/check/ast/go"
	"github.com/Bitspark/bittype/internal/lex"
	"github.com/Bitspark/bittype/contracts/text/go"
)

type Parser[L any] struct {
	Input      text.Input
	All        []lex.Token
	Tokens     []lex.Token
	Pos, Depth int
	Tree       ast.AST[L]
	Err        error
}

func New[L any](input text.Input, options lex.Options) (*Parser[L], error) {
	all, err := lex.Scan(input, options)
	if err != nil {
		return nil, err
	}
	p := &Parser[L]{Input: input, All: all, Tree: ast.AST[L]{ID: input.TreeID}}
	for _, t := range all {
		if t.Class != "whitespace" && t.Class != "comment" {
			p.Tokens = append(p.Tokens, t)
		}
	}
	p.Tokens = append(p.Tokens, lex.Token{Class: "eof", Start: len(input.Text), End: len(input.Text)})
	return p, nil
}
func (p *Parser[L]) Peek() string     { return p.Tokens[p.Pos].Text }
func (p *Parser[L]) Token() lex.Token { return p.Tokens[p.Pos] }
func (p *Parser[L]) At(n int) string {
	if p.Pos+n >= len(p.Tokens) {
		return ""
	}
	return p.Tokens[p.Pos+n].Text
}
func (p *Parser[L]) EOF() bool { return p.Token().Class == "eof" }
func (p *Parser[L]) Take() lex.Token {
	t := p.Token()
	if !p.EOF() {
		p.Pos++
	}
	return t
}
func (p *Parser[L]) Eat(s string) bool {
	if p.Peek() != s || p.EOF() {
		return false
	}
	p.Take()
	return true
}
func (p *Parser[L]) Want(s string) {
	if p.Err == nil && !p.Eat(s) {
		p.Fail(text.SyntaxError, "expected "+s)
	}
}
func (p *Parser[L]) Fail(code text.Code, message string) {
	if p.Err != nil {
		return
	}
	t := p.Token()
	p.Err = &text.Diagnostic{Code: code, Message: message, Span: &ast.Span{Source: p.Input.SourceID, Start: uint64(t.Start), End: uint64(t.End)}}
}
func (p *Parser[L]) Ident() string {
	t := p.Token()
	if t.Class != "identifier" {
		p.Fail(text.SyntaxError, "expected identifier")
		return ""
	}
	p.Take()
	return t.Text
}
func (p *Parser[L]) Enter() bool {
	if p.Err != nil {
		return false
	}
	p.Depth++
	if p.Depth > check.MaxDepth/2 {
		p.Fail(text.LimitExceeded, "syntax nesting exceeds limit")
		p.Depth--
		return false
	}
	return true
}
func (p *Parser[L]) Leave() { p.Depth-- }
func (p *Parser[L]) Add(kind ast.Kind[L], start int, attrs []ast.Attribute, children ...ast.Child) ast.NodeID {
	end := start
	if p.Pos > 0 {
		end = p.Tokens[p.Pos-1].End
	}
	return p.AddSpan(kind, start, end, attrs, children...)
}
func (p *Parser[L]) AddSpan(kind ast.Kind[L], start, end int, attrs []ast.Attribute, children ...ast.Child) ast.NodeID {
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	p.Tree.Nodes = append(p.Tree.Nodes, ast.Node[L]{Kind: kind, Attributes: attrs, Children: children,
		Origin: ast.Origin{Spans: []ast.Span{{Source: p.Input.SourceID, Start: uint64(start), End: uint64(end)}}}})
	return ast.NodeID(len(p.Tree.Nodes))
}
func (p *Parser[L]) Set(id ast.NodeID, name, value string) {
	if id == 0 {
		return
	}
	n := &p.Tree.Nodes[id-1]
	for i := range n.Attributes {
		if n.Attributes[i].Name == name {
			n.Attributes[i].Text = value
			return
		}
	}
	n.Attributes = append(n.Attributes, ast.Attribute{Name: name, Text: value})
}
func (p *Parser[L]) End() error {
	if p.Err == nil && !p.EOF() {
		p.Fail(text.Unsupported, fmt.Sprintf("unsupported syntax %q", p.Peek()))
	}
	return p.Err
}

// Compact drops speculative nodes and assigns fresh addresses. Parsers do not
// publish resolved references; each retained subtree gets independent ownership.
func (p *Parser[L]) Compact(root ast.NodeID) ast.AST[L] {
	out := ast.AST[L]{ID: p.Tree.ID}
	depth := 0
	var copyNode func(ast.NodeID) ast.NodeID
	copyNode = func(id ast.NodeID) ast.NodeID {
		if id == 0 || p.Err != nil {
			return 0
		}
		depth++
		defer func() { depth-- }()
		if depth > check.MaxDepth {
			p.Fail(text.LimitExceeded, "abstract nesting exceeds limit")
			return 0
		}
		n := p.Tree.Nodes[id-1]
		n.Children = append([]ast.Child(nil), n.Children...)
		n.Attributes = append([]ast.Attribute(nil), n.Attributes...)
		n.Origin.Spans = append([]ast.Span(nil), n.Origin.Spans...)
		for i := range n.Children {
			n.Children[i].Node = copyNode(n.Children[i].Node)
		}
		out.Nodes = append(out.Nodes, n)
		return ast.NodeID(len(out.Nodes))
	}
	out.Root = copyNode(root)
	return out
}
