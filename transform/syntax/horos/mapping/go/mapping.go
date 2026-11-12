// Package mapping holds the identities of the two generation profiles that the
// horos mapping side declares. It holds nothing else: no adapter, no rule and
// no list of profiles.
//
// A profile is a promise about the native type, so the mapping library that
// makes the promise states it. Each builtin structural mapping implements
// generation.MappingLibrary itself and writes RepresentationProfile into its
// own description. Its constructor takes no profile. The generation registry
// selects the libraries registered under the profile that a request names. A
// contribution of another profile states its own identity in its own package,
// and nothing here changes.
//
// The source patterns, the declaration recipe and the naming function that a
// contribution builds its rules from live in the sibling packages pattern,
// declarations and naming.
package mapping

// RepresentationProfile is the structural representation profile. It promises
// native types that represent every value of the horos type. It enforces no
// kind constraint: a native generic parameter carries no bound.
const RepresentationProfile = "structural-representation/1"

// EnforcementProfile promises the same representation and native enforcement
// of kind constraints. It is declared and not implemented: no mapping library
// in this module states it, so the registry finds no mapping for a request
// that names it and fails with generation.ErrUnsupported.
const EnforcementProfile = "structural-enforcement/1"
