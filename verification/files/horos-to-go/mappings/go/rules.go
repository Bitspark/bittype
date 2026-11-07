// Package mappings supplies independently selectable static Horos/Go representation laws.
package mappings

import (
 "go/types"
 "github.com/Bitspark/bittype/semantics/go"
 n "github.com/Bitspark/bittype/semantics/golang/go"
 v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
 decl "github.com/Bitspark/horos/core/go/declarations"
 core "github.com/bitspark/ontos/core/go"
)

const (
 ScalarLaw="horos-go/static-scalars-1"
 ListLaw="horos-go/list-slice-1"
 OptionalLaw="horos-go/optional-presence-1"
 RecordLaw="horos-go/record-fields-1"
)
func key(name string)semantics.Key{return semantics.Key{Owner:"horos-go/static",Name:name}}

func scalar(name string,kind types.BasicKind) v.Rule{
 return v.Rule{Key:key(name),Dependency:ScalarLaw,Project:func(p v.Subject)([]v.Subject,bool,error){
  if p.Optional() || p.Source().Kind()!=decl.Descriptor || p.Source().Constructor()!="scalar"{return nil,false,nil}
  value,ok:=p.Source().Value().(core.Tuple)
  if !ok || value.Len()!=2{return nil,false,nil}
  atom,ok:=value.At(1).(core.Atom)
  return nil,ok && string(atom.Bytes())==name && p.Target().IsBasic(kind),nil
 }}
}
func Text() v.Rule{return scalar("utf8-text",types.String)}
func Boolean() v.Rule{return scalar("bool",types.Bool)}

func List() v.Rule{return v.Rule{Key:key("list-slice"),Dependency:ListLaw,Project:list}}
func list(p v.Subject)([]v.Subject,bool,error){
 element,ok:=p.Target().SliceElement()
 if p.Optional() || !ok || p.Source().Kind()!=decl.Descriptor || p.Source().Constructor()!="list-of" || len(p.Source().Operands())!=1 {return nil,false,nil}
 child,err:=p.Operand(0,element)
 return []v.Subject{child},true,err
}

type Optional struct{
 Constructor n.Type
 PresenceField,ValueField string
}
func(o Optional)Rule()v.Rule{return v.Rule{Key:key("optional-presence"),Dependency:OptionalLaw,Project:o.project}}
func(o Optional)project(p v.Subject)([]v.Subject,bool,error){
 if !p.Optional() || !p.Target().Origin().Same(o.Constructor){return nil,false,nil}
 params:=o.Constructor.Parameters()
 fields,ok:=o.Constructor.Fields()
 if !o.Constructor.IsNamed() || !ok || len(params)!=1 || !params[0].Unconstrained() || len(fields)!=2 || o.PresenceField==o.ValueField || len(p.Target().Arguments())!=1{return nil,false,nil}
 present,value:=false,false
 for _,f:=range fields{
  if f.Embedded || f.Tag!=""{return nil,false,nil}
  if f.Name==o.PresenceField{present=f.Type.IsBasic(types.Bool)}
  if f.Name==o.ValueField{value=f.Type.Same(params[0])}
 }
 if !present || !value{return nil,false,nil}
 child,err:=p.Required(0)
 return []v.Subject{child},true,err
}

func Record()v.Rule{return v.Rule{Key:key("closed-record"),Dependency:RecordLaw,Project:record}}
func record(p v.Subject)([]v.Subject,bool,error){
 if p.Optional() || !p.Declaration().Valid() || !p.Target().IsNamed(){return nil,false,nil}
 source,err:=p.Declaration().Fields()
 native,ok:=p.Target().Fields()
 bindings:=p.Bindings()
 if err!=nil || !ok || len(bindings)!=len(source) || len(bindings)!=len(native){return nil,false,nil}
 for _,f:=range native{if f.Embedded || f.Tag!=""{return nil,false,nil}}
 usedSource,usedTarget:=map[string]bool{},map[string]bool{}
 children:=make([]v.Subject,0,len(bindings))
 for _,b:=range bindings{
  if usedSource[b.Source] || usedTarget[b.Target]{return nil,false,nil}
  usedSource[b.Source],usedTarget[b.Target]=true,true
  child,err:=p.Field(b.Source,b.Target)
  if err!=nil{return nil,false,nil}
  children=append(children,child)
 }
 return children,true,nil
}
