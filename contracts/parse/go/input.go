package parse

import "github.com/Bitspark/bittype/trees/ast/go"

// Input IDs are caller-owned opaque addresses, not paths to open or hashes.
//
// bittype:dimension Syntax
type Input struct {
	SourceID ast.SourceID
	TreeID   ast.TreeID
	Text     string
}
