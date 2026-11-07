// Package pair is a consumer-owned fixture. Its two native arguments occur in
// the reverse order of the Horos tuple's operands.
package pair

import (
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	"github.com/Bitspark/logos-contract/contract/go/typed"
)

const Law = "example/pair-reversed-1"

func Rule(helper n.Type) v.Rule {
	return v.Rule{Key: semantics.Key{Owner: "example", Name: "tuple-pair"}, Dependency: Law, Law: ReversePremises{}, Project: func(p v.Subject) ([]v.Subject, bool, error) {
		if p.Optional() || p.Source().Kind() != decl.Descriptor || p.Source().Constructor() != "tuple" || len(p.Source().Operands()) != 2 || !helper.Same(p.Target().Origin()) {
			return nil, false, nil
		}
		params, args := helper.Parameters(), p.Target().Arguments()
		fields, ok := helper.Fields()
		if !ok || len(params) != 2 || len(args) != 2 || !params[0].Unconstrained() || !params[1].Unconstrained() || len(fields) != 2 {
			return nil, false, nil
		}
		for i, name := range []string{"First", "Second"} {
			f := fields[i]
			if f.Name != name || f.Embedded || f.Tag != "" || !f.Type.Same(params[1-i]) {
				return nil, false, nil
			}
		}
		left, err := p.Operand(0, args[1])
		if err != nil {
			return nil, false, err
		}
		right, err := p.Operand(1, args[0])
		if err != nil {
			return nil, false, err
		}
		return []v.Subject{left, right}, true, nil
	}}
}

// ReversePremises consumes the right witness first, matching the native helper's
// argument order. Logos will reject a candidate applying the two in source order.
type ReversePremises struct{}

func apply(f typed.Expr, arguments ...typed.Expr) typed.Expr {
	for _, argument := range arguments {
		f = typed.App(f, argument)
	}
	return f
}
func relation(h, n typed.Expr) typed.Expr { return apply(typed.Const("Rep"), h, n) }
func (ReversePremises) Signature(c v.Constructors) typed.Expr {
	if c.Arity != 2 {
		return typed.Expr{}
	}
	// Under h0,h1,n0,n1 the right premise uses h1,n1. Under that premise
	// the left uses h0,n0. Both proof binders are in scope at the conclusion.
	right := relation(typed.Var(2), typed.Var(0))
	left := relation(typed.Var(4), typed.Var(2))
	result := relation(apply(typed.Const(c.Source), typed.Var(5), typed.Var(4)),
		apply(typed.Const(c.Target), typed.Var(3), typed.Var(2)))
	body := typed.Pi(right, typed.Pi(left, result))
	for _, domain := range []string{"N", "N", "H", "H"} {
		body = typed.Pi(typed.Const(domain), body)
	}
	return body
}
func (ReversePremises) Apply(law typed.Expr, children []v.Evidence) typed.Expr {
	if len(children) != 2 {
		return typed.Expr{}
	}
	left, right := children[0], children[1]
	return apply(law, left.Source, right.Source, left.Target, right.Target, right.Proof, left.Proof)
}
