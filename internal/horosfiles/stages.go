// Package horosfiles composes explicitly supplied stages for the Horos Files
// profile. Canonical meaning remains in the Horos resolver and analyzer.
package horosfiles

import (
	"context"
	"fmt"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/horos/go"
	"github.com/Bitspark/bittype/source/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"reflect"
	"sort"
)

type Stages struct {
	Language                           representation.Language
	Parser                             parse.Provider[parse.Input, cst.CST[s.Language]]
	Formatter                          format.Provider[ast.View[s.Language], cst.CST[s.Language]]
	Printer                            print.Provider[cst.CST[s.Language], string]
	Checker                            check.Checker
	Resolver                           resolution.Resolver
	Analyzer                           analysis.Analyzer
	ResolutionOptions, AnalysisOptions semantic.Options
}

func (p Stages) Capture() (Stages, error) {
	if err := representation.CheckMarker[s.Language](p.Language); err != nil {
		return Stages{}, err
	}
	if representation.IsNil(p.Parser) || representation.IsNil(p.Printer) {
		return Stages{}, source.ErrInvalid
	}
	if !representation.IsNil(p.Formatter) && p.Formatter.ID() != p.Language.ID() {
		return Stages{}, representation.ErrLanguageMismatch
	}
	if p.Parser.ID() != p.Language.ID() || p.Printer.ID() != p.Language.ID() {
		return Stages{}, representation.ErrLanguageMismatch
	}
	if err := p.Checker.Validate(); err != nil {
		return Stages{}, err
	}
	if err := p.Resolver.Validate(); err != nil {
		return Stages{}, err
	}
	if err := p.Analyzer.Validate(); err != nil {
		return Stages{}, err
	}
	for _, definition := range []representation.Language{p.Checker.Language(), p.Resolver.Language(), p.Analyzer.Language()} {
		if !p.Language.Compatible(definition) {
			return Stages{}, representation.ErrLanguageMismatch
		}
	}
	// This Files composition requires the fixed Horos semantic contracts. Configured
	// resource options remain explicit and are captured with the matching receipts.
	if p.Resolver.Contract().Semantic.Profile != n.ResolutionProfile || p.Analyzer.Contract().Semantic.Profile != n.AnalysisProfile {
		return Stages{}, semantic.ErrInvalidOperation
	}
	for _, options := range []*semantic.Options{&p.ResolutionOptions, &p.AnalysisOptions} {
		settings, err := semantic.CaptureSettings(options.Configuration, options.Limits)
		if err != nil {
			return Stages{}, err
		}
		options.Configuration = settings
		if len(options.Requested) != 0 {
			return Stages{}, semantic.ErrInvalidRequest
		}
	}
	return p, nil
}
func (p Stages) Evidence() semantic.EvidenceRequirement {
	c := p.Analyzer.Contract().Semantic
	return semantic.EvidenceRequirement{Operation: analysis.Kind(), Profile: c.Profile, Configuration: p.AnalysisOptions.Configuration, Reports: c.SemanticsRequirements, Semantics: true}.Clone()
}

func attr(tree ast.AST[s.Language], id ast.NodeID, name string) string {
	for _, a := range tree.Nodes[id-1].Attributes {
		if a.Name == name {
			return a.Text
		}
	}
	return ""
}
func ModulePins(tree ast.AST[s.Language]) []string {
	var pins []string
	for _, child := range tree.Nodes[tree.Root-1].Children {
		pins = append(pins, attr(tree, child.Node, s.AttributePin))
	}
	sort.Strings(pins)
	return pins
}
func RootPin(tree ast.AST[s.Language]) string {
	for _, ref := range tree.Nodes[tree.Root-1].References {
		if ref.Role == s.RoleRoot {
			return attr(tree, ref.Target.Node, s.AttributePin)
		}
	}
	return ""
}
func EntryModuleSetting() semantics.Key {
	return semantics.Key{Owner: "horos/files", Name: "entry-module"}
}
func SelectEntry(tree *ast.AST[s.Language], settings []semantic.Setting) error {
	if len(settings) == 0 {
		return nil
	}
	if len(settings) != 1 || settings[0].Key != EntryModuleSetting() {
		return semantic.ErrInvalidRequest
	}
	for _, child := range tree.Nodes[tree.Root-1].Children {
		if attr(*tree, child.Node, s.AttributePin) == string(settings[0].Value) {
			tree.Nodes[tree.Root-1].References = []ast.Reference{{Role: s.RoleRoot, Target: ast.NodeRef{Tree: tree.ID, Node: child.Node}}}
			return nil
		}
	}
	return fmt.Errorf("%w: entry module is absent from selected file", semantic.ErrInvalidRequest)
}
func (p Stages) Establish(ctx context.Context, units []source.Unit[s.Language], entry files.ID, limits semantic.Limits) (source.Publication[s.Language], error) {
	input := semantic.Input{Limits: limits, Files: &files.Context{}}
	nodes, bytes := limits.MaxInputNodes, limits.MaxInputBytes
	for _, unit := range units {
		if err := ctx.Err(); err != nil {
			return source.Publication[s.Language]{}, err
		}
		// A limit that the units of a Files stage exceed is a Files limit. The
		// semantic value returned before stays joined to it.
		if len(unit.Syntax.Abstract.Nodes) > nodes {
			return source.Publication[s.Language]{}, source.WithCategory(semantic.ErrLimit)
		}
		nodes -= len(unit.Syntax.Abstract.Nodes)
		text, err := p.Printer.Print(ctx, unit.Syntax)
		if err != nil {
			return source.Publication[s.Language]{}, err
		}
		if len(text) > bytes {
			return source.Publication[s.Language]{}, source.WithCategory(semantic.ErrLimit)
		}
		bytes -= len(text)
		envelope, err := representation.Pack(p.Language, unit.Syntax)
		if err != nil {
			return source.Publication[s.Language]{}, err
		}
		input.Trees = append(input.Trees, envelope)
		input.Files.Files.Entries = append(input.Files.Files.Entries, files.File{ID: unit.File, Path: unit.Path, Text: text})
		input.Files.Links = append(input.Files.Links, files.SyntaxLink{File: unit.File, Tree: unit.Syntax.Abstract.ID, Source: unit.Source})
		if unit.File == entry {
			input.Entry = ast.NodeRef{Tree: unit.Syntax.Abstract.ID, Node: unit.Syntax.Abstract.Root}
		}
	}
	snapshot, err := semantic.Capture(ctx, p.Checker, input)
	if err != nil {
		// The capture answers with a files value or a semantic limit.
		return source.Publication[s.Language]{}, source.WithCategory(err)
	}
	ro := p.ResolutionOptions
	ro.Limits = limits
	resolved, err := p.Resolver.Resolve(ctx, snapshot, ro)
	if err != nil {
		return source.Publication[s.Language]{}, err
	}
	ao := p.AnalysisOptions
	ao.Limits = limits
	analyzed, err := p.Analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: ao})
	if err != nil {
		return source.Publication[s.Language]{}, err
	}
	if _, err := semantic.AdmitEvidence(snapshot, []semantic.EvidenceRequirement{p.Evidence()}, []semantic.Result{analyzed}); err != nil {
		return source.Publication[s.Language]{}, &semantic.Failure{Result: analyzed}
	}
	groups, relations, err := Organization(units)
	if err != nil {
		return source.Publication[s.Language]{}, err
	}
	return source.Publication[s.Language]{Snapshot: snapshot, Analysis: analyzed, Units: units, Groups: groups, Relations: relations}, nil
}

// Organization derives only physical assertions from explicit module/import syntax.
func Organization(units []source.Unit[s.Language]) ([]source.Group, []source.Relation, error) {
	var groups []source.Group
	var relations []source.Relation
	locations := map[string]files.ID{}
	for _, unit := range units {
		tree := unit.Syntax.Abstract
		for _, module := range tree.Nodes[tree.Root-1].Children {
			pin := attr(tree, module.Node, s.AttributePin)
			if pin == "" || locations[pin] != "" {
				return nil, nil, source.ErrInvalid
			}
			locations[pin] = unit.File
			groups = append(groups, source.Group{ID: source.GroupID(pin), Kind: n.ModuleKind(), Members: []files.ID{unit.File}, Sites: []ast.NodeRef{{Tree: tree.ID, Node: module.Node}}})
		}
	}
	for _, unit := range units {
		tree := unit.Syntax.Abstract
		for _, module := range tree.Nodes[tree.Root-1].Children {
			for _, child := range tree.Nodes[module.Node-1].Children {
				if child.Role == s.RoleImport {
					to := locations[attr(tree, child.Node, s.AttributePin)]
					if to == "" {
						if attr(tree, child.Node, s.AttributePin) == n.PreludePin() {
							continue
						}
						return nil, nil, source.ErrInvalid
					}
					at := ast.NodeRef{Tree: tree.ID, Node: child.Node}
					relations = append(relations, source.Relation{Kind: n.ImportRelation(), From: unit.File, To: to, At: &at})
				}
			}
		}
	}
	return groups, relations, nil
}

func (p Stages) CheckPublication(ctx context.Context, publication source.Publication[s.Language]) error {
	groups, relations, err := Organization(publication.Units)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(groups, publication.Groups) || !reflect.DeepEqual(relations, publication.Relations) {
		return source.ErrInvalid
	}
	graph, _ := publication.Analysis.Semantics()
	moduleSites := map[ast.NodeRef]bool{}
	for _, entity := range graph.Data().Entities {
		if entity.Kind == n.ModuleKind() {
			for _, site := range entity.Declarations {
				moduleSites[site] = true
			}
		}
	}
	// Every supplied module must be inside the admitted closure; no file is silently ignored.
	for _, group := range groups {
		if len(group.Sites) != 1 || !moduleSites[group.Sites[0]] {
			return fmt.Errorf("%w: supplied module is outside the entry closure", source.ErrInvalid)
		}
	}
	for _, unit := range publication.Units {
		text, err := p.Printer.Print(ctx, unit.Syntax)
		if err != nil {
			return err
		}
		parsed, err := p.Parser.Parse(ctx, parse.Input{SourceID: unit.Source, TreeID: unit.Syntax.Abstract.ID, Text: text})
		if err != nil {
			return err
		}
		// Language-owned admission compares canonical module content. This is an
		// agreement check, never the formatter's AST-to-CST construction algorithm.
		if !reflect.DeepEqual(ModulePins(parsed.Abstract), ModulePins(unit.Syntax.Abstract)) {
			return source.ErrInvalid
		}
	}
	return nil
}
