package generationfixture

import (
	"context"
	c "github.com/Bitspark/bittype/generation/syntax/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type catalog struct { f *Fixture }
func(catalog catalog) Declarations() []g.Declaration[SourceLanguage] { var out []g.Declaration[SourceLanguage];for _,p:=range catalog.f.config.Plans{if p.Declaration.Valid(){out=append(out,p.Declaration)}};return out }
func(v catalog) DeclarationSelector(d g.Declaration[SourceLanguage])(c.Selector,bool,error){for _,p:=range v.f.config.Plans{if p.Declaration==d{key,err:=c.NewSelector(v.f.source.snapshot.Language(),p.Key);return key,true,err}};return c.Selector{},false,g.ErrScope}
func(v catalog) UseSelector(use g.UseData[SourceLanguage])(c.Selector,error){name,err:=v.f.source.selectPlan(use);if err!=nil{return c.Selector{},err};return c.NewSelector(v.f.source.snapshot.Language(),name)}

// Automatic exercises the production composer with an unrelated source language,
// independent association owners, an external generic and a native-only helper.
func(f *Fixture) Automatic(ctx context.Context)(representation.ASTSet,error){
	var bindings []c.Binding[SourceLanguage,s.Language];var helpers []g.Plan[SourceLanguage,s.Language]
	for _,plan:=range f.config.Plans{
		if plan.Key=="cursor-helper"{helpers=append(helpers,plan);continue}
		key,err:=c.NewSelector(f.source.snapshot.Language(),plan.Key);if err!=nil{return representation.ASTSet{},err}
		bindings=append(bindings,c.Binding[SourceLanguage,s.Language]{Selector:key,Mapping:c.Mapping[SourceLanguage,s.Language]{ID:plan.Key,PluginID:"independent-fixture",Plan:func(d g.Declaration[SourceLanguage])(g.Plan[SourceLanguage,s.Language],error){if d!=plan.Declaration{return g.Plan[SourceLanguage,s.Language]{},g.ErrScope};return plan,nil}}})
	}
	table,err:=c.NewBindings(f.source.snapshot.Language(),bindings,nil,nil);if err!=nil{return representation.ASTSet{},err}
	config:=f.Config();config.Plans=nil;config.Select=nil
	return c.Run(ctx,f.Input(),c.Composition[SourceLanguage,s.Language]{Config:config,Catalog:catalog{f},Bindings:table,Roots:f.Roots(),Helpers:helpers,Limits:g.DefaultLimits()})
}
