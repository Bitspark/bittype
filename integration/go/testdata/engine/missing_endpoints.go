package independent

import (
	"context"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Incomplete struct{}

func (Incomplete) Transform(context.Context, transform.Input) (representation.ASTSet, error) {
	return representation.ASTSet{}, nil
}

var _ transform.ASTTransformer = Incomplete{}
