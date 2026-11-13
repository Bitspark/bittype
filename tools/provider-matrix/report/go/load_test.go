package report

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadRejectsInvalidFlagsAndNegativeLimits(t *testing.T) {
	for _, flags := range [][]string{{"-overlay=elsewhere.json"}, {"-test"}, {"-json"}, {"-C", "elsewhere"}, {"./different/..."}, {"-tags"}, {"-modfile=elsewhere.mod"}} {
		if err := validateBuildFlags(flags); err == nil {
			t.Fatal("unsafe source selection accepted", flags)
		}
	}
	if err := validateBuildFlags([]string{"-race", "-tags=integration", "-mod", "readonly"}); err != nil {
		t.Fatal(err)
	}
	if _, err := (LoadLimits{MaxFiles: -1}).effective(); err == nil {
		t.Fatal("negative load limit accepted")
	}
}

func TestCapturedSourceUsesSharedBudgetAndRejectsExternalPaths(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "first.go")
	if err := os.WriteFile(path, []byte("one"), 0600); err != nil {
		t.Fatal(err)
	}
	p := &program{root: root, ctx: context.Background(), limits: LoadLimits{MaxFiles: 2, MaxSourceBytes: 5}, sources: map[string]string{}}
	if got, err := p.readSource(path); err != nil || got != "one" {
		t.Fatal(got, err)
	}
	if err := os.WriteFile(path, []byte("changed"), 0600); err != nil {
		t.Fatal(err)
	}
	if got, err := p.readSource(path); err != nil || got != "one" || p.sourceBytes != 3 {
		t.Fatal("capture changed or was charged twice", got, err)
	}
	second := filepath.Join(root, "second.go")
	if err := os.WriteFile(second, []byte("two"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := p.readSource(second); err == nil || !strings.Contains(err.Error(), "source bytes") {
		t.Fatal("source-byte budget ignored", err)
	}
	if _, err := p.readSource(filepath.Join(root, "..", "outside.go")); err == nil || !strings.Contains(err.Error(), "outside root") {
		t.Fatal("external source path accepted", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p.ctx = ctx
	if _, err := p.readSource(path); err != context.Canceled {
		t.Fatal("cancellation ignored", err)
	}
}

func TestCapturedSourceRejectsSymlinkEscape(t *testing.T) {
	root, outside := t.TempDir(), t.TempDir()
	target := filepath.Join(outside, "external.go")
	if err := os.WriteFile(target, []byte("package external"), 0600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "inside.go")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	p := &program{root: root, ctx: context.Background(), limits: DefaultLoadLimits(), sources: map[string]string{}}
	if _, err := p.readSource(link); err == nil || !strings.Contains(err.Error(), "symlink resolves outside root") {
		t.Fatal("source symlink escaped selected module", err)
	}
}
