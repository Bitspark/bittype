// Package composition is what a consumer's composition does for any pair of
// a source and a target: it holds the contributions, payload codecs, facets
// and products that the consumer selects, and registers them. It names no
// language. A fixture of one source or one pair instantiates it with its own
// types and adds what is its own: its source provider, its roots, its profile.
package composition

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Bundle is what a publisher offers for one contribution: the contribution
// and the codec of its payload schema, each by the constructor of its own
// package. Only a composition imports both.
type Bundle[F, T, R any] struct {
	Contribution func(from representation.Language) (g.Contribution[F, T, R], error)
	Codec        func() (g.PayloadCodec, error)
}

// Set is what a consumer registers for one request.
type Set[F, T, R any] struct {
	Contributions []g.Contribution[F, T, R]
	Codecs        []g.PayloadCodec
	Facets        []g.Facet
	Products      []g.Product
}

// Build makes the contribution and the codec of each bundle.
func Build[F, T, R any](from representation.Language, bundles ...Bundle[F, T, R]) (Set[F, T, R], error) {
	var set Set[F, T, R]
	for _, bundle := range bundles {
		contribution, err := bundle.Contribution(from)
		if err != nil {
			return Set[F, T, R]{}, err
		}
		codec, err := bundle.Codec()
		if err != nil {
			return Set[F, T, R]{}, err
		}
		set.Contributions, set.Codecs = append(set.Contributions, contribution), append(set.Codecs, codec)
	}
	return set, nil
}

// Register registers the set, as a consumer does: each contribution, each
// codec, each facet and each product. The first failure is returned. A
// registry holds one codec for each schema reference, so a set names each
// codec once.
func (s Set[F, T, R]) Register(registry *c.Registry, from, to representation.Language) error {
	for _, contribution := range s.Contributions {
		if err := c.RegisterContribution(registry, from, to, contribution); err != nil {
			return err
		}
	}
	for _, codec := range s.Codecs {
		if err := registry.RegisterCodec(codec); err != nil {
			return err
		}
	}
	for _, facet := range s.Facets {
		if err := registry.RegisterFacet(facet); err != nil {
			return err
		}
	}
	for _, product := range s.Products {
		if err := registry.RegisterProduct(product); err != nil {
			return err
		}
	}
	return nil
}

// Select returns the entries of a table of bundles under the given names, in
// that order. Registration order is no preference: the order only makes a
// test's composition readable.
func Select[B any](table map[string]B, names ...string) []B {
	out := make([]B, 0, len(names))
	for _, name := range names {
		out = append(out, table[name])
	}
	return out
}
