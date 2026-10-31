package consumer

import (
	"github.com/Bitspark/bittype/contracts/pipeline/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type left struct{}
type right struct{}

var stage pipeline.Stage[ast.AST[left], ast.AST[right]]
var result, failure = stage(ast.AST[left]{})
var incompatible string = result
