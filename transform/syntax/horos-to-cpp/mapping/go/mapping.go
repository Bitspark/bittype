package mapping

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/cpp/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func primitive(profile, sourceName, targetName string) rules.Library[s.Language] {
	return rules.Library[s.Language]{Profile: profile, DescribeRules: func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := rules.Pattern(from, "scalar", sourceName)
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: sourceName + "/cpp/1", Input: pattern, Output: model.Primitive(targetName)}}}, nil
	}}
}
func Text(profile string) rules.Library[s.Language] { return primitive(profile, "utf8-text", "string") }
func Bool(profile string) rules.Library[s.Language] { return primitive(profile, "bool", "bool") }
func List(profile string) rules.Library[s.Language] {
	slot := g.TypeSlot("E")
	return rules.Library[s.Language]{Profile: profile, DescribeRules: func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := rules.Element(from, slot)
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: "list/cpp/1", Input: pattern, Output: model.Collection(g.Mapped[s.Language](slot))}}}, nil
	}}
}
func Record(profile string) rules.Library[s.Language] {
	return rules.Library[s.Language]{Profile: profile, DescribeRules: func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := rules.Pattern(from, "record", "closed")
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: "record/cpp/1", Input: pattern, Expand: func(ctx context.Context, r g.Expansion[h.Language, s.Language, source.Reader]) (g.Expression[s.Language], error) {
			src, err := r.Source()
			if err != nil {
				return g.Expression[s.Language]{}, err
			}
			use, err := r.Use()
			if err != nil {
				return g.Expression[s.Language]{}, err
			}
			if !src.IsBody(use.Environment, use.Use) {
				return g.Expression[s.Language]{}, src.Unsupported(use.Use, "nested records require a named declaration in this profile")
			}
			view, err := src.Read(use.Use)
			if err != nil {
				return g.Expression[s.Language]{}, err
			}
			fields, err := src.Fields(use.Use)
			if err != nil {
				return g.Expression[s.Language]{}, err
			}
			var members []g.Expression[s.Language]
			for _, f := range fields {

				ty, err := r.Mapped(f.Use)
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				if f.Optional {
					ty = model.External("std", "optional", ty)
				}
				members = append(members, g.At(model.Field(f.Name, f.Optional, ty), f.Origin))
			}
			return g.At(model.Record(members...), view.Origin), nil
		}}}}, nil
	}}
}
func Declarations(profile string) rules.Library[s.Language] {
	return rules.Declarations(profile, model.Form("struct"), model.Form("alias"))
}
func Libraries(profile string) []g.MappingLibrary[h.Language, s.Language, source.Reader] {
	return []g.MappingLibrary[h.Language, s.Language, source.Reader]{Text(profile), Bool(profile), List(profile), Record(profile), Declarations(profile)}
}
