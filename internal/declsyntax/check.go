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
	spelling "github.com/Bitspark/bittype/internal/horos"
	"github.com/Bitspark/bittype/model/syntax/go"
	predicate "github.com/Bitspark/bittype/model/syntax/horos/lexical/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

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

// evaluators binds the names of the horos predicates that the table uses. The
// identifier rule is the lexical helper's, so it has one evaluator.
var evaluators = g.Evaluators{
	predicate.Identifier():  spelling.Identifier,
	predicate.ModuleName():  moduleName,
	predicate.Hexadecimal(): hexadecimal,
	predicate.Pin():         pin,
}

var expressions = []h.Kind{h.KindDescriptor, h.KindParameter, h.KindReference, h.KindApplication}
var fragments = []h.Kind{h.KindDescriptor, h.KindParameter, h.KindReference, h.KindApplication, h.KindAtom, h.KindTuple}
var sorts = []h.Kind{h.KindProperSort, h.KindArrowSort}
var kinds = []h.Kind{h.KindTop, h.KindBottom, h.KindKindParameterReference, h.KindKindReference, h.KindKindApplication, h.KindUnion, h.KindIntersection, h.KindArrow, h.KindNodePattern, h.KindRecursiveKind, h.KindRecursiveVariable}
var sequences = []h.Kind{h.KindSequenceEpsilon, h.KindSequenceItem, h.KindSequenceConcat, h.KindSequenceAlternative, h.KindSequenceStar}
var rules = syntax.Rules[h.Language]{
	h.KindProgram:                {Roles: []syntax.Role[h.Language]{{Name: RoleModule, Min: 1, Max: -1, Kinds: []h.Kind{h.KindModule}}}},
	h.KindModule:                 {Attributes: map[string]syntax.Attribute{AttributeName: {Required: true, Predicate: predicate.ModuleName()}, AttributePin: {Required: true, Predicate: predicate.Pin()}}, Roles: []syntax.Role[h.Language]{{Name: RoleImport, Max: -1, Kinds: []h.Kind{h.KindImport}}, {Name: RoleDeclaration, Max: -1, Kinds: []h.Kind{h.KindTypeDefinition, h.KindKindDefinition}}}},
	h.KindImport:                 {Attributes: map[string]syntax.Attribute{AttributeAlias: {Required: true, Predicate: predicate.Identifier()}, AttributeName: {Required: true, Predicate: predicate.ModuleName()}, AttributePin: {Required: true, Predicate: predicate.Pin()}}},
	h.KindTypeDefinition:         {Attributes: map[string]syntax.Attribute{AttributeName: {Required: true, Predicate: predicate.Identifier()}}, Roles: []syntax.Role[h.Language]{{Name: RoleParameter, Max: -1, Kinds: []h.Kind{h.KindTypeParameter}}, {Name: RoleResult, Min: 1, Max: 1, Kinds: kinds}, {Name: RoleBody, Min: 1, Max: 1, Kinds: expressions}}},
	h.KindKindDefinition:         {Attributes: map[string]syntax.Attribute{AttributeName: {Required: true, Predicate: predicate.Identifier()}}, Roles: []syntax.Role[h.Language]{{Name: RoleParameter, Max: -1, Kinds: []h.Kind{h.KindKindParameter}}, {Name: RoleResult, Min: 1, Max: 1, Kinds: sorts}, {Name: RoleBody, Min: 1, Max: 1, Kinds: kinds}}},
	h.KindTypeParameter:          {Attributes: map[string]syntax.Attribute{AttributeName: {Required: true, Predicate: predicate.Identifier()}, AttributeIndex: {Required: true, Predicate: syntax.Decimal()}}, Roles: []syntax.Role[h.Language]{{Name: RoleContract, Min: 1, Max: 1, Kinds: kinds}}},
	h.KindKindParameter:          {Attributes: map[string]syntax.Attribute{AttributeName: {Required: true, Predicate: predicate.Identifier()}, AttributeIndex: {Required: true, Predicate: syntax.Decimal()}}, Roles: []syntax.Role[h.Language]{{Name: RoleSort, Min: 1, Max: 1, Kinds: sorts}}},
	h.KindParameter:              {Attributes: map[string]syntax.Attribute{AttributeIndex: {Required: true, Predicate: syntax.Decimal()}}},
	h.KindReference:              {Attributes: map[string]syntax.Attribute{AttributeScope: {Required: true, OneOf: []string{h.ScopeLocal, h.ScopeImported}}, AttributeName: {Required: true, Predicate: predicate.Identifier()}, AttributeAlias: {Predicate: predicate.Identifier()}}},
	h.KindApplication:            {Roles: []syntax.Role[h.Language]{{Name: RoleCallee, Min: 1, Max: 1, Kinds: expressions}, {Name: RoleArgument, Min: 1, Max: -1, Kinds: expressions}}},
	h.KindDescriptor:             {Attributes: map[string]syntax.Attribute{AttributeConstructor: {Required: true, Predicate: syntax.Nonempty()}}, Roles: []syntax.Role[h.Language]{{Name: RoleOperand, Max: -1, Kinds: fragments}}},
	h.KindAtom:                   {Attributes: map[string]syntax.Attribute{AttributeHex: {Required: true, Predicate: predicate.Hexadecimal()}}},
	h.KindTuple:                  {Roles: []syntax.Role[h.Language]{{Name: RoleElement, Max: -1, Kinds: fragments}}},
	h.KindProperSort:             {},
	h.KindArrowSort:              {Roles: []syntax.Role[h.Language]{{Name: RoleInput, Min: 1, Max: -1, Kinds: sorts}, {Name: RoleResult, Min: 1, Max: 1, Kinds: sorts}}},
	h.KindTop:                    {Roles: []syntax.Role[h.Language]{{Name: RoleSort, Min: 1, Max: 1, Kinds: sorts}}},
	h.KindBottom:                 {Roles: []syntax.Role[h.Language]{{Name: RoleSort, Min: 1, Max: 1, Kinds: sorts}}},
	h.KindKindParameterReference: {Attributes: map[string]syntax.Attribute{AttributeIndex: {Required: true, Predicate: syntax.Decimal()}}},
	h.KindKindReference:          {Attributes: map[string]syntax.Attribute{AttributeScope: {Required: true, OneOf: []string{h.ScopeLocal, h.ScopeImported}}, AttributeName: {Required: true, Predicate: predicate.Identifier()}, AttributeAlias: {Predicate: predicate.Identifier()}}},
	h.KindKindApplication:        {Roles: []syntax.Role[h.Language]{{Name: RoleCallee, Min: 1, Max: 1, Kinds: []h.Kind{h.KindKindReference}}, {Name: RoleArgument, Min: 1, Max: -1, Kinds: kinds}}},
	h.KindUnion:                  {Roles: []syntax.Role[h.Language]{{Name: RoleSort, Min: 1, Max: 1, Kinds: sorts}, {Name: RoleOperand, Min: 2, Max: -1, Kinds: kinds}}},
	h.KindIntersection:           {Roles: []syntax.Role[h.Language]{{Name: RoleSort, Min: 1, Max: 1, Kinds: sorts}, {Name: RoleOperand, Min: 2, Max: -1, Kinds: kinds}}},
	h.KindArrow:                  {Roles: []syntax.Role[h.Language]{{Name: RoleInput, Min: 1, Max: -1, Kinds: kinds}, {Name: RoleResult, Min: 1, Max: 1, Kinds: kinds}}},
	h.KindNodePattern:            {Attributes: map[string]syntax.Attribute{AttributeConstructor: {Required: true, OneOf: []string{"any", "atom", "atom-exact", "atom-len", "atom-in", "text-len", "key-grammar", "text-non-blank", "lowercase-uuid", "scalar", "exact", "wire", "seq", "list-of", "set-of", "not", "map-of", "tuple", "compound", "union", "and", "record"}}}, Roles: []syntax.Role[h.Language]{{Name: RoleSelector, Max: -1, Kinds: []h.Kind{h.KindSelector}}, {Name: RolePattern, Min: 1, Max: 1, Kinds: sequences}}},
	h.KindSelector:               {Attributes: map[string]syntax.Attribute{AttributeName: {Required: true, OneOf: []string{"selector", "mode"}}, AttributeValue: {Required: true, OneOf: []string{"int", "utf8-text", "bool", "open", "closed"}}}},
	h.KindRecursiveKind:          {Attributes: map[string]syntax.Attribute{AttributeIndex: {Required: true, Predicate: syntax.Decimal()}}, Roles: []syntax.Role[h.Language]{{Name: RoleEquation, Min: 1, Max: -1, Kinds: []h.Kind{h.KindEquation}}}},
	h.KindEquation:               {Attributes: map[string]syntax.Attribute{AttributeName: {Predicate: predicate.Identifier()}, AttributeIndex: {Required: true, Predicate: syntax.Decimal()}}, Roles: []syntax.Role[h.Language]{{Name: RoleBody, Min: 1, Max: 1, Kinds: kinds}}},
	h.KindRecursiveVariable:      {Attributes: map[string]syntax.Attribute{AttributeDepth: {Required: true, Predicate: syntax.Decimal()}, AttributeIndex: {Required: true, Predicate: syntax.Decimal()}}},
	h.KindSequenceEpsilon:        {},
	h.KindSequenceItem:           {Roles: []syntax.Role[h.Language]{{Name: RoleElement, Min: 1, Max: 1, Kinds: kinds}}},
	h.KindSequenceConcat:         {Roles: []syntax.Role[h.Language]{{Name: RoleElement, Min: 2, Max: -1, Kinds: sequences}}},
	h.KindSequenceAlternative:    {Roles: []syntax.Role[h.Language]{{Name: RoleElement, Min: 2, Max: -1, Kinds: sequences}}},
	h.KindSequenceStar:           {Roles: []syntax.Role[h.Language]{{Name: RoleElement, Min: 1, Max: 1, Kinds: sequences}}},
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
	if err := g.Check(shape, h.KindProgram, rules, evaluators); err != nil {
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
