package declarations

// PROTOTYPE (bittype#196): this mapping under the contribution contract. The
// Library of the mapping contract stays beside it on this branch, as the
// baseline that the differential test compares with. The migration deletes it.

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/declarations/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[h.Language, s.Language, source.Reader]
type decide = g.DecideStep[h.Language, source.Reader]
type name = g.NameStep[h.Language, source.Reader]

// Target is the stated target of the contribution (O4).
var Target = string(s.Language{}.ID())

// Contribution returns the declaration contribution: the default native
// declaration policy for every named declaration without an exact mapping. It
// applies the horos declaration policy with the two native forms of this
// target and records the policy's decision in this target's payload type.
func Contribution(representation.Language) (contribution, error) {
	return contribution{
		Reference: vocabulary.Reference, Target: Target, Profile: rules.RepresentationProfile,
		Domain: g.ForDeclarationPolicy(), Schema: vocabulary.Schema,
		Rule: g.Rule[h.Language, source.Reader]{
			Name: func(_ context.Context, step name) ([]g.Request[h.Language], error) {
				return recipe.Requests(step, model.EvidenceRecordForm), nil
			},
			Decide: func(_ context.Context, step decide) (g.Decision, error) {
				d, err := recipe.Decide(step, model.EvidenceRecordForm, vocabulary.Reference.Identity.Publisher)
				return g.Decision{Payload: vocabulary.Payload{Name: d.Name, Nominal: d.Nominal, Parameters: d.Parameters}}, err
			},
		},
		Declare: func(_ context.Context, p g.Projection[s.Language]) (g.DeclarationDescription[s.Language], error) {
			payload, ok := p.Payload().(vocabulary.Payload)
			if !ok {
				return g.DeclarationDescription[s.Language]{}, g.ErrInvalid
			}
			return recipe.Describe(p, recipe.Decision{Name: payload.Name, Nominal: payload.Nominal, Parameters: payload.Parameters}, model.Form(model.FormInterface), model.Form(model.FormAlias))
		},
	}, nil
}
