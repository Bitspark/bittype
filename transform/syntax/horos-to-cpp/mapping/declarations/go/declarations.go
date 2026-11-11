// Package declarations is the builtin horos-to-cpp mapping library for named
// horos declarations. It selects the C++ forms struct and alias for them.
package declarations

import (
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	model "github.com/Bitspark/bittype/transform/syntax/cpp/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
)

// Library applies the horos declaration recipe with the forms struct and
// alias.
func Library(profile string) rules.Library[s.Language] {
	return rules.Declarations(profile, model.Form("struct"), model.Form("alias"))
}
