package integration_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
)

// O46 item 2 (the finding F4 of the two other pairs), gate line 2.
//
// Shows, on the independent fixture source, whose judgments state what is
// wrong and no site: the kernel locates a source judgment. One refuted term is
// the body of two declarations, so the graph does not say which site an
// occurrence of it has. Planning either declaration fails with the source's
// verdict, rejected, and the failure is located at the declaring site of the
// declaration that this occurrence lies in, a different site for each. A
// declaration that the source rejects as a whole is located at its own
// declaring site. The source's explanation stays reachable in the failure.
// Does not show: a term that the graph records at exactly one site (the
// kernel then reports that site; TestKernelLocatesUseFailuresAtTheUse and the
// horos tests of unsupported forms show it), or a source whose graph shares
// one term among the occurrences inside ONE declaration: the kernel reports
// the declaration's site for each of them.
// Red under: O46-judgment-use-not-located, O46-judgment-declaration-not-located.
func TestGateKernelLocatesASourceJudgment(t *testing.T) {
	fixture := gateFixture(t)
	composition := indComposition(t, fixture, indOptions{})
	for _, tc := range []struct{ root, explanation string }{
		{"FaultA", "is refuted"},
		{"FaultB", "is refuted"},
		{"Refused", "is refuted"},
	} {
		plan, err := fixture.PlanContributions(context.Background(), composition, tc.root)
		if err == nil || plan.Valid() {
			t.Fatalf("%s was planned", tc.root)
		}
		if !errors.Is(err, g.ErrRejected) {
			t.Errorf("%s: the failure is not the source's verdict: %v", tc.root, err)
		}
		var located *g.SiteError
		if !errors.As(err, &located) {
			t.Errorf("%s: the kernel did not locate the judgment: %v", tc.root, err)
			continue
		}
		if want := fixture.Site(tc.root); located.Site != want || want.Tree == "" {
			t.Errorf("%s: located at %v, want its declaring site %v", tc.root, located.Site, want)
		}
		if !strings.Contains(located.Message, tc.explanation) || !strings.Contains(err.Error(), tc.explanation) {
			t.Errorf("%s: the source's explanation is lost: %q in %v", tc.root, located.Message, err)
		}
		if matched := gateRefusalClasses(err); len(matched) != 0 {
			t.Errorf("%s: a source judgment matches the refusal values %v", tc.root, matched)
		}
	}
	if fixture.Site("FaultA") == fixture.Site("FaultB") {
		t.Fatal("the two declarations of the shared term have one site: the fixture shows nothing")
	}
}
