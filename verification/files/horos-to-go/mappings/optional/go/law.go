// Package optional supplies the independently selectable static Horos/Go
// optional-presence representation law. The consumer supplies the helper type
// of its captured Go model and names its presence and value fields.
package optional

import (
	"go/types"

	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
)

const Law = "horos-go/optional-presence-1"

func key(name string) semantics.Key { return semantics.Key{Owner: "horos-go/static", Name: name} }

type helper struct {
	constructor               n.Type
	presenceField, valueField string
}

func Rule(constructor n.Type, presenceField, valueField string) v.Rule {
	return v.Rule{Key: key("optional-presence"), Dependency: Law, Law: v.Congruence{}, Project: helper{constructor, presenceField, valueField}.project}
}
func (o helper) project(p v.Subject) ([]v.Subject, bool, error) {
	if !p.Optional() || !p.Target().Origin().Same(o.constructor) {
		return nil, false, nil
	}
	params := o.constructor.Parameters()
	fields, ok := o.constructor.Fields()
	if !o.constructor.IsNamed() || !ok || len(params) != 1 || !params[0].Unconstrained() || len(fields) != 2 || o.presenceField == o.valueField || len(p.Target().Arguments()) != 1 {
		return nil, false, nil
	}
	present, value := false, false
	for _, f := range fields {
		if f.Embedded || f.Tag != "" {
			return nil, false, nil
		}
		if f.Name == o.presenceField {
			present = f.Type.IsBasic(types.Bool)
		}
		if f.Name == o.valueField {
			value = f.Type.Same(params[0])
		}
	}
	if !present || !value {
		return nil, false, nil
	}
	child, err := p.Payload(p.Target().Arguments()[0])
	return []v.Subject{child}, true, err
}
