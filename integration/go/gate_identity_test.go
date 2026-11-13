package integration_test

// Identities and nodes under the contribution contract (bittype#196,
// docs/contribution-contract.md, "Identities" and "Optional fields"): gate
// line 2, rows 5, 12, 14, 15, 18, 19, 20, 21, 25 and 32, with O5 and O31.

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	gomodel "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
)

// gateNodesWith counts the nodes of a plan that have the identity id.
func gateNodesWith(plan g.RepresentationPlan, id string) int {
	count := 0
	for _, node := range plan.Nodes() {
		if node.ID == id {
			count++
		}
	}
	return count
}

// Rows 5 and 19 (O4, O13, O20), gate line 2.
//
// Shows: two contributions that differ only in their representation revision
// give byte-identical payloads and different node identities, and so do two
// that differ only in the revision of the family's meaning. An invariant of a
// revision is carried by the reference and needs no field. The node of the
// list that depends on the text is another node too, with an equal payload,
// and so is the declaration above it: a link and an answer are bound to the
// dependency's reference. Two different contributions that produce the same
// native type do not share a node: the generated text of the plans is equal.
// The answers that a rule received are part of its node's identity: with an
// equal payload and equal links, the answer "unknown" and the answer
// "exhausted" give two nodes.
// Does not show: what a revision means, which is its publisher's. That the
// identity of a node covers the target configuration (finding F-B3,
// TestGateNodeIdentityCoversTheTargetConfiguration).
// Red under: R5-identity-ignores-representation-revision,
// R5-identity-ignores-meaning-revision, R19-identity-ignores-contribution,
// R19-identity-ignores-answers, R19-identity-ignores-links.
func TestGateRevisionIsIdentityAndNoPayload(t *testing.T) {
	from, input := gateSource(t, "module example.revision;\n\ntype Names = list<text>;\n")
	type planned struct {
		declaration, list, text g.PlanNode
		native                  string
	}
	plan := func(representation, meaning string) planned {
		bundle := gateGolang(t, from).replace(gateEvidenceList(t, from, gateEvidenceListOptions{})).replace(gateNilFreeText(t, from, representation, meaning)).withCodecs(gateDirectCodec(t), gateEmptyCodec(t))
		generation, err := gatePlanGo(t, input, bundle, gateRequest{})
		if err != nil {
			t.Fatal(err)
		}
		p := generation.Representation()
		list := gateBody(t, p, "Names")
		if len(list.Links) != 1 {
			t.Fatalf("the list has the links %+v", list.Links)
		}
		text, _ := p.Node(list.Links[0].Node)
		return planned{gateDeclaration(t, p, "Names"), list, text, gateText[gs.Language](t, golangText(), generation)}
	}
	one := plan("golang/1", "1")
	for _, c := range []struct {
		name  string
		other planned
		part  func(g.ContributionReference) string
	}{
		{"the representation revision", plan("golang/2", "1"), func(r g.ContributionReference) string { return r.Representation }},
		{"the revision of the family's meaning", plan("golang/1", "2"), func(r g.ContributionReference) string { return r.Meaning.Revision }},
	} {
		two := c.other
		if one.text.Contribution == nil || two.text.Contribution == nil || c.part(*one.text.Contribution) == c.part(*two.text.Contribution) {
			t.Fatalf("%s: the two plans do not name two revisions: %v and %v", c.name, one.text.Contribution, two.text.Contribution)
		}
		if one.text.Payload != two.text.Payload || one.list.Payload != two.list.Payload || one.declaration.Payload != two.declaration.Payload || one.text.Payload == "" {
			t.Errorf("%s: the payloads differ, or the text has none: %q and %q", c.name, one.text.Payload, two.text.Payload)
		}
		if one.text.ID == two.text.ID {
			t.Errorf("%s: the text is one node %s under two revisions", c.name, one.text.ID)
		}
		if one.list.ID == two.list.ID || one.declaration.ID == two.declaration.ID {
			t.Errorf("%s: a node that depends on the text is the same node under two revisions of the text", c.name)
		}
		if one.native != two.native {
			t.Errorf("%s: the two plans generate different native text", c.name)
		}
	}
	// Two different contributions never share a node, whatever they produce.
	builtin, err := gatePlanGo(t, input, gateGolang(t, from).replace(gateEvidenceList(t, from, gateEvidenceListOptions{})).withCodecs(gateDirectCodec(t)), gateRequest{})
	if err != nil {
		t.Fatal(err)
	}
	builtinList := gateBody(t, builtin.Representation(), "Names")
	builtinText, _ := builtin.Representation().Node(builtinList.Links[0].Node)
	if builtinText.ID == one.text.ID || builtinText.Contribution == nil || *builtinText.Contribution == *one.text.Contribution {
		t.Errorf("the builtin text contribution and the fixture's share the node %s", builtinText.ID)
	}
	// The answers are part of the identity. On the independent source the probe
	// of Word is answered "unknown" where the word states nothing and
	// "exhausted" where it states that: the payload and the linked node are the
	// same, and the node of the probe is another.
	probe := func(facts []g.Fact) g.PlanNode {
		generation, err := indPlan(t, indOptions{wordFacts: facts}, false, "ProbeWord")
		if err != nil {
			t.Fatal(err)
		}
		node, _ := gateEvidenceLink(t, generation.Representation(), "ProbeWord")
		return node
	}
	unknown, exhausted := probe([]g.Fact{}), probe([]g.Fact{gateNilFree.Exhausted()})
	if unknown.Payload != exhausted.Payload || unknown.Links[0].Node != exhausted.Links[0].Node || unknown.Links[0].Answers[0].Outcome == exhausted.Links[0].Answers[0].Outcome {
		t.Fatalf("the two probes do not differ in their answer alone: %+v and %+v", unknown.Links, exhausted.Links)
	}
	if unknown.ID == exhausted.ID {
		t.Errorf("two selections with different answers are one node %s", unknown.ID)
	}
}

// Rows 19 and 35, gate line 2: the identity of a node is over "everything
// that makes two selections represent differently", and the page lists the
// target configuration among what a shorter key would miss.
//
// Shows: one request planned under two target configurations gives two nodes
// for a selection.
// Does not show: which settings a configuration has.
// Red under: no mutant. The test is skipped as PENDING: it is red against the
// kernel as built (finding F-B3), so no mutant can show it red.
func TestGateNodeIdentityCoversTheTargetConfiguration(t *testing.T) {
	t.Skip("PENDING F-B3: the plan records the target configuration and binds nothing to it; a node has the same identity under two target configurations")
	from, input := gateSource(t, "module example.configuration;\n\ntype Names = list<text>;\n")
	node := func(value string) g.PlanNode {
		generation, err := gatePlanGo(t, input, gateGolang(t, from), gateRequest{configuration: []g.Setting{{Name: "example.test/gate/setting", Value: value}}})
		if err != nil {
			t.Fatal(err)
		}
		if got := generation.Representation().Target().Configuration; len(got) != 1 || got[0].Value != value {
			t.Fatalf("the plan records the configuration %+v", got)
		}
		return gateBody(t, generation.Representation(), "Names")
	}
	if one, two := node("1"), node("2"); one.ID == two.ID {
		t.Errorf("the list is one node %s under two target configurations", one.ID)
	}
}

// Row 18 (and row 19), gate line 2.
//
// Shows, on the independent fixture source: two uses of Page<Person> in two
// declarations are two use entries, with two identities, two sites and two
// scopes, on one representation node, which the plan holds once. The node
// holds no site. The open declaration Page is one node with one binder, the
// application links that binder to what was selected for Person, and inside
// the open declaration the element is the binder and no argument.
// Does not show: the procedure by which #195 finds that two uses share a
// node. Here a node is identified by a digest of its content.
// Red under: R18-use-entry-not-recorded, R18-node-per-use,
// R18-node-held-per-use.
func TestGateTwoUsesShareOneNode(t *testing.T) {
	generation, err := indPlan(t, indOptions{}, false, "People", "PeopleAgain")
	if err != nil {
		t.Fatal(err)
	}
	plan := generation.Representation()
	people, again := gateBody(t, plan, "People"), gateBody(t, plan, "PeopleAgain")
	if people.ID != again.ID {
		t.Fatalf("the two uses of Page<Person> are the nodes %s and %s", people.ID, again.ID)
	}
	if held := gateNodesWith(plan, people.ID); held != 1 {
		t.Errorf("the plan holds the node of Page<Person> %d times", held)
	}
	entries := gateUsesOf(plan, people.ID)
	if len(entries) != 2 || entries[0].ID == entries[1].ID || entries[0].Site == entries[1].Site || entries[0].Scope == entries[1].Scope {
		t.Fatalf("the node of Page<Person> has the use entries %+v, want two with two identities, sites and scopes", entries)
	}
	scopes := []string{entries[0].Scope, entries[1].Scope}
	slices.Sort(scopes)
	if want := []string{gateDeclaration(t, plan, "People").Symbol, gateDeclaration(t, plan, "PeopleAgain").Symbol}; !slices.Equal(scopes, want) {
		t.Errorf("the two use entries are in the scopes %v, want %v", scopes, want)
	}
	page := gateDeclaration(t, plan, "Page")
	person, _ := gateSymbol(plan, "Person")
	if people.Kind != "application" || people.Application == nil || people.Application.Declaration != page.ID || page.Kind != "declaration" || page.Binders != 1 || len(people.Application.Target) != 1 || people.Application.Target[0].Binder != 0 {
		t.Fatalf("Page<Person> is %+v over %+v", people.Application, page)
	}
	if argument, _ := plan.Node(people.Application.Target[0].Node); argument.Kind != "application" || argument.Symbol != person.ID {
		t.Errorf("the argument of Page<Person> is %+v, want the application of Person", argument)
	}
	body := gateBody(t, plan, "Page")
	if len(body.Links) != 1 {
		t.Fatalf("the body of Page has the links %+v", body.Links)
	}
	if element, _ := plan.Node(body.Links[0].Node); element.Kind != "binder" || element.Binder == nil || element.Binder.Scope != page.Symbol || element.Binder.Index != 0 {
		t.Errorf("inside Page<T> the element is %+v, want the binder of Page", element)
	}
}

// Rows 19, 20 and 32 (O5, O31), gate line 2.
//
// Shows, on the independent fixture source: Ignore<T> omits its parameter in
// the target. Asked for bindings, the application Ignore<Person> has no
// target argument and one adapter argument, tied to binder 0 and to what was
// selected for Person, and Ignore<Person> and Ignore<Page<Person>>, whose
// native types are equal, are two nodes. Asked for types alone, the omitted
// argument is not planned: no declaration of Person results, the application
// has no adapter argument, and the two applications are one node. The
// generated text is the same for both requests. An exact contribution that
// names its argument for conversion only is resolved when bindings are asked
// for and left out for types alone, and the need is recorded on the link
// either way.
// Does not show: an adapter. No facet is involved: the plan records which
// node an adapter would be built from.
// Red under: R20-omitted-argument-dropped, R20-types-plan-omitted-argument,
// R19-identity-ignores-adapters, R20-conversion-always-resolved,
// R20-conversion-not-recorded.
func TestGateOmittedArgumentIsKeptForConversion(t *testing.T) {
	bindings, err := indPlan(t, indOptions{}, true, "Kept", "KeptPage")
	if err != nil {
		t.Fatal(err)
	}
	plan := bindings.Representation()
	kept, keptPage := gateBody(t, plan, "Kept"), gateBody(t, plan, "KeptPage")
	person, planned := gateSymbol(plan, "Person")
	if !planned {
		t.Fatal("asked for bindings, the omitted argument Person is not planned")
	}
	if kept.Application == nil || len(kept.Application.Target) != 0 {
		t.Errorf("Ignore<Person> has the target arguments %+v", kept.Application)
	} else if adapters := kept.Application.Adapters; len(adapters) != 1 || adapters[0].Binder != 0 {
		t.Errorf("Ignore<Person> has the adapter arguments %+v, want one for binder 0", adapters)
	} else if argument, _ := plan.Node(adapters[0].Node); argument.Kind != "application" || argument.Symbol != person.ID {
		t.Errorf("the adapter argument of Ignore<Person> is %+v, want what was selected for Person", argument)
	}
	if kept.ID == keptPage.ID {
		t.Error("asked for bindings, Ignore<Person> and Ignore<Page<Person>> are one node")
	}
	if got := plan.Products(); !slices.Equal(got, []string{g.TypesProduct, gateBindings.Identity}) {
		t.Errorf("the plan names the products %v", got)
	}
	// Types alone.
	types, err := indPlan(t, indOptions{}, false, "Kept", "KeptPage")
	if err != nil {
		t.Fatal(err)
	}
	plan = types.Representation()
	kept, keptPage = gateBody(t, plan, "Kept"), gateBody(t, plan, "KeptPage")
	if _, planned := gateSymbol(plan, "Person"); planned {
		t.Error("asked for types alone, the omitted argument Person is planned")
	}
	if kept.Application == nil || len(kept.Application.Adapters) != 0 || len(kept.Application.Target) != 0 {
		t.Errorf("asked for types alone, Ignore<Person> has the arguments %+v", kept.Application)
	}
	if kept.ID != keptPage.ID {
		t.Error("asked for types alone, two applications that differ in an argument that is not planned are two nodes")
	}
	native := gateParseGo(t, gateText[gs.Language](t, golangText(), types))
	if native.aliases["Kept"] != "Ignore" || native.aliases["KeptPage"] != "Ignore" {
		t.Errorf("the two applications do not have equal native types: %v", native.aliases)
	}
	if withBindings := gateText[gs.Language](t, golangText(), bindings); !strings.Contains(withBindings, "type Kept = Ignore") || !strings.Contains(withBindings, "type KeptPage = Ignore") {
		t.Errorf("asked for bindings, the two applications do not have equal native types:\n%s", withBindings)
	}
	// An exact contribution that names its argument for conversion only (O31).
	for _, c := range []struct {
		bindings bool
		need     string
	}{{true, "conversion"}, {false, "conversion, not demanded"}} {
		generation, err := indPlan(t, indOptions{cellConversionOnly: true}, c.bindings, "CellOfPerson")
		if err != nil {
			t.Fatal(err)
		}
		plan := generation.Representation()
		cell := gateBody(t, plan, "CellOfPerson")
		if len(cell.Links) != 1 {
			t.Fatalf("bindings %v: the use of Cell has the links %+v", c.bindings, cell.Links)
		}
		link := cell.Links[0]
		_, planned := gateSymbol(plan, "Person")
		if planned != c.bindings || (link.Node != "") != c.bindings || link.Need != c.need || link.Kind != "child" {
			t.Errorf("bindings %v: the argument that only conversion needs is linked as %+v, and Person is planned: %v", c.bindings, link, planned)
		}
		if got := gateParseGo(t, gateText[gs.Language](t, golangText(), generation)).aliases["CellOfPerson"]; got != "string" {
			t.Errorf("bindings %v: CellOfPerson is generated as %q", c.bindings, got)
		}
	}
}

// gateSmuggling makes the Go record contribution keep the operand and the
// role selector of the first optional field it meets and present them later,
// in one of four ways: "operand" presents the kept operand with its role at
// the next optional field, "role" presents the next field's own operand with
// the kept role, "forged" presents each optional field's own operand with a
// role that another issuer than the source issued for the same pair, and
// "none" changes nothing.
func gateSmuggling(base goContribution, mode string) goContribution {
	base.Reference.Identity = g.ContributionIdentity{Publisher: gatePublisher, Key: "smuggling-record"}
	name := base.Rule.Name
	forger := g.NewRoleIssuer[h.Language]()
	var kept *source.Field
	base.Rule.Name = func(ctx context.Context, step horosNameStep) ([]horosRequest, error) {
		requests, err := name(ctx, step)
		if err != nil {
			return nil, err
		}
		reader, err := step.Reader()
		if err != nil {
			return nil, err
		}
		fields, err := reader.Fields(step.Subject.Use)
		if err != nil || len(fields) != len(requests) {
			return nil, errors.Join(errors.New("the record's requests are not one for each field"), err)
		}
		for i, field := range fields {
			switch {
			case !field.Role.Valid() || mode == "none":
			case mode == "forged":
				requests[i] = g.Field(field.Use, forger.Issue(source.OptionalFieldRole, step.Subject.Use.Use, field.Use), "Smuggled")
			case kept == nil:
				kept = &field
			case mode == "operand":
				requests[i] = g.Field(kept.Use, kept.Role, "Smuggled")
			case mode == "role":
				requests[i] = g.Field(field.Use, kept.Role, "Smuggled")
			}
		}
		return requests, nil
	}
	return base
}

// Row 12, gate line 2.
//
// Shows: the kernel keeps its scope check at a field position. A request with
// the operand of another record is refused, because the operand is no child
// of the asking record's use. A request with the record's own operand and a
// role that the source issued for another record's field, or for another
// field of the same record, is refused, because the role selector was not
// issued for this parent and operand. A role that another issuer than the
// source issued for the right pair under the same key is refused too. Each is
// a scope error, matches no refusal value, and no plan results. The same
// contribution without the smuggling plans.
// Does not show: that the source issues a role only for an optional field,
// which is the source integration's own fact. The comparison of the role's
// parent alone: an operand has one parent, so a role whose operand is the
// request's is a role of its parent for every role that the source issues.
// The mutant R12-role-ignores-parent therefore survives.
// Red under: R12-field-ignores-child, R12-role-ignores-operand,
// R12-role-ignores-issuer. Not red under R12-role-ignores-parent, for the
// reason above.
func TestGateFieldPositionKeepsItsScope(t *testing.T) {
	const two, siblings = "module example.scope;\n\ntype A = { p?: text; };\ntype B = { q?: text; };\n", "module example.scope;\n\ntype C = { p?: text; q?: text; };\n"
	const notChild, notIssued = "a request must name a child the kernel minted for the subject", "the role selector was not issued for this parent and operand"
	plan := func(module, mode string) (g.GenerationPlan, error) {
		from, input := gateSource(t, module)
		b := gateGolang(t, from)
		return gatePlanGo(t, input, b.replace(gateSmuggling(b.find(t, "horos-to-golang/record"), mode)), gateRequest{})
	}
	for _, module := range []string{two, siblings} {
		if generation, err := plan(module, "none"); err != nil || !generation.Representation().Valid() {
			t.Fatalf("the contribution without smuggling does not plan: %v", err)
		}
	}
	for _, c := range []struct {
		name, module, mode, want string
	}{
		{"the operand of another record, with its role", two, "operand", notChild},
		{"a role issued for a field of another record", two, "role", notIssued},
		{"a role issued for another field of the same record", siblings, "role", notIssued},
		{"a role of another issuer for the right pair", two, "forged", notIssued},
	} {
		generation, err := plan(c.module, c.mode)
		if err == nil || generation.Representation().Valid() {
			t.Errorf("a request at a field position with %s is planned (error: %v)", c.name, err)
			continue
		}
		if !errors.Is(err, g.ErrScope) || !strings.Contains(err.Error(), c.want) || len(gateRefusalClasses(err)) != 0 {
			t.Errorf("a request at a field position with %s: got %v, want the scope error %q", c.name, err, c.want)
		}
	}
}

// Row 14 (and row 13, O40), gate line 2: the second half of the row.
// TestGateMemberFragmentIsTheTargetsFieldForm shows the first.
//
// Shows, on the independent fixture source: a target that states no member
// vocabulary refuses a field position. The mock record, whose fields are
// answered at a field position, is refused as unsupported, with the kernel as
// owner and the contribution selected for the position named, and no plan
// results. The same target plans a declaration that has no field position,
// and the Go target, which states a member vocabulary, plans the mock record.
// Does not show: a fragment of more than one member.
// Red under: R14-field-position-without-vocabulary,
// R10-no-member-vocabulary-by-design.
func TestGateTargetWithoutMemberVocabularyRefusesAFieldPosition(t *testing.T) {
	fixture := gateFixture(t)
	bare := gateBareTarget{ggtarget.Provider{Package: "generated"}}
	if _, has := any(bare).(g.MemberVocabulary[gs.Language]); has {
		t.Fatal("the bare target states a member vocabulary")
	}
	if _, has := any(ggtarget.Provider{}).(g.MemberVocabulary[gs.Language]); !has {
		t.Fatal("the Go target states no member vocabulary")
	}
	plan := func(target g.Target[gs.Language], root string) (g.GenerationPlan, error) {
		composition := indComposition(t, fixture, indOptions{})
		composition.Target = target
		return fixture.PlanContributions(context.Background(), composition, root)
	}
	if generation, err := plan(nil, "Outcome"); err != nil || !generation.Representation().Valid() {
		t.Fatalf("the Go target does not plan the mock record: %v", err)
	}
	if generation, err := plan(bare, "People"); err != nil || !generation.Representation().Valid() {
		t.Errorf("the target without a member vocabulary does not plan a declaration without a field position: %v", err)
	}
	generation, err := plan(bare, "Outcome")
	if err == nil || generation.Representation().Valid() {
		t.Fatalf("the target without a member vocabulary plans a field position (error: %v)", err)
	}
	refusal := gateExpectRefusal(t, "a field position for a target without a member vocabulary", err, "unsupported")
	if refusal.Owner != g.KernelOwner || refusal.Reason != g.ReasonNoMemberVocabulary || refusal.Contribution != gateReference("mock-optional-nullable", "golang/1") || refusal.Subject != "role/"+f.OptionalNullableRole {
		t.Errorf("the refusal is %+v", refusal)
	}
}

// Row 15 (O4), gate line 2.
//
// Shows: the identity that a registry compares is the publisher and the
// publisher's local key. Two contributions with one local key under two
// publishers register side by side. One publisher with one key twice fails as
// an invalid description, for the same target and for a second target: the
// target is a stated field and no part of the key. A contribution registers
// for the target that it states and for no other.
// Does not show: that a publisher's name is its owner's. A name does not
// authenticate, and the composition lock is deferred (row 15a).
// Red under: R15-key-without-publisher, R15-key-twice-admitted,
// R15-key-per-target, R15-stated-target-unchecked.
func TestGatePublisherAndLocalKey(t *testing.T) {
	from, _ := gateSource(t, "module example.identity;\n\ntype Name = text;\n")
	goLanguage, tsLanguage := nativeChecker(t, golangText()).Language(), nativeChecker(t, typescriptText()).Language()
	var registry c.Registry
	text := gateGolang(t, from).find(t, "horos-to-golang/text")
	if err := c.RegisterContribution(&registry, from, goLanguage, text); err != nil {
		t.Fatal(err)
	}
	other := text
	other.Reference.Identity.Publisher = "example.org/another-publisher"
	if err := c.RegisterContribution(&registry, from, goLanguage, other); err != nil {
		t.Errorf("one local key under two publishers does not register: %v", err)
	}
	if err := c.RegisterContribution(&registry, from, goLanguage, text); !errors.Is(err, g.ErrInvalid) {
		t.Errorf("one publisher registers one key twice: %v", err)
	}
	// Another contribution of the same publisher under the same key, for a
	// second target.
	sameKey := gateTypescript(t, from).find(t, "horos-to-typescript/text")
	sameKey.Reference.Identity = text.Reference.Identity
	if err := c.RegisterContribution(&registry, from, tsLanguage, sameKey); !errors.Is(err, g.ErrInvalid) {
		t.Errorf("one publisher registers one key for two targets: %v", err)
	}
	// The control: the TypeScript contribution under its own key registers.
	if err := c.RegisterContribution(&registry, from, tsLanguage, gateTypescript(t, from).find(t, "horos-to-typescript/text")); err != nil {
		t.Errorf("the TypeScript contribution under its own key does not register: %v", err)
	}
	// A contribution is registered for the target it states.
	misstated := gateTypescript(t, from).find(t, "horos-to-typescript/list")
	misstated.Target = goTarget
	if err := c.RegisterContribution(&registry, from, tsLanguage, misstated); !errors.Is(err, g.ErrScope) {
		t.Errorf("a contribution that states another target than it is registered for: %v", err)
	}
}

// Row 21 (#166), gate line 2.
//
// Shows, on the independent fixture source: a contribution whose domain is
// one application, Page<Person>, makes that application refused as
// unsupported, with the kernel as owner, the application named by its
// declaration and its actuals, and the contribution named. Its rule does not
// run. Another application of the same declaration, Page<Cell<Person>>,
// plans, and so does the open declaration. An actual that is itself an
// application is matched by its callee: a specialization for Page<Cell>
// refuses Page<Cell<Person>> and names it so. Without the specialization
// Page<Person> plans.
// Does not show: a design of #166. A specialization is never planned.
// Red under: R21-specialization-planned-as-substitution,
// R21-application-not-named, R21-specialization-matches-any-actual,
// R10-specialization-by-design.
func TestGateSpecializationIsRefusedWithTheApplicationNamed(t *testing.T) {
	fixture := gateFixture(t)
	ran := 0
	compose := func(key string, actuals ...string) f.Composition {
		composition := indComposition(t, fixture, indOptions{})
		selectors := make([]g.Selector, len(actuals))
		for i, actual := range actuals {
			selectors[i] = indSelector(t, fixture, actual)
		}
		specialization := indEmpty(key, g.ForSpecialization(indSelector(t, fixture, "declaration/Page"), selectors...), nil,
			func(context.Context, goProjection) (goProjected, error) {
				ran++
				return goProjected{Type: gomodel.Primitive(gomodel.PrimitiveString)}, nil
			})
		specialization.Rule.Decide = func(context.Context, indDecideStep) (g.Decision, error) {
			ran++
			return g.Decision{Payload: f.Empty{}}, nil
		}
		composition.Contributions = append(composition.Contributions, specialization)
		return composition
	}
	plans := func(composition f.Composition, root string) error {
		generation, err := fixture.PlanContributions(context.Background(), composition, root)
		if err == nil && !generation.Representation().Valid() {
			t.Fatal("the planner returned no plan and no error")
		}
		return err
	}
	if err := plans(indComposition(t, fixture, indOptions{}), "People"); err != nil {
		t.Fatalf("without a specialization Page<Person> does not plan: %v", err)
	}
	forPerson := compose("specialized-page-of-person", "declaration/Person")
	for _, root := range []string{"Cells", "Page", "Outer"} {
		if err := plans(forPerson, root); err != nil {
			t.Errorf("with a specialization for Page<Person>, %s does not plan: %v", root, err)
		}
	}
	for _, c := range []struct {
		name, key, actual, root, application string
	}{
		{"Page<Person>", "specialized-page-of-person", "declaration/Person", "People", "declaration/Page<declaration/Person>"},
		{"Page<Cell<Person>>, by the callee of its actual", "specialized-page-of-cell", "declaration/Cell", "Cells", "declaration/Page<declaration/Cell>"},
	} {
		err := plans(compose(c.key, c.actual), c.root)
		if err == nil {
			t.Errorf("the specialization for %s is planned", c.name)
			continue
		}
		refusal := gateExpectRefusal(t, "a specialization for "+c.name, err, "unsupported")
		if refusal.Owner != g.KernelOwner || refusal.Reason != g.ReasonSpecialization || refusal.Application != c.application || refusal.Contribution != gateReference(c.key, "golang/1") || refusal.Subject != "declaration/Page" {
			t.Errorf("the refusal for %s is %+v, want the application %q", c.name, refusal, c.application)
		}
		if !strings.Contains(err.Error(), "application "+c.application) {
			t.Errorf("the refusal for %s does not name the application in its text: %v", c.name, err)
		}
	}
	if ran != 0 {
		t.Errorf("the rule or the projector of a specialization ran %d times", ran)
	}
}

// Row 25 (O34), gate line 2.
//
// Shows, on the independent fixture source: two contributions that each own a
// declaration and spell it alike get two logical symbols. The key of an owned
// declaration is its owner's contribution reference and its local name, and
// the spelling is not the key. One contribution that refers to its owned
// declaration at two uses gets one symbol. A request for a shared helper is
// in the contract and is refused as unsupported, with the kernel as owner.
// Does not show: what a Files stage does with two declarations of one name in
// one scope. The Go target refuses such a program, so the test plans with a
// target that states no program rule. Whose key a shared helper request uses:
// it is not specified until a contribution needs one (O34).
// Red under: R25-owned-shared-by-spelling, R25-owned-per-use,
// R10-shared-helper-by-design.
func TestGateOwnedDeclarationsAreKeyedByTheirOwner(t *testing.T) {
	fixture := gateFixture(t)
	const spelling, local = "Helper", "helper"
	composition := indComposition(t, fixture, indOptions{})
	composition.Target = ruleless{ggtarget.Provider{Package: "generated"}}
	uses := map[string]int{}
	owning := func(key string, body g.Expression[gs.Language], record func(helper g.Expression[gs.Language], p goProjection) goProjected) {
		contribution := &composition.Contributions[indFind(t, composition, key)]
		contribution.Owned = []g.OwnedDescription[gs.Language]{{Name: local, Declaration: g.TargetDeclaration[gs.Language]{Name: spelling, Form: gomodel.Form(gomodel.FormAlias), Body: body}}}
		contribution.Project = func(_ context.Context, p goProjection) (goProjected, error) {
			helper, err := p.Owned(local)
			uses[key]++
			return record(helper, p), err
		}
	}
	owning("independent-person", gomodel.Primitive(gomodel.PrimitiveBool), func(helper g.Expression[gs.Language], _ goProjection) goProjected {
		return goProjected{Type: gomodel.Record(gomodel.Field("Name", false, gomodel.Primitive(gomodel.PrimitiveString)), gomodel.Field("Mark", false, helper)), Members: []g.Member{{Name: "Name"}, {Name: "Mark"}}}
	})
	owning("independent-page", gomodel.Primitive(gomodel.PrimitiveString), func(helper g.Expression[gs.Language], p goProjection) goProjected {
		return goProjected{Type: gomodel.Record(gomodel.Field("Items", false, gomodel.Collection(p.Dependency(0).Type)), gomodel.Field("Cursor", false, helper)), Members: []g.Member{{Name: "Items"}, {Name: "Cursor"}}}
	})
	// Person and the first field of Outcome are two uses of the person form.
	generation, err := fixture.PlanContributions(context.Background(), composition, "Person", "Page", "Outcome")
	if err != nil {
		t.Fatal(err)
	}
	if uses["independent-person"] != 2 || uses["independent-page"] != 1 {
		t.Fatalf("the two contributions referred to their owned declarations %v times, want the person twice and the page once", uses)
	}
	var owned []string
	for _, symbol := range generation.Representation().Symbols() {
		if symbol.Role == "owned" {
			if symbol.Intent != spelling {
				t.Errorf("the owned declaration %s is spelled %q", symbol.ID, symbol.Intent)
			}
			owned = append(owned, symbol.ID)
		}
	}
	slices.Sort(owned)
	want := []string{"owned:" + gateReference("independent-page", "golang/1").String() + "/" + local, "owned:" + gateReference("independent-person", "golang/1").String() + "/" + local}
	if !slices.Equal(owned, want) {
		t.Errorf("the declarations spelled %s have the symbols %v, want %v", spelling, owned, want)
	}
	// A shared helper is asked for by a request, which is refused.
	asking := indComposition(t, fixture, indOptions{})
	asking.Contributions[indFind(t, asking, "independent-word")].Rule.Name = func(context.Context, indNameStep) ([]indRequest, error) {
		return []indRequest{g.SharedHelper[f.SourceLanguage]("example.org/shared/helper")}, nil
	}
	refused, err := fixture.PlanContributions(context.Background(), asking, "Word")
	if err == nil || refused.Representation().Valid() {
		t.Fatalf("a request for a shared helper is planned (error: %v)", err)
	}
	if refusal := gateExpectRefusal(t, "a request for a shared helper", err, "unsupported"); refusal.Owner != g.KernelOwner || refusal.Reason != g.ReasonSharedHelper || refusal.Contribution != gateReference("independent-word", "golang/1") {
		t.Errorf("the refusal of a shared helper request is %+v", refusal)
	}
}
