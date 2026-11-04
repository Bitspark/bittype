package semantics_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func vocabulary(name string) semantics.Key {
	return semantics.Key{Owner: "independent-fixture", Name: name}
}

func occurrence(node ast.NodeID, role string, ordinal uint32) semantics.Occurrence {
	return semantics.Occurrence{
		Site: ast.NodeRef{Tree: "current", Node: node},
		Role: vocabulary(role), Ordinal: ordinal,
	}
}

func currentSite(site ast.NodeRef) bool {
	return site.Tree == "current" && site.Node > 0 && site.Node <= 4
}

// This independent producer supplies an ordinary generic declaration and an
// application. Its byte payload deliberately contains NUL and non-UTF-8 bytes;
// neither the carrier nor its consumer gives the provider's encoding a meaning.
func exampleData() semantics.Data {
	use := occurrence(4, "use", 0)
	return semantics.Data{
		Entities: []semantics.Entity{
			{Kind: vocabulary("defined"), Name: "Sequence", Declarations: []ast.NodeRef{{Tree: "current", Node: 1}}, Parameters: []semantics.Parameter{{Entity: 2, Sort: vocabulary("type")}}, Body: 2},
			{Kind: vocabulary("parameter"), Name: "T", Owner: 1, Declarations: []ast.NodeRef{{Tree: "current", Node: 2}}},
			{Kind: vocabulary("builtin"), Name: "Scalar"},
			{Kind: vocabulary("alias"), Name: "Applied", Declarations: []ast.NodeRef{{Tree: "current", Node: 3}}, Body: 4},
		},
		Terms: []semantics.Term{
			{Kind: semantics.Reference, Target: 2},
			{Kind: semantics.ProviderExpression, Operator: vocabulary("sequence"), Encoding: vocabulary("opaque-v1"), Data: []byte{'A', 0, 255}, Operands: []semantics.Operand{{Locator: []uint32{7, 0}, Term: 1}}},
			{Kind: semantics.Reference, Target: 3},
			{Kind: semantics.Application, Callee: 6, Arguments: []semantics.TermID{3}},
			{Kind: semantics.Application, Callee: 3},
			{Kind: semantics.Reference, Target: 1},
		},
		Bindings: []semantics.Binding{
			{At: use, Target: 1},
			{At: occurrence(4, "use", 1), Target: 3},
		},
		Expressions: []semantics.Expression{{At: occurrence(4, "expression", 0), Term: 4}},
		Dependencies: []semantics.Dependency{
			{From: 4, To: 1, Relation: vocabulary("mentions"), At: &use},
			{From: 4, To: 3, Relation: vocabulary("argument")},
		},
		Facts: []semantics.Fact{
			{Subject: semantics.Subject{Entity: 1}, Relation: vocabulary("representation"), Term: 2},
			{Subject: semantics.Subject{Entity: 1}, Relation: vocabulary("representation"), Term: 6},
			{Subject: semantics.Subject{Term: 4}, Relation: vocabulary("result"), Term: 3},
			{Subject: semantics.Subject{Entity: 4}, Relation: vocabulary("result"), Term: 6},
			{Subject: semantics.Subject{Entity: 3}, Relation: vocabulary("intrinsic"), Term: 3},
		},
	}
}

func requireModel(t *testing.T, data semantics.Data) semantics.Semantics[struct{}] {
	t.Helper()
	model, err := semantics.New[struct{}](data, semantics.DefaultLimits(), currentSite)
	if err != nil || !model.Valid() {
		t.Fatalf("independent valid graph was rejected: %v", err)
	}
	return model
}

func TestSemanticModelQueriesPreserveIdentity(t *testing.T) {
	model := requireModel(t, exampleData())
	if target, ok := model.Resolve(occurrence(4, "use", 0)); !ok || target != 1 {
		t.Fatalf("declaration use: target %d, found %v", target, ok)
	}
	if target, ok := model.Resolve(occurrence(4, "use", 1)); !ok || target != 3 {
		t.Fatalf("second semantic position in same syntax node: %d, %v", target, ok)
	}
	if _, ok := model.Resolve(occurrence(4, "unrecorded", 0)); ok {
		t.Fatal("absence became an established binding")
	}
	id, ok := model.TermAt(occurrence(4, "expression", 0))
	if !ok || id != 4 {
		t.Fatalf("application expression: %d, %v", id, ok)
	}
	bindings, ok := model.DirectParameterBindings(id)
	want := []semantics.ArgumentBinding{{Parameter: 2, Argument: 3}}
	if !ok || !reflect.DeepEqual(bindings, want) {
		t.Fatalf("effective generic binding: %#v, %v", bindings, ok)
	}
	if zero, ok := model.DirectParameterBindings(5); !ok || len(zero) != 0 {
		t.Fatalf("zero-argument application lost its application form: %#v, %v", zero, ok)
	}
	if _, ok := model.DirectParameterBindings(6); ok {
		t.Fatal("reference to an open declaration became an application")
	}
	dependencies := model.DependenciesFrom(4, vocabulary("mentions"))
	if len(dependencies) != 1 || dependencies[0].To != 1 {
		t.Fatalf("dependency query mixed relation vocabularies: %#v", dependencies)
	}
	if got := model.DependenciesFrom(1, vocabulary("mentions")); len(got) != 0 {
		t.Fatalf("dependency query invented reverse edges: %#v", got)
	}
	builtin, ok := model.Entity(3)
	if !ok || len(builtin.Declarations) != 0 {
		t.Fatal("source-less builtin needs no invented syntax site")
	}
	term, _ := model.Term(2)
	if !reflect.DeepEqual(term.Data, []byte{'A', 0, 255}) || !reflect.DeepEqual(term.Operands[0].Locator, []uint32{7, 0}) {
		t.Fatal("provider payload or operand locator was interpreted or normalized")
	}
	if got := model.Facts(semantics.Subject{Entity: 1}, vocabulary("representation")); !reflect.DeepEqual(got, []semantics.TermID{2, 6}) {
		t.Fatalf("multi-valued relation lost provider order: %v", got)
	}
	if got := model.Facts(semantics.Subject{Term: 4}, vocabulary("result")); !reflect.DeepEqual(got, []semantics.TermID{3}) {
		t.Fatalf("term subject confused with same-numbered entity: %v", got)
	}
	if got := model.Facts(semantics.Subject{Entity: 3}, vocabulary("intrinsic")); !reflect.DeepEqual(got, []semantics.TermID{3}) {
		t.Fatalf("source-less subject needs no invented syntax: %v", got)
	}
	if got := model.Facts(semantics.Subject{Entity: 1}, vocabulary("unknown")); len(got) != 0 {
		t.Fatal("unrecorded relation acquired a fact")
	}
}

func mutateGraph(data semantics.Data) {
	data.Entities[0].Name = "changed"
	data.Entities[0].Declarations[0].Node = 99
	data.Entities[0].Parameters[0].Entity = 99
	data.Entities[0].Parameters[0].Sort.Name = "changed"
	data.Terms[1].Data[0] = 9
	data.Terms[1].Operands[0].Locator[0] = 99
	data.Terms[1].Operands[0].Term = 99
	data.Terms[3].Arguments[0] = 99
	data.Terms[3].Callee = 99
	data.Bindings[0].Target = 99
	data.Expressions[0].Term = 99
	data.Dependencies[0].To = 99
	data.Dependencies[0].At.Site.Node = 99
	data.Facts[0].Subject = semantics.Subject{Term: 99}
	data.Facts[0].Relation.Name = "changed"
	data.Facts[0].Term = 99
}

func TestSemanticModelOwnsEveryMutableIngressAndAccessor(t *testing.T) {
	input := exampleData()
	model := requireModel(t, input)
	expected := exampleData()
	mutateGraph(input)
	if !reflect.DeepEqual(model.Data(), expected) {
		t.Fatal("producer retained mutable storage in the published graph")
	}
	mutateGraph(model.Data())
	entity, _ := model.Entity(1)
	entity.Declarations[0].Node = 88
	entity.Parameters[0].Entity = 88
	term, _ := model.Term(2)
	term.Data[0] = 88
	term.Operands[0].Locator[0] = 88
	term.Operands[0].Term = 88
	application, _ := model.Term(4)
	application.Arguments[0] = 88
	application.Callee = 88
	dependencies := model.DependenciesFrom(4, vocabulary("mentions"))
	dependencies[0].To = 88
	dependencies[0].At.Site.Node = 88
	bindings, _ := model.DirectParameterBindings(4)
	bindings[0].Argument = 88
	facts := model.Facts(semantics.Subject{Entity: 1}, vocabulary("representation"))
	facts[0] = 88
	copyOfModel := model
	if !reflect.DeepEqual(copyOfModel.Data(), expected) {
		t.Fatal("accessor or copied model handle exposed private mutable storage")
	}
}

func TestSemanticModelRejectsBrokenGraphs(t *testing.T) {
	cases := []struct {
		name   string
		change func(*semantics.Data)
	}{
		{"missing entity vocabulary", func(d *semantics.Data) { d.Entities[0].Kind.Owner = "" }},
		{"invalid display encoding", func(d *semantics.Data) { d.Entities[0].Name = string([]byte{255}) }},
		{"unknown owner", func(d *semantics.Data) { d.Entities[0].Owner = 99 }},
		{"self owner", func(d *semantics.Data) { d.Entities[0].Owner = 1 }},
		{"unknown meaning", func(d *semantics.Data) { d.Entities[0].Body = 99 }},
		{"duplicate declaration site", func(d *semantics.Data) {
			d.Entities[0].Declarations = append(d.Entities[0].Declarations, d.Entities[0].Declarations[0])
		}},
		{"historical-only declaration site", func(d *semantics.Data) { d.Entities[0].Declarations[0].Tree = "historical" }},
		{"wrong parameter owner", func(d *semantics.Data) { d.Entities[1].Owner = 4 }},
		{"parameter without a sort", func(d *semantics.Data) { d.Entities[0].Parameters[0].Sort = semantics.Key{} }},
		{"unknown term variant", func(d *semantics.Data) { d.Terms[0].Kind = 0 }},
		{"unknown reference target", func(d *semantics.Data) { d.Terms[0].Target = 99 }},
		{"reference with arguments", func(d *semantics.Data) { d.Terms[0].Arguments = []semantics.TermID{3} }},
		{"reference with callee", func(d *semantics.Data) { d.Terms[0].Callee = 3 }},
		{"reference with opaque payload", func(d *semantics.Data) { d.Terms[0].Data = []byte{1} }},
		{"application with entity target", func(d *semantics.Data) { d.Terms[3].Target = 1 }},
		{"application missing callee", func(d *semantics.Data) { d.Terms[3].Callee = 0 }},
		{"application unknown callee", func(d *semantics.Data) { d.Terms[3].Callee = 99 }},
		{"application with opaque payload", func(d *semantics.Data) { d.Terms[3].Data = []byte{1} }},
		{"application with provider operator", func(d *semantics.Data) { d.Terms[3].Operator = vocabulary("call") }},
		{"application with provider encoding", func(d *semantics.Data) { d.Terms[3].Encoding = vocabulary("opaque") }},
		{"application with provider operands", func(d *semantics.Data) { d.Terms[3].Operands = []semantics.Operand{{Term: 3}} }},
		{"application unknown argument", func(d *semantics.Data) { d.Terms[3].Arguments[0] = 99 }},
		{"provider term with entity target", func(d *semantics.Data) { d.Terms[1].Target = 1 }},
		{"provider term with callee", func(d *semantics.Data) { d.Terms[1].Callee = 3 }},
		{"provider term without encoding", func(d *semantics.Data) { d.Terms[1].Encoding = semantics.Key{} }},
		{"provider term without operator", func(d *semantics.Data) { d.Terms[1].Operator = semantics.Key{} }},
		{"provider term unknown operand", func(d *semantics.Data) { d.Terms[1].Operands[0].Term = 99 }},
		{"duplicate operand locator", func(d *semantics.Data) { d.Terms[1].Operands = append(d.Terms[1].Operands, d.Terms[1].Operands[0]) }},
		{"unknown binding target", func(d *semantics.Data) { d.Bindings[0].Target = 99 }},
		{"binding has no role", func(d *semantics.Data) { d.Bindings[0].At.Role = semantics.Key{} }},
		{"duplicate binding occurrence", func(d *semantics.Data) { d.Bindings = append(d.Bindings, d.Bindings[0]) }},
		{"unknown expression term", func(d *semantics.Data) { d.Expressions[0].Term = 99 }},
		{"duplicate expression occurrence", func(d *semantics.Data) { d.Expressions = append(d.Expressions, d.Expressions[0]) }},
		{"unknown dependency source", func(d *semantics.Data) { d.Dependencies[0].From = 99 }},
		{"unknown dependency target", func(d *semantics.Data) { d.Dependencies[0].To = 99 }},
		{"dependency without relation", func(d *semantics.Data) { d.Dependencies[0].Relation = semantics.Key{} }},
		{"dependency with unadmitted occurrence", func(d *semantics.Data) { d.Dependencies[0].At.Site.Tree = "historical" }},
		{"fact without subject", func(d *semantics.Data) { d.Facts[0].Subject = semantics.Subject{} }},
		{"fact with both subject variants", func(d *semantics.Data) { d.Facts[0].Subject.Term = 2 }},
		{"fact with unknown entity subject", func(d *semantics.Data) { d.Facts[0].Subject.Entity = 99 }},
		{"fact with unknown term subject", func(d *semantics.Data) { d.Facts[2].Subject.Term = 99 }},
		{"fact without relation", func(d *semantics.Data) { d.Facts[0].Relation = semantics.Key{} }},
		{"fact without value", func(d *semantics.Data) { d.Facts[0].Term = 0 }},
		{"fact with unknown value", func(d *semantics.Data) { d.Facts[0].Term = 99 }},
		{"duplicate fact triple", func(d *semantics.Data) { d.Facts = append(d.Facts, d.Facts[0]) }},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			data := exampleData()
			test.change(&data)
			model, err := semantics.New[struct{}](data, semantics.DefaultLimits(), currentSite)
			if !errors.Is(err, semantics.ErrInvalid) || model.Valid() {
				t.Fatalf("malformed graph published: valid=%v, err=%v", model.Valid(), err)
			}
		})
	}
}

func TestSemanticModelRequiresOneBasedCurrentSites(t *testing.T) {
	for _, site := range []ast.NodeRef{{Tree: "current", Node: 0}, {Tree: "", Node: 1}} {
		data := exampleData()
		data.Entities[0].Declarations[0] = site
		model, err := semantics.New[struct{}](data, semantics.DefaultLimits(), func(ast.NodeRef) bool { return true })
		if !errors.Is(err, semantics.ErrInvalid) || model.Valid() {
			t.Fatalf("invalid address admitted by permissive catalogue: %#v, %v", site, err)
		}
	}
	if model, err := semantics.New[struct{}](exampleData(), semantics.DefaultLimits(), nil); !errors.Is(err, semantics.ErrInvalid) || model.Valid() {
		t.Fatalf("source sites admitted without a catalogue: %v", err)
	}
	model := requireModel(t, exampleData())
	if _, ok := model.Entity(0); ok {
		t.Fatal("zero semantic entity ID was accepted")
	}
	if _, ok := model.Term(0); ok {
		t.Fatal("zero semantic term ID was accepted")
	}
}

func TestSemanticModelLimitsAndNoPartialOutput(t *testing.T) {
	base := semantics.Data{
		Entities: []semantics.Entity{{Kind: semantics.Key{Owner: "k", Name: "d"}}},
		Terms:    []semantics.Term{{Kind: semantics.Reference, Target: 1}},
	}
	exact := semantics.Limits{MaxRecords: 2, MaxEdges: 1, MaxBytes: 2, MaxDepth: 1}
	if model, err := semantics.New[struct{}](base, exact, nil); err != nil || !model.Valid() {
		t.Fatalf("exact record/byte/depth boundary: %v", err)
	}
	for _, test := range []struct {
		name   string
		limits semantics.Limits
	}{
		{"records", semantics.Limits{MaxRecords: 1, MaxEdges: 1, MaxBytes: 2, MaxDepth: 1}},
		{"bytes", semantics.Limits{MaxRecords: 2, MaxEdges: 1, MaxBytes: 1, MaxDepth: 1}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if model, err := semantics.New[struct{}](base, test.limits, nil); !errors.Is(err, semantics.ErrLimit) || model.Valid() {
				t.Fatalf("exhaustion published graph or became semantic invalidity: %v", err)
			}
		})
	}
	for _, limits := range []semantics.Limits{{}, {MaxRecords: 2, MaxEdges: 1, MaxBytes: 2, MaxDepth: 513}} {
		if model, err := semantics.New[struct{}](base, limits, nil); !errors.Is(err, semantics.ErrInvalid) || model.Valid() {
			t.Fatalf("invalid configuration confused with exhaustion: %v", err)
		}
	}
	data := exampleData()
	limits := semantics.DefaultLimits()
	limits.MaxRecords = 1
	called := false
	model, err := semantics.New[struct{}](data, limits, func(ast.NodeRef) bool { called = true; return true })
	if !errors.Is(err, semantics.ErrLimit) || model.Valid() || called {
		t.Fatalf("oversized graph reached provenance validation: valid=%v called=%v err=%v", model.Valid(), called, err)
	}
}

func TestSemanticModelCountsOperandLocatorStorage(t *testing.T) {
	data := semantics.Data{
		Entities: []semantics.Entity{{Kind: vocabulary("builtin")}},
		Terms: []semantics.Term{
			{Kind: semantics.Reference, Target: 1},
			{Kind: semantics.ProviderExpression, Operator: vocabulary("extension"), Encoding: vocabulary("opaque"), Operands: []semantics.Operand{{Locator: []uint32{4, 3, 2, 1}, Term: 1}}},
		},
	}
	limits := semantics.DefaultLimits()
	limits.MaxEdges = 5 // one operand plus four locator entries
	if _, err := semantics.New[struct{}](data, limits, nil); err != nil {
		t.Fatalf("exact locator-storage boundary rejected: %v", err)
	}
	limits.MaxEdges = 4
	if model, err := semantics.New[struct{}](data, limits, nil); !errors.Is(err, semantics.ErrLimit) || model.Valid() {
		t.Fatalf("unbounded locator storage: %v", err)
	}
}

func TestSemanticModelDAGDepthAndEntityRecursion(t *testing.T) {
	wrap := func(children ...semantics.TermID) semantics.Term {
		term := semantics.Term{Kind: semantics.ProviderExpression, Operator: vocabulary("wrap"), Encoding: vocabulary("opaque")}
		for i, child := range children {
			term.Operands = append(term.Operands, semantics.Operand{Locator: []uint32{uint32(i)}, Term: child})
		}
		return term
	}
	// Earlier terms are visited first. A later parent reuses a memoized subtree;
	// checking only newly visited nodes would miss its longer containment path.
	data := semantics.Data{
		Entities:     []semantics.Entity{{Kind: vocabulary("recursive"), Name: "Node", Body: 4}},
		Terms:        []semantics.Term{{Kind: semantics.Reference, Target: 1}, wrap(1), wrap(2), wrap(1, 3)},
		Dependencies: []semantics.Dependency{{From: 1, To: 1, Relation: vocabulary("mentions")}},
	}
	limits := semantics.DefaultLimits()
	limits.MaxDepth = 4
	if model, err := semantics.New[struct{}](data, limits, nil); err != nil || !model.Valid() {
		t.Fatalf("finite term DAG with recursive entity meaning rejected: %v", err)
	}
	limits.MaxDepth = 3
	if model, err := semantics.New[struct{}](data, limits, nil); !errors.Is(err, semantics.ErrLimit) || model.Valid() {
		t.Fatalf("shared-subtree height ignored: %v", err)
	}
	data.Terms[1] = wrap(2)
	if model, err := semantics.New[struct{}](data, semantics.DefaultLimits(), nil); !errors.Is(err, semantics.ErrInvalid) || model.Valid() {
		t.Fatalf("term containment cycle accepted: %v", err)
	}
}

func TestSemanticZeroModelQueriesAreEmpty(t *testing.T) {
	var model semantics.Semantics[struct{}]
	if model.Valid() || !reflect.DeepEqual(model.Data(), semantics.Data{}) {
		t.Fatal("zero model became a published graph")
	}
	if _, ok := model.Entity(1); ok {
		t.Fatal("zero model entity")
	}
	if _, ok := model.Term(1); ok {
		t.Fatal("zero model term")
	}
	if _, ok := model.Resolve(occurrence(1, "use", 0)); ok {
		t.Fatal("zero model binding")
	}
	if _, ok := model.TermAt(occurrence(1, "expression", 0)); ok {
		t.Fatal("zero model expression")
	}
	if _, ok := model.DirectParameterBindings(1); ok {
		t.Fatal("zero model application")
	}
	if len(model.DependenciesFrom(1, vocabulary("mentions"))) != 0 {
		t.Fatal("zero model dependencies")
	}
	if len(model.Facts(semantics.Subject{Entity: 1}, vocabulary("representation"))) != 0 {
		t.Fatal("zero model facts")
	}
}
