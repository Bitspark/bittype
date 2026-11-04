package integration_test

import (
 "github.com/Bitspark/bittype/contracts/check/go"
 "github.com/Bitspark/bittype/contracts/parse/go"
 "github.com/Bitspark/bittype/contracts/format/go"
 "github.com/Bitspark/bittype/contracts/print/go"
)

// Native conformance fixtures exercise four capabilities of one implementation.
// Engine registration and production contracts do not require this combination.
type nativeProvider[L any] interface {
 parse.Provider[L]
 check.Provider[L]
 format.Provider[L]
 print.Provider[L]
}
