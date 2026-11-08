package consumer_test

import (
	"context"
	"testing"

	fixture "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
)

func TestInstalledIndependentGenerationPlan(t *testing.T) {
	f, err := fixture.New(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	plan, err := f.Plan(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Steps()) == 0 {
		t.Fatal("empty plan")
	}
	out, err := plan.Run(context.Background())
	if err != nil || !out.Snapshot().Valid() {
		t.Fatal(err)
	}
}
