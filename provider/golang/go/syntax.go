package provider

import (
	"context"
	checking "github.com/Bitspark/bittype/check/syntax/golang/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	formatting "github.com/Bitspark/bittype/format/syntax/golang/go"
	parsing "github.com/Bitspark/bittype/parse/syntax/golang/go"
	printing "github.com/Bitspark/bittype/print/syntax/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func (Provider) Parse(ctx context.Context, input parse.Input) (cst.CST[s.Language], error) {
	return parsing.Parse(ctx, input)
}
func (Provider) Check(view ast.View[s.Language]) error { return checking.Check(view) }
func (Provider) Format(ctx context.Context, view ast.View[s.Language]) (cst.CST[s.Language], error) {
	return formatting.Format(ctx, view)
}
func (Provider) Print(ctx context.Context, tree cst.CST[s.Language]) (string, error) {
	return printing.Print(ctx, tree)
}

func NewSyntaxParser(definition representation.Language) (parse.SyntaxParser, error) {
	return parse.NewSyntaxParser[s.Language](definition, parsing.Provider{})
}
func NewSyntaxChecker(definition representation.Language) (check.Checker, error) {
	return check.NewChecker[s.Language](definition, checking.Provider{})
}
func NewSyntaxFormatter(definition representation.Language) (format.SyntaxFormatter, error) {
	return format.NewSyntaxFormatter[s.Language](definition, formatting.Provider{})
}
func NewSyntaxPrinter(definition representation.Language) (print.SyntaxPrinter, error) {
	provider, err := printing.New[s.Language](s.Language{})
	if err != nil {
		return print.SyntaxPrinter{}, err
	}
	return print.NewSyntaxPrinter[s.Language](definition, provider)
}

var _ parse.Provider[parse.Input, cst.CST[s.Language]] = Provider{}
var _ check.Provider[s.Language] = Provider{}
var _ format.Provider[ast.View[s.Language], cst.CST[s.Language]] = Provider{}
var _ print.Provider[cst.CST[s.Language], string] = Provider{}
