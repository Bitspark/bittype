package target

import (
	"context"
	"errors"
	"slices"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type fixtureLanguage struct{}

// plannedFixture is a planned declaration with the provenance the planner
// states and the origin its mapping states.
func plannedFixture(id uint32, name string, origin ast.NodeID, provenance ...ast.NodeRef) g.PlannedDeclaration[fixtureLanguage] {
	return g.PlannedDeclaration[fixtureLanguage]{ID: id, Provenance: provenance, Description: g.TargetDeclaration[fixtureLanguage]{Name: name, Origin: ast.NodeRef{Tree: "origin", Node: origin}}}
}

func sourceSite(tree ast.TreeID, node ast.NodeID) ast.NodeRef {
	return ast.NodeRef{Tree: tree, Node: node}
}

// Two planned declarations of one name are ordered, and the same way in both
// directions.
func TestCompareDeclarationsOrdersTwoDeclarationsOfOneName(t *testing.T) {
	a, b := plannedFixture(1, "Page", 3, sourceSite("source", 3)), plannedFixture(2, "Page", 3, sourceSite("source", 3))
	if ab, ba := compareDeclarations(a, b), compareDeclarations(b, a); ab != -ba || ab == 0 {
		t.Fatal("two declarations of one name have no order", ab, ba)
	}
}

// The key is the name, then the provenance, then the planned identity. The
// result is 0 only for one and the same planned declaration. The origins and
// the planned identities of this list run against its order wherever an earlier
// part of the key decides, so neither can stand in for the provenance.
func TestCompareDeclarationsKey(t *testing.T) {
	ordered := []g.PlannedDeclaration[fixtureLanguage]{
		plannedFixture(12, "A", 90, sourceSite("t", 5)),                    // the name decides before everything else
		plannedFixture(11, "B", 80),                                        // then the provenance: none, as a helper has
		plannedFixture(10, "B", 70, sourceSite("s", 9)),                    // the tree of its first site
		plannedFixture(9, "B", 60, sourceSite("t", 2)),                     // the node of its first site
		plannedFixture(8, "B", 50, sourceSite("t", 2), sourceSite("a", 1)), // a list after the list it continues
		plannedFixture(7, "B", 40, sourceSite("t", 2), sourceSite("a", 4)), // the node of a later site
		plannedFixture(6, "B", 30, sourceSite("t", 2), sourceSite("b", 1)), // the tree of a later site
		plannedFixture(4, "B", 20, sourceSite("t", 3)),                     // then the planned identity
		plannedFixture(5, "B", 10, sourceSite("t", 3)),
		plannedFixture(1, "C", 5, sourceSite("a", 1)),
	}
	for i, a := range ordered {
		for j, b := range ordered {
			want := 0
			if i < j {
				want = -1
			} else if i > j {
				want = 1
			}
			if got := compareDeclarations(a, b); got != want {
				t.Fatalf("compare(%d, %d) = %d, want %d", a.ID, b.ID, got, want)
			}
			if compareDeclarations(a, b) != -compareDeclarations(b, a) {
				t.Fatalf("compare(%d, %d) is not the opposite of compare(%d, %d)", a.ID, b.ID, b.ID, a.ID)
			}
		}
	}
	// What the key does not name does not order: the role, and the origin that
	// a mapping may state or leave out.
	same := ordered[7]
	same.Role, same.Description.Origin = g.HelperDeclaration, ast.NodeRef{}
	if compareDeclarations(ordered[7], same) != 0 {
		t.Fatal("the comparator reads more than its key")
	}
	// Every arrangement sorts to the one order.
	shuffled := []g.PlannedDeclaration[fixtureLanguage]{ordered[5], ordered[8], ordered[2], ordered[9], ordered[6], ordered[0], ordered[4], ordered[7], ordered[3], ordered[1]}
	slices.SortFunc(shuffled, compareDeclarations[fixtureLanguage])
	for i := range ordered {
		if shuffled[i].ID != ordered[i].ID {
			t.Fatal("sorting gave another order at", i)
		}
	}
}

func TestConstructionFragmentsAndSymbolsBelongToTheirRun(t *testing.T) {
	a := New[fixtureLanguage](context.Background(), g.DefaultLimits())
	b := New[fixtureLanguage](context.Background(), g.DefaultLimits())
	af, err := a.Node(ast.Node[fixtureLanguage]{Kind: "name"})
	if err != nil {
		t.Fatal(err)
	}
	bf, err := b.Node(ast.Node[fixtureLanguage]{Kind: "name"})
	if err != nil {
		t.Fatal(err)
	}
	symbol := Symbol[fixtureLanguage]{owner: a, declaration: 1, parameter: -1, name: "Example"}
	if _, err = a.Link(af, symbol, "target"); err != nil {
		t.Fatal(err)
	}
	if _, err = b.Link(bf, symbol, "target"); !errors.Is(err, g.ErrScope) {
		t.Fatal("foreign symbol accepted", err)
	}
	if _, err = b.Node(ast.Node[fixtureLanguage]{Kind: "array"}, Child[fixtureLanguage]{Role: "element", Fragment: af}); !errors.Is(err, g.ErrScope) {
		t.Fatal("foreign fragment accepted", err)
	}
	if _, err = a.Node(ast.Node[fixtureLanguage]{Kind: "name", References: []ast.Reference{{Role: "target", Target: ast.NodeRef{Tree: "fake", Node: 1}}}}); !errors.Is(err, g.ErrInvalid) {
		t.Fatal("unplanned reference accepted", err)
	}
}
