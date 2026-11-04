package parse

import (
	"errors"
	"strings"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var ErrInvalidParser = errors.New("invalid parse provider")

// Parser binds one typed provider to a captured language/profile/marker.
// It owns boundary snapshots and never installs another capability.
type Parser struct {
	definition representation.Language
	run func(Input) (representation.CST, error)
}

func NewParser[L any](definition representation.Language, provider Provider[L]) (Parser, error) {
	if representation.IsNil(provider) {
		return Parser{}, ErrInvalidParser
	}
	if err := representation.CheckMarker[L](definition); err != nil {
		return Parser{}, err
	}
	if provider.ID() != definition.ID() || strings.TrimSpace(provider.Name()) == "" {
		return Parser{}, ErrInvalidParser
	}
	result := Parser{definition: definition}
	result.run = func(input Input) (representation.CST, error) {
		parsed, err := provider.Parse(input)
		if err != nil {
			return representation.CST{}, err
		}
		return representation.PackCST(definition, parsed)
	}
	return result, nil
}

func (p Parser) ID() language.ID { return p.definition.ID() }
func (p Parser) Name() string { return p.definition.Name() }
func (p Parser) Language() representation.Language { return p.definition }
func (p Parser) Validate() error {
	if p.run == nil {
		return ErrInvalidParser
	}
	return p.definition.Validate()
}
func (p Parser) Parse(input Input) (representation.CST, error) {
	if err := p.Validate(); err != nil {
		return representation.CST{}, err
	}
	return p.run(input)
}

var _ language.Provider = Parser{}
