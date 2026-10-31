package consumer

import (
	"github.com/Bitspark/bittype/ast/go"
	"github.com/Bitspark/bittype/cst/go"
	"github.com/Bitspark/bittype/language/go"
	"github.com/Bitspark/bittype/pipeline/go"
	"github.com/Bitspark/bittype/text/go"
)
type Language struct{}
type Independent struct{}
func (Independent) ID() language.ID { return "independent-language" }
func (Independent) Name() string { return "Independent language" }
func (Independent) Parse(text.Input) (cst.CST[Language],error) {return cst.CST[Language]{},nil}
func (Independent) Format(ast.View[Language]) (cst.CST[Language],error) {return cst.CST[Language]{},nil}
func (Independent) Check(ast.View[Language]) error {return nil}
func (Independent) Print(cst.CST[Language]) (string,error) {return "",nil}

var _ language.LanguageProvider=Independent{}
var _ text.Provider[Language]=Independent{}
var _ ast.View[Language]=ast.AST[Language]{}
var _ ast.View[Language]=cst.CST[Language]{}
var _ pipeline.Stage[text.Input,cst.CST[Language]]=Independent{}.Parse
var _ pipeline.Stage[ast.View[Language],cst.CST[Language]]=Independent{}.Format
var _ pipeline.Stage[cst.CST[Language],string]=Independent{}.Print
