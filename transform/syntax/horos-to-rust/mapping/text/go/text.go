// Package text is the builtin horos-to-rust mapping library for the horos
// scalar utf8-text. A consumer registers it on its own.
package text

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/rust/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func primitive(profile, sourceName, targetName string) rules.Library[s.Language] {
	return rules.Library[s.Language]{Profile: profile, DescribeRules: func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := patterns.Head(from, "scalar", sourceName)
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: sourceName + "/rust/1", Input: pattern, Output: model.Primitive(targetName)}}}, nil
	}}
}

// Library maps the horos scalar utf8-text to the Rust model primitive string.
func Library(profile string) rules.Library[s.Language] {
	return primitive(profile, "utf8-text", "string")
}
