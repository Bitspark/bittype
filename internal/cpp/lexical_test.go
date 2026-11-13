package cpp_test

import (
	"testing"

	native "github.com/Bitspark/bittype/internal/cpp"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/model/syntax/cpp/lexical/go"
)

// The evaluators read the lists that the lexical part declares. The first
// half asks of every declared word what the evaluator says of it. The second
// half states expected answers that are written here and not read from a list.
func TestEvaluatorsReadTheDeclaredLists(t *testing.T) {
	for _, word := range lexical.Reserved() {
		if !g.Identifier(word) || native.Identifier(word) {
			t.Errorf("reserved word %q: shaped like an identifier %v, admitted as one %v", word, g.Identifier(word), native.Identifier(word))
		}
	}
	for _, word := range lexical.BuiltinWords() {
		if !native.BuiltinWord(word) {
			t.Errorf("builtin word %q is not evaluated as one", word)
		}
	}
	for text, want := range map[string]bool{
		"alignas": false, "class": false, "xor_eq": false, "co_yield": false, "wchar_t": false,
		"classy": true, "Class": true, "final": true, "override": true, "_": true, "x1": true, "": false, "1x": false,
	} {
		if got := native.Identifier(text); got != want {
			t.Errorf("Identifier(%q) = %v, want %v", text, got, want)
		}
	}
	for text, want := range map[string]bool{
		"void": true, "double": true, "char8_t": true, "unsigned long long": true, "long double": true, "signed char": true, "wchar_t": true,
		"string": false, "size_t": false, "long long long": false, "unsigned  int": false, "": false, "auto": false,
	} {
		if got := native.Builtin(text); got != want {
			t.Errorf("Builtin(%q) = %v, want %v", text, got, want)
		}
	}
}
