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

var types = []string{"path-type", "qualified-path-type", "primitive-type", "reference-type", "raw-pointer-type", "slice-type", "array-type", "tuple-type", "never-type", "function-type", "trait-object-type"}
var bounds = []string{"path-type", "lifetime-reference"}
var typeBounds = []string{"path-type", "lifetime-reference", "relaxed-bound"}
var generics = []string{"type-parameter", "lifetime-parameter"}
var items = []string{"module", "use", "type-alias", "struct", "tuple-struct", "enum", "trait", "trait-impl"}
var visibility = g.Set("private", "public", "crate")
var primitive = g.Set("bool", "char", "str", "f32", "f64", "i8", "i16", "i32", "i64", "i128", "isize", "u8", "u16", "u32", "u64", "u128", "usize")

func named() map[string]g.Attribute {
	return g.Attributes("name", g.Attr(true, native.Declaration), "visibility", g.Attr(false, visibility))
}
func prefix() []g.Role {
	return []g.Role{g.Many("generic-parameter", 0, generics...), g.Many("where", 0, "where-predicate")}
}
func with(rs []g.Role, more ...g.Role) []g.Role { return append(rs, more...) }

var rules = g.Rules{
	"file":                    {Roles: []g.Role{g.Many("item", 0, items...)}},
	"module":                  {Attributes: named(), Roles: []g.Role{g.Many("item", 0, items...)}},
	"use":                     {Attributes: g.Attributes("alias", g.Attr(false, identifier)), Roles: []g.Role{g.One("path", "path-type")}},
	"type-alias":              {Attributes: named(), Roles: with(prefix(), g.One("type", types...))},
	"struct":                  {Attributes: named(), Roles: with(prefix(), g.Many("field", 0, "field"))},
	"tuple-struct":            {Attributes: named(), Roles: with(prefix(), g.Many("field", 0, "tuple-field"))},
	"enum":                    {Attributes: named(), Roles: with(prefix(), g.Many("variant", 0, "unit-variant", "tuple-variant", "record-variant"))},
	"unit-variant":            {Attributes: g.Attributes("name", g.Attr(true, native.Declaration))},
	"tuple-variant":           {Attributes: g.Attributes("name", g.Attr(true, native.Declaration)), Roles: []g.Role{g.Many("field", 0, "tuple-field")}},
	"record-variant":          {Attributes: g.Attributes("name", g.Attr(true, native.Declaration)), Roles: []g.Role{g.Many("field", 0, "field")}},
	"field":                   {Attributes: named(), Roles: []g.Role{g.One("type", types...)}},
	"tuple-field":             {Attributes: g.Attributes("visibility", g.Attr(false, visibility)), Roles: []g.Role{g.One("type", types...)}},
	"trait":                   {Attributes: named(), Roles: []g.Role{g.Many("generic-parameter", 0, generics...), g.Many("bound", 0, bounds...), g.Many("where", 0, "where-predicate"), g.Many("member", 0, "method-signature", "associated-type")}},
	"trait-impl":              {Roles: []g.Role{g.Many("generic-parameter", 0, generics...), g.One("trait", "path-type"), g.One("subject", types...), g.Many("where", 0, "where-predicate"), g.Many("member", 0, "type-alias")}},
	"associated-type":         {Attributes: g.Attributes("name", g.Attr(true, native.Declaration)), Roles: []g.Role{g.Many("generic-parameter", 0, generics...), g.Many("bound", 0, typeBounds...), g.Many("where", 0, "where-predicate")}},
	"method-signature":        {Attributes: g.Attributes("name", g.Attr(true, native.Declaration)), Roles: with(prefix(), g.Opt("receiver", "receiver"), g.Many("parameter", 0, "parameter"), g.Opt("result", types...))},
	"receiver":                {Attributes: g.Attributes("receiver-form", g.Attr(true, g.Set("value", "shared", "mutable"))), Roles: []g.Role{g.Opt("lifetime", "lifetime-reference")}},
	"parameter":               {Attributes: g.Attributes("name", g.Attr(false, identifier)), Roles: []g.Role{g.One("type", types...)}},
	"type-parameter":          {Attributes: g.Attributes("name", g.Attr(true, native.Declaration)), Roles: []g.Role{g.Many("bound", 0, typeBounds...), g.Opt("default", types...)}},
	"lifetime-parameter":      {Attributes: g.Attributes("name", g.Attr(true, native.LifetimeParameter)), Roles: []g.Role{g.Many("bound", 0, "lifetime-reference")}},
	"lifetime-reference":      {Attributes: g.Attributes("name", g.Attr(true, native.LifetimeReference)), Reference: true},
	"where-predicate":         {Roles: []g.Role{g.One("subject", append(types, "lifetime-reference")...), g.Many("bound", 1, typeBounds...)}},
	"relaxed-bound":           {Roles: []g.Role{g.One("path", "path-type")}},
	"path-type":               {Attributes: g.Attributes("absolute", g.Attr(false, g.Bool)), Roles: []g.Role{g.Many("segment", 1, "path-segment")}},
	"qualified-path-type":     {Roles: []g.Role{g.One("subject", types...), g.Opt("trait", "path-type"), g.Many("segment", 1, "path-segment")}},
	"path-segment":            {Attributes: g.Attributes("name", g.Attr(true, native.PathSegment)), Roles: []g.Role{g.Many("argument", 0, append(types, "lifetime-reference", "associated-type-binding")...)}, Reference: true},
	"associated-type-binding": {Roles: []g.Role{g.One("path", "path-segment"), g.One("type", types...)}},
	"primitive-type":          {Attributes: g.Attributes("name", g.Attr(true, primitive))},
	"reference-type":          {Attributes: g.Attributes("mutable", g.Attr(false, g.Bool)), Roles: []g.Role{g.One("element", types...), g.Opt("lifetime", "lifetime-reference")}},
	"raw-pointer-type":        {Attributes: g.Attributes("mutable", g.Attr(false, g.Bool)), Roles: []g.Role{g.One("element", types...)}},
	"slice-type":              {Roles: []g.Role{g.One("element", types...)}},
	"array-type":              {Attributes: g.Attributes("length", g.Attr(true, g.Decimal)), Roles: []g.Role{g.One("element", types...)}},
	"tuple-type":              {Roles: []g.Role{g.Many("element", 0, types...)}},
	"never-type":              {},
	"function-type":           {Attributes: g.Attributes("unsafe", g.Attr(false, g.Bool), "abi", g.Attr(false, native.ABI)), Roles: []g.Role{g.Many("parameter", 0, "parameter"), g.Opt("result", types...)}},
	"trait-object-type":       {Roles: []g.Role{g.Many("bound", 1, bounds...)}},
}

// bittype:grammar rules
// bittype:scope docs/syntax.md#rust
func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, "file", rules); err != nil {
		return err
	}
	v := tree.View[s.Language]{Tree: t}
	noDefaults := func(owner ast.NodeID) error {
		for _, p := range v.Children(owner, "generic-parameter") {
			if v.Child(p, "default") != 0 {
				return text.AtNode(t, p, "generic parameter defaults are unsupported in this declaration")
			}
		}
		return nil
	}
	for i, n := range t.Nodes {
		id := ast.NodeID(i + 1)
		bad := func(m string) error { return text.AtNode(t, id, m) }
		for _, child := range n.Children {
			if v.Kind(child.Node) != "path-type" || child.Role == "bound" || child.Role == "trait" || n.Kind == s.KindUse || n.Kind == s.KindRelaxedBound {
				continue
			}
			segments := v.Children(child.Node, "segment")
			if !v.Flag(child.Node, "absolute") && len(segments) == 1 && len(v.Children(segments[0], "argument")) == 0 && primitive(v.Attr(segments[0], "name")) {
				return text.AtNode(t, child.Node, "bare primitive spelling requires primitive-type in a type expression")
			}
		}
		seenType := false
		for _, p := range v.Children(id, "generic-parameter") {
			if v.Kind(p) == "type-parameter" {
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
			for _, member := range v.Children(id, "member") {
				if vis := v.Attr(member, "visibility"); vis != "" && vis != "private" {
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
			if v.Kind(v.Child(id, "subject")) == "lifetime-reference" {
				for _, bound := range v.Children(id, "bound") {
					if v.Kind(bound) != "lifetime-reference" {
						return text.AtNode(t, bound, "lifetime predicates require lifetime bounds")
					}
				}
			}
		case s.KindPathSegment:
			seenType, seenBinding := false, false
			for _, arg := range v.Children(id, "argument") {
				switch v.Kind(arg) {
				case "lifetime-reference":
					if seenType || seenBinding {
						return text.AtNode(t, arg, "lifetime arguments must precede type arguments and associated-type bindings")
					}
				case "associated-type-binding":
					seenBinding = true
				default:
					if seenBinding {
						return text.AtNode(t, arg, "type arguments must precede associated-type bindings")
					}
					seenType = true
				}
			}
		case s.KindAssociatedTypeBinding:
			if !native.Declaration(v.Attr(v.Child(id, "path"), "name")) {
				return bad("associated-type binding requires an unqualified associated name")
			}
		case s.KindUse:
			for _, seg := range v.Children(v.Child(id, "path"), "segment") {
				if len(v.Children(seg, "argument")) > 0 {
					return bad("use paths cannot have generic arguments")
				}
			}
		case s.KindReceiver:
			if v.Attr(id, "receiver-form") == "value" && v.Child(id, "lifetime") != 0 {
				return bad("value receiver has no lifetime annotation")
			}
		case s.KindMethodSignature:
			for _, p := range v.Children(id, "parameter") {
				if !v.Has(p, "name") {
					return text.AtNode(t, p, "method parameter requires a name")
				}
			}
			if err := noDefaults(id); err != nil {
				return err
			}
		case s.KindTupleVariant, s.KindRecordVariant:
			for _, field := range v.Children(id, "field") {
				if vis := v.Attr(field, "visibility"); vis != "" && vis != "private" {
					return bad("variant fields inherit visibility")
				}
			}
		case s.KindTraitObjectType:
			hasTrait := false
			for _, b := range v.Children(id, "bound") {
				hasTrait = hasTrait || v.Kind(b) == "path-type"
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
