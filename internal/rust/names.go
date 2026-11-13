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
func Identifier(v string) bool {
	if strings.HasPrefix(v, "r#") {
		name := v[2:]
		return g.Identifier(name) && name != "_" && name != "self" && name != "Self" && name != "super" && name != "crate"
	}
	return g.Identifier(v) && !reserved(v)
}
func Declaration(v string) bool       { return v != "_" && Identifier(v) }
func LifetimeParameter(v string) bool { return v != "_" && g.Identifier(v) && !reserved(v) }
func LifetimeReference(v string) bool {
	return v == "static" || v == "_" || LifetimeParameter(v)
}
func PathSegment(v string) bool {
	return v == "self" || v == "Self" || v == "super" || v == "crate" || Identifier(v)
}
