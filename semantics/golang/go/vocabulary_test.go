package golang

import (
	"testing"

	"github.com/Bitspark/bittype/semantics/go"
)

// The kind of a predeclared basic type and the role of a name binding are
// names of the published graph. A recorded predeclared destination is compared
// with an entity of that kind at a binding of that role, so a reader depends
// on both spellings.
func TestPublishedNamesOfThePredeclaredEntityAndTheNameBinding(t *testing.T) {
	for _, c := range []struct {
		name string
		got  semantics.Key
		want string
	}{
		{"BuiltinKind", BuiltinKind(), "builtin"},
		{"NameRole", NameRole(), "name"},
	} {
		if c.got != (semantics.Key{Owner: "go", Name: c.want}) {
			t.Errorf("%s is %v, want go/%s", c.name, c.got, c.want)
		}
	}
}
