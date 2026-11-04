// Package transform defines directed AST transformations and their typed bridge.
// Concrete mappings belong to separate transforms/<source>-to-<target> packages.
package transform

import (
	"context"
	"errors"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var ErrInvalidTransformer = errors.New("invalid AST transformer")

// Input supplies captured syntax and ordered receipts explicitly. Callers must
// not mutate Evidence during invocation. No provider runs implicitly.
type Input struct {
	Source   semantic.Snapshot
	Evidence []semantic.Result
}

// Contract declares exactly the evidence required by a typed mapping.
type Contract struct {
	Evidence []semantic.EvidenceRequirement
}

// ASTTransformer is a directed, potentially partial AST mapping. Endpoint IDs
// are stable. Implementations return an error for unsupported input, preserve
// necessary provenance, perform no implicit IO, and never return usable partial
// output on failure. Transforming syntax is not proof of semantic equivalence.
type ASTTransformer interface {
	FromID() language.ID
	ToID() language.ID
	Transform(context.Context, Input) (representation.ASTSet, error)
}

// DescribedTransformer additionally advertises typed/profile endpoints so a
// registry can reject incompatibility before execution. Minimal implementations
// are still checked against registered endpoints at every execution boundary.
type DescribedTransformer interface {
	ASTTransformer
	FromLanguage() representation.Language
	ToLanguage() representation.Language
}
