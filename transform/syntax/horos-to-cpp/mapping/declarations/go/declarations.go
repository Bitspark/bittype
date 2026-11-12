// Package declarations is the builtin horos-to-cpp mapping library for named
// horos declarations. It selects the C++ forms struct and alias for them.
package declarations

import (
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	model "github.com/Bitspark/bittype/transform/syntax/cpp/model/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
)

// Library applies the horos declaration recipe with the forms struct and
// alias.
func Library() rules.Library[s.Language] {
	return recipe.Library(model.Form(model.FormStruct), model.Form(model.FormAlias))
}
