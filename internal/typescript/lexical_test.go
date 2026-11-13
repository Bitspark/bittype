package typescript_test

import (
	"testing"

	lit "github.com/Bitspark/bittype/internal/typescript"
	"github.com/Bitspark/bittype/model/syntax/typescript/lexical/go"
)

// The evaluators read the lists and the patterns that the lexical part
// declares. The first half asks of every declared word what the evaluator says
// of it. The second half states expected answers that are written here and not
// read from a list.
func TestEvaluatorsReadTheDeclaredLists(t *testing.T) {
	for _, word := range lexical.Reserved() {
		if lit.Identifier(word) || lit.TypeName(word) {
			t.Errorf("reserved word %q is admitted as an identifier", word)
		}
	}
	for _, word := range lexical.PredefinedTypes() {
		if lit.TypeName(word) {
			t.Errorf("predefined type word %q is admitted as a type name", word)
		}
	}
	for text, want := range map[string]bool{
		"break": false, "this": false, "await": false, "interface": false, "typeof": false,
		"type": true, "of": true, "number": true, "$": true, "a$b": true, "Break": true, "": false, "1a": false,
	} {
		if got := lit.Identifier(text); got != want {
			t.Errorf("Identifier(%q) = %v, want %v", text, got, want)
		}
	}
	for text, want := range map[string]bool{
		"any": false, "undefined": false, "object": false, "bigint": false, "void": false,
		"Any": true, "Page": true, "null_": true, "this": false, "type": true,
	} {
		if got := lit.TypeName(text); got != want {
			t.Errorf("TypeName(%q) = %v, want %v", text, got, want)
		}
	}
	for text, want := range map[string]bool{
		"0": true, "-1": true, "1_000": true, "0x1F": true, "0b101": true, "0o17": true, "1.5e-3": true, ".5": true, "1.": true,
		"": false, "01": false, "1__0": false, "0x": false, "1n": false, "+1": false, "1e": false,
	} {
		if got := lit.Number(text); got != want {
			t.Errorf("Number(%q) = %v, want %v", text, got, want)
		}
	}
	for text, want := range map[string]bool{
		"0n": true, "-12n": true, "1_000n": true, "0xFFn": true, "0b1n": true, "0o7n": true,
		"": false, "1": false, "01n": false, "1.0n": false, "n": false, "1N": false,
	} {
		if got := lit.Bigint(text); got != want {
			t.Errorf("Bigint(%q) = %v, want %v", text, got, want)
		}
	}
}
