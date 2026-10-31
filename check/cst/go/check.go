// Package cstcheck checks concrete carrier structure, not token grammar.
package cstcheck

import (
	"github.com/Bitspark/bittype/check/go"
	"github.com/Bitspark/bittype/cst/go"
	"github.com/Bitspark/bittype/text/go"
	"unicode/utf8"
)

func Check[L any](tree cst.CST[L]) error {
	if err := check.Structure(tree.Abstract); err != nil {
		return err
	}
	bad := func(message string) error { return &text.Diagnostic{Code: text.InvalidCST, Message: message} }
	if tree.Root == 0 || int(tree.Root) > len(tree.Nodes) {
		return bad("invalid concrete root")
	}
	if len(tree.Nodes) > check.MaxNodes {
		return &text.Diagnostic{Code: text.LimitExceeded, Message: "too many concrete nodes"}
	}
	bytes, entries := check.MaxTextBytes, check.MaxEntries
	take := func(left *int, count int) bool {
		if count > *left {
			return false
		}
		*left -= count
		return true
	}
	limit := func() error { return &text.Diagnostic{Code: text.LimitExceeded, Message: "concrete storage exceeds resource limit"} }
	// Preflight before UTF-8 scanning or the printer allocates expanded output.
	for _, n := range tree.Nodes {
		if !take(&bytes, len(n.Kind)) || !take(&entries, len(n.Children)) || !take(&entries, len(n.Syntax)) || !take(&entries, len(n.Origin.Spans)) || !take(&entries, len(n.Origin.Nodes)) {
			return limit()
		}
		if n.Text != nil && !take(&bytes, len(*n.Text)) {
			return limit()
		}
		for _, s := range n.Origin.Spans {
			if !take(&bytes, len(s.Source)) {
				return limit()
			}
		}
		for _, r := range n.Origin.Nodes {
			if !take(&bytes, len(r.Tree)) {
				return limit()
			}
		}
	}
	seen := make([]bool, len(tree.Nodes))
	type entry struct {
		id    cst.NodeID
		depth int
	}
	stack := []entry{{tree.Root, 0}}
	for len(stack) > 0 {
		e := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if e.id == 0 || int(e.id) > len(tree.Nodes) {
			return bad("concrete child outside arena")
		}
		if seen[e.id-1] {
			return bad("concrete cycle or multiple parents")
		}
		if e.depth > check.MaxDepth+2 {
			return &text.Diagnostic{Code: text.LimitExceeded, Message: "concrete nesting exceeds limit"}
		}
		seen[e.id-1] = true
		n := tree.Nodes[e.id-1]
		if n.Kind == "" {
			return bad("empty concrete kind")
		}
		if n.Text != nil && (len(n.Children) != 0 || !utf8.ValidString(*n.Text) || *n.Text == "") {
			return bad("invalid concrete leaf")
		}
		for _, id := range n.Syntax {
			if id == 0 || int(id) > len(tree.Abstract.Nodes) {
				return bad("invalid abstract view address")
			}
		}
		for _, s := range n.Origin.Spans {
			if s.Source == "" || s.Start > s.End {
				return bad("invalid concrete source span")
			}
		}
		for _, r := range n.Origin.Nodes {
			if r.Tree == "" || r.Node == 0 {
				return bad("invalid concrete origin address")
			}
		}
		for _, id := range n.Children {
			stack = append(stack, entry{id, e.depth + 1})
		}
	}
	for _, ok := range seen {
		if !ok {
			return bad("unreachable concrete node")
		}
	}
	return nil
}
