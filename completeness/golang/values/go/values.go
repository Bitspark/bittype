// Package values is the Go value vocabulary of the completeness test backend
// (prototype, bittype#196, row 29, O8). It is a test instrument: only test
// compositions, the Go completeness backend and the Go facet packages import
// it. bittype's native vocabularies stay types-only, and the production owner
// of value-level Go is the future Go runtime library.
//
// The product this vocabulary serves is "Go conversion code that compiles and
// runs". Its facet kinds are owned here, because the contract of a kind
// belongs to the product that the facet serves. The kernel knows none of them.
//
// What the vocabulary is, and what it is not. A lowering returns Go text with
// typed holes. Every native type, member, binder, adapter argument, generated
// declaration and owned helper enters the text through a hole that holds a
// reference of the plan, which the backend checks with the plan and resolves.
// The text between the holes is the facet's own Go. So this vocabulary shows
// that every reference category of the contract can be issued, checked and
// resolved, and that the result compiles and runs. It is not a typed
// expression vocabulary: "source-code strings never stand in for typed target
// expressions" is not met for the text between the holes.
package values

import (
	"fmt"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
)

// Target is the target that the facets of this vocabulary are for.
var Target = string(s.Language{}.ID())

// Issue makes the facet of one kind for a contribution of this target. The
// publisher of a facet package is the publisher of its contribution.
func Issue(kind g.FacetKind, contribution g.ContributionReference, lowering Lowering) (g.Facet, error) {
	return kind.Issue(contribution.Identity.Publisher, contribution.Identity, Target, lowering)
}

// Product is the identity of the product that a request names to ask for
// executable Go conversion.
const Product = "bittype.completeness/golang-conversion"

// Hole is one typed hole of lowered Go text. Exactly one of its parts is set.
type Hole struct {
	// Reference is a reference of the plan, and Want the category that the
	// place in the text needs: a type, a member, an owned helper.
	Reference g.PlanReference
	Want      g.ReferenceCategory
	// Dependency is 1 + the index of a link of the node: the expression that
	// the backend lowered for that dependency. Declared says that the
	// dependency is the body of the nominal declaration being lowered.
	Dependency int
	Declared   bool
	// Helper is the local name of a helper that this lowering introduces.
	Helper string
	// SelfType is the Go type of the nominal declaration whose body the node is.
	SelfType bool
}

// Code is Go text with holes. The marker of hole i is "«i»".
type Code struct {
	Text  string
	Holes []Hole
}

// Go builds code: each verb %s of format stands for the next hole.
func Go(format string, holes ...Hole) Code {
	markers := make([]any, len(holes))
	for i := range holes {
		markers[i] = fmt.Sprintf("«%d»", i)
	}
	return Code{fmt.Sprintf(format, markers...), holes}
}

func Dependency(link int) Hole { return Hole{Dependency: link + 1} }
func DeclaredBody() Hole       { return Hole{Dependency: 1, Declared: true} }
func Helper(name string) Hole  { return Hole{Helper: name} }
func SelfType() Hole           { return Hole{SelfType: true} }
func Reference(r g.PlanReference, want g.ReferenceCategory) Hole {
	return Hole{Reference: r, Want: want}
}

// HelperCode is a top-level Go declaration that a lowering introduces (row
// 29: "A helper that the lowering itself introduces"). Its key is the
// contribution reference and Name; the backend emits it once and issues its
// spelling.
type HelperCode struct {
	Name string
	// Code declares the helper. The hole Helper(Name) in it is its own name.
	Code Code
}

// Context is what the backend gives a lowering for one node. It has the typed
// view of the node's frozen payload, decoded by the codec that the backend's
// composition registered, the references the plan issues for the node, and
// what the source semantics says about the node's children.
type Context interface {
	// Payload is the typed view of the node's frozen decision.
	Payload() any
	// Links is the number of resolved requests of the node, and Stated what
	// the rule stated in request i.
	Links() int
	Stated(link int) string
	// Fragment reports whether link i is a field position.
	Fragment(link int) bool
	// Member refers to the member of the node's native type that link i
	// supplies; a field position supplies the members of its fragment.
	Member(link int) (g.PlanReference, error)
	// Owned refers to a declaration that the node's contribution owns.
	Owned(name string) (g.PlanReference, error)
	// Canonical is the name that the source gives the child of link i: the
	// field name of a record. It is source meaning, read from the snapshot.
	Canonical(link int) (string, error)
}

// Lowering is the value of the lowering facet kind: how one contribution's
// recorded decisions become Go conversion code. For a use the code is an
// expression of the runtime type Adapter[N]. For a field position it is an
// expression of the runtime type FieldState[M] (the field-state protocol).
// For the declaration policy it is the adapter of the declared type.
type Lowering struct {
	Lower   func(Context) (Code, error)
	Helpers func(Context) ([]HelperCode, error)
}

func kind(identity string) g.FacetKind {
	k, err := g.NewFacetKind(g.FacetKindDescription{Identity: identity, Version: "1", Validate: func(value any) error {
		if l, ok := value.(Lowering); !ok || l.Lower == nil {
			return fmt.Errorf("the value is no lowering")
		}
		return nil
	}})
	if err != nil {
		panic(err)
	}
	return k
}

// ValueLowering is the facet kind of a contribution that converts a value on
// its own: a use or a declaration. FieldStateLowering is the facet kind of a
// contribution at a field position, which has a field-state protocol and no
// standalone value adapter (row 14).
var (
	ValueLowering      = kind("bittype.completeness/golang/value-lowering")
	FieldStateLowering = kind("bittype.completeness/golang/field-state-lowering")
)

// Conversion is the product: what it requires of every selected contribution,
// by the kind of its subject.
func Conversion() g.Product {
	return g.Product{Identity: Product, Requires: map[g.SubjectKind][]g.FacetKind{
		g.UseSubject:         {ValueLowering},
		g.DeclarationSubject: {ValueLowering},
		g.FieldSubject:       {FieldStateLowering},
	}}
}
