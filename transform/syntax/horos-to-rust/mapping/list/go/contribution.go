// Package list is the builtin horos-to-rust contribution for horos
// lists (docs/contribution-contract.md). A consumer's composition registers it
// on its own, with the codec of its payload schema.
package list

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/vocabulary/list/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/rust/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[h.Language, s.Language, source.Reader]
type decide = g.DecideStep[h.Language, source.Reader]
type name = g.NameStep[h.Language, source.Reader]

// target is the target that the contribution states (O4).
var target = string(s.Language{}.ID())

// Contribution returns the list contribution: a horos list becomes the
// Rust model collection of what was selected for its element. The element is
// a link of the frozen record. Nothing else varies, so the payload is empty.
func Contribution(from representation.Language) (g.Contribution[h.Language, s.Language, source.Reader], error) {
	selector, err := source.Builtin(from, "list-of", "")
	if err != nil {
		return contribution{}, err
	}
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.RepresentationProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
		Rule: g.Rule[h.Language, source.Reader]{
			Name: func(_ context.Context, step name) ([]g.Request[h.Language], error) {
				reader, err := step.Reader()
				if err != nil {
					return nil, err
				}
				element, err := reader.Element(step.Subject.Use)
				if err != nil {
					return nil, err
				}
				return []g.Request[h.Language]{g.Reference(element)}, nil
			},
			Decide: func(context.Context, decide) (g.Decision, error) {
				return g.Decision{Payload: vocabulary.Payload{}}, nil
			},
		},
		Project: func(_ context.Context, p g.Projection[s.Language]) (g.Projected[s.Language], error) {
			if p.Dependencies() != 1 {
				return g.Projected[s.Language]{}, g.ErrScope
			}
			return g.Projected[s.Language]{Type: model.Collection(p.Dependency(0).Type)}, nil
		},
	}, nil
}
