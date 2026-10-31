package typescript

import "unicode"

// ECMAScript WhiteSpace and LineTerminator, excluding unrelated Unicode spaces.
func Space(r rune)bool{
	return unicode.Is(unicode.Zs,r)||r=='\t'||r=='\v'||r=='\f'||r=='\r'||r=='\n'||r=='\u2028'||r=='\u2029'||r=='\ufeff'
}
