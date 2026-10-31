package consumer

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/cst/go"
	rust "github.com/Bitspark/bittype/syntax/rust/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
)

var _ ast.View[ts.Language] = cst.CST[rust.Language]{}
