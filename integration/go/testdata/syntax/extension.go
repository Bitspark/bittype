package consumer

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/language/go"
)

// A separately owned provider requires no edit or import in the identity package.
type ExampleLanguage struct{}

func (ExampleLanguage) ID() language.ID { return "example.external" }
func (ExampleLanguage) Name() string    { return "External example" }

var provider language.LanguageProvider = ExampleLanguage{}
var tree ast.AST[ExampleLanguage]
var node ast.Node[ExampleLanguage]

const declaration ast.Kind[ExampleLanguage] = "example-declaration"
