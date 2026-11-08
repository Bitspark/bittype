package generation

import (
	"context"
	"slices"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// SymbolIssuer is a target library's authority to issue and interpret immutable
// external destinations. Its key is opaque to the generation and Files kernels.
type SymbolIssuer struct{ owner *identity }
type ExternalSymbol struct {
	owner    *identity
	language representation.Language
	key      string
}

func NewSymbolIssuer() SymbolIssuer { return SymbolIssuer{fresh()} }
func (i SymbolIssuer) Issue(language representation.Language, key string) (ExternalSymbol, error) {
	if i.owner == nil || language.Validate() != nil || key == "" {
		return ExternalSymbol{}, ErrInvalid
	}
	return ExternalSymbol{i.owner, language, key}, nil
}
func (i SymbolIssuer) Read(language representation.Language, s ExternalSymbol) (string, bool) {
	if i.owner == nil || i.owner != s.owner || !language.Compatible(s.language) {
		return "", false
	}
	return s.key, true
}

type ExternalUse struct {
	At     semantics.Occurrence
	Symbol ExternalSymbol
}

// Assembly is a target construction record. Exposed contains the realized
// declaration site of each publication in the supplied publication order.
type Assembly struct {
	Syntax       representation.ASTSet
	Exposed      []ast.NodeRef
	ExternalUses []ExternalUse
}

// Output can only be produced by a successfully finished generation session.
// Sites and immutable target symbols always travel with their own checked capture.
type Output struct {
	snapshot semantic.Snapshot
	exposed  []ast.NodeRef
	external []ExternalUse
}

func (o Output) Snapshot() semantic.Snapshot { return o.snapshot }
func (o Output) Exposed() []ast.NodeRef      { return slices.Clone(o.exposed) }
func (o Output) ExternalUses() []ExternalUse { return slices.Clone(o.external) }

func seal(ctx context.Context, checker check.Checker, assembly Assembly, publications int, limits Limits) (Output, error) {
	set := assembly.Syntax
	if len(set.Trees) == 0 || len(set.Trees) > 256 || len(assembly.Exposed) != publications {
		return Output{}, ErrIncomplete
	}
	if len(assembly.ExternalUses) > limits.References || len(assembly.Exposed) > limits.References {
		return Output{}, ErrLimit
	}
	// Freeze target-owned slices before invoking the independently supplied checker.
	set.Trees = slices.Clone(set.Trees)
	assembly.Exposed = slices.Clone(assembly.Exposed)
	assembly.ExternalUses = slices.Clone(assembly.ExternalUses)
	trees := map[ast.TreeID]int{}
	nodes, bytes := limits.OutputNodes, limits.OutputBytes
	for _, tree := range set.Trees {
		info, err := tree.Inspect()
		if err != nil {
			return Output{}, err
		}
		if info.Nodes > nodes || info.TextBytes > bytes {
			return Output{}, ErrLimit
		}
		nodes -= info.Nodes
		bytes -= info.TextBytes
		trees[info.ID] = info.Nodes
	}
	contains := func(ref ast.NodeRef) bool { return ref.Node != 0 && uint64(ref.Node) <= uint64(trees[ref.Tree]) }
	exposed := map[ast.NodeRef]bool{}
	for _, site := range assembly.Exposed {
		if !contains(site) || exposed[site] {
			return Output{}, ErrScope
		}
		exposed[site] = true
	}
	seen := map[semantics.Occurrence]bool{}
	for _, use := range assembly.ExternalUses {
		if !contains(use.At.Site) || !use.At.Role.Valid() || seen[use.At] || use.Symbol.owner == nil || !checker.Language().Compatible(use.Symbol.language) {
			return Output{}, ErrScope
		}
		size := len(use.Symbol.key) + len(use.At.Role.Owner) + len(use.At.Role.Name)
		if size > bytes {
			return Output{}, ErrLimit
		}
		bytes -= size
		seen[use.At] = true
	}
	snapshot, err := semantic.Capture(ctx, checker, semantic.Input{Entry: set.Entry, Trees: set.Trees})
	if err != nil {
		return Output{}, err
	}
	return Output{snapshot, assembly.Exposed, assembly.ExternalUses}, nil
}
