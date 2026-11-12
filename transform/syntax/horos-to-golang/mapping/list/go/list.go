// Package list is the horos-to-golang mapping library for lists: a horos
// list-of becomes a Go collection of its mapped element.
package list

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// library is this package's one mapping library.
type library struct{}

// Library returns the list mapping. It states the structural representation
// profile itself.
func Library() g.MappingLibrary[h.Language, s.Language, source.Reader] { return library{} }

// Describe states the structural representation profile and the one binding.
func (library) Describe(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
	slot := g.TypeSlot("E")
	pattern, err := patterns.Element(from, slot)
	if err != nil {
		return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
	}
	return g.MappingDescription[h.Language, s.Language, source.Reader]{Profile: rules.RepresentationProfile, Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: "list/golang/1", Input: pattern, Output: model.Collection(g.Mapped[s.Language](slot))}}}, nil
}
