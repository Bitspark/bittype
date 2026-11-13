package integration_test

// PROTOTYPE (bittype#196). The test composition of the contribution contract:
// the consumer that registers bittype's own bundles (O11). This file is the
// one place of the package that lists the contributions the prototype moved.

import (
	"context"
	"slices"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	gobooleancodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/boolean/go"
	godeclarationscodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/declarations/go"
	golistcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/list/go"
	gooptionalcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/optional/go"
	gorecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/record/go"
	gotextcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/text/go"
	goboolean "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/boolean/go"
	godeclarations "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/declarations/go"
	golist "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/list/go"
	gooptional "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/optional/go"
	gorecord "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/record/go"
	gotext "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/text/go"
	tsdeclarationscodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/declarations/go"
	tslistcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/list/go"
	tsoptionalcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/optional/go"
	tsrecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/record/go"
	tstextcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/text/go"
	tsdeclarations "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/declarations/go"
	tslist "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/list/go"
	tsoptional "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/optional/go"
	tsrecord "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/record/go"
	tstext "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/text/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	tgtarget "github.com/Bitspark/bittype/transform/syntax/typescript/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type golangContribution = g.Contribution[h.Language, gs.Language, source.Reader]
type typescriptContribution = g.Contribution[h.Language, ts.Language, source.Reader]
type horosRule = g.Rule[h.Language, source.Reader]
type horosName = g.NameStep[h.Language, source.Reader]
type horosDecide = g.DecideStep[h.Language, source.Reader]
type horosRequest = g.Request[h.Language]

// bundle is what a composition registers for one target: contributions, the
// codecs of their payload schemas, and what a test adds for products.
type bundle[T any] struct {
	contributions []g.Contribution[h.Language, T, source.Reader]
	codecs        []g.PayloadCodec
	facets        []g.Facet
	products      []g.Product
}

// replace puts another contribution in the place of the one with the same
// domain: the composition edit of a consumer that selects it instead.
func (b bundle[T]) replace(with g.Contribution[h.Language, T, source.Reader]) bundle[T] {
	out := b
	out.contributions = nil
	for _, old := range b.contributions {
		if old.Domain.Key() != with.Domain.Key() {
			out.contributions = append(out.contributions, old)
		}
	}
	out.contributions = append(out.contributions, with)
	return out
}
func (b bundle[T]) with(more ...g.Contribution[h.Language, T, source.Reader]) bundle[T] {
	out := b
	out.contributions = append(append([]g.Contribution[h.Language, T, source.Reader](nil), b.contributions...), more...)
	return out
}
func (b bundle[T]) withCodecs(more ...g.PayloadCodec) bundle[T] {
	out := b
	out.codecs = append(append([]g.PayloadCodec(nil), b.codecs...), more...)
	return out
}
func (b bundle[T]) without(key string) bundle[T] {
	out := b
	out.contributions = nil
	for _, old := range b.contributions {
		if old.Reference.Identity.Key != key {
			out.contributions = append(out.contributions, old)
		}
	}
	return out
}
func (b bundle[T]) find(t *testing.T, key string) g.Contribution[h.Language, T, source.Reader] {
	t.Helper()
	for _, old := range b.contributions {
		if old.Reference.Identity.Key == key {
			return old
		}
	}
	t.Fatalf("no contribution %s in the bundle", key)
	panic("unreachable")
}

func protoMust[V any](t *testing.T) func(V, error) V {
	return func(v V, err error) V {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
}

// golangBundle is bittype's own bundle for Go as the prototype moved it: text,
// booleans, lists, closed records, the declaration contribution and the
// optional contribution.
func golangBundle(t *testing.T, from representation.Language) bundle[gs.Language] {
	t.Helper()
	contribution, codec := protoMust[golangContribution](t), protoMust[g.PayloadCodec](t)
	return bundle[gs.Language]{
		contributions: []golangContribution{
			contribution(gotext.Contribution(from)), contribution(goboolean.Contribution(from)), contribution(golist.Contribution(from)),
			contribution(gorecord.Contribution(from)), contribution(godeclarations.Contribution(from)), contribution(gooptional.Contribution(from)),
		},
		codecs: []g.PayloadCodec{
			codec(gotextcodec.Codec()), codec(gobooleancodec.Codec()), codec(golistcodec.Codec()),
			codec(gorecordcodec.Codec()), codec(godeclarationscodec.Codec()), codec(gooptionalcodec.Codec()),
		},
	}
}

// typescriptBundle is bittype's own bundle for TypeScript as the prototype
// moved it. It has no boolean contribution: that package is not moved.
func typescriptBundle(t *testing.T, from representation.Language) bundle[ts.Language] {
	t.Helper()
	contribution, codec := protoMust[typescriptContribution](t), protoMust[g.PayloadCodec](t)
	return bundle[ts.Language]{
		contributions: []typescriptContribution{
			contribution(tstext.Contribution(from)), contribution(tslist.Contribution(from)), contribution(tsrecord.Contribution(from)),
			contribution(tsdeclarations.Contribution(from)), contribution(tsoptional.Contribution(from)),
		},
		codecs: []g.PayloadCodec{
			codec(tstextcodec.Codec()), codec(tslistcodec.Codec()), codec(tsrecordcodec.Codec()),
			codec(tsdeclarationscodec.Codec()), codec(tsoptionalcodec.Codec()),
		},
	}
}

// protoSource parses horos text and analyses it into a generation input.
func protoSource(t *testing.T, text string) (representation.Language, transform.Input) {
	t.Helper()
	parsed, err := parsehoros.Parse(context.Background(), parse.Input{SourceID: "prototype.horos", TreeID: "canonical", Text: text})
	if err != nil {
		t.Fatal(err)
	}
	_, from, input := generationSource(t, parsed.Abstract)
	return from, input
}

type protoRequest struct {
	products      []string
	configuration []g.Setting
	roots         []string // names of root declarations; all of the entry module when empty
}

// protoPlanWith registers a bundle in a fresh registry and plans the request.
func protoPlanWith[T any](t *testing.T, input transform.Input, checker check.Checker, target g.Target[T], b bundle[T], q protoRequest) (g.ContributionPlan, error) {
	t.Helper()
	var registry c.Registry
	from, to := input.Source.Language(), checker.Language()
	if err := c.RegisterSource[h.Language, source.Reader](&registry, from, source.Provider{}); err != nil {
		t.Fatal(err)
	}
	if err := c.RegisterTarget(&registry, to, target); err != nil {
		t.Fatal(err)
	}
	for _, contribution := range b.contributions {
		if err := c.RegisterContribution(&registry, from, to, contribution); err != nil {
			return g.ContributionPlan{}, err
		}
	}
	for _, codec := range b.codecs {
		if err := registry.RegisterCodec(codec); err != nil {
			return g.ContributionPlan{}, err
		}
	}
	for _, facet := range b.facets {
		if err := registry.RegisterFacet(facet); err != nil {
			return g.ContributionPlan{}, err
		}
	}
	for _, product := range b.products {
		if err := registry.RegisterProduct(product); err != nil {
			return g.ContributionPlan{}, err
		}
	}
	roots, err := genfixture.Roots(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(q.roots) > 0 {
		roots = protoRoots(t, input, q.roots)
	}
	request := c.ContributionRequest{Request: c.Request{Input: input, To: to, Output: "generated", Roots: roots, Profile: rules.RepresentationProfile, Limits: g.DefaultLimits()}, Products: q.products, Configuration: q.configuration}
	return registry.PlanContributions(context.Background(), request, checker)
}

func protoPlanGo(t *testing.T, input transform.Input, b bundle[gs.Language], q protoRequest) (g.ContributionPlan, error) {
	t.Helper()
	return protoPlanWith[gs.Language](t, input, nativeChecker(t, golangText()), ggtarget.Provider{Package: "generated"}, b, q)
}
func protoPlanTS(t *testing.T, input transform.Input, b bundle[ts.Language], q protoRequest) (g.ContributionPlan, error) {
	t.Helper()
	return protoPlanWith[ts.Language](t, input, nativeChecker(t, typescriptText()), tgtarget.Provider{}, b, q)
}

// protoText runs a plan and prints the generated native text.
func protoText[L any](t *testing.T, p nativeProvider[L], plan g.ContributionPlan) string {
	t.Helper()
	out, err := plan.Generation.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return formatted(t, p, outputTree[L](t, out))
}

// protoRoots finds the declaring sites of the named type definitions.
func protoRoots(t *testing.T, input transform.Input, names []string) []ast.NodeRef {
	t.Helper()
	var out []ast.NodeRef
	for _, envelope := range input.Source.Inputs() {
		tree, err := representation.Unpack[h.Language](input.Source.Language(), envelope)
		if err != nil {
			t.Fatal(err)
		}
		for i, node := range tree.Nodes {
			if node.Kind == h.KindTypeDefinition && slices.Contains(names, nodeAttribute(node, h.AttributeName)) {
				out = append(out, ast.NodeRef{Tree: tree.ID, Node: ast.NodeID(i + 1)})
			}
		}
	}
	if len(out) != len(names) {
		t.Fatalf("found %d of the roots %v", len(out), names)
	}
	return out
}
