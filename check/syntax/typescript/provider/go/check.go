package check

import (
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/internal/tree"
	lit "github.com/Bitspark/bittype/internal/typescript"
	"github.com/Bitspark/bittype/model/syntax/go"
	"github.com/Bitspark/bittype/model/syntax/typescript/lexical/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

var types = []s.Kind{s.KindTypeReference, s.KindQualifiedType, s.KindTypeApplication, s.KindKeywordType, s.KindLiteralType, s.KindArrayType, s.KindTupleType, s.KindObjectType, s.KindUnionType, s.KindIntersectionType, s.KindConditionalType, s.KindInferType, s.KindMappedType, s.KindFunctionType, s.KindTypeOperator, s.KindIndexedAccessType}
var members = []s.Kind{s.KindProperty, s.KindMethod, s.KindCallSignature, s.KindIndexSignature}
var keywords = []string{"any", "unknown", "never", "void", "undefined", "boolean", "number", "bigint", "string", "symbol", "object", "this"}

// evaluators binds the names of TypeScript's own predicates that the table uses.
var evaluators = g.Evaluators{
	lexical.Identifier(): lit.Identifier,
	lexical.TypeName():   lit.TypeName,
}

var rules = syntax.Rules[s.Language]{
	s.KindSourceFile:        {Roles: []syntax.Role[s.Language]{{Name: s.RoleImport, Max: -1, Kinds: []s.Kind{s.KindImport}}, {Name: s.RoleDeclaration, Max: -1, Kinds: []s.Kind{s.KindTypeAlias, s.KindInterface}}}},
	s.KindImport:            {Attributes: map[string]syntax.Attribute{s.AttributeModule: {Required: true, Predicate: syntax.Nonempty()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleBinding, Min: 1, Max: -1, Kinds: []s.Kind{s.KindImportBinding, s.KindNamespaceImport}}}},
	s.KindImportBinding:     {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeAlias: {Predicate: lexical.Identifier()}}},
	s.KindNamespaceImport:   {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}},
	s.KindTypeAlias:         {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.TypeName()}, s.AttributeExported: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleTypeParameter, Max: -1, Kinds: []s.Kind{s.KindTypeParameter}}, {Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindInterface:         {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.TypeName()}, s.AttributeExported: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleTypeParameter, Max: -1, Kinds: []s.Kind{s.KindTypeParameter}}, {Name: s.RoleBase, Max: -1, Kinds: []s.Kind{s.KindTypeReference, s.KindQualifiedType, s.KindTypeApplication}}, {Name: s.RoleMember, Max: -1, Kinds: members}}},
	s.KindTypeParameter:     {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.TypeName()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleConstraint, Max: 1, Kinds: types}, {Name: s.RoleDefault, Max: 1, Kinds: types}}},
	s.KindTypeReference:     {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.TypeName()}}, Reference: true},
	s.KindQualifiedType:     {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleQualifier, Min: 1, Max: 1, Kinds: []s.Kind{s.KindTypeReference, s.KindQualifiedType}}}, Reference: true},
	s.KindTypeApplication:   {Roles: []syntax.Role[s.Language]{{Name: s.RoleConstructor, Min: 1, Max: 1, Kinds: []s.Kind{s.KindTypeReference, s.KindQualifiedType}}, {Name: s.RoleArgument, Min: 1, Max: -1, Kinds: types}}},
	s.KindKeywordType:       {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, OneOf: keywords}}},
	s.KindLiteralType:       {Attributes: map[string]syntax.Attribute{s.AttributeLiteralKind: {Required: true, OneOf: []string{"string", "number", "bigint", "boolean", "null"}}, s.AttributeText: {Required: true, Predicate: syntax.Nonempty()}}},
	s.KindArrayType:         {Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindTupleType:         {Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Max: -1, Kinds: []s.Kind{s.KindTupleElement}}}},
	s.KindTupleElement:      {Attributes: map[string]syntax.Attribute{s.AttributeName: {Predicate: lexical.Identifier()}, s.AttributeOptional: {Predicate: syntax.Bool()}, s.AttributeRest: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindObjectType:        {Roles: []syntax.Role[s.Language]{{Name: s.RoleMember, Max: -1, Kinds: members}}},
	s.KindUnionType:         {Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 2, Max: -1, Kinds: types}}},
	s.KindIntersectionType:  {Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 2, Max: -1, Kinds: types}}},
	s.KindConditionalType:   {Roles: []syntax.Role[s.Language]{{Name: s.RoleCheck, Min: 1, Max: 1, Kinds: types}, {Name: s.RoleConstraint, Min: 1, Max: 1, Kinds: types}, {Name: s.RoleTrue, Min: 1, Max: 1, Kinds: types}, {Name: s.RoleFalse, Min: 1, Max: 1, Kinds: types}}},
	s.KindInferType:         {Roles: []syntax.Role[s.Language]{{Name: s.RoleTypeParameter, Min: 1, Max: 1, Kinds: []s.Kind{s.KindTypeParameter}}}},
	s.KindMappedType:        {Attributes: map[string]syntax.Attribute{s.AttributeReadonlyModifier: {OneOf: []string{"add", "remove"}}, s.AttributeOptionalModifier: {OneOf: []string{"add", "remove"}}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleTypeParameter, Min: 1, Max: 1, Kinds: []s.Kind{s.KindTypeParameter}}, {Name: s.RoleNameType, Max: 1, Kinds: types}, {Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindFunctionType:      {Roles: []syntax.Role[s.Language]{{Name: s.RoleTypeParameter, Max: -1, Kinds: []s.Kind{s.KindTypeParameter}}, {Name: s.RoleParameter, Max: -1, Kinds: []s.Kind{s.KindParameter}}, {Name: s.RoleResult, Min: 1, Max: 1, Kinds: types}}},
	s.KindParameter:         {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeOptional: {Predicate: syntax.Bool()}, s.AttributeRest: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindProperty:          {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: syntax.String()}, s.AttributeOptional: {Predicate: syntax.Bool()}, s.AttributeReadonly: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindMethod:            {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: syntax.String()}, s.AttributeOptional: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleSignature, Min: 1, Max: 1, Kinds: []s.Kind{s.KindFunctionType}}}},
	s.KindCallSignature:     {Roles: []syntax.Role[s.Language]{{Name: s.RoleSignature, Min: 1, Max: 1, Kinds: []s.Kind{s.KindFunctionType}}}},
	s.KindIndexSignature:    {Attributes: map[string]syntax.Attribute{s.AttributeReadonly: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleParameter, Min: 1, Max: 1, Kinds: []s.Kind{s.KindParameter}}, {Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindTypeOperator:      {Attributes: map[string]syntax.Attribute{s.AttributeOperator: {Required: true, OneOf: []string{"keyof", "readonly", "unique"}}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindIndexedAccessType: {Roles: []syntax.Role[s.Language]{{Name: s.RoleObject, Min: 1, Max: 1, Kinds: types}, {Name: s.RoleIndex, Min: 1, Max: 1, Kinds: types}}},
}

// bittype:grammar rules
// bittype:scope docs/syntax.md#typescript
func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, s.KindSourceFile, rules, evaluators); err != nil {
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
