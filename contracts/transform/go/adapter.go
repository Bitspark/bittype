package transform

import (
	"github.com/Bitspark/bittype/trees/representation/go"

	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/pipeline/go"
)

// NewTransformer adapts a typed stage without changing its AST representation.
// The stage receives owned storage and its result is checked and snapshotted.
// This adapter contains no type mapping or behavioral contract binding.
func NewTransformer[F, T any](from, to representation.Language, stage pipeline.Stage[ast.AST[F], ast.AST[T]]) (ASTTransformer, error) {
	if stage == nil {
		return nil, ErrInvalidTransformer
	}
	if err := representation.CheckMarker[F](from); err != nil {
		return nil, err
	}
	if err := representation.CheckMarker[T](to); err != nil {
		return nil, err
	}
	return bound[F, T]{from: from, to: to, stage: stage}, nil
}

type bound[F, T any] struct {
	from, to representation.Language
	stage    pipeline.Stage[ast.AST[F], ast.AST[T]]
}

func (b bound[F, T]) FromID() language.ID    { return b.from.ID() }
func (b bound[F, T]) ToID() language.ID      { return b.to.ID() }
func (b bound[F, T]) FromLanguage() representation.Language { return b.from }
func (b bound[F, T]) ToLanguage() representation.Language   { return b.to }
func (b bound[F, T]) Transform(input representation.AST) (representation.AST, error) {
	if err := b.from.Check(input); err != nil {
		return representation.AST{}, err
	}
	source, err := representation.Unpack[F](b.from, input)
	if err != nil {
		return representation.AST{}, err
	}
	result, err := b.stage(source)
	if err != nil {
		return representation.AST{}, err
	}
	return representation.Pack(b.to, result)
}
