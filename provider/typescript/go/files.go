package provider

import (
	"context"
	analyzing "github.com/Bitspark/bittype/analyze/files/typescript/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	filesprint "github.com/Bitspark/bittype/print/files/go"
	syntaxprint "github.com/Bitspark/bittype/print/syntax/go"
	resolving "github.com/Bitspark/bittype/resolve/files/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
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

func (p Provider) Resolve(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	return (resolving.Provider{Compiler: p.Compiler}).Resolve(ctx, request)
}
func (Provider) Analyze(ctx context.Context, request analysis.Request[s.Language]) (semantic.Result, error) {
	return analyzing.Analyze(ctx, request)
}

// NewFilesResolver takes an explicitly provisioned compiler; it performs no IO.
func NewFilesResolver(definition representation.Language, compiler resolving.Compiler) (resolution.Resolver, error) {
	if err := compiler.Validate(); err != nil {
		return resolution.Resolver{}, err
	}
	return resolution.NewResolver[s.Language](definition, resolving.Contract(), resolving.Provider{Compiler: compiler})
}
func NewFilesAnalyzer(definition representation.Language) (analysis.Analyzer, error) {
	return analysis.NewAnalyzer[s.Language](definition, analyzing.Contract(), analyzing.Provider{})
}

var _ resolution.Provider[s.Language] = Provider{}
var _ analysis.Provider[s.Language] = Provider{}
