package generation

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type outputLanguage struct{}

func (outputLanguage) ID() language.ID { return "output-test" }
func (outputLanguage) Name() string    { return "Output test" }

type outputChecker struct{ outputLanguage }

func (outputChecker) Check(ast.View[outputLanguage]) error { return nil }

// sealTest seals an assembly of the test language.
func sealTest(checker check.Checker, assembly Assembly, declarations []planned, exposed []uint32, limits Limits) (Output, error) {
	return seal[outputLanguage](context.Background(), checker, assembly, declarations, exposed, limits)
}

// outputChecking is the test language and its checker.
func outputChecking(t *testing.T) (representation.Language, check.Checker) {
	t.Helper()
	lang, err := representation.NewLanguage[outputLanguage](outputLanguage{}, "one")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[outputLanguage](lang, outputChecker{})
	if err != nil {
		t.Fatal(err)
	}
	return lang, checker
}

// outputNode is one node of a test tree with its children in order.
func outputNode(kind string, children ...ast.NodeID) ast.Node[outputLanguage] {
	node := ast.Node[outputLanguage]{Kind: ast.Kind[outputLanguage](kind)}
	for _, child := range children {
		node.Children = append(node.Children, ast.Child{Role: "child", Node: child})
	}
	return node
}

// outputAssembly is one unit with one requested declaration at node 2.
func outputAssembly(t *testing.T) (check.Checker, Assembly) {
	t.Helper()
	checker, assembly := inventoryAssembly(t, 1)
	return checker, assembly
}

// onePlanned is the plan of outputAssembly: its one declaration is requested.
func onePlanned() ([]planned, []uint32) {
	return []planned{{role: RequestedDeclaration}}, []uint32{1}
}

// inventoryAssembly is one unit with count declarations at the nodes 2 to
// count+1, each with one child, and an assembly that lists every declaration.
func inventoryAssembly(t *testing.T, count int) (check.Checker, Assembly) {
	t.Helper()
	lang, err := representation.NewLanguage[outputLanguage](outputLanguage{}, "one")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[outputLanguage](lang, outputChecker{})
	if err != nil {
		t.Fatal(err)
	}
	tree := ast.AST[outputLanguage]{ID: "output", Root: 1, Nodes: []ast.Node[outputLanguage]{{Kind: "unit"}}}
	var sites []ast.NodeRef
	for i := range count {
		node := ast.NodeID(i + 2)
		tree.Nodes[0].Children = append(tree.Nodes[0].Children, ast.Child{Role: "declaration", Node: node})
		tree.Nodes = append(tree.Nodes, ast.Node[outputLanguage]{Kind: "declaration"})
		sites = append(sites, ast.NodeRef{Tree: tree.ID, Node: node})
	}
	for i := range count {
		tree.Nodes[i+1].Children = []ast.Child{{Role: "body", Node: ast.NodeID(len(tree.Nodes) + 1)}}
		tree.Nodes = append(tree.Nodes, ast.Node[outputLanguage]{Kind: "body"})
	}
	packed, err := representation.PackSet(lang, ast.Set[outputLanguage]{Entry: ast.NodeRef{Tree: tree.ID, Node: 1}, Trees: []ast.AST[outputLanguage]{tree}})
	if err != nil {
		t.Fatal(err)
	}
	return checker, Assembly{Syntax: packed, Declarations: sites}
}
func TestSealOwnsItsSitesAndRejectsIncompleteMetadata(t *testing.T) {
	checker, assembly := outputAssembly(t)
	declarations, exposed := onePlanned()
	out, err := sealTest(checker, assembly, declarations, exposed, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	assembly.Declarations[0] = ast.NodeRef{Tree: "foreign", Node: 1}
	sites := out.Exposed()
	sites[0] = ast.NodeRef{}
	if out.Exposed()[0].Tree != "output" {
		t.Fatal("mutable exposure metadata")
	}
	if bad, err := sealTest(checker, assembly, declarations, exposed, DefaultLimits()); !errors.Is(err, ErrScope) || bad.Snapshot().Valid() {
		t.Fatal(err)
	}
	_, assembly = outputAssembly(t)
	two := []planned{{role: RequestedDeclaration}, {role: RequestedDeclaration}}
	if bad, err := sealTest(checker, assembly, two, []uint32{1, 2}, DefaultLimits()); !errors.Is(err, ErrIncomplete) || bad.Snapshot().Valid() {
		t.Fatal(err)
	}
}
func TestSealRejectsForeignDuplicateAndOversizedExternalUses(t *testing.T) {
	checker, base := outputAssembly(t)
	symbol, err := NewSymbolIssuer().Issue(checker.Language(), "module/type")
	if err != nil {
		t.Fatal(err)
	}
	at := semantics.Occurrence{Site: base.Declarations[0], Role: semantics.Key{Owner: "fixture", Name: "type"}}
	declarations, exposed := onePlanned()
	for name, uses := range map[string][]ExternalUse{
		"foreign":   {{At: semantics.Occurrence{Site: ast.NodeRef{Tree: "other", Node: 2}, Role: at.Role}, Symbol: symbol}},
		"duplicate": {{At: at, Symbol: symbol}, {At: at, Symbol: symbol}},
		"unissued":  {{At: at}},
	} {
		t.Run(name, func(t *testing.T) {
			a := base
			a.ExternalUses = uses
			if out, err := sealTest(checker, a, declarations, exposed, DefaultLimits()); err == nil || out.Snapshot().Valid() {
				t.Fatal("invalid external use accepted")
			}
		})
	}
	base.ExternalUses = []ExternalUse{{At: at, Symbol: symbol}}
	limits := DefaultLimits()
	limits.OutputBytes = 2
	if out, err := sealTest(checker, base, declarations, exposed, limits); !errors.Is(err, ErrLimit) || out.Snapshot().Valid() {
		t.Fatal(err)
	}
}

// fivePlanned is a plan of five declarations whose planned order differs from
// the order of their sites: planned declaration i+1 is realized at the site
// realized[i] of inventoryAssembly(t, 5). The request lists the planned
// declarations 5, 1 and 4, in that order.
func fivePlanned(t *testing.T, assembly Assembly) ([]planned, []uint32, Assembly, ValueIssuer) {
	t.Helper()
	issuer := NewValueIssuer()
	intent := func(content string) []IssuedValue {
		value, err := issuer.Issue(content)
		if err != nil {
			t.Fatal(err)
		}
		return []IssuedValue{value}
	}
	source := func(nodes ...ast.NodeID) []ast.NodeRef {
		var out []ast.NodeRef
		for _, node := range nodes {
			out = append(out, ast.NodeRef{Tree: "source", Node: node})
		}
		return out
	}
	declarations := []planned{
		{RequestedDeclaration, source(11), intent("first")},
		{HelperDeclaration, nil, intent("")},
		{DependencyDeclaration, source(13, 14), nil},
		{RequestedDeclaration, source(15), intent("fourth")},
		{RequestedDeclaration, source(16), intent("fifth\x00\xff")},
	}
	sites := assembly.Declarations
	realized := []int{2, 0, 4, 1, 3}
	assembly.Declarations = nil
	for _, site := range realized {
		assembly.Declarations = append(assembly.Declarations, sites[site])
	}
	return declarations, []uint32{5, 1, 4}, assembly, issuer
}

func TestSealedInventoryListsEveryPlannedDeclaration(t *testing.T) {
	checker, base := inventoryAssembly(t, 5)
	declarations, exposed, assembly, issuer := fivePlanned(t, base)
	out, err := sealTest(checker, assembly, declarations, exposed, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	inventory := out.Inventory()
	if len(inventory) != 5 {
		t.Fatal("inventory size", len(inventory))
	}
	// One record per planned declaration, in the order of the sites, with what
	// the plan states about the declaration realized at that site.
	at := []int{1, 3, 0, 4, 2}
	handles := map[DeclarationHandle]bool{}
	var requested []ast.NodeRef
	for i, record := range inventory {
		want := declarations[at[i]]
		if record.Site != base.Declarations[i] || record.Role != want.role || !slices.Equal(record.Provenance, want.provenance) || !slices.Equal(record.Values, want.values) {
			t.Fatalf("record %d: %+v", i, record)
		}
		if !record.Handle.Valid() || handles[record.Handle] {
			t.Fatal("handle is not one per record", i)
		}
		handles[record.Handle] = true
		resolved, err := out.Resolve(record.Handle)
		if err != nil || !reflect.DeepEqual(resolved, record) {
			t.Fatal("handle does not resolve to its record", i, err)
		}
		if record.Role == RequestedDeclaration {
			requested = append(requested, record.Site)
		}
	}
	// The requested records are exactly the exposed sites. Exposed keeps the
	// order of the request.
	exposure := out.Exposed()
	if !slices.Equal(exposure, []ast.NodeRef{base.Declarations[3], base.Declarations[2], base.Declarations[1]}) {
		t.Fatal("exposed sites", exposure)
	}
	slices.SortFunc(exposure, func(a, b ast.NodeRef) int { return int(a.Node) - int(b.Node) })
	if !slices.Equal(requested, exposure) {
		t.Fatal("requested records differ from the exposed sites", requested, exposure)
	}
	// The kernel hands a value back to its issuer as it was issued.
	if content, ok := issuer.Read(inventory[3].Values[0]); !ok || content != "fifth\x00\xff" {
		t.Fatal("value changed in the inventory", content, ok)
	}
	if content, ok := issuer.Read(inventory[0].Values[0]); !ok || content != "" {
		t.Fatal("empty value changed in the inventory", content, ok)
	}
	// Neither a returned record nor the plan's own slices reach the output.
	inventory[2].Provenance[0] = ast.NodeRef{}
	inventory[2].Values[0] = IssuedValue{}
	inventory[0] = OwnedDeclaration{}
	declarations[0].provenance[0] = ast.NodeRef{}
	declarations[0].values[0] = IssuedValue{}
	again := out.Inventory()
	if again[2].Provenance[0] != (ast.NodeRef{Tree: "source", Node: 11}) || again[2].Values[0] == (IssuedValue{}) || again[0].Site != base.Declarations[0] {
		t.Fatal("mutable inventory")
	}
	resolved, err := out.Resolve(again[2].Handle)
	if err != nil {
		t.Fatal(err)
	}
	resolved.Provenance[0] = ast.NodeRef{}
	resolved.Values[0] = IssuedValue{}
	if again = out.Inventory(); again[2].Provenance[0] != (ast.NodeRef{Tree: "source", Node: 11}) || again[2].Values[0] == (IssuedValue{}) {
		t.Fatal("mutable resolved record")
	}
}

func TestSealRefusesAnIncompleteInventory(t *testing.T) {
	checker, base := inventoryAssembly(t, 5)
	declarations, exposed, complete, _ := fivePlanned(t, base)
	body := ast.NodeRef{Tree: "output", Node: 7}
	for name, c := range map[string]struct {
		sites func([]ast.NodeRef) []ast.NodeRef
		want  error
	}{
		"one declaration omitted": {func(s []ast.NodeRef) []ast.NodeRef { return s[:4] }, ErrIncomplete},
		"no declaration listed":   {func(s []ast.NodeRef) []ast.NodeRef { return nil }, ErrIncomplete},
		"one record too many":     {func(s []ast.NodeRef) []ast.NodeRef { return append(s, body) }, ErrIncomplete},
		"site in another tree":    {func(s []ast.NodeRef) []ast.NodeRef { s[1].Tree = "other"; return s }, ErrScope},
		"site past the last node": {func(s []ast.NodeRef) []ast.NodeRef { s[1].Node = 12; return s }, ErrScope},
		"site without a node":     {func(s []ast.NodeRef) []ast.NodeRef { s[1].Node = 0; return s }, ErrScope},
		"site listed twice":       {func(s []ast.NodeRef) []ast.NodeRef { s[4] = s[0]; return s }, ErrScope},
	} {
		t.Run(name, func(t *testing.T) {
			assembly := complete
			assembly.Declarations = c.sites(slices.Clone(complete.Declarations))
			out, err := sealTest(checker, assembly, declarations, exposed, DefaultLimits())
			if !errors.Is(err, c.want) || out.Snapshot().Valid() || len(out.Inventory()) != 0 || len(out.Exposed()) != 0 {
				t.Fatal("incomplete inventory sealed", err)
			}
		})
	}
	// The last node of the tree is inside the snapshot. It is the body of the
	// declaration at node 6, whose record it replaces here.
	assembly := complete
	assembly.Declarations = slices.Clone(complete.Declarations)
	assembly.Declarations[2] = ast.NodeRef{Tree: "output", Node: 11}
	if _, err := sealTest(checker, assembly, declarations, exposed, DefaultLimits()); err != nil {
		t.Fatal(err)
	}
}

// A handle carries its output. Two seals of one assembly and one plan have
// equal sites, and neither resolves a handle of the other.
func TestHandleResolvesOnlyAgainstItsOutput(t *testing.T) {
	checker, base := inventoryAssembly(t, 5)
	declarations, exposed, assembly, _ := fivePlanned(t, base)
	first, err := sealTest(checker, assembly, declarations, exposed, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	second, err := sealTest(checker, assembly, declarations, exposed, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	a, b := first.Inventory(), second.Inventory()
	for i := range a {
		if a[i].Site != b[i].Site || a[i].Role != b[i].Role {
			t.Fatal("the two outputs differ in more than their identity", i)
		}
		if a[i].Handle == b[i].Handle {
			t.Fatal("two outputs issued one handle", i)
		}
		if record, err := second.Resolve(a[i].Handle); !errors.Is(err, ErrScope) || record.Handle.Valid() || record.Site != (ast.NodeRef{}) {
			t.Fatal("a handle of the first output resolved against the second", i, err)
		}
		if _, err := first.Resolve(b[i].Handle); !errors.Is(err, ErrScope) {
			t.Fatal("a handle of the second output resolved against the first", i, err)
		}
		if record, err := first.Resolve(a[i].Handle); err != nil || record.Site != a[i].Site {
			t.Fatal("a handle did not resolve against its own output", i, err)
		}
	}
	// A copy of an output is the same output.
	copied := first
	if _, err := copied.Resolve(a[0].Handle); err != nil {
		t.Fatal(err)
	}
	for name, handle := range map[string]DeclarationHandle{
		"zero":             {},
		"no output":        {index: 1},
		"no declaration":   {output: first.owner},
		"past the records": {output: first.owner, index: 6},
	} {
		if handle.Valid() != (name == "past the records") {
			t.Fatal("validity of the handle:", name)
		}
		if _, err := first.Resolve(handle); !errors.Is(err, ErrScope) {
			t.Fatal("resolved a handle the output did not issue:", name, err)
		}
	}
	if _, err := (Output{}).Resolve(a[0].Handle); !errors.Is(err, ErrScope) {
		t.Fatal("a failed output resolved a handle", err)
	}
	if _, err := (Output{}).Resolve(DeclarationHandle{}); !errors.Is(err, ErrScope) {
		t.Fatal("a failed output resolved the zero handle", err)
	}
}

func TestIssuedValueIsReadOnlyByItsIssuer(t *testing.T) {
	owner, other := NewValueIssuer(), NewValueIssuer()
	for _, content := range []string{"", "Person", "\x00\xff not text"} {
		value, err := owner.Issue(content)
		if err != nil {
			t.Fatal(err)
		}
		if got, ok := owner.Read(value); !ok || got != content {
			t.Fatal("the issuer cannot read its value", got, ok)
		}
		if got, ok := other.Read(value); ok || got != "" {
			t.Fatal("another issuer read the value", got)
		}
		if got, ok := (ValueIssuer{}).Read(value); ok || got != "" {
			t.Fatal("no issuer read the value", got)
		}
	}
	if _, ok := owner.Read(IssuedValue{}); ok {
		t.Fatal("an issuer read a value nobody issued")
	}
	if _, ok := (ValueIssuer{}).Read(IssuedValue{}); ok {
		t.Fatal("no issuer read a value nobody issued")
	}
	if value, err := (ValueIssuer{}).Issue("Person"); !errors.Is(err, ErrInvalid) || value != (IssuedValue{}) {
		t.Fatal("a value without an issuer", err)
	}
}

// A copy of a description owns its values: the kernel keeps copies, and what a
// library or a target does to its own slice afterwards does not reach them.
func TestDescriptionCopiesOwnTheirValues(t *testing.T) {
	value, err := NewValueIssuer().Issue("Person")
	if err != nil {
		t.Fatal(err)
	}
	native := NativeDeclaration[outputLanguage]{Values: []IssuedValue{value}}
	described := DeclarationDescription[outputLanguage]{Values: []IssuedValue{value}}
	nativeCopy, describedCopy := native.Clone(), described.Clone()
	native.Values[0], described.Values[0] = IssuedValue{}, IssuedValue{}
	if nativeCopy.Values[0] != value || describedCopy.Values[0] != value {
		t.Fatal("a copy shares its values with the original")
	}
}

// The records follow the trees in the snapshot's order, then the nodes of one
// tree. The first tree of this set has the greater identity, and its second
// declaration has a greater node number than the declaration of the second tree.
func TestSealedInventoryFollowsTheTreesOfTheSnapshot(t *testing.T) {
	lang, checker := outputChecking(t)
	first := ast.AST[outputLanguage]{ID: "z-unit", Root: 1, Nodes: []ast.Node[outputLanguage]{outputNode("unit", 2, 3), outputNode("declaration"), outputNode("declaration")}}
	second := ast.AST[outputLanguage]{ID: "a-unit", Root: 1, Nodes: []ast.Node[outputLanguage]{outputNode("unit", 2), outputNode("declaration")}}
	packed, err := representation.PackSet(lang, ast.Set[outputLanguage]{Entry: ast.NodeRef{Tree: first.ID, Node: 1}, Trees: []ast.AST[outputLanguage]{first, second}})
	if err != nil {
		t.Fatal(err)
	}
	// Planned order: the declaration of the second tree, then the two of the
	// first tree from the back.
	listed := []ast.NodeRef{{Tree: "a-unit", Node: 2}, {Tree: "z-unit", Node: 3}, {Tree: "z-unit", Node: 2}}
	declarations := []planned{{role: RequestedDeclaration}, {role: DependencyDeclaration}, {role: HelperDeclaration}}
	out, err := sealTest(checker, Assembly{Syntax: packed, Declarations: listed}, declarations, []uint32{1}, DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	var order []ast.TreeID
	for _, tree := range out.Snapshot().Inputs() {
		info, err := tree.Inspect()
		if err != nil {
			t.Fatal(err)
		}
		order = append(order, info.ID)
	}
	if !slices.Equal(order, []ast.TreeID{"z-unit", "a-unit"}) {
		t.Fatal("the snapshot has another order of trees", order)
	}
	inventory := out.Inventory()
	want := []OwnedDeclaration{
		{Site: ast.NodeRef{Tree: "z-unit", Node: 2}, Role: HelperDeclaration},
		{Site: ast.NodeRef{Tree: "z-unit", Node: 3}, Role: DependencyDeclaration},
		{Site: ast.NodeRef{Tree: "a-unit", Node: 2}, Role: RequestedDeclaration},
	}
	if len(inventory) != len(want) {
		t.Fatal("inventory size", len(inventory))
	}
	for i, record := range inventory {
		if record.Site != want[i].Site || record.Role != want[i].Role {
			t.Fatalf("record %d is %v with role %d", i, record.Site, record.Role)
		}
		if resolved, err := out.Resolve(record.Handle); err != nil || resolved.Site != record.Site {
			t.Fatal("handle does not resolve to its record", i, err)
		}
	}
	if exposed := out.Exposed(); len(exposed) != 1 || exposed[0] != want[2].Site {
		t.Fatal("exposed sites", exposed)
	}
}

// No listed site may contain another listed site. The unit of this test holds
// its two declarations one level down, as the C++ target does:
//
//	1 unit
//	  2 namespace
//	    3 declaration
//	      5 body
//	        6 leaf
//	    4 declaration
//	      7 body
//	        8 leaf
func TestSealRefusesAnInventoryWithNestedSites(t *testing.T) {
	lang, checker := outputChecking(t)
	tree := ast.AST[outputLanguage]{ID: "output", Root: 1, Nodes: []ast.Node[outputLanguage]{
		outputNode("unit", 2), outputNode("namespace", 3, 4), outputNode("declaration", 5), outputNode("declaration", 7),
		outputNode("body", 6), outputNode("leaf"), outputNode("body", 8), outputNode("leaf"),
	}}
	packed, err := representation.PackSet(lang, ast.Set[outputLanguage]{Entry: ast.NodeRef{Tree: tree.ID, Node: 1}, Trees: []ast.AST[outputLanguage]{tree}})
	if err != nil {
		t.Fatal(err)
	}
	declarations := []planned{{role: RequestedDeclaration}, {role: RequestedDeclaration}}
	for name, c := range map[string]struct {
		sites [2]ast.NodeID
		want  error
	}{
		"two declarations one level down":                {[2]ast.NodeID{3, 4}, nil},
		"nodes of two subtrees":                          {[2]ast.NodeID{3, 7}, nil},
		"the root and a declaration":                     {[2]ast.NodeID{1, 3}, ErrScope},
		"a declaration and the root":                     {[2]ast.NodeID{3, 1}, ErrScope},
		"the parent of a declaration":                    {[2]ast.NodeID{2, 4}, ErrScope},
		"a declaration and its child":                    {[2]ast.NodeID{3, 5}, ErrScope},
		"a child listed before its declaration":          {[2]ast.NodeID{5, 3}, ErrScope},
		"a declaration and a node two levels inside it":  {[2]ast.NodeID{4, 8}, ErrScope},
		"a node two levels inside, before its ancestor":  {[2]ast.NodeID{6, 3}, ErrScope},
		"the root and a node four levels below the root": {[2]ast.NodeID{8, 1}, ErrScope},
	} {
		t.Run(name, func(t *testing.T) {
			listed := []ast.NodeRef{{Tree: tree.ID, Node: c.sites[0]}, {Tree: tree.ID, Node: c.sites[1]}}
			out, err := sealTest(checker, Assembly{Syntax: packed, Declarations: listed}, declarations, []uint32{1, 2}, DefaultLimits())
			if c.want == nil {
				if err != nil || len(out.Inventory()) != 2 {
					t.Fatal("sites that do not contain each other were refused", err)
				}
				return
			}
			if !errors.Is(err, c.want) || out.Snapshot().Valid() || len(out.Inventory()) != 0 {
				t.Fatal("a site inside another listed site was sealed", err)
			}
		})
	}
}

// Formatting a value does not print its content, alone or inside a record or a
// description.
func TestIssuedValueDoesNotPrintItsContent(t *testing.T) {
	const content = "intent-that-stays-with-its-owner"
	value, err := NewValueIssuer().Issue(content)
	if err != nil {
		t.Fatal(err)
	}
	record := OwnedDeclaration{Values: []IssuedValue{value}}
	native := NativeDeclaration[outputLanguage]{Values: []IssuedValue{value}}
	described := DeclarationDescription[outputLanguage]{Values: []IssuedValue{value}}
	for _, verb := range []string{"%v", "%+v", "%#v", "%s", "%q", "%x", "%d"} {
		for name, subject := range map[string]any{"value": value, "pointer": &value, "values": []IssuedValue{value}, "record": record, "native declaration": native, "description": described} {
			printed := fmt.Sprintf(verb, subject)
			if strings.Contains(printed, content) || strings.Contains(printed, fmt.Sprintf("%x", content)) {
				t.Fatalf("%s of a %s prints the content: %s", verb, name, printed)
			}
		}
	}
	// What is printed says only whether the value was issued.
	other, err := NewValueIssuer().Issue("another content")
	if err != nil {
		t.Fatal(err)
	}
	if got, same, unissued := fmt.Sprint(value), fmt.Sprint(other), fmt.Sprint(IssuedValue{}); got == unissued || got != same {
		t.Fatal("what a value prints depends on more than whether it was issued:", got, same, unissued)
	}
}
