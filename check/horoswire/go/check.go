// Package check admits declaration syntax without resolving names or constructors.
package check
import (
 common "github.com/Bitspark/bittype/internal/declsyntax"
 s "github.com/Bitspark/bittype/syntax/horoswire/go"
 "github.com/Bitspark/bittype/trees/ast/go"
)
func Check(view ast.View[s.Language]) error { return common.Check(view) }
