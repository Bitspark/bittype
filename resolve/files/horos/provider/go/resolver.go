// Package resolve composes the fixed Horos semantic frontend.
package resolve

import (
	"context"

	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	bridge "github.com/Bitspark/bittype/internal/horosmodel"
	n "github.com/Bitspark/bittype/semantics/horos/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
)

var frontend = bridge.Frontend[s.Language]{
	DeclarationPrefix: n.DeclarationTagPrefix, DescriptorPrefix: n.DescriptorTagPrefix,
	ResolveModel: func(root core.Value, deps []core.Value, limits decl.Limits) (bridge.Program, error) {
		return decl.ResolveWithLimits(root, deps, limits)
	},
	PinModel: decl.PinWithLimits,
}

// bittype:scope docs/resolution.md#horos-composition
func ResolveReferences(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	return frontend.ResolveReferences(ctx, request)
}

// Provider supplies this role independently of other language capabilities.
type Provider struct{ s.Language }

func (Provider) Resolve(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	return ResolveReferences(ctx, request)
}

var _ resolution.Provider[s.Language] = Provider{}
