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
func namedProof(r record,h,n,witness,hName,nName,sourceFact,nativeFact expression)expression{
 a,b:=r.body.terms(h,n)
 return app(con("law.named"),hName,a,nName,b,sourceFact,nativeFact,r.body.proof(h,n,witness))
}
func familyResult(h, n expression) expression {
	return conjunction(member(app(con("h.family"), h)), rep(app(con("h.family"), h), app(con("n.family"), n)))
}
func familyLogic(r record) (signature, typed.Expr, typed.Expr) {
 d:=signature{declarations:domains()}
 d.add(namedDeclaration(),r.body.sites()...)
 r.body.describe(&d)
 d.add(assume("h.family",arrow(con("H"),con("H")),ModelBasis),r.body.sourceSites...)
 d.add(assume("n.family",arrow(con("N"),con("N")),ModelBasis),r.body.nativeSites...)
	source := pi(con("H"), func(h expression) expression {
		a, _ := r.body.terms(h,nil)
		return arrow(member(h), sourceBody(app(con("h.family"), h), a))
	})
	native := pi(con("N"), func(n expression) expression {
		_, b := r.body.terms(nil,n)
		return nativeBody(app(con("n.family"), n), b)
	})
	guarantee := pi(con("H"), func(h expression) expression { return arrow(member(h), member(app(con("h.family"), h))) })
	d.add(assume("source.family.body",source,SourceBasis),r.body.sourceSites...)
 d.add(assume("native.family.body",native,NativeBasis),r.body.nativeSites...)
 d.add(assume("source.family.kind",guarantee,SourceBasis),site(sourceKey(),r.source.ParameterSite(0)))
	expected := pi(con("H"), func(h expression) expression {
		return pi(con("N"), func(n expression) expression { return arrow(member(h), arrow(rep(h, n), familyResult(h, n))) })
	})
	candidate := lambda(con("H"), func(h expression) expression {
		return lambda(con("N"), func(n expression) expression {
			return lambda(member(h), func(m expression) expression {
				return lambda(rep(h, n), func(w expression) expression {
					proof := namedProof(r, h, n, w, app(con("h.family"), h), app(con("n.family"), n), app(con("source.family.body"), h, m), app(con("native.family.body"), n))
					return pair(familyResult(h, n), app(con("source.family.kind"), h, m), proof)
				})
			})
		})
	})
	return d, expected(0), candidate(0)
}
func argumentLogic(r record)(signature,typed.Expr){
 var d signature
 r.body.describe(&d)
 d.add(assume("h.argument",con("H"),ModelBasis),r.body.sourceSites...)
 d.add(assume("n.argument",con("N"),ModelBasis),r.body.nativeSites...)
 // A nullary argument has no binder. Its projected tree is closed.
 a,b:=r.body.terms(nil,nil)
 d.add(assume("source.argument.body",sourceBody(con("h.argument"),a),SourceBasis),r.body.sourceSites...)
 d.add(assume("native.argument.body",nativeBody(con("n.argument"),b),NativeBasis),r.body.nativeSites...)
 d.add(assume("argument.member",member(con("h.argument")),SourceBasis),r.body.sourceSites...)
 proof:=namedProof(r,nil,nil,nil,con("h.argument"),con("n.argument"),con("source.argument.body"),con("native.argument.body"))
 return d,proof(0)
}
