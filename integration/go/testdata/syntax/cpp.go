// This fixture is type-checked separately against the public Go API.
// It is an authored target AST, not generated source or a grammar-check result.
package consumer

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/language/go"
	"github.com/Bitspark/bittype/pipeline/go"
	syntax "github.com/Bitspark/bittype/syntax/cpp/go"
)

var _ language.LanguageProvider = syntax.Language{}

// namespace example: template<class T> Page with std::vector<T> items,
// Person, PeoplePage = Page<Person>, and an abstract Reader<T> declaration.
var example = syntax.AST{
	ID:   "cpp-example",
	Root: 1,
	Nodes: []syntax.Node{
		// 1
		{Kind: syntax.KindTranslationUnit, Children: []ast.Child{{Role: syntax.RoleInclude, Node: 2}, {Role: syntax.RoleInclude, Node: 3}, {Role: syntax.RoleDeclaration, Node: 4}}},
		// 2
		{Kind: syntax.KindInclude, Attributes: []ast.Attribute{{Name: syntax.AttributePath, Text: "vector"}, {Name: syntax.AttributeSystem, Text: "true"}}},
		// 3
		{Kind: syntax.KindInclude, Attributes: []ast.Attribute{{Name: syntax.AttributePath, Text: "string"}, {Name: syntax.AttributeSystem, Text: "true"}}},
		// 4
		{Kind: syntax.KindNamespace, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "example"}}, Children: []ast.Child{{Role: syntax.RoleDeclaration, Node: 5}, {Role: syntax.RoleDeclaration, Node: 13}, {Role: syntax.RoleDeclaration, Node: 18}, {Role: syntax.RoleDeclaration, Node: 23}}},
		// 5
		{Kind: syntax.KindRecord, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Page"}, {Name: syntax.AttributeRecordKey, Text: "struct"}}, Children: []ast.Child{{Role: syntax.RoleTemplateParameter, Node: 6}, {Role: syntax.RoleMember, Node: 7}}},
		// 6
		{Kind: syntax.KindTypeParameter, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}},
		// 7
		{Kind: syntax.KindField, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "items"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 8}}},
		// 8
		{Kind: syntax.KindNamedType, Children: []ast.Child{{Role: syntax.RoleSegment, Node: 9}, {Role: syntax.RoleSegment, Node: 10}}},
		// 9
		{Kind: syntax.KindNameSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "std"}}},
		// 10
		{Kind: syntax.KindNameSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "vector"}}, Children: []ast.Child{{Role: syntax.RoleArgument, Node: 11}}},
		// 11
		{Kind: syntax.KindNamedType, Children: []ast.Child{{Role: syntax.RoleSegment, Node: 12}}},
		// 12
		{Kind: syntax.KindNameSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "cpp-example", Node: 6}}}},
		// 13
		{Kind: syntax.KindRecord, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Person"}, {Name: syntax.AttributeRecordKey, Text: "struct"}}, Children: []ast.Child{{Role: syntax.RoleMember, Node: 14}}},
		// 14
		{Kind: syntax.KindField, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "name"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 15}}},
		// 15
		{Kind: syntax.KindNamedType, Children: []ast.Child{{Role: syntax.RoleSegment, Node: 16}, {Role: syntax.RoleSegment, Node: 17}}},
		// 16
		{Kind: syntax.KindNameSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "std"}}},
		// 17
		{Kind: syntax.KindNameSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "string"}}},
		// 18
		{Kind: syntax.KindTypeAlias, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "PeoplePage"}}, Children: []ast.Child{{Role: syntax.RoleType, Node: 19}}},
		// 19
		{Kind: syntax.KindNamedType, Children: []ast.Child{{Role: syntax.RoleSegment, Node: 20}}},
		// 20
		{Kind: syntax.KindNameSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Page"}}, Children: []ast.Child{{Role: syntax.RoleArgument, Node: 21}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "cpp-example", Node: 5}}}},
		// 21
		{Kind: syntax.KindNamedType, Children: []ast.Child{{Role: syntax.RoleSegment, Node: 22}}},
		// 22
		{Kind: syntax.KindNameSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Person"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "cpp-example", Node: 13}}}},
		// 23
		{Kind: syntax.KindRecord, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "Reader"}, {Name: syntax.AttributeRecordKey, Text: "class"}}, Children: []ast.Child{{Role: syntax.RoleTemplateParameter, Node: 24}, {Role: syntax.RoleMember, Node: 25}, {Role: syntax.RoleMember, Node: 28}}},
		// 24
		{Kind: syntax.KindTypeParameter, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}},
		// 25
		{Kind: syntax.KindMethod, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "read"}, {Name: syntax.AttributeAccess, Text: "public"}, {Name: syntax.AttributeVirtual, Text: "true"}, {Name: syntax.AttributePure, Text: "true"}}, Children: []ast.Child{{Role: syntax.RoleResult, Node: 26}}},
		// 26
		{Kind: syntax.KindNamedType, Children: []ast.Child{{Role: syntax.RoleSegment, Node: 27}}},
		// 27
		{Kind: syntax.KindNameSegment, Attributes: []ast.Attribute{{Name: syntax.AttributeName, Text: "T"}}, References: []ast.Reference{{Role: syntax.RoleTarget, Target: ast.NodeRef{Tree: "cpp-example", Node: 24}}}},
		// 28
		{Kind: syntax.KindDestructor, Attributes: []ast.Attribute{{Name: syntax.AttributeAccess, Text: "public"}, {Name: syntax.AttributeVirtual, Text: "true"}}},
	},
}

var carrier ast.AST[syntax.Language] = example

// Signature witnesses only; no stage is implemented or invoked.
var operation func(syntax.AST) (syntax.AST, error)
var stage pipeline.Stage[syntax.AST, syntax.AST] = operation
