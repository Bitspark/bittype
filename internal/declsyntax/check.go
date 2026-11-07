// Package declsyntax admits Horos declaration syntax without resolving names.
// Semantic rules and the truth of binding annotations belong to resolution.
package declsyntax

import (
	"reflect"
	"slices"
	"strconv"

	structure "github.com/Bitspark/bittype/check/syntax/shared/ast/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	g "github.com/Bitspark/bittype/internal/grammar"
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

var expressions = []string{KindDescriptor, KindParameter, KindReference, KindApplication}
var fragments = []string{KindDescriptor, KindParameter, KindReference, KindApplication, KindAtom, KindTuple}
var sorts = []string{KindProperSort, KindArrowSort}
var kinds = []string{KindTop, KindBottom, KindKindParameterReference, KindKindReference, KindKindApplication, KindUnion, KindIntersection, KindArrow, KindNodePattern, KindRecursiveKind, KindRecursiveVariable}
var sequences = []string{KindSequenceEpsilon, KindSequenceItem, KindSequenceConcat, KindSequenceAlternative, KindSequenceStar}
var rules = g.Rules{
	KindProgram:                {Roles: []g.Role{g.Many(RoleModule, 1, KindModule)}},
	KindModule:                 {Attributes: g.Attributes(AttributeName, g.Attr(true, moduleName), AttributePin, g.Attr(true, pin)), Roles: []g.Role{g.Many(RoleImport, 0, KindImport), g.Many(RoleDeclaration, 0, KindTypeDefinition, KindKindDefinition)}},
	KindImport:                 {Attributes: g.Attributes(AttributeAlias, g.Attr(true, identifier), AttributeName, g.Attr(true, moduleName), AttributePin, g.Attr(true, pin))},
	KindTypeDefinition:         {Attributes: g.Attributes(AttributeName, g.Attr(true, identifier)), Roles: []g.Role{g.Many(RoleParameter, 0, KindTypeParameter), g.One(RoleResult, kinds...), g.One(RoleBody, expressions...)}},
	KindKindDefinition:         {Attributes: g.Attributes(AttributeName, g.Attr(true, identifier)), Roles: []g.Role{g.Many(RoleParameter, 0, KindKindParameter), g.One(RoleResult, sorts...), g.One(RoleBody, kinds...)}},
	KindTypeParameter:          {Attributes: g.Attributes(AttributeName, g.Attr(true, identifier), AttributeIndex, g.Attr(true, g.Decimal)), Roles: []g.Role{g.One(RoleContract, kinds...)}},
	KindKindParameter:          {Attributes: g.Attributes(AttributeName, g.Attr(true, identifier), AttributeIndex, g.Attr(true, g.Decimal)), Roles: []g.Role{g.One(RoleSort, sorts...)}},
	KindParameter:              {Attributes: g.Attributes(AttributeIndex, g.Attr(true, g.Decimal))},
	KindReference:              {Attributes: g.Attributes(AttributeScope, g.Attr(true, g.Set("local", "imported")), AttributeName, g.Attr(true, identifier), AttributeAlias, g.Attr(false, identifier))},
	KindApplication:            {Roles: []g.Role{g.One(RoleCallee, expressions...), g.Many(RoleArgument, 1, expressions...)}},
	KindDescriptor:             {Attributes: g.Attributes(AttributeConstructor, g.Attr(true, g.Nonempty)), Roles: []g.Role{g.Many(RoleOperand, 0, fragments...)}},
	KindAtom:                   {Attributes: g.Attributes(AttributeHex, g.Attr(true, hexadecimal))},
	KindTuple:                  {Roles: []g.Role{g.Many(RoleElement, 0, fragments...)}},
	KindProperSort:             {},
	KindArrowSort:              {Roles: []g.Role{g.Many(RoleInput, 1, sorts...), g.One(RoleResult, sorts...)}},
	KindTop:                    {Roles: []g.Role{g.One(RoleSort, sorts...)}},
	KindBottom:                 {Roles: []g.Role{g.One(RoleSort, sorts...)}},
	KindKindParameterReference: {Attributes: g.Attributes(AttributeIndex, g.Attr(true, g.Decimal))},
	KindKindReference:          {Attributes: g.Attributes(AttributeScope, g.Attr(true, g.Set("local", "imported")), AttributeName, g.Attr(true, identifier), AttributeAlias, g.Attr(false, identifier))},
	KindKindApplication:        {Roles: []g.Role{g.One(RoleCallee, KindKindReference), g.Many(RoleArgument, 1, kinds...)}},
	KindUnion:                  {Roles: []g.Role{g.One(RoleSort, sorts...), g.Many(RoleOperand, 2, kinds...)}},
	KindIntersection:           {Roles: []g.Role{g.One(RoleSort, sorts...), g.Many(RoleOperand, 2, kinds...)}},
	KindArrow:                  {Roles: []g.Role{g.Many(RoleInput, 1, kinds...), g.One(RoleResult, kinds...)}},
	KindNodePattern:            {Attributes: g.Attributes(AttributeConstructor, g.Attr(true, g.Set("any", "atom", "atom-exact", "atom-len", "atom-in", "text-len", "key-grammar", "text-non-blank", "lowercase-uuid", "scalar", "exact", "wire", "seq", "list-of", "set-of", "not", "map-of", "tuple", "compound", "union", "and", "record"))), Roles: []g.Role{g.Many(RoleSelector, 0, KindSelector), g.One(RolePattern, sequences...)}},
	KindSelector:               {Attributes: g.Attributes(AttributeName, g.Attr(true, g.Set("selector", "mode")), AttributeValue, g.Attr(true, g.Set("int", "utf8-text", "bool", "open", "closed")))},
	KindRecursiveKind:          {Attributes: g.Attributes(AttributeIndex, g.Attr(true, g.Decimal)), Roles: []g.Role{g.Many(RoleEquation, 1, KindEquation)}},
	KindEquation:               {Attributes: g.Attributes(AttributeName, g.Attr(false, identifier), AttributeIndex, g.Attr(true, g.Decimal)), Roles: []g.Role{g.One(RoleBody, kinds...)}},
	KindRecursiveVariable:      {Attributes: g.Attributes(AttributeDepth, g.Attr(true, g.Decimal), AttributeIndex, g.Attr(true, g.Decimal))},
	KindSequenceEpsilon:        {},
	KindSequenceItem:           {Roles: []g.Role{g.One(RoleElement, kinds...)}},
	KindSequenceConcat:         {Roles: []g.Role{g.Many(RoleElement, 2, sequences...)}},
	KindSequenceAlternative:    {Roles: []g.Role{g.Many(RoleElement, 2, sequences...)}},
	KindSequenceStar:           {Roles: []g.Role{g.One(RoleElement, sequences...)}},
}

func attribute[L any](node ast.Node[L], name string) string {
	for _, value := range node.Attributes {
		if value.Name == name {
			return value.Text
		}
	}
	return ""
}

func Check[L any](view ast.View[L]) error {
	if view == nil {
		return &text.Diagnostic{Code: text.InvalidAST, Message: "missing declaration syntax"}
	}
	v := reflect.ValueOf(view)
	switch v.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface, reflect.Func, reflect.Chan:
		if v.IsNil() {
			return &text.Diagnostic{Code: text.InvalidAST, Message: "missing declaration syntax"}
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
	rootReferences := 0
	for i, node := range tree.Nodes {
		id := ast.NodeID(i + 1)
		if node.Kind == KindReference || node.Kind == KindKindReference {
			imported := attribute(node, AttributeScope) == "imported"
			if imported != (attribute(node, AttributeAlias) != "") {
				return text.AtNode(tree, id, "only imported references require an alias")
			}
		}
		if node.Kind == KindNodePattern {
			seen := map[string]bool{}
			for _, child := range node.Children {
				if child.Role != RoleSelector {
					continue
				}
				selector := tree.Nodes[child.Node-1]
				key, value := attribute(selector, AttributeName), attribute(selector, AttributeValue)
				valid := false
				switch attribute(node, AttributeConstructor) {
				case "scalar":
					valid = key == "selector" && (value == "int" || value == "utf8-text" || value == "bool")
				case "record":
					valid = key == "mode" && (value == "open" || value == "closed")
				}
				if !valid || seen[key] {
					return text.AtNode(tree, child.Node, "invalid or repeated node selector")
				}
				seen[key] = true
			}
		}
		for _, ref := range node.References {
			allowed := false
			switch node.Kind {
			case KindProgram:
				allowed = ref.Role == RoleRoot && ref.Target.Tree == tree.ID
				if allowed {
					rootTarget = ref.Target.Node
					rootReferences++
				}
			case KindImport, KindParameter, KindKindParameterReference, KindRecursiveVariable:
				allowed = ref.Role == RoleTarget
			case KindReference, KindKindReference:
				allowed = ref.Role == RoleTarget || (ref.Role == RoleImport && attribute(node, AttributeScope) == "imported")
			}
			if !allowed {
				return text.AtNode(tree, id, "unsupported declaration annotation role")
			}
		}
	}
	if rootReferences == 1 {
		for _, child := range tree.Nodes[tree.Root-1].Children {
			if child.Node == rootTarget {
				return lexical(tree, tree.Root, 0, 0, nil)
			}
		}
	}
	return text.AtNode(tree, tree.Root, "root reference must select a contained module")
}

// Lexical slots are syntax. Declaration lookup, kind formation and application
// admission remain the authoritative resolver's responsibility.
func lexical[L any](tree ast.AST[L], id ast.NodeID, typeParameters, kindParameters int, groups []int) error {
	node := tree.Nodes[id-1]
	index := func(name string, count int) bool {
		n, err := strconv.Atoi(attribute(node, name))
		return err == nil && n >= 0 && n < count
	}
	switch node.Kind {
	case KindModule:
		imports, declarations := map[string]bool{}, map[string]bool{}
		for _, child := range node.Children {
			value := tree.Nodes[child.Node-1]
			name, seen := attribute(value, AttributeName), declarations
			if child.Role == RoleImport {
				name, seen = attribute(value, AttributeAlias), imports
			}
			if seen[name] {
				return text.AtNode(tree, child.Node, "duplicate declaration name or import alias")
			}
			seen[name] = true
		}
	case KindTypeDefinition, KindKindDefinition:
		count := 0
		seen := map[string]bool{}
		for _, child := range node.Children {
			if child.Role != RoleParameter {
				continue
			}
			parameter := tree.Nodes[child.Node-1]
			name := attribute(parameter, AttributeName)
			if attribute(parameter, AttributeIndex) != strconv.Itoa(count) || seen[name] {
				return text.AtNode(tree, child.Node, "parameters require distinct names and ordered contiguous indices")
			}
			seen[name] = true
			count++
		}
		typeParameters, kindParameters, groups = 0, 0, nil
		if node.Kind == KindTypeDefinition {
			typeParameters = count
		} else {
			kindParameters = count
		}
	case KindParameter:
		if !index(AttributeIndex, typeParameters) {
			return text.AtNode(tree, id, "type parameter index has no source binder")
		}
	case KindKindParameterReference:
		if !index(AttributeIndex, kindParameters) {
			return text.AtNode(tree, id, "kind parameter index has no source binder")
		}
	case KindRecursiveKind:
		if !index(AttributeIndex, len(node.Children)) {
			return text.AtNode(tree, id, "recursive kind entry has no equation")
		}
		for i, child := range node.Children {
			if attribute(tree.Nodes[child.Node-1], AttributeIndex) != strconv.Itoa(i) {
				return text.AtNode(tree, child.Node, "equations require ordered contiguous indices")
			}
		}
		groups = append(groups, len(node.Children))
	case KindRecursiveVariable:
		depth, err := strconv.Atoi(attribute(node, AttributeDepth))
		if err != nil || depth < 0 || depth >= len(groups) || !index(AttributeIndex, groups[len(groups)-1-depth]) {
			return text.AtNode(tree, id, "recursive variable has no lexical equation")
		}
	}
	for _, child := range node.Children {
		if err := lexical(tree, child.Node, typeParameters, kindParameters, groups); err != nil {
			return err
		}
	}
	return nil
}
