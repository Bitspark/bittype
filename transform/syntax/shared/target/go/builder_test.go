package target

import (
	"context"
	"errors"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type fixtureLanguage struct{}

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
