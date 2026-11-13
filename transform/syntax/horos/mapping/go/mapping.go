// Package mapping holds the identities of the two generation profiles that the
// horos mapping side declares, and the function type through which a horos
// mapping library is given as its description function.
//
// A profile is a promise about the native type, so the mapping library that
// makes the promise states it. Each builtin structural mapping writes
// RepresentationProfile into its own description, and its constructor takes no
// profile. Nothing here writes a profile, compares one or refuses one. The
// generation registry selects the libraries registered under the profile that
// a request names. A contribution of another profile states its own identity
// in its own package, and nothing here changes.
//
// The source patterns, the declaration recipe and the naming function that a
// contribution builds its rules from live in the sibling packages pattern,
// declarations and naming.
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
// of kind constraints. It is declared and not implemented: no mapping library
// in this module states it, so the registry finds no mapping for a request
// that names it and fails with generation.ErrUnsupported.
const EnforcementProfile = "structural-enforcement/1"

// Publisher names the publisher of the contributions of the structural
// profiles: the library that publishes them, which is bittype (row 15). The
// representation owner of the structural families is "bittype's structural
// profile" (docs/representation-families.md), so the package that declares
// that profile declares its publisher. The vocabulary package of each
// contribution states it in the contribution's identity. The name
// authenticates nothing: until a composition lock binds it (row 15a),
// registration in one process is the binding.
const Publisher = "github.com/Bitspark/bittype"

// ReasonNestedRecord is the structural profile's reason for refusing a record
// that is no body of a named declaration: "Records must be named
// declarations" (docs/semantic-generation.md). Each target's record
// contribution refuses such a record as unsupported with this reason.
const ReasonNestedRecord = "nested records require a named declaration in this profile"

// Library is a horos mapping library given as its description function. The
// function is the whole library: it states the profile, the bindings, the
// declaration recipes and the helpers. The type adds one thing, the method
// that makes the function a generation.MappingLibrary.
type Library[T any] func(from, to representation.Language) (g.MappingDescription[h.Language, T, source.Reader], error)

// Describe returns what the function returns, unchanged. A nil function fails
// with generation.ErrInvalid.
func (p Library[T]) Describe(from, to representation.Language) (g.MappingDescription[h.Language, T, source.Reader], error) {
	if p == nil {
		return g.MappingDescription[h.Language, T, source.Reader]{}, g.ErrInvalid
	}
	return p(from, to)
}
