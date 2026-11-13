package typescript

import (
	"testing"

	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
)

// The names of the body of an interface with base types are names of the
// published graph: the compiler helper writes them and Model.Extension reads
// them.
func TestPublishedNamesOfTheExtension(t *testing.T) {
	if got := ExtensionOperator(); got != (semantics.Key{Owner: "typescript", Name: "extension"}) {
		t.Errorf("ExtensionOperator is %v, want typescript/extension", got)
	}
	if got := JSONEncoding(); got != (semantics.Key{Owner: "typescript", Name: "json-1"}) {
		t.Errorf("JSONEncoding is %v, want typescript/json-1", got)
	}
}

// The entities of every graph below: an interface, its base and a property.
const (
	baseEntity semantics.EntityID = iota + 1
	propertyEntity
)

func expression(operator semantics.Key, data string, operands ...semantics.TermID) semantics.Term {
	term := semantics.Term{Kind: semantics.ProviderExpression, Operator: operator, Encoding: JSONEncoding(), Data: []byte(data)}
	for i, operand := range operands {
		term.Operands = append(term.Operands, semantics.Operand{Locator: []uint32{uint32(i)}, Term: operand})
	}
	return term
}

func modelOf(t *testing.T, terms ...semantics.Term) Model {
	t.Helper()
	graph, err := semantics.New[s.Language](semantics.Data{
		Entities: []semantics.Entity{{Kind: InterfaceKind(), Name: "Base"}, {Kind: Key("property"), Name: "active"}},
		Terms:    terms,
	}, semantics.DefaultLimits(), nil)
	if err != nil {
		t.Fatal(err)
	}
	return Model{&model{graph: graph}}
}

// Model.Extension reads the published form and has no answer for another one.
// The first three terms of every graph are a reference to the base, the type
// of the one own member and the object of the own members.
func TestExtensionReader(t *testing.T) {
	base := semantics.Term{Kind: semantics.Reference, Target: baseEntity}
	boolean := expression(Key("primitive"), `"boolean"`)
	object := expression(ObjectOperator(), `[{"Name":"active","Optional":true,"Readonly":false,"Entity":2}]`, 2)
	const baseTerm, booleanTerm, objectTerm, subject = 1, 2, 3, 4

	m := modelOf(t, base, boolean, object, expression(ExtensionOperator(), `{"Bases":1}`, baseTerm, objectTerm))
	extension, ok := m.Extension(subject)
	if !ok || len(extension.Bases) != 1 || extension.Bases[0] != baseTerm || len(extension.Members) != 1 {
		t.Fatal("one base and one member", extension, ok)
	}
	if member := extension.Members[0]; member.Name != "active" || !member.Optional || member.Readonly || member.Type.ID() != booleanTerm {
		t.Fatal("the member", member)
	}
	m = modelOf(t, base, boolean, object, expression(ExtensionOperator(), `{"Bases":2}`, baseTerm, baseTerm, objectTerm))
	if extension, ok := m.Extension(subject); !ok || len(extension.Bases) != 2 || len(extension.Members) != 1 {
		t.Fatal("two bases", extension, ok)
	}
	// A base with type arguments is an application.
	application := semantics.Term{Kind: semantics.Application, Callee: baseTerm, Arguments: []semantics.TermID{booleanTerm}}
	m = modelOf(t, base, boolean, object, application, expression(ExtensionOperator(), `{"Bases":2}`, baseTerm, subject, objectTerm))
	if extension, ok := m.Extension(subject + 1); !ok || len(extension.Bases) != 2 || extension.Bases[1] != subject {
		t.Fatal("a base that is an application", extension, ok)
	}
	encoded := func(term semantics.Term, encoding semantics.Key) semantics.Term {
		term.Encoding = encoding
		return term
	}

	for name, term := range map[string]semantics.Term{
		"an object":                           expression(ObjectOperator(), `[]`),
		"another operator":                    expression(Key("intersection"), `{"Bases":1}`, baseTerm, objectTerm),
		"no base":                             expression(ExtensionOperator(), `{"Bases":0}`, objectTerm),
		"a negative number of bases":          expression(ExtensionOperator(), `{"Bases":-1}`),
		"data that states no number":          expression(ExtensionOperator(), `null`, baseTerm, objectTerm),
		"data that is no object":              expression(ExtensionOperator(), `[1]`, baseTerm, objectTerm),
		"data that is no JSON":                expression(ExtensionOperator(), `{"Bases":1`, baseTerm, objectTerm),
		"fewer operands than the data says":   expression(ExtensionOperator(), `{"Bases":2}`, baseTerm, objectTerm),
		"more operands than the data says":    expression(ExtensionOperator(), `{"Bases":1}`, baseTerm, baseTerm, objectTerm),
		"an operand after the members object": expression(ExtensionOperator(), `{"Bases":1}`, baseTerm, objectTerm, objectTerm),
		"no members object":                   expression(ExtensionOperator(), `{"Bases":1}`, baseTerm),
		"a last operand that is no object":    expression(ExtensionOperator(), `{"Bases":1}`, baseTerm, booleanTerm),
		"the members object first":            expression(ExtensionOperator(), `{"Bases":1}`, objectTerm, baseTerm),
		"a reference":                         {Kind: semantics.Reference, Target: baseEntity},
		// The forms that decoding alone would accept.
		"another encoding with the same bytes": encoded(expression(ExtensionOperator(), `{"Bases":1}`, baseTerm, objectTerm), Key("utf8")),
		"the member in lower case":             expression(ExtensionOperator(), `{"bases":1}`, baseTerm, objectTerm),
		"the member in upper case":             expression(ExtensionOperator(), `{"BASES":1}`, baseTerm, objectTerm),
		"a further member":                     expression(ExtensionOperator(), `{"Bases":1,"Members":7}`, baseTerm, objectTerm),
		"the member twice":                     expression(ExtensionOperator(), `{"Bases":2,"Bases":1}`, baseTerm, objectTerm),
		"white space around the data":          expression(ExtensionOperator(), " {\n\"Bases\" : 1 }\n", baseTerm, objectTerm),
		"the number with a fraction":           expression(ExtensionOperator(), `{"Bases":1.0}`, baseTerm, objectTerm),
		"the number as a text":                 expression(ExtensionOperator(), `{"Bases":"1"}`, baseTerm, objectTerm),
		"a number beyond an int":               expression(ExtensionOperator(), `{"Bases":92233720368547758070}`, baseTerm, objectTerm),
		"an empty object as data":              expression(ExtensionOperator(), `{}`, baseTerm, objectTerm),
		"empty data":                           expression(ExtensionOperator(), ``, baseTerm, objectTerm),
		"a base that is an object":             expression(ExtensionOperator(), `{"Bases":1}`, objectTerm, objectTerm),
		"a base that is a primitive":           expression(ExtensionOperator(), `{"Bases":1}`, booleanTerm, objectTerm),
	} {
		if extension, ok := modelOf(t, base, boolean, object, term).Extension(subject); ok {
			t.Errorf("%s was read as an extension: %+v", name, extension)
		}
	}
	if _, ok := m.Extension(0); ok {
		t.Error("no term was read as an extension")
	}
	if _, ok := m.Extension(99); ok {
		t.Error("a term outside the graph was read as an extension")
	}
	if _, ok := (Model{}).Extension(subject); ok {
		t.Error("an invalid model read an extension")
	}
}
