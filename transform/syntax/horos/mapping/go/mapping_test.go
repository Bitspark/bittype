package mapping

import (
	"errors"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

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

// native stands for a target language. Library never looks at it.
type native struct{}

// Library adds a method to a description function and nothing else. It
// returns the function's description and error unchanged: it writes no
// profile, completes none and refuses none. A nil function is invalid.
func TestLibraryReturnsWhatItsFunctionDescribes(t *testing.T) {
	type description = g.MappingDescription[h.Language, native, source.Reader]
	sentinel := errors.New("the function failed")
	for _, profile := range []string{RepresentationProfile, EnforcementProfile, "mapping-test/another/1", ""} {
		calls := 0
		library := Library[native](func(representation.Language, representation.Language) (description, error) {
			calls++
			return description{Profile: profile}, sentinel
		})
		d, err := library.Describe(representation.Language{}, representation.Language{})
		if d.Profile != profile || !errors.Is(err, sentinel) || calls != 1 {
			t.Errorf("stated %q: described %q with %v after %d calls", profile, d.Profile, err, calls)
		}
	}
	if d, err := Library[native](nil).Describe(representation.Language{}, representation.Language{}); !errors.Is(err, g.ErrInvalid) || d.Profile != "" {
		t.Errorf("a nil function described %q with %v", d.Profile, err)
	}
}
