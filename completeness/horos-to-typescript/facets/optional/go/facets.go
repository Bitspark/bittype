// Package optional is the facet package of the optional contribution for
// TypeScript (prototype, bittype#196, O8, O16): the reference to its run-time
// factory in the stand-in. It is test code of bittype's own bundle. It imports
// the contribution's vocabulary and the stand-in's vocabulary, and neither
// the mapping package nor the codec.
package optional

import (
	"github.com/Bitspark/bittype/completeness/typescript/values/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/optional/go"
)

// Facets returns the facets of the contribution for the stand-in product.
func Facets() ([]g.Facet, error) {
	f, err := values.Issue(values.FieldStateFactory, vocabulary.Reference, "optionalField")
	return []g.Facet{f}, err
}
