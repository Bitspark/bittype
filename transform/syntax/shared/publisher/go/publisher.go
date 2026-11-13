// Package publisher names the publisher of the contributions that bittype
// itself publishes (prototype, bittype#196, row 15). A publisher is the library
// that publishes a mapping. The name authenticates nothing: until a
// composition lock binds it (row 15a), registration in one process is the
// binding.
package publisher

// Name is the publisher of bittype's builtin contributions.
const Name = "github.com/Bitspark/bittype"
