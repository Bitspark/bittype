// Package list is the builtin horos-to-rust mapping library for horos lists.
// A consumer registers it on its own.
package list

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/rust/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// library is this package's one mapping library.
type library struct{}

// Library maps a horos list to the Rust model collection of its mapped
// element type.
func Library() g.MappingLibrary[h.Language, s.Language, source.Reader] { return library{} }

// Describe states the structural representation profile and the one binding.
func (library) Describe(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
	slot := g.TypeSlot("E")
	pattern, err := patterns.Element(from, slot)
	if err != nil {
		return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
	}
	return g.MappingDescription[h.Language, s.Language, source.Reader]{Profile: rules.RepresentationProfile, Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: "list/rust/1", Input: pattern, Output: model.Collection(g.Mapped[s.Language](slot))}}}, nil
}
