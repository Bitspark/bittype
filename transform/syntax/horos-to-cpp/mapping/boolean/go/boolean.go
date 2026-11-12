// Package boolean is the builtin horos-to-cpp mapping library for the horos
// scalar bool. It describes the C++ primitive bool for it.
package boolean

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/cpp/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// primitive is repeated in the text package. A mapping package does not
// import another mapping package.
func primitive(profile, sourceName, targetName string) rules.Library[s.Language] {
	return rules.Library[s.Language]{Profile: profile, DescribeRules: func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := patterns.Head(from, "scalar", sourceName)
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: sourceName + "/cpp/1", Input: pattern, Output: model.Primitive(targetName)}}}, nil
	}}
}

// Library maps bool to bool under the binding bool/cpp/1.
func Library(profile string) rules.Library[s.Language] {
	return primitive(profile, "bool", model.PrimitiveBool)
}
