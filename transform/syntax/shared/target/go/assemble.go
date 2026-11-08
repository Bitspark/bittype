package target

import (
	"context"
	"fmt"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"slices"
)

func (s *Store[F, T]) Assemble(ctx context.Context, publications []g.Publication[F, T], definitions []g.Definition[F, T], limits g.Limits) (g.Assembly, error) {
	return s.assemble(ctx, publications, definitions, limits, nil, nil)
}

// AssembleWithSites reports precisely the emitted occurrences of fragments owned
// by this store. A target can record the intended destination at these
// sites without encoding semantic metadata in identifier text or source origins.
func (s *Store[F, T]) AssembleWithSites(ctx context.Context, publications []g.Publication[F, T], definitions []g.Definition[F, T], limits g.Limits, watch []g.Fragment[T]) (g.Assembly, map[ast.NodeID]g.Fragment[T], error) {
	if ctx == nil || !limits.Valid() {
		return g.Assembly{}, nil, g.ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return g.Assembly{}, nil, err
	}
	if len(watch) > limits.OutputNodes {
		return g.Assembly{}, nil, g.ErrLimit
	}
	tracked := make(map[g.Fragment[T]]bool, len(watch))
	for _, f := range watch {
		if _, ok := s.fragmentRecords[f]; !ok {
			return g.Assembly{}, nil, g.ErrScope
		}
		tracked[f] = true
	}
	sites := map[ast.NodeID]g.Fragment[T]{}
	result, err := s.assemble(ctx, publications, definitions, limits, tracked, sites)
	if err != nil {
		return g.Assembly{}, nil, err
	}
	return result, sites, nil
}
func (s *Store[F, T]) assemble(ctx context.Context, publications []g.Publication[F, T], definitions []g.Definition[F, T], limits g.Limits, tracked map[g.Fragment[T]]bool, sites map[ast.NodeID]g.Fragment[T]) (g.Assembly, error) {
	fail := func(err error) (g.Assembly, error) { return g.Assembly{}, err }
	if s.Check == nil || s.Finish == nil {
		return fail(g.ErrInvalid)
	}
	defs := map[g.Destination[T]]g.Definition[F, T]{}
	keys := make([]g.Destination[T], 0, len(definitions))
	for _, def := range definitions {
		d := def.Callee().Destination()
		if _, ok := defs[d]; ok {
			return fail(g.ErrInvalid)
		}
		defs[d] = def
		keys = append(keys, d)
	}
	for _, p := range publications {
		if _, ok := defs[p.Callee().Destination()]; !ok {
			return fail(g.ErrIncomplete)
		}
	}
	slices.SortFunc(keys, func(a, b g.Destination[T]) int {
		if s.destinationRecords[a].name < s.destinationRecords[b].name {
			return -1
		}
		return 1
	})
	state := map[g.Destination[T]]byte{}
	var order []g.Destination[T]
	var visit func(g.Destination[T], int) error
	visit = func(d g.Destination[T], depth int) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		if depth > limits.Depth {
			return g.ErrLimit
		}
		if state[d] == 2 {
			return nil
		}
		if state[d] == 1 {
			return fmt.Errorf("%w: recursive native definition", g.ErrUnsupported)
		}
		def, ok := defs[d]
		if !ok {
			return g.ErrIncomplete
		}
		state[d] = 1
		seen := map[g.Fragment[T]]bool{}
		var walk func(g.Fragment[T], int) error
		walk = func(f g.Fragment[T], level int) error {
			if level > limits.Depth {
				return g.ErrLimit
			}
			if seen[f] {
				return nil
			}
			seen[f] = true
			r, ok := s.fragmentRecords[f]
			if !ok {
				return g.ErrScope
			}
			if r.referenceRole != "" && !r.symbol.IsFormal() {
				if err := visit(r.symbol.destination, depth+1); err != nil {
					return err
				}
			}
			for _, c := range r.children {
				if err := walk(c.Fragment, level+1); err != nil {
					return err
				}
			}
			return nil
		}
		if err := walk(def.Body(), 0); err != nil {
			return err
		}
		state[d] = 2
		order = append(order, d)
		return nil
	}
	for _, d := range keys {
		if err := visit(d, 0); err != nil {
			return fail(err)
		}
	}
	anchors := map[g.Fragment[T]]Symbol[T]{}
	var roots []g.Fragment[T]
	for _, d := range order {
		def := defs[d]
		l := s.layoutRecords[def.Callee().Layout()]
		record := s.destinationRecords[d]
		for i, p := range l.parameters {
			symbol, _ := s.Formal(def.Callee().Layout(), i)
			anchors[p] = symbol
		}
		root, err := record.build(slices.Clone(l.parameters), def.Body())
		if err != nil {
			return fail(err)
		}
		anchors[root] = Symbol[T]{destination: d, name: record.name}
		roots = append(roots, root)
	}
	root, err := s.Finish(roots)
	if err != nil {
		return fail(err)
	}
	out := ast.AST[T]{ID: s.id}
	symbols := map[Symbol[T]]ast.NodeID{}
	type relocation struct {
		node   ast.NodeID
		symbol Symbol[T]
		role   string
	}
	var relocations []relocation
	bytes := 0
	var emit func(g.Fragment[T], int) (ast.NodeID, error)
	emit = func(f g.Fragment[T], depth int) (ast.NodeID, error) {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		if depth > limits.Depth || len(out.Nodes) >= limits.OutputNodes {
			return 0, g.ErrLimit
		}
		r, ok := s.fragmentRecords[f]
		if !ok {
			return 0, g.ErrScope
		}
		n := r.node
		n.Attributes = slices.Clone(n.Attributes)
		n.Origin.Nodes = slices.Clone(n.Origin.Nodes)
		n.Origin.Spans = slices.Clone(n.Origin.Spans)
		bytes += len(n.Kind)
		for _, a := range n.Attributes {
			bytes += len(a.Name) + len(a.Text)
		}
		if bytes > limits.OutputBytes {
			return 0, g.ErrLimit
		}
		id := ast.NodeID(len(out.Nodes) + 1)
		out.Nodes = append(out.Nodes, n)
		if tracked[f] {
			sites[id] = f
		}
		if anchor, ok := anchors[f]; ok {
			if symbols[anchor] != 0 {
				return 0, g.ErrInvalid
			}
			symbols[anchor] = id
		}
		for _, child := range r.children {
			cid, err := emit(child.Fragment, depth+1)
			if err != nil {
				return 0, err
			}
			out.Nodes[id-1].Children = append(out.Nodes[id-1].Children, ast.Child{Role: child.Role, Node: cid})
		}
		if r.referenceRole != "" {
			relocations = append(relocations, relocation{id, r.symbol, r.referenceRole})
		}
		return id, nil
	}
	out.Root, err = emit(root, 0)
	if err != nil {
		return fail(err)
	}
	for _, r := range relocations {
		target := symbols[r.symbol]
		if target == 0 {
			return fail(g.ErrIncomplete)
		}
		out.Nodes[r.node-1].References = append(out.Nodes[r.node-1].References, ast.Reference{Role: r.role, Target: ast.NodeRef{Tree: s.id, Node: target}})
	}
	if err = s.Check(out); err != nil {
		return fail(err)
	}
	syntax, err := representation.PackSet(s.language, ast.Set[T]{Entry: ast.NodeRef{Tree: out.ID, Node: out.Root}, Trees: []ast.AST[T]{out}})
	if err != nil {
		return fail(err)
	}
	assembly := g.Assembly{Syntax: syntax}
	for _, pub := range publications {
		symbol, err := s.Named(pub.Callee())
		if err != nil {
			return fail(err)
		}
		id := symbols[symbol]
		if id == 0 {
			return fail(g.ErrIncomplete)
		}
		assembly.Exposed = append(assembly.Exposed, ast.NodeRef{Tree: out.ID, Node: id})
	}
	return assembly, nil
}
