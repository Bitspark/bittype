// Package cst holds language-parameterized concrete syntax and its abstract view.
// It knows no language grammar, token spelling or formatting convention.
package cst

import "github.com/Bitspark/bittype/ast/go"

type NodeID uint32
type Kind[L any] string

// Node is either a leaf (Text != nil, no Children) or a container (Text == nil).
// Kind belongs to the language. Leaves include punctuation, comments and spaces.
// Syntax addresses zero or more nodes of this CST's stored abstract view.
type Node[L any] struct {
	Kind     Kind[L]
	Text     *string
	Children []NodeID
	Syntax   []ast.NodeID
	Origin   ast.Origin
}

// CST owns both arenas. Stage inputs and borrowed views must not be mutated.
type CST[L any] struct {
	Abstract ast.AST[L]
	Root     NodeID
	Nodes    []Node[L]
}

func (tree CST[L]) ASTView() ast.AST[L] { return tree.Abstract }
