// This fixture is type-checked separately against the public Go API.
// It is an authored target AST, not generated source or a grammar-check result.
package consumer

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/pipeline/go"
	syntax "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

var _ language.LanguageProvider = syntax.Language{}

// Page[T any] struct { Items []T }, Person, PeoplePage = Page[Person],
// and Reader[T any] interface { Read() T }. The two T binders stay distinct.
var example = syntax.AST{
	ID:   "go-example",
	Root: 1,
	Nodes: []syntax.Node{
		// 1
		{Kind: syntax.KindFile, Attributes: []ast.Attribute{{Name: syntax.AttributePackage, Text: "example"}}, Children: []ast.Child{{Role: syntax.RoleDeclaration, Node: 2}, {Role: syntax.RoleDeclaration, Node: 9}, {Role: syntax.RoleDeclaration, Node: 13}, {Role: syntax.RoleDeclaration, Node: 17}}},
		// 2
		{Kind: syntax.KindTypeDeclaration, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Page"}}, Children: []ast.Child{{Role: syntax.RoleTypeParameter, Node: 3}, {Role: syntax.RoleType, Node: 5}}},
		// 3
		{Kind: syntax.KindTypeParameter, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}, Children: []ast.Child{{Role: syntax.RoleConstraint, Node: 4}}},
		// 4
		{Kind: syntax.KindNamedType, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "any"}}},
		// 5
		{Kind: syntax.KindStructType, Children: []ast.Child{{Role: syntax.RoleField, Node: 6}}},
		// 6
		{Kind: syntax.KindField, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Items"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 7}}},
		// 7
		{Kind: syntax.KindSliceType, Children: []ast.Child{{Role: syntax.RoleElement, Node: 8}}},
		// 8
		{Kind: syntax.KindNamedType, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "go-example", Node: 3}}}},
		// 9
		{Kind: syntax.KindTypeDeclaration, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Person"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 10}}},
		// 10
		{Kind: syntax.KindStructType, Children: []ast.Child{{Role: syntax.RoleField, Node: 11}}},
		// 11
		{Kind: syntax.KindField, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Name"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 12}}},
		// 12
		{Kind: syntax.KindNamedType, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "string"}}},
		// 13
		{Kind: syntax.KindTypeAlias, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "PeoplePage"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 14}}},
		// 14
		{Kind: syntax.KindTypeApplication, Children: []ast.Child{{Role: syntax.RoleConstructor, Node: 15}, {Role: syntax.RoleArgument, Node: 16}}},
		// 15
		{Kind: syntax.KindNamedType, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Page"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "go-example", Node: 2}}}},
		// 16
		{Kind: syntax.KindNamedType, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Person"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "go-example", Node: 9}}}},
		// 17
		{Kind: syntax.KindTypeDeclaration, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Reader"}}, Children: []ast.Child{{Role: syntax.RoleTypeParameter, Node: 18}, {Role: syntax.RoleType, Node: 20}}},
		// 18
		{Kind: syntax.KindTypeParameter, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}, Children: []ast.Child{{Role: syntax.RoleConstraint, Node: 19}}},
		// 19
		{Kind: syntax.KindNamedType, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "any"}}},
		// 20
		{Kind: syntax.KindInterfaceType, Children: []ast.Child{{Role: syntax.RoleMember, Node: 21}}},
		// 21
		{Kind: syntax.KindMethod, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Read"}}, Children: []ast.Child{{Role: syntax.RoleSignature, Node: 22}}},
		// 22
		{Kind: syntax.KindFunctionType, Children: []ast.Child{{Role: syntax.RoleResult, Node: 23}}},
		// 23
		{Kind: syntax.KindParameter, Children: []ast.Child{{Role: syntax.RoleType, Node: 24}}},
		// 24
		{Kind: syntax.KindNamedType, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "go-example", Node: 18}}}},
	},
}

var carrier ast.AST[syntax.Language] = example

// Signature witnesses only; no stage is implemented or invoked.
var operation func(syntax.AST) (syntax.AST, error)
var stage pipeline.Stage[syntax.AST, syntax.AST] = operation
