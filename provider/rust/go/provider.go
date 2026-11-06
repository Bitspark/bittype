// Package provider composes independently registered Rust capabilities.
// Constructors take an explicit language/profile descriptor and register nothing.
package provider

import (
	resolve "github.com/Bitspark/bittype/resolve/files/rust/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
)

// Provider exposes the language's implemented roles for direct typed composition.
// Each role remains independently constructible and independently registered.
type Provider struct {
	s.Language
	Compiler resolve.Compiler
}
