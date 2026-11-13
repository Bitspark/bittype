package target

import (
	"cmp"
	"context"
	"slices"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// compareDeclarations orders the declarations that seed construction. Its key
// is the name, then the provenance, then the planned identity. Equal names
// compare equal and the next part of the key decides, so
// compareDeclarations(a, b) == -compareDeclarations(b, a) for every pair. It
// returns 0 only for one and the same planned declaration.
//
// The provenance is the planner's own fact: the declaring sites of the source
// declarations a declaration represents. Two provenances compare site by site,
// each site by its tree and then its node, and a list sorts before every list
// that continues it. No provenance, which is what a helper has, sorts first.
// The origin that a mapping may state on its description is not part of the
// key, because a mapping may leave it out.
//
// Limit: two declarations with equal names and equal provenance are ordered by
// their planned identity, which follows the order of the requested roots. Two
// native declarations of one source declaration are such a pair, and so are
// two helpers of one name.
func compareDeclarations[T any](a, b g.PlannedDeclaration[T]) int {
	return cmp.Or(
		cmp.Compare(a.Description.Name, b.Description.Name),
		slices.CompareFunc(a.Provenance, b.Provenance, compareSites),
		cmp.Compare(a.ID, b.ID),
	)
}

// compareSites orders two source sites by their tree and then their node.
func compareSites(a, b ast.NodeRef) int {
	return cmp.Or(cmp.Compare(a.Tree, b.Tree), cmp.Compare(a.Node, b.Node))
}

// Compile builds the tree of a closed program. It holds no rule about names
// across declarations: it identifies a declaration by its planned identity and
// a parameter by its position, and it reads a name only to write it. Two
// declarations of one name get two sites, and a parameter may be named like a
// declaration. A target that needs one flat scope states that in its program
// rule, and a language whose scopes follow from a placement decides its
// collisions in its Files stage. Compile keeps the one rule whose meaning is
// inside a declaration: the parameters of one declaration have names, and no
// two of them share one.
func Compile[T any](ctx context.Context, in g.BuildInput[T], store *Store[T], lower Lowering[T]) (g.Assembly, error) {
	fail := func(e error) (g.Assembly, error) { return g.Assembly{}, e }
	declarations := in.Program.Declarations()
	if ctx == nil || store == nil || len(declarations) == 0 {
		return fail(g.ErrInvalid)
	}
	anchors := map[Fragment[T]]Symbol[T]{}
	parameters := map[uint32][]Fragment[T]{}
	for _, d := range declarations {
		if d.ID == 0 || int(d.ID) > len(declarations) {
			return fail(g.ErrScope)
		}
	}
	for _, d := range declarations {
		seen := map[string]bool{}
		for i, p := range d.Description.Parameters {
			if p.Name == "" || seen[p.Name] {
				return fail(g.ErrScope)
			}
			seen[p.Name] = true
			f, e := lower.Parameter(p.Name, p.Origin)
			if e != nil {
				return fail(e)
			}
			parameters[d.ID] = append(parameters[d.ID], f)
			anchors[f] = Symbol[T]{store, d.ID, i, p.Name}
		}
	}
	state := map[uint32]byte{}
	roots := []Fragment[T]{}
	operations := 0
	var declaration func(uint32, int) error
	var expression func(g.Expression[T], uint32, int) (Fragment[T], error)
	expression = func(e g.Expression[T], scope uint32, depth int) (Fragment[T], error) {
		if err := ctx.Err(); err != nil {
			return Fragment[T]{}, err
		}
		operations++
		if depth > in.Limits.Depth || operations > in.Limits.Operations {
			return Fragment[T]{}, g.ErrLimit
		}
		args := []Fragment[T]{}
		for _, a := range e.Arguments() {
			f, err := expression(a, scope, depth+1)
			if err != nil {
				return f, err
			}
			args = append(args, f)
		}
		switch e.Kind() {
		case g.ParameterExpression:
			if e.Declaration() != scope || e.Index() < 0 || e.Index() >= len(parameters[scope]) || len(args) != 0 {
				return Fragment[T]{}, g.ErrScope
			}
			return lower.Reference(Symbol[T]{store, scope, e.Index(), declarations[scope-1].Description.Parameters[e.Index()].Name}, nil, e.Origin())
		case g.DeclarationExpression:
			id := e.Declaration()
			if id == 0 || int(id) > len(declarations) {
				return Fragment[T]{}, g.ErrScope
			}
			if len(args) != len(declarations[id-1].Description.Parameters) {
				return Fragment[T]{}, g.ErrScope
			}
			if err := declaration(id, depth+1); err != nil {
				return Fragment[T]{}, err
			}
			return lower.Reference(Symbol[T]{store, id, -1, declarations[id-1].Description.Name}, args, e.Origin())
		case g.NativeExpression:
			return lower.Expression(e, args)
		default:
			return Fragment[T]{}, g.ErrScope
		}
	}
	declaration = func(id uint32, depth int) error {
		if state[id] == 2 {
			return nil
		}
		if state[id] == 1 {
			return g.ErrUnsupported
		}
		if depth > in.Limits.Depth {
			return g.ErrLimit
		}
		state[id] = 1
		d := declarations[id-1]
		body, err := expression(d.Description.Body, id, depth+1)
		if err != nil {
			return err
		}
		root, err := lower.Declaration(d.Description, parameters[id], body)
		if err != nil {
			return err
		}
		anchors[root] = Symbol[T]{store, id, -1, d.Description.Name}
		roots = append(roots, root)
		state[id] = 2
		return nil
	}
	order := slices.Clone(declarations)
	slices.SortFunc(order, compareDeclarations[T])
	for _, d := range order {
		if err := declaration(d.ID, 0); err != nil {
			return fail(err)
		}
	}
	root, err := lower.Finish(roots)
	if err != nil {
		return fail(err)
	}
	tree := ast.AST[T]{ID: in.Output}
	symbols := map[Symbol[T]]ast.NodeID{}
	type relocation struct {
		site   ast.NodeID
		symbol Symbol[T]
		role   string
	}
	var refs []relocation
	var external []g.ExternalUse
	bytes := 0
	var emit func(Fragment[T], int) (ast.NodeID, error)
	emit = func(f Fragment[T], depth int) (ast.NodeID, error) {
		if err := ctx.Err(); err != nil {
			return 0, err
		}
		if depth > in.Limits.Depth || len(tree.Nodes) >= in.Limits.OutputNodes {
			return 0, g.ErrLimit
		}
		rec, err := store.read(f)
		if err != nil {
			return 0, err
		}
		n := rec.node
		n.Attributes = slices.Clone(n.Attributes)
		n.Origin.Nodes = slices.Clone(n.Origin.Nodes)
		n.Origin.Spans = slices.Clone(n.Origin.Spans)
		bytes += len(n.Kind)
		for _, a := range n.Attributes {
			bytes += len(a.Name) + len(a.Text)
		}
		if bytes > in.Limits.OutputBytes {
			return 0, g.ErrLimit
		}
		id := ast.NodeID(len(tree.Nodes) + 1)
		tree.Nodes = append(tree.Nodes, n)
		if sym, ok := anchors[f]; ok {
			if symbols[sym] != 0 {
				return 0, g.ErrScope
			}
			symbols[sym] = id
		}
		for _, child := range rec.children {
			cid, err := emit(child.Fragment, depth+1)
			if err != nil {
				return 0, err
			}
			tree.Nodes[id-1].Children = append(tree.Nodes[id-1].Children, ast.Child{Role: child.Role, Node: cid})
		}
		if rec.external != nil {
			use := *rec.external
			use.At.Site = ast.NodeRef{Tree: in.Output, Node: id}
			external = append(external, use)
		} else if rec.role != "" {
			refs = append(refs, relocation{id, rec.symbol, rec.role})
		}
		return id, nil
	}
	tree.Root, err = emit(root, 0)
	if err != nil {
		return fail(err)
	}
	for _, ref := range refs {
		target := symbols[ref.symbol]
		if target == 0 {
			return fail(g.ErrIncomplete)
		}
		tree.Nodes[ref.site-1].References = append(tree.Nodes[ref.site-1].References, ast.Reference{Role: ref.role, Target: ast.NodeRef{Tree: in.Output, Node: target}})
	}
	syntax, err := representation.PackSet(in.Language, ast.Set[T]{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: []ast.AST[T]{tree}})
	if err != nil {
		return fail(err)
	}
	// The inventory: for every planned declaration, in planned order, the node
	// that the lowering returned for it. That node was emitted once. Where it
	// sits in the unit is the lowering's decision. The kernel adds what the plan
	// knows about each declaration.
	out := g.Assembly{Syntax: syntax, ExternalUses: external}
	for _, d := range declarations {
		node := symbols[Symbol[T]{store, d.ID, -1, d.Description.Name}]
		if node == 0 {
			return fail(g.ErrIncomplete)
		}
		out.Declarations = append(out.Declarations, ast.NodeRef{Tree: in.Output, Node: node})
	}
	return out, nil
}
