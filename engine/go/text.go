package engine

import (
    "fmt"
    "github.com/Bitspark/bittype/contracts/language/go"
    "github.com/Bitspark/bittype/contracts/text/go"
    "github.com/Bitspark/bittype/trees/representation/go"
)

func (e *Engine) textProvider(id language.ID) (Provider, error) {
	registration, err := e.lookup(id)
	if err != nil {
		return Provider{}, err
	}
	if registration.text == nil {
		return Provider{}, fmt.Errorf("%w: %s", ErrMissingCapability, id)
	}
	return *registration.text, nil
}
func (e *Engine) Parse(id language.ID, input text.Input) (representation.CST, error) {
	provider, err := e.textProvider(id)
	if err != nil {
		return representation.CST{}, err
	}
	return provider.parse(input)
}
func (e *Engine) Format(input representation.AST) (representation.CST, error) {
	provider, err := e.textProvider(input.LanguageID())
	if err != nil {
		return representation.CST{}, err
	}
	if err := provider.definition.Check(input); err != nil {
		return representation.CST{}, err
	}
	return provider.format(input)
}
func (e *Engine) Print(input representation.CST) (string, error) {
	provider, err := e.textProvider(input.LanguageID())
	if err != nil {
		return "", err
	}
	return provider.print(input)
}
func (e *Engine) Render(input representation.AST) (string, error) {
	formatted, err := e.Format(input)
	if err != nil {
		return "", err
	}
	return e.Print(formatted)
}
