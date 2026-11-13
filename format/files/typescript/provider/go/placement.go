// Package format places analyzed TypeScript declarations into concrete source units.
package format

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"strings"

	f "github.com/Bitspark/bittype/contracts/format/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
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

// Organization is the placement of existing source. Its sites belong to the
// snapshot that travels in the same input, and it travels as a setting. It
// names no entry: no TypeScript reader has one. The formatter chooses the entry
// of the fresh snapshot, which is the file of least identity.
type Organization struct {
	Units []Placement
}

// OrganizationSetting is the key of the vocabulary without an entry. The
// vocabulary that named one was published as declaration-placement-1. A
// setting under that key is another setting and is refused, whatever it holds:
// a vocabulary does not change under its key.
func OrganizationSetting() semantics.Key { return n.Key("declaration-placement-2") }
func (o Organization) Setting() (semantic.Setting, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return semantic.Setting{}, err
	}
	return semantic.Setting{Key: OrganizationSetting(), Value: data}, nil
}

// GeneratedPlacement is the placement of the declarations of one sealed output
// into TypeScript files. Files gives each file identity its path. Definitions
// gives each declaration the file that holds it. A declaration is named by the
// handle that the output issued for it, so a placement belongs to one output
// and to no other run of the same plan.
//
// Both tables are keyed. A placement therefore states no order, cannot name a
// file twice and cannot place a declaration twice. It names no entry. The
// formatter owns the order: the files by their identities, and the
// declarations of one file by their sites in the sealed output. The entry of
// the fresh snapshot is the file of least identity. That entry is an address.
// It names no root of a TypeScript program and limits no check.
type GeneratedPlacement struct {
	Files       map[files.ID]files.Path
	Definitions map[g.DeclarationHandle]files.ID
}

var placements = f.NewGeneratedPlacementIssuer[GeneratedPlacement]()

// Issue returns the value that a generated-Files request carries. It owns a
// copy of both tables.
func (p GeneratedPlacement) Issue() f.GeneratedPlacement {
	return placements.Issue(GeneratedPlacement{Files: maps.Clone(p.Files), Definitions: maps.Clone(p.Definitions)})
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
	public   map[ast.NodeRef]bool
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
					// Another import form is TypeScript that this profile does not cover.
					return nil, fmt.Errorf("%w: only named type imports are supported", source.ErrUnsupported)
				}
			}
		default:
			return nil, source.ErrInvalid
		}
	}
	return out, nil
}
func readProgram(ctx context.Context, snapshot semantic.Snapshot, graph semantics.Semantics[s.Language]) (*program, error) {
	p, err := readDeclarations(snapshot)
	if err != nil {
		return nil, err
	}
	if err := p.readSemantics(ctx, graph); err != nil {
		return nil, err
	}
	return p, nil
}

// readDeclarations reads what the captured trees say alone: the files, their
// trees and the top-level declarations of each. A placement is judged against
// this much, before the graph is read.
func readDeclarations(snapshot semantic.Snapshot) (*program, error) {
	fileContext := snapshot.Files()
	if fileContext == nil {
		// A capture without files is well formed, and this profile reads none.
		return nil, fmt.Errorf("%w: TypeScript placement requires captured files", source.ErrUnsupported)
	}
	p := &program{public: map[ast.NodeRef]bool{}, trees: map[ast.TreeID]s.AST{}, tops: map[ast.NodeRef]bool{}, bindings: map[ast.NodeRef]symbol{}}
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
		for _, site := range sites {
			p.tops[site] = true
		}
	}
	return p, nil
}

// readSemantics reads what the graph says of the declarations: that every file
// is a module, what each module exports and where each name is bound. A script
// is refused here, as unsupported coverage.
func (p *program) readSemantics(ctx context.Context, graph semantics.Semantics[s.Language]) error {
	data := graph.Data()
	modules := map[ast.NodeRef]bool{}
	for _, e := range data.Entities {
		if e.Kind == n.ModuleKind() && len(e.Declarations) == 1 {
			modules[e.Declarations[0]] = true
		}
	}
	for _, tree := range p.trees {
		if !modules[ref(tree.ID, tree.Root)] {
			// A script is TypeScript that this profile does not cover.
			return fmt.Errorf("%w: placement requires TypeScript modules with file-local scope", source.ErrUnsupported)
		}
	}
	for _, fact := range data.Facts {
		if fact.Relation != n.NamedExportRelation() {
			continue
		}
		if fact.Subject.Entity == 0 {
			return source.ErrInvalid
		}
		owner := data.Entities[fact.Subject.Entity-1]
		term := data.Terms[fact.Term-1]
		if term.Kind != semantics.Reference || owner.Kind != n.ModuleKind() || len(owner.Declarations) != 1 {
			return source.ErrInvalid
		}
		declaration := data.Entities[term.Target-1]
		if len(declaration.Declarations) != 1 || declaration.Owner != fact.Subject.Entity || !p.tops[declaration.Declarations[0]] {
			return source.ErrInvalid
		}
		p.public[declaration.Declarations[0]] = true
	}
	for _, b := range data.Bindings {
		if err := ctx.Err(); err != nil {
			return err
		}
		if b.At.Role != n.NameRole() && b.At.Role != n.ImportRole() {
			continue
		}
		entity := data.Entities[b.Target-1]
		var dest symbol
		switch {
		case len(entity.Declarations) == 1:
			dest.site = entity.Declarations[0]
		case entity.Kind == n.PackageTypeKind() && entity.Owner != 0:
			owner := data.Entities[entity.Owner-1]
			if owner.Kind != n.PackageKind() {
				return source.ErrInvalid
			}
			dest.external = imports.Export{Module: owner.Name, Name: entity.Name}
		case entity.Kind == n.BuiltinKind():
			dest.builtin = entity.Name
		case len(entity.Declarations) > 1:
			// The compiler keeps a member and a type parameter of one name of an
			// interface in one symbol, so a use of the parameter has both as its
			// declarations. The file is valid TypeScript that this profile does
			// not tell apart, and a later slice can: unsupported coverage.
			return fmt.Errorf("%w: the compiler declares %s at %d captured sites", source.ErrUnsupported, entity.Name, len(entity.Declarations))
		default:
			// A target without a captured declaration that is no package type
			// and no builtin: the graph does not say what the name refers to.
			return fmt.Errorf("%w: ambiguous native binding %s", source.ErrInvalid, entity.Name)
		}
		if _, dup := p.bindings[b.At.Site]; dup {
			return source.ErrInvalid
		}
		p.bindings[b.At.Site] = dest
	}
	return nil
}

// decodeOrganization reads the placement of existing source from the settings
// of a request. It reports false when the request states none. The form of the
// setting is decided here, before the captured program is read: a setting that
// is malformed under its vocabulary is invalid input whatever the program is.
func decodeOrganization(settings []semantic.Setting) (Organization, bool, error) {
	if len(settings) == 0 {
		return Organization{}, false, nil
	}
	if len(settings) != 1 || settings[0].Key != OrganizationSetting() {
		return Organization{}, false, source.ErrInvalid
	}
	var o Organization
	// A setting contains one strict object.
	decoder := json.NewDecoder(bytes.NewReader(settings[0].Value))
	decoder.DisallowUnknownFields()
	if !json.Valid(settings[0].Value) || decoder.Decode(&o) != nil || repeatedMember(settings[0].Value) {
		return Organization{}, false, source.ErrInvalid
	}
	return o, true, nil
}

// admit decides the validity of a placement of existing source against the
// captured program. The lists of the setting become the two keyed tables that
// files.Cover reads, so the common rules have one owner: every declaration
// exactly once, every file declared, and an inventory that the files contract
// admits. A list can state what a keyed table cannot, a file identity or a
// declaration that is named twice, and that is refused here.
func (p *program) admit(o Organization, limits source.Limits) error {
	// The number of files is counted before a list is read. The count answers
	// with files.ErrLimit, or with files.ErrInvalid for an invalid limit record.
	if err := limits.Files.CheckFiles(len(o.Units)); err != nil {
		return source.WithCategory(err)
	}
	table := make(map[files.ID]files.Path, len(o.Units))
	assignment := map[ast.NodeRef]files.ID{}
	for _, u := range o.Units {
		if _, twice := table[u.File]; twice {
			// The files contract refuses an identity that two files share.
			return fmt.Errorf("%w: file identity %q is named twice: %w", source.ErrInvalid, u.File, files.ErrInvalid)
		}
		table[u.File] = u.Path
		for _, site := range u.Declarations {
			if _, twice := assignment[site]; twice {
				return fmt.Errorf("%w: declaration %v is placed twice", source.ErrInvalid, site)
			}
			assignment[site] = u.File
		}
	}
	tops := make([]ast.NodeRef, 0, len(p.tops))
	for site := range p.tops {
		tops = append(tops, site)
	}
	// Cover answers with files.ErrInvalid or files.ErrLimit.
	_, err := files.Cover(table, assignment, tops, limits.Files)
	return source.WithCategory(err)
}

// repeatedMember reports a JSON object that names one member twice. The decoder
// keeps the last value, so the first would be accepted and ignored. Names are
// compared as the decoder matches them, without regard to case. The caller has
// decoded the value into its type, which bounds the nesting.
func repeatedMember(data []byte) bool {
	decoder := json.NewDecoder(bytes.NewReader(data))
	repeated := false
	var value func()
	value = func() {
		token, err := decoder.Token()
		if err != nil {
			return
		}
		switch token {
		case json.Delim('{'):
			var names []string
			for decoder.More() {
				key, _ := decoder.Token()
				name, _ := key.(string)
				for _, seen := range names {
					repeated = repeated || strings.EqualFold(seen, name)
				}
				names = append(names, name)
				value()
			}
			decoder.Token()
		case json.Delim('['):
			for decoder.More() {
				value()
			}
			decoder.Token()
		}
	}
	value()
	return repeated
}
