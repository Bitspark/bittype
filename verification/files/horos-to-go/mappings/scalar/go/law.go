// Package scalar supplies the independently selectable static Horos/Go scalar
// representation law. Text and Boolean are two rules of that one law.
package scalar

import (
	"go/types"

	"github.com/Bitspark/bittype/semantics/go"
	v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
)

const Law = "horos-go/static-scalars-1"

func key(name string) semantics.Key { return semantics.Key{Owner: "horos-go/static", Name: name} }

func scalar(name string, kind types.BasicKind) v.Rule {
	return v.Rule{Key: key(name), Dependency: Law, Law: v.Congruence{}, Project: func(p v.Subject) ([]v.Subject, bool, error) {
		if p.Optional() || p.Source().Kind() != decl.Descriptor || p.Source().Constructor() != "scalar" {
			return nil, false, nil
		}
		value, ok := p.Source().Value().(core.Tuple)
		if !ok || value.Len() != 2 {
			return nil, false, nil
		}
		atom, ok := value.At(1).(core.Atom)
		return nil, ok && string(atom.Bytes()) == name && p.Target().IsBasic(kind), nil
	}}
}
func Text() v.Rule    { return scalar("utf8-text", types.String) }
func Boolean() v.Rule { return scalar("bool", types.Bool) }
