package check

import (
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/internal/tree"
	lit "github.com/Bitspark/bittype/internal/typescript"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func identifier(v string) bool { return lit.Identifier(v) }

var types = []s.Kind{s.KindTypeReference, s.KindQualifiedType, s.KindTypeApplication, s.KindKeywordType, s.KindLiteralType, s.KindArrayType, s.KindTupleType, s.KindObjectType, s.KindUnionType, s.KindIntersectionType, s.KindConditionalType, s.KindInferType, s.KindMappedType, s.KindFunctionType, s.KindTypeOperator, s.KindIndexedAccessType}
var members = []s.Kind{s.KindProperty, s.KindMethod, s.KindCallSignature, s.KindIndexSignature}
var keywords = g.Set("any", "unknown", "never", "void", "undefined", "boolean", "number", "bigint", "string", "symbol", "object", "this")
var rules = g.Rules[s.Language]{
	s.KindSourceFile:        {Roles: []g.Role[s.Language]{g.Many(s.RoleImport, 0, s.KindImport), g.Many(s.RoleDeclaration, 0, s.KindTypeAlias, s.KindInterface)}},
	s.KindImport:            {Attributes: g.Attributes(s.AttributeModule, g.Attr(true, g.Nonempty)), Roles: []g.Role[s.Language]{g.Many(s.RoleBinding, 1, s.KindImportBinding, s.KindNamespaceImport)}},
	s.KindImportBinding:     {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier), s.AttributeAlias, g.Attr(false, identifier))},
	s.KindNamespaceImport:   {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier))},
	s.KindTypeAlias:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, lit.TypeName), s.AttributeExported, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.Many(s.RoleTypeParameter, 0, s.KindTypeParameter), g.One(s.RoleType, types...)}},
	s.KindInterface:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, lit.TypeName), s.AttributeExported, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.Many(s.RoleTypeParameter, 0, s.KindTypeParameter), g.Many(s.RoleBase, 0, s.KindTypeReference, s.KindQualifiedType, s.KindTypeApplication), g.Many(s.RoleMember, 0, members...)}},
	s.KindTypeParameter:     {Attributes: g.Attributes(s.AttributeName, g.Attr(true, lit.TypeName)), Roles: []g.Role[s.Language]{g.Opt(s.RoleConstraint, types...), g.Opt(s.RoleDefault, types...)}},
	s.KindTypeReference:     {Attributes: g.Attributes(s.AttributeName, g.Attr(true, lit.TypeName)), Reference: true},
	s.KindQualifiedType:     {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier)), Roles: []g.Role[s.Language]{g.One(s.RoleQualifier, s.KindTypeReference, s.KindQualifiedType)}, Reference: true},
	s.KindTypeApplication:   {Roles: []g.Role[s.Language]{g.One(s.RoleConstructor, s.KindTypeReference, s.KindQualifiedType), g.Many(s.RoleArgument, 1, types...)}},
	s.KindKeywordType:       {Attributes: g.Attributes(s.AttributeName, g.Attr(true, keywords))},
	s.KindLiteralType:       {Attributes: g.Attributes(s.AttributeLiteralKind, g.Attr(true, g.Set("string", "number", "bigint", "boolean", "null")), s.AttributeText, g.Attr(true, g.Nonempty))},
	s.KindArrayType:         {Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindTupleType:         {Roles: []g.Role[s.Language]{g.Many(s.RoleElement, 0, s.KindTupleElement)}},
	s.KindTupleElement:      {Attributes: g.Attributes(s.AttributeName, g.Attr(false, identifier), s.AttributeOptional, g.Attr(false, g.Bool), s.AttributeRest, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindObjectType:        {Roles: []g.Role[s.Language]{g.Many(s.RoleMember, 0, members...)}},
	s.KindUnionType:         {Roles: []g.Role[s.Language]{g.Many(s.RoleType, 2, types...)}},
	s.KindIntersectionType:  {Roles: []g.Role[s.Language]{g.Many(s.RoleType, 2, types...)}},
	s.KindConditionalType:   {Roles: []g.Role[s.Language]{g.One(s.RoleCheck, types...), g.One(s.RoleConstraint, types...), g.One(s.RoleTrue, types...), g.One(s.RoleFalse, types...)}},
	s.KindInferType:         {Roles: []g.Role[s.Language]{g.One(s.RoleTypeParameter, s.KindTypeParameter)}},
	s.KindMappedType:        {Attributes: g.Attributes(s.AttributeReadonlyModifier, g.Attr(false, g.Set("add", "remove")), s.AttributeOptionalModifier, g.Attr(false, g.Set("add", "remove"))), Roles: []g.Role[s.Language]{g.One(s.RoleTypeParameter, s.KindTypeParameter), g.Opt(s.RoleNameType, types...), g.One(s.RoleType, types...)}},
	s.KindFunctionType:      {Roles: []g.Role[s.Language]{g.Many(s.RoleTypeParameter, 0, s.KindTypeParameter), g.Many(s.RoleParameter, 0, s.KindParameter), g.One(s.RoleResult, types...)}},
	s.KindParameter:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, identifier), s.AttributeOptional, g.Attr(false, g.Bool), s.AttributeRest, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindProperty:          {Attributes: g.Attributes(s.AttributeName, g.Attr(true, g.String), s.AttributeOptional, g.Attr(false, g.Bool), s.AttributeReadonly, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindMethod:            {Attributes: g.Attributes(s.AttributeName, g.Attr(true, g.String), s.AttributeOptional, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.One(s.RoleSignature, s.KindFunctionType)}},
	s.KindCallSignature:     {Roles: []g.Role[s.Language]{g.One(s.RoleSignature, s.KindFunctionType)}},
	s.KindIndexSignature:    {Attributes: g.Attributes(s.AttributeReadonly, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.One(s.RoleParameter, s.KindParameter), g.One(s.RoleType, types...)}},
	s.KindTypeOperator:      {Attributes: g.Attributes(s.AttributeOperator, g.Attr(true, g.Set("keyof", "readonly", "unique"))), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindIndexedAccessType: {Roles: []g.Role[s.Language]{g.One(s.RoleObject, types...), g.One(s.RoleIndex, types...)}},
}

// bittype:grammar rules
// bittype:scope docs/syntax.md#typescript
func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, s.KindSourceFile, rules); err != nil {
		return err
	}
	v := tree.View[s.Language]{Tree: t}
	// Infer's position is a syntactic restriction. This walk neither resolves
	// its binder nor decides what a conditional or mapped type evaluates to.
	var context func(ast.NodeID, bool) error
	context = func(id ast.NodeID, inConstraint bool) error {
		n := v.Node(id)
		if n.Kind == s.KindInferType && !inConstraint {
			return text.AtNode(t, id, "infer is only supported within a conditional extends operand")
		}
		for _, child := range n.Children {
			if err := context(child.Node, inConstraint || (n.Kind == s.KindConditionalType && child.Role == s.RoleConstraint)); err != nil {
				return err
			}
		}
		return nil
	}
	if err := context(t.Root, false); err != nil {
		return err
	}
	for i, n := range t.Nodes {
		id := ast.NodeID(i + 1)
		bad := func(m string) error { return text.AtNode(t, id, m) }
		switch n.Kind {
		case s.KindMappedType, s.KindInferType:
			parameter := v.Child(id, s.RoleTypeParameter)
			if v.Child(parameter, s.RoleDefault) != 0 {
				return text.AtNode(t, parameter, "mapped and inferred parameters cannot have defaults")
			}
			if n.Kind == s.KindMappedType && v.Child(parameter, s.RoleConstraint) == 0 {
				return text.AtNode(t, parameter, "mapped parameter requires an in type")
			}
		case s.KindImport:
			b := v.Children(id, s.RoleBinding)
			for _, x := range b {
				if v.Node(x).Kind == s.KindNamespaceImport && len(b) != 1 {
					return bad("namespace import must be the only binding")
				}
			}
		case s.KindLiteralType:
			value := v.Attr(id, s.AttributeText)
			valid := false
			switch v.Attr(id, s.AttributeLiteralKind) {
			case "string":
				_, err := lit.Unquote(value)
				valid = err == nil
			case "number":
				valid = lit.Number(value)
			case "bigint":
				valid = lit.Bigint(value)
			case "boolean":
				valid = value == "true" || value == "false"
			case "null":
				valid = value == "null"
			}
			if !valid {
				return bad("literal kind and token do not agree")
			}
		case s.KindTypeOperator:
			child := v.Child(id, s.RoleType)
			op := v.Attr(id, s.AttributeOperator)
			if op == "readonly" && v.Node(child).Kind != s.KindArrayType && v.Node(child).Kind != s.KindTupleType {
				return bad("readonly applies to array or tuple syntax")
			}
			if op == "unique" && (v.Node(child).Kind != s.KindKeywordType || v.Attr(child, s.AttributeName) != "symbol") {
				return bad("unique applies to symbol")
			}
		case s.KindFunctionType, s.KindTupleType:
			role := s.RoleParameter
			if n.Kind == s.KindTupleType {
				role = s.RoleElement
			}
			list := v.Children(id, role)
			optional, rest := false, false
			named := -1
			for j, p := range list {
				o, r := v.Flag(p, s.AttributeOptional), v.Flag(p, s.AttributeRest)
				if r && o {
					return text.AtNode(t, p, "rest cannot be optional")
				}
				if n.Kind == s.KindFunctionType && r && j != len(list)-1 {
					return text.AtNode(t, p, "rest parameter must be last")
				}
				if r && rest {
					return text.AtNode(t, p, "multiple rest elements")
				}
				if o && rest {
					return text.AtNode(t, p, "optional element after rest")
				}
				if !o && !r && optional {
					return text.AtNode(t, p, "required element after optional")
				}
				if o {
					optional = true
				}
				if r {
					rest = true
				}
				if n.Kind == s.KindTupleType {
					mode := 0
					if v.Has(p, s.AttributeName) {
						mode = 1
					}
					if named >= 0 && named != mode {
						return text.AtNode(t, p, "tuple labels must be consistent")
					}
					named = mode
				}
			}
		case s.KindIndexSignature:
			p := v.Child(id, s.RoleParameter)
			if v.Flag(p, s.AttributeOptional) || v.Flag(p, s.AttributeRest) {
				return bad("index parameter cannot be optional or rest")
			}
		}
	}
	return nil
}

type Provider struct{ s.Language }

func (Provider) Check(view ast.View[s.Language]) error { return Check(view) }

var _ checkcontract.Provider[s.Language] = Provider{}
