package source

import (
	"fmt"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/semantics/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/bitspark/ontos/codec/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// Descriptor is a source-owned reading of a built-in descriptor occurrence.
// Native rules live in supplied mappings. Literal payloads are never traversed.
type Descriptor struct {
	Constructor   string
	Discriminator string
	Origin        ast.NodeRef
}
type Field struct {
	Name     string
	Optional bool
	Use      g.Use[h.Language]
	Origin   ast.NodeRef
}
type Declaration struct {
	Name      string
	ModulePin string
	Root      bool
	Origin    ast.NodeRef
}
type Parameter struct {
	Name   string
	Origin ast.NodeRef
}

// reader is the only value representation libraries receive. It exposes the
// interpretation's reading methods, not its kernel-facing judgments.
type reader struct{ s *interpretation }

func (r reader) Declaration(d g.Declaration[h.Language]) (Declaration, error) {
	e, ok, err := r.s.definition(d)
	if err != nil {
		return Declaration{}, err
	}
	if !ok {
		return Declaration{}, g.ErrScope
	}
	pin, err := r.s.modulePin(e.Owner)
	if err != nil {
		return Declaration{}, err
	}
	module := r.s.data.Entities[e.Owner-1]
	if len(module.Declarations) != 1 {
		return Declaration{}, semantic.ErrInvalidRequest
	}
	return Declaration{Name: e.Name, ModulePin: pin, Root: module.Declarations[0] == r.s.rootModule, Origin: e.Declarations[0]}, nil
}
func (r reader) Parameter(b g.Binder[h.Language]) (Parameter, error) {
	id, ok := b.In(r.s.graph)
	if !ok {
		return Parameter{}, g.ErrScope
	}
	e := r.s.data.Entities[id-1]
	if len(e.Declarations) != 1 {
		return Parameter{}, semantic.ErrInvalidRequest
	}
	return Parameter{e.Name, e.Declarations[0]}, nil
}
func (r reader) Descriptor(u g.Use[h.Language]) (Descriptor, bool, error) { return r.s.descriptor(u) }

// Fields reads a closed or open record payload (horos/record, kind, field...).
// Each field is (horos/field, name, req|opt, type); field k's type is the
// operand at locator [k+2, 3].
func (r reader) Fields(v g.UseView[h.Language]) ([]Field, error) {
	id, ok := v.Use.In(r.s.graph)
	if !ok || v.Kind != g.ExpressionUse {
		return nil, g.ErrScope
	}
	p, err := r.s.payload(id)
	if err != nil {
		return nil, err
	}
	if p.Len() < 2 || atom(p.At(0)) != "horos/record" {
		return nil, semantic.ErrInvalidRequest
	}
	var fields []Field
	for n := 2; n < p.Len(); n++ {
		f, ok := p.At(n).(core.Tuple)
		if !ok || f.Len() != 4 {
			return nil, semantic.ErrInvalidRequest
		}
		name, err := data.ReadText(f.At(1))
		if err != nil {
			return nil, err
		}
		u, ok := v.Operand(uint32(n), 3)
		if !ok {
			return nil, semantic.ErrInvalidRequest
		}
		term, _ := u.In(r.s.graph)
		fields = append(fields, Field{Name: name, Optional: atom(f.At(2)) == "opt", Use: u, Origin: r.s.sites[term]})
	}
	return fields, nil
}

// Element is the operand at locator [1] of a collection descriptor.
func (r reader) Element(v g.UseView[h.Language]) (g.Use[h.Language], error) {
	if _, ok := v.Use.In(r.s.graph); !ok || v.Kind != g.ExpressionUse {
		return g.Use[h.Language]{}, g.ErrScope
	}
	u, ok := v.Operand(1)
	if !ok {
		return g.Use[h.Language]{}, semantic.ErrInvalidRequest
	}
	return u, nil
}
func (r reader) Unsupported(u g.Use[h.Language], message string) error {
	id, _ := u.In(r.s.graph)
	site := r.s.sites[id]
	return &text.Diagnostic{Code: text.Unsupported, Message: message, Node: &site}
}

func atom(v core.Value) string {
	a, ok := v.(core.Atom)
	if !ok {
		return ""
	}
	return string(a.Bytes())
}
func (s *interpretation) payload(id semantics.TermID) (core.Tuple, error) {
	t := s.data.Terms[id-1]
	if t.Encoding != (semantics.Key{Owner: "ontos", Name: "codec-v1"}) {
		return core.Tuple{}, semantic.ErrInvalidRequest
	}
	v, err := codec.Decode(t.Data)
	if err != nil {
		return core.Tuple{}, err
	}
	p, ok := v.(core.Tuple)
	if !ok {
		return core.Tuple{}, semantic.ErrInvalidRequest
	}
	return p, nil
}
func (s *interpretation) modulePin(id semantics.EntityID) (string, error) {
	term := s.identities[id]
	if term == 0 {
		return "", semantic.ErrInvalidRequest
	}
	t := s.data.Terms[term-1]
	if t.Operator != (semantics.Key{Owner: "horos", Name: "identity/module"}) || t.Encoding != (semantics.Key{Owner: "ontos", Name: "codec-v1"}) {
		return "", semantic.ErrInvalidRequest
	}
	v, err := codec.Decode(t.Data)
	if err != nil {
		return "", err
	}
	a, ok := v.(core.Atom)
	if !ok || len(a.Bytes()) != 32 {
		return "", semantic.ErrInvalidRequest
	}
	return fmt.Sprintf("%x", a.Bytes()), nil
}
