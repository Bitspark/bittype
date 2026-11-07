// Package target constructs typescript syntax. It contains no canonical type mappings.
package target

import (
	"context"
	"errors"

	check "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	lit "github.com/Bitspark/bittype/internal/typescript"
	imports "github.com/Bitspark/bittype/internal/typescriptimports"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Config struct{ TreeID ast.TreeID }
type Backend[F any] struct {
	*m.Store[F, s.Language]
	config   Config
	language representation.Language
	imports  map[g.Fragment[s.Language]]imports.Export
}

func New[F any](ctx context.Context, to representation.Language, config Config, limits g.Limits) (*Backend[F], error) {
	if to.ID() != (s.Language{}).ID() {
		return nil, g.ErrInvalid
	}
	store, err := m.New[F, s.Language](ctx, to, config.TreeID, limits)
	if err != nil {
		return nil, err
	}
	b := &Backend[F]{Store: store, config: config, language: to, imports: map[g.Fragment[s.Language]]imports.Export{}}
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

// ImportedType refers to a named type from a selected public package export.
func (b *Backend[F]) ImportedType(module, export string, args []g.Fragment[s.Language], origin ast.NodeRef) (g.Fragment[s.Language], error) {
	if !lit.PackageSpecifier(module) || !lit.TypeName(export) {
		return g.Fragment[s.Language]{}, g.ErrInvalid
	}
	ref, err := b.Node(m.Node(s.KindTypeReference, m.Attributes("name", export), origin))
	if err != nil {
		return ref, err
	}
	b.imports[ref] = imports.Export{Module: module, Name: export}
	return b.apply(ref, args, origin)
}
func (b *Backend[F]) Assemble(ctx context.Context, p []g.Publication[F, s.Language], d []g.Definition[F, s.Language], limits g.Limits) (representation.ASTSet, error) {
	if len(b.imports) == 0 {
		return b.Store.Assemble(ctx, p, d, limits)
	}
	watched := make([]g.Fragment[s.Language], 0, len(b.imports))
	for f := range b.imports {
		watched = append(watched, f)
	}
	result, sites, err := b.Store.AssembleWithSites(ctx, p, d, limits, watched)
	if err != nil {
		return representation.ASTSet{}, err
	}
	tree, err := representation.Unpack[s.Language](b.language, result.EntryAST())
	if err != nil {
		return representation.ASTSet{}, err
	}
	uses := map[ast.NodeID]imports.Export{}
	for id, f := range sites {
		uses[id] = b.imports[f]
	}
	if err = imports.Emit(ctx, &tree, uses, limits.OutputNodes, limits.OutputBytes); errors.Is(err, imports.ErrLimit) {
		return representation.ASTSet{}, g.ErrLimit
	}
	if err != nil {
		return representation.ASTSet{}, err
	}
	if err = check.Check(tree); err != nil {
		return representation.ASTSet{}, err
	}
	return representation.PackSet(b.language, ast.Set[s.Language]{Entry: result.Entry, Trees: []s.AST{tree}})
}
