// Package horos contains private Horos spelling rules, without semantic resolution.
package horos

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Identifier(s string) bool {
	if s == "" {
		return false
	}
	for i, c := range []byte(s) {
		if c != '_' && !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') && !(i > 0 && c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

func Reserved(s string) bool {
	switch s {
	case "module", "import", "from", "pin", "type", "self", "Type", "Data", "open", "any", "atom", "text", "int", "bool", "tuple", "compound", "list", "seq", "set", "map", "union", "and", "not", "wire", "exact", "rec", "ref", "true", "false", "null", "hex":
		return true
	}
	return false
}

func Name(s string) string {
	if Reserved(s) {
		return "`" + s + "`"
	}
	return s
}

func ModuleName(s string) string {
	parts := strings.Split(s, ".")
	for i := range parts {
		parts[i] = Name(parts[i])
	}
	return strings.Join(parts, ".")
}

// Strings use Go's double-quoted Unicode escapes, excluding byte/octal escapes
// that could denote malformed UTF-8. Raw atoms use #hex for arbitrary bytes.
func Unquote(s string) (string, error) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return "", fmt.Errorf("expected double-quoted string")
	}
	for i := 1; i < len(s)-1; i++ {
		if s[i] == '\\' {
			i++
			if i >= len(s)-1 || !strings.ContainsRune("abfnrtv\\\"uU", rune(s[i])) {
				return "", fmt.Errorf("unsupported string escape")
			}
		}
	}
	value, err := strconv.Unquote(s)
	if err != nil || !utf8.ValidString(value) {
		return "", fmt.Errorf("invalid Unicode string")
	}
	return value, nil
}

func Quote(s string) string {
	// strconv.Quote uses byte escapes for some controls; choose Unicode escapes
	// instead so every formatter-produced token belongs to our source grammar.
	q := strconv.Quote(s)
	var out strings.Builder
	for i := 0; i < len(q); i++ {
		if q[i] == '\\' && i+1 < len(q) {
			if q[i+1] == 'x' && i+3 < len(q) {
				out.WriteString("\\u00" + q[i+2:i+4])
				i += 3
			} else {
				out.WriteByte(q[i])
				i++
				out.WriteByte(q[i])
			}
		} else {
			out.WriteByte(q[i])
		}
	}
	return out.String()
}

// Bounded decimal spelling, independent of machine word size.
const MaxIntegerDigits = 4096

func Integer(s string) bool {
	if strings.HasPrefix(s, "-") {
		s = s[1:]
	}
	if len(s) == 0 || len(s) > MaxIntegerDigits || (len(s) > 1 && s[0] == '0') {
		return false
	}
	for _, c := range []byte(s) {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
