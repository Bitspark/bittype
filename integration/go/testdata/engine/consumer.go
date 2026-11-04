package independent

import (
	"context"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
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

var _ check.Provider[First] = checker[First]{}

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
	firstChecker, err := check.NewChecker[First](first, checker[First]{first.ID()})
	if err != nil {
		return representation.AST{}, err
	}
	if err := e.RegisterChecker(firstChecker); err != nil {
		return representation.AST{}, err
	}
	if err := e.RegisterLanguage(second); err != nil {
		return representation.AST{}, err
	}
	secondChecker, err := check.NewChecker[Second](second, checker[Second]{second.ID()})
	if err != nil {
		return representation.AST{}, err
	}
	if err := e.RegisterChecker(secondChecker); err != nil {
		return representation.AST{}, err
	}
	transformer, err := transform.NewTransformer[First, Second](first, second, transform.Contract{}, func(_ context.Context, r transform.Request[First]) (ast.AST[Second], error) { return stage(r.Source()) })
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
	snapshot, err := e.Capture(context.Background(), semantic.Input{Entry: ast.NodeRef{Tree: input.ID, Node: input.Root}, Trees: []representation.AST{packed}})
	if err != nil {
		return representation.AST{}, err
	}
	return e.TransformTo(context.Background(), transform.Input{Source: snapshot}, second.ID())
}
