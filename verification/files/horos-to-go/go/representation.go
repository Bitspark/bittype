// Package verification composes static representation laws over captured Horos and Go subjects.
package verification

import (
 "errors"
 "fmt"
 "slices"
 association "github.com/Bitspark/bittype/contracts/verification/go"
 "github.com/Bitspark/bittype/semantics/go"
 n "github.com/Bitspark/bittype/semantics/golang/go"
 "github.com/Bitspark/bittype/trees/ast/go"
 h "github.com/Bitspark/bittype/verification/files/source/horos/go"
 joint "github.com/Bitspark/bittype/verification/logos/go"
 decl "github.com/Bitspark/horos/core/go/declarations"
 "github.com/Bitspark/logos-contract/contract/go/typed"
)

var ErrCoverage=errors.New("no selected representation covers the captured pair")
var ErrAmbiguous=errors.New("several selected representations cover the same pair")

type FieldBinding struct{ Source, Target string }
type binding struct { source decl.ParameterID; target n.Type }
type subject struct {
 parent *subject
 declaration h.Declaration
 expression *decl.Expression
 native n.Type
 optional bool
 fields []FieldBinding
 binder binding
 sourceSite,nativeSite ast.NodeRef
}
// Subject is issued for one occurrence under the current binder. Its constructors
// expose only operands, fields and arguments of that occurrence.
type Subject struct{ s *subject }
func (p Subject) Source() *decl.Expression { return p.s.expression }
func (p Subject) Target() n.Type { return p.s.native }
func (p Subject) Optional() bool { return p.s.optional }
func (p Subject) Declaration() h.Declaration { return p.s.declaration }
func (p Subject) Bindings() []FieldBinding { return slices.Clone(p.s.fields) }
func (p Subject) child(e *decl.Expression,t n.Type,optional bool,s,n ast.NodeRef) Subject {
 return Subject{&subject{parent:p.s,expression:e,native:t,optional:optional,binder:p.s.binder,sourceSite:s,nativeSite:n}}
}
// Operand pairs one actual Horos operand with a direct operand or argument
// of the captured native type. A type from another capture cannot be substituted.
func (p Subject) Operand(index int, target n.Type) (Subject,error) {
 ops:=p.Source().Operands()
 if index<0 || index>=len(ops) || !target.Valid() || !p.Target().Model().Same(target.Model()) { return Subject{},joint.ErrInvalid }
 term,ok:=p.Target().Model().Graph().Term(p.Target().ID())
 if !ok {return Subject{},joint.ErrInvalid}
 found:=slices.Contains(term.Arguments,target.ID())
 for _,op:=range term.Operands { if op.Term==target.ID(){found=true;break} }
 if !found{return Subject{},joint.ErrInvalid}
 return p.child(ops[index].Expression(),target,false,p.s.sourceSite,p.s.nativeSite),nil
}
// Required pairs the field's original type with the selected helper argument.
func (p Subject) Required(index int) (Subject,error) {
 args:=p.Target().Arguments()
 if !p.Optional() || index<0 || index>=len(args){return Subject{},joint.ErrInvalid}
 return p.child(p.Source(),args[index],false,p.s.sourceSite,p.s.nativeSite),nil
}
// Field selects actual fields by name. A library's record law owns the
// correspondence, including completeness, optionality, tags and embedding.
func (p Subject) Field(source,target string)(Subject,error){
 sf,err:=p.Declaration().Fields()
 nf,ok:=p.Target().Fields()
 if err!=nil || !ok{return Subject{},ErrCoverage}
 for _,a:=range sf { if a.Name==source {
  for _,b:=range nf {if b.Name==target {
   return p.child(a.Type,b.Type,a.Optional,a.Site,b.Site),nil
  }}
 }}
 return Subject{},ErrCoverage
}

// Rule contributes a structural law: representations of all returned child
// pairs entail representation of this pair. Selection never admits a dependency.
type Rule struct {
 Key semantics.Key
 Dependency string
 Project func(Subject)(children []Subject, applies bool, err error)
}
type Mappings struct{ rules []Rule }
func Select(rules ...Rule)(Mappings,error){
 if len(rules)>128{return Mappings{},association.ErrLimit}
 keys:=make(map[semantics.Key]bool,len(rules))
 for _,r:=range rules{
  if !r.Key.Valid() || keys[r.Key] || r.Dependency=="" || r.Project==nil{return Mappings{},joint.ErrInvalid}
  if len(r.Key.Owner)+len(r.Key.Name)+len(r.Dependency)>4096{return Mappings{},association.ErrLimit}
  keys[r.Key]=true
 }
 return Mappings{slices.Clone(rules)},nil
}

type node struct{
 symbol,dependency string
 children []*node
 sourceSites,nativeSites []association.Site
 bound bool
}
func (m Mappings) project(root Subject,prefix string)(*node,error){
 count:=0
 var visit func(Subject,int)(*node,error)
 visit=func(p Subject,depth int)(*node,error){
  count++
  if count>512 || depth>32{return nil,association.ErrLimit}
  if p.s==nil || p.Source()==nil || !p.Target().Valid(){return nil,joint.ErrInvalid}
  result:=&node{sourceSites:[]association.Site{site(sourceKey(),p.s.sourceSite)},nativeSites:[]association.Site{site(targetKey(),p.s.nativeSite)}}
  if loc:=p.Target().Origin().Site();loc.Node!=0 && loc!=p.s.nativeSite{ result.nativeSites=append(result.nativeSites,site(targetKey(),loc)) }
  if id,ok:=p.Source().Parameter();ok && !p.Optional(){
   if id!=p.s.binder.source || !p.Target().Same(p.s.binder.target){return nil,ErrCoverage}
   result.bound=true
   return result,nil
  }
  var children []Subject
  for _,rule:=range m.rules {
   inputs,applies,err:=rule.Project(p)
   if err!=nil{return nil,err}
   if !applies{continue}
   if result.symbol!="" {return nil,ErrAmbiguous}
   if len(inputs)>64 {return nil,association.ErrLimit}
   // A position identifies the captured occurrence. Length framing keeps
   // distinct owner/name pairs from aliasing in a Logos statement.
   result.symbol=fmt.Sprintf("%s.%d.%d.%s.%s",prefix,count,len(rule.Key.Owner),rule.Key.Owner,rule.Key.Name)
   result.dependency=rule.Dependency
   children=slices.Clone(inputs)
  }
  if result.symbol==""{return nil,ErrCoverage}
  for _,child:=range children{
   if child.s==nil || child.s.parent!=p.s {return nil,joint.ErrInvalid}
   projected,err:=visit(child,depth+1)
   if err!=nil{return nil,err}
   result.children=append(result.children,projected)
  }
  return result,nil
 }
 return visit(root,0)
}

// signature records the exact sites supplied at declaration construction.
type signature struct { declarations []typed.Declaration; origins []joint.Origin }
func(s *signature)add(d typed.Declaration,sites ...association.Site){
 if len(sites)>0{s.origins=append(s.origins,joint.Origin{Location:typed.Location{Phase:typed.DeclarationPhase,Index:uint32(len(s.declarations)),HasIndex:true},Sites:slices.Clone(sites)})}
 s.declarations=append(s.declarations,d)
}
func(s *signature)include(part signature){
 offset:=uint32(len(s.declarations))
 s.declarations=append(s.declarations,part.declarations...)
 for _,o:=range part.origins { o.Location.Index+=offset;s.origins=append(s.origins,o) }
}
func(s signature)problem(input association.Input,req []association.Requirement,expected typed.Expr,sites []association.Site) joint.Problem {
 origins:=append([]joint.Origin{{Location:typed.Location{Phase:typed.ExpectedPhase},Sites:sites}},s.origins...)
 return joint.Problem{Input:input,Requirements:req,Statement:typed.Statement{Foundation:typed.Foundation,Declarations:s.declarations,Expected:expected},Origins:origins}
}
func(x *node)sites()[]association.Site{return append(slices.Clone(x.sourceSites),x.nativeSites...)}
func(x *node)describe(s *signature){
 if x.bound{return}
 for _,c:=range x.children{c.describe(s)}
 count:=len(x.children)
 hs,ns,law:=x.symbol+".source",x.symbol+".target",x.symbol+".law"
 s.add(assume(hs,domainsN(con("H"),count,con("H")),ModelBasis),x.sourceSites...)
 s.add(assume(ns,domainsN(con("N"),count,con("N")),ModelBasis),x.nativeSites...)
 typ:=quantify(con("H"),count,func(h []expression)expression{
  return quantify(con("N"),count,func(n []expression)expression{
   result:=rep(app(con(hs),h...),app(con(ns),n...))
   for i:=count-1;i>=0;i--{result=arrow(rep(h[i],n[i]),result)}
   return result
  })
 })
 s.add(assume(law,typ,x.dependency),x.sites()...)
}
func(x *node)terms(h,n expression)(expression,expression){
 if x.bound {return h,n}
 var hs,ns []expression
 for _,c:=range x.children{a,b:=c.terms(h,n);hs=append(hs,a);ns=append(ns,b)}
 return app(con(x.symbol+".source"),hs...),app(con(x.symbol+".target"),ns...)
}
func(x *node)proof(h,n,witness expression)expression{
 if x.bound{return witness}
 var hs,ns,proofs []expression
 for _,c:=range x.children{a,b:=c.terms(h,n);hs=append(hs,a);ns=append(ns,b);proofs=append(proofs,c.proof(h,n,witness))}
 return app(con(x.symbol+".law"),append(append(hs,ns...),proofs...)...)
}
