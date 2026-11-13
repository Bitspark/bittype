// Package lexical is the lexical part of the Rust syntax model: the names of
// the predicates over Rust spelling, its word lists and the table of the
// characters of a crate name. It knows no node kind, role or attribute.
//
// Each name is a function without input that returns one fixed name. The
// comment says what the predicate is meant to admit; its evaluator states it.
package lexical

import "github.com/Bitspark/bittype/model/syntax/go"

// owner is the owner of the names below: the language's identifier.
const owner = "rust"

// Identifier names a Rust identifier that is no reserved word, or a raw
// identifier.
func Identifier() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "identifier"} }

// Declaration names an identifier that can name a declaration: any but a
// lone underscore.
func Declaration() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "declaration"} }

// LifetimeParameter names the name of a lifetime parameter, without its
// apostrophe.
func LifetimeParameter() syntax.Predicate {
	return syntax.Predicate{Owner: owner, Name: "lifetime-parameter"}
}

// LifetimeReference names the name of a lifetime in use, without its
// apostrophe: a parameter's name, "static" or a lone underscore.
func LifetimeReference() syntax.Predicate {
	return syntax.Predicate{Owner: owner, Name: "lifetime-reference"}
}

// PathSegment names one segment of a path: an identifier or one of the path
// keywords.
func PathSegment() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "path-segment"} }

// ABI names the name of a calling convention.
func ABI() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "abi"} }
