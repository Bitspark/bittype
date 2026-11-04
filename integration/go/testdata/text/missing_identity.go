package consumer

import (
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

type Language struct{}
type MissingName struct{}

func (MissingName) ID() language.ID                              { return "missing" }
func (MissingName) Parse(parse.Input) (cst.CST[Language], error) { return cst.CST[Language]{}, nil }

var _ parse.Provider[Language] = MissingName{}

func (MissingName) Format(ast.View[Language]) (cst.CST[Language], error) {
	return cst.CST[Language]{}, nil
}
func (MissingName) Print(cst.CST[Language]) (string, error) { return "", nil }

var _ format.Provider[Language] = MissingName{}
var _ print.Provider[Language] = MissingName{}
