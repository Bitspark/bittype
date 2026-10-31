package independent

import (
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type First struct{}
type Second struct{}

func Wrong(definition representation.Language) {
	_, _ = representation.Pack[First](definition, ast.AST[Second]{})
}
