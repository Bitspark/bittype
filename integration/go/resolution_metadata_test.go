package integration_test

import (
	"context"
	"errors"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	"testing"
)

func TestSemanticInvocationMetadataUsesAggregateBudgets(t *testing.T) {
	_, d, snapshot, resolver := bindingFixture(t, nil)
	settings := []semantic.Setting{{Key: analysisKey("config"), Value: make([]byte, 512)}}
	result, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{Configuration: settings})
	if err != nil {
		t.Fatal(err)
	}
	requirement := bindingRequirement()
	requirement.Configuration = settings
	dependency := requirement.Dependency(result)
	called := false
	op, err := semantic.NewOperation[amber](d, bindingContract().Semantic, analysis.Kind(), bindingProvider[amber]{id: d.ID()}, func(_ context.Context, request semantic.Request[amber]) (semantic.Result, error) {
		called = true
		return request.Publish(semantic.Publication{Reports: bindingReports()})
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"dependency bytes", "setting entries", "report entries", "own settings and requested keys"} {
		t.Run(name, func(t *testing.T) {
			limits := semantic.DefaultLimits()
			options := semantic.Options{}
			dependencies := []semantic.Dependency{dependency, dependency}
			switch name {
			case "dependency bytes":
				limits.MaxConfigurationBytes = 1024
			case "setting entries":
				limits.MaxConfigurationEntries = 1
			case "report entries":
				limits.MaxRequested = 2
			case "own settings and requested keys":
				limits.MaxConfigurationBytes = 1024
				dependencies = nil
				options.Configuration = settings
				options.Requested = []semantics.Key{{Owner: "test", Name: string(make([]byte, 600))}}
			}
			options.Limits = limits
			called = false
			out, err := op.Run(context.Background(), snapshot, options, dependencies)
			if !errors.Is(err, semantic.ErrLimit) || out.Valid() || called {
				t.Fatalf("aggregate budget did not reject before callback: %v, called=%v", err, called)
			}
		})
	}
	if _, err := op.Run(context.Background(), snapshot, semantic.Options{}, []semantic.Dependency{dependency, dependency}); err != nil {
		t.Fatal("bounded multiple dependencies rejected", err)
	}
}

func TestReusedModelChecksReceivingDepthEdgesAndBytes(t *testing.T) {
	_, d, s, r := bindingFixture(t, func(_ context.Context, request semantic.Request[amber]) (semantic.Result, error) {
		facts := bindingFacts()
		facts.Entities[0].Body = 1
		facts.Terms = []semantics.Term{
			{Kind: semantics.ProviderExpression, Operator: analysisKey("list"), Encoding: analysisKey("bytes"), Data: make([]byte, 512), Operands: []semantics.Operand{{Locator: []uint32{0}, Term: 2}}},
			{Kind: semantics.Reference, Target: 1},
		}
		return request.Publish(semantic.Publication{Facts: &facts, Reports: bindingReports()})
	})
	resolved, err := r.Resolve(context.Background(), s, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	a := bindingAnalysis(t, d, bindingRequirement(), nil)
	for _, name := range []string{"depth", "edges", "bytes"} {
		t.Run(name, func(t *testing.T) {
			limits := semantic.DefaultLimits()
			switch name {
			case "depth":
				limits.Semantics.MaxDepth = 1
			case "edges":
				limits.Semantics.MaxEdges = 1
			case "bytes":
				limits.Semantics.MaxBytes = 256
			}
			out, err := a.Analyze(context.Background(), analysis.Input{Snapshot: s, Resolution: &resolved, Options: semantic.Options{Limits: limits}})
			if !errors.Is(err, semantic.ErrLimit) || out.Valid() {
				t.Fatalf("reused model escaped %s limit: %v", name, err)
			}
		})
	}
}

func TestAnalyzerInvocationIdentityWithinOneSnapshot(t *testing.T) {
	_, d, s, r := bindingFixture(t, nil)
	resolved, err := r.Resolve(context.Background(), s, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	var retained semantic.Result
	a := bindingAnalysis(t, d, bindingRequirement(), func(_ context.Context, request analysis.Request[amber]) (semantic.Result, error) {
		if retained.Valid() {
			return retained, nil
		}
		dependency, _ := request.Resolution()
		var err error
		retained, err = request.Publish(semantic.Publication{ReuseSemanticsFrom: &dependency, Reports: bindingReports()})
		return retained, err
	})
	input := analysis.Input{Snapshot: s, Resolution: &resolved}
	if _, err := a.Analyze(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	if out, err := a.Analyze(context.Background(), input); !errors.Is(err, semantic.ErrStaleResult) || out.Valid() {
		t.Fatal("analysis returned a prior invocation on the same snapshot", err)
	}
}

func TestReusedModelPreservesCalleesAndSubjectFactBounds(t *testing.T) {
	_, definition, snapshot, resolver := bindingFixture(t, func(_ context.Context, request semantic.Request[amber]) (semantic.Result, error) {
		facts := semantics.Data{
			Entities: []semantics.Entity{{Kind: semantics.Key{Owner: "k", Name: "d"}}},
			Terms: []semantics.Term{
				{Kind: semantics.Reference, Target: 1},
				{Kind: semantics.Application, Callee: 1, Arguments: []semantics.TermID{1}},
			},
			Facts: []semantics.Fact{{Subject: semantics.Subject{Term: 2}, Relation: semantics.Key{Owner: "r", Name: "f"}, Term: 1}},
		}
		return request.Publish(semantic.Publication{Facts: &facts, Reports: bindingReports()})
	})
	resolved, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	analyzer := bindingAnalysis(t, definition, bindingRequirement(), nil)
	limits := semantic.DefaultLimits()
	limits.Semantics = semantics.Limits{MaxRecords: 4, MaxEdges: 4, MaxBytes: 4, MaxDepth: 2}
	result, err := analyzer.Analyze(context.Background(), analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Limits: limits}})
	if err != nil || !result.HasSemantics() {
		t.Fatalf("exact callee/fact bounds rejected reuse: %v", err)
	}
	before, _ := resolved.Semantics()
	after, _ := result.Semantics()
	if before != after {
		t.Fatal("callee/fact reuse changed sealed model identity")
	}
	term, _ := after.Term(2)
	values := after.Facts(semantics.Subject{Term: 2}, semantics.Key{Owner: "r", Name: "f"})
	if term.Kind != semantics.Application || term.Callee != 1 || len(values) != 1 || values[0] != 1 {
		t.Fatal("reuse lost callee or fact identity")
	}
	values[0] = 99
	if got := before.Facts(semantics.Subject{Term: 2}, semantics.Key{Owner: "r", Name: "f"}); len(got) != 1 || got[0] != 1 {
		t.Fatal("reused fact query aliases receipt storage")
	}
	for name, reduce := range map[string]func(*semantics.Limits){
		"fact record":           func(l *semantics.Limits) { l.MaxRecords-- },
		"fact and callee edges": func(l *semantics.Limits) { l.MaxEdges-- },
		"relation bytes":        func(l *semantics.Limits) { l.MaxBytes-- },
		"callee depth":          func(l *semantics.Limits) { l.MaxDepth-- },
	} {
		t.Run(name, func(t *testing.T) {
			tighter := limits
			reduce(&tighter.Semantics)
			out, err := analyzer.Analyze(context.Background(), analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Limits: tighter}})
			if !errors.Is(err, semantic.ErrLimit) || out.Valid() || out.HasSemantics() {
				t.Fatalf("receiving budget ignored %s: %v", name, err)
			}
		})
	}
}
