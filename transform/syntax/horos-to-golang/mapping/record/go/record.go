// Package record is the horos-to-golang mapping library for records: a closed
// horos record becomes a Go record of its mapped fields. The type of an
// optional field is wrapped in a helper that the consumer names. This package
// neither creates the key of that helper nor supplies its declaration.
package record

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	naming "github.com/Bitspark/bittype/transform/syntax/horos/mapping/naming/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Library returns the record mapping. It states the structural representation
// profile itself. optional is the
// key of the one-parameter helper that wraps the type of an optional field.
// The consumer creates the key and registers a library that supplies a helper
// under it.
//
// An invalid key makes Describe fail with generation.ErrInvalid, so the library
// cannot be registered. A valid key that no registered library supplies makes
// planning fail with generation.ErrUnsupported at the first optional field.
// An optional field is never emitted unwrapped.
func Library(optional g.HelperKey) rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		if !optional.Valid() {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, g.ErrInvalid
		}
		pattern, err := patterns.Head(from, "record", h.RecordClosed)
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Profile: rules.RepresentationProfile, Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: "record/golang/1", Input: pattern, Expand: func(ctx context.Context, r g.Expansion[h.Language, s.Language, source.Reader]) (g.Expression[s.Language], error) {
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
				f.Name, err = naming.Exported(f.Name)
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				ty, err := r.Mapped(f.Use)
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				if f.Optional {
					ty = g.Helper(optional, ty)
				}
				members = append(members, g.At(model.Field(f.Name, f.Optional, ty), f.Origin))
			}
			return g.At(model.Record(members...), view.Origin), nil
		}}}}, nil
	}
}
