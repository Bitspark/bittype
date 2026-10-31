package independent

import "github.com/Bitspark/bittype/transform/go"

type Incomplete struct{}

func (Incomplete) Transform(a transform.AST) (transform.AST, error) { return a, nil }

var _ transform.ASTTransformer = Incomplete{}
