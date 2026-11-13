// Package stringtype is the golang-to-typescript contribution for the Go
// basic type string (docs/contribution-contract.md). It belongs to the two
// other pairs of bittype#196 (decision O43). A consumer's composition
// registers it on its own, with the codec of its payload schema.
package stringtype

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/vocabulary/stringtype/go"
	rules "github.com/Bitspark/bittype/transform/syntax/golang/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/golang/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[gs.Language, s.Language, source.Reader]
type decide = g.DecideStep[gs.Language, source.Reader]

// target is the target that the contribution states (O4).
var target = string(s.Language{}.ID())

// Contribution returns the string contribution: the Go type string becomes
// the TypeScript model primitive string. Nothing varies between two
// selections, so the payload is empty. The target type is an invariant that
// the representation revision fixes (row 5).
//
// COPY (AC-11, finding AG-1 of the two other pairs): a rule without requests
// that decides the empty payload and a projector that returns one constant
// target expression are the same in every contribution of a scalar, of every
// pair. Only the selector, the reference and the constant are this
// contribution's. The mechanism has no agnostic implementation, so this
// package repeats it. That the pair is green does not show that the
// mechanism has an owner.
func Contribution(from representation.Language) (g.Contribution[gs.Language, s.Language, source.Reader], error) {
	selector, err := source.Basic(from, source.BasicString)
	if err != nil {
		return contribution{}, err
	}
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.TypesProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
		Rule: g.Rule[gs.Language, source.Reader]{Decide: func(context.Context, decide) (g.Decision, error) {
			return g.Decision{Payload: vocabulary.Payload{}}, nil
		}},
		Project: func(context.Context, g.Projection[s.Language]) (g.Projected[s.Language], error) {
			return g.Projected[s.Language]{Type: model.Primitive(model.PrimitiveString)}, nil
		},
	}, nil
}
