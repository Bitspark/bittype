// Package naming holds the rule by which a horos name becomes an exported
// native identifier. The horos side of a mapping owns this rule: the
// declaration recipe and a mapping that names fields apply the same function,
// so both spell a name alike. It is one rule for every target.
package naming

import (
	"fmt"
	"strings"
)

// Exported returns the exported native identifier for a horos name, or an
// error for a name it does not support.
func Exported(name string) (string, error) {
	if name == "" || name == "_" || strings.Contains(name, "__") {
		return "", fmt.Errorf("unsupported identifier %q", name)
	}
	for i, b := range []byte(name) {
		if !(b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b == '_' || i > 0 && b >= '0' && b <= '9') {
			return "", fmt.Errorf("unsupported identifier %q", name)
		}
	}
	if name[0] == '_' {
		name = "Type" + name
	} else if name[0] >= 'a' && name[0] <= 'z' {
		name = string(name[0]-32) + name[1:]
	}
	if name == "Self" {
		name += "_"
	}
	return name, nil
}
