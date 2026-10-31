package concrete

import (
	"github.com/Bitspark/bittype/ast/go"
	cstcheck "github.com/Bitspark/bittype/check/cst/go"
	"github.com/Bitspark/bittype/cst/go"
	"github.com/Bitspark/bittype/internal/lex"
	"github.com/Bitspark/bittype/internal/tree"
	"github.com/Bitspark/bittype/text/go"
	"sort"
	"strings"
)

type Vocabulary[L any] struct{ Layout, Token, Whitespace, Comment cst.Kind[L] }
type Doc[L any] struct {
	Kind   cst.Kind[L]
	Text   *string
	Syntax []ast.NodeID
	Parts  []Doc[L]
}
type Builder[L any] struct{ Vocabulary[L] }

func (b Builder[L]) T(s string) Doc[L]     { return Doc[L]{Kind: b.Token, Text: &s} }
func (b Builder[L]) Space() Doc[L]         { return b.White(" ") }
func (b Builder[L]) White(s string) Doc[L] { return Doc[L]{Kind: b.Whitespace, Text: &s} }
func (b Builder[L]) Line(indent string, level int) Doc[L] {
	return b.White("\n" + strings.Repeat(indent, level))
}
func (b Builder[L]) Seq(parts ...Doc[L]) Doc[L] { return Doc[L]{Kind: b.Layout, Parts: parts} }
func (b Builder[L]) Group(id ast.NodeID, parts ...Doc[L]) Doc[L] {
	d := b.Seq(parts...)
	if id != 0 {
		d.Syntax = []ast.NodeID{id}
	}
	return d
}
func (b Builder[L]) Join(parts []Doc[L], sep ...Doc[L]) Doc[L] {
	var all []Doc[L]
	for i, part := range parts {
		if i > 0 {
			all = append(all, sep...)
		}
		all = append(all, part)
	}
	return b.Seq(all...)
}
func Build[L any](abstract ast.AST[L], doc Doc[L]) (cst.CST[L], error) {
	out := cst.CST[L]{Abstract: tree.Clone(abstract)}
	var visit func(Doc[L]) cst.NodeID
	visit = func(d Doc[L]) cst.NodeID {
		n := cst.Node[L]{Kind: d.Kind, Syntax: append([]ast.NodeID(nil), d.Syntax...)}
		if d.Text != nil {
			s := *d.Text
			n.Text = &s
		}
		for _, part := range d.Parts {
			n.Children = append(n.Children, visit(part))
		}
		out.Nodes = append(out.Nodes, n)
		return cst.NodeID(len(out.Nodes))
	}
	out.Root = visit(doc)
	if err := cstcheck.Check(out); err != nil {
		return cst.CST[L]{}, err
	}
	return out, nil
}

// Parsed groups retained lexical leaves using native source intervals. Several
// abstract nodes may share one concrete spelling (for example grouped Go fields).
// All source bytes are retained once; no spelling is synthesized.
func Parsed[L any](abstract ast.AST[L], input text.Input, tokens []lex.Token, vocab Vocabulary[L]) (cst.CST[L], error) {
	type interval struct {
		start, end int
		kind       cst.Kind[L]
		ids        []ast.NodeID
		children   []*interval
	}
	root := &interval{start: 0, end: len(input.Text), kind: cst.Kind[L](abstract.Nodes[abstract.Root-1].Kind)}
	groups := map[[2]int]*interval{{0, len(input.Text)}: root}
	for i, n := range abstract.Nodes {
		for _, span := range n.Origin.Spans {
			if span.Source != input.SourceID {
				continue
			}
			if span.End > uint64(len(input.Text)) || span.Start > span.End {
				return cst.CST[L]{}, text.AtNode(abstract, ast.NodeID(i+1), "source span outside input")
			}
			if span.Start == span.End {
				continue
			}
			key := [2]int{int(span.Start), int(span.End)}
			g := groups[key]
			if g == nil {
				g = &interval{start: key[0], end: key[1], kind: cst.Kind[L](n.Kind)}
				groups[key] = g
			}
			g.ids = append(g.ids, ast.NodeID(i+1))
		}
	}
	var ordered []*interval
	for _, g := range groups {
		if g != root {
			ordered = append(ordered, g)
		}
	}
	sort.Slice(ordered, func(i, j int) bool {
		if ordered[i].start != ordered[j].start {
			return ordered[i].start < ordered[j].start
		}
		return ordered[i].end > ordered[j].end
	})
	stack := []*interval{root}
	for _, g := range ordered {
		for len(stack) > 1 && g.start >= stack[len(stack)-1].end {
			stack = stack[:len(stack)-1]
		}
		parent := stack[len(stack)-1]
		if g.start < parent.start || g.end > parent.end {
			return cst.CST[L]{}, &text.Diagnostic{Code: text.InvalidCST, Message: "crossing native source intervals"}
		}
		parent.children = append(parent.children, g)
		stack = append(stack, g)
	}
	out := cst.CST[L]{Abstract: abstract}
	index := 0
	var visit func(*interval) (cst.NodeID, error)
	visit = func(g *interval) (cst.NodeID, error) {
		n := cst.Node[L]{Kind: g.kind, Syntax: g.ids, Origin: ast.Origin{Spans: []ast.Span{{Source: input.SourceID, Start: uint64(g.start), End: uint64(g.end)}}}}
		leaf := func(t lex.Token) {
			kind := vocab.Token
			if t.Class == "whitespace" {
				kind = vocab.Whitespace
			}
			if t.Class == "comment" {
				kind = vocab.Comment
			}
			s := t.Text
			out.Nodes = append(out.Nodes, cst.Node[L]{Kind: kind, Text: &s, Origin: ast.Origin{Spans: []ast.Span{{Source: input.SourceID, Start: uint64(t.Start), End: uint64(t.End)}}}})
			n.Children = append(n.Children, cst.NodeID(len(out.Nodes)))
		}
		for _, child := range g.children {
			for index < len(tokens) && tokens[index].Start < child.start {
				leaf(tokens[index])
				index++
			}
			if index < len(tokens) && tokens[index].Start != child.start {
				return 0, &text.Diagnostic{Code: text.InvalidCST, Message: "native interval splits token"}
			}
			id, err := visit(child)
			if err != nil {
				return 0, err
			}
			n.Children = append(n.Children, id)
		}
		for index < len(tokens) && tokens[index].Start < g.end {
			if tokens[index].End > g.end {
				return 0, &text.Diagnostic{Code: text.InvalidCST, Message: "native interval splits token"}
			}
			leaf(tokens[index])
			index++
		}
		out.Nodes = append(out.Nodes, n)
		return cst.NodeID(len(out.Nodes)), nil
	}
	id, err := visit(root)
	if err != nil {
		return cst.CST[L]{}, err
	}
	out.Root = id
	if index != len(tokens) {
		return cst.CST[L]{}, &text.Diagnostic{Code: text.InvalidCST, Message: "unconsumed concrete tokens"}
	}
	if err := cstcheck.Check(out); err != nil {
		return cst.CST[L]{}, err
	}
	return out, nil
}
