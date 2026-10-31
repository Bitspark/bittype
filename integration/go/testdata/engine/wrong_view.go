package independent

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/transform/go"
)

type First struct{}
type Second struct{}

func Wrong(definition transform.Language) {
	_, _ = transform.Pack[First](definition, ast.AST[Second]{})
}
