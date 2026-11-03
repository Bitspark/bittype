package declsyntax

import (
	"strings"
	"testing"
)

func TestModuleNameSpellingAndAllocationBound(t *testing.T) {
	for _, name := range []string{"a", "example.pagination", "_private.Type2"} {
		if !moduleName(name) {
			t.Fatalf("valid module name rejected: %q", name)
		}
	}
	for _, name := range []string{"", ".a", "a.", "a..b", "a.2b", "a/b", "a.\u00e9"} {
		if moduleName(name) {
			t.Fatalf("invalid module name accepted: %q", name)
		}
	}
	// A name can occupy most of the AST text budget. Rejecting empty components
	// must not allocate a string-header array proportional to the number of dots.
	name := strings.Repeat(".", 64<<10)
	accepted := false
	allocations := testing.AllocsPerRun(3, func() { accepted = moduleName(name) })
	if accepted || allocations != 0 {
		t.Fatalf("dot-heavy name rejection allocated: accepted=%v allocations=%v", accepted, allocations)
	}
}
