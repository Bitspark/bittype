package rust_test

import (
	"testing"

	native "github.com/Bitspark/bittype/internal/rust"
	"github.com/Bitspark/bittype/model/syntax/rust/lexical/go"
)

// The evaluators read the lists and the table that the lexical part declares.
// The first half asks of every declared entry what the evaluator says of it.
// The second half states expected answers that are written here and not read
// from a list.
func TestEvaluatorsReadTheDeclaredLists(t *testing.T) {
	for _, word := range lexical.Reserved() {
		if !native.Reserved(word) || native.Identifier(word) || native.LifetimeParameter(word) {
			t.Errorf("reserved word %q: reserved %v, an identifier %v, a lifetime parameter %v", word, native.Reserved(word), native.Identifier(word), native.LifetimeParameter(word))
		}
	}
	for _, name := range lexical.ABIs() {
		if !native.ABI(name) {
			t.Errorf("calling convention %q is not evaluated as one", name)
		}
	}
	// Every range of the table: its two ends are characters of a crate name.
	// The ranges ascend and do not touch, so the character before a range and
	// the one after it are none.
	table := lexical.CrateNameCharacters()
	for i, r := range table {
		if r[0] > r[1] || i > 0 && table[i-1][1]+1 >= r[0] {
			t.Fatalf("range %d (%#x to %#x) does not ascend apart from its neighbour", i, r[0], r[1])
		}
		if !native.CrateName(string(r[0])) || !native.CrateName(string(r[1])) || native.CrateName(string(r[0]-1)) || native.CrateName(string(r[1]+1)) {
			t.Errorf("range %d (%#x to %#x) is not evaluated as declared", i, r[0], r[1])
		}
	}
	for text, want := range map[string]bool{
		"as": false, "gen": false, "Self": false, "yield": false, "try": false, "while": false,
		"union": true, "auto": true, "r#type": true, "r#self": false, "_": true, "self_": true, "": false,
	} {
		if got := native.Identifier(text); got != want {
			t.Errorf("Identifier(%q) = %v, want %v", text, got, want)
		}
	}
	for text, want := range map[string]bool{
		"Rust": true, "C": true, "C-unwind": true, "efiapi": true, "sysv64": true,
		"c": false, "rust": false, "C-Unwind": false, "": false, "wasm": false,
	} {
		if got := native.ABI(text); got != want {
			t.Errorf("ABI(%q) = %v, want %v", text, got, want)
		}
	}
	for text, want := range map[string]bool{
		"pages": true, "_": true, "x²": true, "ⅷ": true, "a٣": true, "über": true, "0": true, "z": true, "Z": true,
		// A letter followed by a combining mark: the mark is no character of a crate name.
		"": false, "my-pages": false, "a b": false, "a\u0301": false, "a.b": false, "`": false, "{": false,
	} {
		if got := native.CrateName(text); got != want {
			t.Errorf("CrateName(%q) = %v, want %v", text, got, want)
		}
	}
}
