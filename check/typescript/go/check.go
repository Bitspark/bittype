package typescriptcheck

import (
	"github.com/Bitspark/bittype/trees/ast/go"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/internal/tree"
	lit "github.com/Bitspark/bittype/internal/typescript"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/contracts/text/go"
)

func identifier(v string) bool { return lit.Identifier(v) }

var types = []string{"type-reference", "qualified-type", "type-application", "keyword-type", "literal-type", "array-type", "tuple-type", "object-type", "union-type", "intersection-type", "function-type", "type-operator", "indexed-access-type"}
var members = []string{"property", "method", "call-signature", "index-signature"}
var keywords = g.Set("any", "unknown", "never", "void", "undefined", "boolean", "number", "bigint", "string", "symbol", "object", "this")
var rules = g.Rules{
	"source-file":         {Roles: []g.Role{g.Many("import", 0, "import"), g.Many("declaration", 0, "type-alias", "interface")}},
	"import":              {Attributes: g.Attributes("module", g.Attr(true, g.Nonempty)), Roles: []g.Role{g.Many("binding", 1, "import-binding", "namespace-import")}},
	"import-binding":      {Attributes: g.Attributes("name", g.Attr(true, identifier), "alias", g.Attr(false, identifier))},
	"namespace-import":    {Attributes: g.Attributes("name", g.Attr(true, identifier))},
	"type-alias":          {Attributes: g.Attributes("name", g.Attr(true, lit.TypeName), "exported", g.Attr(false, g.Bool)), Roles: []g.Role{g.Many("type-parameter", 0, "type-parameter"), g.One("type", types...)}},
	"interface":           {Attributes: g.Attributes("name", g.Attr(true, lit.TypeName), "exported", g.Attr(false, g.Bool)), Roles: []g.Role{g.Many("type-parameter", 0, "type-parameter"), g.Many("base", 0, "type-reference", "qualified-type", "type-application"), g.Many("member", 0, members...)}},
	"type-parameter":      {Attributes: g.Attributes("name", g.Attr(true, lit.TypeName)), Roles: []g.Role{g.Opt("constraint", types...), g.Opt("default", types...)}},
	"type-reference":      {Attributes: g.Attributes("name", g.Attr(true, lit.TypeName)), Reference: true},
	"qualified-type":      {Attributes: g.Attributes("name", g.Attr(true, identifier)), Roles: []g.Role{g.One("qualifier", "type-reference", "qualified-type")}, Reference: true},
	"type-application":    {Roles: []g.Role{g.One("constructor", "type-reference", "qualified-type"), g.Many("argument", 1, types...)}},
	"keyword-type":        {Attributes: g.Attributes("name", g.Attr(true, keywords))},
	"literal-type":        {Attributes: g.Attributes("literal-kind", g.Attr(true, g.Set("string", "number", "bigint", "boolean", "null")), "text", g.Attr(true, g.Nonempty))},
	"array-type":          {Roles: []g.Role{g.One("element", types...)}},
	"tuple-type":          {Roles: []g.Role{g.Many("element", 0, "tuple-element")}},
	"tuple-element":       {Attributes: g.Attributes("name", g.Attr(false, identifier), "optional", g.Attr(false, g.Bool), "rest", g.Attr(false, g.Bool)), Roles: []g.Role{g.One("type", types...)}},
	"object-type":         {Roles: []g.Role{g.Many("member", 0, members...)}},
	"union-type":          {Roles: []g.Role{g.Many("type", 2, types...)}},
	"intersection-type":   {Roles: []g.Role{g.Many("type", 2, types...)}},
	"function-type":       {Roles: []g.Role{g.Many("type-parameter", 0, "type-parameter"), g.Many("parameter", 0, "parameter"), g.One("result", types...)}},
	"parameter":           {Attributes: g.Attributes("name", g.Attr(true, identifier), "optional", g.Attr(false, g.Bool), "rest", g.Attr(false, g.Bool)), Roles: []g.Role{g.One("type", types...)}},
	"property":            {Attributes: g.Attributes("name", g.Attr(true, g.String), "optional", g.Attr(false, g.Bool), "readonly", g.Attr(false, g.Bool)), Roles: []g.Role{g.One("type", types...)}},
	"method":              {Attributes: g.Attributes("name", g.Attr(true, g.String), "optional", g.Attr(false, g.Bool)), Roles: []g.Role{g.One("signature", "function-type")}},
	"call-signature":      {Roles: []g.Role{g.One("signature", "function-type")}},
	"index-signature":     {Attributes: g.Attributes("readonly", g.Attr(false, g.Bool)), Roles: []g.Role{g.One("parameter", "parameter"), g.One("type", types...)}},
	"type-operator":       {Attributes: g.Attributes("operator", g.Attr(true, g.Set("keyof", "readonly", "unique"))), Roles: []g.Role{g.One("type", types...)}},
	"indexed-access-type": {Roles: []g.Role{g.One("object", types...), g.One("index", types...)}},
}

func Check(view ast.View[s.Language]) error {
	t := view.ASTView()
	if err := g.Check(t, "source-file", rules); err != nil {
		return err
	}
	v := tree.View[s.Language]{Tree: t}
	for i, n := range t.Nodes {
		id := ast.NodeID(i + 1)
		bad := func(m string) error { return text.AtNode(t, id, m) }
		switch n.Kind {
		case s.KindImport:
			b := v.Children(id, "binding")
			for _, x := range b {
				if v.Kind(x) == "namespace-import" && len(b) != 1 {
					return bad("namespace import must be the only binding")
				}
			}
		case s.KindLiteralType:
			value := v.Attr(id, "text")
			valid := false
			switch v.Attr(id, "literal-kind") {
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
			child := v.Child(id, "type")
			op := v.Attr(id, "operator")
			if op == "readonly" && v.Kind(child) != "array-type" && v.Kind(child) != "tuple-type" {
				return bad("readonly applies to array or tuple syntax")
			}
			if op == "unique" && (v.Kind(child) != "keyword-type" || v.Attr(child, "name") != "symbol") {
				return bad("unique applies to symbol")
			}
		case s.KindFunctionType, s.KindTupleType:
			role := "parameter"
			if n.Kind == s.KindTupleType {
				role = "element"
			}
			list := v.Children(id, role)
			optional, rest := false, false
			named := -1
			for j, p := range list {
				o, r := v.Flag(p, "optional"), v.Flag(p, "rest")
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
					if v.Has(p, "name") {
						mode = 1
					}
					if named >= 0 && named != mode {
						return text.AtNode(t, p, "tuple labels must be consistent")
					}
					named = mode
				}
			}
		case s.KindIndexSignature:
			p := v.Child(id, "parameter")
			if v.Flag(p, "optional") || v.Flag(p, "rest") {
				return bad("index parameter cannot be optional or rest")
			}
		}
	}
	return nil
}

type Provider struct{ s.Language }

func (Provider) Check(view ast.View[s.Language]) error { return Check(view) }

var _ text.CheckerProvider[s.Language] = Provider{}
