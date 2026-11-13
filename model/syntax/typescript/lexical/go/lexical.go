// Package lexical is the lexical part of the TypeScript syntax model: the
// names of the predicates over TypeScript spelling. It knows no node kind, role
// or attribute.
//
// Each name is a function without input that returns one fixed name. The
// comment says what the predicate is meant to admit; its evaluator states it.
package lexical

import "github.com/Bitspark/bittype/model/syntax/go"

// owner is the owner of the names below: the language's identifier.
const owner = "typescript"

// Identifier names a TypeScript identifier that is no reserved word.
func Identifier() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "identifier"} }

// TypeName names an identifier that can name a type: one that is no predefined
// type word.
func TypeName() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "type-name"} }
