package format

import (
	"context"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/source/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"strings"
)

// FilesLimits separately bounds concrete storage and semantic work.
type FilesLimits struct {
	Source   source.Limits
	Semantic semantic.Limits
}

func DefaultFilesLimits() FilesLimits {
	return FilesLimits{source.DefaultLimits(), semantic.DefaultLimits()}
}
func (l FilesLimits) normalized() (FilesLimits, error) {
	if l == (FilesLimits{}) {
		l = DefaultFilesLimits()
	}
	if !l.Source.Valid() || !l.Semantic.Valid() {
		return FilesLimits{}, source.ErrInvalid
	}
	return l, nil
}

type FilesInput[L any] struct {
	Abstract      semantics.Semantics[L]
	Snapshot      semantic.Snapshot
	Evidence      semantic.Result
	Configuration []semantic.Setting
	Limits        FilesLimits
}
type FilesContract[L any] struct {
	Evidence semantic.EvidenceRequirement
	Print    func(context.Context, cst.CST[L]) (string, error)
	Check    func(context.Context, FilesInput[L], source.Publication[L]) error
}
type FilesRequest[L any] struct{ state *filesRequestState[L] }
type filesRequestState[L any] struct {
	input     FilesInput[L]
	publisher source.Publisher[L]
}

func (r FilesRequest[L]) Input() FilesInput[L] {
	if r.state == nil {
		return FilesInput[L]{}
	}
	out := r.state.input
	out.Configuration, _ = semantic.CaptureSettings(out.Configuration, out.Limits.Semantic)
	return out
}
func (r FilesRequest[L]) Publish(p source.Publication[L]) (source.Source[L], error) {
	if r.state == nil {
		return source.Source[L]{}, source.ErrInvalid
	}
	return r.state.publisher.Publish(p)
}

// FilesValue is construction input for the heterogeneous engine. PackFilesInput
// retains the real marker; admission still requires exact snapshot/receipt identity.
type FilesValue struct {
	Abstract      representation.Semantics
	Snapshot      semantic.Snapshot
	Evidence      semantic.Result
	Configuration []semantic.Setting
	Limits        FilesLimits
}

func PackFilesInput[L any](definition representation.Language, input FilesInput[L]) (FilesValue, error) {
	abstract, err := representation.PackSemantics(definition, input.Abstract)
	if err != nil {
		return FilesValue{}, err
	}
	return FilesValue{abstract, input.Snapshot, input.Evidence, input.Configuration, input.Limits}, nil
}

type FilesFormatter[L any] struct {
	definition representation.Language
	provider   Provider[FilesRequest[L], source.Source[L]]
	contract   FilesContract[L]
}
type FilesRegistration interface {
	language.Provider
	Language() representation.Language
	Validate() error
	FormatFiles(context.Context, FilesValue) (source.Envelope, error)
}

func NewFilesFormatter[L any](definition representation.Language, contract FilesContract[L], provider Provider[FilesRequest[L], source.Source[L]]) (FilesFormatter[L], error) {
	if representation.IsNil(provider) || provider.ID() != definition.ID() || strings.TrimSpace(provider.Name()) == "" || contract.Print == nil || contract.Check == nil || !contract.Evidence.Semantics {
		return FilesFormatter[L]{}, ErrInvalidFormatter
	}
	if err := representation.CheckMarker[L](definition); err != nil {
		return FilesFormatter[L]{}, err
	}
	if err := semantic.ValidateEvidenceRequirements([]semantic.EvidenceRequirement{contract.Evidence}); err != nil {
		return FilesFormatter[L]{}, err
	}
	contract.Evidence = contract.Evidence.Clone()
	return FilesFormatter[L]{definition, provider, contract}, nil
}
func (p FilesFormatter[L]) ID() language.ID                   { return p.definition.ID() }
func (p FilesFormatter[L]) Name() string                      { return p.definition.Name() }
func (p FilesFormatter[L]) Language() representation.Language { return p.definition }
func (p FilesFormatter[L]) Validate() error {
	if representation.IsNil(p.provider) {
		return ErrInvalidFormatter
	}
	return p.definition.Validate()
}
func (p FilesFormatter[L]) Format(ctx context.Context, input FilesInput[L]) (source.Source[L], error) {
	if err := p.Validate(); err != nil {
		return source.Source[L]{}, err
	}
	if ctx == nil {
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
	if _, err := semantic.AdmitEvidence(input.Snapshot, []semantic.EvidenceRequirement{p.contract.Evidence}, []semantic.Result{input.Evidence}); err != nil {
		return source.Source[L]{}, err
	}
	supplied, err := representation.PackSemantics(p.definition, input.Abstract)
	if err != nil {
		return source.Source[L]{}, err
	}
	established, ok := input.Evidence.Semantics()
	if !ok || !supplied.Same(established) {
		return source.Source[L]{}, semantic.ErrInvalidRequest
	}
	if err := supplied.CheckLimits(limits.Semantic.Semantics); err != nil {
		return source.Source[L]{}, err
	}
	if err := input.Snapshot.CheckLimits(limits.Semantic); err != nil {
		return source.Source[L]{}, err
	}
	settings, err := semantic.CaptureSettings(input.Configuration, limits.Semantic)
	if err != nil {
		return source.Source[L]{}, err
	}
	input.Configuration = settings
	policy := source.Policy[L]{Evidence: p.contract.Evidence, SemanticLimits: limits.Semantic, Limits: limits.Source, Print: p.contract.Print, Check: func(ctx context.Context, output source.Publication[L]) error {
		if input.Snapshot.Same(output.Snapshot) {
			return source.ErrStale
		}
		guardInput := input
		guardInput.Configuration, _ = semantic.CaptureSettings(input.Configuration, limits.Semantic)
		return p.contract.Check(ctx, guardInput, output)
	}}
	return publishFiles(ctx, p.definition, policy, func(publisher source.Publisher[L]) (source.Source[L], error) {
		return p.provider.Format(ctx, FilesRequest[L]{&filesRequestState[L]{input, publisher}})
	})
}
func publishFiles[L any](ctx context.Context, language representation.Language, policy source.Policy[L], run func(source.Publisher[L]) (source.Source[L], error)) (source.Source[L], error) {
	publisher, err := source.NewPublisher(ctx, language, policy)
	if err != nil {
		return source.Source[L]{}, err
	}
	defer publisher.Close()
	out, err := run(publisher)
	if err != nil {
		return source.Source[L]{}, err
	}
	if err = publisher.Accept(out); err != nil {
		return source.Source[L]{}, err
	}
	return out, nil
}
func (p FilesFormatter[L]) FormatFiles(ctx context.Context, input FilesValue) (source.Envelope, error) {
	abstract, err := representation.UnpackSemantics[L](p.definition, input.Abstract)
	if err != nil {
		return source.Envelope{}, err
	}
	output, err := p.Format(ctx, FilesInput[L]{abstract, input.Snapshot, input.Evidence, input.Configuration, input.Limits})
	if err != nil {
		return source.Envelope{}, err
	}
	return source.Pack(p.definition, output)
}
