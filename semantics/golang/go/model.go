package golang

import (
	"encoding/json"
	"go/types"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type model struct {
	receipt semantic.Result
	graph   semantics.Semantics[s.Language]
}
type Model struct{ state *model }
type Type struct {
	owner *model
	id    semantics.TermID
}
type Field struct {
	Name     string
	Type     Type
	Site     ast.NodeRef
	Embedded bool
	Tag      string
}

// FieldInfo is the Go JSON field-metadata encoding. Operands hold the field types;
// Entity identifies the field's declaration in this same graph.
type FieldInfo struct {
	Name     string
	Tag      string
	Embedded bool
	Entity   semantics.EntityID
}

// Open requires an admitted analysis receipt and never invokes a compiler.
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
func (m Model) Same(n Model) bool {
	return m.Valid() && n.Valid() && m.state.receipt == n.state.receipt
}
func (m Model) Snapshot() semantic.Snapshot { return m.Receipt().Snapshot() }
func (m Model) Receipt() semantic.Result {
	if !m.Valid() {
		return semantic.Result{}
	}
	return m.state.receipt
}
func (m Model) Graph() semantics.Semantics[s.Language] {
	if !m.Valid() {
		return semantics.Semantics[s.Language]{}
	}
	return m.state.graph
}
func (m Model) Type(id semantics.TermID) (Type, error) {
	if !m.Valid() {
		return Type{}, ErrForeign
	}
	// Only terms explicitly carrying native type identity are type handles.
	if len(m.state.graph.Facts(semantics.Subject{Term: id}, Key("identical"))) != 1 {
		return Type{}, ErrForeign
	}
	return Type{m.state, id}, nil
}

// Lookup finds a package-level type declared by captured code. It fails when
// more than one captured package declares the name; use LookupIn then.
func (m Model) Lookup(name string) (Type, error) {
	return m.lookup("", name)
}

// LookupIn finds a package-level type of the package at importPath. A supplied
// package's declaration is present only where captured code uses it.
func (m Model) LookupIn(importPath, name string) (Type, error) {
	if importPath == "" {
		return Type{}, ErrForeign
	}
	return m.lookup(importPath, name)
}
func (m Model) lookup(importPath, name string) (Type, error) {
	if !m.Valid() {
		return Type{}, ErrForeign
	}
	var found []semantics.TermID
	for _, e := range m.state.graph.Data().Entities {
		if e.Name != name || (e.Kind != Key("defined") && e.Kind != Key("alias")) {
			continue
		}
		owner, _ := m.state.graph.Entity(e.Owner)
		// Captured packages have declaration sites; supplied packages have none.
		if owner.Kind == Key("package") && (importPath == "" && len(owner.Declarations) > 0 || importPath != "" && owner.Name == importPath) {
			found = append(found, e.Body)
		}
	}
	if len(found) != 1 {
		return Type{}, ErrForeign
	}
	return m.Type(found[0])
}
func (m Model) Basic(kind types.BasicKind) (Type, error) {
	if !m.Valid() || (kind != types.String && kind != types.Bool) {
		return Type{}, ErrUnsupported
	}
	for i := range m.state.graph.Data().Terms {
		t := Type{m.state, semantics.TermID(i + 1)}
		if t.IsBasic(kind) {
			return t, nil
		}
	}
	return Type{}, ErrForeign
}
func (t Type) Valid() bool          { return t.owner != nil && t.id != 0 }
func (t Type) Model() Model         { return Model{t.owner} }
func (t Type) ID() semantics.TermID { return t.id }
func (t Type) term() semantics.Term {
	if !t.Valid() {
		return semantics.Term{}
	}
	v, _ := t.owner.graph.Term(t.id)
	return v
}
func (t Type) fact(name string) Type {
	if !t.Valid() {
		return Type{}
	}
	v := t.owner.graph.Facts(semantics.Subject{Term: t.id}, Key(name))
	if len(v) != 1 {
		return Type{}
	}
	return Type{t.owner, v[0]}
}
func (t Type) entity() semantics.Entity {
	v := t.term()
	if v.Kind == semantics.Application {
		return (Type{t.owner, v.Callee}).entity()
	}
	if v.Kind != semantics.Reference {
		return semantics.Entity{}
	}
	e, _ := t.owner.graph.Entity(v.Target)
	return e
}
func (t Type) Same(u Type) bool {
	return t.Valid() && u.Valid() && t.Model().Same(u.Model()) && t.fact("identical").id == u.fact("identical").id
}
func (t Type) IsBasic(kind types.BasicKind) bool {
	v := t.term()
	var k types.BasicKind
	return v.Operator == Key("basic") && json.Unmarshal(v.Data, &k) == nil && k == kind
}
func (t Type) IsNamed() bool { return t.entity().Kind == Key("defined") }
func (t Type) IsAlias() bool { return t.entity().Kind == Key("alias") }
func (t Type) Site() ast.NodeRef {
	e := t.entity()
	if len(e.Declarations) > 0 {
		return e.Declarations[0]
	}
	return ast.NodeRef{}
}
func (t Type) Parameters() []Type {
	if !t.IsNamed() && !t.IsAlias() {
		return nil
	}
	var out []Type
	for _, p := range t.entity().Parameters {
		e, _ := t.owner.graph.Entity(p.Entity)
		out = append(out, Type{t.owner, e.Body})
	}
	return out
}
func (t Type) Arguments() []Type {
	var out []Type
	for _, id := range t.term().Arguments {
		out = append(out, Type{t.owner, id})
	}
	return out
}
func (t Type) Origin() Type {
	if !t.IsNamed() && !t.IsAlias() {
		return Type{}
	}
	if v := t.term(); v.Kind == semantics.Application {
		return Type{t.owner, v.Callee}
	}
	return t
}
func (t Type) Underlying() Type  { return t.fact("underlying") }
func (t Type) Constraint() Type  { return t.fact("constraint") }
func (t Type) AliasTarget() Type { return t.fact("alias-target") }
func (t Type) SliceElement() (Type, bool) {
	v := t.term()
	if v.Operator != Key("slice") || len(v.Operands) != 1 {
		return Type{}, false
	}
	return Type{t.owner, v.Operands[0].Term}, true
}
func (t Type) Fields() ([]Field, bool) {
	v := t.Underlying().term()
	if v.Operator != Key("struct") {
		return nil, false
	}
	var metadata []FieldInfo
	if json.Unmarshal(v.Data, &metadata) != nil || len(metadata) != len(v.Operands) {
		return nil, false
	}
	out := make([]Field, len(metadata))
	for i, f := range metadata {
		e, _ := t.owner.graph.Entity(f.Entity)
		var site ast.NodeRef
		if len(e.Declarations) > 0 {
			site = e.Declarations[0]
		}
		out[i] = Field{f.Name, Type{t.owner, v.Operands[i].Term}, site, f.Embedded, f.Tag}
	}
	return out, true
}

// Unconstrained reports whether the graph establishes that the parameter is
// unconstrained. A supplied parameter's constraint is not projected, so it is
// never reported unconstrained; callers must treat false as "not established".
func (t Type) Unconstrained() bool {
	if t.entity().Kind != Key("parameter") {
		return false
	}
	v := t.Constraint().Underlying().fact("empty").term()
	return string(v.Data) == "true"
}
func (t Type) String() string { return string(t.fact("display").term().Data) }
