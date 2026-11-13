package integration_test

// Gate line 7, the kernel's part (bittype#196, docs/contribution-contract.md,
// "The gate", line 7, and "The extension domain"): what a replacement of one
// registered contribution changes in a plan and what it leaves, and one
// recorded refusal for each item of the extension domain that is out or not
// built. These tests state what each shows. They do not show that line 7 is
// met: the line also asks for a contribution added with no shared package
// changed, out of tree and in the test tree, which is not here (O45).

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	shared "github.com/Bitspark/bittype/integration/go/testdata/composition"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	gomodel "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
)

// gateDependsOn reports, for each node of a plan, whether it or a node that it
// reaches through its links, its application's declaration or its
// application's arguments, at any depth, was selected with the contribution.
func gateDependsOn(plan g.RepresentationPlan, identity g.ContributionIdentity) map[string]bool {
	nodes := map[string]g.PlanNode{}
	for _, node := range plan.Nodes() {
		nodes[node.ID] = node
	}
	memo := map[string]bool{}
	var visit func(string) bool
	visit = func(id string) bool {
		if known, ok := memo[id]; ok {
			return known
		}
		memo[id] = false
		node := nodes[id]
		depends := node.Contribution != nil && node.Contribution.Identity == identity
		for _, link := range node.Links {
			depends = depends || link.Node != "" && visit(link.Node)
		}
		if application := node.Application; application != nil {
			depends = depends || visit(application.Declaration)
			for _, argument := range append(slices.Clone(application.Target), application.Adapters...) {
				depends = depends || visit(argument.Node)
			}
		}
		memo[id] = depends
		return depends
	}
	for id := range nodes {
		visit(id)
	}
	return memo
}

// gateSelection is one selection as the comparison of two plans counts it:
// what was selected, for which source identity, and the frozen payload.
type gateSelection struct{ kind, source, contribution, payload string }

// gateSelections counts the selections of a plan that were not made with the
// contribution skip. A declaration that a contribution owns is left out: it
// belongs to its owner and goes with it.
func gateSelections(plan g.RepresentationPlan, skip g.ContributionIdentity) map[gateSelection]int {
	out := map[gateSelection]int{}
	for _, node := range plan.Nodes() {
		if node.Contribution != nil && node.Contribution.Identity != skip && node.Kind != "owned" {
			out[gateSelection{node.Kind, node.Source, node.Contribution.String(), node.Payload}]++
		}
	}
	return out
}

// gateComposed plans the request of a module's type definitions for Go in a
// fresh registry, as gatePlan does, after compose has done what a consumer's
// composition does beyond registering: an explicit replacement.
func gateComposed(t *testing.T, input transform.Input, b gateBundle[gs.Language], compose func(*c.Registry) error) (g.GenerationPlan, error) {
	t.Helper()
	var registry c.Registry
	checker := nativeChecker(t, golangText())
	from, to := input.Source.Language(), checker.Language()
	if err := c.RegisterSource[h.Language, source.Reader](&registry, from, source.Provider{}); err != nil {
		t.Fatal(err)
	}
	if err := c.RegisterTarget(&registry, to, ggtarget.Provider{Package: "generated"}); err != nil {
		t.Fatal(err)
	}
	if err := (shared.Set[h.Language, gs.Language, source.Reader]{Contributions: b.contributions, Codecs: b.codecs}).Register(&registry, from, to); err != nil {
		t.Fatal(err)
	}
	if err := compose(&registry); err != nil {
		return g.GenerationPlan{}, err
	}
	roots, err := genfixture.Roots(input)
	if err != nil {
		t.Fatal(err)
	}
	return registry.PlanGeneration(context.Background(), c.Request{Input: input, To: to, Output: "generated", Roots: roots, Profile: rules.RepresentationProfile, Limits: g.DefaultLimits()}, checker)
}

// Row 1 (O20, O26), gate line 7: "One registered contribution is replaced by
// another for the same domain."
//
// Shows: a consumer's composition has a fixture contribution in the place of
// the builtin Go contribution for the role of an optional field. The fixture
// represents the field as a slice, so the generated text differs. Between the
// plan with the builtin and the plan with the fixture, the payloads of every
// other selection are byte-identical. Every node that does not depend on the
// replaced contribution is the same node in both plans. Every node that
// depends on it is another node (O26): an answer and a link are bound to the
// dependency. Both groups are not empty. With both contributions registered
// the domain is ambiguous, whatever the order of registration, until the
// consumer replaces explicitly. After the explicit replacement the serialized
// plan is byte for byte the plan of the composition that never registered the
// builtin. A replacement names two registered contributions of one domain: one
// that names itself, an unregistered contribution or a contribution of
// another domain is a scope error, and a second replacement of the same
// contribution is refused.
// Does not show: that no shared package changes when a contribution is added.
// That is the diff of a commit and the module out of tree, which are not
// here. A replacement of a contribution for a use or of the declaration
// policy: the replaced domain here is a role.
// Red under: L7-replacement-ignored, L7-replacement-sets-aside-the-other,
// L7-two-contributions-first-wins, L7-replace-across-domains,
// L7-replace-twice-admitted, L7-node-depends-on-the-composition,
// R19-identity-ignores-links.
func TestGateReplacedContributionLeavesTheOthers(t *testing.T) {
	from, input := gateSource(t, `module example.replace;

type Person = { name: text; active: bool; };
type Page<T> = { items: list<T>; next?: text; };
type People = Page<Person>;
type Names = list<text>;
`)
	builtin := gateGolang(t, from)
	replaced := builtin.find(t, "horos-to-golang/optional").Reference.Identity
	replacement := gateFragment("slice-presence", goTarget, func(name string, operand g.Expression[gs.Language]) g.Expression[gs.Language] {
		return gomodel.Field(name, false, gomodel.Collection(operand))
	})
	if replacement.Domain.Key() != builtin.find(t, "horos-to-golang/optional").Domain.Key() {
		t.Fatal("the fixture does not answer the domain of the builtin optional contribution")
	}
	with := builtin.replace(replacement).withCodecs(gateEmptyCodec(t))
	beforeGeneration, err := gatePlanGo(t, input, builtin, gateRequest{})
	if err != nil {
		t.Fatal(err)
	}
	afterGeneration, err := gatePlanGo(t, input, with, gateRequest{})
	if err != nil {
		t.Fatal(err)
	}
	before, after := beforeGeneration.Representation(), afterGeneration.Representation()
	if _, selected := before.Contribution(builtin.find(t, "horos-to-golang/optional").Reference); !selected {
		t.Fatal("the module does not select the contribution that is replaced")
	}
	if _, selected := after.Contribution(replacement.Reference); !selected {
		t.Fatal("the composition with the fixture does not select it")
	}
	if a, b := gateText[gs.Language](t, golangText(), beforeGeneration), gateText[gs.Language](t, golangText(), afterGeneration); a == b {
		t.Error("the replacement generates the same text: it is no other representation")
	}
	// The payloads of every other selection.
	if a, b := gateSelections(before, replaced), gateSelections(after, replacement.Reference.Identity); !reflect.DeepEqual(a, b) || len(a) == 0 {
		t.Errorf("the payloads of the other selections changed:\n%v\n%v", a, b)
	} else {
		total := 0
		for _, count := range a {
			total += count
		}
		t.Logf("the payloads of the %d other selections are byte-identical", total)
	}
	// The frozen records.
	depends := gateDependsOn(before, replaced)
	same, other := 0, 0
	for _, node := range before.Nodes() {
		_, kept := after.Node(node.ID)
		switch {
		case node.Kind == "owned":
			// A declaration that the replaced contribution owns goes with it.
		case !depends[node.ID] && !kept:
			t.Errorf("the frozen record of %s %s changed, and it does not depend on the replaced contribution", node.Kind, node.Source)
		case !depends[node.ID]:
			same++
		case kept:
			t.Errorf("the node of %s %s depends on the replaced contribution and is the same node", node.Kind, node.Source)
		default:
			other++
		}
	}
	if same == 0 || other == 0 {
		t.Errorf("%d nodes do not depend on the replaced contribution and %d depend on it: the comparison needs both", same, other)
	}
	t.Logf("%d frozen records that do not depend on the replaced contribution are identical, and %d that depend on it are other nodes", same, other)
	for _, name := range []string{"Person", "Names"} {
		if node := gateDeclaration(t, before, name); depends[node.ID] {
			t.Errorf("%s depends on the replaced contribution", name)
		}
	}
	for _, name := range []string{"Page", "People"} {
		if node := gateDeclaration(t, before, name); !depends[node.ID] {
			t.Errorf("%s does not depend on the replaced contribution", name)
		}
	}
	// Both registered: ambiguous, whatever the order, until the consumer
	// replaces explicitly.
	both := builtin.with(replacement).withCodecs(gateEmptyCodec(t))
	reversed := both
	reversed.contributions = slices.Clone(both.contributions)
	slices.Reverse(reversed.contributions)
	for name, bundle := range map[string]gateBundle[gs.Language]{"the builtin registered first": both, "the fixture registered first": reversed} {
		generation, err := gatePlanGo(t, input, bundle, gateRequest{})
		if err == nil || generation.Representation().Valid() {
			t.Errorf("%s: two contributions for one role are planned without a replacement (error: %v)", name, err)
			continue
		}
		if refusal := gateExpectRefusal(t, name, err, "ambiguous contribution"); refusal.Owner != g.KernelOwner || refusal.Reason != g.ReasonSeveralContributions+" role:"+source.OptionalFieldRole {
			t.Errorf("%s: the refusal is %+v", name, refusal)
		}
		composed, err := gateComposed(t, input, bundle, func(registry *c.Registry) error {
			return registry.ReplaceContribution(replaced, replacement.Reference.Identity)
		})
		if err != nil {
			t.Errorf("%s: the explicit replacement does not plan: %v", name, err)
			continue
		}
		if written := gateMust[[]byte](t)(composed.Representation().Serialize()); !bytes.Equal(written, gateMust[[]byte](t)(after.Serialize())) {
			t.Errorf("%s: the plan after the explicit replacement is not the plan of the composition without the replaced contribution", name)
		}
	}
	// What a replacement may name.
	list := builtin.find(t, "horos-to-golang/list").Reference.Identity
	for _, named := range []struct {
		name           string
		expected, with g.ContributionIdentity
		want           error
	}{
		{"a contribution by itself", replaced, replaced, g.ErrScope},
		{"a contribution that is not registered", replaced, g.ContributionIdentity{Publisher: gatePublisher, Key: "not-registered"}, g.ErrScope},
		{"a contribution of another domain", list, replacement.Reference.Identity, g.ErrScope},
	} {
		if _, err := gateComposed(t, input, both, func(registry *c.Registry) error { return registry.ReplaceContribution(named.expected, named.with) }); !errors.Is(err, named.want) {
			t.Errorf("a replacement of %s: got %v, want %v", named.name, err, named.want)
		}
	}
	if _, err := gateComposed(t, input, both, func(registry *c.Registry) error {
		if err := registry.ReplaceContribution(replaced, replacement.Reference.Identity); err != nil {
			t.Fatal(err)
		}
		return registry.ReplaceContribution(replaced, replacement.Reference.Identity)
	}); !errors.Is(err, g.ErrAmbiguous) {
		t.Errorf("a second replacement of one contribution: got %v, want an ambiguity", err)
	}
}

// gateExtensionItem is one item of the extension domain that is out or not
// built, with the refusal that is recorded for it.
type gateExtensionItem struct {
	item string
	// scope: the kernel refuses the request as out of scope. It is no refusal
	// value: the contract has no way to state the request.
	scope string
	// class, owner and reason are those of a refusal.
	class, owner, reason string
	drive                func() (g.GenerationPlan, error)
}

// Row 33 (O17, O34), gate line 7: "one recorded refusal for each item that is
// out or not yet built".
//
// Shows, on the independent fixture source, the refusal for each item of the
// table of the extension domain that has one. A further family that is not
// built: a missing contribution. A further product whose facet a selected
// contribution lacks: a missing facet. A contextual representation at a field
// position whose fragment the target's member vocabulary does not admit:
// unsupported, with that vocabulary as owner. Cross-field layout coordination,
// which is out: a contribution at a field position that requests the operand
// of a sibling field, which the record's rule handed to it, or that requests a
// field position itself, is refused as out of scope, so the contract gives it
// no way to coordinate fields. Application-specific specialization:
// unsupported, with the application named. A shared helper request:
// unsupported. Each refusal matches exactly one value, and no plan results.
// Does not show: the fixture for each item that is in, added with the shared
// packages unchanged. That is the other half of the test of row 33 and is not
// here. The authoring route, for which nothing is refused.
// Red under: L7-field-position-requests-a-sibling,
// L7-field-position-requests-a-field-position, R10-missing-as-unsupported,
// G1-missing-facet-as-unsupported, R10-member-fragment-by-design,
// R10-specialization-by-design, R10-shared-helper-by-design.
func TestGateExtensionDomainRefusals(t *testing.T) {
	fixture := gateFixture(t)
	plan := func(change func(*f.Composition), root string) func() (g.GenerationPlan, error) {
		return func() (g.GenerationPlan, error) {
			composition := indComposition(t, fixture, indOptions{})
			if change != nil {
				change(&composition)
			}
			return fixture.PlanContributions(context.Background(), composition, root)
		}
	}
	// A record's rule sees every operand. It hands them to the contribution at
	// the field position through state of the test, as a contribution could
	// through state of its own package.
	coordinating := func(request func(operands []g.Operand[f.SourceLanguage], roles []g.Role[f.SourceLanguage], own indNameStep) indRequest) func(*f.Composition) {
		return func(composition *f.Composition) {
			var operands []g.Operand[f.SourceLanguage]
			var roles []g.Role[f.SourceLanguage]
			record := &composition.Contributions[indFind(t, *composition, "independent-outcome")]
			name := record.Rule.Name
			record.Rule.Name = func(ctx context.Context, step indNameStep) ([]indRequest, error) {
				reader, err := step.Reader()
				if err != nil {
					return nil, err
				}
				if roles, err = reader.Roles(step.Subject.Use); err != nil {
					return nil, err
				}
				operands = step.Subject.Use.Operands
				return name(ctx, step)
			}
			composition.Contributions[indFind(t, *composition, "mock-optional-nullable")].Rule.Name = func(_ context.Context, step indNameStep) ([]indRequest, error) {
				if len(operands) != 2 || len(roles) != 2 {
					return nil, errors.New("the record handed over no operands")
				}
				return []indRequest{g.Reference(step.Subject.Field.Operand), request(operands, roles, step)}, nil
			}
		}
	}
	sibling := func(operands []g.Operand[f.SourceLanguage], _ []g.Role[f.SourceLanguage], own indNameStep) indRequest {
		for _, operand := range operands {
			if operand.Use != own.Subject.Field.Operand {
				return g.Reference(operand.Use)
			}
		}
		return g.Reference(own.Subject.Field.Operand)
	}
	position := func(operands []g.Operand[f.SourceLanguage], roles []g.Role[f.SourceLanguage], own indNameStep) indRequest {
		for i, operand := range operands {
			if operand.Use != own.Subject.Field.Operand {
				return g.Field(operand.Use, roles[i], "Sibling")
			}
		}
		return g.Field(own.Subject.Field.Operand, own.Subject.Field.Role, "Own")
	}
	kind := gateMust[g.FacetKind](t)(g.NewFacetKind(g.FacetKindDescription{Identity: gatePublisher + "/facet/extension", Version: "1"}))
	const notChild = "a request must name a child the kernel minted for the subject"
	items := []gateExtensionItem{
		{item: "a further family that is not built", class: "missing contribution", owner: g.KernelOwner, reason: g.ReasonNoContribution + " form:expression/unbuilt-family", drive: plan(nil, "Unbuilt")},
		{item: "a further product whose facet a selected contribution lacks", class: "missing facet", owner: g.KernelOwner, reason: g.ReasonFacet + ": " + kind.Reference().Identity,
			drive: plan(func(composition *f.Composition) {
				product := g.Product{Identity: gatePublisher + "/extension", Requires: map[g.SubjectKind][]g.FacetKind{g.UseSubject: {kind}}}
				composition.Products, composition.Requested = []g.Product{product}, []string{product.Identity}
			}, "Word")},
		{item: "a contextual representation at a field position whose fragment is no member", class: "unsupported", owner: "member vocabulary of " + goTarget, reason: g.ReasonMemberFragment,
			drive: plan(func(composition *f.Composition) {
				composition.Contributions[indFind(t, *composition, "mock-optional-nullable")].Project = func(_ context.Context, p goProjection) (goProjected, error) {
					return goProjected{Fragment: []g.Expression[gs.Language]{gomodel.Collection(p.Dependency(0).Type)}, Members: []g.Member{{Name: p.Stated()}}}, nil
				}
			}, "Outcome")},
		{item: "cross-field layout: the operand of a sibling field", scope: notChild, drive: plan(coordinating(sibling), "Outcome")},
		{item: "cross-field layout: a field position requested from a field position", scope: notChild, drive: plan(coordinating(position), "Outcome")},
		{item: "application-specific specialization", class: "unsupported", owner: g.KernelOwner, reason: g.ReasonSpecialization,
			drive: plan(func(composition *f.Composition) {
				specialization := indEmpty("extension-specialization", g.ForSpecialization(indSelector(t, fixture, "declaration/Page"), indSelector(t, fixture, "declaration/Person")), nil,
					func(context.Context, goProjection) (goProjected, error) {
						return goProjected{Type: gomodel.Primitive(gomodel.PrimitiveString)}, nil
					})
				composition.Contributions = append(composition.Contributions, specialization)
			}, "People")},
		{item: "a shared helper request", class: "unsupported", owner: g.KernelOwner, reason: g.ReasonSharedHelper,
			drive: plan(func(composition *f.Composition) {
				composition.Contributions[indFind(t, *composition, "independent-word")].Rule.Name = func(context.Context, indNameStep) ([]indRequest, error) {
					return []indRequest{g.SharedHelper[f.SourceLanguage]("example.org/shared/helper")}, nil
				}
			}, "Word")},
	}
	// The control: the mock record plans as it is, so each refusal below is of
	// the change and not of the composition.
	if generation, err := plan(nil, "Outcome")(); err != nil || !generation.Representation().Valid() {
		t.Fatalf("the mock record does not plan: %v", err)
	}
	for _, item := range items {
		generation, err := item.drive()
		if err == nil || generation.Representation().Valid() {
			t.Errorf("%s is planned (error: %v)", item.item, err)
			continue
		}
		if item.scope != "" {
			if !errors.Is(err, g.ErrScope) || !strings.Contains(err.Error(), item.scope) || len(gateRefusalClasses(err)) != 0 {
				t.Errorf("%s: got %v, want the scope error %q", item.item, err, item.scope)
			}
			continue
		}
		refusal := gateExpectRefusal(t, item.item, err, item.class)
		if refusal.Owner != item.owner || !strings.HasPrefix(refusal.Reason, item.reason) {
			t.Errorf("%s: the refusal has the owner %q and the reason %q, want %q and %q", item.item, refusal.Owner, refusal.Reason, item.owner, item.reason)
		}
		if item.class == "unsupported" && item.reason == g.ReasonSpecialization && refusal.Application != "declaration/Page<declaration/Person>" {
			t.Errorf("%s: the refusal names the application %q", item.item, refusal.Application)
		}
	}
}
