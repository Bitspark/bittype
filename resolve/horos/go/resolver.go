// Package resolve composes the fixed horos semantic frontend.
package resolve

import (
	"context"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	bridge "github.com/Bitspark/bittype/internal/horosmodel"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
)

const Profile = "horos-declarations-3/resolution-1"

var frontend = bridge.Frontend[s.Language]{
	DeclarationPrefix: "horos/decl/", DescriptorPrefix: "horos/", Owner: "horos", ResolutionProfile: Profile,
	ResolveModel: func(root core.Value, deps []core.Value, limits decl.Limits) (bridge.Program, error) {
		return decl.ResolveWithLimits(root, deps, limits)
	},
	PinModel: decl.PinWithLimits,
}

func ReferenceResolution() semantics.Key { return frontend.ReferenceResolution() }
func Validity() semantics.Key            { return frontend.Validity() }
func Indexing() semantics.Key            { return frontend.Indexing() }
func Bindings() semantics.Key            { return frontend.Bindings() }
func Applications() semantics.Key        { return frontend.Applications() }
func Dependencies() semantics.Key        { return frontend.Dependencies() }
func TypeFacts() semantics.Key           { return frontend.TypeFacts() }
func TargetRole() semantics.Key          { return frontend.TargetRole() }
func ExpressionRole() semantics.Key      { return frontend.ExpressionRole() }
func ReferenceDependency() semantics.Key { return frontend.ReferenceDependency() }
func KindRelation() semantics.Key        { return frontend.KindRelation() }
func InterfaceRelation() semantics.Key   { return frontend.InterfaceRelation() }
func DataStatusRelation() semantics.Key  { return frontend.DataStatusRelation() }
func DataSummaryRelation() semantics.Key { return frontend.DataSummaryRelation() }
func IdentityRelation() semantics.Key    { return frontend.IdentityRelation() }
func LimitsSetting() semantics.Key       { return frontend.LimitsSetting() }
func Contract() resolution.Contract      { return frontend.Contract() }
func ResolveReferences(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	return frontend.ResolveReferences(ctx, request)
}
