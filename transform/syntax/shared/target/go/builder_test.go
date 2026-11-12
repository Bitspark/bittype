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

func plannedFixture(id uint32, name string, tree ast.TreeID, node ast.NodeID) g.PlannedDeclaration[fixtureLanguage] {
	return g.PlannedDeclaration[fixtureLanguage]{ID: id, Description: g.NativeDeclaration[fixtureLanguage]{Name: name, Origin: ast.NodeRef{Tree: tree, Node: node}}}
}

// Two planned declarations of one name are ordered, and the same way in both
// directions.
func TestCompareDeclarationsOrdersTwoDeclarationsOfOneName(t *testing.T) {
	a, b := plannedFixture(1, "Page", "source", 3), plannedFixture(2, "Page", "source", 3)
	if ab, ba := compareDeclarations(a, b), compareDeclarations(b, a); ab != -ba || ab == 0 {
		t.Fatal("two declarations of one name have no order", ab, ba)
	}
}

// The key is the name, then the origin, then the planned identity. The result
// is 0 only for one and the same planned declaration.
func TestCompareDeclarationsKey(t *testing.T) {
	ordered := []g.PlannedDeclaration[fixtureLanguage]{
		plannedFixture(9, "A", "t", 5), // the name decides before everything else
		plannedFixture(8, "B", "", 0),  // then the origin: none, as a helper has
		plannedFixture(7, "B", "s", 9), // its tree
		plannedFixture(6, "B", "t", 2), // its node
		plannedFixture(4, "B", "t", 3), // then the planned identity
		plannedFixture(5, "B", "t", 3),
		plannedFixture(1, "C", "a", 1),
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
	// What the key does not name does not order: the role and the provenance.
	same := ordered[4]
	same.Role, same.Provenance = g.HelperDeclaration, []ast.NodeRef{{Tree: "other", Node: 1}}
	if compareDeclarations(ordered[4], same) != 0 {
		t.Fatal("the comparator reads more than its key")
	}
	// Every arrangement sorts to the one order.
	shuffled := []g.PlannedDeclaration[fixtureLanguage]{ordered[5], ordered[2], ordered[6], ordered[0], ordered[4], ordered[3], ordered[1]}
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
