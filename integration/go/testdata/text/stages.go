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

func (Language) ID() language.ID { return "external-stages" }
func (Language) Name() string    { return "External stages" }

// A consumer may supply any single stage, without the other stage methods or
// the aggregate text contract.
type Parser struct{ Language }
type Formatter struct{ Language }
type Printer struct{ Language }

func (Parser) Parse(parse.Input) (cst.CST[Language], error) {
	return cst.CST[Language]{}, nil
}
func (Formatter) Format(ast.View[Language]) (cst.CST[Language], error) {
	return cst.CST[Language]{}, nil
}
func (Printer) Print(cst.CST[Language]) (string, error) { return "", nil }

var _ parse.Provider[Language] = Parser{}
var _ format.Provider[Language] = Formatter{}
var _ print.Provider[Language] = Printer{}
