package integration_test

import (
	"context"
	"errors"
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	"math/big"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/horos/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/engine/go"
	provider "github.com/Bitspark/bittype/provider/horos/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	"github.com/Bitspark/bittype/semantics/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// A test-only capability catalogue qualifies every semantic provider explicitly.
// It is not imported by any production registry or language package.
func semanticConformanceCases() []struct {
	directory string
	qualify   func(*testing.T)
} {
	return []struct {
		directory string
		qualify   func(*testing.T)
	}{{"horos", qualifyHoros}}
}

func TestSemanticProviders(t *testing.T) {
	for _, fixture := range semanticConformanceCases() {
		t.Run(fixture.directory, fixture.qualify)
	}
}

func qualifyHoros(t *testing.T) {
	provider := provider.Provider{}
	if provider.ID() != "horos" || provider.Name() != "Horos" {
		t.Fatal("unstable Horos identity")
	}
	if _, ok := any(provider).(parse.Provider[h.Language]); !ok {
		t.Fatal("Horos provider lacks parsing")
	}
	if _, ok := any(provider).(format.Provider[h.Language]); !ok {
		t.Fatal("Horos provider lacks formatting")
	}
	// Independent consumer: a generic list and its explicit text specialization.
	atom := func(s string) core.Value { return core.NewAtom([]byte(s)) }
	tuple := func(tag string, args ...core.Value) core.Value {
		return core.NewTuple(append([]core.Value{atom(tag)}, args...)...)
	}
	txt := data.MustEncodeText
	textType := tuple("horos/scalar", atom("utf8-text"))
	parameter := tuple("horos/decl/param", data.EncodeInt(big.NewInt(0)))
	page := tuple("horos/decl/type", data.EncodeList([]core.Value{tuple("horos/decl/parameter", txt("T"), tuple("horos/decl/proper", atom("data")))}), atom("type"), tuple("horos/list-of", parameter))
	applied := tuple("horos/decl/type", data.EncodeList(nil), atom("type"), tuple("horos/decl/apply", tuple("horos/decl/local", txt("Page")), data.EncodeList([]core.Value{textType})))
	declarations, err := data.EncodeMap([][2]core.Value{{txt("Page"), page}, {txt("TextPage"), applied}})
	if err != nil {
		t.Fatal(err)
	}
	empty, _ := data.EncodeMap(nil)
	module := tuple("horos/decl/module", data.EncodeInt(big.NewInt(3)), txt("consumer.pagination"), empty, declarations)
	tree, err := resolve.Resolve(module, nil, "consumer")
	if err != nil {
		t.Fatal(err)
	}
	definition, err := representation.NewLanguage[h.Language](provider, h.Profile)
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := representation.Pack[h.Language](definition, tree)
	if err != nil {
		t.Fatal(err)
	}
	tree.Nodes[tree.Root-1].Kind = "broken" // caller storage is not the envelope
	var registry engine.Engine
	if err := registry.RegisterLanguage(definition); err != nil {
		t.Fatal(err)
	}
	checker, err := checkcontract.NewChecker[h.Language](definition, provider)
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterChecker(checker); err != nil {
		t.Fatal(err)
	}
	snapshot, err := registry.Capture(context.Background(), semantic.Input{Entry: ast.NodeRef{Tree: "consumer", Node: tree.Root}, Trees: []representation.AST{envelope}})
	if err != nil {
		t.Fatal(err)
	}
	input := analysis.Input{Snapshot: snapshot}
	if _, err := registry.Analyze(context.Background(), input); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatalf("syntax registration silently installs an analyzer: %v", err)
	}
	analyzer, err := analysis.NewAnalyzer[h.Language](definition, analyze.Contract(), provider)
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterAnalyzer(analyzer); err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Resolve(context.Background(), snapshot, semantic.Options{}); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal(err)
	}
	resolver, err := resolution.NewResolver[h.Language](definition, resolve.Contract(), provider)
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterResolver(resolver); err != nil {
		t.Fatal(err)
	}
	resolved, err := registry.Resolve(context.Background(), snapshot, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	input.Resolution = &resolved
	result, err := registry.Analyze(context.Background(), input)
	if err != nil || !result.HasModel() {
		t.Fatalf("explicitly registered Horos analysis: %v", err)
	}
	for _, obligation := range analyze.Contract().Semantic.Obligations {
		report, found := result.Report(obligation.Key)
		want := semantic.Undetermined
		if obligation.Key == resolve.Validity() {
			want = semantic.Holds
		}
		if !found || report.Status != semantic.Complete || report.Verdict != want {
			t.Fatalf("Horos obligation not established: %+v", report)
		}
	}
	qualifyHorosFacts(t, result)
	resolvedModel, _ := resolved.Model()
	analyzedModel, _ := result.Model()
	if resolvedModel != analyzedModel {
		t.Fatal("analysis rebuilt the resolved model instead of reusing it")
	}
	output, err := registry.TransformTo(context.Background(), transform.Input{Source: input.Snapshot}, provider.ID())
	if err != nil {
		t.Fatal(err)
	}
	checked, err := representation.Unpack[h.Language](definition, output)
	if err != nil {
		t.Fatal(err)
	}
	model, err := resolve.Model(checked)
	if err != nil {
		t.Fatal(err)
	}
	specialized, err := model.Specialize(decl.DeclarationID{Module: model.Root(), Name: "TextPage"}, nil)
	if err != nil || !core.Equal(specialized, tuple("horos/list-of", textType)) {
		t.Fatalf("specialization through registered AST: %v", err)
	}
	checked.Nodes[checked.Root-1].Kind = "mutated"
	if fresh, err := representation.Unpack[h.Language](definition, output); err != nil || fresh.Nodes[fresh.Root-1].Kind != h.KindProgram {
		t.Fatal("unpacked AST aliases envelope")
	}
	if _, err := registry.Parse(provider.ID(), parse.Input{}); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatalf("unregistered parse: %v", err)
	}
	if _, err := registry.Render(envelope); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatalf("unregistered render: %v", err)
	}
}

// This consumer uses shared semantic types and queries only. Horos-specific
// decoding, specialization and AST node vocabulary are unnecessary here.
func qualifyHorosFacts(t *testing.T, result semantic.Result) {
	t.Helper()
	model, ok := result.Model()
	if !ok {
		t.Fatal("no semantic model")
	}
	facts := model.Data()
	var page, applied semantics.EntityID
	for i, entity := range facts.Entities {
		switch entity.Name {
		case "Page":
			page = semantics.EntityID(i + 1)
		case "TextPage":
			applied = semantics.EntityID(i + 1)
		}
	}
	pageEntity, hasPage := model.Entity(page)
	appliedEntity, hasApplied := model.Entity(applied)
	if !hasPage || !hasApplied || len(pageEntity.Parameters) != 1 {
		t.Fatal("declarations or parameter identity lost")
	}
	term, exists := model.Term(appliedEntity.Body)
	bindings, bound := model.DirectParameterBindings(appliedEntity.Body)
	callee, hasCallee := model.Term(term.Callee)
	if !exists || term.Kind != semantics.Application || !hasCallee || callee.Kind != semantics.Reference || callee.Target != page || !bound || len(bindings) != 1 || bindings[0].Parameter != pageEntity.Parameters[0].Entity || bindings[0].Argument != term.Arguments[0] {
		t.Fatal("generic parameter bindings do not preserve the declared application")
	}
	if len(facts.Bindings) == 0 || len(facts.Expressions) == 0 {
		t.Fatal("missing occurrence facts")
	}
	for _, binding := range facts.Bindings {
		if target, found := model.Resolve(binding.At); !found || target != binding.Target {
			t.Fatal("occurrence target query disagrees with exported binding")
		}
	}
	for _, expression := range facts.Expressions {
		if meaning, found := model.TermAt(expression.At); !found || meaning != expression.Term {
			t.Fatal("expression query disagrees with exported term")
		}
	}
	foundDependency := false
	for _, dependency := range facts.Dependencies {
		if dependency.From == applied && dependency.To == page {
			for _, queried := range model.DependenciesFrom(applied, dependency.Relation) {
				foundDependency = foundDependency || queried.To == page
			}
		}
	}
	if !foundDependency {
		t.Fatal("application dependency was not exported")
	}
}
