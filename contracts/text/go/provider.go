// Package text composes text capabilities and supplies shared diagnostics.
package text

import (
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
)

// Provider composes independently usable stage contracts for native text tooling.
type Provider[L any] interface {
	parse.Provider[L]
	format.Provider[L]
	check.Provider[L]
	print.Provider[L]
}
