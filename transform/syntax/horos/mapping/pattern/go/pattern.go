// Package pattern holds the source patterns of horos mapping bindings: which
// uses of a horos type a binding selects and which operands it captures. The
// horos source owns the selectors and the reading of an operand; this package
// states them as generation patterns once, so that every target's mapping
// selects the same uses. It decides nothing about a native type, a name or a
// profile.
package pattern

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Head selects the uses of the builtin horos constructor with the given head
// and discriminator. It captures no operand.
func Head(from representation.Language, head, discriminator string) (g.Pattern[h.Language, source.Reader], error) {
	key, err := source.Builtin(from, head, discriminator)
	return g.Pattern[h.Language, source.Reader]{Selector: key}, err
}

// Element selects the uses of a horos list and captures the element type in
// slot.
func Element(from representation.Language, slot g.Slot) (g.Pattern[h.Language, source.Reader], error) {
	p, err := Head(from, "list-of", "")
	if err != nil {
		return p, err
	}
	p.Slots = []g.Slot{slot}
	p.Capture = func(ctx context.Context, src source.Reader, u g.UseView[h.Language]) ([]g.Capture[h.Language], error) {
		operand, err := src.Element(u)
		if err != nil {
			return nil, err
		}
		return []g.Capture[h.Language]{{Slot: slot, Use: operand}}, nil
	}
	return p, nil
}
