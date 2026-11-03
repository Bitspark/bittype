// Package lift implements an explicitly selected partial Horos-to-Horos/Wire map.
// Horos owns lifting semantics; this package preserves AST correspondence and
// supplies the existing transformer SPI. It has no native or runtime adapter.
package lift

import (
	"errors"
	"maps"
	"slices"

	structure "github.com/Bitspark/bittype/check/ast/go"
	check "github.com/Bitspark/bittype/check/horoswire/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	bridge "github.com/Bitspark/bittype/internal/horosmodel"
	source "github.com/Bitspark/bittype/resolve/horos/go"
	target "github.com/Bitspark/bittype/resolve/horoswire/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	w "github.com/Bitspark/bittype/syntax/horoswire/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	ground "github.com/Bitspark/horos/core/go/declarations"
	wire "github.com/Bitspark/horos/core/go/wire"
)

var ErrOptions = errors.New("invalid Horos lift selection")

// Entry names the source module whose derived import closure is the output AST.
// Every selected root must be in that closure; independent roots require separate
// invocations (or Horos's multi-module Derivation API). No selected root is dropped.
type Options struct {
	Roots  []ground.DeclarationID
	Entry  ground.Digest
	TreeID ast.TreeID
	Lift   wire.LiftOptions
}

// Result retains canonical derivation evidence alongside the tooling projection.
// Tree is caller-owned storage; Derivation owns its immutable canonical evidence.
type Result struct {
	Tree       w.AST
	Derivation *wire.Derivation
}

// Selection is a separate bounded input. Check known collection sizes before
// copying them or invoking model resolution; name bytes are charged by occurrence.
func selection(options Options) error {
	l := options.Lift.Limits
	if options.TreeID == "" || len(options.Roots) == 0 || l.MaxDepth <= 0 || l.MaxDepth > 512 || l.MaxNodes <= 0 || l.MaxBytes <= 0 || l.MaxModules <= 0 {
		return ErrOptions
	}
	entries, bytes := min(l.MaxNodes, structure.MaxEntries), min(l.MaxBytes, structure.MaxTextBytes)
	take := func(left *int, count int) bool {
		if count > *left {
			return false
		}
		*left -= count
		return true
	}
	exhausted := func() error {
		return &text.Diagnostic{Code: text.LimitExceeded, Message: "lift selection exceeds its resource budget"}
	}
	if !take(&entries, len(options.Roots)) || !take(&entries, len(options.Lift.ModuleNames)) || !take(&entries, len(options.Lift.DeclarationNames)) || !take(&bytes, len(options.TreeID)) {
		return exhausted()
	}
	for _, root := range options.Roots {
		if !take(&bytes, len(root.Module)) || !take(&bytes, len(root.Name)) {
			return exhausted()
		}
	}
	for pin, name := range options.Lift.ModuleNames {
		if !take(&bytes, len(pin)) || !take(&bytes, len(name)) {
			return exhausted()
		}
	}
	for id, name := range options.Lift.DeclarationNames {
		if !take(&bytes, len(id.Module)) || !take(&bytes, len(id.Name)) || !take(&bytes, len(name)) {
			return exhausted()
		}
	}
	return nil
}

func Derive(input h.AST, options Options) (Result, error) {
	if err := selection(options); err != nil {
		return Result{}, err
	}
	if options.TreeID == "" || options.TreeID == input.ID {
		return Result{}, ErrOptions
	}
	program, err := source.Model(input)
	if err != nil {
		return Result{}, err
	}
	derivation, err := wire.Lift(program, options.Roots, options.Lift)
	if err != nil {
		return Result{}, err
	}
	var entry wire.Digest
	found := false
	for _, module := range derivation.Modules() {
		if module.Source() == options.Entry {
			entry = module.Pin()
			found = true
		}
	}
	if !found {
		return Result{}, ErrOptions
	}
	derived, err := derivation.Program(entry, options.Lift.Limits)
	if err != nil {
		return Result{}, err
	}
	for _, root := range derivation.Roots() {
		if _, ok := derived.Declaration(root.Declaration); !ok {
			return Result{}, ErrOptions
		}
	}
	output, err := target.Project(derived, options.TreeID)
	if err != nil {
		return Result{}, err
	}
	from, err := bridge.Index(input)
	if err != nil {
		return Result{}, err
	}
	to, err := bridge.Index(output)
	if err != nil {
		return Result{}, err
	}
	// Bound additional origin storage before appending, then check the combined
	// carrier budget. These are operational bounds, not type equivalence claims.
	entries, bytes := structure.MaxEntries, structure.MaxTextBytes
	attach := func(destination, origin ast.NodeRef) error {
		if destination.Tree != output.ID || destination.Node == 0 || int(destination.Node) > len(output.Nodes) || origin.Tree != input.ID || origin.Node == 0 || int(origin.Node) > len(input.Nodes) {
			return ErrOptions
		}
		inherited := input.Nodes[origin.Node-1].Origin
		count := len(inherited.Nodes) + len(inherited.Spans) + 1
		if count > entries {
			return &text.Diagnostic{Code: text.LimitExceeded, Message: "lift origin entries exhausted"}
		}
		take := func(n int) bool {
			if n > bytes {
				return false
			}
			bytes -= n
			return true
		}
		if !take(len(origin.Tree)) {
			return &text.Diagnostic{Code: text.LimitExceeded, Message: "lift origin bytes exhausted"}
		}
		for _, node := range inherited.Nodes {
			if !take(len(node.Tree)) {
				return &text.Diagnostic{Code: text.LimitExceeded, Message: "lift origin bytes exhausted"}
			}
		}
		for _, span := range inherited.Spans {
			if !take(len(span.Source)) {
				return &text.Diagnostic{Code: text.LimitExceeded, Message: "lift origin bytes exhausted"}
			}
		}
		entries -= count
		out := &output.Nodes[destination.Node-1].Origin
		out.Nodes = append(out.Nodes, origin)
		out.Nodes = append(out.Nodes, inherited.Nodes...)
		out.Spans = append(out.Spans, inherited.Spans...)
		return nil
	}
	if err := attach(ast.NodeRef{Tree: output.ID, Node: output.Root}, ast.NodeRef{Tree: input.ID, Node: input.Root}); err != nil {
		return Result{}, err
	}
	for _, module := range derivation.Modules() {
		if err := attach(to.ModuleSite(module.Pin()), from.ModuleSite(module.Source())); err != nil {
			return Result{}, err
		}
		model, _ := derived.Module(module.Pin())
		for _, imported := range model.Imports() {
			if err := attach(to.ImportSite(module.Pin(), imported.Alias()), from.ImportSite(module.Source(), imported.Alias())); err != nil {
				return Result{}, err
			}
		}
	}
	for _, declaration := range derivation.Declarations() {
		if err := attach(to.DeclarationSite(declaration.Target.Declaration), from.DeclarationSite(declaration.Source)); err != nil {
			return Result{}, err
		}
	}
	for _, parameter := range derivation.Parameters() {
		destination := ground.ParameterID{Declaration: parameter.Target.Declaration, Index: parameter.Index}
		if err := attach(to.ParameterSite(destination), from.ParameterSite(parameter.Source)); err != nil {
			return Result{}, err
		}
	}
	for _, origin := range derivation.Origins() {
		old := from.OperandSite(from.BodySite(origin.Source), origin.SourcePath)
		next := to.OperandSite(to.BodySite(origin.Target.Declaration), origin.TargetPath)
		if err := attach(next, old); err != nil {
			return Result{}, err
		}
	}
	if err := check.Check(output); err != nil {
		return Result{}, err
	}
	return Result{Tree: output, Derivation: derivation}, nil
}

// NewTransformer captures a specific selection; nothing is automatically
// registered. A failed lift never retries closed embedding or another mapping.
func NewTransformer(from, to representation.Language, options Options) (transform.ASTTransformer, error) {
	if from.ID() != (h.Language{}).ID() || from.Profile() != h.Profile || to.ID() != (w.Language{}).ID() || to.Profile() != w.Profile {
		return nil, ErrOptions
	}
	if err := selection(options); err != nil {
		return nil, err
	}
	options.Roots = slices.Clone(options.Roots)
	options.Lift.ModuleNames = maps.Clone(options.Lift.ModuleNames)
	options.Lift.DeclarationNames = maps.Clone(options.Lift.DeclarationNames)
	return transform.NewTransformer[h.Language, w.Language](from, to, func(input h.AST) (w.AST, error) {
		result, err := Derive(input, options)
		if err != nil {
			return w.AST{}, err
		}
		return result.Tree, nil
	})
}
