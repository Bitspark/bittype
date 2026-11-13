package source

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	n "github.com/Bitspark/bittype/semantics/horos/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
)

// Reader is the horos interpretation API for representation libraries. It
// reads kernel handles against the captured horos graph and supplies no
// navigation: declarations, binders, actuals and children come from the
// kernel's views. Neither the kernel nor native targets interpret these values.
type Reader interface {
	// Declaration reads a horos type definition's name, module pin, root-module
	// membership and declaration site.
	Declaration(g.Declaration[h.Language]) (Declaration, error)
	// Parameter reads a binder's name and declaration site.
	Parameter(g.Binder[h.Language]) (Parameter, error)
	// Descriptor reads a built-in descriptor occurrence. ok is false for a
	// reference or application, which has no descriptor.
	Descriptor(g.Use[h.Language]) (Descriptor, bool, error)
	// Fields reads a record descriptor's fields; each field type is the
	// kernel's operand occurrence at the horos field locator.
	Fields(g.UseView[h.Language]) ([]Field, error)
	// Element is a collection descriptor's element operand.
	Element(g.UseView[h.Language]) (g.Use[h.Language], error)
	// Unsupported is a diagnostic at an occurrence's horos syntax site.
	Unsupported(g.Use[h.Language], string) error
}

type Provider struct{}

func (Provider) Contract() transform.Contract { return Contract() }
func (Provider) Open(ctx context.Context, c g.Captured[h.Language]) (g.Source[h.Language, Reader], error) {
	s, err := open(ctx, c)
	if err != nil {
		return g.Source[h.Language, Reader]{}, err
	}
	return g.Source[h.Language, Reader]{Interpretation: s, Reader: reader{s}, Roles: s.roles, Semantics: Semantics}, nil
}

// Semantics names the source interpretation that this integration reads: the
// source-semantics version axis that a plan records (bittype#196, row 16). It
// is the analysis profile that the horos vocabulary publishes, because that
// profile is the interpretation the checked references use. This package
// spells no identity of its own for it.
const Semantics = n.AnalysisProfile

// OptionalFieldRole is the key of the role selector that this integration
// issues for an optional field of a record (row 12). A contribution whose
// domain is that role names this key.
const OptionalFieldRole = "horos/optional-record-field"

var _ g.SourceProvider[h.Language, Reader] = Provider{}
