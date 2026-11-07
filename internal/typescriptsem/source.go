// Package typescriptsem owns the private TypeScript compiler bridge.
package typescriptsem

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	parser "github.com/Bitspark/bittype/parse/syntax/typescript/provider/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
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
	Path  files.Path
	Text  string
	Nodes []Node
}
type Query struct {
	Source uint32
	Target uint32
}
type input struct {
	NamedImports bool
	Packages     []n.PackageTypes
	Version      int
	Units        []Unit
	Query        *Query
}

// The syntax checker admits explicit false flags. The parser omits a flag
// when its token is absent. All other attributes must match exactly.
func sameAttributes(actual, parsed []ast.Attribute) bool {
	count := 0
	for _, a := range actual {
		if a.Text == "false" && (a.Name == s.AttributeOptional || a.Name == s.AttributeReadonly || a.Name == s.AttributeExported || a.Name == s.AttributeRest) {
			continue
		}
		if !slices.Contains(parsed, a) {
			return false
		}
		count++
	}
	return count == len(parsed)
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
		if link.Tree == "" || file.ID == "" {
			return nil, semantic.ErrInvalidRequest
		}
		if !utf8.ValidString(file.Text) {
			return nil, semantic.ErrInvalidRequest
		}
		if !strings.HasSuffix(string(file.Path), ".ts") || strings.HasSuffix(string(file.Path), ".d.ts") {
			return nil, fmt.Errorf("%w: captured files must be .ts source units", semantic.ErrInvalidRequest)
		}
		totalNodes += len(tree.Nodes)
		totalBytes += len(file.Text)
		if totalNodes > 4096 || totalBytes > 1<<20 {
			return nil, semantic.ErrLimit
		}
		// Reparse the supplied original text to establish correspondence. No
		// formatter or printer participates, and caller origins are not trusted.
		parsed, err := parser.Parse(ctx, parse.Input{Text: file.Text, SourceID: link.Source, TreeID: link.Tree})
		if err != nil {
			return nil, fmt.Errorf("%w: captured source/AST mismatch: %v", semantic.ErrInvalidRequest, err)
		}
		expected := parsed.ASTView()
		// Compare by containment position. Parser and generators can use different
		// node numbering and attribute order for exactly the same syntax.
		if len(expected.Nodes) != len(tree.Nodes) {
			return nil, semantic.ErrInvalidRequest
		}
		parsedFor := make(map[ast.NodeID]ast.NodeID, len(tree.Nodes))
		var pair func(ast.NodeID, ast.NodeID) error
		pair = func(actualID, parsedID ast.NodeID) error {
			actual, parsed := tree.Nodes[actualID-1], expected.Nodes[parsedID-1]
			if actual.Kind != parsed.Kind || !sameAttributes(actual.Attributes, parsed.Attributes) || len(actual.Children) != len(parsed.Children) {
				return fmt.Errorf("%w: node %d (%s) differs from parsed %s in attributes %v / %v or children", semantic.ErrInvalidRequest, actualID, actual.Kind, parsed.Kind, actual.Attributes, parsed.Attributes)
			}
			parsedFor[actualID] = parsedID
			for i, child := range actual.Children {
				if child.Role != parsed.Children[i].Role {
					return semantic.ErrInvalidRequest
				}
				if err := pair(child.Node, parsed.Children[i].Node); err != nil {
					return err
				}
			}
			return nil
		}
		if err := pair(tree.Root, expected.Root); err != nil {
			return nil, fmt.Errorf("%w: captured source/AST mismatch in %s", err, tree.ID)
		}
		// TypeScript offsets use UTF-16; Bittype source spans use UTF-8 bytes.
		offsets := make([]int, len(file.Text)+1)
		count := 0
		for offset, r := range file.Text {
			offsets[offset] = count
			count++
			if r > 0xffff {
				count++
			}
			offsets[offset+utf8.RuneLen(r)] = count
		}
		unit := Unit{Tree: tree.ID, Path: file.Path, Text: file.Text}
		for i, node := range tree.Nodes {
			spans := expected.Nodes[parsedFor[ast.NodeID(i+1)]-1].Origin.Spans
			if len(spans) != 1 {
				return nil, semantic.ErrInvalidRequest
			}
			span := spans[0]
			unit.Nodes = append(unit.Nodes, Node{ID: ast.NodeID(i + 1), Kind: node.Kind, Start: offsets[span.Start], End: offsets[span.End], Attributes: node.Attributes, Children: node.Children, References: tree.Nodes[i].References})
		}
		units = append(units, unit)
	}
	return units, nil
}
