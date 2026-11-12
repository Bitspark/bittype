// Package record is the builtin horos-to-cpp mapping library for closed horos
// records. It describes a C++ record whose fields have their mapped types.
package record

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/cpp/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Library maps a closed record to a record of its mapped fields under the
// binding record/cpp/1. An optional field is wrapped in std::optional here.
func Library() rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := patterns.Head(from, "record", "closed")
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Profile: rules.RepresentationProfile, Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: "record/cpp/1", Input: pattern, Expand: func(ctx context.Context, r g.Expansion[h.Language, s.Language, source.Reader]) (g.Expression[s.Language], error) {
			src, err := r.Reader()
			if err != nil {
				return g.Expression[s.Language]{}, err
			}
			use, err := r.Use()
			if err != nil {
				return g.Expression[s.Language]{}, err
			}
			if !use.Body {
				return g.Expression[s.Language]{}, src.Unsupported(use.Use, "nested records require a named declaration in this profile")
			}
			view, _, err := src.Descriptor(use.Use)
			if err != nil {
				return g.Expression[s.Language]{}, err
			}
			fields, err := src.Fields(use)
			if err != nil {
				return g.Expression[s.Language]{}, err
			}
			var members []g.Expression[s.Language]
			for _, f := range fields {

				ty, err := r.Mapped(f.Use)
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				if f.Optional {
					ty = model.External("std", "optional", ty)
				}
				members = append(members, g.At(model.Field(f.Name, f.Optional, ty), f.Origin))
			}
			return g.At(model.Record(members...), view.Origin), nil
		}}}}, nil
	}
}
