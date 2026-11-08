// Package target constructs typescript syntax. It contains no canonical type mappings.
package target

import (
	"context"

	check "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	lit "github.com/Bitspark/bittype/internal/typescript"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/transform/syntax/typescript/symbol/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Config struct{ TreeID ast.TreeID }
type Backend[F any] struct {
	*m.Store[F, s.Language]
	config   Config
	language representation.Language
	external map[g.Fragment[s.Language]]symbol.Type
}

func New[F any](ctx context.Context, to representation.Language, config Config, limits g.Limits) (*Backend[F], error) {
	if to.ID() != (s.Language{}).ID() {
		return nil, g.ErrInvalid
	}
	store, err := m.New[F, s.Language](ctx, to, config.TreeID, limits)
	if err != nil {
		return nil, err
	}
	b := &Backend[F]{Store: store, config: config, language: to, external: map[g.Fragment[s.Language]]symbol.Type{}}
	store.Check = check.Check
	store.Finish = b.assemble
	return b, nil
}
func (b *Backend[F]) Parameter(name string, origin ast.NodeRef) (g.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindTypeParameter, m.Attributes("name", name), origin))
}
func (b *Backend[F]) Reference(symbol m.Symbol[s.Language], args []g.Fragment[s.Language], origin ast.NodeRef) (g.Fragment[s.Language], error) {
	node, err := b.Node(m.Node(s.KindTypeReference, m.Attributes("name", symbol.Name()), origin))
	if err != nil {
		return node, err
	}
	node, err = b.Link(node, symbol, s.RoleTarget)
	if err != nil {
		return node, err
	}
	return b.apply(node, args, origin)
}

func (b *Backend[F]) apply(node g.Fragment[s.Language], args []g.Fragment[s.Language], origin ast.NodeRef) (g.Fragment[s.Language], error) {
	if len(args) == 0 {
		return node, nil
	}
	children := append([]m.Child[s.Language]{{Role: s.RoleConstructor, Fragment: node}}, edges(s.RoleArgument, args)...)
	return b.Node(m.Node(s.KindTypeApplication, nil, origin), children...)
}
func edges(role string, fragments []g.Fragment[s.Language]) []m.Child[s.Language] {
	out := make([]m.Child[s.Language], len(fragments))
	for i, f := range fragments {
		out[i] = m.Child[s.Language]{Role: role, Fragment: f}
	}
	return out
}

// PackageType refers to a named export of a selected package.
func (b *Backend[F]) PackageType(module, export string, args []g.Fragment[s.Language], origin ast.NodeRef) (g.Fragment[s.Language], error) {
	if !lit.PackageSpecifier(module) || !lit.TypeName(export) {
		return g.Fragment[s.Language]{}, g.ErrInvalid
	}
	return b.externalType(module, export, args, origin)
}

// GlobalType identifies a public type supplied by the target library environment.
func (b *Backend[F]) GlobalType(name string, args []g.Fragment[s.Language], origin ast.NodeRef) (g.Fragment[s.Language], error) {
	if !lit.TypeName(name) {
		return g.Fragment[s.Language]{}, g.ErrInvalid
	}
	return b.externalType("", name, args, origin)
}
func (b *Backend[F]) externalType(module, export string, args []g.Fragment[s.Language], origin ast.NodeRef) (g.Fragment[s.Language], error) {
	ref, err := b.Node(m.Node(s.KindTypeReference, m.Attributes("name", export), origin))
	if err != nil {
		return ref, err
	}
	b.external[ref] = symbol.Type{Module: module, Name: export}
	return b.apply(ref, args, origin)
}
func (b *Backend[F]) Assemble(ctx context.Context, p []g.Publication[F, s.Language], d []g.Definition[F, s.Language], limits g.Limits) (g.Assembly, error) {
	watched := make([]g.Fragment[s.Language], 0, len(b.external))
	for f := range b.external {
		watched = append(watched, f)
	}
	result, sites, err := b.Store.AssembleWithSites(ctx, p, d, limits, watched)
	if err != nil {
		return g.Assembly{}, err
	}

	tree, err := representation.Unpack[s.Language](b.language, result.Syntax.EntryAST())
	if err != nil {
		return g.Assembly{}, err
	}
	for i, node := range tree.Nodes {
		if err := ctx.Err(); err != nil {
			return g.Assembly{}, err
		}
		if node.Kind != s.KindTypeReference {
			continue
		}
		id := ast.NodeID(i + 1)
		f, external := sites[id]
		if external {
			if len(node.References) != 0 {
				return g.Assembly{}, g.ErrScope
			}
			use, err := symbol.ExternalUse(b.language, ast.NodeRef{Tree: tree.ID, Node: id}, b.external[f])
			if err != nil {
				return g.Assembly{}, err
			}
			result.ExternalUses = append(result.ExternalUses, use)
		} else if len(node.References) != 1 || node.References[0].Role != s.RoleTarget {
			return g.Assembly{}, g.ErrIncomplete
		}
	}
	return result, nil
}
