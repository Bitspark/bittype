// Package record holds the policy by which a closed horos record becomes a
// record of a target: which requests its rule names and what it records in
// its payload. The projector that pairs the two reads no source and is the
// shared mechanism of contributions
// (docs/contribution-contract.md, rows 5, 11 to 13, O3, O29, O38). The horos
// side of a mapping owns this policy, and it is the same for every target. A
// target's record package supplies the spelling of a member, its model's
// field and record forms and its payload type, and it is the contribution.
//
// The policy lays out a required field itself and delegates the presence of
// an optional field: it requests the field position with the role that the
// horos source issued and states the member's name in the request. It reads
// no payload of the contribution that answers, and it states no profile.
package record

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// ReasonNested is the reason for refusing a record that is no body of a named
// declaration: "Records must be named declarations"
// (docs/semantic-generation.md). The refusal is unsupported.
const ReasonNested = "nested records require a named declaration in this profile"

// Spelling is a target's spelling of a member for a field's name. A target
// that keeps the name as the source states it passes Verbatim. A spelling
// that does not cover a name refuses it.
type Spelling func(field string) (string, error)

// Verbatim spells a member as the source names the field.
func Verbatim(field string) (string, error) { return field, nil }

// Requests is step 1: one request for each field, in the order of the source.
// A required field is a reference to what is selected for its type. An
// optional field is a field position, with the member's name stated in the
// request. owner qualifies the reason of the refusal of a nested record.
func Requests(step g.NameStep[h.Language, source.Reader], owner string, spell Spelling) ([]g.Request[h.Language], error) {
	use := step.Subject.Use
	if !use.Body {
		return nil, g.Refuse(g.RefusalUnsupported, owner, ReasonNested)
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
		member, err := spell(f.Name)
		if err != nil {
			return nil, err
		}
		requests = append(requests, g.Field(f.Use, f.Role, member))
	}
	return requests, nil
}

// Members is step 3: the names of the members that the record lays out
// itself, in order. Each target's package records them in the payload type of
// its own schema. The name of a delegated field is in its request.
func Members(step g.DecideStep[h.Language, source.Reader], spell Spelling) ([]string, error) {
	reader, err := step.Reader()
	if err != nil {
		return nil, err
	}
	fields, err := reader.Fields(step.Subject.Use)
	if err != nil {
		return nil, err
	}
	var members []string
	for _, f := range fields {
		if f.Role.Valid() {
			continue
		}
		member, err := spell(f.Name)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}
