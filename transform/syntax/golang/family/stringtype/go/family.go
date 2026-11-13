// Package stringtype is the vocabulary of the family of the Go type string: the one
// home of the family's meaning (docs/contribution-contract.md, O4, O19, O24,
// O27). It names no target. The vocabulary package of each contribution of
// this family imports it, so that the meaning identity is spelled once.
//
// The family's meaning is the application's promised domain (O24). The
// accepted pages define that domain for a representation family of horos, as
// a set of canonical horos values. This source is Go, and no accepted
// page defines a family for a source that is not horos. So the meaning is
// stated here in Go's own terms: the values of the Go type string. Nothing checks the
// statement. It is one of the findings of the two other pairs (bittype#196,
// O43).
package stringtype

import g "github.com/Bitspark/bittype/contracts/generation/go"

// Meaning is the meaning identity of the family and its revision.
var Meaning = g.Meaning{Identity: "bittype.golang/string-type", Revision: "1"}
