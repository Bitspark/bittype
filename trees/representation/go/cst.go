package representation

import (
	cstcheck "github.com/Bitspark/bittype/check/cst/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/internal/tree"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

// CST is a checked concrete snapshot with an already constructed AST view.
// Concrete storage is opaque; UnpackCST returns an independent typed copy.
type CST struct {
	abstract AST
	concrete any
}

func (c CST) LanguageID() language.ID { return c.abstract.LanguageID() }
func (c CST) Profile() string         { return c.abstract.Profile() }
func (c CST) ASTView() AST            { return c.abstract }

func PackCST[L any](definition Language, input cst.CST[L]) (CST, error) {
	if err := cstcheck.Check(input); err != nil {
		return CST{}, err
	}
	abstract, err := Pack(definition, input)
	if err != nil {
		return CST{}, err
	}
	return CST{abstract: abstract, concrete: cloneConcrete(input)}, nil
}

func UnpackCST[L any](definition Language, input CST) (cst.CST[L], error) {
	if _, err := Unpack[L](definition, input.abstract); err != nil {
		return cst.CST[L]{}, err
	}
	concrete, ok := input.concrete.(cst.CST[L])
	if !ok {
		return cst.CST[L]{}, ErrTypeMismatch
	}
	return cloneConcrete(concrete), nil
}
func cloneConcrete[L any](input cst.CST[L]) cst.CST[L] {
	result := input
	result.Abstract = tree.Clone(input.Abstract)
	result.Nodes = make([]cst.Node[L], len(input.Nodes))
	for i, n := range input.Nodes {
		result.Nodes[i] = n
		if n.Text != nil {
			value := *n.Text
			result.Nodes[i].Text = &value
		}
		result.Nodes[i].Children = append([]cst.NodeID(nil), n.Children...)
		result.Nodes[i].Syntax = append([]ast.NodeID(nil), n.Syntax...)
		result.Nodes[i].Origin.Spans = append([]ast.Span(nil), n.Origin.Spans...)
		result.Nodes[i].Origin.Nodes = append([]ast.NodeRef(nil), n.Origin.Nodes...)
	}
	return result
}
