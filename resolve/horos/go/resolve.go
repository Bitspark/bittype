// Package resolve projects canonical Ontos module values into a checked Horos
// AST. Horos owns resolution and substitution; this package performs no IO.
package resolve

import (
	bridge "github.com/Bitspark/bittype/internal/horos"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
)

// Resolve retains the root's complete pinned import closure in one AST.
// treeID is a tooling address, separate from the canonical module pins.
func Resolve(root core.Value, dependencies []core.Value, treeID ast.TreeID) (s.AST, error) {
	program, err := decl.Resolve(root, dependencies)
	if err != nil {
		return s.AST{}, err
	}
	return Project(program, treeID)
}

// Project projects a checked model without specializing its open declarations.
func Project(program *decl.Program, treeID ast.TreeID) (s.AST, error) {
	return bridge.Project(program, treeID)
}

// Model checks the AST and recovers the Horos model for inspection or explicit
// ground specialization. It does not repair bindings or mutate the borrowed AST.
func Model(view ast.View[s.Language]) (*decl.Program, error) {
	return bridge.Read(view)
}
