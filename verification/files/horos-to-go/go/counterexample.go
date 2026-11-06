package verification

import (
	nativequery "github.com/Bitspark/bittype/analyze/files/golang/go"
	"go/types"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	association "github.com/Bitspark/bittype/contracts/verification/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	h "github.com/Bitspark/bittype/verification/files/source/horos/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	"github.com/Bitspark/logos-contract/contract/go/typed"
)

type counterexample struct {
	source    h.Declaration
	target    n.Type
	nonmember h.Judgment
	admission nativequery.Admission
}

// Counterexample fixes the original source parameter contract, the native
// declaration and an explicitly checked bool witness. The mapping relates only
// this scalar pair; it makes no assertion about the declaration's fields.
func Counterexample(source h.Declaration, target n.Type, nonmember h.Judgment, admission nativequery.Admission) (Claim, error) {
	if !source.Valid() || !target.IsNamed() || !source.Model().Same(nonmember.Model()) || !target.Model().Same(admission.Model()) || !admission.Constructor().Same(target) {
		return Claim{}, ErrInvalid
	}
	sp, np := source.Definition().Parameters(), target.Parameters()
	if len(sp) != 1 || len(np) != 1 || nonmember.Kind().Module() != source.Definition().ID().Module || !sameKind(sp[0].Kind, nonmember.Kind().Expression()) {
		return Claim{}, ErrCoverage
	}
	args := admission.Arguments()
	if len(args) != 1 || !args[0].IsBasic(types.Bool) || scalar(nonmember.Value().Expression()) != "bool" {
		return Claim{}, ErrCoverage
	}
	if nonmember.Verdict() != decl.Refuted || !admission.Admitted() {
		return Claim{}, semantic.ErrEvidenceUnavailable
	}
	d := domains()
	d = append(d, assume("h.bool", con("H"), ModelBasis), assume("n.bool", con("N"), ModelBasis), assume("D", universe(), ModelBasis), assume("native.declaration", con("D"), ModelBasis), assume("Admits", arrow(con("D"), arrow(con("N"), universe())), ModelBasis))
	empty := fixed(typed.ID(typed.Bool(), typed.True(), typed.False()))
	d = append(d, define("Empty", universe(), empty), assume("source.nonmember", arrow(member(con("h.bool")), con("Empty")), SourceBasis), assume("native.admits", app(con("Admits"), con("native.declaration"), con("n.bool")), NativeBasis), assume("law.bool", rep(con("h.bool"), con("n.bool")), ScalarLaw))
	expected := counterexampleType()
	input, req := associationInput(source.Model(), target.Model())
	sites := []association.Site{site(sourceKey(), source.ParameterSite(0)), site(targetKey(), np[0].Site())}
	return Claim{&claim{input: input, requirements: req, question: question("admission-counterexample"), statement: typed.Statement{Foundation: typed.Foundation, Declarations: d, Expected: expected(0)}, origins: logicalOrigins(d, sites), counterexample: &counterexample{source, target, nonmember, admission}}}, nil
}
func counterexampleBody(h, n expression) expression {
	return conjunction(app(con("Admits"), con("native.declaration"), n), conjunction(rep(h, n), arrow(member(h), con("Empty"))))
}
func counterexampleType() expression {
	return sigma(con("H"), func(h expression) expression {
		return sigma(con("N"), func(n expression) expression { return counterexampleBody(h, n) })
	})
}
func CounterexampleCandidate() typed.Expr {
	h, n := con("h.bool"), con("n.bool")
	inner := conjunction(rep(h, n), arrow(member(h), con("Empty")))
	proof := pair(counterexampleBody(h, n), con("native.admits"), pair(inner, con("law.bool"), con("source.nonmember")))
	second := sigma(con("N"), func(n expression) expression { return counterexampleBody(h, n) })
	return pair(counterexampleType(), h, pair(second, n, proof))(0)
}
