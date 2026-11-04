package format

import (
	"errors"
	"strings"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var ErrInvalidFormatter = errors.New("invalid format provider")

// Formatter binds one typed provider to a captured language/profile/marker.
// It owns boundary snapshots and never installs another capability.
type Formatter struct {
	definition representation.Language
	run func(representation.AST) (representation.CST, error)
}

func NewFormatter[L any](definition representation.Language, provider Provider[L]) (Formatter, error) {
	if representation.IsNil(provider) {
		return Formatter{}, ErrInvalidFormatter
	}
	if err := representation.CheckMarker[L](definition); err != nil {
		return Formatter{}, err
	}
	if provider.ID() != definition.ID() || strings.TrimSpace(provider.Name()) == "" {
		return Formatter{}, ErrInvalidFormatter
	}
	result := Formatter{definition: definition}
	result.run = func(input representation.AST) (representation.CST, error) {
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
	return result, nil
}

func (p Formatter) ID() language.ID { return p.definition.ID() }
func (p Formatter) Name() string { return p.definition.Name() }
func (p Formatter) Language() representation.Language { return p.definition }
func (p Formatter) Validate() error {
	if p.run == nil {
		return ErrInvalidFormatter
	}
	return p.definition.Validate()
}
func (p Formatter) Format(input representation.AST) (representation.CST, error) {
	if err := p.Validate(); err != nil {
		return representation.CST{}, err
	}
	return p.run(input)
}

var _ language.Provider = Formatter{}
