// Package declarations is the typescript-to-horos declaration contribution:
// the default declaration policy for every TypeScript interface and type
// alias without an exact mapping (docs/contribution-contract.md, row 2, O12).
// It belongs to the two other pairs of bittype#196 (decision O43). A
// consumer's composition registers it on its own, with the codec of its
// payload schema.
//
// It covers an interface and a type alias of an object type, and it makes a
// horos type definition of it, named as TypeScript spells the type. It
// refuses, as unsupported, a type alias of a type that is no object type. It
// states no generic convention: the TypeScript source reader reads no generic
// declaration.
//
// It asks no evidence about its body. A declaration contribution of a pair
// horos-to-<language> asks whether its body is selected as the target's
// record form, to choose between a defined declaration and an alias. horos
// has one declaration form, a type definition, which names any type
// expression, so there is nothing to choose and the horos model declares no
// evidence kind. What this contribution needs to know about its body is a
// fact of the SOURCE (is it an object type), and a rule reads source facts
// from the source view.
//
// Not covered and not refused here: a TypeScript name that is no horos
// identifier (for example one with a dollar sign). The horos target refuses
// it as an invalid description when the planner checks the declaration: the
// spelling rule of horos is in a private helper that a mapping may not
// import.
package declarations

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	model "github.com/Bitspark/bittype/transform/syntax/horos/model/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/vocabulary/declarations/go"
	rules "github.com/Bitspark/bittype/transform/syntax/typescript/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/typescript/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[ts.Language, s.Language, source.Reader]
type decide = g.DecideStep[ts.Language, source.Reader]
type name = g.NameStep[ts.Language, source.Reader]

// target is the target that the contribution states (O4).
var target = string(s.Language{}.ID())

// ReasonNoObjectType is the reason of this contribution's refusal of a type
// alias of a type that is no object type.
const ReasonNoObjectType = "a type alias of a type that is no object type is not covered"

// Contribution returns the declaration contribution. It refers to the
// declaration's body and decides the name of the type definition.
//
// It states that the declaration is transparent. The contract's word for
// that is "an alias and no defined declaration in the target" (O35, O41): a
// horos type definition is interchangeable with the expression it names, so
// what is selected for the body is what a use of the definition is. This
// rests on reading horos as structural. No rule of the pair asks evidence
// about a declaration, so the statement changes no output here.
func Contribution(representation.Language) (g.Contribution[ts.Language, s.Language, source.Reader], error) {
	return contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.CanonicalizationProfile,
		Domain: g.ForDeclarationPolicy(), Schema: vocabulary.Schema,
		Rule: g.Rule[ts.Language, source.Reader]{
			Name: func(_ context.Context, step name) ([]g.Request[ts.Language], error) {
				return []g.Request[ts.Language]{g.Reference(step.Subject.Declaration.Body)}, nil
			},
			Decide: decideName,
		},
		Declare: func(_ context.Context, p g.Projection[s.Language]) (g.DeclarationDescription[s.Language], error) {
			payload, ok := p.Payload().(vocabulary.Payload)
			if !ok || p.Dependencies() != 1 || p.Binders() != 0 {
				return g.DeclarationDescription[s.Language]{}, g.ErrInvalid
			}
			return g.DeclarationDescription[s.Language]{Name: payload.Name, Form: model.Form(model.FormDefinition), Origin: p.Origin(), Body: p.Dependency(0).Type}, nil
		},
	}, nil
}

// decideName is step 3: the name of the type definition, as TypeScript
// spells the type.
func decideName(_ context.Context, step decide) (g.Decision, error) {
	reader, err := step.Reader()
	if err != nil {
		return g.Decision{}, err
	}
	object, err := reader.Object(step.Subject.Declaration.Body)
	if err != nil {
		return g.Decision{}, err
	}
	if !object {
		return g.Decision{}, g.Refuse(g.RefusalUnsupported, vocabulary.Reference.Identity.Publisher, ReasonNoObjectType)
	}
	declaration, err := reader.Declaration(step.Subject.Declaration.Declaration)
	if err != nil {
		return g.Decision{}, err
	}
	return g.Decision{Payload: vocabulary.Payload{Name: declaration.Name}, Transparent: true}, nil
}
