package horos_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"testing"

	horos "github.com/Bitspark/bittype/internal/horos"
)

// This file spells every name of tags.go a second time, on purpose. The
// parser writes each of them and the formatter reads it, both from tags.go, so
// they follow a changed value together. A canonical value that was written
// before holds the old one, and the horos library reads the old one. Changing
// a value is a decision: make it here as well.
var canonicalNames = []struct {
	name, got, want string
}{
	{"FieldTag", horos.FieldTag, "horos/field"},
	{"RecordOpen", horos.RecordOpen, "open"},
	{"RecordClosed", horos.RecordClosed, "closed"},
	{"FieldRequired", horos.FieldRequired, "req"},
	{"FieldOptional", horos.FieldOptional, "opt"},
}

func TestCanonicalNamesAreUnchanged(t *testing.T) {
	for _, c := range canonicalNames {
		if c.got != c.want {
			t.Errorf("%s = %q, canonical as %q", c.name, c.got, c.want)
		}
	}
}

// Every exported name of tags.go has a pin above. A new one fails here until
// it gets one.
func TestEveryCanonicalNameIsPinned(t *testing.T) {
	var pinned []string
	for _, c := range canonicalNames {
		pinned = append(pinned, c.name)
	}
	file, err := parser.ParseFile(token.NewFileSet(), "tags.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	var exported []string
	for _, d := range file.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			if d.Recv == nil && d.Name.IsExported() {
				exported = append(exported, d.Name.Name)
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch spec := spec.(type) {
				case *ast.ValueSpec:
					for _, name := range spec.Names {
						if name.IsExported() {
							exported = append(exported, name.Name)
						}
					}
				case *ast.TypeSpec:
					if spec.Name.IsExported() {
						exported = append(exported, spec.Name.Name)
					}
				}
			}
		}
	}
	if len(exported) == 0 {
		t.Fatal("found no exported name; the test must run in the package directory")
	}
	for _, name := range exported {
		if !slices.Contains(pinned, name) {
			t.Errorf("exported %s has no pin in this file", name)
		}
	}
	for _, name := range pinned {
		if !slices.Contains(exported, name) {
			t.Errorf("this file pins %s, which tags.go does not declare", name)
		}
	}
}
