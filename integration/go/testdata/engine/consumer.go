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

func Compose(stage func(ast.AST[First]) (ast.AST[Second], error), input ast.AST[First]) (representation.ASTSet, error) {
	first, err := representation.NewLanguage[First](checker[First]{"external-first"}, "v1")
	if err != nil {
		return representation.ASTSet{}, err
	}
	second, err := representation.NewLanguage[Second](checker[Second]{"external-second"}, "v1")
	if err != nil {
		return representation.ASTSet{}, err
	}
	e := engine.New()
	if err := e.RegisterLanguage(first); err != nil {
		return representation.ASTSet{}, err
	}
	firstChecker, err := check.NewChecker[First](first, checker[First]{first.ID()})
	if err != nil {
		return representation.ASTSet{}, err
	}
	if err := e.RegisterChecker(firstChecker); err != nil {
		return representation.ASTSet{}, err
	}
	if err := e.RegisterLanguage(second); err != nil {
		return representation.ASTSet{}, err
	}
	secondChecker, err := check.NewChecker[Second](second, checker[Second]{second.ID()})
	if err != nil {
		return representation.ASTSet{}, err
	}
	if err := e.RegisterChecker(secondChecker); err != nil {
		return representation.ASTSet{}, err
	}
	transformer, err := transform.NewTransformer[First, Second](first, second, transform.Contract{}, func(_ context.Context, r transform.Request[First]) (ast.Set[Second], error) {
		return singleResult(stage(r.Source()))
	})
	if err != nil {
		return representation.ASTSet{}, err
	}
	if err := e.RegisterTransformer(transformer); err != nil {
		return representation.ASTSet{}, err
	}
	packed, err := representation.Pack(first, input)
	if err != nil {
		return representation.ASTSet{}, err
	}
	snapshot, err := e.Capture(context.Background(), semantic.Input{Entry: ast.NodeRef{Tree: input.ID, Node: input.Root}, Trees: []representation.AST{packed}})
	if err != nil {
		return representation.ASTSet{}, err
	}
	return e.TransformTo(context.Background(), transform.Input{Source: snapshot}, second.ID())
}

func singleResult[L any](tree ast.AST[L], err error) (ast.Set[L], error) {
	if err != nil {
		return ast.Set[L]{}, err
	}
	return ast.Set[L]{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: []ast.AST[L]{tree}}, nil
}
