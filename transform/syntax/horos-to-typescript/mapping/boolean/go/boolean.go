// Package boolean is the horos-to-typescript mapping library for the horos
// scalar bool. A consumer registers it on its own; it imports no other mapping
// library.
package boolean

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

// primitive is repeated in the text mapping library: a mapping library does
// not import another one.
func primitive(sourceName, targetName string) rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		pattern, err := patterns.Head(from, "scalar", sourceName)
		if err != nil {
			return g.MappingDescription[h.Language, s.Language, source.Reader]{}, err
		}
		return g.MappingDescription[h.Language, s.Language, source.Reader]{Profile: rules.RepresentationProfile, Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{{ID: sourceName + "/typescript/1", Input: pattern, Output: model.Primitive(targetName)}}}, nil
	}
}

// Library binds bool to the native primitive named "bool" under the binding ID
// "bool/typescript/1". Both names are unchanged from the package this library
// was split from.
func Library() rules.Library[s.Language] {
	return primitive("bool", model.PrimitiveBool)
}
