package transform

import (
 "reflect"

 "github.com/Bitspark/bittype/ast/go"
 "github.com/Bitspark/bittype/language/go"
 "github.com/Bitspark/bittype/pipeline/go"
)

// NewTransformer adapts a typed stage without changing its AST representation.
// The stage receives owned storage and its result is checked and snapshotted.
// This adapter contains no type mapping or behavioral contract binding.
func NewTransformer[F, T any](from, to Language, stage pipeline.Stage[ast.AST[F], ast.AST[T]]) (ASTTransformer, error) {
 if stage == nil { return nil, ErrInvalidTransformer }
 if err := from.acceptMarker(reflect.TypeFor[F]()); err != nil { return nil, err }
 if err := to.acceptMarker(reflect.TypeFor[T]()); err != nil { return nil, err }
 return bound[F, T]{from: from, to: to, stage: stage}, nil
}
type bound[F, T any] struct {
 from, to Language
 stage pipeline.Stage[ast.AST[F], ast.AST[T]]
}
func (b bound[F, T]) FromID() language.ID { return b.from.ID() }
func (b bound[F, T]) ToID() language.ID { return b.to.ID() }
func (b bound[F, T]) FromLanguage() Language { return b.from }
func (b bound[F, T]) ToLanguage() Language { return b.to }
func (b bound[F, T]) Transform(input AST) (AST, error) {
 if err := b.from.Check(input); err != nil { return AST{}, err }
 source, err := Unpack[F](b.from, input)
 if err != nil { return AST{}, err }
 result, err := b.stage(source)
 if err != nil { return AST{}, err }
 return Pack(b.to, result)
}
