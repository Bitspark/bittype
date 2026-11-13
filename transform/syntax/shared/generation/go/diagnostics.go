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
		// A source diagnostic that states what is wrong and not where is
		// reported at the site at which the kernel located the failure.
		var located *g.SiteError
		if diagnostic.Node == nil && diagnostic.Span == nil && errors.As(err, &located) && located.Site.Tree != "" {
			site := located.Site
			return fmt.Errorf("%w: %w", &text.Diagnostic{Code: diagnostic.Code, Message: diagnostic.Message, Node: &site}, err)
		}
		return err
	}
	// A refusal of a contribution or of the planner is reported at the site of
	// its subject, with its owner's reason. The class stays the refusal's.
	var refusal *g.Refusal
	if errors.As(err, &refusal) && refusal.Site.Tree != "" {
		code := text.Unsupported
		if refusal.Class == g.RefusalLimit {
			code = text.LimitExceeded
		}
		site := refusal.Site
		return fmt.Errorf("%w: %w", &text.Diagnostic{Code: code, Message: refusal.Reason, Node: &site}, err)
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
