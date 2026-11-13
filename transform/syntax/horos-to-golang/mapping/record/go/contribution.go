package record

// PROTOTYPE (bittype#196): this mapping under the contribution contract. The
// Library of the mapping contract stays beside it on this branch, as the
// baseline that the differential test compares with. The migration deletes it.

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/record/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	naming "github.com/Bitspark/bittype/transform/syntax/horos/mapping/naming/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[h.Language, s.Language, source.Reader]
type decide = g.DecideStep[h.Language, source.Reader]
type name = g.NameStep[h.Language, source.Reader]

// Target is the stated target of the contribution (O4).
var Target = string(s.Language{}.ID())

// Contribution returns the record contribution: a closed horos record becomes
// a native record. It lays out its required fields itself. For an optional
// field it delegates presence (row 11, D3): it requests the field position
// with the role that the source issued and states the member's name in the
// request (O3). It imports no optional contribution and reads no payload of
// one. It declares at registration the fact that every selection of it is the
// model's record form, which its revision fixes (O23).
func Contribution(from representation.Language) (contribution, error) {
	selector, err := source.Builtin(from, "record", "closed")
	return contribution{
		Reference: vocabulary.Reference, Target: Target, Profile: rules.RepresentationProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
		Facts: []g.Fact{model.EvidenceRecordForm.Proved("")},
		Rule: g.Rule[h.Language, source.Reader]{
			Name: func(_ context.Context, step name) ([]g.Request[h.Language], error) {
				use := step.Subject.Use
				if !use.Body {
					return nil, g.Refuse(g.RefusalByDesign, vocabulary.Reference.Identity.Publisher, "nested records require a named declaration in this profile")
				}
				fields, err := step.Reader.Fields(use)
				if err != nil {
					return nil, err
				}
				var requests []g.Request[h.Language]
				for _, f := range fields {
					if !f.Role.Valid() {
						requests = append(requests, g.Reference(f.Use))
						continue
					}
					member, err := memberName(f.Name)
					if err != nil {
						return nil, err
					}
					requests = append(requests, g.Field(f.Use, f.Role, member))
				}
				return requests, nil
			},
			Decide: func(_ context.Context, step decide) (g.Decision, error) {
				fields, err := step.Reader.Fields(step.Subject.Use)
				if err != nil {
					return g.Decision{}, err
				}
				payload := vocabulary.Payload{}
				for _, f := range fields {
					if f.Role.Valid() {
						continue
					}
					member, err := memberName(f.Name)
					if err != nil {
						return g.Decision{}, err
					}
					payload.Members = append(payload.Members, member)
				}
				return g.Decision{Payload: payload}, nil
			},
		},
		Project: func(_ context.Context, p g.Projection[s.Language]) (g.Projected[s.Language], error) {
			payload, ok := p.Payload().(vocabulary.Payload)
			if !ok {
				return g.Projected[s.Language]{}, g.ErrInvalid
			}
			var members []g.Expression[s.Language]
			var table []g.Member
			own := 0
			for i := 0; i < p.Dependencies(); i++ {
				d := p.Dependency(i)
				if len(d.Fragment) > 0 {
					members = append(members, d.Fragment...)
					table = append(table, g.FragmentMembers(i))
					continue
				}
				if own >= len(payload.Members) {
					return g.Projected[s.Language]{}, g.ErrInvalid
				}
				members = append(members, g.At(model.Field(payload.Members[own], false, d.Type), d.Origin))
				table = append(table, g.Member{Name: payload.Members[own]})
				own++
			}
			return g.Projected[s.Language]{Type: g.At(model.Record(members...), p.Origin()), Members: table}, nil
		},
	}, err
}

// memberName is the spelling of a member: the exported native identifier.
func memberName(field string) (string, error) {
	member, err := naming.Exported(field)
	if err != nil {
		return "", g.Refuse(g.RefusalUnsupported, vocabulary.Reference.Identity.Publisher, err.Error())
	}
	return member, nil
}
