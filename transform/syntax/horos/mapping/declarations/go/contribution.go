// Package declarations holds the policy by which a named horos declaration
// becomes a native declaration: its name, the module prefix of a declaration
// outside the root module, the choice between the record form and the alias
// form, and its generic parameters (docs/contribution-contract.md, row 2,
// O12). It is the part of the division of generics that a contribution owns:
// the native parameter order, and alias or defined declaration. The horos
// side of a mapping owns the policy, and it is the same for every target. A
// target's declarations package supplies its two native forms, its evidence
// kind and its payload type, and it is the contribution.
//
// The policy states no profile. A contribution states the profile. The policy
// gives a native declaration bare generic parameters: it carries no kind
// constraint into the native type.
package declarations

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	naming "github.com/Bitspark/bittype/transform/syntax/horos/mapping/naming/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// Decision is what the policy decides for one declaration. Each target's
// package records it in the payload type of its own schema.
type Decision struct {
	// Name is the native name: the exported identifier of the horos name, with
	// the module prefix of a declaration outside the root module.
	Name string
	// Nominal: the declaration becomes a defined native declaration. Otherwise
	// it becomes a native alias.
	Nominal bool
	// Parameters are the names of the native parameters, in the order of the
	// source parameters.
	Parameters []string
}

// Requests is step 1: the policy needs to know whether the representation
// selected for the body is the target's record form. It asks for that as
// evidence about the body and reads no payload of another contribution.
func Requests(step g.NameStep[h.Language, source.Reader], recordForm g.EvidenceKind) []g.Request[h.Language] {
	return []g.Request[h.Language]{g.Evidence(step.Subject.Declaration.Body, g.Demand{Kind: recordForm})}
}

// Decide is step 3. owner qualifies the reason of a refusal. A name that the
// naming rule does not cover is unsupported: a later slice of that rule can
// admit the same name.
func Decide(step g.DecideStep[h.Language, source.Reader], recordForm g.EvidenceKind, owner string) (Decision, error) {
	open := step.Subject.Declaration
	reader, err := step.Reader()
	if err != nil {
		return Decision{}, err
	}
	decl, err := reader.Declaration(open.Declaration)
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
		parameter, err := reader.Parameter(f)
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
// decision and the native type of the body. It reads no source. The native
// declaration has bare generic parameters: it carries no kind constraint.
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
