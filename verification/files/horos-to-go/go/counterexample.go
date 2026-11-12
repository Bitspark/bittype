package verification

import (
	nativequery "github.com/Bitspark/bittype/analyze/files/golang/queries/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	h "github.com/Bitspark/bittype/verification/files/source/horos/go"
	joint "github.com/Bitspark/bittype/verification/logos/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	"github.com/Bitspark/logos-contract/contract/go/typed"
)

// Admissions interprets admission counterexamples under one selected mapping set.
type Admissions struct {
	mappings       Mappings
	interpretation joint.Interpretation[*witness]
}

func NewAdmissions(m Mappings) Admissions {
	return Admissions{m, joint.Interpret(question("admission-counterexample"), witnessProblem)}
}
func (l Admissions) Authority() joint.Authority { return l.interpretation.Authority() }

type witness struct {
	nonmember h.Judgment
	admission nativequery.Admission
	source    h.Declaration
	target    n.Type
	pair      *node
}
type Witness struct {
	Claim     joint.Claim
	Candidate typed.Expr
}

func (l Admissions) Counterexample(source h.Declaration, target n.Type, nonmember h.Judgment, admission nativequery.Admission) (Witness, error) {
	if !source.Valid() || !target.IsNamed() || !source.Model().Same(nonmember.Model()) || !target.Model().Same(admission.Model()) || !admission.Constructor().Same(target) {
		return Witness{}, joint.ErrInvalid
	}
	sp, np := source.Definition().Parameters(), target.Parameters()
	if len(sp) != 1 || len(np) != 1 || nonmember.Kind().Module() != source.Definition().ID().Module || !sameKind(sp[0].Kind, nonmember.Kind().Expression()) {
		return Witness{}, ErrCoverage
	}
	args := admission.Arguments()
	if len(args) != 1 {
		return Witness{}, ErrCoverage
	}
	if nonmember.Verdict() != decl.Refuted || !admission.Admitted() {
		return Witness{}, semantic.ErrEvidenceUnavailable
	}
	pair, err := l.mappings.project(Subject{&subject{expression: nonmember.Value().Expression(), native: args[0], sourceSite: source.ParameterSite(0), nativeSite: np[0].Site()}}, "witness")
	if err != nil {
		return Witness{}, err
	}
	w := &witness{nonmember, admission, source, target, pair}
	claim, err := l.interpretation.Ask(w)
	if err != nil {
		return Witness{}, err
	}
	return Witness{claim, w.candidate()}, nil
}
func witnessProblem(w *witness) (joint.Problem, error) {
	d := signature{declarations: domains()}
	w.pair.describe(&d)
	h, n := w.pair.terms(nil, nil)
	d.add(assume("D", universe(), ModelBasis))
	d.add(assume("native.declaration", con("D"), ModelBasis), site(targetKey(), w.target.Site()))
	d.add(assume("Admits", arrow(con("D"), arrow(con("N"), universe())), ModelBasis))
	d.add(define("Empty", universe(), fixed(typed.ID(typed.Bool(), typed.True(), typed.False()))))
	d.add(assume("source.nonmember", arrow(member(h), con("Empty")), SourceBasis), w.pair.sourceSite())
	d.add(assume("native.admits", app(con("Admits"), con("native.declaration"), n), NativeBasis), w.pair.nativeSites()...)
	input, req := associationInput(w.nonmember.Model(), w.admission.Model())
	sites := append(w.pair.sites(), site(sourceKey(), w.source.Site()), site(targetKey(), w.target.Site()))
	return d.problem(input, req, counterexampleType()(0), sites), nil
}
func counterexampleBody(h, n expression) expression {
	return conjunction(app(con("Admits"), con("native.declaration"), n), conjunction(rep(h, n), arrow(member(h), con("Empty"))))
}
func counterexampleType() expression {
	return sigma(con("H"), func(h expression) expression {
		return sigma(con("N"), func(n expression) expression { return counterexampleBody(h, n) })
	})
}
func (w *witness) candidate() typed.Expr {
	h, n := w.pair.terms(nil, nil)
	inner := conjunction(rep(h, n), arrow(member(h), con("Empty")))
	proof := pair(counterexampleBody(h, n), con("native.admits"), pair(inner, w.pair.proof(nil, nil, nil), con("source.nonmember")))
	second := sigma(con("N"), func(n expression) expression { return counterexampleBody(h, n) })
	return pair(counterexampleType(), h, pair(second, n, proof))(0)
}
