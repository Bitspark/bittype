// Package text is the facet package of the text contribution for Go
// (prototype, bittype#196, O8, O16): its lowering to Go conversion code.
// It is test code of bittype's own bundle. It imports the contribution's
// vocabulary and the Go value vocabulary, and neither the mapping package nor
// the codec. A consumer's composition registers it apart from the mapping.
package text

import (
	"github.com/Bitspark/bittype/completeness/golang/values/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/text/go"
)

// Facets returns the facets of the contribution for the Go conversion product.
func Facets() ([]g.Facet, error) {
	f, err := values.Issue(values.ValueLowering, vocabulary.Reference, values.Lowering{
		// The lowering introduces its adapter as a helper of its own.
		Helpers: func(values.Context) ([]values.HelperCode, error) {
			return []values.HelperCode{{Name: "adapter", Code: values.Go(`func %s() Adapter[string] {
	return Adapter[string]{
		Decode: func(v Value) (string, error) {
			n, ok := v.(string)
			if !ok {
				return n, fmt.Errorf("structural mismatch: text is expected, got %%T", v)
			}
			return n, nil
		},
		Encode: func(n string) (Value, error) { return n, nil },
	}
}`, values.Helper("adapter"))}}, nil
		},
		Lower: func(values.Context) (values.Code, error) { return values.Go("%s()", values.Helper("adapter")), nil },
	})
	return []g.Facet{f}, err
}
