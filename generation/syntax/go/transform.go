package generation

import (
	"context"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Factory creates fresh invocation stores from the exact admitted request.
// Caller-owned mapping policy remains stable across invocations.
type Factory[F, T any] func(context.Context, transform.Request[F]) (Composition[F, T], error)

func NewTransformer[F, T any](from, to representation.Language, contract transform.Contract, factory Factory[F, T]) (transform.ASTTransformer, error) {
	if factory == nil {
		return nil, transform.ErrInvalidTransformer
	}
	return transform.NewTransformer[F, T](from, to, contract, func(ctx context.Context, r transform.Request[F]) (ast.Set[T], error) {
		c, err := factory(ctx, r)
		if err != nil {
			return ast.Set[T]{}, err
		}
		if !to.Compatible(c.Config.Target) {
			return ast.Set[T]{}, g.ErrScope
		}
		var evidence []semantic.Result
		for i := 0; ; i++ {
			e, ok := r.Evidence(i)
			if !ok {
				break
			}
			evidence = append(evidence, e)
		}
		output, err := Run(ctx, g.Input{Snapshot: r.Snapshot(), Evidence: evidence}, c)
		if err != nil {
			return ast.Set[T]{}, err
		}
		result := ast.Set[T]{Entry: output.Entry}
		for _, envelope := range output.Trees {
			tree, err := representation.Unpack[T](to, envelope)
			if err != nil {
				return ast.Set[T]{}, err
			}
			result.Trees = append(result.Trees, tree)
		}
		return result, nil
	})
}
