package generation

import (
	"context"
	"errors"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type outputLanguage struct{}

func (outputLanguage) ID() language.ID { return "output-test" }
func (outputLanguage) Name() string    { return "Output test" }

type outputChecker struct{ outputLanguage }

func (outputChecker) Check(ast.View[outputLanguage]) error { return nil }
func outputAssembly(t *testing.T) (check.Checker, Assembly) {
	t.Helper()
	lang, err := representation.NewLanguage[outputLanguage](outputLanguage{}, "one")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[outputLanguage](lang, outputChecker{})
	if err != nil {
		t.Fatal(err)
	}
	tree := ast.AST[outputLanguage]{ID: "output", Root: 1, Nodes: []ast.Node[outputLanguage]{{Kind: "unit", Children: []ast.Child{{Role: "declaration", Node: 2}}}, {Kind: "declaration"}}}
	packed, err := representation.PackSet(lang, ast.Set[outputLanguage]{Entry: ast.NodeRef{Tree: tree.ID, Node: 1}, Trees: []ast.AST[outputLanguage]{tree}})
	if err != nil {
		t.Fatal(err)
	}
	return checker, Assembly{Syntax: packed, Exposed: []ast.NodeRef{{Tree: tree.ID, Node: 2}}}
}
func TestSealOwnsItsSitesAndRejectsIncompleteMetadata(t *testing.T) {
	checker, assembly := outputAssembly(t)
	out, err := seal(context.Background(), checker, assembly, 1, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	assembly.Exposed[0] = ast.NodeRef{Tree: "foreign", Node: 1}
	sites := out.Exposed()
	sites[0] = ast.NodeRef{}
	if out.Exposed()[0].Tree != "output" {
		t.Fatal("mutable exposure metadata")
	}
	if bad, err := seal(context.Background(), checker, assembly, 1, DefaultLimits()); !errors.Is(err, ErrScope) || bad.Snapshot().Valid() {
		t.Fatal(err)
	}
	_, assembly = outputAssembly(t)
	if bad, err := seal(context.Background(), checker, assembly, 2, DefaultLimits()); !errors.Is(err, ErrIncomplete) || bad.Snapshot().Valid() {
		t.Fatal(err)
	}
}
func TestSealRejectsForeignDuplicateAndOversizedExternalUses(t *testing.T) {
	checker, base := outputAssembly(t)
	symbol, err := NewSymbolIssuer().Issue(checker.Language(), "module/type")
	if err != nil {
		t.Fatal(err)
	}
	at := semantics.Occurrence{Site: base.Exposed[0], Role: semantics.Key{Owner: "fixture", Name: "type"}}
	for name, uses := range map[string][]ExternalUse{
		"foreign":   {{At: semantics.Occurrence{Site: ast.NodeRef{Tree: "other", Node: 2}, Role: at.Role}, Symbol: symbol}},
		"duplicate": {{At: at, Symbol: symbol}, {At: at, Symbol: symbol}},
		"unissued":  {{At: at}},
	} {
		t.Run(name, func(t *testing.T) {
			a := base
			a.ExternalUses = uses
			if out, err := seal(context.Background(), checker, a, 1, DefaultLimits()); err == nil || out.Snapshot().Valid() {
				t.Fatal("invalid external use accepted")
			}
		})
	}
	base.ExternalUses = []ExternalUse{{At: at, Symbol: symbol}}
	limits := DefaultLimits()
	limits.OutputBytes = 2
	if out, err := seal(context.Background(), checker, base, 1, limits); !errors.Is(err, ErrLimit) || out.Snapshot().Valid() {
		t.Fatal(err)
	}
}
