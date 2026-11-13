package cpp

import (
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/model/syntax/cpp/lexical/go"
	"strings"
)

// builtinWord evaluates the words of a builtin type that the lexical part
// declares.
var builtinWord = g.Words(lexical.BuiltinWords())

func BuiltinWord(s string) bool { return builtinWord(s) }
func Builtin(s string) bool {
	words := strings.Fields(s)
	if len(words) == 0 || strings.Join(words, " ") != s {
		return false
	}
	counts := map[string]int{}
	for _, w := range words {
		if !BuiltinWord(w) {
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
