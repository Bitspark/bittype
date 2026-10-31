package integration_test

import (
	"go/types"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/language/go"
	"github.com/Bitspark/bittype/syntax/cpp/go"
	"github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/syntax/typescript/go"
)

func TestLanguageIdentityInterface(t *testing.T) {
	// These calls exercise only the common interface, as a future registry would.
	cases := []struct {
		provider language.Language
		id language.ID
		name string
	}{
		{typescript.Language{}, "typescript", "TypeScript"},
		{golang.Language{}, "go", "Go"},
		{rust.Language{}, "rust", "Rust"},
		{cpp.Language{}, "cpp", "C++"},
	}
	seen := map[language.ID]bool{}
	for _, test := range cases {
		t.Run(string(test.id), func(t *testing.T) {
			if got := test.provider.ID(); got != test.id {
				t.Fatalf("registration identity: got %q, want %q", got, test.id)
			}
			if got := test.provider.Name(); got != test.name {
				t.Fatalf("display name: got %q, want %q", got, test.name)
			}
			if seen[test.provider.ID()] {
				t.Fatalf("duplicate registration identity: %q", test.provider.ID())
			}
			seen[test.provider.ID()] = true
		})
	}
}

func TestSyntaxConsumers(t *testing.T) {
	packages := loadContracts(t)
	for _, name := range []string{"typescript", "golang", "rust", "cpp", "extension"} {
		t.Run(name, func(t *testing.T) {
			if err := checkConsumer(t, packages, "syntax/"+name); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestTargetLanguageTypesRemainDistinct(t *testing.T) {
	packages := loadContracts(t)
	languages := []string{"typescript", "golang", "rust", "cpp"}
	for _, from := range languages {
		for _, to := range languages {
			if from == to {
				continue
			}
			for _, name := range []string{"Language", "AST", "Node", "Kind"} {
				source := packages[modulePath+"/syntax/"+from+"/go"].Scope().Lookup(name)
				target := packages[modulePath+"/syntax/"+to+"/go"].Scope().Lookup(name)
				if source == nil || target == nil {
					t.Fatalf("%s or %s does not expose %s", from, to, name)
				}
				if types.Identical(source.Type(), target.Type()) || types.AssignableTo(source.Type(), target.Type()) {
					t.Errorf("%s.%s can be confused with %s.%s", from, name, to, name)
				}
			}
		}
	}
}

func TestMixedTargetConsumerIsRejected(t *testing.T) {
	diagnostics := consumerDiagnostics(t, loadContracts(t), "syntax/mixed_languages")
	if len(diagnostics) != 4 {
		t.Fatalf("want four incompatible AST assignments, got %v", diagnostics)
	}
	for _, diagnostic := range diagnostics {
		if !strings.Contains(diagnostic.Error(), "cannot use") {
			t.Errorf("unexpected consumer diagnostic: %v", diagnostic)
		}
	}
}
