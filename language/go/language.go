// Package language defines language identity for syntax providers.
//
// It contains no language list or concrete syntax imports. The agnostic engine
// accepts explicitly configured providers supplied by its composition root.
// No registry or automatic registration is implemented here.
package language

// ID is an exact, case-sensitive registration key for a language family.
// An implementation must return a stable, nonempty value. IDs are metadata:
// equal IDs do not establish Go type identity or grammar compatibility.
type ID string

// LanguageProvider identifies a syntax provider without depending on its concrete type.
// Configured instances return stable, nonempty metadata without IO or registration
// side effects. Name is for display only. Native marker providers have usable
// zero values; configured registration descriptors require explicit construction.
//
// This interface intentionally supplies no AST construction, validation, parsing,
// transformation or formatting operations. It does not erase or convert AST[L].
type LanguageProvider interface {
	ID() ID
	Name() string
}
