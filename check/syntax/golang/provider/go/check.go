package check

import (
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/internal/tree"
	"github.com/Bitspark/bittype/model/syntax/go"
	"github.com/Bitspark/bittype/model/syntax/golang/lexical/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"go/token"
)

func identifier(v string) bool { return g.Identifier(v) && !token.Lookup(v).IsKeyword() }
func importName(v string) bool { return v == "." || identifier(v) }

// evaluators binds the names of Go's own predicates that the table uses.
var evaluators = g.Evaluators{
	lexical.Identifier(): identifier,
	lexical.ImportName(): importName,
}

var types = []s.Kind{s.KindNamedType, s.KindQualifiedType, s.KindTypeApplication, s.KindStructType, s.KindInterfaceType, s.KindFunctionType, s.KindPointerType, s.KindSliceType, s.KindArrayType, s.KindMapType, s.KindChannelType}
var constraints = append(append([]s.Kind{}, types...), s.KindTypeUnion, s.KindTypeApproximation)
var rules = syntax.Rules[s.Language]{
	s.KindFile:              {Attributes: map[string]syntax.Attribute{s.AttributePackage: {Required: true, Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleImport, Max: -1, Kinds: []s.Kind{s.KindImport}}, {Name: s.RoleDeclaration, Max: -1, Kinds: []s.Kind{s.KindTypeDeclaration, s.KindTypeAlias}}}},
	s.KindImport:            {Attributes: map[string]syntax.Attribute{s.AttributePath: {Required: true, Predicate: syntax.Nonempty()}, s.AttributeAlias: {Predicate: lexical.ImportName()}}},
	s.KindTypeDeclaration:   {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleTypeParameter, Max: -1, Kinds: []s.Kind{s.KindTypeParameter}}, {Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindTypeAlias:         {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleTypeParameter, Max: -1, Kinds: []s.Kind{s.KindTypeParameter}}, {Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindTypeParameter:     {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleConstraint, Min: 1, Max: 1, Kinds: constraints}}},
	s.KindNamedType:         {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}, Reference: true},
	s.KindQualifiedType:     {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeQualifier: {Required: true, Predicate: lexical.Identifier()}}, Reference: true},
	s.KindTypeApplication:   {Roles: []syntax.Role[s.Language]{{Name: s.RoleConstructor, Min: 1, Max: 1, Kinds: []s.Kind{s.KindNamedType, s.KindQualifiedType}}, {Name: s.RoleArgument, Min: 1, Max: -1, Kinds: types}}},
	s.KindStructType:        {Roles: []syntax.Role[s.Language]{{Name: s.RoleField, Max: -1, Kinds: []s.Kind{s.KindField, s.KindEmbeddedField}}}},
	s.KindField:             {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeTag: {Predicate: syntax.String()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindEmbeddedField:     {Attributes: map[string]syntax.Attribute{s.AttributeTag: {Predicate: syntax.String()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: []s.Kind{s.KindNamedType, s.KindQualifiedType, s.KindTypeApplication, s.KindPointerType}}}},
	s.KindInterfaceType:     {Roles: []syntax.Role[s.Language]{{Name: s.RoleMember, Max: -1, Kinds: []s.Kind{s.KindMethod, s.KindEmbeddedType}}}},
	s.KindMethod:            {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleSignature, Min: 1, Max: 1, Kinds: []s.Kind{s.KindFunctionType}}}},
	s.KindEmbeddedType:      {Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: constraints}}},
	s.KindFunctionType:      {Roles: []syntax.Role[s.Language]{{Name: s.RoleParameter, Max: -1, Kinds: []s.Kind{s.KindParameter}}, {Name: s.RoleResult, Max: -1, Kinds: []s.Kind{s.KindParameter}}}},
	s.KindParameter:         {Attributes: map[string]syntax.Attribute{s.AttributeName: {Predicate: lexical.Identifier()}, s.AttributeVariadic: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindPointerType:       {Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindSliceType:         {Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindArrayType:         {Attributes: map[string]syntax.Attribute{s.AttributeLength: {Required: true, Predicate: syntax.Decimal()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindMapType:           {Roles: []syntax.Role[s.Language]{{Name: s.RoleKey, Min: 1, Max: 1, Kinds: types}, {Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindChannelType:       {Attributes: map[string]syntax.Attribute{s.AttributeDirection: {OneOf: []string{"both", "send", "receive"}}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindTypeUnion:         {Roles: []syntax.Role[s.Language]{{Name: s.RoleTerm, Min: 2, Max: -1, Kinds: append(types, s.KindTypeApproximation)}}},
	s.KindTypeApproximation: {Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
}

// bittype:grammar rules
// bittype:scope docs/syntax.md#go
func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, s.KindFile, rules, evaluators); err != nil {
		return err
	}
	v := tree.View[s.Language]{Tree: t}
	for i, n := range t.Nodes {
		id := ast.NodeID(i + 1)
		if n.Kind == s.KindFunctionType {
			for _, role := range []string{s.RoleParameter, s.RoleResult} {
				params := v.Children(id, role)
				named := -1
				for j, p := range params {
					mode := 0
					if v.Has(p, s.AttributeName) {
						mode = 1
					}
					if named >= 0 && named != mode {
						return text.AtNode(t, p, "cannot mix named and unnamed parameters")
					}
					named = mode
					if v.Flag(p, s.AttributeVariadic) && (role != s.RoleParameter || j != len(params)-1) {
						return text.AtNode(t, p, "variadic parameter must be final input")
					}
				}
			}
		}
		if n.Kind == s.KindEmbeddedField {
			e := v.Child(id, s.RoleType)
			if v.Node(e).Kind == s.KindPointerType {
				kind := v.Node(v.Child(e, s.RoleElement)).Kind
				if kind != s.KindNamedType && kind != s.KindQualifiedType && kind != s.KindTypeApplication {
					return text.AtNode(t, e, "embedded pointer requires a named type")
				}
			}
		}
	}
	return nil
}

type Provider struct{ s.Language }

func (Provider) Check(v ast.View[s.Language]) error { return Check(v) }

var _ checkcontract.Provider[s.Language] = Provider{}
