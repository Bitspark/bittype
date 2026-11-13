package integration_test

// The pair typescript-to-horos: the second of the two other pairs of
// bittype#196 (decision O43). Its purpose is evidence that the contribution
// contract and its kernel assume no native target. It covers one object type
// with string and boolean properties, declared by an interface or by a type
// alias, which becomes a closed horos record of text and booleans. It is
// done-when 1 of bittype#173 in the smallest form. Each test says what it
// shows and what it does not.

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	pairs "github.com/Bitspark/bittype/integration/go/testdata/pairgeneration"
	tsn "github.com/Bitspark/bittype/semantics/typescript/go"
	hs "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	htarget "github.com/Bitspark/bittype/transform/syntax/horos/target/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	tshbooleancodec "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/codec/booleantype/go"
	tshdeclarationscodec "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/codec/declarations/go"
	tshnumbercodec "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/codec/numbertype/go"
	tshobjectcodec "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/codec/objecttype/go"
	tshstringcodec "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/codec/stringtype/go"
	tshboolean "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/mapping/booleantype/go"
	tshdeclarations "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/mapping/declarations/go"
	tshnumber "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/mapping/numbertype/go"
	tshobject "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/mapping/objecttype/go"
	tshstring "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/mapping/stringtype/go"
	tshbooleanvocabulary "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/vocabulary/booleantype/go"
	tshdeclarationsvocabulary "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/vocabulary/declarations/go"
	tshnumbervocabulary "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/vocabulary/numbertype/go"
	tshobjectvocabulary "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/vocabulary/objecttype/go"
	tshstringvocabulary "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/vocabulary/stringtype/go"
	tsbooleanfamily "github.com/Bitspark/bittype/transform/syntax/typescript/family/booleantype/go"
	tsdeclarationsfamily "github.com/Bitspark/bittype/transform/syntax/typescript/family/declarations/go"
	tsobjectfamily "github.com/Bitspark/bittype/transform/syntax/typescript/family/objecttype/go"
	tsstringfamily "github.com/Bitspark/bittype/transform/syntax/typescript/family/stringtype/go"
	tsrules "github.com/Bitspark/bittype/transform/syntax/typescript/mapping/go"
	tssource "github.com/Bitspark/bittype/transform/syntax/typescript/source/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type tsHBundle = pairs.Bundle[ts.Language, hs.Language, tssource.Reader]

// tsHBundles is the whole pair, as a consumer's composition registers it.
// Registration order is no preference.
func tsHBundles() map[string]tsHBundle {
	return map[string]tsHBundle{
		"objecttype":   {Contribution: tshobject.Contribution, Codec: tshobjectcodec.Codec},
		"stringtype":   {Contribution: tshstring.Contribution, Codec: tshstringcodec.Codec},
		"booleantype":  {Contribution: tshboolean.Contribution, Codec: tshbooleancodec.Codec},
		"numbertype":   {Contribution: tshnumber.Contribution, Codec: tshnumbercodec.Codec},
		"declarations": {Contribution: tshdeclarations.Contribution, Codec: tshdeclarationscodec.Codec},
	}
}

// tsHComposition is one composition of the pair over one TypeScript source text.
type tsHComposition struct {
	registry *c.Registry
	request  c.Request
	checker  check.Checker
	source   ts.AST
}

// tsSourceDeclarations returns the sites of the named interfaces and type aliases
// of a TypeScript tree.
func tsSourceDeclarations(t *testing.T, tree ts.AST, names ...string) []ast.NodeRef {
	t.Helper()
	var out []ast.NodeRef
	for _, name := range names {
		found := false
		for i, node := range tree.Nodes {
			if (node.Kind == ts.KindInterface || node.Kind == ts.KindTypeAlias) && nodeAttribute(node, "name") == name {
				out, found = append(out, ast.NodeRef{Tree: tree.ID, Node: ast.NodeID(i + 1)}), true
			}
		}
		if !found {
			t.Fatalf("the TypeScript source declares no type %s", name)
		}
	}
	return out
}

// horosTarget is the horos language as the target of the pair, under the
// profile of horos declaration syntax, and its own checker.
func horosTarget(t *testing.T) check.Checker {
	t.Helper()
	to, err := representation.NewLanguage[hs.Language](horosText(), hs.Profile)
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[hs.Language](to, horosText())
	if err != nil {
		t.Fatal(err)
	}
	return checker
}

const tsHModule = "example.people"

// composeTSH parses and analyses a TypeScript source text with the TypeScript
// stages, and registers the TypeScript source reader, the horos target and
// the bundles. A nil bundle of extra leaves the bundle of that name out.
func composeTSH(t *testing.T, source string, roots []string, extra map[string]*tsHBundle) *tsHComposition {
	t.Helper()
	snapshot, trees, _ := tsCapture(t, source)
	_, analysed := tsResults(t, snapshot)
	input := transform.Input{Source: snapshot, Evidence: []semantic.Result{analysed}}
	bundles := tsHBundles()
	for name, bundle := range extra {
		if bundle == nil {
			delete(bundles, name)
		} else {
			bundles[name] = *bundle
		}
	}
	pair := pairs.Pair[ts.Language, hs.Language, tssource.Reader]{Source: tssource.Provider{}, Target: htarget.Provider{Module: tsHModule}, Profile: tsrules.CanonicalizationProfile}
	for _, bundle := range bundles {
		pair.Bundles = append(pair.Bundles, bundle)
	}
	q := &tsHComposition{registry: &c.Registry{}, checker: horosTarget(t), source: trees[0]}
	var err error
	if q.request, err = pairs.Compose(q.registry, pair, input, q.checker.Language(), "generated-horos", tsSourceDeclarations(t, trees[0], roots...)); err != nil {
		t.Fatal(err)
	}
	return q
}

func (q *tsHComposition) plan() (g.GenerationPlan, error) {
	return q.registry.PlanGeneration(context.Background(), q.request, q.checker)
}

const tsHSource = "export interface Person { name: string; active: boolean }\nexport type Flag = { on: boolean };\n"

// The end to end test of the pair.
//
// Shows: a TypeScript source text is parsed and analysed by the TypeScript
// stages (the analysis runs the pinned compiler in a node process); the
// request names two roots, an interface and a type alias of an object type;
// the kernel plans them through the TypeScript source reader and the
// contributions of the pair; the horos target builds a horos tree; the horos
// syntax checker admits it when the output is sealed; the horos formatter,
// which checks the module's pin against the canonical bytes of the tree's
// own declarations, prints it, and the text is asserted byte for byte. The
// printed text parses again with the horos parser, and the parser assigns the
// module the pin that the target gave it. The output's inventory records
// each definition as requested, with the site of its TypeScript declaration
// as its provenance.
// Does not show: that the horos type has the values of the TypeScript type.
// The canonicalization is a statement of the profile, and nothing here
// converts or compares a value. It does not show that the generated module
// resolves and analyses as a horos program: the horos resolver is not run on
// it. It does not show a reference from one definition to another: the pair
// refuses it.
func TestTypescriptToHorosEndToEnd(t *testing.T) {
	q := composeTSH(t, tsHSource, []string{"Person", "Flag"}, nil)
	out, err := q.registry.Generate(context.Background(), q.request, q.checker)
	if err != nil {
		t.Fatal(err)
	}
	tree := outputTree[hs.Language](t, out)
	const want = "module example.people;\n\ntype Flag = {\n    on: bool;\n};\n\ntype Person = {\n    name: text;\n    active: bool;\n};\n"
	printed := formatted(t, horosText(), tree)
	if printed != want {
		t.Fatalf("printed:\n%s\nwant:\n%s", printed, want)
	}
	module := func(tree hs.AST) hs.Node {
		for _, node := range tree.Nodes {
			if node.Kind == hs.KindModule {
				return node
			}
		}
		t.Fatal("the tree has no module")
		return hs.Node{}
	}
	parsed, err := horosText().Parse(context.Background(), parse.Input{SourceID: "generated.horos", TreeID: "reparsed", Text: printed})
	if err != nil {
		t.Fatal("the printed text does not parse:", err)
	}
	pin := nodeAttribute(module(tree), hs.AttributePin)
	if pin == strings.Repeat("0", 64) || pin != nodeAttribute(module(parsed.ASTView()), hs.AttributePin) {
		t.Fatalf("the target gave the module the pin %s, and the parser assigns %s", pin, nodeAttribute(module(parsed.ASTView()), hs.AttributePin))
	}
	source := tsSourceDeclarations(t, q.source, "Person", "Flag")
	provenance := map[ast.NodeRef]string{}
	for _, record := range out.Inventory() {
		if len(record.Provenance) != 1 || record.Role != g.RequestedDeclaration || record.Site.Tree != tree.ID {
			t.Fatalf("a declaration of the output has the record %+v", record)
		}
		node := tree.Nodes[record.Site.Node-1]
		if node.Kind != hs.KindTypeDefinition {
			t.Fatalf("the inventory names a node of kind %s", node.Kind)
		}
		provenance[record.Provenance[0]] = nodeAttribute(node, hs.AttributeName)
	}
	if len(provenance) != 2 || provenance[source[0]] != "Person" || provenance[source[1]] != "Flag" {
		t.Fatalf("the inventory relates the TypeScript declarations to %v", provenance)
	}
}

// A type alias of an interface.
//
// Shows: a type alias whose right side names an interface plans as that
// interface's object type. The TypeScript analysis gives the alias the type
// of the interface, with the interface's properties, so the reader meets an
// alias of an object type: the definition is written out, and it does not
// refer to a definition of the interface. That is what a TypeScript alias
// means, and it is why this case is no refusal although the pair refuses a
// reference to a named type in a property.
// Does not show: an alias of a generic interface or of a package type. The
// reader judges both unsupported, and neither is run here.
func TestTypescriptToHorosAliasOfInterface(t *testing.T) {
	q := composeTSH(t, "export interface P { name: string }\nexport type T = P;\n", []string{"T"}, nil)
	out, err := q.registry.Generate(context.Background(), q.request, q.checker)
	if err != nil {
		t.Fatal(err)
	}
	const want = "module example.people;\n\ntype T = {\n    name: text;\n};\n"
	if printed := formatted(t, horosText(), outputTree[hs.Language](t, out)); printed != want {
		t.Fatalf("printed:\n%s\nwant:\n%s", printed, want)
	}
}

// tsHRefusal is one case of what the pair does not cover.
type tsHRefusal struct {
	name, source string
	extra        map[string]*tsHBundle
	// A refusal of a contribution or of the kernel: its class, the owner and
	// the reason. For a judgment of the source reader class is zero, and
	// reason is the reader's diagnostic. invalid is a description that the
	// horos target does not admit.
	class         g.RefusalClass
	invalid       bool
	owner, reason string
	contribution  g.ContributionReference
}

func tsHRefusals() []tsHRefusal {
	publisher := tsrules.Publisher
	object, declarations, number := tshobjectvocabulary.Reference, tshdeclarationsvocabulary.Reference, tshnumbervocabulary.Reference
	reads := "this TypeScript source reader reads object types and primitive types, and not the type "
	return []tsHRefusal{
		// Numbers first: the pair's own refusal, with the accepted page's reason.
		{name: "a number property", source: "export interface T { count: number }", class: g.RefusalUnsupported, owner: publisher, reason: tshnumber.ReasonNumber, contribution: number},
		// What a contribution of the pair recognizes and does not cover.
		{name: "an optional property", source: "export interface T { name?: string }", class: g.RefusalUnsupported, owner: publisher, reason: tshobject.ReasonOptionalProperty, contribution: object},
		{name: "a readonly property", source: "export interface T { readonly name: string }", class: g.RefusalUnsupported, owner: publisher, reason: tshobject.ReasonReadonlyProperty, contribution: object},
		{name: "an object type in a property", source: "export interface T { home: { city: string } }", class: g.RefusalUnsupported, owner: publisher, reason: tshobject.ReasonNestedObject, contribution: object},
		{name: "a type alias of a string", source: "export type T = string;", class: g.RefusalUnsupported, owner: publisher, reason: tshdeclarations.ReasonNoObjectType, contribution: declarations},
		// A primitive type for which the pair has no contribution: the reader
		// names it, and nothing answers for it.
		{name: "a bigint property", source: "export interface T { big: bigint }", class: g.RefusalMissingContribution, owner: g.KernelOwner, reason: g.ReasonNoContribution + " form:type/primitive/bigint"},
		{name: "a null property", source: "export interface T { nothing: null }", class: g.RefusalMissingContribution, owner: g.KernelOwner, reason: g.ReasonNoContribution + " form:type/primitive/null"},
		{name: "an any property", source: "export interface T { thing: any }", class: g.RefusalMissingContribution, owner: g.KernelOwner, reason: g.ReasonNoContribution + " form:type/primitive/any"},
		// What the TypeScript source reader does not read: a judgment of the
		// source, with the verdict unsupported.
		{name: "an array property", source: "export interface T { names: string[] }", reason: reads + "string[]"},
		{name: "a union property", source: "export interface T { either: string | boolean }", reason: reads + "string | boolean"},
		{name: "a nullable property", source: "export interface T { name: string | null }", reason: reads + "string | null"},
		{name: "a literal property", source: "export interface T { kind: \"a\" }", reason: reads + "\"a\""},
		{name: "a tuple property", source: "export interface T { pair: [string, boolean] }", reason: reads + "[string, boolean]"},
		{name: "a reference to a named type", source: "export interface P { name: string }\nexport interface T { lead: P }", reason: "this TypeScript source reader reads no reference to a named type or a type parameter"},
		{name: "a generic declaration", source: "export interface T<E> { name: string }", reason: "this TypeScript source reader reads no generic declaration"},
		// A name that TypeScript admits and horos does not: no refusal of the
		// pair, the horos target does not admit the description.
		{name: "a name that is no horos identifier", source: "export interface $T { name: string }", invalid: true},
	}
}

// Each refusal of the pair.
//
// Shows, for each case: planning fails and gives no plan; a refusal of a
// contribution or of the kernel has exactly the class, the owner, the reason
// and the selected contribution listed, matches the contract value of its
// class and no other of the eight; a judgment of the TypeScript source reader
// fails with the value of the verdict unsupported, is no refusal of a
// contribution, and carries the reader's diagnostic. The number is refused by
// the pair's number contribution as unsupported, with the sentence of
// docs/generation.md as its reason. A TypeScript name that is no horos
// identifier is no refusal: the horos target's check of the declaration
// fails as an invalid description.
// Does not show: that the list is complete. TypeScript has more type
// expressions than the listed ones (intersections, readonly arrays,
// conditional and mapped types, package types), and none of them is run
// here. It does not show a refusal for a string that is no well-formed
// Unicode text: the profile states that limit, and no type-level rule can
// see it.
func TestTypescriptToHorosRefusals(t *testing.T) {
	values := []error{g.ErrUnsupported, g.ErrRefusedByDesign, g.ErrRequirementUnsatisfied, g.ErrRequirementUnresolved, g.ErrMissingFacet, g.ErrAmbiguous, g.ErrMissingContribution, g.ErrLimit}
	for _, refusal := range tsHRefusals() {
		t.Run(refusal.name, func(t *testing.T) {
			root := "T"
			if refusal.invalid {
				root = "$T"
			}
			q := composeTSH(t, refusal.source, []string{root}, refusal.extra)
			plan, err := q.plan()
			if err == nil || plan.Valid() {
				t.Fatal("the request planned")
			}
			var got *g.Refusal
			if refusal.invalid {
				if errors.As(err, &got) || !errors.Is(err, g.ErrInvalid) {
					t.Fatalf("not an invalid description: %v", err)
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
			if got.Class != refusal.class || got.Owner != refusal.owner || got.Reason != refusal.reason || got.Contribution != refusal.contribution || got.Target != string(hs.Language{}.ID()) || got.Profile != tsrules.CanonicalizationProfile {
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
// the same bytes; the plan names TypeScript as its source with the analysis
// profile of the TypeScript vocabulary as the source semantics, horos as its
// TARGET and the pair's profile; its contribution entries are the four
// contributions that were selected, each with the pair's publisher, a local
// key under the pair's name, the meaning of its family of the TypeScript
// source and the schema of its vocabulary. The number contribution is
// registered and not in the plan: nothing selected it. Each declaration node
// records that it is transparent.
// Does not show: a reader of the plan in another process or another
// language. The plan has no binding product.
func TestTypescriptToHorosPlan(t *testing.T) {
	q := composeTSH(t, tsHSource, []string{"Person", "Flag"}, nil)
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
	if read.Source() != (g.PlanSource{Language: string(ts.Language{}.ID()), Semantics: tsn.AnalysisProfile}) || read.Target().Language != string(hs.Language{}.ID()) || read.Profile() != tsrules.CanonicalizationProfile {
		t.Fatalf("the plan names the source %+v, the target %+v and the profile %q", read.Source(), read.Target(), read.Profile())
	}
	want := map[string]struct {
		reference g.ContributionReference
		meaning   g.Meaning
		schema    g.SchemaReference
	}{
		"objecttype":   {tshobjectvocabulary.Reference, tsobjectfamily.Meaning, tshobjectvocabulary.Schema},
		"stringtype":   {tshstringvocabulary.Reference, tsstringfamily.Meaning, tshstringvocabulary.Schema},
		"booleantype":  {tshbooleanvocabulary.Reference, tsbooleanfamily.Meaning, tshbooleanvocabulary.Schema},
		"declarations": {tshdeclarationsvocabulary.Reference, tsdeclarationsfamily.Meaning, tshdeclarationsvocabulary.Schema},
	}
	entries := read.Contributions()
	if len(entries) != len(want) {
		t.Fatalf("the plan names %d contributions", len(entries))
	}
	for _, entry := range entries {
		name, found := strings.CutPrefix(entry.Reference.Identity.Key, "typescript-to-horos/")
		expected, known := want[name]
		if !found || !known || entry.Reference != expected.reference || entry.Reference.Identity.Publisher != tsrules.Publisher || entry.Reference.Meaning != expected.meaning || entry.Schema != expected.schema || entry.Target != string(hs.Language{}.ID()) || entry.Profile != tsrules.CanonicalizationProfile {
			t.Errorf("the plan names the contribution %+v", entry)
		}
		delete(want, name)
	}
	roots := read.Roots()
	if len(roots) != 2 || roots[0] != "source:declaration/0.ts/Person" || roots[1] != "source:declaration/0.ts/Flag" {
		t.Fatalf("roots %v", roots)
	}
	declarations := 0
	for _, node := range read.Nodes() {
		if node.Kind == "declaration" {
			declarations++
			if !node.Transparent {
				t.Errorf("the declaration node of %s is not transparent", node.Symbol)
			}
		}
	}
	if declarations != 2 {
		t.Fatalf("%d declaration nodes", declarations)
	}
}
