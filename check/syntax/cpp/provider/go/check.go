package check

import (
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	native "github.com/Bitspark/bittype/internal/cpp"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/internal/tree"
	"github.com/Bitspark/bittype/model/syntax/cpp/lexical/go"
	"github.com/Bitspark/bittype/model/syntax/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"strings"
)

var types = []s.Kind{s.KindNamedType, s.KindBuiltinType, s.KindPointerType, s.KindLValueReferenceType, s.KindRValueReferenceType, s.KindCVQualifiedType, s.KindArrayType, s.KindFunctionType, s.KindMemberPointerType}
var declarations = []s.Kind{s.KindNamespace, s.KindRecord, s.KindTypeAlias, s.KindEnum}
var members = []s.Kind{s.KindField, s.KindMethod, s.KindDestructor, s.KindTypeAlias, s.KindRecord, s.KindEnum}
var access = []string{"public", "protected", "private"}

func headerName(v string) bool { return g.Nonempty(v) && !strings.ContainsAny(v, "\r\n") }

// evaluators binds the names of C++'s own predicates that the table uses.
var evaluators = g.Evaluators{
	lexical.Identifier():      native.Identifier,
	lexical.FundamentalType(): native.FundamentalType,
	lexical.HeaderName():      headerName,
}

var rules = syntax.Rules[s.Language]{
	s.KindTranslationUnit:     {Roles: []syntax.Role[s.Language]{{Name: s.RoleInclude, Max: -1, Kinds: []s.Kind{s.KindInclude}}, {Name: s.RoleDeclaration, Max: -1, Kinds: declarations}}},
	s.KindInclude:             {Attributes: map[string]syntax.Attribute{s.AttributePath: {Required: true, Predicate: lexical.HeaderName()}, s.AttributeSystem: {Predicate: syntax.Bool()}}},
	s.KindNamespace:           {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleDeclaration, Max: -1, Kinds: declarations}}},
	s.KindRecord:              {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeAccess: {OneOf: access}, s.AttributeRecordKey: {Required: true, OneOf: []string{"struct", "class"}}, s.AttributeForward: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleTemplateParameter, Max: -1, Kinds: []s.Kind{s.KindTypeParameter}}, {Name: s.RoleBase, Max: -1, Kinds: []s.Kind{s.KindBase}}, {Name: s.RoleMember, Max: -1, Kinds: members}}},
	s.KindTypeAlias:           {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeAccess: {OneOf: access}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleTemplateParameter, Max: -1, Kinds: []s.Kind{s.KindTypeParameter}}, {Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindEnum:                {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeAccess: {OneOf: access}, s.AttributeScoped: {Predicate: syntax.Bool()}, s.AttributeForward: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleUnderlying, Max: 1, Kinds: types}, {Name: s.RoleEnumerator, Max: -1, Kinds: []s.Kind{s.KindEnumerator}}}},
	s.KindEnumerator:          {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeValue: {Predicate: syntax.Signed()}}},
	s.KindBase:                {Attributes: map[string]syntax.Attribute{s.AttributeAccess: {OneOf: access}, s.AttributeVirtual: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: []s.Kind{s.KindNamedType}}}},
	s.KindField:               {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeAccess: {OneOf: access}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindMethod:              {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeAccess: {OneOf: access}, s.AttributeConst: {Predicate: syntax.Bool()}, s.AttributeVolatile: {Predicate: syntax.Bool()}, s.AttributeVirtual: {Predicate: syntax.Bool()}, s.AttributePure: {Predicate: syntax.Bool()}, s.AttributeNoexcept: {Predicate: syntax.Bool()}, s.AttributeRefQualifier: {OneOf: []string{"none", "lvalue", "rvalue"}}, s.AttributeVariadic: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleParameter, Max: -1, Kinds: []s.Kind{s.KindParameter}}, {Name: s.RoleResult, Min: 1, Max: 1, Kinds: types}}},
	s.KindDestructor:          {Attributes: map[string]syntax.Attribute{s.AttributeAccess: {OneOf: access}, s.AttributeVirtual: {Predicate: syntax.Bool()}, s.AttributePure: {Predicate: syntax.Bool()}, s.AttributeNoexcept: {Predicate: syntax.Bool()}}},
	s.KindParameter:           {Attributes: map[string]syntax.Attribute{s.AttributeName: {Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindTypeParameter:       {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeParameterKey: {OneOf: []string{"typename", "class"}}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleDefault, Max: 1, Kinds: types}}},
	s.KindNamedType:           {Attributes: map[string]syntax.Attribute{s.AttributeGlobal: {Predicate: syntax.Bool()}, s.AttributeTypename: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleSegment, Min: 1, Max: -1, Kinds: []s.Kind{s.KindNameSegment}}}},
	s.KindNameSegment:         {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeTemplate: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleArgument, Max: -1, Kinds: types}}, Reference: true},
	s.KindBuiltinType:         {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.FundamentalType()}}},
	s.KindPointerType:         {Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindLValueReferenceType: {Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindRValueReferenceType: {Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindCVQualifiedType:     {Attributes: map[string]syntax.Attribute{s.AttributeConst: {Predicate: syntax.Bool()}, s.AttributeVolatile: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindArrayType:           {Attributes: map[string]syntax.Attribute{s.AttributeLength: {Required: true, Predicate: syntax.Positive()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindFunctionType:        {Attributes: map[string]syntax.Attribute{s.AttributeNoexcept: {Predicate: syntax.Bool()}, s.AttributeVariadic: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleParameter, Max: -1, Kinds: []s.Kind{s.KindParameter}}, {Name: s.RoleResult, Min: 1, Max: 1, Kinds: types}}},
	s.KindMemberPointerType:   {Roles: []syntax.Role[s.Language]{{Name: s.RoleOwner, Min: 1, Max: 1, Kinds: []s.Kind{s.KindNamedType}}, {Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
}

// bittype:grammar rules
// bittype:scope docs/syntax.md#c
func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, s.KindTranslationUnit, rules, evaluators); err != nil {
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
