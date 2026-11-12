package check

import (
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	g "github.com/Bitspark/bittype/internal/grammar"
	native "github.com/Bitspark/bittype/internal/rust"
	"github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func identifier(v string) bool { return native.Identifier(v) }

var types = []s.Kind{s.KindPathType, s.KindQualifiedPathType, s.KindPrimitiveType, s.KindReferenceType, s.KindRawPointerType, s.KindSliceType, s.KindArrayType, s.KindTupleType, s.KindNeverType, s.KindFunctionType, s.KindTraitObjectType}
var bounds = []s.Kind{s.KindPathType, s.KindLifetimeReference}
var typeBounds = []s.Kind{s.KindPathType, s.KindLifetimeReference, s.KindRelaxedBound}
var generics = []s.Kind{s.KindTypeParameter, s.KindLifetimeParameter}
var items = []s.Kind{s.KindModule, s.KindUse, s.KindTypeAlias, s.KindStruct, s.KindTupleStruct, s.KindEnum, s.KindTrait, s.KindTraitImpl}
var visibility = g.Set("private", "public", "crate")
var primitive = g.Set("bool", "char", "str", "f32", "f64", "i8", "i16", "i32", "i64", "i128", "isize", "u8", "u16", "u32", "u64", "u128", "usize")

func named() map[string]g.Attribute {
	return g.Attributes(s.AttributeName, g.Attr(true, native.Declaration), s.AttributeVisibility, g.Attr(false, visibility))
}
func prefix() []g.Role[s.Language] {
	return []g.Role[s.Language]{g.Many(s.RoleGenericParameter, 0, generics...), g.Many(s.RoleWhere, 0, s.KindWherePredicate)}
}
func with(rs []g.Role[s.Language], more ...g.Role[s.Language]) []g.Role[s.Language] {
	return append(rs, more...)
}

var rules = g.Rules[s.Language]{
	s.KindFile:                  {Roles: []g.Role[s.Language]{g.Many(s.RoleItem, 0, items...)}},
	s.KindModule:                {Attributes: named(), Roles: []g.Role[s.Language]{g.Many(s.RoleItem, 0, items...)}},
	s.KindUse:                   {Attributes: g.Attributes(s.AttributeAlias, g.Attr(false, identifier)), Roles: []g.Role[s.Language]{g.One(s.RolePath, s.KindPathType)}},
	s.KindTypeAlias:             {Attributes: named(), Roles: with(prefix(), g.One(s.RoleType, types...))},
	s.KindStruct:                {Attributes: named(), Roles: with(prefix(), g.Many(s.RoleField, 0, s.KindField))},
	s.KindTupleStruct:           {Attributes: named(), Roles: with(prefix(), g.Many(s.RoleField, 0, s.KindTupleField))},
	s.KindEnum:                  {Attributes: named(), Roles: with(prefix(), g.Many(s.RoleVariant, 0, s.KindUnitVariant, s.KindTupleVariant, s.KindRecordVariant))},
	s.KindUnitVariant:           {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.Declaration))},
	s.KindTupleVariant:          {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.Declaration)), Roles: []g.Role[s.Language]{g.Many(s.RoleField, 0, s.KindTupleField)}},
	s.KindRecordVariant:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.Declaration)), Roles: []g.Role[s.Language]{g.Many(s.RoleField, 0, s.KindField)}},
	s.KindField:                 {Attributes: named(), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindTupleField:            {Attributes: g.Attributes(s.AttributeVisibility, g.Attr(false, visibility)), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindTrait:                 {Attributes: named(), Roles: []g.Role[s.Language]{g.Many(s.RoleGenericParameter, 0, generics...), g.Many(s.RoleBound, 0, bounds...), g.Many(s.RoleWhere, 0, s.KindWherePredicate), g.Many(s.RoleMember, 0, s.KindMethodSignature, s.KindAssociatedType)}},
	s.KindTraitImpl:             {Roles: []g.Role[s.Language]{g.Many(s.RoleGenericParameter, 0, generics...), g.One(s.RoleTrait, s.KindPathType), g.One(s.RoleSubject, types...), g.Many(s.RoleWhere, 0, s.KindWherePredicate), g.Many(s.RoleMember, 0, s.KindTypeAlias)}},
	s.KindAssociatedType:        {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.Declaration)), Roles: []g.Role[s.Language]{g.Many(s.RoleGenericParameter, 0, generics...), g.Many(s.RoleBound, 0, typeBounds...), g.Many(s.RoleWhere, 0, s.KindWherePredicate)}},
	s.KindMethodSignature:       {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.Declaration)), Roles: with(prefix(), g.Opt(s.RoleReceiver, s.KindReceiver), g.Many(s.RoleParameter, 0, s.KindParameter), g.Opt(s.RoleResult, types...))},
	s.KindReceiver:              {Attributes: g.Attributes(s.AttributeReceiverForm, g.Attr(true, g.Set("value", "shared", "mutable"))), Roles: []g.Role[s.Language]{g.Opt(s.RoleLifetime, s.KindLifetimeReference)}},
	s.KindParameter:             {Attributes: g.Attributes(s.AttributeName, g.Attr(false, identifier)), Roles: []g.Role[s.Language]{g.One(s.RoleType, types...)}},
	s.KindTypeParameter:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.Declaration)), Roles: []g.Role[s.Language]{g.Many(s.RoleBound, 0, typeBounds...), g.Opt(s.RoleDefault, types...)}},
	s.KindLifetimeParameter:     {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.LifetimeParameter)), Roles: []g.Role[s.Language]{g.Many(s.RoleBound, 0, s.KindLifetimeReference)}},
	s.KindLifetimeReference:     {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.LifetimeReference)), Reference: true},
	s.KindWherePredicate:        {Roles: []g.Role[s.Language]{g.One(s.RoleSubject, append(types, s.KindLifetimeReference)...), g.Many(s.RoleBound, 1, typeBounds...)}},
	s.KindRelaxedBound:          {Roles: []g.Role[s.Language]{g.One(s.RolePath, s.KindPathType)}},
	s.KindPathType:              {Attributes: g.Attributes(s.AttributeAbsolute, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.Many(s.RoleSegment, 1, s.KindPathSegment)}},
	s.KindQualifiedPathType:     {Roles: []g.Role[s.Language]{g.One(s.RoleSubject, types...), g.Opt(s.RoleTrait, s.KindPathType), g.Many(s.RoleSegment, 1, s.KindPathSegment)}},
	s.KindPathSegment:           {Attributes: g.Attributes(s.AttributeName, g.Attr(true, native.PathSegment)), Roles: []g.Role[s.Language]{g.Many(s.RoleArgument, 0, append(types, s.KindLifetimeReference, s.KindAssociatedTypeBinding)...)}, Reference: true},
	s.KindAssociatedTypeBinding: {Roles: []g.Role[s.Language]{g.One(s.RolePath, s.KindPathSegment), g.One(s.RoleType, types...)}},
	s.KindPrimitiveType:         {Attributes: g.Attributes(s.AttributeName, g.Attr(true, primitive))},
	s.KindReferenceType:         {Attributes: g.Attributes(s.AttributeMutable, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...), g.Opt(s.RoleLifetime, s.KindLifetimeReference)}},
	s.KindRawPointerType:        {Attributes: g.Attributes(s.AttributeMutable, g.Attr(false, g.Bool)), Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindSliceType:             {Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindArrayType:             {Attributes: g.Attributes(s.AttributeLength, g.Attr(true, g.Decimal)), Roles: []g.Role[s.Language]{g.One(s.RoleElement, types...)}},
	s.KindTupleType:             {Roles: []g.Role[s.Language]{g.Many(s.RoleElement, 0, types...)}},
	s.KindNeverType:             {},
	s.KindFunctionType:          {Attributes: g.Attributes(s.AttributeUnsafe, g.Attr(false, g.Bool), s.AttributeABI, g.Attr(false, native.ABI)), Roles: []g.Role[s.Language]{g.Many(s.RoleParameter, 0, s.KindParameter), g.Opt(s.RoleResult, types...)}},
	s.KindTraitObjectType:       {Roles: []g.Role[s.Language]{g.Many(s.RoleBound, 1, bounds...)}},
}

// bittype:grammar rules
// bittype:scope docs/syntax.md#rust
func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, s.KindFile, rules); err != nil {
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
			if !v.Flag(child.Node, s.AttributeAbsolute) && len(segments) == 1 && len(v.Children(segments[0], s.RoleArgument)) == 0 && primitive(v.Attr(segments[0], s.AttributeName)) {
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
			if !native.Declaration(v.Attr(v.Child(id, s.RolePath), s.AttributeName)) {
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
