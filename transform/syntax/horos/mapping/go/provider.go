package generation

import (
	"context"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// Rule opens one Horos mapping for an explicitly selected representation profile.
type Rule[T any, N Native[T]] struct {
	profile string
	open    func(source.Reader, N) (g.Contribution[h.Language, T], error)
}

func (p Rule[T, N]) Open(ctx context.Context, r g.MappingRequest[source.Reader, N]) (g.Contribution[h.Language, T], error) {
	zero := g.Contribution[h.Language, T]{}
	src, native, err := r.Opened()
	if err != nil {
		return zero, err
	}
	if p.profile != RepresentationProfile {
		return zero, g.ErrUnsupported
	}
	if p.open == nil {
		return zero, g.ErrInvalid
	}
	if err = ctx.Err(); err != nil {
		return zero, err
	}
	return p.open(src, native)
}
func BuiltinRule[T any, N Native[T]](profile, head, discriminator string, build func(source.Reader, N) g.Mapping[h.Language, T]) Rule[T, N] {
	return Rule[T, N]{profile, func(src source.Reader, native N) (g.Contribution[h.Language, T], error) {
		key, err := source.Builtin(src.Snapshot().Language(), head, discriminator)
		if err != nil {
			return g.Contribution[h.Language, T]{}, err
		}
		return g.Contribution[h.Language, T]{Bindings: []g.RepresentationBinding[h.Language, T]{{Selector: key, Mapping: build(src, native)}}}, nil
	}}
}
func BinderRule[T any, N Native[T]](profile string) Rule[T, N] {
	return Rule[T, N]{profile, func(src source.Reader, _ N) (g.Contribution[h.Language, T], error) {
		key, err := source.Bound(src.Snapshot().Language())
		if err != nil {
			return g.Contribution[h.Language, T]{}, err
		}
		return g.Contribution[h.Language, T]{Bindings: []g.RepresentationBinding[h.Language, T]{{Selector: key, Mapping: Bound[T]()}}}, nil
	}}
}
func DeclarationRule[T any, N Native[T]](profile string) Rule[T, N] {
	return Rule[T, N]{profile, func(src source.Reader, native N) (g.Contribution[h.Language, T], error) {
		m := Structural[T](src, native)
		return g.Contribution[h.Language, T]{Structural: &m}, nil
	}}
}
func HelperRule[T any, N Native[T]](profile string, build func(N) g.Plan[h.Language, T]) Rule[T, N] {
	return Rule[T, N]{profile, func(_ source.Reader, native N) (g.Contribution[h.Language, T], error) {
		return g.Contribution[h.Language, T]{Helpers: []g.Plan[h.Language, T]{build(native)}}, nil
	}}
}
