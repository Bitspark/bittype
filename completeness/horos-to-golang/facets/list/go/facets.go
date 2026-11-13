// Package list is the facet package of the list contribution for Go
// (prototype, bittype#196, O8, O16): its lowering to Go conversion code.
// It is test code of bittype's own bundle. It imports the contribution's
// vocabulary and the Go value vocabulary, and neither the mapping package nor
// the codec. A consumer's composition registers it apart from the mapping.
package list

import (
	"github.com/Bitspark/bittype/completeness/golang/values/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/list/go"
)

// Facets returns the facets of the contribution for the Go conversion product.
func Facets() ([]g.Facet, error) {
	f, err := values.Issue(values.ValueLowering, vocabulary.Reference, values.Lowering{
		Helpers: func(values.Context) ([]values.HelperCode, error) {
			return []values.HelperCode{{Name: "adapter", Code: values.Go(`func %s[E any](element Adapter[E]) Adapter[[]E] {
	return Adapter[[]E]{
		Decode: func(v Value) ([]E, error) {
			items, ok := v.([]any)
			if !ok {
				return nil, fmt.Errorf("structural mismatch: a list is expected, got %%T", v)
			}
			out := make([]E, len(items))
			for i, item := range items {
				n, err := element.Decode(item)
				if err != nil {
					return nil, err
				}
				out[i] = n
			}
			return out, nil
		},
		Encode: func(n []E) (Value, error) {
			out := make([]any, len(n))
			for i, item := range n {
				v, err := element.Encode(item)
				if err != nil {
					return nil, err
				}
				out[i] = v
			}
			return out, nil
		},
	}
}`, values.Helper("adapter"))}}, nil
		},
		// The element's adapter is the adapter of the dependency: a generic
		// adapter argument where the element is bound to a binder.
		Lower: func(values.Context) (values.Code, error) {
			return values.Go("%s(%s)", values.Helper("adapter"), values.Dependency(0)), nil
		},
	})
	return []g.Facet{f}, err
}
