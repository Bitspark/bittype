// Package optional is the builtin horos-to-typescript contribution for the
// role "optional record field" (docs/contribution-contract.md, row 11, O38).
// A consumer's composition registers it on its own, with the codec of its
// payload schema. A record contribution delegates the presence of an optional
// field to it through the planner. It imports no record contribution.
package optional

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/optional/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[h.Language, s.Language, source.Reader]
type decide = g.DecideStep[h.Language, source.Reader]
type name = g.NameStep[h.Language, source.Reader]

// target is the target that the contribution states (O4).
var target = string(s.Language{}.ID())

// Contribution returns the optional contribution in the delegated form (row
// 11, O38). Its domain is the role that the horos source issues for an
// optional field of a record. It answers a field position with a member
// fragment in the member vocabulary of the TypeScript model: one field, named as the
// request states, whose type is the operand's type and which carries the optional flag of the field form, so that the property is optional. The operand is a link of the frozen record, so the
// payload is empty. It has no standalone value adapter (row 14).
func Contribution(representation.Language) (g.Contribution[h.Language, s.Language, source.Reader], error) {
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.RepresentationProfile,
		Domain: g.ForRole(source.OptionalFieldRole), Schema: vocabulary.Schema,
		Rule: g.Rule[h.Language, source.Reader]{
			Name: func(_ context.Context, step name) ([]g.Request[h.Language], error) {
				return []g.Request[h.Language]{g.Reference(step.Subject.Field.Operand)}, nil
			},
			Decide: func(context.Context, decide) (g.Decision, error) {
				return g.Decision{Payload: vocabulary.Payload{}}, nil
			},
		},
		Project: func(_ context.Context, p g.Projection[s.Language]) (g.Projected[s.Language], error) {
			if p.Dependencies() != 1 {
				return g.Projected[s.Language]{}, g.ErrScope
			}
			operand := p.Dependency(0)
			wrapped := operand.Type
			return g.Projected[s.Language]{Fragment: []g.Expression[s.Language]{g.At(model.Field(p.Stated(), true, wrapped), operand.Origin)}, Members: []g.Member{{Name: p.Stated()}}}, nil
		},
	}, nil
}
