// Package source reads established Horos semantics for supplied native mappings.
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
func Bound(language representation.Language) (g.Selector, error) {
	return g.NewSelector(language, "bound")
}

func (s *Source) DeclarationSelector(d g.Declaration[h.Language]) (g.Selector, bool, error) {
	i, ok := s.declarations.Index(d)
	if !ok || s.entities[i] == 0 {
		return g.Selector{}, false, g.ErrScope
	}
	e := s.facts.Entities[s.entities[i]-1]
	pin, err := s.modulePin(e.Owner)
	if err != nil {
		return g.Selector{}, false, err
	}
	key, err := Named(s.snapshot.Language(), pin, e.Name)
	return key, true, err
}
func (s *Source) UseSelector(use g.UseData[h.Language]) (g.Selector, error) {
	if use.Bound.Valid() {
		return Bound(s.snapshot.Language())
	}
	view, err := s.Read(use.Use)
	if err != nil {
		return g.Selector{}, err
	}
	return Builtin(s.snapshot.Language(), view.Constructor, view.Discriminator)
}
