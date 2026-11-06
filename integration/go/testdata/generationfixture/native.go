package generationfixture

import (
	"context"
	"fmt"
	"maps"
	"reflect"
	"slices"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	gp "github.com/Bitspark/bittype/provider/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type nativeIdentity struct {
	run       g.RunID
	config    g.ConfigID
	frame     g.FrameID
	signature g.SignatureID
}

func identity(site g.Site[SourceLanguage, s.Language]) nativeIdentity {
	return nativeIdentity{site.Run(), site.Config(), site.Frame(), site.Signature()}
}

type destinationRecord struct {
	identity nativeIdentity
	name     string
	owned    bool
	alias    bool
	// An external export has a supplied artifact identity and native arity. It
	// has no fabricated declaration NodeRef in the generated AST.
	artifact      string
	externalArity int
}

type layoutRecord struct {
	identity    nativeIdentity
	destination g.Destination[s.Language]
	parameters  []string
}

type symbol struct {
	destination g.Destination[s.Language]
	layout      g.Layout[s.Language]
	formal      int
}

// A fragment is an immutable existing Go AST subtree with typed symbolic
// relocations. It is not an alternative grammar or caller-provided free list.
type fragmentRecord struct {
	tree  s.AST
	links map[ast.NodeID]symbol
}

type Native struct {
	language           representation.Language
	destinations       g.Space[g.DestinationRole[s.Language]]
	layouts            g.Space[g.LayoutRole[s.Language]]
	fragments          g.Space[g.FragmentRole[s.Language]]
	destinationRecords map[g.Destination[s.Language]]destinationRecord
	layoutRecords      map[g.Layout[s.Language]]layoutRecord
	fragmentRecords    map[g.Fragment[s.Language]]fragmentRecord
	definitions        map[g.Destination[s.Language]]g.Layout[s.Language]
	storedNodes        int
}

func newNative() (*Native, error) {
	language, err := representation.NewLanguage[s.Language](s.Language{}, "native-types-v1")
	if err != nil {
		return nil, err
	}
	return &Native{
		language:           language,
		destinations:       g.NewSpace[g.DestinationRole[s.Language]](),
		layouts:            g.NewSpace[g.LayoutRole[s.Language]](),
		fragments:          g.NewSpace[g.FragmentRole[s.Language]](),
		destinationRecords: make(map[g.Destination[s.Language]]destinationRecord),
		layoutRecords:      make(map[g.Layout[s.Language]]layoutRecord),
		fragmentRecords:    make(map[g.Fragment[s.Language]]fragmentRecord),
		definitions:        make(map[g.Destination[s.Language]]g.Layout[s.Language]),
	}, nil
}

func (n *Native) reserve(site g.Site[SourceLanguage, s.Language], name string, owned, alias bool, artifact string, arity int) (g.Destination[s.Language], bool, error) {
	if !site.Valid() || !site.Signature().Valid() || (!owned && artifact == "") {
		return g.Destination[s.Language]{}, false, g.ErrInvalid
	}
	d, _ := n.destinations.Ref(uint64(len(n.destinationRecords) + 1))
	n.destinationRecords[d] = destinationRecord{identity: identity(site), name: name, owned: owned, alias: alias, artifact: artifact, externalArity: arity}
	return d, owned, nil
}

func (n *Native) freeze(site g.Site[SourceLanguage, s.Language], d g.Destination[s.Language], names ...string) (g.Layout[s.Language], error) {
	r, ok := n.destinationRecords[d]
	if !ok || r.identity != identity(site) || n.definitions[d].Valid() {
		return g.Layout[s.Language]{}, g.ErrScope
	}
	if !r.owned && len(names) != r.externalArity {
		return g.Layout[s.Language]{}, g.ErrInvalid
	}
	seen := make(map[string]bool)
	for _, name := range names {
		if name == "" || seen[name] {
			return g.Layout[s.Language]{}, g.ErrInvalid
		}
		seen[name] = true
	}
	l, _ := n.layouts.Ref(uint64(len(n.layoutRecords) + 1))
	n.layoutRecords[l] = layoutRecord{identity: identity(site), destination: d, parameters: slices.Clone(names)}
	n.definitions[d] = l
	return l, nil
}

func (n *Native) CheckDestination(ctx context.Context, site g.Site[SourceLanguage, s.Language], d g.Destination[s.Language], owned bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r, ok := n.destinationRecords[d]
	if !ok || r.identity != identity(site) || r.owned != owned {
		return g.ErrScope
	}
	return nil
}

func (n *Native) CheckLayout(ctx context.Context, site g.Site[SourceLanguage, s.Language], d g.Destination[s.Language], l g.Layout[s.Language]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r, ok := n.layoutRecords[l]
	if !ok || r.identity != identity(site) || r.destination != d || n.definitions[d] != l {
		return g.ErrScope
	}
	return nil
}

func (n *Native) CheckProduct(ctx context.Context, want g.Want[SourceLanguage, s.Language], f g.Fragment[s.Language]) error {
	return n.visible(ctx, f, identity(want.Site()))
}

func (n *Native) CheckBody(ctx context.Context, callee g.Callee[SourceLanguage, s.Language], f g.Fragment[s.Language]) error {
	l, ok := n.layoutRecords[callee.Layout()]
	if !ok || l.destination != callee.Destination() || l.identity.frame != callee.Frame() || l.identity.signature != callee.Signature() {
		return g.ErrScope
	}
	return n.visible(ctx, f, l.identity)
}

func (n *Native) CheckExposure(ctx context.Context, site g.Site[SourceLanguage, s.Language], root g.Declaration[SourceLanguage], callee g.Callee[SourceLanguage, s.Language]) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	d, ok := n.destinationRecords[callee.Destination()]
	if !ok || !d.owned || d.identity.run != site.Run() || d.identity.config != site.Config() || callee.Open().Declaration != root {
		return g.ErrScope
	}
	return nil
}

func (n *Native) visible(ctx context.Context, f g.Fragment[s.Language], in nativeIdentity) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r, ok := n.fragmentRecords[f]
	if !ok {
		return g.ErrScope
	}
	for _, link := range r.links {
		if link.layout.Valid() {
			l, exists := n.layoutRecords[link.layout]
			if !exists || l.identity.run != in.run || l.identity.config != in.config || l.identity.frame != in.frame || link.formal < 0 || link.formal >= len(l.parameters) {
				return g.ErrScope
			}
		} else {
			d, exists := n.destinationRecords[link.destination]
			if !exists || d.identity.run != in.run || d.identity.config != in.config {
				return g.ErrScope
			}
		}
	}
	return nil
}

func cloneNode(node s.Node) s.Node {
	node.Attributes = slices.Clone(node.Attributes)
	node.Children = slices.Clone(node.Children)
	node.References = slices.Clone(node.References)
	node.Origin.Nodes = slices.Clone(node.Origin.Nodes)
	node.Origin.Spans = slices.Clone(node.Origin.Spans)
	return node
}

func (n *Native) store(tree s.AST, links map[ast.NodeID]symbol) (g.Fragment[s.Language], error) {
	if len(n.fragmentRecords) >= 4096 || len(tree.Nodes) > 200000-n.storedNodes {
		return g.Fragment[s.Language]{}, g.ErrLimit
	}
	f, _ := n.fragments.Ref(uint64(len(n.fragmentRecords) + 1))
	tree.ID = ast.TreeID(fmt.Sprintf("fragment-%d", len(n.fragmentRecords)+1))
	owned := make([]s.Node, len(tree.Nodes))
	for i, node := range tree.Nodes {
		owned[i] = cloneNode(node)
	}
	tree.Nodes = owned
	n.fragmentRecords[f] = fragmentRecord{tree: tree, links: maps.Clone(links)}
	n.storedNodes += len(tree.Nodes)
	return f, nil
}

// CheckCaptureIsolation exercises the target fixture's ownership boundary using
// caller-retained mutable construction storage. It emits no generated artifact.
func CheckCaptureIsolation() error {
	native, err := newNative()
	if err != nil {
		return err
	}
	tree := s.AST{Root: 2, Nodes: []s.Node{
		{Kind: s.KindNamedType, Attributes: []ast.Attribute{{Name: s.AttributeName, Text: "string"}}, Origin: ast.Origin{Nodes: []ast.NodeRef{{Tree: "input", Node: 1}}, Spans: []ast.Span{{Source: "input", Start: 1, End: 2}}}},
		{Kind: s.KindSliceType, Children: []ast.Child{{Role: s.RoleElement, Node: 1}}},
	}}
	links := map[ast.NodeID]symbol{1: {formal: 7}}
	fragment, err := native.store(tree, links)
	if err != nil {
		return err
	}
	tree.Nodes[0].Attributes[0].Text = "mutated"
	tree.Nodes[0].Origin.Nodes[0].Node = 999
	tree.Nodes[0].Origin.Spans[0].End = 999
	tree.Nodes[1].Children[0].Node = 999
	links[1] = symbol{formal: 999}
	stored := native.fragmentRecords[fragment]
	if stored.tree.Nodes[0].Attributes[0].Text != "string" || stored.tree.Nodes[0].Origin.Nodes[0].Node != 1 || stored.tree.Nodes[0].Origin.Spans[0].End != 2 || stored.tree.Nodes[1].Children[0].Node != 1 || stored.links[1].formal != 7 {
		return fmt.Errorf("native fragment retained mutable construction storage")
	}
	return nil
}

func (n *Native) named(name string, link symbol) (g.Fragment[s.Language], error) {
	tree := s.AST{Root: 1, Nodes: []s.Node{{Kind: s.KindNamedType, Attributes: []ast.Attribute{{Name: s.AttributeName, Text: name}}}}}
	links := make(map[ast.NodeID]symbol)
	if link.destination.Valid() || link.layout.Valid() {
		links[1] = link
	}
	return n.store(tree, links)
}

func (n *Native) text() (g.Fragment[s.Language], error) { return n.named("string", symbol{}) }

func (n *Native) formal(l g.Layout[s.Language], index int) (g.Fragment[s.Language], error) {
	r, ok := n.layoutRecords[l]
	if !ok || index < 0 || index >= len(r.parameters) {
		return g.Fragment[s.Language]{}, g.ErrScope
	}
	return n.named(r.parameters[index], symbol{layout: l, formal: index})
}

type fragmentPart struct {
	role  string
	value g.Fragment[s.Language]
}

func (n *Native) compose(kind s.Kind, attributes []ast.Attribute, parts ...fragmentPart) (g.Fragment[s.Language], error) {
	tree := s.AST{}
	links := make(map[ast.NodeID]symbol)
	var children []ast.Child
	for _, part := range parts {
		r, ok := n.fragmentRecords[part.value]
		if !ok {
			return g.Fragment[s.Language]{}, g.ErrScope
		}
		if len(r.tree.Nodes) > 200000-len(tree.Nodes) {
			return g.Fragment[s.Language]{}, g.ErrLimit
		}
		offset := ast.NodeID(len(tree.Nodes))
		for _, node := range r.tree.Nodes {
			copy := cloneNode(node)
			for i := range copy.Children {
				copy.Children[i].Node += offset
			}
			tree.Nodes = append(tree.Nodes, copy)
		}
		for at, link := range r.links {
			links[at+offset] = link
		}
		children = append(children, ast.Child{Role: part.role, Node: r.tree.Root + offset})
	}
	tree.Nodes = append(tree.Nodes, s.Node{Kind: kind, Attributes: slices.Clone(attributes), Children: children})
	tree.Root = ast.NodeID(len(tree.Nodes))
	return n.store(tree, links)
}

func (n *Native) slice(element g.Fragment[s.Language]) (g.Fragment[s.Language], error) {
	return n.compose(s.KindSliceType, nil, fragmentPart{s.RoleElement, element})
}

// sameExpression compares this target's captured native syntax and exact
// symbolic identities. It establishes no equivalence between source types.
func (n *Native) sameExpression(a, b g.Fragment[s.Language]) bool {
	left, leftOK := n.fragmentRecords[a]
	right, rightOK := n.fragmentRecords[b]
	return leftOK && rightOK && left.tree.Root == right.tree.Root && reflect.DeepEqual(left.tree.Nodes, right.tree.Nodes) && maps.Equal(left.links, right.links)
}

func (n *Native) reference(callee g.Callee[SourceLanguage, s.Language], args ...g.Fragment[s.Language]) (g.Fragment[s.Language], error) {
	d, ok := n.destinationRecords[callee.Destination()]
	l, found := n.layoutRecords[callee.Layout()]
	if !ok || !found || l.destination != callee.Destination() || len(args) != len(l.parameters) {
		return g.Fragment[s.Language]{}, g.ErrScope
	}
	f, err := n.named(d.name, symbol{destination: callee.Destination()})
	if err != nil || len(args) == 0 {
		return f, err
	}
	parts := []fragmentPart{{s.RoleConstructor, f}}
	for _, arg := range args {
		parts = append(parts, fragmentPart{s.RoleArgument, arg})
	}
	return n.compose(s.KindTypeApplication, nil, parts...)
}

type nativeField struct {
	name  string
	value g.Fragment[s.Language]
}

func (n *Native) record(fields ...nativeField) (g.Fragment[s.Language], error) {
	var members []fragmentPart
	for _, field := range fields {
		f, err := n.compose(s.KindField, []ast.Attribute{{Name: s.AttributeName, Text: field.name}}, fragmentPart{s.RoleType, field.value})
		if err != nil {
			return g.Fragment[s.Language]{}, err
		}
		members = append(members, fragmentPart{s.RoleField, f})
	}
	return n.compose(s.KindStructType, nil, members...)
}

func (n *Native) Assemble(ctx context.Context, pubs []g.Publication[SourceLanguage, s.Language], defs []g.Definition[SourceLanguage, s.Language], limits g.Limits) (representation.ASTSet, error) {
	if err := ctx.Err(); err != nil {
		return representation.ASTSet{}, err
	}
	tree := s.AST{ID: "generated-fixture"}
	declarations := make(map[g.Destination[s.Language]]ast.NodeID)
	formals := make(map[g.Layout[s.Language]][]ast.NodeID)
	add := func(node s.Node) (ast.NodeID, error) {
		if len(tree.Nodes) >= limits.OutputNodes {
			return 0, g.ErrLimit
		}
		tree.Nodes = append(tree.Nodes, node)
		return ast.NodeID(len(tree.Nodes)), nil
	}
	for _, def := range defs {
		c := def.Callee()
		d, ok := n.destinationRecords[c.Destination()]
		l, found := n.layoutRecords[c.Layout()]
		if !ok || !found || !d.owned || declarations[c.Destination()] != 0 {
			return representation.ASTSet{}, g.ErrScope
		}
		kind := s.KindTypeDeclaration
		if d.alias {
			kind = s.KindTypeAlias
		}
		id, err := add(s.Node{Kind: kind, Attributes: []ast.Attribute{{Name: s.AttributeName, Text: d.name}}})
		if err != nil {
			return representation.ASTSet{}, err
		}
		declarations[c.Destination()] = id
		for _, name := range l.parameters {
			constraint, err := add(s.Node{Kind: s.KindInterfaceType})
			if err != nil {
				return representation.ASTSet{}, err
			}
			parameter, err := add(s.Node{Kind: s.KindTypeParameter, Attributes: []ast.Attribute{{Name: s.AttributeName, Text: name}}, Children: []ast.Child{{Role: s.RoleConstraint, Node: constraint}}})
			if err != nil {
				return representation.ASTSet{}, err
			}
			formals[c.Layout()] = append(formals[c.Layout()], parameter)
			tree.Nodes[id-1].Children = append(tree.Nodes[id-1].Children, ast.Child{Role: s.RoleTypeParameter, Node: parameter})
		}
	}
	for _, pub := range pubs {
		if declarations[pub.Callee().Destination()] == 0 {
			return representation.ASTSet{}, g.ErrIncomplete
		}
	}
	var roots []ast.Child
	for _, def := range defs {
		if err := n.CheckBody(ctx, def.Callee(), def.Body()); err != nil {
			return representation.ASTSet{}, err
		}
		r := n.fragmentRecords[def.Body()]
		if len(r.tree.Nodes) > limits.OutputNodes-len(tree.Nodes) {
			return representation.ASTSet{}, g.ErrLimit
		}
		offset := ast.NodeID(len(tree.Nodes))
		for i, node := range r.tree.Nodes {
			if err := ctx.Err(); err != nil {
				return representation.ASTSet{}, err
			}
			copy := cloneNode(node)
			for j := range copy.Children {
				copy.Children[j].Node += offset
			}
			if link, ok := r.links[ast.NodeID(i+1)]; ok {
				var target ast.NodeID
				if link.layout.Valid() {
					parameters := formals[link.layout]
					if link.formal >= len(parameters) {
						return representation.ASTSet{}, g.ErrScope
					}
					target = parameters[link.formal]
				} else if n.destinationRecords[link.destination].owned {
					target = declarations[link.destination]
					if target == 0 {
						return representation.ASTSet{}, g.ErrIncomplete
					}
				}
				if target != 0 {
					copy.References = []ast.Reference{{Role: s.RoleTarget, Target: ast.NodeRef{Tree: tree.ID, Node: target}}}
				}
			}
			if _, err := add(copy); err != nil {
				return representation.ASTSet{}, err
			}
		}
		id := declarations[def.Callee().Destination()]
		tree.Nodes[id-1].Children = append(tree.Nodes[id-1].Children, ast.Child{Role: s.RoleType, Node: r.tree.Root + offset})
		roots = append(roots, ast.Child{Role: s.RoleDeclaration, Node: id})
	}
	root, err := add(s.Node{Kind: s.KindFile, Attributes: []ast.Attribute{{Name: s.AttributePackage, Text: "fixture"}}, Children: roots})
	if err != nil {
		return representation.ASTSet{}, err
	}
	tree.Root = root
	if err := (gp.Provider{}).Check(tree); err != nil {
		return representation.ASTSet{}, err
	}
	return representation.PackSet(n.language, ast.Set[s.Language]{Entry: ast.NodeRef{Tree: tree.ID, Node: root}, Trees: []s.AST{tree}})
}
