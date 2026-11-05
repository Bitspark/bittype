// Package format organizes admitted Horos semantics into concrete files.
package format

import (
	"context"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/internal/horosfiles"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/source/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"slices"
)

type Provider struct {
	s.Language
	stages horosfiles.Stages
}

func New(definition representation.Language, parser parse.Provider[parse.Input, cst.CST[s.Language]], syntax format.Provider[ast.View[s.Language], cst.CST[s.Language]], printer print.Provider[cst.CST[s.Language], string], checker check.Checker, resolver resolution.Resolver, analyzer analysis.Analyzer, resolutionOptions, analysisOptions semantic.Options) (format.FilesFormatter[s.Language], error) {
	if representation.IsNil(syntax) {
		return format.FilesFormatter[s.Language]{}, format.ErrInvalidFormatter
	}
	stages, err := (horosfiles.Stages{Language: definition, Parser: parser, Formatter: syntax, Printer: printer, Checker: checker, Resolver: resolver, Analyzer: analyzer, ResolutionOptions: resolutionOptions, AnalysisOptions: analysisOptions}).Capture()
	if err != nil {
		return format.FilesFormatter[s.Language]{}, err
	}
	contract := format.FilesContract[s.Language]{Evidence: stages.Evidence(), Print: printer.Print, Check: func(ctx context.Context, input format.FilesInput[s.Language], output source.Publication[s.Language]) error {
		before, pins, err := horosfiles.ProgramIdentity(input.Snapshot)
		if err != nil {
			return err
		}
		after, newPins, err := horosfiles.ProgramIdentity(output.Snapshot)
		if err != nil {
			return err
		}
		if before != after || !slices.Equal(pins, newPins) {
			return source.ErrInvalid
		}
		return stages.CheckPublication(ctx, output)
	}}
	return format.NewFilesFormatter(definition, contract, Provider{stages: stages})
}
func (p Provider) Format(ctx context.Context, request format.FilesRequest[s.Language]) (source.Source[s.Language], error) {
	if ctx == nil || representation.IsNil(p.stages.Formatter) {
		return source.Source[s.Language]{}, source.ErrInvalid
	}
	input := request.Input()
	placements, trees, entry, err := horosfiles.Materialize(input.Snapshot, input.Configuration, input.Limits.Source)
	if err != nil {
		return source.Source[s.Language]{}, err
	}
	var units []source.Unit[s.Language]
	for i, tree := range trees {
		if err := ctx.Err(); err != nil {
			return source.Source[s.Language]{}, err
		}
		concrete, err := p.stages.Formatter.Format(ctx, tree)
		if err != nil {
			return source.Source[s.Language]{}, err
		}
		units = append(units, source.Unit[s.Language]{File: placements[i].ID, Path: placements[i].Path, Source: ast.SourceID(placements[i].ID), Syntax: concrete})
	}
	publication, err := p.stages.Establish(ctx, units, entry, input.Limits.Semantic)
	if err != nil {
		return source.Source[s.Language]{}, err
	}
	return request.Publish(publication)
}

var _ format.Provider[format.FilesRequest[s.Language], source.Source[s.Language]] = Provider{}

func OrganizationSetting() semantics.Key { return horosfiles.OrganizationSetting() }
