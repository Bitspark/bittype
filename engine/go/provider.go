package engine

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Provider packages existing typed text capabilities for heterogeneous
// registration. No parser, formatter, checker or printer is implemented here.
type Provider struct {
	definition representation.Language
	parse      func(text.Input) (representation.CST, error)
	format     func(representation.AST) (representation.CST, error)
	print      func(representation.CST) (string, error)
}

func NewProvider[L any](provider text.Provider[L], profile string) (Provider, error) {
	definition, err := representation.NewLanguage[L](provider, profile)
	if err != nil {
		return Provider{}, err
	}
	result := Provider{definition: definition}
	result.parse = func(input text.Input) (representation.CST, error) {
		parsed, err := provider.Parse(input)
		if err != nil {
			return representation.CST{}, err
		}
		return representation.PackCST(definition, parsed)
	}
	result.format = func(input representation.AST) (representation.CST, error) {
		abstract, err := representation.Unpack[L](definition, input)
		if err != nil {
			return representation.CST{}, err
		}
		formatted, err := provider.Format(abstract)
		if err != nil {
			return representation.CST{}, err
		}
		return representation.PackCST(definition, formatted)
	}
	result.print = func(input representation.CST) (string, error) {
		concrete, err := representation.UnpackCST[L](definition, input)
		if err != nil {
			return "", err
		}
		output, err := provider.Print(concrete)
		if err != nil {
			return "", err
		}
		return output, nil
	}
	return result, nil
}
func (p Provider) ID() language.ID                   { return p.definition.ID() }
func (p Provider) Name() string                      { return p.definition.Name() }
func (p Provider) Language() representation.Language { return p.definition }
func (p Provider) validate() error {
	if err := p.definition.Validate(); err != nil {
		return err
	}
	if p.parse == nil || p.format == nil || p.print == nil {
		return ErrInvalidProvider
	}
	return nil
}

var _ language.LanguageProvider = Provider{}
