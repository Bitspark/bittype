// Package provider composes independently registered Horos capabilities.
// Constructors take an explicit language/profile descriptor and register nothing.
package provider

import s "github.com/Bitspark/bittype/syntax/horos/go"

// Provider exposes the language's implemented roles for direct typed composition.
// Each role remains independently constructible and independently registered.
type Provider struct{ s.Language }
