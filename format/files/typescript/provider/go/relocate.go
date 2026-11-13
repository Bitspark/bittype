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
	public    map[ast.NodeRef]bool
}

func materialize(ctx context.Context, input fm.FilesInput[s.Language]) (*result, error) {
	// The form of the setting is decided before the captured program is read.
	organization, stated, err := decodeOrganization(input.Configuration)
	if err != nil {
		return nil, err
	}
	program, err := readDeclarations(input.Snapshot)
	if err != nil {
		return nil, err
	}
	if !stated {
		organization = program.original
	}
	// The placement is judged against the top-level declarations before the
	// graph is read: its files, its items and its paths need nothing else, so
	// a script does not hide a placement that no slice admits.
	if err := program.admit(organization, input.Limits.Source); err != nil {
		return nil, err
	}
	if _, err := sourcePaths(organization.Units); err != nil {
		return nil, err
	}
	if err := program.readSemantics(ctx, input.Abstract); err != nil {
		return nil, err
	}
	return placeProgram(ctx, program, organization.Units, input.Limits.Source)
}

// sourcePaths judges the paths of a placement. A path that is no TypeScript
// path is invalid input, and it is the error. A path that TypeScript reads and
// this profile excludes is unsupported coverage, and the first such path is
// the first result: the caller reports it after everything invalid.
func sourcePaths(units []Placement) (uncovered error, invalid error) {
	for _, u := range units {
		switch {
		case imports.SourcePath(u.Path):
		case !imports.ExcludedPath(u.Path):
			// A path that is no TypeScript path is invalid input.
			return nil, fmt.Errorf("%w: %q is no TypeScript source path", source.ErrInvalid, u.Path)
		case uncovered == nil:
			uncovered = fmt.Errorf("%w: path %q", source.ErrUnsupported, u.Path)
		}
	}
	return uncovered, nil
}

// placeProgram realizes a placement whose common rules are admitted: every
// declaration has one declared file and the inventory of the files is valid.
// units are in the order that the formatter realizes. Validity is decided
// before coverage, and both before anything is constructed:
//
//   - a path that is no TypeScript path, two declarations of one name in one
//     file and a use of an unexposed declaration of another file are invalid
//     input;
//   - a path that TypeScript reads and this profile excludes, a type parameter
//     that would capture a use of a declaration of its file, and a file
//     without a declaration, are unsupported coverage.
//
// A placement with a defect of both kinds is invalid input.
func placeProgram(ctx context.Context, program *program, units []Placement, limits source.Limits) (*result, error) {
	uncovered, err := sourcePaths(units)
	if err != nil {
		return nil, err
	}
	captured, err := program.relocatable(ctx, units)
	if err != nil {
		return nil, err
	}
	if uncovered == nil {
		uncovered = captured
	}
	for _, u := range units {
		if len(u.Declarations) == 0 && uncovered == nil {
			uncovered = fmt.Errorf("%w: file %q has no declaration", source.ErrUnsupported, u.File)
		}
	}
	if uncovered != nil {
		return nil, uncovered
	}
	// The entry of the fresh snapshot is the file of least identity. No
	// placement names one, and the rule reads no order of the input. There is a
	// first unit: the admission refuses a placement without a file.
	layout := &result{public: map[ast.NodeRef]bool{}, entry: units[0].File, moved: map[ast.NodeRef]ast.NodeRef{}}
	filesByTree := map[ast.TreeID]Placement{}
	for _, placement := range units {
		if placement.File < layout.entry {
			layout.entry = placement.File
		}
		tree := s.AST{ID: ast.TreeID("typescript/files/" + string(placement.File))}
		layout.units = append(layout.units, placed{placement: placement, tree: tree, expected: map[ast.NodeID]symbol{}})
		filesByTree[tree.ID] = placement
	}
	if err := layout.copyDeclarations(ctx, program, limits); err != nil {
		return nil, err
	}
	nodesLeft := limits.MaxNodes
	for _, unit := range layout.units {
		nodesLeft -= len(unit.tree.Nodes)
	}
	for i := range layout.units {
		unit := &layout.units[i]
		before := len(unit.tree.Nodes)
		if err := layout.bindUnit(ctx, program, unit, filesByTree, before+nodesLeft, limits.MaxBytes); err != nil {
			return nil, err
		}
		nodesLeft -= len(unit.tree.Nodes) - before
	}
	return layout, nil
}

// relocatable decides what only the declarations can decide about a placement,
// before a tree is constructed. One file cannot hold two declarations of one
// name. A declaration that another file uses has to be exposed: Files never
// exports a declaration without authorization. Both are invalid input. Every
// placed site is a top-level declaration of the program: the admission of an
// existing placement and the reading of a generated one established that.
//
// The first result is a refusal of coverage, which the caller reports after
// everything invalid: a use of a declaration of the same file that a type
// parameter would capture. A use is written by the name of its declaration,
// and inside a declaration that has a type parameter of that name the compiler
// binds the name to the parameter. A later slice may rename the parameter
// where that is permitted, so the placement is unsupported coverage. Placing
// the declaration in another file is admitted today: the import gets an alias.
// The check reads every type parameter below the declaration that holds the
// use, which is exact for the parameters of an alias and of an interface, the
// only binders of this profile.
func (program *program) relocatable(ctx context.Context, units []Placement) (captured error, err error) {
	type located struct {
		site    ast.NodeRef
		unit    int
		binders map[string]bool
	}
	home := map[ast.NodeRef]int{}
	var references []located
	for i, u := range units {
		names := map[string]bool{}
		for _, site := range u.Declarations {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			tree := program.trees[site.Tree]
			name := imports.Attr(tree.Nodes[site.Node-1], s.AttributeName)
			if names[name] {
				return nil, fmt.Errorf("%w: duplicate declaration %s in %s", source.ErrInvalid, name, u.Path)
			}
			names[name] = true
			binders := map[string]bool{}
			// The trees are captured, so they hold no cycle.
			for stack := []ast.NodeID{site.Node}; len(stack) > 0; {
				id := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				node := tree.Nodes[id-1]
				home[ref(site.Tree, id)] = i
				if node.Kind == s.KindTypeParameter {
					binders[imports.Attr(node, s.AttributeName)] = true
				}
				if node.Kind == s.KindTypeReference {
					references = append(references, located{ref(site.Tree, id), i, binders})
				}
				for _, child := range node.Children {
					stack = append(stack, child.Node)
				}
			}
		}
	}
	for _, use := range references {
		dest, bound := program.bindings[use.site]
		if !bound {
			return nil, fmt.Errorf("%w: type reference lacks native binding at %v", source.ErrInvalid, use.site)
		}
		if dest.external.Module != "" || dest.builtin != "" {
			continue
		}
		at, placed := home[dest.site]
		if !placed {
			return nil, fmt.Errorf("%w: the target of the type reference at %v is not placed", source.ErrInvalid, use.site)
		}
		name := imports.Attr(program.trees[dest.site.Tree].Nodes[dest.site.Node-1], s.AttributeName)
		if at != use.unit && (!program.tops[dest.site] || !program.public[dest.site]) {
			return nil, fmt.Errorf("%w: cross-unit type %s must be exported", source.ErrInvalid, name)
		}
		if at == use.unit && program.tops[dest.site] && use.binders[name] && captured == nil {
			captured = fmt.Errorf("%w: a type parameter %s in %s would capture a use of the declaration %s of that file", source.ErrUnsupported, name, units[at].Path, name)
		}
	}
	return captured, nil
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
		if program.tops[old] && program.public[old] && imports.Attr(node, s.AttributeExported) != "true" {
			node.Attributes = append(node.Attributes, ast.Attribute{Name: s.AttributeExported, Text: "true"})
		}
		id := imports.AddNode(&unit.tree, node)
		layout.moved[old] = ref(unit.tree.ID, id)
		if program.tops[old] {
			layout.public[ref(unit.tree.ID, id)] = program.public[old]
		}
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
		for _, site := range unit.placement.Declarations {
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
