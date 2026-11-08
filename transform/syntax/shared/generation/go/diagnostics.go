package generation

import (
	"errors"
	"fmt"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
)

func diagnostic(input transform.Input, err error) error {
	var diagnostic *text.Diagnostic
	if errors.As(err, &diagnostic) {
		return err
	}
	// A kernel limitation located at a source site is reported there.
	var located *g.SiteError
	if errors.As(err, &located) && errors.Is(err, g.ErrUnsupported) {
		site := located.Site
		return fmt.Errorf("%w: %w", &text.Diagnostic{Code: text.Unsupported, Message: located.Message, Node: &site}, err)
	}
	code := text.Unsupported
	if errors.Is(err, g.ErrLimit) {
		code = text.LimitExceeded
	} else if !errors.Is(err, g.ErrUnsupported) && !errors.Is(err, g.ErrInvalid) {
		return err
	}
	site := input.Source.Entry()
	return fmt.Errorf("%w: %w", &text.Diagnostic{Code: code, Message: "generation failed", Node: &site}, err)
}
