// Package cpp configures the public generation composer for Horos and cpp.
package cpp

import (
	"github.com/Bitspark/bittype/contracts/transform/go"
	c "github.com/Bitspark/bittype/generation/syntax/go"
	source "github.com/Bitspark/bittype/generation/syntax/source/horos/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// New installs no builtin or library mappings. The consumer supplies the complete factory.
func New(from, to representation.Language, factory c.Factory[h.Language, s.Language]) (transform.ASTTransformer, error) {
	if from.ID() != (h.Language{}).ID() || from.Profile() != h.Profile || to.ID() != (s.Language{}).ID() {
		return nil, transform.ErrInvalidTransformer
	}
	return c.NewTransformer(from, to, source.Contract(), factory)
}
