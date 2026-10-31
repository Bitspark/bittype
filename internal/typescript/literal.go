package typescript

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

var number = regexp.MustCompile("^-?(?:(?:0[xX][0-9a-fA-F](?:_?[0-9a-fA-F])*|0[bB][01](?:_?[01])*|0[oO][0-7](?:_?[0-7])*)|(?:(?:0|[1-9](?:_?[0-9])*)(?:\\.(?:[0-9](?:_?[0-9])*)?)?|\\.[0-9](?:_?[0-9])*)(?:[eE][+-]?[0-9](?:_?[0-9])*)?)$")
var bigint = regexp.MustCompile("^-?(?:0|[1-9](?:_?[0-9])*|0[xX][0-9a-fA-F](?:_?[0-9a-fA-F])*|0[bB][01](?:_?[01])*|0[oO][0-7](?:_?[0-7])*)n$")

func Number(s string) bool  { return number.MatchString(s) }
func Bigint(s string) bool  { return bigint.MatchString(s) }
func Quote(s string) string { b, _ := json.Marshal(s); return string(b) }
func Unquote(s string) (string, error) {
	bad := func() (string, error) { return "", fmt.Errorf("invalid string literal") }
	if len(s) < 2 || (s[0] != '"' && s[0] != '\'') || s[len(s)-1] != s[0] {
		return bad()
	}
	quote := s[0]
	body := s[1 : len(s)-1]
	var out strings.Builder
	for i := 0; i < len(body); {
		if body[i] == quote || body[i] == '\r' || body[i] == '\n' {
			return bad()
		}
		if body[i] != '\\' {
			r, n := utf8.DecodeRuneInString(body[i:])
			if r == utf8.RuneError && n == 1 {
				return bad()
			}
			out.WriteRune(r)
			i += n
			continue
		}
		i++
		if i == len(body) {
			return bad()
		}
		c := body[i]
		i++
		switch c {
		case '\n':
			continue
		case '\r':
			if i < len(body) && body[i] == '\n' {
				i++
			}
			continue
		case 'b':
			out.WriteByte('\b')
		case 'f':
			out.WriteByte('\f')
		case 'n':
			out.WriteByte('\n')
		case 'r':
			out.WriteByte('\r')
		case 't':
			out.WriteByte('\t')
		case 'v':
			out.WriteByte('\v')
		case '0':
			if i < len(body) && body[i] >= '0' && body[i] <= '9' {
				return bad()
			}
			out.WriteByte(0)
		case 'x', 'u':
			count := 2
			if c == 'u' {
				count = 4
			}
			start := i
			braced := false
			if c == 'u' && i < len(body) && body[i] == '{' {
				braced = true
				i++
				start = i
				for i < len(body) && body[i] != '}' {
					i++
				}
				count = i - start
				if i == len(body) || count < 1 || count > 6 {
					return bad()
				}
			}
			if start+count > len(body) {
				return bad()
			}
			value, err := strconv.ParseUint(body[start:start+count], 16, 32)
			if err != nil {
				return bad()
			}
			if braced {
				i++
			} else {
				i = start + count
			}
			r := rune(value)
			if r >= 0xd800 && r <= 0xdbff && !braced && strings.HasPrefix(body[i:], "\\u") && i+6 <= len(body) {
				second, err := strconv.ParseUint(body[i+2:i+6], 16, 16)
				if err != nil {
					return bad()
				}
				r = utf16.DecodeRune(r, rune(second))
				i += 6
			}
			if !utf8.ValidRune(r) || r == utf8.RuneError {
				return bad()
			}
			out.WriteRune(r)
		default:
			if c >= '1' && c <= '9' {
				return bad()
			}
			out.WriteByte(c)
		}
	}
	return out.String(), nil
}
