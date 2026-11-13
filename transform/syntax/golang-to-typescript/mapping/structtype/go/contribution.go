// Package structtype is the golang-to-typescript contribution for a Go struct
// type (docs/contribution-contract.md). It belongs to the two other pairs of
// bittype#196 (decision O43), which show that the contribution contract and
// its kernel assume no horos source. A consumer's composition registers it on
// its own, with the codec of its payload schema.
//
// It covers a struct type that is the underlying type of a defined type and
// whose fields are neither embedded nor tagged. It refuses, as unsupported:
// a struct type that is no body of a defined type, an embedded field, and a
// field with a tag. What a field's type may be is not this contribution's
// question: the planner selects a contribution for it.
package structtype

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/vocabulary/structtype/go"
	rules "github.com/Bitspark/bittype/transform/syntax/golang/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/golang/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[gs.Language, s.Language, source.Reader]
type decide = g.DecideStep[gs.Language, source.Reader]
type name = g.NameStep[gs.Language, source.Reader]

// target is the target that the contribution states (O4).
var target = string(s.Language{}.ID())

// The reasons of this contribution's refusals. Each is unsupported: a later
// slice of this contribution can admit the same input.
const (
	ReasonNestedStruct  = "a struct type that is no underlying type of a defined type is not covered"
	ReasonEmbeddedField = "an embedded field is not covered"
	ReasonFieldTag      = "a field with a tag is not covered: a tag may rename the member"
)

// Contribution returns the struct contribution: a Go struct type becomes a
// TypeScript model record with one member for each field, in the order of the
// source, named as the field is spelled. It records the member names in its
// payload. It declares at registration the fact that every selection of it is
// the model's record form, which its revision fixes (O23).
func Contribution(from representation.Language) (g.Contribution[gs.Language, s.Language, source.Reader], error) {
	selector, err := source.StructType(from)
	if err != nil {
		return contribution{}, err
	}
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.TypesProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
		Facts:   []g.Fact{model.EvidenceRecordForm.Proved("")},
		Rule:    g.Rule[gs.Language, source.Reader]{Name: nameFields, Decide: decideMembers},
		Project: project,
	}, nil
}

func refuse(reason string) error {
	return g.Refuse(g.RefusalUnsupported, vocabulary.Reference.Identity.Publisher, reason)
}

// nameFields is step 1: one reference for each field, in the order of the
// source, to what is selected for the field's type.
func nameFields(_ context.Context, step name) ([]g.Request[gs.Language], error) {
	use := step.Subject.Use
	if !use.Body {
		return nil, refuse(ReasonNestedStruct)
	}
	reader, err := step.Reader()
	if err != nil {
		return nil, err
	}
	fields, err := reader.Fields(use)
	if err != nil {
		return nil, err
	}
	requests := make([]g.Request[gs.Language], 0, len(fields))
	for _, f := range fields {
		switch {
		case f.Embedded:
			return nil, refuse(ReasonEmbeddedField)
		case f.Tag != "":
			return nil, refuse(ReasonFieldTag)
		}
		requests = append(requests, g.Reference(f.Use))
	}
	return requests, nil
}

// decideMembers is step 3: the names of the members, in order.
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
		payload.Members = append(payload.Members, f.Name)
	}
	return g.Decision{Payload: payload}, nil
}

// project builds the target record from the frozen decision and the views of
// the dependencies: member k has the k-th recorded name and the type selected
// for dependency k.
//
// COPY (AC-11, finding AG-3 of the two other pairs): this is the projector of
// the record contributions of the horos pairs without their delegated fields.
// It reads no source, so nothing in it is Go's: the member names come from
// the frozen payload and the types from the dependency views. It has no
// agnostic implementation, so this package repeats it. That the pair is
// green does not show that the mechanism has an owner.
func project(_ context.Context, p g.Projection[s.Language]) (g.Projected[s.Language], error) {
	payload, ok := p.Payload().(vocabulary.Payload)
	if !ok || len(payload.Members) != p.Dependencies() {
		return g.Projected[s.Language]{}, g.ErrInvalid
	}
	members := make([]g.Expression[s.Language], len(payload.Members))
	table := make([]g.Member, len(payload.Members))
	for k, member := range payload.Members {
		d := p.Dependency(k)
		members[k] = g.At(model.Field(member, false, d.Type), d.Origin)
		table[k] = g.Member{Name: member}
	}
	return g.Projected[s.Language]{Type: g.At(model.Record(members...), p.Origin()), Members: table}, nil
}
