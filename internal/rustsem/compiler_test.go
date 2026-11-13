package rustsem

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	parser "github.com/Bitspark/bittype/parse/syntax/rust/provider/go"
	n "github.com/Bitspark/bittype/semantics/rust/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
)

func compilerFixture(t *testing.T) ([]s.AST, *files.Context) {
	t.Helper()
	text := "pub type Byte = u8;"
	parsed, err := parser.Parse(context.Background(), parse.Input{Text: text, SourceID: "test", TreeID: "test"})
	if err != nil {
		t.Fatal(err)
	}
	return []s.AST{parsed.ASTView()}, &files.Context{Files: files.Files{Entries: []files.File{{ID: "test", Path: "test.rs", Text: text}}}, Links: []files.SyntaxLink{{File: "test", Tree: "test", Source: "test"}}}
}

var environment = Environment{Edition: n.Edition2024, Target: n.TargetLinuxX64, CrateName: "pages", CrateKind: n.KindLibrary}

func TestCompilerFailureAndProtocol(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("qualified helper execution platform is Linux")
	}
	trees, source := compilerFixture(t)
	ctx := context.Background()
	sysroot := t.TempDir()
	for _, response := range []string{`{}`, `{"Version":1,"Compiler":"rustc 0.0.0 (x)","Status":"ok","Facts":{}}`, `{"Version":1,"Compiler":"rustc 1.90.0 (x)","Status":"ok"}`, `{"Version":1,"Compiler":"rustc 1.90.0 (x)","Status":"unsupported","Facts":{}}`, `{"Version":1,"Compiler":"rustc 1.90.0 (x)","Status":"ok","Facts":{},"unexpected":true}`} {
		path := filepath.Join(t.TempDir(), "helper")
		if err := os.WriteFile(path, []byte("#!/bin/sh\ncat >/dev/null\nprintf '%s' '"+response+"'\n"), 0700); err != nil {
			t.Fatal(err)
		}
		if _, err := Check(ctx, path, sysroot, environment, trees, source); !errors.Is(err, n.ErrCompiler) {
			t.Fatal(response, err)
		}
	}
	if _, err := Check(ctx, filepath.Join(t.TempDir(), "missing"), sysroot, environment, trees, source); !errors.Is(err, n.ErrCompiler) {
		t.Fatal(err)
	}
	if _, err := Check(ctx, os.Getenv("BITTYPE_RUST_HELPER"), sysroot, environment, trees, source); !errors.Is(err, n.ErrCompiler) {
		t.Fatal("wrong sysroot", err)
	}
	t.Setenv("BITTYPE_SHOULD_NOT_LEAK", "secret-placeholder")
	path := filepath.Join(t.TempDir(), "helper")
	script := "#!/bin/sh\ncat >/dev/null\ntest -z \"$BITTYPE_SHOULD_NOT_LEAK\" || exit 1\nprintf '%s' '{\"Version\":1,\"Compiler\":\"rustc 1.90.0 (x)\",\"Status\":\"incomplete\"}'\n"
	if err := os.WriteFile(path, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	if out, err := Check(ctx, path, sysroot, environment, trees, source); err != nil || out.Status != "incomplete" {
		t.Fatal(out, err)
	}
}
func TestCaptureLimitsAndOffsets(t *testing.T) {
	trees, source := compilerFixture(t)
	source.Files.Entries[0].Text = "\xff"
	if _, err := prepare(context.Background(), trees, source); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal(err)
	}
	text := "// café\r\npub type Byte = u8;\r\n"
	parsed, err := parser.Parse(context.Background(), parse.Input{Text: text, TreeID: "test", SourceID: "test"})
	if err != nil {
		t.Fatal(err)
	}
	source.Files.Entries[0].Text = text
	units, err := prepare(context.Background(), []s.AST{parsed.ASTView()}, source)
	if err != nil {
		t.Fatal(err)
	}
	if units[0].Text != "// café\npub type Byte = u8;\n" {
		t.Fatal(units[0].Text)
	}
	ctx, cancel := context.WithCancel(context.Background())
	b := boundedBuffer{remaining: 3, cancel: cancel}
	if _, err := b.Write([]byte("four")); !errors.Is(err, semantic.ErrLimit) || !b.exceeded || ctx.Err() != context.Canceled {
		t.Fatal("unbounded output")
	}
}
