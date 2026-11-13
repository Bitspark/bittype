package integration_test

// PROTOTYPE (bittype#196), stage 4: the plan and its serialized form. Rows
// 16, 18, 19, 20, 22, 23 and 32 of gate lines 2, 3, 4 and 6.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/format/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	bfiles "github.com/Bitspark/bittype/files/go"
	tsfiles "github.com/Bitspark/bittype/format/files/typescript/provider/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// bindingsProduct is a product that asks for bindings and requires no facet.
// It is enough to make a request "types and bindings" (O5).
var bindingsProduct = g.Product{Identity: "example.org/fixtures/bindings"}

func independentPlan(t *testing.T, o independentOptions, bindings bool, roots ...string) (g.ContributionPlan, error) {
	t.Helper()
	fixture := protoFixture(t)
	composition := independentComposition(t, fixture, o)
	if bindings {
		composition.Products, composition.Requested = []g.Product{bindingsProduct}, []string{bindingsProduct.Identity}
	}
	return fixture.PlanContributions(context.Background(), composition, roots...)
}

func symbolNamed(plan g.RepresentationPlan, name string) (g.PlanSymbol, bool) {
	for _, s := range plan.Symbols() {
		if strings.HasSuffix(s.ID, "/"+name) {
			return s, true
		}
	}
	return g.PlanSymbol{}, false
}

// usesOf returns the use entries that refer to a node.
func usesOf(plan g.RepresentationPlan, node string) []g.PlanUse {
	var out []g.PlanUse
	for _, u := range plan.Uses() {
		if u.Node == node {
			out = append(out, u)
		}
	}
	return out
}

// Row 18: "Two uses of `Page<Person>` are two use entries on one
// representation node."
//
// Shows: a use entry records one selection at one use, with its provenance,
// and the uses that selected the same representation share one node. The
// open declaration Page<T> is one node with one binder, and Page<Person> is
// an application that links the binder to what was selected for Person.
// Does not show: the procedure #195 will choose to find that two uses share a
// node. The prototype identifies a node by a digest of its content.
func TestPrototypeStage4TwoUsesOneNode(t *testing.T) {
	plan, err := independentPlan(t, independentOptions{}, false, "People", "PeopleAgain")
	if err != nil {
		t.Fatal(err)
	}
	people, again := bodyNode(t, plan.Plan, "People"), bodyNode(t, plan.Plan, "PeopleAgain")
	entries := usesOf(plan.Plan, people.ID)
	if people.ID != again.ID || len(entries) != 2 || entries[0].Site == entries[1].Site || entries[0].Scope == entries[1].Scope {
		t.Fatalf("the two uses of Page<Person> are the nodes %s and %s with the use entries %+v", people.ID, again.ID, entries)
	}
	t.Logf("GREEN two use entries %d and %d, at two sites, on the node %s", entries[0].ID, entries[1].ID, people.ID)
	page := symbolNode(t, plan.Plan, "Page")
	person, _ := symbolNamed(plan.Plan, "Person")
	if people.Kind != "application" || people.Application.Declaration != page.ID || page.Kind != "declaration" || page.Binders != 1 || len(people.Application.Native) != 1 || people.Application.Native[0].Binder != 0 {
		t.Errorf("Page<Person> is %+v over %+v", people.Application, page)
	}
	argument, _ := plan.Plan.Node(people.Application.Native[0].Node)
	if argument.Kind != "application" || argument.Symbol != person.ID {
		t.Errorf("the native argument of Page<Person> is %+v", argument)
	}
	// Inside the open declaration the element is the binder, not an argument.
	body, _ := plan.Plan.Node(page.Links[0].Node)
	element, _ := plan.Plan.Node(body.Links[0].Node)
	if element.Kind != "binder" || element.Binder.Scope != page.Symbol || element.Binder.Index != 0 {
		t.Errorf("inside Page<T> the element is %+v", element)
	}
}

// Rows 19, 20 and 32, O5: "Equal native types with different adapters stay
// apart", and "An omitted native parameter still yields an adapter argument".
//
// Shows: Ignore<T> omits its parameter natively. Asked for bindings, the
// application Ignore<Person> has no native argument and one adapter argument,
// and Ignore<Person> and Ignore<Page<Person>>, whose native types are equal,
// are two nodes. Asked for types alone, the omitted argument is not planned:
// no contribution for it is needed. The same holds for an exact contribution
// that names its argument for conversion only. Red under a kernel that drops
// the omitted argument, and under a kernel whose node identity leaves out the
// links.
// Does not show: an adapter. No value facet is involved; the plan records
// which node an adapter would be built from.
func TestPrototypeStage4OmittedArgumentIsKeptForConversion(t *testing.T) {
	check := func() error {
		plan, err := independentPlan(t, independentOptions{}, true, "Kept", "KeptPage")
		if err != nil {
			return err
		}
		kept, keptPage := bodyNode(t, plan.Plan, "Kept"), bodyNode(t, plan.Plan, "KeptPage")
		if kept.Application == nil || len(kept.Application.Native) != 0 {
			return fmt.Errorf("Ignore<Person> has the native arguments %+v", kept.Application)
		}
		if len(kept.Application.Adapters) != 1 || kept.Application.Adapters[0].Binder != 0 {
			return fmt.Errorf("Ignore<Person> has the adapter arguments %+v", kept.Application.Adapters)
		}
		if kept.ID == keptPage.ID {
			return fmt.Errorf("Ignore<Person> and Ignore<Page<Person>> are one node")
		}
		return nil
	}
	redGreen(t, "asked for bindings, an omitted native parameter yields an adapter argument, and equal native types with different adapters stay apart", check, "application-drops-omitted-argument", "node-identity-ignores-links")
	// Types alone.
	types, err := independentPlan(t, independentOptions{}, false, "Kept", "KeptPage")
	if err != nil {
		t.Fatal(err)
	}
	if _, planned := symbolNamed(types.Plan, "Person"); planned || len(bodyNode(t, types.Plan, "Kept").Application.Adapters) != 0 {
		t.Error("asked for types alone, the omitted argument was planned")
	} else {
		t.Logf("GREEN asked for types alone, the omitted argument is not planned; products %v", types.Plan.Products())
	}
	text := protoText(t, golangText(), types)
	if !strings.Contains(text, "type Kept = Ignore") || !strings.Contains(text, "type KeptPage = Ignore") {
		t.Errorf("the two applications do not have equal native types:\n%s", text)
	}
	// The exact contribution that omits its parameter natively.
	for _, bindings := range []bool{false, true} {
		plan, err := independentPlan(t, independentOptions{cellConversionOnly: true}, bindings, "CellOfPerson")
		if err != nil {
			t.Fatal(err)
		}
		link := bodyNode(t, plan.Plan, "CellOfPerson").Links[0]
		_, planned := symbolNamed(plan.Plan, "Person")
		if planned != bindings || (link.Node != "") != bindings {
			t.Errorf("bindings %v: the conversion-only argument is linked as %+v", bindings, link)
		} else {
			t.Logf("GREEN bindings %v: the exact contribution's conversion-only argument is linked as %q", bindings, link.Need)
		}
	}
}

// memoSource is what a reader derives a memo key from, read from the
// serialized plan alone: the node of a use. A node's identity is over its
// decisions, its application and the nodes selected for its dependencies, so
// it is "representation nodes, applications and selected dependencies".
func memoSource(t *testing.T, serialized []byte, declaration string) string {
	t.Helper()
	plan, err := g.ParsePlan(serialized)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte("memo\x00" + bodyNode(t, plan, declaration).ID))
	return hex.EncodeToString(sum[:8])
}

// Row 22, O20: a memo key "from the serialized plan alone: two uses with
// equal nodes, applications and selected dependencies get equal keys, and a
// different selected dependency gives a different key".
//
// Shows: the plan carries what a token is derived from. A reader that parsed
// the serialized plan, and holds nothing of the planner, derives equal keys
// for two uses of Page<Person> and another key once another contribution is
// selected for Person, although the native type of that use is unchanged.
// Does not show: the memo, its token or its states. They are the service
// library's (bitspark-service#40). The stand-in repeats this in stage 7.
func TestPrototypeStage4MemoKeyFromTheSerializedPlan(t *testing.T) {
	serialize := func(o independentOptions) []byte {
		plan, err := independentPlan(t, o, true, "People", "PeopleAgain")
		if err != nil {
			t.Fatal(err)
		}
		serialized, err := plan.Plan.Serialize()
		if err != nil {
			t.Fatal(err)
		}
		return serialized
	}
	one := serialize(independentOptions{})
	// The same composition with another representation revision of the
	// contribution selected for Person's body: a different selected dependency.
	fixture := protoFixture(t)
	composition := independentComposition(t, fixture, independentOptions{})
	for i := range composition.Contributions {
		if composition.Contributions[i].Reference.Identity.Key == "independent-person" {
			composition.Contributions[i].Reference.Representation = "golang/2"
		}
	}
	composition.Products, composition.Requested = []g.Product{bindingsProduct}, []string{bindingsProduct.Identity}
	other, err := fixture.PlanContributions(context.Background(), composition, "People", "PeopleAgain")
	if err != nil {
		t.Fatal(err)
	}
	two, err := other.Plan.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	a, b, c := memoSource(t, one, "People"), memoSource(t, one, "PeopleAgain"), memoSource(t, two, "People")
	if a != b || a == c {
		t.Errorf("memo keys %s, %s and, with another dependency, %s", a, b, c)
	} else {
		t.Logf("GREEN equal keys %s for two uses of Page<Person>; %s with another selected dependency", a, c)
	}
}

// Row 16: "Each axis is exercised alone: an unknown required value of one
// axis is refused while the other four are unchanged."
//
// Shows: the plan states the envelope protocol, the payload schema of each
// selected contribution, both parts of the representation semantics and the
// source semantics, each in its own place, and a reader refuses an unknown
// value of each under that axis's name while it accepts the plan that
// differs in nothing else.
// Does not show: the fourth axis, generated-code and runtime compatibility.
// It is stated in the facet that names a factory, and the stand-in of stage 7
// exercises it.
func TestPrototypeStage4VersionAxesAlone(t *testing.T) {
	from, input := protoSource(t, "module example.axes;\n\ntype Names = list<text>;\n")
	plan := func(representation, meaning, schemaVersion string) (g.RepresentationPlan, []byte) {
		text, _ := nilFreeText(t, from, representation, meaning)
		// The fixture's schema in the given version, with a codec of that version.
		schema, codec := emptyFixtureCodec(t, "nil_free_text")
		schema.Version = schemaVersion
		text.Schema = schema
		p, err := protoPlanGo(t, input, golangBundle(t, from).replace(text).withCodecs(versioned{codec, schema}), protoRequest{})
		if err != nil {
			t.Fatal(err)
		}
		serialized, err := p.Plan.Serialize()
		if err != nil {
			t.Fatal(err)
		}
		return p.Plan, serialized
	}
	base, serialized := plan("golang/1", "1", "1")
	knowledge := g.ReaderKnowledge{SourceSemantics: []string{source.Semantics}}
	for _, c := range base.Contributions() {
		knowledge.Schemas, knowledge.Contributions = append(knowledge.Schemas, c.Schema), append(knowledge.Contributions, c.Reference)
	}
	read := func(serialized []byte) error {
		parsed, err := g.ParsePlan(serialized)
		if err != nil {
			return err
		}
		return parsed.CheckKnown(knowledge)
	}
	if err := read(serialized); err != nil {
		t.Fatalf("the reader refuses the plan it knows: %v", err)
	}
	edit := func(serialized []byte, change func(map[string]any)) []byte {
		var data map[string]any
		if err := json.Unmarshal(serialized, &data); err != nil {
			t.Fatal(err)
		}
		change(data)
		out, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	_, otherSchema := plan("golang/1", "1", "2")
	_, otherMeaning := plan("golang/1", "2", "1")
	_, otherRepresentation := plan("golang/2", "1", "1")
	for _, axis := range []struct {
		name, want string
		serialized []byte
	}{
		{"the envelope protocol", "envelope protocol", edit(serialized, func(d map[string]any) { d["envelope"] = "bittype.plan/9" })},
		{"the payload schema", "payload schema", otherSchema},
		{"representation semantics: the meaning revision", "meaning revision", otherMeaning},
		{"representation semantics: the representation revision", "representation revision", otherRepresentation},
		{"the source semantics", "source semantics", edit(serialized, func(d map[string]any) { d["source"].(map[string]any)["semantics"] = "horos/declarations-5" })},
	} {
		err := read(axis.serialized)
		if !errors.Is(err, g.ErrVersion) || !strings.Contains(err.Error(), axis.want) {
			t.Errorf("%s: an unknown value is not refused under its axis: %v", axis.name, err)
		} else {
			t.Logf("GREEN an unknown value of %s is refused: %v", axis.name, err)
		}
	}
}

// versioned presents a codec under another schema reference.
type versioned struct {
	g.PayloadCodec
	schema g.SchemaReference
}

func (v versioned) Schema() g.SchemaReference { return v.schema }

// realizationTS writes the realization manifest of one TypeScript output
// under one placement. It runs the TypeScript Files formatter for the
// placement and checks that each declaration is in the file it names.
//
// A stand-in: main has no reader on the formatter's result (#273 is open), so
// the manifest takes the module path from the placement that the consumer
// asked for and the name from the generation output, and only checks them
// against the text the formatter wrote. It does not read the realization.
func realizationTS(t *testing.T, plan g.ContributionPlan, output g.Output, placementName string, moduleOf func(name string) string) g.RealizationManifest {
	t.Helper()
	tree := outputTree[ts.Language](t, output)
	digest, err := plan.Plan.Digest()
	if err != nil {
		t.Fatal(err)
	}
	manifest := g.RealizationManifest{Plan: digest, Run: string(tree.ID), Target: plan.Plan.Target(), Placement: placementName}
	units := map[string]*tsfiles.Placement{}
	var order []string
	for _, record := range output.Inventory() {
		name := nodeAttribute(tree.Nodes[record.Site.Node-1], ts.AttributeName)
		if record.Symbol == "" || name == "" {
			t.Fatalf("an inventory record without the planner's identity or a name: %+v", record)
		}
		module := moduleOf(name)
		if units[module] == nil {
			units[module] = &tsfiles.Placement{File: bfiles.ID(strings.ReplaceAll(strings.TrimSuffix(module, ".ts"), "/", "-")), Path: bfiles.Path(module)}
			order = append(order, module)
		}
		units[module].Declarations = append(units[module].Declarations, record.Site)
		manifest.Symbols = append(manifest.Symbols, g.RealizedSymbol{Symbol: record.Symbol, Name: name, Module: module})
	}
	organization := tsfiles.Organization{Entry: units[order[0]].File}
	for _, module := range order {
		organization.Units = append(organization.Units, *units[module])
	}
	setting, err := organization.Setting()
	if err != nil {
		t.Fatal(err)
	}
	formatter, err := tsfiles.NewGenerated(placementStages(t))
	if err != nil {
		t.Fatal(err)
	}
	result, err := formatter.Format(context.Background(), format.GeneratedFilesInput{Output: output, Configuration: []semantic.Setting{setting}})
	if err != nil {
		t.Fatal(err)
	}
	texts := map[string]string{}
	for _, unit := range result.Units() {
		texts[string(unit.Path)] = formatted(t, typescriptText(), unit.Syntax.Abstract)
	}
	for _, s := range manifest.Symbols {
		if text := texts[s.Module]; !strings.Contains(text, "interface "+s.Name) && !strings.Contains(text, "type "+s.Name) {
			t.Fatalf("the formatter did not place %s in %s:\n%s", s.Name, s.Module, text)
		}
	}
	return manifest
}

// Rows 23 and 32: "A plan read with the realization manifest of another
// output fails as a linking error, not as an incomplete plan", and "One plan
// under two TypeScript placements resolves each logical symbol to another
// module path and the same declaration".
//
// Shows: a plan entry refers to a declaration by the planner's logical symbol
// identity; the inventory record of the sealed output carries that identity;
// a manifest relates plan, output and realization; one plan links with the
// manifests of two placements and resolves every symbol to another module
// path and the same name; the manifest of another plan's output is a linking
// error. Red under a reader that does not compare the plan's digest: it then
// reports an incomplete plan.
// Does not show: a reader on the formatter's result. main has none (#273), so
// the third test of gate line 6 is not run: "No native name reaches emitted
// glue except through the reader, and a formatter mutant that realizes
// another spelling is followed" (row 24).
func TestPrototypeStage4RealizationManifest(t *testing.T) {
	from, input := protoSource(t, protoPagination)
	plan, err := protoPlanTS(t, input, typescriptBundle(t, from), protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	output, err := plan.Generation.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// The planner's identity on the inventory record.
	symbols := map[string]bool{}
	for _, record := range output.Inventory() {
		symbols[record.Symbol] = true
	}
	for _, s := range plan.Plan.Symbols() {
		if !symbols[s.ID] {
			t.Errorf("no inventory record carries the logical symbol %s", s.ID)
		}
	}
	single := realizationTS(t, plan, output, "one module", func(string) string { return "generated/types.ts" })
	split := realizationTS(t, plan, output, "one module for each declaration", func(name string) string { return "generated/" + strings.ToLower(name) + ".ts" })
	one, err := g.LinkRealization(plan.Plan, single)
	if err != nil {
		t.Fatal(err)
	}
	two, err := g.LinkRealization(plan.Plan, split)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range plan.Plan.Symbols() {
		if one[s.ID].Name != two[s.ID].Name || one[s.ID].Module == two[s.ID].Module || one[s.ID].Name != s.Intent {
			t.Errorf("symbol %s resolves to %+v and %+v", s.ID, one[s.ID], two[s.ID])
		}
	}
	t.Logf("GREEN one plan under two placements: %d logical symbols, each at another module path with the same declaration", len(one))
	// The manifest of another plan's output.
	_, otherInput := protoSource(t, "module example.other;\n\ntype Other = { name: text; };\n")
	otherPlan, err := protoPlanTS(t, otherInput, typescriptBundle(t, from), protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	otherOutput, err := otherPlan.Generation.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	foreign := realizationTS(t, otherPlan, otherOutput, "one module", func(string) string { return "generated/types.ts" })
	redGreen(t, "a plan read with the realization manifest of another output fails as a linking error, not as an incomplete plan", func() error {
		_, err := g.LinkRealization(plan.Plan, foreign)
		if !errors.Is(err, g.ErrLinking) || errors.Is(err, g.ErrIncomplete) {
			return fmt.Errorf("got %v", err)
		}
		return nil
	}, "link-ignores-plan-digest")
}

var _ = f.Profile
var _ ast.NodeRef
