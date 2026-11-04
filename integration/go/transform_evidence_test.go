package integration_test

import (
	"context"
	"errors"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/engine/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"testing"
)

func TestTransformEvidenceAdmission(t *testing.T) {
	e, from, source, resolver := bindingFixture(t, nil)
	ctx := context.Background()
	settings := []semantic.Setting{{Key: analysisKey("configuration"), Value: []byte("selected")}}
	result, err := resolver.Resolve(ctx, source, semantic.Options{Configuration: settings})
	if err != nil {
		t.Fatal(err)
	}
	to := definition[blue](t, "output")
	requirement := semantic.EvidenceRequirement{Operation: resolution.Kind(), Profile: bindingContract().Semantic.Profile, Configuration: settings, Reports: bindingContract().Semantic.ModelRequirements, Model: true}
	calls := 0
	mapper, err := transform.NewTransformer[amber, blue](from, to, transform.Contract{Evidence: []semantic.EvidenceRequirement{requirement}}, func(_ context.Context, r transform.Request[amber]) (ast.AST[blue], error) {
		calls++
		evidence, ok := r.Evidence(0)
		if !ok || !evidence.HasModel() {
			t.Fatal("missing admitted graph")
		}
		if _, ok := r.Evidence(1); ok {
			t.Fatal("unexpected evidence")
		}
		one := r.Source()
		one.Nodes[0].Kind = "mutation"
		if r.Source().Nodes[0].Kind != "document" {
			t.Fatal("request source aliases storage")
		}
		return mapExample[amber, blue](r.Source())
	})
	if err != nil {
		t.Fatal(err)
	}
	// Constructor owns nested metadata. Caller mutation cannot change admission.
	requirement.Configuration[0].Value[0] = 'x'
	requirement.Reports[0].Verdict = semantic.Rejected
	valid := transform.Input{Source: source, Evidence: []semantic.Result{result}}
	if _, err = mapper.Transform(ctx, valid); err != nil || calls != 1 {
		t.Fatal(err)
	}
	fresh, err := e.Capture(ctx, semantic.Input{Entry: source.Entry(), Trees: source.Inputs()})
	if err != nil {
		t.Fatal(err)
	}
	changed, err := resolver.Resolve(ctx, source, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	for name, input := range map[string]transform.Input{
		"missing": {Source: source}, "extra": {Source: source, Evidence: []semantic.Result{result, result}},
		"stale": {Source: fresh, Evidence: []semantic.Result{result}}, "wrong-config": {Source: source, Evidence: []semantic.Result{changed}},
		"empty": {},
	} {
		t.Run(name, func(t *testing.T) {
			out, err := mapper.Transform(ctx, input)
			if err == nil || out.LanguageID() != "" || calls != 1 {
				t.Fatalf("invalid evidence reached mapping: %v", err)
			}
		})
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if out, err := mapper.Transform(canceled, valid); !errors.Is(err, context.Canceled) || out.LanguageID() != "" || calls != 1 {
		t.Fatal(err)
	}
	if out, err := mapper.Transform(nil, valid); err == nil || out.LanguageID() != "" || calls != 1 {
		t.Fatal("nil context accepted")
	}
	// Check wrong operation/profile and undisclosed requirements through the shared admission API.
	freshRequirement := semantic.EvidenceRequirement{Operation: resolution.Kind(), Profile: bindingContract().Semantic.Profile, Reports: bindingContract().Semantic.ModelRequirements, Model: true}
	for name, mutate := range map[string]func(*semantic.EvidenceRequirement){
		"operation": func(r *semantic.EvidenceRequirement) { r.Operation = analysisKey("other") },
		"profile":   func(r *semantic.EvidenceRequirement) { r.Profile = "another-profile" },
		"report": func(r *semantic.EvidenceRequirement) {
			r.Reports = []semantic.Requirement{{Obligation: analysisKey("unadvertised"), Verdict: semantic.Holds}}
		},
	} {
		t.Run(name, func(t *testing.T) {
			r := freshRequirement.Clone()
			mutate(&r)
			if _, err := semantic.AdmitEvidence(source, []semantic.EvidenceRequirement{r}, []semantic.Result{changed}); err == nil {
				t.Fatal("incompatible requirement accepted")
			}
		})
	}
}

func TestTransformRejectsUnestablishedEvidence(t *testing.T) {
	for _, tc := range []struct {
		name    string
		status  semantic.Completion
		verdict semantic.Verdict
	}{
		{"incomplete", semantic.Incomplete, semantic.Undetermined}, {"unsupported", semantic.Unsupported, semantic.Undetermined}, {"rejected", semantic.Complete, semantic.Rejected}, {"missing-model", semantic.Complete, semantic.Holds},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, from, source, resolver := bindingFixture(t, func(_ context.Context, r semantic.Request[amber]) (semantic.Result, error) {
				return r.Publish(semantic.Publication{Reports: []semantic.Report{{Obligation: bindingKey(), Status: tc.status, Verdict: tc.verdict}}})
			})
			result, err := resolver.Resolve(context.Background(), source, semantic.Options{})
			if err != nil {
				t.Fatal(err)
			}
			requirement := semantic.EvidenceRequirement{Operation: resolution.Kind(), Profile: bindingContract().Semantic.Profile, Reports: bindingContract().Semantic.ModelRequirements, Model: true}
			called := false
			mapper, err := transform.NewTransformer[amber, blue](from, definition[blue](t, "output"), transform.Contract{Evidence: []semantic.EvidenceRequirement{requirement}}, func(context.Context, transform.Request[amber]) (ast.AST[blue], error) {
				called = true
				return ast.AST[blue]{}, nil
			})
			if err != nil {
				t.Fatal(err)
			}
			output, err := mapper.Transform(context.Background(), transform.Input{Source: source, Evidence: []semantic.Result{result}})
			if !errors.Is(err, semantic.ErrEvidenceUnavailable) || called || output.LanguageID() != "" {
				t.Fatalf("unestablished receipt admitted: %v", err)
			}
		})
	}
}

func TestTransformRouteDoesNotForwardEvidence(t *testing.T) {
	e, from, source, resolver := bindingFixture(t, nil)
	ctx := context.Background()
	resolved, err := resolver.Resolve(ctx, source, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	middle, last := definition[blue](t, "middle"), definition[coral](t, "last")
	registerChecked[blue](t, e, middle)
	registerChecked[coral](t, e, last)
	contract := transform.Contract{Evidence: []semantic.EvidenceRequirement{{Operation: resolution.Kind(), Profile: bindingContract().Semantic.Profile, Reports: bindingContract().Semantic.ModelRequirements, Model: true}}}
	first, err := transform.NewTransformer[amber, blue](from, middle, contract, func(_ context.Context, r transform.Request[amber]) (ast.AST[blue], error) {
		return mapExample[amber, blue](r.Source())
	})
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	second, err := transform.NewTransformer[blue, coral](middle, last, contract, func(_ context.Context, r transform.Request[blue]) (ast.AST[coral], error) {
		calls++
		return mapExample[blue, coral](r.Source())
	})
	if err != nil {
		t.Fatal(err)
	}
	connect(t, e, first, second)
	output, err := e.TransformTo(ctx, transform.Input{Source: source, Evidence: []semantic.Result{resolved}}, last.ID())
	if err == nil || calls != 0 || output.LanguageID() != "" {
		t.Fatal("source evidence crossed an intermediate transformation")
	}
}

func TestTransformCancellationAfterCallback(t *testing.T) {
	a, b := definition[amber](t, "a"), definition[blue](t, "b")
	ctx, cancel := context.WithCancel(context.Background())
	stage, err := transform.NewTransformer[amber, blue](a, b, transform.Contract{}, func(_ context.Context, r transform.Request[amber]) (ast.AST[blue], error) {
		cancel()
		return mapExample[amber, blue](r.Source())
	})
	if err != nil {
		t.Fatal(err)
	}
	output, err := stage.Transform(ctx, captureForTransform[amber](t, a, packed[amber](t, a)))
	if !errors.Is(err, context.Canceled) || output.LanguageID() != "" {
		t.Fatal(err)
	}
}

func TestTypedTransformationRetainsInputCollection(t *testing.T) {
	from, to := definition[amber](t, "source"), definition[blue](t, "target")
	e := engine.New()
	registerChecked[amber](t, e, from)
	first := exampleTree[amber]()
	second := exampleTree[amber]()
	second.ID = "second"
	one, err := representation.Pack(from, first)
	if err != nil {
		t.Fatal(err)
	}
	two, err := representation.Pack(from, second)
	if err != nil {
		t.Fatal(err)
	}
	source, err := e.Capture(context.Background(), semantic.Input{Entry: ast.NodeRef{Tree: second.ID, Node: second.Root}, Trees: []representation.AST{one, two}})
	if err != nil {
		t.Fatal(err)
	}
	mapper, err := transform.NewTransformer[amber, blue](from, to, transform.Contract{}, func(_ context.Context, r transform.Request[amber]) (ast.AST[blue], error) {
		if r.Source().ID != "second" || r.Entry() != source.Entry() || len(r.Inputs()) != 2 {
			t.Fatal("entry tree or source collection lost")
		}
		inputs := r.Inputs()
		inputs[0] = representation.AST{}
		if r.Inputs()[0].LanguageID() == "" {
			t.Fatal("input slice aliases request")
		}
		return mapExample[amber, blue](r.Source())
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = mapper.Transform(context.Background(), transform.Input{Source: source}); err != nil {
		t.Fatal(err)
	}
}
