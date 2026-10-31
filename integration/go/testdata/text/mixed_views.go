package consumer

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/cst/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	rust "github.com/Bitspark/bittype/syntax/rust/go"
)
var _ ast.View[ts.Language]=cst.CST[rust.Language]{}
