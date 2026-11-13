package source

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
)

// WithCategory is the one place that says which Files category a value of
// another package carries. Every formatter and the publisher call it.
func TestWithCategoryJoinsTheCategoryOfAForeignValue(t *testing.T) {
	other := errors.New("an error of no category")
	categories := []error{ErrInvalid, ErrUnsupported, semantic.ErrEvidenceUnavailable, ErrLimit}
	for _, example := range []struct {
		name string
		in   error
		want error
	}{
		{"a files limit", files.ErrLimit, ErrLimit},
		{"a wrapped files limit", fmt.Errorf("capture: %w", files.ErrLimit), ErrLimit},
		{"an invalid files value", files.ErrInvalid, ErrInvalid},
		{"a wrapped invalid files value", fmt.Errorf("capture: %w", files.ErrInvalid), ErrInvalid},
		{"a semantic limit", semantic.ErrLimit, ErrLimit},
		{"a wrapped semantic limit", fmt.Errorf("capture: %w", semantic.ErrLimit), ErrLimit},
		{"an invalid semantic request", semantic.ErrInvalidRequest, nil},
		{"another error", other, nil},
		{"a cancellation", fmt.Errorf("stopped: %w", other), nil},
	} {
		out := WithCategory(example.in)
		if !errors.Is(out, example.in) {
			t.Errorf("%s: the value is no longer matched: %q", example.name, out)
		}
		for _, category := range categories {
			if errors.Is(out, category) != (category == example.want) {
				t.Errorf("%s: errors.Is(%q, %q) = %v", example.name, out, category, errors.Is(out, category))
			}
		}
		if example.want == nil && out != example.in {
			t.Errorf("%s: an error without a category was changed: %q", example.name, out)
		}
	}
	if WithCategory(nil) != nil {
		t.Error("no error became an error")
	}
	// A refusal that has its category keeps it and is not joined again.
	for _, category := range categories {
		in := fmt.Errorf("%w: %w", category, files.ErrLimit)
		if out := WithCategory(in); out != in {
			t.Errorf("a refusal with the category %q was changed: %q", category, out)
		}
	}
}
