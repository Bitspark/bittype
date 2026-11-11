// Package mapping holds the adapter through which a horos mapping contribution
// is a generation mapping library, and the names of the two profiles. The
// adapter owns one decision: which profile a library answers under. It refuses
// to describe a library under any profile but RepresentationProfile.
//
// The source patterns, the declaration recipe and the naming function that a
// contribution builds its rules from are separate decisions and live in the
// sibling packages pattern, declarations and naming.
package mapping

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

const RepresentationProfile = "pagination-representation/1"
const EnforcementProfile = "pagination-enforcement/1"

type Library[T any] struct {
	Profile       string
	DescribeRules func(representation.Language, representation.Language) (g.MappingDescription[h.Language, T, source.Reader], error)
}

func (p Library[T]) Describe(from, to representation.Language) (g.MappingDescription[h.Language, T, source.Reader], error) {
	if p.Profile != RepresentationProfile {
		return g.MappingDescription[h.Language, T, source.Reader]{}, g.ErrUnsupported
	}
	if p.DescribeRules == nil {
		return g.MappingDescription[h.Language, T, source.Reader]{}, g.ErrInvalid
	}
	d, err := p.DescribeRules(from, to)
	d.Profile = p.Profile
	return d, err
}
