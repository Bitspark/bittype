package provider

import (
	"context"
	analyzing "github.com/Bitspark/bittype/analyze/files/horos/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	filesformat "github.com/Bitspark/bittype/format/files/horos/go"
	syntaxformat "github.com/Bitspark/bittype/format/syntax/horos/go"
	filesparse "github.com/Bitspark/bittype/parse/files/horos/go"
	syntaxparse "github.com/Bitspark/bittype/parse/syntax/horos/go"
	filesprint "github.com/Bitspark/bittype/print/files/go"
	syntaxprint "github.com/Bitspark/bittype/print/syntax/go"
	resolving "github.com/Bitspark/bittype/resolve/files/horos/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// NewFilesPrinter composes only printing; it requires no semantic provider.
func NewFilesPrinter(definition representation.Language) (print.FilesPrinter[s.Language], error) {
	syntax, err := syntaxprint.New[s.Language](s.Language{})
	if err != nil {
		return print.FilesPrinter[s.Language]{}, err
	}
	provider, err := filesprint.New[s.Language](s.Language{}, syntax)
	if err != nil {
		return print.FilesPrinter[s.Language]{}, err
	}
	return print.NewFilesPrinter[s.Language](definition, provider)
}

func (Provider) Resolve(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	return resolving.ResolveReferences(ctx, request)
}
func (Provider) Analyze(ctx context.Context, request analysis.Request[s.Language]) (semantic.Result, error) {
	return analyzing.Analyze(ctx, request)
}

func NewFilesResolver(definition representation.Language) (resolution.Resolver, error) {
	return resolution.NewResolver[s.Language](definition, resolving.Contract(), resolving.Provider{})
}
func NewFilesAnalyzer(definition representation.Language) (analysis.Analyzer, error) {
	return analysis.NewAnalyzer[s.Language](definition, analyzing.Contract(), analyzing.Provider{})
}

func fileStages(definition representation.Language) (check.Checker, resolution.Resolver, analysis.Analyzer, error) {
	checker, err := NewSyntaxChecker(definition)
	if err != nil {
		return checker, resolution.Resolver{}, analysis.Analyzer{}, err
	}
	resolver, err := NewFilesResolver(definition)
	if err != nil {
		return checker, resolver, analysis.Analyzer{}, err
	}
	analyzer, err := NewFilesAnalyzer(definition)
	return checker, resolver, analyzer, err
}

func NewFilesParser(definition representation.Language) (parse.FilesParser[s.Language], error) {
	checker, resolver, analyzer, err := fileStages(definition)
	if err != nil {
		return parse.FilesParser[s.Language]{}, err
	}
	printer, err := syntaxprint.New[s.Language](s.Language{})
	if err != nil {
		return parse.FilesParser[s.Language]{}, err
	}
	return filesparse.New(definition, syntaxparse.Provider{}, printer, checker, resolver, analyzer, semantic.Options{}, semantic.Options{})
}
func NewFilesFormatter(definition representation.Language) (format.FilesFormatter[s.Language], error) {
	checker, resolver, analyzer, err := fileStages(definition)
	if err != nil {
		return format.FilesFormatter[s.Language]{}, err
	}
	printer, err := syntaxprint.New[s.Language](s.Language{})
	if err != nil {
		return format.FilesFormatter[s.Language]{}, err
	}
	return filesformat.New(definition, syntaxparse.Provider{}, syntaxformat.Provider{}, printer, checker, resolver, analyzer, semantic.Options{}, semantic.Options{})
}

var _ resolution.Provider[s.Language] = Provider{}
var _ analysis.Provider[s.Language] = Provider{}
