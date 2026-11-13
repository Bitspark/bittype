package mapping

import "testing"

// The identities are recorded outside this module: a request names a profile,
// and a plan records the profile and the publisher of what was selected. This
// test spells them on purpose, beside their declaration and nowhere else, so
// that a change of a value is a decision and not an accident. Every other test
// and fixture takes the identities from this package.
func TestProfileIdentitiesArePinned(t *testing.T) {
	if RepresentationProfile != "structural-representation/1" {
		t.Errorf("RepresentationProfile is %q", RepresentationProfile)
	}
	if EnforcementProfile != "structural-enforcement/1" {
		t.Errorf("EnforcementProfile is %q", EnforcementProfile)
	}
	if Publisher != "github.com/Bitspark/bittype" {
		t.Errorf("Publisher is %q", Publisher)
	}
}
