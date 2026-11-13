// Package symbol owns the Rust external destinations shared by generation and Files.
package symbol

import (
	"errors"
	"strings"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	native "github.com/Bitspark/bittype/internal/rust"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Crate identifies a dependency that the Rust environment selects. It is a
// dependency and not a spelling: the name under which a crate root sees a
// dependency can differ from the crate's name, and Rust Files chooses that
// local name. The zero Crate is no dependency. Crates are comparable.
type Crate struct{ name string }

// Std is the crate std of the standard library of the pinned compiler. It is
// the one dependency of today. A crate of a library needs that crate's
// semantics.
//
// The name is spelled here for what this package answers: the name that the
// compiler gives the crate, which Name returns and Dependency resolves. The
// Rust model declares the same name for the external form, because a mapping
// may import the model and not this package, and this package may not import
// the model. It is one name with two declarations, and a test keeps them
// equal.
func Std() Crate { return Crate{"std"} }

// Name returns the crate's own name, which is the name the Rust bridge gives
// the crate's root module.
func (c Crate) Name() string { return c.name }

// segment judges one segment of an item's path. A key holds a name and never
// a spelling.
//
// A raw identifier as written text (r#name) is one spelling of a name and not
// the name, so it is invalid and gives no second symbol for one item. The
// underscore names no item and is invalid.
//
// A keyword that the language lets be written raw, such as fn, is a name:
// r#fn is its spelling, and Rust Files chooses a raw spelling where the
// language needs one. This profile covers no such name, so it is unsupported
// and a later slice can admit it. Which keywords have a raw form is the
// lexical helper's rule, the language's: every keyword but self, Self, super
// and crate. Those four name no item and are invalid, as is every text that
// is no identifier.
func segment(text string) error {
	switch {
	case strings.HasPrefix(text, "r#"), text == "_":
		return g.ErrInvalid
	case native.Identifier(text):
		return nil
	case native.Identifier("r#" + text):
		return g.ErrUnsupported
	}
	return g.ErrInvalid
}

// Dependency returns the dependency whose crate has the given name. The name
// is judged by the compiler's one rule for the name of a crate,
// native.CrateName, as the Rust environment judges the name of a captured
// crate. A text that the compiler does not take as the name of a crate is
// refused with generation.ErrInvalid: the empty text, a raw identifier, a
// text with a hyphen or a path separator. A name that the compiler takes and
// that is no dependency of today is refused with generation.ErrUnsupported:
// every crate but Std, a keyword such as fn and the underscore among them.
func Dependency(name string) (Crate, error) {
	if !native.CrateName(name) {
		return Crate{}, g.ErrInvalid
	}
	if name != Std().name {
		return Crate{}, g.ErrUnsupported
	}
	return Std(), nil
}

// Symbol is a destination outside the declarations of a generated output. It
// states its category. The zero Symbol is no destination. Symbols are
// comparable.
type Symbol struct {
	category string
	crate    Crate
	name     string
}

// The categories of today.
const (
	item      = "item"
	primitive = "primitive"
)

// separator joins the path of an item. No Rust identifier contains it.
const separator = "::"

// Item identifies the public item at path below the root of crate. The crate
// is separate from the path. The destination is that export, and never the
// definition that the export leads to: the item string::String of Std is
// another destination than the definition in the crate alloc.
//
// A crate that is no dependency, an empty path and a segment that names no
// item are refused with generation.ErrInvalid: a raw identifier, the
// underscore, self, Self, super, crate, and a text that is no identifier. A
// path whose segments are names and that holds a keyword with a raw form is
// refused with generation.ErrUnsupported. Validity is judged over every
// segment before coverage.
func Item(crate Crate, path ...string) (Symbol, error) {
	if crate != Std() || len(path) == 0 {
		return Symbol{}, g.ErrInvalid
	}
	var uncovered error
	for _, text := range path {
		if err := segment(text); errors.Is(err, g.ErrInvalid) {
			return Symbol{}, err
		} else if err != nil {
			uncovered = err
		}
	}
	if uncovered != nil {
		return Symbol{}, uncovered
	}
	return Symbol{item, crate, strings.Join(path, separator)}, nil
}

// StdString is the item string::String of Std.
func StdString() Symbol { return Symbol{item, Std(), "string" + separator + "String"} }

// StdVec is the item vec::Vec of Std.
func StdVec() Symbol { return Symbol{item, Std(), "vec" + separator + "Vec"} }

// Bool is Rust's primitive type bool. It is the one primitive type that is
// issued today. bool is no keyword, so a declaration or a generic parameter of
// that name captures a use, and the destination is the primitive type itself.
//
// The Rust helper binds a primitive type node to what the compiler resolved
// the name to. A use reaches this destination when that binding names the
// entity of the primitive type: an entity of the primitive kind of the Rust
// vocabulary whose name is the one that Primitive returns. A use that a
// declaration or a generic parameter captures binds to that declaration
// instead. docs/destinations.md states the rule.
func Bool() Symbol { return Symbol{primitive, Crate{}, "bool"} }

// Item returns the dependency and the path of an item, and false for a symbol
// of another category.
func (y Symbol) Item() (Crate, []string, bool) {
	if y.category != item {
		return Crate{}, nil, false
	}
	return y.crate, strings.Split(y.name, separator), true
}

// Primitive returns the name of a primitive type, and false for a symbol of
// another category.
func (y Symbol) Primitive() (string, bool) { return y.name, y.category == primitive }

var issuer = g.NewSymbolIssuer()

// role is the role of an external use on its node: the role that a reference
// to a declaration of the output has on such a node. A name use carries one of
// the two.
var role = semantics.Key{Owner: string(s.Language{}.ID()), Name: s.RoleTarget}

// ExternalUse records symbol as the destination of the name use at site.
func ExternalUse(language representation.Language, site ast.NodeRef, symbol Symbol) (g.ExternalUse, error) {
	if representation.CheckMarker[s.Language](language) != nil || language.ID() != (s.Language{}).ID() {
		return g.ExternalUse{}, g.ErrScope
	}
	// A symbol of a category was made by a constructor of this package. Only
	// the zero symbol has none.
	if symbol.category == "" {
		return g.ExternalUse{}, g.ErrInvalid
	}
	issued, err := issuer.Issue(language, symbol.category+"\x00"+symbol.crate.name+"\x00"+symbol.name)
	if err != nil {
		return g.ExternalUse{}, err
	}
	return g.ExternalUse{At: semantics.Occurrence{Site: site, Role: role}, Symbol: issued}, nil
}

// ReadExternal returns the destination that an external use of this package
// records. It refuses a use of another issuer, of another language or at
// another role with generation.ErrScope.
func ReadExternal(language representation.Language, use g.ExternalUse) (Symbol, error) {
	key, ok := issuer.Read(language, use.Symbol)
	if !ok || use.At.Role != role || use.At.Ordinal != 0 {
		return Symbol{}, g.ErrScope
	}
	// Every key of this issuer is a category, a crate and a name.
	category, rest, _ := strings.Cut(key, "\x00")
	crate, name, _ := strings.Cut(rest, "\x00")
	if category == primitive {
		return Bool(), nil
	}
	return Item(Crate{crate}, strings.Split(name, separator)...)
}
