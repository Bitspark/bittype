// Package provider composes independently registered Horos syntax and semantic stages.
package provider

import (
	"context"

	analyze "github.com/Bitspark/bittype/analyze/horos/go"
	check "github.com/Bitspark/bittype/check/horos/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	formatcontract "github.com/Bitspark/bittype/contracts/format/go"
	parsecontract "github.com/Bitspark/bittype/contracts/parse/go"
	printcontract "github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	format "github.com/Bitspark/bittype/format/horos/go"
	parse "github.com/Bitspark/bittype/parse/horos/go"
	print "github.com/Bitspark/bittype/print/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

type Provider struct{ s.Language }

func (Provider) Parse(input parsecontract.Input) (cst.CST[s.Language], error) {
	return parse.Parse(input)
}
func (Provider) Format(view ast.View[s.Language]) (cst.CST[s.Language], error) {
	return format.Format(view)
}
func (Provider) Print(tree cst.CST[s.Language]) (string, error) { return print.Print(tree) }

var _ parsecontract.Provider[s.Language] = Provider{}
var _ formatcontract.Provider[s.Language] = Provider{}
var _ printcontract.Provider[s.Language] = Provider{}

func (Provider) Check(view ast.View[s.Language]) error { return check.Check(view) }
func (Provider) Analyze(ctx context.Context, request analysis.Request[s.Language]) (semantic.Result, error) {
	return analyze.Analyze(ctx, request)
}

var _ checkcontract.Provider[s.Language] = Provider{}
var _ analysis.Provider[s.Language] = Provider{}

func (Provider) Resolve(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	return resolve.ResolveReferences(ctx, request)
}

var _ resolution.Provider[s.Language] = Provider{}
