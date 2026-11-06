package typescriptsem

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	parser "github.com/Bitspark/bittype/parse/syntax/typescript/go"
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
	if _, err := Check(context.Background(), node, module, trees, source, nil); !errors.Is(err, n.ErrCompiler) {
		t.Fatal("foreign compiler", err)
	}
	if _, err := Check(context.Background(), filepath.Join(t.TempDir(), "missing-node"), module, trees, source, nil); !errors.Is(err, n.ErrCompiler) {
		t.Fatal("missing executable", err)
	}
	if _, err := Check(context.Background(), node, filepath.Join(t.TempDir(), "missing.cjs"), trees, source, nil); !errors.Is(err, n.ErrCompiler) {
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
	ctx, cancel := context.WithCancel(context.Background())
	b := boundedBuffer{remaining: 3, cancel: cancel}
	if _, err := b.Write([]byte("four")); !errors.Is(err, semantic.ErrLimit) || !b.exceeded || ctx.Err() != context.Canceled {
		t.Fatal("unbounded compiler response")
	}
}
