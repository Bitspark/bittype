// Package provider supplies an explicitly registered, AST-only Horos provider.
// It has no text parser, formatter, printer or native-language transformer.
package provider

import (
	check "github.com/Bitspark/bittype/check/horos/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type Provider struct{ s.Language }

func (Provider) Check(view ast.View[s.Language]) error { return check.Check(view) }

var _ text.CheckerProvider[s.Language] = Provider{}
