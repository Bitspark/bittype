// Package print emits a concrete tree exactly as supplied. It knows no language
// names, token kinds, indentation, escaping or precedence rules.
package print

import (
	cstcheck "github.com/Bitspark/bittype/check/cst/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"strings"
)

func Print[L any](tree cst.CST[L]) (string, error) {
	if err := cstcheck.Check(tree); err != nil {
		return "", err
	}
	var result strings.Builder
	stack := []cst.NodeID{tree.Root}
	for len(stack) > 0 {
		id := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		n := tree.Nodes[id-1]
		if n.Text != nil {
			result.WriteString(*n.Text)
			continue
		}
		for i := len(n.Children) - 1; i >= 0; i-- {
			stack = append(stack, n.Children[i])
		}
	}
	return result.String(), nil
}
