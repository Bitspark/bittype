// Package optional is the facet package of the optional contribution for Go
// (prototype, bittype#196, O8, O16): the field-state protocol of the optional field on the Go side.
// It is test code of bittype's own bundle. It imports the contribution's
// vocabulary and the Go value vocabulary, and neither the mapping package nor
// the codec. A consumer's composition registers it apart from the mapping.
package optional

import (
	"github.com/Bitspark/bittype/completeness/golang/values/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/optional/go"
)

// Facets returns the facets of the contribution for the Go conversion product.
// The contribution has a field-state protocol and no standalone value adapter
// (row 14). The record's adapter invokes the protocol for the field and never
// reads this contribution's payload.
//
// What the protocol needs on the Go side: the adapter selected for the
// operand, and the declaration this contribution owns. It does not need the
// record, the member's name or the field's canonical name: the native member
// is a value of its own type, BittypeOptional[T], so the record hands the
// member's value over and takes one back.
func Facets() ([]g.Facet, error) {
	f, err := values.Issue(values.FieldStateLowering, vocabulary.Reference, values.Lowering{
		Helpers: func(ctx values.Context) ([]values.HelperCode, error) {
			owned, err := ctx.Owned("optional")
			if err != nil {
				return nil, err
			}
			helper := values.Reference(owned, g.OwnedReference)
			// Present and Value are members of the declaration that this
			// contribution owns: its own spelling, fixed by its revision.
			return []values.HelperCode{{Name: "state", Code: values.Go(`func %s[T any](operand Adapter[T]) FieldState[%s[T]] {
	return FieldState[%s[T]]{
		Construct: func(present bool, value Value) (%s[T], error) {
			if !present {
				return %s[T]{}, nil
			}
			n, err := operand.Decode(value)
			return %s[T]{Present: true, Value: n}, err
		},
		State: func(member %s[T]) (bool, Value, error) {
			if !member.Present {
				return false, nil, nil
			}
			v, err := operand.Encode(member.Value)
			return true, v, err
		},
	}
}`, values.Helper("state"), helper, helper, helper, helper, helper, helper)}}, nil
		},
		Lower: func(values.Context) (values.Code, error) {
			return values.Go("%s(%s)", values.Helper("state"), values.Dependency(0)), nil
		},
	})
	return []g.Facet{f}, err
}
