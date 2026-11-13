package format

import (
	"context"
	"fmt"
	f "github.com/Bitspark/bittype/contracts/format/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
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
	evidence, err := n.EvidenceFor(p.Configuration)
	if err != nil {
		return f.GeneratedFilesFormatter[s.Language]{}, err
	}
	return f.NewGeneratedFilesFormatter(p.Language, f.GeneratedFilesContract[s.Language]{Evidence: evidence, Print: p.Printer.Print, Check: checkGenerated}, GeneratedProvider{stages: p})
}

// bittype:scope docs/typescript-type-imports.md#example-and-scope
func (p GeneratedProvider) Format(ctx context.Context, request f.GeneratedFilesRequest[s.Language]) (source.Source[s.Language], error) {
	input := request.Input()
	layout, err := placeGenerated(ctx, input)
	if err != nil {
		return source.Source[s.Language]{}, err
	}
	return p.stages.organize(ctx, layout, input.Limits.Semantic, request.Publish)
}
func checkGenerated(ctx context.Context, input f.GeneratedFilesInput, publication source.Publication[s.Language]) error {
	layout, err := placeGenerated(ctx, input)
	if err != nil {
		return err
	}
	return checkLayout(ctx, layout, publication)
}
func placeGenerated(ctx context.Context, input f.GeneratedFilesInput) (*result, error) {
	p, err := generatedProgram(ctx, input.Output)
	if err != nil {
		return nil, err
	}
	return placeProgram(ctx, p, input.Configuration, input.Limits.Source)
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
