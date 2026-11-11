package integration_test

import (
	n "github.com/Bitspark/bittype/semantics/golang/go"
	v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
	m "github.com/Bitspark/bittype/verification/files/horos-to-go/mappings/go"
)

// This file is the one place in this package that names the horos-to-go
// verification laws. Tests select each law and its dependency name through it,
// so a change to where a law lives edits this file only.
const (
	scalarLawName   = m.ScalarLaw
	listLawName     = m.ListLaw
	optionalLawName = m.OptionalLaw
	recordLawName   = m.RecordLaw
)

func textLaw() v.Rule    { return m.Text() }
func booleanLaw() v.Rule { return m.Boolean() }
func listLaw() v.Rule    { return m.List() }
func recordLaw() v.Rule  { return m.Record() }
func optionalLaw(constructor n.Type, presenceField, valueField string) v.Rule {
	return (m.Optional{Constructor: constructor, PresenceField: presenceField, ValueField: valueField}).Rule()
}
