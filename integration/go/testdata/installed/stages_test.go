package consumer_test

import (
"github.com/Bitspark/bittype/contracts/transform/go"
"github.com/Bitspark/bittype/contracts/semantic/go"
"context"
	"errors"
	nativecheck "github.com/Bitspark/bittype/check/typescript/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/engine/go"
	nativeformat "github.com/Bitspark/bittype/format/typescript/go"
	nativeparse "github.com/Bitspark/bittype/parse/typescript/go"
	nativeprint "github.com/Bitspark/bittype/print/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"testing"
)

type Parser struct{ ts.Language }

func (Parser) Parse(input parse.Input) (cst.CST[ts.Language], error) { return nativeparse.Parse(input) }

type Checker struct{ ts.Language }

func (Checker) Check(input ast.View[ts.Language]) error { return nativecheck.Check(input) }

type Formatter struct{ ts.Language }

func (Formatter) Format(input ast.View[ts.Language]) (cst.CST[ts.Language], error) {
	return nativeformat.Format(input)
}

type Printer struct{ ts.Language }

func (Printer) Print(input cst.CST[ts.Language]) (string, error) { return nativeprint.Print(input) }

func TestIndependentNativeStages(t *testing.T) {
	definition, err := representation.NewLanguage[ts.Language](ts.Language{}, "native-types-v1")
	if err != nil {
		t.Fatal(err)
	}
	parser, err := parse.NewParser[ts.Language](definition, Parser{})
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[ts.Language](definition, Checker{})
	if err != nil {
		t.Fatal(err)
	}
	formatter, err := format.NewFormatter[ts.Language](definition, Formatter{})
	if err != nil {
		t.Fatal(err)
	}
	printer, err := print.NewPrinter[ts.Language](definition, Printer{})
	if err != nil {
		t.Fatal(err)
	}
	e := engine.New()
	if err := e.RegisterLanguage(definition); err != nil {
		t.Fatal(err)
	}
	if err := e.RegisterParser(parser); err != nil {
		t.Fatal(err)
	}
	source := "// owned by consumer\ntype Page<T> = T[];\n"
	parsed, err := e.Parse(definition.ID(), parse.Input{SourceID: "independent", TreeID: "independent", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	if err := e.Check(parsed.ASTView()); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal(err)
	}
	if err := e.RegisterPrinter(printer); err != nil {
		t.Fatal(err)
	}
	if got, err := e.Print(parsed); err != nil || got != source {
		t.Fatal(got, err)
	}
	if err := e.RegisterFormatter(formatter); err != nil {
		t.Fatal(err)
	}
	if _, err := e.Render(parsed.ASTView()); err != nil {
		t.Fatal(err)
	}
	if err := e.RegisterChecker(checker); err != nil {
		t.Fatal(err)
	}
	info, err := parsed.ASTView().Inspect()
if err != nil { t.Fatal(err) }
snapshot, err := e.Capture(context.Background(), semantic.Input{Entry: ast.NodeRef{Tree: info.ID, Node: info.Root}, Trees: []representation.AST{parsed.ASTView()}})
if err != nil { t.Fatal(err) }
if _, err := e.TransformTo(context.Background(), transform.Input{Source: snapshot}, definition.ID()); err != nil {
		t.Fatal(err)
	}
}
