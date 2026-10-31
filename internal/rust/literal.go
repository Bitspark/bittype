package rust

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Space(r rune) bool {
	return (r >= '\t' && r <= '\r') || r == ' ' || r == '\u0085' || r == '\u200e' || r == '\u200f' || r == '\u2028' || r == '\u2029'
}

// UnquoteABI implements the admitted Rust string escapes, not Go escapes.
func UnquoteABI(s string) (string, error) {
	bad := func() (string, error) { return "", fmt.Errorf("invalid Rust ABI string") }
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return bad()
	}
	body := s[1 : len(s)-1]
	var out strings.Builder
	for i := 0; i < len(body); {
		r, n := utf8.DecodeRuneInString(body[i:])
		i += n
		if r == '"' || r == '\r' || r == '\n' || r == utf8.RuneError && n == 1 {
			return bad()
		}
		if r != '\\' {
			out.WriteRune(r)
			continue
		}
		if i == len(body) {
			return bad()
		}
		c := body[i]
		i++
		switch c {
		case '"', '\'', '\\':
			out.WriteByte(c)
		case 'n':
			out.WriteByte('\n')
		case 'r':
			out.WriteByte('\r')
		case 't':
			out.WriteByte('\t')
		case '0':
			out.WriteByte(0)
		case 'x':
			if i+2 > len(body) {
				return bad()
			}
			v, err := strconv.ParseUint(body[i:i+2], 16, 8)
			if err != nil || v > 127 {
				return bad()
			}
			out.WriteByte(byte(v))
			i += 2
		case 'u':
			if i >= len(body) || body[i] != '{' {
				return bad()
			}
			i++
			start := i
			for i < len(body) && body[i] != '}' {
				i++
			}
			if i == len(body) {
				return bad()
			}
			digits := strings.ReplaceAll(body[start:i], "_", "")
			i++
			if len(digits) < 1 || len(digits) > 6 {
				return bad()
			}
			v, err := strconv.ParseUint(digits, 16, 32)
			if err != nil || !utf8.ValidRune(rune(v)) {
				return bad()
			}
			out.WriteRune(rune(v))
		default:
			return bad()
		}
	}
	return out.String(), nil
}
func ABI(s string) bool {
	switch s {
	case "Rust", "C", "C-unwind", "system", "system-unwind", "cdecl", "stdcall", "fastcall", "vectorcall", "thiscall", "aapcs", "win64", "sysv64", "efiapi":
		return true
	}
	return false
}
