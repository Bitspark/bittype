package rust

import (
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
type Field struct {
	Name   string
	Type   Type
	Public bool
	Sites  []ast.NodeRef
}
type fieldInfo struct {
	Name   string
	Public bool
	Entity semantics.EntityID
}

// Open admits the exact analysis receipt; it invokes no compiler.
func Open(result semantic.Result) (Model, error) {
	if _, err := semantic.AdmitEvidence(result.Snapshot(), []semantic.EvidenceRequirement{Evidence()}, []semantic.Result{result}); err != nil {
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
	if !m.Valid() || len(m.Graph().Facts(semantics.Subject{Term: id}, Key("native-type"))) != 1 {
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
	terms := d.owner.graph.Facts(semantics.Subject{Entity: d.id}, Key("declared-type"))
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
	return d.owner.graph.Facts(semantics.Subject{Entity: d.id}, Key("predicate"))
}
func (d Declaration) Default() Type {
	if !d.Valid() {
		return Type{}
	}
	terms := d.owner.graph.Facts(semantics.Subject{Entity: d.id}, Key("default"))
	if len(terms) != 1 {
		return Type{}
	}
	return Type{d.owner, terms[0]}
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
func (t Type) fact(name string) semantics.TermID {
	if !t.Valid() {
		return 0
	}
	terms := t.owner.graph.Facts(semantics.Subject{Term: t.id}, Key(name))
	if len(terms) != 1 {
		return 0
	}
	return terms[0]
}
func (t Type) String() string {
	if !t.Valid() {
		return ""
	}
	term, _ := t.owner.graph.Term(t.fact("display"))
	return string(term.Data)
}
func (t Type) Fields() ([]Field, bool) {
	if !t.Valid() {
		return nil, false
	}
	term, _ := t.owner.graph.Term(t.fact("fields"))
	if term.Operator != Key("fields") {
		return nil, false
	}
	var metadata []fieldInfo
	if json.Unmarshal(term.Data, &metadata) != nil || len(metadata) != len(term.Operands) {
		return nil, false
	}
	out := make([]Field, len(metadata))
	for i, f := range metadata {
		e, _ := t.owner.graph.Entity(f.Entity)
		out[i] = Field{Name: f.Name, Type: Type{t.owner, term.Operands[i].Term}, Public: f.Public, Sites: e.Declarations}
	}
	return out, true
}
