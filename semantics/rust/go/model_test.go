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
		"no fact":                      {[]semantics.Term{public}, nil},
		"two facts":                    {[]semantics.Term{public, public}, []semantics.TermID{1, 2}},
		"the fact of another relation": {nil, nil},
		// The two forms carry no data: the helper writes null.
		"public with data":                          {[]semantics.Term{expression(PublicOperator(), `{"anything":1}`)}, []semantics.TermID{1}},
		"public without data":                       {[]semantics.Term{expression(PublicOperator(), ``)}, []semantics.TermID{1}},
		"restricted with data":                      {[]semantics.Term{reference(rootEntity), expression(RestrictedOperator(), `{"Scope":1}`, 1)}, []semantics.TermID{2}},
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
		// The order of the members and the space between them are no part of the form.
		`{"Nameable":true,"Path":"definition","Audience":"another-crate"}`:                true,
		` { "Audience" : "another-crate" , "Path" : "definition" , "Nameable" : false } `: false,
	} {
		want := Effective{Audience: AudienceAnotherCrate, Path: PathDefinition, Nameable: nameable}
		if e, err := stated(data).Effective(); err != nil || e != want {
			t.Fatal(data, e, err)
		}
	}
	for name, data := range map[string]string{
		"no answer":        `{"Audience":"another-crate","Path":"definition"}`,
		"no audience":      `{"Path":"definition","Nameable":true}`,
		"no path":          `{"Audience":"another-crate","Nameable":true}`,
		"another audience": `{"Audience":"this-crate","Path":"definition","Nameable":true}`,
		"another path":     `{"Audience":"another-crate","Path":"re-export","Nameable":true}`,
		// The form is the three members by exact name, each once, and nothing after them.
		"a repeated member":          `{"Audience":"another-crate","Path":"definition","Nameable":false,"Nameable":true}`,
		"a member repeated alike":    `{"Audience":"another-crate","Audience":"another-crate","Path":"definition","Nameable":true}`,
		"a member in lower case":     `{"audience":"another-crate","Path":"definition","Nameable":true}`,
		"the path in lower case":     `{"Audience":"another-crate","path":"definition","Nameable":true}`,
		"the answer in lower case":   `{"Audience":"another-crate","Path":"definition","nameable":true}`,
		"every member in lower case": `{"audience":"another-crate","path":"definition","nameable":true}`,
		"a second value":             `{"Audience":"another-crate","Path":"definition","Nameable":true}{"Nameable":false}`,
		"text after the value":       `{"Audience":"another-crate","Path":"definition","Nameable":true} trailing`,
		"a second closing brace":     `{"Audience":"another-crate","Path":"definition","Nameable":true}}`,
		"an answer that is null":     `{"Audience":"another-crate","Path":"definition","Nameable":null}`,
		"an answer that is text":     `{"Audience":"another-crate","Path":"definition","Nameable":"true"}`,
		"an answer that is a number": `{"Audience":"another-crate","Path":"definition","Nameable":1}`,
		"an audience that is null":   `{"Audience":null,"Path":"definition","Nameable":true}`,
		"a path that is a list":      `{"Audience":"another-crate","Path":["definition"],"Nameable":true}`,
		"a list":                     `[{"Audience":"another-crate","Path":"definition","Nameable":true}]`,
		"an unknown member":          `{"Audience":"another-crate","Path":"definition","Nameable":true,"Through":1}`,
		"no object":                  `true`,
		"no data":                    ``,
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

// The fields term describes each field by exactly its name and its entity.
// A type without the fact has no fields, which is no error. A fact that is
// present twice or in another form fails, as the visibility readers do.
func TestFieldsReader(t *testing.T) {
	// fields reads the type of a graph whose fields facts carry the given terms.
	fields := func(terms ...semantics.Term) ([]Field, bool, error) {
		t.Helper()
		data := semantics.Data{
			Entities: []semantics.Entity{
				{Kind: ModuleKind(), Name: "pages"},
				{Kind: Key("struct"), Name: "pages::A", Owner: rootEntity},
				{Kind: Key("field"), Name: "pages::A::value", Owner: itemEntity},
			},
			Terms: []semantics.Term{reference(itemEntity), reference(rootEntity)},
		}
		for _, term := range terms {
			data.Terms = append(data.Terms, term)
			data.Facts = append(data.Facts, semantics.Fact{Subject: semantics.Subject{Term: 1}, Relation: FieldsRelation(), Term: semantics.TermID(len(data.Terms))})
		}
		graph, err := semantics.New[s.Language](data, semantics.DefaultLimits(), nil)
		if err != nil {
			t.Fatal(err)
		}
		return Type{&model{graph: graph}, 1}.Fields()
	}
	one := func(data string, operands ...semantics.TermID) semantics.Term {
		return expression(FieldsOperator(), data, operands...)
	}
	read, has, err := fields(one(`[{"Name":"value","Entity":3}]`, 2))
	if err != nil || !has || len(read) != 1 || read[0].Name != "value" || read[0].Declaration.ID() != otherEntity || read[0].Type.ID() != 2 {
		t.Fatal(read, has, err)
	}
	if read, has, err := fields(one(` [ { "Entity" : 3 , "Name" : "value" } ] `, 2)); err != nil || !has || len(read) != 1 || read[0].Name != "value" || read[0].Declaration.ID() != otherEntity {
		t.Fatal("another order and spacing", read, has, err)
	}
	// A record without a field has the fact, with an empty list.
	if read, has, err := fields(one(`[]`)); err != nil || !has || len(read) != 0 {
		t.Fatal("a record without a field", read, has, err)
	}
	// A type without the fact has no fields. That is an answer and no error.
	if read, has, err := fields(); err != nil || has || read != nil {
		t.Fatal("a type without the fact", read, has, err)
	}
	if read, has, err := (Type{}).Fields(); err != nil || has || read != nil {
		t.Fatal("the zero type", read, has, err)
	}
	// A malformed fact is an error and never reads as "no fields".
	for name, terms := range map[string][]semantics.Term{
		"the member of the earlier profile":      {one(`[{"Name":"value","Public":true,"Entity":3}]`, 2)},
		"members in lower case":                  {one(`[{"name":"value","entity":3}]`, 2)},
		"a repeated member":                      {one(`[{"Name":"other","Name":"value","Entity":3}]`, 2)},
		"no entity":                              {one(`[{"Name":"value"}]`, 2)},
		"no name":                                {one(`[{"Entity":3}]`, 2)},
		"a name that is null":                    {one(`[{"Name":null,"Entity":3}]`, 2)},
		"an entity that is text":                 {one(`[{"Name":"value","Entity":"3"}]`, 2)},
		"an entity that is null":                 {one(`[{"Name":"value","Entity":null}]`, 2)},
		"an entity that is none":                 {one(`[{"Name":"value","Entity":0}]`, 2)},
		"an entity the graph does not have":      {one(`[{"Name":"value","Entity":9}]`, 2)},
		"text after the list":                    {one(`[{"Name":"value","Entity":3}] trailing`, 2)},
		"a second list":                          {one(`[{"Name":"value","Entity":3}][]`, 2)},
		"null":                                   {one(`null`)},
		"no data":                                {one(``)},
		"an object":                              {one(`{"Name":"value","Entity":3}`, 2)},
		"a field that is no object":              {one(`["value"]`, 2)},
		"more fields than types":                 {one(`[{"Name":"value","Entity":3},{"Name":"other","Entity":3}]`, 2)},
		"more types than fields":                 {one(`[{"Name":"value","Entity":3}]`, 2, 2)},
		"another operator":                       {expression(PublicOperator(), `[{"Name":"value","Entity":3}]`, 2)},
		"a reference":                            {reference(otherEntity)},
		"two facts":                              {one(`[{"Name":"value","Entity":3}]`, 2), one(`[{"Name":"value","Entity":3}]`, 2)},
		"a well-formed fact and a malformed one": {one(`[{"Name":"value","Entity":3}]`, 2), one(`null`)},
	} {
		if read, has, err := fields(terms...); !errors.Is(err, ErrMissingFact) || has || read != nil {
			t.Errorf("%s: read as %+v, %v, %v", name, read, has, err)
		}
	}
}

// members is the strict decoder of both readers: exactly the named members,
// each once and spelled so, and nothing after the object.
func TestStrictMembers(t *testing.T) {
	values, ok := members([]byte(` {"B": [1, {"x": null}], "A": "text"} `), "A", "B")
	if !ok || len(values) != 2 || string(values[0]) != `"text"` || string(values[1]) != `[1, {"x": null}]` {
		t.Fatal(values, ok)
	}
	if values, ok := members([]byte(`{}`)); !ok || len(values) != 0 {
		t.Fatal("no member", values, ok)
	}
	if values, ok := members([]byte(`{"A":null}`), "A"); !ok || string(values[0]) != "null" {
		t.Fatal("a member that is null is present", values, ok)
	}
	for name, data := range map[string]string{
		"a missing member":             `{"A":1}`,
		"no member at all":             `{}`,
		"an unknown member":            `{"A":1,"B":2,"C":3}`,
		"a member in another case":     `{"a":1,"B":2}`,
		"a repeated member":            `{"A":1,"B":2,"A":1}`,
		"a member repeated first":      `{"A":1,"A":2,"B":3}`,
		"a second object":              `{"A":1,"B":2}{}`,
		"text after the object":        `{"A":1,"B":2} x`,
		"a second closing brace":       `{"A":1,"B":2}}`,
		"an object that is not closed": `{"A":1,"B":2`,
		"a value that is cut off":      `{"A":1,"B":`,
		"a list":                       `[{"A":1,"B":2}]`,
		"a list of names and values":   `["A",1,"B",2]`,
		"a text":                       `"A"`,
		"null":                         `null`,
		"nothing":                      ``,
	} {
		if values, ok := members([]byte(data), "A", "B"); ok || values != nil {
			t.Errorf("%s: read as %s", name, values)
		}
	}
}
