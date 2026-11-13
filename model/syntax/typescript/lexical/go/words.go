package lexical

// The lists below state rules of the language's spelling as data. Each is
// published by a function without input that returns a fresh value, so that
// no importer can change it for another. Their evaluators are operations.

// Reserved lists the words that are no identifier.
func Reserved() []string {
	return []string{
		"break", "case", "catch", "class", "const", "continue", "debugger", "default", "delete", "do",
		"else", "enum", "export", "extends", "false", "finally", "for", "function", "if", "import",
		"in", "instanceof", "new", "null", "return", "super", "switch", "throw", "true", "try",
		"typeof", "var", "void", "while", "with", "implements", "interface", "let", "package", "private",
		"protected", "public", "static", "yield", "await", "this",
	}
}

// PredefinedTypes lists the predefined type words. None names a type in a
// binding or an unqualified use.
func PredefinedTypes() []string {
	return []string{
		"any", "unknown", "never", "number", "bigint", "boolean", "string", "symbol", "void", "object", "undefined",
	}
}

// The two patterns of a number, as regular expressions in the syntax of Go's
// regexp package.
const (
	// NumberPattern is the text of a number literal, with an optional sign.
	NumberPattern = "^-?(?:(?:0[xX][0-9a-fA-F](?:_?[0-9a-fA-F])*|0[bB][01](?:_?[01])*|0[oO][0-7](?:_?[0-7])*)|(?:(?:0|[1-9](?:_?[0-9])*)(?:\\.(?:[0-9](?:_?[0-9])*)?)?|\\.[0-9](?:_?[0-9])*)(?:[eE][+-]?[0-9](?:_?[0-9])*)?)$"
	// BigintPattern is the text of a bigint literal, with an optional sign.
	BigintPattern = "^-?(?:0|[1-9](?:_?[0-9])*|0[xX][0-9a-fA-F](?:_?[0-9a-fA-F])*|0[bB][01](?:_?[01])*|0[oO][0-7](?:_?[0-7])*)n$"
)
