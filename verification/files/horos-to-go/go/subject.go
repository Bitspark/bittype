package verification

import (
	n "github.com/Bitspark/bittype/semantics/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	h "github.com/Bitspark/bittype/verification/files/source/horos/go"
	joint "github.com/Bitspark/bittype/verification/logos/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	"slices"
)

type FieldBinding struct{ Source, Target string }
type binding struct {
	source decl.ParameterID
	target n.Type
}
type subject struct {
	parent                 *subject
	declaration            h.Declaration
 sourceModel h.Model
	expression             *decl.Expression
	native                 n.Type
	optional               bool
	fields                 []FieldBinding
	binder                 binding
	sourceSite, nativeSite ast.NodeRef
}

// Subject is issued for one occurrence under the current binder. Its constructors
// expose only operands, fields and arguments of that occurrence.
type Subject struct{ s *subject }

func(p Subject)Source()*decl.Expression{
 if p.s==nil || p.s.expression==nil{return nil}
 // Horos exposes immutable operand handles and copying accessors. Isolate the
 // outer expression as well: a contributor can replace a value it receives.
 view:=*p.s.expression
 return &view
}
func (p Subject) Target() n.Type             { return p.s.native }
func (p Subject) Optional() bool             { return p.s.optional }
func (p Subject) Declaration() h.Declaration { return p.s.declaration }
func (p Subject) Bindings() []FieldBinding   { return slices.Clone(p.s.fields) }
func (p Subject) child(e *decl.Expression, t n.Type, optional bool, s, n ast.NodeRef) Subject {
	return Subject{&subject{parent: p.s, sourceModel:p.s.sourceModel, expression: e, native: t, optional: optional, binder: p.s.binder, sourceSite: s, nativeSite: n}}
}

// Operand pairs one actual Horos operand with a direct operand or argument
// of the captured native type. A type from another capture cannot be substituted.
func (p Subject) Operand(index int, target n.Type) (Subject, error) {
	ops := p.Source().Operands()
	if index < 0 || index >= len(ops) || !target.Valid() || !p.Target().Model().Same(target.Model()) {
		return Subject{}, joint.ErrInvalid
	}
 if !p.directNative(target){return Subject{},joint.ErrInvalid}
 site:=p.s.sourceSite
 if p.s.sourceModel.Valid(){
  operand:=p.s.sourceModel.OperandSite(site,ops[index].Path())
  if operand.Node==0{return Subject{},joint.ErrInvalid}
  site=operand
 }
 return p.child(ops[index].Expression(),target,false,site,p.s.nativeSite),nil
}

// Payload uses the same source value under a projected native child. For
// a field this consumes the field occurrence's optionality. The mapping owns
// how the native representation exposes that value.
func(p Subject)Payload(target n.Type)(Subject,error){
 same:=p.Target().Same(target)
 if same && !p.Optional() || !same && !p.directNative(target){return Subject{},joint.ErrInvalid}
 return p.child(p.Source(),target,false,p.s.sourceSite,p.s.nativeSite),nil
}
func(p Subject)directNative(target n.Type)bool{
 if !target.Valid() || !p.Target().Model().Same(target.Model()){return false}
 term,ok:=p.Target().Model().Graph().Term(p.Target().ID())
 if !ok{return false}
 if slices.Contains(term.Arguments,target.ID()){return true}
 for _,op:=range term.Operands{if op.Term==target.ID(){return true}}
 return false
}

// Field selects actual fields by name. A library's record law owns the
// correspondence, including completeness, optionality, tags and embedding.
func (p Subject) Field(source, target string) (Subject, error) {
	sf, err := p.Declaration().Fields()
	nf, ok := p.Target().Fields()
	if err != nil || !ok {
		return Subject{}, ErrCoverage
	}
	for _, a := range sf {
		if a.Name == source {
			for _, b := range nf {
				if b.Name == target {
					return p.child(a.Type, b.Type, a.Optional, a.Site, b.Site), nil
				}
			}
		}
	}
	return Subject{}, ErrCoverage
}
