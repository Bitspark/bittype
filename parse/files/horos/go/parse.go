// Package horosparse composes per-file syntax parsing with Horos semantics.
package horosparse

import (
	"context"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
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
)

type Provider struct {
	s.Language
	stages horosfiles.Stages
}

func New(definition representation.Language, syntax parse.Provider[parse.Input, cst.CST[s.Language]], printer print.Provider[cst.CST[s.Language], string], checker check.Checker, resolver resolution.Resolver, analyzer analysis.Analyzer, resolutionOptions, analysisOptions semantic.Options) (parse.FilesParser[s.Language], error) {
	stages, err := (horosfiles.Stages{Language: definition, Parser: syntax, Printer: printer, Checker: checker, Resolver: resolver, Analyzer: analyzer, ResolutionOptions: resolutionOptions, AnalysisOptions: analysisOptions}).Capture()
	if err != nil {
		return parse.FilesParser[s.Language]{}, err
	}
	contract := parse.FilesContract[s.Language]{Evidence: stages.Evidence(), Print: printer.Print, Check: func(ctx context.Context, input parse.FilesInput, output source.Publication[s.Language]) error {
		for _, unit := range output.Units {
			if unit.File == input.Entry {
				copy := unit.Syntax.Abstract
				if err := horosfiles.SelectEntry(&copy, input.Configuration); err != nil {
					return err
				}
				if horosfiles.RootPin(copy) != horosfiles.RootPin(unit.Syntax.Abstract) {
					return source.ErrInvalid
				}
			}
		}
		return stages.CheckPublication(ctx, output)
	}}
	return parse.NewFilesParser(definition, contract, Provider{stages: stages})
}
func (p Provider) Parse(ctx context.Context, request parse.FilesRequest[s.Language]) (source.Source[s.Language], error) {
	if ctx == nil || representation.IsNil(p.stages.Parser) {
		return source.Source[s.Language]{}, source.ErrInvalid
	}
	input := request.Input()
	var units []source.Unit[s.Language]
	nodes := input.Limits.Semantic.MaxInputNodes
	for _, file := range input.Files.Entries {
		if err := ctx.Err(); err != nil {
			return source.Source[s.Language]{}, err
		}
		id := ast.TreeID("horos/files/" + string(file.ID))
		sourceID := ast.SourceID(file.ID)
		concrete, err := p.stages.Parser.Parse(ctx, parse.Input{SourceID: sourceID, TreeID: id, Text: file.Text})
		if err != nil {
			return source.Source[s.Language]{}, err
		}
		if len(concrete.Abstract.Nodes) > nodes {
			return source.Source[s.Language]{}, semantic.ErrLimit
		}
		nodes -= len(concrete.Abstract.Nodes)
		if file.ID == input.Entry {
			if err := horosfiles.SelectEntry(&concrete.Abstract, input.Configuration); err != nil {
				return source.Source[s.Language]{}, err
			}
		}
		units = append(units, source.Unit[s.Language]{File: file.ID, Path: file.Path, Source: sourceID, Syntax: concrete})
	}
	publication, err := p.stages.Establish(ctx, units, input.Entry, input.Limits.Semantic)
	if err != nil {
		return source.Source[s.Language]{}, err
	}
	return request.Publish(publication)
}

var _ parse.Provider[parse.FilesRequest[s.Language], source.Source[s.Language]] = Provider{}

func EntryModuleSetting() semantics.Key { return horosfiles.EntryModuleSetting() }
