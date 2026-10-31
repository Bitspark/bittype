package integration_test

import (
	goformat "github.com/Bitspark/bittype/format/golang/go"
	goparse "github.com/Bitspark/bittype/parse/golang/go"
	printer "github.com/Bitspark/bittype/print/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"testing"
)

func TestGoConcretePipeline(t *testing.T) {
	source := "package sample\r\n\r\n// pagination\r\ntype Page[T any] struct { Items []T; Next *string }\r\ntype PeoplePage = Page[string]\r\ntype Reader[T any] interface { Read(key string) (T, error) }\r\n"
	parsed, err := goparse.Parse(text.Input{SourceID: "go-input", TreeID: "go-tree", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	printed, err := printer.Print(parsed)
	if err != nil || printed != source {
		t.Fatalf("lossless print: %v\n%q", err, printed)
	}
	formatted, err := goformat.Format(parsed)
	if err != nil {
		t.Fatal(err)
	}
	output, err := printer.Print(formatted)
	if err != nil {
		t.Fatal(err)
	}
	again, err := goparse.Parse(text.Input{SourceID: "go-output", TreeID: "go-again", Text: output})
	if err != nil {
		t.Fatalf("%v\n%s", err, output)
	}
	stable, err := goformat.Format(again)
	if err != nil {
		t.Fatal(err)
	}
	output2, err := printer.Print(stable)
	if err != nil || output2 != output {
		t.Fatalf("unstable formatting: %v\n%s\n%s", err, output, output2)
	}
	parsed.Abstract.Nodes[parsed.Abstract.Root-1].Attributes[0].Text = "changed"
	if formatted.Abstract.Nodes[formatted.Abstract.Root-1].Attributes[0].Text != "sample" {
		t.Fatal("formatter aliases its input")
	}
	t.Log(output)
}
