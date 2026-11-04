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

func (Language) ID() language.ID { return "external-stages" }
func (Language) Name() string    { return "External stages" }

// A consumer may supply any single stage, without the other stage methods or
// the aggregate text contract.
type Parser struct{ Language }
type Formatter struct{ Language }
type Printer struct{ Language }

func (Parser) Parse(context.Context, parse.Input) (cst.CST[Language], error) {
	return cst.CST[Language]{}, nil
}
func (Formatter) Format(context.Context, ast.View[Language]) (cst.CST[Language], error) {
	return cst.CST[Language]{}, nil
}
func (Printer) Print(context.Context, cst.CST[Language]) (string, error) { return "", nil }

var _ parse.Provider[parse.Input, cst.CST[Language]] = Parser{}
var _ format.Provider[ast.View[Language], cst.CST[Language]] = Formatter{}
var _ print.Provider[cst.CST[Language], string] = Printer{}
