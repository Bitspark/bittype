// Package cppprovider composes the native text capabilities explicitly.
package cppprovider

import (
	check "github.com/Bitspark/bittype/check/cpp/go"
	parsecontract "github.com/Bitspark/bittype/contracts/parse/go"
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	formatcontract "github.com/Bitspark/bittype/contracts/format/go"
	printcontract "github.com/Bitspark/bittype/contracts/print/go"
	format "github.com/Bitspark/bittype/format/cpp/go"
	parse "github.com/Bitspark/bittype/parse/cpp/go"
	print "github.com/Bitspark/bittype/print/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

// Provider has no initialization, discovery or registration side effects.
type Provider struct{ s.Language }

func (Provider) Parse(input parsecontract.Input) (cst.CST[s.Language], error) {
	return parse.Parse(input)
}
func (Provider) Format(view ast.View[s.Language]) (cst.CST[s.Language], error) {
	return format.Format(view)
}
func (Provider) Check(view ast.View[s.Language]) error          { return check.Check(view) }
func (Provider) Print(tree cst.CST[s.Language]) (string, error) { return print.Print(tree) }

var _ parsecontract.Provider[s.Language] = Provider{}
var _ checkcontract.Provider[s.Language] = Provider{}
var _ formatcontract.Provider[s.Language] = Provider{}
var _ printcontract.Provider[s.Language] = Provider{}
