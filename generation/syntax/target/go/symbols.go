package target

import (
	"context"
	"slices"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// UsesFormal reads actual symbolic syntax. It grants no source membership or
// representation theorem; an owning mapping decides what this observation means.
func (s *Store[F,T]) UsesFormal(f g.Fragment[T],l g.Layout[T],index int) bool {
	seen:=map[g.Fragment[T]]bool{}
	var visit func(g.Fragment[T])bool
	visit=func(f g.Fragment[T])bool{if seen[f]{return false};seen[f]=true;r,ok:=s.fragmentRecords[f];if !ok{return false};if r.symbol.layout==l&&r.symbol.index==index&&r.referenceRole!=""{return true};for _,c:=range r.children{if visit(c.Fragment){return true}};return false};return visit(f)
}

func (s *Store[F,T]) Reserve(site g.Site[F,T],name string,origin ast.NodeRef,build func([]g.Fragment[T],g.Fragment[T]) (g.Fragment[T],error)) (g.Destination[T],bool,error) {
	if !site.Valid()||!site.Signature().Valid()||name==""||build==nil||s.occupied[name]{return g.Destination[T]{},false,g.ErrInvalid}
	if len(s.destinationRecords)>=s.limits.Reservations{return g.Destination[T]{},false,g.ErrLimit}
	s.occupied[name]=true;d,_:=s.destinations.Ref(uint64(len(s.destinationRecords)+1));s.destinationRecords[d]=destination[T]{identity:at(site),name:name,origin:origin,build:build};return d,true,nil
}
func (s *Store[F,T]) Freeze(site g.Site[F,T],d g.Destination[T],names []string,parameters []g.Fragment[T]) (g.Layout[T],error) {
	r,ok:=s.destinationRecords[d];if !ok||r.identity!=at(site)||s.headers[d].Valid()||len(names)!=len(parameters)||len(names)>s.limits.Formals{return g.Layout[T]{},g.ErrScope}
	seen:=map[string]bool{};for i,name:=range names{if name==""||seen[name]||s.occupied[name]{return g.Layout[T]{},g.ErrInvalid};if _,ok:=s.fragmentRecords[parameters[i]];!ok{return g.Layout[T]{},g.ErrScope};seen[name]=true}
	l,_:=s.layouts.Ref(uint64(len(s.layoutRecords)+1));s.layoutRecords[l]=layout[T]{identity:at(site),destination:d,names:slices.Clone(names),parameters:slices.Clone(parameters)};s.headers[d]=l;return l,nil
}
func (s *Store[F,T]) Formal(l g.Layout[T],index int) (Symbol[T],error) {
	r,ok:=s.layoutRecords[l];if !ok||index<0||index>=len(r.names){return Symbol[T]{},g.ErrScope};return Symbol[T]{layout:l,index:index,name:r.names[index]},nil
}
func (s *Store[F,T]) Named(c g.Callee[F,T]) (Symbol[T],error) {
	d,ok:=s.destinationRecords[c.Destination()];if !ok||s.headers[c.Destination()]!=c.Layout(){return Symbol[T]{},g.ErrScope};return Symbol[T]{destination:c.Destination(),name:d.name},nil
}
func (s *Store[F,T]) symbol(symbol Symbol[T]) (identity,error) {
	if symbol.layout.Valid(){l,ok:=s.layoutRecords[symbol.layout];if !ok||symbol.index<0||symbol.index>=len(l.names)||symbol.name!=l.names[symbol.index]{return identity{},g.ErrScope};return l.identity,nil}
	d,ok:=s.destinationRecords[symbol.destination];if !ok||symbol.name!=d.name{return identity{},g.ErrScope};return d.identity,nil
}
func (s *Store[F,T]) CheckDestination(ctx context.Context,site g.Site[F,T],d g.Destination[T],owned bool) error {
	if err:=ctx.Err();err!=nil{return err};r,ok:=s.destinationRecords[d];if !ok||r.identity!=at(site)||!owned{return g.ErrScope};return nil
}
func (s *Store[F,T]) CheckLayout(ctx context.Context,site g.Site[F,T],d g.Destination[T],l g.Layout[T]) error {
	if err:=ctx.Err();err!=nil{return err};r,ok:=s.layoutRecords[l];if !ok||r.identity!=at(site)||r.destination!=d||s.headers[d]!=l{return g.ErrScope};return nil
}
func (s *Store[F,T]) visible(ctx context.Context,f g.Fragment[T],in identity) error {
	seen:=map[g.Fragment[T]]bool{}
	var walk func(g.Fragment[T],int) error
	walk=func(f g.Fragment[T],depth int) error {
		if err:=ctx.Err();err!=nil{return err};if depth>s.limits.Depth{return g.ErrLimit};if seen[f]{return nil};seen[f]=true
		r,ok:=s.fragmentRecords[f];if !ok{return g.ErrScope}
		if r.referenceRole!=""{id,err:=s.symbol(r.symbol);if err!=nil{return err};if id.run!=in.run||id.config!=in.config||r.symbol.IsFormal()&&(id.frame!=in.frame||id.signature!=in.signature){return g.ErrScope}}
		for _,c:=range r.children{if err:=walk(c.Fragment,depth+1);err!=nil{return err}};return nil
	};return walk(f,0)
}
func (s *Store[F,T]) CheckProduct(ctx context.Context,w g.Want[F,T],f g.Fragment[T]) error { return s.visible(ctx,f,at(w.Site())) }
func (s *Store[F,T]) CheckBody(ctx context.Context,c g.Callee[F,T],f g.Fragment[T]) error {
	l,ok:=s.layoutRecords[c.Layout()];if !ok||l.destination!=c.Destination()||l.identity.frame!=c.Frame(){return g.ErrScope};return s.visible(ctx,f,l.identity)
}
func (s *Store[F,T]) CheckExposure(ctx context.Context,site g.Site[F,T],root g.Declaration[F],c g.Callee[F,T]) error {
	if err:=ctx.Err();err!=nil{return err};d,ok:=s.destinationRecords[c.Destination()];if !ok||d.identity.run!=site.Run()||d.identity.config!=site.Config()||root!=c.Open().Declaration{return g.ErrScope};return nil
}
