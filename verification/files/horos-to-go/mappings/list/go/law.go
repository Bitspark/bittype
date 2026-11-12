// Package list supplies the independently selectable static Horos/Go
// list-slice representation law.
package list

import (
	"github.com/Bitspark/bittype/semantics/go"
	v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

const Law = "horos-go/list-slice-1"

func key(name string) semantics.Key { return semantics.Key{Owner: "horos-go/static", Name: name} }

func Rule() v.Rule {
	return v.Rule{Key: key("list-slice"), Dependency: Law, Law: v.Congruence{}, Project: project}
}
func project(p v.Subject) ([]v.Subject, bool, error) {
	element, ok := p.Target().SliceElement()
	if p.Optional() || !ok || p.Source().Kind() != decl.Descriptor || p.Source().Constructor() != "list-of" || len(p.Source().Operands()) != 1 {
		return nil, false, nil
	}
	child, err := p.Operand(0, element)
	return []v.Subject{child}, true, err
}
