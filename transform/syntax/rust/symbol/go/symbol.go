// Package symbol owns the Rust external destinations shared by generation and Files.
package symbol

import (
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
func Std() Crate { return Crate{"std"} }

// Name returns the crate's own name, which is the name the Rust bridge gives
// the crate's root module.
func (c Crate) Name() string { return c.name }

// Dependency returns the dependency whose crate has the given name. A text
// that is no Rust identifier is refused with generation.ErrInvalid, and every
// crate but Std with generation.ErrUnsupported.
func Dependency(name string) (Crate, error) {
	if !native.Identifier(name) {
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
// A crate that is no dependency, an empty path and a segment that is no Rust
// identifier are refused with generation.ErrInvalid.
func Item(crate Crate, path ...string) (Symbol, error) {
	if crate != Std() || len(path) == 0 {
		return Symbol{}, g.ErrInvalid
	}
	for _, segment := range path {
		if !native.Identifier(segment) {
			return Symbol{}, g.ErrInvalid
		}
	}
	return Symbol{item, crate, strings.Join(path, separator)}, nil
}

// StdString is the item string::String of Std.
func StdString() Symbol { return Symbol{item, Std(), "string" + separator + "String"} }

// StdVec is the item vec::Vec of Std.
func StdVec() Symbol { return Symbol{item, Std(), "vec" + separator + "Vec"} }

// Bool is Rust's primitive type bool. It is the one primitive type that is
// issued today. bool is no keyword, so a declaration of that name captures a
// use, and the destination is the primitive type itself.
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
	var key string
	switch symbol.category {
	case item:
		if _, err := Item(symbol.crate, strings.Split(symbol.name, separator)...); err != nil {
			return g.ExternalUse{}, err
		}
		key = item + "\x00" + symbol.crate.name + "\x00" + symbol.name
	case primitive:
		if symbol != Bool() {
			return g.ExternalUse{}, g.ErrInvalid
		}
		key = primitive + "\x00" + symbol.name
	default:
		return g.ExternalUse{}, g.ErrInvalid
	}
	issued, err := issuer.Issue(language, key)
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
	parts := strings.Split(key, "\x00")
	switch {
	case len(parts) == 3 && parts[0] == item:
		crate, err := Dependency(parts[1])
		if err != nil {
			return Symbol{}, g.ErrInvalid
		}
		return Item(crate, strings.Split(parts[2], separator)...)
	case len(parts) == 2 && parts[0] == primitive && parts[1] == Bool().name:
		return Bool(), nil
	}
	return Symbol{}, g.ErrInvalid
}
