// Package record is the builtin horos-to-golang contribution for closed
// horos records (docs/contribution-contract.md). A consumer's composition
// registers it on its own, with the codec of its payload schema. It delegates
// the presence of an optional field to whichever contribution the composition
// registers for that role, and imports none.
package record

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/record/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	naming "github.com/Bitspark/bittype/transform/syntax/horos/mapping/naming/go"
	policy "github.com/Bitspark/bittype/transform/syntax/horos/mapping/record/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	shared "github.com/Bitspark/bittype/transform/syntax/shared/contribution/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[h.Language, s.Language, source.Reader]
type decide = g.DecideStep[h.Language, source.Reader]
type name = g.NameStep[h.Language, source.Reader]

// target is the target that the contribution states (O4).
var target = string(s.Language{}.ID())

// publisher qualifies the reason of a refusal of this contribution.
var publisher = vocabulary.Reference.Identity.Publisher

// Contribution returns the record contribution: a closed horos record becomes
// a Go model record, by the horos record policy. What this target supplies
// is the spelling of a member, the exported native identifier of the field's name, the
// model's field and record forms, and its payload type, which records the
// names of the members that the record lays out itself. It declares at
// registration the fact that every selection of it is the model's record
// form, which its revision fixes (O23).
func Contribution(from representation.Language) (g.Contribution[h.Language, s.Language, source.Reader], error) {
	selector, err := source.Builtin(from, "record", h.RecordClosed)
	if err != nil {
		return contribution{}, err
	}
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.RepresentationProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
		Facts: []g.Fact{model.EvidenceRecordForm.Proved("")},
		Rule: g.Rule[h.Language, source.Reader]{
			Name: func(_ context.Context, step name) ([]g.Request[h.Language], error) {
				return policy.Requests(step, publisher, memberName)
			},
			Decide: func(_ context.Context, step decide) (g.Decision, error) {
				members, err := policy.Members(step, memberName)
				if err != nil {
					return g.Decision{}, err
				}
				return g.Decision{Payload: vocabulary.Payload{Members: members}}, nil
			},
		},
		Project: func(_ context.Context, p g.Projection[s.Language]) (g.Projected[s.Language], error) {
			payload, ok := p.Payload().(vocabulary.Payload)
			if !ok {
				return g.Projected[s.Language]{}, g.ErrInvalid
			}
			return shared.Members(p, payload.Members, field, model.Record)
		},
	}, nil
}

// field is a member that the record lays out itself: a required field of the
// model.
func field(name string, ty g.Expression[s.Language]) g.Expression[s.Language] {
	return model.Field(name, false, ty)
}

// memberName is the spelling of a member: the exported native identifier of
// the field's name. A name that the naming rule does not cover is unsupported:
// a later slice of that rule can admit it.
func memberName(field string) (string, error) {
	member, err := naming.Exported(field)
	if err != nil {
		return "", g.Refuse(g.RefusalUnsupported, publisher, err.Error())
	}
	return member, nil
}
