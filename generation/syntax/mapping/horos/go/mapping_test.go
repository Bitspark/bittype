package generation_test

import (
	"testing"

	rules "github.com/Bitspark/bittype/generation/syntax/mapping/horos/go"
)

func TestStructuralNamingPolicy(t *testing.T) {
	for _, name := range []string{"", "_", "a__b", "a-b", "Ã©"} {
		if _, err := rules.Exported(name); err == nil {
			t.Fatal("unsupported identifier accepted", name)
		}
	}
	for _, pair := range [][2]string{{"page", "Page"}, {"_page", "Type_page"}, {"Self", "Self_"}, {"Page", "Page"}} {
		if got, err := rules.Exported(pair[0]); err != nil || got != pair[1] {
			t.Fatal(got, err)
		}
	}
}
