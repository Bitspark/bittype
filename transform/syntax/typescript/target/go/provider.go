package target

import (
	"context"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
)

// Builder is the native syntax API exposed to independently supplied mappings.
type Builder[F any] interface {
	m.Builder[F, s.Language]
}
type Provider[F any] struct{}

func (p Provider[F]) Open(ctx context.Context, r g.TargetRequest) (Builder[F], error) {
	info, err := r.Info()
	if err != nil {
		return nil, err
	}
	return New[F](ctx, info.To, Config{TreeID: info.Output}, info.Limits)
}
