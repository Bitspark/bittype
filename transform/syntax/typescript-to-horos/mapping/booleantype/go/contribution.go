// Package booleantype is the typescript-to-horos contribution for the TypeScript
// primitive type boolean (docs/contribution-contract.md). It belongs to the two
// other pairs of bittype#196 (decision O43). A consumer's composition
// registers it on its own, with the codec of its payload schema.
package booleantype

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	model "github.com/Bitspark/bittype/transform/syntax/horos/model/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/vocabulary/booleantype/go"
	rules "github.com/Bitspark/bittype/transform/syntax/typescript/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/typescript/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[ts.Language, s.Language, source.Reader]
type decide = g.DecideStep[ts.Language, source.Reader]

// target is the target that the contribution states (O4).
var target = string(s.Language{}.ID())

// Contribution returns the boolean contribution: the TypeScript type boolean
// becomes the horos scalar bool. Nothing varies between two selections, so
// the payload is empty. The target type is an invariant that the
// representation revision fixes (row 5).
//
// COPY (AC-11, finding AG-1 of the two other pairs): a rule without requests
// that decides the empty payload and a projector that returns one constant
// target expression are the same in every contribution of a scalar, of every
// pair. Only the selector, the reference and the constant are this
// contribution's. The mechanism has no agnostic implementation, so this
// package repeats it. That the pair is green does not show that the
// mechanism has an owner.
func Contribution(from representation.Language) (g.Contribution[ts.Language, s.Language, source.Reader], error) {
	selector, err := source.Primitive(from, source.PrimitiveBoolean)
	if err != nil {
		return contribution{}, err
	}
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.CanonicalizationProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
		Rule: g.Rule[ts.Language, source.Reader]{Decide: func(context.Context, decide) (g.Decision, error) {
			return g.Decision{Payload: vocabulary.Payload{}}, nil
		}},
		Project: func(context.Context, g.Projection[s.Language]) (g.Projected[s.Language], error) {
			return g.Projected[s.Language]{Type: model.Scalar(model.ScalarBool)}, nil
		},
	}, nil
}
