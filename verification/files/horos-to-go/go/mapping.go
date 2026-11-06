package verification

import (
	"fmt"
	"go/types"
	"slices"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	association "github.com/Bitspark/bittype/contracts/verification/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	h "github.com/Bitspark/bittype/verification/files/source/horos/go"
	n "github.com/Bitspark/bittype/verification/files/target/golang/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	"github.com/Bitspark/logos/contract/go/typed"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

type FieldBinding struct { Source, Target string }
// Optional selects the actual ordinary helper and its two named fields.
// This is a static representation law; it supplies no runtime conversion.
type Optional struct { Constructor n.Type; PresenceField, ValueField string }
type pattern struct { shape string; element *pattern }
type record struct { source h.Declaration; target n.Type; fields []FieldBinding; patterns []pattern; sites []association.Site }
type family struct { record record; kind h.Kind; guarantee h.Judgment; optional Optional; declarations []typed.Declaration; expected typed.Expr; candidate typed.Expr }
type Family struct { state *family }
type argument struct { family *family; record record; member h.Judgment; declarations []typed.Declaration; candidate typed.Expr }
type Argument struct { state *argument }

func sourceKey()semantics.Key{return semantics.Key{Owner:"horos-go",Name:"source"}}
func targetKey()semantics.Key{return semantics.Key{Owner:"horos-go",Name:"target"}}
func question(name string)semantics.Key{return semantics.Key{Owner:"horos-go/static-1",Name:name}}
func site(key semantics.Key,node ast.NodeRef)association.Site{return association.Site{Participant:key,Node:node}}
func associationInput(source h.Model,target n.Model)(association.Input,[]association.Requirement){
	return association.Input{Source:sourceKey(),Target:targetKey(),Participants:[]association.Participant{
		{Key:sourceKey(),Snapshot:source.Snapshot(),Evidence:[]semantic.Result{source.Receipt()}},
		{Key:targetKey(),Snapshot:target.Snapshot(),Evidence:[]semantic.Result{target.Receipt()}},
	}},[]association.Requirement{
		{Key:sourceKey(),Language:source.Snapshot().Language(),Evidence:[]semantic.EvidenceRequirement{h.Evidence()}},
		{Key:targetKey(),Language:target.Snapshot().Language(),Evidence:[]semantic.EvidenceRequirement{n.Evidence()}},
	}
}
func scalar(e *decl.Expression)string{
	if e==nil || e.Kind()!=decl.Descriptor || e.Constructor()!="scalar"{return ""}
	t,ok:=e.Value().(core.Tuple);if !ok || t.Len()!=2{return ""};a,ok:=t.At(1).(core.Atom);if !ok{return ""}
	switch string(a.Bytes()){case "utf8-text":return "text";case "bool":return "bool"};return ""
}
func optionalValid(o Optional)bool{
	if !o.Constructor.IsNamed() || o.PresenceField=="" || o.ValueField=="" || o.PresenceField==o.ValueField{return false}
	p:=o.Constructor.Parameters();fields,ok:=o.Constructor.Fields();if !ok || len(p)!=1 || !p[0].Unconstrained() || len(fields)!=2{return false}
	present,value:=false,false
	for _,f:=range fields{if f.Embedded || f.Tag!=""{return false};switch f.Name{case o.PresenceField:present=f.Type.IsBasic(types.Bool);case o.ValueField:value=f.Type.Same(p[0])}}
	return present && value
}
func match(e *decl.Expression,target n.Type,sourceParameter decl.ParameterID,nativeParameter n.Type,depth int)(pattern,error){
	if depth>32{return pattern{},association.ErrLimit}
	if p,ok:=e.Parameter();ok{if p!=sourceParameter || !target.Same(nativeParameter){return pattern{},ErrCoverage};return pattern{shape:"parameter"},nil}
	if name:=scalar(e);name!=""{kind:=types.String;if name=="bool"{kind=types.Bool};if !target.IsBasic(kind){return pattern{},ErrCoverage};return pattern{shape:name},nil}
	if e.Kind()==decl.Descriptor && e.Constructor()=="list-of"{
		x,ok:=target.SliceElement();ops:=e.Operands();if !ok || len(ops)!=1{return pattern{},ErrCoverage}
		p,err:=match(ops[0].Expression(),x,sourceParameter,nativeParameter,depth+1);return pattern{shape:"list",element:&p},err
	}
	return pattern{},ErrCoverage
}
func bindRecord(source h.Declaration,target n.Type,fields []FieldBinding,optional Optional)(record,error){
	if !source.Valid() || !target.IsNamed() || len(fields)>64{return record{},ErrCoverage}
	sf,err:=source.Fields();if err!=nil{return record{},fmt.Errorf("%w: %v",ErrCoverage,err)}
	nf,ok:=target.Fields();if !ok || len(sf)!=len(fields) || len(nf)!=len(fields){return record{},ErrCoverage}
	sp,np:=source.Definition().Parameters(),target.Parameters();if len(sp)!=len(np) || len(sp)>1{return record{},ErrCoverage}
	var sourceParameter decl.ParameterID;var nativeParameter n.Type
	if len(sp)==1{sourceParameter=decl.ParameterID{Declaration:source.Definition().ID(),Index:0};nativeParameter=np[0];if !nativeParameter.Unconstrained(){return record{},ErrCoverage}}
	r:=record{source:source,target:target,fields:slices.Clone(fields),sites:[]association.Site{site(sourceKey(),source.Site()),site(targetKey(),target.Site())}}
	usedS,usedN:=map[string]bool{},map[string]bool{}
	for _,binding:=range fields{
		if usedS[binding.Source] || usedN[binding.Target]{return record{},ErrCoverage};usedS[binding.Source]=true;usedN[binding.Target]=true
		var a *h.Field;var b *n.Field
		for i:=range sf{if sf[i].Name==binding.Source{a=&sf[i]}}
		for i:=range nf{if nf[i].Name==binding.Target{b=&nf[i]}}
		if a==nil || b==nil || b.Embedded || b.Tag!=""{return record{},ErrCoverage}
		ty:=b.Type
		if a.Optional{
			if !optionalValid(optional) || !optional.Constructor.Model().Same(target.Model()) || !ty.Origin().Same(optional.Constructor){return record{},ErrCoverage}
			args:=ty.Arguments();if len(args)!=1{return record{},ErrCoverage};ty=args[0]
			r.sites=append(r.sites,site(targetKey(),optional.Constructor.Site()))
		}
		p,err:=match(a.Type,ty,sourceParameter,nativeParameter,0);if err!=nil{return record{},err}
		if a.Optional{p=pattern{shape:"optional",element:&p}}
		r.patterns=append(r.patterns,p);r.sites=append(r.sites,site(sourceKey(),a.Site),site(targetKey(),b.Site))
	}
	return r,nil
}
func sameKind(a,b *decl.Kind)bool{return a!=nil && b!=nil && a.Value().Equal(b.Value())}
func declarationValue(v h.Value,d h.Declaration)bool{
	if !v.Model().Same(d.Model()){return false};e:=v.Expression();if e==nil{return false};id,ok:=e.Target();return ok && id==d.Definition().ID()
}

// BindFamily admits the explicit constructor-kind query under the declaration's
// original contracts. No named kind (including Data) is recognized here.
func BindFamily(source h.Declaration,target n.Type,fields []FieldBinding,optional Optional,kind h.Kind,guarantee h.Judgment)(Family,error){
	if !source.Valid() || !source.Model().Same(kind.Model()) || !source.Model().Same(guarantee.Model()){return Family{},ErrInvalid}
	if guarantee.Verdict()!=decl.Proven{return Family{},semantic.ErrEvidenceUnavailable}
	if !declarationValue(guarantee.Value(),source) || len(source.Definition().Parameters())!=1{return Family{},ErrCoverage}
	k:=kind.Expression();q:=guarantee.Kind().Expression();if k==nil || q==nil{return Family{},ErrCoverage}
	want:=core.NewTuple(core.NewAtom([]byte("horos/decl/arrow")),data.EncodeList([]core.Value{k.Value()}),k.Value())
	if !q.Value().Equal(want){return Family{},ErrCoverage}
	r,err:=bindRecord(source,target,fields,optional);if err!=nil{return Family{},err}
	f:=&family{record:r,kind:kind,guarantee:guarantee,optional:optional}
	f.declarations,f.expected,f.candidate=familyLogic(r)
	return Family{f},nil
}
func (f Family) Theorem()(Claim,error){
	if f.state==nil{return Claim{},ErrInvalid};s:=f.state;input,req:=associationInput(s.record.source.Model(),s.record.target.Model())
	origins:=[]Origin{{Location:typed.Location{Phase:typed.ExpectedPhase},Sites:slices.Clone(s.record.sites)}}
	return Claim{&claim{input:input,requirements:req,question:question("conditional-representation"),statement:typed.Statement{Foundation:typed.Foundation,Declarations:slices.Clone(s.declarations),Expected:s.expected},origins:origins,family:s}},nil
}
// Candidate returns ordinary proof syntax. It carries no acceptance authority.
func (f Family) Candidate()typed.Expr{if f.state==nil{return typed.Expr{}};return f.state.candidate}

func (f Family) BindArgument(source h.Declaration,target n.Type,fields []FieldBinding,membership h.Judgment)(Argument,error){
	if f.state==nil || !source.Valid() || !source.Model().Same(f.state.record.source.Model()) || !target.Model().Same(f.state.record.target.Model()) || !source.Model().Same(membership.Model()){return Argument{},ErrInvalid}
	if membership.Verdict()!=decl.Proven{return Argument{},semantic.ErrEvidenceUnavailable}
	if !membership.Kind().Same(f.state.kind) || len(source.Definition().Parameters())!=0{return Argument{},ErrCoverage}
	// Value() of a nullary declaration is its body, so compare the authoritative
	// descriptor with this declaration's exact descriptor after its explicit query.
	descriptor,err:=membership.Value().Descriptor();if err!=nil{return Argument{},err}
	if !descriptor.Equal(source.Definition().Expression().Value()){return Argument{},ErrCoverage}
	r,err:=bindRecord(source,target,fields,f.state.optional);if err!=nil{return Argument{},err}
	d,proof:=argumentLogic(r)
	return Argument{&argument{f.state,r,membership,d,proof}},nil
}

// Application retains the checked closed theorem as a definition and rechecks
// it with each application. It never turns an acceptance into an assumption.
func (f Family) Application(theorem Result,arg Argument,depth int)(Claim,error){
	if f.state==nil || arg.state==nil || arg.state.family!=f.state || !theorem.Valid() || theorem.Request().state.claim.family!=f.state || theorem.Request().Question()!=question("conditional-representation") || !theorem.Outcome().IsAccepted(){return Claim{},ErrInvalid}
	if depth<1 || depth>4{return Claim{},association.ErrLimit}
	d:=slices.Clone(f.state.declarations)
	d=append(d,typed.Define("family.theorem",theorem.Request().Statement().Expected,theorem.Candidate()))
	d=append(d,arg.state.declarations...)
	h,n:=con("h.argument"),con("n.argument")
	for range depth{h=app(con("h.family"),h);n=app(con("n.family"),n)}
	input,req:=associationInput(f.state.record.source.Model(),f.state.record.target.Model())
	sites:=append(slices.Clone(f.state.record.sites),arg.state.record.sites...)
	return Claim{&claim{input:input,requirements:req,question:question("application-representation"),statement:typed.Statement{Foundation:typed.Foundation,Declarations:d,Expected:conjunction(member(h),rep(h,n))(0)},origins:[]Origin{{Location:typed.Location{Phase:typed.ExpectedPhase},Sites:sites}},family:f.state}},nil
}
func (a Argument) Candidate(depth int)typed.Expr{
	if a.state==nil || depth<1 || depth>4{return typed.Expr{}}
	h,n:=con("h.argument"),con("n.argument")
	p:=pair(conjunction(member(h),rep(h,n)),con("argument.member"),fixed(a.state.candidate))
	for range depth{p=app(con("family.theorem"),h,n,fst(p),snd(p));h=app(con("h.family"),h);n=app(con("n.family"),n)}
	return p(0)
}
