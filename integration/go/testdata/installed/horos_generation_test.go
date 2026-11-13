package consumer_test

import (
	"context"
	"strings"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/files/horos/provider/go"
	checkhoros "github.com/Bitspark/bittype/check/syntax/horos/provider/go"
	checktypescript "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/engine/go"
	formattypescript "github.com/Bitspark/bittype/format/syntax/typescript/provider/go"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	syntaxprint "github.com/Bitspark/bittype/print/syntax/shared/provider/go"
	resolve "github.com/Bitspark/bittype/resolve/files/horos/provider/go"
	"github.com/Bitspark/bittype/semantics/horos/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	tsbooleancodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/boolean/go"
	tsdeclarationscodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/declarations/go"
	tslistcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/list/go"
	tsoptionalcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/optional/go"
	tsrecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/record/go"
	tstextcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/text/go"
	tsboolean "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/boolean/go"
	tsdeclarations "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/declarations/go"
	tslist "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/list/go"
	tsoptional "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/optional/go"
	tsrecord "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/record/go"
	tstext "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/text/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	reader "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	target "github.com/Bitspark/bittype/transform/syntax/typescript/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func TestInstalledSemanticGeneration(t *testing.T) {
	ctx := context.Background()
	source, err := parsehoros.Parse(context.Background(), parse.Input{SourceID: "input.horos", TreeID: "input", Text: "module installed; type Page<T> = { items: list<T>; next?: text; }; type People = Page<text>;"})
	if err != nil {
		t.Fatal(err)
	}
	from, err := representation.NewLanguage[h.Language](h.Language{}, h.Profile)
	if err != nil {
		t.Fatal(err)
	}
	to, err := representation.NewLanguage[ts.Language](ts.Language{}, "native-types-v1")
	if err != nil {
		t.Fatal(err)
	}
	e := engine.New()
	for _, d := range []representation.Language{from, to} {
		if err = e.RegisterLanguage(d); err != nil {
			t.Fatal(err)
		}
	}
	hc, err := check.NewChecker[h.Language](from, checkhoros.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	tc, err := check.NewChecker[ts.Language](to, checktypescript.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []check.Checker{hc, tc} {
		if err = e.RegisterChecker(c); err != nil {
			t.Fatal(err)
		}
	}
	resolver, err := resolution.NewResolver[h.Language](from, horos.ResolutionContract(), resolve.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterResolver(resolver); err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[h.Language](from, horos.AnalysisContract(), analyze.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterAnalyzer(analyzer); err != nil {
		t.Fatal(err)
	}
	packed, err := representation.Pack(from, source)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := e.Capture(ctx, semantic.Input{Entry: ast.NodeRef{Tree: source.Abstract.ID, Node: source.Abstract.Root}, Trees: []representation.AST{packed}})
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := e.Resolve(ctx, snapshot, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	analyzed, err := e.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved})
	if err != nil {
		t.Fatal(err)
	}
	registry := e.Generation()

	if err = c.RegisterSource[h.Language, reader.Reader](registry, from, reader.Provider{}); err != nil {
		t.Fatal(err)
	}
	if err = c.RegisterTarget(registry, to, target.Provider{}); err != nil {
		t.Fatal(err)
	}
	// The consumer's composition: each builtin contribution with the codec of
	// its payload schema. The record contribution delegates the optional
	// field to the optional contribution through the planner.
	for _, bundle := range []struct {
		contribution func(representation.Language) (g.Contribution[h.Language, ts.Language, reader.Reader], error)
		codec        func() (g.PayloadCodec, error)
	}{
		{tstext.Contribution, tstextcodec.Codec},
		{tsboolean.Contribution, tsbooleancodec.Codec},
		{tslist.Contribution, tslistcodec.Codec},
		{tsrecord.Contribution, tsrecordcodec.Codec},
		{tsdeclarations.Contribution, tsdeclarationscodec.Codec},
		{tsoptional.Contribution, tsoptionalcodec.Codec},
	} {
		contribution, err := bundle.contribution(from)
		if err != nil {
			t.Fatal(err)
		}
		if err := c.RegisterContribution(registry, from, to, contribution); err != nil {
			t.Fatal(err)
		}
		codec, err := bundle.codec()
		if err != nil {
			t.Fatal(err)
		}
		if err := registry.RegisterCodec(codec); err != nil {
			t.Fatal(err)
		}
	}
	var roots []ast.NodeRef
	for i, node := range source.Abstract.Nodes {
		if node.Kind == h.KindTypeDefinition {
			roots = append(roots, ast.NodeRef{Tree: source.Abstract.ID, Node: ast.NodeID(i + 1)})
		}
	}
	request := c.Request{Input: transform.Input{Source: snapshot, Evidence: []semantic.Result{analyzed}}, To: to, Roots: roots, Output: "generated-types", Limits: g.DefaultLimits(), Profile: rules.RepresentationProfile}
	plan, err := e.PlanGeneration(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Steps()) == 0 {
		t.Fatal("empty plan")
	}
	output, err := plan.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
	typed, err := representation.Unpack[ts.Language](to, output.Snapshot().Inputs()[0])
	if err != nil {
		t.Fatal(err)
	}
	concrete, err := formattypescript.Format(context.Background(), typed)
	if err != nil {
		t.Fatal(err)
	}
	text, err := syntaxprint.Print[ts.Language](context.Background(), concrete)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "interface Page<T>") || !strings.Contains(text, "next?: string;") || !strings.Contains(text, "type People = Page<string>;") {
		t.Fatal(text)
	}
}
