package symbol_test

import (
	"errors"
	"slices"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	contract "github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/semantics/go"
	rs "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/transform/syntax/rust/symbol/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// foreign is a language of this test only.
type foreign struct{}

func (foreign) ID() contract.ID { return "foreign" }
func (foreign) Name() string    { return "Foreign" }

func language(t *testing.T, profile string) representation.Language {
	t.Helper()
	l, err := representation.NewLanguage[rs.Language](rs.Language{}, profile)
	if err != nil {
		t.Fatal(err)
	}
	return l
}

// The standard library is the one dependency, and a crate is named apart from
// the path of an item.
func TestItemNamesADependencyAndAPath(t *testing.T) {
	std, err := symbol.Dependency("std")
	if err != nil || std != symbol.Std() || std.Name() != "std" {
		t.Fatal(std, err)
	}
	// Crates of a library, and the other crates of the sysroot: a later slice
	// can admit them.
	for _, name := range []string{"core", "alloc", "serde_json", "r#type"} {
		if _, err := symbol.Dependency(name); !errors.Is(err, g.ErrUnsupported) {
			t.Errorf("%s: %v, want unsupported", name, err)
		}
	}
	for _, name := range []string{"", "std::option", "::std", "a-b", "crate", "fn"} {
		if _, err := symbol.Dependency(name); !errors.Is(err, g.ErrInvalid) {
			t.Errorf("%q: %v, want invalid", name, err)
		}
	}
	option, err := symbol.Item(std, "option", "Option")
	if err != nil {
		t.Fatal(err)
	}
	crate, path, ok := option.Item()
	if !ok || crate != symbol.Std() || !slices.Equal(path, []string{"option", "Option"}) {
		t.Fatal(crate, path, ok)
	}
	if _, ok := option.Primitive(); ok {
		t.Error("an item reads as a primitive type")
	}
	// A path is not joined with its crate: the same segments under another
	// split are another destination, or none.
	if other, err := symbol.Item(std, "option"); err != nil || other == option {
		t.Error("a shorter path is the same destination", err)
	}
	for _, path := range [][]string{nil, {""}, {"option", ""}, {"option::Option"}, {"option", "fn"}, {"std", "crate"}, {"a b"}} {
		if _, err := symbol.Item(std, path...); !errors.Is(err, g.ErrInvalid) {
			t.Errorf("%q: %v, want invalid", path, err)
		}
	}
	if _, err := symbol.Item(symbol.Crate{}, "option", "Option"); !errors.Is(err, g.ErrInvalid) {
		t.Error("no dependency", err)
	}
	// The two items that the Rust target itself selects.
	for want, y := range map[string]symbol.Symbol{"string::String": symbol.StdString(), "vec::Vec": symbol.StdVec()} {
		crate, path, ok := y.Item()
		again, err := symbol.Item(crate, path...)
		if !ok || crate != symbol.Std() || err != nil || again != y || path[0]+"::"+path[1] != want {
			t.Errorf("%s: %v %v %v %v", want, crate, path, ok, err)
		}
	}
	if symbol.StdString() == symbol.StdVec() {
		t.Error("two items are one destination")
	}
}

// bool is the one primitive type that has a symbol.
func TestBoolIsAPrimitiveDestination(t *testing.T) {
	name, ok := symbol.Bool().Primitive()
	if !ok || name != "bool" {
		t.Fatal(name, ok)
	}
	if _, _, ok := symbol.Bool().Item(); ok {
		t.Error("a primitive type reads as an item")
	}
	if _, ok := (symbol.Symbol{}).Primitive(); ok {
		t.Error("the zero symbol reads as a primitive type")
	}
	if _, _, ok := (symbol.Symbol{}).Item(); ok {
		t.Error("the zero symbol reads as an item")
	}
}

// A use is read back by the issuer of this package, under the language it was
// issued for, at the role it was recorded with.
func TestExternalDestinationUsesIssuerLanguageAndRole(t *testing.T) {
	native, other := language(t, "native"), language(t, "other")
	at := ast.NodeRef{Tree: "generated", Node: 1}
	option, err := symbol.Item(symbol.Std(), "option", "Option")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []symbol.Symbol{option, symbol.StdString(), symbol.StdVec(), symbol.Bool()} {
		use, err := symbol.ExternalUse(native, at, want)
		if err != nil {
			t.Fatal(err)
		}
		if use.At.Site != at || !use.At.Role.Valid() || use.At.Ordinal != 0 {
			t.Fatal("the use is not at the site", use.At)
		}
		if got, err := symbol.ReadExternal(native, use); err != nil || got != want {
			t.Fatal(got, err)
		}
		if _, err := symbol.ReadExternal(other, use); !errors.Is(err, g.ErrScope) {
			t.Fatal("another profile", err)
		}
		moved := use
		moved.At.Role = semantics.Key{Owner: "rust", Name: "name"}
		if _, err := symbol.ReadExternal(native, moved); !errors.Is(err, g.ErrScope) {
			t.Fatal("another role", err)
		}
		moved = use
		moved.At.Ordinal = 1
		if _, err := symbol.ReadExternal(native, moved); !errors.Is(err, g.ErrScope) {
			t.Fatal("another ordinal", err)
		}
	}
	// The same key from another issuer is another symbol.
	for _, key := range []string{"item\x00std\x00option::Option", "primitive\x00bool", "std::option::Option"} {
		foreign, err := g.NewSymbolIssuer().Issue(native, key)
		if err != nil {
			t.Fatal(err)
		}
		use, err := symbol.ExternalUse(native, at, option)
		if err != nil {
			t.Fatal(err)
		}
		use.Symbol = foreign
		if _, err := symbol.ReadExternal(native, use); !errors.Is(err, g.ErrScope) {
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
	if _, err := symbol.ExternalUse(another, at, option); !errors.Is(err, g.ErrScope) {
		t.Fatal("another language", err)
	}
	if _, err := symbol.ExternalUse(representation.Language{}, at, option); !errors.Is(err, g.ErrScope) {
		t.Fatal("no language", err)
	}
}
