package integration_test

import (
	"fmt"
	"testing"

	"github.com/Bitspark/bittype/semantics/go"
	golangsyntax "github.com/Bitspark/bittype/syntax/golang/go"
	typescriptsyntax "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// The tests in this file read a published graph as the generation kernel's
// walk does (contracts/generation/go/occurrence.go): a declaration's body and
// its binders, the variant of a term, the target of a reference, the direct
// parameter bindings of an application and the operands of a provider
// expression. The reader below reads no fact, no entity kind, no operator, no
// payload and no locator. A name is read once, to find a declaration to start
// at. The reader is one function for every language (bittype#310, done-when 5).
//
// The contract package exports no way to run its walk over a graph without a
// source provider and a planning request, so the reader is written here and
// follows the walk's own steps.

// structuralUse is what the walk sees of one occurrence.
type structuralUse struct {
	Term semantics.TermID
	// Exactly one of the following three readings holds.
	Bound    semantics.EntityID // a reference to a binder of the declaration being read
	Callee   semantics.EntityID // a reference to, or a direct application of, a declaration
	Operands []structuralUse    // a provider expression
	// Actuals pairs the callee's binders with the occurrences that fill them.
	Actuals    []structuralActual
	Expression bool
}

type structuralActual struct {
	Binder semantics.EntityID
	Use    structuralUse
}

// structuralDeclaration is what the walk sees of one declaration.
type structuralDeclaration struct {
	Entity  semantics.EntityID
	Formals []semantics.EntityID
	Body    structuralUse
}

func structuralBinders[L any](graph semantics.Semantics[L]) map[semantics.EntityID]semantics.EntityID {
	binders := map[semantics.EntityID]semantics.EntityID{}
	for i, e := range graph.Data().Entities {
		for _, p := range e.Parameters {
			binders[p.Entity] = semantics.EntityID(i + 1)
		}
	}
	return binders
}

// readStructurally opens a declaration as the walk does and reads its whole
// body. It fails where the walk fails: a declaration without a body, a
// reference to a binder of another declaration, a generic declaration used
// without arguments and an application without direct parameter bindings.
func readStructurally[L any](graph semantics.Semantics[L], declaration semantics.EntityID) (structuralDeclaration, error) {
	e, ok := graph.Entity(declaration)
	if !ok {
		return structuralDeclaration{}, fmt.Errorf("entity %d is not in the graph", declaration)
	}
	if e.Body == 0 {
		return structuralDeclaration{}, fmt.Errorf("declaration %d has no body", declaration)
	}
	out := structuralDeclaration{Entity: declaration}
	for _, p := range e.Parameters {
		out.Formals = append(out.Formals, p.Entity)
	}
	body, err := readUse(graph, structuralBinders(graph), declaration, e.Body, 0)
	out.Body = body
	return out, err
}

func readUse[L any](graph semantics.Semantics[L], binders map[semantics.EntityID]semantics.EntityID, declaration semantics.EntityID, id semantics.TermID, depth int) (structuralUse, error) {
	use := structuralUse{Term: id}
	if depth > 64 {
		return use, fmt.Errorf("term %d is deeper than 64", id)
	}
	t, ok := graph.Term(id)
	if !ok {
		return use, fmt.Errorf("term %d is not in the graph", id)
	}
	switch t.Kind {
	case semantics.Reference:
		if owner, binder := binders[t.Target]; binder {
			if owner != declaration {
				return use, fmt.Errorf("term %d refers to a binder of another declaration", id)
			}
			use.Bound = t.Target
			return use, nil
		}
		callee, ok := graph.Entity(t.Target)
		if !ok {
			return use, fmt.Errorf("term %d refers to no entity", id)
		}
		if len(callee.Parameters) != 0 {
			return use, fmt.Errorf("term %d uses a generic declaration without arguments", id)
		}
		use.Callee = t.Target
	case semantics.Application:
		bindings, ok := graph.DirectParameterBindings(id)
		if !ok {
			return use, fmt.Errorf("term %d is an application without direct parameter bindings", id)
		}
		callee, _ := graph.Term(t.Callee)
		use.Callee = callee.Target
		for _, b := range bindings {
			actual, err := readUse(graph, binders, declaration, b.Argument, depth+1)
			if err != nil {
				return use, err
			}
			use.Actuals = append(use.Actuals, structuralActual{Binder: b.Parameter, Use: actual})
		}
	case semantics.ProviderExpression:
		use.Expression = true
		for _, operand := range t.Operands {
			child, err := readUse(graph, binders, declaration, operand.Term, depth+1)
			if err != nil {
				return use, err
			}
			use.Operands = append(use.Operands, child)
		}
	default:
		return use, fmt.Errorf("term %d has no variant", id)
	}
	return use, nil
}

// declarationNamed finds the one entity of that name that has a declaration
// site. This is the only place where a name is read. The fixtures below give
// no field and no parameter the name of a declaration.
func declarationNamed[L any](t *testing.T, graph semantics.Semantics[L], name string) semantics.EntityID {
	t.Helper()
	var found []semantics.EntityID
	for i, e := range graph.Data().Entities {
		if e.Name == name && len(e.Declarations) > 0 {
			found = append(found, semantics.EntityID(i+1))
		}
	}
	if len(found) != 1 {
		t.Fatalf("declarations named %s: %v", name, found)
	}
	return found[0]
}

// definitionBodyCase is the same five declarations in each language.
type definitionBodyCase struct {
	address, person, tree, page, people semantics.EntityID
}

// checkDefinitionBodies asserts what the walk must reach, one subtest for each:
//
//   - every field of a struct or object type (Address has one, Person three),
//   - a reference from one declaration to another (Person's third field is
//     Address),
//   - a recursive declaration (Tree's second field holds a reference to Tree,
//     and the read ends),
//   - a generic declaration with a use of its parameter (Page's first field
//     holds a reference to Page's binder),
//   - a use of the generic declaration whose argument is paired with that
//     binder (People is Page of Person).
//
// It does not show that a fact still reaches its consumer, and it does not
// show what an operand means: the reader cannot tell a field from an element.
func checkDefinitionBodies[L any](t *testing.T, graph semantics.Semantics[L], c definitionBodyCase) {
	read := func(t *testing.T, name string, id semantics.EntityID) structuralDeclaration {
		t.Helper()
		d, err := readStructurally(graph, id)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		return d
	}
	t.Run("every field", func(t *testing.T) {
		address, person := read(t, "Address", c.address), read(t, "Person", c.person)
		if !address.Body.Expression || len(address.Body.Operands) != 1 || !address.Body.Operands[0].Expression {
			t.Fatalf("Address: the walk does not reach its one field: %+v", address.Body)
		}
		if !person.Body.Expression || len(person.Body.Operands) != 3 {
			t.Fatalf("Person: the walk does not reach its three fields: %+v", person.Body)
		}
		if name, active := person.Body.Operands[0], person.Body.Operands[1]; !name.Expression || !active.Expression || name.Term == active.Term || name.Term != address.Body.Operands[0].Term {
			t.Fatalf("Person: the first two fields are not the two scalar types: %+v", person.Body)
		}
	})
	t.Run("reference to another declaration", func(t *testing.T) {
		person := read(t, "Person", c.person)
		if len(person.Body.Operands) != 3 {
			t.Fatalf("Person: the walk does not reach its three fields: %+v", person.Body)
		}
		if home := person.Body.Operands[2]; home.Callee != c.address || home.Expression || home.Bound != 0 || len(home.Actuals) != 0 {
			t.Fatalf("Person: the third field is no reference to the declaration Address: %+v", home)
		}
	})
	t.Run("recursive declaration", func(t *testing.T) {
		tree := read(t, "Tree", c.tree)
		if !tree.Body.Expression || len(tree.Body.Operands) != 2 {
			t.Fatalf("Tree: the walk does not reach its two fields: %+v", tree.Body)
		}
		if children := tree.Body.Operands[1]; !children.Expression || len(children.Operands) != 1 || children.Operands[0].Callee != c.tree {
			t.Fatalf("Tree: the second field does not hold a reference to Tree itself: %+v", children)
		}
	})
	t.Run("generic declaration with a use of its parameter", func(t *testing.T) {
		page := read(t, "Page", c.page)
		if len(page.Formals) != 1 || !page.Body.Expression || len(page.Body.Operands) != 2 {
			t.Fatalf("Page: the walk does not reach one binder and two fields: %+v", page)
		}
		if items := page.Body.Operands[0]; !items.Expression || len(items.Operands) != 1 || items.Operands[0].Bound != page.Formals[0] {
			t.Fatalf("Page: the first field does not hold a use of Page's parameter: %+v", items)
		}
	})
	t.Run("application paired with the binder", func(t *testing.T) {
		people := read(t, "People", c.people)
		formals, _ := graph.Entity(c.page)
		if people.Body.Expression || people.Body.Callee != c.page || len(people.Body.Actuals) != 1 || len(formals.Parameters) != 1 {
			t.Fatalf("People: the body is no application of Page with one actual: %+v", people.Body)
		}
		if actual := people.Body.Actuals[0]; actual.Binder != formals.Parameters[0].Entity || actual.Use.Callee != c.person {
			t.Fatalf("People: the actual is not Person paired with Page's parameter: %+v", actual)
		}
	})
}

func TestGoDeclarationBodyIsItsDefinition(t *testing.T) {
	snapshot, _ := goCapture(t, "type Address struct { City string }; type Person struct { Name string; Active bool; Home Address }; type Tree struct { Value string; Children []Tree }; type Page[T any] struct { Items []T; Next string }; type People = Page[Person]")
	_, result := goResults(t, snapshot)
	envelope, ok := result.Semantics()
	if !ok {
		t.Fatal("no graph", result.Diagnostics())
	}
	graph, err := representation.UnpackSemantics[golangsyntax.Language](snapshot.Language(), envelope)
	if err != nil {
		t.Fatal(err)
	}
	checkDefinitionBodies(t, graph, definitionBodyCase{
		address: declarationNamed(t, graph, "Address"), person: declarationNamed(t, graph, "Person"), tree: declarationNamed(t, graph, "Tree"),
		page: declarationNamed(t, graph, "Page"), people: declarationNamed(t, graph, "People"),
	})
}

func TestTypeScriptDeclarationBodyIsItsDefinition(t *testing.T) {
	snapshot, _, _ := tsCapture(t, "export interface Address { city: string }\nexport interface Person { name: string; active: boolean; home: Address }\nexport type Tree = { value: string; children: Tree[] };\nexport interface Page<T> { items: T[]; next: string }\nexport type People = Page<Person>;\nexport interface Employee extends Person, Page<Person> { desk: Address; role: string }\n")
	_, result := tsResults(t, snapshot)
	envelope, ok := result.Semantics()
	if !ok {
		t.Fatal("no graph", result.Diagnostics())
	}
	graph, err := representation.UnpackSemantics[typescriptsyntax.Language](snapshot.Language(), envelope)
	if err != nil {
		t.Fatal(err)
	}
	checkDefinitionBodies(t, graph, definitionBodyCase{
		address: declarationNamed(t, graph, "Address"), person: declarationNamed(t, graph, "Person"), tree: declarationNamed(t, graph, "Tree"),
		page: declarationNamed(t, graph, "Page"), people: declarationNamed(t, graph, "People"),
	})
	// An interface with base types: the walk reaches each base declaration,
	// the actual of a base with type arguments, and each member that the
	// interface declares itself. It does not reach an inherited member here:
	// that one is reached by opening the base.
	t.Run("interface with base types", func(t *testing.T) {
		person, address, page := declarationNamed(t, graph, "Person"), declarationNamed(t, graph, "Address"), declarationNamed(t, graph, "Page")
		employee, err := readStructurally(graph, declarationNamed(t, graph, "Employee"))
		if err != nil {
			t.Fatal(err)
		}
		if !employee.Body.Expression || len(employee.Body.Operands) != 3 {
			t.Fatalf("Employee: the walk does not reach two bases and the own members: %+v", employee.Body)
		}
		if first := employee.Body.Operands[0]; first.Callee != person || first.Expression || len(first.Actuals) != 0 {
			t.Fatalf("Employee: the first base is no reference to the declaration Person: %+v", first)
		}
		formals, _ := graph.Entity(page)
		if second := employee.Body.Operands[1]; second.Callee != page || len(second.Actuals) != 1 || second.Actuals[0].Binder != formals.Parameters[0].Entity || second.Actuals[0].Use.Callee != person {
			t.Fatalf("Employee: the second base is no application of Page to Person: %+v", second)
		}
		own := employee.Body.Operands[2]
		if !own.Expression || len(own.Operands) != 2 || own.Operands[0].Callee != address || !own.Operands[1].Expression || len(own.Operands[1].Operands) != 0 {
			t.Fatalf("Employee: the walk does not reach its two own members: %+v", own)
		}
	})
}
