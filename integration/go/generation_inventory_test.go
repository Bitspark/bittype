package integration_test

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	cs "github.com/Bitspark/bittype/syntax/cpp/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	rs "github.com/Bitspark/bittype/syntax/rust/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// These qualify the declaration inventory of a sealed output and the handles it
// issues (#272). No test in this file reads a name attribute: a declaration is
// selected by what its inventory record states.

// inventoryRoots need one dependency, Page, and one helper, Cursor.
var inventoryRoots = []string{"Person", "Outer", "Nested"}

func inventoryOutput(t *testing.T, fixture *f.Fixture, roots ...string) g.Output {
	t.Helper()
	plan, err := fixture.PlanRoots(context.Background(), f.Standard, roots...)
	if err != nil {
		t.Fatal(err)
	}
	out, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// childSites are the children of one node, optionally of one role.
func childSites[L any](tree ast.AST[L], node ast.NodeID, role string) []ast.NodeRef {
	var out []ast.NodeRef
	for _, child := range tree.Nodes[node-1].Children {
		if role == "" || child.Role == role {
			out = append(out, ast.NodeRef{Tree: tree.ID, Node: child.Node})
		}
	}
	return out
}

// refersTo reports whether a node below from refers to the node to.
func refersTo[L any](tree ast.AST[L], from, to ast.NodeRef) bool {
	for _, reference := range tree.Nodes[from.Node-1].References {
		if reference.Target == to {
			return true
		}
	}
	for _, child := range tree.Nodes[from.Node-1].Children {
		if refersTo(tree, ast.NodeRef{Tree: tree.ID, Node: child.Node}, to) {
			return true
		}
	}
	return false
}

func TestGenerationInventoryListsRequestedDependencyAndHelper(t *testing.T) {
	fixture := independent(t)
	out := inventoryOutput(t, fixture, inventoryRoots...)
	tree := outputTree[gs.Language](t, out)
	inventory := out.Inventory()
	if len(inventory) != 5 {
		t.Fatal("inventory size", len(inventory))
	}
	// Every record is one top-level declaration of the snapshot, and every
	// top-level declaration has its record.
	top := childSites(tree, tree.Root, "")
	if len(top) != 5 {
		t.Fatal("top-level declarations", len(top))
	}
	roles := map[g.DeclarationRole]int{}
	bySource := map[ast.NodeRef]g.OwnedDeclaration{}
	var helper g.OwnedDeclaration
	for i, record := range inventory {
		if record.Site != top[i] {
			t.Fatal("record is not the top-level declaration at its position", i, record.Site)
		}
		resolved, err := out.Resolve(record.Handle)
		if err != nil || !reflect.DeepEqual(resolved, record) {
			t.Fatal("handle does not resolve to its record", i, err)
		}
		roles[record.Role]++
		if record.Role == g.HelperDeclaration {
			if len(record.Provenance) != 0 {
				t.Fatal("a helper has a source declaration", record.Provenance)
			}
			helper = record
			continue
		}
		if len(record.Provenance) != 1 {
			t.Fatal("provenance of a source declaration", i, record.Provenance)
		}
		if _, twice := bySource[record.Provenance[0]]; twice {
			t.Fatal("two records for one source declaration", record.Provenance[0])
		}
		bySource[record.Provenance[0]] = record
	}
	if roles[g.RequestedDeclaration] != 3 || roles[g.DependencyDeclaration] != 1 || roles[g.HelperDeclaration] != 1 {
		t.Fatal("roles", roles)
	}
	// Each requested root is selected by the source site the request named. Its
	// record is the exposed site of that root and carries what the fixture's
	// recipe stated.
	exposed := out.Exposed()
	if len(exposed) != len(inventoryRoots) {
		t.Fatal("exposed sites", exposed)
	}
	for i, name := range inventoryRoots {
		record, ok := bySource[fixture.Site(name)]
		if !ok || record.Role != g.RequestedDeclaration || record.Site != exposed[i] {
			t.Fatal("requested root is not selected by its source site", name)
		}
		if intent, ok := f.Intent(record); !ok || intent != name {
			t.Fatal("the recipe's value did not reach the record", name, intent, ok)
		}
		if _, ok := f.Access(record); ok {
			t.Fatal("a value of the declaration recipe read as the helper's", name)
		}
	}
	page, ok := bySource[fixture.Site("Page")]
	if !ok || page.Role != g.DependencyDeclaration || slices.Contains(exposed, page.Site) {
		t.Fatal("the dependency is not selected by its source site")
	}
	if intent, ok := f.Intent(page); !ok || intent != "Page" {
		t.Fatal("the recipe's value did not reach the dependency", intent, ok)
	}
	if access, ok := f.Access(helper); !ok || access != f.HelperAccess {
		t.Fatal("the helper recipe's value did not reach the helper", access, ok)
	}
	if _, ok := f.Intent(helper); ok || len(helper.Values) != 1 || slices.Contains(exposed, helper.Site) {
		t.Fatal("the helper record", helper)
	}
	// The selected records are the declarations the tree links: both aliases
	// use the dependency, the dependency uses the helper.
	outer, nested, person := bySource[fixture.Site("Outer")], bySource[fixture.Site("Nested")], bySource[fixture.Site("Person")]
	if !refersTo(tree, outer.Site, page.Site) || !refersTo(tree, nested.Site, page.Site) || !refersTo(tree, page.Site, helper.Site) {
		t.Fatal("the selected records are not the linked declarations")
	}
	if refersTo(tree, person.Site, page.Site) || refersTo(tree, person.Site, helper.Site) || refersTo(tree, helper.Site, page.Site) {
		t.Fatal("an unrelated declaration links the dependency or the helper")
	}
}

// The inventory follows the sites, and the sites do not follow the order of the
// requested roots. Only Exposed keeps the order of the request.
func TestGenerationInventoryDoesNotFollowTheOrderOfTheRoots(t *testing.T) {
	fixture := independent(t)
	first := inventoryOutput(t, fixture, "Person", "Outer", "Nested")
	second := inventoryOutput(t, fixture, "Nested", "Person", "Outer")
	if !reflect.DeepEqual(outputTree[gs.Language](t, first), outputTree[gs.Language](t, second)) {
		t.Fatal("the order of the roots changed the tree")
	}
	a, b := first.Inventory(), second.Inventory()
	if len(a) != 5 || len(b) != 5 {
		t.Fatal("inventory sizes", len(a), len(b))
	}
	for i := range a {
		a[i].Handle, b[i].Handle = g.DeclarationHandle{}, g.DeclarationHandle{}
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("the order of the roots changed the inventory")
	}
	x, y := first.Exposed(), second.Exposed()
	if x[0] != y[1] || x[1] != y[2] || x[2] != y[0] {
		t.Fatal("exposed sites do not follow the request", x, y)
	}
}

// One plan run twice gives two outputs with equal sites. A handle of one is
// refused by the other.
func TestGenerationHandleIsRefusedByAnotherRunOfThePlan(t *testing.T) {
	fixture := independent(t)
	plan, err := fixture.PlanRoots(context.Background(), f.Standard, inventoryRoots...)
	if err != nil {
		t.Fatal(err)
	}
	first, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	second, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	a, b := first.Inventory(), second.Inventory()
	if len(a) != 5 || len(b) != 5 {
		t.Fatal("inventory sizes", len(a), len(b))
	}
	for i := range a {
		// Tree ID and node number are equal in both runs, and so is everything
		// else the record states.
		if a[i].Site != b[i].Site || a[i].Role != b[i].Role || !slices.Equal(a[i].Provenance, b[i].Provenance) || !slices.Equal(a[i].Values, b[i].Values) {
			t.Fatal("the two runs differ in more than their identity", i)
		}
		if a[i].Handle == b[i].Handle {
			t.Fatal("two runs issued one handle", i)
		}
		if record, err := second.Resolve(a[i].Handle); !errors.Is(err, g.ErrScope) || record.Handle.Valid() {
			t.Fatal("a handle of the first run resolved against the second", i, err)
		}
		if _, err := first.Resolve(b[i].Handle); !errors.Is(err, g.ErrScope) {
			t.Fatal("a handle of the second run resolved against the first", i, err)
		}
		if record, err := first.Resolve(a[i].Handle); err != nil || record.Site != a[i].Site {
			t.Fatal("a handle did not resolve against its own run", i, err)
		}
		if record, err := second.Resolve(b[i].Handle); err != nil || record.Site != b[i].Site {
			t.Fatal("a handle did not resolve against its own run", i, err)
		}
	}
}

// The kernel seals what the target built only if the assembly lists every
// planned declaration once, inside the snapshot.
func TestGenerationRunRefusesAnAssemblyWithAnIncompleteInventory(t *testing.T) {
	fixture := independent(t)
	standard := inventoryOutput(t, fixture, inventoryRoots...).Inventory()
	for i := range standard {
		standard[i].Handle = g.DeclarationHandle{}
	}
	for name, c := range map[string]struct {
		tamper func([]ast.NodeRef) []ast.NodeRef
		want   error
	}{
		"complete":                {func(s []ast.NodeRef) []ast.NodeRef { return s }, nil},
		"one declaration omitted": {func(s []ast.NodeRef) []ast.NodeRef { return s[:len(s)-1] }, g.ErrIncomplete},
		"one record too many":     {func(s []ast.NodeRef) []ast.NodeRef { return append(s, ast.NodeRef{Tree: s[0].Tree, Node: 1}) }, g.ErrIncomplete},
		"site in another tree":    {func(s []ast.NodeRef) []ast.NodeRef { s[0].Tree = "another"; return s }, g.ErrScope},
		"site listed twice":       {func(s []ast.NodeRef) []ast.NodeRef { s[1] = s[0]; return s }, g.ErrScope},
	} {
		t.Run(name, func(t *testing.T) {
			built := 0
			plan, err := fixture.PlanBuilt(context.Background(), func(a g.Assembly) g.Assembly {
				built = len(a.Declarations)
				a.Declarations = c.tamper(slices.Clone(a.Declarations))
				return a
			}, inventoryRoots...)
			if err != nil {
				t.Fatal(err)
			}
			out, err := plan.Run(context.Background())
			if built != 5 {
				t.Fatal("the target listed another number of declarations", built)
			}
			if c.want == nil {
				if err != nil {
					t.Fatal("the complete assembly was not sealed", err)
				}
				// The plan is run a second time. What the target overwrote
				// in its copy of the program is in neither inventory.
				again, err := plan.Run(context.Background())
				if err != nil {
					t.Fatal(err)
				}
				for _, run := range []g.Output{out, again} {
					inventory := run.Inventory()
					for i := range inventory {
						inventory[i].Handle = g.DeclarationHandle{}
					}
					if !reflect.DeepEqual(inventory, standard) {
						t.Fatal("the target changed what the plan states about its declarations")
					}
				}
				return
			}
			if !errors.Is(err, c.want) || out.Snapshot().Valid() || len(out.Inventory()) != 0 {
				t.Fatal("an incomplete inventory was sealed", err)
			}
		})
	}
}

// A value in a description must come from an issuer, and its content counts
// against the output budget. Both are refused when the plan is made.
func TestGenerationPlanRefusesUnissuedAndOversizedValues(t *testing.T) {
	fixture := independent(t)
	if err := independentFails(t, fixture, f.UnissuedValue, "Person"); !errors.Is(err, g.ErrInvalid) {
		t.Fatal("a declaration recipe stated a value nobody issued", err)
	}
	if err := independentFails(t, fixture, f.UnissuedHelperValue, "Page"); !errors.Is(err, g.ErrInvalid) {
		t.Fatal("a helper recipe stated a value nobody issued", err)
	}
	limits := g.DefaultLimits()
	limits.OutputBytes = f.LargeValueBytes
	plan, err := fixture.PlanRootsWithin(context.Background(), f.Standard, limits, "Person")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = plan.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if plan, err = fixture.PlanRootsWithin(context.Background(), f.LargeValue, limits, "Person"); !errors.Is(err, g.ErrLimit) || plan.Valid() {
		t.Fatal("a value larger than the output budget was planned", err)
	}
}

// The shared builder lists the declarations once for all four targets. Every
// root of the pagination module is requested, so no record is a dependency.
// A target without a helper recipe states nothing, and its records carry no
// value.
func TestGenerationInventoryCoversEveryDeclarationInFourTargets(t *testing.T) {
	e, _, input := generationInput(t, paginationModule())
	roots, err := genfixture.Roots(input)
	if err != nil {
		t.Fatal(err)
	}
	check := func(t *testing.T, out g.Output, declarations []ast.NodeRef, helpers int) {
		t.Helper()
		inventory := out.Inventory()
		if len(declarations) != len(roots)+helpers || len(inventory) != len(declarations) {
			t.Fatal("inventory size", len(inventory), "declarations", len(declarations))
		}
		var requested, sources []ast.NodeRef
		for i, record := range inventory {
			if record.Site != declarations[i] {
				t.Fatal("record is not the declaration at its position", i, record.Site)
			}
			if resolved, err := out.Resolve(record.Handle); err != nil || resolved.Site != record.Site {
				t.Fatal("handle does not resolve to its record", i, err)
			}
			switch record.Role {
			case g.RequestedDeclaration:
				requested = append(requested, record.Site)
				sources = append(sources, record.Provenance...)
			case g.HelperDeclaration:
				helpers--
				if len(record.Provenance) != 0 {
					t.Fatal("a helper has a source declaration", record.Provenance)
				}
			default:
				t.Fatal("unexpected role", record.Role)
			}
			if len(record.Values) != 0 {
				t.Fatal("a record carries a value no owner stated", record.Values)
			}
		}
		if helpers != 0 {
			t.Fatal("helper records", helpers)
		}
		order := func(a, b ast.NodeRef) int { return int(a.Node) - int(b.Node) }
		exposed, wanted := out.Exposed(), slices.Clone(roots)
		slices.SortFunc(exposed, order)
		slices.SortFunc(sources, order)
		slices.SortFunc(wanted, order)
		if !slices.Equal(requested, exposed) {
			t.Fatal("requested records differ from the exposed sites", requested, exposed)
		}
		if !slices.Equal(sources, wanted) {
			t.Fatal("provenance differs from the requested roots", sources, wanted)
		}
	}
	t.Run("typescript", func(t *testing.T) {
		out := generatedOutput(t, e, input, typescriptText(), tsRequest)
		tree := outputTree[ts.Language](t, out)
		check(t, out, childSites(tree, tree.Root, ts.RoleDeclaration), 0)
	})
	t.Run("golang", func(t *testing.T) {
		out := generatedOutput(t, e, input, golangText(), goRequest)
		tree := outputTree[gs.Language](t, out)
		check(t, out, childSites(tree, tree.Root, gs.RoleDeclaration), 1)
	})
	t.Run("rust", func(t *testing.T) {
		out := generatedOutput(t, e, input, rustText(), rustRequest)
		tree := outputTree[rs.Language](t, out)
		check(t, out, childSites(tree, tree.Root, rs.RoleItem), 0)
	})
	t.Run("cpp", func(t *testing.T) {
		out := generatedOutput(t, e, input, cppText(), cppRequest)
		tree := outputTree[cs.Language](t, out)
		// C++ declarations are the children of the one namespace of the unit.
		var declarations []ast.NodeRef
		for _, site := range childSites(tree, tree.Root, cs.RoleDeclaration) {
			if tree.Nodes[site.Node-1].Kind == cs.KindNamespace {
				declarations = append(declarations, childSites(tree, site.Node, cs.RoleDeclaration)...)
			}
		}
		check(t, out, declarations, 0)
	})
}
