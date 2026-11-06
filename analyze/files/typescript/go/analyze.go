// Package analyze consumes exact TypeScript resolution receipts.
package analyze

import (
	"context"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	resolve "github.com/Bitspark/bittype/resolve/files/typescript/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
)

func Contract() analysis.Contract {
	c:=resolve.Contract().Semantic; c.Profile=n.AnalysisProfile
	return analysis.Contract{Semantic:c,Resolution:&resolution.Requirement{Profile:n.ResolutionProfile,Configuration:n.Configuration()}}
}
func Analyze(ctx context.Context,request analysis.Request[s.Language]) (semantic.Result,error) {
	if ctx==nil || request.Contract().Profile!=n.AnalysisProfile || !n.MatchesConfiguration(request.Configuration()) { return semantic.Result{},semantic.ErrInvalidRequest }
	if err:=ctx.Err(); err!=nil { return semantic.Result{},err }
	resolved,ok:=request.Resolution(); if !ok { return semantic.Result{},semantic.ErrInvalidRequest }
	publication:=semantic.Publication{Diagnostics:resolved.Diagnostics()}; known:=map[semantics.Key]bool{}
	for _,obligation:=range Contract().Semantic.Obligations {
		report,found:=resolved.Report(obligation.Key); if !found { report=semantic.Report{Obligation:obligation.Key,Status:semantic.Incomplete} }
		publication.Reports=append(publication.Reports,report); known[obligation.Key]=true
	}
	for _,key:=range request.Requested() { if !known[key] { publication.Reports=append(publication.Reports,semantic.Report{Obligation:key,Status:semantic.Unsupported}) } }
	if resolved.HasSemantics() { publication.ReuseSemanticsFrom=&resolved }; return request.Publish(publication)
}
type Provider struct { s.Language }
func (Provider) Analyze(ctx context.Context,request analysis.Request[s.Language]) (semantic.Result,error) { return Analyze(ctx,request) }
var _ analysis.Provider[s.Language]=Provider{}
