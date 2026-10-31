package independent

import (
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Incomplete struct{}

func (Incomplete) Transform(a representation.AST) (representation.AST, error) { return a, nil }

var _ transform.ASTTransformer = Incomplete{}
