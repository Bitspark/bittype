// Package generation composes supplied mappings without interpreting source or target types.
package generation

import (
	"fmt"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// bindingTable contains the effective selection after provider-qualified replacement.
type bindingTable[F, T any] struct {
	language   representation.Language
	entries    []g.RepresentationBinding[F, T]
	structural *g.Mapping[F, T]
}

func (b bindingTable[F, T]) selectMapping(key g.Selector, structural bool) (g.Mapping[F, T], error) {
	if !key.Valid() || !b.language.Compatible(key.Language()) {
		return g.Mapping[F, T]{}, g.ErrInvalid
	}
	for _, entry := range b.entries {
		if entry.Selector.Key() == key.Key() {
			return entry.Mapping, nil
		}
	}
	if structural && b.structural != nil {
		return *b.structural, nil
	}
	return g.Mapping[F, T]{}, fmt.Errorf("%w: no representation for %s", g.ErrUnsupported, key.Key())
}
