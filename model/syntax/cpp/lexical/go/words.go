package lexical

// The lists below state rules of the language's spelling as data. Each is
// published by a function without input that returns a fresh value, so that
// no importer can change it for another. Their evaluators are operations.

// Reserved lists the words that are no identifier.
func Reserved() []string {
	return []string{
		"alignas", "alignof", "and", "and_eq", "asm", "auto", "bitand", "bitor", "bool", "break",
		"case", "catch", "char", "char8_t", "char16_t", "char32_t", "class", "compl", "concept", "const",
		"consteval", "constexpr", "constinit", "const_cast", "continue", "co_await", "co_return", "co_yield", "decltype", "default",
		"delete", "do", "double", "dynamic_cast", "else", "enum", "explicit", "export", "extern", "false",
		"float", "for", "friend", "goto", "if", "inline", "int", "long", "mutable", "namespace",
		"new", "noexcept", "not", "not_eq", "nullptr", "operator", "or", "or_eq", "private", "protected",
		"public", "register", "reinterpret_cast", "requires", "return", "short", "signed", "sizeof", "static", "static_assert",
		"static_cast", "struct", "switch", "template", "this", "thread_local", "throw", "true", "try", "typedef",
		"typeid", "typename", "union", "unsigned", "using", "virtual", "void", "volatile", "wchar_t", "while",
		"xor", "xor_eq",
	}
}

// BuiltinWords lists the words of which the name of a builtin type is made.
func BuiltinWords() []string {
	return []string{
		"void", "bool", "char", "wchar_t", "char8_t", "char16_t", "char32_t", "signed", "unsigned", "short", "int", "long", "float", "double",
	}
}
