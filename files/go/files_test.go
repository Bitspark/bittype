package files_test

import (
	"errors"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"slices"
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

// TestCover qualifies the admission of a placement of items into files. The
// item is a string here and a declaration handle or a module pin elsewhere.
func TestCover(t *testing.T) {
	table := func() map[files.ID]files.Path {
		return map[files.ID]files.Path{"b": "out/b.txt", "a": "out/a.txt", "empty": "out/empty.txt"}
	}
	assignment := func() map[string]files.ID { return map[string]files.ID{"x": "b", "y": "a", "z": "b"} }
	items := []string{"x", "y", "z"}
	// A file without an item is admitted, and the inventory comes in the order
	// of the identities whatever order the table was built in.
	inventory, err := files.Cover(table(), assignment(), items, files.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	want := []files.File{{ID: "a", Path: "out/a.txt"}, {ID: "b", Path: "out/b.txt"}, {ID: "empty", Path: "out/empty.txt"}}
	if !slices.Equal(inventory.Entries, want) {
		t.Fatal("inventory", inventory.Entries)
	}
	refused := []struct {
		name   string
		want   error
		change func(map[files.ID]files.Path, map[string]files.ID, *[]string, *files.Limits)
	}{
		{"an item without an assignment", files.ErrInvalid, func(_ map[files.ID]files.Path, a map[string]files.ID, _ *[]string, _ *files.Limits) {
			delete(a, "y")
		}},
		{"an assigned item that is not placed", files.ErrInvalid, func(_ map[files.ID]files.Path, a map[string]files.ID, _ *[]string, _ *files.Limits) {
			a["foreign"] = "a"
		}},
		{"an item in place of another", files.ErrInvalid, func(_ map[files.ID]files.Path, a map[string]files.ID, _ *[]string, _ *files.Limits) {
			delete(a, "y")
			a["foreign"] = "a"
		}},
		{"an assignment to an undeclared file", files.ErrInvalid, func(_ map[files.ID]files.Path, a map[string]files.ID, _ *[]string, _ *files.Limits) {
			a["y"] = "absent"
		}},
		{"no file", files.ErrInvalid, func(f map[files.ID]files.Path, a map[string]files.ID, i *[]string, _ *files.Limits) {
			clear(f)
			clear(a)
			*i = nil
		}},
		{"an empty file identity", files.ErrInvalid, func(f map[files.ID]files.Path, _ map[string]files.ID, _ *[]string, _ *files.Limits) {
			f[""] = "out/c.txt"
		}},
		{"one path for two files", files.ErrInvalid, func(f map[files.ID]files.Path, _ map[string]files.ID, _ *[]string, _ *files.Limits) {
			f["empty"] = "out/a.txt"
		}},
		{"a path that is also a directory", files.ErrInvalid, func(f map[files.ID]files.Path, _ map[string]files.ID, _ *[]string, _ *files.Limits) {
			f["empty"] = "out/a.txt/c.txt"
		}},
		{"a path that leaves the root", files.ErrInvalid, func(f map[files.ID]files.Path, _ map[string]files.ID, _ *[]string, _ *files.Limits) {
			f["empty"] = "../c.txt"
		}},
		{"more files than the file limit", files.ErrLimit, func(_ map[files.ID]files.Path, _ map[string]files.ID, _ *[]string, l *files.Limits) {
			l.MaxFiles = 2
		}},
		{"more files than the link limit", files.ErrLimit, func(_ map[files.ID]files.Path, _ map[string]files.ID, _ *[]string, l *files.Limits) {
			l.MaxLinks = 2
		}},
		{"identities and paths over the byte limit", files.ErrLimit, func(_ map[files.ID]files.Path, _ map[string]files.ID, _ *[]string, l *files.Limits) {
			l.MaxBytes = 8
		}},
		{"an invalid limit record", files.ErrInvalid, func(_ map[files.ID]files.Path, _ map[string]files.ID, _ *[]string, l *files.Limits) {
			l.MaxLinks = 0
		}},
		// The number of files is counted before an assignment is read.
		{"too many files and an undeclared file", files.ErrLimit, func(_ map[files.ID]files.Path, a map[string]files.ID, _ *[]string, l *files.Limits) {
			l.MaxFiles = 2
			a["y"] = "absent"
		}},
		// An assignment is read before the paths are.
		{"an undeclared file and paths over the byte limit", files.ErrInvalid, func(_ map[files.ID]files.Path, a map[string]files.ID, _ *[]string, l *files.Limits) {
			l.MaxBytes = 8
			a["y"] = "absent"
		}},
	}
	for _, c := range refused {
		f, a, i, l := table(), assignment(), slices.Clone(items), files.DefaultLimits()
		c.change(f, a, &i, &l)
		inventory, err := files.Cover(f, a, i, l)
		other := files.ErrLimit
		if c.want == files.ErrLimit {
			other = files.ErrInvalid
		}
		if !errors.Is(err, c.want) || errors.Is(err, other) || len(inventory.Entries) != 0 {
			t.Errorf("%s: %v, inventory %v", c.name, err, inventory.Entries)
		}
	}
	// Exactly the limit is admitted.
	limits := files.DefaultLimits()
	limits.MaxFiles, limits.MaxLinks = 3, 3
	if _, err := files.Cover(table(), assignment(), items, limits); err != nil {
		t.Fatal("three files under a limit of three", err)
	}
}
