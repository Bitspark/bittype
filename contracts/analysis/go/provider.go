// Package analysis defines analysis with an optional explicit resolution input.
package analysis

import (
	"context"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Contract struct {
	Semantic   semantic.Contract
	Resolution *resolution.Requirement
}

type Input struct {
	Snapshot   semantic.Snapshot
	Options    semantic.Options
	Resolution *semantic.Result
}

type Request[L any] struct{ semantic.Request[L] }

func (r Request[L]) Resolution() (semantic.Result, bool) { return r.Dependency(0) }

type Provider[L any] interface {
	language.Provider
	Analyze(context.Context, Request[L]) (semantic.Result, error)
}

func Kind() semantics.Key { return semantics.Key{Owner: "bittype", Name: "analysis"} }

type Analyzer struct {
	operation  semantic.Operation
	resolution *resolution.Requirement
}

func NewAnalyzer[L any](syntax representation.Language, contract Contract, provider Provider[L]) (Analyzer, error) {
	if representation.IsNil(provider) {
		return Analyzer{}, semantic.ErrInvalidOperation
	}
	var requirement *resolution.Requirement
	if contract.Resolution != nil {
		if err := contract.Resolution.Validate(); err != nil {
			return Analyzer{}, err
		}
		copy := contract.Resolution.Clone()
		requirement = &copy
	}
	operation, err := semantic.NewOperation[L](syntax, contract.Semantic, Kind(), provider, func(ctx context.Context, request semantic.Request[L]) (semantic.Result, error) {
		return provider.Analyze(ctx, Request[L]{Request: request})
	})
	if err != nil {
		return Analyzer{}, err
	}
	return Analyzer{operation: operation, resolution: requirement}, nil
}

func (a Analyzer) ID() language.ID                   { return a.Language().ID() }
func (a Analyzer) Name() string                      { return a.Language().Name() }
func (a Analyzer) Language() representation.Language { return a.operation.Language() }
func (a Analyzer) Contract() Contract {
	c := Contract{Semantic: a.operation.Contract()}
	if a.resolution != nil {
		copy := a.resolution.Clone()
		c.Resolution = &copy
	}
	return c
}
func (a Analyzer) Validate() error { return a.operation.Validate() }
func (a Analyzer) Analyze(ctx context.Context, input Input) (semantic.Result, error) {
	var dependencies []semantic.Dependency
	if (a.resolution == nil) != (input.Resolution == nil) {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	if a.resolution != nil {
		dependencies = append(dependencies, a.resolution.Dependency(*input.Resolution))
	}
	return a.operation.Run(ctx, input.Snapshot, input.Options, dependencies)
}
