package integration_test

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	print "github.com/Bitspark/bittype/print/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"strings"
	"testing"
)

func TestIndependentTextProviderContracts(t *testing.T) {
	packages := loadContracts(t)
	if err := checkConsumer(t, packages, "text/provider"); err != nil {
		t.Fatal(err)
	}
	if err := checkConsumer(t, packages, "text/missing_identity"); err == nil || !strings.Contains(err.Error(), "Name") {
		t.Fatalf("missing Name not rejected: %v", err)
	}
	if err := checkConsumer(t, packages, "text/mixed_views"); err == nil {
		t.Fatal("cross-language AST view accepted")
	}
}

type unfamiliarLanguage struct{}

func (unfamiliarLanguage) ID() language.ID { return "unfamiliar" }
func (unfamiliarLanguage) Name() string    { return "Unfamiliar" }
func TestPrinterNeedsNoKnownLanguage(t *testing.T) {
	var _ language.LanguageProvider = unfamiliarLanguage{}
	value := "hello"
	space := "\r\n"
	concrete := cst.CST[unfamiliarLanguage]{
		Abstract: ast.AST[unfamiliarLanguage]{ID: "foreign-tree", Root: 1, Nodes: []ast.Node[unfamiliarLanguage]{{Kind: "foreign-declaration"}}},
		Root:     3,
		Nodes: []cst.Node[unfamiliarLanguage]{
			{Kind: "foreign-word", Text: &value},
			{Kind: "foreign-layout", Text: &space},
			{Kind: "foreign-sentence", Children: []cst.NodeID{1, 2}, Syntax: []ast.NodeID{1}},
		},
	}
	output, err := print.Print(concrete)
	if err != nil || output != "hello\r\n" {
		t.Fatalf("%v: %q", err, output)
	}
}
