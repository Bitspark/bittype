// Package mapping holds the identities that the Go mapping side declares: the
// one generation profile of a pair whose source is Go, and the publisher of
// the contributions that state it. It holds identities only and imports
// nothing.
//
// It exists for the two other pairs of bittype#196 (decision O43), which show
// that the contribution contract and its kernel assume no horos source. Go as
// a source is not supported beyond what the pair golang-to-typescript covers.
package mapping

// TypesProfile is the profile of the contributions whose source is Go. It
// promises a target type that has one member for each field of a Go struct
// type, in the order of the source, with the field's name as the source spells
// it and with the target type selected for the field's type. It promises
// nothing about values: no binding is defined for a pair whose source is not
// horos, because a binding is defined between canonical horos values and
// native values.
const TypesProfile = "golang-declared-types/1"

// Publisher names the publisher of the contributions that state TypesProfile:
// the library that publishes them, which is bittype (row 15). The name
// authenticates nothing: until a composition lock binds it (row 15a),
// registration in one process is the binding.
//
// COPY (AC-11, finding AG-6 of the two other pairs): the same name is spelled
// in transform/syntax/horos/mapping/go, because the vocabulary of a
// contribution may import the mapping identities of its own source language
// only. The publisher is bittype whatever the source language is: one name
// with a spelling per source language, and no agnostic home.
const Publisher = "github.com/Bitspark/bittype"
