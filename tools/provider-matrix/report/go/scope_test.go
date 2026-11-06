package report

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScopeDocumentsAndFragments(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "scope.md"), []byte("# Scope\n## Repeated\n## Repeated\n<a id=\"owned\"></a>\n```\n# Hidden\n```\n"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, ref := range []string{"scope.md", "scope.md#scope", "scope.md#repeated-1", "scope.md#owned"} {
		if _, err := scopeDocument(root, ref); err != nil {
			t.Fatal(ref, err)
		}
	}
	for _, ref := range []string{"scope.md#missing", "scope.md#hidden", "../scope.md", "/scope.md", "C:/scope.md", "missing.md"} {
		if _, err := scopeDocument(root, ref); err == nil {
			t.Fatal("accepted", ref)
		}
	}
	out := filepath.Join(t.TempDir(), "outside.md")
	if err := os.WriteFile(out, []byte("# Outside"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(out, filepath.Join(root, "link.md")); err == nil {
		if _, err := scopeDocument(root, "link.md"); err == nil {
			t.Fatal("accepted symlink escape")
		}
	}
}
