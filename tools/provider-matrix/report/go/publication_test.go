package report

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestPairPublicationStopsOnSecondWriteFailure(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "docs"), 0700); err != nil {
		t.Fatal(err)
	}
	docs := Documents{root: root, matrixPath: "docs/providers.md", detailsPath: "docs/provider-capabilities.md", Matrix: "matrix", Capabilities: "details"}
	failed := errors.New("second write failed")
	calls := 0
	err := writeDocuments(docs, func(path string, data []byte, mode os.FileMode) error {
		calls++
		if calls == 2 {
			return failed
		}
		return os.WriteFile(path, data, mode)
	})
	if !errors.Is(err, failed) || calls != 2 {
		t.Fatal("publication did not stop", err, calls)
	}
	if _, err := os.Stat(filepath.Join(root, docs.detailsPath)); !os.IsNotExist(err) {
		t.Fatal("failed details write was hidden")
	}
}

func TestRenderedPairRequiresValidBothRegions(t *testing.T) {
	if _, err := replaceRegion("<!-- generated:providers:begin -->\n<!-- generated:providers:end -->", "transformers", "new"); err == nil {
		t.Fatal("missing region accepted")
	}
	if err := Write(Documents{Matrix: "one", Capabilities: "two"}); err == nil {
		t.Fatal("unbound output paths accepted")
	}
}
