package verification

import (
	"slices"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	association "github.com/Bitspark/bittype/contracts/verification/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	h "github.com/Bitspark/bittype/verification/files/source/horos/go"
	joint "github.com/Bitspark/bittype/verification/logos/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	"github.com/Bitspark/logos-contract/contract/go/typed"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

type namedBody struct {
	source         h.Declaration
	target         n.Type
	representation *node
}
type family struct {
	guarantee           h.Judgment
	library             Families
	body                namedBody
	kind                h.Kind
	sig                 signature
	expected, candidate typed.Expr
	theorem             joint.Claim
}
type Family struct{ state *family }
type argument struct {
	member    h.Judgment
	family    *family
	sig       signature
	body      *node
	candidate typed.Expr
}
type Argument struct{ state *argument }
type application struct {
	family   *family
	argument *argument
	theorem  joint.Result
	depth    int
}

// Families owns the named-declaration questions. Its receiver authorities
// bind the same mapping selection used to interpret every field and argument.
type Families struct {
	mappings    Mappings
	theorem     joint.Interpretation[*family]
	application joint.Interpretation[application]
}

func NewFamilies(m Mappings) Families {
	return Families{m, joint.Interpret(question("conditional-representation"), familyProblem), joint.Interpret(question("application-representation"), applicationProblem)}
}
func (l Families) Authorities() []joint.Authority {
	return []joint.Authority{l.theorem.Authority(), l.application.Authority()}
}
func sourceKey() semantics.Key { return semantics.Key{Owner: "horos-go", Name: "source"} }
func targetKey() semantics.Key { return semantics.Key{Owner: "horos-go", Name: "target"} }
func question(name string) semantics.Key {
	return semantics.Key{Owner: "horos-go/static-1", Name: name}
}
func site(key semantics.Key, node ast.NodeRef) association.Site {
	return association.Site{Participant: key, Node: node}
}
func associationInput(source h.Model, target n.Model) (association.Input, []association.Requirement) {
	return association.Input{Source: sourceKey(), Target: targetKey(), Participants: []association.Participant{
		{Key: sourceKey(), Snapshot: source.Snapshot(), Evidence: []semantic.Result{source.Receipt()}},
		{Key: targetKey(), Snapshot: target.Snapshot(), Evidence: []semantic.Result{target.Receipt()}},
	}}, []association.Requirement{
		{Key: sourceKey(), Language: source.Snapshot().Language(), Evidence: []semantic.EvidenceRequirement{h.Evidence()}},
		{Key: targetKey(), Language: target.Snapshot().Language(), Evidence: []semantic.EvidenceRequirement{n.Evidence()}},
	}
}
func (m Mappings) bindBody(source h.Declaration, target n.Type, fields []FieldBinding, prefix string) (namedBody, error) {
	if len(fields) > 64 {
		return namedBody{}, association.ErrLimit
	}
	if !source.Valid() || !target.IsNamed() {
		return namedBody{}, ErrCoverage
	}
	sp, np := source.Definition().Parameters(), target.Parameters()
	if len(sp) != len(np) || len(sp) > 1 {
		return namedBody{}, ErrCoverage
	}
	var binder binding
	if len(sp) == 1 {
		if !np[0].Unconstrained() {
			return namedBody{}, ErrCoverage
		}
		binder = binding{decl.ParameterID{Declaration: source.Definition().ID(), Index: 0}, np[0]}
	}
	root := Subject{&subject{declaration: source, sourceModel: source.Model(), expression: source.Definition().Expression(), native: target, fields: slices.Clone(fields), binder: binder, sourceSite: source.Site(), nativeSite: target.Site()}}
	body, err := m.project(root, prefix)
	if err != nil {
		return namedBody{}, err
	}
	return namedBody{source, target, body}, nil
}
func sameKind(a, b *decl.Kind) bool { return a != nil && b != nil && a.Value().Equal(b.Value()) }
func declarationValue(v h.Value, d h.Declaration) bool {
	if !v.Model().Same(d.Model()) {
		return false
	}
	e := v.Expression()
	if e == nil {
		return false
	}
	id, ok := e.Target()
	return ok && id == d.Definition().ID()
}

// Bind admits the explicit constructor-kind query under the declaration's
// original contracts. No named kind (including Data) is recognized here.
func (l Families) Bind(source h.Declaration, target n.Type, fields []FieldBinding, kind h.Kind, guarantee h.Judgment) (Family, error) {
	if !source.Valid() || !source.Model().Same(kind.Model()) || !source.Model().Same(guarantee.Model()) {
		return Family{}, joint.ErrInvalid
	}
	if guarantee.Verdict() != decl.Proven {
		return Family{}, semantic.ErrEvidenceUnavailable
	}
	if !declarationValue(guarantee.Value(), source) || len(source.Definition().Parameters()) != 1 {
		return Family{}, ErrCoverage
	}
	k := kind.Expression()
	q := guarantee.Kind().Expression()
	if k == nil || q == nil {
		return Family{}, ErrCoverage
	}
	if !k.Sort().IsProper() || !source.Definition().Parameters()[0].Kind.Sort().IsProper() {
		return Family{}, ErrCoverage
	}
	want := core.NewTuple(core.NewAtom([]byte("horos/decl/arrow")), data.EncodeList([]core.Value{k.Value()}), k.Value())
	if guarantee.Kind().Module() != kind.Module() || !q.Value().Equal(want) {
		return Family{}, ErrCoverage
	}
	r, err := l.mappings.bindBody(source, target, fields, "family")
	if err != nil {
		return Family{}, err
	}
	f := &family{body: r, kind: kind, library: l, guarantee: guarantee}
	f.sig, f.expected, f.candidate = familyLogic(r)
	f.theorem, err = l.theorem.Ask(f)
	if err != nil {
		return Family{}, err
	}
	return Family{f}, nil
}
func (f Family) Theorem() (joint.Claim, error) {
	if f.state == nil {
		return joint.Claim{}, joint.ErrInvalid
	}
	return f.state.theorem, nil
}
func familyProblem(f *family) (joint.Problem, error) {
	input, req := associationInput(f.guarantee.Model(), f.body.target.Model())
	return f.sig.problem(input, req, f.expected, f.body.representation.sites()), nil
}

// Candidate returns ordinary proof syntax. It carries no acceptance authority.
func (f Family) Candidate() typed.Expr {
	if f.state == nil {
		return typed.Expr{}
	}
	return f.state.candidate
}

func (f Family) BindArgument(source h.Declaration, target n.Type, fields []FieldBinding, membership h.Judgment) (Argument, error) {
	if f.state == nil || !source.Valid() || !source.Model().Same(f.state.body.source.Model()) || !target.Model().Same(f.state.body.target.Model()) || !source.Model().Same(membership.Model()) {
		return Argument{}, joint.ErrInvalid
	}
	if membership.Verdict() != decl.Proven {
		return Argument{}, semantic.ErrEvidenceUnavailable
	}
	if !membership.Kind().Same(f.state.kind) || len(source.Definition().Parameters()) != 0 {
		return Argument{}, ErrCoverage
	}
	// Value() of a nullary declaration is its body, so compare the authoritative
	// descriptor with this declaration's exact descriptor after its explicit query.
	descriptor, err := membership.Value().Descriptor()
	if err != nil {
		return Argument{}, err
	}
	if !descriptor.Equal(source.Definition().Expression().Value()) {
		return Argument{}, ErrCoverage
	}
	r, err := f.state.library.mappings.bindBody(source, target, fields, "argument")
	if err != nil {
		return Argument{}, err
	}
	d, proof := argumentLogic(r)
	return Argument{&argument{membership, f.state, d, r.representation, proof}}, nil
}

// Application retains the checked closed theorem as a definition and rechecks
// it with each application. It never turns an acceptance into an assumption.
func (f Family) Application(theorem joint.Result, arg Argument, depth int) (joint.Claim, error) {
	if f.state == nil || arg.state == nil || arg.state.family != f.state || !theorem.Request().Answers(f.state.theorem) || !theorem.Outcome().IsAccepted() {
		return joint.Claim{}, joint.ErrInvalid
	}
	if depth < 1 || depth > 4 {
		return joint.Claim{}, association.ErrLimit
	}
	return f.state.library.application.Ask(application{f.state, arg.state, theorem, depth})
}
func applicationProblem(a application) (joint.Problem, error) {
	f := a.family
	var s signature
	s.include(f.sig)
	s.add(typed.Define("family.theorem", a.theorem.Request().Statement().Expected, a.theorem.Candidate()), f.body.representation.sites()...)
	s.include(a.argument.sig)
	h, n := con("h.argument"), con("n.argument")
	for range a.depth {
		h = app(con("h.family"), h)
		n = app(con("n.family"), n)
	}
	input, req := associationInput(a.argument.member.Model(), f.body.target.Model())
	return s.problem(input, req, conjunction(member(h), rep(h, n))(0), append(f.body.representation.sites(), a.argument.body.sites()...)), nil
}
func (a Argument) Candidate(depth int) typed.Expr {
	if a.state == nil || depth < 1 || depth > 4 {
		return typed.Expr{}
	}
	h, n := con("h.argument"), con("n.argument")
	p := pair(conjunction(member(h), rep(h, n)), con("argument.member"), fixed(a.state.candidate))
	for range depth {
		p = app(con("family.theorem"), h, n, fst(p), snd(p))
		h = app(con("h.family"), h)
		n = app(con("n.family"), n)
	}
	return p(0)
}
