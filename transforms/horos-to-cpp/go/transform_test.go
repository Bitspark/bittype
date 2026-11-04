package cpp

import (
	"errors"
	"testing"

	"github.com/Bitspark/bittype/contracts/transform/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func TestRejectStandardNamespace(t *testing.T) {
	from, err := representation.NewLanguage[h.Language](h.Language{}, h.Profile)
	if err != nil {
		t.Fatal(err)
	}
	to, err := representation.NewLanguage[s.Language](s.Language{}, "native-types-v1")
	if err != nil {
		t.Fatal(err)
	}
	// C++ [namespace.std]/1 does not permit ordinary generated declarations in std.
	if mapper, err := New(from, to, Config{TreeID: "generated", Namespace: "std"}); !errors.Is(err, transform.ErrInvalidTransformer) || mapper != nil {
		t.Fatalf("reserved standard namespace accepted: %v", err)
	}
	if _, err := New(from, to, Config{TreeID: "generated", Namespace: "pagination"}); err != nil {
		t.Fatal(err)
	}
}
