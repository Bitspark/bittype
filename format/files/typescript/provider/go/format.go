package format

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	f "github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	"github.com/Bitspark/bittype/source/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Stages are the explicitly supplied collaborators for realizing and checking
// TypeScript units. Configuration is the native compiler's receipt-bound import
// environment. Organization is separate and belongs to each Format request.
type Stages struct {
	Language      representation.Language
	Syntax        f.Provider[ast.View[s.Language], cst.CST[s.Language]]
	Printer       print.Provider[cst.CST[s.Language], string]
	Checker       check.Checker
	Resolver      resolution.Resolver
	Analyzer      analysis.Analyzer
	Configuration []semantic.Setting
}

func (p Stages) capture() (Stages, error) {
	if representation.CheckMarker[s.Language](p.Language) != nil || representation.IsNil(p.Syntax) || representation.IsNil(p.Printer) {
		return Stages{}, source.ErrInvalid
	}
	if p.Syntax.ID() != p.Language.ID() || p.Printer.ID() != p.Language.ID() || !p.Language.Compatible(p.Checker.Language()) || !p.Language.Compatible(p.Resolver.Language()) || !p.Language.Compatible(p.Analyzer.Language()) {
		return Stages{}, source.ErrInvalid
	}
	_, enabled, err := n.Imports(p.Configuration)
	if err != nil || !enabled {
		return Stages{}, fmt.Errorf("%w: named TypeScript import environment required", source.ErrInvalid)
	}
	p.Configuration, err = semantic.CaptureSettings(p.Configuration, semantic.DefaultLimits())
	return p, err
}

// Capture establishes facts for existing TypeScript files.
func (p Stages) Capture(ctx context.Context, units []source.Unit[s.Language], entry files.ID, limits semantic.Limits) (f.FilesInput[s.Language], error) {
	stages, err := p.capture()
	if err != nil {
		return f.FilesInput[s.Language]{}, err
	}
	return stages.establish(ctx, units, entry, limits)
}
func (p Stages) establish(ctx context.Context, units []source.Unit[s.Language], entry files.ID, limits semantic.Limits) (f.FilesInput[s.Language], error) {
	zero := f.FilesInput[s.Language]{}
	if ctx == nil {
		return zero, source.ErrInvalid
	}
	if limits == (semantic.Limits{}) {
		limits = semantic.DefaultLimits()
	}
	if !limits.Valid() || len(units) == 0 {
		return zero, source.ErrInvalid
	}
	if len(units) > limits.MaxInputs {
		return zero, source.ErrLimit
	}
	captured := semantic.Input{Limits: limits, Files: &files.Context{}}
	nodes, bytes := limits.MaxInputNodes, limits.MaxInputBytes
	for _, u := range units {
		if err := ctx.Err(); err != nil {
			return zero, err
		}
		nodes -= len(u.Syntax.Abstract.Nodes)
		if nodes < 0 {
			return zero, source.ErrLimit
		}
		text, err := p.Printer.Print(ctx, u.Syntax)
		if err != nil {
			return zero, err
		}
		bytes -= len(text)
		if bytes < 0 {
			return zero, source.ErrLimit
		}
		tree, err := representation.Pack(p.Language, u.Syntax)
		if err != nil {
			return zero, err
		}
		captured.Trees = append(captured.Trees, tree)
		captured.Files.Links = append(captured.Files.Links, files.SyntaxLink{File: u.File, Tree: u.Syntax.Abstract.ID, Source: u.Source})
		captured.Files.Files.Entries = append(captured.Files.Files.Entries, files.File{ID: u.File, Path: u.Path, Text: text})
		if u.File == entry {
			captured.Entry = ref(u.Syntax.Abstract.ID, u.Syntax.Abstract.Root)
		}
	}
	snap, err := semantic.Capture(ctx, p.Checker, captured)
	if err != nil {
		// The capture answers with a files value or a semantic limit.
		return zero, source.WithCategory(err)
	}
	options := semantic.Options{Limits: limits, Configuration: p.Configuration}
	resolved, err := p.Resolver.Resolve(ctx, snap, options)
	if err != nil {
		return zero, err
	}
	result, err := p.Analyzer.Analyze(ctx, analysis.Input{Snapshot: snap, Resolution: &resolved, Options: options})
	if err != nil {
		return zero, err
	}
	required, err := n.EvidenceFor(p.Configuration)
	if err != nil {
		return zero, err
	}
	if _, err = semantic.AdmitEvidence(snap, []semantic.EvidenceRequirement{required}, []semantic.Result{result}); err != nil {
		return zero, fmt.Errorf("%w: native TypeScript: %v", err, result.Diagnostics())
	}
	envelope, _ := result.Semantics()
	graph, err := representation.UnpackSemantics[s.Language](p.Language, envelope)
	if err != nil {
		return zero, err
	}
	return f.FilesInput[s.Language]{Abstract: graph, Snapshot: snap, Evidence: result}, nil
}

type Provider struct {
	s.Language
	stages Stages
}

func New(stages Stages) (f.FilesFormatter[s.Language], error) {
	p, err := stages.capture()
	if err != nil {
		return f.FilesFormatter[s.Language]{}, err
	}
	evidence, err := n.EvidenceFor(p.Configuration)
	if err != nil {
		return f.FilesFormatter[s.Language]{}, err
	}
	return f.NewFilesFormatter(p.Language, f.FilesContract[s.Language]{Evidence: evidence, Print: p.Printer.Print, Check: checkPublication}, Provider{stages: p})
}

// bittype:scope docs/typescript-type-imports.md#example-and-scope
func (p Provider) Format(ctx context.Context, r f.FilesRequest[s.Language]) (source.Source[s.Language], error) {
	input := r.Input()
	layout, err := materialize(ctx, input)
	if err != nil {
		return source.Source[s.Language]{}, err
	}
	return p.stages.organize(ctx, layout, input.Limits.Semantic, r.Publish)
}
func (p Stages) organize(ctx context.Context, layout *result, limits semantic.Limits, publish func(source.Publication[s.Language]) (source.Source[s.Language], error)) (source.Source[s.Language], error) {
	var units []source.Unit[s.Language]
	for _, u := range layout.units {
		syntax, err := p.Syntax.Format(ctx, u.tree)
		if err != nil {
			return source.Source[s.Language]{}, err
		}
		units = append(units, source.Unit[s.Language]{File: u.placement.File, Path: u.placement.Path, Source: ast.SourceID(u.placement.File), Syntax: syntax})
	}
	fresh, err := p.establish(ctx, units, layout.entry, limits)
	if err != nil {
		return source.Source[s.Language]{}, err
	}
	return publish(source.Publication[s.Language]{Snapshot: fresh.Snapshot, Analysis: fresh.Evidence, Units: units, Groups: layout.groups(), Relations: layout.relations})
}
func (x *result) groups() []source.Group {
	groups := make([]source.Group, 0, len(x.units))
	for _, u := range x.units {
		groups = append(groups, source.Group{ID: source.GroupID(u.placement.File), Kind: n.ModuleKind(), Members: []files.ID{u.placement.File}, Sites: []ast.NodeRef{ref(u.tree.ID, u.tree.Root)}})
	}
	return groups
}
func checkPublication(ctx context.Context, input f.FilesInput[s.Language], out source.Publication[s.Language]) error {
	expected, err := materialize(ctx, input)
	if err != nil {
		return err
	}
	return checkLayout(ctx, expected, out)
}
func checkLayout(ctx context.Context, expected *result, out source.Publication[s.Language]) error {
	if len(expected.units) != len(out.Units) || !reflect.DeepEqual(expected.groups(), out.Groups) || !reflect.DeepEqual(expected.relations, out.Relations) {
		return source.ErrInvalid
	}
	graphEnvelope, ok := out.Analysis.Semantics()
	if !ok {
		return semantic.ErrEvidenceUnavailable
	}
	graph, err := representation.UnpackSemantics[s.Language](out.Snapshot.Language(), graphEnvelope)
	if err != nil {
		return err
	}
	realized, err := readProgram(ctx, out.Snapshot, graph)
	if err != nil {
		return err
	}
	for site, want := range expected.public {
		if !realized.tops[site] || realized.public[site] != want {
			return fmt.Errorf("%w: public exposure changed at %v", source.ErrInvalid, site)
		}
	}
	for i, u := range expected.units {
		actual := out.Units[i]
		if u.placement.File != actual.File || u.placement.Path != actual.Path || actual.Source != ast.SourceID(actual.File) || !reflect.DeepEqual(u.tree, actual.Syntax.Abstract) {
			return source.ErrInvalid
		}
		if u.placement.File == expected.entry && out.Snapshot.Entry() != ref(u.tree.ID, u.tree.Root) {
			return source.ErrInvalid
		}
		for site, original := range u.expected {
			want := original
			if want.site.Tree != "" {
				want.site = expected.moved[want.site]
			}
			if got, ok := realized.bindings[ref(u.tree.ID, site)]; !ok || got != want {
				return fmt.Errorf("%w: moved reference changed its compiler binding in %s", source.ErrInvalid, u.placement.Path)
			}
		}
	}
	return nil
}

var _ f.Provider[f.FilesRequest[s.Language], source.Source[s.Language]] = Provider{}
