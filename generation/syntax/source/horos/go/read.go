package source

import (
	"fmt"
	"slices"
	"strings"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/bitspark/ontos/codec/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// Expression is a source-owned reading of descriptor operands. Native rules
// live in supplied mappings. Literal payloads are never traversed as types.
type Expression struct {
	Constructor string
	Discriminator string
	Origin ast.NodeRef
	Operands []Operand
}
type Operand struct { Path []uint32; Use g.Use[h.Language] }
type Field struct { Name string; Optional bool; Use g.Use[h.Language]; Origin ast.NodeRef }
type Declaration struct { Name string; ModulePin string; Root bool; Origin ast.NodeRef; Body g.Use[h.Language] }
type Parameter struct { Name string; Origin ast.NodeRef }

func (s *Source) Declaration(d g.Declaration[h.Language]) (Declaration,error) {
	i,ok:=s.declarations.Index(d);if !ok||s.entities[i]==0{return Declaration{},g.ErrScope}
	e:=s.facts.Entities[i-1];pin,err:=s.modulePin(e.Owner);if err!=nil{return Declaration{},err}
	return Declaration{Name:e.Name,ModulePin:pin,Root:slices.Contains(s.roots,d),Origin:e.Declarations[0],Body:s.use(e.Body)},nil
}
func (s *Source) Parameter(b g.Binder[h.Language]) (Parameter,error) {
	i,ok:=s.binders.Index(b);if !ok||s.bound[semantics.EntityID(i)]!=b{return Parameter{},g.ErrScope};e:=s.facts.Entities[i-1];if len(e.Declarations)!=1{return Parameter{},semantic.ErrInvalidRequest};return Parameter{e.Name,e.Declarations[0]},nil
}
func (s *Source) IsBody(env g.Environment[h.Language],use g.Use[h.Language]) bool {
	i,ok:=s.environments.Index(env);if !ok||i<=1||s.entities[i-1]==0{return false};return s.use(s.facts.Entities[i-2].Body)==use
}
func (s *Source) Read(use g.Use[h.Language]) (Expression,error) {
	i,formal,err:=s.useIndex(use);if err!=nil{return Expression{},err};if formal!=0{return Expression{},g.ErrUnsupported}
	t:=s.facts.Terms[i-1];v:=Expression{Origin:s.sites[semantics.TermID(i)]}
	if t.Kind==semantics.Reference||t.Kind==semantics.Application { v.Constructor="reference";return v,nil }
	if t.Kind!=semantics.ProviderExpression||t.Operator.Owner!="horos" { return v,semantic.ErrInvalidRequest }
	v.Constructor=strings.TrimPrefix(t.Operator.Name,"descriptor/")
	for _,o:=range t.Operands { v.Operands=append(v.Operands,Operand{Path:slices.Clone(o.Locator),Use:s.use(o.Term)}) }
	if v.Constructor=="scalar"||v.Constructor=="record" {
		p,err:=s.payload(semantics.TermID(i));if err!=nil{return v,err};if p.Len()<2{return v,semantic.ErrInvalidRequest};v.Discriminator=atom(p.At(1))
	}
	return v,nil
}
func (s *Source) Operand(use g.Use[h.Language],path ...uint32) (g.Use[h.Language],error) {
	v,err:=s.Read(use);if err!=nil{return g.Use[h.Language]{},err};for _,o:=range v.Operands { if slices.Equal(o.Path,path){return o.Use,nil} };return g.Use[h.Language]{},semantic.ErrInvalidRequest
}
func (s *Source) Fields(use g.Use[h.Language]) ([]Field,error) {
	i,formal,err:=s.useIndex(use);if err!=nil||formal!=0{return nil,g.ErrScope};p,err:=s.payload(semantics.TermID(i));if err!=nil{return nil,err}
	if p.Len()<2||atom(p.At(0))!="horos/record"{return nil,semantic.ErrInvalidRequest}
	var fields []Field
	for n:=2;n<p.Len();n++ {
		f,ok:=p.At(n).(core.Tuple);if !ok||f.Len()!=4{return nil,semantic.ErrInvalidRequest};name,err:=data.ReadText(f.At(1));if err!=nil{return nil,err}
		u,err:=s.Operand(use,uint32(n),3);if err!=nil{return nil,err};view,err:=s.Read(u);if err!=nil{return nil,err}
		fields=append(fields,Field{Name:name,Optional:atom(f.At(2))=="opt",Use:u,Origin:view.Origin})
	}
	return fields,nil
}
func atom(v core.Value) string { a,ok:=v.(core.Atom);if !ok{return ""};return string(a.Bytes()) }
func (s *Source) payload(id semantics.TermID) (core.Tuple,error) {
	t:=s.facts.Terms[id-1];if t.Encoding!=(semantics.Key{Owner:"ontos",Name:"codec-v1"}){return core.Tuple{},semantic.ErrInvalidRequest};v,err:=codec.Decode(t.Data);if err!=nil{return core.Tuple{},err};p,ok:=v.(core.Tuple);if !ok{return core.Tuple{},semantic.ErrInvalidRequest};return p,nil
}
func (s *Source) modulePin(id semantics.EntityID) (string,error) {
	term:=s.identities[id];if term==0{return "",semantic.ErrInvalidRequest};t:=s.facts.Terms[term-1]
	if t.Operator!=(semantics.Key{Owner:"horos",Name:"identity/module"})||t.Encoding!=(semantics.Key{Owner:"ontos",Name:"codec-v1"}){return "",semantic.ErrInvalidRequest}
	v,err:=codec.Decode(t.Data);if err!=nil{return "",err};a,ok:=v.(core.Atom);if !ok||len(a.Bytes())!=32{return "",semantic.ErrInvalidRequest};return fmt.Sprintf("%x",a.Bytes()),nil
}
