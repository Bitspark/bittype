// Package cppprovider composes the native text capabilities explicitly.
package cppprovider

import (
	check "github.com/Bitspark/bittype/check/cpp/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	format "github.com/Bitspark/bittype/format/cpp/go"
	parse "github.com/Bitspark/bittype/parse/cpp/go"
	print "github.com/Bitspark/bittype/print/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

// Provider has no initialization, discovery or registration side effects.
type Provider struct{ s.Language }

func (Provider) Parse(input text.Input) (cst.CST[s.Language], error) { return parse.Parse(input) }
func (Provider) Format(view ast.View[s.Language]) (cst.CST[s.Language], error) {
	return format.Format(view)
}
func (Provider) Check(view ast.View[s.Language]) error          { return check.Check(view) }
func (Provider) Print(tree cst.CST[s.Language]) (string, error) { return print.Print(tree) }

var _ text.Provider[s.Language] = Provider{}
