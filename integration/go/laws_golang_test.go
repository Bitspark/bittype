package integration_test

import (
	n "github.com/Bitspark/bittype/semantics/golang/go"
	v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
	"github.com/Bitspark/bittype/verification/files/horos-to-go/mappings/list/go"
	"github.com/Bitspark/bittype/verification/files/horos-to-go/mappings/optional/go"
	"github.com/Bitspark/bittype/verification/files/horos-to-go/mappings/record/go"
	"github.com/Bitspark/bittype/verification/files/horos-to-go/mappings/scalar/go"
)

// This file is the one place in this package that names the horos-to-go
// verification laws. Tests select each law and its dependency name through it,
// so a change to where a law lives edits this file only.
const (
	scalarLawName   = scalar.Law
	listLawName     = list.Law
	optionalLawName = optional.Law
	recordLawName   = record.Law
)

func textLaw() v.Rule    { return scalar.Text() }
func booleanLaw() v.Rule { return scalar.Boolean() }
func listLaw() v.Rule    { return list.Rule() }
func recordLaw() v.Rule  { return record.Rule() }
func optionalLaw(constructor n.Type, presenceField, valueField string) v.Rule {
	return optional.Rule(constructor, presenceField, valueField)
}
