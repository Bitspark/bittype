// Package symbol owns the Go external destinations shared by generation and Files.
package symbol

import (
	"go/types"
	"strings"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Symbol is a destination outside the declarations of a generated output. It
// states its category: a predeclared type is never a package symbol with an
// empty package. The zero Symbol is no destination. Symbols are comparable.
type Symbol struct{ category, name string }

// The one category of today. A type of another package is a further category
// and needs that package's semantics.
const predeclared = "predeclared"

// Predeclared identifies the type that name denotes in Go's universe scope.
// The destination is that predeclared entity itself, and never a declaration
// or a type parameter that has the same name.
//
// A name that the universe scope does not declare as a type is refused with
// generation.ErrInvalid. Only a basic type is issued today. error, comparable
// and any are refused with generation.ErrUnsupported: the Go bridge exports
// them with the kinds of a declaration, without a marker that tells them from
// one, so their destination cannot be established yet.
func Predeclared(name string) (Symbol, error) {
	object, ok := types.Universe.Lookup(name).(*types.TypeName)
	if !ok {
		return Symbol{}, g.ErrInvalid
	}
	if _, basic := object.Type().(*types.Basic); !basic {
		return Symbol{}, g.ErrUnsupported
	}
	return Symbol{predeclared, name}, nil
}

// Predeclared returns the identifier of the universe scope, and false for a
// symbol of another category.
func (y Symbol) Predeclared() (string, bool) {
	return y.name, y.category == predeclared && y.name != ""
}

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
	name, ok := symbol.Predeclared()
	if !ok {
		return g.ExternalUse{}, g.ErrInvalid
	}
	issued, err := issuer.Issue(language, predeclared+"\x00"+name)
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
	category, name, ok := strings.Cut(key, "\x00")
	if !ok || category != predeclared {
		return Symbol{}, g.ErrInvalid
	}
	return Predeclared(name)
}
