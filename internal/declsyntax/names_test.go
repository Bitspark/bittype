package declsyntax

import (
	"errors"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/text/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// The table names horos/identifier and horos/pin, and this package binds them:
// the first to the identifier rule of the horos lexical helper, which is
// ASCII, and the second to a pin of 64 lowercase hexadecimal digits. The
// shared identifier shape would admit a letter outside ASCII, and the rule of
// hexadecimal text would admit a shorter pin.
func TestIdentifierAndPinAreTheHorosRules(t *testing.T) {
	pin := strings.Repeat("0a", 32)
	// A program of one module that imports one module under an alias.
	program := func(alias, modulePin string) h.AST {
		return h.AST{ID: "program", Root: 3, Nodes: []h.Node{
			{Kind: h.KindImport, Attributes: []ast.Attribute{{Name: AttributeAlias, Text: alias}, {Name: AttributeName, Text: "example.other"}, {Name: AttributePin, Text: pin}}},
			{Kind: h.KindModule, Attributes: []ast.Attribute{{Name: AttributeName, Text: "example.pages"}, {Name: AttributePin, Text: modulePin}}, Children: []ast.Child{{Role: RoleImport, Node: 1}}},
			{Kind: h.KindProgram, Children: []ast.Child{{Role: RoleModule, Node: 2}}, References: []ast.Reference{{Role: RoleRoot, Target: ast.NodeRef{Tree: "program", Node: 2}}}},
		}}
	}
	if err := Check(program("other", pin)); err != nil {
		t.Fatal(err)
	}
	refused := func(what string, tree h.AST, attribute string) {
		t.Helper()
		err := Check(tree)
		var diagnostic *text.Diagnostic
		if !errors.As(err, &diagnostic) || diagnostic.Code != text.InvalidAST || !strings.HasSuffix(err.Error(), "unknown or invalid attribute: "+attribute) {
			t.Errorf("%s: got %v, want an invalid attribute %s", what, err, attribute)
		}
	}
	for _, alias := range []string{"é", "naïve", "名前", "1a", ""} {
		refused("alias "+alias, program(alias, pin), AttributeAlias)
	}
	for _, alias := range []string{"a", "_other", "Other2"} {
		if err := Check(program(alias, pin)); err != nil {
			t.Errorf("alias %s: %v", alias, err)
		}
	}
	for what, modulePin := range map[string]string{
		"62 digits":        strings.Repeat("0a", 31),
		"66 digits":        strings.Repeat("0a", 33),
		"no digit":         "",
		"64 in upper case": strings.Repeat("0A", 32),
		"63 and no hex":    strings.Repeat("0a", 31) + "0g",
		"64 with an x":     strings.Repeat("0a", 31) + "0x",
	} {
		refused("a pin of "+what, program("other", modulePin), AttributePin)
	}
}
