package report

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPublicationStopsOnEachWriteFailure(t *testing.T) {
	for _, failAt := range []int{1, 2, 3} {
		t.Run(fmt.Sprint(failAt), func(t *testing.T) {
			root := t.TempDir()
			if err := os.Mkdir(filepath.Join(root, "docs"), 0700); err != nil {
				t.Fatal(err)
			}
			docs := Documents{root: root, matrixPath: "docs/providers.md", detailsPath: "docs/provider-capabilities.md", overviewPath: "docs/language-support.md", Matrix: "matrix", Capabilities: "details", Overview: "overview"}
			failed := errors.New("write failed")
			calls := 0
			err := writeDocuments(docs, func(path string, data []byte, mode os.FileMode) error {
				calls++
				if calls == failAt {
					return failed
				}
				return os.WriteFile(path, data, mode)
			})
			if !errors.Is(err, failed) || calls != failAt {
				t.Fatal("publication did not stop", err, calls)
			}
			if _, err := os.Stat(filepath.Join(root, docs.overviewPath)); !os.IsNotExist(err) {
				t.Fatal("failed write was hidden or publication continued")
			}
		})
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
