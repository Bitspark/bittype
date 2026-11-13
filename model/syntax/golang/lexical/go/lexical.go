// Package lexical is the lexical part of the Go syntax model: the names of the
// predicates over Go spelling. It knows no node kind, role or attribute.
//
// It holds no word list. Go's keywords are those that go/token of the toolchain
// that builds bittype reports, so the model names the rule and does not state
// it.
//
// Each name is a function without input that returns one fixed name. The
// comment says what the predicate is meant to admit; its evaluator states it.
package lexical

import "github.com/Bitspark/bittype/model/syntax/go"

// owner is the owner of the names below: the language's identifier.
const owner = "go"

// Identifier names a Go identifier that is no keyword.
func Identifier() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "identifier"} }

// ImportName names what an import declaration may write before its path: a
// period or an identifier.
func ImportName() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "import-name"} }
