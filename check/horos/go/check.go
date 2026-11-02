// Package check admits Horos declaration syntax without resolving its names.
// Semantic rules and the truth of optional binding annotations belong to analysis.
package check

import (
	"reflect"
	"slices"

	structure "github.com/Bitspark/bittype/check/ast/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	g "github.com/Bitspark/bittype/internal/grammar"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func identifier(value string) bool {
	if value == "" {
		return false
	}
	for i, c := range []byte(value) {
		if c != '_' && !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') && !(i > 0 && c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

func moduleName(value string) bool {
	first := true
	for i := 0; i < len(value); i++ {
		c := value[i]
		if c == '.' {
			if first {
				return false
			}
			first = true
			continue
		}
		if !(c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || !first && c >= '0' && c <= '9') {
			return false
		}
		first = false
	}
	return !first
}

func hexadecimal(value string) bool {
	if len(value)%2 != 0 {
		return false
	}
	for _, c := range []byte(value) {
		if !(c >= '0' && c <= '9') && !(c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}

func pin(value string) bool { return len(value) == 64 && hexadecimal(value) }

var expressions = []string{"descriptor", "parameter", "application"}
var fragments = []string{"descriptor", "parameter", "application", "atom", "tuple"}
var rules = g.Rules{
	"program":        {Roles: []g.Role{g.Many("module", 1, "module")}},
	"module":         {Attributes: g.Attributes("name", g.Attr(true, moduleName), "pin", g.Attr(true, pin)), Roles: []g.Role{g.Many("import", 0, "import"), g.Many("declaration", 0, "declaration")}},
	"import":         {Attributes: g.Attributes("alias", g.Attr(true, identifier), "name", g.Attr(true, moduleName), "pin", g.Attr(true, pin))},
	"declaration":    {Attributes: g.Attributes("name", g.Attr(true, identifier)), Roles: []g.Role{g.Many("parameter", 0, "type-parameter"), g.One("body", expressions...)}},
	"type-parameter": {Attributes: g.Attributes("name", g.Attr(true, identifier), "index", g.Attr(true, g.Decimal))},
	"parameter":      {Attributes: g.Attributes("index", g.Attr(true, g.Decimal))},
	"application":    {Attributes: g.Attributes("scope", g.Attr(true, g.Set("local", "imported")), "name", g.Attr(true, identifier), "alias", g.Attr(false, identifier)), Roles: []g.Role{g.Many("argument", 0, expressions...)}},
	"descriptor":     {Attributes: g.Attributes("constructor", g.Attr(true, g.Nonempty)), Roles: []g.Role{g.Many("operand", 0, fragments...)}},
	"atom":           {Attributes: g.Attributes("hex", g.Attr(true, hexadecimal))},
	"tuple":          {Roles: []g.Role{g.Many("element", 0, fragments...)}},
}

func attribute(node s.Node, name string) string {
	for _, value := range node.Attributes {
		if value.Name == name {
			return value.Text
		}
	}
	return ""
}

func Check(view ast.View[s.Language]) error {
	if view == nil {
		return &text.Diagnostic{Code: text.InvalidAST, Message: "missing Horos syntax"}
	}
	v := reflect.ValueOf(view)
	switch v.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface, reflect.Func, reflect.Chan:
		if v.IsNil() {
			return &text.Diagnostic{Code: text.InvalidAST, Message: "missing Horos syntax"}
		}
	}
	tree := view.ASTView()
	if err := structure.Structure(tree); err != nil {
		return err
	}
	// The shared grammar checker supports the native target role. Horos has
	// several annotation roles, checked below without changing caller storage.
	shape := tree
	shape.Nodes = slices.Clone(tree.Nodes)
	for i := range shape.Nodes {
		shape.Nodes[i].References = nil
	}
	if err := g.Check(shape, "program", rules); err != nil {
		return err
	}
	rootTarget := ast.NodeID(0)
	for i, node := range tree.Nodes {
		id := ast.NodeID(i + 1)
		if node.Kind == s.KindApplication {
			imported := attribute(node, s.AttributeScope) == "imported"
			if imported != (attribute(node, s.AttributeAlias) != "") {
				return text.AtNode(tree, id, "only imported applications require an alias")
			}
		}
		for _, ref := range node.References {
			allowed := false
			switch node.Kind {
			case s.KindProgram:
				allowed = ref.Role == s.RoleRoot && ref.Target.Tree == tree.ID
				if allowed {
					rootTarget = ref.Target.Node
				}
			case s.KindImport, s.KindParameter:
				allowed = ref.Role == s.RoleTarget
			case s.KindApplication:
				allowed = ref.Role == s.RoleTarget || (ref.Role == s.RoleImport && attribute(node, s.AttributeScope) == "imported")
			}
			if !allowed {
				return text.AtNode(tree, id, "unsupported Horos annotation role")
			}
		}
	}
	for _, child := range tree.Nodes[tree.Root-1].Children {
		if child.Node == rootTarget {
			return nil
		}
	}
	return text.AtNode(tree, tree.Root, "root reference must select a contained module")
}
