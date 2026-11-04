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
	_, _ = transform.NewTransformer[First, Second](from, to, transform.Contract{}, func(context.Context, transform.Request[First]) (ast.AST[Third], error) { return ast.AST[Third]{}, nil })
}
