// Package numbertype is the typescript-to-horos contribution for the TypeScript
// primitive type number (docs/contribution-contract.md). It belongs to the two
// other pairs of bittype#196 (decision O43). A consumer's composition
// registers it on its own, with the codec of its payload schema.
//
// It refuses: the pair covers no number.
package numbertype

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/vocabulary/numbertype/go"
	rules "github.com/Bitspark/bittype/transform/syntax/typescript/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/typescript/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[ts.Language, s.Language, source.Reader]
type decide = g.DecideStep[ts.Language, source.Reader]

// target is the target that the contribution states (O4).
var target = string(s.Language{}.ID())

// ReasonNumber is the reason of the refusal. It quotes the accepted page.
const ReasonNumber = "number is not covered: \"native numeric and nullable forms do not uniquely identify source domains or absence semantics\" (docs/generation.md), and no semantic choice for number is made"

// Contribution returns the number contribution. It answers for the TypeScript
// type number and refuses every selection as unsupported: a later slice can
// admit the same input once a canonicalization states what a number is (an
// integer, a binary64 value without NaN and the infinities, or another
// domain). It is registered so that the refusal is the pair's own statement
// and carries its reason. Without it a use of number would be refused by the
// kernel as a missing contribution, which says that nothing was built and not
// why.
func Contribution(from representation.Language) (g.Contribution[ts.Language, s.Language, source.Reader], error) {
	selector, err := source.Primitive(from, source.PrimitiveNumber)
	if err != nil {
		return contribution{}, err
	}
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.CanonicalizationProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
		Rule: g.Rule[ts.Language, source.Reader]{Decide: func(context.Context, decide) (g.Decision, error) {
			return g.Decision{}, g.Refuse(g.RefusalUnsupported, vocabulary.Reference.Identity.Publisher, ReasonNumber)
		}},
		// The contract requires a projector of every contribution. This one is
		// never reached: the rule refuses before anything is frozen.
		Project: func(context.Context, g.Projection[s.Language]) (g.Projected[s.Language], error) {
			return g.Projected[s.Language]{}, g.ErrInvalid
		},
	}, nil
}
