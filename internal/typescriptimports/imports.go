// Package typescriptimports owns the named type-import vocabulary shared by
// TypeScript AST generation and Files formatting.
package typescriptimports

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/Bitspark/bittype/files/go"
	tsname "github.com/Bitspark/bittype/internal/typescript"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type Export struct{ Module, Name string }

var ErrLimit = errors.New("TypeScript import limit exceeded")

func SourcePath(p files.Path) bool {
	return files.ValidPath(p) && strings.HasSuffix(string(p), ".ts") && !strings.HasSuffix(string(p), ".d.ts") && !strings.ContainsAny(string(p), "?#%")
}

// ExcludedPath reports a valid path that TypeScript itself reads as a source
// file and that SourcePath excludes: a declaration file, another TypeScript
// extension, or a URL delimiter in a path with such an extension.
func ExcludedPath(p files.Path) bool {
	if !files.ValidPath(p) || SourcePath(p) {
		return false
	}
	for _, extension := range []string{".ts", ".mts", ".cts", ".tsx"} {
		if strings.HasSuffix(string(p), extension) {
			return true
		}
	}
	return false
}
func Relative(from, to files.Path) (string, error) {
	if !SourcePath(from) || !SourcePath(to) || from == to {
		return "", fmt.Errorf("invalid TypeScript import route %q to %q", from, to)
	}
	parent := strings.Split(string(from), "/")
	parent = parent[:len(parent)-1]
	dest := strings.Split(strings.TrimSuffix(string(to), ".ts")+".js", "/")
	i := 0
	for i < len(parent) && i < len(dest)-1 && parent[i] == dest[i] {
		i++
	}
	result := strings.Repeat("../", len(parent)-i) + strings.Join(dest[i:], "/")
	if i == len(parent) {
		result = "./" + result
	}
	return result, nil
}

func Attr(n ts.Node, name string) string {
	for _, a := range n.Attributes {
		if a.Name == name {
			return a.Text
		}
	}
	return ""
}
func SetName(n *ts.Node, name string) {
	for i := range n.Attributes {
		if n.Attributes[i].Name == ts.AttributeName {
			n.Attributes[i].Text = name
			return
		}
	}
}
func AddNode(tree *ts.AST, n ts.Node) ast.NodeID {
	tree.Nodes = append(tree.Nodes, n)
	return ast.NodeID(len(tree.Nodes))
}

// Emit binds actual occurrence sites to their import bindings and adds the
// declarations to this unit. All aliases are chosen after local names are known.
func Emit(ctx context.Context, tree *ts.AST, uses map[ast.NodeID]Export, maxNodes, maxBytes int) error {
	if len(uses) == 0 {
		return ctx.Err()
	}
	if len(uses) > maxNodes || len(tree.Nodes) > maxNodes {
		return ErrLimit
	}
	byExport := map[Export][]ast.NodeID{}
	occupied := map[string]bool{}
	for i, n := range tree.Nodes {
		id := ast.NodeID(i + 1)
		if _, isImport := uses[id]; isImport {
			continue
		}
		switch n.Kind {
		case ts.KindTypeAlias, ts.KindInterface, ts.KindTypeParameter, ts.KindTypeReference:
			occupied[Attr(n, ts.AttributeName)] = true
		case ts.KindImport:
			return fmt.Errorf("expected a unit with its former imports removed")
		}
	}
	for id, e := range uses {
		if id == 0 || int(id) > len(tree.Nodes) || tree.Nodes[id-1].Kind != ts.KindTypeReference || !tsname.TypeName(e.Name) || e.Module == "" {
			return fmt.Errorf("invalid TypeScript import occurrence")
		}
		byExport[e] = append(byExport[e], id)
	}
	exports := make([]Export, 0, len(byExport))
	for e := range byExport {
		exports = append(exports, e)
	}
	sort.Slice(exports, func(i, j int) bool {
		a, b := exports[i], exports[j]
		return a.Module < b.Module || a.Module == b.Module && a.Name < b.Name
	})
	var imports []ast.Child
	var previous string
	var module ast.NodeID
	suffix := map[string]int{}
	for _, e := range exports {
		if err := ctx.Err(); err != nil {
			return err
		}
		name := e.Name
		for occupied[name] {
			suffix[e.Name]++
			name = fmt.Sprintf("%s_%d", e.Name, suffix[e.Name])
		}
		occupied[name] = true
		attrs := []ast.Attribute{{Name: ts.AttributeName, Text: e.Name}}
		if name != e.Name {
			attrs = append(attrs, ast.Attribute{Name: ts.AttributeAlias, Text: name})
		}
		if len(tree.Nodes) >= maxNodes {
			return ErrLimit
		}
		binding := AddNode(tree, ts.Node{Kind: ts.KindImportBinding, Attributes: attrs})
		if e.Module != previous {
			if len(tree.Nodes) >= maxNodes {
				return ErrLimit
			}
			module = AddNode(tree, ts.Node{Kind: ts.KindImport, Attributes: []ast.Attribute{{Name: ts.AttributeModule, Text: e.Module}}})
			imports = append(imports, ast.Child{Role: ts.RoleImport, Node: module})
			previous = e.Module
		}
		tree.Nodes[module-1].Children = append(tree.Nodes[module-1].Children, ast.Child{Role: ts.RoleBinding, Node: binding})
		for _, site := range byExport[e] {
			node := &tree.Nodes[site-1]
			SetName(node, name)
			node.References = []ast.Reference{{Role: ts.RoleTarget, Target: ast.NodeRef{Tree: tree.ID, Node: binding}}}
		}
	}
	tree.Nodes[tree.Root-1].Children = append(imports, tree.Nodes[tree.Root-1].Children...)
	bytes := 0
	for _, n := range tree.Nodes {
		bytes += len(n.Kind)
		for _, a := range n.Attributes {
			bytes += len(a.Name) + len(a.Text)
		}
		if bytes > maxBytes {
			return ErrLimit
		}
	}
	return nil
}
