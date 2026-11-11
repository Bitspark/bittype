// Package declarations is the horos-to-typescript mapping library for named
// horos declarations. A consumer registers it on its own; it imports no other
// mapping library.
package declarations

import (
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
)

// Library applies the horos declaration recipe with the TypeScript forms
// "interface" for records and "alias" for every other declaration.
func Library(profile string) rules.Library[s.Language] {
	return recipe.Library(profile, model.Form("interface"), model.Form("alias"))
}
