package integration_test

import (
	"errors"
	"go/types"
	"testing"

	"github.com/Bitspark/bittype/semantics/go"
	golang "github.com/Bitspark/bittype/semantics/golang/go"
	typescript "github.com/Bitspark/bittype/semantics/typescript/go"
	typescriptsyntax "github.com/Bitspark/bittype/syntax/typescript/go"
)

// The tests in this file state, for each kind of declaration, what the body of
// its entity is and where the type of its name is found (bittype#310). They
// read the two vocabularies, which the tests in definition_body_test.go do
// not.

func entityNamed[L any](t *testing.T, graph semantics.Semantics[L], kind semantics.Key, name string) (semantics.EntityID, semantics.Entity) {
	t.Helper()
	var found []semantics.EntityID
	data := graph.Data()
	for i, e := range data.Entities {
		if e.Kind == kind && e.Name == name {
			found = append(found, semantics.EntityID(i+1))
		}
	}
	if len(found) != 1 {
		t.Fatalf("entities of kind %v named %s: %v", kind, name, found)
	}
	return found[0], data.Entities[found[0]-1]
}

func bodyTerm[L any](t *testing.T, graph semantics.Semantics[L], e semantics.Entity) semantics.Term {
	t.Helper()
	term, ok := graph.Term(e.Body)
	if !ok {
		t.Fatalf("%s has no body", e.Name)
	}
	return term
}

func TestGoBodyOfEachDeclarationKind(t *testing.T) {
	m := goModel(t, "type Celsius float64; type Temp Celsius; type Names []string; type Alias = Temp; type Box[T any] struct { Value T }; type Boxed = Box[Temp]")
	graph := m.Graph()
	declared := func(id semantics.EntityID) golang.Type {
		t.Helper()
		ty, err := m.DeclaredType(id)
		if err != nil {
			t.Fatalf("entity %d has no declared type: %v", id, err)
		}
		return ty
	}
	celsiusID, celsius := entityNamed(t, graph, golang.DefinedKind(), "Celsius")
	tempID, temp := entityNamed(t, graph, golang.DefinedKind(), "Temp")
	_, names := entityNamed(t, graph, golang.DefinedKind(), "Names")
	aliasID, alias := entityNamed(t, graph, golang.AliasKind(), "Alias")
	boxID, box := entityNamed(t, graph, golang.DefinedKind(), "Box")
	_, boxed := entityNamed(t, graph, golang.AliasKind(), "Boxed")

	// A defined type whose right-hand side is a type literal or a predeclared
	// type: the body is that type's term, which is also its underlying type.
	if body := bodyTerm(t, graph, celsius); body.Operator != golang.BasicOperator() || celsius.Body != declared(celsiusID).Underlying().ID() {
		t.Fatalf("Celsius: body %+v", body)
	}
	if body := bodyTerm(t, graph, names); body.Operator != golang.SliceOperator() || len(body.Operands) != 1 {
		t.Fatalf("Names: body %+v", body)
	}
	// A defined type whose right-hand side is another defined type: the body
	// is a reference to that declaration, and the underlying type stays a fact.
	if body := bodyTerm(t, graph, temp); body.Kind != semantics.Reference || body.Target != celsiusID || temp.Body != declared(celsiusID).ID() {
		t.Fatalf("Temp: body %+v", body)
	}
	if !declared(tempID).Underlying().IsBasic(types.Float64) || declared(tempID).ID() == temp.Body {
		t.Fatal("Temp: the underlying type is no longer the fact of the declared type")
	}
	if own, _ := graph.Term(declared(tempID).ID()); own.Kind != semantics.Reference || own.Target != tempID {
		t.Fatalf("Temp: the declared type is no reference to Temp: %+v", own)
	}
	// An alias: the body is the term of the given type, which is also its
	// alias target.
	if body := bodyTerm(t, graph, alias); body.Kind != semantics.Reference || body.Target != tempID || alias.Body != declared(aliasID).AliasTarget().ID() {
		t.Fatalf("Alias: body %+v", body)
	}
	// A generic declaration: the body is its struct, and its parameter has no
	// body. The parameter's type and constraint are reached through the fact.
	if body := bodyTerm(t, graph, box); body.Operator != golang.StructOperator() || len(body.Operands) != 1 || len(box.Parameters) != 1 {
		t.Fatalf("Box: body %+v", body)
	}
	parameter, _ := graph.Entity(box.Parameters[0].Entity)
	if parameter.Kind != golang.ParameterKind() || parameter.Body != 0 {
		t.Fatalf("Box's parameter has a body: %+v", parameter)
	}
	if p := declared(box.Parameters[0].Entity); !p.Unconstrained() || len(declared(boxID).Parameters()) != 1 || declared(boxID).Parameters()[0].ID() != p.ID() || bodyTerm(t, graph, box).Operands[0].Term != p.ID() {
		t.Fatal("Box's parameter: the fact does not lead to its type and constraint")
	}
	// A use of a generic declaration as a right-hand side: an application.
	if body := bodyTerm(t, graph, boxed); body.Kind != semantics.Application {
		t.Fatalf("Boxed: body %+v", body)
	}
	if bindings, ok := graph.DirectParameterBindings(boxed.Body); !ok || len(bindings) != 1 || bindings[0].Parameter != box.Parameters[0].Entity || bindings[0].Argument != declared(tempID).ID() {
		t.Fatalf("Boxed: bindings %+v %v", bindings, ok)
	}
	// A predeclared type has no body. Its type is reached through the fact.
	stringID, predeclared := entityNamed(t, graph, golang.BuiltinKind(), "string")
	if predeclared.Body != 0 || !declared(stringID).IsBasic(types.String) {
		t.Fatalf("string: %+v", predeclared)
	}
	// Every entity of a type name carries exactly one declared-type fact, a
	// package and a field carry none, and Lookup reads the fact.
	for i, e := range graph.Data().Entities {
		id := semantics.EntityID(i + 1)
		facts := graph.Facts(semantics.Subject{Entity: id}, golang.DeclaredTypeRelation())
		named := e.Kind == golang.DefinedKind() || e.Kind == golang.AliasKind() || e.Kind == golang.ParameterKind() || e.Kind == golang.BuiltinKind()
		if named != (len(facts) == 1) || !named && len(facts) != 0 {
			t.Fatalf("%s of kind %v carries %d declared-type facts", e.Name, e.Kind, len(facts))
		}
		if _, err := m.DeclaredType(id); named != (err == nil) || !named && !errors.Is(err, golang.ErrForeign) {
			t.Fatalf("%s of kind %v: DeclaredType %v", e.Name, e.Kind, err)
		}
		if (e.Kind == golang.DefinedKind() || e.Kind == golang.AliasKind()) && len(e.Declarations) > 0 {
			if found, err := m.Lookup(e.Name); err != nil || found.ID() != facts[0] {
				t.Fatalf("Lookup(%s) does not read the declared-type fact: %v", e.Name, err)
			}
		}
	}
}

func TestGoSuppliedDeclarationHasNoBody(t *testing.T) {
	model, _ := goImportModel(t, map[string]string{
		"pages/page.go": "package pages\n\nimport \"example.org/cell\"\n\ntype Local cell.Cell[string]\n",
	}, goImportSettings(t, cellPackage))
	graph := model.Graph()
	cellID, cell := entityNamed(t, graph, golang.DefinedKind(), "Cell")
	if cell.Body != 0 || len(cell.Declarations) != 0 || len(cell.Parameters) != 1 {
		t.Fatalf("the supplied declaration has a body or a site: %+v", cell)
	}
	if parameter, _ := graph.Entity(cell.Parameters[0].Entity); parameter.Body != 0 {
		t.Fatalf("the supplied parameter has a body: %+v", parameter)
	}
	declared, err := model.DeclaredType(cellID)
	if err != nil || declared.Underlying().Valid() {
		t.Fatal("the supplied declaration's type is not reached through the fact, or its structure is published", err)
	}
	if found, err := model.LookupIn("example.org/cell", "Cell"); err != nil || found.ID() != declared.ID() {
		t.Fatal("LookupIn does not read the declared-type fact of a supplied declaration", err)
	}
	if parameters := declared.Parameters(); len(parameters) != 1 || !parameters[0].Valid() {
		t.Fatal("the supplied declaration's parameter is not reached through the fact")
	}
	// The captured declaration that is defined by the supplied one has the
	// application as its body: a use of the supplied declaration, no structure.
	_, local := entityNamed(t, graph, golang.DefinedKind(), "Local")
	if body := bodyTerm(t, graph, local); body.Kind != semantics.Application {
		t.Fatalf("Local: body %+v", body)
	}
	if bindings, ok := graph.DirectParameterBindings(local.Body); !ok || len(bindings) != 1 || bindings[0].Parameter != cell.Parameters[0].Entity {
		t.Fatalf("Local: bindings %+v %v", bindings, ok)
	}
}

func TestTypeScriptBodyOfEachDeclarationKind(t *testing.T) {
	snapshot, trees, _ := tsCapture(t, `export type Text = string;
export type Alias = Text;
export type Flag = { on: boolean };
export interface Person { name: string }
export interface Derived extends Person { active: boolean }
export interface Box<T extends Person> { value: T }
export type Boxed = Box<Person>;
export type Choice = "a" | "b";
export type Items = Person[];`)
	_, result := tsResults(t, snapshot)
	m, err := typescript.Open(result)
	if err != nil {
		t.Fatal(err, result.Diagnostics())
	}
	// The compiler helper is JavaScript and reads two child roles by spelling
	// to find a definition: the right-hand side of an alias and the base of an
	// interface. syntax/typescript/go declares them. With another spelling of
	// the first no alias below has a body, and with another spelling of the
	// second Derived has one.
	roles := map[string]int{}
	for _, node := range trees[0].Nodes {
		for _, child := range node.Children {
			if node.Kind == typescriptsyntax.KindTypeAlias && child.Role == typescriptsyntax.RoleType || node.Kind == typescriptsyntax.KindInterface && child.Role == typescriptsyntax.RoleBase {
				roles[child.Role]++
			}
		}
	}
	if roles[typescriptsyntax.RoleType] != 6 || roles[typescriptsyntax.RoleBase] != 1 {
		t.Fatal("the capture lacks a right-hand side or a base", roles)
	}
	graph := m.Graph()
	properties := func(d typescript.Declaration) semantics.TermID {
		t.Helper()
		term, ok := graph.OnlyFact(semantics.Subject{Term: d.Type().ID()}, typescript.PropertiesRelation())
		if !ok {
			t.Fatalf("%s: the declared type has no properties fact", d.Entity().Name)
		}
		return term
	}
	text, alias, flag := tsDecl(t, m, "0", "Text"), tsDecl(t, m, "0", "Alias"), tsDecl(t, m, "0", "Flag")
	person, derived, box := tsDecl(t, m, "0", "Person"), tsDecl(t, m, "0", "Derived"), tsDecl(t, m, "0", "Box")
	boxed, choice, items := tsDecl(t, m, "0", "Boxed"), tsDecl(t, m, "0", "Choice"), tsDecl(t, m, "0", "Items")

	// An alias of a keyword type: the body is the primitive's term.
	if body := bodyTerm(t, graph, text.Entity()); body.Operator != typescript.Key("primitive") || text.Entity().Body != text.Type().ID() {
		t.Fatalf("Text: body %+v", body)
	}
	// An alias of a name: the body is a reference to that declaration. The
	// declared type, which the compiler resolved, stays the fact.
	if body := bodyTerm(t, graph, alias.Entity()); body.Kind != semantics.Reference || body.Target != text.ID() || !alias.Type().Same(text.Type()) {
		t.Fatalf("Alias: body %+v", body)
	}
	// An alias of a type literal: the body is the object, whose operands are
	// the member types. The declared type is a reference to the alias, and
	// its properties fact relates the same object.
	if body := bodyTerm(t, graph, flag.Entity()); body.Operator != typescript.ObjectOperator() || len(body.Operands) != 1 || flag.Entity().Body != properties(flag) {
		t.Fatalf("Flag: body %+v", body)
	}
	if own := flag.Type().Term(); own.Kind != semantics.Reference || own.Target != flag.ID() {
		t.Fatalf("Flag: the declared type is no reference to Flag: %+v", own)
	}
	// An interface without an extends clause: the object of its members.
	if body := bodyTerm(t, graph, person.Entity()); body.Operator != typescript.ObjectOperator() || len(body.Operands) != 1 || person.Entity().Body != properties(person) {
		t.Fatalf("Person: body %+v", body)
	}
	// An interface with an extends clause has no body until bittype#310
	// decides how its base types are operands of the definition. Its type and
	// all its members, inherited ones too, are reached as before.
	if derived.Entity().Body != 0 {
		t.Fatalf("Derived has a body: %+v", bodyTerm(t, graph, derived.Entity()))
	}
	if members, ok := derived.Type().Properties(); !ok || len(members) != 2 {
		t.Fatalf("Derived: members %+v", members)
	}
	// A generic declaration: the body is its object, and its parameter has no
	// body. The parameter's type and constraint are entity facts.
	parameter := box.Parameters()[0]
	if body := bodyTerm(t, graph, box.Entity()); body.Operator != typescript.ObjectOperator() || len(body.Operands) != 1 || body.Operands[0].Term != parameter.Type().ID() {
		t.Fatalf("Box: body %+v", body)
	}
	if parameter.Entity().Body != 0 || !parameter.Type().Valid() || !parameter.Constraint().Valid() {
		t.Fatalf("Box's parameter: %+v", parameter.Entity())
	}
	// An alias of a generic type reference: an application.
	if body := bodyTerm(t, graph, boxed.Entity()); body.Kind != semantics.Application {
		t.Fatalf("Boxed: body %+v", body)
	}
	if bindings, ok := graph.DirectParameterBindings(boxed.Entity().Body); !ok || len(bindings) != 1 || bindings[0].Parameter != parameter.ID() || bindings[0].Argument != person.Type().ID() {
		t.Fatalf("Boxed: bindings %+v %v", bindings, ok)
	}
	// An alias of a union and of an array: the structural expression.
	if body := bodyTerm(t, graph, choice.Entity()); body.Operator != typescript.Key("union") || len(body.Operands) != 2 {
		t.Fatalf("Choice: body %+v", body)
	}
	if body := bodyTerm(t, graph, items.Entity()); body.Operator != typescript.Key("array") || len(body.Operands) != 1 || body.Operands[0].Term != person.Type().ID() {
		t.Fatalf("Items: body %+v", body)
	}
	// No declaration's body is a reference to the declaration itself.
	for _, d := range []typescript.Declaration{text, alias, flag, person, box, boxed, choice, items} {
		if body := bodyTerm(t, graph, d.Entity()); body.Kind == semantics.Reference && body.Target == d.ID() {
			t.Fatalf("%s: the body is a reference to the declaration itself", d.Entity().Name)
		}
	}
}
