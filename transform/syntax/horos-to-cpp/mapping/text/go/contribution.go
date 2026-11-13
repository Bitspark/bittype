package text

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/cpp/model/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/vocabulary/text/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[h.Language, s.Language, source.Reader]
type decide = g.DecideStep[h.Language, source.Reader]

// target is the target that the contribution states (O4).
var target = string(s.Language{}.ID())

// Contribution returns the text contribution: the horos scalar utf8-text
// becomes the C++ model primitive string. Nothing varies between two
// selections, so the payload is empty. The native type is an invariant that
// the representation revision fixes (row 5).
func Contribution(from representation.Language) (g.Contribution[h.Language, s.Language, source.Reader], error) {
	selector, err := source.Builtin(from, "scalar", "utf8-text")
	if err != nil {
		return contribution{}, err
	}
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.RepresentationProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
		Rule: g.Rule[h.Language, source.Reader]{Decide: func(context.Context, decide) (g.Decision, error) {
			return g.Decision{Payload: vocabulary.Payload{}}, nil
		}},
		Project: func(context.Context, g.Projection[s.Language]) (g.Projected[s.Language], error) {
			return g.Projected[s.Language]{Type: model.Primitive(model.PrimitiveString)}, nil
		},
	}, nil
}
