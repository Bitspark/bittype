package consumer

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

type Language struct{}
type MissingName struct{}

func (MissingName) ID() language.ID                             { return "missing" }
func (MissingName) Parse(text.Input) (cst.CST[Language], error) { return cst.CST[Language]{}, nil }

var _ text.ParserProvider[Language] = MissingName{}
