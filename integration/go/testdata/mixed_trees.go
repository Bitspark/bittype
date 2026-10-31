package consumer

import "github.com/Bitspark/bittype/trees/ast/go"

type left struct{}
type right struct{}

var incompatible ast.AST[left] = ast.AST[right]{}
