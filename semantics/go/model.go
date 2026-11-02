// Package semantics stores established, provider-owned semantic facts.
// It checks graph integrity, never a language's lookup or type rules.
package semantics

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Bitspark/bittype/trees/ast/go"
)

var ErrInvalid = errors.New("invalid semantic model")
var ErrLimit = errors.New("semantic model resource limit")

// IDs are one-based addresses local to one Model, never persistent identities.
type EntityID uint32
type TermID uint32

// Key names vocabulary owned by a provider. It causes no registration or lookup.
type Key struct{ Owner, Name string }

func (k Key) Valid() bool {
	return strings.TrimSpace(k.Owner) != "" && strings.TrimSpace(k.Name) != "" && utf8.ValidString(k.Owner) && utf8.ValidString(k.Name)
}

type Entity struct {
	Kind         Key
	Name         string
	Owner        EntityID
	Declarations []ast.NodeRef
	Parameters   []Parameter
	Body         TermID
}

type Parameter struct {
	Entity EntityID
	Sort   Key
}

type TermKind uint8

const (
	Reference TermKind = iota + 1
	Application
	ProviderExpression
)

// Term has exactly one variant. Application arguments are effective positional
// arguments, not a mixture of written, inferred, defaulted or packed arguments.
// A provider expression's encoding defines its bytes and operand locators.
type Term struct {
	Kind      TermKind
	Target    EntityID
	Arguments []TermID
	Operator  Key
	Encoding  Key
	Data      []byte
	Operands  []Operand
}

type Operand struct {
	Locator []uint32
	Term    TermID
}

// Occurrence identifies a semantic position within a current syntax node.
// Role and Ordinal are provider-owned; an occurrence is not a candidate set.
type Occurrence struct {
	Site    ast.NodeRef
	Role    Key
	Ordinal uint32
}

type Binding struct {
	At     Occurrence
	Target EntityID
}

type Expression struct {
	At   Occurrence
	Term TermID
}

type Dependency struct {
	From     EntityID
	To       EntityID
	Relation Key
	At       *Occurrence
}

// Data is construction storage. New bounds, validates and copies it before
// publication. The producer must not mutate it concurrently with New.
type Data struct {
	Entities     []Entity
	Terms        []Term
	Bindings     []Binding
	Expressions  []Expression
	Dependencies []Dependency
}

type Limits struct {
	MaxRecords int
	MaxEdges   int
	MaxBytes   int
	MaxDepth   int
}

func DefaultLimits() Limits {
	return Limits{MaxRecords: 100000, MaxEdges: 800000, MaxBytes: 16 << 20, MaxDepth: 256}
}

func (l Limits) Valid() bool {
	return l.MaxRecords > 0 && uint64(l.MaxRecords) <= uint64(^uint32(0)) && l.MaxEdges > 0 && l.MaxBytes > 0 && l.MaxDepth > 0 && l.MaxDepth <= 512
}

// Model is an immutable owned graph. Its zero value is invalid. Accessors return
// independent mutable copies; copying the Model handle is safe.
type Model struct{ state *Data }

func (m Model) Valid() bool { return m.state != nil }

type budget struct{ records, edges, bytes int }

func take(left *int, count int) bool {
	if count < 0 || count > *left {
		return false
	}
	*left -= count
	return true
}

func (b *budget) key(k Key) bool {
	return take(&b.bytes, len(k.Owner)) && take(&b.bytes, len(k.Name))
}

func (b *budget) site(s ast.NodeRef) bool      { return take(&b.bytes, len(s.Tree)) }
func (b *budget) occurrence(o Occurrence) bool { return b.site(o.Site) && b.key(o.Role) }

func bound(d Data, l Limits) error {
	b := budget{records: l.MaxRecords, edges: l.MaxEdges, bytes: l.MaxBytes}
	for _, count := range []int{len(d.Entities), len(d.Terms), len(d.Bindings), len(d.Expressions), len(d.Dependencies)} {
		if !take(&b.records, count) {
			return ErrLimit
		}
	}
	for _, e := range d.Entities {
		if !b.key(e.Kind) || !take(&b.bytes, len(e.Name)) || !take(&b.edges, len(e.Declarations)) || !take(&b.edges, len(e.Parameters)) {
			return ErrLimit
		}
		for _, site := range e.Declarations {
			if !b.site(site) {
				return ErrLimit
			}
		}
		for _, parameter := range e.Parameters {
			if !b.key(parameter.Sort) {
				return ErrLimit
			}
		}
	}
	for _, term := range d.Terms {
		if !b.key(term.Operator) || !b.key(term.Encoding) || !take(&b.bytes, len(term.Data)) || !take(&b.edges, len(term.Arguments)) || !take(&b.edges, len(term.Operands)) {
			return ErrLimit
		}
		for _, operand := range term.Operands {
			if !take(&b.edges, len(operand.Locator)) {
				return ErrLimit
			}
		}
	}
	for _, binding := range d.Bindings {
		if !b.occurrence(binding.At) {
			return ErrLimit
		}
	}
	for _, expression := range d.Expressions {
		if !b.occurrence(expression.At) {
			return ErrLimit
		}
	}
	for _, dependency := range d.Dependencies {
		if !b.key(dependency.Relation) || (dependency.At != nil && !b.occurrence(*dependency.At)) {
			return ErrLimit
		}
	}
	return nil
}

// New creates a complete immutable graph. contains must recognize direct sites
// in the captured input catalogue; historical origins are not that catalogue.
// The callback must not perform IO or mutate input. Nil admits no source sites.
// Bounds are checked before copying. Only term containment must be acyclic;
// entity references and labeled dependency edges may form cycles.
func New(data Data, limits Limits, contains func(ast.NodeRef) bool) (Model, error) {
	return newModel(data, limits, contains, true)
}

// CheckLimits re-admits sealed storage without copying or changing its identity.
func (m Model) CheckLimits(limits Limits) error {
	if !m.Valid() { return ErrInvalid }
	_, err := newModel(*m.state, limits, func(ast.NodeRef) bool { return true }, false)
	return err
}

func newModel(data Data, limits Limits, contains func(ast.NodeRef) bool, copyStorage bool) (Model, error) {
	if !limits.Valid() {
		return Model{}, fmt.Errorf("%w: invalid limits", ErrInvalid)
	}
	if err := bound(data, limits); err != nil {
		return Model{}, err
	}
	validEntity := func(id EntityID) bool { return id > 0 && uint64(id) <= uint64(len(data.Entities)) }
	validTerm := func(id TermID) bool { return id > 0 && uint64(id) <= uint64(len(data.Terms)) }
	validSite := func(site ast.NodeRef) bool {
		return site.Tree != "" && site.Node > 0 && contains != nil && contains(site)
	}
	validOccurrence := func(o Occurrence) bool { return o.Role.Valid() && validSite(o.Site) }
	bad := func(reason string) (Model, error) { return Model{}, fmt.Errorf("%w: %s", ErrInvalid, reason) }
	parameterOwners := make(map[EntityID]EntityID)
	for i, e := range data.Entities {
		id := EntityID(i + 1)
		if !e.Kind.Valid() || !utf8.ValidString(e.Name) || (e.Owner != 0 && (!validEntity(e.Owner) || e.Owner == id)) || (e.Body != 0 && !validTerm(e.Body)) {
			return bad("invalid entity")
		}
		sites := make(map[ast.NodeRef]bool)
		for _, site := range e.Declarations {
			if !validSite(site) || sites[site] {
				return bad("invalid or duplicate declaration site")
			}
			sites[site] = true
		}
		for _, parameter := range e.Parameters {
			if !validEntity(parameter.Entity) || parameter.Entity == id || !parameter.Sort.Valid() || parameterOwners[parameter.Entity] != 0 || data.Entities[parameter.Entity-1].Owner != id {
				return bad("invalid parameter ownership")
			}
			parameterOwners[parameter.Entity] = id
		}
	}
	for _, term := range data.Terms {
		switch term.Kind {
		case Reference, Application:
			if !validEntity(term.Target) || term.Operator != (Key{}) || term.Encoding != (Key{}) || len(term.Data) != 0 || len(term.Operands) != 0 {
				return bad("invalid reference/application fields")
			}
			if term.Kind == Reference && len(term.Arguments) != 0 {
				return bad("reference has arguments")
			}
			if term.Kind == Application && len(term.Arguments) != len(data.Entities[term.Target-1].Parameters) {
				return bad("application arity does not match effective parameters")
			}
			for _, argument := range term.Arguments {
				if !validTerm(argument) {
					return bad("unknown application argument")
				}
			}
		case ProviderExpression:
			if term.Target != 0 || len(term.Arguments) != 0 || !term.Operator.Valid() || !term.Encoding.Valid() {
				return bad("invalid provider expression fields")
			}
			locations := make(map[string]bool)
			for _, operand := range term.Operands {
				if !validTerm(operand.Term) {
					return bad("unknown provider operand")
				}
				var path strings.Builder
				for _, part := range operand.Locator {
					path.WriteString(strconv.FormatUint(uint64(part), 10))
					path.WriteByte('/')
				}
				key := path.String()
				if locations[key] {
					return bad("duplicate provider operand locator")
				}
				locations[key] = true
			}
		default:
			return bad("unknown term variant")
		}
	}
	bindings := make(map[Occurrence]bool)
	for _, binding := range data.Bindings {
		if !validOccurrence(binding.At) || !validEntity(binding.Target) || bindings[binding.At] {
			return bad("invalid or duplicate binding")
		}
		bindings[binding.At] = true
	}
	expressions := make(map[Occurrence]bool)
	for _, expression := range data.Expressions {
		if !validOccurrence(expression.At) || !validTerm(expression.Term) || expressions[expression.At] {
			return bad("invalid or duplicate expression association")
		}
		expressions[expression.At] = true
	}
	for _, dependency := range data.Dependencies {
		if !validEntity(dependency.From) || !validEntity(dependency.To) || !dependency.Relation.Valid() || (dependency.At != nil && !validOccurrence(*dependency.At)) {
			return bad("invalid dependency")
		}
	}
	// Memoized heights enforce the longest containment path, including DAG joins.
	state := make([]uint8, len(data.Terms))
	height := make([]int, len(data.Terms))
	var visit func(TermID, int) (int, error)
	visit = func(id TermID, depth int) (int, error) {
		if depth > limits.MaxDepth {
			return 0, ErrLimit
		}
		index := int(id - 1)
		if state[index] == 1 {
			return 0, fmt.Errorf("%w: term containment cycle", ErrInvalid)
		}
		if state[index] == 2 {
			if height[index]+depth-1 > limits.MaxDepth {
				return 0, ErrLimit
			}
			return height[index], nil
		}
		state[index] = 1
		maximum := 1
		child := func(next TermID) error {
			h, err := visit(next, depth+1)
			if err != nil {
				return err
			}
			if h+1 > maximum {
				maximum = h + 1
			}
			return nil
		}
		for _, argument := range data.Terms[index].Arguments {
			if err := child(argument); err != nil {
				return 0, err
			}
		}
		for _, operand := range data.Terms[index].Operands {
			if err := child(operand.Term); err != nil {
				return 0, err
			}
		}
		state[index], height[index] = 2, maximum
		return maximum, nil
	}
	for i := range data.Terms {
		if _, err := visit(TermID(i+1), 1); err != nil {
			return Model{}, err
		}
	}
	owned := data
	if copyStorage { owned = clone(data) }
	return Model{state: &owned}, nil
}

func copyEntity(e Entity) Entity {
	e.Declarations = slices.Clone(e.Declarations)
	e.Parameters = slices.Clone(e.Parameters)
	return e
}

func copyTerm(t Term) Term {
	t.Arguments = slices.Clone(t.Arguments)
	t.Data = slices.Clone(t.Data)
	t.Operands = slices.Clone(t.Operands)
	for i := range t.Operands {
		t.Operands[i].Locator = slices.Clone(t.Operands[i].Locator)
	}
	return t
}

func copyDependency(d Dependency) Dependency {
	if d.At != nil {
		site := *d.At
		d.At = &site
	}
	return d
}

func clone(data Data) Data {
	result := Data{Entities: slices.Clone(data.Entities), Terms: slices.Clone(data.Terms), Bindings: slices.Clone(data.Bindings), Expressions: slices.Clone(data.Expressions), Dependencies: slices.Clone(data.Dependencies)}
	for i := range result.Entities {
		result.Entities[i] = copyEntity(result.Entities[i])
	}
	for i := range result.Terms {
		result.Terms[i] = copyTerm(result.Terms[i])
	}
	for i := range result.Dependencies {
		result.Dependencies[i] = copyDependency(result.Dependencies[i])
	}
	return result
}

func (m Model) Data() Data {
	if !m.Valid() {
		return Data{}
	}
	return clone(*m.state)
}

func (m Model) Entity(id EntityID) (Entity, bool) {
	if !m.Valid() || id == 0 || uint64(id) > uint64(len(m.state.Entities)) {
		return Entity{}, false
	}
	return copyEntity(m.state.Entities[id-1]), true
}

func (m Model) Term(id TermID) (Term, bool) {
	if !m.Valid() || id == 0 || uint64(id) > uint64(len(m.state.Terms)) {
		return Term{}, false
	}
	return copyTerm(m.state.Terms[id-1]), true
}

func (m Model) Resolve(at Occurrence) (EntityID, bool) {
	if m.Valid() {
		for _, binding := range m.state.Bindings {
			if binding.At == at {
				return binding.Target, true
			}
		}
	}
	return 0, false
}

func (m Model) TermAt(at Occurrence) (TermID, bool) {
	if m.Valid() {
		for _, expression := range m.state.Expressions {
			if expression.At == at {
				return expression.Term, true
			}
		}
	}
	return 0, false
}

func (m Model) DependenciesFrom(from EntityID, relation Key) []Dependency {
	var result []Dependency
	if m.Valid() {
		for _, dependency := range m.state.Dependencies {
			if dependency.From == from && dependency.Relation == relation {
				result = append(result, copyDependency(dependency))
			}
		}
	}
	return result
}

type ArgumentBinding struct {
	Parameter EntityID
	Argument  TermID
}

// ParameterBindings derives the sole positional binding view. A zero-argument
// application returns an empty successful view; a reference is not an application.
func (m Model) ParameterBindings(id TermID) ([]ArgumentBinding, bool) {
	if !m.Valid() || id == 0 || uint64(id) > uint64(len(m.state.Terms)) {
		return nil, false
	}
	term := m.state.Terms[id-1]
	if term.Kind != Application {
		return nil, false
	}
	parameters := m.state.Entities[term.Target-1].Parameters
	result := make([]ArgumentBinding, len(parameters))
	for i, parameter := range parameters {
		result[i] = ArgumentBinding{Parameter: parameter.Entity, Argument: term.Arguments[i]}
	}
	return result, true
}
