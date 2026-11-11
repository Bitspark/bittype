package typescript

import (
	"path/filepath"

	"github.com/Bitspark/bittype/contracts/semantic/go"
)

// Compiler names explicitly provisioned executables. Constructors do no discovery
// or installation; TypeScriptPath is the pinned package's lib/typescript.js.
type Compiler struct {
	NodePath       string
	TypeScriptPath string
}

// Validate requires both paths to be absolute. It reads no file and starts no
// process, so the resolver, the query service and the compiler helper share
// this one check.
func (c Compiler) Validate() error {
	if !filepath.IsAbs(c.NodePath) || !filepath.IsAbs(c.TypeScriptPath) {
		return semantic.ErrInvalidOperation
	}
	return nil
}
