// Package transform defines directed AST transformations and their typed bridge.
// Concrete mappings belong to separate transforms/<source>-to-<target> packages.
package transform

import (
"errors"
"github.com/Bitspark/bittype/contracts/language/go"
"github.com/Bitspark/bittype/trees/representation/go"
)

var ErrInvalidTransformer = errors.New("invalid AST transformer")

// ASTTransformer is a directed, potentially partial AST mapping. Endpoint IDs
// are stable. Implementations return an error for unsupported input, preserve
// necessary provenance, perform no implicit IO, and never return usable partial
// output on failure. Transforming syntax is not proof of semantic equivalence.
type ASTTransformer interface {
	FromID() language.ID
	ToID() language.ID
	Transform(representation.AST) (representation.AST, error)
}

// DescribedTransformer additionally advertises typed/profile endpoints so a
// registry can reject incompatibility before execution. Minimal implementations
// are still checked against registered endpoints at every execution boundary.
type DescribedTransformer interface {
	ASTTransformer
	FromLanguage() representation.Language
	ToLanguage() representation.Language
}
