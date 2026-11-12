package mapping

import (
	"errors"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// native stands for a target language. The adapter never looks at it.
type native struct{}

type description = g.MappingDescription[h.Language, native, source.Reader]

// counted returns a library under the given profile whose rules count their
// calls and return the given error.
func counted(profile string, calls *int, err error) Library[native] {
	return Library[native]{Profile: profile, DescribeRules: func(representation.Language, representation.Language) (description, error) {
		*calls++
		return description{}, err
	}}
}

// The adapter admits a profile by no list. A library answers under the profile
// that its author or its consumer gives it, also under one that this package
// does not declare.
func TestLibraryAnswersUnderTheProfileItIsGiven(t *testing.T) {
	for _, profile := range []string{RepresentationProfile, "mapping-test/another-profile/1"} {
		calls := 0
		d, err := counted(profile, &calls, nil).Describe(representation.Language{}, representation.Language{})
		if err != nil || d.Profile != profile || calls != 1 {
			t.Fatalf("profile %q: described under %q with %v after %d calls", profile, d.Profile, err, calls)
		}
	}
}

// The enforcement profile is declared here and not implemented. The adapter
// refuses it before it runs a rule, so no description exists under that name.
func TestLibraryRefusesTheEnforcementProfile(t *testing.T) {
	calls := 0
	d, err := counted(EnforcementProfile, &calls, nil).Describe(representation.Language{}, representation.Language{})
	if !errors.Is(err, g.ErrUnsupported) || d.Profile != "" || calls != 0 {
		t.Fatalf("described under %q with %v after %d calls", d.Profile, err, calls)
	}
}

func TestLibraryNeedsRulesAndReturnsTheirError(t *testing.T) {
	if _, err := (Library[native]{Profile: RepresentationProfile}).Describe(representation.Language{}, representation.Language{}); !errors.Is(err, g.ErrInvalid) {
		t.Fatal("a library without rules:", err)
	}
	sentinel := errors.New("rules failed")
	calls := 0
	if _, err := counted(RepresentationProfile, &calls, sentinel).Describe(representation.Language{}, representation.Language{}); !errors.Is(err, sentinel) || calls != 1 {
		t.Fatal("the rules' error was lost:", err, calls)
	}
}

// The two declared profiles are two promises and need two identities.
func TestDeclaredProfilesAreDistinct(t *testing.T) {
	if RepresentationProfile == "" || EnforcementProfile == "" || RepresentationProfile == EnforcementProfile {
		t.Fatal(RepresentationProfile, EnforcementProfile)
	}
}
