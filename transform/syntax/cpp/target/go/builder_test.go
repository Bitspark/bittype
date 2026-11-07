package target_test

import (
	"context"
	"errors"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	target "github.com/Bitspark/bittype/transform/syntax/cpp/target/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func TestRejectReservedNamespace(t *testing.T) {
	language, _ := representation.NewLanguage[s.Language](s.Language{}, "native-types-v1")
	for _, namespace := range []string{"", "std", "_reserved", "class", "a-b"} {
		backend, err := target.New[struct{}](context.Background(), language, target.Config{TreeID: "output", Namespace: namespace}, g.DefaultLimits())
		if backend != nil || !errors.Is(err, g.ErrInvalid) {
			t.Fatal("reserved namespace accepted", namespace, err)
		}
	}
	if _, err := target.New[struct{}](context.Background(), language, target.Config{TreeID: "output", Namespace: "pagination"}, g.DefaultLimits()); err != nil {
		t.Fatal(err)
	}
}
