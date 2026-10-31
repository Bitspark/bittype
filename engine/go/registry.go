package engine

import (
	"fmt"
	"github.com/Bitspark/bittype/trees/representation/go"
	"sort"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
)

// RegisterLanguage installs an AST-only descriptor. RegisterProvider installs a
// descriptor and text capabilities together instead; their keys share one space.
func (e *Engine) RegisterLanguage(definition representation.Language) error {
	if err := definition.Validate(); err != nil {
		return err
	}
	return e.register(registration{language: definition})
}
func (e *Engine) RegisterProvider(provider Provider) error {
	if err := provider.validate(); err != nil {
		return err
	}
	return e.register(registration{language: provider.definition, text: &provider})
}
func (e *Engine) register(r registration) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	id := r.language.ID()
	if _, ok := e.languages[id]; ok {
		return fmt.Errorf("%w: %s", ErrDuplicateLanguage, id)
	}
	if e.languages == nil {
		e.languages = make(map[language.ID]registration)
	}
	e.languages[id] = r
	return nil
}

func (e *Engine) RegisterTransformer(t transform.ASTTransformer) error {
	if representation.IsNil(t) {
		return transform.ErrInvalidTransformer
	}
	from, to := t.FromID(), t.ToID()
	if from == "" || to == "" || from == to {
		return transform.ErrInvalidTransformer
	}
	var described bool
	var fromDefinition, toDefinition representation.Language
	if d, ok := t.(transform.DescribedTransformer); ok {
		described = true
		fromDefinition, toDefinition = d.FromLanguage(), d.ToLanguage()
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	f, ok := e.languages[from]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownLanguage, from)
	}
	dest, ok := e.languages[to]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownLanguage, to)
	}
	if described && (!f.language.Compatible(fromDefinition) || !dest.language.Compatible(toDefinition)) {
		return fmt.Errorf("%w: %w", transform.ErrInvalidTransformer, representation.ErrLanguageMismatch)
	}
	if _, ok := e.edges[from][to]; ok {
		return fmt.Errorf("%w: %s -> %s", ErrDuplicateTransformer, from, to)
	}
	if e.edges == nil {
		e.edges = make(map[language.ID]map[language.ID]edge)
	}
	if e.edges[from] == nil {
		e.edges[from] = make(map[language.ID]edge)
	}
	e.edges[from][to] = edge{from: from, to: to, transformer: t}
	return nil
}

// Languages returns a sorted, independent metadata slice.
func (e *Engine) Languages() []representation.Language {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]representation.Language, 0, len(e.languages))
	for _, r := range e.languages {
		result = append(result, r.language)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID() < result[j].ID() })
	return result
}
func (e *Engine) Language(id language.ID) (representation.Language, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	r, ok := e.languages[id]
	return r.language, ok
}
func (e *Engine) lookup(id language.ID) (registration, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	r, ok := e.languages[id]
	if !ok {
		return registration{}, fmt.Errorf("%w: %s", ErrUnknownLanguage, id)
	}
	return r, nil
}
