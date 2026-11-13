package typescript

import (
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/model/syntax/typescript/lexical/go"
	"strings"
)

// reserved and predefinedType evaluate the two word lists that the lexical
// part declares.
var reserved, predefinedType = g.Words(lexical.Reserved()), g.Words(lexical.PredefinedTypes())

func Identifier(v string) bool { return g.Identifier(strings.ReplaceAll(v, "$", "_")) && !reserved(v) }

// TypeName excludes predefined type words in binding and unqualified use positions.
// Property names and qualified segments follow their own native rules.
func TypeName(v string) bool {
	return Identifier(v) && !predefinedType(v)
}
