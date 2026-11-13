package rust

import (
	"github.com/Bitspark/bittype/model/syntax/rust/lexical/go"
	"sort"
)

// CrateName reports whether the pinned compiler takes name as the name of a
// crate: it is not empty, and every character is alphanumeric or an
// underscore. That is the rule of rustc 1.90.0,
// compiler/rustc_session/src/output.rs:57-79 (validate_crate_name):
//
//	if c.is_alphanumeric() || c == '_' { continue; }
//
// with char::is_alphanumeric of library/core/src/char/methods.rs:922-924,
// is_alphabetic (778-783) and is_numeric (1003-1008), which read the
// compiler's own Unicode tables.
//
// It is the one evaluator of that rule. The Rust environment judges the
// name of a captured crate and of a dependency by it, and the Rust symbol
// package the crate of a destination.
func CrateName(name string) bool {
	for _, r := range name {
		at := sort.Search(len(crateNameCharacters), func(i int) bool { return crateNameCharacters[i][1] >= r })
		if at == len(crateNameCharacters) || crateNameCharacters[at][0] > r {
			return false
		}
	}
	return name != ""
}

// crateNameCharacters is the table of that rule, which the lexical part of the
// Rust syntax model declares: ranges of code points in ascending order, each
// with its first and its last code point. The lexical part says how the table
// is regenerated for a new compiler pin.
var crateNameCharacters = lexical.CrateNameCharacters()
