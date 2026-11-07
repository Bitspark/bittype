package generation

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/language/go"

	"github.com/Bitspark/bittype/trees/representation/go"
	"testing"
)

type source struct{}

func (source) ID() language.ID { return "independent" }
func (source) Name() string    { return "Independent" }

type renamed struct{}

func (renamed) ID() language.ID { return "independent" }
func (renamed) Name() string    { return "A different display name" }
func TestQualifiedReplacement(t *testing.T) {
	lang, _ := representation.NewLanguage[source](source{}, "v1")
	otherName, _ := representation.NewLanguage[source](renamed{}, "v1")
	key, _ := g.NewSelector(lang, "same")
	equivalent, _ := g.NewSelector(otherName, "same")
	missing, _ := g.NewSelector(lang, "different")
	r := &Registry{}
	base := MappingRef{r, 1}
	selected := MappingRef{r, 2}
	part := func(k g.Selector, id string) g.Contribution[source, source] {
		m := g.Mapping[source, source]{ID: id, PluginID: "owner", Plan: func(g.Declaration[source]) (g.Plan[source, source], error) { return g.Plan[source, source]{}, nil }}
		return g.Contribution[source, source]{Bindings: []g.RepresentationBinding[source, source]{{Selector: k, Mapping: m}}}
	}
	parts := []g.Contribution[source, source]{part(key, "same-id"), part(equivalent, "same-id")}
	replacement := Replacement{Selector: key, Expected: BindingRef{base, "same-id"}, With: BindingRef{selected, "same-id"}}
	selection := Selection{Mappings: []MappingRef{base, selected}, Replace: []Replacement{replacement}}
	table, _, err := collect(lang, selection, parts)
	if err != nil || len(table.entries) != 1 {
		t.Fatal(table, err)
	}
	for _, bad := range [][]Replacement{nil, {replacement, replacement}, {{Selector: missing, Expected: replacement.Expected, With: replacement.With}}, {{Selector: key, Expected: BindingRef{base, "stale"}, With: replacement.With}}, {{Selector: key, Expected: replacement.With, With: replacement.With}}} {
		selection.Replace = bad
		if _, _, err := collect(lang, selection, parts); err == nil {
			t.Fatal("accepted ambiguous or stale replacement", bad)
		}
	}
	wrong, _ := representation.NewLanguage[source](source{}, "other")
	foreign, _ := g.NewSelector(wrong, "same")
	if _, _, err := collect(lang, Selection{Mappings: []MappingRef{base}}, []g.Contribution[source, source]{part(foreign, "x")}); err == nil {
		t.Fatal("accepted foreign profile")
	}
}
