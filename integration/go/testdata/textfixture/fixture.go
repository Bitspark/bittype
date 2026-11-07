// Package textfixture groups explicitly supplied stages for conformance tests.
// Production consumers register the individual stage adapters.
package textfixture

import (
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	printing "github.com/Bitspark/bittype/print/syntax/shared/provider/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

type Parser[L any] interface {
	parse.Provider[parse.Input, cst.CST[L]]
}
type Checker[L any] interface{ check.Provider[L] }
type Formatter[L any] interface {
	format.Provider[ast.View[L], cst.CST[L]]
}
type Printer[L any] interface {
	print.Provider[cst.CST[L], string]
}
type Fixture[L any] struct {
	Parser[L]
	Checker[L]
	Formatter[L]
	Printer[L]
}

func (f Fixture[L]) ID() language.ID { return f.Parser.ID() }
func (f Fixture[L]) Name() string    { return f.Parser.Name() }
func New[L any](p Parser[L], c Checker[L], f Formatter[L]) Fixture[L] {
	return Fixture[L]{p, c, f, Printing[L](p)}
}
func Printing[L any](identity language.Provider) printing.Provider[L] {
	p, err := printing.New[L](identity)
	if err != nil {
		panic(err)
	}
	return p
}
