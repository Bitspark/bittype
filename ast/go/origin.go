package ast

// NodeRef addresses a node in an explicitly identified AST.
// Tree must name that AST and Node must name one of its nodes.
// The address carries no proof of the target's language, grammar or identity;
// a resolver/checker must establish those properties when they are required.
type NodeRef struct {
	Tree TreeID
	Node NodeID
}

// SourceID identifies an input text in a caller-supplied source collection.
// It is opaque to ast: it is not opened as a file or interpreted as a URL.
type SourceID string

// Span is a half-open range of UTF-8 byte offsets [Start, End) in Source.
// End must be at least Start, and both must fit within the supplied source.
// A zero-length range is permitted; absent provenance is an empty Spans slice.
type Span struct {
	Source SourceID
	Start  uint64
	End    uint64
}

// Origin records optional source spans and contributing nodes from earlier ASTs.
// Multiple contributors can explain a generated or combined node. Node addresses
// are intentionally language-independent so target syntax can refer back to its
// canonical input. This is provenance, not a semantic identity or trust claim.
type Origin struct {
	Spans []Span
	Nodes []NodeRef
}
