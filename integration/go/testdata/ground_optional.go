package consumer

import (
	ground "github.com/Bitspark/bittype/syntax/horos/go"
	optional "github.com/Bitspark/bittype/syntax/horoswire/go"
)

// Matching declaration vocabulary does not make the Go markers interchangeable.
var _ ground.AST = optional.AST{}
