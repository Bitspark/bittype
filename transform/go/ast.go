// Package transform bridges typed AST stages into an explicitly registered graph.
// It contains no language vocabulary or concrete transformer.
package transform

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/check/go"
	"github.com/Bitspark/bittype/internal/tree"
	"github.com/Bitspark/bittype/language/go"
)

var (
	ErrInvalidLanguage    = errors.New("invalid language descriptor")
	ErrLanguageMismatch   = errors.New("language or profile mismatch")
	ErrTypeMismatch       = errors.New("AST marker type mismatch")
	ErrInvalidAST         = errors.New("invalid AST envelope")
	ErrInvalidTransformer = errors.New("invalid AST transformer")
)

// AST is an owned, checked, opaque snapshot. Copying the envelope is safe;
// Unpack returns fresh mutable storage. Its zero value is invalid.
type AST struct {
	id      language.ID
	profile string
	payload payload
}

type payload interface {
	marker() reflect.Type
	clone() payload
}
type typed[L any] struct{ value ast.AST[L] }

func (p typed[L]) marker() reflect.Type { return reflect.TypeFor[L]() }
func (p typed[L]) clone() payload       { return typed[L]{tree.Clone(p.value)} }

func (a AST) LanguageID() language.ID { return a.id }
func (a AST) Profile() string         { return a.profile }

// Snapshot owns a deep copy, including attributes, references and origin slices.
func (a AST) Snapshot() (AST, error) {
	if a.payload == nil || a.id == "" || a.profile == "" {
		return AST{}, ErrInvalidAST
	}
	return AST{id: a.id, profile: a.profile, payload: a.payload.clone()}, nil
}

// Pack accepts either a typed AST or its existing CST view. It validates and
// snapshots storage; it never erases grammar, resolves names or prints source.
func Pack[L any](definition Language, view ast.View[L]) (AST, error) {
	if err := definition.acceptMarker(reflect.TypeFor[L]()); err != nil {
		return AST{}, err
	}
	if IsNil(view) {
		return AST{}, ErrInvalidAST
	}
	input := view.ASTView()
	if err := check.Structure(input); err != nil {
		return AST{}, fmt.Errorf("%w: %w", ErrInvalidAST, err)
	}
	result := AST{id: definition.id, profile: definition.profile, payload: typed[L]{tree.Clone(input)}}
	if err := definition.Check(result); err != nil {
		return AST{}, err
	}
	return result, nil
}

// Unpack checks the identity/profile/actual Go type and returns a deep copy.
// Grammar checking was performed when the envelope was constructed.
func Unpack[L any](definition Language, input AST) (ast.AST[L], error) {
	if err := definition.acceptMarker(reflect.TypeFor[L]()); err != nil {
		return ast.AST[L]{}, err
	}
	if err := definition.match(input); err != nil {
		return ast.AST[L]{}, err
	}
	return tree.Clone(input.payload.(typed[L]).value), nil
}

// IsNil recognizes nil interfaces and typed nil providers without invoking them.
func IsNil(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	}
	return false
}
