// Package syntax is the agnostic syntax model: the types in which a language
// declares its grammar table, and the names of the predicates that every
// language can use. It holds plain data and names. It evaluates nothing: the
// operation that reads a table decides what a tree must be to satisfy it.
package syntax

import "github.com/Bitspark/bittype/trees/ast/go"

// Predicate names a constraint on the text of an attribute by its owner, as a
// semantic key names a fact. Naming causes no registration or lookup. The zero
// value names no predicate. What a named predicate admits is stated by its
// evaluator, which an operation binds to the name.
type Predicate struct{ Owner, Name string }

// Attribute declares one attribute of a node kind.
type Attribute struct {
	// Required says that a node of the kind has the attribute.
	Required bool
	// Predicate, if not zero, names the constraint on the attribute's text.
	Predicate Predicate
	// OneOf, if present, lists the texts the attribute may have. With a
	// predicate, both hold. An absent set, which in Go is a nil slice,
	// restricts nothing. A set that is present and empty, which in Go is a
	// slice that is not nil and has no element, is a defect of the table: it
	// admits no text, and an operation refuses the table that holds it.
	OneOf []string
}

// Role declares one child role of a node kind: its name, how many children
// may have it and of which kinds they are. Max -1 means no upper bound.
type Role[L any] struct {
	Name     string
	Min, Max int
	Kinds    []ast.Kind[L]
}

// Rule declares one node kind: its attributes by name, its child roles in
// their order, and whether a node of the kind may hold a reference of the role
// "target".
type Rule[L any] struct {
	Attributes map[string]Attribute
	Roles      []Role[L]
	Reference  bool
}

// Rules is one language's grammar table. Its keys and the kinds of its roles
// are that language's node kinds, so a table of one language cannot be applied
// to a tree of another. Role and attribute names are strings.
type Rules[L any] map[ast.Kind[L]]Rule[L]
