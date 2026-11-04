package provider

import (
	analysisstage "github.com/Bitspark/bittype/analyze/horos/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	filesformat "github.com/Bitspark/bittype/format/files/horos/go"
	filesparse "github.com/Bitspark/bittype/parse/files/horos/go"
	filesprint "github.com/Bitspark/bittype/print/files/go"
	resolutionstage "github.com/Bitspark/bittype/resolve/horos/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func fileStages() (representation.Language, check.Checker, resolution.Resolver, analysis.Analyzer, error) {
	p := Provider{}
	definition, err := representation.NewLanguage[s.Language](p, s.Profile)
	if err != nil {
		return definition, check.Checker{}, resolution.Resolver{}, analysis.Analyzer{}, err
	}
	checker, err := check.NewChecker[s.Language](definition, p)
	if err != nil {
		return definition, checker, resolution.Resolver{}, analysis.Analyzer{}, err
	}
	resolver, err := resolution.NewResolver[s.Language](definition, resolutionstage.Contract(), p)
	if err != nil {
		return definition, checker, resolver, analysis.Analyzer{}, err
	}
	analyzer, err := analysis.NewAnalyzer[s.Language](definition, analysisstage.Contract(), p)
	return definition, checker, resolver, analyzer, err
}

// Each constructor returns exactly one capability; nothing is registered implicitly.
func NewFilesParser() (parse.FilesParser[s.Language], error) {
	definition, checker, resolver, analyzer, err := fileStages()
	if err != nil {
		return parse.FilesParser[s.Language]{}, err
	}
	return filesparse.New(definition, Provider{}, Provider{}, checker, resolver, analyzer, semantic.Options{}, semantic.Options{})
}
func NewFilesFormatter() (format.FilesFormatter[s.Language], error) {
	definition, checker, resolver, analyzer, err := fileStages()
	if err != nil {
		return format.FilesFormatter[s.Language]{}, err
	}
	return filesformat.New(definition, Provider{}, Provider{}, Provider{}, checker, resolver, analyzer, semantic.Options{}, semantic.Options{})
}
func NewFilesPrinter() (print.FilesPrinter[s.Language], error) {
	definition, err := representation.NewLanguage[s.Language](Provider{}, s.Profile)
	if err != nil {
		return print.FilesPrinter[s.Language]{}, err
	}
	provider, err := filesprint.New[s.Language](Provider{}, Provider{})
	if err != nil {
		return print.FilesPrinter[s.Language]{}, err
	}
	return print.NewFilesPrinter(definition, provider)
}
