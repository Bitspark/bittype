package verification

import "github.com/Bitspark/logos/contract/go/typed"

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

func assume(name string, typ expression, deps ...string) typed.Declaration {
	return typed.Assume(name, typ(0), deps...)
}
func define(name string, typ, value expression) typed.Declaration {
	return typed.Define(name, typ(0), value(0))
}

func domains() []typed.Declaration {
	return []typed.Declaration{
		assume("H", universe()), assume("N", universe()),
		assume("Member", arrow(con("H"), universe())),
		assume("Rep", arrow(con("H"), arrow(con("N"), universe()))),
		assume("SourceBodyOf", arrow(con("H"), arrow(con("H"), universe()))),
		assume("NativeBodyOf", arrow(con("N"), arrow(con("N"), universe()))),
	}
}

// Dependency names are contracts of this selected integration, not ambient trust
// or registrations. The receiver explicitly supplies the policy allowing them.
const (
	SourceBasis = "horos/declarations-4/procedure"
	NativeBasis = "go1.25/captured-types-1"
	ScalarLaw   = "horos-go/static-scalars-1"
	ListLaw     = "horos-go/list-slice-1"
	OptionalLaw = "horos-go/optional-presence-1"
	RecordLaw   = "horos-go/record-fields-1"
	NamedLaw    = "horos-go/named-representation-1"
)

func MappingDependencies() []string {
	return []string{SourceBasis, NativeBasis, ScalarLaw, ListLaw, OptionalLaw, RecordLaw, NamedLaw}
}

func mappingBasis() []typed.Declaration {
	d := domains()
	for _, scalar := range []string{"text", "bool"} {
		d = append(d, assume("h."+scalar, con("H")), assume("n."+scalar, con("N")), assume("law."+scalar, rep(con("h."+scalar), con("n."+scalar)), ScalarLaw))
	}
	for _, shape := range []struct{ source, target, law, dependency string }{{"list", "slice", "list", ListLaw}, {"optional", "optional", "optional", OptionalLaw}} {
		d = append(d, assume("h."+shape.source, arrow(con("H"), con("H"))), assume("n."+shape.target, arrow(con("N"), con("N"))))
		t := pi(con("H"), func(h expression) expression {
			return pi(con("N"), func(n expression) expression {
				return arrow(rep(h, n), rep(app(con("h."+shape.source), h), app(con("n."+shape.target), n)))
			})
		})
		d = append(d, assume("law."+shape.law, t, shape.dependency))
	}
	// Named declarations remain distinct from their structural bodies.
	named := pi(con("H"), func(a expression) expression {
		return pi(con("H"), func(b expression) expression {
			return pi(con("N"), func(x expression) expression {
				return pi(con("N"), func(y expression) expression {
					return arrow(sourceBody(a, b), arrow(nativeBody(x, y), arrow(rep(b, y), rep(a, x))))
				})
			})
		})
	})
	return append(d, assume("law.named", named, NamedLaw))
}
