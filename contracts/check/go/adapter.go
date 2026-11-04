package check

import (
	"errors"
	"strings"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var ErrInvalidChecker = errors.New("invalid check provider")

// Checker binds one typed provider to a captured language/profile/marker.
// It owns boundary snapshots and never installs another capability.
type Checker struct {
	definition representation.Language
	run func(representation.AST) error
}

func NewChecker[L any](definition representation.Language, provider Provider[L]) (Checker, error) {
	if representation.IsNil(provider) {
		return Checker{}, ErrInvalidChecker
	}
	if err := representation.CheckMarker[L](definition); err != nil {
		return Checker{}, err
	}
	if provider.ID() != definition.ID() || strings.TrimSpace(provider.Name()) == "" {
		return Checker{}, ErrInvalidChecker
	}
	result := Checker{definition: definition}
	result.run = func(input representation.AST) error {
		abstract, err := representation.Unpack[L](definition, input)
		if err != nil {
			return err
		}
		return provider.Check(abstract)
	}
	return result, nil
}

func (p Checker) ID() language.ID { return p.definition.ID() }
func (p Checker) Name() string { return p.definition.Name() }
func (p Checker) Language() representation.Language { return p.definition }
func (p Checker) Validate() error {
	if p.run == nil {
		return ErrInvalidChecker
	}
	return p.definition.Validate()
}
func (p Checker) Check(input representation.AST) error {
	if err := p.Validate(); err != nil {
		return err
	}
	return p.run(input)
}

var _ language.Provider = Checker{}
