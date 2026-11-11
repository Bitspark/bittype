// Package boolean is the horos-to-golang mapping library for booleans: the
// horos scalar bool becomes the Go type bool.
package boolean

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
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
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: sourceName + "/golang/1", Input: pattern, Output: model.Primitive(targetName)}}}, nil
	}}
}

// Library returns the boolean mapping under the given profile.
func Library(profile string) rules.Library[s.Language] {
	return primitive(profile, "bool", "bool")
}
