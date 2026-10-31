package consumer

import (
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/contracts/pipeline/go"
)

type left struct{}
type right struct{}

var stage pipeline.Stage[ast.AST[right], string]
var result, failure = stage(ast.AST[left]{})
