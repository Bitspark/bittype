package integration_test

import (
	"context"
	"github.com/Bitspark/bittype/contracts/language/go"
	print "github.com/Bitspark/bittype/print/syntax/go"
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
	if err := checkConsumer(t, packages, "text/stages"); err != nil {
		t.Fatal(err)
	}
	for _, fixture := range []struct{ name, reason string }{
		{"missing_identity", "Name"},
		{"mixed_views", "cannot use"},
	} {
		diagnostics := consumerDiagnostics(t, packages, "text/"+fixture.name)
		boundaries := []string{"parse.Provider", "format.Provider", "print.Provider"}
		if fixture.name == "mixed_views" {
			boundaries = append(boundaries, "ast.View")
		}
		for _, provider := range boundaries {
			found := false
			for _, diagnostic := range diagnostics {
				message := diagnostic.Error()
				if strings.Contains(message, provider) && strings.Contains(message, fixture.reason) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("%s did not reject %s: %v", provider, fixture.name, diagnostics)
			}
		}
	}
}

type unfamiliarLanguage struct{}

func (unfamiliarLanguage) ID() language.ID { return "unfamiliar" }
func (unfamiliarLanguage) Name() string    { return "Unfamiliar" }
func TestPrinterNeedsNoKnownLanguage(t *testing.T) {
	var _ language.Provider = unfamiliarLanguage{}
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
	printer, err := print.New[unfamiliarLanguage](unfamiliarLanguage{})
	if err != nil {
		t.Fatal(err)
	}
	output, err := printer.Print(context.Background(), concrete)
	if err != nil || output != "hello\r\n" {
		t.Fatalf("%v: %q", err, output)
	}
}
