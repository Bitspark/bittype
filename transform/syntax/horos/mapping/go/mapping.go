// Package mapping supplies Horos-owned source patterns and declaration recipes.
package mapping

import (
	"context"
	"fmt"
	"strings"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

const RepresentationProfile = "pagination-representation/1"
const EnforcementProfile = "pagination-enforcement/1"

type Library[T any] struct {
	Profile       string
	DescribeRules func(representation.Language, representation.Language) (g.MappingDescription[h.Language, T, source.Reader], error)
}

func (p Library[T]) Describe(from, to representation.Language) (g.MappingDescription[h.Language, T, source.Reader], error) {
	if p.Profile != RepresentationProfile {
		return g.MappingDescription[h.Language, T, source.Reader]{}, g.ErrUnsupported
	}
	if p.DescribeRules == nil {
		return g.MappingDescription[h.Language, T, source.Reader]{}, g.ErrInvalid
	}
	d, err := p.DescribeRules(from, to)
	d.Profile = p.Profile
	return d, err
}
func Pattern(from representation.Language, head, discriminator string) (g.Pattern[h.Language, source.Reader], error) {
	key, err := source.Builtin(from, head, discriminator)
	return g.Pattern[h.Language, source.Reader]{Selector: key}, err
}
func Element(from representation.Language, slot g.Slot) (g.Pattern[h.Language, source.Reader], error) {
	p, err := Pattern(from, "list-of", "")
	if err != nil {
		return p, err
	}
	p.Slots = []g.Slot{slot}
	p.Capture = func(ctx context.Context, src source.Reader, u g.UseView[h.Language]) ([]g.Capture[h.Language], error) {
		operand, err := src.Element(u)
		if err != nil {
			return nil, err
		}
		return []g.Capture[h.Language]{{Slot: slot, Use: operand}}, nil
	}
	return p, nil
}
func Declarations[T any](profile string, record, alias g.Form[T]) Library[T] {
	return Library[T]{profile, func(from, to representation.Language) (g.MappingDescription[h.Language, T, source.Reader], error) {
		return g.MappingDescription[h.Language, T, source.Reader]{Declarations: []g.DeclarationRecipe[h.Language, T, source.Reader]{{ID: "declarations", Expand: func(ctx context.Context, r g.DeclarationRequest[h.Language, T, source.Reader]) (g.DeclarationDescription[T], error) {
			zero := g.DeclarationDescription[T]{}
			src, err := r.Reader()
			if err != nil {
				return zero, err
			}
			open, err := r.Declaration()
			if err != nil {
				return zero, err
			}
			decl, err := src.Declaration(open.Declaration)
			if err != nil {
				return zero, err
			}
			name, err := Exported(decl.Name)
			if err != nil {
				return zero, src.Unsupported(open.Body, err.Error())
			}
			if !decl.Root {
				name = "M" + decl.ModulePin[:16] + "_" + name
			}
			body, descriptor, err := src.Descriptor(open.Body)
			if err != nil {
				return zero, err
			}
			form := alias
			if descriptor && body.Constructor == "record" {
				form = record
			}
			desc := g.DeclarationDescription[T]{Name: name, Form: form, Origin: decl.Origin}
			for i, f := range open.Formals {
				parameter, err := src.Parameter(f)
				if err != nil {
					return zero, err
				}
				pname, err := Exported(parameter.Name)
				if err != nil {
					return zero, err
				}
				desc.Generic.Parameters = append(desc.Generic.Parameters, g.ParameterDescription{Name: pname, Origin: parameter.Origin, RequiredUse: true})
				desc.Generic.Bindings = append(desc.Generic.Bindings, g.Parameter[T](i))
				desc.Generic.Arguments = append(desc.Generic.Arguments, g.Argument[T](i))
			}
			desc.Body, err = r.Mapped(open.Body)
			return desc, err
		}}}}, nil
	}}
}
func Exported(name string) (string, error) {
	if name == "" || name == "_" || strings.Contains(name, "__") {
		return "", fmt.Errorf("unsupported identifier %q", name)
	}
	for i, b := range []byte(name) {
		if !(b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b == '_' || i > 0 && b >= '0' && b <= '9') {
			return "", fmt.Errorf("unsupported identifier %q", name)
		}
	}
	if name[0] == '_' {
		name = "Type" + name
	} else if name[0] >= 'a' && name[0] <= 'z' {
		name = string(name[0]-32) + name[1:]
	}
	if name == "Self" {
		name += "_"
	}
	return name, nil
}
