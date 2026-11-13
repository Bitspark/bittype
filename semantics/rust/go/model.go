package rust

import (
	"bytes"
	"encoding/json"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type model struct {
	receipt semantic.Result
	graph   semantics.Semantics[s.Language]
}
type Model struct{ state *model }
type Declaration struct {
	owner *model
	id    semantics.EntityID
}
type Type struct {
	owner *model
	id    semantics.TermID
}

// Field is one field of a struct type. Its visibility is the declared fact of
// Declaration, read with Declaration.Visibility like that of an item.
type Field struct {
	Name        string
	Type        Type
	Declaration Declaration
	Sites       []ast.NodeRef
}
type fieldInfo struct {
	Name   string
	Entity semantics.EntityID
}

// Visibility is a declared visibility in the compiler's resolved form: public,
// or restricted to the module Scope. Scope is invalid when Public holds. The
// crate root is a module, and it is the scope of pub(crate).
type Visibility struct {
	Public bool
	Scope  Declaration
}

// Effective says whether Audience can name a declaration through Path. In this
// profile the audience is another crate and the path is the definition's, and
// Nameable holds exactly when the declaration and every module around it are
// declared public. It is not what the compiler calls reachable, which includes
// types that another crate meets and cannot name.
type Effective struct {
	Audience string
	Path     string
	Nameable bool
}

// Open admits the exact analysis receipt of a supported environment; it invokes no compiler.
func Open(result semantic.Result) (Model, error) {
	evidence, err := EvidenceFor(result.Configuration())
	if err != nil {
		return Model{}, err
	}
	if _, err := semantic.AdmitEvidence(result.Snapshot(), []semantic.EvidenceRequirement{evidence}, []semantic.Result{result}); err != nil {
		return Model{}, err
	}
	envelope, _ := result.Semantics()
	graph, err := representation.UnpackSemantics[s.Language](result.Snapshot().Language(), envelope)
	if err != nil {
		return Model{}, err
	}
	return Model{&model{result, graph}}, nil
}
func (m Model) Valid() bool { return m.state != nil }
func (m Model) Same(other Model) bool {
	return m.Valid() && other.Valid() && m.state.receipt == other.state.receipt
}
func (m Model) Receipt() semantic.Result {
	if !m.Valid() {
		return semantic.Result{}
	}
	return m.state.receipt
}
func (m Model) Snapshot() semantic.Snapshot { return m.Receipt().Snapshot() }
func (m Model) Graph() semantics.Semantics[s.Language] {
	if !m.Valid() {
		return semantics.Semantics[s.Language]{}
	}
	return m.state.graph
}
func (m Model) Declaration(id semantics.EntityID) (Declaration, error) {
	if !m.Valid() {
		return Declaration{}, ErrForeign
	}
	if _, ok := m.state.graph.Entity(id); !ok {
		return Declaration{}, ErrForeign
	}
	return Declaration{m.state, id}, nil
}

// Lookup selects an unambiguous compiler path, for example captured::models::Person.
// IDs remain authoritative; paths can coincide across native namespaces.
func (m Model) Lookup(path string) (Declaration, error) {
	if !m.Valid() {
		return Declaration{}, ErrForeign
	}
	var found semantics.EntityID
	for i, e := range m.Graph().Data().Entities {
		if e.Name == path {
			if found != 0 {
				return Declaration{}, ErrForeign
			}
			found = semantics.EntityID(i + 1)
		}
	}
	return m.Declaration(found)
}
func (m Model) Type(id semantics.TermID) (Type, error) {
	if !m.Valid() || len(m.Graph().Facts(semantics.Subject{Term: id}, NativeTypeRelation())) != 1 {
		return Type{}, ErrForeign
	}
	return Type{m.state, id}, nil
}
func (d Declaration) Valid() bool            { return d.owner != nil && d.id != 0 }
func (d Declaration) ID() semantics.EntityID { return d.id }
func (d Declaration) Model() Model           { return Model{d.owner} }
func (d Declaration) Entity() semantics.Entity {
	if !d.Valid() {
		return semantics.Entity{}
	}
	e, _ := d.owner.graph.Entity(d.id)
	return e
}
func (d Declaration) Type() Type {
	if !d.Valid() {
		return Type{}
	}
	terms := d.owner.graph.Facts(semantics.Subject{Entity: d.id}, DeclaredTypeRelation())
	if len(terms) != 1 {
		return Type{}
	}
	return Type{d.owner, terms[0]}
}
func (d Declaration) Parameters() []Declaration {
	var out []Declaration
	for _, p := range d.Entity().Parameters {
		out = append(out, Declaration{d.owner, p.Entity})
	}
	return out
}

// Predicates returns declared compiler predicates, not proofs of arbitrary uses.
func (d Declaration) Predicates() []semantics.TermID {
	if !d.Valid() {
		return nil
	}
	return d.owner.graph.Facts(semantics.Subject{Entity: d.id}, PredicateRelation())
}
func (d Declaration) Default() Type {
	if !d.Valid() {
		return Type{}
	}
	terms := d.owner.graph.Facts(semantics.Subject{Entity: d.id}, DefaultRelation())
	if len(terms) != 1 {
		return Type{}
	}
	return Type{d.owner, terms[0]}
}

// Visibility reads the declared visibility of a captured item or field. It
// fails with ErrMissingFact when the graph has no such fact for the
// declaration, more than one, or one in another form. A declaration without
// the fact is never reported as private.
func (d Declaration) Visibility() (Visibility, error) {
	term, err := d.fact(DeclaredVisibilityRelation())
	if err != nil {
		return Visibility{}, err
	}
	switch {
	case term.Operator == PublicOperator() && len(term.Operands) == 0:
		return Visibility{Public: true}, nil
	case term.Operator == RestrictedOperator() && len(term.Operands) == 1:
		scope, _ := d.owner.graph.Term(term.Operands[0].Term)
		module, found := d.owner.graph.Entity(scope.Target)
		if scope.Kind != semantics.Reference || !found || module.Kind != ModuleKind() {
			return Visibility{}, ErrMissingFact
		}
		return Visibility{Scope: Declaration{d.owner, scope.Target}}, nil
	}
	return Visibility{}, ErrMissingFact
}

// Effective reads the effective fact of a captured item. It fails with
// ErrMissingFact like Visibility, and also for an audience or a path that this
// reader does not know.
func (d Declaration) Effective() (Effective, error) {
	term, err := d.fact(EffectiveVisibilityRelation())
	if err != nil {
		return Effective{}, err
	}
	var data struct {
		Audience *string
		Path     *string
		Nameable *bool
	}
	decoder := json.NewDecoder(bytes.NewReader(term.Data))
	decoder.DisallowUnknownFields()
	if term.Operator != EffectiveVisibilityOperator() || len(term.Operands) != 0 || decoder.Decode(&data) != nil {
		return Effective{}, ErrMissingFact
	}
	if data.Audience == nil || data.Path == nil || data.Nameable == nil || *data.Audience != AudienceAnotherCrate || *data.Path != PathDefinition {
		return Effective{}, ErrMissingFact
	}
	return Effective{Audience: *data.Audience, Path: *data.Path, Nameable: *data.Nameable}, nil
}

// PrivateInterface returns the declarations that the compiler observed in the
// interface of this one although they are more private than it: for a field,
// the type of `pub field: Hidden` beside a private `struct Hidden`. Another
// crate cannot use such an interface. The result is empty when the compiler
// observed none.
func (d Declaration) PrivateInterface() ([]Declaration, error) {
	if !d.Valid() {
		return nil, ErrForeign
	}
	var out []Declaration
	for _, id := range d.owner.graph.Facts(semantics.Subject{Entity: d.id}, PrivateInterfaceRelation()) {
		term, _ := d.owner.graph.Term(id)
		if _, found := d.owner.graph.Entity(term.Target); term.Kind != semantics.Reference || !found {
			return nil, ErrMissingFact
		}
		out = append(out, Declaration{d.owner, term.Target})
	}
	return out, nil
}

// fact returns the term of the one fact that a declaration must have for relation.
func (d Declaration) fact(relation semantics.Key) (semantics.Term, error) {
	if !d.Valid() {
		return semantics.Term{}, ErrForeign
	}
	terms := d.owner.graph.Facts(semantics.Subject{Entity: d.id}, relation)
	if len(terms) != 1 {
		return semantics.Term{}, ErrMissingFact
	}
	term, found := d.owner.graph.Term(terms[0])
	if !found || term.Kind != semantics.ProviderExpression {
		return semantics.Term{}, ErrMissingFact
	}
	return term, nil
}
func (t Type) Valid() bool          { return t.owner != nil && t.id != 0 }
func (t Type) ID() semantics.TermID { return t.id }
func (t Type) Model() Model         { return Model{t.owner} }

// Same is compiler type identity with its generic binder environment in one receipt.
// It makes no claim about coercions, trait satisfaction or representation preservation.
func (t Type) Same(other Type) bool {
	return t.Valid() && other.Valid() && t.Model().Same(other.Model()) && t.id == other.id
}
func (t Type) Term() semantics.Term {
	if !t.Valid() {
		return semantics.Term{}
	}
	term, _ := t.owner.graph.Term(t.id)
	return term
}
func (t Type) fact(relation semantics.Key) semantics.TermID {
	if !t.Valid() {
		return 0
	}
	terms := t.owner.graph.Facts(semantics.Subject{Term: t.id}, relation)
	if len(terms) != 1 {
		return 0
	}
	return terms[0]
}
func (t Type) String() string {
	if !t.Valid() {
		return ""
	}
	term, _ := t.owner.graph.Term(t.fact(DisplayRelation()))
	return string(term.Data)
}
func (t Type) Fields() ([]Field, bool) {
	if !t.Valid() {
		return nil, false
	}
	term, _ := t.owner.graph.Term(t.fact(FieldsRelation()))
	if term.Operator != FieldsOperator() {
		return nil, false
	}
	var metadata []fieldInfo
	if json.Unmarshal(term.Data, &metadata) != nil || len(metadata) != len(term.Operands) {
		return nil, false
	}
	out := make([]Field, len(metadata))
	for i, f := range metadata {
		e, _ := t.owner.graph.Entity(f.Entity)
		out[i] = Field{Name: f.Name, Type: Type{t.owner, term.Operands[i].Term}, Declaration: Declaration{t.owner, f.Entity}, Sites: e.Declarations}
	}
	return out, true
}
