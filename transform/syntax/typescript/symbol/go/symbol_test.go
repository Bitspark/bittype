package symbol_test

import (
	"errors"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/semantics/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/transform/syntax/typescript/symbol/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"testing"
)

func TestExternalDestinationUsesIssuerAndProfile(t *testing.T) {
	language, err := representation.NewLanguage[ts.Language](ts.Language{}, "native")
	if err != nil {
		t.Fatal(err)
	}
	other, err := representation.NewLanguage[ts.Language](ts.Language{}, "other")
	if err != nil {
		t.Fatal(err)
	}
	want := symbol.Type{Module: "@example/types", Name: "Person"}
	use, err := symbol.ExternalUse(language, ast.NodeRef{Tree: "generated", Node: 1}, want)
	if err != nil {
		t.Fatal(err)
	}
	// The role is this package's own name. It is sealed into every output.
	if use.At.Role != (semantics.Key{Owner: "typescript", Name: "name"}) {
		t.Fatal("role", use.At.Role)
	}
	got, err := symbol.ReadExternal(language, use)
	if err != nil || got != want {
		t.Fatal(got, err)
	}
	if _, err := symbol.ReadExternal(other, use); !errors.Is(err, g.ErrScope) {
		t.Fatal("foreign profile", err)
	}
	foreign, err := g.NewSymbolIssuer().Issue(language, "@example/types\x00Person")
	if err != nil {
		t.Fatal(err)
	}
	use.Symbol = foreign
	if _, err := symbol.ReadExternal(language, use); !errors.Is(err, g.ErrScope) {
		t.Fatal("foreign issuer", err)
	}
}
