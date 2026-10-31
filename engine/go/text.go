package engine

import (
 "fmt"

 "github.com/Bitspark/bittype/ast/go"
 cstcheck "github.com/Bitspark/bittype/check/cst/go"
 "github.com/Bitspark/bittype/cst/go"
 "github.com/Bitspark/bittype/internal/tree"
 "github.com/Bitspark/bittype/language/go"
 "github.com/Bitspark/bittype/text/go"
 "github.com/Bitspark/bittype/transform/go"
)

// CST is a checked concrete snapshot with an already constructed AST view.
// Concrete storage is opaque; UnpackCST returns an independent typed copy.
type CST struct {
 abstract transform.AST
 concrete any
}
func (c CST) LanguageID() language.ID { return c.abstract.LanguageID() }
func (c CST) Profile() string { return c.abstract.Profile() }
func (c CST) ASTView() transform.AST { return c.abstract }

func PackCST[L any](definition transform.Language, input cst.CST[L]) (CST,error) {
 if err := cstcheck.Check(input); err != nil { return CST{},err }
 abstract,err := transform.Pack(definition,input)
 if err != nil { return CST{},err }
 return CST{abstract:abstract,concrete:cloneConcrete(input)},nil
}

func UnpackCST[L any](definition transform.Language, input CST) (cst.CST[L],error) {
 if _,err := transform.Unpack[L](definition,input.abstract); err != nil { return cst.CST[L]{},err }
 concrete,ok := input.concrete.(cst.CST[L])
 if !ok { return cst.CST[L]{},transform.ErrTypeMismatch }
 return cloneConcrete(concrete),nil
}
func cloneConcrete[L any](input cst.CST[L]) cst.CST[L] {
 result := input
 result.Abstract = tree.Clone(input.Abstract)
 result.Nodes = make([]cst.Node[L],len(input.Nodes))
 for i,n := range input.Nodes {
  result.Nodes[i] = n
  if n.Text != nil { value := *n.Text; result.Nodes[i].Text = &value }
  result.Nodes[i].Children = append([]cst.NodeID(nil),n.Children...)
  result.Nodes[i].Syntax = append([]ast.NodeID(nil),n.Syntax...)
  result.Nodes[i].Origin.Spans = append([]ast.Span(nil),n.Origin.Spans...)
  result.Nodes[i].Origin.Nodes = append([]ast.NodeRef(nil),n.Origin.Nodes...)
 }
 return result
}

func (e *Engine) textProvider(id language.ID) (Provider,error) {
 registration,err := e.lookup(id)
 if err != nil { return Provider{},err }
 if registration.text == nil { return Provider{},fmt.Errorf("%w: %s",ErrMissingCapability,id) }
 return *registration.text,nil
}
func (e *Engine) Parse(id language.ID,input text.Input) (CST,error) {
 provider,err := e.textProvider(id)
 if err != nil { return CST{},err }
 return provider.parse(input)
}
func (e *Engine) Format(input transform.AST) (CST,error) {
 provider,err := e.textProvider(input.LanguageID())
 if err != nil { return CST{},err }
 if err := provider.definition.Check(input); err != nil { return CST{},err }
 return provider.format(input)
}
func (e *Engine) Print(input CST) (string,error) {
 provider,err := e.textProvider(input.LanguageID())
 if err != nil { return "",err }
 return provider.print(input)
}
func (e *Engine) Render(input transform.AST) (string,error) {
 formatted,err := e.Format(input)
 if err != nil { return "",err }
 return e.Print(formatted)
}
