package consumer_test

import (
	"context"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	formatter "github.com/Bitspark/bittype/format/horos/go"
	parser "github.com/Bitspark/bittype/parse/horos/go"
	provider "github.com/Bitspark/bittype/provider/horos/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

func TestInstalledHorosTextPagination(t *testing.T) {
	p := provider.Provider{}
	source := "module installed.source; import h from horos.kinds pin \"sha256:" + decl.PreludePin().String() + "\"; type Page<T> = { items: list<T>; next?: text; }; type DataPage<D: h.Data> = Page<D>; type People = DataPage<text>;"
	parsed, err := (parser.Provider{}).Parse(context.Background(), parse.Input{SourceID: "installed.horos", TreeID: "installed", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	printed, err := p.Print(context.Background(), parsed)
	if err != nil || printed != source {
		t.Fatalf("source round trip: %v", err)
	}
	model, err := resolve.Model(parsed)
	if err != nil {
		t.Fatal(err)
	}
	value, err := model.Value(decl.DeclarationID{Module: model.Root(), Name: "People"})
	if err != nil { t.Fatal(err) }
	if _, err := model.GroundDescriptor(value); err != nil {
		t.Fatal(err)
	}
	formatted, err := (formatter.Provider{}).Format(context.Background(), parsed)
	if err != nil {
		t.Fatal(err)
	}
	output, err := p.Print(context.Background(), formatted)
	if err != nil {
		t.Fatal(err)
	}
	reparsed, err := p.Parse(context.Background(), parse.Input{SourceID: "output.horos", TreeID: "output", Text: output})
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := resolve.Model(reparsed)
	if err != nil || recovered.Root() != model.Root() {
		t.Fatalf("canonical round trip: %v", err)
	}
}
