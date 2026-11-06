// Package rust configures the public generation composer for Horos and rust.
package rust

import (
	"github.com/Bitspark/bittype/contracts/transform/go"
	c "github.com/Bitspark/bittype/generation/syntax/go"
	source "github.com/Bitspark/bittype/generation/syntax/source/horos/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// New installs no builtin or library mappings. The consumer supplies the complete factory.
// bittype:scope docs/generation-composition.md#construction-and-selection
func New(from, to representation.Language, factory c.Factory[h.Language, s.Language]) (transform.ASTTransformer, error) {
	if from.ID() != (h.Language{}).ID() || from.Profile() != h.Profile || to.ID() != (s.Language{}).ID() {
		return nil, transform.ErrInvalidTransformer
	}
	return c.NewTransformer(from, to, source.Contract(), factory)
}
