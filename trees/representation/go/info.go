package representation

import "github.com/Bitspark/bittype/trees/ast/go"

// Info is bounded structural metadata from an already admitted envelope. It
// exposes no syntax storage, semantic identity, content hash or unchecked view.
// Entries counts attributes, containment/reference edges and origin entries.
// TextBytes counts their textual storage, including names, roles and tree IDs.
type Info struct {
	ID        ast.TreeID
	Root      ast.NodeID
	Nodes     int
	Entries   int
	TextBytes int
}

// Inspect reads private owned storage without cloning or invoking a checker.
// Admission already bounded each envelope. Batch consumers can use these counts
// to enforce aggregate budgets before allocating additional copies.
func (a AST) Inspect() (Info, error) {
	if a.payload == nil || a.id == "" || a.profile == "" {
		return Info{}, ErrInvalidAST
	}
	return a.payload.inspect(), nil
}

func (p typed[L]) inspect() Info {
	input := p.value
	result := Info{ID: input.ID, Root: input.Root, Nodes: len(input.Nodes), TextBytes: len(input.ID)}
	for _, node := range input.Nodes {
		result.TextBytes += len(node.Kind)
		result.Entries += len(node.Attributes) + len(node.Children) + len(node.References) + len(node.Origin.Spans) + len(node.Origin.Nodes)
		for _, attribute := range node.Attributes {
			result.TextBytes += len(attribute.Name) + len(attribute.Text)
		}
		for _, child := range node.Children {
			result.TextBytes += len(child.Role)
		}
		for _, reference := range node.References {
			result.TextBytes += len(reference.Role) + len(reference.Target.Tree)
		}
		for _, span := range node.Origin.Spans {
			result.TextBytes += len(span.Source)
		}
		for _, origin := range node.Origin.Nodes {
			result.TextBytes += len(origin.Tree)
		}
	}
	return result
}
