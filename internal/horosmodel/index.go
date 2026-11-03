package horosmodel

import (
	"encoding/hex"

	s "github.com/Bitspark/bittype/internal/declsyntax"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

// Index records syntax addresses only. Callers first establish semantic truth
// through their concrete resolver. No pin, annotation or type fact is inferred.
func Index[L any](tree ast.AST[L]) (*Decoded[L], error) {
	if err := s.Check(tree); err != nil {
		return nil, err
	}
	d := &Decoded[L]{tree: tree, modules: make(map[decl.Digest]ast.NodeID), declarations: make(map[decl.DeclarationID]ast.NodeID), parameters: make(map[decl.ParameterID]ast.NodeID), imports: make(map[importID]ast.NodeID)}
	r := reader[L]{tree: tree}
	for _, child := range tree.Nodes[tree.Root-1].Children {
		bytes, err := hex.DecodeString(r.attr(child.Node, s.AttributePin))
		if err != nil || len(bytes) != 32 {
			return nil, ErrInvalid
		}
		var pin decl.Digest
		copy(pin[:], bytes)
		d.indexModule(child.Node, pin, r)
	}
	return d, nil
}

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

func (d *Decoded[L]) ImportSite(module decl.Digest, alias string) ast.NodeRef {
	return d.site(d.imports[importID{module: module, alias: alias}])
}
