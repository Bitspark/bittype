package tree

import "github.com/Bitspark/bittype/trees/ast/go"

func Attrs(pairs ...string) []ast.Attribute {
	var out []ast.Attribute
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, ast.Attribute{Name: pairs[i], Text: pairs[i+1]})
	}
	return out
}
func Child(role string, id ast.NodeID) ast.Child { return ast.Child{Role: role, Node: id} }
func Children(role string, ids []ast.NodeID) []ast.Child {
	out := make([]ast.Child, len(ids))
	for i, id := range ids {
		out[i] = Child(role, id)
	}
	return out
}

type View[L any] struct{ Tree ast.AST[L] }

func (v View[L]) Node(id ast.NodeID) ast.Node[L] { return v.Tree.Nodes[id-1] }
func (v View[L]) Kind(id ast.NodeID) string      { return string(v.Node(id).Kind) }
func (v View[L]) Attr(id ast.NodeID, name string) string {
	for _, a := range v.Node(id).Attributes {
		if a.Name == name {
			return a.Text
		}
	}
	return ""
}
func (v View[L]) Has(id ast.NodeID, name string) bool {
	for _, a := range v.Node(id).Attributes {
		if a.Name == name {
			return true
		}
	}
	return false
}
func (v View[L]) Flag(id ast.NodeID, name string) bool { return v.Attr(id, name) == "true" }
func (v View[L]) Children(id ast.NodeID, role string) []ast.NodeID {
	var out []ast.NodeID
	for _, c := range v.Node(id).Children {
		if c.Role == role {
			out = append(out, c.Node)
		}
	}
	return out
}
func (v View[L]) Child(id ast.NodeID, role string) ast.NodeID {
	for _, c := range v.Node(id).Children {
		if c.Role == role {
			return c.Node
		}
	}
	return 0
}
func Clone[L any](input ast.AST[L]) ast.AST[L] {
	out := input
	out.Nodes = make([]ast.Node[L], len(input.Nodes))
	for i, n := range input.Nodes {
		out.Nodes[i] = n
		out.Nodes[i].Attributes = append([]ast.Attribute(nil), n.Attributes...)
		out.Nodes[i].Children = append([]ast.Child(nil), n.Children...)
		out.Nodes[i].References = append([]ast.Reference(nil), n.References...)
		out.Nodes[i].Origin.Spans = append([]ast.Span(nil), n.Origin.Spans...)
		out.Nodes[i].Origin.Nodes = append([]ast.NodeRef(nil), n.Origin.Nodes...)
	}
	return out
}
