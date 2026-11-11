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

// Library returns the list mapping under the given profile.
func Library(profile string) rules.Library[s.Language] {
	slot := g.TypeSlot("E")
	return rules.Library[s.Language]{Profile: profile, DescribeRules: func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := patterns.Element(from, slot)
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: "list/golang/1", Input: pattern, Output: model.Collection(g.Mapped[s.Language](slot))}}}, nil
	}}
}
