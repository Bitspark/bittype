// Package lexical is the lexical part of the TypeScript syntax model: the
// names of the predicates over TypeScript spelling, its word lists and the
// patterns of a number. It knows no node kind, role or attribute.
//
// Each name is a function without input that returns one fixed name. The
// comment says what the predicate is meant to admit; its evaluator states it.
package lexical

import "github.com/Bitspark/bittype/model/syntax/go"

// owner is the owner of the names below: the language's identifier.
const owner = "typescript"

// Identifier names a TypeScript identifier that is no reserved word.
func Identifier() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "identifier"} }

// TypeName names an identifier that can name a user-defined type: one that is
// no predefined type word. It is not the TypeName of TypeScript's grammar,
// which is a reference that may be qualified by a namespace. The grammar has
// no term for what this name admits: the name that a type declaration binds
// is its BindingIdentifier, and a sentence beside the grammar excludes the
// predefined type words: "The predefined type keywords are reserved and
// cannot be used as names of user defined types." (TypeScript language
// specification, 3.8.1.)
func TypeName() syntax.Predicate { return syntax.Predicate{Owner: owner, Name: "type-name"} }
