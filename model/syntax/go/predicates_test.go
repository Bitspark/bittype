package syntax_test

import (
	"testing"

	"github.com/Bitspark/bittype/model/syntax/go"
)

// A name is an identity: a table of any language may hold it, and an operation
// binds an evaluator to it. The six shared names are pinned here, so a changed
// string fails a test and is no silent change of every table that names it.
func TestSharedNamesArePinned(t *testing.T) {
	names := map[string]syntax.Predicate{
		"bool":     syntax.Bool(),
		"decimal":  syntax.Decimal(),
		"nonempty": syntax.Nonempty(),
		"positive": syntax.Positive(),
		"signed":   syntax.Signed(),
		"string":   syntax.String(),
	}
	if syntax.Owner != "bittype" {
		t.Errorf("the owner of the shared names is %q", syntax.Owner)
	}
	seen := map[syntax.Predicate]string{}
	for want, name := range names {
		if name != (syntax.Predicate{Owner: "bittype", Name: want}) {
			t.Errorf("%s is %s/%s", want, name.Owner, name.Name)
		}
		if other, twice := seen[name]; twice {
			t.Errorf("%s and %s are one name", want, other)
		}
		seen[name] = want
	}
}
