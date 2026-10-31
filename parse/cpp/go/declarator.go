package cppparse

import (
	"github.com/Bitspark/bittype/ast/go"
	tr "github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	"github.com/Bitspark/bittype/text/go"
)

type qualifier struct {
	id    ast.NodeID
	attrs []ast.Attribute
}
type declarator struct {
	name      string
	apply     func(ast.NodeID) ast.NodeID
	qualified []*qualifier
}

func (p *parser) memberPointerAhead() bool {
	depth := 0
	for i := p.Pos; i < len(p.Tokens); i++ {
		t := p.Tokens[i].Text
		if t == "<" {
			depth++
		}
		if t == ">" {
			if depth == 0 {
				return false
			}
			depth--
		}
		if depth == 0 && t == "::" && i+1 < len(p.Tokens) && p.Tokens[i+1].Text == "*" {
			return true
		}
		if depth == 0 && (t == "(" || t == ")" || t == "[" || t == "]" || t == ";" || t == "," || t == "=" || t == "*" || t == "&" || t == "&&" || t == "") {
			return false
		}
	}
	return false
}
func (p *parser) declarator(named bool) declarator {
	identity := func(id ast.NodeID) ast.NodeID { return id }
	result := declarator{apply: identity}
	if !p.Enter() {
		return result
	}
	defer p.Leave()
	start := p.Token().Start
	var prefixes []func(ast.NodeID) ast.NodeID
	for p.Err == nil {
		kind := s.KindPointerType
		var owner ast.NodeID
		if p.Eat("*") {
		} else if p.Eat("&") {
			kind = s.KindLValueReferenceType
		} else if p.Eat("&&") {
			kind = s.KindRValueReferenceType
		} else if p.memberPointerAhead() {
			owner = p.named()
			p.Want("::")
			p.Want("*")
			kind = s.KindMemberPointerType
		} else {
			break
		}
		attrs := p.cv()
		end := p.Tokens[p.Pos-1].End
		prefixes = append(prefixes, func(base ast.NodeID) ast.NodeID {
			children := []ast.Child{}
			if owner != 0 {
				children = append(children, tr.Child("owner", owner))
			}
			children = append(children, tr.Child("element", base))
			id := p.AddSpan(kind, start, end, nil, children...)
			if len(attrs) > 0 {
				id = p.AddSpan(s.KindCVQualifiedType, start, end, attrs, tr.Child("element", id))
			}
			return id
		})
	}
	if p.Peek() == "(" && (p.At(1) == "*" || p.At(1) == "&" || p.At(1) == "&&" || p.At(1) == "(" || (named && p.At(2) == ")") || p.groupedMemberPointer()) {
		p.Take()
		result = p.declarator(named)
		p.Want(")")
	} else if named && p.Token().Class == "identifier" {
		result.name = p.Ident()
	}
	for p.Err == nil {
		if p.Eat("[") {
			length := p.Take().Text
			p.Want("]")
			old := result.apply
			end := p.Tokens[p.Pos-1].End
			result.apply = func(base ast.NodeID) ast.NodeID {
				return old(p.AddSpan(s.KindArrayType, start, end, tr.Attrs("length", length), tr.Child("element", base)))
			}
		} else if p.Eat("(") {
			var params []ast.Child
			variadic := false
			for p.Err == nil && !p.EOF() && p.Peek() != ")" {
				if p.Eat("...") {
					variadic = true
					break
				}
				begin := p.Token().Start
				base := p.base()
				d := p.declarator(true)
				value := d.apply(base)
				if len(d.qualified) > 0 {
					p.Fail(text.Unsupported, "qualified function parameter type")
					break
				}
				var attrs []ast.Attribute
				if d.name != "" {
					attrs = tr.Attrs("name", d.name)
				}
				params = append(params, tr.Child("parameter", p.Add(s.KindParameter, begin, attrs, tr.Child("type", value))))
				if !p.Eat(",") {
					break
				}
				if p.Peek() == ")" {
					p.Fail(text.SyntaxError, "trailing comma in function parameters")
				}
			}
			p.Want(")")
			if p.Err == nil && len(params) == 1 {
				param := p.Tree.Nodes[params[0].Node-1]
				typ := p.Tree.Nodes[param.Children[0].Node-1]
				if len(param.Attributes) == 0 && typ.Kind == s.KindBuiltinType && len(typ.Attributes) == 1 && typ.Attributes[0].Text == "void" {
					params = nil
				}
			}
			q := &qualifier{attrs: p.cv()}
			if p.Eat("&") {
				q.attrs = append(q.attrs, tr.Attrs("ref-qualifier", "lvalue")...)
			} else if p.Eat("&&") {
				q.attrs = append(q.attrs, tr.Attrs("ref-qualifier", "rvalue")...)
			}
			var attrs []ast.Attribute
			if variadic {
				attrs = tr.Attrs("variadic", "true")
			}
			if p.Eat("noexcept") {
				attrs = append(attrs, tr.Attrs("noexcept", "true")...)
			}
			if len(q.attrs) > 0 {
				result.qualified = append(result.qualified, q)
			}
			end := p.Tokens[p.Pos-1].End
			old := result.apply
			result.apply = func(base ast.NodeID) ast.NodeID {
				children := append(append([]ast.Child{}, params...), tr.Child("result", base))
				q.id = p.AddSpan(s.KindFunctionType, start, end, attrs, children...)
				return old(q.id)
			}
		} else {
			break
		}
	}
	direct := result.apply
	result.apply = func(base ast.NodeID) ast.NodeID {
		if p.Err != nil {
			return 0
		}
		for _, prefix := range prefixes {
			base = prefix(base)
		}
		return direct(base)
	}
	return result
}
func (p *parser) groupedMemberPointer() bool {
	if p.Peek() != "(" {
		return false
	}
	p.Pos++
	ok := p.memberPointerAhead()
	p.Pos--
	return ok
}
