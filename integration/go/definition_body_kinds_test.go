package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/types"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	resolvetypescript "github.com/Bitspark/bittype/resolve/files/typescript/provider/go"
	"github.com/Bitspark/bittype/trees/representation/go"

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
	m := goModel(t, "type Celsius float64; type Temp Celsius; type Names []string; type Alias = Temp; type Box[T any] struct { Value T }; type Boxed = Box[Temp]; type Held Box[Temp]; type GA[T any] = []T; type Number interface { ~int | ~float64 }")
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
	// A field has no body. Its type is the declared-type fact of its entity,
	// and the same term is the field's operand of the struct term.
	var fields []golang.FieldInfo
	if err := json.Unmarshal(bodyTerm(t, graph, box).Data, &fields); err != nil || len(fields) != 1 {
		t.Fatal("Box: field metadata", err, fields)
	}
	if field, _ := graph.Entity(fields[0].Entity); field.Kind != golang.Key("field") || field.Name != "Value" || field.Body != 0 {
		t.Fatalf("Box's field has a body: %+v", field)
	}
	if declared(fields[0].Entity).ID() != bodyTerm(t, graph, box).Operands[0].Term {
		t.Fatal("Box's field: the fact does not name the type that the struct's operand names")
	}
	if read, ok := declared(boxID).Fields(); !ok || len(read) != 1 || read[0].Type.ID() != declared(fields[0].Entity).ID() {
		t.Fatal("Box's field: Fields does not read the same type")
	}
	// A use of a generic declaration as a right-hand side: an application.
	if body := bodyTerm(t, graph, boxed); body.Kind != semantics.Application {
		t.Fatalf("Boxed: body %+v", body)
	}
	if bindings, ok := graph.DirectParameterBindings(boxed.Body); !ok || len(bindings) != 1 || bindings[0].Parameter != box.Parameters[0].Entity || bindings[0].Argument != declared(tempID).ID() {
		t.Fatalf("Boxed: bindings %+v %v", bindings, ok)
	}
	// A defined type over an instantiation of a captured generic type: the
	// body is the application, as for the alias. The struct of the instance is
	// the underlying type and stays a fact.
	heldID, held := entityNamed(t, graph, golang.DefinedKind(), "Held")
	if body := bodyTerm(t, graph, held); body.Kind != semantics.Application || held.Body == declared(heldID).Underlying().ID() || !declared(heldID).Underlying().Valid() {
		t.Fatalf("Held: body %+v", body)
	}
	if bindings, ok := graph.DirectParameterBindings(held.Body); !ok || len(bindings) != 1 || bindings[0].Parameter != box.Parameters[0].Entity || bindings[0].Argument != declared(tempID).ID() {
		t.Fatalf("Held: bindings %+v %v", bindings, ok)
	}
	// A generic alias: the body is the term of its right-hand side, and the
	// use of its parameter is a reference to the parameter's entity.
	gaID, ga := entityNamed(t, graph, golang.AliasKind(), "GA")
	if body := bodyTerm(t, graph, ga); body.Operator != golang.SliceOperator() || len(body.Operands) != 1 || len(ga.Parameters) != 1 || body.Operands[0].Term != declared(ga.Parameters[0].Entity).ID() || ga.Body != declared(gaID).AliasTarget().ID() {
		t.Fatalf("GA: body %+v", body)
	}
	// An interface that states a type set: the body is the interface, whose
	// operand is the union of its terms.
	numberID, number := entityNamed(t, graph, golang.DefinedKind(), "Number")
	if body := bodyTerm(t, graph, number); body.Operator != golang.Key("interface") || len(body.Operands) != 1 || number.Body != declared(numberID).Underlying().ID() {
		t.Fatalf("Number: body %+v", body)
	}
	if union, _ := graph.Term(bodyTerm(t, graph, number).Operands[0].Term); union.Operator != golang.Key("union") || len(union.Operands) != 2 {
		t.Fatalf("Number: the operand is no union of two terms: %+v", union)
	}
	// A predeclared type has no body. Its type is reached through the fact.
	stringID, predeclared := entityNamed(t, graph, golang.BuiltinKind(), "string")
	if predeclared.Body != 0 || !declared(stringID).IsBasic(types.String) {
		t.Fatalf("string: %+v", predeclared)
	}
	// Every entity of a type name and of a field carries exactly one
	// declared-type fact, a package carries none, and Lookup reads the fact.
	// Only a captured declaration or alias has a body.
	kinds := map[semantics.Key]int{}
	for i, e := range graph.Data().Entities {
		id := semantics.EntityID(i + 1)
		kinds[e.Kind]++
		if declaration := (e.Kind == golang.DefinedKind() || e.Kind == golang.AliasKind()) && len(e.Declarations) > 0; declaration != (e.Body != 0) {
			t.Fatalf("%s of kind %v: body %d", e.Name, e.Kind, e.Body)
		}
		facts := graph.Facts(semantics.Subject{Entity: id}, golang.DeclaredTypeRelation())
		named := e.Kind == golang.DefinedKind() || e.Kind == golang.AliasKind() || e.Kind == golang.ParameterKind() || e.Kind == golang.BuiltinKind() || e.Kind == golang.Key("field")
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
	// The loop met a field, also the field of the instantiated struct, and a
	// package.
	if kinds[golang.Key("field")] < 2 || kinds[golang.PackageKind()] != 1 {
		t.Fatal("the graph lacks a kind that the loop judges", kinds)
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
export interface Staff extends Person, Box<Person> {}
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
	// second the body of Derived and of Staff is the object of all members.
	roles := map[string]int{}
	for _, node := range trees[0].Nodes {
		for _, child := range node.Children {
			if node.Kind == typescriptsyntax.KindTypeAlias && child.Role == typescriptsyntax.RoleType || node.Kind == typescriptsyntax.KindInterface && child.Role == typescriptsyntax.RoleBase {
				roles[child.Role]++
			}
		}
	}
	if roles[typescriptsyntax.RoleType] != 6 || roles[typescriptsyntax.RoleBase] != 3 {
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
	staff := tsDecl(t, m, "0", "Staff")

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
	if _, ok := m.Extension(person.Entity().Body); ok {
		t.Fatal("Person: an interface without a base was read as an extension")
	}
	// A property has no body. Its type is the declared-type fact of its
	// entity, and the same term is the property's operand of the object term.
	var members []typescript.PropertyInfo
	if err := json.Unmarshal(bodyTerm(t, graph, person.Entity()).Data, &members); err != nil || len(members) != 1 {
		t.Fatal("Person: member metadata", err, members)
	}
	if property, _ := graph.Entity(members[0].Entity); property.Kind != typescript.Key("property") || property.Name != "name" || property.Body != 0 {
		t.Fatalf("Person's property has a body: %+v", property)
	}
	if declared, ok := graph.OnlyFact(semantics.Subject{Entity: members[0].Entity}, typescript.DeclaredTypeRelation()); !ok || declared != bodyTerm(t, graph, person.Entity()).Operands[0].Term {
		t.Fatal("Person's property: the fact does not name the type that the object's operand names")
	}
	if read, ok := person.Type().Properties(); !ok || len(read) != 1 || read[0].Type.ID() != bodyTerm(t, graph, person.Entity()).Operands[0].Term {
		t.Fatal("Person's property: Properties does not read the same type")
	}
	// An interface with an extends clause: an extension. Its operands are the
	// base type terms as written and then the object of its own members. Its
	// data states the number of bases. The declared type still has all
	// members, inherited ones too.
	body := bodyTerm(t, graph, derived.Entity())
	var info typescript.ExtensionInfo
	if body.Operator != typescript.ExtensionOperator() || json.Unmarshal(body.Data, &info) != nil || info.Bases != 1 || len(body.Operands) != 2 {
		t.Fatalf("Derived: body %+v", body)
	}
	extension, ok := m.Extension(derived.Entity().Body)
	if !ok || len(extension.Bases) != 1 || extension.Bases[0] != body.Operands[0].Term || len(extension.Members) != 1 || extension.Members[0].Name != "active" || extension.Members[0].Type.String() != "boolean" {
		t.Fatalf("Derived: extension %+v %v", extension, ok)
	}
	if base, _ := graph.Term(extension.Bases[0]); base.Kind != semantics.Reference || base.Target != person.ID() {
		t.Fatalf("Derived: the base is no reference to Person: %+v", base)
	}
	if own, _ := graph.Term(body.Operands[1].Term); own.Operator != typescript.ObjectOperator() || len(own.Operands) != 1 || body.Operands[1].Term == properties(derived) {
		t.Fatalf("Derived: the last operand is not the object of its own members: %+v", own)
	}
	if all, ok := derived.Type().Properties(); !ok || len(all) != 2 {
		t.Fatalf("Derived: members %+v", all)
	}
	// Two bases, the second with type arguments, and no member of its own:
	// the bases in written order, an application as elsewhere, and an object
	// without operands.
	body = bodyTerm(t, graph, staff.Entity())
	extension, ok = m.Extension(staff.Entity().Body)
	if !ok || body.Operator != typescript.ExtensionOperator() || len(body.Operands) != 3 || len(extension.Bases) != 2 || len(extension.Members) != 0 {
		t.Fatalf("Staff: extension %+v %v", extension, ok)
	}
	if first, _ := graph.Term(extension.Bases[0]); first.Kind != semantics.Reference || first.Target != person.ID() {
		t.Fatalf("Staff: the first base is no reference to Person: %+v", first)
	}
	if bindings, ok := graph.DirectParameterBindings(extension.Bases[1]); !ok || len(bindings) != 1 || bindings[0].Parameter != box.Parameters()[0].ID() || bindings[0].Argument != person.Type().ID() {
		t.Fatalf("Staff: the second base is no application of Box to Person: %+v %v", bindings, ok)
	}
	if all, ok := staff.Type().Properties(); !ok || len(all) != 2 {
		t.Fatalf("Staff: members %+v", all)
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
	for _, d := range []typescript.Declaration{text, alias, flag, person, derived, box, staff, boxed, choice, items} {
		if body := bodyTerm(t, graph, d.Entity()); body.Kind == semantics.Reference && body.Target == d.ID() {
			t.Fatalf("%s: the body is a reference to the declaration itself", d.Entity().Name)
		}
	}
	// Only an alias and an interface have a body. An alias, an interface, a
	// type parameter and a property carry exactly one declared-type fact, and
	// no other entity carries one.
	kinds := map[semantics.Key]int{}
	for i, e := range graph.Data().Entities {
		kinds[e.Kind]++
		if declaration := e.Kind == typescript.AliasKind() || e.Kind == typescript.InterfaceKind(); declaration != (e.Body != 0) {
			t.Fatalf("%s of kind %v: body %d", e.Name, e.Kind, e.Body)
		}
		typed := e.Kind == typescript.AliasKind() || e.Kind == typescript.InterfaceKind() || e.Kind == typescript.ParameterKind() || e.Kind == typescript.Key("property")
		if facts := graph.Facts(semantics.Subject{Entity: semantics.EntityID(i + 1)}, typescript.DeclaredTypeRelation()); typed != (len(facts) == 1) || !typed && len(facts) != 0 {
			t.Fatalf("%s of kind %v carries %d declared-type facts", e.Name, e.Kind, len(facts))
		}
	}
	if kinds[typescript.Key("property")] < 4 || kinds[typescript.ModuleKind()] != 1 || kinds[typescript.ParameterKind()] != 1 {
		t.Fatal("the graph lacks a kind that the loop judges", kinds)
	}
}

// A chain of declarations nests no frame for each link: the bridge sets the
// bodies in a pass of their own, after the declarations have their entities.
// Both sources are analysed on main, with the same number of entities.
func TestGoChainOfDeclarationsNestsNoFramePerLink(t *testing.T) {
	// 200 defined types, each over the next.
	var chain []string
	for i := 1; i < 200; i++ {
		chain = append(chain, fmt.Sprintf("type C%d C%d", i, i+1))
	}
	chain = append(chain, "type C200 string")
	graph := goModel(t, strings.Join(chain, "; ")).Graph()
	if entities := len(graph.Data().Entities); entities != 202 {
		t.Fatal("the chain: entities", entities)
	}
	for i := 1; i < 200; i++ {
		_, link := entityNamed(t, graph, golang.DefinedKind(), fmt.Sprintf("C%d", i))
		next, _ := entityNamed(t, graph, golang.DefinedKind(), fmt.Sprintf("C%d", i+1))
		if body := bodyTerm(t, graph, link); body.Kind != semantics.Reference || body.Target != next {
			t.Fatalf("C%d: body %+v", i, body)
		}
	}
	if _, last := entityNamed(t, graph, golang.DefinedKind(), "C200"); bodyTerm(t, graph, last).Operator != golang.BasicOperator() {
		t.Fatal("C200: body", bodyTerm(t, graph, last))
	}
	// 65 structs whose one field is a defined type over the next struct.
	var pairs []string
	for i := 1; i < 65; i++ {
		pairs = append(pairs, fmt.Sprintf("type A%d struct { F B%d }; type B%d A%d", i, i, i, i+1))
	}
	pairs = append(pairs, "type A65 struct { F string }")
	graph = goModel(t, strings.Join(pairs, "; ")).Graph()
	if entities := len(graph.Data().Entities); entities != 196 {
		t.Fatal("the pairs: entities", entities)
	}
	for i := 1; i < 65; i++ {
		_, a := entityNamed(t, graph, golang.DefinedKind(), fmt.Sprintf("A%d", i))
		b, definition := entityNamed(t, graph, golang.DefinedKind(), fmt.Sprintf("B%d", i))
		next, _ := entityNamed(t, graph, golang.DefinedKind(), fmt.Sprintf("A%d", i+1))
		body := bodyTerm(t, graph, a)
		if body.Operator != golang.StructOperator() || len(body.Operands) != 1 {
			t.Fatalf("A%d: body %+v", i, body)
		}
		if field, _ := graph.Term(body.Operands[0].Term); field.Kind != semantics.Reference || field.Target != b {
			t.Fatalf("A%d: its field is no reference to B%d: %+v", i, i, field)
		}
		if body := bodyTerm(t, graph, definition); body.Kind != semantics.Reference || body.Target != next {
			t.Fatalf("B%d: body %+v", i, body)
		}
	}
}

// Where the compiler reduces the right-hand side of an alias to another named
// type, the body is a use of that type and not its structure. Inside a body
// the operands are the compiler's types: an alias of a primitive or of an
// interface is not kept there. A top-level reference, the callee of an
// application and the bases of an extension are as written.
func TestTypeScriptBodyOfReducedAndErasedForms(t *testing.T) {
	m := tsModel(t, `export type A = { x: string };
export interface I { y: string }
export interface Box<T> { value: T }
export type M2 = A & unknown;
export type M4 = A | never;
export type M5 = A & {};
export type N = I & unknown;
export type RB = Box<string> & unknown;
export type Text = string;
export type T2 = Text;
export type Who = I;
export type P = Box<Text>;
export type Q = { name: Text };
export interface H { who: Who; t: T2 }
export interface Parent { name: string; id: string | number }
export interface Child extends Parent { id: string }`)
	graph := m.Graph()
	d := func(name string) typescript.Declaration { t.Helper(); return tsDecl(t, m, "0", name) }
	a, i, box, text := d("A"), d("I"), d("Box"), d("Text")
	// Reduced to an alias of a type literal: a reference to A. It is the term
	// of A's declared type, and not A's structure, which is A's own body.
	for _, name := range []string{"M2", "M4", "M5"} {
		reduced := d(name)
		if body := bodyTerm(t, graph, reduced.Entity()); body.Kind != semantics.Reference || body.Target != a.ID() || reduced.Entity().Body != a.Type().ID() || reduced.Entity().Body == a.Entity().Body {
			t.Fatalf("%s: body %+v", name, body)
		}
		if !reduced.Type().Same(a.Type()) {
			t.Fatalf("%s: the declared type is not the type of A", name)
		}
	}
	if body := bodyTerm(t, graph, a.Entity()); body.Operator != typescript.ObjectOperator() {
		t.Fatalf("A: body %+v", body)
	}
	// Reduced to an interface: a reference to it.
	if n := d("N"); bodyTerm(t, graph, n.Entity()).Kind != semantics.Reference || bodyTerm(t, graph, n.Entity()).Target != i.ID() || n.Entity().Body == i.Entity().Body {
		t.Fatalf("N: body %+v", bodyTerm(t, graph, n.Entity()))
	}
	// Reduced to an instantiation: its application.
	rb := d("RB")
	if body := bodyTerm(t, graph, rb.Entity()); body.Kind != semantics.Application || rb.Entity().Body != rb.Type().ID() {
		t.Fatalf("RB: body %+v", body)
	}
	if bindings, ok := graph.DirectParameterBindings(rb.Entity().Body); !ok || len(bindings) != 1 || bindings[0].Parameter != box.Parameters()[0].ID() || bindings[0].Argument != text.Type().ID() {
		t.Fatalf("RB: bindings %+v %v", bindings, ok)
	}
	// A top-level reference is as written: T2 is a reference to Text, and
	// Who a reference to I.
	if body := bodyTerm(t, graph, d("T2").Entity()); body.Kind != semantics.Reference || body.Target != text.ID() {
		t.Fatalf("T2: body %+v", body)
	}
	// The callee of an application is as written, and its argument is the
	// compiler's type: the alias Text is not kept, the argument is string.
	p := d("P")
	application := bodyTerm(t, graph, p.Entity())
	if callee, _ := graph.Term(application.Callee); application.Kind != semantics.Application || callee.Kind != semantics.Reference || callee.Target != box.ID() {
		t.Fatalf("P: body %+v", application)
	}
	if len(application.Arguments) != 1 || application.Arguments[0] != text.Type().ID() || text.Type().Term().Operator != typescript.Key("primitive") {
		t.Fatalf("P: the argument is not the primitive type: %+v", application)
	}
	// A member's type is the compiler's type: string for Text and for T2, and
	// a reference to I for Who.
	if body := bodyTerm(t, graph, d("Q").Entity()); body.Operator != typescript.ObjectOperator() || len(body.Operands) != 1 || body.Operands[0].Term != text.Type().ID() {
		t.Fatalf("Q: body %+v", body)
	}
	h := bodyTerm(t, graph, d("H").Entity())
	if h.Operator != typescript.ObjectOperator() || len(h.Operands) != 2 || h.Operands[1].Term != text.Type().ID() {
		t.Fatalf("H: body %+v", h)
	}
	if who, _ := graph.Term(h.Operands[0].Term); who.Kind != semantics.Reference || who.Target != i.ID() {
		t.Fatalf("H: the member who is no reference to I: %+v", who)
	}
	// A member that declares an inherited name again is the interface's own:
	// it is in the members object of the extension, with its own type.
	child := d("Child")
	extension, ok := m.Extension(child.Entity().Body)
	if !ok || len(extension.Bases) != 1 || len(extension.Members) != 1 || extension.Members[0].Name != "id" || extension.Members[0].Type.String() != "string" {
		t.Fatalf("Child: extension %+v %v", extension, ok)
	}
	if all, ok := child.Type().Properties(); !ok || len(all) != 2 {
		t.Fatalf("Child: members %+v", all)
	}
	// Model.DeclaredType reads the fact by entity: for a declaration the type
	// that Declaration.Type returns, for a property the type of its operand,
	// and nothing for an entity without the fact.
	if ty, err := m.DeclaredType(a.ID()); err != nil || !ty.Same(a.Type()) {
		t.Fatal("DeclaredType of A", err)
	}
	var members []typescript.PropertyInfo
	if err := json.Unmarshal(bodyTerm(t, graph, a.Entity()).Data, &members); err != nil || len(members) != 1 {
		t.Fatal("A: member metadata", err)
	}
	if ty, err := m.DeclaredType(members[0].Entity); err != nil || ty.ID() != bodyTerm(t, graph, a.Entity()).Operands[0].Term || ty.String() != "string" {
		t.Fatal("DeclaredType of the property x", err)
	}
	if _, err := m.Declaration(members[0].Entity); !errors.Is(err, typescript.ErrForeign) {
		t.Fatal("a property was opened as a declaration", err)
	}
	module, _ := entityNamed(t, graph, typescript.ModuleKind(), "0.ts")
	if _, err := m.DeclaredType(module); !errors.Is(err, typescript.ErrForeign) {
		t.Fatal("DeclaredType of the module", err)
	}
	if _, err := m.DeclaredType(0); !errors.Is(err, typescript.ErrForeign) {
		t.Fatal("DeclaredType of no entity", err)
	}
}

// The same rule holds for a supplied package type: a right-hand side that is
// reduced to it is a reference, or an application for an instantiation. An
// interface whose base is a supplied type without members has an extension
// whose base is a reference to the package type. One whose base has members
// is refused, as its declared type was before: the profile publishes no
// structure of a supplied type.
func TestTypeScriptBodyOverSuppliedTypes(t *testing.T) {
	settings, err := typescript.ImportConfiguration([]typescript.PackageTypes{{Module: "@bitspark/pkg", Declaration: "export type Obj = { a: string };\nexport interface Cell<T> { value: T }\nexport interface Empty {}\nexport interface Base { a: string }\n"}})
	if err != nil {
		t.Fatal(err)
	}
	resolve := func(text string) semantic.Result {
		t.Helper()
		snapshot, _, _ := tsCapture(t, text)
		resolver, err := resolvetypescript.New(snapshot.Language(), tsCompiler(t))
		if err != nil {
			t.Fatal(err)
		}
		resolved, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{Configuration: settings})
		if err != nil {
			t.Fatal(err)
		}
		return resolved
	}
	resolved := resolve(`import type { Obj, Cell, Empty } from "@bitspark/pkg";
export type Mine = Obj;
export type Unk = Obj & unknown;
export type Held = Cell<string> & unknown;
export interface D extends Empty { own: string }`)
	envelope, ok := resolved.Semantics()
	if !ok {
		t.Fatal("no graph", resolved.Diagnostics())
	}
	graph, err := representation.UnpackSemantics[typescriptsyntax.Language](resolved.Snapshot().Language(), envelope)
	if err != nil {
		t.Fatal(err)
	}
	obj, _ := entityNamed(t, graph, typescript.PackageTypeKind(), "Obj")
	cell, _ := entityNamed(t, graph, typescript.PackageTypeKind(), "Cell")
	empty, supplied := entityNamed(t, graph, typescript.PackageTypeKind(), "Empty")
	if supplied.Body != 0 || len(graph.Facts(semantics.Subject{Entity: empty}, typescript.DeclaredTypeRelation())) != 0 {
		t.Fatalf("the supplied type has a body or a declared type: %+v", supplied)
	}
	for _, name := range []string{"Mine", "Unk"} {
		_, alias := entityNamed(t, graph, typescript.AliasKind(), name)
		if body := bodyTerm(t, graph, alias); body.Kind != semantics.Reference || body.Target != obj {
			t.Fatalf("%s: body %+v", name, body)
		}
	}
	_, held := entityNamed(t, graph, typescript.AliasKind(), "Held")
	if callee, _ := graph.Term(bodyTerm(t, graph, held).Callee); bodyTerm(t, graph, held).Kind != semantics.Application || callee.Target != cell || len(bodyTerm(t, graph, held).Arguments) != 1 {
		t.Fatalf("Held: body %+v", bodyTerm(t, graph, held))
	}
	_, d := entityNamed(t, graph, typescript.InterfaceKind(), "D")
	body := bodyTerm(t, graph, d)
	if body.Operator != typescript.ExtensionOperator() || len(body.Operands) != 2 {
		t.Fatalf("D: body %+v", body)
	}
	if base, _ := graph.Term(body.Operands[0].Term); base.Kind != semantics.Reference || base.Target != empty {
		t.Fatalf("D: the base is no reference to the supplied type: %+v", base)
	}
	if own, _ := graph.Term(body.Operands[1].Term); own.Operator != typescript.ObjectOperator() || len(own.Operands) != 1 {
		t.Fatalf("D: the own members: %+v", own)
	}
	// A supplied base with members: unsupported coverage and no graph.
	refused := resolve("import type { Base } from \"@bitspark/pkg\";\nexport interface D extends Base { own: string }")
	if report, _ := refused.Report(typescript.Coverage()); refused.HasSemantics() || report.Status != semantic.Unsupported {
		t.Fatal("an interface over a supplied base with members", report)
	}
}
