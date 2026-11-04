package ast

// Set is construction storage for one language's complete multi-unit output.
// Entry is explicit, even for a one-tree set. It carries no file organization.
type Set[L any] struct {
	Entry NodeRef
	Trees []AST[L]
}
