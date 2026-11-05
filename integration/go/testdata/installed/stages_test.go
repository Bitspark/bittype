package consumer_test

import (
	"context"
	"errors"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/engine/go"
	tsprovider "github.com/Bitspark/bittype/provider/typescript/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"testing"
)

func TestIndependentNativeStages(t *testing.T) {
	definition, err := representation.NewLanguage[ts.Language](ts.Language{}, "native-types-v1")
	if err != nil {
		t.Fatal(err)
	}
	parser, err := tsprovider.NewSyntaxParser(definition)
	if err != nil {
		t.Fatal(err)
	}
	checker, err := tsprovider.NewSyntaxChecker(definition)
	if err != nil {
		t.Fatal(err)
	}
	formatter, err := tsprovider.NewSyntaxFormatter(definition)
	if err != nil {
		t.Fatal(err)
	}
	printer, err := tsprovider.NewSyntaxPrinter(definition)
	if err != nil {
		t.Fatal(err)
	}
	e := engine.New()
	if err := e.RegisterLanguage(definition); err != nil {
		t.Fatal(err)
	}
	if err := e.RegisterSyntaxParser(parser); err != nil {
		t.Fatal(err)
	}
	source := "// owned by consumer\ntype Page<T> = T[];\n"
	parsed, err := e.ParseSyntax(context.Background(), definition.ID(), parse.Input{SourceID: "independent", TreeID: "independent", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	if err := e.Check(parsed.ASTView()); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal(err)
	}
	if err := e.RegisterSyntaxPrinter(printer); err != nil {
		t.Fatal(err)
	}
	if got, err := e.PrintSyntax(context.Background(), parsed); err != nil || got != source {
		t.Fatal(got, err)
	}
	if err := e.RegisterSyntaxFormatter(formatter); err != nil {
		t.Fatal(err)
	}
	if _, err := e.RenderSyntax(context.Background(), parsed.ASTView()); err != nil {
		t.Fatal(err)
	}
	if err := e.RegisterChecker(checker); err != nil {
		t.Fatal(err)
	}
	info, err := parsed.ASTView().Inspect()
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := e.Capture(context.Background(), semantic.Input{Entry: ast.NodeRef{Tree: info.ID, Node: info.Root}, Trees: []representation.AST{parsed.ASTView()}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.TransformTo(context.Background(), transform.Input{Source: snapshot}, definition.ID()); err != nil {
		t.Fatal(err)
	}
}
