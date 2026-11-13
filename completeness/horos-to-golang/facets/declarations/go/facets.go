// Package declarations is the facet package of the declarations contribution for Go
// (prototype, bittype#196, O8, O16): its lowering to Go conversion code.
// It is test code of bittype's own bundle. It imports the contribution's
// vocabulary and the Go value vocabulary, and neither the mapping package nor
// the codec. A consumer's composition registers it apart from the mapping.
package declarations

import (
	"github.com/Bitspark/bittype/completeness/golang/values/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/declarations/go"
)

// Facets returns the facets of the contribution for the Go conversion product.
// The adapter of a declaration is the adapter of its body. Where the recorded
// decision is the nominal form, the body is lowered as the declared type.
func Facets() ([]g.Facet, error) {
	f, err := values.Issue(values.ValueLowering, vocabulary.Reference, values.Lowering{
		Lower: func(ctx values.Context) (values.Code, error) {
			payload, ok := ctx.Payload().(vocabulary.Payload)
			if !ok {
				return values.Code{}, g.ErrInvalid
			}
			if payload.Nominal {
				return values.Go("%s", values.DeclaredBody()), nil
			}
			return values.Go("%s", values.Dependency(0)), nil
		},
	})
	return []g.Facet{f}, err
}
