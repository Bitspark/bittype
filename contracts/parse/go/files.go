package parse

import (
	"context"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
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

type FilesInput struct {
	Files         files.Files
	Entry         files.ID
	Configuration []semantic.Setting
	Limits        FilesLimits
}
type FilesContract[L any] struct {
	Evidence semantic.EvidenceRequirement
	Print    func(context.Context, cst.CST[L]) (string, error)
	Check    func(context.Context, FilesInput, source.Publication[L]) error
}
type FilesRequest[L any] struct{ state *filesRequestState[L] }
type filesRequestState[L any] struct {
	input     FilesInput
	publisher source.Publisher[L]
}

func (r FilesRequest[L]) Input() FilesInput {
	if r.state == nil {
		return FilesInput{}
	}
	out := r.state.input
	out.Files = out.Files.Clone()
	out.Configuration, _ = semantic.CaptureSettings(out.Configuration, out.Limits.Semantic)
	return out
}
func (r FilesRequest[L]) Publish(p source.Publication[L]) (source.Source[L], error) {
	if r.state == nil {
		return source.Source[L]{}, source.ErrInvalid
	}
	return r.state.publisher.Publish(p)
}

type FilesParser[L any] struct {
	definition representation.Language
	provider   Provider[FilesRequest[L], source.Source[L]]
	contract   FilesContract[L]
}

// FilesRegistration is the heterogeneous engine bridge. Typed callers use Parse.
type FilesRegistration interface {
	language.Provider
	Language() representation.Language
	Validate() error
	ParseFiles(context.Context, FilesInput) (source.Envelope, error)
}

func NewFilesParser[L any](definition representation.Language, contract FilesContract[L], provider Provider[FilesRequest[L], source.Source[L]]) (FilesParser[L], error) {
	if representation.IsNil(provider) || provider.ID() != definition.ID() || strings.TrimSpace(provider.Name()) == "" || contract.Print == nil || contract.Check == nil || !contract.Evidence.Semantics {
		return FilesParser[L]{}, ErrInvalidParser
	}
	if err := representation.CheckMarker[L](definition); err != nil {
		return FilesParser[L]{}, err
	}
	if err := semantic.ValidateEvidenceRequirements([]semantic.EvidenceRequirement{contract.Evidence}); err != nil {
		return FilesParser[L]{}, err
	}
	contract.Evidence = contract.Evidence.Clone()
	return FilesParser[L]{definition, provider, contract}, nil
}
func (p FilesParser[L]) ID() language.ID                   { return p.definition.ID() }
func (p FilesParser[L]) Name() string                      { return p.definition.Name() }
func (p FilesParser[L]) Language() representation.Language { return p.definition }
func (p FilesParser[L]) Validate() error {
	if representation.IsNil(p.provider) {
		return ErrInvalidParser
	}
	return p.definition.Validate()
}
func (p FilesParser[L]) Parse(ctx context.Context, input FilesInput) (source.Source[L], error) {
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
	owned, err := input.Files.Capture(limits.Source.Files)
	if err != nil {
		return source.Source[L]{}, err
	}
	input.Files = owned
	found := false
	for _, file := range owned.Entries {
		found = found || file.ID == input.Entry
	}
	if !found {
		return source.Source[L]{}, source.ErrInvalid
	}
	settings, err := semantic.CaptureSettings(input.Configuration, limits.Semantic)
	if err != nil {
		return source.Source[L]{}, err
	}
	input.Configuration = settings
	publisher, err := source.NewPublisher(ctx, p.definition, source.Policy[L]{Evidence: p.contract.Evidence, SemanticLimits: limits.Semantic, Limits: limits.Source, Print: p.contract.Print, Check: func(ctx context.Context, output source.Publication[L]) error {

		captured := output.Snapshot.Files()
		if captured == nil || !source.EqualFiles(input.Files, captured.Files) {
			return source.ErrInvalid
		}
		selected := false
		for _, link := range captured.Links {
			selected = selected || link.File == input.Entry && link.Tree == output.Snapshot.Entry().Tree
		}
		if !selected {
			return source.ErrInvalid
		}
		guardInput := input
		guardInput.Configuration, _ = semantic.CaptureSettings(input.Configuration, limits.Semantic)
		guardInput.Files = input.Files.Clone()
		return p.contract.Check(ctx, guardInput, output)
	}})
	if err != nil {
		return source.Source[L]{}, err
	}
	defer publisher.Close()
	output, err := p.provider.Parse(ctx, FilesRequest[L]{&filesRequestState[L]{input, publisher}})
	if err != nil {
		return source.Source[L]{}, err
	}
	if err := publisher.Accept(output); err != nil {
		return source.Source[L]{}, err
	}
	return output, nil
}
func (p FilesParser[L]) ParseFiles(ctx context.Context, input FilesInput) (source.Envelope, error) {
	output, err := p.Parse(ctx, input)
	if err != nil {
		return source.Envelope{}, err
	}
	return source.Pack(p.definition, output)
}
