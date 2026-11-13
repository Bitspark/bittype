package integration_test

import (
	"context"
	"errors"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	gotarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// These qualify the shared builder without a rule about names across
// declarations (#273). A target that needs one flat scope states it in its own
// program rule. The builder identifies a declaration by its planned identity
// and a parameter by its position, and reads a name only to write it.

// ruleless is the fixture's Go target without a program rule.
type ruleless struct{ gotarget.Provider }

func (ruleless) CheckProgram(context.Context, g.Program[gs.Language]) error { return nil }

// A target that states no program rule builds two declarations of one name. It
// gets two sites, and each reference is bound to its own declaration: People is
// Page of Person, and both Page and Person are named Same.
func TestBuilderBuildsTwoDeclarationsOfOneName(t *testing.T) {
	ctx := context.Background()
	for _, mode := range []f.Mode{f.SameName, f.ParameterNamedSame} {
		fixture := independent(t)
		plan, err := fixture.PlanTargetIn(ctx, mode, ruleless{gotarget.Provider{Package: "generated"}}, "People")
		if err != nil {
			t.Fatal("planning two declarations of one name", mode, err)
		}
		out, err := plan.Run(ctx)
		if err != nil {
			t.Fatal("the builder refused two declarations of one name", mode, err)
		}
		tree := outputTree[gs.Language](t, out)
		bySource := map[ast.NodeRef]g.OwnedDeclaration{}
		for _, record := range out.Inventory() {
			if len(record.Provenance) == 1 {
				bySource[record.Provenance[0]] = record
			}
		}
		person, page, people := bySource[fixture.Site("Person")], bySource[fixture.Site("Page")], bySource[fixture.Site("People")]
		if !person.Handle.Valid() || !page.Handle.Valid() || !people.Handle.Valid() || person.Site == page.Site {
			t.Fatal("the two declarations of one name do not have two sites", mode, person.Site, page.Site)
		}
		for _, record := range []g.OwnedDeclaration{person, page} {
			if name := nodeAttribute(tree.Nodes[record.Site.Node-1], gs.AttributeName); name != "Same" {
				t.Fatal("the declaration is not named Same", mode, name)
			}
		}
		// The one application below People: its constructor is Page, and its first
		// argument is Person.
		var application ast.NodeID
		var find func(ast.NodeID)
		find = func(id ast.NodeID) {
			if tree.Nodes[id-1].Kind == gs.KindTypeApplication && application == 0 {
				application = id
			}
			for _, child := range tree.Nodes[id-1].Children {
				find(child.Node)
			}
		}
		find(people.Site.Node)
		if application == 0 {
			t.Fatal("People holds no application", mode)
		}
		target := func(role string) ast.NodeRef {
			for _, child := range tree.Nodes[application-1].Children {
				if child.Role == role {
					node := tree.Nodes[child.Node-1]
					if len(node.References) != 1 || nodeAttribute(node, gs.AttributeName) != "Same" {
						t.Fatal("the use is not one reference under the name Same", mode, role, node.References)
					}
					return node.References[0].Target
				}
			}
			t.Fatal("the application has no child in the role", mode, role)
			return ast.NodeRef{}
		}
		if constructor := target(gs.RoleConstructor); constructor != page.Site {
			t.Fatal("the constructor is not bound to the declaration Page", mode, constructor, page.Site)
		}
		if argument := target(gs.RoleArgument); argument != person.Site {
			t.Fatal("the argument is not bound to the declaration Person", mode, argument, person.Site)
		}
	}
}

// The flat scope is a rule of the target that states it. The Go target refuses
// the same two programs before anything is built, as it did.
func TestGoTargetStillStatesOneFlatScope(t *testing.T) {
	ctx := context.Background()
	for _, mode := range []f.Mode{f.SameName, f.ParameterNamedSame, f.SameParameters} {
		fixture := independent(t)
		plan, err := fixture.PlanTargetIn(ctx, mode, gotarget.Provider{Package: "generated"}, "People")
		if err == nil {
			var out g.Output
			out, err = plan.Run(ctx)
			if out.Snapshot().Valid() {
				t.Fatal("the Go target built names that share its one scope", mode)
			}
		}
		if !errors.Is(err, g.ErrUnsupported) {
			t.Fatal("the Go target's program rule refuses with the unsupported value", mode, err)
		}
	}
}

// The builder keeps the one rule whose meaning is inside a declaration: two
// parameters of one declaration cannot share a name, because a parameter is
// written by its name.
func TestBuilderRefusesTwoParametersOfOneNameInOneDeclaration(t *testing.T) {
	ctx := context.Background()
	fixture := independent(t)
	plan, err := fixture.PlanTargetIn(ctx, f.SameParameters, ruleless{gotarget.Provider{Package: "generated"}}, "People")
	if err == nil {
		var out g.Output
		out, err = plan.Run(ctx)
		if out.Snapshot().Valid() {
			t.Fatal("the builder built a declaration with two parameters of one name")
		}
	}
	if !errors.Is(err, g.ErrScope) {
		t.Fatal("two parameters of one name in one declaration", err)
	}
}
