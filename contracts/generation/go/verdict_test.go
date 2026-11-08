package generation

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// A source diagnostic explains a verdict; it cannot add a second error class.
func TestVerdictKeepsOneFailureClass(t *testing.T) {
	impostor := fmt.Errorf("site 3: %w", ErrLimit)
	for verdictValue, want := range map[Verdict]error{Unsupported: ErrUnsupported, Rejected: ErrRejected, Undetermined: semantic.ErrEvidenceUnavailable} {
		err := verdict(Judgment{Verdict: verdictValue, Diagnostic: impostor})
		if !errors.Is(err, want) || errors.Is(err, ErrLimit) || err.Error() != want.Error()+": site 3: "+ErrLimit.Error() {
			t.Fatal(verdictValue, err)
		}
	}
	// A located diagnostic keeps its site even when it carries a kernel class.
	site := &SiteError{Err: ErrLimit, Message: "deep", Site: ast.NodeRef{Tree: "source", Node: 7}}
	var atSite *SiteError
	if err := verdict(Judgment{Verdict: Unsupported, Diagnostic: site}); !errors.As(err, &atSite) || atSite.Site.Node != 7 || errors.Is(err, ErrLimit) {
		t.Fatal(err)
	}
	if err := verdict(Judgment{Verdict: Admitted, Diagnostic: impostor}); err != nil {
		t.Fatal(err)
	}
	if err := verdict(Judgment{}); !errors.Is(err, ErrInvalid) {
		t.Fatal(err)
	}
	// An ordinary diagnostic stays inspectable through the wrapped chain.
	type located struct{ error }
	plain := located{errors.New("at site 4")}
	var found located
	if err := verdict(Judgment{Verdict: Rejected, Diagnostic: plain}); !errors.As(err, &found) {
		t.Fatal(err)
	}
}
