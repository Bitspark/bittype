package transform

import (
	"context"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/internal/tree"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"slices"
)

// Request supplies an owned typed source and admitted immutable evidence.
type Request[L any] struct {
	source   ast.AST[L]
	entry    ast.NodeRef
	evidence []semantic.Result
	inputs   []representation.AST
}

func (r Request[L]) Source() ast.AST[L] { return tree.Clone(r.source) }
func (r Request[L]) Entry() ast.NodeRef { return r.entry }

// Inputs retains the complete captured source collection for multi-tree mappings.
func (r Request[L]) Inputs() []representation.AST { return slices.Clone(r.inputs) }
func (r Request[L]) Evidence(index int) (semantic.Result, bool) {
	if index < 0 || index >= len(r.evidence) {
		return semantic.Result{}, false
	}
	return r.evidence[index], true
}

// NewTransformer captures requirements and shares semantic receipt admission.
// Endpoint syntax admission remains the engine's checker responsibility.
func NewTransformer[F, T any](from, to representation.Language, contract Contract, stage func(context.Context, Request[F]) (ast.AST[T], error)) (ASTTransformer, error) {
	if stage == nil {
		return nil, ErrInvalidTransformer
	}
	if err := representation.CheckMarker[F](from); err != nil {
		return nil, err
	}
	if err := representation.CheckMarker[T](to); err != nil {
		return nil, err
	}
	if err := semantic.ValidateEvidenceRequirements(contract.Evidence); err != nil {
		return nil, err
	}
	requirements := make([]semantic.EvidenceRequirement, len(contract.Evidence))
	for i, requirement := range contract.Evidence {
		requirements[i] = requirement.Clone()
	}
	return bound[F, T]{from: from, to: to, requirements: requirements, stage: stage}, nil
}

type bound[F, T any] struct {
	from, to     representation.Language
	requirements []semantic.EvidenceRequirement
	stage        func(context.Context, Request[F]) (ast.AST[T], error)
}

func (b bound[F, T]) FromID() language.ID                   { return b.from.ID() }
func (b bound[F, T]) ToID() language.ID                     { return b.to.ID() }
func (b bound[F, T]) FromLanguage() representation.Language { return b.from }
func (b bound[F, T]) ToLanguage() representation.Language   { return b.to }
func (b bound[F, T]) Transform(ctx context.Context, input Input) (representation.AST, error) {
	if ctx == nil {
		return representation.AST{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return representation.AST{}, err
	}
	evidence, err := semantic.AdmitEvidence(input.Source, b.requirements, input.Evidence)
	if err != nil {
		return representation.AST{}, err
	}
	envelope, err := input.Source.EntryTree()
	if err != nil {
		return representation.AST{}, err
	}
	source, err := representation.Unpack[F](b.from, envelope)
	if err != nil {
		return representation.AST{}, err
	}
	result, err := b.stage(ctx, Request[F]{source: source, entry: input.Source.Entry(), evidence: evidence, inputs: input.Source.Inputs()})
	if err != nil {
		return representation.AST{}, err
	}
	if err := ctx.Err(); err != nil {
		return representation.AST{}, err
	}
	return representation.Pack(b.to, result)
}
