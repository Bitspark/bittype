// This fixture is type-checked separately against the public Go API.
// It is an authored target AST, not generated source or a grammar-check result.
package consumer

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/language/go"
	"github.com/Bitspark/bittype/pipeline/go"
	syntax "github.com/Bitspark/bittype/syntax/rust/go"
)

var _ language.Language = syntax.Language{}

// Page<T> { items: Vec<T> }, Person, PeoplePage = Page<Person>,
// and trait Reader<T> { fn read(&self) -> T; }. Vec and String use the prelude.
var example = syntax.AST{
	ID: "rust-example",
	Root: 1,
	Nodes: []syntax.Node{
		// 1
		{Kind: syntax.KindFile, Children: []ast.Child{{Role: syntax.RoleItem, Node: 2}, {Role: syntax.RoleItem, Node: 9}, {Role: syntax.RoleItem, Node: 13}, {Role: syntax.RoleItem, Node: 18}}},
		// 2
		{Kind: syntax.KindStruct, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Page"}}, Children: []ast.Child{{Role: syntax.RoleGenericParameter, Node: 3}, {Role: syntax.RoleField, Node: 4}}},
		// 3
		{Kind: syntax.KindTypeParameter, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}},
		// 4
		{Kind: syntax.KindField, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "items"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 5}}},
		// 5
		{Kind: syntax.KindPathType, Children: []ast.Child{{Role: syntax.RoleSegment, Node: 6}}},
		// 6
		{Kind: syntax.KindPathSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Vec"}}, Children: []ast.Child{{Role: syntax.RoleArgument, Node: 7}}},
		// 7
		{Kind: syntax.KindPathType, Children: []ast.Child{{Role: syntax.RoleSegment, Node: 8}}},
		// 8
		{Kind: syntax.KindPathSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "rust-example", Node: 3}}}},
		// 9
		{Kind: syntax.KindStruct, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Person"}}, Children: []ast.Child{{Role: syntax.RoleField, Node: 10}}},
		// 10
		{Kind: syntax.KindField, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "name"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 11}}},
		// 11
		{Kind: syntax.KindPathType, Children: []ast.Child{{Role: syntax.RoleSegment, Node: 12}}},
		// 12
		{Kind: syntax.KindPathSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "String"}}},
		// 13
		{Kind: syntax.KindTypeAlias, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "PeoplePage"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 14}}},
		// 14
		{Kind: syntax.KindPathType, Children: []ast.Child{{Role: syntax.RoleSegment, Node: 15}}},
		// 15
		{Kind: syntax.KindPathSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Page"}}, Children: []ast.Child{{Role: syntax.RoleArgument, Node: 16}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "rust-example", Node: 2}}}},
		// 16
		{Kind: syntax.KindPathType, Children: []ast.Child{{Role: syntax.RoleSegment, Node: 17}}},
		// 17
		{Kind: syntax.KindPathSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Person"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "rust-example", Node: 9}}}},
		// 18
		{Kind: syntax.KindTrait, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Reader"}}, Children: []ast.Child{{Role: syntax.RoleGenericParameter, Node: 19}, {Role: syntax.RoleMember, Node: 20}}},
		// 19
		{Kind: syntax.KindTypeParameter, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}},
		// 20
		{Kind: syntax.KindMethodSignature, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "read"}}, Children: []ast.Child{{Role: syntax.RoleReceiver, Node: 21}, {Role: syntax.RoleResult, Node: 22}}},
		// 21
		{Kind: syntax.KindReceiver, Attributes: []ast.Attribute{{Name: syntax.AttributeReceiverForm, Text: "shared"}}},
		// 22
		{Kind: syntax.KindPathType, Children: []ast.Child{{Role: syntax.RoleSegment, Node: 23}}},
		// 23
		{Kind: syntax.KindPathSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "rust-example", Node: 19}}}},
	},
}

var carrier ast.AST[syntax.Language] = example

// Signature witnesses only; no stage is implemented or invoked.
var operation func(syntax.AST) (syntax.AST, error)
var stage pipeline.Stage[syntax.AST, syntax.AST] = operation
