package lexical

// The lists below state rules of the language's spelling as data. Each is
// published by a function without input that returns a fresh value, so that
// no importer can change it for another. Their evaluators are operations.

// Reserved lists the words of the language. A name that is one of them is
// written between backquotes.
func Reserved() []string {
	return []string{
		"module", "import", "from", "pin", "type", "kind", "self", "top", "bottom", "node",
		"nu", "in", "at", "epsilon", "item", "concat", "alternative", "star", "open", "any",
		"atom", "text", "int", "bool", "tuple", "compound", "list", "seq", "set", "map",
		"union", "and", "not", "wire", "exact", "rec", "ref", "true", "false", "null",
		"hex",
	}
}

// MaxIntegerDigits bounds the decimal digits of an integer, independent of
// machine word size.
const MaxIntegerDigits = 4096
