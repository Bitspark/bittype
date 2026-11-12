// Package mapping holds the adapter through which a horos mapping contribution
// is a generation mapping library, and the identities of the two profiles that
// the horos mapping side declares.
//
// The adapter admits a profile by no list. A library answers under the profile
// that it is given, and the generation registry selects the libraries
// registered under the profile that a request names. A further profile needs
// its own libraries and its consumer's composition, and no change here.
//
// The adapter refuses one profile, EnforcementProfile, which this package
// declares and does not implement.
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

// RepresentationProfile is the structural representation profile. It promises
// native types that represent every value of the horos type. It enforces no
// kind constraint: a native generic parameter carries no bound.
const RepresentationProfile = "structural-representation/1"

// EnforcementProfile promises the same representation and native enforcement
// of kind constraints. It is declared and not implemented. Nothing that is
// built on the common horos mapping policy enforces a kind constraint, so
// Library refuses to answer under this identity.
const EnforcementProfile = "structural-enforcement/1"

// Library is a horos mapping contribution as a generation mapping library. The
// description that DescribeRules returns is published under Profile.
type Library[T any] struct {
	Profile       string
	DescribeRules func(representation.Language, representation.Language) (g.MappingDescription[h.Language, T, source.Reader], error)
}

// Describe returns the library's rules under its profile. Under
// EnforcementProfile it fails with generation.ErrUnsupported before it runs a
// rule, so representation output never carries the enforcement identity.
func (p Library[T]) Describe(from, to representation.Language) (g.MappingDescription[h.Language, T, source.Reader], error) {
	if p.Profile == EnforcementProfile {
		return g.MappingDescription[h.Language, T, source.Reader]{}, g.ErrUnsupported
	}
	if p.DescribeRules == nil {
		return g.MappingDescription[h.Language, T, source.Reader]{}, g.ErrInvalid
	}
	d, err := p.DescribeRules(from, to)
	d.Profile = p.Profile
	return d, err
}
