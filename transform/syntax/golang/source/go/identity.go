// Package source interprets established Go semantics for the generation
// kernel and for the mapping contributions of a pair whose source is Go.
//
// It is the Go source reader of the two other pairs of bittype#196 (decision
// O43). It reads the analysed vocabulary of Go (semantics/golang/go) through
// the kernel's walk, and it reads a small fragment: a defined type without
// type parameters whose underlying type is a struct type, the struct type's
// fields, the basic types string and bool, and a use of a defined type.
// Everything else it judges as unsupported.
package source

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// The names of the basic types this reader reads, as Go spells them.
const (
	BasicString = "string"
	BasicBool   = "bool"
)

// selector issues a source identity of Go. The key describes Go's grammar or
// names a Go declaration, never a selected mapping. That a selector of another
// language selects nothing is the kernel's check; this function only refuses a
// language value that is no Go.
func selector(language representation.Language, key string) (g.Selector, error) {
	if err := representation.CheckMarker[s.Language](language); err != nil {
		return g.Selector{}, g.ErrInvalid
	}
	return g.NewSelector(language, key)
}

// StructType is the source identity of the form "struct type".
func StructType(language representation.Language) (g.Selector, error) {
	return selector(language, "type/struct")
}

// Basic is the source identity of one basic type, by its Go name.
func Basic(language representation.Language, name string) (g.Selector, error) {
	if name != BasicString && name != BasicBool {
		return g.Selector{}, g.ErrInvalid
	}
	return selector(language, "type/basic/"+name)
}

// Defined is the source identity of a defined type: the import path of its
// package and its name.
func Defined(language representation.Language, packagePath, name string) (g.Selector, error) {
	if packagePath == "" || name == "" {
		return g.Selector{}, g.ErrInvalid
	}
	return selector(language, "declaration/"+packagePath+"/"+name)
}
