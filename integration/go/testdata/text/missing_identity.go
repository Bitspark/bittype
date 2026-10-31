package consumer

import (
	"github.com/Bitspark/bittype/cst/go"
	"github.com/Bitspark/bittype/language/go"
	"github.com/Bitspark/bittype/text/go"
)
type Language struct{}
type MissingName struct{}
func (MissingName) ID() language.ID {return "missing"}
func (MissingName) Parse(text.Input)(cst.CST[Language],error){return cst.CST[Language]{},nil}
var _ text.ParserProvider[Language]=MissingName{}
