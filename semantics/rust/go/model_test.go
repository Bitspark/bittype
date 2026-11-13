package rust

import (
	"errors"
	"testing"

	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
)

// The entities of every graph below: the crate root, a struct and a parameter.
const (
	rootEntity semantics.EntityID = iota + 1
	itemEntity
	otherEntity
)

func expression(operator semantics.Key, data string, operands ...semantics.TermID) semantics.Term {
	term := semantics.Term{Kind: semantics.ProviderExpression, Operator: operator, Encoding: Key("json-1"), Data: []byte(data)}
	for i, operand := range operands {
		term.Operands = append(term.Operands, semantics.Operand{Locator: []uint32{uint32(i)}, Term: operand})
	}
	return term
}
func reference(target semantics.EntityID) semantics.Term {
	return semantics.Term{Kind: semantics.Reference, Target: target}
}

// declaration is the struct of a graph that holds terms and, about the
// struct, one fact of relation for each of the listed terms.
func declaration(t *testing.T, relation semantics.Key, terms []semantics.Term, facts ...semantics.TermID) Declaration {
	t.Helper()
	data := semantics.Data{
		Entities: []semantics.Entity{
			{Kind: ModuleKind(), Name: "pages"},
			{Kind: Key("struct"), Name: "pages::A", Owner: rootEntity},
			{Kind: Key("parameter"), Name: "pages::A::T", Owner: itemEntity},
		},
		Terms: terms,
	}
	for _, term := range facts {
		data.Facts = append(data.Facts, semantics.Fact{Subject: semantics.Subject{Entity: itemEntity}, Relation: relation, Term: term})
	}
	graph, err := semantics.New[s.Language](data, semantics.DefaultLimits(), nil)
	if err != nil {
		t.Fatal(err)
	}
	return Declaration{&model{graph: graph}, itemEntity}
}

func TestDeclaredVisibilityReader(t *testing.T) {
	public := expression(PublicOperator(), "null")
	restricted := expression(RestrictedOperator(), "null", 1)
	d := declaration(t, DeclaredVisibilityRelation(), []semantics.Term{reference(rootEntity), public}, 2)
	if v, err := d.Visibility(); err != nil || !v.Public || v.Scope.Valid() {
		t.Fatal("public", v, err)
	}
	d = declaration(t, DeclaredVisibilityRelation(), []semantics.Term{reference(rootEntity), restricted}, 2)
	if v, err := d.Visibility(); err != nil || v.Public || v.Scope.ID() != rootEntity {
		t.Fatal("restricted", v, err)
	}
	for name, c := range map[string]struct {
		terms []semantics.Term
		facts []semantics.TermID
	}{
		"no fact":                                   {[]semantics.Term{public}, nil},
		"two facts":                                 {[]semantics.Term{public, public}, []semantics.TermID{1, 2}},
		"the fact of another relation":              {nil, nil},
		"an unknown operator":                       {[]semantics.Term{expression(Key("crate"), "null")}, []semantics.TermID{1}},
		"a reference instead of a form":             {[]semantics.Term{reference(rootEntity)}, []semantics.TermID{1}},
		"public with an operand":                    {[]semantics.Term{reference(rootEntity), expression(PublicOperator(), "null", 1)}, []semantics.TermID{2}},
		"restricted without a scope":                {[]semantics.Term{expression(RestrictedOperator(), "null")}, []semantics.TermID{1}},
		"restricted to two scopes":                  {[]semantics.Term{reference(rootEntity), expression(RestrictedOperator(), "null", 1, 1)}, []semantics.TermID{2}},
		"restricted to no reference":                {[]semantics.Term{public, expression(RestrictedOperator(), "null", 1)}, []semantics.TermID{2}},
		"restricted to an entity that is no module": {[]semantics.Term{reference(otherEntity), expression(RestrictedOperator(), "null", 1)}, []semantics.TermID{2}},
	} {
		relation := DeclaredVisibilityRelation()
		terms, facts := c.terms, c.facts
		if name == "the fact of another relation" {
			relation, terms, facts = EffectiveVisibilityRelation(), []semantics.Term{public}, []semantics.TermID{1}
		}
		if v, err := declaration(t, relation, terms, facts...).Visibility(); !errors.Is(err, ErrMissingFact) || v != (Visibility{}) {
			t.Errorf("%s: read as %+v, %v", name, v, err)
		}
	}
	if _, err := (Declaration{}).Visibility(); !errors.Is(err, ErrForeign) {
		t.Error("the zero declaration", err)
	}
}

func TestEffectiveVisibilityReader(t *testing.T) {
	stated := func(data string) Declaration {
		return declaration(t, EffectiveVisibilityRelation(), []semantics.Term{expression(EffectiveVisibilityOperator(), data)}, 1)
	}
	for data, nameable := range map[string]bool{
		`{"Audience":"another-crate","Path":"definition","Nameable":true}`:  true,
		`{"Audience":"another-crate","Path":"definition","Nameable":false}`: false,
	} {
		want := Effective{Audience: AudienceAnotherCrate, Path: PathDefinition, Nameable: nameable}
		if e, err := stated(data).Effective(); err != nil || e != want {
			t.Fatal(data, e, err)
		}
	}
	for name, data := range map[string]string{
		"no answer":         `{"Audience":"another-crate","Path":"definition"}`,
		"no audience":       `{"Path":"definition","Nameable":true}`,
		"no path":           `{"Audience":"another-crate","Nameable":true}`,
		"another audience":  `{"Audience":"this-crate","Path":"definition","Nameable":true}`,
		"another path":      `{"Audience":"another-crate","Path":"re-export","Nameable":true}`,
		"an unknown member": `{"Audience":"another-crate","Path":"definition","Nameable":true,"Through":1}`,
		"no object":         `true`,
		"no data":           ``,
	} {
		if e, err := stated(data).Effective(); !errors.Is(err, ErrMissingFact) || e != (Effective{}) {
			t.Errorf("%s: read as %+v, %v", name, e, err)
		}
	}
	whole := `{"Audience":"another-crate","Path":"definition","Nameable":true}`
	for name, d := range map[string]Declaration{
		"no fact":             declaration(t, EffectiveVisibilityRelation(), []semantics.Term{expression(EffectiveVisibilityOperator(), whole)}),
		"two facts":           declaration(t, EffectiveVisibilityRelation(), []semantics.Term{expression(EffectiveVisibilityOperator(), whole), expression(EffectiveVisibilityOperator(), whole)}, 1, 2),
		"an unknown operator": declaration(t, EffectiveVisibilityRelation(), []semantics.Term{expression(PublicOperator(), whole)}, 1),
		"an operand":          declaration(t, EffectiveVisibilityRelation(), []semantics.Term{reference(rootEntity), expression(EffectiveVisibilityOperator(), whole, 1)}, 2),
		"a reference":         declaration(t, EffectiveVisibilityRelation(), []semantics.Term{reference(rootEntity)}, 1),
	} {
		if e, err := d.Effective(); !errors.Is(err, ErrMissingFact) || e != (Effective{}) {
			t.Errorf("%s: read as %+v, %v", name, e, err)
		}
	}
	if _, err := (Declaration{}).Effective(); !errors.Is(err, ErrForeign) {
		t.Error("the zero declaration", err)
	}
}

func TestPrivateInterfaceReader(t *testing.T) {
	if observed, err := declaration(t, PrivateInterfaceRelation(), nil).PrivateInterface(); err != nil || len(observed) != 0 {
		t.Fatal("no observation", observed, err)
	}
	d := declaration(t, PrivateInterfaceRelation(), []semantics.Term{reference(rootEntity), reference(otherEntity)}, 1, 2)
	if observed, err := d.PrivateInterface(); err != nil || len(observed) != 2 || observed[0].ID() != rootEntity || observed[1].ID() != otherEntity {
		t.Fatal("two observations", observed, err)
	}
	d = declaration(t, PrivateInterfaceRelation(), []semantics.Term{reference(rootEntity), expression(PublicOperator(), "null")}, 1, 2)
	if observed, err := d.PrivateInterface(); !errors.Is(err, ErrMissingFact) || observed != nil {
		t.Error("an observation that names no declaration", observed, err)
	}
	if _, err := (Declaration{}).PrivateInterface(); !errors.Is(err, ErrForeign) {
		t.Error("the zero declaration", err)
	}
}
