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

// Library adds a method to a description function and nothing else. Whatever
// profile the function states, Describe returns the function's description
// and its error unchanged: it writes no profile, completes none and refuses
// none. That holds when the function succeeds and when it fails.
func TestLibraryReturnsWhatItsFunctionDescribes(t *testing.T) {
	type description = g.MappingDescription[h.Language, native, source.Reader]
	sentinel := errors.New("the function failed")
	for _, profile := range []string{RepresentationProfile, EnforcementProfile, "mapping-test/another/1", ""} {
		for _, failure := range []error{nil, sentinel} {
			calls := 0
			library := Library[native](func(representation.Language, representation.Language) (description, error) {
				calls++
				// One binding marks the description, so that an emptied
				// description is seen.
				return description{Profile: profile, Bindings: []g.BindingDescription[h.Language, native, source.Reader]{{ID: "mapping-test/binding"}}}, failure
			})
			d, err := library.Describe(representation.Language{}, representation.Language{})
			if d.Profile != profile || len(d.Bindings) != 1 || d.Bindings[0].ID != "mapping-test/binding" || err != failure || calls != 1 {
				t.Errorf("stated %q with %v: described %q with %d bindings and %v after %d calls", profile, failure, d.Profile, len(d.Bindings), err, calls)
			}
		}
	}
}

// A nil function is no library.
func TestLibraryWithoutAFunctionIsInvalid(t *testing.T) {
	d, err := Library[native](nil).Describe(representation.Language{}, representation.Language{})
	if !errors.Is(err, g.ErrInvalid) || d.Profile != "" {
		t.Errorf("a nil function described %q with %v", d.Profile, err)
	}
}
