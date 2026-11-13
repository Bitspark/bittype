package integration_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
)

// The Go record contribution delegates the presence of an optional field to
// whichever contribution the composition registers for the role that the
// horos source issues (row 11, O38). The pagination module has one optional
// field, Page.next. With the optional contribution registered the field is
// wrapped in the helper that this contribution owns. Without one, the role's
// domain has no contribution and the generation fails as "missing
// contribution": the field is never emitted with its unwrapped type. No
// consumer creates a helper key any more.
func TestGolangRecordDelegatesItsOptionalField(t *testing.T) {
	_, _, input := generationInput(t, paginationModule())
	checker := nativeChecker(t, golangText())
	plan := func(names ...string) (g.GenerationPlan, error) {
		var registry c.Registry
		var bundles []genfixture.Bundle[gs.Language]
		for _, name := range names {
			bundles = append(bundles, golangBundleTable()[name])
		}
		request, err := genfixture.Compose(&registry, input, checker.Language(), "generated-go", ggtarget.Provider{Package: "pagination"}, bundles...)
		if err != nil {
			return g.GenerationPlan{}, err
		}
		return registry.PlanGeneration(context.Background(), request, checker)
	}
	without := []string{"text", "boolean", "list", "record", "declarations"}
	p, err := plan(append(without, "optional")...)
	if err != nil || !p.Valid() {
		t.Fatal("the builtin set did not plan", err)
	}
	out, err := p.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if text := formatted(t, golangText(), outputTree[gs.Language](t, out)); !strings.Contains(text, "Next BittypeOptional[string]") || strings.Count(text, "type BittypeOptional[") != 1 {
		t.Fatal("the optional field is not wrapped in the one helper:", text)
	}
	// No contribution answers the role.
	p, err = plan(without...)
	var refusal *g.Refusal
	if !errors.Is(err, g.ErrMissingContribution) || !errors.As(err, &refusal) || p.Valid() {
		t.Fatal("an optional field planned without a contribution for its role", err)
	}
	if refusal.Class != g.RefusalMissingContribution || refusal.Owner != g.KernelOwner || !strings.Contains(refusal.Reason, source.OptionalFieldRole) {
		t.Fatalf("the refusal is %+v", refusal)
	}
	for _, other := range []error{g.ErrUnsupported, g.ErrAmbiguous, g.ErrRefusedByDesign, g.ErrMissingFacet, g.ErrLimit} {
		if errors.Is(err, other) {
			t.Errorf("the refusal also matches %v", other)
		}
	}
}
