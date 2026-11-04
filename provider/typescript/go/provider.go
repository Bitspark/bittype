// Package typescriptprovider composes the native text capabilities explicitly.
package typescriptprovider

import (
	"context"
	check "github.com/Bitspark/bittype/check/typescript/go"
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	formatcontract "github.com/Bitspark/bittype/contracts/format/go"
	parsecontract "github.com/Bitspark/bittype/contracts/parse/go"
	printcontract "github.com/Bitspark/bittype/contracts/print/go"
	format "github.com/Bitspark/bittype/format/typescript/go"
	parse "github.com/Bitspark/bittype/parse/typescript/go"
	print "github.com/Bitspark/bittype/print/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

// Provider has no initialization, discovery or registration side effects.
type Provider struct{ s.Language }

func (Provider) Parse(ctx context.Context, input parsecontract.Input) (cst.CST[s.Language], error) {
	return parse.Parse(ctx, input)
}
func (Provider) Format(ctx context.Context, view ast.View[s.Language]) (cst.CST[s.Language], error) {
	return format.Format(ctx, view)
}
func (Provider) Check(view ast.View[s.Language]) error { return check.Check(view) }
func (Provider) Print(ctx context.Context, tree cst.CST[s.Language]) (string, error) {
	return print.Print(ctx, tree)
}

var _ parsecontract.Provider[parsecontract.Input, cst.CST[s.Language]] = Provider{}
var _ checkcontract.Provider[s.Language] = Provider{}
var _ formatcontract.Provider[ast.View[s.Language], cst.CST[s.Language]] = Provider{}
var _ printcontract.Provider[cst.CST[s.Language], string] = Provider{}
