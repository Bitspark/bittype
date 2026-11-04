package consumer

import (
	"context"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

type Language struct{}
type MissingName struct{}

func (MissingName) ID() language.ID { return "missing" }
func (MissingName) Parse(context.Context, parse.Input) (cst.CST[Language], error) {
	return cst.CST[Language]{}, nil
}

var _ parse.Provider[parse.Input, cst.CST[Language]] = MissingName{}

func (MissingName) Format(context.Context, ast.View[Language]) (cst.CST[Language], error) {
	return cst.CST[Language]{}, nil
}
func (MissingName) Print(context.Context, cst.CST[Language]) (string, error) { return "", nil }

var _ format.Provider[ast.View[Language], cst.CST[Language]] = MissingName{}
var _ print.Provider[cst.CST[Language], string] = MissingName{}
