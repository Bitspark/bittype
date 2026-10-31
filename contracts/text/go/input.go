// Package text defines in-memory native text boundaries and diagnostics.
package text

import "github.com/Bitspark/bittype/trees/ast/go"

// Input IDs are caller-owned opaque addresses, not paths to open or hashes.
type Input struct {
	SourceID ast.SourceID
	TreeID   ast.TreeID
	Text     string
}
