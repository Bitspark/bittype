package horos_test

import (
	"testing"

	horos "github.com/Bitspark/bittype/syntax/horos/go"
)

// The two values of the scope attribute are spelled here a second time, on
// purpose. The parser and the model projection write them, and the checker,
// the canonical reader and the formatter read them, all from the vocabulary.
// They follow a changed value together, so no other test notices one. A tree
// that was written before holds the old value.
func TestScopeValuesAreUnchanged(t *testing.T) {
	if horos.ScopeLocal != "local" {
		t.Errorf("ScopeLocal = %q, published as %q", horos.ScopeLocal, "local")
	}
	if horos.ScopeImported != "imported" {
		t.Errorf("ScopeImported = %q, published as %q", horos.ScopeImported, "imported")
	}
}
