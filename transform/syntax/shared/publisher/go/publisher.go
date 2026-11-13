// Package publisher names the publisher of the contributions that bittype
// itself publishes (docs/contribution-contract.md, row 15). A publisher is
// the library that publishes a contribution, and bittype publishes
// contributions of more than one pair, so the name has one home that names no
// language: the vocabulary of each of bittype's contributions states it in
// the contribution's identity. The name authenticates nothing: until a
// composition lock binds it (row 15a), registration in one process is the
// binding. A contribution of another publisher states its own name in its own
// package and does not import this one.
package publisher

// Name is the publisher of bittype's own contributions.
const Name = "github.com/Bitspark/bittype"
