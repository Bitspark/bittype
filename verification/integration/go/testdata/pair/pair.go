// Package pair is a consumer-owned fixture. Its two native arguments occur in
// the reverse order of the Horos tuple's operands.
package pair

import (
 "github.com/Bitspark/bittype/semantics/go"
 n "github.com/Bitspark/bittype/semantics/golang/go"
 v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
 decl "github.com/Bitspark/horos/core/go/declarations"
)
const Law="example/pair-reversed-1"

func Rule(helper n.Type)v.Rule{
 return v.Rule{Key:semantics.Key{Owner:"example",Name:"tuple-pair"},Dependency:Law,Project:func(p v.Subject)([]v.Subject,bool,error){
  if p.Optional() || p.Source().Kind()!=decl.Descriptor || p.Source().Constructor()!="tuple" || len(p.Source().Operands())!=2 || !helper.Same(p.Target().Origin()){return nil,false,nil}
  params,args:=helper.Parameters(),p.Target().Arguments()
  fields,ok:=helper.Fields()
  if !ok || len(params)!=2 || len(args)!=2 || !params[0].Unconstrained() || !params[1].Unconstrained() || len(fields)!=2{return nil,false,nil}
  for i,name:=range []string{"First","Second"}{
   f:=fields[i]
   if f.Name!=name || f.Embedded || f.Tag!="" || !f.Type.Same(params[1-i]){return nil,false,nil}
  }
  left,err:=p.Operand(0,args[1])
  if err!=nil{return nil,false,err}
  right,err:=p.Operand(1,args[0])
  if err!=nil{return nil,false,err}
  return []v.Subject{left,right},true,nil
 }}
}
