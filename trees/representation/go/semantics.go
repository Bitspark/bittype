package representation

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/internal/tree"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"reflect"
)

// Semantics retains the language/profile and actual marker of an immutable graph.
// Admission never converts a graph to Semantics[any].
type Semantics struct {
	definition Language
	payload    semanticPayload
}

type semanticPayload interface {
	marker() reflect.Type
	Data() semantics.Data
	CheckLimits(semantics.Limits) error
	Entity(semantics.EntityID) (semantics.Entity, bool)
	Term(semantics.TermID) (semantics.Term, bool)
	Resolve(semantics.Occurrence) (semantics.EntityID, bool)
	TermAt(semantics.Occurrence) (semantics.TermID, bool)
	DependenciesFrom(semantics.EntityID, semantics.Key) []semantics.Dependency
	Facts(semantics.Subject, semantics.Key) []semantics.TermID
	DirectParameterBindings(semantics.TermID) ([]semantics.ArgumentBinding, bool)
}
type semanticTyped[L any] struct{ semantics.Semantics[L] }

func (semanticTyped[L]) marker() reflect.Type { return reflect.TypeFor[L]() }

func PackSemantics[L any](definition Language, view semantics.View[L]) (Semantics, error) {
	if err := CheckMarker[L](definition); err != nil {
		return Semantics{}, err
	}
	if IsNil(view) {
		return Semantics{}, semantics.ErrInvalid
	}
	graph := view.SemanticsView()
	if !graph.Valid() {
		return Semantics{}, semantics.ErrInvalid
	}
	return Semantics{definition, semanticTyped[L]{graph}}, nil
}
func UnpackSemantics[L any](definition Language, input Semantics) (semantics.Semantics[L], error) {
	if err := CheckMarker[L](definition); err != nil {
		return semantics.Semantics[L]{}, err
	}
	if !input.Valid() || !definition.Compatible(input.definition) {
		return semantics.Semantics[L]{}, ErrLanguageMismatch
	}
	value, ok := input.payload.(semanticTyped[L])
	if !ok {
		return semantics.Semantics[L]{}, ErrTypeMismatch
	}
	return value.Semantics, nil
}
func (s Semantics) Valid() bool             { return s.payload != nil }
func (s Semantics) LanguageID() language.ID { return s.definition.ID() }
func (s Semantics) Language() Language      { return s.definition }
func (s Semantics) Same(other Semantics) bool {
	return s.Valid() && s.definition.Compatible(other.definition) && s.payload == other.payload
}
func (s Semantics) Data() semantics.Data {
	if !s.Valid() {
		return semantics.Data{}
	}
	return s.payload.Data()
}
func (s Semantics) CheckLimits(l semantics.Limits) error {
	if !s.Valid() {
		return semantics.ErrInvalid
	}
	return s.payload.CheckLimits(l)
}
func (s Semantics) Entity(id semantics.EntityID) (semantics.Entity, bool) {
	if !s.Valid() {
		return semantics.Entity{}, false
	}
	return s.payload.Entity(id)
}
func (s Semantics) Term(id semantics.TermID) (semantics.Term, bool) {
	if !s.Valid() {
		return semantics.Term{}, false
	}
	return s.payload.Term(id)
}
func (s Semantics) Resolve(at semantics.Occurrence) (semantics.EntityID, bool) {
	if !s.Valid() {
		return 0, false
	}
	return s.payload.Resolve(at)
}
func (s Semantics) TermAt(at semantics.Occurrence) (semantics.TermID, bool) {
	if !s.Valid() {
		return 0, false
	}
	return s.payload.TermAt(at)
}
func (s Semantics) DependenciesFrom(id semantics.EntityID, key semantics.Key) []semantics.Dependency {
	if !s.Valid() {
		return nil
	}
	return s.payload.DependenciesFrom(id, key)
}
func (s Semantics) Facts(subject semantics.Subject, key semantics.Key) []semantics.TermID {
	if !s.Valid() {
		return nil
	}
	return s.payload.Facts(subject, key)
}
func (s Semantics) DirectParameterBindings(id semantics.TermID) ([]semantics.ArgumentBinding, bool) {
	if !s.Valid() {
		return nil, false
	}
	return s.payload.DirectParameterBindings(id)
}

// SameAST compares the complete captured tree, not just its tooling address.
func SameAST[L any](definition Language, input AST, view ast.View[L]) bool {
	value, err := Unpack[L](definition, input)
	return err == nil && !IsNil(view) && reflect.DeepEqual(value, tree.Clone(view.ASTView()))
}
