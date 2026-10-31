package integration_test

import (
 "errors"
 "fmt"
 "reflect"
 "sync"
 "testing"

 "github.com/Bitspark/bittype/ast/go"
 "github.com/Bitspark/bittype/engine/go"
 "github.com/Bitspark/bittype/language/go"
 "github.com/Bitspark/bittype/transform/go"
)

// These languages and mappings belong only to independent tests. No production
// language transformer, canonical model or registry entry supplies their names.
type amber struct{}
type blue struct{}
type coral struct{}
type dusk struct{}
type echo struct{}

type testChecker[L any] struct { id language.ID; visit func(ast.AST[L]) error }
func (p testChecker[L]) ID() language.ID { return p.id }
func (p testChecker[L]) Name() string { return "test "+string(p.id) }
func (p testChecker[L]) Check(v ast.View[L]) error {
 a:=v.ASTView()
 if p.visit != nil { return p.visit(a) }
 if a.Nodes[a.Root-1].Kind != "document" { return errors.New("unknown test grammar root") }
 return nil
}
func definition[L any](t *testing.T,id language.ID) transform.Language {
 t.Helper()
 result,err:=transform.NewLanguage[L](testChecker[L]{id:id},"test-v1")
 if err!=nil { t.Fatal(err) }
 return result
}
func exampleTree[L any]() ast.AST[L] {
 return ast.AST[L]{ID:"input",Root:1,Nodes:[]ast.Node[L]{
  {Kind:"document",Children:[]ast.Child{{Role:"declaration",Node:2}}},
  {Kind:"declaration",Attributes:attr("name","Page"),Children:[]ast.Child{{Role:"parameter",Node:3}},
   References:[]ast.Reference{{Role:"library",Target:ast.NodeRef{Tree:"external",Node:7}}},
   Origin:ast.Origin{Spans:[]ast.Span{{Source:"source",Start:2,End:8}},Nodes:[]ast.NodeRef{{Tree:"earlier",Node:4}}}},
  {Kind:"parameter",Attributes:attr("name","T")},
 }}
}
func packed[L any](t *testing.T,l transform.Language) transform.AST {
 t.Helper()
 result,err:=transform.Pack(l,exampleTree[L]())
 if err!=nil { t.Fatal(err) }
 return result
}
func mapExample[F,T any](input ast.AST[F]) (ast.AST[T],error) {
 result:=ast.AST[T]{ID:input.ID,Root:input.Root,Nodes:make([]ast.Node[T],len(input.Nodes))}
 for i,n:=range input.Nodes {
  result.Nodes[i]=ast.Node[T]{Kind:ast.Kind[T](n.Kind),Attributes:n.Attributes,Children:n.Children,References:n.References,Origin:n.Origin}
 }
 return result,nil
}
func bind[F,T any](t *testing.T,from,to transform.Language) transform.ASTTransformer {
 t.Helper()
 result,err:=transform.NewTransformer[F,T](from,to,mapExample[F,T])
 if err!=nil { t.Fatal(err) }
 return result
}
func register(t *testing.T,e *engine.Engine,languages ...transform.Language) {
 t.Helper()
 for _,l:=range languages { if err:=e.RegisterLanguage(l); err!=nil { t.Fatal(err) } }
}
func connect(t *testing.T,e *engine.Engine,edges ...transform.ASTTransformer) {
 t.Helper()
 for _,edge:=range edges { if err:=e.RegisterTransformer(edge); err!=nil { t.Fatal(err) } }
}

func TestTransformerSnapshotsAndIdentity(t *testing.T) {
 a,b:=definition[amber](t,"amber"),definition[blue](t,"blue")
 original:=exampleTree[amber]()
 saved:=exampleTree[amber]()
 envelope,err:=transform.Pack(a,original)
 if err!=nil { t.Fatal(err) }
 original.Nodes[1].Attributes[0].Text="changed"
 original.Nodes[1].References[0].Target.Node=100
 original.Nodes[1].Origin.Spans[0].Source="changed"
 original.Nodes[1].Origin.Nodes[0].Tree="changed"
 original.Nodes[0].Children[0].Node=99
 extracted,err:=transform.Unpack[amber](a,envelope)
 if err!=nil || !reflect.DeepEqual(extracted,saved) { t.Fatalf("packing aliases input: %v",err) }
 extracted.Nodes[1].Attributes[0].Text="changed again"
 again,err:=transform.Unpack[amber](a,envelope)
 if err!=nil || !reflect.DeepEqual(again,saved) { t.Fatalf("unpacking exposes private storage: %v",err) }
 stage:=bind[amber,blue](t,a,b)
 result,err:=stage.Transform(envelope)
 if err!=nil { t.Fatal(err) }
 output,err:=transform.Unpack[blue](b,result)
 expected,_:=mapExample[amber,blue](saved)
 if err!=nil || !reflect.DeepEqual(output,expected) { t.Fatalf("metadata lost: %v",err) }
 if _,err:=transform.Unpack[blue](a,envelope); !errors.Is(err,transform.ErrTypeMismatch) { t.Fatal(err) }
 sameID:=definition[blue](t,"amber")
 if _,err:=transform.Unpack[blue](sameID,envelope); !errors.Is(err,transform.ErrTypeMismatch) { t.Fatal(err) }
 otherProfile,err:=transform.NewLanguage[amber](testChecker[amber]{id:"amber"},"test-v2")
 if err!=nil { t.Fatal(err) }
 if err:=otherProfile.Check(envelope); !errors.Is(err,transform.ErrLanguageMismatch) { t.Fatal(err) }
 if _,err:=transform.NewTransformer[blue,amber](a,b,mapExample[blue,amber]); !errors.Is(err,transform.ErrTypeMismatch) { t.Fatal(err) }
}

func TestTransformerRejectsInvalidBoundaries(t *testing.T) {
 a,b:=definition[amber](t,"amber"),definition[blue](t,"blue")
 var nilChecker *testChecker[amber]
 if _,err:=transform.NewLanguage[amber](nilChecker,"v1"); !errors.Is(err,transform.ErrInvalidLanguage) { t.Fatal(err) }
 if _,err:=transform.NewLanguage[amber](testChecker[amber]{},"v1"); !errors.Is(err,transform.ErrInvalidLanguage) { t.Fatal(err) }
 if _,err:=transform.NewLanguage[amber](testChecker[amber]{id:"x"}," "); !errors.Is(err,transform.ErrInvalidLanguage) { t.Fatal(err) }
 if _,err:=transform.NewTransformer[amber,blue](a,b,nil); !errors.Is(err,transform.ErrInvalidTransformer) { t.Fatal(err) }
 var view *ast.AST[amber]
 if _,err:=transform.Pack(a,view); !errors.Is(err,transform.ErrInvalidAST) { t.Fatal(err) }
 broken:=exampleTree[amber](); broken.Root=99
 if _,err:=transform.Pack(a,broken); !errors.Is(err,transform.ErrInvalidAST) { t.Fatal(err) }
 wrongGrammar:=exampleTree[amber](); wrongGrammar.Nodes[0].Kind="wrong"
 if _,err:=transform.Pack(a,wrongGrammar); err==nil { t.Fatal("grammar check skipped") }
 if _,err:=(transform.AST{}).Snapshot(); !errors.Is(err,transform.ErrInvalidAST) { t.Fatal(err) }
 input:=packed[amber](t,a)
 failed:=errors.New("unsupported declaration")
 stage,err:=transform.NewTransformer[amber,blue](a,b,func(ast.AST[amber])(ast.AST[blue],error){ return exampleTree[blue](),failed })
 if err!=nil { t.Fatal(err) }
 output,err:=stage.Transform(input)
 if !errors.Is(err,failed) || output.LanguageID()!="" { t.Fatal("partial transformer output") }
 stage,err=transform.NewTransformer[amber,blue](a,b,func(ast.AST[amber])(ast.AST[blue],error){ result:=exampleTree[blue](); result.Root=99; return result,nil })
 if err!=nil { t.Fatal(err) }
 if output,err:=stage.Transform(input); err==nil || output.LanguageID()!="" { t.Fatal("invalid stage output accepted") }
}

func TestTransformerCallbacksCannotMutateEnvelopes(t *testing.T) {
 checker:=testChecker[amber]{id:"amber",visit:func(a ast.AST[amber])error{ a.Nodes[1].Attributes[0].Text="checker mutation"; return nil }}
 a,err:=transform.NewLanguage[amber](checker,"v1")
 if err!=nil { t.Fatal(err) }
 b:=definition[blue](t,"blue")
 input:=packed[amber](t,a)
 var retained ast.AST[blue]
 adapter,err:=transform.NewTransformer[amber,blue](a,b,func(source ast.AST[amber])(ast.AST[blue],error){
  source.Nodes[1].Attributes[0].Text="stage mutation"
  retained,_=mapExample[amber,blue](source)
  return retained,nil
 })
 if err!=nil { t.Fatal(err) }
 result,err:=adapter.Transform(input)
 if err!=nil { t.Fatal(err) }
 retained.Nodes[1].Attributes[0].Text="later mutation"
 original,_:=transform.Unpack[amber](a,input)
 output,_:=transform.Unpack[blue](b,result)
 if original.Nodes[1].Attributes[0].Text!="Page" || output.Nodes[1].Attributes[0].Text!="stage mutation" { t.Fatal("callback storage escaped") }
}

type directTransformer struct {
 from,to language.ID
 run func(transform.AST)(transform.AST,error)
}
func (d *directTransformer) FromID() language.ID { return d.from }
func (d *directTransformer) ToID() language.ID { return d.to }
func (d *directTransformer) Transform(input transform.AST)(transform.AST,error) { return d.run(input) }

func TestEngineRegistration(t *testing.T) {
 e:=engine.New()
 a,b:=definition[amber](t,"amber"),definition[blue](t,"blue")
 if err:=e.RegisterLanguage(transform.Language{}); !errors.Is(err,transform.ErrInvalidLanguage) { t.Fatal(err) }
 if err:=e.RegisterTransformer(bind[amber,blue](t,a,b)); !errors.Is(err,engine.ErrUnknownLanguage) { t.Fatal(err) }
 register(t,e,a,b)
 if err:=e.RegisterLanguage(definition[blue](t,"amber")); !errors.Is(err,engine.ErrDuplicateLanguage) { t.Fatal(err) }
 var nilTransformer *directTransformer
 if err:=e.RegisterTransformer(nilTransformer); !errors.Is(err,transform.ErrInvalidTransformer) { t.Fatal(err) }
 connect(t,e,bind[amber,blue](t,a,b))
 if err:=e.RegisterTransformer(bind[amber,blue](t,a,b)); !errors.Is(err,engine.ErrDuplicateTransformer) { t.Fatal(err) }
 if err:=e.RegisterTransformer(bind[amber,amber](t,a,a)); !errors.Is(err,transform.ErrInvalidTransformer) { t.Fatal(err) }
 wrong,err:=transform.NewLanguage[blue](testChecker[blue]{id:"blue"},"different")
 if err!=nil { t.Fatal(err) }
 if err:=e.RegisterTransformer(bind[blue,amber](t,wrong,a)); !errors.Is(err,transform.ErrLanguageMismatch) { t.Fatal(err) }
 languages:=e.Languages()
 if len(languages)!=2 || languages[0].ID()!="amber" || languages[1].ID()!="blue" { t.Fatal("unstable registry metadata") }
 languages[0]=transform.Language{}
 if l,ok:=e.Language("amber"); !ok || l.Name()=="" { t.Fatal("registry slice escaped") }
 if _,ok:=e.Language("absent"); ok { t.Fatal("invented language") }
}

func TestEngineRoutesAndExecution(t *testing.T) {
 a,b,c,d,f:=definition[amber](t,"amber"),definition[blue](t,"blue"),definition[coral](t,"coral"),definition[dusk](t,"dusk"),definition[echo](t,"echo")
 e:=engine.New()
 register(t,e,a,b,c,d,f)
 connect(t,e,bind[amber,blue](t,a,b),bind[blue,dusk](t,b,d),bind[dusk,amber](t,d,a))
 plan,err:=e.Plan("amber","dusk")
 if err!=nil || !reflect.DeepEqual(plan.Languages(),[]language.ID{"amber","blue","dusk"}) { t.Fatalf("route: %v %v",plan.Languages(),err) }
 path:=plan.Languages(); path[0]="changed"
 if plan.FromID()!="amber" || plan.ToID()!="dusk" { t.Fatal("plan storage escaped") }
 input:=packed[amber](t,a)
 output,err:=e.Execute(input,plan)
 if err!=nil || output.LanguageID()!="dusk" { t.Fatalf("execution: %v",err) }
 got,_:=transform.Unpack[dusk](d,output); want,_:=mapExample[amber,dusk](exampleTree[amber]())
 if !reflect.DeepEqual(got,want) { t.Fatal("route lost metadata") }
 if _,err:=e.Plan("amber","echo"); !errors.Is(err,engine.ErrNoPath) { t.Fatal(err) }
 if _,err:=e.Plan("unknown","amber"); !errors.Is(err,engine.ErrUnknownLanguage) { t.Fatal(err) }
 if _,err:=e.PlanVia("amber","blue","amber"); !errors.Is(err,engine.ErrInvalidRoute) { t.Fatal(err) }
 if _,err:=e.PlanVia(); !errors.Is(err,engine.ErrInvalidRoute) { t.Fatal(err) }
 if _,err:=e.Execute(input,engine.Plan{}); !errors.Is(err,engine.ErrInvalidRoute) { t.Fatal(err) }
 if _,err:=e.PlanVia("amber","dusk"); !errors.Is(err,engine.ErrNoPath) { t.Fatal(err) }
 identity,err:=e.TransformTo(input,"amber")
 if err!=nil { t.Fatal(err) }
 identical,_:=transform.Unpack[amber](a,identity)
 if !reflect.DeepEqual(identical,exampleTree[amber]()) { t.Fatal("identity changed syntax") }
 connect(t,e,bind[amber,coral](t,a,c),bind[coral,dusk](t,c,d))
 if _,err:=e.Plan("amber","dusk"); !errors.Is(err,engine.ErrAmbiguousPath) { t.Fatal(err) }
 explicit,err:=e.PlanVia("amber","coral","dusk")
 if err!=nil { t.Fatal(err) }
 if _,err:=e.Execute(input,explicit); err!=nil { t.Fatal(err) }
 connect(t,e,bind[amber,dusk](t,a,d))
 direct,err:=e.Plan("amber","dusk")
 if err!=nil || len(direct.Languages())!=2 { t.Fatal("direct route was not preferred") }
 // Previously inspected explicit plans remain valid after additive registration.
 if _,err:=e.Execute(input,plan); err!=nil { t.Fatal(err) }
}

func TestEngineAmbiguityPropagatesThroughSharedTail(t *testing.T) {
 e:=engine.New()
 a,b,c,d,f:=definition[amber](t,"a"),definition[blue](t,"b"),definition[coral](t,"c"),definition[dusk](t,"d"),definition[echo](t,"f")
 register(t,e,a,b,c,d,f)
 connect(t,e,bind[amber,blue](t,a,b),bind[amber,coral](t,a,c),bind[blue,dusk](t,b,d),bind[coral,dusk](t,c,d),bind[dusk,echo](t,d,f))
 if _,err:=e.Plan("a","f"); !errors.Is(err,engine.ErrAmbiguousPath) { t.Fatal(err) }
}

func TestEngineChecksCustomTransformersAndDoesNotFallback(t *testing.T) {
 a,b,c:=definition[amber](t,"a"),definition[blue](t,"b"),definition[coral](t,"c")
 e:=engine.New(); register(t,e,a,b,c)
 input,good,wrong:=packed[amber](t,a),packed[blue](t,b),packed[coral](t,c)
 rejected:=errors.New("requires explicit mapping")
 calls:=0
 edge:=&directTransformer{from:"a",to:"b",run:func(transform.AST)(transform.AST,error){ calls++; return good,rejected }}
 connect(t,e,edge,bind[amber,coral](t,a,c),bind[coral,blue](t,c,b))
 output,err:=e.TransformTo(input,"b")
 var step *engine.StepError
 if !errors.Is(err,rejected) || !errors.As(err,&step) || step.Index!=0 || step.From!="a" || step.To!="b" || output.LanguageID()!="" || calls!=1 { t.Fatalf("failure boundary: %v",err) }
 edge.run=func(transform.AST)(transform.AST,error){ return wrong,nil }
 if _,err:=e.TransformTo(input,"b"); !errors.Is(err,transform.ErrLanguageMismatch) { t.Fatal(err) }
 edge.run=func(transform.AST)(transform.AST,error){ return transform.AST{},nil }
 if _,err:=e.TransformTo(input,"b"); !errors.Is(err,transform.ErrInvalidAST) { t.Fatal(err) }
 sameID:=definition[coral](t,"b"); wrong=packed[coral](t,sameID)
 edge.run=func(transform.AST)(transform.AST,error){ return wrong,nil }
 if _,err:=e.TransformTo(input,"b"); !errors.Is(err,transform.ErrTypeMismatch) { t.Fatal(err) }
 edge.to="changed"
 if _,err:=e.TransformTo(input,"b"); !errors.Is(err,transform.ErrInvalidTransformer) { t.Fatal(err) }
 foreign:=engine.New(); register(t,foreign,definition[coral](t,"a"),b)
 foreignPlan,err:=foreign.Plan("a","a"); if err!=nil { t.Fatal(err) }
 if _,err:=foreign.Execute(input,foreignPlan); !errors.Is(err,transform.ErrTypeMismatch) { t.Fatal(err) }
}

func TestEngineConcurrentRegistrationAndReentry(t *testing.T) {
 e:=engine.New(); a,b:=definition[amber](t,"a"),definition[blue](t,"b"); register(t,e,a,b)
 good:=packed[blue](t,b)
 connect(t,e,&directTransformer{from:"a",to:"b",run:func(transform.AST)(transform.AST,error){
  _=e.Languages()
  if _,ok:=e.Language("b"); !ok { return transform.AST{},errors.New("missing reentrant lookup") }
  return good,nil
 }})
 input:=packed[amber](t,a)
 var wg sync.WaitGroup
 failures:=make(chan error,32)
 for i:=0;i<16;i++ {
  wg.Add(1); go func(i int){ defer wg.Done()
   name:=language.ID(fmt.Sprintf("extra-%d",i))
   l,err:=transform.NewLanguage[echo](testChecker[echo]{id:name},"test-v1")
   if err==nil { err=e.RegisterLanguage(l) }
   if err!=nil { failures<-err; return }
   for j:=0;j<5;j++ { if _,err:=e.TransformTo(input,"b"); err!=nil { failures<-err; return } }
  }(i)
 }
 wg.Wait(); close(failures)
 for err:=range failures { t.Error(err) }
 if len(e.Languages())!=18 { t.Fatal("lost registrations") }
}
