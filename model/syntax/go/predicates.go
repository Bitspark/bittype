package syntax

// Owner is the owner of the predicate names that no language owns. The keys of
// the semantic graph that no language owns have the same owner.
const Owner = "bittype"

// The names below are the predicates that every language's table can use.
// Each is a function without input that returns one fixed name. The comment
// says what the predicate is meant to admit; its evaluator states it.

// Bool names the texts "true" and "false".
func Bool() Predicate { return Predicate{Owner: Owner, Name: "bool"} }

// Decimal names an unsigned decimal integer without a leading zero.
func Decimal() Predicate { return Predicate{Owner: Owner, Name: "decimal"} }

// Nonempty names a text that is not empty and is valid UTF-8.
func Nonempty() Predicate { return Predicate{Owner: Owner, Name: "nonempty"} }

// Positive names a decimal integer other than zero.
func Positive() Predicate { return Predicate{Owner: Owner, Name: "positive"} }

// Signed names a decimal integer with an optional minus sign.
func Signed() Predicate { return Predicate{Owner: Owner, Name: "signed"} }

// String names a text that is valid UTF-8 and holds no NUL.
func String() Predicate { return Predicate{Owner: Owner, Name: "string"} }
