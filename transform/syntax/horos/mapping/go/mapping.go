// Package mapping holds the identities that the horos mapping side declares:
// the two generation profiles, the publisher of the contributions of the
// structural profile, and the profile's reason for refusing a nested record.
// It holds identities only and imports nothing.
//
// A profile is a promise about the native type, so the contribution that
// makes the promise states it. Each builtin structural contribution writes
// RepresentationProfile into its own description, and its constructor takes no
// profile. Nothing here writes a profile, compares one or refuses one. The
// generation registry selects the contributions registered under the profile
// that a request names. A contribution of another profile states its own
// identity in its own package, and nothing here changes.
//
// The declaration policy and the naming function that a contribution builds
// its rule from live in the sibling packages declarations and naming.
package mapping

// RepresentationProfile is the structural representation profile. It promises
// native types that represent every value of the horos type. It enforces no
// kind constraint: a native generic parameter carries no bound.
const RepresentationProfile = "structural-representation/1"

// EnforcementProfile promises the same representation and native enforcement
// of kind constraints. It is declared and not implemented: no contribution in
// this module states it, so the registry finds no contribution for a request
// that names it and refuses the request as a missing contribution.
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
