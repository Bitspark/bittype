// Package engine composes explicitly registered syntax providers and directed
// AST transformations. Its zero value is ready to register. It has no built-in
// languages, canonical route, filesystem IO, plugin discovery or mapping rules.
package engine

import (
	"errors"
	"fmt"
	"github.com/Bitspark/bittype/trees/representation/go"
	"sync"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	generation "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
)

var (
	ErrDuplicateFilesParser     = errors.New("Files parser is already registered")
	ErrDuplicateFilesFormatter  = errors.New("Files formatter is already registered")
	ErrDuplicateFilesPrinter    = errors.New("Files printer is already registered")
	ErrUnknownLanguage          = errors.New("language is not registered")
	ErrDuplicateLanguage        = errors.New("language is already registered")
	ErrDuplicateTransformer     = errors.New("directed transformer is already registered")
	ErrNoPath                   = errors.New("no transformation path")
	ErrAmbiguousPath            = errors.New("multiple shortest transformation paths")
	ErrInvalidRoute             = errors.New("invalid transformation route")
	ErrMissingCapability        = errors.New("language has no registered capability")
	ErrDuplicateSyntaxParser    = errors.New("parser is already registered")
	ErrDuplicateChecker         = errors.New("checker is already registered")
	ErrDuplicateSyntaxFormatter = errors.New("formatter is already registered")
	ErrDuplicateSyntaxPrinter   = errors.New("printer is already registered")
)

// Engine supports concurrent registration and operations. Registration is
// additive: no silent replacement or unregistration. Each operation captures a
// consistent registry snapshot and invokes user code outside the lock. Providers
// and transformers must support concurrent calls if the application makes them.
// An Engine must not be copied after use.
type Engine struct {
	mu         sync.RWMutex
	languages  map[language.ID]registration
	edges      map[language.ID]map[language.ID]edge
	generation generation.Registry
}
type registration struct {
	filesParser             parse.FilesRegistration
	filesFormatter          format.FilesRegistration
	generatedFilesFormatter format.GeneratedFilesRegistration
	filesPrinter            print.FilesRegistration
	language                representation.Language
	parser                  *parse.SyntaxParser
	checker                 *check.Checker
	formatter               *format.SyntaxFormatter
	printer                 *print.SyntaxPrinter
	analyzer                *analysis.Analyzer
	resolver                *resolution.Resolver
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
