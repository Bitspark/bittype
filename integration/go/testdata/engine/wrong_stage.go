package independent

import (
	"context"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type First struct{}
type Second struct{}
type Third struct{}

func Wrong(from, to representation.Language) {
	_, _ = transform.NewTransformer[First, Second](from, to, transform.Contract{}, func(context.Context, transform.Request[First]) (ast.Set[Third], error) {
		return singleResult(ast.AST[Third]{}, nil)
	})
}

func singleResult[L any](tree ast.AST[L], err error) (ast.Set[L], error) {
	if err != nil {
		return ast.Set[L]{}, err
	}
	return ast.Set[L]{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: []ast.AST[L]{tree}}, nil
}
