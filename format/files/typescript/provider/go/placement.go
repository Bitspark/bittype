// Package format places analyzed TypeScript declarations into concrete source units.
package format

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	imports "github.com/Bitspark/bittype/internal/typescriptimports"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	"github.com/Bitspark/bittype/source/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Placement assigns declaration sites in one admitted input capture to a file.
type Placement struct {
	File         files.ID
	Path         files.Path
	Declarations []ast.NodeRef
}
type Organization struct {
	Entry files.ID
	Units []Placement
}

func OrganizationSetting() semantics.Key { return n.Key("declaration-placement-1") }
func (o Organization) Setting() (semantic.Setting, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return semantic.Setting{}, err
	}
	return semantic.Setting{Key: OrganizationSetting(), Value: data}, nil
}

type symbol struct {
	site     ast.NodeRef
	external imports.Export
	// A builtin compiler symbol has no application-owned declaration.
	builtin string
}
type program struct {
	trees    map[ast.TreeID]s.AST
	tops     map[ast.NodeRef]bool
	bindings map[ast.NodeRef]symbol
	original Organization
}

func ref(tree ast.TreeID, node ast.NodeID) ast.NodeRef { return ast.NodeRef{Tree: tree, Node: node} }
func declarations(t s.AST) ([]ast.NodeRef, error) {
	var out []ast.NodeRef
	for _, child := range t.Nodes[t.Root-1].Children {
		switch child.Role {
		case s.RoleDeclaration:
			out = append(out, ref(t.ID, child.Node))
		case s.RoleImport:
			for _, b := range t.Nodes[child.Node-1].Children {
				if t.Nodes[b.Node-1].Kind != s.KindImportBinding {
					return nil, fmt.Errorf("only named type imports are supported")
				}
			}
		default:
			return nil, source.ErrInvalid
		}
	}
	return out, nil
}
func readProgram(ctx context.Context, snapshot semantic.Snapshot, graph semantics.Semantics[s.Language]) (*program, error) {
	fileContext := snapshot.Files()
	if fileContext == nil {
		return nil, fmt.Errorf("%w: TypeScript placement requires captured files", source.ErrInvalid)
	}
	p := &program{trees: map[ast.TreeID]s.AST{}, tops: map[ast.NodeRef]bool{}, bindings: map[ast.NodeRef]symbol{}}
	byFile := map[files.ID]ast.TreeID{}
	for _, link := range fileContext.Links {
		byFile[link.File] = link.Tree
	}
	for _, envelope := range snapshot.Inputs() {
		tree, err := representation.Unpack[s.Language](snapshot.Language(), envelope)
		if err != nil {
			return nil, err
		}
		p.trees[tree.ID] = tree
	}
	for _, file := range fileContext.Files.Entries {
		tree, ok := p.trees[byFile[file.ID]]
		if !ok {
			return nil, source.ErrInvalid
		}
		sites, err := declarations(tree)
		if err != nil {
			return nil, err
		}
		p.original.Units = append(p.original.Units, Placement{file.ID, file.Path, sites})
		if snapshot.Entry() == ref(tree.ID, tree.Root) {
			p.original.Entry = file.ID
		}
		for _, site := range sites {
			p.tops[site] = true
		}
	}
	data := graph.Data()
	modules := map[ast.NodeRef]bool{}
	for _, e := range data.Entities {
		if e.Kind == n.Key("module") && len(e.Declarations) == 1 {
			modules[e.Declarations[0]] = true
		}
	}
	for _, tree := range p.trees {
		if !modules[ref(tree.ID, tree.Root)] {
			return nil, fmt.Errorf("%w: placement requires TypeScript modules with file-local scope", source.ErrInvalid)
		}
	}
	for _, b := range data.Bindings {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if b.At.Role != n.Key("name") && b.At.Role != n.Key("import") {
			continue
		}
		entity := data.Entities[b.Target-1]
		var dest symbol
		switch {
		case len(entity.Declarations) == 1:
			dest.site = entity.Declarations[0]
		case entity.Kind == n.Key("package-type") && entity.Owner != 0:
			owner := data.Entities[entity.Owner-1]
			if owner.Kind != n.Key("package") {
				return nil, source.ErrInvalid
			}
			dest.external = imports.Export{Module: owner.Name, Name: entity.Name}
		case entity.Kind == n.Key("builtin"):
			dest.builtin = entity.Name
		default:
			return nil, fmt.Errorf("%w: ambiguous native binding %s", source.ErrInvalid, entity.Name)
		}
		if _, dup := p.bindings[b.At.Site]; dup {
			return nil, source.ErrInvalid
		}
		p.bindings[b.At.Site] = dest
	}
	return p, nil
}
func (p *program) organization(settings []semantic.Setting, limits source.Limits) (Organization, error) {
	o := p.original
	if len(settings) > 0 {
		if len(settings) != 1 || settings[0].Key != OrganizationSetting() {
			return Organization{}, source.ErrInvalid
		}
		o = Organization{}
		// A setting contains one strict object, including the explicit entry.
		decoder := json.NewDecoder(bytes.NewReader(settings[0].Value))
		decoder.DisallowUnknownFields()
		if !json.Valid(settings[0].Value) || decoder.Decode(&o) != nil {
			return Organization{}, source.ErrInvalid
		}
	}
	if len(o.Units) > limits.Files.MaxFiles {
		return Organization{}, source.ErrLimit
	}
	inventory := files.Files{}
	assigned := map[ast.NodeRef]bool{}
	entry := false
	for _, u := range o.Units {
		if !imports.SourcePath(u.Path) || len(u.Declarations) == 0 {
			return Organization{}, source.ErrInvalid
		}
		inventory.Entries = append(inventory.Entries, files.File{ID: u.File, Path: u.Path})
		if u.File == o.Entry {
			entry = true
		}
		for _, site := range u.Declarations {
			if !p.tops[site] || assigned[site] {
				return Organization{}, fmt.Errorf("%w: duplicate or unknown declaration %v", source.ErrInvalid, site)
			}
			assigned[site] = true
		}
	}
	if !entry || len(assigned) != len(p.tops) {
		return Organization{}, fmt.Errorf("%w: incomplete placement or missing entry", source.ErrInvalid)
	}
	if _, err := inventory.Capture(limits.Files); err != nil {
		return Organization{}, err
	}
	return o, nil
}
