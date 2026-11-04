package files_test

import (
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"testing"
)

func TestFilesCapture(t *testing.T) {
	input := files.Files{Entries: []files.File{{ID: "a", Path: "nested/a.horos", Text: "original"}}}
	owned, err := input.Capture(files.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	input.Entries[0].Text = "changed"
	if owned.Entries[0].Text != "original" {
		t.Fatal("borrowed files")
	}
	for _, path := range []files.Path{"", "/abs", "C:/abs", "//host/path", "a\\b", "a/../b", "a/./b", "a//b", "nul\x00"} {
		if files.ValidPath(path) {
			t.Fatal("invalid path admitted", path)
		}
	}
	for _, entries := range [][]files.File{
		{{ID: "a", Path: "a"}, {ID: "a", Path: "b"}},
		{{ID: "a", Path: "a"}, {ID: "b", Path: "a"}},
		{{ID: "a", Path: "parent"}, {ID: "b", Path: "parent/child"}},
	} {
		if _, err := (files.Files{Entries: entries}).Capture(files.DefaultLimits()); err == nil {
			t.Fatal("duplicate or conflicting inventory admitted")
		}
	}
	c := files.Context{Files: owned, Links: []files.SyntaxLink{{File: "a", Tree: "tree", Source: "source"}}}
	if _, err := c.Capture(files.DefaultLimits(), map[ast.TreeID]ast.NodeID{"tree": 1}); err != nil {
		t.Fatal(err)
	}
	c.Links[0].Tree = "other"
	if _, err := c.Capture(files.DefaultLimits(), map[ast.TreeID]ast.NodeID{"tree": 1}); err == nil {
		t.Fatal("unknown tree admitted")
	}
}
