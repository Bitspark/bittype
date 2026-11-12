package integration_test

import (
	"context"
	"errors"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/boolean/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/declarations/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/list/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/optional/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/record/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/text/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
)

// The Go record mapping wraps an optional field in the helper whose key its
// consumer passes. The pagination module has one optional field, Page.next.
// Without a helper registered under exactly that key the generation fails; the
// field is never emitted with its unwrapped type.
func TestGolangRecordRequiresItsOptionalHelper(t *testing.T) {
	type libraries = []g.MappingLibrary[h.Language, gs.Language, source.Reader]
	with := func(key g.HelperKey, more ...g.MappingLibrary[h.Language, gs.Language, source.Reader]) libraries {
		return append(libraries{text.Library(), boolean.Library(), list.Library(), record.Library(key), declarations.Library()}, more...)
	}
	_, _, input := generationInput(t, paginationModule())
	checker := nativeChecker(t, golangText())
	plan := func(set libraries) (g.GenerationPlan, error) {
		var registry c.Registry
		request, err := genfixture.Request(&registry, input, checker.Language(), "generated-go", ggtarget.Provider{Package: "pagination"}, set...)
		if err != nil {
			return g.GenerationPlan{}, err
		}
		return registry.PlanGeneration(context.Background(), request, checker)
	}
	key := g.NewHelperKey("horos/optional")
	if p, err := plan(with(key, optional.Library(key))); err != nil || !p.Valid() {
		t.Fatal("the shared key did not plan", err)
	}
	// An invalid key is refused when the record library is registered.
	if p, err := plan(with(g.HelperKey{}, optional.Library(key))); !errors.Is(err, g.ErrInvalid) || p.Valid() {
		t.Fatal("record library accepted an invalid optional key", err)
	}
	// No library supplies the helper.
	if p, err := plan(with(key)); !errors.Is(err, g.ErrUnsupported) || p.Valid() {
		t.Fatal("optional field planned without its helper", err)
	}
	// A key with the same name and another identity is another helper.
	if p, err := plan(with(key, optional.Library(g.NewHelperKey("horos/optional")))); !errors.Is(err, g.ErrUnsupported) || p.Valid() {
		t.Fatal("optional field planned with a helper under another key", err)
	}
}
