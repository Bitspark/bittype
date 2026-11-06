package report

import (
	"slices"
	"strings"
	"testing"
)

func TestProviderGridPreservesDimensionsContractsAndSparseCells(t *testing.T) {
	p := fixture(t, map[string]string{
		"source/ledger/go": `package ledger
// bittype:dimension Ledger
type Document[L any] struct {}`,
		// An independent contract deliberately shares a display package/type name
		// with print.Provider. Go identity, not that spelling, owns its column.
		"contracts/index/go": `package print
import "example.test/fixture/contracts/language/go"
type Provider[I any] interface {language.Provider; Index(I)error}`,
		"index/ledger/ice/go": `package indexing
import (s "example.test/fixture/syntax/ice/go"; d "example.test/fixture/source/ledger/go")
type Indexer struct{s.Language}
func(Indexer)Index(d.Document[s.Language])error{return nil}`,
		"print/syntax/ember/go": `package alternative
import (s "example.test/fixture/syntax/ember/go"; c "example.test/fixture/trees/cst/go")
type Binary struct{s.Language}
func(Binary)Print(c.CST[s.Language])([]byte,error){return nil,nil}`,
	})
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	grid, err := collectProviderGrid(m)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(grid.dimensions, []string{"Files", "Ledger", "Syntax"}) || len(grid.columns) != 2 {
		t.Fatalf("lost or invented axis: %#v", grid)
	}
	if grid.columns[0].label != "example.test/fixture/contracts/index/go.Provider" || grid.columns[1].label != "example.test/fixture/contracts/print/go.Provider" {
		t.Fatal("colliding display names conflated contracts", grid.columns)
	}
	if len(grid.columns[1].rows) != 3 {
		t.Fatal("generic instantiations were discarded", grid.columns[1].rows)
	}
	for _, lang := range m.languages {
		identity, table, err := p.renderProviderGrid(grid, lang, "https://example.test/source", "details.md")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(identity, "Language") || strings.Contains(table, "language.Provider") || strings.Contains(table, "| Shared |") {
			t.Fatal("identity was made a pipeline stage", identity, table)
		}
		cells := tableCells(table)
		if lang.id == "ember" {
			if cells["Ledger"][0] != "/" || cells["Ledger"][1] != "/" || cells["Syntax"][0] != "/" {
				t.Fatal("invented empty-cell implementation", cells)
			}
			if strings.Count(cells["Syntax"][1], "details.md#cap-") != 2 {
				t.Fatal("distinct syntax implementations lost", cells["Syntax"])
			}
		} else {
			if cells["Files"][0] != "/" || cells["Files"][1] != "/" || cells["Syntax"][1] != "/" || !strings.Contains(cells["Ledger"][0], "Index") {
				t.Fatal("incorrect sparse cells", cells)
			}
		}
	}
}

func TestTransformerMatrixKeepsAllEndpointsAndOnlyDirectEdges(t *testing.T) {
	p := fixture(t, map[string]string{
		"syntax/snow/go":     `package snow; import "example.test/fixture/contracts/language/go"; type Language struct{};func(Language)ID()language.ID{return "snow"};func(Language)Name()string{return "Snow"}`,
		"transforms/next/go": `package next; import(t "example.test/fixture/contracts/transform/go"; a "example.test/fixture/syntax/ice/go"; b "example.test/fixture/syntax/snow/go");func New()t.ASTTransformer{return t.NewTransformer[a.Language,b.Language]()}`,
	})
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	table, err := p.renderTransformers(m, "https://example.test/source", "details.md")
	if err != nil {
		t.Fatal(err)
	}
	cells := tableCells(table)
	if len(cells) != 3 {
		t.Fatal("language row omitted", table)
	}
	for _, name := range []string{"Ember", "Ice", "Snow"} {
		if len(cells[name]) != 3 {
			t.Fatal("language column omitted", table)
		}
	}
	if cells["Ember"][0] != "/" || cells["Ember"][2] != "/" || cells["Ice"][0] != "/" || cells["Ice"][1] != "/" || !slices.Equal(cells["Snow"], []string{"/", "/", "/"}) {
		t.Fatal("inferred identity, reverse or composed edge", table)
	}
	if !strings.Contains(cells["Ember"][1], "details.md#cap-") || !strings.Contains(cells["Ice"][2], "details.md#cap-") {
		t.Fatal("direct edge evidence lost", table)
	}
	_, existing, err := p.render(m, "https://example.test/source", "details.md")
	if err != nil || existing != table {
		t.Fatal("matrix views disagree", err)
	}
}

func tableCells(table string) map[string][]string {
	rows := map[string][]string{}
	for _, line := range strings.Split(strings.TrimSpace(table), "\n")[2:] {
		parts := strings.Split(line, "|")
		var cells []string
		for _, part := range parts[2 : len(parts)-1] {
			cells = append(cells, strings.TrimSpace(part))
		}
		rows[strings.TrimSpace(parts[1])] = cells
	}
	return rows
}
