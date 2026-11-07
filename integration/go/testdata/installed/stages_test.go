package consumer_test

import (
	checktypescript "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	check "github.com/Bitspark/bittype/contracts/check/go"
	format "github.com/Bitspark/bittype/contracts/format/go"
	print "github.com/Bitspark/bittype/contracts/print/go"
	formattypescript "github.com/Bitspark/bittype/format/syntax/typescript/provider/go"
	textfixture "github.com/Bitspark/bittype/integration/go/testdata/textfixture"
	parsetypescript "github.com/Bitspark/bittype/parse/syntax/typescript/provider/go"

	"context"
	"errors"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/engine/go"

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
	parser, err := parse.NewSyntaxParser[ts.Language](definition, parsetypescript.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[ts.Language](definition, checktypescript.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	formatter, err := format.NewSyntaxFormatter[ts.Language](definition, formattypescript.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	printer, err := print.NewSyntaxPrinter[ts.Language](definition, textfixture.Printing[ts.Language](ts.Language{}))
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
