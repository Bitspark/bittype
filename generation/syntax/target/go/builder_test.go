package target_test

import (
	"context"
	"errors"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	m "github.com/Bitspark/bittype/generation/syntax/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type native struct{}

func (native) ID() language.ID { return "independent-native" }
func (native) Name() string    { return "Independent native" }

func TestNativeStorageCapturesWithoutSharing(t *testing.T) {
	language, _ := representation.NewLanguage[native](native{}, "v1")
	store, err := m.New[native, native](context.Background(), language, "output", g.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	node := m.Node(ast.Kind[native]("type"), m.Attributes("name", "Original"), ast.NodeRef{Tree: "source", Node: 1})
	fragment, err := store.Node(node)
	if err != nil {
		t.Fatal(err)
	}
	node.Attributes[0].Text = "Mutated"
	node.Origin.Nodes[0].Node = 2
	first, _, err := store.Inspect(fragment)
	if err != nil || first.Attributes[0].Text != "Original" || first.Origin.Nodes[0].Node != 1 {
		t.Fatal("captured storage changed", err)
	}
	first.Attributes[0].Text = "Again"
	second, _, _ := store.Inspect(fragment)
	if second.Attributes[0].Text != "Original" {
		t.Fatal("inspection exposed retained storage")
	}
	foreign, _ := m.New[native, native](context.Background(), language, "other", g.DefaultLimits())
	if _, err := foreign.Node(m.Node(ast.Kind[native]("container"), nil, ast.NodeRef{}), m.Child[native]{Role: "element", Fragment: fragment}); !errors.Is(err, g.ErrScope) {
		t.Fatal("foreign fragment accepted", err)
	}
	node.References = []ast.Reference{{Role: "target", Target: ast.NodeRef{Tree: "invented", Node: 1}}}
	if _, err := store.Node(node); !errors.Is(err, g.ErrInvalid) {
		t.Fatal("raw arena reference accepted", err)
	}
}

func TestNativeStorageLimitsAndCancellation(t *testing.T) {
	language, _ := representation.NewLanguage[native](native{}, "v1")
	limits := g.DefaultLimits()
	limits.OutputNodes = 1
	store, _ := m.New[native, native](context.Background(), language, "output", limits)
	node := m.Node(ast.Kind[native]("type"), nil, ast.NodeRef{})
	if _, err := store.Node(node); err != nil {
		t.Fatal(err)
	}
	if fragment, err := store.Node(node); fragment.Valid() || !errors.Is(err, g.ErrLimit) {
		t.Fatal("node limit accepted output", err)
	}
	limits.OutputBytes = 1
	small, _ := m.New[native, native](context.Background(), language, "small", limits)
	if fragment, err := small.Node(node); fragment.Valid() || !errors.Is(err, g.ErrLimit) {
		t.Fatal("byte limit accepted output", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	canceled, _ := m.New[native, native](ctx, language, "canceled", g.DefaultLimits())
	cancel()
	if fragment, err := canceled.Node(node); fragment.Valid() || !errors.Is(err, context.Canceled) {
		t.Fatal("canceled allocation succeeded", err)
	}
}
