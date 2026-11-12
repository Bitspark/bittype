// Package record supplies the independently selectable static Horos/Go
// closed-record representation law.
package record

import (
	"github.com/Bitspark/bittype/semantics/go"
	v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
)

const Law = "horos-go/record-fields-1"

func key(name string) semantics.Key { return semantics.Key{Owner: "horos-go/static", Name: name} }

func Rule() v.Rule {
	return v.Rule{Key: key("closed-record"), Dependency: Law, Law: v.Congruence{}, Project: project}
}
func project(p v.Subject) ([]v.Subject, bool, error) {
	if p.Optional() || !p.Declaration().Valid() || !p.Target().IsNamed() {
		return nil, false, nil
	}
	source, err := p.Declaration().Fields()
	native, ok := p.Target().Fields()
	bindings := p.Bindings()
	if err != nil || !ok || len(bindings) != len(source) || len(bindings) != len(native) {
		return nil, false, nil
	}
	for _, f := range native {
		if f.Embedded || f.Tag != "" {
			return nil, false, nil
		}
	}
	usedSource, usedTarget := map[string]bool{}, map[string]bool{}
	children := make([]v.Subject, 0, len(bindings))
	for _, b := range bindings {
		if usedSource[b.Source] || usedTarget[b.Target] {
			return nil, false, nil
		}
		usedSource[b.Source], usedTarget[b.Target] = true, true
		child, err := p.Field(b.Source, b.Target)
		if err != nil {
			return nil, false, nil
		}
		children = append(children, child)
	}
	return children, true, nil
}
