package list

// PROTOTYPE (bittype#196): this mapping under the contribution contract. The
// Library of the mapping contract stays beside it on this branch, as the
// baseline that the differential test compares with. The migration deletes it.

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/list/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[h.Language, s.Language, source.Reader]
type decide = g.DecideStep[h.Language, source.Reader]
type name = g.NameStep[h.Language, source.Reader]

// Target is the stated target of the contribution (O4).
var Target = string(s.Language{}.ID())

// Contribution returns the list contribution: a horos list-of becomes a native
// collection of what was selected for its element. The element is a link of
// the frozen record. Nothing else varies, so the payload is empty.
func Contribution(from representation.Language) (contribution, error) {
	selector, err := source.Builtin(from, "list-of", "")
	return contribution{
		Reference: vocabulary.Reference, Target: Target, Profile: rules.RepresentationProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
		Rule: g.Rule[h.Language, source.Reader]{
			Name: func(_ context.Context, step name) ([]g.Request[h.Language], error) {
				element, err := step.Reader.Element(step.Subject.Use)
				return []g.Request[h.Language]{g.Reference(element)}, err
			},
			Decide: func(context.Context, decide) (g.Decision, error) {
				return g.Decision{Payload: vocabulary.Payload{}}, nil
			},
		},
		Project: func(_ context.Context, p g.Projection[s.Language]) (g.Projected[s.Language], error) {
			return g.Projected[s.Language]{Type: model.Collection(p.Dependency(0).Type)}, nil
		},
	}, err
}
