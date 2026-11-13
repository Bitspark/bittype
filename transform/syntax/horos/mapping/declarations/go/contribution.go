package declarations

// PROTOTYPE (bittype#196): the declaration policy under the contribution
// contract. It is the part of the division of generics that a contribution
// owns (row 2, O12): the native parameter order, and alias or defined
// declaration. The horos side owns it and it is the same for every target. A
// target's declarations package supplies its two native forms, its evidence
// kind and its payload type.

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	naming "github.com/Bitspark/bittype/transform/syntax/horos/mapping/naming/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// Decision is what the policy decides for one declaration. Each target's
// package records it in the payload type of its own schema.
type Decision struct {
	Name       string
	Nominal    bool
	Parameters []string
}

// Requests is step 1: the policy needs to know whether the representation
// selected for the body is the target's record form. It asks for that as
// evidence about the body, and reads no payload of another family.
func Requests(step g.NameStep[h.Language, source.Reader], recordForm g.EvidenceKind) []g.Request[h.Language] {
	return []g.Request[h.Language]{g.Evidence(step.Subject.Declaration.Body, g.Demand{Kind: recordForm})}
}

// Decide is step 3. owner qualifies the reason of a refusal.
func Decide(step g.DecideStep[h.Language, source.Reader], recordForm g.EvidenceKind, owner string) (Decision, error) {
	open := step.Subject.Declaration
	decl, err := step.Reader.Declaration(open.Declaration)
	if err != nil {
		return Decision{}, err
	}
	name, err := naming.Exported(decl.Name)
	if err != nil {
		return Decision{}, g.Refuse(g.RefusalUnsupported, owner, err.Error())
	}
	if !decl.Root {
		name = "M" + decl.ModulePin[:16] + "_" + name
	}
	outcome, _, _ := step.Evidence(0, recordForm)
	d := Decision{Name: name, Nominal: outcome == g.Proved}
	for _, f := range open.Formals {
		parameter, err := step.Reader.Parameter(f)
		if err != nil {
			return Decision{}, err
		}
		pname, err := naming.Exported(parameter.Name)
		if err != nil {
			return Decision{}, g.Refuse(g.RefusalUnsupported, owner, err.Error())
		}
		d.Parameters = append(d.Parameters, pname)
	}
	return d, nil
}

// Describe is the projector's part: the native declaration from the recorded
// decision and the body's native type. It reads no source.
func Describe[T any](p g.Projection[T], d Decision, record, alias g.Form[T]) (g.DeclarationDescription[T], error) {
	if len(d.Parameters) != p.Binders() || p.Dependencies() != 1 {
		return g.DeclarationDescription[T]{}, g.ErrScope
	}
	form := alias
	if d.Nominal {
		form = record
	}
	desc := g.DeclarationDescription[T]{Name: d.Name, Form: form, Origin: p.Origin(), Body: p.Dependency(0).Type}
	for i, name := range d.Parameters {
		desc.Generic.Parameters = append(desc.Generic.Parameters, g.ParameterDescription{Name: name, Origin: p.BinderOrigin(i), RequiredUse: true})
		desc.Generic.Bindings = append(desc.Generic.Bindings, g.Parameter[T](i))
		desc.Generic.Arguments = append(desc.Generic.Arguments, g.Argument[T](i))
	}
	return desc, nil
}
