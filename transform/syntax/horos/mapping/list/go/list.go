// Package list holds the horos side of a list contribution: a horos list
// refers to what is selected for its element. It is the same for every
// target. A target's list package supplies its model's collection form and
// its payload type, and it is the contribution.
package list

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// Requests is step 1: one reference, to the element of the list.
func Requests(step g.NameStep[h.Language, source.Reader]) ([]g.Request[h.Language], error) {
	reader, err := step.Reader()
	if err != nil {
		return nil, err
	}
	element, err := reader.Element(step.Subject.Use)
	if err != nil {
		return nil, err
	}
	return []g.Request[h.Language]{g.Reference(element)}, nil
}

// Element is the projector's part: the target type that was selected for the
// element. It reads no source.
func Element[T any](p g.Projection[T]) (g.Expression[T], error) {
	if p.Dependencies() != 1 {
		return g.Expression[T]{}, g.ErrScope
	}
	return p.Dependency(0).Type, nil
}
