package consumer

import (
	"context"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/engine/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Language struct{}
type Provider struct{}

func (Provider) ID() language.ID                { return "consumer" }
func (Provider) Name() string                   { return "Consumer" }
func (Provider) Check(ast.View[Language]) error { return nil }
func (Provider) Resolve(_ context.Context, request semantic.Request[Language]) (semantic.Result, error) {
	return request.Publish(semantic.Publication{Reports: []semantic.Report{{Obligation: semantics.Key{Owner: "consumer", Name: "valid"}, Status: semantic.Complete, Verdict: semantic.Holds}}})
}

var _ resolution.Provider[Language] = Provider{}

func Compose() (*engine.Engine, error) {
	syntax, err := representation.NewLanguage[Language](Provider{}, "syntax-v1")
	if err != nil {
		return nil, err
	}
	key := semantics.Key{Owner: "consumer", Name: "valid"}
	contract := resolution.Contract{Bindings: key, Semantic: semantic.Contract{Profile: "facts-v1", Obligations: []semantic.Obligation{{Key: key, Domain: semantics.Key{Owner: "consumer", Name: "entry"}}}, ModelRequirements: []semantic.Requirement{{Obligation: key, Verdict: semantic.Holds}}}}
	resolver, err := resolution.NewResolver[Language](syntax, contract, Provider{})
	if err != nil {
		return nil, err
	}
	kernel := engine.New()
	if err := kernel.RegisterLanguage(syntax); err != nil {
		return nil, err
	}
	if err := kernel.RegisterResolver(resolver); err != nil {
		return nil, err
	}
	return kernel, nil
}
