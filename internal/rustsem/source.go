// Package rustsem owns the private, pinned Rust compiler bridge.
package rustsem

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	parser "github.com/Bitspark/bittype/parse/syntax/rust/provider/go"
	n "github.com/Bitspark/bittype/semantics/rust/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type Node struct {
	ID         ast.NodeID
	Kind       s.Kind
	Start      int
	End        int
	Attributes []ast.Attribute
	Children   []ast.Child
	References []ast.Reference
}
type Unit struct {
	Tree  ast.TreeID
	Text  string
	Nodes []Node
}
type input struct {
	Version int
	Units   []Unit
}

func prepare(ctx context.Context, trees []s.AST, source *files.Context) ([]Unit, error) {
	if source == nil {
		return nil, nil
	}
	if len(trees) > 16 {
		return nil, semantic.ErrLimit
	}
	var units []Unit
	totalNodes, totalBytes := 0, 0
	for _, tree := range trees {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		var link files.SyntaxLink
		var file files.File
		for _, candidate := range source.Links {
			if candidate.Tree == tree.ID {
				link = candidate
				break
			}
		}
		for _, candidate := range source.Files.Entries {
			if candidate.ID == link.File {
				file = candidate
				break
			}
		}
		if link.Tree == "" || file.ID == "" || !utf8.ValidString(file.Text) || !strings.HasSuffix(string(file.Path), ".rs") {
			return nil, semantic.ErrInvalidRequest
		}
		totalNodes += len(tree.Nodes)
		totalBytes += len(file.Text)
		if totalNodes > 4096 || totalBytes > 1<<20 {
			return nil, semantic.ErrLimit
		}
		parsed, err := parser.Parse(ctx, parse.Input{Text: file.Text, SourceID: link.Source, TreeID: link.Tree})
		if err != nil {
			return nil, fmt.Errorf("%w: captured source/AST mismatch: %v", semantic.ErrInvalidRequest, err)
		}
		expected := parsed.ASTView()
		if expected.Root != tree.Root || len(expected.Nodes) != len(tree.Nodes) {
			return nil, semantic.ErrInvalidRequest
		}
		for i, node := range expected.Nodes {
			actual := tree.Nodes[i]
			if node.Kind != actual.Kind || !reflect.DeepEqual(node.Attributes, actual.Attributes) || !reflect.DeepEqual(node.Children, actual.Children) {
				return nil, fmt.Errorf("%w: captured source/AST mismatch at %s:%d", semantic.ErrInvalidRequest, tree.ID, i+1)
			}
		}
		// rustc normalizes BOM/CRLF. Preserve correspondence to the original
		// captured bytes while passing the same normalized source to its driver.
		offsets := make([]int, len(file.Text)+1)
		normalized := make([]byte, 0, len(file.Text))
		for i := 0; i < len(file.Text); i++ {
			offsets[i] = len(normalized)
			if i < 3 && strings.HasPrefix(file.Text, "\xef\xbb\xbf") {
				continue
			}
			if file.Text[i] == '\r' && i+1 < len(file.Text) && file.Text[i+1] == '\n' {
				continue
			}
			normalized = append(normalized, file.Text[i])
		}
		offsets[len(file.Text)] = len(normalized)
		unit := Unit{Tree: tree.ID, Text: string(normalized)}
		for i, node := range expected.Nodes {
			if len(node.Origin.Spans) != 1 {
				return nil, semantic.ErrInvalidRequest
			}
			span := node.Origin.Spans[0]
			unit.Nodes = append(unit.Nodes, Node{ID: ast.NodeID(i + 1), Kind: node.Kind, Start: offsets[span.Start], End: offsets[span.End], Attributes: node.Attributes, Children: node.Children, References: tree.Nodes[i].References})
		}
		units = append(units, unit)
	}
	return units, nil
}
func diagnostic(status, message string) semantic.Diagnostic {
	return semantic.Diagnostic{Code: n.Key(status), Subject: n.Coverage(), Message: message}
}
