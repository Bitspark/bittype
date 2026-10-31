package cppcheck

import (
	"github.com/Bitspark/bittype/contracts/text/go"
	native "github.com/Bitspark/bittype/internal/cpp"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"strings"
)

var types = []string{"named-type", "builtin-type", "pointer-type", "lvalue-reference-type", "rvalue-reference-type", "cv-qualified-type", "array-type", "function-type", "member-pointer-type"}
var declarations = []string{"namespace", "record", "type-alias", "enum"}
var members = []string{"field", "method", "destructor", "type-alias", "record", "enum"}
var access = g.Set("public", "protected", "private")

func named() map[string]g.Attribute {
	return g.Attributes("name", g.Attr(true, native.Identifier), "access", g.Attr(false, access))
}
func extra(base map[string]g.Attribute, pairs ...any) map[string]g.Attribute {
	for k, v := range g.Attributes(pairs...) {
		base[k] = v
	}
	return base
}

var rules = g.Rules{
	"translation-unit":      {Roles: []g.Role{g.Many("include", 0, "include"), g.Many("declaration", 0, declarations...)}},
	"include":               {Attributes: g.Attributes("path", g.Attr(true, func(v string) bool { return g.Nonempty(v) && !strings.ContainsAny(v, "\r\n") }), "system", g.Attr(false, g.Bool))},
	"namespace":             {Attributes: g.Attributes("name", g.Attr(true, native.Identifier)), Roles: []g.Role{g.Many("declaration", 0, declarations...)}},
	"record":                {Attributes: extra(named(), "record-key", g.Attr(true, g.Set("struct", "class")), "forward", g.Attr(false, g.Bool)), Roles: []g.Role{g.Many("template-parameter", 0, "type-parameter"), g.Many("base", 0, "base"), g.Many("member", 0, members...)}},
	"type-alias":            {Attributes: named(), Roles: []g.Role{g.Many("template-parameter", 0, "type-parameter"), g.One("type", types...)}},
	"enum":                  {Attributes: extra(named(), "scoped", g.Attr(false, g.Bool), "forward", g.Attr(false, g.Bool)), Roles: []g.Role{g.Opt("underlying", types...), g.Many("enumerator", 0, "enumerator")}},
	"enumerator":            {Attributes: g.Attributes("name", g.Attr(true, native.Identifier), "value", g.Attr(false, g.Signed))},
	"base":                  {Attributes: g.Attributes("access", g.Attr(false, access), "virtual", g.Attr(false, g.Bool)), Roles: []g.Role{g.One("type", "named-type")}},
	"field":                 {Attributes: named(), Roles: []g.Role{g.One("type", types...)}},
	"method":                {Attributes: extra(named(), "const", g.Attr(false, g.Bool), "volatile", g.Attr(false, g.Bool), "virtual", g.Attr(false, g.Bool), "pure", g.Attr(false, g.Bool), "noexcept", g.Attr(false, g.Bool), "ref-qualifier", g.Attr(false, g.Set("none", "lvalue", "rvalue")), "variadic", g.Attr(false, g.Bool)), Roles: []g.Role{g.Many("parameter", 0, "parameter"), g.One("result", types...)}},
	"destructor":            {Attributes: g.Attributes("access", g.Attr(false, access), "virtual", g.Attr(false, g.Bool), "pure", g.Attr(false, g.Bool), "noexcept", g.Attr(false, g.Bool))},
	"parameter":             {Attributes: g.Attributes("name", g.Attr(false, native.Identifier)), Roles: []g.Role{g.One("type", types...)}},
	"type-parameter":        {Attributes: g.Attributes("name", g.Attr(true, native.Identifier), "parameter-key", g.Attr(false, g.Set("typename", "class"))), Roles: []g.Role{g.Opt("default", types...)}},
	"named-type":            {Attributes: g.Attributes("global", g.Attr(false, g.Bool), "typename", g.Attr(false, g.Bool)), Roles: []g.Role{g.Many("segment", 1, "name-segment")}},
	"name-segment":          {Attributes: g.Attributes("name", g.Attr(true, native.Identifier), "template", g.Attr(false, g.Bool)), Roles: []g.Role{g.Many("argument", 0, types...)}, Reference: true},
	"builtin-type":          {Attributes: g.Attributes("name", g.Attr(true, native.Builtin))},
	"pointer-type":          {Roles: []g.Role{g.One("element", types...)}},
	"lvalue-reference-type": {Roles: []g.Role{g.One("element", types...)}},
	"rvalue-reference-type": {Roles: []g.Role{g.One("element", types...)}},
	"cv-qualified-type":     {Attributes: g.Attributes("const", g.Attr(false, g.Bool), "volatile", g.Attr(false, g.Bool)), Roles: []g.Role{g.One("element", types...)}},
	"array-type":            {Attributes: g.Attributes("length", g.Attr(true, g.Positive)), Roles: []g.Role{g.One("element", types...)}},
	"function-type":         {Attributes: g.Attributes("noexcept", g.Attr(false, g.Bool), "variadic", g.Attr(false, g.Bool)), Roles: []g.Role{g.Many("parameter", 0, "parameter"), g.One("result", types...)}},
	"member-pointer-type":   {Roles: []g.Role{g.One("owner", "named-type"), g.One("element", types...)}},
}

func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, "translation-unit", rules); err != nil {
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
		if v.Has(id, "access") && n.Kind != s.KindBase {
			parent := parents[id]
			if parent == 0 || v.Kind(parent) != "record" {
				return bad("access belongs to record members")
			}
		}
		if v.Flag(id, "pure") && !v.Flag(id, "virtual") {
			return bad("pure declaration must be virtual")
		}
		switch n.Kind {
		case s.KindInclude:
			path := v.Attr(id, "path")
			if v.Flag(id, "system") {
				if strings.ContainsRune(path, '>') {
					return bad("invalid system header name")
				}
			} else if strings.ContainsRune(path, '"') {
				return bad("invalid quoted header name")
			}
		case s.KindRecord:
			if v.Flag(id, "forward") && (len(v.Children(id, "base")) > 0 || len(v.Children(id, "member")) > 0) {
				return bad("forward record cannot have bases or members")
			}
		case s.KindEnum:
			if v.Flag(id, "forward") && (len(v.Children(id, "enumerator")) > 0 || (!v.Flag(id, "scoped") && v.Child(id, "underlying") == 0)) {
				return bad("invalid forward enum")
			}
		case s.KindCVQualifiedType:
			if !v.Flag(id, "const") && !v.Flag(id, "volatile") {
				return bad("qualification requires const or volatile")
			}
			e := v.Kind(v.Child(id, "element"))
			if e != "named-type" && e != "builtin-type" && e != "pointer-type" && e != "member-pointer-type" {
				return bad("qualifiers require a directly representable object or pointer type")
			}
		case s.KindFunctionType, s.KindMethod:
			r := v.Kind(v.Child(id, "result"))
			if r == "array-type" || r == "function-type" {
				return bad("function cannot return an array or function type")
			}
		case s.KindArrayType:
			e := v.Kind(v.Child(id, "element"))
			if e == "function-type" || e == "lvalue-reference-type" || e == "rvalue-reference-type" {
				return bad("array element cannot be a function or reference")
			}
		case s.KindPointerType, s.KindMemberPointerType, s.KindLValueReferenceType, s.KindRValueReferenceType:
			e := v.Kind(v.Child(id, "element"))
			if e == "lvalue-reference-type" || e == "rvalue-reference-type" {
				return bad("nested reference declarator is not representable")
			}
		}
	}
	return nil
}

type Provider struct{ s.Language }

func (Provider) Check(view ast.View[s.Language]) error { return Check(view) }

var _ text.CheckerProvider[s.Language] = Provider{}
