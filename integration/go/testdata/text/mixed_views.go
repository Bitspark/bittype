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

var parser parse.Provider[parse.Input, cst.CST[ts.Language]]
var formatter format.Provider[ast.View[ts.Language], cst.CST[ts.Language]]
var printer print.Provider[cst.CST[ts.Language], string]

var _ parse.Provider[parse.Input, cst.CST[rust.Language]] = parser
var _ format.Provider[ast.View[rust.Language], cst.CST[rust.Language]] = formatter
var _ print.Provider[cst.CST[rust.Language], string] = printer
