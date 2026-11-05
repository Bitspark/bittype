// Package files describes supplied text files. It performs no filesystem IO.
package files

import (
	"errors"
	"github.com/Bitspark/bittype/trees/ast/go"
	"slices"
	"strings"
	"unicode/utf8"
)

type ID string
type Path string
type File struct {
	ID   ID
	Path Path
	Text string
}

// Files is a supplied collection of file paths and text.
//
// bittype:dimension Files
type Files struct{ Entries []File }
type SyntaxLink struct {
	File   ID
	Tree   ast.TreeID
	Source ast.SourceID
}
type Context struct {
	Files Files
	Links []SyntaxLink
}
type Limits struct{ MaxFiles, MaxBytes, MaxLinks int }

func DefaultLimits() Limits  { return Limits{256, 16 << 20, 256} }
func (l Limits) Valid() bool { return l.MaxFiles > 0 && l.MaxBytes > 0 && l.MaxLinks > 0 }

var ErrInvalid = errors.New("invalid Files context")
var ErrLimit = errors.New("Files resource limit")

func ValidPath(path Path) bool {
	s := string(path)
	if s == "" || !utf8.ValidString(s) || strings.ContainsAny(s, "\\\x00:") {
		return false
	}
	for _, part := range strings.Split(s, "/") {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}
func (f Files) Clone() Files { return Files{Entries: slices.Clone(f.Entries)} }
func (f Files) Capture(l Limits) (Files, error) {
	if !l.Valid() || len(f.Entries) == 0 {
		return Files{}, ErrInvalid
	}
	if len(f.Entries) > l.MaxFiles {
		return Files{}, ErrLimit
	}
	ids, paths := map[ID]bool{}, map[Path]bool{}
	bytes := l.MaxBytes
	for _, entry := range f.Entries {
		for _, s := range []string{string(entry.ID), string(entry.Path), entry.Text} {
			if len(s) > bytes {
				return Files{}, ErrLimit
			}
			bytes -= len(s)
		}
		if entry.ID == "" || !utf8.ValidString(string(entry.ID)) || strings.ContainsRune(string(entry.ID), 0) || !ValidPath(entry.Path) || ids[entry.ID] || paths[entry.Path] {
			return Files{}, ErrInvalid
		}
		ids[entry.ID], paths[entry.Path] = true, true
	}
	for path := range paths {
		parts := strings.Split(string(path), "/")
		for i := 1; i < len(parts); i++ {
			if paths[Path(strings.Join(parts[:i], "/"))] {
				return Files{}, ErrInvalid
			}
		}
	}
	return f.Clone(), nil
}
func (c Context) Clone() Context {
	return Context{Files: c.Files.Clone(), Links: slices.Clone(c.Links)}
}
func (c Context) Capture(l Limits, trees map[ast.TreeID]ast.NodeID) (Context, error) {
	owned, err := c.Files.Capture(l)
	if err != nil {
		return Context{}, err
	}
	if len(c.Links) > l.MaxLinks {
		return Context{}, ErrLimit
	}
	if len(c.Links) != len(owned.Entries) || len(c.Links) != len(trees) {
		return Context{}, ErrInvalid
	}
	ids, sources, seen := map[ID]bool{}, map[ast.SourceID]bool{}, map[ast.TreeID]bool{}
	for _, entry := range owned.Entries {
		ids[entry.ID] = true
	}
	bytes := l.MaxBytes
	for _, link := range c.Links {
		for _, s := range []string{string(link.File), string(link.Tree), string(link.Source)} {
			if len(s) > bytes {
				return Context{}, ErrLimit
			}
			bytes -= len(s)
		}
		if !ids[link.File] || link.Source == "" || sources[link.Source] || seen[link.Tree] || trees[link.Tree] == 0 {
			return Context{}, ErrInvalid
		}
		ids[link.File] = false
		sources[link.Source], seen[link.Tree] = true, true
	}
	return Context{Files: owned, Links: slices.Clone(c.Links)}, nil
}
