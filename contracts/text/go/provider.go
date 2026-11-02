package text

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

// Every specialized provider extends the common identity contract.
type Parser[L any] interface {
	language.Provider
	Parse(Input) (cst.CST[L], error)
}
type Formatter[L any] interface {
	language.Provider
	Format(ast.View[L]) (cst.CST[L], error)
}
type Checker[L any] interface {
	language.Provider
	Check(ast.View[L]) error
}
type Printer[L any] interface {
	language.Provider
	Print(cst.CST[L]) (string, error)
}
type Provider[L any] interface {
	Parser[L]
	Formatter[L]
	Checker[L]
	Printer[L]
}
