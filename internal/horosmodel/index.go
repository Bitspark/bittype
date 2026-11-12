package horosmodel

import (
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

func (d *Decoded[L]) indexModule(node ast.NodeID, pin decl.Digest, r reader[L]) {
	d.modules[pin] = node
	for _, member := range r.node(node).Children {
		if member.Role == s.RoleImport {
			d.imports[importID{module: pin, alias: r.attr(member.Node, s.AttributeAlias)}] = member.Node
			continue
		}
		id := decl.DeclarationID{Module: pin, Name: r.attr(member.Node, s.AttributeName)}
		d.declarations[id] = member.Node
		index := 0
		for _, parameter := range r.node(member.Node).Children {
			if parameter.Role == s.RoleParameter {
				d.parameters[decl.ParameterID{Declaration: id, Index: index}] = parameter.Node
				index++
			}
		}
	}
}
