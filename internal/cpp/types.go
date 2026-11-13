package cpp

import (
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/model/syntax/cpp/lexical/go"
	"strings"
)

// fundamentalTypeWord evaluates the words of a fundamental type that the
// lexical part declares.
var fundamentalTypeWord = g.Words(lexical.FundamentalTypeWords())

// FundamentalTypeWord reports one word of the spelling of a fundamental type.
func FundamentalTypeWord(s string) bool { return fundamentalTypeWord(s) }

// FundamentalType reports the spelling of a fundamental type.
func FundamentalType(s string) bool {
	words := strings.Fields(s)
	if len(words) == 0 || strings.Join(words, " ") != s {
		return false
	}
	counts := map[string]int{}
	for _, w := range words {
		if !FundamentalTypeWord(w) {
			return false
		}
		counts[w]++
	}
	for w, n := range counts {
		if (w == "long" && n > 2) || (w != "long" && n > 1) {
			return false
		}
	}
	if counts["signed"]+counts["unsigned"] > 1 || counts["short"] > 0 && counts["long"] > 0 {
		return false
	}
	for _, solo := range []string{"void", "bool", "wchar_t", "char8_t", "char16_t", "char32_t", "float"} {
		if counts[solo] > 0 {
			return len(words) == 1
		}
	}
	if counts["double"] > 0 {
		return len(words) == 1 || (len(words) == 2 && counts["long"] == 1)
	}
	if counts["char"] > 0 {
		return len(words) == 1 || (len(words) == 2 && (counts["signed"]+counts["unsigned"] == 1))
	}
	return true
}
