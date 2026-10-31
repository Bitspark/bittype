package independent

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/engine/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type First struct{}
type Second struct{}
type checker[L any] struct{ id language.ID }

func (c checker[L]) ID() language.ID         { return c.id }
func (c checker[L]) Name() string            { return string(c.id) }
func (c checker[L]) Check(ast.View[L]) error { return nil }

func Compose(stage func(ast.AST[First]) (ast.AST[Second], error), input ast.AST[First]) (representation.AST, error) {
	first, err := representation.NewLanguage[First](checker[First]{"external-first"}, "v1")
	if err != nil {
		return representation.AST{}, err
	}
	second, err := representation.NewLanguage[Second](checker[Second]{"external-second"}, "v1")
	if err != nil {
		return representation.AST{}, err
	}
	e := engine.New()
	if err := e.RegisterLanguage(first); err != nil {
		return representation.AST{}, err
	}
	if err := e.RegisterLanguage(second); err != nil {
		return representation.AST{}, err
	}
	transformer, err := transform.NewTransformer[First, Second](first, second, stage)
	if err != nil {
		return representation.AST{}, err
	}
	if err := e.RegisterTransformer(transformer); err != nil {
		return representation.AST{}, err
	}
	packed, err := representation.Pack(first, input)
	if err != nil {
		return representation.AST{}, err
	}
	return e.TransformTo(packed, second.ID())
}
