package resolve_test

import (
	"testing"

	provider "github.com/Bitspark/bittype/provider/horos/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Syntax admission must not require resolution to have succeeded already.
func TestUnresolvedHorosSyntaxReachesAnalysis(t *testing.T) {
	tree := fixture(t)
	for i := range tree.Nodes {
		if tree.Nodes[i].Kind != s.KindProgram {
			tree.Nodes[i].References = nil
		}
	}
	application := find(tree, s.KindApplication, "Page")
	setAttr(&tree.Nodes[application-1], s.AttributeName, "Missing")
	definition, err := representation.NewLanguage[s.Language](provider.Provider{}, s.Profile)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := representation.Pack(definition, tree); err != nil {
		t.Fatalf("unresolved syntax must reach explicit analysis: %v", err)
	}
	if model, err := resolve.Model(tree); err == nil || model != nil {
		t.Fatal("explicit semantic resolution accepted the unknown declaration")
	}
}
