// Package source interprets established TypeScript semantics for the
// generation kernel and for the mapping contributions of a pair whose source
// is TypeScript.
//
// It is the TypeScript source reader of the two other pairs of bittype#196
// (decision O43). It reads the analysed vocabulary of TypeScript
// (semantics/typescript/go) through the kernel's walk, and it reads a small
// fragment: an interface or a type alias without type parameters of a
// captured module, an object type and its properties, and the primitive
// types. Everything else it judges as unsupported: a generic declaration, a
// reference to a named type, a type of a package, and every other type
// expression (literals, arrays, tuples, unions, intersections).
package source

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// The names of the primitive types that contributions of this module name, as
// TypeScript spells them. The reader names every primitive type the analysis
// publishes by its TypeScript spelling, these three among them.
const (
	PrimitiveString  = "string"
	PrimitiveBoolean = "boolean"
	PrimitiveNumber  = "number"
)

// selector issues a source identity of TypeScript. The key describes
// TypeScript's grammar or names a TypeScript declaration, never a selected
// mapping.
func selector(language representation.Language, key string) (g.Selector, error) {
	if err := representation.CheckMarker[s.Language](language); err != nil {
		return g.Selector{}, g.ErrInvalid
	}
	return g.NewSelector(language, key)
}

// ObjectType is the source identity of the form "object type": the members of
// an interface or of an object type literal.
func ObjectType(language representation.Language) (g.Selector, error) {
	return selector(language, "type/object")
}

// Primitive is the source identity of one primitive type, by its TypeScript
// spelling.
func Primitive(language representation.Language, name string) (g.Selector, error) {
	if name == "" {
		return g.Selector{}, g.ErrInvalid
	}
	return selector(language, "type/primitive/"+name)
}

// Named is the source identity of an interface or a type alias: the path of
// its module and its name.
func Named(language representation.Language, module, name string) (g.Selector, error) {
	if module == "" || name == "" {
		return g.Selector{}, g.ErrInvalid
	}
	return selector(language, "declaration/"+module+"/"+name)
}
