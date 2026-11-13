package integration_test

// The pair golang-to-typescript: the first of the two other pairs of
// bittype#196 (decision O43). Its purpose is evidence that the contribution
// contract and its kernel assume no horos source. It covers one struct type
// with string and bool fields and one use of a defined type by another. Each
// test says what it shows and what it does not.

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	pairs "github.com/Bitspark/bittype/integration/go/testdata/pairgeneration"
	gon "github.com/Bitspark/bittype/semantics/golang/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	gotsboolcodec "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/codec/booltype/go"
	gotsdeclarationscodec "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/codec/declarations/go"
	gotsstringcodec "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/codec/stringtype/go"
	gotsstructcodec "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/codec/structtype/go"
	gotsbool "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/mapping/booltype/go"
	gotsdeclarations "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/mapping/declarations/go"
	gotsstring "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/mapping/stringtype/go"
	gotsstruct "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/mapping/structtype/go"
	gotsboolvocabulary "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/vocabulary/booltype/go"
	gotsdeclarationsvocabulary "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/vocabulary/declarations/go"
	gotsstringvocabulary "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/vocabulary/stringtype/go"
	gotsstructvocabulary "github.com/Bitspark/bittype/transform/syntax/golang-to-typescript/vocabulary/structtype/go"
	goboolfamily "github.com/Bitspark/bittype/transform/syntax/golang/family/booltype/go"
	godeclarationsfamily "github.com/Bitspark/bittype/transform/syntax/golang/family/declarations/go"
	gostringfamily "github.com/Bitspark/bittype/transform/syntax/golang/family/stringtype/go"
	gostructfamily "github.com/Bitspark/bittype/transform/syntax/golang/family/structtype/go"
	gorules "github.com/Bitspark/bittype/transform/syntax/golang/mapping/go"
	gosource "github.com/Bitspark/bittype/transform/syntax/golang/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	tsmodel "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	tgtarget "github.com/Bitspark/bittype/transform/syntax/typescript/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type goTSBundle = pairs.Bundle[gs.Language, ts.Language, gosource.Reader]
type goTSContribution = g.Contribution[gs.Language, ts.Language, gosource.Reader]

// goTSBundles is the whole pair, as a consumer's composition registers it.
// Registration order is no preference.
func goTSBundles() map[string]goTSBundle {
	return map[string]goTSBundle{
		"structtype":   {Contribution: gotsstruct.Contribution, Codec: gotsstructcodec.Codec},
		"stringtype":   {Contribution: gotsstring.Contribution, Codec: gotsstringcodec.Codec},
		"booltype":     {Contribution: gotsbool.Contribution, Codec: gotsboolcodec.Codec},
		"declarations": {Contribution: gotsdeclarations.Contribution, Codec: gotsdeclarationscodec.Codec},
	}
}

// goTSComposition is one composition of the pair over one Go source text.
type goTSComposition struct {
	registry *c.Registry
	request  c.Request
	checker  check.Checker
	source   gs.AST
}

// goDeclarations returns the sites of the named type declarations of a Go tree.
func goDeclarations(t *testing.T, tree gs.AST, names ...string) []ast.NodeRef {
	t.Helper()
	var out []ast.NodeRef
	for _, name := range names {
		found := false
		for i, node := range tree.Nodes {
			if (node.Kind == gs.KindTypeDeclaration || node.Kind == gs.KindTypeAlias) && nodeAttribute(node, gs.AttributeName) == name {
				out, found = append(out, ast.NodeRef{Tree: tree.ID, Node: ast.NodeID(i + 1)}), true
			}
		}
		if !found {
			t.Fatalf("the Go source declares no type %s", name)
		}
	}
	return out
}

// composeGoTS parses and analyses a Go source text with the Go stages, and
// registers the Go source reader, the TypeScript target and the bundles. A
// bundle of extra replaces the bundle of the same name, and a nil bundle
// leaves it out. A fixture bundle has no codec of its own: its contribution
// states the schema of the string contribution, whose codec is registered.
func composeGoTS(t *testing.T, source string, roots []string, extra map[string]*goTSBundle) *goTSComposition {
	t.Helper()
	snapshot, trees := goCapture(t, source)
	_, analysed := goResults(t, snapshot)
	input := transform.Input{Source: snapshot, Evidence: []semantic.Result{analysed}}
	bundles := goTSBundles()
	for name, bundle := range extra {
		if bundle == nil {
			delete(bundles, name)
		} else {
			bundles[name] = *bundle
		}
	}
	pair := pairs.Pair[gs.Language, ts.Language, gosource.Reader]{Source: gosource.Provider{}, Target: tgtarget.Provider{}, Profile: gorules.TypesProfile}
	for _, bundle := range bundles {
		pair.Bundles = append(pair.Bundles, bundle)
	}
	q := &goTSComposition{registry: &c.Registry{}, checker: nativeChecker(t, typescriptText()), source: trees[0]}
	var err error
	if q.request, err = pairs.Compose(q.registry, pair, input, q.checker.Language(), "generated-ts", goDeclarations(t, trees[0], roots...)); err != nil {
		t.Fatal(err)
	}
	return q
}

func (q *goTSComposition) plan() (g.GenerationPlan, error) {
	return q.registry.PlanGeneration(context.Background(), q.request, q.checker)
}

const goTSSource = "type Person struct { Name string; Active bool }; type Team struct { Lead Person }"

// The end to end test of the pair.
//
// Shows: a Go source text is parsed and analysed by the Go stages; the
// request names one root, Team; the kernel plans it through the Go source
// reader and the four contributions of the pair; the TypeScript target builds
// the tree; the TypeScript checker admits it when the output is sealed; the
// TypeScript formatter prints it, and the text is asserted byte for byte.
// Person is in the output although it is no root: the kernel reached it
// through the use in Team. The type of Team's member is a type reference
// whose target is the declaration node of Person in the output tree, and the
// output's inventory records that node as a dependency declaration whose
// provenance is the site of the Go declaration of Person.
// Does not show: that the TypeScript type has the values of the Go type. No
// binding is defined for this pair, and nothing here converts a value. It
// does not show a Go name that TypeScript does not admit, a second package
// or an exported and an unexported field treated differently: the pair reads
// one package and spells every field as Go does.
func TestGolangToTypescriptEndToEnd(t *testing.T) {
	q := composeGoTS(t, goTSSource, []string{"Team"}, nil)
	out, err := q.registry.Generate(context.Background(), q.request, q.checker)
	if err != nil {
		t.Fatal(err)
	}
	tree := outputTree[ts.Language](t, out)
	const want = "interface Person {\n  Name: string;\n  Active: boolean;\n}\n\ninterface Team {\n  Lead: Person;\n}\n"
	if printed := formatted(t, typescriptText(), tree); printed != want {
		t.Fatalf("printed:\n%s\nwant:\n%s", printed, want)
	}
	person, reference := ast.NodeID(0), ast.NodeID(0)
	for i, node := range tree.Nodes {
		switch {
		case node.Kind == ts.KindInterface && nodeAttribute(node, "name") == "Person":
			person = ast.NodeID(i + 1)
		case node.Kind == ts.KindTypeReference:
			if reference != 0 {
				t.Fatal("the output has more than one type reference")
			}
			reference = ast.NodeID(i + 1)
		}
	}
	if person == 0 || reference == 0 {
		t.Fatal("the output has no interface Person or no type reference")
	}
	references := tree.Nodes[reference-1].References
	if len(references) != 1 || references[0].Role != ts.RoleTarget || references[0].Target != (ast.NodeRef{Tree: tree.ID, Node: person}) {
		t.Fatalf("the use of Person is not linked to the declaration of Person: %+v", references)
	}
	if uses := out.ExternalUses(); len(uses) != 0 {
		t.Fatalf("the output names %d external types", len(uses))
	}
	source := goDeclarations(t, q.source, "Person", "Team")
	roles := map[ast.NodeRef]g.DeclarationRole{}
	for _, record := range out.Inventory() {
		if len(record.Provenance) != 1 {
			t.Fatalf("a declaration of the output has the provenance %v", record.Provenance)
		}
		roles[record.Provenance[0]] = record.Role
		if record.Provenance[0] == source[0] && record.Site != (ast.NodeRef{Tree: tree.ID, Node: person}) {
			t.Fatalf("the declaration that comes from the Go type Person is at %v, and the interface Person at node %d", record.Site, person)
		}
	}
	if len(roles) != 2 || roles[source[0]] != g.DependencyDeclaration || roles[source[1]] != g.RequestedDeclaration {
		t.Fatalf("the inventory has the roles %v", roles)
	}
}

// goTSNamedPerson is a fixture contribution for the exact declaration Person
// that names a global TypeScript type "Person" and plans no declaration.
func goTSNamedPerson(from representation.Language) (goTSContribution, error) {
	selector, err := gosource.Defined(from, gon.PackagePath, "Person")
	if err != nil {
		return goTSContribution{}, err
	}
	return goTSContribution{
		Reference: g.ContributionReference{Identity: g.ContributionIdentity{Publisher: "bittype/test", Key: "golang-to-typescript/person-by-name"}, Meaning: godeclarationsfamily.Meaning, Representation: "typescript/test"},
		Target:    string(ts.Language{}.ID()), Profile: gorules.TypesProfile, Domain: g.ForExact(selector), Schema: gotsstringvocabulary.Schema,
		Rule: g.Rule[gs.Language, gosource.Reader]{Decide: func(context.Context, g.DecideStep[gs.Language, gosource.Reader]) (g.Decision, error) {
			return g.Decision{Payload: gotsstringvocabulary.Payload{}}, nil
		}},
		Project: func(context.Context, g.Projection[ts.Language]) (g.Projected[ts.Language], error) {
			return g.Projected[ts.Language]{Type: tsmodel.External("", "Person")}, nil
		},
	}, nil
}

// The control of the link assertion.
//
// Shows: when a contribution writes the name "Person" where the pair links
// the declaration, the member of Team prints as the same text, "Lead:
// Person;", and the assertions of the end to end test tell the two apart:
// the type reference has no target, the output has no declaration of Person,
// and the output records one external use. So the end to end test's
// assertion on the reference is an assertion about a link and not about a
// spelling.
// Does not show: anything about the pair's own contributions. The naming
// contribution is a fixture of this test.
func TestGolangToTypescriptNameIsNoLink(t *testing.T) {
	q := composeGoTS(t, goTSSource, []string{"Team"}, map[string]*goTSBundle{"person-by-name": {Contribution: goTSNamedPerson}})
	out, err := q.registry.Generate(context.Background(), q.request, q.checker)
	if err != nil {
		t.Fatal(err)
	}
	tree := outputTree[ts.Language](t, out)
	const want = "interface Team {\n  Lead: Person;\n}\n"
	if printed := formatted(t, typescriptText(), tree); printed != want {
		t.Fatalf("printed:\n%s\nwant:\n%s", printed, want)
	}
	for _, node := range tree.Nodes {
		if node.Kind == ts.KindTypeReference && len(node.References) != 0 {
			t.Fatalf("the named type is linked: %+v", node.References)
		}
	}
	if len(out.ExternalUses()) != 1 || len(out.Inventory()) != 1 {
		t.Fatalf("%d external uses and %d declarations", len(out.ExternalUses()), len(out.Inventory()))
	}
}

// goTSStructWithout is a fixture contribution for the struct type form that
// projects a primitive and states the given facts about the record form.
func goTSStructWithout(facts ...g.Fact) func(representation.Language) (goTSContribution, error) {
	return func(from representation.Language) (goTSContribution, error) {
		selector, err := gosource.StructType(from)
		if err != nil {
			return goTSContribution{}, err
		}
		return goTSContribution{
			Reference: g.ContributionReference{Identity: g.ContributionIdentity{Publisher: "bittype/test", Key: "golang-to-typescript/struct-as-text"}, Meaning: gostructfamily.Meaning, Representation: "typescript/test"},
			Target:    string(ts.Language{}.ID()), Profile: gorules.TypesProfile, Domain: g.ForForm(selector), Schema: gotsstringvocabulary.Schema,
			Rule: g.Rule[gs.Language, gosource.Reader]{Decide: func(context.Context, g.DecideStep[gs.Language, gosource.Reader]) (g.Decision, error) {
				return g.Decision{Payload: gotsstringvocabulary.Payload{}, Facts: facts}, nil
			}},
			Project: func(context.Context, g.Projection[ts.Language]) (g.Projected[ts.Language], error) {
				return g.Projected[ts.Language]{Type: tsmodel.Primitive(tsmodel.PrimitiveString)}, nil
			},
		}, nil
	}
}

// goTSRefusal is one case of what the pair does not cover.
type goTSRefusal struct {
	name, source string
	roots        []string
	extra        map[string]*goTSBundle
	// A refusal of a contribution or of the kernel: its class, the owner and
	// the reason. For a judgment of the source reader class is zero, and
	// reason is the reader's diagnostic. uncovered is a source that the Go
	// analysis does not cover: the evidence is not admitted.
	class         g.RefusalClass
	uncovered     bool
	owner, reason string
	contribution  g.ContributionReference
}

func goTSRefusals() []goTSRefusal {
	publisher := gorules.Publisher
	structure, declarations := gotsstructvocabulary.Reference, gotsdeclarationsvocabulary.Reference
	return []goTSRefusal{
		// What the Go source reader does not read: a judgment of the source,
		// with the verdict unsupported.
		{name: "a number field", source: "type T struct { Count int }", roots: []string{"T"}, reason: "this Go source reader reads struct types, string and bool, and not the type int"},
		{name: "a slice field", source: "type T struct { Names []string }", roots: []string{"T"}, reason: "this Go source reader reads struct types, string and bool, and not the type []string"},
		{name: "a generic declaration", source: "type T[E any] struct { Name string }", roots: []string{"T"}, reason: "this Go source reader reads no generic declaration"},
		{name: "an alias declaration", source: "type P struct { Name string }; type T = P", roots: []string{"T"}, reason: "this Go source reader reads no alias declaration"},
		{name: "a use of an alias", source: "type P struct { Name string }; type A = P; type T struct { Lead A }", roots: []string{"T"}, reason: "this Go source reader reads a use of a defined type, and of no alias or type parameter"},
		// What the Go analysis of the base profile does not cover. It publishes
		// its graph and reports that its coverage does not hold, so the kernel
		// admits no evidence and neither the reader nor a contribution is asked.
		{name: "a pointer field", source: "type P struct { Name string }; type T struct { Next *P }", roots: []string{"T"}, uncovered: true},
		{name: "a map field", source: "type T struct { Index map[string]bool }", roots: []string{"T"}, uncovered: true},
		{name: "an embedded field", source: "type P struct { Name string }; type T struct { P }", roots: []string{"T"}, uncovered: true},
		// What a contribution of the pair recognizes and does not cover.
		{name: "a field with a tag", source: "type T struct { Name string `json:\"name\"` }", roots: []string{"T"}, class: g.RefusalUnsupported, owner: publisher, reason: gotsstruct.ReasonFieldTag, contribution: structure},
		{name: "a struct type in a field", source: "type T struct { Home struct { City string } }", roots: []string{"T"}, class: g.RefusalUnsupported, owner: publisher, reason: gotsstruct.ReasonNestedStruct, contribution: structure},
		// A defined type whose body is not selected as the record form: the
		// declaration contribution needs that fact and nothing states it.
		{name: "a defined string", source: "type T string", roots: []string{"T"}, class: g.RefusalRequirementUnresolved, owner: publisher, reason: gotsdeclarations.ReasonNoRecordForm, contribution: declarations},
		{name: "a use of a defined string", source: "type N string; type T struct { Name N }", roots: []string{"T"}, class: g.RefusalRequirementUnresolved, owner: publisher, reason: gotsdeclarations.ReasonNoRecordForm, contribution: declarations},
		{name: "a struct selected in a form that states nothing", source: "type T struct { Name string }", roots: []string{"T"}, extra: map[string]*goTSBundle{"structtype": {Contribution: goTSStructWithout()}}, class: g.RefusalRequirementUnresolved, owner: publisher, reason: gotsdeclarations.ReasonNoRecordForm, contribution: declarations},
		{name: "a struct selected in a form that refutes the record form", source: "type T struct { Name string }", roots: []string{"T"}, extra: map[string]*goTSBundle{"structtype": {Contribution: goTSStructWithout(tsmodel.EvidenceRecordForm.Refuted(""))}}, class: g.RefusalRequirementUnsatisfied, owner: publisher, reason: gotsdeclarations.ReasonNoRecordForm, contribution: declarations},
		// A composition that lacks a contribution of the pair.
		{name: "no contribution for bool", source: "type T struct { Active bool }", roots: []string{"T"}, extra: map[string]*goTSBundle{"booltype": nil}, class: g.RefusalMissingContribution, owner: g.KernelOwner, reason: g.ReasonNoContribution + " form:type/basic/bool"},
	}
}

// Each refusal of the pair.
//
// Shows, for each case: planning fails and gives no plan; a refusal of a
// contribution or of the kernel has exactly the class, the owner, the reason
// and the selected contribution listed, matches the contract value of its
// class and no other of the eight; a judgment of the Go source reader fails
// with the value of the verdict unsupported, is no refusal of a contribution,
// and carries the reader's diagnostic; a source that the Go analysis does not
// cover fails with the value of evidence that is not established, before the
// reader is opened.
// Does not show: that the list is complete. It has one case for each kind of
// thing the reader judges unsupported and for each refusal that a
// contribution of the pair states and that a source can reach; a Go type
// expression that is not listed (a channel, a function type, an interface, an
// array) is not run here. It does not show the struct contribution's refusal
// of an embedded field: the Go analysis of the base profile does not cover an
// embedded field (docs/native-semantics.md, "First Go profile"), so no source
// reaches that line, and the contract gives a test no way to run a rule
// without the planner. The kernel's refusal of a recursive reference is not
// reachable through this pair either: Go admits a recursive type only
// through a pointer, a slice or a map.
func TestGolangToTypescriptRefusals(t *testing.T) {
	values := []error{g.ErrUnsupported, g.ErrRefusedByDesign, g.ErrRequirementUnsatisfied, g.ErrRequirementUnresolved, g.ErrMissingFacet, g.ErrAmbiguous, g.ErrMissingContribution, g.ErrLimit}
	for _, refusal := range goTSRefusals() {
		t.Run(refusal.name, func(t *testing.T) {
			q := composeGoTS(t, refusal.source, refusal.roots, refusal.extra)
			plan, err := q.plan()
			if err == nil || plan.Valid() {
				t.Fatal("the request planned")
			}
			var got *g.Refusal
			if refusal.uncovered {
				if errors.As(err, &got) || !errors.Is(err, semantic.ErrEvidenceUnavailable) || errors.Is(err, g.ErrUnsupported) {
					t.Fatalf("not the failure of evidence that is not established: %v", err)
				}
				return
			}
			if refusal.class == 0 {
				var diagnostic *text.Diagnostic
				if errors.As(err, &got) || !errors.Is(err, g.ErrUnsupported) || !errors.As(err, &diagnostic) || diagnostic.Code != text.Unsupported || diagnostic.Message != refusal.reason {
					t.Fatalf("not the source reader's judgment %q: %v", refusal.reason, err)
				}
				return
			}
			if !errors.As(err, &got) {
				t.Fatalf("no refusal: %v", err)
			}
			if got.Class != refusal.class || got.Owner != refusal.owner || got.Reason != refusal.reason || got.Contribution != refusal.contribution || got.Target != string(ts.Language{}.ID()) || got.Profile != gorules.TypesProfile {
				t.Fatalf("refused as class %d by %q for %q with the contribution %v: %v", got.Class, got.Owner, got.Reason, got.Contribution, err)
			}
			for _, value := range values {
				if errors.Is(err, value) != (value == refusal.class.Value()) {
					t.Fatalf("the refusal of class %d matches the value %q wrongly: %v", refusal.class, value, err)
				}
			}
		})
	}
}

// The plan of the pair.
//
// Shows: the plan of the end to end request serializes; the serialized form
// reads back in ParsePlan to a plan with the same digest that serializes to
// the same bytes; the plan names Go as its source with the analysis profile
// of the Go vocabulary as the source semantics, TypeScript as its target and
// the pair's profile; its contribution entries are exactly the four
// contributions of the pair, each with the pair's publisher, a local key
// under the pair's name, the meaning of its family of the Go source and the
// schema of its vocabulary; the root is the symbol of the Go declaration
// Team, and the symbol of Person is a dependency.
// Does not show: a reader of the plan in another process or another
// language. The plan has no binding product, so no reader has anything to
// convert.
func TestGolangToTypescriptPlan(t *testing.T) {
	q := composeGoTS(t, goTSSource, []string{"Team"}, nil)
	generation, err := q.plan()
	if err != nil {
		t.Fatal(err)
	}
	plan := generation.Representation()
	serialized, err := plan.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	read, err := g.ParsePlan(serialized)
	if err != nil {
		t.Fatal(err)
	}
	again, err := read.Serialize()
	if err != nil || !bytes.Equal(again, serialized) {
		t.Fatal("the plan that was read back serializes to other bytes:", err)
	}
	first, err1 := plan.Digest()
	second, err2 := read.Digest()
	if err1 != nil || err2 != nil || first != second || first == "" {
		t.Fatalf("digests %q and %q: %v %v", first, second, err1, err2)
	}
	if read.Source() != (g.PlanSource{Language: string(gs.Language{}.ID()), Semantics: gon.AnalysisProfile}) || read.Target().Language != "typescript" || read.Profile() != gorules.TypesProfile {
		t.Fatalf("the plan names the source %+v, the target %+v and the profile %q", read.Source(), read.Target(), read.Profile())
	}
	want := map[string]struct {
		reference g.ContributionReference
		meaning   g.Meaning
		schema    g.SchemaReference
	}{
		"structtype":   {gotsstructvocabulary.Reference, gostructfamily.Meaning, gotsstructvocabulary.Schema},
		"stringtype":   {gotsstringvocabulary.Reference, gostringfamily.Meaning, gotsstringvocabulary.Schema},
		"booltype":     {gotsboolvocabulary.Reference, goboolfamily.Meaning, gotsboolvocabulary.Schema},
		"declarations": {gotsdeclarationsvocabulary.Reference, godeclarationsfamily.Meaning, gotsdeclarationsvocabulary.Schema},
	}
	entries := read.Contributions()
	if len(entries) != len(want) {
		t.Fatalf("the plan names %d contributions", len(entries))
	}
	for _, entry := range entries {
		name, found := strings.CutPrefix(entry.Reference.Identity.Key, "golang-to-typescript/")
		expected, known := want[name]
		if !found || !known || entry.Reference != expected.reference || entry.Reference.Identity.Publisher != gorules.Publisher || entry.Reference.Meaning != expected.meaning || entry.Schema != expected.schema || entry.Target != "typescript" || entry.Profile != gorules.TypesProfile {
			t.Errorf("the plan names the contribution %+v", entry)
		}
		delete(want, name)
	}
	roots, symbols := read.Roots(), read.Symbols()
	team := "source:declaration/" + gon.PackagePath + "/Team"
	if len(roots) != 1 || roots[0] != team || len(symbols) != 2 {
		t.Fatalf("roots %v and %d symbols", roots, len(symbols))
	}
	for _, symbol := range symbols {
		role := map[string]string{team: "requested", "source:declaration/" + gon.PackagePath + "/Person": "dependency"}[symbol.ID]
		if role == "" || symbol.Role != role || symbol.Intent != strings.TrimPrefix(symbol.ID, "source:declaration/"+gon.PackagePath+"/") {
			t.Errorf("the plan has the symbol %+v", symbol)
		}
	}
}
