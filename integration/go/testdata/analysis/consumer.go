package consumer

import (
	"context"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/language/go"
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
func (Provider) Analyze(_ context.Context, request analysis.Request[Language]) (analysis.Result, error) {
	return request.Publish(analysis.Publication{Reports: []analysis.Report{{Obligation: semantics.Key{Owner: "consumer", Name: "valid"}, Status: analysis.Complete, Verdict: analysis.Holds}}})
}

var _ analysis.Provider[Language] = Provider{}

func Compose() (*engine.Engine, error) {
	syntax, err := representation.NewLanguage[Language](Provider{}, "syntax-v1")
	if err != nil {
		return nil, err
	}
	key := semantics.Key{Owner: "consumer", Name: "valid"}
	contract := analysis.Contract{Profile: "facts-v1", Obligations: []analysis.Obligation{{Key: key, Domain: semantics.Key{Owner: "consumer", Name: "entry"}}}}
	analyzer, err := analysis.NewAnalyzer[Language](syntax, contract, Provider{})
	if err != nil {
		return nil, err
	}
	kernel := engine.New()
	if err := kernel.RegisterLanguage(syntax); err != nil {
		return nil, err
	}
	if err := kernel.RegisterAnalyzer(analyzer); err != nil {
		return nil, err
	}
	return kernel, nil
}
