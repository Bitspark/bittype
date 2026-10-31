package typescript

import (
	"strings"
	"unicode"
)

// ECMAScript WhiteSpace and LineTerminator, excluding unrelated Unicode spaces.
func Space(r rune) bool {
	return unicode.Is(unicode.Zs, r) || r == '\t' || r == '\v' || r == '\f' || r == '\r' || r == '\n' || r == '\u2028' || r == '\u2029' || r == '\ufeff'
}

func LineTerminator(r rune) bool { return r == '\r' || r == '\n' || r == '\u2028' || r == '\u2029' }

func Directive(comment string) bool {
	if strings.HasPrefix(comment, "///") {
		return true
	}
	body := strings.TrimPrefix(strings.TrimPrefix(comment, "//"), "/*")
	for _, line := range strings.FieldsFunc(body, LineTerminator) {
		line = strings.TrimLeftFunc(line, Space)
		line = strings.TrimLeftFunc(strings.TrimPrefix(line, "*"), Space)
		if strings.HasPrefix(line, "@ts-") {
			return true
		}
	}
	return false
}
