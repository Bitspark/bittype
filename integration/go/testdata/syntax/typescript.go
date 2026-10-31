// This fixture is type-checked separately against the public Go API.
// It is an authored target AST, not generated source or a grammar-check result.
package consumer

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/language/go"
	"github.com/Bitspark/bittype/pipeline/go"
	syntax "github.com/Bitspark/bittype/syntax/typescript/go"
)

var _ language.Language = syntax.Language{}

// Page<T> { items: T[] }, Person, PeoplePage = Page<Person>, Reader<T> { read(): T }.
var example = syntax.AST{
	ID:   "typescript-example",
	Root: 1,
	Nodes: []syntax.Node{
		// 1
		{Kind: syntax.KindSourceFile, Children: []ast.Child{{Role: syntax.RoleDeclaration, Node: 2}, {Role: syntax.RoleDeclaration, Node: 7}, {Role: syntax.RoleDeclaration, Node: 10}, {Role: syntax.RoleDeclaration, Node: 14}}},
		// 2
		{Kind: syntax.KindInterface, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Page"}}, Children: []ast.Child{{Role: syntax.RoleTypeParameter, Node: 3}, {Role: syntax.RoleMember, Node: 4}}},
		// 3
		{Kind: syntax.KindTypeParameter, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}},
		// 4
		{Kind: syntax.KindProperty, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "items"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 5}}},
		// 5
		{Kind: syntax.KindArrayType, Children: []ast.Child{{Role: syntax.RoleElement, Node: 6}}},
		// 6
		{Kind: syntax.KindTypeReference, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "typescript-example", Node: 3}}}},
		// 7
		{Kind: syntax.KindInterface, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Person"}}, Children: []ast.Child{{Role: syntax.RoleMember, Node: 8}}},
		// 8
		{Kind: syntax.KindProperty, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "name"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 9}}},
		// 9
		{Kind: syntax.KindKeywordType, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "string"}}},
		// 10
		{Kind: syntax.KindTypeAlias, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "PeoplePage"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 11}}},
		// 11
		{Kind: syntax.KindTypeApplication, Children: []ast.Child{{Role: syntax.RoleConstructor, Node: 12}, {Role: syntax.RoleArgument, Node: 13}}},
		// 12
		{Kind: syntax.KindTypeReference, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Page"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "typescript-example", Node: 2}}}},
		// 13
		{Kind: syntax.KindTypeReference, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Person"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "typescript-example", Node: 7}}}},
		// 14
		{Kind: syntax.KindInterface, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Reader"}}, Children: []ast.Child{{Role: syntax.RoleTypeParameter, Node: 15}, {Role: syntax.RoleMember, Node: 16}}},
		// 15
		{Kind: syntax.KindTypeParameter, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}},
		// 16
		{Kind: syntax.KindMethod, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "read"}}, Children: []ast.Child{{Role: syntax.RoleSignature, Node: 17}}},
		// 17
		{Kind: syntax.KindFunctionType, Children: []ast.Child{{Role: syntax.RoleResult, Node: 18}}},
		// 18
		{Kind: syntax.KindTypeReference, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "typescript-example", Node: 15}}}},
	},
}

var carrier ast.AST[syntax.Language] = example

// Signature witnesses only; no stage is implemented or invoked.
var operation func(syntax.AST) (syntax.AST, error)
var stage pipeline.Stage[syntax.AST, syntax.AST] = operation
