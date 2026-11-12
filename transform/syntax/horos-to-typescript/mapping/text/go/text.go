// Package text is the horos-to-typescript mapping library for the horos scalar
// utf8-text. A consumer registers it on its own; it imports no other mapping
// library.
package text

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// primitive is repeated in the boolean mapping library: a mapping library does
// not import another one.
func primitive(profile, sourceName, targetName string) rules.Library[s.Language] {
	return rules.Library[s.Language]{Profile: profile, DescribeRules: func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := patterns.Head(from, "scalar", sourceName)
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: sourceName + "/typescript/1", Input: pattern, Output: model.Primitive(targetName)}}}, nil
	}}
}

// Library binds utf8-text to the TypeScript primitive string under the binding
// ID "utf8-text/typescript/1".
func Library(profile string) rules.Library[s.Language] {
	return primitive(profile, "utf8-text", "string")
}
