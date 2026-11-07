package typescript

import "strings"

// PackageSpecifier accepts a bare package or exported subpath in the named-import profile.
func PackageSpecifier(s string) bool {
	if len(s) == 0 || len(s) > 1024 || strings.ContainsAny(s, "\\:#%?\x00") {
		return false
	}
	parts := strings.Split(s, "/")
	if s[0] == '@' && len(parts) < 2 {
		return false
	}
	for i, part := range parts {
		if i == 0 && strings.HasPrefix(part, "@") {
			part = part[1:]
		}
		if part == "" || part == "." || part == ".." || strings.HasPrefix(part, ".") {
			return false
		}
		for _, ch := range part {
			if !(ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' || strings.ContainsRune("-._", ch)) {
				return false
			}
		}
	}
	return true
}
