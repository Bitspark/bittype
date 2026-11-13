package source

import (
	"context"
	"encoding/json"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// Reader is the TypeScript interpretation API for the mapping contributions
// of a pair whose source is TypeScript. It reads kernel handles against the
// captured TypeScript graph and supplies no navigation: declarations and the
// uses of property types come from the kernel's views.
type Reader interface {
	// Declaration reads the name of an interface or a type alias, the path of
	// its module, which of the two it is, and its declaration site.
	Declaration(g.Declaration[s.Language]) (Declaration, error)
	// Object reports whether a use is an object type.
	Object(g.Use[s.Language]) (bool, error)
	// Properties reads the properties of an object type use. Each property's
	// type is the kernel's operand occurrence of that property.
	Properties(g.UseView[s.Language]) ([]Property, error)
}

// Declaration is the reading of an interface or a type alias. It repeats what
// the analysed vocabulary publishes of the declaration: the entity's name and
// kind, its owner's name and its declaration site.
type Declaration struct {
	Name      string
	Module    string
	Interface bool
	Origin    ast.NodeRef
}

// Property is the reading of one property of an object type. It is the
// Property of semantics/typescript/go with the property's type given as the
// kernel's occurrence in place of the vocabulary's type handle, and with one
// site.
type Property struct {
	Name     string
	Optional bool
	Readonly bool
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
	return Declaration{Name: e.Name, Module: owner.Name, Interface: e.Kind == n.InterfaceKind(), Origin: e.Declarations[0]}, nil
}

func (r reader) Object(u g.Use[s.Language]) (bool, error) {
	id, ok := u.In(r.i.captured)
	if !ok {
		return false, g.ErrScope
	}
	term, _ := r.i.graph.Term(id)
	return term.Operator == n.ObjectOperator(), nil
}

// Properties decodes the metadata of an object term with the metadata type
// that the analysed vocabulary publishes. Property k's type is the operand at
// locator [k]: the vocabulary pairs its metadata with the term's operands by
// position, and this reader checks that the kernel minted an operand there.
// The vocabulary's own Properties reads the same term, but through the type
// that a fact relates to it, and the kernel's occurrence is the term itself.
func (r reader) Properties(v g.UseView[s.Language]) ([]Property, error) {
	id, ok := v.Use.In(r.i.captured)
	if !ok || v.Kind != g.ExpressionUse {
		return nil, g.ErrScope
	}
	term, _ := r.i.graph.Term(id)
	var metadata []n.PropertyInfo
	if term.Operator != n.ObjectOperator() || json.Unmarshal(term.Data, &metadata) != nil || len(metadata) != len(term.Operands) {
		return nil, semantic.ErrInvalidRequest
	}
	out := make([]Property, len(metadata))
	for k, p := range metadata {
		u, ok := v.Operand(uint32(k))
		if !ok {
			return nil, semantic.ErrInvalidRequest
		}
		property := Property{Name: p.Name, Optional: p.Optional, Readonly: p.Readonly, Use: u}
		if e, ok := r.i.graph.Entity(p.Entity); ok && len(e.Declarations) > 0 {
			property.Origin = e.Declarations[0]
		}
		out[k] = property
	}
	return out, nil
}

// Provider is the TypeScript source of a generation: a consumer's composition
// registers it for the TypeScript language.
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
// analysis profile that the TypeScript vocabulary publishes. This package
// spells no identity of its own for it. The reader issues no field role.
const Semantics = n.AnalysisProfile

var _ g.SourceProvider[s.Language, Reader] = Provider{}
