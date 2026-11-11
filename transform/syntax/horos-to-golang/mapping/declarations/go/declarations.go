// Package declarations is the horos-to-golang mapping library for named
// declarations: a horos record declaration becomes a Go defined type and any
// other declaration a Go alias.
package declarations

import (
	s "github.com/Bitspark/bittype/syntax/golang/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
)

// Library returns the declaration mapping under the given profile.
func Library(profile string) rules.Library[s.Language] {
	return rules.Declarations(profile, model.Form("defined"), model.Form("alias"))
}
