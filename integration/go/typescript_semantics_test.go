package integration_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/files/typescript/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/engine/go"
	"github.com/Bitspark/bittype/files/go"
	p "github.com/Bitspark/bittype/provider/typescript/go"
	resolve "github.com/Bitspark/bittype/resolve/files/typescript/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func tsCompiler(t *testing.T) resolve.Compiler {
	t.Helper()
	c:=resolve.Compiler{NodePath:os.Getenv("BITTYPE_TYPESCRIPT_NODE"),TypeScriptPath:os.Getenv("BITTYPE_TYPESCRIPT_MODULE")}
	if err:=c.Validate(); err!=nil { t.Fatal("Provision the pinned TypeScript compiler and set BITTYPE_TYPESCRIPT_NODE and BITTYPE_TYPESCRIPT_MODULE:",err) }
	return c
}
func tsCapture(t *testing.T,texts ...string) (semantic.Snapshot,[]s.AST,*files.Context) {
	t.Helper(); var trees []s.AST; source:=&files.Context{}
	for i,text:=range texts {
		id:=fmt.Sprint(i); parsed,err:=(p.Provider{}).Parse(context.Background(),parse.Input{Text:text,SourceID:ast.SourceID(id),TreeID:ast.TreeID(id)})
		if err!=nil { t.Fatal(err) }; trees=append(trees,parsed.ASTView())
		source.Files.Entries=append(source.Files.Entries,files.File{ID:files.ID(id),Path:files.Path(id+".ts"),Text:text})
		source.Links=append(source.Links,files.SyntaxLink{File:files.ID(id),Tree:ast.TreeID(id),Source:ast.SourceID(id)})
	}
	return tsCaptureTrees(t,trees,source),trees,source
}
func tsCaptureTrees(t *testing.T,trees []s.AST,source *files.Context) semantic.Snapshot {
	t.Helper(); definition,err:=representation.NewLanguage[s.Language](s.Language{},"typescript-syntax-1"); if err!=nil { t.Fatal(err) }
	checker,err:=p.NewSyntaxChecker(definition); if err!=nil { t.Fatal(err) }; var envelopes []representation.AST
	for _,tree:=range trees { env,err:=representation.Pack[s.Language](definition,tree); if err!=nil { t.Fatal(err) }; envelopes=append(envelopes,env) }
	snapshot,err:=semantic.Capture(context.Background(),checker,semantic.Input{Entry:ast.NodeRef{Tree:trees[0].ID,Node:trees[0].Root},Trees:envelopes,Files:source}); if err!=nil { t.Fatal(err) }; return snapshot
}
func tsResults(t *testing.T,snapshot semantic.Snapshot) (semantic.Result,semantic.Result) {
	t.Helper(); resolver,err:=p.NewFilesResolver(snapshot.Language(),tsCompiler(t)); if err!=nil { t.Fatal(err) }
	resolved,err:=resolver.Resolve(context.Background(),snapshot,semantic.Options{Configuration:n.Configuration()}); if err!=nil { t.Fatal(err) }
	analyzer,err:=p.NewFilesAnalyzer(snapshot.Language()); if err!=nil { t.Fatal(err) }
	result,err:=analyzer.Analyze(context.Background(),analysis.Input{Snapshot:snapshot,Resolution:&resolved,Options:semantic.Options{Configuration:n.Configuration()}}); if err!=nil { t.Fatal(err) }
	return resolved,result
}
func tsModel(t *testing.T,texts ...string) n.Model {
	t.Helper(); snapshot,_,_:=tsCapture(t,texts...); _,result:=tsResults(t,snapshot); model,err:=n.Open(result)
	if err!=nil { t.Fatalf("model: %v; diagnostics: %+v",err,result.Diagnostics()) }; return model
}
func tsDecl(t *testing.T,m n.Model,tree ast.TreeID,name string) n.Declaration {
	t.Helper(); d,err:=m.Lookup(tree,name); if err!=nil { t.Fatal(name,err) }; return d
}
func tsType(t *testing.T,m n.Model,name string) n.Type { t.Helper(); return tsDecl(t,m,"0",name).Type() }

func qualifyTypeScript(t *testing.T) {
	ctx:=context.Background(); snapshot,_,_:=tsCapture(t,"export interface Person {name: string} export interface Page<T> {items: T[]; next?: string} export type People = Page<Person>;")
	var e engine.Engine; if err:=e.RegisterLanguage(snapshot.Language()); err!=nil { t.Fatal(err) }
	if _,err:=e.Resolve(ctx,snapshot,semantic.Options{}); !errors.Is(err,engine.ErrMissingCapability) { t.Fatal(err) }
	resolver,err:=p.NewFilesResolver(snapshot.Language(),tsCompiler(t)); if err!=nil { t.Fatal(err) }; if err=e.RegisterResolver(resolver); err!=nil { t.Fatal(err) }
	resolved,err:=e.Resolve(ctx,snapshot,semantic.Options{Configuration:n.Configuration()}); if err!=nil { t.Fatal(err) }
	input:=analysis.Input{Snapshot:snapshot,Resolution:&resolved,Options:semantic.Options{Configuration:n.Configuration()}}
	if _,err=e.Analyze(ctx,input); !errors.Is(err,engine.ErrMissingCapability) { t.Fatal(err) }
	analyzer,err:=p.NewFilesAnalyzer(snapshot.Language()); if err!=nil { t.Fatal(err) }; if err=e.RegisterAnalyzer(analyzer); err!=nil { t.Fatal(err) }
	result,err:=e.Analyze(ctx,input); if err!=nil { t.Fatal(err) }
	rg,_:=resolved.Semantics(); ag,_:=result.Semantics(); if !rg.Same(ag) { t.Fatal("analysis rebuilt the compiler model",result.Diagnostics()) }
	m,err:=n.Open(result); if err!=nil { t.Fatal(err,result.Diagnostics()) }
	if _,err=n.Open(resolved); err==nil { t.Fatal("resolution became analysis evidence") }
	if len(m.Graph().Data().Bindings)==0 || len(m.Graph().Data().Expressions)==0 { t.Fatal("missing native occurrences") }
	page:=tsDecl(t,m,"0","Page"); if len(page.Parameters())!=1 || page.Parameters()[0].Entity().Owner!=page.ID() { t.Fatal("lost binder ownership") }
	people:=tsType(t,m,"People"); fields,ok:=people.Properties(); if !ok || len(fields)!=2 || !fields[1].Optional { t.Fatal(fields) }
	if fields[0].Type.String()!="Person[]" { t.Fatal("lost substitution",fields[0].Type.String()) }
	if pairs,ok:=m.Graph().DirectParameterBindings(people.ID()); !ok || len(pairs)!=1 || pairs[0].Parameter!=page.Parameters()[0].ID() { t.Fatal("generic binding",pairs,ok) }
	data:=m.Graph().Data(); data.Entities[0].Name="mutated"; fields[0].Sites=nil
	if m.Graph().Data().Entities[0].Name=="mutated" { t.Fatal("mutable graph") }
}

func TestTypeScriptSemanticGenericsAndNativeFacts(t *testing.T) {
	m:=tsModel(t,`export type Text = string;
export type Alias = Text;
export interface Person { readonly name: string; age?: number }
export interface Page<T> { items: T[]; next?: string }
export type Nested = Page<Page<Person>>;
export interface Named<T extends {name: string}> { value: T }
export type Good = Named<Person>;
export interface Other<T> { value: T }
export interface Recursive { children: Recursive[] }
export interface Derived extends Person { active: boolean }
export type Combined = Person & { active: boolean };
export type Choice = "a" | "b";
export type Pair = readonly [string, number?];
export type ReadonlyItems = readonly Person[];
export type StandardItems = Array<Person>;
export type StandardReadonly = ReadonlyArray<Person>;`)
	if !tsType(t,m,"Text").Same(tsType(t,m,"Alias")) || tsDecl(t,m,"0","Text").ID()==tsDecl(t,m,"0","Alias").ID() { t.Fatal("alias declaration vs native identity") }
	page,named,other:=tsDecl(t,m,"0","Page"),tsDecl(t,m,"0","Named"),tsDecl(t,m,"0","Other")
	if page.Parameters()[0].Type().Same(other.Parameters()[0].Type()) || page.Parameters()[0].Constraint().Valid() || !named.Parameters()[0].Constraint().Valid() { t.Fatal("parameter identity/constraint") }
	fields,ok:=tsType(t,m,"Person").Properties(); if !ok || !fields[0].Readonly || !fields[1].Optional || !strings.Contains(fields[1].Type.String(),"undefined") { t.Fatal(fields) }
	fields,ok=tsType(t,m,"Nested").Properties(); if !ok || fields[0].Type.String()!="Page<Person>[]" { t.Fatal("nested substitution",fields) }
	fields,ok=tsType(t,m,"Derived").Properties(); if !ok || len(fields)!=3 { t.Fatal("inherited properties",fields) }
	// Recursive definitions retain a named anchor in the admitted acyclic graph.
	if tsType(t,m,"Recursive").Term().Kind!=semantics.Reference { t.Fatal("lost recursive identity") }
}

func TestTypeScriptSemanticReports(t *testing.T) {
	for _,tc:=range []struct{name,text string; coverage semantic.Completion; validity semantic.Verdict}{
		{"missing-name","type Box = Missing;",semantic.Complete,semantic.Rejected},
		{"constraint","type Box<T extends string> = {value:T}; type Bad = Box<boolean>;",semantic.Complete,semantic.Rejected},
		{"arity","type Box<T> = {value:T}; type Bad = Box;",semantic.Complete,semantic.Rejected},
		{"duplicate","type Box = {value:string; value:boolean};",semantic.Complete,semantic.Rejected},
		{"default","type Box<T = string> = {value:T};",semantic.Unsupported,semantic.Undetermined},
		{"import","import {Box} from 'ambient'; type X = Box;",semantic.Unsupported,semantic.Undetermined},
		{"function","type F = (x: string) => boolean;",semantic.Unsupported,semantic.Undetermined},
		{"method","interface Box { value(): string }",semantic.Unsupported,semantic.Undetermined},
		{"keyof","type Box = {value:string}; type Key = keyof Box;",semantic.Unsupported,semantic.Undetermined},
		{"merging","interface Box {x:string} interface Box {y:string}",semantic.Unsupported,semantic.Undetermined},
		{"external","type DateValue = Date;",semantic.Unsupported,semantic.Undetermined},
	} { t.Run(tc.name,func(t *testing.T){ snapshot,_,_:=tsCapture(t,tc.text); resolved,result:=tsResults(t,snapshot)
		if resolved.HasSemantics() || result.HasSemantics() { t.Fatal("partial graph") }; coverage,_:=result.Report(n.Coverage()); validity,_:=result.Report(n.Validity())
		if coverage.Status!=tc.coverage || validity.Verdict!=tc.validity { t.Fatal(coverage,validity,result.Diagnostics()) }
		if len(result.Diagnostics())==0 { t.Fatal("missing diagnostic") }; if _,err:=n.Open(result); err==nil { t.Fatal("invalid model admitted") }
	}) }
}

func TestTypeScriptSemanticScopeAndSourceCorrespondence(t *testing.T) {
	m:=tsModel(t,"export interface Box<T> {value:T}","export interface Box<T> {value:T}")
	a,b:=tsDecl(t,m,"0","Box"),tsDecl(t,m,"1","Box"); if a.ID()==b.ID() || a.Parameters()[0].ID()==b.Parameters()[0].ID() { t.Fatal("module identities collapsed") }
	global:=tsModel(t,"type People = Person[];", "interface Person {name:string}")
	if !strings.Contains(tsType(t,global,"People").String(),"People") { t.Fatal("cross-file script scope") }
	// The leading non-BMP character exercises UTF-8 -> UTF-16 correspondence.
	snapshot,trees,source:=tsCapture(t,"// 😀\nexport type Text = string; export type Alias = Text;")
	_,result:=tsResults(t,snapshot); if !result.HasSemantics() { t.Fatal(result.Diagnostics()) }
	resolver,_:=p.NewFilesResolver(snapshot.Language(),tsCompiler(t))
	changed:=source.Clone(); changed.Files.Entries[0].Text=strings.ReplaceAll(changed.Files.Entries[0].Text,"string","boolean")
	forged:=tsCaptureTrees(t,trees,&changed)
	if _,err:=resolver.Resolve(context.Background(),forged,semantic.Options{Configuration:n.Configuration()}); !errors.Is(err,semantic.ErrInvalidRequest) { t.Fatal("changed source",err) }
	for i,node:=range trees[0].Nodes { if node.Kind==s.KindTypeReference { trees[0].Nodes[i].References=[]ast.Reference{{Role:s.RoleTarget,Target:ast.NodeRef{Tree:trees[0].ID,Node:trees[0].Root}}}; break } }
	forged=tsCaptureTrees(t,trees,source)
	if _,err:=resolver.Resolve(context.Background(),forged,semantic.Options{Configuration:n.Configuration()}); !errors.Is(err,semantic.ErrInvalidRequest) { t.Fatal("forged binding",err) }
	// Origin text is not the evidence: original source is reparsed independently.
	for i:=range trees[0].Nodes { trees[0].Nodes[i].References=nil; trees[0].Nodes[i].Origin=ast.Origin{} }
	_,result=tsResults(t,tsCaptureTrees(t,trees,source)); if !result.HasSemantics() { t.Fatal("trusted caller origin",result.Diagnostics()) }
}

func TestTypeScriptSemanticReceiptsAndLimits(t *testing.T) {
	snapshot,trees,source:=tsCapture(t,"type Box<T> = {value:T};")
	resolved,result:=tsResults(t,snapshot); analyzer,_:=p.NewFilesAnalyzer(snapshot.Language()); resolver,_:=p.NewFilesResolver(snapshot.Language(),tsCompiler(t))
	for name,input:=range map[string]analysis.Input{
		"missing":{Snapshot:snapshot,Options:semantic.Options{Configuration:n.Configuration()}},
		"foreign":{Snapshot:tsCaptureTrees(t,trees,source),Resolution:&resolved,Options:semantic.Options{Configuration:n.Configuration()}},
		"wrong-role":{Snapshot:snapshot,Resolution:&result,Options:semantic.Options{Configuration:n.Configuration()}},
	} { t.Run(name,func(t *testing.T){ if r,err:=analyzer.Analyze(context.Background(),input); err==nil || r.Valid() { t.Fatal("bad receipt",err) } }) }
	missing:=tsCaptureTrees(t,trees,nil); _,r:=tsResults(t,missing); coverage,_:=r.Report(n.Coverage()); if r.HasSemantics() || coverage.Status!=semantic.Incomplete { t.Fatal("missing source context",r.Diagnostics()) }
	for i:=range n.Configuration() { settings:=n.Configuration(); settings[i].Value=[]byte("foreign"); if _,err:=resolver.Resolve(context.Background(),snapshot,semantic.Options{Configuration:settings}); !errors.Is(err,semantic.ErrInvalidRequest) { t.Fatal("foreign environment",err) } }
	if _,err:=p.NewFilesResolver(snapshot.Language(),resolve.Compiler{}); err==nil { t.Fatal("implicit compiler") }
	canceled,cancel:=context.WithCancel(context.Background()); cancel(); if _,err:=resolver.Resolve(canceled,snapshot,semantic.Options{Configuration:n.Configuration()}); !errors.Is(err,context.Canceled) { t.Fatal(err) }
	var texts []string; for i:=0;i<17;i++ { texts=append(texts,fmt.Sprintf("type X%d = string;",i)) }; large,_,_:=tsCapture(t,texts...); _,r=tsResults(t,large); coverage,_=r.Report(n.Coverage()); if coverage.Status!=semantic.Incomplete || r.HasSemantics() { t.Fatal("input bounds") }
	unknown:=n.Key("unknown-query"); r,err:=resolver.Resolve(context.Background(),snapshot,semantic.Options{Configuration:n.Configuration(),Requested:[]semantics.Key{unknown}}); if err!=nil { t.Fatal(err) }; report,_:=r.Report(unknown); if report.Status!=semantic.Unsupported { t.Fatal(report) }
}

func TestTypeScriptSemanticAssignability(t *testing.T) {
	m:=tsModel(t,`type Narrow = {x:string; y:boolean}; type Wide = {x:string};
type Optional = {x?:string}; type ExplicitUndefined = {x:string|undefined};
type Letters = "a" | "b"; type One = "a"; type Text = string;`)
	q,err:=analyze.NewQueries(m,tsCompiler(t)); if err!=nil { t.Fatal(err) }
	for _,tc:=range []struct{from,to string; want bool}{{"Narrow","Wide",true},{"Wide","Narrow",false},{"ExplicitUndefined","Optional",false},{"One","Letters",true},{"Text","Letters",false}} {
		from,to:=tsType(t,m,tc.from),tsType(t,m,tc.to); a,err:=q.Assignable(context.Background(),from,to)
		if err!=nil || a.Holds()!=tc.want || !a.Model().Same(m) || !a.Source().Same(from) || !a.Target().Same(to) { t.Fatal(tc,a.Holds(),err) }
	}
	foreign:=tsModel(t,"type Text = string;"); if _,err:=q.Assignable(context.Background(),tsType(t,foreign,"Text"),tsType(t,m,"Text")); !errors.Is(err,n.ErrForeign) { t.Fatal("foreign argument",err) }
	if _,err:=q.Assignable(context.Background(),n.Type{},tsType(t,m,"Text")); !errors.Is(err,n.ErrForeign) { t.Fatal("zero argument",err) }
	if _,err:=analyze.NewQueries(n.Model{},tsCompiler(t)); !errors.Is(err,n.ErrForeign) { t.Fatal("zero model",err) }
	canceled,cancel:=context.WithCancel(context.Background()); cancel(); if _,err:=q.Assignable(canceled,tsType(t,m,"Text"),tsType(t,m,"Text")); !errors.Is(err,context.Canceled) { t.Fatal(err) }
}

func TestTypeScriptSemanticConcurrentQueries(t *testing.T) {
	m:=tsModel(t,"type Box = {value:string}; type Read = {value:string};"); q,err:=analyze.NewQueries(m,tsCompiler(t)); if err!=nil { t.Fatal(err) }
	from,to:=tsType(t,m,"Box"),tsType(t,m,"Read"); original:=m.Graph().Data(); var wg sync.WaitGroup
	for i:=0;i<3;i++ { wg.Add(1); go func(){ defer wg.Done(); a,err:=q.Assignable(context.Background(),from,to); if err!=nil || !a.Holds() { t.Error(err) }; if !reflect.DeepEqual(original,m.Graph().Data()) { t.Error("mutable semantics") } }() }; wg.Wait()
}
