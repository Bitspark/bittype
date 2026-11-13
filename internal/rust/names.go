package rust

import (
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/model/syntax/rust/lexical/go"
	"strings"
)

// reserved evaluates the reserved words that the lexical part declares.
var reserved = g.Words(lexical.Reserved())

// Reserved reports a strict or reserved keyword.
func Reserved(v string) bool { return reserved(v) }

// IdentifierOrUnderscore reports an identifier or a lone underscore.
func IdentifierOrUnderscore(v string) bool {
	if strings.HasPrefix(v, "r#") {
		name := v[2:]
		return g.Identifier(name) && name != "_" && name != "self" && name != "Self" && name != "super" && name != "crate"
	}
	return g.Identifier(v) && !reserved(v)
}

// Identifier reports a Rust identifier. A lone underscore is none.
func Identifier(v string) bool        { return v != "_" && IdentifierOrUnderscore(v) }
func LifetimeParameter(v string) bool { return v != "_" && g.Identifier(v) && !reserved(v) }
func Lifetime(v string) bool {
	return v == "static" || v == "_" || LifetimeParameter(v)
}
func PathSegment(v string) bool {
	return v == "self" || v == "Self" || v == "super" || v == "crate" || IdentifierOrUnderscore(v)
}
