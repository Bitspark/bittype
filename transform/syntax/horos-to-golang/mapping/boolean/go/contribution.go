package boolean

// PROTOTYPE (bittype#196): this mapping under the contribution contract. The
// Library of the mapping contract stays beside it on this branch, as the
// baseline that the differential test compares with. The migration deletes it.

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/boolean/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[h.Language, s.Language, source.Reader]
type decide = g.DecideStep[h.Language, source.Reader]

// Target is the stated target of the contribution (O4).
var Target = string(s.Language{}.ID())

// Contribution returns the boolean contribution: the horos scalar bool
// becomes the native primitive PrimitiveBool. Nothing varies between selections,
// so the payload is empty and the native type is an invariant of the
// representation revision.
func Contribution(from representation.Language) (contribution, error) {
	selector, err := source.Builtin(from, "scalar", "bool")
	return contribution{
		Reference: vocabulary.Reference, Target: Target, Profile: rules.RepresentationProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
		Rule: g.Rule[h.Language, source.Reader]{Decide: func(context.Context, decide) (g.Decision, error) {
			return g.Decision{Payload: vocabulary.Payload{}}, nil
		}},
		Project: func(context.Context, g.Projection[s.Language]) (g.Projected[s.Language], error) {
			return g.Projected[s.Language]{Type: model.Primitive(model.PrimitiveBool)}, nil
		},
	}, err
}
