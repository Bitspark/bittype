package parse

import (
	"strconv"

	"github.com/Bitspark/bittype/contracts/text/go"
	tr "github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func (p *parser) properSort(start int) ast.NodeID {
	return p.AddSpan(s.KindProperSort, start, start, nil)
}

func (p *parser) properTop(start int) ast.NodeID {
	return p.AddSpan(s.KindTop, start, start, nil, tr.Child(s.RoleSort, p.properSort(start)))
}

// A comma at the outer parenthesis level distinguishes a nonempty n-ary
// domain from ordinary grouping. Nested applications, patterns and nu groups
// retain their own delimiters. This lookahead neither binds nor resolves names.
func (p *parser) multipleDomain() bool {
	if p.Peek() != "(" {
		return false
	}
	parentheses, angles, braces := 0, 0, 0
	for i := p.Pos; i < len(p.Tokens); i++ {
		switch p.Tokens[i].Text {
		case "(":
			parentheses++
		case ")":
			parentheses--
			if parentheses == 0 {
				return false
			}
		case "<":
			angles++
		case ">":
			angles--
		case "{":
			braces++
		case "}":
			braces--
		case ",":
			if parentheses == 1 && angles == 0 && braces == 0 {
				return true
			}
		}
	}
	return false
}

func (p *parser) grouped(id ast.NodeID, start int) ast.NodeID {
	if p.Err == nil && id != 0 && p.Pos > 0 {
		p.Tree.Nodes[id-1].Origin.Spans = []ast.Span{{Source: p.Input.SourceID, Start: uint64(start), End: uint64(p.Tokens[p.Pos-1].End)}}
	}
	return id
}

func (p *parser) sortExpression() ast.NodeID {
	if !p.Enter() {
		return 0
	}
	defer p.Leave()
	start := p.Token().Start
	var inputs []ast.NodeID
	if p.multipleDomain() {
		p.Want("(")
		inputs = append(inputs, p.sortExpression())
		for p.Err == nil && p.Eat(",") {
			inputs = append(inputs, p.sortExpression())
		}
		p.Want(")")
		p.Want("->")
	} else {
		var input ast.NodeID
		if p.Eat("(") {
			input = p.sortExpression()
			p.Want(")")
			p.grouped(input, start)
		} else {
			p.Want("P")
			input = p.Add(s.KindProperSort, start, nil)
		}
		if !p.Eat("->") {
			return input
		}
		inputs = append(inputs, input)
	}
	children := tr.Children(s.RoleInput, inputs)
	children = append(children, tr.Child(s.RoleResult, p.sortExpression()))
	return p.Add(s.KindArrowSort, start, nil, children...)
}

func (p *parser) kindExpression() ast.NodeID {
	if !p.Enter() {
		return 0
	}
	defer p.Leave()
	start := p.Token().Start
	var inputs []ast.NodeID
	if p.multipleDomain() {
		p.Want("(")
		inputs = append(inputs, p.kindExpression())
		for p.Err == nil && p.Eat(",") {
			inputs = append(inputs, p.kindExpression())
		}
		p.Want(")")
		p.Want("->")
	} else {
		input := p.kindUnion()
		if !p.Eat("->") {
			return input
		}
		inputs = append(inputs, input)
	}
	children := tr.Children(s.RoleInput, inputs)
	children = append(children, tr.Child(s.RoleResult, p.kindExpression()))
	return p.Add(s.KindArrow, start, nil, children...)
}

func (p *parser) kindBoolean(kind s.Kind, start int, operands []ast.NodeID) ast.NodeID {
	if len(operands) == 1 {
		return operands[0]
	}
	children := []ast.Child{tr.Child(s.RoleSort, p.properSort(start))}
	children = append(children, tr.Children(s.RoleOperand, operands)...)
	return p.Add(kind, start, nil, children...)
}

func (p *parser) kindUnion() ast.NodeID {
	start := p.Token().Start
	operands := []ast.NodeID{p.kindIntersection()}
	for p.Err == nil && p.Eat("|") {
		operands = append(operands, p.kindIntersection())
	}
	return p.kindBoolean(s.KindUnion, start, operands)
}

func (p *parser) kindIntersection() ast.NodeID {
	start := p.Token().Start
	operands := []ast.NodeID{p.kindPrimary()}
	for p.Err == nil && p.Eat("&") {
		operands = append(operands, p.kindPrimary())
	}
	return p.kindBoolean(s.KindIntersection, start, operands)
}

func (p *parser) kindPrimary() ast.NodeID {
	if p.Err != nil {
		return 0
	}
	start := p.Token().Start
	if p.Eat("(") {
		id := p.kindExpression()
		p.Want(")")
		if p.Err == nil && p.Eat("at") {
			kind := p.Tree.Nodes[id-1].Kind
			if kind != s.KindUnion && kind != s.KindIntersection {
				p.Fail(text.SyntaxError, "sort ascription requires a union or intersection")
				return 0
			}
			if p.ascribed[id] {
				p.Fail(text.SyntaxError, "Boolean kind already has an explicit sort")
				return 0
			}
			p.Want("(")
			sort := p.sortExpression()
			p.Want(")")
			if p.Err == nil {
				p.Tree.Nodes[id-1].Children[0] = tr.Child(s.RoleSort, sort)
				if p.ascribed == nil {
					p.ascribed = map[ast.NodeID]bool{}
				}
				p.ascribed[id] = true
			}
		}
		return p.grouped(id, start)
	}
	switch p.Peek() {
	case "top", "bottom":
		word := p.Take().Text
		p.Want("(")
		sort := p.sortExpression()
		p.Want(")")
		kind := s.KindTop
		if word == "bottom" {
			kind = s.KindBottom
		}
		return p.Add(kind, start, nil, tr.Child(s.RoleSort, sort))
	case "node":
		return p.nodePattern()
	case "nu":
		return p.recursiveKind()
	}
	local := p.Eat("self")
	if local {
		p.Want(".")
	}
	name := p.name()
	attrs := tr.Attrs(s.AttributeScope, "local", s.AttributeName, name)
	var callee ast.NodeID
	if !local && p.Eat(".") {
		attrs = tr.Attrs(s.AttributeScope, "imported", s.AttributeAlias, name, s.AttributeName, p.name())
	} else if !local {
		for i := len(p.groups) - 1; i >= 0; i-- {
			if index, ok := p.groups[i][name]; ok {
				callee = p.Add(s.KindRecursiveVariable, start, tr.Attrs(s.AttributeDepth, strconv.Itoa(len(p.groups)-1-i), s.AttributeIndex, strconv.Itoa(index)))
				break
			}
		}
		if callee == 0 {
			if index, ok := p.kindParameters[name]; ok {
				callee = p.Add(s.KindKindParameterReference, start, tr.Attrs(s.AttributeIndex, strconv.Itoa(index)))
			}
		}
	}
	if callee == 0 {
		callee = p.Add(s.KindKindReference, start, attrs)
	}
	if p.Err != nil || !p.Eat("<") {
		return callee
	}
	if p.Tree.Nodes[callee-1].Kind != s.KindKindReference {
		p.Fail(text.SyntaxError, "kind applications require a named kind definition")
		return 0
	}
	children := []ast.Child{tr.Child(s.RoleCallee, callee)}
	for p.Err == nil && !p.EOF() && p.Peek() != ">" {
		children = append(children, tr.Child(s.RoleArgument, p.kindExpression()))
		if !p.Eat(",") {
			break
		}
	}
	if len(children) == 1 {
		p.Fail(text.SyntaxError, "empty kind application argument list; use a reference")
	}
	p.Want(">")
	return p.Add(s.KindKindApplication, start, nil, children...)
}

func (p *parser) nodePattern() ast.NodeID {
	start := p.Take().Start
	p.Want("(")
	form := p.quoted()
	p.Want(",")
	p.Want("{")
	var children []ast.Child
	seen := map[string]bool{}
	for p.Err == nil && !p.EOF() && p.Peek() != "}" {
		begin := p.Token().Start
		key := p.quoted()
		if seen[key] {
			p.Fail(text.SyntaxError, "duplicate node selector")
		}
		seen[key] = true
		p.Want(":")
		value := p.quoted()
		children = append(children, tr.Child(s.RoleSelector, p.Add(s.KindSelector, begin, tr.Attrs(s.AttributeName, key, s.AttributeValue, value))))
		if !p.Eat(",") {
			break
		}
	}
	p.Want("}")
	p.Want(",")
	children = append(children, tr.Child(s.RolePattern, p.sequencePattern()))
	p.Want(")")
	return p.Add(s.KindNodePattern, start, tr.Attrs(s.AttributeConstructor, form), children...)
}

func (p *parser) sequencePattern() ast.NodeID {
	if !p.Enter() {
		return 0
	}
	defer p.Leave()
	start := p.Token().Start
	word := p.Take().Text
	if word == "epsilon" {
		return p.Add(s.KindSequenceEpsilon, start, nil)
	}
	var kind s.Kind
	switch word {
	case "item":
		kind = s.KindSequenceItem
	case "concat":
		kind = s.KindSequenceConcat
	case "alternative":
		kind = s.KindSequenceAlternative
	case "star":
		kind = s.KindSequenceStar
	default:
		p.Fail(text.SyntaxError, "expected a child sequence pattern")
		return 0
	}
	p.Want("(")
	var children []ast.Child
	if kind == s.KindSequenceItem {
		children = append(children, tr.Child(s.RoleElement, p.kindExpression()))
	} else {
		children = append(children, tr.Child(s.RoleElement, p.sequencePattern()))
		for kind != s.KindSequenceStar && p.Err == nil && p.Eat(",") {
			children = append(children, tr.Child(s.RoleElement, p.sequencePattern()))
		}
		if kind != s.KindSequenceStar && len(children) < 2 {
			p.Fail(text.SyntaxError, "sequence concatenation and alternative require at least two patterns")
		}
	}
	p.Want(")")
	return p.Add(kind, start, nil, children...)
}

// Collect the group's binders before parsing any body, so forward references
// and mutual recursion use the same lexical slots as backward references.
func (p *parser) equationNames() map[string]int {
	names := map[string]int{}
	braces, boundary := 0, true
	for i := p.Pos; i < len(p.Tokens); i++ {
		token := p.Tokens[i]
		if token.Class == "eof" {
			break
		}
		if braces == 0 && token.Text == "}" {
			return names
		}
		if boundary {
			name, valid := sourceName(token.Text)
			if !valid || i+1 >= len(p.Tokens) || p.Tokens[i+1].Text != "=" {
				p.Pos = i
				p.Fail(text.SyntaxError, "expected a named kind equation")
				return nil
			}
			if _, exists := names[name]; exists {
				p.Pos = i
				p.Fail(text.SyntaxError, "duplicate recursive equation name")
				return nil
			}
			names[name] = len(names)
			boundary = false
		}
		switch token.Text {
		case "{":
			braces++
		case "}":
			braces--
		case ";":
			if braces == 0 {
				boundary = true
			}
		}
	}
	p.Fail(text.SyntaxError, "unterminated recursive kind group")
	return nil
}

func (p *parser) recursiveKind() ast.NodeID {
	start := p.Take().Start
	p.Want("{")
	names := p.equationNames()
	if p.Err != nil {
		return 0
	}
	if len(names) == 0 {
		p.Fail(text.SyntaxError, "recursive kinds require at least one equation")
		return 0
	}
	p.groups = append(p.groups, names)
	defer func() { p.groups = p.groups[:len(p.groups)-1] }()
	var children []ast.Child
	for p.Err == nil && !p.EOF() && p.Peek() != "}" {
		begin := p.Token().Start
		name := p.name()
		p.Want("=")
		body := p.kindExpression()
		p.Want(";")
		equation := p.Add(s.KindEquation, begin, tr.Attrs(s.AttributeName, name, s.AttributeIndex, strconv.Itoa(len(children))), tr.Child(s.RoleBody, body))
		children = append(children, tr.Child(s.RoleEquation, equation))
	}
	p.Want("}")
	p.Want("in")
	entry, ok := names[p.name()]
	if !ok {
		p.Fail(text.SyntaxError, "recursive kind entry must name an equation in this group")
	}
	return p.Add(s.KindRecursiveKind, start, tr.Attrs(s.AttributeIndex, strconv.Itoa(entry)), children...)
}
