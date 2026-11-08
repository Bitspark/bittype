// Package source interprets established horos semantics for the generation
// kernel and for horos representation libraries.
package source

import (
	"encoding/hex"
	"fmt"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Builtin includes literal discriminators (for example scalar/utf8-text).
// These identities describe Horos grammar, never a selected native mapping.
func Builtin(language representation.Language, constructor, discriminator string) (g.Selector, error) {
	if language.ID() != (h.Language{}).ID() || language.Profile() != h.Profile {
		return g.Selector{}, g.ErrInvalid
	}
	return g.NewSelector(language, fmt.Sprintf("builtin/%d:%s/%d:%s", len(constructor), constructor, len(discriminator), discriminator))
}
func Named(language representation.Language, modulePin, name string) (g.Selector, error) {
	_, err := hex.DecodeString(modulePin)
	if err != nil || len(modulePin) != 64 || name == "" || language.ID() != (h.Language{}).ID() || language.Profile() != h.Profile {
		return g.Selector{}, g.ErrInvalid
	}
	return g.NewSelector(language, "declaration/"+modulePin+"/"+name)
}
