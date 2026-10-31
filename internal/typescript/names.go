package typescript

import (
	g "github.com/Bitspark/bittype/internal/grammar"
	"strings"
)

var reserved = map[string]bool{}

func init() {
	for _, word := range strings.Fields("break case catch class const continue debugger default delete do else enum export extends false finally for function if import in instanceof new null return super switch throw true try typeof var void while with implements interface let package private protected public static yield await this") {
		reserved[word] = true
	}
}
func Identifier(v string) bool { return g.Identifier(strings.ReplaceAll(v, "$", "_")) && !reserved[v] }

// TypeName excludes predefined type words in binding and unqualified use positions.
// Property names and qualified segments follow their own native rules.
func TypeName(v string) bool {
	if !Identifier(v) {
		return false
	}
	switch v {
	case "any", "unknown", "never", "number", "bigint", "boolean", "string", "symbol", "void", "object", "undefined":
		return false
	}
	return true
}
