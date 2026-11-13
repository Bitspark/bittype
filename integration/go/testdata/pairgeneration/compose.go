// Package pairgeneration is the consumer composition of the two other pairs
// of bittype#196 (decision O43): golang-to-typescript and typescript-to-horos.
// It registers a source, a target and the bundles of one ordered pair and
// returns the request. It names no language: the source type, the target type
// and the reader type are its type parameters.
//
// The composition helper of the horos pairs, testdata/horosgeneration.Compose,
// does the same steps with the source fixed to horos. It is not changed here.
// That is finding AG-5 of the two other pairs (AC-11): registering a source,
// a target, each contribution and each codec for one request needs no
// knowledge of a language.
package pairgeneration

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Bundle is what a consumer's composition registers for one contribution: the
// contribution and the codec of its payload schema. Each is one package, and
// only a composition imports both. Codec is nil for a contribution that
// states a schema whose codec another bundle registers.
type Bundle[F, T, R any] struct {
	Contribution func(from representation.Language) (g.Contribution[F, T, R], error)
	Codec        func() (g.PayloadCodec, error)
}

// Pair is one ordered pair as a composition states it: the source, the
// target, the profile that its contributions state, and the bundles.
type Pair[F, T, R any] struct {
	Source  g.SourceProvider[F, R]
	Target  g.Target[T]
	Profile string
	Bundles []Bundle[F, T, R]
}

// Compose registers the pair in a fresh registry and returns the request for
// the given roots. Each bundle's codec is registered once: a registry holds
// one codec for each schema reference.
func Compose[F, T, R any](registry *c.Registry, pair Pair[F, T, R], input transform.Input, to representation.Language, output ast.TreeID, roots []ast.NodeRef) (c.Request, error) {
	zero := c.Request{}
	from := input.Source.Language()
	if err := c.RegisterSource[F, R](registry, from, pair.Source); err != nil {
		return zero, err
	}
	if err := c.RegisterTarget(registry, to, pair.Target); err != nil {
		return zero, err
	}
	for _, bundle := range pair.Bundles {
		contribution, err := bundle.Contribution(from)
		if err != nil {
			return zero, err
		}
		if err := c.RegisterContribution(registry, from, to, contribution); err != nil {
			return zero, err
		}
		if bundle.Codec == nil {
			continue
		}
		codec, err := bundle.Codec()
		if err != nil {
			return zero, err
		}
		if err := registry.RegisterCodec(codec); err != nil {
			return zero, err
		}
	}
	return c.Request{Input: input, To: to, Output: output, Roots: roots, Profile: pair.Profile, Limits: g.DefaultLimits()}, nil
}
