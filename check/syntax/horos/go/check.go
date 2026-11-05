// Package check admits declaration syntax without resolving names or constructors.
package check

import (
	checkcontract "github.com/Bitspark/bittype/contracts/check/go"
	common "github.com/Bitspark/bittype/internal/declsyntax"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func Check(view ast.View[s.Language]) error { return common.Check(view) }

// Provider supplies syntax admission independently of every other role.
type Provider struct{ s.Language }

func (Provider) Check(view ast.View[s.Language]) error { return Check(view) }

var _ checkcontract.Provider[s.Language] = Provider{}
