package integration_test

import (
	"context"
	"errors"
	"testing"

	analyzegolang "github.com/Bitspark/bittype/analyze/files/golang/provider/go"
	analyzetypescript "github.com/Bitspark/bittype/analyze/files/typescript/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	resolution "github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	resolvegolang "github.com/Bitspark/bittype/resolve/files/golang/provider/go"
	resolvetypescript "github.com/Bitspark/bittype/resolve/files/typescript/provider/go"
	"github.com/Bitspark/bittype/semantics/go"
	golang "github.com/Bitspark/bittype/semantics/golang/go"
	typescript "github.com/Bitspark/bittype/semantics/typescript/go"
	golangsyntax "github.com/Bitspark/bittype/syntax/golang/go"
	typescriptsyntax "github.com/Bitspark/bittype/syntax/typescript/go"
)

// The Go and the TypeScript profile have their second version since
// bittype#310, because a published member changed: the body of an entity. A
// receipt of the first version is no receipt of the second. The check below is
// one function for both languages and runs no compiler.

// profileResolver publishes an empty graph under the contract it is registered
// with. It stands for a resolver of another version of a profile.
type profileResolver[L any] struct {
	language.Provider
	reports []semantics.Key
}

func (p profileResolver[L]) Resolve(_ context.Context, request semantic.Request[L]) (semantic.Result, error) {
	publication := semantic.Publication{Facts: &semantics.Data{}}
	for _, key := range p.reports {
		publication.Reports = append(publication.Reports, semantic.Report{Obligation: key, Status: semantic.Complete, Verdict: semantic.Holds})
	}
	return request.Publish(publication)
}

// profileAnalyzer publishes its resolution's graph as an analysis.
type profileAnalyzer[L any] struct {
	language.Provider
	reports []semantics.Key
}

func (p profileAnalyzer[L]) Analyze(_ context.Context, request analysis.Request[L]) (semantic.Result, error) {
	resolved, _ := request.Resolution()
	publication := semantic.Publication{ReuseSemanticsFrom: &resolved}
	for _, key := range p.reports {
		publication.Reports = append(publication.Reports, semantic.Report{Obligation: key, Status: semantic.Complete, Verdict: semantic.Holds})
	}
	return request.Publish(publication)
}

// earlierProfile is what the check needs to know of one language.
type earlierProfile[L any] struct {
	snapshot                           semantic.Snapshot
	settings                           []semantic.Setting
	reports                            []semantics.Key
	resolution                         resolution.Contract
	analysis                           analysis.Contract
	earlierResolution, earlierAnalysis string
	// The language's own providers and reader.
	resolver resolution.Provider[L]
	analyzer analysis.Provider[L]
	open     func(semantic.Result) error
}

func checkEarlierProfileIsRefused[L any](t *testing.T, c earlierProfile[L]) {
	t.Helper()
	ctx := context.Background()
	if c.resolution.Semantic.Profile == c.earlierResolution || c.analysis.Semantic.Profile == c.earlierAnalysis {
		t.Fatal("the profile did not change")
	}
	resolve := func(profile string, provider resolution.Provider[L]) (semantic.Result, error) {
		t.Helper()
		contract := c.resolution
		contract.Semantic.Profile = profile
		resolver, err := resolution.NewResolver[L](c.snapshot.Language(), contract, provider)
		if err != nil {
			t.Fatal(err)
		}
		return resolver.Resolve(ctx, c.snapshot, semantic.Options{Configuration: c.settings})
	}
	analyze := func(profile string, provider analysis.Provider[L], resolved semantic.Result) (semantic.Result, error) {
		t.Helper()
		contract := c.analysis
		contract.Semantic.Profile = profile
		analyzer, err := analysis.NewAnalyzer[L](c.snapshot.Language(), contract, provider)
		if err != nil {
			t.Fatal(err)
		}
		return analyzer.Analyze(ctx, analysis.Input{Snapshot: c.snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: c.settings}})
	}
	publishing := profileResolver[L]{c.snapshot.Language(), c.reports}
	reusing := profileAnalyzer[L]{c.snapshot.Language(), c.reports}

	// The same graph and configuration under this version are admitted.
	current, err := resolve(c.resolution.Semantic.Profile, publishing)
	if err != nil {
		t.Fatal(err)
	}
	analyzed, err := analyze(c.analysis.Semantic.Profile, reusing, current)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.open(analyzed); err != nil {
		t.Fatal("a receipt of this profile", err)
	}
	// Analysis does not take a resolution receipt of the earlier profile.
	earlierResolved, err := resolve(c.earlierResolution, publishing)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := analyze(c.analysis.Semantic.Profile, reusing, earlierResolved); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Error("a resolution receipt of the earlier profile was admitted", err)
	}
	// The reader does not open an analysis receipt of the earlier profile.
	earlierAnalyzed, err := analyze(c.earlierAnalysis, reusing, current)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.open(earlierAnalyzed); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Error("an analysis receipt of the earlier profile was opened", err)
	}
	// The language's own providers serve no request under the earlier profile.
	if _, err := resolve(c.earlierResolution, c.resolver); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Error("the resolver served the earlier profile", err)
	}
	if _, err := analyze(c.earlierAnalysis, c.analyzer, current); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Error("the analyzer served the earlier profile", err)
	}
}

func TestGoSemanticEarlierProfileIsRefused(t *testing.T) {
	if golang.ResolutionProfile != "go1.25/captured-types-2/resolution-1" || golang.AnalysisProfile != "go1.25/captured-types-2/analysis-1" {
		t.Fatal("the profile strings", golang.ResolutionProfile, golang.AnalysisProfile)
	}
	snapshot, _ := goCapture(t, "type A struct{}")
	checkEarlierProfileIsRefused(t, earlierProfile[golangsyntax.Language]{
		snapshot: snapshot, settings: golang.Configuration(), reports: []semantics.Key{golang.Coverage(), golang.Bindings(), golang.Validity()},
		resolution: golang.ResolutionContract(), analysis: golang.AnalysisContract(),
		earlierResolution: "go1.25/captured-types-1/resolution-1", earlierAnalysis: "go1.25/captured-types-1/analysis-1",
		resolver: resolvegolang.Provider{}, analyzer: analyzegolang.Provider{},
		open: func(result semantic.Result) error { _, err := golang.Open(result); return err },
	})
}

func TestTypeScriptSemanticEarlierProfileIsRefused(t *testing.T) {
	if typescript.ResolutionProfile != "typescript-5.9.3/captured-types-2/resolution-1" || typescript.AnalysisProfile != "typescript-5.9.3/captured-types-2/analysis-1" || typescript.HelperProfile != "captured-types-2" {
		t.Fatal("the profile strings", typescript.ResolutionProfile, typescript.AnalysisProfile, typescript.HelperProfile)
	}
	snapshot, _, _ := tsCapture(t, "export interface A {}")
	compiler := tsCompiler(t)
	checkEarlierProfileIsRefused(t, earlierProfile[typescriptsyntax.Language]{
		snapshot: snapshot, settings: typescript.Configuration(), reports: []semantics.Key{typescript.Coverage(), typescript.Bindings(), typescript.Validity()},
		resolution: typescript.ResolutionContract(), analysis: typescript.AnalysisContract(),
		earlierResolution: "typescript-5.9.3/captured-types-1/resolution-1", earlierAnalysis: "typescript-5.9.3/captured-types-1/analysis-1",
		resolver: resolvetypescript.Provider{Compiler: compiler}, analyzer: analyzetypescript.Provider{},
		open: func(result semantic.Result) error { _, err := typescript.Open(result); return err },
	})
	// The configuration of the earlier helper profile is no configuration of
	// this one: the helper setting is part of every receipt.
	earlier := typescript.Configuration()
	found := false
	for i, setting := range earlier {
		if setting.Key == typescript.Key("helper") {
			earlier[i].Value, found = []byte("captured-types-1"), true
		}
	}
	if !found || typescript.MatchesConfiguration(earlier) {
		t.Fatal("the earlier helper setting is a configuration of this profile")
	}
	if _, err := typescript.EvidenceFor(earlier); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Error("evidence for the earlier helper setting", err)
	}
	resolver, err := resolvetypescript.New(snapshot.Language(), compiler)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{Configuration: earlier}); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Error("resolution with the earlier helper setting", err)
	}
}
