package ast

// View exposes an already constructed abstract syntax view. It performs no
// parsing, resolution or conversion. The returned storage is borrowed read-only.
//
// bittype:dimension Syntax
type View[L any] interface {
	ASTView() AST[L]
}

func (tree AST[L]) ASTView() AST[L] { return tree }
