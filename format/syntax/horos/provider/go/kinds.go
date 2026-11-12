package format

import (
	"strconv"

	lit "github.com/Bitspark/bittype/internal/horos"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func (f *formatter) properTop(id ast.NodeID) bool {
	return f.Node(id).Kind == s.KindTop && f.Node(f.Child(id, s.RoleSort)).Kind == s.KindProperSort
}

func (f *formatter) sortExpression(id ast.NodeID) doc {
	if f.Node(id).Kind == s.KindProperSort {
		return f.Group(id, f.T("P"))
	}
	var inputs []doc
	nodes := f.Children(id, s.RoleInput)
	for _, input := range nodes {
		part := f.sortExpression(input)
		if len(nodes) == 1 && f.Node(input).Kind == s.KindArrowSort {
			part = f.Seq(f.T("("), part, f.T(")"))
		}
		inputs = append(inputs, part)
	}
	domain := f.Join(inputs, f.T(","), f.Space())
	if len(nodes) > 1 {
		domain = f.Seq(f.T("("), domain, f.T(")"))
	}
	return f.Group(id, domain, f.Space(), f.T("->"), f.Space(), f.sortExpression(f.Child(id, s.RoleResult)))
}

func (f *formatter) kindName(id ast.NodeID) doc {
	name := f.Attr(id, s.AttributeName)
	var parts []doc
	if f.Attr(id, s.AttributeScope) == s.ScopeImported {
		parts = append(parts, f.T(lit.Name(f.Attr(id, s.AttributeAlias))), f.T("."))
	} else {
		shadowed := false
		for _, parameter := range f.kindParameters {
			shadowed = shadowed || parameter == name
		}
		for _, group := range f.groups {
			for _, variable := range group {
				shadowed = shadowed || variable == name
			}
		}
		if shadowed {
			parts = append(parts, f.T("self"), f.T("."))
		}
	}
	return f.Seq(append(parts, f.T(lit.Name(name)))...)
}

func (f *formatter) kindExpression(id ast.NodeID, parent, level int) doc {
	if f.err != nil {
		return f.Seq()
	}
	prec := 4
	var result doc
	switch f.Node(id).Kind {
	case s.KindTop, s.KindBottom:
		name := "top"
		if f.Node(id).Kind == s.KindBottom {
			name = "bottom"
		}
		result = f.call(name, []doc{f.sortExpression(f.Child(id, s.RoleSort))})
	case s.KindKindParameterReference:
		index, err := strconv.Atoi(f.Attr(id, s.AttributeIndex))
		if err != nil || index < 0 || index >= len(f.kindParameters) {
			return f.bad(id, "kind parameter index has no source binder")
		}
		result = f.T(lit.Name(f.kindParameters[index]))
	case s.KindKindReference:
		result = f.kindName(id)
	case s.KindKindApplication:
		var arguments []doc
		for _, argument := range f.Children(id, s.RoleArgument) {
			arguments = append(arguments, f.kindExpression(argument, 0, level))
		}
		result = f.Seq(f.kindExpression(f.Child(id, s.RoleCallee), 0, level), f.T("<"), f.Join(arguments, f.T(","), f.Space()), f.T(">"))
	case s.KindUnion, s.KindIntersection:
		operator := "|"
		prec = 2
		if f.Node(id).Kind == s.KindIntersection {
			operator, prec = "&", 3
		}
		var operands []doc
		for _, operand := range f.Children(id, s.RoleOperand) {
			// Equal precedence also needs grouping: nested canonical Boolean
			// nodes are not silently flattened into a single operand list.
			operands = append(operands, f.kindExpression(operand, prec, level))
		}
		result = f.Join(operands, f.Space(), f.T(operator), f.Space())
		sort := f.Child(id, s.RoleSort)
		if f.Node(sort).Kind != s.KindProperSort {
			result = f.Seq(f.T("("), result, f.T(")"), f.Space(), f.T("at"), f.Space(), f.T("("), f.sortExpression(sort), f.T(")"))
			prec = 4
		}
	case s.KindArrow:
		prec = 1
		var inputs []doc
		nodes := f.Children(id, s.RoleInput)
		for _, input := range nodes {
			inputParent := 0
			if len(nodes) == 1 {
				inputParent = prec
			}
			inputs = append(inputs, f.kindExpression(input, inputParent, level))
		}
		domain := f.Join(inputs, f.T(","), f.Space())
		if len(nodes) > 1 {
			domain = f.Seq(f.T("("), domain, f.T(")"))
		}
		result = f.Seq(domain, f.Space(), f.T("->"), f.Space(), f.kindExpression(f.Child(id, s.RoleResult), 0, level))
	case s.KindNodePattern:
		var selectors []doc
		for _, selector := range f.Children(id, s.RoleSelector) {
			selectors = append(selectors, f.Group(selector, f.T(lit.Quote(f.Attr(selector, s.AttributeName))), f.T(":"), f.Space(), f.T(lit.Quote(f.Attr(selector, s.AttributeValue)))))
		}
		mapping := f.Seq(f.T("{"), f.Join(selectors, f.T(","), f.Space()), f.T("}"))
		result = f.call("node", []doc{f.T(lit.Quote(f.Attr(id, s.AttributeConstructor))), mapping, f.sequencePattern(f.Child(id, s.RolePattern), level)})
	case s.KindRecursiveKind:
		result = f.recursiveKind(id, level)
	case s.KindRecursiveVariable:
		depth, depthErr := strconv.Atoi(f.Attr(id, s.AttributeDepth))
		index, indexErr := strconv.Atoi(f.Attr(id, s.AttributeIndex))
		if depthErr != nil || indexErr != nil || depth < 0 || depth >= len(f.groups) || index < 0 || index >= len(f.groups[len(f.groups)-1-depth]) {
			return f.bad(id, "recursive variable has no lexical source binder")
		}
		result = f.T(lit.Name(f.groups[len(f.groups)-1-depth][index]))
	default:
		return f.bad(id, "expected kind expression")
	}
	if prec <= parent {
		result = f.Seq(f.T("("), result, f.T(")"))
	}
	return f.Group(id, result)
}

func (f *formatter) sequencePattern(id ast.NodeID, level int) doc {
	var name string
	switch f.Node(id).Kind {
	case s.KindSequenceEpsilon:
		return f.Group(id, f.T("epsilon"))
	case s.KindSequenceItem:
		return f.Group(id, f.call("item", []doc{f.kindExpression(f.Child(id, s.RoleElement), 0, level)}))
	case s.KindSequenceConcat:
		name = "concat"
	case s.KindSequenceAlternative:
		name = "alternative"
	case s.KindSequenceStar:
		name = "star"
	default:
		return f.bad(id, "expected child sequence pattern")
	}
	var children []doc
	for _, child := range f.Children(id, s.RoleElement) {
		children = append(children, f.sequencePattern(child, level))
	}
	return f.Group(id, f.call(name, children))
}

func (f *formatter) recursiveKind(id ast.NodeID, level int) doc {
	equations := f.Children(id, s.RoleEquation)
	used, supplied := map[string]bool{}, map[string]bool{}
	for _, parameter := range f.kindParameters {
		used[parameter] = true
	}
	for _, group := range f.groups {
		for _, name := range group {
			used[name] = true
		}
	}
	for _, equation := range equations {
		if name := f.Attr(equation, s.AttributeName); name != "" {
			supplied[name] = true
		}
	}
	names := make([]string, len(equations))
	next := 0
	for i, equation := range equations {
		name := f.Attr(equation, s.AttributeName)
		if name == "" || used[name] {
			for {
				name = "K" + strconv.Itoa(next)
				next++
				if !used[name] && !supplied[name] {
					break
				}
			}
		}
		names[i], used[name] = name, true
	}
	f.groups = append(f.groups, names)
	defer func() { f.groups = f.groups[:len(f.groups)-1] }()
	parts := []doc{f.T("nu"), f.Space(), f.T("{")}
	for i, equation := range equations {
		parts = append(parts, f.Line("    ", level+1), f.Group(equation, f.T(lit.Name(names[i])), f.Space(), f.T("="), f.Space(), f.kindExpression(f.Child(equation, s.RoleBody), 0, level+1), f.T(";")))
	}
	entry, err := strconv.Atoi(f.Attr(id, s.AttributeIndex))
	if err != nil || entry < 0 || entry >= len(names) {
		return f.bad(id, "recursive kind entry has no equation")
	}
	return f.Seq(append(parts, f.Line("    ", level), f.T("}"), f.Space(), f.T("in"), f.Space(), f.T(lit.Name(names[entry])))...)
}
