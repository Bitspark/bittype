// Package generation composes supplied mappings without interpreting source or target types.
package generation

import (
	"fmt"
	"slices"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Selector is an exact source-owned identity, qualified by language/profile/marker.
// Its key is never interpreted by the composer. Display names and graph-local IDs
// are not substitutes for a source integration's stable declaration identity.
type Selector struct {
	language representation.Language
	key      string
}

func NewSelector(language representation.Language, key string) (Selector, error) {
	if language.Validate() != nil || key == "" {
		return Selector{}, g.ErrInvalid
	}
	return Selector{language, key}, nil
}
func (s Selector) Key() string                       { return s.key }
func (s Selector) Valid() bool                       { return s.language.Validate() == nil && s.key != "" }
func (s Selector) Language() representation.Language { return s.language }

// Mapping creates a plan for a real declaration, or a direct expression when
// declaration is zero. Factories allocate invocation-local collaborators only;
// representation work is demanded through the returned SPI operations.
type Mapping[F, T any] struct {
	ID       string
	PluginID string
	Plan     func(g.Declaration[F]) (g.Plan[F, T], error)
}
type Binding[F, T any] struct {
	Selector Selector
	Mapping  Mapping[F, T]
}
type Override[F, T any] struct {
	Binding    Binding[F, T]
	ExpectedID string
}

// Bindings captures one effective table. Overrides replace an existing exact
// binding once, with its expected configuration ID; they are never a fallback.
type Bindings[F, T any] struct {
	language   representation.Language
	entries    []Binding[F, T]
	structural *Mapping[F, T]
}

func NewBindings[F, T any](from representation.Language, entries []Binding[F, T], overrides []Override[F, T], structural *Mapping[F, T]) (Bindings[F, T], error) {
	if err := representation.CheckMarker[F](from); err != nil {
		return Bindings[F, T]{}, err
	}
	out := Bindings[F, T]{language: from, entries: slices.Clone(entries)}
	indexes := map[string]int{}
	validMapping := func(m Mapping[F, T]) bool { return m.ID != "" && m.PluginID != "" && m.Plan != nil }
	for i, b := range out.entries {
		if !b.Selector.Valid() || !from.Compatible(b.Selector.language) || !validMapping(b.Mapping) {
			return Bindings[F, T]{}, g.ErrInvalid
		}
		if _, exists := indexes[b.Selector.key]; exists {
			return Bindings[F, T]{}, fmt.Errorf("%w: duplicate selector %s", g.ErrInvalid, b.Selector.key)
		}
		indexes[b.Selector.key] = i
	}
	seen := map[string]bool{}
	for _, override := range overrides {
		b := override.Binding
		i, exists := indexes[b.Selector.key]
		if !b.Selector.Valid() || !from.Compatible(b.Selector.language) || !exists || seen[b.Selector.key] || override.ExpectedID == "" || out.entries[i].Mapping.ID != override.ExpectedID || !validMapping(b.Mapping) || b.Mapping.ID == override.ExpectedID {
			return Bindings[F, T]{}, fmt.Errorf("%w: invalid override %s", g.ErrInvalid, b.Selector.key)
		}
		seen[b.Selector.key] = true
		out.entries[i] = b
	}
	ids := map[string]bool{}
	for _, b := range out.entries {
		if ids[b.Mapping.ID] {
			return Bindings[F, T]{}, g.ErrInvalid
		}
		ids[b.Mapping.ID] = true
	}
	if structural != nil {
		if !validMapping(*structural) || ids[structural.ID] {
			return Bindings[F, T]{}, g.ErrInvalid
		}
		copy := *structural
		out.structural = &copy
	}
	return out, nil
}

func (b Bindings[F, T]) selectMapping(key Selector, structural bool) (Mapping[F, T], error) {
	if !key.Valid() || !b.language.Compatible(key.language) {
		return Mapping[F, T]{}, g.ErrInvalid
	}
	for _, entry := range b.entries {
		if entry.Selector.key == key.key {
			return entry.Mapping, nil
		}
	}
	if structural && b.structural != nil {
		return *b.structural, nil
	}
	return Mapping[F, T]{}, fmt.Errorf("%w: no representation for %s", g.ErrUnsupported, key.key)
}
