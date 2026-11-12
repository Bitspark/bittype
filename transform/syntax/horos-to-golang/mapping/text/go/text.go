// Package text is the horos-to-golang mapping library for text: the horos
// scalar utf8-text becomes the Go type string.
package text

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func primitive(sourceName, targetName string) rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := patterns.Head(from, "scalar", sourceName)
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Profile: rules.RepresentationProfile, Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: sourceName + "/golang/1", Input: pattern, Output: model.Primitive(targetName)}}}, nil
	}
}

// Library returns the text mapping. It states the structural representation
// profile itself.
func Library() rules.Library[s.Language] {
	return primitive("utf8-text", model.PrimitiveString)
}
