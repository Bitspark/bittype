package verification

import (
	n "github.com/Bitspark/bittype/semantics/golang/go"
	"github.com/Bitspark/logos-contract/contract/go/typed"
)

// Local syntax construction uses lexical closures solely to calculate de Bruijn
// indices. Logos owns all substitution, conversion, inference and checking.
type expression func(int) typed.Expr

func fixed(e typed.Expr) expression { return func(int) typed.Expr { return e } }
func con(s string) expression       { return fixed(typed.Const(s)) }
func app(f expression, args ...expression) expression {
	return func(d int) typed.Expr {
		e := f(d)
		for _, a := range args {
			e = typed.App(e, a(d))
		}
		return e
	}
}
func bind(kind string, domain expression, body func(expression) expression) expression {
	return func(d int) typed.Expr {
		v := func(at int) typed.Expr { return typed.Var(uint32(at - d - 1)) }
		b := body(v)(d + 1)
		switch kind {
		case "pi":
			return typed.Pi(domain(d), b)
		case "lambda":
			return typed.Lam(domain(d), b)
		default:
			return typed.Sigma(domain(d), b)
		}
	}
}
func pi(domain expression, body func(expression) expression) expression {
	return bind("pi", domain, body)
}
func lambda(domain expression, body func(expression) expression) expression {
	return bind("lambda", domain, body)
}
func sigma(domain expression, body func(expression) expression) expression {
	return bind("sigma", domain, body)
}
func pair(typ, a, b expression) expression {
	return func(d int) typed.Expr { return typed.Pair(typ(d), a(d), b(d)) }
}
func fst(x expression) expression      { return func(d int) typed.Expr { return typed.Fst(x(d)) } }
func snd(x expression) expression      { return func(d int) typed.Expr { return typed.Snd(x(d)) } }
func arrow(a, b expression) expression { return pi(a, func(expression) expression { return b }) }
func conjunction(a, b expression) expression {
	return sigma(a, func(expression) expression { return b })
}
func universe() expression                  { return fixed(typed.Sort(0)) }
func rep(h, n expression) expression        { return app(con("Rep"), h, n) }
func member(h expression) expression        { return app(con("Member"), h) }
func sourceBody(a, b expression) expression { return app(con("SourceBodyOf"), a, b) }
func nativeBody(a, b expression) expression { return app(con("NativeBodyOf"), a, b) }

func assume(name string, typ expression, dependency string) typed.Declaration {
	return typed.Assume(name, typ(0), dependency)
}
func define(name string, typ, value expression) typed.Declaration {
	return typed.Define(name, typ(0), value(0))
}

func domains() []typed.Declaration {
	return []typed.Declaration{
		assume("H", universe(), ModelBasis), assume("N", universe(), ModelBasis),
		assume("Member", arrow(con("H"), universe()), ModelBasis),
		assume("Rep", arrow(con("H"), arrow(con("N"), universe())), ModelBasis),
		assume("SourceBodyOf", arrow(con("H"), arrow(con("H"), universe())), ModelBasis),
		assume("NativeBodyOf", arrow(con("N"), arrow(con("N"), universe())), ModelBasis),
	}
}

// Dependency names are contracts of this selected integration, not ambient trust
// or registrations. The receiver explicitly supplies the policy allowing them.
const (
	ModelBasis  = "horos-go/captured-model-1"
	SourceBasis = "horos/declarations-4/procedure"
	NativeBasis = n.Profile
	NamedLaw    = "horos-go/named-representation-1"
)

// The named declaration interpretation consumes independently projected body evidence.
func namedDeclaration() typed.Declaration {
	named := pi(con("H"), func(a expression) expression {
		return pi(con("H"), func(b expression) expression {
			return pi(con("N"), func(x expression) expression {
				return pi(con("N"), func(y expression) expression {
					return arrow(sourceBody(a, b), arrow(nativeBody(x, y), arrow(rep(b, y), rep(a, x))))
				})
			})
		})
	})
	return assume("law.named", named, NamedLaw)
}
