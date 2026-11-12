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
	h "github.com/Bitspark/bittype/syntax/horos/go"
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

var expressions = []h.Kind{h.KindDescriptor, h.KindParameter, h.KindReference, h.KindApplication}
var fragments = []h.Kind{h.KindDescriptor, h.KindParameter, h.KindReference, h.KindApplication, h.KindAtom, h.KindTuple}
var sorts = []h.Kind{h.KindProperSort, h.KindArrowSort}
var kinds = []h.Kind{h.KindTop, h.KindBottom, h.KindKindParameterReference, h.KindKindReference, h.KindKindApplication, h.KindUnion, h.KindIntersection, h.KindArrow, h.KindNodePattern, h.KindRecursiveKind, h.KindRecursiveVariable}
var sequences = []h.Kind{h.KindSequenceEpsilon, h.KindSequenceItem, h.KindSequenceConcat, h.KindSequenceAlternative, h.KindSequenceStar}
var rules = g.Rules[h.Language]{
	h.KindProgram:                {Roles: []g.Role[h.Language]{g.Many(RoleModule, 1, h.KindModule)}},
	h.KindModule:                 {Attributes: g.Attributes(AttributeName, g.Attr(true, moduleName), AttributePin, g.Attr(true, pin)), Roles: []g.Role[h.Language]{g.Many(RoleImport, 0, h.KindImport), g.Many(RoleDeclaration, 0, h.KindTypeDefinition, h.KindKindDefinition)}},
	h.KindImport:                 {Attributes: g.Attributes(AttributeAlias, g.Attr(true, identifier), AttributeName, g.Attr(true, moduleName), AttributePin, g.Attr(true, pin))},
	h.KindTypeDefinition:         {Attributes: g.Attributes(AttributeName, g.Attr(true, identifier)), Roles: []g.Role[h.Language]{g.Many(RoleParameter, 0, h.KindTypeParameter), g.One(RoleResult, kinds...), g.One(RoleBody, expressions...)}},
	h.KindKindDefinition:         {Attributes: g.Attributes(AttributeName, g.Attr(true, identifier)), Roles: []g.Role[h.Language]{g.Many(RoleParameter, 0, h.KindKindParameter), g.One(RoleResult, sorts...), g.One(RoleBody, kinds...)}},
	h.KindTypeParameter:          {Attributes: g.Attributes(AttributeName, g.Attr(true, identifier), AttributeIndex, g.Attr(true, g.Decimal)), Roles: []g.Role[h.Language]{g.One(RoleContract, kinds...)}},
	h.KindKindParameter:          {Attributes: g.Attributes(AttributeName, g.Attr(true, identifier), AttributeIndex, g.Attr(true, g.Decimal)), Roles: []g.Role[h.Language]{g.One(RoleSort, sorts...)}},
	h.KindParameter:              {Attributes: g.Attributes(AttributeIndex, g.Attr(true, g.Decimal))},
	h.KindReference:              {Attributes: g.Attributes(AttributeScope, g.Attr(true, g.Set(h.ScopeLocal, h.ScopeImported)), AttributeName, g.Attr(true, identifier), AttributeAlias, g.Attr(false, identifier))},
	h.KindApplication:            {Roles: []g.Role[h.Language]{g.One(RoleCallee, expressions...), g.Many(RoleArgument, 1, expressions...)}},
	h.KindDescriptor:             {Attributes: g.Attributes(AttributeConstructor, g.Attr(true, g.Nonempty)), Roles: []g.Role[h.Language]{g.Many(RoleOperand, 0, fragments...)}},
	h.KindAtom:                   {Attributes: g.Attributes(AttributeHex, g.Attr(true, hexadecimal))},
	h.KindTuple:                  {Roles: []g.Role[h.Language]{g.Many(RoleElement, 0, fragments...)}},
	h.KindProperSort:             {},
	h.KindArrowSort:              {Roles: []g.Role[h.Language]{g.Many(RoleInput, 1, sorts...), g.One(RoleResult, sorts...)}},
	h.KindTop:                    {Roles: []g.Role[h.Language]{g.One(RoleSort, sorts...)}},
	h.KindBottom:                 {Roles: []g.Role[h.Language]{g.One(RoleSort, sorts...)}},
	h.KindKindParameterReference: {Attributes: g.Attributes(AttributeIndex, g.Attr(true, g.Decimal))},
	h.KindKindReference:          {Attributes: g.Attributes(AttributeScope, g.Attr(true, g.Set(h.ScopeLocal, h.ScopeImported)), AttributeName, g.Attr(true, identifier), AttributeAlias, g.Attr(false, identifier))},
	h.KindKindApplication:        {Roles: []g.Role[h.Language]{g.One(RoleCallee, h.KindKindReference), g.Many(RoleArgument, 1, kinds...)}},
	h.KindUnion:                  {Roles: []g.Role[h.Language]{g.One(RoleSort, sorts...), g.Many(RoleOperand, 2, kinds...)}},
	h.KindIntersection:           {Roles: []g.Role[h.Language]{g.One(RoleSort, sorts...), g.Many(RoleOperand, 2, kinds...)}},
	h.KindArrow:                  {Roles: []g.Role[h.Language]{g.Many(RoleInput, 1, kinds...), g.One(RoleResult, kinds...)}},
	h.KindNodePattern:            {Attributes: g.Attributes(AttributeConstructor, g.Attr(true, g.Set("any", "atom", "atom-exact", "atom-len", "atom-in", "text-len", "key-grammar", "text-non-blank", "lowercase-uuid", "scalar", "exact", "wire", "seq", "list-of", "set-of", "not", "map-of", "tuple", "compound", "union", "and", "record"))), Roles: []g.Role[h.Language]{g.Many(RoleSelector, 0, h.KindSelector), g.One(RolePattern, sequences...)}},
	h.KindSelector:               {Attributes: g.Attributes(AttributeName, g.Attr(true, g.Set("selector", "mode")), AttributeValue, g.Attr(true, g.Set("int", "utf8-text", "bool", "open", "closed")))},
	h.KindRecursiveKind:          {Attributes: g.Attributes(AttributeIndex, g.Attr(true, g.Decimal)), Roles: []g.Role[h.Language]{g.Many(RoleEquation, 1, h.KindEquation)}},
	h.KindEquation:               {Attributes: g.Attributes(AttributeName, g.Attr(false, identifier), AttributeIndex, g.Attr(true, g.Decimal)), Roles: []g.Role[h.Language]{g.One(RoleBody, kinds...)}},
	h.KindRecursiveVariable:      {Attributes: g.Attributes(AttributeDepth, g.Attr(true, g.Decimal), AttributeIndex, g.Attr(true, g.Decimal))},
	h.KindSequenceEpsilon:        {},
	h.KindSequenceItem:           {Roles: []g.Role[h.Language]{g.One(RoleElement, kinds...)}},
	h.KindSequenceConcat:         {Roles: []g.Role[h.Language]{g.Many(RoleElement, 2, sequences...)}},
	h.KindSequenceAlternative:    {Roles: []g.Role[h.Language]{g.Many(RoleElement, 2, sequences...)}},
	h.KindSequenceStar:           {Roles: []g.Role[h.Language]{g.One(RoleElement, sequences...)}},
}

func attribute(node h.Node, name string) string {
	for _, value := range node.Attributes {
		if value.Name == name {
			return value.Text
		}
	}
	return ""
}

func Check(view ast.View[h.Language]) error {
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
	if err := g.Check(shape, h.KindProgram, rules); err != nil {
		return err
	}
	rootTarget := ast.NodeID(0)
	rootReferences := 0
	for i, node := range tree.Nodes {
		id := ast.NodeID(i + 1)
		if node.Kind == h.KindReference || node.Kind == h.KindKindReference {
			imported := attribute(node, AttributeScope) == h.ScopeImported
			if imported != (attribute(node, AttributeAlias) != "") {
				return text.AtNode(tree, id, "only imported references require an alias")
			}
		}
		if node.Kind == h.KindNodePattern {
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
			case h.KindProgram:
				allowed = ref.Role == RoleRoot && ref.Target.Tree == tree.ID
				if allowed {
					rootTarget = ref.Target.Node
					rootReferences++
				}
			case h.KindImport, h.KindParameter, h.KindKindParameterReference, h.KindRecursiveVariable:
				allowed = ref.Role == RoleTarget
			case h.KindReference, h.KindKindReference:
				allowed = ref.Role == RoleTarget || (ref.Role == RoleImport && attribute(node, AttributeScope) == h.ScopeImported)
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
func lexical(tree h.AST, id ast.NodeID, typeParameters, kindParameters int, groups []int) error {
	node := tree.Nodes[id-1]
	index := func(name string, count int) bool {
		n, err := strconv.Atoi(attribute(node, name))
		return err == nil && n >= 0 && n < count
	}
	switch node.Kind {
	case h.KindModule:
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
	case h.KindTypeDefinition, h.KindKindDefinition:
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
		if node.Kind == h.KindTypeDefinition {
			typeParameters = count
		} else {
			kindParameters = count
		}
	case h.KindParameter:
		if !index(AttributeIndex, typeParameters) {
			return text.AtNode(tree, id, "type parameter index has no source binder")
		}
	case h.KindKindParameterReference:
		if !index(AttributeIndex, kindParameters) {
			return text.AtNode(tree, id, "kind parameter index has no source binder")
		}
	case h.KindRecursiveKind:
		if !index(AttributeIndex, len(node.Children)) {
			return text.AtNode(tree, id, "recursive kind entry has no equation")
		}
		for i, child := range node.Children {
			if attribute(tree.Nodes[child.Node-1], AttributeIndex) != strconv.Itoa(i) {
				return text.AtNode(tree, child.Node, "equations require ordered contiguous indices")
			}
		}
		groups = append(groups, len(node.Children))
	case h.KindRecursiveVariable:
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
