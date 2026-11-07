package integration_test

import (
	"context"
	"github.com/Bitspark/bittype/contracts/parse/go"
	goformat "github.com/Bitspark/bittype/format/syntax/golang/provider/go"
	goparse "github.com/Bitspark/bittype/parse/syntax/golang/provider/go"
	printer "github.com/Bitspark/bittype/print/syntax/shared/provider/go"
	"testing"
)

// bittype:case target=github.com/Bitspark/bittype/parse/syntax/golang/provider/go#Parse feature=go-concrete-pipeline case=pagination-and-interface expects=round-trip-preserved
// bittype:case target=github.com/Bitspark/bittype/format/syntax/golang/provider/go#Format feature=go-concrete-pipeline case=pagination-and-interface expects=stable-formatting
// bittype:case target=github.com/Bitspark/bittype/print/syntax/shared/provider/go#Print feature=go-concrete-pipeline case=pagination-and-interface expects=source-retained language=github.com/Bitspark/bittype/syntax/golang/go#Language
func TestGoConcretePipeline(t *testing.T) {
	source := "package sample\r\n\r\n// pagination\r\ntype Page[T any] struct { Items []T; Next *string }\r\ntype PeoplePage = Page[string]\r\ntype Reader[T any] interface { Read(key string) (T, error) }\r\n"
	parsed, err := goparse.Parse(context.Background(), parse.Input{SourceID: "go-input", TreeID: "go-tree", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	printed, err := printer.Print(context.Background(), parsed)
	if err != nil || printed != source {
		t.Fatalf("lossless print: %v\n%q", err, printed)
	}
	formatted, err := goformat.Format(context.Background(), parsed)
	if err != nil {
		t.Fatal(err)
	}
	output, err := printer.Print(context.Background(), formatted)
	if err != nil {
		t.Fatal(err)
	}
	again, err := goparse.Parse(context.Background(), parse.Input{SourceID: "go-output", TreeID: "go-again", Text: output})
	if err != nil {
		t.Fatalf("%v\n%s", err, output)
	}
	stable, err := goformat.Format(context.Background(), again)
	if err != nil {
		t.Fatal(err)
	}
	output2, err := printer.Print(context.Background(), stable)
	if err != nil || output2 != output {
		t.Fatalf("unstable formatting: %v\n%s\n%s", err, output, output2)
	}
	parsed.Abstract.Nodes[parsed.Abstract.Root-1].Attributes[0].Text = "changed"
	if formatted.Abstract.Nodes[formatted.Abstract.Root-1].Attributes[0].Text != "sample" {
		t.Fatal("formatter aliases its input")
	}
	t.Log(output)
}
