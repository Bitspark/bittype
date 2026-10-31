package ast

// Kind is a language-owned syntax discriminator. The same spelling in two
// different named language markers does not give the kinds the same Go type.
// There are no reserved kinds or built-in language/contract names.
type Kind[L any] string

// Attribute is an uninterpreted textual syntax atom.
//
// A grammar defines the name, allowed text and meaning. Text can preserve an
// identifier or an exact literal spelling without rounding or normalization.
// It is not an escape hatch for executable source or a runtime value codec.
type Attribute struct {
	Name string
	Text string
}

// Child is an ordered containment edge to another node in the same AST.
// The grammar owns Role names and cardinality; repeated roles can form a list.
type Child struct {
	Role string
	Node NodeID
}

// Reference is a resolved, non-containment link.
// A grammar owns Role and the allowed target. Unresolved spellings belong in
// syntax attributes, not in a fabricated target with a missing node.
type Reference struct {
	Role   string
	Target NodeRef
}

// Node contains language-owned syntax and uninterpreted metadata.
//
// Attributes and Children preserve their supplied order. A grammar decides
// which order matters and whether duplicate attributes or roles are allowed.
// No grammar, name resolution or normalization is performed by these fields.
type Node[L any] struct {
	Kind       Kind[L]
	Attributes []Attribute
	Children   []Child
	References []Reference
	Origin     Origin
}
