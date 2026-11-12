package check

import (
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	native "github.com/Bitspark/bittype/internal/cpp"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"strings"
)

var types = []s.Kind{s.KindNamedType, s.KindBuiltinType, s.KindPointerType, s.KindLValueReferenceType, s.KindRValueReferenceType, s.KindCVQualifiedType, s.KindArrayType, s.KindFunctionType, s.KindMemberPointerType}
var declarations = []s.Kind{s.KindNamespace, s.KindRecord, s.KindTypeAlias, s.KindEnum}
var members = []s.Kind{s.KindField, s.KindMethod, s.KindDestructor, s.KindTypeAlias, s.KindRecord, s.KindEnum}
var access = g.Set("public", "protected", "private")

func named() map[string]g.Attribute {
	return g.Attributes(s.AttributeName, g.Attr(true, native.Identifier), s.AttributeAccess, g.Attr(false, access))
}
func extra(base map[string]g.Attribute, pairs ...any) map[string]g.Attribute {
	for k, v := range g.Attributes(pairs...) {
		base[k] = v
	}
	return base
}

var rules = g.Rules[s.Language]{
	s.KindTranslationUnit:     {Roles: []g.Role[s.Language]{g.Many(s.RoleInclude, 0, s.KindInclude), g.Many(s.RoleDeclaration, 0, declarations...)}},
	s.KindInclude:             {Attributes: g.Attributes(s.AttributePath, g.Attr(true, func(v string) bool { return g.Nonempty(v) && !strings.ContainsAny(v, "\r\n") }), s.AttributeSystem, g.Attr(false, g.Bool))},
	s.KindNamespace:           {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.Identifier)), Roles: []g.Role[s.Language]{g.Many(s.RoleDeclaration, 0, declarations...)}},
	s.KindRecord:              {Attributes: extra(named(), s.AttributeRecordKey, g.Attr(true, g.Set("struct", "class")), s.AttributeForward, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.Many(s.RoleTemplateParameter, 0, s.KindTypeParameter), g.Many(s.RoleBase, 0, s.KindBase), g.Many(s.RoleMember, 0, members...)}},
	s.KindTypeAlias:           {Attributes: named(), Roles: []g.Role[s.Language]{g.Many(s.RoleTemplateParameter, 0, s.KindTypeParameter), g.One(s.RoleType, types...)}},
	s.KindEnum:                {Attributes: extra(named(), s.AttributeScoped, g.Attr(false, g.Bool), s.AttributeForward, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.Opt(s.RoleUnderlying, types...), g.Many(s.RoleEnumerator, 0, s.KindEnumerator)}},
	s.KindEnumerator:          {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.Identifier), s.AttributeValue, g.Attr(false, g.Signed))},
	s.KindBase:                {Attributes: g.Attributes(s.AttributeAccess, g.Attr(false, access), s.AttributeVirtual, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.One(s.RoleType, s.KindNamedType)}},
	s.KindField:               {Attributes: named(), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindMethod:              {Attributes: extra(named(), s.AttributeConst, g.Attr(false, g.Bool), s.AttributeVolatile, g.Attr(false, g.Bool), s.AttributeVirtual, g.Attr(false, g.Bool), s.AttributePure, g.Attr(false, g.Bool), s.AttributeNoexcept, g.Attr(false, g.Bool), s.AttributeRefQualifier, g.Attr(false, g.Set("none", "lvalue", "rvalue")), s.AttributeVariadic, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.Many(s.RoleParameter, 0, s.KindParameter), g.One(s.RoleResult, types...)}},
	s.KindDestructor:          {Attributes: g.Attributes(s.AttributeAccess, g.Attr(false, access), s.AttributeVirtual, g.Attr(false, g.Bool), s.AttributePure, g.Attr(false, g.Bool), s.AttributeNoexcept, g.Attr(false, g.Bool))},
	s.KindParameter:           {Attributes: g.Attributes(s.AttributeName, g.Attr(false, native.Identifier)), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindTypeParameter:       {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.Identifier), s.AttributeParameterKey, g.Attr(false, g.Set("typename", "class"))), Roles: []g.Role[s.Language]{g.Opt(s.RoleDefault, types...)}},
	s.KindNamedType:           {Attributes: g.Attributes(s.AttributeGlobal, g.Attr(false, g.Bool), s.AttributeTypename, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.Many(s.RoleSegment, 1, s.KindNameSegment)}},
	s.KindNameSegment:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.Identifier), s.AttributeTemplate, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.Many(s.RoleArgument, 0, types...)}, Reference: true},
	s.KindBuiltinType:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.Builtin))},
	s.KindPointerType:         {Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindLValueReferenceType: {Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindRValueReferenceType: {Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindCVQualifiedType:     {Attributes: g.Attributes(s.AttributeConst, g.Attr(false, g.Bool), s.AttributeVolatile, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindArrayType:           {Attributes: g.Attributes(s.AttributeLength, g.Attr(true, g.Positive)), Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindFunctionType:        {Attributes: g.Attributes(s.AttributeNoexcept, g.Attr(false, g.Bool), s.AttributeVariadic, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.Many(s.RoleParameter, 0, s.KindParameter), g.One(s.RoleResult, types...)}},
	s.KindMemberPointerType:   {Roles: []g.Role[s.Language]{g.One(s.RoleOwner, s.KindNamedType), g.One(s.RoleElement, types...)}},
}

// bittype:grammar rules
// bittype:scope docs/syntax.md#c
func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, s.KindTranslationUnit, rules); err != nil {
		return err
	}
	v := tree.View[s.Language]{Tree: t}
	parents := map[ast.NodeID]ast.NodeID{}
	for i, n := range t.Nodes {
		for _, c := range n.Children {
			parents[c.Node] = ast.NodeID(i + 1)
		}
	}
	for i, n := range t.Nodes {
		id := ast.NodeID(i + 1)
		bad := func(m string) error { return text.AtNode(t, id, m) }
		if v.Has(id, s.AttributeAccess) && n.Kind != s.KindBase {
			parent := parents[id]
			if parent == 0 || v.Node(parent).Kind != s.KindRecord {
				return bad("access belongs to record members")
			}
		}
		if v.Flag(id, s.AttributePure) && !v.Flag(id, s.AttributeVirtual) {
			return bad("pure declaration must be virtual")
		}
		switch n.Kind {
		case s.KindInclude:
			path := v.Attr(id, s.AttributePath)
			if v.Flag(id, s.AttributeSystem) {
				if strings.ContainsRune(path, '>') {
					return bad("invalid system header name")
				}
			} else if strings.ContainsRune(path, '"') {
				return bad("invalid quoted header name")
			}
		case s.KindRecord:
			if v.Flag(id, s.AttributeForward) && (len(v.Children(id, s.RoleBase)) > 0 || len(v.Children(id, s.RoleMember)) > 0) {
				return bad("forward record cannot have bases or members")
			}
		case s.KindEnum:
			if v.Flag(id, s.AttributeForward) && (len(v.Children(id, s.RoleEnumerator)) > 0 || (!v.Flag(id, s.AttributeScoped) && v.Child(id, s.RoleUnderlying) == 0)) {
				return bad("invalid forward enum")
			}
		case s.KindCVQualifiedType:
			if !v.Flag(id, s.AttributeConst) && !v.Flag(id, s.AttributeVolatile) {
				return bad("qualification requires const or volatile")
			}
			e := v.Node(v.Child(id, s.RoleElement)).Kind
			if e != s.KindNamedType && e != s.KindBuiltinType && e != s.KindPointerType && e != s.KindMemberPointerType {
				return bad("qualifiers require a directly representable object or pointer type")
			}
		case s.KindFunctionType, s.KindMethod:
			r := v.Node(v.Child(id, s.RoleResult)).Kind
			if r == s.KindArrayType || r == s.KindFunctionType {
				return bad("function cannot return an array or function type")
			}
		case s.KindArrayType:
			e := v.Node(v.Child(id, s.RoleElement)).Kind
			if e == s.KindFunctionType || e == s.KindLValueReferenceType || e == s.KindRValueReferenceType {
				return bad("array element cannot be a function or reference")
			}
		case s.KindPointerType, s.KindMemberPointerType, s.KindLValueReferenceType, s.KindRValueReferenceType:
			e := v.Node(v.Child(id, s.RoleElement)).Kind
			if e == s.KindLValueReferenceType || e == s.KindRValueReferenceType {
				return bad("nested reference declarator is not representable")
			}
		}
	}
	return nil
}

type Provider struct{ s.Language }

func (Provider) Check(view ast.View[s.Language]) error { return Check(view) }

var _ checkcontract.Provider[s.Language] = Provider{}
