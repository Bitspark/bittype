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

// The relation from the entity of a type name to its own type is a name of
// the published graph: the bridge writes it and Lookup, LookupIn, Parameters
// and DeclaredType read it. An entity's body is not that type.
func TestPublishedNameOfTheDeclaredTypeRelation(t *testing.T) {
	if got := DeclaredTypeRelation(); got != (semantics.Key{Owner: "go", Name: "declared-type"}) {
		t.Errorf("DeclaredTypeRelation is %v, want go/declared-type", got)
	}
}
