// Package resolution defines independently registered reference resolution.
package resolution

import (
	"context"
	"slices"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type ResolverProvider[L any] interface {
	language.LanguageProvider
	Resolve(context.Context, semantic.Request[L]) (semantic.Result, error)
}

// Bindings identifies the obligation asserting definitive targets for every
// required reference in its provider-declared domain. Holds does not imply
// complete language validity. A model requires Complete/Holds for this key.
type Contract struct {
	Semantic semantic.Contract
	Bindings semantics.Key
}

func Kind() semantics.Key { return semantics.Key{Owner: "bittype", Name: "resolution"} }

func (c Contract) Validate() error {
	if err := c.Semantic.Validate(); err != nil {
		return err
	}
	if !c.Bindings.Valid() {
		return semantic.ErrInvalidOperation
	}
	for _, requirement := range c.Semantic.ModelRequirements {
		if requirement.Obligation == c.Bindings && requirement.Verdict == semantic.Holds {
			return nil
		}
	}
	return semantic.ErrInvalidOperation
}

type Requirement struct {
	Profile       string
	Configuration []semantic.Setting
	Reports       []semantic.Requirement
}

func (r Requirement) Validate() error {
	return semantic.ValidateDependencySpecification(Kind(), r.Profile, r.Configuration, r.Reports)
}
func (r Requirement) Clone() Requirement {
	r.Configuration = slices.Clone(r.Configuration)
	for i := range r.Configuration {
		r.Configuration[i].Value = slices.Clone(r.Configuration[i].Value)
	}
	r.Reports = slices.Clone(r.Reports)
	return r
}
func (r Requirement) Dependency(result semantic.Result) semantic.Dependency {
	return semantic.Dependency{Result: result, Operation: Kind(), Profile: r.Profile, Configuration: r.Configuration, Reports: r.Reports}
}

type Resolver struct {
	operation semantic.Operation
	bindings  semantics.Key
}

func NewResolver[L any](syntax representation.Language, contract Contract, provider ResolverProvider[L]) (Resolver, error) {
	if err := contract.Validate(); err != nil {
		return Resolver{}, err
	}
	if representation.IsNil(provider) {
		return Resolver{}, semantic.ErrInvalidOperation
	}
	operation, err := semantic.NewOperation[L](syntax, contract.Semantic, Kind(), provider, provider.Resolve)
	if err != nil {
		return Resolver{}, err
	}
	return Resolver{operation: operation, bindings: contract.Bindings}, nil
}

func (r Resolver) ID() language.ID                   { return r.Language().ID() }
func (r Resolver) Name() string                      { return r.Language().Name() }
func (r Resolver) Language() representation.Language { return r.operation.Language() }
func (r Resolver) Contract() Contract {
	return Contract{Semantic: r.operation.Contract(), Bindings: r.bindings}
}
func (r Resolver) Validate() error { return r.operation.Validate() }
func (r Resolver) Resolve(ctx context.Context, source semantic.Snapshot, options semantic.Options) (semantic.Result, error) {
	return r.operation.Run(ctx, source, options, nil)
}
