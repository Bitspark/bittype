// Package resolve composes the fixed horos/wire semantic frontend.
package resolve
import (
 "context"
 "github.com/Bitspark/bittype/contracts/resolution/go"
 "github.com/Bitspark/bittype/contracts/semantic/go"
 bridge "github.com/Bitspark/bittype/internal/horosmodel"
 "github.com/Bitspark/bittype/semantics/go"
 s "github.com/Bitspark/bittype/syntax/horoswire/go"
 decl "github.com/Bitspark/horos/core/go/wire"
 core "github.com/bitspark/ontos/core/go"
)
const Profile = "horos-wire-1/resolution-1"
var frontend = bridge.Frontend[s.Language]{
 DeclarationPrefix: "horos/wire/", DescriptorPrefix: "horos/wire/", Owner: "horos/wire", ResolutionProfile: Profile,
 ResolveModel: func(root core.Value, deps []core.Value, limits decl.Limits) (bridge.Program,error) { return decl.ResolveWithLimits(root,deps,limits) },
 PinModel: decl.PinWithLimits,
}
func ReferenceResolution() semantics.Key { return frontend.ReferenceResolution() }
func Validity() semantics.Key { return frontend.Validity() }
func Indexing() semantics.Key { return frontend.Indexing() }
func Bindings() semantics.Key { return frontend.Bindings() }
func Applications() semantics.Key { return frontend.Applications() }
func Dependencies() semantics.Key { return frontend.Dependencies() }
func TargetRole() semantics.Key { return frontend.TargetRole() }
func ExpressionRole() semantics.Key { return frontend.ExpressionRole() }
func ApplicationDependency() semantics.Key { return frontend.ApplicationDependency() }
func LimitsSetting() semantics.Key { return frontend.LimitsSetting() }
func Contract() resolution.Contract { return frontend.Contract() }
func ResolveReferences(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result,error) { return frontend.ResolveReferences(ctx,request) }
