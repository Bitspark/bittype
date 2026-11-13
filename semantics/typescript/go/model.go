package typescript

import (
	"encoding/json"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
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
type Property struct {
	Name     string
	Type     Type
	Optional bool
	Readonly bool
	Sites    []ast.NodeRef
}

// PropertyInfo is TypeScript's json-1 object metadata; operands contain declared
// non-missing value types. Optional presence is separate, not an added undefined.
type PropertyInfo struct {
	Name     string
	Optional bool
	Readonly bool
	Entity   semantics.EntityID
}

// ExtensionInfo is the json-1 data of an extension term: the number of base
// type terms. They are the term's first operands, in written order. Exactly
// one more operand follows them: the object term of the members that the
// interface declares itself.
type ExtensionInfo struct {
	Bases int
}

// Extension is the written definition of an interface with base types.
type Extension struct {
	// Bases are the base type terms as written, in written order: a Reference
	// to the base's entity, or an Application for a base with type arguments.
	Bases []semantics.TermID
	// Members are the members that the interface declares itself. Inherited
	// members are reached through the bases, or all together through
	// Type.Properties of the declared type.
	Members []Property
}

// Open admits the exact analysis receipt without running a compiler.
func Open(result semantic.Result) (Model, error) {
	requirement, err := EvidenceFor(result.Configuration())
	if err != nil {
		return Model{}, err
	}
	if _, err := semantic.AdmitEvidence(result.Snapshot(), []semantic.EvidenceRequirement{requirement}, []semantic.Result{result}); err != nil {
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
	e, ok := m.state.graph.Entity(id)
	if !ok || (e.Kind != AliasKind() && e.Kind != InterfaceKind() && e.Kind != ParameterKind()) {
		return Declaration{}, ErrForeign
	}
	return Declaration{m.state, id}, nil
}

// Lookup uses a captured tree address because module-local names may repeat.
func (m Model) Lookup(tree ast.TreeID, name string) (Declaration, error) {
	if !m.Valid() {
		return Declaration{}, ErrForeign
	}
	for i, e := range m.state.graph.Data().Entities {
		if e.Name != name || (e.Kind != AliasKind() && e.Kind != InterfaceKind()) {
			continue
		}
		for _, site := range e.Declarations {
			if site.Tree == tree {
				return m.Declaration(semantics.EntityID(i + 1))
			}
		}
	}
	return Declaration{}, ErrForeign
}
func (m Model) Type(id semantics.TermID) (Type, error) {
	if !m.Valid() {
		return Type{}, ErrForeign
	}
	if _, ok := m.state.graph.OnlyFact(semantics.Subject{Term: id}, NativeTypeRelation()); !ok {
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
func (d Declaration) fact(relation semantics.Key) Type {
	if !d.Valid() {
		return Type{}
	}
	id, ok := d.owner.graph.OnlyFact(semantics.Subject{Entity: d.id}, relation)
	if !ok {
		return Type{}
	}
	return Type{d.owner, id}
}
func (d Declaration) Type() Type       { return d.fact(DeclaredTypeRelation()) }
func (d Declaration) Constraint() Type { return d.fact(ConstraintRelation()) }
func (d Declaration) Parameters() []Declaration {
	var out []Declaration
	for _, p := range d.Entity().Parameters {
		out = append(out, Declaration{d.owner, p.Entity})
	}
	return out
}
func (t Type) Valid() bool          { return t.owner != nil && t.id != 0 }
func (t Type) ID() semantics.TermID { return t.id }
func (t Type) Model() Model         { return Model{t.owner} }

// Same means the same compiler type in this receipt, not mutual assignability.
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
	id, _ := t.owner.graph.OnlyFact(semantics.Subject{Term: t.id}, relation)
	return id
}
func (t Type) String() string {
	if !t.Valid() {
		return ""
	}
	term, _ := t.owner.graph.Term(t.fact(DisplayRelation()))
	return string(term.Data)
}
func (t Type) Properties() ([]Property, bool) {
	if !t.Valid() {
		return nil, false
	}
	return t.owner.object(t.fact(PropertiesRelation()))
}

// object reads an object term: its members with their value types.
func (m *model) object(id semantics.TermID) ([]Property, bool) {
	term, _ := m.graph.Term(id)
	if term.Operator != ObjectOperator() {
		return nil, false
	}
	var fields []PropertyInfo
	if json.Unmarshal(term.Data, &fields) != nil || len(fields) != len(term.Operands) {
		return nil, false
	}
	out := make([]Property, len(fields))
	for i, f := range fields {
		e, _ := m.graph.Entity(f.Entity)
		out[i] = Property{Name: f.Name, Type: Type{m, term.Operands[i].Term}, Optional: f.Optional, Readonly: f.Readonly, Sites: e.Declarations}
	}
	return out, true
}

// Extension reads an extension term, the body of the entity of an interface
// that has an extends clause. It has no answer for a term in another form:
// another operator, data that is no ExtensionInfo, fewer than one base, a
// number of operands other than the bases and one more, or a last operand
// that is no object term. The body of an interface without an extends clause
// is an object term and no extension.
func (m Model) Extension(id semantics.TermID) (Extension, bool) {
	if !m.Valid() {
		return Extension{}, false
	}
	term, _ := m.state.graph.Term(id)
	if term.Operator != ExtensionOperator() {
		return Extension{}, false
	}
	var info ExtensionInfo
	if json.Unmarshal(term.Data, &info) != nil || info.Bases < 1 || len(term.Operands) != info.Bases+1 {
		return Extension{}, false
	}
	members, ok := m.state.object(term.Operands[info.Bases].Term)
	if !ok {
		return Extension{}, false
	}
	out := Extension{Members: members}
	for _, operand := range term.Operands[:info.Bases] {
		out.Bases = append(out.Bases, operand.Term)
	}
	return out, true
}
