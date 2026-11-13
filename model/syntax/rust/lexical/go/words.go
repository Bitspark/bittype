package lexical

// The lists below state rules of the language's spelling as data. Each is
// published by a function without input that returns a fresh value, so that
// no importer can change it for another. Their evaluators are operations.

// Reserved lists the strict and the reserved keywords of edition 2024. None is
// an identifier, and none can start a path that names a crate.
func Reserved() []string {
	return []string{
		"as", "break", "const", "continue", "crate", "else", "enum", "extern", "false", "fn",
		"for", "if", "impl", "in", "let", "loop", "match", "mod", "move", "mut",
		"pub", "ref", "return", "self", "Self", "static", "struct", "super", "trait", "true",
		"type", "unsafe", "use", "where", "while", "async", "await", "dyn", "abstract", "become",
		"box", "do", "final", "macro", "override", "priv", "typeof", "unsized", "virtual", "yield",
		"try", "gen",
	}
}

// ABIs lists the names of the calling conventions.
func ABIs() []string {
	return []string{
		"Rust", "C", "C-unwind", "system", "system-unwind", "cdecl", "stdcall",
		"fastcall", "vectorcall", "thiscall", "aapcs", "win64", "sysv64", "efiapi",
	}
}
