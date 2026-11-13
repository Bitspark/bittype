package integration_test

// Row 10 of the contribution contract (bittype#196,
// docs/contribution-contract.md, "Refusals"): every refusal of a builtin rule
// and every refusal of the kernel matches exactly one of the eight values. The
// test finds the packages by path, finds what constructs a refusal in each
// file and compares what it finds with its table, so that a new refusal
// cannot stay undriven.

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	gomodel "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// gateRefusalDirectories finds the packages whose refusals row 10 is about,
// by path below the root of the repository: the mapping packages of every
// pair, the mapping packages of every source, the shared parts of
// contributions and of targets, the target of every language, the contract
// and the registry. A new pair or a new language is found without an edit of
// a list.
func gateRefusalDirectories(t *testing.T, root string) []string {
	t.Helper()
	var found []string
	for _, pattern := range [][]string{
		{"transform", "syntax", "*-to-*", "mapping", "*", "go"},
		{"transform", "syntax", "*", "mapping", "*", "go"},
		{"transform", "syntax", "*", "mapping", "go"},
		{"transform", "syntax", "shared", "contribution", "go"},
		{"transform", "syntax", "*", "target", "go"},
		{"contracts", "generation", "go"},
		{"transform", "syntax", "shared", "generation", "go"},
	} {
		matches, err := filepath.Glob(filepath.Join(append([]string{root}, pattern...)...))
		if err != nil {
			t.Fatal(err)
		}
		for _, match := range matches {
			if info, err := os.Stat(match); err == nil && info.IsDir() && !slices.Contains(found, match) {
				found = append(found, match)
			}
		}
	}
	slices.Sort(found)
	return found
}

// gateRefusalValueNames are the names of the eight values in the contract
// package.
var gateRefusalValueNames = []string{"ErrUnsupported", "ErrRefusedByDesign", "ErrRequirementUnsatisfied", "ErrRequirementUnresolved", "ErrMissingFacet", "ErrAmbiguous", "ErrMissingContribution", "ErrLimit"}

// gateConstruction is one place that constructs a refusal: its file below the
// root, its form, its line and the text of that line. The forms are "Refuse"
// (a call of the contract's Refuse), "literal" (a composite literal of the
// type Refusal) and "refusal" (a call of the planner's method refusal).
type gateConstruction struct {
	file, form string
	line       int
	text       string
}

// gateRefusalScan is what the scan finds in the files that are no tests: the
// constructions of a refusal, and for each file the lines that use one of the
// eight values as a bare value.
type gateRefusalScan struct {
	constructions []gateConstruction
	bare          map[string][]int
	packages      []string
}

func gateScanRefusals(t *testing.T, root string) gateRefusalScan {
	t.Helper()
	scan := gateRefusalScan{bare: map[string][]int{}}
	fset := token.NewFileSet()
	for _, directory := range gateRefusalDirectories(t, root) {
		files, err := filepath.Glob(filepath.Join(directory, "*.go"))
		if err != nil {
			t.Fatal(err)
		}
		relativeDirectory, err := filepath.Rel(root, directory)
		if err != nil {
			t.Fatal(err)
		}
		scan.packages = append(scan.packages, filepath.ToSlash(relativeDirectory))
		for _, file := range files {
			if strings.HasSuffix(file, "_test.go") {
				continue
			}
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
			parsed, err := parser.ParseFile(fset, file, content, 0)
			if err != nil {
				t.Fatal(err)
			}
			relative, err := filepath.Rel(root, file)
			if err != nil {
				t.Fatal(err)
			}
			relative = filepath.ToSlash(relative)
			construct := func(form string, at token.Pos) {
				line := fset.Position(at).Line
				scan.constructions = append(scan.constructions, gateConstruction{relative, form, line, strings.TrimSpace(lines[line-1])})
			}
			name := func(e ast.Expr) string {
				switch e := e.(type) {
				case *ast.Ident:
					return e.Name
				case *ast.SelectorExpr:
					return e.Sel.Name
				}
				return ""
			}
			declared := map[token.Pos]bool{}
			ast.Inspect(parsed, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.ValueSpec:
					// The declaration of a value is no use of it.
					for _, declaredName := range n.Names {
						declared[declaredName.Pos()] = true
					}
				case *ast.CallExpr:
					switch name(n.Fun) {
					case "Refuse":
						construct("Refuse", n.Pos())
					case "refusal":
						construct("refusal", n.Pos())
					}
				case *ast.CompositeLit:
					if name(n.Type) == "Refusal" {
						construct("literal", n.Pos())
					}
				case *ast.SelectorExpr:
					if slices.Contains(gateRefusalValueNames, n.Sel.Name) {
						scan.bare[relative] = append(scan.bare[relative], fset.Position(n.Pos()).Line)
					}
					return false
				case *ast.Ident:
					if slices.Contains(gateRefusalValueNames, n.Name) && !declared[n.Pos()] {
						scan.bare[relative] = append(scan.bare[relative], fset.Position(n.Pos()).Line)
					}
				}
				return true
			})
		}
	}
	return scan
}

// gateRefusalSite is one row of the table of row 10: the places of one file
// that construct one refusal, found by a text that the line of each holds,
// and the cases that drive them.
type gateRefusalSite struct {
	id, file, form string
	// fragment is a text on the line of each place of this row, and on the line
	// of no other construction of the same form in the file.
	fragment string
	count    int
	// constructor: the place makes the value for its callers and refuses
	// nothing itself. It has no case.
	constructor bool
	// cases are the cases of TestGateEveryRefusalMatchesExactlyOneValue that
	// drive the row. Each case belongs to one row.
	cases []string
	// elsewhere are further tests of this package that drive the row.
	elsewhere []string
}

const (
	gateContractFile     = "contracts/generation/go/contribution.go"
	gatePlannerFile      = "contracts/generation/go/contribution_plan.go"
	gateRegistryFile     = "transform/syntax/shared/generation/go/registry.go"
	gateRecordPolicyFile = "transform/syntax/horos/mapping/record/go/record.go"
	gateDeclarationsFile = "transform/syntax/horos/mapping/declarations/go/contribution.go"
	gateGoRecordFile     = "transform/syntax/horos-to-golang/mapping/record/go/contribution.go"
)

// gateRowTenPairs are the four pairs from horos whose builtin rules the test
// drives. This is a list because driving a rule needs the Go types of a
// pair's target, bundle and checker (O46). The scan of the refusals is no
// list: it finds every package by path, so a refusal of a fifth pair is found
// and fails the comparison until it has a row and a case.
var gateRowTenPairs = []string{"golang", "typescript", "rust", "cpp"}

func gatePerPair(prefix string) []string {
	var cases []string
	for _, pair := range gateRowTenPairs {
		cases = append(cases, prefix+"/"+pair)
	}
	return cases
}

// gateRefusalSites is the table of row 10: every place that constructs a
// refusal in the packages that gateRefusalDirectories finds.
var gateRefusalSites = []gateRefusalSite{
	// The rules of the builtin contributions.
	{id: "record policy: a nested record", file: gateRecordPolicyFile, form: "Refuse", fragment: "ReasonNested", count: 1, cases: gatePerPair("nested-record")},
	{id: "declaration policy: a name that the naming rule does not cover", file: gateDeclarationsFile, form: "Refuse", fragment: "g.RefusalUnsupported, owner, err.Error()", count: 2,
		cases: append(gatePerPair("declaration-name"), gatePerPair("parameter-name")...)},
	{id: "golang record: a field name that the naming rule does not cover", file: gateGoRecordFile, form: "Refuse", fragment: "g.RefusalUnsupported, publisher, err.Error()", count: 1,
		cases: []string{"field-name/step-3", "field-name/step-1"}},
	// The kernel.
	{id: "contract: Refuse", file: gateContractFile, form: "literal", fragment: "Class: class, Owner: owner, Reason: reason}", count: 1, constructor: true},
	{id: "planner: refusal", file: gatePlannerFile, form: "literal", fragment: "Class: class, Owner: owner, Reason: reason, Subject: subject", count: 1, constructor: true},
	{id: "planner: a missing facet", file: gatePlannerFile, form: "literal", fragment: "Class: RefusalMissingFacet", count: 1, cases: []string{"missing-facet"},
		elsewhere: []string{"TestGuarantee1RequiredFacetsBelongToTheSelectedContribution"}},
	{id: "planner: no contribution answers a domain", file: gatePlannerFile, form: "refusal", fragment: "RefusalMissingContribution", count: 1,
		cases: []string{"missing/unbuilt-family", "missing/unanswered-role", "missing/no-declaration-policy"}},
	{id: "planner: several contributions answer a domain", file: gatePlannerFile, form: "refusal", fragment: "RefusalAmbiguousContribution", count: 1,
		cases: []string{"ambiguous/form", "ambiguous/role", "ambiguous/declaration-policy"}},
	{id: "planner: a specialization", file: gatePlannerFile, form: "refusal", fragment: "ReasonSpecialization", count: 1, cases: []string{"specialization"}},
	{id: "planner: a request for a shared helper", file: gatePlannerFile, form: "refusal", fragment: "ReasonSharedHelper", count: 1, cases: []string{"shared-helper"},
		elsewhere: []string{"TestPlanningProjectionScopeAndNoFallback"}},
	{id: "planner: a target without a member vocabulary", file: gatePlannerFile, form: "refusal", fragment: "ReasonNoMemberVocabulary", count: 1, cases: []string{"no-member-vocabulary"}},
	{id: "planner: a member fragment that the target does not admit", file: gatePlannerFile, form: "refusal", fragment: "ReasonMemberFragment", count: 1, cases: []string{"member-fragment"},
		elsewhere: []string{"TestGateMemberFragmentIsTheTargetsFieldForm"}},
	{id: "planner: an actual that does not discharge a requirement", file: gatePlannerFile, form: "refusal", fragment: "ReasonRequirement", count: 1,
		cases: []string{"requirement/unresolved", "requirement/unsatisfied"}, elsewhere: []string{"TestGateGenericAndWrittenOutAgree"}},
	{id: "planner: a cycle", file: gatePlannerFile, form: "refusal", fragment: "RefusalUnsupported, KernelOwner, reason,", count: 1,
		cases: []string{"cycle/evidence", "cycle/recursive-reference"}, elsewhere: []string{"TestGateEvidenceCycleAndRecursiveReference"}},
	{id: "planner: an exact mapping where a declaration is published", file: gatePlannerFile, form: "refusal", fragment: "ReasonExactWithoutPolicy", count: 1, cases: []string{"exact-without-policy"}},
	{id: "registry: no contribution states the profile", file: gateRegistryFile, form: "literal", fragment: "Class: g.RefusalMissingContribution", count: 1, cases: []string{"missing/no-profile"}},
}

// gateBareValues is, for each file that the scan finds, the number of places
// that use one of the eight values as a bare value, and what they are. None
// constructs a refusal. The count is compared so that a new return of a bare
// value is seen. Where a place returns the value to the caller of the planner,
// the error matches the value and is no refusal: it has no owner, reason or
// subject.
var gateBareValues = map[string]struct {
	count int
	what  string
}{
	gateContractFile:                                         {8, "RefusalClass.Value: the value of each of the eight classes"},
	gatePlannerFile:                                          {9, "seven returns of the bare limit, the wrapped ambiguity of two codecs of one schema, and the comparison that turns a missing facet into its refusal"},
	"contracts/generation/go/facet.go":                       {1, "matchFacet answers the planner, which makes the refusal of a missing facet from it"},
	"contracts/generation/go/occurrence.go":                  {1, "a limitation of the kernel that it locates at a site: the bare unsupported value in a SiteError"},
	"contracts/generation/go/output.go":                      {3, "three returns of the bare limit when an output is sealed, after planning"},
	"contracts/generation/go/planning.go":                    {4, "three returns of the bare limit (operations and depth, reservations, carried values) and the bare unsupported value for a target parameter that its declaration must use"},
	"contracts/generation/go/source.go":                      {1, "the verdict Unsupported of a source judgment, which stays the source's"},
	"transform/syntax/cpp/target/go/provider.go":             {2, "the target's own refusal of intent that it cannot build: the bare unsupported value"},
	"transform/syntax/golang/target/go/provider.go":          {1, "the target's own refusal of intent that it cannot build: the bare unsupported value"},
	"transform/syntax/rust/target/go/provider.go":            {1, "the target's own refusal of intent that it cannot build: the bare unsupported value"},
	"transform/syntax/typescript/target/go/provider.go":      {1, "the target's own refusal of intent that it cannot build: the bare unsupported value"},
	"transform/syntax/shared/generation/go/contributions.go": {6, "four refusals of a second registration (a codec, a facet, a product, a replacement) with the bare ambiguous value, a product that is not registered as the wrapped unsupported value, and the bare limit for more contributions than the request's limit"},
	"transform/syntax/shared/generation/go/diagnostics.go":   {3, "comparisons that choose the code of a diagnostic"},
	"transform/syntax/shared/generation/go/registry.go":      {3, "two refusals of a second registration (a source, a target) with the bare ambiguous value, and the bare limit for more roots than the request's limit"},
	"transform/syntax/shared/target/go/builder.go":           {2, "two returns of the bare limit while a target builds, after planning"},
	"transform/syntax/shared/target/go/check.go":             {3, "the name rules of a target's program: the bare unsupported value"},
	"transform/syntax/shared/target/go/compile.go":           {5, "four returns of the bare limit and one of the bare unsupported value while a target builds, after planning"},
}

// gateDriven is one driven case: what the planner returned and what the case
// expects of it.
type gateDriven struct {
	id string
	// class is the one refusal value that the error matches.
	class string
	// owner and reason are the refusal's. The reason is compared as a prefix:
	// the kernel appends the domain or the requirement to some of its reasons.
	owner, reason string
	// contribution is the selected contribution that the refusal names, where
	// the case states one.
	contribution *g.ContributionReference
	// bare: the error matches the value and is no refusal.
	bare bool
	err  error
	// planned: the planner returned a plan.
	planned bool
}

// gateBareTarget is a Go target that states no member vocabulary: it has the
// methods of a target and no other.
type gateBareTarget struct{ g.Target[gs.Language] }

// gateDriveBuiltin drives the refusals of the builtin rules of one pair from
// horos: a nested record, a declaration name and a parameter name that the
// naming rule does not cover. For the pair whose record spells a member by
// the naming rule it also drives a field name, in step 3 and, for an optional
// field, in step 1. For the pairs that keep a field's name as the source
// states it, it checks that the same modules plan.
func gateDriveBuiltin[T any](t *testing.T, pair string, spellsByNamingRule bool, bundle func(*testing.T, representation.Language) gateBundle[T], plan func(*testing.T, transform.Input, gateBundle[T], gateRequest) (g.GenerationPlan, error), driven *[]gateDriven) {
	t.Helper()
	run := func(id, module, key, reason string) {
		from, input := gateSource(t, module)
		b := bundle(t, from)
		reference := b.find(t, "horos-to-"+pair+"/"+key).Reference
		generation, err := plan(t, input, b, gateRequest{})
		*driven = append(*driven, gateDriven{id: id, class: "unsupported", owner: reference.Identity.Publisher, reason: reason, contribution: &reference, err: err, planned: generation.Representation().Valid()})
	}
	run("nested-record/"+pair, "module example.nested;\n\ntype Outer = { inner: { name: text; }; };\n", "record", "nested records require a named declaration in this profile")
	run("declaration-name/"+pair, "module example.names;\n\ntype a__b = { x: text; };\n", "declarations", `unsupported identifier "a__b"`)
	run("parameter-name/"+pair, "module example.names;\n\ntype Page<t__u> = { items: list<t__u>; };\n", "declarations", `unsupported identifier "t__u"`)
	const required, optional = "module example.names;\n\ntype R = { a__b: text; };\n", "module example.names;\n\ntype R = { a__b?: text; };\n"
	if spellsByNamingRule {
		run("field-name/step-3", required, "record", `unsupported identifier "a__b"`)
		run("field-name/step-1", optional, "record", `unsupported identifier "a__b"`)
		return
	}
	for _, module := range []string{required, optional} {
		from, input := gateSource(t, module)
		if generation, err := plan(t, input, bundle(t, from), gateRequest{}); err != nil || !generation.Representation().Valid() {
			t.Errorf("%s keeps a field's name as the source states it and does not plan %q: %v", pair, module, err)
		}
	}
}

// gateDriveKernel drives the refusals of the kernel on the independent
// fixture source.
func gateDriveKernel(t *testing.T, driven *[]gateDriven) {
	t.Helper()
	fixture := gateFixture(t)
	kernel := func(id, class, reason string, contribution *g.ContributionReference, composition f.Composition, roots ...string) {
		generation, err := fixture.PlanContributions(context.Background(), composition, roots...)
		*driven = append(*driven, gateDriven{id: id, class: class, owner: g.KernelOwner, reason: reason, contribution: contribution, err: err, planned: generation.Representation().Valid()})
	}
	base := func(o indOptions) f.Composition { return indComposition(t, fixture, o) }
	without := func(key string) f.Composition {
		composition := base(indOptions{})
		composition.Contributions = slices.Delete(composition.Contributions, indFind(t, composition, key), indFind(t, composition, key)+1)
		return composition
	}
	twice := func(key string) f.Composition {
		composition := base(indOptions{})
		second := composition.Contributions[indFind(t, composition, key)]
		second.Reference.Identity.Key += "-again"
		composition.Contributions = append(composition.Contributions, second)
		return composition
	}
	reference := func(key string) *g.ContributionReference {
		r := gateReference(key, "golang/1")
		return &r
	}

	kernel("missing/unbuilt-family", "missing contribution", g.ReasonNoContribution+" form:expression/unbuilt-family", nil, base(indOptions{}), "Unbuilt")
	kernel("missing/unanswered-role", "missing contribution", g.ReasonNoContribution+" role:"+f.OptionalNullableRole, nil, without("mock-optional-nullable"), "Outcome")
	kernel("missing/no-declaration-policy", "missing contribution", g.ReasonNoContribution+" policy", nil, without("independent-declarations"), "Word")
	kernel("ambiguous/form", "ambiguous contribution", g.ReasonSeveralContributions+" form:expression/word", nil, twice("independent-word"), "Word")
	kernel("ambiguous/role", "ambiguous contribution", g.ReasonSeveralContributions+" role:"+f.OptionalNullableRole, nil, twice("mock-optional-nullable"), "Outcome")
	kernel("ambiguous/declaration-policy", "ambiguous contribution", g.ReasonSeveralContributions+" policy", nil, twice("independent-declarations"), "Word")

	specialized := base(indOptions{})
	specialization := indEmpty("specialized-page", g.ForSpecialization(indSelector(t, fixture, "declaration/Page"), indSelector(t, fixture, "declaration/Person")), nil,
		func(context.Context, goProjection) (goProjected, error) {
			return goProjected{Type: gomodel.Primitive(gomodel.PrimitiveString)}, nil
		})
	specialized.Contributions = append(specialized.Contributions, specialization)
	kernel("specialization", "unsupported", g.ReasonSpecialization, reference("specialized-page"), specialized, "People")

	helping := base(indOptions{})
	helping.Contributions[indFind(t, helping, "independent-word")].Rule.Name = func(context.Context, indNameStep) ([]indRequest, error) {
		return []indRequest{g.SharedHelper[f.SourceLanguage]("a helper of another contribution")}, nil
	}
	kernel("shared-helper", "unsupported", g.ReasonSharedHelper, reference("independent-word"), helping, "Word")

	bare := base(indOptions{})
	bare.Target = gateBareTarget{ggtarget.Provider{Package: "generated"}}
	kernel("no-member-vocabulary", "unsupported", g.ReasonNoMemberVocabulary, reference("mock-optional-nullable"), bare, "Outcome")

	kernel("requirement/unresolved", "requirement unresolved", g.ReasonRequirement, reference("independent-declarations"), base(indOptions{probeRequires: true}), "ProbedPerson")
	kernel("requirement/unsatisfied", "requirement unsatisfied", g.ReasonRequirement, reference("independent-word"), base(indOptions{probeRequires: true, wordFacts: []g.Fact{gateNilFree.Refuted("")}}), "ProbedWord")
	kernel("cycle/evidence", "unsupported", g.ReasonEvidenceCycle, nil, base(indOptions{}), "Loop")
	kernel("cycle/recursive-reference", "unsupported", g.ReasonRecursiveReference, nil, base(indOptions{}), "Rec")
	kernel("exact-without-policy", "unsupported", g.ReasonExactWithoutPolicy, reference("independent-cell"), base(indOptions{}), "Cell")

	kind, err := g.NewFacetKind(g.FacetKindDescription{Identity: gatePublisher + "/facet/row-ten", Version: "1"})
	if err != nil {
		t.Fatal(err)
	}
	product := g.Product{Identity: gatePublisher + "/row-ten", Requires: map[g.SubjectKind][]g.FacetKind{g.UseSubject: {kind}}}
	faceted := base(indOptions{})
	faceted.Products, faceted.Requested = []g.Product{product}, []string{product.Identity}
	kernel("missing-facet", "missing facet", g.ReasonFacet, reference("independent-word"), faceted, "Word")
}

// gateDriveHoros drives, on the horos source and the Go target, the two
// refusals that need them: a member fragment that the target's member
// vocabulary does not admit, and a request whose profile no registered
// contribution states.
func gateDriveHoros(t *testing.T, driven *[]gateDriven) {
	t.Helper()
	_, from, input := generationInput(t, paginationModule())
	fragment := gateFragment("row-ten-member", goTarget, func(_ string, operand g.Expression[gs.Language]) g.Expression[gs.Language] {
		return gomodel.Collection(operand)
	})
	generation, err := gatePlanGo(t, input, gateGolang(t, from).replace(fragment).withCodecs(gateEmptyCodec(t)), gateRequest{})
	*driven = append(*driven, gateDriven{id: "member-fragment", class: "unsupported", owner: "member vocabulary of " + goTarget, reason: g.ReasonMemberFragment, contribution: &fragment.Reference, err: err, planned: generation.Representation().Valid()})

	other := gateGolang(t, from)
	for i := range other.contributions {
		other.contributions[i].Profile = "example.test/another-profile/1"
	}
	generation, err = gatePlanGo(t, input, other, gateRequest{})
	*driven = append(*driven, gateDriven{id: "missing/no-profile", class: "missing contribution", owner: g.KernelOwner, reason: g.ReasonNoProfile + " " + rules.RepresentationProfile, err: err, planned: generation.Representation().Valid()})
}

// gateDriveLimits drives the bounds of a request on the horos source and the
// Go target. Each returns the bare limit.
func gateDriveLimits(t *testing.T, driven *[]gateDriven) {
	t.Helper()
	from, input := gateSource(t, "module example.limits;\n\ntype Person = { name: text; active: bool; };\ntype Pair<A, B> = { first: A; second: B; };\ntype People = Pair<Person, Person>;\n")
	for _, c := range []struct {
		id    string
		roots []string
		bound func(*g.Limits)
	}{
		{"limit/roots", nil, func(l *g.Limits) { l.References = 2 }},
		{"limit/contributions", nil, func(l *g.Limits) { l.Plans = 1 }},
		{"limit/requests-of-a-rule", []string{"Person"}, func(l *g.Limits) { l.References = 1 }},
		{"limit/operations", nil, func(l *g.Limits) { l.Operations = 3 }},
		{"limit/depth", nil, func(l *g.Limits) { l.Depth = 1 }},
		{"limit/binders", []string{"Pair"}, func(l *g.Limits) { l.Formals = 1 }},
		{"limit/declarations", []string{"People"}, func(l *g.Limits) { l.Reservations = 1 }},
		{"limit/frozen-bytes", []string{"Person"}, func(l *g.Limits) { l.OutputBytes = 8 }},
	} {
		limits := g.DefaultLimits()
		c.bound(&limits)
		generation, err := gatePlanGo(t, input, gateGolang(t, from), gateRequest{roots: c.roots, limits: &limits})
		*driven = append(*driven, gateDriven{id: c.id, class: "limit", bare: true, err: err, planned: generation.Representation().Valid()})
	}
}

// Row 10 (O17, O26, O40), gate line 2: "Every refusal of the mapping
// packages matches exactly one value."
//
// Shows: every place that constructs a refusal in the mapping packages of
// every pair and of every source, in the shared parts of contributions and of
// targets, in the target of every language, in the contract and in the
// registry is a row of the table gateRefusalSites, and each row that refuses
// is driven by at least one case. The places are found by path and by the
// syntax of the files, so a new refusal has no row and fails the test. Each
// driven refusal matches exactly one of the eight values, is a refusal with
// the owner and the reason that its row states, names the selected
// contribution where one was selected, and no plan results. The cases: for
// each of the four pairs from horos a nested record, a declaration name and a
// parameter name that the naming rule does not cover; for the Go pair a field
// name that the rule does not cover, in step 3 and in step 1; and of the
// kernel a missing contribution (an unbuilt family, an unanswered role, no
// declaration policy, no contribution of the profile), an ambiguous
// contribution (a form, a role, the declaration policy), a specialization, a
// request for a shared helper, a target without a member vocabulary, a member
// fragment that the vocabulary does not admit, an unresolved and an
// unsatisfied requirement, the two cycles, an exact mapping where a
// declaration is published, and a missing facet. Eight bounds of a request
// each return an error that matches exactly the value "limit". The places
// that use one of the eight values as a bare value are counted by file and
// compared with the table gateBareValues.
// An unexpected error of a rule is a rule failure of its contribution and
// matches no value: an error of its own, one of the eight values returned
// bare, a refusal of a class that only the kernel may use, and a refusal of
// no class.
// Does not show: the four classes that a rule may use, which no builtin rule
// but "unsupported" uses: TestGateSelectedContributionRefusalIsNotDropped
// drives them with a fixture. The bare values of a target's own checks and of
// a source judgment, which are counted and not driven here. That a limit is a
// refusal with an owner, a reason and a subject: it is the bare value
// (finding F-B2, TestGateLimitIsARefusalWithOwnerReasonAndSubject). Whether a
// rule may wrap its refusal: a wrapped refusal is asserted to match at most
// one value.
// Red under: R10-new-refusal-unseen, R10-new-bare-value-unseen,
// R10-nested-record-by-design, R10-declaration-name-is-an-error,
// R10-parameter-name-is-an-error, R10-golang-member-name-is-an-error,
// R10-missing-as-unsupported, R10-ambiguous-as-missing,
// R10-no-profile-as-unsupported, R10-specialization-by-design,
// R10-shared-helper-by-design, R10-no-member-vocabulary-by-design,
// R10-member-fragment-by-design, R10-exact-without-policy-by-design,
// R10-limit-as-unsupported, R10-refusal-matches-two,
// R10-rule-may-use-kernel-class, R10-rule-failure-unwraps.
func TestGateEveryRefusalMatchesExactlyOneValue(t *testing.T) {
	// The table against what the scan finds.
	scan := gateScanRefusals(t, filepath.Join("..", ".."))
	pairs := 0
	for _, directory := range scan.packages {
		if matched, _ := filepath.Match("transform/syntax/*-to-*/mapping/*/go", directory); matched {
			pairs++
		}
	}
	if pairs == 0 || !slices.Contains(scan.packages, "contracts/generation/go") || !slices.Contains(scan.packages, "transform/syntax/shared/generation/go") {
		t.Fatalf("the scan finds %d mapping packages of pairs in %v", pairs, scan.packages)
	}
	t.Logf("the scan reads %d packages, %d of them mapping packages of a pair, and finds %d constructions of a refusal", len(scan.packages), pairs, len(scan.constructions))
	found := map[string][]int{}
	for _, construction := range scan.constructions {
		var rows []string
		for _, site := range gateRefusalSites {
			if site.file == construction.file && site.form == construction.form && strings.Contains(construction.text, site.fragment) {
				rows = append(rows, site.id)
			}
		}
		if len(rows) != 1 {
			t.Errorf("%s:%d constructs a refusal (%s) that %d rows of the table name %v: %s", construction.file, construction.line, construction.form, len(rows), rows, construction.text)
			continue
		}
		found[rows[0]] = append(found[rows[0]], construction.line)
	}
	perFile := map[string]int{}
	for _, construction := range scan.constructions {
		perFile[construction.file]++
	}
	for _, file := range slices.Sorted(maps.Keys(perFile)) {
		t.Logf("constructions in %s: %d", file, perFile[file])
	}
	for file, lines := range scan.bare {
		if want, ok := gateBareValues[file]; !ok || want.count != len(lines) {
			t.Errorf("%s uses one of the eight values as a bare value at %d places (lines %v) and the table states %d", file, len(lines), lines, want.count)
		}
	}
	for _, file := range slices.Sorted(maps.Keys(gateBareValues)) {
		want := gateBareValues[file]
		if _, ok := scan.bare[file]; !ok {
			t.Errorf("the table states %d bare values in %s and the scan finds none", want.count, file)
		}
		t.Logf("bare values in %s: %d at lines %v (%s)", file, len(scan.bare[file]), scan.bare[file], want.what)
	}

	// The cases.
	var driven []gateDriven
	gateDriveBuiltin(t, "golang", true, gateGolang, gatePlanGo, &driven)
	gateDriveBuiltin(t, "typescript", false, gateTypescript, gatePlanTS, &driven)
	gateDriveBuiltin(t, "rust", false, gateRust, gatePlanRust, &driven)
	gateDriveBuiltin(t, "cpp", false, gateCpp, gatePlanCpp, &driven)
	gateDriveKernel(t, &driven)
	gateDriveHoros(t, &driven)
	gateDriveLimits(t, &driven)
	ran := map[string]int{}
	bounds := 0
	for _, d := range driven {
		if d.bare {
			bounds++
		}
	}
	t.Logf("the test drives %d cases: %d refusals and %d bounds of a request", len(driven), len(driven)-bounds, bounds)
	for _, d := range driven {
		ran[d.id]++
		if d.err == nil || d.planned {
			t.Errorf("%s: the planner returned a plan (error: %v)", d.id, d.err)
			continue
		}
		if d.bare {
			if matched := gateRefusalClasses(d.err); len(matched) != 1 || matched[0] != d.class {
				t.Errorf("%s: got %v, which matches %v; want exactly %q", d.id, d.err, matched, d.class)
			}
			continue
		}
		refusal := gateExpectRefusal(t, d.id, d.err, d.class)
		if refusal.Owner != d.owner || !strings.HasPrefix(refusal.Reason, d.reason) {
			t.Errorf("%s: the refusal has the owner %q and the reason %q, want %q and %q", d.id, refusal.Owner, refusal.Reason, d.owner, d.reason)
		}
		if d.contribution != nil && refusal.Contribution != *d.contribution {
			t.Errorf("%s: the refusal names the contribution %s, want %s", d.id, refusal.Contribution, *d.contribution)
		}
	}
	// Every row has its places and its cases, and every case has its row.
	tests := gateTestFunctions(t)
	rowOf := map[string]string{}
	for _, site := range gateRefusalSites {
		if len(found[site.id]) != site.count {
			t.Errorf("row %q: the scan finds %d places in %s (lines %v) and the table states %d", site.id, len(found[site.id]), site.file, found[site.id], site.count)
		}
		if site.constructor != (len(site.cases) == 0) {
			t.Errorf("row %q: a row has cases unless it is a constructor", site.id)
		}
		for _, id := range site.cases {
			if rowOf[id] != "" {
				t.Errorf("the case %s belongs to the rows %q and %q", id, rowOf[id], site.id)
			}
			rowOf[id] = site.id
			if ran[id] != 1 {
				t.Errorf("row %q: its case %s ran %d times", site.id, id, ran[id])
			}
		}
		for _, name := range site.elsewhere {
			if !slices.Contains(tests, name) {
				t.Errorf("row %q names the test %s, which this package does not have", site.id, name)
			}
		}
		t.Logf("row %q: %s lines %v (%s); cases %v; also %v", site.id, site.file, found[site.id], site.form, site.cases, site.elsewhere)
	}
	for _, d := range driven {
		if !d.bare && rowOf[d.id] == "" {
			t.Errorf("the case %s drives no row of the table", d.id)
		}
	}

	// An unexpected error of a rule is a rule failure and matches no value.
	fixture := gateFixture(t)
	word := gateReference("independent-word", "golang/1")
	returns := func(result error) error {
		composition := indComposition(t, fixture, indOptions{})
		composition.Contributions[indFind(t, composition, "independent-word")].Rule.Decide = func(context.Context, indDecideStep) (g.Decision, error) {
			return g.Decision{}, result
		}
		generation, err := fixture.PlanContributions(context.Background(), composition, "Word")
		if err == nil || generation.Representation().Valid() {
			t.Errorf("a rule that returns %v: the planner returned a plan", result)
		}
		return err
	}
	for _, c := range []struct {
		name   string
		result error
	}{
		{"an error of its own", errors.New("the rule fails")},
		{"the bare unsupported value", g.ErrUnsupported},
		{"the bare limit", g.ErrLimit},
		{"the bare missing contribution", g.ErrMissingContribution},
		{"a refusal of the class missing facet", g.Refuse(g.RefusalMissingFacet, gatePublisher, "not a rule's class")},
		{"a refusal of the class ambiguous contribution", g.Refuse(g.RefusalAmbiguousContribution, gatePublisher, "not a rule's class")},
		{"a refusal of the class missing contribution", g.Refuse(g.RefusalMissingContribution, gatePublisher, "not a rule's class")},
		{"a refusal of the class limit", g.Refuse(g.RefusalLimit, gatePublisher, "not a rule's class")},
		{"a refusal of no class", g.Refuse(0, gatePublisher, "no class")},
	} {
		err := returns(c.result)
		var failure *g.RuleFailure
		if !errors.As(err, &failure) || failure.Contribution != word || failure.Err != c.result {
			t.Errorf("a rule that returns %s: got %v, want the rule failure of %s with that cause", c.name, err, word)
		}
		if matched := gateRefusalClasses(err); len(matched) != 0 {
			t.Errorf("a rule that returns %s: the failure matches the refusal values %v", c.name, matched)
		}
	}
	if matched := gateRefusalClasses(returns(fmt.Errorf("the rule wraps: %w", g.Refuse(g.RefusalByDesign, gatePublisher, "wrapped")))); len(matched) > 1 {
		t.Errorf("a refusal that a rule wraps matches the values %v", matched)
	}
}

// gateTestFunctions lists the test functions of this package, read from its
// files.
func gateTestFunctions(t *testing.T) []string {
	t.Helper()
	files, err := filepath.Glob("*_test.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	var names []string
	for _, file := range files {
		parsed, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, declaration := range parsed.Decls {
			if function, ok := declaration.(*ast.FuncDecl); ok && function.Recv == nil && strings.HasPrefix(function.Name.Name, "Test") {
				names = append(names, function.Name.Name)
			}
		}
	}
	return names
}

// Row 10, gate line 2: a refusal carries "an owner-qualified reason and the
// subject", and "limit" is one of the eight classes of the one refusal type.
//
// Shows: a request whose bound is reached while a rule's requests are
// resolved is refused with a refusal of the class "limit" that has an owner,
// a reason and the subject.
// Does not show: which bound a reason names.
// Red under: no mutant. The test is skipped as PENDING: it is red against the
// kernel as built (finding F-B2), so no mutant can show it red.
func TestGateLimitIsARefusalWithOwnerReasonAndSubject(t *testing.T) {
	t.Skip("PENDING F-B2: a limit is returned as the bare value ErrLimit; no place constructs a refusal of the class limit, so a limit has no owner, no reason and no subject")
	from, input := gateSource(t, "module example.limits;\n\ntype Person = { name: text; active: bool; };\n")
	limits := g.DefaultLimits()
	limits.References = 1
	generation, err := gatePlanGo(t, input, gateGolang(t, from), gateRequest{limits: &limits})
	if err == nil || generation.Representation().Valid() {
		t.Fatalf("the request is planned within one reference: %v", err)
	}
	refusal := gateExpectRefusal(t, "a record with two fields within one reference", err, "limit")
	if refusal.Class != g.RefusalLimit || refusal.Owner == "" || refusal.Reason == "" || refusal.Subject == "" {
		t.Errorf("the limit is %+v, want a refusal of the class limit with an owner, a reason and the subject", refusal)
	}
}
