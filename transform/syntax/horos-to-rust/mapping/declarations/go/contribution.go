// Package declarations is the builtin horos-to-rust declaration
// contribution: the default native declaration policy for every named horos
// declaration without an exact mapping (docs/contribution-contract.md, row 2,
// O12). A consumer's composition registers it on its own, with the codec of
// its payload schema.
package declarations

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/vocabulary/declarations/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
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

// Contribution returns the declaration contribution. It applies the horos
// declaration policy with the two native forms of this target, struct for a
// body whose selected representation is the model's record form and alias for
// every other body, and records the policy's decision in this target's payload
// type. It asks for the record form as evidence about the body and reads no
// payload of another contribution (row 3). Where it decides an alias it states
// that the declaration is transparent (O35).
func Contribution(representation.Language) (g.Contribution[h.Language, s.Language, source.Reader], error) {
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.RepresentationProfile,
		Domain: g.ForDeclarationPolicy(), Schema: vocabulary.Schema,
		Rule: g.Rule[h.Language, source.Reader]{
			Name: func(_ context.Context, step name) ([]g.Request[h.Language], error) {
				return recipe.Requests(step, model.EvidenceRecordForm), nil
			},
			Decide: func(_ context.Context, step decide) (g.Decision, error) {
				d, err := recipe.Decide(step, model.EvidenceRecordForm, vocabulary.Reference.Identity.Publisher)
				if err != nil {
					return g.Decision{}, err
				}
				return g.Decision{Payload: vocabulary.Payload{Name: d.Name, Nominal: d.Nominal, Parameters: d.Parameters}, Transparent: !d.Nominal}, nil
			},
		},
		Declare: func(_ context.Context, p g.Projection[s.Language]) (g.DeclarationDescription[s.Language], error) {
			payload, ok := p.Payload().(vocabulary.Payload)
			if !ok {
				return g.DeclarationDescription[s.Language]{}, g.ErrInvalid
			}
			return recipe.Describe(p, recipe.Decision{Name: payload.Name, Nominal: payload.Nominal, Parameters: payload.Parameters}, model.Form(model.FormStruct), model.Form(model.FormAlias))
		},
	}, nil
}
