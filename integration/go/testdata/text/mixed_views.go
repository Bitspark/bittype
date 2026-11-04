package consumer

import (
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	rust "github.com/Bitspark/bittype/syntax/rust/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

var _ ast.View[ts.Language] = cst.CST[rust.Language]{}

var parser parse.Provider[ts.Language]
var formatter format.Provider[ts.Language]
var printer print.Provider[ts.Language]

var _ parse.Provider[rust.Language] = parser
var _ format.Provider[rust.Language] = formatter
var _ print.Provider[rust.Language] = printer
