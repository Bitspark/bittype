// Package verification composes static representation laws over captured Horos and Go subjects.
package verification

import (
	"errors"
	"fmt"
	association "github.com/Bitspark/bittype/contracts/verification/go"
	"github.com/Bitspark/bittype/semantics/go"
	joint "github.com/Bitspark/bittype/verification/logos/go"
	"github.com/Bitspark/logos-contract/contract/go/typed"
	"slices"
)

var ErrCoverage = errors.New("no selected representation covers the captured pair")
var ErrAmbiguous = errors.New("several selected representations cover the same pair")

// Rule selects actual child pairs and contributes a law under an explicit
// dependency. A successful projection grants no logical acceptance.
type Rule struct {
	Key        semantics.Key
	Dependency string
	Law Law
	Project    func(Subject) (children []Subject, applies bool, err error)
}
type Mappings struct{ rules []Rule }

func Select(rules ...Rule) (Mappings, error) {
	if len(rules) > 128 {
		return Mappings{}, association.ErrLimit
	}
	keys := make(map[semantics.Key]bool, len(rules))
	for _, r := range rules {
		if !r.Key.Valid() || keys[r.Key] || r.Dependency == "" || r.Project == nil || r.Law == nil {
			return Mappings{}, joint.ErrInvalid
		}
		if len(r.Key.Owner)+len(r.Key.Name)+len(r.Dependency) > 4096 {
			return Mappings{}, association.ErrLimit
		}
		keys[r.Key] = true
	}
	return Mappings{slices.Clone(rules)}, nil
}

type node struct {
 law Law
	symbol, dependency       string
	children                 []*node
	sourceSites, nativeSites []association.Site
	bound                    bool
}

func (m Mappings) project(root Subject, prefix string) (*node, error) {
	count := 0
	var visit func(Subject, int) (*node, error)
	visit = func(p Subject, depth int) (*node, error) {
		count++
		if count > 512 || depth > 32 {
			return nil, association.ErrLimit
		}
		if p.s == nil || p.Source() == nil || !p.Target().Valid() {
			return nil, joint.ErrInvalid
		}
		result := &node{sourceSites: []association.Site{site(sourceKey(), p.s.sourceSite)}, nativeSites: []association.Site{site(targetKey(), p.s.nativeSite)}}
		if loc := p.Target().Origin().Site(); loc.Node != 0 && loc != p.s.nativeSite {
			result.nativeSites = append(result.nativeSites, site(targetKey(), loc))
		}
		if id, ok := p.Source().Parameter(); ok {
			if id != p.s.binder.source || !p.s.binder.target.Valid() {
				return nil, ErrCoverage
			}
			if !p.Optional() && p.Target().Same(p.s.binder.target){result.bound=true;return result,nil}
		}
		var children []Subject
		for _, rule := range m.rules {
			inputs, applies, err := rule.Project(p)
			if err != nil {
				return nil, err
			}
			if !applies {
				continue
			}
			if result.symbol != "" {
				return nil, ErrAmbiguous
			}
			if len(inputs) > 64 {
				return nil, association.ErrLimit
			}
			// A position identifies the captured occurrence. Length framing keeps
			// distinct owner/name pairs from aliasing in a Logos statement.
			result.symbol = fmt.Sprintf("%s.%d.%d.%s.%s", prefix, count, len(rule.Key.Owner), rule.Key.Owner, rule.Key.Name)
			result.dependency = rule.Dependency
 result.law=rule.Law
			children = slices.Clone(inputs)
		}
		if result.symbol == "" {
			return nil, ErrCoverage
		}
		for _, child := range children {
			if child.s == nil || child.s.parent != p.s {
				return nil, joint.ErrInvalid
			}
			projected, err := visit(child, depth+1)
			if err != nil {
				return nil, err
			}
			result.children = append(result.children, projected)
		}
		return result, nil
	}
	return visit(root, 0)
}

func (x *node) sites() []association.Site {
	return append(slices.Clone(x.sourceSites), x.nativeSites...)
}
func (x *node) describe(s *signature) {
	if x.bound {
		return
	}
	for _, c := range x.children {
		c.describe(s)
	}
	count := len(x.children)
	hs, ns, law := x.symbol+".source", x.symbol+".target", x.symbol+".law"
	s.add(assume(hs, domainsN(con("H"), count, con("H")), ModelBasis), x.sourceSites...)
	s.add(assume(ns, domainsN(con("N"), count, con("N")), ModelBasis), x.nativeSites...)
 spec:=Constructors{Source:hs,Target:ns,Arity:count}
 s.add(typed.Assume(law,x.law.Signature(spec),x.dependency),x.sites()...)
}
func (x *node) terms(h, n expression) (expression, expression) {
	if x.bound {
		return h, n
	}
	var hs, ns []expression
	for _, c := range x.children {
		a, b := c.terms(h, n)
		hs = append(hs, a)
		ns = append(ns, b)
	}
	return app(con(x.symbol+".source"), hs...), app(con(x.symbol+".target"), ns...)
}
func (x *node) proof(h, n, witness expression) expression {
	if x.bound {
		return witness
	}
 // Each child's terms and proof share the same lexical depth and occurrence.
 var sources,targets,proofs []expression
 for _,c:=range x.children{
  a,b:=c.terms(h,n)
  sources=append(sources,a);targets=append(targets,b);proofs=append(proofs,c.proof(h,n,witness))
 }
 return func(depth int)typed.Expr{
  inputs:=make([]Evidence,len(proofs))
  for i:=range inputs{inputs[i]=Evidence{sources[i](depth),targets[i](depth),proofs[i](depth)}}
  return x.law.Apply(typed.Const(x.symbol+".law"),inputs)
 }
}
