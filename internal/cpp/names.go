package cpp

import (
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/model/syntax/cpp/lexical/go"
)

// reserved evaluates the reserved words that the lexical part declares.
var reserved = g.Words(lexical.Reserved())

func Identifier(v string) bool { return g.Identifier(v) && !reserved(v) }
