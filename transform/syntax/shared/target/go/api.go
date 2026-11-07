package target

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// Builder is the symbolic construction API, independent of canonical meanings.
// Native libraries extend this interface with their own syntax conveniences.
type Builder[F, T any] interface {
	g.Native[F, T]
	Node(ast.Node[T], ...Child[T]) (g.Fragment[T], error)
	Inspect(g.Fragment[T]) (ast.Node[T], []Child[T], error)
	Parameter(string, ast.NodeRef) (g.Fragment[T], error)
	Reserve(g.Site[F, T], string, bool, ast.NodeRef) (g.Destination[T], bool, error)
	Header(g.Site[F, T], g.Destination[T], []string, []ast.NodeRef) (g.Layout[T], error)
	Formal(g.Layout[T], int) (Symbol[T], error)
	Named(g.Callee[F, T]) (Symbol[T], error)
	Reference(Symbol[T], []g.Fragment[T], ast.NodeRef) (g.Fragment[T], error)
	UsesFormal(g.Fragment[T], g.Layout[T], int) bool
}
