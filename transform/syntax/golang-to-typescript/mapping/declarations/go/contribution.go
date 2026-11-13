// Package declarations is the golang-to-typescript declaration contribution:
// the default declaration policy for every Go defined type without an exact
// mapping (docs/contribution-contract.md, row 2, O12). It belongs to the two
// other pairs of bittype#196 (decision O43). A consumer's composition
// registers it on its own, with the codec of its payload schema.
//
// It covers a defined type whose body is selected as the TypeScript model's
// record form, and it makes a TypeScript interface of it, named as Go spells
// the type. It refuses a declaration whose body is not proved to be that
// form: as requirement unresolved where nothing is stated, as requirement
// unsatisfied where the form is refuted. It states no generic convention:
// the Go source reader reads no generic declaration.
//
// Not covered and not refused here: a Go name that is no TypeScript type name
// (for example a type named number). The TypeScript checker refuses the
// output when it is sealed.
package declarations

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/vocabulary/declarations/go"
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

// ReasonNoRecordForm is the reason of this contribution's refusal of a
// declaration whose body is not proved to be the model's record form.
const ReasonNoRecordForm = "a defined type becomes an interface, which needs a body in the record form of the TypeScript model"

// Contribution returns the declaration contribution. It asks for the record
// form as evidence about the declaration's body, through the evidence kind of
// the target's model, and reads no payload of another contribution (row 3).
// The declaration it decides is a defined declaration of the target, so it is
// not transparent (O35).
func Contribution(representation.Language) (g.Contribution[gs.Language, s.Language, source.Reader], error) {
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.TypesProfile,
		Domain: g.ForDeclarationPolicy(), Schema: vocabulary.Schema,
		Rule: g.Rule[gs.Language, source.Reader]{
			Name: func(_ context.Context, step name) ([]g.Request[gs.Language], error) {
				return []g.Request[gs.Language]{g.Evidence(step.Subject.Declaration.Body, g.Demand{Kind: model.EvidenceRecordForm})}, nil
			},
			Decide: decideName,
		},
		Declare: func(_ context.Context, p g.Projection[s.Language]) (g.DeclarationDescription[s.Language], error) {
			payload, ok := p.Payload().(vocabulary.Payload)
			if !ok || p.Dependencies() != 1 || p.Binders() != 0 {
				return g.DeclarationDescription[s.Language]{}, g.ErrInvalid
			}
			return g.DeclarationDescription[s.Language]{Name: payload.Name, Form: model.Form(model.FormInterface), Origin: p.Origin(), Body: p.Dependency(0).Type}, nil
		},
	}, nil
}

// decideName is step 3: the name of the interface, as Go spells the type.
func decideName(_ context.Context, step decide) (g.Decision, error) {
	publisher := vocabulary.Reference.Identity.Publisher
	switch outcome, _, _ := step.Evidence(0, model.EvidenceRecordForm); outcome {
	case g.Proved:
	case g.Refuted:
		return g.Decision{}, g.Refuse(g.RefusalRequirementUnsatisfied, publisher, ReasonNoRecordForm)
	default:
		return g.Decision{}, g.Refuse(g.RefusalRequirementUnresolved, publisher, ReasonNoRecordForm)
	}
	reader, err := step.Reader()
	if err != nil {
		return g.Decision{}, err
	}
	declaration, err := reader.Declaration(step.Subject.Declaration.Declaration)
	if err != nil {
		return g.Decision{}, err
	}
	return g.Decision{Payload: vocabulary.Payload{Name: declaration.Name}}, nil
}
