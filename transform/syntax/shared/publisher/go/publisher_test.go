package publisher

import "testing"

// The name is recorded outside this module: a serialized plan records the
// publisher of what was selected. This test spells it on purpose, beside its
// declaration and nowhere else, so that a change of the value is a decision
// and not an accident.
func TestPublisherNameIsPinned(t *testing.T) {
	if Name != "github.com/Bitspark/bittype" {
		t.Errorf("Name is %q", Name)
	}
}
