// Package lexical is the lexical part of the horos syntax model: the names of
// the predicates over horos spelling. It knows no node kind, role or attribute.
//
// Each name is a function without input that returns one fixed name. The
// comment says what the predicate is meant to admit; its evaluator states it.
package lexical

import "github.com/Bitspark/bittype/model/syntax/go"

// owner is the owner of the names below: the language's identifier.
const owner = "horos"

// Identifier names a text of ASCII letters, digits and underscores that does
// not start with a digit.
func Identifier() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "identifier"} }

// ModuleName names identifiers joined by single periods.
func ModuleName() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "module-name"} }

// Hexadecimal names whole bytes in lowercase hexadecimal digits.
func Hexadecimal() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "hexadecimal"} }

// Pin names a content pin: 64 lowercase hexadecimal digits.
func Pin() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "pin"} }
