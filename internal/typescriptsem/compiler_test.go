package typescriptsem

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	parser "github.com/Bitspark/bittype/parse/syntax/typescript/provider/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
)

func TestCompilerFailureAndVersion(t *testing.T) {
	node := os.Getenv("BITTYPE_TYPESCRIPT_NODE")
	if node == "" {
		t.Fatal("BITTYPE_TYPESCRIPT_NODE must be provisioned")
	}
	text := "type Text = string;"
	parsed, err := parser.Parse(context.Background(), parse.Input{Text: text, SourceID: "text", TreeID: "text"})
	if err != nil {
		t.Fatal(err)
	}
	trees := []s.AST{parsed.ASTView()}
	source := &files.Context{Files: files.Files{Entries: []files.File{{ID: "text", Path: "text.ts", Text: text}}}, Links: []files.SyntaxLink{{File: "text", Tree: "text", Source: "text"}}}
	module := filepath.Join(t.TempDir(), "wrong.cjs")
	if err := os.WriteFile(module, []byte("module.exports = {version: '0.0.0'};"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Check(context.Background(), node, module, trees, source, nil, n.Configuration()); !errors.Is(err, n.ErrCompiler) {
		t.Fatal("foreign compiler", err)
	}
	if _, err := Check(context.Background(), filepath.Join(t.TempDir(), "missing-node"), module, trees, source, nil, n.Configuration()); !errors.Is(err, n.ErrCompiler) {
		t.Fatal("missing executable", err)
	}
	if _, err := Check(context.Background(), node, filepath.Join(t.TempDir(), "missing.cjs"), trees, source, nil, n.Configuration()); !errors.Is(err, n.ErrCompiler) {
		t.Fatal("missing module", err)
	}
}
func TestPrepareInvalidTextAndDeadline(t *testing.T) {
	text := "type Text = string;"
	parsed, err := parser.Parse(context.Background(), parse.Input{Text: text, SourceID: "text", TreeID: "text"})
	if err != nil {
		t.Fatal(err)
	}
	source := &files.Context{Files: files.Files{Entries: []files.File{{ID: "text", Path: "text.ts", Text: "\xff"}}}, Links: []files.SyntaxLink{{File: "text", Tree: "text", Source: "text"}}}
	if _, err := prepare(context.Background(), []s.AST{parsed.ASTView()}, source); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal(err)
	}
}

// The helper writes six statuses and spells them itself. Tests of the resolver
// fail when Check reads ok, invalid, unsupported or mismatch under another
// spelling. No other test reads the two below: incomplete, when the helper stops
// at its own depth bound, and environment, which Check turns into ErrCompiler
// together with the helper's diagnostic.
func TestHelperStatusesIncompleteAndEnvironment(t *testing.T) {
	node, module := os.Getenv("BITTYPE_TYPESCRIPT_NODE"), os.Getenv("BITTYPE_TYPESCRIPT_MODULE")
	if node == "" || module == "" {
		t.Fatal("BITTYPE_TYPESCRIPT_NODE and BITTYPE_TYPESCRIPT_MODULE must be provisioned")
	}
	check := func(module, text string) (Output, error) {
		parsed, err := parser.Parse(context.Background(), parse.Input{Text: text, SourceID: "text", TreeID: "text"})
		if err != nil {
			t.Fatal(err)
		}
		source := &files.Context{Files: files.Files{Entries: []files.File{{ID: "text", Path: "text.ts", Text: text}}}, Links: []files.SyntaxLink{{File: "text", Tree: "text", Source: "text"}}}
		return Check(context.Background(), node, module, []s.AST{parsed.ASTView()}, source, nil, n.Configuration())
	}
	// Each alias is an array of the next one, so the first use the helper meets
	// reaches 130 levels of native types. That is within the bounds of prepare
	// and beyond the helper's depth of 128. A diagnostic of the bridge itself
	// would have the code limit.
	var chain strings.Builder
	for i := 130; i > 0; i-- {
		fmt.Fprintf(&chain, "type A%d = A%d[];\n", i, i-1)
	}
	out, err := check(module, chain.String()+"type A0 = string;")
	if err != nil || out.Status != statusIncomplete || out.Facts != nil || len(out.Diagnostics) != 1 || out.Diagnostics[0].Code == n.Key("limit") {
		t.Fatal("the helper's depth bound", out.Status, out.Diagnostics, err)
	}
	foreign := filepath.Join(t.TempDir(), "wrong.cjs")
	if err := os.WriteFile(foreign, []byte("module.exports = {version: '0.0.0'};"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := check(foreign, "type Text = string;"); !errors.Is(err, n.ErrCompiler) || err.Error() == n.ErrCompiler.Error() {
		t.Fatal("the helper's environment diagnostic did not reach the caller", err)
	}
}
