package engine

import (
 "github.com/Bitspark/bittype/language/go"
 "github.com/Bitspark/bittype/text/go"
 "github.com/Bitspark/bittype/transform/go"
)

// Provider packages existing typed text capabilities for heterogeneous
// registration. No parser, formatter, checker or printer is implemented here.
type Provider struct {
 definition transform.Language
 parse func(text.Input) (CST,error)
 format func(transform.AST) (CST,error)
 print func(CST) (string,error)
}

func NewProvider[L any](provider text.Provider[L], profile string) (Provider,error) {
 definition,err := transform.NewLanguage[L](provider,profile)
 if err != nil { return Provider{},err }
 result := Provider{definition: definition}
 result.parse = func(input text.Input) (CST,error) {
  parsed,err := provider.Parse(input)
  if err != nil { return CST{},err }
  return PackCST(definition,parsed)
 }
 result.format = func(input transform.AST) (CST,error) {
  abstract,err := transform.Unpack[L](definition,input)
  if err != nil { return CST{},err }
  formatted,err := provider.Format(abstract)
  if err != nil { return CST{},err }
  return PackCST(definition,formatted)
 }
 result.print = func(input CST) (string,error) {
  concrete,err := UnpackCST[L](definition,input)
  if err != nil { return "",err }
  output,err := provider.Print(concrete)
  if err != nil { return "",err }
  return output,nil
 }
 return result,nil
}
func (p Provider) ID() language.ID { return p.definition.ID() }
func (p Provider) Name() string { return p.definition.Name() }
func (p Provider) Language() transform.Language { return p.definition }
func (p Provider) validate() error {
 if err := p.definition.Validate(); err != nil { return err }
 if p.parse == nil || p.format == nil || p.print == nil { return ErrInvalidProvider }
 return nil
}

var _ language.LanguageProvider = Provider{}
