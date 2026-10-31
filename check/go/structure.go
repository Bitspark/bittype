// Package check validates carrier structure without knowing any native grammar.
package check

import (
	"fmt"
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/text/go"
)

// MaxDepth bounds recursive consumers; it is a resource limit, not a type rule.
const MaxDepth = 256
const MaxNodes = 200000
const MaxTextBytes = 16 << 20
const MaxEntries = MaxNodes * 8

// budget counts aggregate storage before consumers scan or copy text. Repeated
// references to the same backing string still count once per occurrence.
func budget[L any](tree ast.AST[L]) error {
	bytes, entries := MaxTextBytes, MaxEntries
	take := func(left *int, count int) bool {
		if count > *left {
			return false
		}
		*left -= count
		return true
	}
	bad := func() error { return &text.Diagnostic{Code: text.LimitExceeded, Message: "AST storage exceeds resource limit"} }
	if !take(&bytes, len(tree.ID)) {
		return bad()
	}
	for _, n := range tree.Nodes {
		if !take(&bytes, len(n.Kind)) || !take(&entries, len(n.Attributes)) || !take(&entries, len(n.Children)) || !take(&entries, len(n.References)) || !take(&entries, len(n.Origin.Spans)) || !take(&entries, len(n.Origin.Nodes)) {
			return bad()
		}
		for _, a := range n.Attributes {
			if !take(&bytes, len(a.Name)) || !take(&bytes, len(a.Text)) {
				return bad()
			}
		}
		for _, c := range n.Children {
			if !take(&bytes, len(c.Role)) {
				return bad()
			}
		}
		for _, r := range n.References {
			if !take(&bytes, len(r.Role)) || !take(&bytes, len(r.Target.Tree)) {
				return bad()
			}
		}
		for _, s := range n.Origin.Spans {
			if !take(&bytes, len(s.Source)) {
				return bad()
			}
		}
		for _, r := range n.Origin.Nodes {
			if !take(&bytes, len(r.Tree)) {
				return bad()
			}
		}
	}
	return nil
}

func Structure[L any](tree ast.AST[L]) error {
	bad := func(id ast.NodeID, message string) error { return text.AtNode(tree, id, message) }
	if tree.ID == "" || tree.Root == 0 || int(tree.Root) > len(tree.Nodes) {
		return bad(tree.Root, "nonempty tree ID and valid root are required")
	}
	if len(tree.Nodes) > MaxNodes {
		return &text.Diagnostic{Code: text.LimitExceeded, Message: "too many AST nodes"}
	}
	if err := budget(tree); err != nil {
		return err
	}
	seen := make([]bool, len(tree.Nodes))
	type entry struct {
		id    ast.NodeID
		depth int
	}
	stack := []entry{{tree.Root, 0}}
	for len(stack) > 0 {
		e := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if e.id == 0 || int(e.id) > len(tree.Nodes) {
			return bad(e.id, "child address outside arena")
		}
		if e.depth > MaxDepth {
			return &text.Diagnostic{Code: text.LimitExceeded, Message: "AST nesting exceeds limit"}
		}
		if seen[e.id-1] {
			return bad(e.id, "containment cycle or multiple parents")
		}
		seen[e.id-1] = true
		n := tree.Nodes[e.id-1]
		if n.Kind == "" {
			return bad(e.id, "empty node kind")
		}
		attrs := map[string]bool{}
		for _, a := range n.Attributes {
			if a.Name == "" || attrs[a.Name] {
				return bad(e.id, "empty or duplicate attribute: "+a.Name)
			}
			attrs[a.Name] = true
		}
		refs := map[string]bool{}
		for _, r := range n.References {
			if r.Role == "" || refs[r.Role] || r.Target.Tree == "" || r.Target.Node == 0 ||
				(r.Target.Tree == tree.ID && int(r.Target.Node) > len(tree.Nodes)) {
				return bad(e.id, "invalid or duplicate reference")
			}
			refs[r.Role] = true
		}
		for _, s := range n.Origin.Spans {
			if s.Source == "" || s.Start > s.End {
				return bad(e.id, "invalid source span")
			}
		}
		for _, r := range n.Origin.Nodes {
			if r.Tree == "" || r.Node == 0 {
				return bad(e.id, "invalid origin node address")
			}
		}
		for _, child := range n.Children {
			if child.Role == "" {
				return bad(e.id, "empty child role")
			}
			stack = append(stack, entry{child.Node, e.depth + 1})
		}
	}
	for i, ok := range seen {
		if !ok {
			return bad(ast.NodeID(i+1), fmt.Sprintf("unreachable node %d", i+1))
		}
	}
	return nil
}
