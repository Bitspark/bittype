package format

import (
	"context"
	"fmt"
	f "github.com/Bitspark/bittype/contracts/format/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/files/go"
	imports "github.com/Bitspark/bittype/internal/typescriptimports"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	"github.com/Bitspark/bittype/source/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	tssymbol "github.com/Bitspark/bittype/transform/syntax/typescript/symbol/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type GeneratedProvider struct {
	s.Language
	stages Stages
}

// NewGenerated constructs Files formatting for generated TypeScript declaration units.
func NewGenerated(stages Stages) (f.GeneratedFilesFormatter[s.Language], error) {
	p, err := stages.capture()
	if err != nil {
		return f.GeneratedFilesFormatter[s.Language]{}, err
	}
	contract, err := p.GeneratedContract()
	if err != nil {
		return f.GeneratedFilesFormatter[s.Language]{}, err
	}
	return f.NewGeneratedFilesFormatter(p.Language, contract, GeneratedProvider{stages: p})
}

// GeneratedContract returns what TypeScript supplies to the generated-Files
// contract beside its provider: the Syntax printer and the preparation step.
// NewGenerated pairs it with the TypeScript provider. A test pairs it with a
// provider of its own, to hand the TypeScript guard a publication that the
// TypeScript layout does not produce.
//
// The step reads a GeneratedPlacement, has its two tables admitted and then
// decides what TypeScript asks of a placement. The evidence requirement is the
// one of the import environment of the stages, the same for every request: a
// TypeScript placement changes no compilation root. The guard compares the
// publication with the layout that the step derived from the admitted
// placement.
func (p Stages) GeneratedContract() (f.GeneratedFilesContract[s.Language], error) {
	stages, err := p.capture()
	if err != nil {
		return f.GeneratedFilesContract[s.Language]{}, err
	}
	evidence, err := n.EvidenceFor(stages.Configuration)
	if err != nil {
		return f.GeneratedFilesContract[s.Language]{}, err
	}
	prepare := func(ctx context.Context, preparation f.GeneratedFilesPreparation) (f.GeneratedFilesPrepared[s.Language], error) {
		zero := f.GeneratedFilesPrepared[s.Language]{}
		input := preparation.Input()
		placement, issued := placements.Read(input.Placement)
		if !issued {
			// The request carries no placement, or one that another language issued.
			return zero, fmt.Errorf("%w: the request carries no TypeScript placement", source.ErrInvalid)
		}
		admitted, err := preparation.Admit(f.GeneratedPlacementView{Files: placement.Files, Definitions: placement.Definitions})
		if err != nil {
			return zero, err
		}
		layout, err := placeGenerated(ctx, input.Output, admitted, input.Limits.Source)
		if err != nil {
			return zero, err
		}
		return f.GeneratedFilesPrepared[s.Language]{Placement: admitted, Evidence: evidence, Check: func(ctx context.Context, publication source.Publication[s.Language]) error {
			return checkLayout(ctx, layout, publication)
		}}, nil
	}
	return f.GeneratedFilesContract[s.Language]{Print: stages.Printer.Print, Prepare: prepare}, nil
}

// bittype:scope docs/typescript-type-imports.md#example-and-scope
func (p GeneratedProvider) Format(ctx context.Context, request f.GeneratedFilesRequest[s.Language]) (source.Source[s.Language], error) {
	input := request.Input()
	layout, err := placeGenerated(ctx, input.Output, request.Placement(), input.Limits.Source)
	if err != nil {
		return source.Source[s.Language]{}, err
	}
	return p.stages.organize(ctx, layout, input.Limits.Semantic, request.Publish)
}

// placeGenerated realizes an admitted placement of a sealed output. The order
// is the formatter's: the files by their identities, and the declarations of
// one file by their sites in the sealed output.
func placeGenerated(ctx context.Context, output g.Output, admitted f.AdmittedPlacement, limits source.Limits) (*result, error) {
	if !admitted.Valid() {
		return nil, source.ErrInvalid
	}
	p, err := generatedProgram(ctx, output)
	if err != nil {
		return nil, err
	}
	admittedFiles := admitted.Files()
	units := make([]Placement, len(admittedFiles))
	position := make(map[files.ID]int, len(admittedFiles))
	for i, file := range admittedFiles {
		units[i] = Placement{File: file.ID, Path: file.Path}
		position[file.ID] = i
	}
	definitions := admitted.Definitions()
	// The declarations that the output owns are the top-level declarations of
	// its trees, each once. Anything else is a defect of the output.
	if len(definitions) != len(p.tops) {
		return nil, fmt.Errorf("%w: the output owns %d declarations and its trees hold %d", source.ErrInvalid, len(definitions), len(p.tops))
	}
	for _, definition := range definitions {
		site := definition.Declaration.Site
		if !p.tops[site] {
			return nil, fmt.Errorf("%w: the owned declaration at %v is no top-level declaration", source.ErrInvalid, site)
		}
		unit := &units[position[definition.File]]
		unit.Declarations = append(unit.Declarations, site)
	}
	return placeProgram(ctx, p, units, limits)
}

// generatedProgram reads intent from the sealed output. Local destinations are
// syntax addresses and external destinations must come from the TypeScript issuer.
func generatedProgram(ctx context.Context, output g.Output) (*program, error) {
	snap := output.Snapshot()
	p := &program{trees: map[ast.TreeID]s.AST{}, tops: map[ast.NodeRef]bool{}, public: map[ast.NodeRef]bool{}, bindings: map[ast.NodeRef]symbol{}}
	for _, env := range snap.Inputs() {
		t, err := representation.Unpack[s.Language](snap.Language(), env)
		if err != nil {
			return nil, err
		}
		p.trees[t.ID] = t
		for _, c := range t.Nodes[t.Root-1].Children {
			if c.Role != s.RoleDeclaration {
				return nil, source.ErrInvalid
			}
			site := ref(t.ID, c.Node)
			if imports.Attr(t.Nodes[c.Node-1], s.AttributeExported) != "" {
				return nil, source.ErrInvalid
			}
			p.tops[site] = true
		}
		for i, node := range t.Nodes {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if node.Kind == s.KindTypeReference && len(node.References) > 0 {
				if len(node.References) != 1 || node.References[0].Role != s.RoleTarget {
					return nil, source.ErrInvalid
				}
				p.bindings[ref(t.ID, ast.NodeID(i+1))] = symbol{site: node.References[0].Target}
			}
		}
	}
	for _, site := range output.Exposed() {
		if !p.tops[site] {
			return nil, source.ErrInvalid
		}
		p.public[site] = true
	}
	for _, use := range output.ExternalUses() {
		site := use.At.Site
		ext, err := tssymbol.ReadExternal(snap.Language(), use)
		if err != nil {
			// A use that the TypeScript issuer cannot read is invalid input. The
			// generation value stays joined to the Files category.
			return nil, fmt.Errorf("%w: external use at %v: %w", source.ErrInvalid, site, err)
		}
		node := p.trees[site.Tree].Nodes[site.Node-1]
		if node.Kind != s.KindTypeReference || len(node.References) != 0 || imports.Attr(node, s.AttributeName) != ext.Name {
			return nil, fmt.Errorf("%w: external destination at %v", source.ErrInvalid, site)
		}
		if ext.Module == "" {
			p.bindings[site] = symbol{builtin: ext.Name}
		} else {
			p.bindings[site] = symbol{external: imports.Export{Module: ext.Module, Name: ext.Name}}
		}
	}
	return p, nil
}
