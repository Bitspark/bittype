package consumer

import (
	rust "github.com/Bitspark/bittype/syntax/rust/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

var _ ast.View[ts.Language] = cst.CST[rust.Language]{}
