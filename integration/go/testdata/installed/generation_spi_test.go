package consumer_test

import (
	formatgolang "github.com/Bitspark/bittype/format/syntax/golang/provider/go"
	syntaxprint "github.com/Bitspark/bittype/print/syntax/shared/provider/go"

	"context"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	fixture "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"

	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func TestInstalledPublicGenerationSPI(t *testing.T) {
	ctx := context.Background()
	f, err := fixture.New(ctx)
	if err != nil {
		t.Fatal(err)
	}
	session, err := g.Begin[fixture.SourceLanguage, s.Language](ctx, f.Input(), f.Config(), g.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for _, root := range f.Roots() {
		publication, err := session.Publish(root)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = session.Define(publication.Callee()); err != nil {
			t.Fatal(err)
		}
	}
	helper, err := session.AcquireHelper("cursor-helper")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = session.Define(helper); err != nil {
		t.Fatal(err)
	}
	output, err := session.Finish()
	if err != nil {
		t.Fatal(err)
	}
	tree, err := representation.Unpack[s.Language](f.Target(), output.Snapshot().Inputs()[0])
	if err != nil {
		t.Fatal(err)
	}
	concrete, err := formatgolang.Format(ctx, tree)
	if err != nil {
		t.Fatal(err)
	}
	text, err := syntaxprint.Print[s.Language](ctx, concrete)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "Page[Page[V, SV], []V]") || !strings.Contains(text, "type InlinePeople = []Person") {
		t.Fatal(text)
	}
}
