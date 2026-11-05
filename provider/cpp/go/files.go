package provider

import (
	"github.com/Bitspark/bittype/contracts/print/go"
	filesprint "github.com/Bitspark/bittype/print/files/go"
	syntaxprint "github.com/Bitspark/bittype/print/syntax/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
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
