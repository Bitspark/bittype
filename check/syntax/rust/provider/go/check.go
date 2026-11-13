package check

import (
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	g "github.com/Bitspark/bittype/internal/grammar"
	native "github.com/Bitspark/bittype/internal/rust"
	"github.com/Bitspark/bittype/internal/tree"
	"github.com/Bitspark/bittype/model/syntax/go"
	"github.com/Bitspark/bittype/model/syntax/rust/lexical/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"slices"
)

var types = []s.Kind{s.KindPathType, s.KindQualifiedPathType, s.KindPrimitiveType, s.KindReferenceType, s.KindRawPointerType, s.KindSliceType, s.KindArrayType, s.KindTupleType, s.KindNeverType, s.KindFunctionType, s.KindTraitObjectType}
var bounds = []s.Kind{s.KindPathType, s.KindLifetimeReference}
var typeBounds = []s.Kind{s.KindPathType, s.KindLifetimeReference, s.KindRelaxedBound}
var generics = []s.Kind{s.KindTypeParameter, s.KindLifetimeParameter}
var items = []s.Kind{s.KindModule, s.KindUse, s.KindTypeAlias, s.KindStruct, s.KindTupleStruct, s.KindEnum, s.KindTrait, s.KindTraitImpl}
var visibility = []string{"private", "public", "crate"}

// abis are the names of the calling conventions.
var abis = []string{"Rust", "C", "C-unwind", "system", "system-unwind", "cdecl", "stdcall", "fastcall", "vectorcall", "thiscall", "aapcs", "win64", "sysv64", "efiapi"}
var primitive = []string{"bool", "char", "str", "f32", "f64", "i8", "i16", "i32", "i64", "i128", "isize", "u8", "u16", "u32", "u64", "u128", "usize"}

// evaluators binds the names of Rust's own predicates that the table uses.
var evaluators = g.Evaluators{
	lexical.Identifier():             native.Identifier,
	lexical.IdentifierOrUnderscore(): native.IdentifierOrUnderscore,
	lexical.LifetimeParameter():      native.LifetimeParameter,
	lexical.Lifetime():               native.Lifetime,
	lexical.PathSegment():            native.PathSegment,
}

var rules = syntax.Rules[s.Language]{
	s.KindFile:                  {Roles: []syntax.Role[s.Language]{{Name: s.RoleItem, Max: -1, Kinds: items}}},
	s.KindModule:                {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeVisibility: {OneOf: visibility}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleItem, Max: -1, Kinds: items}}},
	s.KindUse:                   {Attributes: map[string]syntax.Attribute{s.AttributeAlias: {Predicate: lexical.IdentifierOrUnderscore()}}, Roles: []syntax.Role[s.Language]{{Name: s.RolePath, Min: 1, Max: 1, Kinds: []s.Kind{s.KindPathType}}}},
	s.KindTypeAlias:             {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeVisibility: {OneOf: visibility}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleGenericParameter, Max: -1, Kinds: generics}, {Name: s.RoleWhere, Max: -1, Kinds: []s.Kind{s.KindWherePredicate}}, {Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindStruct:                {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeVisibility: {OneOf: visibility}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleGenericParameter, Max: -1, Kinds: generics}, {Name: s.RoleWhere, Max: -1, Kinds: []s.Kind{s.KindWherePredicate}}, {Name: s.RoleField, Max: -1, Kinds: []s.Kind{s.KindField}}}},
	s.KindTupleStruct:           {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeVisibility: {OneOf: visibility}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleGenericParameter, Max: -1, Kinds: generics}, {Name: s.RoleWhere, Max: -1, Kinds: []s.Kind{s.KindWherePredicate}}, {Name: s.RoleField, Max: -1, Kinds: []s.Kind{s.KindTupleField}}}},
	s.KindEnum:                  {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeVisibility: {OneOf: visibility}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleGenericParameter, Max: -1, Kinds: generics}, {Name: s.RoleWhere, Max: -1, Kinds: []s.Kind{s.KindWherePredicate}}, {Name: s.RoleVariant, Max: -1, Kinds: []s.Kind{s.KindUnitVariant, s.KindTupleVariant, s.KindRecordVariant}}}},
	s.KindUnitVariant:           {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}},
	s.KindTupleVariant:          {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleField, Max: -1, Kinds: []s.Kind{s.KindTupleField}}}},
	s.KindRecordVariant:         {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleField, Max: -1, Kinds: []s.Kind{s.KindField}}}},
	s.KindField:                 {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeVisibility: {OneOf: visibility}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindTupleField:            {Attributes: map[string]syntax.Attribute{s.AttributeVisibility: {OneOf: visibility}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindTrait:                 {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}, s.AttributeVisibility: {OneOf: visibility}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleGenericParameter, Max: -1, Kinds: generics}, {Name: s.RoleBound, Max: -1, Kinds: bounds}, {Name: s.RoleWhere, Max: -1, Kinds: []s.Kind{s.KindWherePredicate}}, {Name: s.RoleMember, Max: -1, Kinds: []s.Kind{s.KindMethodSignature, s.KindAssociatedType}}}},
	s.KindTraitImpl:             {Roles: []syntax.Role[s.Language]{{Name: s.RoleGenericParameter, Max: -1, Kinds: generics}, {Name: s.RoleTrait, Min: 1, Max: 1, Kinds: []s.Kind{s.KindPathType}}, {Name: s.RoleSubject, Min: 1, Max: 1, Kinds: types}, {Name: s.RoleWhere, Max: -1, Kinds: []s.Kind{s.KindWherePredicate}}, {Name: s.RoleMember, Max: -1, Kinds: []s.Kind{s.KindTypeAlias}}}},
	s.KindAssociatedType:        {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleGenericParameter, Max: -1, Kinds: generics}, {Name: s.RoleBound, Max: -1, Kinds: typeBounds}, {Name: s.RoleWhere, Max: -1, Kinds: []s.Kind{s.KindWherePredicate}}}},
	s.KindMethodSignature:       {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleGenericParameter, Max: -1, Kinds: generics}, {Name: s.RoleWhere, Max: -1, Kinds: []s.Kind{s.KindWherePredicate}}, {Name: s.RoleReceiver, Max: 1, Kinds: []s.Kind{s.KindReceiver}}, {Name: s.RoleParameter, Max: -1, Kinds: []s.Kind{s.KindParameter}}, {Name: s.RoleResult, Max: 1, Kinds: types}}},
	s.KindReceiver:              {Attributes: map[string]syntax.Attribute{s.AttributeReceiverForm: {Required: true, OneOf: []string{"value", "shared", "mutable"}}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleLifetime, Max: 1, Kinds: []s.Kind{s.KindLifetimeReference}}}},
	s.KindParameter:             {Attributes: map[string]syntax.Attribute{s.AttributeName: {Predicate: lexical.IdentifierOrUnderscore()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindTypeParameter:         {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Identifier()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleBound, Max: -1, Kinds: typeBounds}, {Name: s.RoleDefault, Max: 1, Kinds: types}}},
	s.KindLifetimeParameter:     {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.LifetimeParameter()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleBound, Max: -1, Kinds: []s.Kind{s.KindLifetimeReference}}}},
	s.KindLifetimeReference:     {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.Lifetime()}}, Reference: true},
	s.KindWherePredicate:        {Roles: []syntax.Role[s.Language]{{Name: s.RoleSubject, Min: 1, Max: 1, Kinds: append(types, s.KindLifetimeReference)}, {Name: s.RoleBound, Min: 1, Max: -1, Kinds: typeBounds}}},
	s.KindRelaxedBound:          {Roles: []syntax.Role[s.Language]{{Name: s.RolePath, Min: 1, Max: 1, Kinds: []s.Kind{s.KindPathType}}}},
	s.KindPathType:              {Attributes: map[string]syntax.Attribute{s.AttributeAbsolute: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleSegment, Min: 1, Max: -1, Kinds: []s.Kind{s.KindPathSegment}}}},
	s.KindQualifiedPathType:     {Roles: []syntax.Role[s.Language]{{Name: s.RoleSubject, Min: 1, Max: 1, Kinds: types}, {Name: s.RoleTrait, Max: 1, Kinds: []s.Kind{s.KindPathType}}, {Name: s.RoleSegment, Min: 1, Max: -1, Kinds: []s.Kind{s.KindPathSegment}}}},
	s.KindPathSegment:           {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, Predicate: lexical.PathSegment()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleArgument, Max: -1, Kinds: append(types, s.KindLifetimeReference, s.KindAssociatedTypeBinding)}}, Reference: true},
	s.KindAssociatedTypeBinding: {Roles: []syntax.Role[s.Language]{{Name: s.RolePath, Min: 1, Max: 1, Kinds: []s.Kind{s.KindPathSegment}}, {Name: s.RoleType, Min: 1, Max: 1, Kinds: types}}},
	s.KindPrimitiveType:         {Attributes: map[string]syntax.Attribute{s.AttributeName: {Required: true, OneOf: primitive}}},
	s.KindReferenceType:         {Attributes: map[string]syntax.Attribute{s.AttributeMutable: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}, {Name: s.RoleLifetime, Max: 1, Kinds: []s.Kind{s.KindLifetimeReference}}}},
	s.KindRawPointerType:        {Attributes: map[string]syntax.Attribute{s.AttributeMutable: {Predicate: syntax.Bool()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindSliceType:             {Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindArrayType:             {Attributes: map[string]syntax.Attribute{s.AttributeLength: {Required: true, Predicate: syntax.Decimal()}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Min: 1, Max: 1, Kinds: types}}},
	s.KindTupleType:             {Roles: []syntax.Role[s.Language]{{Name: s.RoleElement, Max: -1, Kinds: types}}},
	s.KindNeverType:             {},
	s.KindFunctionType:          {Attributes: map[string]syntax.Attribute{s.AttributeUnsafe: {Predicate: syntax.Bool()}, s.AttributeABI: {OneOf: abis}}, Roles: []syntax.Role[s.Language]{{Name: s.RoleParameter, Max: -1, Kinds: []s.Kind{s.KindParameter}}, {Name: s.RoleResult, Max: 1, Kinds: types}}},
	s.KindTraitObjectType:       {Roles: []syntax.Role[s.Language]{{Name: s.RoleBound, Min: 1, Max: -1, Kinds: bounds}}},
}

// bittype:grammar rules
// bittype:scope docs/syntax.md#rust
func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, s.KindFile, rules, evaluators); err != nil {
		return err
	}
	v := tree.View[s.Language]{Tree: t}
	noDefaults := func(owner ast.NodeID) error {
		for _, p := range v.Children(owner, s.RoleGenericParameter) {
			if v.Child(p, s.RoleDefault) != 0 {
				return text.AtNode(t, p, "generic parameter defaults are unsupported in this declaration")
			}
		}
		return nil
	}
	for i, n := range t.Nodes {
		id := ast.NodeID(i + 1)
		bad := func(m string) error { return text.AtNode(t, id, m) }
		for _, child := range n.Children {
			if v.Node(child.Node).Kind != s.KindPathType || child.Role == s.RoleBound || child.Role == s.RoleTrait || n.Kind == s.KindUse || n.Kind == s.KindRelaxedBound {
				continue
			}
			segments := v.Children(child.Node, s.RoleSegment)
			if !v.Flag(child.Node, s.AttributeAbsolute) && len(segments) == 1 && len(v.Children(segments[0], s.RoleArgument)) == 0 && slices.Contains(primitive, v.Attr(segments[0], s.AttributeName)) {
				return text.AtNode(t, child.Node, "bare primitive spelling requires primitive-type in a type expression")
			}
		}
		seenType := false
		for _, p := range v.Children(id, s.RoleGenericParameter) {
			if v.Node(p).Kind == s.KindTypeParameter {
				seenType = true
			} else if seenType {
				return bad("lifetime parameters must precede type parameters")
			}
		}
		switch n.Kind {
		case s.KindTraitImpl:
			if err := noDefaults(id); err != nil {
				return err
			}
			for _, member := range v.Children(id, s.RoleMember) {
				if vis := v.Attr(member, s.AttributeVisibility); vis != "" && vis != "private" {
					return text.AtNode(t, member, "trait impl associated types inherit visibility")
				}
				if err := noDefaults(member); err != nil {
					return err
				}
			}
		case s.KindAssociatedType:
			if err := noDefaults(id); err != nil {
				return err
			}
		case s.KindWherePredicate:
			if v.Node(v.Child(id, s.RoleSubject)).Kind == s.KindLifetimeReference {
				for _, bound := range v.Children(id, s.RoleBound) {
					if v.Node(bound).Kind != s.KindLifetimeReference {
						return text.AtNode(t, bound, "lifetime predicates require lifetime bounds")
					}
				}
			}
		case s.KindPathSegment:
			seenType, seenBinding := false, false
			for _, arg := range v.Children(id, s.RoleArgument) {
				switch v.Node(arg).Kind {
				case s.KindLifetimeReference:
					if seenType || seenBinding {
						return text.AtNode(t, arg, "lifetime arguments must precede type arguments and associated-type bindings")
					}
				case s.KindAssociatedTypeBinding:
					seenBinding = true
				default:
					if seenBinding {
						return text.AtNode(t, arg, "type arguments must precede associated-type bindings")
					}
					seenType = true
				}
			}
		case s.KindAssociatedTypeBinding:
			if !native.Identifier(v.Attr(v.Child(id, s.RolePath), s.AttributeName)) {
				return bad("associated-type binding requires an unqualified associated name")
			}
		case s.KindUse:
			for _, seg := range v.Children(v.Child(id, s.RolePath), s.RoleSegment) {
				if len(v.Children(seg, s.RoleArgument)) > 0 {
					return bad("use paths cannot have generic arguments")
				}
			}
		case s.KindReceiver:
			if v.Attr(id, s.AttributeReceiverForm) == "value" && v.Child(id, s.RoleLifetime) != 0 {
				return bad("value receiver has no lifetime annotation")
			}
		case s.KindMethodSignature:
			for _, p := range v.Children(id, s.RoleParameter) {
				if !v.Has(p, s.AttributeName) {
					return text.AtNode(t, p, "method parameter requires a name")
				}
			}
			if err := noDefaults(id); err != nil {
				return err
			}
		case s.KindTupleVariant, s.KindRecordVariant:
			for _, field := range v.Children(id, s.RoleField) {
				if vis := v.Attr(field, s.AttributeVisibility); vis != "" && vis != "private" {
					return bad("variant fields inherit visibility")
				}
			}
		case s.KindTraitObjectType:
			hasTrait := false
			for _, b := range v.Children(id, s.RoleBound) {
				hasTrait = hasTrait || v.Node(b).Kind == s.KindPathType
			}
			if !hasTrait {
				return bad("trait object requires a trait path")
			}
		}
	}
	return nil
}

type Provider struct{ s.Language }

func (Provider) Check(view ast.View[s.Language]) error { return Check(view) }

var _ checkcontract.Provider[s.Language] = Provider{}
