package semantics_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Bitspark/bittype/semantics/go"
)

func TestApplicationAdmissionLeavesApplicabilityToProviders(t *testing.T) {
	for name, application := range map[string]semantics.Term{
		"fewer arguments":         {Kind: semantics.Application, Callee: 6},
		"more arguments":          {Kind: semantics.Application, Callee: 6, Arguments: []semantics.TermID{3, 3}},
		"parameter":               {Kind: semantics.Application, Callee: 1, Arguments: []semantics.TermID{3}},
		"zero-argument parameter": {Kind: semantics.Application, Callee: 1},
		"provider callee":         {Kind: semantics.Application, Callee: 2, Arguments: []semantics.TermID{3}},
		"computed callee":         {Kind: semantics.Application, Callee: 4, Arguments: []semantics.TermID{3}},
	} {
		t.Run(name, func(t *testing.T) {
			data := exampleData()
			data.Terms = append(data.Terms, application)
			id := semantics.TermID(len(data.Terms))
			model := requireModel(t, data)
			got, ok := model.Term(id)
			if !ok || !reflect.DeepEqual(got, application) {
				t.Fatalf("callee expression or arguments changed: %#v", got)
			}
			if bindings, ok := model.DirectParameterBindings(id); ok || bindings != nil {
				t.Fatalf("non-direct or unmatched application invented formal bindings: %#v", bindings)
			}
			parameter, _ := model.Entity(2)
			if len(parameter.Parameters) != 0 {
				t.Fatal("constructor parameter acquired fabricated formal entities")
			}
		})
	}
	model := requireModel(t, exampleData())
	for _, id := range []semantics.TermID{0, 1, 2, 6, 99} {
		if bindings, ok := model.DirectParameterBindings(id); ok || bindings != nil {
			t.Fatalf("non-application %d acquired bindings: %#v", id, bindings)
		}
	}
	data := exampleData()
	data.Entities[2].Owner = 1
	model = requireModel(t, data)
	if bindings, ok := model.DirectParameterBindings(5); !ok || len(bindings) != 0 {
		t.Fatal("ordinary entity ownership was mistaken for parameter membership")
	}
}

func TestCalleeContainmentDepthCyclesAndSharedSubtrees(t *testing.T) {
	data := semantics.Data{
		Entities: []semantics.Entity{{Kind: vocabulary("intrinsic")}},
		Terms: []semantics.Term{
			{Kind: semantics.Reference, Target: 1},
			{Kind: semantics.Application, Callee: 1},
			{Kind: semantics.Application, Callee: 2, Arguments: []semantics.TermID{1}},
			{Kind: semantics.Application, Callee: 1, Arguments: []semantics.TermID{3}},
		},
	}
	limits := semantics.DefaultLimits()
	limits.MaxDepth = 4
	if model, err := semantics.New[struct{}](data, limits, nil); err != nil || !model.Valid() {
		t.Fatalf("finite computed-callee DAG rejected: %v", err)
	}
	limits.MaxDepth = 3
	if model, err := semantics.New[struct{}](data, limits, nil); !errors.Is(err, semantics.ErrLimit) || model.Valid() {
		t.Fatalf("memoized callee height lost at an argument join: %v", err)
	}
	for _, target := range []semantics.TermID{2, 3} {
		data.Terms[1].Callee = target
		if model, err := semantics.New[struct{}](data, semantics.DefaultLimits(), nil); !errors.Is(err, semantics.ErrInvalid) || model.Valid() {
			t.Fatalf("callee cycle through term %d admitted: %v", target, err)
		}
	}
	data.Terms[1] = semantics.Term{Kind: semantics.ProviderExpression, Operator: vocabulary("computed"), Encoding: vocabulary("opaque"), Operands: []semantics.Operand{{Term: 3}}}
	if model, err := semantics.New[struct{}](data, semantics.DefaultLimits(), nil); !errors.Is(err, semantics.ErrInvalid) || model.Valid() {
		t.Fatalf("mixed callee/provider-operand cycle admitted: %v", err)
	}
}

func TestCalleeAndFactStorageBudgets(t *testing.T) {
	data := semantics.Data{
		Entities: []semantics.Entity{{Kind: semantics.Key{Owner: "k", Name: "d"}}},
		Terms: []semantics.Term{
			{Kind: semantics.Reference, Target: 1},
			{Kind: semantics.Application, Callee: 1, Arguments: []semantics.TermID{1}},
		},
		Facts: []semantics.Fact{{Subject: semantics.Subject{Entity: 1}, Relation: semantics.Key{Owner: "r", Name: "f"}, Term: 2}},
	}
	// Four records, one callee, one argument, two fact links, and two keys.
	exact := semantics.Limits{MaxRecords: 4, MaxEdges: 4, MaxBytes: 4, MaxDepth: 2}
	model, err := semantics.New[struct{}](data, exact, nil)
	if err != nil || !model.Valid() {
		t.Fatalf("exact fact/callee boundary rejected: %v", err)
	}
	if err := model.CheckLimits(exact); err != nil {
		t.Fatalf("sealed model changed its exact budget: %v", err)
	}
	for name, limits := range map[string]semantics.Limits{
		"fact record":  {MaxRecords: 3, MaxEdges: 4, MaxBytes: 4, MaxDepth: 2},
		"linked edges": {MaxRecords: 4, MaxEdges: 3, MaxBytes: 4, MaxDepth: 2},
		"relation key": {MaxRecords: 4, MaxEdges: 4, MaxBytes: 3, MaxDepth: 2},
		"callee depth": {MaxRecords: 4, MaxEdges: 4, MaxBytes: 4, MaxDepth: 1},
	} {
		t.Run(name, func(t *testing.T) {
			if result, err := semantics.New[struct{}](data, limits, nil); !errors.Is(err, semantics.ErrLimit) || result.Valid() {
				t.Fatalf("over-budget graph published: %v", err)
			}
			if err := model.CheckLimits(limits); !errors.Is(err, semantics.ErrLimit) {
				t.Fatalf("reuse bypassed new storage bounds: %v", err)
			}
		})
	}
}

func TestSubjectFactsAreAssociationsNotContainment(t *testing.T) {
	data := semantics.Data{
		Terms: []semantics.Term{
			{Kind: semantics.ProviderExpression, Operator: vocabulary("left"), Encoding: vocabulary("opaque")},
			{Kind: semantics.ProviderExpression, Operator: vocabulary("right"), Encoding: vocabulary("opaque")},
		},
		Facts: []semantics.Fact{
			{Subject: semantics.Subject{Term: 1}, Relation: vocabulary("related"), Term: 2},
			{Subject: semantics.Subject{Term: 2}, Relation: vocabulary("related"), Term: 1},
			{Subject: semantics.Subject{Term: 1}, Relation: vocabulary("related"), Term: 1},
		},
	}
	limits := semantics.DefaultLimits()
	limits.MaxDepth = 1
	model, err := semantics.New[struct{}](data, limits, nil)
	if err != nil || !model.Valid() {
		t.Fatalf("associations became term containment: %v", err)
	}
	if got := model.Facts(semantics.Subject{Term: 1}, vocabulary("related")); !reflect.DeepEqual(got, []semantics.TermID{2, 1}) {
		t.Fatalf("association order or identity changed: %v", got)
	}
	for _, subject := range []semantics.Subject{{}, {Entity: 1}, {Term: 99}, {Entity: 1, Term: 1}} {
		if got := model.Facts(subject, vocabulary("related")); len(got) != 0 {
			t.Fatalf("invalid subject matched established facts: %#v", subject)
		}
	}
}
