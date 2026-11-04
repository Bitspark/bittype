package consumer_test

import (
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	provider "github.com/Bitspark/bittype/provider/horos/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

func TestInstalledHorosTextPagination(t *testing.T) {
	p := provider.Provider{}
	source := "module installed.source; type Page<T> = { items: list<T>; next?: text; }; type DataPage<D: Data> = Page<D>; type People = DataPage<text>;"
	parsed, err := p.Parse(parse.Input{SourceID: "installed.horos", TreeID: "installed", Text: source})
	if err != nil { t.Fatal(err) }
	printed, err := p.Print(parsed)
	if err != nil || printed != source { t.Fatalf("source round trip: %v", err) }
	model, err := resolve.Model(parsed)
	if err != nil { t.Fatal(err) }
	if _, err := model.DataDescriptor(decl.DeclarationID{Module: model.Root(), Name: "People"}, nil); err != nil { t.Fatal(err) }
	formatted, err := p.Format(parsed)
	if err != nil { t.Fatal(err) }
	output, err := p.Print(formatted)
	if err != nil { t.Fatal(err) }
	reparsed, err := p.Parse(parse.Input{SourceID: "output.horos", TreeID: "output", Text: output})
	if err != nil { t.Fatal(err) }
	recovered, err := resolve.Model(reparsed)
	if err != nil || recovered.Root() != model.Root() { t.Fatalf("canonical round trip: %v", err) }
}
