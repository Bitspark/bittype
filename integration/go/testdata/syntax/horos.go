package consumer

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

var _ language.Provider = h.Language{}
var _ ast.View[h.Language] = h.AST{}
var slot = h.Node{
	Kind:       h.KindParameter,
	Attributes: []ast.Attribute{{Name: h.AttributeIndex, Text: "0"}},
	References: []ast.Reference{{Role: h.RoleTarget, Target: ast.NodeRef{Tree: "pagination", Node: 3}}},
}
