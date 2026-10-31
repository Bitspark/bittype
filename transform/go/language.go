package transform

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Bitspark/bittype/language/go"
	"github.com/Bitspark/bittype/text/go"
)

// Language binds an identity and explicit syntax/profile revision to a Go marker
// and checker. Metadata is captured at construction. Profiles are caller-owned
// compatibility promises, not hashes or automatically negotiated versions.
type Language struct {
	id            language.ID
	name, profile string
	marker        reflect.Type
	check         func(payload) error
}

func NewLanguage[L any](provider text.CheckerProvider[L], profile string) (Language, error) {
	if IsNil(provider) {
		return Language{}, ErrInvalidLanguage
	}
	id, name := provider.ID(), provider.Name()
	if strings.TrimSpace(string(id)) == "" || strings.TrimSpace(name) == "" || strings.TrimSpace(profile) == "" {
		return Language{}, ErrInvalidLanguage
	}
	return Language{id: id, name: name, profile: profile, marker: reflect.TypeFor[L](),
		check: func(p payload) error { return provider.Check(p.(typed[L]).value) }}, nil
}

func (l Language) ID() language.ID { return l.id }
func (l Language) Name() string    { return l.name }
func (l Language) Profile() string { return l.profile }
func (l Language) Validate() error {
	if l.id == "" || l.name == "" || l.profile == "" || l.marker == nil || l.check == nil {
		return ErrInvalidLanguage
	}
	return nil
}
func (l Language) Compatible(other Language) bool {
	return l.Validate() == nil && other.Validate() == nil && l.id == other.id && l.profile == other.profile && l.marker == other.marker
}
func (l Language) acceptMarker(marker reflect.Type) error {
	if err := l.Validate(); err != nil {
		return err
	}
	if marker != l.marker {
		return fmt.Errorf("%w for %s", ErrTypeMismatch, l.id)
	}
	return nil
}
func (l Language) match(input AST) error {
	if err := l.Validate(); err != nil {
		return err
	}
	if input.payload == nil {
		return ErrInvalidAST
	}
	if input.id != l.id || input.profile != l.profile {
		return fmt.Errorf("%w: expected %s (%s), got %s (%s)", ErrLanguageMismatch, l.id, l.profile, input.id, input.profile)
	}
	return l.acceptMarker(input.payload.marker())
}

// Check lends the checker a private copy. An incorrectly mutating checker cannot
// alter the caller's envelope. Provider code is trusted, not sandboxed.
func (l Language) Check(input AST) error {
	if err := l.match(input); err != nil {
		return err
	}
	return l.check(input.payload.clone())
}

var _ language.LanguageProvider = Language{}
