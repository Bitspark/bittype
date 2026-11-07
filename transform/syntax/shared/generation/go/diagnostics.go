package generation

import (
	"context"
	"errors"
	"fmt"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// withDiagnostics associates operation failures with the captured input. It does not
// turn failed mappings into a retry or a weaker representation.
func withDiagnostics(mapper transform.ASTTransformer, from, to representation.Language) transform.ASTTransformer {
	return diagnosticTransformer{ASTTransformer: mapper, from: from, to: to}
}

type diagnosticTransformer struct {
	transform.ASTTransformer
	from, to representation.Language
}

func (d diagnosticTransformer) FromLanguage() representation.Language { return d.from }
func (d diagnosticTransformer) ToLanguage() representation.Language   { return d.to }
func (d diagnosticTransformer) Transform(ctx context.Context, input transform.Input) (representation.ASTSet, error) {
	out, err := d.ASTTransformer.Transform(ctx, input)
	if err == nil {
		return out, nil
	}
	var diagnostic *text.Diagnostic
	if errors.As(err, &diagnostic) {
		return representation.ASTSet{}, err
	}
	code := text.Unsupported
	if errors.Is(err, g.ErrLimit) {
		code = text.LimitExceeded
	} else if !errors.Is(err, g.ErrUnsupported) && !errors.Is(err, g.ErrInvalid) {
		return representation.ASTSet{}, err
	}
	site := input.Source.Entry()
	return representation.ASTSet{}, fmt.Errorf("%w: %w", &text.Diagnostic{Code: code, Message: "generation failed", Node: &site}, err)
}
