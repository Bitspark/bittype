package optional

// PROTOTYPE (bittype#196): this mapping under the contribution contract. The
// Library of the mapping contract stays beside it on this branch, as the
// baseline that the differential test compares with. The migration deletes it.

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/optional/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[h.Language, s.Language, source.Reader]
type decide = g.DecideStep[h.Language, source.Reader]
type name = g.NameStep[h.Language, source.Reader]

// Target is the stated target of the contribution (O4).
var Target = string(s.Language{}.ID())

// Contribution returns the optional contribution in the delegated form (row
// 11, D3). Its domain is the role "optional record field". It answers a field
// position with a member fragment in the member vocabulary of the Go model:
// one field, named as the request states, whose type is the operand's type in
// the helper BittypeOptional. It owns that helper: the key is this
// contribution's reference and the local name, and no consumer creates a key.
// It has no standalone value adapter.
func Contribution(representation.Language) (contribution, error) {
	return contribution{
		Reference: vocabulary.Reference, Target: Target, Profile: rules.RepresentationProfile,
		Domain: g.ForRole(source.OptionalFieldRole), Schema: vocabulary.Schema,
		Rule: g.Rule[h.Language, source.Reader]{
			Name: func(_ context.Context, step name) ([]g.Request[h.Language], error) {
				return []g.Request[h.Language]{g.Reference(step.Subject.Field.Operand)}, nil
			},
			Decide: func(context.Context, decide) (g.Decision, error) {
				return g.Decision{Payload: vocabulary.Payload{}}, nil
			},
		},
		Owned: []g.OwnedDescription[s.Language]{{Name: Helper, Declaration: g.NativeDeclaration[s.Language]{
			Name: "BittypeOptional", Form: model.Form(model.FormDefined), Parameters: []g.NativeParameter{{Name: "T"}},
			Body: model.Record(model.Field("Present", false, model.Primitive(model.PrimitiveBool)), model.Field("Value", false, g.Parameter[s.Language](0))),
		}}},
		Project: func(_ context.Context, p g.Projection[s.Language]) (g.Projected[s.Language], error) {
			operand := p.Dependency(0)
			wrapped, err := p.Owned(Helper, operand.Type)
			if err != nil {
				return g.Projected[s.Language]{}, err
			}
			return g.Projected[s.Language]{Fragment: []g.Expression[s.Language]{g.At(model.Field(p.Stated(), true, wrapped), operand.Origin)}, Members: []g.Member{{Name: p.Stated()}}}, nil
		},
	}, nil
}

// Helper is the local name of the declaration this contribution owns.
const Helper = "optional"
