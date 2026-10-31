package text

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/cst/go"
	"github.com/Bitspark/bittype/language/go"
)

// Every specialized provider extends the common identity contract.
type ParserProvider[L any] interface {
	language.LanguageProvider
	Parse(Input) (cst.CST[L], error)
}
type FormatterProvider[L any] interface {
	language.LanguageProvider
	Format(ast.View[L]) (cst.CST[L], error)
}
type CheckerProvider[L any] interface {
	language.LanguageProvider
	Check(ast.View[L]) error
}
type PrinterProvider[L any] interface {
	language.LanguageProvider
	Print(cst.CST[L]) (string, error)
}
type Provider[L any] interface {
	ParserProvider[L]
	FormatterProvider[L]
	CheckerProvider[L]
	PrinterProvider[L]
}
