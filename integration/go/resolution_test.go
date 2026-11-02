package integration_test

import (
	"context"
	"errors"
	"sync"
	"testing"
 "time"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/engine/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func bindingKey() semantics.Key { return semantics.Key{Owner:"independent-resolver",Name:"references"} }
func bindingContract() resolution.Contract {
	key := bindingKey()
	return resolution.Contract{Bindings:key,Semantic:semantic.Contract{Profile:"bindings-v1",Obligations:[]semantic.Obligation{{Key:key,Domain:analysisKey("entry")}},ModelRequirements:[]semantic.Requirement{{Obligation:key,Verdict:semantic.Holds}}}}
}
func bindingReports() []semantic.Report { return []semantic.Report{{Obligation:bindingKey(),Status:semantic.Complete,Verdict:semantic.Holds}} }
func bindingFacts() semantics.Data {
	return semantics.Data{Entities:[]semantics.Entity{{Kind:analysisKey("symbol"),Name:"Target",Declarations:[]ast.NodeRef{{Tree:"input",Node:2}}}},Bindings:[]semantics.Binding{{At:semantics.Occurrence{Site:ast.NodeRef{Tree:"input",Node:3},Role:analysisKey("reference")},Target:1}}}
}
type bindingProvider[L any] struct {
	id language.ID
	run func(context.Context,semantic.Request[L])(semantic.Result,error)
}
func (p bindingProvider[L]) ID() language.ID { return p.id }
func (p bindingProvider[L]) Name() string { return "Independent binding provider" }
func (p bindingProvider[L]) Resolve(ctx context.Context,r semantic.Request[L])(semantic.Result,error) {
	if p.run != nil { return p.run(ctx,r) }
	facts:=bindingFacts()
	return r.Publish(semantic.Publication{Facts:&facts,Reports:bindingReports()})
}
func bindingFixture(t *testing.T, run func(context.Context,semantic.Request[amber])(semantic.Result,error)) (*engine.Engine,representation.Language,semantic.Snapshot,resolution.Resolver) {
	t.Helper()
	d:=definition[amber](t,"binding-language")
	e:=engine.New();register(t,e,d)
	input:=analysisFixture[amber](t,d)
	source,err:=e.Capture(context.Background(),input.source());if err!=nil {t.Fatal(err)}
	r,err:=resolution.NewResolver[amber](d,bindingContract(),bindingProvider[amber]{id:d.ID(),run:run});if err!=nil {t.Fatal(err)}
	return e,d,source,r
}
func bindingAnalysis(t *testing.T,d representation.Language, requirement resolution.Requirement, run func(context.Context,analysis.Request[amber])(semantic.Result,error)) analysis.Analyzer {
	t.Helper()
	c:=analysis.Contract{Semantic:bindingContract().Semantic,Resolution:&requirement}
	c.Semantic.Profile="binding-consumer-v1"
	if run==nil { run=func(_ context.Context,r analysis.Request[amber])(semantic.Result,error) {
		dependency,ok:=r.Resolution();if !ok{return semantic.Result{},errors.New("missing resolution")}
		return r.Publish(semantic.Publication{ReuseModelFrom:&dependency,Reports:bindingReports()})
	} }
	a,err:=analysis.NewAnalyzer[amber](d,c,analysisChecker[amber]{id:d.ID(),run:run});if err!=nil {t.Fatal(err)}
	return a
}
func bindingRequirement() resolution.Requirement { return resolution.Requirement{Profile:"bindings-v1",Reports:bindingContract().Semantic.ModelRequirements} }

func TestResolverIndependentRegistrationAndBindingOnlyModel(t *testing.T) {
	e,d,s,r:=bindingFixture(t,nil)
	if _,err:=e.Resolve(context.Background(),s,semantic.Options{});!errors.Is(err,engine.ErrMissingCapability){t.Fatal(err)}
	if err:=engine.New().RegisterResolver(r);!errors.Is(err,engine.ErrUnknownLanguage){t.Fatal(err)}
	if err:=e.RegisterResolver(resolution.Resolver{});!errors.Is(err,semantic.ErrInvalidOperation){t.Fatal(err)}
	wrong:=definition[blue](t,d.ID())
	other,err:=resolution.NewResolver[blue](wrong,bindingContract(),bindingProvider[blue]{id:d.ID()});if err!=nil{t.Fatal(err)}
	if err:=e.RegisterResolver(other);!errors.Is(err,representation.ErrLanguageMismatch){t.Fatal(err)}
	if _,err:=resolution.NewResolver[blue](d,bindingContract(),bindingProvider[blue]{id:d.ID()});!errors.Is(err,representation.ErrTypeMismatch){t.Fatal(err)}
	var absent *bindingProvider[amber]
	if _,err:=resolution.NewResolver[amber](d,bindingContract(),absent);!errors.Is(err,semantic.ErrInvalidOperation){t.Fatal(err)}
	bad:=bindingContract();bad.Semantic.ModelRequirements=nil
	if _,err:=resolution.NewResolver[amber](d,bad,bindingProvider[amber]{id:d.ID()});!errors.Is(err,semantic.ErrInvalidOperation){t.Fatal(err)}
	if err:=e.RegisterResolver(r);err!=nil{t.Fatal(err)}
	if err:=e.RegisterResolver(r);!errors.Is(err,engine.ErrDuplicateResolver){t.Fatal(err)}
	advertised,err:=e.Resolver(d.ID());if err!=nil || advertised.Contract().Bindings!=bindingKey(){t.Fatal(err)}
	result,err:=e.Resolve(context.Background(),s,semantic.Options{});if err!=nil{t.Fatal(err)}
	model,ok:=result.Model();if !ok || len(model.Data().Terms)!=0{t.Fatal("binding-only provider required expression/type analysis")}
	at:=bindingFacts().Bindings[0].At
	if id,ok:=model.Resolve(at);!ok || id!=1{t.Fatal("binding query lost its target")}
	if _,err:=e.Analyze(context.Background(),analysis.Input{Snapshot:s});!errors.Is(err,engine.ErrMissingCapability){t.Fatal("resolution registered analysis implicitly")}
}

func TestResolverInvocationIdentityWithinOneSnapshot(t *testing.T) {
	var retained semantic.Result
	e,_,s,r:=bindingFixture(t,func(_ context.Context,request semantic.Request[amber])(semantic.Result,error){
		if retained.Valid(){return retained,nil};facts:=bindingFacts();var err error
		retained,err=request.Publish(semantic.Publication{Facts:&facts,Reports:bindingReports()});return retained,err
	})
	if err:=e.RegisterResolver(r);err!=nil{t.Fatal(err)}
	if _,err:=e.Resolve(context.Background(),s,semantic.Options{});err!=nil{t.Fatal(err)}
	if result,err:=e.Resolve(context.Background(),s,semantic.Options{});!errors.Is(err,semantic.ErrStaleResult)||result.Valid(){t.Fatalf("same snapshot accepted prior invocation: %v",err)}
}

func TestResolutionCompositionChecksSnapshotProfileConfigurationAndOperation(t *testing.T) {
	e,d,s,r:=bindingFixture(t,nil)
	resolved,err:=r.Resolve(context.Background(),s,semantic.Options{});if err!=nil{t.Fatal(err)}
	a:=bindingAnalysis(t,d,bindingRequirement(),nil)
	result,err:=a.Analyze(context.Background(),analysis.Input{Snapshot:s,Resolution:&resolved});if err!=nil{t.Fatal(err)}
	before,_:=resolved.Model();after,_:=result.Model()
	if before!=after{t.Fatal("sealed model identity changed")}
	copy:=after.Data();copy.Entities[0].Name="changed"
	if entity,_:=before.Entity(1);entity.Name!="Target"{t.Fatal("model reuse aliases mutable storage")}
	fresh,err:=e.Capture(context.Background(),analysisFixture[amber](t,d).source());if err!=nil{t.Fatal(err)}
	for name,input:=range map[string]analysis.Input{
		"missing":{Snapshot:s},"zero":{Snapshot:s,Resolution:&semantic.Result{}},
		"same addresses different capture":{Snapshot:fresh,Resolution:&resolved},
		"wrong operation":{Snapshot:s,Resolution:&result},
	}{t.Run(name,func(t *testing.T){out,err:=a.Analyze(context.Background(),input);if !errors.Is(err,semantic.ErrInvalidRequest)||out.Valid(){t.Fatal(err)}})}
	for _,requirement:=range []resolution.Requirement{
		{Profile:"other-profile",Reports:bindingRequirement().Reports},
		{Profile:"bindings-v1",Configuration:[]semantic.Setting{{Key:analysisKey("mode"),Value:[]byte("other")}},Reports:bindingRequirement().Reports},
		{Profile:"bindings-v1",Reports:[]semantic.Requirement{{Obligation:analysisKey("not-advertised"),Verdict:semantic.Holds}}},
	}{a:=bindingAnalysis(t,d,requirement,nil);if out,err:=a.Analyze(context.Background(),analysis.Input{Snapshot:s,Resolution:&resolved});!errors.Is(err,semantic.ErrInvalidRequest)||out.Valid(){t.Fatal(err)}}
	standalone:=analysisAnalyzer[amber](t,d,nil)
	if _,err:=standalone.Analyze(context.Background(),analysis.Input{Snapshot:s,Resolution:&resolved});!errors.Is(err,semantic.ErrInvalidRequest){t.Fatal("undeclared dependency accepted")}
}

func TestResolutionReuseRequiresExactAdmittedReceiptAndBounds(t *testing.T) {
	_,d,s,r:=bindingFixture(t,nil)
	resolved,err:=r.Resolve(context.Background(),s,semantic.Options{});if err!=nil{t.Fatal(err)}
	unrelated,err:=r.Resolve(context.Background(),s,semantic.Options{});if err!=nil{t.Fatal(err)}
	for name,run:=range map[string]func(context.Context,analysis.Request[amber])(semantic.Result,error){
		"unadmitted receipt":func(_ context.Context,request analysis.Request[amber])(semantic.Result,error){return request.Publish(semantic.Publication{ReuseModelFrom:&unrelated,Reports:bindingReports()})},
		"new and reused":func(_ context.Context,request analysis.Request[amber])(semantic.Result,error){facts:=bindingFacts();return request.Publish(semantic.Publication{Facts:&facts,ReuseModelFrom:&resolved,Reports:bindingReports()})},
		"return child result":func(context.Context,analysis.Request[amber])(semantic.Result,error){return resolved,nil},
	}{t.Run(name,func(t *testing.T){a:=bindingAnalysis(t,d,bindingRequirement(),run);out,err:=a.Analyze(context.Background(),analysis.Input{Snapshot:s,Resolution:&resolved});if err==nil||out.Valid(){t.Fatal("invalid reuse accepted")}})}
	a:=bindingAnalysis(t,d,bindingRequirement(),nil)
	limits:=semantic.DefaultLimits();limits.Model.MaxRecords=1
	if out,err:=a.Analyze(context.Background(),analysis.Input{Snapshot:s,Resolution:&resolved,Options:semantic.Options{Limits:limits}});!errors.Is(err,semantic.ErrLimit)||out.Valid(){t.Fatalf("reused model bypassed tighter bounds: %v",err)}
}

func TestResolutionPrerequisiteFailureCannotPublishAnalysisFacts(t *testing.T) {
	_,d,s,r:=bindingFixture(t,func(_ context.Context,r semantic.Request[amber])(semantic.Result,error){
		reports:=bindingReports();reports[0].Verdict=semantic.Rejected
		return r.Publish(semantic.Publication{Reports:reports,Diagnostics:[]semantic.Diagnostic{{Code:analysisKey("unbound"),Subject:analysisKey("source"),Message:"name has no target"}}})
	})
	resolved,err:=r.Resolve(context.Background(),s,semantic.Options{});if err!=nil{t.Fatal(err)}
	a:=bindingAnalysis(t,d,bindingRequirement(),func(_ context.Context,request analysis.Request[amber])(semantic.Result,error){facts:=bindingFacts();return request.Publish(semantic.Publication{Facts:&facts,Reports:bindingReports()})})
	if out,err:=a.Analyze(context.Background(),analysis.Input{Snapshot:s,Resolution:&resolved});!errors.Is(err,semantic.ErrInvalidResult)||out.Valid(){t.Fatalf("failed prerequisite permitted new facts: %v",err)}
	a=bindingAnalysis(t,d,bindingRequirement(),func(_ context.Context,request analysis.Request[amber])(semantic.Result,error){dependency,_:=request.Resolution();return request.Publish(semantic.Publication{Reports:dependency.Reports(),Diagnostics:dependency.Diagnostics()})})
	out,err:=a.Analyze(context.Background(),analysis.Input{Snapshot:s,Resolution:&resolved})
	if err!=nil||!out.Valid()||out.HasModel()||len(out.Diagnostics())!=1{t.Fatalf("normal resolution rejection became an API error: %v",err)}
}

func TestResolutionConfigurationAndSnapshotOwnership(t *testing.T) {
	e,d,s,r:=bindingFixture(t,nil)
	settings:=[]semantic.Setting{{Key:analysisKey("mode"),Value:[]byte("original")}}
	resolved,err:=r.Resolve(context.Background(),s,semantic.Options{Configuration:settings});if err!=nil{t.Fatal(err)}
	requirement:=bindingRequirement();requirement.Configuration=settings
	a:=bindingAnalysis(t,d,requirement,nil)
	settings[0].Value[0]='X';requirement.Reports[0].Verdict=semantic.Rejected
	if _,err:=a.Analyze(context.Background(),analysis.Input{Snapshot:s,Resolution:&resolved});err!=nil{t.Fatal("consumer requirement aliases construction storage",err)}
	input:=analysisFixture[amber](t,d)
	snapshot,err:=e.Capture(context.Background(),input.source());if err!=nil{t.Fatal(err)}
	input.Trees[0]=representation.AST{};input.Entry.Node=0
	if result,err:=r.Resolve(context.Background(),snapshot,semantic.Options{});err!=nil||!result.HasModel(){t.Fatal("snapshot retained caller slice",err)}
}

func TestResolverCancellationAndErrorDiscard(t *testing.T) {
	ctx,cancel:=context.WithCancel(context.Background())
	_,_,s,r:=bindingFixture(t,func(_ context.Context,request semantic.Request[amber])(semantic.Result,error){facts:=bindingFacts();result,err:=request.Publish(semantic.Publication{Facts:&facts,Reports:bindingReports()});cancel();return result,err})
	if out,err:=r.Resolve(ctx,s,semantic.Options{});!errors.Is(err,context.Canceled)||out.Valid(){t.Fatal(err)}
	if out,err:=r.Resolve(nil,s,semantic.Options{});!errors.Is(err,semantic.ErrInvalidRequest)||out.Valid(){t.Fatal(err)}
	sentinel:=errors.New("resolver failed after publication")
	_,_,s,r=bindingFixture(t,func(_ context.Context,request semantic.Request[amber])(semantic.Result,error){facts:=bindingFacts();result,_:=request.Publish(semantic.Publication{Facts:&facts,Reports:bindingReports()});return result,sentinel})
	if out,err:=r.Resolve(context.Background(),s,semantic.Options{});!errors.Is(err,sentinel)||out.Valid(){t.Fatal(err)}
}

func TestEngineConcurrentResolutionAndAnalysis(t *testing.T) {
	e,d,s,_:=bindingFixture(t,nil)
	r,err:=resolution.NewResolver[amber](d,bindingContract(),bindingProvider[amber]{id:d.ID(),run:func(_ context.Context,request semantic.Request[amber])(semantic.Result,error){
		_ = e.Languages();if _,err:=e.Resolver(d.ID());err!=nil{return semantic.Result{},err};facts:=bindingFacts();return request.Publish(semantic.Publication{Facts:&facts,Reports:bindingReports()})
	}});if err!=nil{t.Fatal(err)}
	if err:=e.RegisterResolver(r);err!=nil{t.Fatal(err)}
	if err:=e.RegisterAnalyzer(bindingAnalysis(t,d,bindingRequirement(),nil));err!=nil{t.Fatal(err)}
	var group sync.WaitGroup
	for range 8 {group.Add(1);go func(){defer group.Done();resolved,err:=e.Resolve(context.Background(),s,semantic.Options{});if err!=nil{t.Error(err);return};out,err:=e.Analyze(context.Background(),analysis.Input{Snapshot:s,Resolution:&resolved});if err!=nil||!out.HasModel(){t.Error(err)}}()}
	done:=make(chan struct{});go func(){group.Wait();close(done)}()
 select {case <-done: case <-time.After(5*time.Second):t.Fatal("resolution/analysis callback deadlocked under registry lock")}
}
