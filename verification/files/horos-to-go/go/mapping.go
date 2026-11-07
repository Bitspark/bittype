package verification

import (
	"slices"
 joint "github.com/Bitspark/bittype/verification/logos/go"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	association "github.com/Bitspark/bittype/contracts/verification/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	h "github.com/Bitspark/bittype/verification/files/source/horos/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	"github.com/Bitspark/logos-contract/contract/go/typed"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

type record struct{
 source h.Declaration
 target n.Type
 body *node
}
type family struct{
 library Families
 record record
 kind h.Kind
 sig signature
 expected,candidate typed.Expr
 theorem joint.Claim
}
type Family struct{state *family}
type argument struct{
 family *family
 sig signature
 body *node
 candidate typed.Expr
}
type Argument struct{state *argument}
type application struct{
 family *family
 argument *argument
 theorem joint.Result
 depth int
}
// Families owns the named-declaration questions. Its receiver authorities
// bind the same mapping selection used to interpret every field and argument.
type Families struct{
 mappings Mappings
 theorem joint.Interpretation[*family]
 application joint.Interpretation[application]
}
func NewFamilies(m Mappings) Families{
 return Families{m,joint.Interpret(question("conditional-representation"),familyProblem),joint.Interpret(question("application-representation"),applicationProblem)}
}
func(l Families) Authorities()[]joint.Authority{return []joint.Authority{l.theorem.Authority(),l.application.Authority()}}
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
func (m Mappings) bindRecord(source h.Declaration,target n.Type,fields []FieldBinding,prefix string)(record,error){
 if !source.Valid() || !target.IsNamed() || len(fields)>64{return record{},ErrCoverage}
 sp,np:=source.Definition().Parameters(),target.Parameters()
 if len(sp)!=len(np) || len(sp)>1{return record{},ErrCoverage}
 var binder binding
 if len(sp)==1{
  if !np[0].Unconstrained(){return record{},ErrCoverage}
  binder=binding{decl.ParameterID{Declaration:source.Definition().ID(),Index:0},np[0]}
 }
 root:=Subject{&subject{declaration:source,expression:source.Definition().Expression(),native:target,fields:slices.Clone(fields),binder:binder,sourceSite:source.Site(),nativeSite:target.Site()}}
 body,err:=m.project(root,prefix)
 if err!=nil{return record{},err}
 return record{source,target,body},nil
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
	r, err := l.mappings.bindRecord(source, target, fields,"family")
	if err != nil {
		return Family{}, err
	}
	f := &family{record:r, kind:kind, library:l}
 f.sig,f.expected,f.candidate=familyLogic(r)
 f.theorem=l.theorem.Ask(f)
	return Family{f}, nil
}
func (f Family) Theorem()(joint.Claim,error){
 if f.state==nil{return joint.Claim{},joint.ErrInvalid}
 return f.state.theorem,nil
}
func familyProblem(f *family)(joint.Problem,error){
 input,req:=associationInput(f.record.source.Model(),f.record.target.Model())
 return f.sig.problem(input,req,f.expected,f.record.body.sites()),nil
}
// Candidate returns ordinary proof syntax. It carries no acceptance authority.
func (f Family) Candidate() typed.Expr {
	if f.state == nil {
		return typed.Expr{}
	}
	return f.state.candidate
}

func (f Family) BindArgument(source h.Declaration, target n.Type, fields []FieldBinding, membership h.Judgment) (Argument, error) {
	if f.state == nil || !source.Valid() || !source.Model().Same(f.state.record.source.Model()) || !target.Model().Same(f.state.record.target.Model()) || !source.Model().Same(membership.Model()) {
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
	r, err := f.state.library.mappings.bindRecord(source, target, fields,"argument")
	if err != nil {
		return Argument{}, err
	}
	d, proof := argumentLogic(r)
	return Argument{&argument{f.state,d,r.body,proof}}, nil
}

// Application retains the checked closed theorem as a definition and rechecks
// it with each application. It never turns an acceptance into an assumption.
func(f Family)Application(theorem joint.Result,arg Argument,depth int)(joint.Claim,error){
 if f.state==nil || arg.state==nil || arg.state.family!=f.state || !theorem.Request().Answers(f.state.theorem) || !theorem.Outcome().IsAccepted(){return joint.Claim{},joint.ErrInvalid}
 if depth<1 || depth>4{return joint.Claim{},association.ErrLimit}
 return f.state.library.application.Ask(application{f.state,arg.state,theorem,depth}),nil
}
func applicationProblem(a application)(joint.Problem,error){
 f:=a.family
 var s signature
 s.include(f.sig)
 s.add(typed.Define("family.theorem",a.theorem.Request().Statement().Expected,a.theorem.Candidate()),f.record.body.sites()...)
 s.include(a.argument.sig)
 h,n:=con("h.argument"),con("n.argument")
 for range a.depth {h=app(con("h.family"),h);n=app(con("n.family"),n)}
 input,req:=associationInput(f.record.source.Model(),f.record.target.Model())
 return s.problem(input,req,conjunction(member(h),rep(h,n))(0),append(f.record.body.sites(),a.argument.body.sites()...)),nil
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
