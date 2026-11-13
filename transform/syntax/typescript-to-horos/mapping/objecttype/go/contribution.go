// Package objecttype is the typescript-to-horos contribution for a TypeScript
// object type (docs/contribution-contract.md). It belongs to the two other
// pairs of bittype#196 (decision O43), which show that the contribution
// contract and its kernel assume no native target. A consumer's composition
// registers it on its own, with the codec of its payload schema.
//
// It covers an object type that an interface or a type alias declares and
// whose properties are all required and none readonly. It refuses, as
// unsupported: an object type that is no body of a declaration, an optional
// property, and a readonly property. What a property's type may be is not
// this contribution's question: the planner selects a contribution for it.
package objecttype

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	model "github.com/Bitspark/bittype/transform/syntax/horos/model/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/vocabulary/objecttype/go"
	rules "github.com/Bitspark/bittype/transform/syntax/typescript/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/typescript/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[ts.Language, s.Language, source.Reader]
type decide = g.DecideStep[ts.Language, source.Reader]
type name = g.NameStep[ts.Language, source.Reader]

// target is the target that the contribution states (O4).
var target = string(s.Language{}.ID())

// The reasons of this contribution's refusals. Each is unsupported: a later
// slice of this contribution can admit the same input.
const (
	ReasonNestedObject     = "an object type that no interface and no type alias declares is not covered"
	ReasonOptionalProperty = "an optional property is not covered: a canonicalization has to choose what absence means, and none is chosen"
	ReasonReadonlyProperty = "a readonly property is not covered"
)

// Contribution returns the object contribution: a TypeScript object type
// becomes a closed horos record with one required field for each property,
// in the order of the source, named as the property is spelled. It records
// the field names in its payload.
func Contribution(from representation.Language) (g.Contribution[ts.Language, s.Language, source.Reader], error) {
	selector, err := source.ObjectType(from)
	if err != nil {
		return contribution{}, err
	}
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.CanonicalizationProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
		Rule:    g.Rule[ts.Language, source.Reader]{Name: nameProperties, Decide: decideMembers},
		Project: project,
	}, nil
}

func refuse(reason string) error {
	return g.Refuse(g.RefusalUnsupported, vocabulary.Reference.Identity.Publisher, reason)
}

// nameProperties is step 1: one reference for each property, in the order of
// the source, to what is selected for the property's type.
func nameProperties(_ context.Context, step name) ([]g.Request[ts.Language], error) {
	use := step.Subject.Use
	if !use.Body {
		return nil, refuse(ReasonNestedObject)
	}
	reader, err := step.Reader()
	if err != nil {
		return nil, err
	}
	properties, err := reader.Properties(use)
	if err != nil {
		return nil, err
	}
	requests := make([]g.Request[ts.Language], 0, len(properties))
	for _, p := range properties {
		switch {
		case p.Optional:
			return nil, refuse(ReasonOptionalProperty)
		case p.Readonly:
			return nil, refuse(ReasonReadonlyProperty)
		}
		requests = append(requests, g.Reference(p.Use))
	}
	return requests, nil
}

// decideMembers is step 3: the names of the fields, in order.
func decideMembers(_ context.Context, step decide) (g.Decision, error) {
	reader, err := step.Reader()
	if err != nil {
		return g.Decision{}, err
	}
	properties, err := reader.Properties(step.Subject.Use)
	if err != nil {
		return g.Decision{}, err
	}
	payload := vocabulary.Payload{}
	for _, p := range properties {
		payload.Members = append(payload.Members, p.Name)
	}
	return g.Decision{Payload: payload}, nil
}

// project builds the horos record from the frozen decision and the views of
// the dependencies: field k has the k-th recorded name and the type selected
// for dependency k.
//
// COPY (AC-11, finding AG-3 of the two other pairs): this is the projector of
// the struct contribution of golang-to-typescript and of the record
// contributions of the horos pairs without their delegated fields, with the
// target model's field and record constructors substituted. It reads no
// source. It has no agnostic implementation, so this package repeats it.
// That the pair is green does not show that the mechanism has an owner.
func project(_ context.Context, p g.Projection[s.Language]) (g.Projected[s.Language], error) {
	payload, ok := p.Payload().(vocabulary.Payload)
	if !ok || len(payload.Members) != p.Dependencies() {
		return g.Projected[s.Language]{}, g.ErrInvalid
	}
	fields := make([]g.Expression[s.Language], len(payload.Members))
	table := make([]g.Member, len(payload.Members))
	for k, member := range payload.Members {
		d := p.Dependency(k)
		fields[k] = g.At(model.Field(member, d.Type), d.Origin)
		table[k] = g.Member{Name: member}
	}
	return g.Projected[s.Language]{Type: g.At(model.Record(fields...), p.Origin()), Members: table}, nil
}
