// Package declarations is the builtin horos-to-rust mapping library for named
// horos declarations. A consumer registers it on its own.
package declarations

import (
	s "github.com/Bitspark/bittype/syntax/rust/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	model "github.com/Bitspark/bittype/transform/syntax/rust/model/go"
)

// Library applies the horos declaration recipe with the Rust model forms
// struct, for a record body, and alias, for every other body.
func Library(profile string) rules.Library[s.Language] {
	return recipe.Library(profile, model.Form("struct"), model.Form("alias"))
}
