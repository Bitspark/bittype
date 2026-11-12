// Package optional is the horos-to-golang mapping library for the optional
// helper: the generic Go declaration BittypeOptional, with the fields Present
// and Value, supplied under the helper key that the consumer names.
package optional

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// library is this package's one mapping library.
type library struct{ key g.HelperKey }

// Library returns the optional helper. It states the structural
// representation profile itself. key is the
// helper key that the consumer created; a mapping that wraps a type in this
// helper receives the same key. The generation contract rejects an invalid key
// with generation.ErrInvalid when the library is registered.
func Library(key g.HelperKey) g.MappingLibrary[h.Language, s.Language, source.Reader] {
	return library{key: key}
}

// Describe states the structural representation profile and the one helper.
func (p library) Describe(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
	key := p.key
	return g.MappingDescription[h.Language, s.Language, source.Reader]{Profile: rules.RepresentationProfile, Helpers: []g.HelperRecipe[s.Language]{{Key: key, Declaration: g.NativeDeclaration[s.Language]{
		Name: "BittypeOptional", Form: model.Form(model.FormDefined), Parameters: []g.NativeParameter{{Name: "T"}},
		Body: model.Record(model.Field("Present", false, model.Primitive(model.PrimitiveBool)), model.Field("Value", false, g.Parameter[s.Language](0))),
	}}}}, nil
}
