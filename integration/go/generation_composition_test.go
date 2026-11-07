package integration_test

import (
	checkgolang "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	formatgolang "github.com/Bitspark/bittype/format/syntax/golang/provider/go"
	textfixture "github.com/Bitspark/bittype/integration/go/testdata/textfixture"
	parsegolang "github.com/Bitspark/bittype/parse/syntax/golang/provider/go"

	"context"
	"errors"
	generationfixture "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"

	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"strings"
	"testing"
)

func TestAutomaticGenerationIndependentComposition(t *testing.T) {
	f, err := generationfixture.New(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	out, err := f.Automatic(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	tree, err := representation.Unpack[s.Language](f.Target(), out.EntryAST())
	if err != nil {
		t.Fatal(err)
	}
	text := formatted(t, textfixture.New(parsegolang.Provider{}, checkgolang.Provider{}, formatgolang.Provider{}), tree)
	for _, want := range []string{"type Page[", "type Outer[", "type Nested[", "type Cursor ", "ExternalBox["} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q:\n%s", want, text)
		}
	}
	if strings.Count(text, "type Cursor ") != 1 {
		t.Fatal("helper defined more than once")
	}
}
func TestAutomaticGenerationKeepsWholeResultAssociation(t *testing.T) {
	f, err := generationfixture.NewArgumentOnlySupport(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	out, err := f.Automatic(context.Background())
	if err == nil || out.LanguageID() != "" {
		t.Fatalf("incorrect nested support accepted: %v", err)
	}
}

func TestAutomaticGenerationCancellationReturnsNoOutput(t *testing.T) {
	f, err := generationfixture.New(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	out, err := f.Automatic(ctx)
	if !errors.Is(err, context.Canceled) || out.LanguageID() != "" {
		t.Fatal("canceled invocation produced output", err)
	}
}
