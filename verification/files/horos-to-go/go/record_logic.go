package verification

import "github.com/Bitspark/logos-contract/contract/go/typed"

func domainsN(domain expression, count int, result expression) expression {
	for range count {
		result = arrow(domain, result)
	}
	return result
}
func quantify(domain expression, count int, body func([]expression) expression) expression {
	var step func([]expression) expression
	step = func(xs []expression) expression {
		if len(xs) == count {
			return body(xs)
		}
		return pi(domain, func(x expression) expression { return step(append(append([]expression(nil), xs...), x)) })
	}
	return step(nil)
}
func (p pattern) terms(h, n expression) (expression, expression) {
	switch p.shape {
	case "parameter":
		return h, n
	case "text", "bool":
		return con("h." + p.shape), con("n." + p.shape)
	case "list":
		a, b := p.element.terms(h, n)
		return app(con("h.list"), a), app(con("n.slice"), b)
	case "optional":
		a, b := p.element.terms(h, n)
		return app(con("h.optional"), a), app(con("n.optional"), b)
	default:
		panic("unadmitted private mapping pattern")
	}
}
func (p pattern) proof(h, n, witness expression) expression {
	switch p.shape {
	case "parameter":
		return witness
	case "text", "bool":
		return con("law." + p.shape)
	case "list", "optional":
		a, b := p.element.terms(h, n)
		return app(con("law."+p.shape), a, b, p.element.proof(h, n, witness))
	default:
		panic("unadmitted private mapping pattern")
	}
}
func recordBasis(prefix string, r record) []typed.Declaration {
	count := len(r.patterns)
	d := []typed.Declaration{assume(prefix+".h.record", domainsN(con("H"), count, con("H"))), assume(prefix+".n.record", domainsN(con("N"), count, con("N")))}
	law := quantify(con("H"), count, func(hs []expression) expression {
		return quantify(con("N"), count, func(ns []expression) expression {
			out := rep(app(con(prefix+".h.record"), hs...), app(con(prefix+".n.record"), ns...))
			for i := count - 1; i >= 0; i-- {
				out = arrow(rep(hs[i], ns[i]), out)
			}
			return out
		})
	})
	return append(d, assume(prefix+".law.record", law, RecordLaw))
}
func bodies(prefix string, r record, h, n expression) (expression, expression) {
	var hs, ns []expression
	for _, p := range r.patterns {
		a, b := p.terms(h, n)
		hs = append(hs, a)
		ns = append(ns, b)
	}
	return app(con(prefix+".h.record"), hs...), app(con(prefix+".n.record"), ns...)
}
func bodyProof(prefix string, r record, h, n, witness expression) expression {
	var hs, ns, ps []expression
	for _, p := range r.patterns {
		a, b := p.terms(h, n)
		hs = append(hs, a)
		ns = append(ns, b)
		ps = append(ps, p.proof(h, n, witness))
	}
	return app(con(prefix+".law.record"), append(append(hs, ns...), ps...)...)
}
func namedProof(prefix string, r record, h, n, witness, hName, nName, sourceFact, nativeFact expression) expression {
	a, b := bodies(prefix, r, h, n)
	return app(con("law.named"), hName, a, nName, b, sourceFact, nativeFact, bodyProof(prefix, r, h, n, witness))
}
func familyResult(h, n expression) expression {
	return conjunction(member(app(con("h.family"), h)), rep(app(con("h.family"), h), app(con("n.family"), n)))
}
func familyLogic(r record) ([]typed.Declaration, typed.Expr, typed.Expr) {
	d := mappingBasis()
	d = append(d, recordBasis("family", r)...)
	d = append(d, assume("h.family", arrow(con("H"), con("H"))), assume("n.family", arrow(con("N"), con("N"))))
	source := pi(con("H"), func(h expression) expression {
		a, _ := bodies("family", r, h, con("n.text"))
		return arrow(member(h), sourceBody(app(con("h.family"), h), a))
	})
	native := pi(con("N"), func(n expression) expression {
		_, b := bodies("family", r, con("h.text"), n)
		return nativeBody(app(con("n.family"), n), b)
	})
	guarantee := pi(con("H"), func(h expression) expression { return arrow(member(h), member(app(con("h.family"), h))) })
	d = append(d, assume("source.family.body", source, SourceBasis), assume("native.family.body", native, NativeBasis), assume("source.family.kind", guarantee, SourceBasis))
	expected := pi(con("H"), func(h expression) expression {
		return pi(con("N"), func(n expression) expression { return arrow(member(h), arrow(rep(h, n), familyResult(h, n))) })
	})
	candidate := lambda(con("H"), func(h expression) expression {
		return lambda(con("N"), func(n expression) expression {
			return lambda(member(h), func(m expression) expression {
				return lambda(rep(h, n), func(w expression) expression {
					proof := namedProof("family", r, h, n, w, app(con("h.family"), h), app(con("n.family"), n), app(con("source.family.body"), h, m), app(con("native.family.body"), n))
					return pair(familyResult(h, n), app(con("source.family.kind"), h, m), proof)
				})
			})
		})
	})
	return d, expected(0), candidate(0)
}
func argumentLogic(r record) ([]typed.Declaration, typed.Expr) {
	d := recordBasis("argument", r)
	d = append(d, assume("h.argument", con("H")), assume("n.argument", con("N")))
	h, n := con("h.text"), con("n.text") // no binders occur in this admitted nullary body
	a, b := bodies("argument", r, h, n)
	d = append(d, assume("source.argument.body", sourceBody(con("h.argument"), a), SourceBasis), assume("native.argument.body", nativeBody(con("n.argument"), b), NativeBasis), assume("argument.member", member(con("h.argument")), SourceBasis))
	proof := namedProof("argument", r, h, n, con("law.text"), con("h.argument"), con("n.argument"), con("source.argument.body"), con("native.argument.body"))
	return d, proof(0)
}
