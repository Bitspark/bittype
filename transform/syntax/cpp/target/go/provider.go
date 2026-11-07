package target

import (
	"context"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// Builder is the native syntax API exposed to independently supplied mappings.
type Builder[F any] interface {
	m.Builder[F, s.Language]
	Path([]string, []g.Fragment[s.Language], ast.NodeRef, bool) (g.Fragment[s.Language], error)
	Include(string) error
}
type Provider[F any] struct{ Namespace string }

func (p Provider[F]) Open(ctx context.Context, r g.TargetRequest) (Builder[F], error) {
	info, err := r.Info()
	if err != nil {
		return nil, err
	}
	return New[F](ctx, info.To, Config{TreeID: info.Output, Namespace: p.Namespace}, info.Limits)
}
