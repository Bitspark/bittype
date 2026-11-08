package integration_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	checkts "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/engine/go"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	mappings "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	target "github.com/Bitspark/bittype/transform/syntax/typescript/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func kernelRequest(t *testing.T) (*engine.Engine, c.Request) {
	t.Helper()
	e, _, input := generationInput(t, paginationModule())
	to := mustLanguage[ts.Language](t, ts.Language{})
	if err := e.RegisterLanguage(to); err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[ts.Language](to, checkts.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterChecker(checker); err != nil {
		t.Fatal(err)
	}
	request, err := genfixture.Request(e.Generation(), input, to, "test-output", target.Provider[h.Language]{}, mappings.Providers(rules.RepresentationProfile)...)
	if err != nil {
		t.Fatal(err)
	}
	return e, request
}

type countingSource struct{ opened *atomic.Int32 }

func (p countingSource) Contract() transform.Contract { return source.Contract() }
func (p countingSource) Open(ctx context.Context, r transform.Request[h.Language]) (source.Reader, error) {
	p.opened.Add(1)
	return (source.Provider{}).Open(ctx, r)
}

type otherReader interface {
	source.Reader
	AdditionalAPI()
}
type otherMapping struct{}

func (otherMapping) Open(context.Context, g.MappingRequest[otherReader, target.Builder[h.Language]]) (g.Contribution[h.Language, ts.Language], error) {
	panic("incompatible API was invoked")
}

type otherTarget interface {
	target.Builder[h.Language]
	AdditionalAPI()
}
type wrongTargetMapping struct{}

func (wrongTargetMapping) Open(context.Context, g.MappingRequest[source.Reader, otherTarget]) (g.Contribution[h.Language, ts.Language], error) {
	panic("incompatible native API was invoked")
}
func TestGenerationPreflightAndEvidence(t *testing.T) {
	e, r := kernelRequest(t)
	from := r.Input.Source.Language()
	var opened atomic.Int32
	ref, err := c.RegisterSource[h.Language, source.Reader](e.Generation(), from, countingSource{&opened})
	if err != nil {
		t.Fatal(err)
	}
	r.Selection.Source = ref
	foreign, _ := c.RegisterSource[h.Language, source.Reader](&c.Registry{}, from, source.Provider{})
	wrong, _ := c.RegisterMapping[h.Language, ts.Language, otherReader, target.Builder[h.Language]](e.Generation(), from, r.To, otherMapping{})
	incompatible, _ := c.RegisterMapping[h.Language, ts.Language, source.Reader, otherTarget](e.Generation(), from, r.To, wrongTargetMapping{})
	wrongProfile, _ := representation.NewLanguage[ts.Language](ts.Language{}, "other-profile")
	_, _, foreignInput := generationInput(t, paginationModule())
	for name, change := range map[string]func(*c.Request){
		"empty roots":          func(r *c.Request) { r.Roots = nil },
		"duplicate roots":      func(r *c.Request) { r.Roots = []ast.NodeRef{r.Roots[0], r.Roots[0]} },
		"foreign site":         func(r *c.Request) { r.Roots = []ast.NodeRef{{Tree: "unrelated", Node: 1}} },
		"output collision":     func(r *c.Request) { r.Output = r.Roots[0].Tree },
		"foreign registration": func(r *c.Request) { r.Selection.Source = foreign },
		"zero target":          func(r *c.Request) { r.Selection.Target = c.TargetRef{} },
		"source API":           func(r *c.Request) { r.Selection.Mappings = []c.MappingRef{wrong} },
		"target API":           func(r *c.Request) { r.Selection.Mappings = []c.MappingRef{incompatible} },
		"wrong profile":        func(r *c.Request) { r.To = wrongProfile },
		"mapping selected twice": func(r *c.Request) {
			r.Selection.Mappings = []c.MappingRef{r.Selection.Mappings[0], r.Selection.Mappings[0]}
		},
		"missing receipt":     func(r *c.Request) { r.Input.Evidence = nil },
		"foreign receipt":     func(r *c.Request) { r.Input.Evidence = foreignInput.Evidence },
		"registration budget": func(r *c.Request) { r.Limits.Plans = 1 },
	} {
		t.Run(name, func(t *testing.T) {
			bad := r
			change(&bad)
			out, err := e.Generate(context.Background(), bad)
			if err == nil || out.Snapshot().Valid() || opened.Load() != 0 {
				t.Fatalf("output %v, error %v, source openings %d", out.Snapshot().Entry(), err, opened.Load())
			}
		})
	}
	q := r
	q.Roots = []ast.NodeRef{r.Input.Source.Entry()}
	if out, err := e.Generate(context.Background(), q); !errors.Is(err, g.ErrScope) || out.Snapshot().Valid() || opened.Load() != 1 {
		t.Fatalf("non-declaration root: %v", err)
	}
	if _, err := c.RegisterTarget[h.Language, ts.Language, target.Builder[h.Language]](e.Generation(), r.To, from, target.Provider[h.Language]{}); err == nil {
		t.Fatal("accepted wrong marker")
	}
}

type observingTarget struct{ saved g.TargetRequest }

func (p *observingTarget) Open(ctx context.Context, r g.TargetRequest) (target.Builder[h.Language], error) {
	p.saved = r
	return (target.Provider[h.Language]{}).Open(ctx, r)
}

type observeMapping func(context.Context, g.MappingRequest[source.Reader, target.Builder[h.Language]]) (g.Contribution[h.Language, ts.Language], error)

func (p observeMapping) Open(ctx context.Context, r g.MappingRequest[source.Reader, target.Builder[h.Language]]) (g.Contribution[h.Language, ts.Language], error) {
	return p(ctx, r)
}
func TestGenerationConstructionScopeAndFailure(t *testing.T) {
	e, r := kernelRequest(t)
	from := r.Input.Source.Language()
	var observed g.MappingRequest[source.Reader, target.Builder[h.Language]]
	native := &observingTarget{}
	targetRef, err := c.RegisterTarget[h.Language, ts.Language, target.Builder[h.Language]](e.Generation(), from, r.To, native)
	if err != nil {
		t.Fatal(err)
	}
	r.Selection.Target = targetRef
	var failure error
	extra := observeMapping(func(ctx context.Context, q g.MappingRequest[source.Reader, target.Builder[h.Language]]) (g.Contribution[h.Language, ts.Language], error) {
		src, n, err := q.Opened()
		if err != nil || src == nil || n == nil || !src.Snapshot().Same(r.Input.Source) {
			t.Fatal("opened instances", err)
		}
		info, err := q.Info()
		if err != nil || info.Output != r.Output || info.Limits != r.Limits {
			t.Fatal("invocation metadata", err)
		}
		observed = q
		if _, err = c.RegisterSource[h.Language, source.Reader](e.Generation(), from, source.Provider{}); err != nil {
			t.Fatal(err)
		}
		return g.Contribution[h.Language, ts.Language]{}, failure
	})
	m, err := c.RegisterMapping[h.Language, ts.Language, source.Reader, target.Builder[h.Language]](e.Generation(), from, r.To, extra)
	if err != nil {
		t.Fatal(err)
	}
	r.Selection.Mappings = append(r.Selection.Mappings, m)
	if _, err = e.Generate(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	for _, req := range []g.TargetRequest{{}, native.saved} {
		if _, err := req.Info(); !errors.Is(err, g.ErrScope) {
			t.Fatal("retained target metadata", err)
		}
	}
	if _, _, err := observed.Opened(); !errors.Is(err, g.ErrScope) {
		t.Fatal("retained mapping request", err)
	}
	failure = errors.New("provider failed")
	if out, err := e.Generate(context.Background(), r); !errors.Is(err, failure) || out.Snapshot().Valid() {
		t.Fatal("failed construction published", err)
	}
	if _, err := observed.Info(); !errors.Is(err, g.ErrScope) {
		t.Fatal("failed ticket still active", err)
	}
	for _, bad := range []g.Contribution[h.Language, ts.Language]{
		{Helpers: []g.Plan[h.Language, ts.Language]{{Key: "shared"}, {Key: "shared"}}},
		{Structural: &g.Mapping[h.Language, ts.Language]{ID: "other", PluginID: "independent", Plan: func(g.Declaration[h.Language]) (g.Plan[h.Language, ts.Language], error) {
			return g.Plan[h.Language, ts.Language]{}, nil
		}}},
	} {
		q := r
		ref, err := c.RegisterMapping[h.Language, ts.Language, source.Reader, target.Builder[h.Language]](e.Generation(), from, r.To, tsContribution(func(source.Reader, target.Builder[h.Language]) (g.Contribution[h.Language, ts.Language], error) {
			return bad, nil
		}))
		if err != nil {
			t.Fatal(err)
		}
		q.Selection.Mappings = append([]c.MappingRef{ref}, r.Selection.Mappings[:len(r.Selection.Mappings)-1]...)
		if out, err := e.Generate(context.Background(), q); !errors.Is(err, g.ErrInvalid) || out.Snapshot().Valid() {
			t.Fatal("conflicting contribution", err)
		}
	}
}
func TestEngineConcurrentGeneration(t *testing.T) {
	e, r := kernelRequest(t)
	r.Roots = r.Roots[:1]
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, err := e.Generate(context.Background(), r)
			if err != nil || out.Snapshot().Entry().Tree != r.Output {
				t.Errorf("concurrent generation: %v %v", out.Snapshot().Entry(), err)
			}
		}()
	}
	wg.Wait()
}

// Two independent constructions can use the same descriptors and exact receipts.
// Their opened handles still belong to different invocations.
func TestGenerationRejectsCrossInvocationHandles(t *testing.T) {
	_, r := kernelRequest(t)
	called := false
	transformer, err := transform.NewTransformer[h.Language, ts.Language](r.Input.Source.Language(), r.To, source.Contract(), func(ctx context.Context, tr transform.Request[h.Language]) (ast.Set[ts.Language], error) {
		info := g.TargetInfo{From: r.Input.Source.Language(), To: r.To, Output: r.Output, Limits: r.Limits}
		first, err := g.NewConstruction(ctx, info)
		if err != nil {
			t.Fatal(err)
		}
		defer first.Close()
		second, err := g.NewConstruction(ctx, info)
		if err != nil {
			t.Fatal(err)
		}
		defer second.Close()
		firstSource, err := g.OpenSource(first, source.Provider{}, tr)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = g.OpenSource(second, source.Provider{}, tr); err != nil {
			t.Fatal(err)
		}
		if _, err = g.OpenTarget(first, target.Provider[h.Language]{}); err != nil {
			t.Fatal(err)
		}
		secondTarget, err := g.OpenTarget(second, target.Provider[h.Language]{})
		if err != nil {
			t.Fatal(err)
		}
		_, err = g.OpenMapping(first, mappings.List(rules.RepresentationProfile), firstSource, secondTarget)
		if !errors.Is(err, g.ErrScope) {
			t.Fatal("constructed mapping with foreign native handle", err)
		}
		if _, err = firstSource.Value(second); !errors.Is(err, g.ErrScope) {
			t.Fatal("foreign construction recovered source", err)
		}
		if _, err = firstSource.Value(first); !errors.Is(err, g.ErrScope) {
			t.Fatal("failed construction recovered source", err)
		}
		called = true
		return ast.Set[ts.Language]{}, g.ErrScope
	})
	if err != nil {
		t.Fatal(err)
	}
	if out, err := transformer.Transform(context.Background(), r.Input); !called || !errors.Is(err, g.ErrScope) || len(out.Trees) != 0 {
		t.Fatal("fixture", err)
	}
}
func TestGenerationPublishesOnlyRequestedPage(t *testing.T) {
	e, r := kernelRequest(t)
	sourceTree, _ := r.Input.Source.EntryTree()
	original, err := representation.Unpack[h.Language](r.Input.Source.Language(), sourceTree)
	if err != nil {
		t.Fatal(err)
	}
	var page ast.NodeRef
	for _, ref := range r.Roots {
		if nodeAttribute(original.Nodes[ref.Node-1], h.AttributeName) == "Page" {
			page = ref
		}
	}
	if page.Tree == "" {
		t.Fatal("source Page missing")
	}
	r.Roots = []ast.NodeRef{page}
	output, err := e.Generate(context.Background(), r)
	if err != nil {
		t.Fatal(err)
	}
	tree, err := representation.Unpack[ts.Language](r.To, output.Snapshot().Inputs()[0])
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, node := range tree.Nodes {
		if node.Kind == ts.KindInterface || node.Kind == ts.KindTypeAlias {
			if nodeAttribute(node, ts.AttributeName) != "Page" {
				t.Fatalf("published an unrequested declaration: %+v", node)
			}
			count++
		}
	}
	if count != 1 {
		t.Fatal("Page was not published exactly once", count)
	}
}
