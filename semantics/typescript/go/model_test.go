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
}

// The entities of every graph below: an interface, its base and a property.
const (
	baseEntity semantics.EntityID = iota + 1
	propertyEntity
)

func expression(operator semantics.Key, data string, operands ...semantics.TermID) semantics.Term {
	term := semantics.Term{Kind: semantics.ProviderExpression, Operator: operator, Encoding: Key("json-1"), Data: []byte(data)}
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

	for name, term := range map[string]semantics.Term{
		"an object":                         expression(ObjectOperator(), `[]`),
		"another operator":                  expression(Key("intersection"), `{"Bases":1}`, baseTerm, objectTerm),
		"no base":                           expression(ExtensionOperator(), `{"Bases":0}`, objectTerm),
		"a negative number of bases":        expression(ExtensionOperator(), `{"Bases":-1}`),
		"data that states no number":        expression(ExtensionOperator(), `null`, baseTerm, objectTerm),
		"data that is no object":            expression(ExtensionOperator(), `[1]`, baseTerm, objectTerm),
		"data that is no JSON":              expression(ExtensionOperator(), `{"Bases":1`, baseTerm, objectTerm),
		"fewer operands than the data says": expression(ExtensionOperator(), `{"Bases":2}`, baseTerm, objectTerm),
		"more operands than the data says":  expression(ExtensionOperator(), `{"Bases":1}`, baseTerm, baseTerm, objectTerm),
		"no members object":                 expression(ExtensionOperator(), `{"Bases":1}`, baseTerm),
		"a last operand that is no object":  expression(ExtensionOperator(), `{"Bases":1}`, baseTerm, booleanTerm),
		"the members object first":          expression(ExtensionOperator(), `{"Bases":1}`, objectTerm, baseTerm),
		"a reference":                       {Kind: semantics.Reference, Target: baseEntity},
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
