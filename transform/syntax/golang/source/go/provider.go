package source

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// Reader is the Go interpretation API for the mapping contributions of a pair
// whose source is Go. It reads kernel handles against the captured Go graph
// and supplies no navigation: declarations and the uses of field types come
// from the kernel's views.
type Reader interface {
	// Declaration reads a defined type's name, the import path of its package
	// and its declaration site.
	Declaration(g.Declaration[s.Language]) (Declaration, error)
	// Fields reads the fields of a struct type use. Each field's type is the
	// kernel's operand occurrence of that field.
	Fields(g.UseView[s.Language]) ([]Field, error)
}

// Declaration is the reading of a defined type. It repeats what the analysed
// vocabulary publishes of the declaration: the entity's name, the package of
// semantics/golang/go.DeclaringPackage and the entity's declaration site.
type Declaration struct {
	Name    string
	Package string
	Origin  ast.NodeRef
}

// Field is the reading of one field of a struct type. It is the Field of
// semantics/golang/go with the field's type given as the kernel's occurrence
// in place of the vocabulary's type handle.
type Field struct {
	Name     string
	Embedded bool
	Tag      string
	Use      g.Use[s.Language]
	Origin   ast.NodeRef
}

// reader is the only value mapping contributions receive. It exposes the
// interpretation's reading methods, not its kernel-facing judgments.
type reader struct{ i *interpretation }

func (r reader) Declaration(d g.Declaration[s.Language]) (Declaration, error) {
	e, owner, named, err := r.i.entity(d)
	if err != nil {
		return Declaration{}, err
	}
	if !named || len(e.Declarations) != 1 {
		return Declaration{}, g.ErrScope
	}
	return Declaration{Name: e.Name, Package: owner.Name, Origin: e.Declarations[0]}, nil
}

// Fields reads the field metadata through the analysed vocabulary. Field k's
// type is the operand at locator [k]: the vocabulary pairs its metadata with
// the struct term's operands by position, and this reader checks that the
// operand the kernel minted there is the term of the field's type.
func (r reader) Fields(v g.UseView[s.Language]) ([]Field, error) {
	id, ok := v.Use.In(r.i.captured)
	if !ok || v.Kind != g.ExpressionUse {
		return nil, g.ErrScope
	}
	structure, err := r.i.model.Type(id)
	if err != nil {
		return nil, semantic.ErrInvalidRequest
	}
	fields, ok := structure.Fields()
	if !ok {
		return nil, semantic.ErrInvalidRequest
	}
	out := make([]Field, len(fields))
	for k, f := range fields {
		u, ok := v.Operand(uint32(k))
		if !ok {
			return nil, semantic.ErrInvalidRequest
		}
		if term, ok := u.In(r.i.captured); !ok || term != f.Type.ID() {
			return nil, semantic.ErrInvalidRequest
		}
		out[k] = Field{Name: f.Name, Embedded: f.Embedded, Tag: f.Tag, Use: u, Origin: f.Site}
	}
	return out, nil
}

// Provider is the Go source of a generation: a consumer's composition
// registers it for the Go language.
type Provider struct{}

func (Provider) Contract() transform.Contract { return Contract() }
func (Provider) Open(_ context.Context, c g.Captured[s.Language]) (g.Source[s.Language, Reader], error) {
	i, err := open(c)
	if err != nil {
		return g.Source[s.Language, Reader]{}, err
	}
	return g.Source[s.Language, Reader]{Interpretation: i, Reader: reader{i}, Semantics: Semantics}, nil
}

// Semantics names the source interpretation that this integration reads: the
// source-semantics version axis that a plan records (row 16). It is the
// analysis profile that the Go vocabulary publishes. This package spells no
// identity of its own for it. The reader issues no field role.
const Semantics = n.AnalysisProfile

var _ g.SourceProvider[s.Language, Reader] = Provider{}
