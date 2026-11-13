package generation

import (
	"cmp"
	"context"
	"fmt"
	"io"
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

// ValueIssuer is an owner's authority to state something about a declaration
// it describes and to read that statement back from the sealed output. A naming
// intent and the access a helper requires are such statements. The generation
// kernel carries an issued value from the description to the inventory record
// of its declaration. It stores the value and hands it back, and it never reads
// the content. One issuer is one meaning: an owner with two kinds of statement
// uses two issuers. The content may be empty, where the issuer alone is the
// statement.
//
// Limit: a value is kept from other readers by this API only, as an external
// symbol is. Reflection reads the content. A value is not bound to one
// declaration or one output: the same value compares equal wherever it is
// carried, and whoever holds a value can state it again in another
// description. The issuer says who gave a value its meaning. It does not say
// who put the value where it is found.
type ValueIssuer struct{ owner *identity }
type IssuedValue struct {
	owner   *identity
	content string
}

// String says whether the value was issued. It never contains the content.
func (v IssuedValue) String() string {
	if v.owner == nil {
		return "unissued value"
	}
	return "issued value"
}

// Format writes String for every verb, so that formatting a value, or a record
// or a description that holds one, does not print the content.
func (v IssuedValue) Format(state fmt.State, _ rune) { io.WriteString(state, v.String()) }

func NewValueIssuer() ValueIssuer { return ValueIssuer{fresh()} }
func (i ValueIssuer) Issue(content string) (IssuedValue, error) {
	if i.owner == nil {
		return IssuedValue{}, ErrInvalid
	}
	return IssuedValue{i.owner, content}, nil
}

// Read returns the content of a value this issuer issued, and false for every
// other value.
func (i ValueIssuer) Read(v IssuedValue) (string, bool) {
	if i.owner == nil || i.owner != v.owner {
		return "", false
	}
	return v.content, true
}

// DeclarationRole says why a sealed output owns a declaration. It is a fact of
// the plan, and no target or mapping states it. The values below keep their
// meaning. They are not a closed set: a further role may be added, for example
// for a declaration that a contribution introduces and that has a source
// declaration. A reader therefore refuses a record whose role it does not know,
// and does not treat it as one of the roles it knows.
type DeclarationRole uint8

const (
	// RequestedDeclaration represents a source declaration that the request
	// lists as a root.
	RequestedDeclaration DeclarationRole = iota + 1
	// DependencyDeclaration represents a source declaration that a planned
	// declaration uses and the request does not list.
	DependencyDeclaration
	// HelperDeclaration is demanded by a mapping through a helper key. It
	// represents no source declaration.
	HelperDeclaration
)

// DeclarationHandle identifies one declaration of one sealed output. Only seal
// issues handles. A handle carries its output, as an external symbol carries
// its issuer, so the handles of two outputs differ even when their sites are
// equal, which is the case for two runs of one plan. Handles are comparable.
type DeclarationHandle struct {
	output *identity
	index  uint32
}

func (h DeclarationHandle) Valid() bool { return h.output != nil && h.index != 0 }

// OwnedDeclaration is the inventory record of one declaration that a sealed
// output owns. Site is the node that the target's construction returned for
// the declaration. Where that node sits in its unit is the target's decision:
// it need not be a child of the root. Role and Provenance are the plan's facts. Provenance lists the declaring sites of the
// source declarations it represents: none for a helper, and possibly those of
// several source declarations. Values are the statements of the declaration's
// owner, which only their issuers can read. A record says what the output
// owns. It does not say that its site is one entity of a native compiler.
type OwnedDeclaration struct {
	Handle     DeclarationHandle
	Site       ast.NodeRef
	Role       DeclarationRole
	Provenance []ast.NodeRef
	Values     []IssuedValue
}

func (d OwnedDeclaration) clone() OwnedDeclaration {
	d.Provenance = slices.Clone(d.Provenance)
	d.Values = slices.Clone(d.Values)
	return d
}

// Assembly is a target construction record. Declarations contains the realized
// site of every planned declaration, in the order of Program.Declarations. No
// listed site may contain another. The kernel cannot tell which declaration a
// site realizes: a target that lists the sites of two declarations in each
// other's place is not detected.
type Assembly struct {
	Syntax       representation.ASTSet
	Declarations []ast.NodeRef
	ExternalUses []ExternalUse
}

// Output can only be produced by a successfully finished generation session.
// Sites and immutable target symbols always travel with their own checked capture.
type Output struct {
	snapshot  semantic.Snapshot
	owner     *identity
	inventory []OwnedDeclaration
	exposed   []ast.NodeRef
	external  []ExternalUse
}

func (o Output) Snapshot() semantic.Snapshot { return o.snapshot }
func (o Output) Exposed() []ast.NodeRef      { return slices.Clone(o.exposed) }
func (o Output) ExternalUses() []ExternalUse { return slices.Clone(o.external) }

// Inventory returns one record for every declaration the output owns and no
// other record. The records are in the order of their sites: the trees in the
// snapshot's order, and the nodes of one tree ascending. No site contains
// another. The records with the role RequestedDeclaration have exactly the
// sites of Exposed.
func (o Output) Inventory() []OwnedDeclaration {
	out := make([]OwnedDeclaration, len(o.inventory))
	for i, d := range o.inventory {
		out[i] = d.clone()
	}
	return out
}

// Resolve returns the record a handle identifies. It refuses a handle that
// this output did not issue with ErrScope. That includes the handle of another
// run of the same plan.
func (o Output) Resolve(h DeclarationHandle) (OwnedDeclaration, error) {
	if h.output != o.owner || h.index == 0 || int(h.index) > len(o.inventory) {
		return OwnedDeclaration{}, ErrScope
	}
	return o.inventory[h.index-1].clone(), nil
}

// planned is what a plan states about one of its declarations, without the
// native type of the program.
type planned struct {
	role       DeclarationRole
	provenance []ast.NodeRef
	values     []IssuedValue
}

// seal checks an assembly against the plan it was built from. declarations and
// exposed are the plan's own: one entry per planned declaration in planned
// order, and the planned IDs of the requested roots in request order. T is the
// target language of the plan, in which the assembly's trees are read for
// their containment.
func seal[T any](ctx context.Context, checker check.Checker, assembly Assembly, declarations []planned, exposed []uint32, limits Limits) (Output, error) {
	set := assembly.Syntax
	if len(set.Trees) == 0 || len(set.Trees) > 256 || len(assembly.Declarations) != len(declarations) {
		return Output{}, ErrIncomplete
	}
	if len(assembly.ExternalUses) > limits.References {
		return Output{}, ErrLimit
	}
	// Freeze target-owned slices before invoking the independently supplied checker.
	set.Trees = slices.Clone(set.Trees)
	assembly.ExternalUses = slices.Clone(assembly.ExternalUses)
	trees := map[ast.TreeID]int{}
	position := map[ast.TreeID]int{}
	nodes, bytes := limits.OutputNodes, limits.OutputBytes
	for i, tree := range set.Trees {
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
		position[info.ID] = i
	}
	contains := func(ref ast.NodeRef) bool { return ref.Node != 0 && uint64(ref.Node) <= uint64(trees[ref.Tree]) }
	// The inventory is complete by count. Each site is a node of this assembly
	// that no other record claims.
	inventory := make([]OwnedDeclaration, len(declarations))
	sites := map[ast.NodeRef]bool{}
	for i, site := range assembly.Declarations {
		if !contains(site) || sites[site] {
			return Output{}, ErrScope
		}
		sites[site] = true
		d := declarations[i]
		inventory[i] = OwnedDeclaration{Site: site, Role: d.role, Provenance: slices.Clone(d.provenance), Values: slices.Clone(d.values)}
	}
	// No listed site contains another. This needs no knowledge of the language:
	// whatever a declaration is, it is not inside another declaration of the
	// same inventory. The parent of each node is read from the trees as packed.
	for _, tree := range set.Trees {
		typed, err := representation.Unpack[T](checker.Language(), tree)
		if err != nil {
			return Output{}, err
		}
		parent := make([]ast.NodeID, len(typed.Nodes)+1)
		for i, node := range typed.Nodes {
			for _, child := range node.Children {
				parent[child.Node] = ast.NodeID(i + 1)
			}
		}
		for site := range sites {
			if site.Tree != typed.ID {
				continue
			}
			for above := parent[site.Node]; above != 0; above = parent[above] {
				if sites[ast.NodeRef{Tree: site.Tree, Node: above}] {
					return Output{}, ErrScope
				}
			}
		}
	}
	// A requested root is exposed at the site of its own record.
	requested := make([]ast.NodeRef, len(exposed))
	for i, id := range exposed {
		requested[i] = assembly.Declarations[id-1]
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
	slices.SortFunc(inventory, func(a, b OwnedDeclaration) int {
		return cmp.Or(cmp.Compare(position[a.Site.Tree], position[b.Site.Tree]), cmp.Compare(a.Site.Node, b.Site.Node))
	})
	owner := fresh()
	for i := range inventory {
		inventory[i].Handle = DeclarationHandle{owner, uint32(i + 1)}
	}
	return Output{snapshot, owner, inventory, requested, assembly.ExternalUses}, nil
}
