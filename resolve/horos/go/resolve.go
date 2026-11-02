// Package resolve projects canonical Ontos module values into a checked Horos
// AST. Horos owns resolution and substitution; this package performs no IO.
package resolve

import (
	"reflect"

	check "github.com/Bitspark/bittype/check/horos/go"
	bridge "github.com/Bitspark/bittype/internal/horos"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
)

// Resolve retains the root's complete pinned import closure in one AST.
// treeID is a tooling address, separate from the canonical module pins.
func Resolve(root core.Value, dependencies []core.Value, treeID ast.TreeID) (s.AST, error) {
	return bridge.Resolve(root, dependencies, treeID)
}

// Project projects a checked model without specializing its open declarations.
// Custom-limit models are re-admitted under this profile's fixed Horos defaults.
func Project(program *decl.Program, treeID ast.TreeID) (s.AST, error) {
	return bridge.Project(program, treeID)
}

// Model checks the AST and recovers the Horos model for inspection or explicit
// ground specialization. It does not repair bindings or mutate the borrowed AST.
func Model(view ast.View[s.Language]) (*decl.Program, error) {
	if view == nil {
		return nil, check.Check(view)
	}
	value := reflect.ValueOf(view)
	switch value.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface, reflect.Func, reflect.Chan:
		if value.IsNil() {
			return nil, check.Check(view)
		}
	}
	// ASTView is an arbitrary implementation. Capture once so validation and
	// reconstruction cannot observe different arenas from successive calls.
	tree := view.ASTView()
	if err := check.Check(tree); err != nil {
		return nil, err
	}
	return bridge.Read(tree)
}
