// Package ast defines language-parameterized syntax carriers.
//
// The language parameter separates ASTs and node kinds in Go's type system.
// A language package supplies its own named marker type, kinds and grammar;
// this package contains no language catalogue or grammar implementation.
//
// These are data types, not checked trees. Structural and grammatical
// validation, allocation, traversal and immutable builders are not implemented.
package ast

// TreeID identifies a particular AST within a set of inputs and outputs.
// It is an opaque provenance/addressing key, not a declaration or content hash.
type TreeID string

// NodeID is a one-based index into an AST's Nodes slice. Zero names no node.
// IDs are local to an AST; use NodeRef when crossing tree boundaries.
type NodeID uint32

// AST holds one syntax tree in a flat node arena.
//
// L is a language-owned named marker type. It is not a runtime registration.
// A populated tree has a nonempty ID and a Root in the range 1..len(Nodes).
// Node i+1 is stored at Nodes[i]. The zero value is an empty carrier.
//
// Children describe syntax containment; References describe resolved links.
// Recursive declarations use references rather than containment cycles.
// Fields remain mutable: callers must treat stage inputs as read-only and
// create separate output storage. No immutability or validity is enforced here.
type AST[L any] struct {
	ID    TreeID
	Root  NodeID
	Nodes []Node[L]
}
