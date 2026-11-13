package horosmodel

import (
	"github.com/Bitspark/bittype/contracts/text/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
)

// SourceModule associates an admitted module's syntax with its canonical value.
// The value retains unresolved references; it is not a resolved Program.
type SourceModule struct {
	Node  ast.NodeID
	Value core.Value
}

// SourceModules lowers admitted syntax in source order without resolving names
// or evaluating types and kinds. Both module output and source pins use this path.
func SourceModules[L any](tree ast.AST[L]) ([]SourceModule, error) {
	r := reader[L]{tree: tree, declarationPrefix: s.DeclarationTagPrefix, descriptorPrefix: s.DescriptorTagPrefix}
	out := make([]SourceModule, 0, len(tree.Nodes[tree.Root-1].Children))
	for _, child := range tree.Nodes[tree.Root-1].Children {
		value, err := r.module(child.Node)
		if err != nil {
			return nil, err
		}
		out = append(out, SourceModule{Node: child.Node, Value: value})
	}
	return out, nil
}

// SourcePins encodes the admitted syntax, without resolving or validating type
// expressions. Pin is only a bounded hash of the exact canonical value bytes.
// The parser assigns these pins; the formatter checks, never repairs, them.
func SourcePins[L any](tree ast.AST[L]) (map[ast.NodeID]string, error) {
	modules, err := SourceModules(tree)
	if err != nil {
		return nil, err
	}
	out := make(map[ast.NodeID]string, len(modules))
	for _, module := range modules {
		pin, err := decl.Pin(module.Value)
		if err != nil {
			return nil, text.AtNode(tree, module.Node, err.Error())
		}
		out[module.Node] = pin.String()
	}
	return out, nil
}

// SourceLiteral is intentionally restricted to literal Atom/Tuple fragments.
// It cannot interpret a parameter-shaped literal as a declaration expression.
func SourceLiteral[L any](tree ast.AST[L], id ast.NodeID) (core.Value, error) {
	var literal func(ast.NodeID) bool
	literal = func(id ast.NodeID) bool {
		n := tree.Nodes[id-1]
		if s.Kind(n.Kind) != s.KindAtom && s.Kind(n.Kind) != s.KindTuple {
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
