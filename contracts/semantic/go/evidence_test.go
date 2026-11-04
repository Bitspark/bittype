package semantic

import (
	"errors"
	"github.com/Bitspark/bittype/semantics/go"
	"strings"
	"testing"
)

func TestEvidenceRequirementsBoundAggregateMetadata(t *testing.T) {
	key := semantics.Key{Owner: "test", Name: "operation"}
	entry := EvidenceRequirement{Operation: key, Profile: strings.Repeat("p", DefaultLimits().MaxConfigurationBytes/2)}
	if err := ValidateEvidenceRequirements([]EvidenceRequirement{entry}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateEvidenceRequirements([]EvidenceRequirement{entry, entry}); !errors.Is(err, ErrLimit) {
		t.Fatal("aggregate metadata escaped budget", err)
	}
	if err := ValidateEvidenceRequirements(make([]EvidenceRequirement, DefaultLimits().MaxInputs+1)); !errors.Is(err, ErrLimit) {
		t.Fatal(err)
	}
}
func TestEvidenceRequirementClone(t *testing.T) {
	r := EvidenceRequirement{Operation: semantics.Key{Owner: "test", Name: "operation"}, Profile: "v1", Configuration: []Setting{{Key: semantics.Key{Owner: "test", Name: "setting"}, Value: []byte("a")}}, Reports: []Requirement{{Obligation: semantics.Key{Owner: "test", Name: "valid"}, Verdict: Holds}}}
	copy := r.Clone()
	copy.Configuration[0].Value[0] = 'b'
	copy.Reports[0].Verdict = Rejected
	if string(r.Configuration[0].Value) != "a" || r.Reports[0].Verdict != Holds {
		t.Fatal("requirements alias")
	}
}
