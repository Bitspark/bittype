// Package lexical is the lexical part of the Rust syntax model: the names of
// the predicates over Rust spelling, its reserved words and the table of the
// characters of a crate name. It knows no node kind, role or attribute.
//
// Each name is a function without input that returns one fixed name. The
// comment says what the predicate is meant to admit; its evaluator states it.
package lexical

import "github.com/Bitspark/bittype/model/syntax/go"

// owner is the owner of the names below: the language's identifier.
const owner = "rust"

// Identifier names a Rust identifier: a word that is no reserved word, or a
// raw identifier. A lone underscore is none.
func Identifier() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "identifier"} }

// IdentifierOrUnderscore names an identifier or a lone underscore.
func IdentifierOrUnderscore() syntax.Predicate {
	return syntax.Predicate{Owner: owner, Name: "identifier-or-underscore"}
}

// LifetimeParameter names the name of a lifetime parameter, without its
// apostrophe.
func LifetimeParameter() syntax.Predicate {
	return syntax.Predicate{Owner: owner, Name: "lifetime-parameter"}
}

// Lifetime names a lifetime, without its apostrophe: the name of a lifetime
// parameter, "static" or a lone underscore.
func Lifetime() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "lifetime"} }

// PathSegment names one segment of a path: an identifier or one of the path
// keywords.
func PathSegment() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "path-segment"} }
