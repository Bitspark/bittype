package golangsem

import (
	"context"
	"encoding/json"
	goast "go/ast"
	"go/token"
	"go/types"
	"sort"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// Checked compiler objects never cross the public semantic result boundary.
type Checked struct { Facts semantics.Data; Types map[semantics.TermID]types.Type }
type exporter struct {
	ctx context.Context
	data semantics.Data
	objects map[types.Object]semantics.EntityID
	types map[types.Type]semantics.TermID
	native map[semantics.TermID]types.Type
	positions map[token.Pos]ast.NodeRef
	err error
	depth int
}
func (e *exporter) term(t semantics.Term) semantics.TermID {
	if e.err!=nil { return 0 }
	if len(e.data.Terms)+len(e.data.Entities)+len(e.data.Facts)>32768 { e.err=semantic.ErrLimit; return 0 }
	e.data.Terms=append(e.data.Terms,t); return semantics.TermID(len(e.data.Terms))
}
func (e *exporter) expression(kind string,data any,operands ...semantics.TermID) semantics.Term {
	encoded,_:=json.Marshal(data)
	t:=semantics.Term{Kind:semantics.ProviderExpression,Operator:n.Key(kind),Encoding:n.Key("json-1"),Data:encoded}
	for i,id:=range operands { t.Operands=append(t.Operands,semantics.Operand{Locator:[]uint32{uint32(i)},Term:id}) }
	return t
}
func (e *exporter) fact(id semantics.TermID,relation string,value semantics.TermID) {
	e.data.Facts=append(e.data.Facts,semantics.Fact{Subject:semantics.Subject{Term:id},Relation:n.Key(relation),Term:value})
}
func (e *exporter) object(o types.Object) semantics.EntityID {
	if id:=e.objects[o]; id!=0 { return id }
	if e.err!=nil { return 0 }
	kind:="field"; owner:=semantics.EntityID(0)
	if _,ok:=o.(*types.TypeName); ok {
		kind="builtin"
		switch t:=o.Type().(type) {
		case *types.Named: kind="defined"; _=t
		case *types.Alias: kind="alias"
		case *types.TypeParam: kind="parameter"
		}
		if o.Parent()!=nil && o.Pkg()!=nil && o.Parent()==o.Pkg().Scope() { owner=1 }
	}
	id:=semantics.EntityID(len(e.data.Entities)+1)
	e.objects[o]=id
	entity:=semantics.Entity{Kind:n.Key(kind),Name:o.Name(),Owner:owner}
	if site:=e.positions[o.Pos()]; site.Tree!="" { entity.Declarations=[]ast.NodeRef{site} }
	e.data.Entities=append(e.data.Entities,entity)
	body:=e.typ(o.Type())
	e.data.Entities[id-1].Body=body
	return id
}
func (e *exporter) parameters(owner semantics.EntityID,params *types.TypeParamList) {
	for i:=0;i<params.Len();i++ {
		id:=e.object(params.At(i).Obj()); if id==0 { return }
		e.data.Entities[id-1].Owner=owner
		e.data.Entities[owner-1].Parameters=append(e.data.Entities[owner-1].Parameters,semantics.Parameter{Entity:id,Sort:n.Key("type")})
	}
}
func (e *exporter) typ(t types.Type) semantics.TermID {
	if id:=e.types[t]; id!=0 { return id }
	if e.err!=nil { return 0 }
	if err:=e.ctx.Err(); err!=nil { e.err=err; return 0 }
	e.depth++; defer func(){e.depth--}()
	if e.depth>128 { e.err=semantic.ErrLimit; return 0 }
	id:=e.term(semantics.Term{}); if id==0 { return 0 }
	e.types[t]=id; e.native[id]=t
	var term semantics.Term
	switch t:=t.(type) {
	case *types.Basic:
		term=e.expression("basic",t.Kind())
	case *types.Named:
		if t.TypeArgs().Len()>0 {
			term=semantics.Term{Kind:semantics.Application,Callee:e.typ(t.Origin())}
			for i:=0;i<t.TypeArgs().Len();i++ { term.Arguments=append(term.Arguments,e.typ(t.TypeArgs().At(i))) }
		} else {
			obj:=e.object(t.Obj()); term=semantics.Term{Kind:semantics.Reference,Target:obj}
			if obj!=0 { e.parameters(obj,t.TypeParams()) }
		}
	case *types.Alias:
		if t.TypeArgs().Len()>0 {
			term=semantics.Term{Kind:semantics.Application,Callee:e.typ(t.Origin())}
			for i:=0;i<t.TypeArgs().Len();i++ { term.Arguments=append(term.Arguments,e.typ(t.TypeArgs().At(i))) }
		} else {
			obj:=e.object(t.Obj()); term=semantics.Term{Kind:semantics.Reference,Target:obj}
			if obj!=0 { e.parameters(obj,t.TypeParams()) }
		}
		e.fact(id,"alias-target",e.typ(t.Rhs()))
	case *types.TypeParam:
		term=semantics.Term{Kind:semantics.Reference,Target:e.object(t.Obj())}
		e.fact(id,"constraint",e.typ(t.Constraint()))
	case *types.Slice:
		term=e.expression("slice",nil,e.typ(t.Elem()))
	case *types.Struct:
		var fields []n.FieldInfo; var operands []semantics.TermID
		for i:=0;i<t.NumFields();i++ {
			f:=t.Field(i); fields=append(fields,n.FieldInfo{Name:f.Name(),Tag:t.Tag(i),Embedded:f.Embedded(),Entity:e.object(f)}); operands=append(operands,e.typ(f.Type()))
		}
		term=e.expression("struct",fields,operands...)
	case *types.Interface:
		t.Complete(); var operands []semantics.TermID
		for i:=0;i<t.NumEmbeddeds();i++ { operands=append(operands,e.typ(t.EmbeddedType(i))) }
		for i:=0;i<t.NumExplicitMethods();i++ { operands=append(operands,e.typ(t.ExplicitMethod(i).Type())) }
		term=e.expression("interface",struct{Embeddeds int; Methods []string; Comparable bool; MethodSet bool}{Embeddeds:t.NumEmbeddeds(),Methods:methodNames(t),Comparable:t.IsComparable(),MethodSet:t.IsMethodSet()},operands...)
		e.fact(id,"empty",e.term(e.expression("boolean",t.Empty())))
	case *types.Union:
		var tilde []bool; var operands []semantics.TermID
		for i:=0;i<t.Len();i++ { tilde=append(tilde,t.Term(i).Tilde()); operands=append(operands,e.typ(t.Term(i).Type())) }
		term=e.expression("union",tilde,operands...)
	case *types.Signature:
		// The predeclared error interface includes a method signature even though
		// authored methods and function syntax remain outside this input profile.
		term=e.expression("signature",t.Variadic(),e.typ(t.Params()),e.typ(t.Results()))
	case *types.Tuple:
		var operands []semantics.TermID
		for i:=0;i<t.Len();i++ { operands=append(operands,e.typ(t.At(i).Type())) }
		term=e.expression("tuple",nil,operands...)
	default:
		e.err=n.ErrUnsupported
	}
	e.data.Terms[id-1]=term
	if e.err!=nil { return id }
	e.fact(id,"underlying",e.typ(t.Underlying()))
	e.fact(id,"display",e.term(semantics.Term{Kind:semantics.ProviderExpression,Operator:n.Key("display"),Encoding:n.Key("utf8"),Data:[]byte(types.TypeString(t,nil))}))
	return id
}
func methodNames(t *types.Interface) []string { var out []string; for i:=0;i<t.NumExplicitMethods();i++ { out=append(out,t.ExplicitMethod(i).Name()) }; return out }
func siteLess(a,b ast.NodeRef) bool { if a.Tree!=b.Tree { return a.Tree<b.Tree }; return a.Node<b.Node }
func export(ctx context.Context,pkg *types.Package,info *types.Info,origins map[*goast.Ident]ast.NodeRef,expressions map[goast.Expr]ast.NodeRef,positions map[token.Pos]ast.NodeRef,inputs []s.AST) (*Checked,error) {
	e:=exporter{ctx:ctx,objects:map[types.Object]semantics.EntityID{},types:map[types.Type]semantics.TermID{},native:map[semantics.TermID]types.Type{},positions:positions}
	packageEntity:=semantics.Entity{Kind:n.Key("package"),Name:pkg.Path()}
	for _,tree:=range inputs { packageEntity.Declarations=append(packageEntity.Declarations,ast.NodeRef{Tree:tree.ID,Node:tree.Root}) }
	e.data.Entities=append(e.data.Entities,packageEntity)
	// Seed public scalar queries even when neither type is mentioned in syntax.
	e.typ(types.Typ[types.Bool]); e.typ(types.Typ[types.String])
	var definitions []*goast.Ident
	for id,obj:=range info.Defs { if obj!=nil { definitions=append(definitions,id) } }
	sort.Slice(definitions,func(i,j int) bool { return siteLess(origins[definitions[i]],origins[definitions[j]]) })
	for _,id:=range definitions { e.object(info.Defs[id]) }
	var uses []*goast.Ident
	for id:=range info.Uses { uses=append(uses,id) }
	sort.Slice(uses,func(i,j int) bool { return siteLess(origins[uses[i]],origins[uses[j]]) })
	for _,id:=range uses { e.data.Bindings=append(e.data.Bindings,semantics.Binding{At:semantics.Occurrence{Site:origins[id],Role:n.Key("name")},Target:e.object(info.Uses[id])}) }
	var exprs []goast.Expr
	for expr:=range info.Types { if _,ok:=expressions[expr]; ok { exprs=append(exprs,expr) } }
	sort.Slice(exprs,func(i,j int) bool { return siteLess(expressions[exprs[i]],expressions[exprs[j]]) })
	for _,expr:=range exprs { e.data.Expressions=append(e.data.Expressions,semantics.Expression{At:semantics.Occurrence{Site:expressions[expr],Role:n.Key("type-expression")},Term:e.typ(info.Types[expr].Type)}) }
	if e.err!=nil { return nil,e.err }
	// Identity is established by go/types, including aliases and structurally
	// identical unnamed types. It is scoped to this exported graph.
	var representatives []semantics.TermID
	for id:=semantics.TermID(1);int(id)<=len(e.data.Terms);id++ {
		t,ok:=e.native[id]; if !ok { continue }
		if err:=ctx.Err(); err!=nil { return nil,err }
		canonical:=id
		for _,other:=range representatives { if types.Identical(t,e.native[other]) { canonical=other; break } }
		if canonical==id { representatives=append(representatives,id) }
		e.fact(id,"identical",canonical)
	}
	return &Checked{Facts:e.data,Types:e.native},nil
}
