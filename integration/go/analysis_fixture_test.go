package integration_test

import (
	"context"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/engine/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Test construction storage keeps the existing input-mutation cases explicit.
// Each helper calls the actual public capture and operation APIs.
type analysisCase struct {
	Entry         ast.NodeRef
	Trees         []representation.AST
	Requested     []semantics.Key
	Configuration []semantic.Setting
	Limits        semantic.Limits
}

func (i analysisCase) source() semantic.Input {
	return semantic.Input{Entry: i.Entry, Trees: i.Trees, Limits: i.Limits}
}
func (i analysisCase) options() semantic.Options {
	return semantic.Options{Requested: i.Requested, Configuration: i.Configuration, Limits: i.Limits}
}
func analyzeCase(e *engine.Engine, ctx context.Context, input analysisCase) (semantic.Result, error) {
	source, err := e.Capture(ctx, input.source())
	if err != nil {
		return semantic.Result{}, err
	}
	return e.Analyze(ctx, analysis.Input{Snapshot: source, Options: input.options()})
}
func runAnalyzerCase(a analysis.Analyzer, ctx context.Context, input analysisCase) (semantic.Result, error) {
	source, err := semantic.Capture(ctx, a.Language(), input.source())
	if err != nil {
		return semantic.Result{}, err
	}
	return a.Analyze(ctx, analysis.Input{Snapshot: source, Options: input.options()})
}
