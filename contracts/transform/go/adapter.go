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
	snapshot semantic.Snapshot
	source   ast.AST[L]
	entry    ast.NodeRef
	evidence []semantic.Result
	inputs   []representation.AST
}

func (r Request[L]) Source() ast.AST[L] { return tree.Clone(r.source) }
func (r Request[L]) Entry() ast.NodeRef { return r.entry }

// Snapshot returns the exact captured input, including when no evidence is required.
// A zero request returns an invalid snapshot.
func (r Request[L]) Snapshot() semantic.Snapshot { return r.snapshot }

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
func NewTransformer[F, T any](from, to representation.Language, contract Contract, stage func(context.Context, Request[F]) (ast.Set[T], error)) (ASTTransformer, error) {
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
	stage        func(context.Context, Request[F]) (ast.Set[T], error)
}

func (b bound[F, T]) FromID() language.ID                   { return b.from.ID() }
func (b bound[F, T]) ToID() language.ID                     { return b.to.ID() }
func (b bound[F, T]) FromLanguage() representation.Language { return b.from }
func (b bound[F, T]) ToLanguage() representation.Language   { return b.to }
func (b bound[F, T]) Transform(ctx context.Context, input Input) (representation.ASTSet, error) {
	if ctx == nil {
		return representation.ASTSet{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return representation.ASTSet{}, err
	}
	request, err := AdmitRequest[F](b.from, Contract{Evidence: b.requirements}, input)
	if err != nil {
		return representation.ASTSet{}, err
	}
	result, err := b.stage(ctx, request)
	if err != nil {
		return representation.ASTSet{}, err
	}
	if err := ctx.Err(); err != nil {
		return representation.ASTSet{}, err
	}
	return representation.PackSet(b.to, result)
}

// AdmitRequest binds the typed source view and the requested evidence to the
// exact input snapshot for a transformation or generation invocation.
func AdmitRequest[L any](language representation.Language, contract Contract, input Input) (Request[L], error) {
	evidence, err := semantic.AdmitEvidence(input.Source, contract.Evidence, input.Evidence)
	if err != nil {
		return Request[L]{}, err
	}
	envelope, err := input.Source.EntryTree()
	if err != nil {
		return Request[L]{}, err
	}
	source, err := representation.Unpack[L](language, envelope)
	if err != nil {
		return Request[L]{}, err
	}
	return Request[L]{snapshot: input.Source, source: source, entry: input.Source.Entry(), evidence: evidence, inputs: input.Source.Inputs()}, nil
}
