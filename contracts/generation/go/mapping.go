package generation

import "github.com/Bitspark/bittype/trees/representation/go"

// Selector is an exact source-owned identity, qualified by language/profile/marker.
// Its key is never interpreted by the composer. Display names and graph-local IDs
// are not substitutes for a source integration's stable declaration identity.
type Selector struct {
	language representation.Language
	key      string
}

func NewSelector(language representation.Language, key string) (Selector, error) {
	if language.Validate() != nil || key == "" {
		return Selector{}, ErrInvalid
	}
	return Selector{language, key}, nil
}
func (s Selector) Key() string                       { return s.key }
func (s Selector) Valid() bool                       { return s.language.Validate() == nil && s.key != "" }
func (s Selector) Language() representation.Language { return s.language }

// Mapping creates a plan for a real declaration, or a direct expression when
// declaration is zero. Factories allocate invocation-local collaborators only;
// representation work is demanded through the returned SPI operations.
type Mapping[F, T any] struct {
	ID       string
	PluginID string
	Plan     func(Declaration[F]) (Plan[F, T], error)
}
type RepresentationBinding[F, T any] struct {
	Selector Selector
	Mapping  Mapping[F, T]
}

// Catalog is source-owned interpretation; enumeration does not demand expansion.
type Catalog[F any] interface {
	Declarations() []Declaration[F]
	DeclarationSelector(Declaration[F]) (Selector, bool, error)
	UseSelector(UseData[F]) (Selector, error)
}
