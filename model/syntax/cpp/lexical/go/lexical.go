// Package lexical is the lexical part of the C++ syntax model: the names of
// the predicates over C++ spelling, and its word lists. It knows no node kind,
// role or attribute.
//
// Each name is a function without input that returns one fixed name. The
// comment says what the predicate is meant to admit; its evaluator states it.
package lexical

import "github.com/Bitspark/bittype/model/syntax/go"

// owner is the owner of the names below: the language's identifier.
const owner = "cpp"

// Identifier names a C++ identifier that is no reserved word.
func Identifier() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "identifier"} }

// FundamentalType names the spelling of a fundamental type: its words in a
// combination the language admits, separated by single spaces, such as
// "unsigned long long".
func FundamentalType() syntax.Predicate {
	return syntax.Predicate{Owner: owner, Name: "fundamental-type"}
}

// HeaderName names a header name without its delimiters: a text that is not
// empty and holds no line break.
func HeaderName() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "header-name"} }
