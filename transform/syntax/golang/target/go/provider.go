package target

import (
	"context"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
)

// Builder is the native syntax API exposed to independently supplied mappings.
type Builder[F any] interface {
	m.Builder[F, s.Language]
}
type Provider[F any] struct{ Package string }

func (p Provider[F]) Open(ctx context.Context, r g.TargetRequest) (Builder[F], error) {
	info, err := r.Info()
	if err != nil {
		return nil, err
	}
	return New[F](ctx, info.To, Config{TreeID: info.Output, Package: p.Package}, info.Limits)
}
