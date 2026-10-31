// Package golangprovider composes the native text capabilities explicitly.
package golangprovider

import (
	"github.com/Bitspark/bittype/trees/ast/go"
	check "github.com/Bitspark/bittype/check/golang/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	format "github.com/Bitspark/bittype/format/golang/go"
	parse "github.com/Bitspark/bittype/parse/golang/go"
	print "github.com/Bitspark/bittype/print/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/contracts/text/go"
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
