package horos_test

import (
	"strings"
	"testing"

	spelling "github.com/Bitspark/bittype/internal/horos"
	"github.com/Bitspark/bittype/model/syntax/horos/lexical/go"
)

// The evaluators read the list and the bound that the lexical part declares.
// The first half asks of every declared word what the evaluator says of it.
// The second half states expected answers that are written here and not read
// from the list.
func TestEvaluatorsReadTheDeclaredLists(t *testing.T) {
	for _, word := range lexical.Reserved() {
		if !spelling.Identifier(word) || !spelling.Reserved(word) || spelling.Name(word) != "`"+word+"`" {
			t.Errorf("reserved word %q: an identifier %v, reserved %v, written %s", word, spelling.Identifier(word), spelling.Reserved(word), spelling.Name(word))
		}
	}
	for text, want := range map[string]bool{
		"module": true, "hex": true, "nu": true, "alternative": true, "null": true,
		"Module": false, "modules": false, "kinds": false, "": false, "page": false,
	} {
		if got := spelling.Reserved(text); got != want {
			t.Errorf("Reserved(%q) = %v, want %v", text, got, want)
		}
	}
	if got := spelling.ModuleName("example.type.page"); got != "example.`type`.page" {
		t.Errorf("ModuleName wrote %s", got)
	}
	// The bound on the digits of an integer is 4096, with or without a sign.
	for digits, want := range map[int]bool{1: true, 4096: true, 4097: false} {
		text := strings.Repeat("7", digits)
		if spelling.Integer(text) != want || spelling.Integer("-"+text) != want {
			t.Errorf("an integer of %d digits: %v and %v, want %v", digits, spelling.Integer(text), spelling.Integer("-"+text), want)
		}
	}
	if spelling.MaxIntegerDigits != 4096 {
		t.Errorf("the bound is %d", spelling.MaxIntegerDigits)
	}
}
