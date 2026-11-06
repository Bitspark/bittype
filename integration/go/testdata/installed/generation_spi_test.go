package consumer_test

import (
	"context"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	fixture "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	gp "github.com/Bitspark/bittype/provider/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func TestInstalledPublicGenerationSPI(t *testing.T) {
	ctx := context.Background()
	f, err := fixture.New(ctx)
	if err != nil { t.Fatal(err) }
	session, err := g.Begin[fixture.SourceLanguage, s.Language](ctx, f.Input(), f.Config(), g.DefaultLimits())
	if err != nil { t.Fatal(err) }
	for _, root := range f.Roots() {
		publication, err := session.Publish(root)
		if err != nil { t.Fatal(err) }
		if _, err = session.Define(publication.Callee()); err != nil { t.Fatal(err) }
	}
	helper, err := session.AcquireHelper("cursor-helper")
	if err != nil { t.Fatal(err) }
	if _, err = session.Define(helper); err != nil { t.Fatal(err) }
	output, err := session.Finish()
	if err != nil { t.Fatal(err) }
	tree, err := representation.Unpack[s.Language](f.Target(), output.EntryAST())
	if err != nil { t.Fatal(err) }
	concrete, err := (gp.Provider{}).Format(ctx, tree)
	if err != nil { t.Fatal(err) }
	text, err := (gp.Provider{}).Print(ctx, concrete)
	if err != nil { t.Fatal(err) }
	if !strings.Contains(text, "Page[Page[V, SV], []V]") || !strings.Contains(text, "type InlinePeople = []Person") { t.Fatal(text) }
}
