package consumer

import "github.com/Bitspark/bittype/ast/go"

type left struct{}
type right struct{}

var incompatible ast.Node[left] = ast.Node[right]{}
