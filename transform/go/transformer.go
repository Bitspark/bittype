package transform

import "github.com/Bitspark/bittype/language/go"

// ASTTransformer is a directed, potentially partial AST mapping. Endpoint IDs
// are stable. Implementations return an error for unsupported input, preserve
// necessary provenance, perform no implicit IO, and never return usable partial
// output on failure. Transforming syntax is not proof of semantic equivalence.
type ASTTransformer interface {
 FromID() language.ID
 ToID() language.ID
 Transform(AST) (AST, error)
}

// DescribedTransformer additionally advertises typed/profile endpoints so a
// registry can reject incompatibility before execution. Minimal implementations
// are still checked against registered endpoints at every execution boundary.
type DescribedTransformer interface {
 ASTTransformer
 FromLanguage() Language
 ToLanguage() Language
}
