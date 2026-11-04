package consumer

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/pipeline/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

type Language struct{}
type Independent struct{}

func (Independent) ID() language.ID                              { return "independent-language" }
func (Independent) Name() string                                 { return "Independent language" }
func (Independent) Parse(parse.Input) (cst.CST[Language], error) { return cst.CST[Language]{}, nil }
func (Independent) Format(ast.View[Language]) (cst.CST[Language], error) {
	return cst.CST[Language]{}, nil
}
func (Independent) Check(ast.View[Language]) error          { return nil }
func (Independent) Print(cst.CST[Language]) (string, error) { return "", nil }

var _ language.Provider = Independent{}
var _ parse.Provider[Language] = Independent{}
var _ check.Provider[Language] = Independent{}
var _ format.Provider[Language] = Independent{}
var _ print.Provider[Language] = Independent{}
var _ ast.View[Language] = ast.AST[Language]{}
var _ ast.View[Language] = cst.CST[Language]{}
var _ pipeline.Stage[parse.Input, cst.CST[Language]] = Independent{}.Parse
var _ pipeline.Stage[ast.View[Language], cst.CST[Language]] = Independent{}.Format
var _ pipeline.Stage[cst.CST[Language], string] = Independent{}.Print
