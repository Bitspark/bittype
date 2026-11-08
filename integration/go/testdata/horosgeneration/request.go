// Package horosgeneration is consumer configuration for integration fixtures.
// It supplies explicit providers and root sites; all assembly is kernel-owned.
package horosgeneration

import (
	"sync"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func Roots(input transform.Input) ([]ast.NodeRef, error) {
	envelope, err := input.Source.EntryTree()
	if err != nil {
		return nil, err
	}
	tree, err := representation.Unpack[h.Language](input.Source.Language(), envelope)
	if err != nil {
		return nil, err
	}
	var root ast.NodeRef
	for _, r := range tree.Nodes[tree.Root-1].References {
		if r.Role == h.RoleRoot {
			root = r.Target
		}
	}
	for _, envelope := range input.Source.Inputs() {
		tree, err := representation.Unpack[h.Language](input.Source.Language(), envelope)
		if err != nil {
			return nil, err
		}
		if tree.ID != root.Tree {
			continue
		}
		if root.Node == 0 || int(root.Node) > len(tree.Nodes) {
			return nil, g.ErrInvalid
		}
		var out []ast.NodeRef
		for _, child := range tree.Nodes[root.Node-1].Children {
			if tree.Nodes[child.Node-1].Kind == h.KindTypeDefinition {
				out = append(out, ast.NodeRef{Tree: tree.ID, Node: child.Node})
			}
		}
		return out, nil
	}
	return nil, g.ErrInvalid
}

var configuredSources sync.Map

func Request[T any](registry *c.Registry, input transform.Input, to representation.Language, output ast.TreeID, target g.Target[T], mappings ...g.MappingLibrary[h.Language, T, source.Reader]) (c.Request, error) {
	zero := c.Request{}
	from := input.Source.Language()
	if _, loaded := configuredSources.LoadOrStore(registry, true); !loaded {
		if err := c.RegisterSource[h.Language, source.Reader](registry, from, source.Provider{}); err != nil {
			return zero, err
		}
	}
	if err := c.RegisterTarget(registry, to, target); err != nil {
		return zero, err
	}
	for _, p := range mappings {
		if _, err := c.RegisterMapping(registry, from, to, p); err != nil {
			return zero, err
		}
	}
	roots, err := Roots(input)
	if err != nil {
		return zero, err
	}
	return c.Request{Input: input, To: to, Output: output, Roots: roots, Profile: "pagination-representation/1", Limits: g.DefaultLimits()}, nil
}
