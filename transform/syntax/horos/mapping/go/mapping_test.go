package mapping

import "testing"

// The two identities are recorded outside this module: a request names one,
// and a plan entry will record the one that was selected. This test spells
// both on purpose, beside their declaration and nowhere else, so that a change
// of either value is a decision and not an accident. Every other test and
// fixture takes the identities from this package.
func TestProfileIdentitiesArePinned(t *testing.T) {
	if RepresentationProfile != "structural-representation/1" {
		t.Errorf("RepresentationProfile is %q", RepresentationProfile)
	}
	if EnforcementProfile != "structural-enforcement/1" {
		t.Errorf("EnforcementProfile is %q", EnforcementProfile)
	}
}
