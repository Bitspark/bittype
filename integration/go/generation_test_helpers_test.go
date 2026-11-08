package integration_test

import "github.com/Bitspark/bittype/trees/ast/go"

func nodeAttribute[L any](node ast.Node[L], name string) string {
	for _, attribute := range node.Attributes {
		if attribute.Name == name {
			return attribute.Text
		}
	}
	return ""
}
