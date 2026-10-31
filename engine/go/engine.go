// Package engine composes explicitly registered syntax providers and directed
// AST transformations. Its zero value is ready to register. It has no built-in
// languages, canonical route, filesystem IO, plugin discovery or mapping rules.
package engine

import (
	"github.com/Bitspark/bittype/trees/representation/go"
	"errors"
	"fmt"
	"sync"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
)

var (
	ErrUnknownLanguage      = errors.New("language is not registered")
	ErrDuplicateLanguage    = errors.New("language is already registered")
	ErrDuplicateTransformer = errors.New("directed transformer is already registered")
	ErrNoPath               = errors.New("no transformation path")
	ErrAmbiguousPath        = errors.New("multiple shortest transformation paths")
	ErrInvalidRoute         = errors.New("invalid transformation route")
	ErrMissingCapability    = errors.New("language has no registered text capability")
	ErrInvalidProvider      = errors.New("invalid text provider")
)

// Engine supports concurrent registration and operations. Registration is
// additive: no silent replacement or unregistration. Each operation captures a
// consistent registry snapshot and invokes user code outside the lock. Providers
// and transformers must support concurrent calls if the application makes them.
// An Engine must not be copied after use.
type Engine struct {
	mu        sync.RWMutex
	languages map[language.ID]registration
	edges     map[language.ID]map[language.ID]edge
}
type registration struct {
	language representation.Language
	text     *Provider
}
type edge struct {
	from, to    language.ID
	transformer transform.ASTTransformer
}

func New() *Engine { return &Engine{} }

// StepError retains the original provider diagnostic through errors.Unwrap.
// Index is zero-based within the selected route's edges.
type StepError struct {
	Index    int
	From, To language.ID
	Err      error
}

func (e *StepError) Error() string {
	return fmt.Sprintf("transform step %d (%s -> %s): %v", e.Index, e.From, e.To, e.Err)
}
func (e *StepError) Unwrap() error { return e.Err }
