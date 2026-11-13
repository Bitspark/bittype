package symbol_test

import (
	"errors"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	contract "github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/semantics/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/transform/syntax/golang/symbol/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// foreign is a language of this test only.
type foreign struct{}

func (foreign) ID() contract.ID { return "foreign" }
func (foreign) Name() string    { return "Foreign" }

func language(t *testing.T, profile string) representation.Language {
	t.Helper()
	l, err := representation.NewLanguage[gs.Language](gs.Language{}, profile)
	if err != nil {
		t.Fatal(err)
	}
	return l
}

// A predeclared symbol is issued for a basic type of the universe scope and
// for nothing else.
func TestPredeclaredNamesABasicTypeOfTheUniverseScope(t *testing.T) {
	for _, name := range []string{"string", "bool", "int", "uint8", "byte", "rune", "float64"} {
		y, err := symbol.Predeclared(name)
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if got, ok := y.Predeclared(); !ok || got != name {
			t.Errorf("%s reads back as %q, %v", name, got, ok)
		}
	}
	// Types of the universe scope that are no basic type: the bridge cannot
	// tell them from a declaration yet.
	for _, name := range []string{"error", "comparable", "any"} {
		if _, err := symbol.Predeclared(name); !errors.Is(err, g.ErrUnsupported) {
			t.Errorf("%s: %v, want unsupported", name, err)
		}
	}
	// No type of the universe scope: a value, a function, a keyword, a
	// declared name, a qualified name and no name.
	for _, name := range []string{"", "nil", "true", "len", "type", "Person", "String", "unsafe.Pointer", "untyped bool"} {
		if _, err := symbol.Predeclared(name); !errors.Is(err, g.ErrInvalid) {
			t.Errorf("%q: %v, want invalid", name, err)
		}
	}
	if _, ok := (symbol.Symbol{}).Predeclared(); ok {
		t.Error("the zero symbol reads as a predeclared type")
	}
	a, _ := symbol.Predeclared("byte")
	b, _ := symbol.Predeclared("uint8")
	if a == b {
		t.Error("byte and uint8 are one destination")
	}
}

// A use is read back by the issuer of this package, under the language it was
// issued for, at the role it was recorded with.
func TestExternalDestinationUsesIssuerLanguageAndRole(t *testing.T) {
	native, other := language(t, "native"), language(t, "other")
	want, err := symbol.Predeclared("string")
	if err != nil {
		t.Fatal(err)
	}
	at := ast.NodeRef{Tree: "generated", Node: 1}
	use, err := symbol.ExternalUse(native, at, want)
	if err != nil {
		t.Fatal(err)
	}
	// The role is this package's own name. It is sealed into every output.
	if use.At.Site != at || use.At.Role != (semantics.Key{Owner: "go", Name: "target"}) || use.At.Ordinal != 0 {
		t.Fatal("the use is not at the site", use.At)
	}
	if got, err := symbol.ReadExternal(native, use); err != nil || got != want {
		t.Fatal(got, err)
	}
	if _, err := symbol.ReadExternal(other, use); !errors.Is(err, g.ErrScope) {
		t.Fatal("another profile", err)
	}
	moved := use
	moved.At.Role = semantics.Key{Owner: "another", Name: "role"}
	if _, err := symbol.ReadExternal(native, moved); !errors.Is(err, g.ErrScope) {
		t.Fatal("another role", err)
	}
	moved = use
	moved.At.Ordinal = 1
	if _, err := symbol.ReadExternal(native, moved); !errors.Is(err, g.ErrScope) {
		t.Fatal("another ordinal", err)
	}
	// The same key from another issuer is another symbol.
	for _, key := range []string{"predeclared\x00string", "string", "\x00string"} {
		foreign, err := g.NewSymbolIssuer().Issue(native, key)
		if err != nil {
			t.Fatal(err)
		}
		forged := use
		forged.Symbol = foreign
		if _, err := symbol.ReadExternal(native, forged); !errors.Is(err, g.ErrScope) {
			t.Fatal("a foreign issuer", key, err)
		}
	}
	// No symbol is recorded for the zero symbol or under another language.
	if _, err := symbol.ExternalUse(native, at, symbol.Symbol{}); !errors.Is(err, g.ErrInvalid) {
		t.Fatal("the zero symbol", err)
	}
	another, err := representation.NewLanguage[foreign](foreign{}, "native")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := symbol.ExternalUse(another, at, want); !errors.Is(err, g.ErrScope) {
		t.Fatal("another language", err)
	}
	if _, err := symbol.ExternalUse(representation.Language{}, at, want); !errors.Is(err, g.ErrScope) {
		t.Fatal("no language", err)
	}
}
