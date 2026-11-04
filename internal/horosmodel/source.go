package horosmodel

import (
	"github.com/Bitspark/bittype/contracts/text/go"
	s "github.com/Bitspark/bittype/internal/declsyntax"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
)

// SourcePins encodes the admitted syntax, without resolving or validating type
// expressions. Pin is only a bounded hash of the exact canonical value bytes.
// The parser assigns these pins; the formatter checks, never repairs, them.
func SourcePins[L any](tree ast.AST[L]) (map[ast.NodeID]string, error) {
	r := reader[L]{tree: tree, declarationPrefix: "horos/decl/", descriptorPrefix: "horos/"}
	out := make(map[ast.NodeID]string)
	for _, child := range tree.Nodes[tree.Root-1].Children {
		value, err := r.module(child.Node)
		if err != nil {
			return nil, err
		}
		pin, err := decl.Pin(value)
		if err != nil {
			return nil, text.AtNode(tree, child.Node, err.Error())
		}
		out[child.Node] = pin.String()
	}
	return out, nil
}

// SourceLiteral is intentionally restricted to literal Atom/Tuple fragments.
// It cannot interpret a parameter-shaped literal as a declaration expression.
func SourceLiteral[L any](tree ast.AST[L], id ast.NodeID) (core.Value, error) {
	var literal func(ast.NodeID) bool
	literal = func(id ast.NodeID) bool {
		n := tree.Nodes[id-1]
		if n.Kind != s.KindAtom && n.Kind != s.KindTuple {
			return false
		}
		for _, child := range n.Children {
			if !literal(child.Node) {
				return false
			}
		}
		return true
	}
	if !literal(id) {
		return nil, text.AtNode(tree, id, "expected a literal Atom/Tuple value")
	}
	return (reader[L]{tree: tree}).value(id)
}
