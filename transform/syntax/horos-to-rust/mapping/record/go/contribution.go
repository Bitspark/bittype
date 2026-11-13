// Package record is the builtin horos-to-rust contribution for closed
// horos records (docs/contribution-contract.md). A consumer's composition
// registers it on its own, with the codec of its payload schema. It delegates
// the presence of an optional field to whichever contribution the composition
// registers for that role, and imports none.
package record

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/vocabulary/record/go"
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

// Contribution returns the record contribution: a closed horos record becomes
// a Rust model record. It lays out its required fields itself and records
// their member names, in order, in its payload. For an optional field it
// delegates presence (row 11, O38): it requests the field position with the
// role that the source issued and states the member's name in the request
// (O3, O29). It reads no payload of the contribution that answers. It declares
// at registration the fact that every selection of it is the model's record
// form, which its revision fixes (O23).
func Contribution(from representation.Language) (g.Contribution[h.Language, s.Language, source.Reader], error) {
	selector, err := source.Builtin(from, "record", h.RecordClosed)
	if err != nil {
		return contribution{}, err
	}
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.RepresentationProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
		Facts:   []g.Fact{model.EvidenceRecordForm.Proved("")},
		Rule:    g.Rule[h.Language, source.Reader]{Name: nameFields, Decide: decideMembers},
		Project: project,
	}, nil
}

// nameFields is step 1: one request for each field, in the order of the
// source. A required field is a reference to what is selected for its type.
// An optional field is a field position.
func nameFields(_ context.Context, step name) ([]g.Request[h.Language], error) {
	use := step.Subject.Use
	if !use.Body {
		return nil, g.Refuse(g.RefusalUnsupported, vocabulary.Reference.Identity.Publisher, rules.ReasonNestedRecord)
	}
	reader, err := step.Reader()
	if err != nil {
		return nil, err
	}
	fields, err := reader.Fields(use)
	if err != nil {
		return nil, err
	}
	requests := make([]g.Request[h.Language], 0, len(fields))
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
}

// decideMembers is step 3: the names of the members that the record lays out
// itself, in order. The name of a delegated field is in its request.
func decideMembers(_ context.Context, step decide) (g.Decision, error) {
	reader, err := step.Reader()
	if err != nil {
		return g.Decision{}, err
	}
	fields, err := reader.Fields(step.Subject.Use)
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
}

// project builds the native record from the frozen decision and the views of
// the dependencies. A dependency with a fragment is a delegated field: its
// members are the fragment. Every other dependency is a field that the record
// lays out itself, named by the next member of the payload.
func project(_ context.Context, p g.Projection[s.Language]) (g.Projected[s.Language], error) {
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
	if own != len(payload.Members) {
		return g.Projected[s.Language]{}, g.ErrInvalid
	}
	return g.Projected[s.Language]{Type: g.At(model.Record(members...), p.Origin()), Members: table}, nil
}

// memberName is the spelling of a member: the field's name as the source
// states it.
func memberName(field string) (string, error) { return field, nil }
