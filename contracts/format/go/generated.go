package format

import (
	"context"
	"fmt"
	"github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/source/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"slices"
)

// GeneratedPlacement is a consumer's placement on its way to the formatter of
// the language that issued it. A placement type belongs to its language: what a
// TypeScript file, a Go package or a Rust module tree states differs. The
// engine and this contract carry the value and do not read it. Only the issuer
// reads it back, as only the issuer of an external symbol reads that symbol, so
// a placement of another language is a refusal of that language's preparation
// step and never a failed type assertion.
type GeneratedPlacement struct {
	owner *placementOwner
	value any
}
type placementOwner struct{ marker byte }

// Valid reports whether an issuer issued the value. The zero value is no
// placement.
func (p GeneratedPlacement) Valid() bool { return p.owner != nil }

// GeneratedPlacementIssuer is a language's authority to issue placements of its
// own type P and to read them back. One issuer is one placement type of one
// language.
type GeneratedPlacementIssuer[P any] struct{ owner *placementOwner }

func NewGeneratedPlacementIssuer[P any]() GeneratedPlacementIssuer[P] {
	return GeneratedPlacementIssuer[P]{&placementOwner{1}}
}

// Issue wraps a placement. The issuer does not copy it: a language whose type
// holds maps or slices issues its own copy.
func (i GeneratedPlacementIssuer[P]) Issue(value P) GeneratedPlacement {
	if i.owner == nil {
		return GeneratedPlacement{}
	}
	return GeneratedPlacement{i.owner, value}
}

// Read returns the placement that this issuer issued, and false for the zero
// value and for the value of every other issuer.
func (i GeneratedPlacementIssuer[P]) Read(p GeneratedPlacement) (P, bool) {
	var zero P
	if i.owner == nil || p.owner != i.owner {
		return zero, false
	}
	value, ok := p.value.(P)
	return value, ok
}

// GeneratedFilesInput carries the sealed output of one generation run, the
// consumer's placement of its declarations and the limits. The placement is a
// typed value built from the handles that the output issued. No placement
// bytes cross this contract. The target supplies its dependency environment
// when its formatter is constructed.
type GeneratedFilesInput struct {
	Output    generation.Output
	Placement GeneratedPlacement
	Limits    FilesLimits
}

// GeneratedPlacementView is a placement in the one form that the common
// admission reads: each file identity with its relative path, and each
// declaration handle with the file that holds its definition. A language
// derives the view from its own placement type. The view states no scope, no
// order and no entry.
type GeneratedPlacementView struct {
	Files       map[files.ID]files.Path
	Definitions map[generation.DeclarationHandle]files.ID
}

// PlacedDeclaration is one declaration of the sealed output with the file that
// the admitted placement gives it.
type PlacedDeclaration struct {
	Declaration generation.OwnedDeclaration
	File        files.ID
}

// AdmittedPlacement is a view that the common admission accepted for one
// request. Only GeneratedFilesPreparation.Admit produces it. It holds:
//
//   - every assigned handle is a handle of the sealed output of the request;
//   - every declaration of the output's inventory has exactly one file;
//   - every assignment names a file of the files table;
//   - the files table satisfies files.Files.Capture under the request's limits,
//     and has no more files than the file limit and the link limit.
//
// It states one definition assignment per owned declaration. It does not state
// that a declaration is one entity of a native compiler, it interprets no
// scope, and it admits a file without a definition and two declarations of one
// name in one file. A language refuses those for itself.
type AdmittedPlacement struct{ state *admittedPlacement }
type admittedPlacement struct {
	files       files.Files
	definitions []PlacedDeclaration
}

func (a AdmittedPlacement) Valid() bool { return a.state != nil }

// Files returns the admitted files, each with its identity and path and
// without text, in the order of their identities. This is the order of files
// that a Files formatter realizes: no placement states one.
func (a AdmittedPlacement) Files() []files.File {
	if a.state == nil {
		return nil
	}
	return a.state.files.Clone().Entries
}

// Definitions returns every declaration of the sealed output with its file, in
// the order of the output's inventory, which is the order of the sites. This
// is the order of declarations that a Files formatter realizes.
func (a AdmittedPlacement) Definitions() []PlacedDeclaration {
	if a.state == nil {
		return nil
	}
	out := make([]PlacedDeclaration, len(a.state.definitions))
	for i, d := range a.state.definitions {
		d.Declaration.Provenance = slices.Clone(d.Declaration.Provenance)
		d.Declaration.Values = slices.Clone(d.Declaration.Values)
		out[i] = d
	}
	return out
}

// admitPlacement is the common admission. Its refusals are invalid input or an
// exceeded limit, and none is unsupported coverage: validity is decided here,
// before a language refuses what it does not cover.
func admitPlacement(output generation.Output, view GeneratedPlacementView, limits source.Limits) (*admittedPlacement, error) {
	// Whether a handle belongs to an output is the output's answer. Its value
	// stays joined to the Files category.
	for handle := range view.Definitions {
		if _, err := output.Resolve(handle); err != nil {
			return nil, fmt.Errorf("%w: a placed declaration is no declaration of this output: %w", source.ErrInvalid, err)
		}
	}
	inventory := output.Inventory()
	handles := make([]generation.DeclarationHandle, len(inventory))
	for i, record := range inventory {
		handles[i] = record.Handle
	}
	// Cover answers with a value of the files package only: files.ErrInvalid or
	// files.ErrLimit.
	admitted, err := files.Cover(view.Files, view.Definitions, handles, limits.Files)
	if err != nil {
		return nil, source.WithCategory(err)
	}
	definitions := make([]PlacedDeclaration, len(inventory))
	for i, record := range inventory {
		definitions[i] = PlacedDeclaration{Declaration: record, File: view.Definitions[record.Handle]}
	}
	return &admittedPlacement{admitted, definitions}, nil
}

// published refuses a publication whose file identities and paths differ from
// the admitted files table.
func published[L any](a *admittedPlacement, units []source.Unit[L]) error {
	if len(units) != len(a.files.Entries) {
		return fmt.Errorf("%w: %d files published for %d admitted", source.ErrInvalid, len(units), len(a.files.Entries))
	}
	paths := make(map[files.ID]files.Path, len(a.files.Entries))
	for _, file := range a.files.Entries {
		paths[file.ID] = file.Path
	}
	for _, unit := range units {
		path, admitted := paths[unit.File]
		if !admitted || path != unit.Path {
			return fmt.Errorf("%w: published file %q at %q is not a file of the admitted placement", source.ErrInvalid, unit.File, unit.Path)
		}
		delete(paths, unit.File)
	}
	return nil
}

// GeneratedFilesPreparation is one request as the language's preparation step
// receives it.
type GeneratedFilesPreparation struct{ state *generatedPreparation }
type generatedPreparation struct {
	input    GeneratedFilesInput
	admitted *admittedPlacement
	closed   bool
}

// Input returns the request with its limits normalized.
func (p GeneratedFilesPreparation) Input() GeneratedFilesInput {
	if p.state == nil {
		return GeneratedFilesInput{}
	}
	return p.state.input
}

// Admit runs the common admission on the view that the language derived from
// its placement. The step calls it once, before it refuses anything as outside
// its coverage, so that validity is decided before coverage. Its refusals have
// a Files category: invalid input or an exceeded limit.
func (p GeneratedFilesPreparation) Admit(view GeneratedPlacementView) (AdmittedPlacement, error) {
	if p.state == nil || p.state.closed || p.state.admitted != nil {
		return AdmittedPlacement{}, ErrInvalidFormatter
	}
	admitted, err := admitPlacement(p.state.input.Output, view, p.state.input.Limits.Source)
	if err != nil {
		return AdmittedPlacement{}, err
	}
	p.state.admitted = admitted
	return AdmittedPlacement{admitted}, nil
}

// GeneratedFilesPrepared is what a language states for one request: the view
// that the common admission accepted, the evidence that the analysis of the
// fresh output has to establish, and the guard of the publication.
//
// Evidence may depend on the placement: the compilation root of a Rust crate
// and the translation-unit contexts of C++ headers are part of what a receipt
// has to be bound to, and both follow from the placement. This contract matches
// the requirement with the receipt and learns no manifest type.
//
// Check compares the published value with the admitted consumer input: where a
// definition was published against where the placement puts it, and what the
// language realized against what the output intends. Running the layout
// function a second time and comparing its result with the publication shows
// that the function is deterministic. It is an integrity check and no more: a
// layout function that ignores a part of the placement agrees with itself.
type GeneratedFilesPrepared[L any] struct {
	Placement AdmittedPlacement
	Evidence  semantic.EvidenceRequirement
	Check     func(context.Context, source.Publication[L]) error
}

// GeneratedFilesContract is what a language supplies beside its provider: the
// Syntax printer and one preparation step. The step reads the language's own
// placement from the request, derives the view, has it admitted, and returns
// what GeneratedFilesPrepared holds. It refuses a placement that is malformed
// under the language's rules as invalid input, and one that the language does
// not cover yet as unsupported coverage, in that order.
type GeneratedFilesContract[L any] struct {
	Print   func(context.Context, cst.CST[L]) (string, error)
	Prepare func(context.Context, GeneratedFilesPreparation) (GeneratedFilesPrepared[L], error)
}

// GeneratedFilesRequest is one admitted request as the provider receives it.
type GeneratedFilesRequest[L any] struct {
	input     GeneratedFilesInput
	placement AdmittedPlacement
	publisher source.Publisher[L]
}

func (r GeneratedFilesRequest[L]) Input() GeneratedFilesInput   { return r.input }
func (r GeneratedFilesRequest[L]) Placement() AdmittedPlacement { return r.placement }
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

// NewGeneratedFilesFormatter pairs a language's preparation step with a
// provider. The two are separate arguments, so a test can pair a language's own
// step with a provider that tampers with the publication. The evidence
// requirement is no part of the construction: the step states it per request.
func NewGeneratedFilesFormatter[L any](language representation.Language, contract GeneratedFilesContract[L], provider Provider[GeneratedFilesRequest[L], source.Source[L]]) (GeneratedFilesFormatter[L], error) {
	if err := representation.CheckMarker[L](language); err != nil {
		return GeneratedFilesFormatter[L]{}, err
	}
	if representation.IsNil(provider) || provider.ID() != language.ID() || contract.Print == nil || contract.Prepare == nil {
		return GeneratedFilesFormatter[L]{}, ErrInvalidFormatter
	}
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

// Format admits the request, asks the language which evidence it requires,
// calls the provider and guards its publication. The provider is not called
// when the step or the common admission refuses. A step that returns without
// an admitted placement, without a guard or with a requirement that asks for no
// semantic graph fails the request with ErrInvalidFormatter: that is a defect
// of the language's contribution and no category of the input.
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
	// The snapshot is valid and the limits are normalized, so the check answers
	// with semantic.ErrLimit, or with a value of the files package for a
	// snapshot that carries files.
	if err = snapshot.CheckLimits(limits.Semantic); err != nil {
		return source.Source[L]{}, source.WithCategory(err)
	}
	preparation := &generatedPreparation{input: input}
	prepared, err := p.contract.Prepare(ctx, GeneratedFilesPreparation{preparation})
	preparation.closed = true
	if err != nil {
		return source.Source[L]{}, err
	}
	if preparation.admitted == nil || prepared.Placement.state != preparation.admitted || prepared.Check == nil || !prepared.Evidence.Semantics {
		return source.Source[L]{}, ErrInvalidFormatter
	}
	if err = semantic.ValidateEvidenceRequirements([]semantic.EvidenceRequirement{prepared.Evidence}); err != nil {
		return source.Source[L]{}, err
	}
	admitted := prepared.Placement
	policy := source.Policy[L]{Evidence: prepared.Evidence.Clone(), SemanticLimits: limits.Semantic, Limits: limits.Source, Print: p.contract.Print, Check: func(ctx context.Context, output source.Publication[L]) error {
		if snapshot.Same(output.Snapshot) {
			return source.ErrStale
		}
		if err := published(admitted.state, output.Units); err != nil {
			return err
		}
		return prepared.Check(ctx, output)
	}}
	return publishFiles(ctx, p.definition, policy, func(publisher source.Publisher[L]) (source.Source[L], error) {
		return p.provider.Format(ctx, GeneratedFilesRequest[L]{input, admitted, publisher})
	})
}
func (p GeneratedFilesFormatter[L]) FormatGeneratedFiles(ctx context.Context, input GeneratedFilesInput) (source.Envelope, error) {
	output, err := p.Format(ctx, input)
	if err != nil {
		return source.Envelope{}, err
	}
	return source.Pack(p.definition, output)
}
