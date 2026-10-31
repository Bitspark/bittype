package rustcheck

import (
	"github.com/Bitspark/bittype/ast/go"
	g "github.com/Bitspark/bittype/internal/grammar"
	native "github.com/Bitspark/bittype/internal/rust"
	"github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/text/go"
)

func identifier(v string) bool { return native.Identifier(v) }

var types = []string{"path-type", "primitive-type", "reference-type", "raw-pointer-type", "slice-type", "array-type", "tuple-type", "never-type", "function-type", "trait-object-type"}
var bounds = []string{"path-type", "lifetime-reference"}
var generics = []string{"type-parameter", "lifetime-parameter"}
var items = []string{"module", "use", "type-alias", "struct", "tuple-struct", "enum", "trait"}
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
	"file":               {Roles: []g.Role{g.Many("item", 0, items...)}},
	"module":             {Attributes: named(), Roles: []g.Role{g.Many("item", 0, items...)}},
	"use":                {Attributes: g.Attributes("alias", g.Attr(false, identifier)), Roles: []g.Role{g.One("path", "path-type")}},
	"type-alias":         {Attributes: named(), Roles: with(prefix(), g.One("type", types...))},
	"struct":             {Attributes: named(), Roles: with(prefix(), g.Many("field", 0, "field"))},
	"tuple-struct":       {Attributes: named(), Roles: with(prefix(), g.Many("field", 0, "tuple-field"))},
	"enum":               {Attributes: named(), Roles: with(prefix(), g.Many("variant", 0, "unit-variant", "tuple-variant", "record-variant"))},
	"unit-variant":       {Attributes: g.Attributes("name", g.Attr(true, identifier))},
	"tuple-variant":      {Attributes: g.Attributes("name", g.Attr(true, identifier)), Roles: []g.Role{g.Many("field", 0, "tuple-field")}},
	"record-variant":     {Attributes: g.Attributes("name", g.Attr(true, identifier)), Roles: []g.Role{g.Many("field", 0, "field")}},
	"field":              {Attributes: named(), Roles: []g.Role{g.One("type", types...)}},
	"tuple-field":        {Attributes: g.Attributes("visibility", g.Attr(false, visibility)), Roles: []g.Role{g.One("type", types...)}},
	"trait":              {Attributes: named(), Roles: []g.Role{g.Many("generic-parameter", 0, generics...), g.Many("bound", 0, bounds...), g.Many("where", 0, "where-predicate"), g.Many("member", 0, "method-signature", "associated-type")}},
	"associated-type":    {Attributes: g.Attributes("name", g.Attr(true, identifier)), Roles: []g.Role{g.Many("bound", 0, bounds...)}},
	"method-signature":   {Attributes: g.Attributes("name", g.Attr(true, identifier)), Roles: with(prefix(), g.Opt("receiver", "receiver"), g.Many("parameter", 0, "parameter"), g.Opt("result", types...))},
	"receiver":           {Attributes: g.Attributes("receiver-form", g.Attr(true, g.Set("value", "shared", "mutable"))), Roles: []g.Role{g.Opt("lifetime", "lifetime-reference")}},
	"parameter":          {Attributes: g.Attributes("name", g.Attr(false, identifier)), Roles: []g.Role{g.One("type", types...)}},
	"type-parameter":     {Attributes: g.Attributes("name", g.Attr(true, identifier)), Roles: []g.Role{g.Many("bound", 0, bounds...), g.Opt("default", types...)}},
	"lifetime-parameter": {Attributes: g.Attributes("name", g.Attr(true, g.Identifier)), Roles: []g.Role{g.Many("bound", 0, "lifetime-reference")}},
	"lifetime-reference": {Attributes: g.Attributes("name", g.Attr(true, g.Identifier)), Reference: true},
	"where-predicate":    {Roles: []g.Role{g.One("subject", append(types, "lifetime-reference")...), g.Many("bound", 1, bounds...)}},
	"path-type":          {Attributes: g.Attributes("absolute", g.Attr(false, g.Bool)), Roles: []g.Role{g.Many("segment", 1, "path-segment")}},
	"path-segment":       {Attributes: g.Attributes("name", g.Attr(true, native.PathSegment)), Roles: []g.Role{g.Many("argument", 0, append(types, "lifetime-reference")...)}, Reference: true},
	"primitive-type":     {Attributes: g.Attributes("name", g.Attr(true, primitive))},
	"reference-type":     {Attributes: g.Attributes("mutable", g.Attr(false, g.Bool)), Roles: []g.Role{g.One("element", types...), g.Opt("lifetime", "lifetime-reference")}},
	"raw-pointer-type":   {Attributes: g.Attributes("mutable", g.Attr(false, g.Bool)), Roles: []g.Role{g.One("element", types...)}},
	"slice-type":         {Roles: []g.Role{g.One("element", types...)}},
	"array-type":         {Attributes: g.Attributes("length", g.Attr(true, g.Decimal)), Roles: []g.Role{g.One("element", types...)}},
	"tuple-type":         {Roles: []g.Role{g.Many("element", 0, types...)}},
	"never-type":         {},
	"function-type":      {Attributes: g.Attributes("unsafe", g.Attr(false, g.Bool), "abi", g.Attr(false, native.ABI)), Roles: []g.Role{g.Many("parameter", 0, "parameter"), g.Opt("result", types...)}},
	"trait-object-type":  {Roles: []g.Role{g.Many("bound", 1, bounds...)}},
}

func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, "file", rules); err != nil {
		return err
	}
	v := tree.View[s.Language]{Tree: t}
	for i, n := range t.Nodes {
		id := ast.NodeID(i + 1)
		bad := func(m string) error { return text.AtNode(t, id, m) }
		seenType := false
		for _, p := range v.Children(id, "generic-parameter") {
			if v.Kind(p) == "type-parameter" {
				seenType = true
			} else if seenType {
				return bad("lifetime parameters must precede type parameters")
			}
		}
		switch n.Kind {
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
			for _, p := range v.Children(id, "generic-parameter") {
				if v.Child(p, "default") != 0 {
					return text.AtNode(t, p, "method generic defaults are unsupported")
				}
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

var _ text.CheckerProvider[s.Language] = Provider{}
