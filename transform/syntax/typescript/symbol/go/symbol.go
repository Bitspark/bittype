// Package symbol owns the TypeScript public type symbols shared by generation and Files.
package symbol

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	lit "github.com/Bitspark/bittype/internal/typescript"
	"github.com/Bitspark/bittype/semantics/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"strings"
)

// Type identifies a named export, or a global type when Module is empty.
type Type struct{ Module, Name string }

var issuer = g.NewSymbolIssuer()
var nameRole = semantics.Key{Owner: "typescript", Name: "name"}

func ExternalUse(language representation.Language, site ast.NodeRef, export Type) (g.ExternalUse, error) {
	if representation.CheckMarker[ts.Language](language) != nil || language.ID() != (ts.Language{}).ID() {
		return g.ExternalUse{}, g.ErrScope
	}
	if (export.Module != "" && !lit.PackageSpecifier(export.Module)) || !lit.TypeName(export.Name) {
		return g.ExternalUse{}, g.ErrInvalid
	}
	symbol, err := issuer.Issue(language, export.Module+"\x00"+export.Name)
	if err != nil {
		return g.ExternalUse{}, err
	}
	return g.ExternalUse{At: semantics.Occurrence{Site: site, Role: nameRole}, Symbol: symbol}, nil
}
func ReadExternal(language representation.Language, use g.ExternalUse) (Type, error) {
	key, ok := issuer.Read(language, use.Symbol)
	if !ok || use.At.Role != nameRole || use.At.Ordinal != 0 {
		return Type{}, g.ErrScope
	}
	module, name, ok := strings.Cut(key, "\x00")
	if !ok {
		return Type{}, g.ErrInvalid
	}
	return Type{Module: module, Name: name}, nil
}
