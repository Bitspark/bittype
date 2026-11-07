package format

import (
	"context"
	"errors"
	"fmt"
	"slices"

	fm "github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/files/go"
	imports "github.com/Bitspark/bittype/internal/typescriptimports"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	"github.com/Bitspark/bittype/source/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type placed struct {
	placement Placement
	tree      s.AST
	expected  map[ast.NodeID]symbol
}

type result struct {
	entry     files.ID
	units     []placed
	moved     map[ast.NodeRef]ast.NodeRef
	relations []source.Relation
}

func materialize(ctx context.Context, input fm.FilesInput[s.Language]) (*result, error) {
	program, err := readProgram(ctx, input.Snapshot, input.Abstract)
	if err != nil {
		return nil, err
	}
	organization, err := program.organization(input.Configuration, input.Limits.Source)
	if err != nil {
		return nil, err
	}
	layout := &result{entry: organization.Entry, moved: map[ast.NodeRef]ast.NodeRef{}}
	filesByTree := map[ast.TreeID]Placement{}
	for _, placement := range organization.Units {
		tree := s.AST{ID: ast.TreeID("typescript/files/" + string(placement.File))}
		layout.units = append(layout.units, placed{placement: placement, tree: tree, expected: map[ast.NodeID]symbol{}})
		filesByTree[tree.ID] = placement
	}
	if err := layout.copyDeclarations(ctx, program, input.Limits.Source); err != nil {
		return nil, err
	}
	nodesLeft := input.Limits.Source.MaxNodes
	for _, unit := range layout.units {
		nodesLeft -= len(unit.tree.Nodes)
	}
	for i := range layout.units {
		unit := &layout.units[i]
		before := len(unit.tree.Nodes)
		if err := layout.bindUnit(ctx, program, unit, filesByTree, before+nodesLeft, input.Limits.Source.MaxBytes); err != nil {
			return nil, err
		}
		nodesLeft -= len(unit.tree.Nodes) - before
	}
	return layout, nil
}

// Copy declarations without their former imports. Retain origins and record
// every original site's new location before binding references across files.
func (layout *result) copyDeclarations(ctx context.Context, program *program, limits source.Limits) error {
	nodes, work := limits.MaxNodes, limits.MaxWork
	var copyNode func(*placed, ast.NodeRef, int) (ast.NodeID, error)
	copyNode = func(unit *placed, old ast.NodeRef, depth int) (ast.NodeID, error) {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		nodes--
		work--
		if nodes < 0 || work < 0 || depth > limits.MaxDepth {
			return 0, source.ErrLimit
		}
		node := program.trees[old.Tree].Nodes[old.Node-1]
		children := node.Children
		node.Children = nil
		node.References = nil
		node.Attributes = slices.Clone(node.Attributes)
		id := imports.AddNode(&unit.tree, node)
		layout.moved[old] = ref(unit.tree.ID, id)
		if node.Kind == s.KindTypeReference {
			dest, ok := program.bindings[old]
			if !ok {
				return 0, fmt.Errorf("%w: type reference lacks native binding at %v", source.ErrInvalid, old)
			}
			unit.expected[id] = dest
		}
		for _, c := range children {
			child, err := copyNode(unit, ref(old.Tree, c.Node), depth+1)
			if err != nil {
				return 0, err
			}
			unit.tree.Nodes[id-1].Children = append(unit.tree.Nodes[id-1].Children, ast.Child{Role: c.Role, Node: child})
		}
		return id, nil
	}
	for i := range layout.units {
		unit := &layout.units[i]
		nodes--
		if nodes < 0 {
			return source.ErrLimit
		}
		unit.tree.Root = imports.AddNode(&unit.tree, s.Node{Kind: s.KindSourceFile})
		names := map[string]bool{}
		for _, site := range unit.placement.Declarations {
			name := imports.Attr(program.trees[site.Tree].Nodes[site.Node-1], s.AttributeName)
			if names[name] {
				return fmt.Errorf("%w: duplicate declaration %s in %s", source.ErrInvalid, name, unit.placement.Path)
			}
			names[name] = true
			id, err := copyNode(unit, site, 1)
			if err != nil {
				return err
			}
			unit.tree.Nodes[unit.tree.Root-1].Children = append(unit.tree.Nodes[unit.tree.Root-1].Children, ast.Child{Role: s.RoleDeclaration, Node: id})
		}
	}
	return nil
}

func (layout *result) bindUnit(ctx context.Context, program *program, unit *placed, filesByTree map[ast.TreeID]Placement, maxNodes, maxBytes int) error {
	uses := map[ast.NodeID]imports.Export{}
	locals := map[string]files.ID{}
	for id, dest := range unit.expected {
		if err := ctx.Err(); err != nil {
			return err
		}
		if dest.external.Module != "" {
			uses[id] = dest.external
			continue
		}
		if dest.builtin != "" {
			continue
		}
		moved, ok := layout.moved[dest.site]
		if !ok {
			return source.ErrInvalid
		}
		original := program.trees[dest.site.Tree].Nodes[dest.site.Node-1]
		name := imports.Attr(original, s.AttributeName)
		if moved.Tree == unit.tree.ID {
			node := &unit.tree.Nodes[id-1]
			imports.SetName(node, name)
			node.References = []ast.Reference{{Role: s.RoleTarget, Target: moved}}
			continue
		}
		if !program.tops[dest.site] || imports.Attr(original, s.AttributeExported) != "true" {
			return fmt.Errorf("%w: cross-unit type %s must be exported", source.ErrInvalid, name)
		}
		target := filesByTree[moved.Tree]
		route, err := imports.Relative(unit.placement.Path, target.Path)
		if err != nil {
			return err
		}
		uses[id] = imports.Export{Module: route, Name: name}
		locals[route] = target.File
	}
	if err := imports.Emit(ctx, &unit.tree, uses, maxNodes, maxBytes); err != nil {
		if errors.Is(err, imports.ErrLimit) {
			return source.ErrLimit
		}
		return err
	}
	for _, child := range unit.tree.Nodes[unit.tree.Root-1].Children {
		if child.Role != s.RoleImport {
			continue
		}
		spec := imports.Attr(unit.tree.Nodes[child.Node-1], s.AttributeModule)
		if to, ok := locals[spec]; ok {
			site := ref(unit.tree.ID, child.Node)
			layout.relations = append(layout.relations, source.Relation{Kind: n.Key("type-import"), From: unit.placement.File, To: to, At: &site})
		}
	}
	return nil
}
