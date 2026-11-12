package check

import (
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"go/token"
)

func identifier(v string) bool { return g.Identifier(v) && !token.Lookup(v).IsKeyword() }

var types = []s.Kind{s.KindNamedType, s.KindQualifiedType, s.KindTypeApplication, s.KindStructType, s.KindInterfaceType, s.KindFunctionType, s.KindPointerType, s.KindSliceType, s.KindArrayType, s.KindMapType, s.KindChannelType}
var constraints = append(append([]s.Kind{}, types...), s.KindTypeUnion, s.KindTypeApproximation)
var rules = g.Rules[s.Language]{
	s.KindFile:              {Attributes: g.Attributes(s.AttributePackage, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.Many(s.RoleImport, 0, s.KindImport), g.Many(s.RoleDeclaration, 0, s.KindTypeDeclaration, s.KindTypeAlias)}},
	s.KindImport:            {Attributes: g.Attributes(s.AttributePath, g.Attr(true, g.Nonempty), s.AttributeAlias, g.Attr(false, func(v string) bool { return v == "." || identifier(v) }))},
	s.KindTypeDeclaration:   {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.Many(s.RoleTypeParameter, 0, s.KindTypeParameter), g.One(s.RoleType, types...)}},
	s.KindTypeAlias:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.Many(s.RoleTypeParameter, 0, s.KindTypeParameter), g.One(s.RoleType, types...)}},
	s.KindTypeParameter:     {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.One(s.RoleConstraint, constraints...)}},
	s.KindNamedType:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Reference: true},
	s.KindQualifiedType:     {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier), s.AttributeQualifier, g.Attr(true, identifier)), Reference: true},
	s.KindTypeApplication:   {Roles: []g.Role[s.Language]{g.One(s.RoleConstructor, s.KindNamedType, s.KindQualifiedType), g.Many(s.RoleArgument, 1, types...)}},
	s.KindStructType:        {Roles: []g.Role[s.Language]{g.Many(s.RoleField, 0, s.KindField, s.KindEmbeddedField)}},
	s.KindField:             {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier), s.AttributeTag, g.Attr(false, g.String)), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindEmbeddedField:     {Attributes: g.Attributes(s.AttributeTag, g.Attr(false, g.String)), Roles: []g.Role[s.Language]{g.One(s.RoleType, s.KindNamedType, s.KindQualifiedType, s.KindTypeApplication, s.KindPointerType)}},
	s.KindInterfaceType:     {Roles: []g.Role[s.Language]{g.Many(s.RoleMember, 0, s.KindMethod, s.KindEmbeddedType)}},
	s.KindMethod:            {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.One(s.RoleSignature, s.KindFunctionType)}},
	s.KindEmbeddedType:      {Roles: []g.Role[s.Language]{g.One(s.RoleType, constraints...)}},
	s.KindFunctionType:      {Roles: []g.Role[s.Language]{g.Many(s.RoleParameter, 0, s.KindParameter), g.Many(s.RoleResult, 0, s.KindParameter)}},
	s.KindParameter:         {Attributes: g.Attributes(s.AttributeName, g.Attr(false, identifier), s.AttributeVariadic, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindPointerType:       {Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindSliceType:         {Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindArrayType:         {Attributes: g.Attributes(s.AttributeLength, g.Attr(true, g.Decimal)), Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindMapType:           {Roles: []g.Role[s.Language]{g.One(s.RoleKey, types...), g.One(s.RoleElement, types...)}},
	s.KindChannelType:       {Attributes: g.Attributes(s.AttributeDirection, g.Attr(false, g.Set("both", "send", "receive"))), Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindTypeUnion:         {Roles: []g.Role[s.Language]{g.Many(s.RoleTerm, 2, append(types, s.KindTypeApproximation)...)}},
	s.KindTypeApproximation: {Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
}

// bittype:grammar rules
// bittype:scope docs/syntax.md#go
func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, s.KindFile, rules); err != nil {
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
