package format

import (
	"context"
	"github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/source/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// GeneratedFilesInput carries intended bindings from a sealed run and the
// consumer's placement settings. The target supplies its dependency environment.
type GeneratedFilesInput struct {
	Output        generation.Output
	Configuration []semantic.Setting
	Limits        FilesLimits
}
type GeneratedFilesContract[L any] struct {
	Evidence semantic.EvidenceRequirement
	Print    func(context.Context, cst.CST[L]) (string, error)
	Check    func(context.Context, GeneratedFilesInput, source.Publication[L]) error
}
type GeneratedFilesRequest[L any] struct {
	input     GeneratedFilesInput
	publisher source.Publisher[L]
}

func (r GeneratedFilesRequest[L]) Input() GeneratedFilesInput {
	out := r.input
	out.Configuration, _ = semantic.CaptureSettings(out.Configuration, out.Limits.Semantic)
	return out
}
func (r GeneratedFilesRequest[L]) Publish(p source.Publication[L]) (source.Source[L], error) {
	return r.publisher.Publish(p)
}

type GeneratedFilesFormatter[L any] struct {
	definition representation.Language
	provider   Provider[GeneratedFilesRequest[L], source.Source[L]]
	contract   GeneratedFilesContract[L]
}
type GeneratedFilesRegistration interface {
	language.Provider
	Language() representation.Language
	Validate() error
	FormatGeneratedFiles(context.Context, GeneratedFilesInput) (source.Envelope, error)
}

func NewGeneratedFilesFormatter[L any](language representation.Language, contract GeneratedFilesContract[L], provider Provider[GeneratedFilesRequest[L], source.Source[L]]) (GeneratedFilesFormatter[L], error) {
	if err := representation.CheckMarker[L](language); err != nil {
		return GeneratedFilesFormatter[L]{}, err
	}
	if representation.IsNil(provider) || provider.ID() != language.ID() || contract.Print == nil || contract.Check == nil || !contract.Evidence.Semantics {
		return GeneratedFilesFormatter[L]{}, ErrInvalidFormatter
	}
	if err := semantic.ValidateEvidenceRequirements([]semantic.EvidenceRequirement{contract.Evidence}); err != nil {
		return GeneratedFilesFormatter[L]{}, err
	}
	contract.Evidence = contract.Evidence.Clone()
	return GeneratedFilesFormatter[L]{language, provider, contract}, nil
}
func (p GeneratedFilesFormatter[L]) ID() language.ID                   { return p.definition.ID() }
func (p GeneratedFilesFormatter[L]) Name() string                      { return p.definition.Name() }
func (p GeneratedFilesFormatter[L]) Language() representation.Language { return p.definition }
func (p GeneratedFilesFormatter[L]) Validate() error {
	if representation.IsNil(p.provider) {
		return ErrInvalidFormatter
	}
	return p.definition.Validate()
}
func (p GeneratedFilesFormatter[L]) Format(ctx context.Context, input GeneratedFilesInput) (source.Source[L], error) {
	if err := p.Validate(); err != nil {
		return source.Source[L]{}, err
	}
	snapshot := input.Output.Snapshot()
	if ctx == nil || !snapshot.Valid() || !p.definition.Compatible(snapshot.Language()) {
		return source.Source[L]{}, source.ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return source.Source[L]{}, err
	}
	limits, err := input.Limits.normalized()
	if err != nil {
		return source.Source[L]{}, err
	}
	input.Limits = limits
	if err = snapshot.CheckLimits(limits.Semantic); err != nil {
		return source.Source[L]{}, err
	}
	input.Configuration, err = semantic.CaptureSettings(input.Configuration, limits.Semantic)
	if err != nil {
		return source.Source[L]{}, err
	}
	policy := source.Policy[L]{Evidence: p.contract.Evidence, SemanticLimits: limits.Semantic, Limits: limits.Source, Print: p.contract.Print, Check: func(ctx context.Context, output source.Publication[L]) error {
		if snapshot.Same(output.Snapshot) {
			return source.ErrStale
		}
		guard := input
		guard.Configuration, _ = semantic.CaptureSettings(input.Configuration, limits.Semantic)
		return p.contract.Check(ctx, guard, output)
	}}
	return publishFiles(ctx, p.definition, policy, func(publisher source.Publisher[L]) (source.Source[L], error) {
		return p.provider.Format(ctx, GeneratedFilesRequest[L]{input, publisher})
	})
}
func (p GeneratedFilesFormatter[L]) FormatGeneratedFiles(ctx context.Context, input GeneratedFilesInput) (source.Envelope, error) {
	output, err := p.Format(ctx, input)
	if err != nil {
		return source.Envelope{}, err
	}
	return source.Pack(p.definition, output)
}
