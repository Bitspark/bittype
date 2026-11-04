package engine

import (
	"fmt"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// attach changes only the selected capability of an existing language. Its
// install callback is engine-owned bookkeeping, never a provider callback.
func (e *Engine) attach(definition representation.Language, install func(*registration) error) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	r, ok := e.languages[definition.ID()]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownLanguage, definition.ID())
	}
	if !r.language.Compatible(definition) {
		return representation.ErrLanguageMismatch
	}
	if err := install(&r); err != nil {
		return err
	}
	e.languages[definition.ID()] = r
	return nil
}

// RegisterParser attaches only this capability; no other stage is inferred.
func (e *Engine) RegisterParser(provider parse.Parser) error {
	if err := provider.Validate(); err != nil {
		return err
	}
	return e.attach(provider.Language(), func(r *registration) error {
		if r.parser != nil {
			return fmt.Errorf("%w: %s", ErrDuplicateParser, provider.ID())
		}
		r.parser = &provider
		return nil
	})
}

// RegisterFormatter attaches only this capability; no other stage is inferred.
func (e *Engine) RegisterFormatter(provider format.Formatter) error {
	if err := provider.Validate(); err != nil {
		return err
	}
	return e.attach(provider.Language(), func(r *registration) error {
		if r.formatter != nil {
			return fmt.Errorf("%w: %s", ErrDuplicateFormatter, provider.ID())
		}
		r.formatter = &provider
		return nil
	})
}

// RegisterPrinter attaches only this capability; no other stage is inferred.
func (e *Engine) RegisterPrinter(provider print.Printer) error {
	if err := provider.Validate(); err != nil {
		return err
	}
	return e.attach(provider.Language(), func(r *registration) error {
		if r.printer != nil {
			return fmt.Errorf("%w: %s", ErrDuplicatePrinter, provider.ID())
		}
		r.printer = &provider
		return nil
	})
}

// RegisterChecker attaches only this capability; no other stage is inferred.
func (e *Engine) RegisterChecker(provider check.Checker) error {
	if err := provider.Validate(); err != nil {
		return err
	}
	return e.attach(provider.Language(), func(r *registration) error {
		if r.checker != nil {
			return fmt.Errorf("%w: %s", ErrDuplicateChecker, provider.ID())
		}
		r.checker = &provider
		return nil
	})
}
