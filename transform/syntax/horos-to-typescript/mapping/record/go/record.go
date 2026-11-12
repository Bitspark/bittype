// Package record is the horos-to-typescript mapping library for closed horos
// records. A consumer registers it on its own; it imports no other mapping
// library.
package record

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Library binds a closed record to a TypeScript record of its fields under the
// binding ID "record/typescript/1". The kernel resolves each field's mapping.
func Library(profile string) rules.Library[s.Language] {
	return rules.Library[s.Language]{Profile: profile, DescribeRules: func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := patterns.Head(from, "record", "closed")
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: "record/typescript/1", Input: pattern, Expand: func(ctx context.Context, r g.Expansion[h.Language, s.Language, source.Reader]) (g.Expression[s.Language], error) {
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

				members = append(members, g.At(model.Field(f.Name, f.Optional, ty), f.Origin))
			}
			return g.At(model.Record(members...), view.Origin), nil
		}}}}, nil
	}}
}
