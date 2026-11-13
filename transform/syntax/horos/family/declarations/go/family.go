// Package declarations is the vocabulary of the family of the declaration contribution: the one
// home of the family's meaning (docs/contribution-contract.md, O4, O19, O24,
// O27). It names no target. The vocabulary package of each target's
// contribution of this family imports it, so that the meaning identity is
// spelled once.
//
// The family's meaning is the application's promised domain with its
// predicates: what a named horos declaration means: the values of its body, under its parameters.
// What one target's representation covers of it, the native equivalence and a
// permitted normalization belong to that contribution's representation
// revision and not to this package.
package declarations

import g "github.com/Bitspark/bittype/contracts/generation/go"

// Meaning is the meaning identity of the family and its revision. They are
// the same in every target.
var Meaning = g.Meaning{Identity: "bittype.structural/declaration", Revision: "1"}
