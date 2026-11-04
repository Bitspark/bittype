package print

import (
	"errors"
	"strings"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var ErrInvalidPrinter = errors.New("invalid print provider")

// Printer binds one typed provider to a captured language/profile/marker.
// It owns boundary snapshots and never installs another capability.
type Printer struct {
	definition representation.Language
	run func(representation.CST) (string, error)
}

func NewPrinter[L any](definition representation.Language, provider Provider[L]) (Printer, error) {
	if representation.IsNil(provider) {
		return Printer{}, ErrInvalidPrinter
	}
	if err := representation.CheckMarker[L](definition); err != nil {
		return Printer{}, err
	}
	if provider.ID() != definition.ID() || strings.TrimSpace(provider.Name()) == "" {
		return Printer{}, ErrInvalidPrinter
	}
	result := Printer{definition: definition}
	result.run = func(input representation.CST) (string, error) {
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

func (p Printer) ID() language.ID { return p.definition.ID() }
func (p Printer) Name() string { return p.definition.Name() }
func (p Printer) Language() representation.Language { return p.definition }
func (p Printer) Validate() error {
	if p.run == nil {
		return ErrInvalidPrinter
	}
	return p.definition.Validate()
}
func (p Printer) Print(input representation.CST) (string, error) {
	if err := p.Validate(); err != nil {
		return "", err
	}
	return p.run(input)
}

var _ language.Provider = Printer{}
