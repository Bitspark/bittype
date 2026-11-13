package check_test

import (
	"errors"
	"strings"
	"testing"

	check "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// The table names go/identifier, and the checker binds it to Go's rule: an
// identifier that is no keyword. The shared identifier shape alone would
// admit a keyword.
func TestIdentifierIsNoKeyword(t *testing.T) {
	file := func(name string) ast.AST[s.Language] {
		return ast.AST[s.Language]{ID: "file", Root: 1, Nodes: []ast.Node[s.Language]{
			{Kind: s.KindFile, Attributes: []ast.Attribute{{Name: s.AttributePackage, Text: name}}},
		}}
	}
	for name, admitted := range map[string]bool{
		"main": true, "funcs": true, "Type": true, "_": true, "é": true,
		"func": false, "type": false, "package": false, "range": false, "": false, "1a": false,
	} {
		err := check.Check(file(name))
		if admitted && err != nil {
			t.Errorf("package name %q: %v", name, err)
		}
		var diagnostic *text.Diagnostic
		if !admitted && (!errors.As(err, &diagnostic) || diagnostic.Code != text.InvalidAST || !strings.HasSuffix(err.Error(), "unknown or invalid attribute: package")) {
			t.Errorf("package name %q: got %v, want an invalid attribute", name, err)
		}
	}
}
