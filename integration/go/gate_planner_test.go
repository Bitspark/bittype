package integration_test

// Three tests of the planner and the registry that the research owner asked
// for after the prototype (bittype#196): a refusal or failure of a selected
// contribution is returned and not dropped, a capability that a contribution
// kept gives nothing after its step has returned, and registration refuses an
// incomplete contribution.

import (
	"context"
	"errors"
	"reflect"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	gomodel "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
)

// Rows 1 and 10 (O17), gate line 2. The fault this test names: the
// prototype's planner dropped what the rule or the projector of an exact
// contribution returned, through a shadowed error variable in the planner's
// function for a use (planner.use in contribution_plan.go). The planning then
// went on with an empty dependency.
//
// Shows: for an exact contribution (the fixture's contribution for the
// declaration Cell) and for a contribution of a form (the fixture's word),
// what the selected contribution returns reaches the caller of the planner as
// it is. A refusal in step 1, in step 3 or in the projector is that refusal:
// its class, owner and reason, with the contribution, the subject and the
// path that the planner adds, and it matches exactly one refusal value. An
// error that is no refusal, from step 3 or from the projector, is a rule
// failure of that contribution with the error as its cause, and it matches no
// refusal value. No plan results, and the step that refuses or fails runs
// once. The same composition without the refusal plans.
// Does not show: the same for a contribution at a field position or for the
// declaration policy. Their results are returned by other statements of the
// planner, which this test does not reach.
// Red under: X-exact-error-shadowed (the exact contribution),
// X-form-error-shadowed (the form).
func TestGateSelectedContributionRefusalIsNotDropped(t *testing.T) {
	const reason = "the fixture refuses this subject"
	sentinel := errors.New("the fixture fails")
	type variant struct {
		name  string
		class string // the refusal value that the error matches, or "" for a rule failure
		apply func(contribution *f.Contribution, result func() error)
	}
	refusal := func(class g.RefusalClass) func() error {
		return func() error { return g.Refuse(class, gatePublisher, reason) }
	}
	inName := func(contribution *f.Contribution, result func() error) {
		contribution.Rule.Name = func(context.Context, indNameStep) ([]indRequest, error) { return nil, result() }
	}
	inDecide := func(contribution *f.Contribution, result func() error) {
		contribution.Rule.Decide = func(context.Context, indDecideStep) (g.Decision, error) { return g.Decision{}, result() }
	}
	inProject := func(contribution *f.Contribution, result func() error) {
		contribution.Project = func(context.Context, goProjection) (goProjected, error) { return goProjected{}, result() }
	}
	variants := []struct {
		variant
		result func() error
	}{
		{variant{"a refusal by design in step 1", "refused by design", inName}, refusal(g.RefusalByDesign)},
		{variant{"an unsupported subject in step 3", "unsupported", inDecide}, refusal(g.RefusalUnsupported)},
		{variant{"an unsatisfied requirement in the projector", "requirement unsatisfied", inProject}, refusal(g.RefusalRequirementUnsatisfied)},
		{variant{"an unresolved requirement in step 3", "requirement unresolved", inDecide}, refusal(g.RefusalRequirementUnresolved)},
		{variant{"an error in step 1", "", inName}, func() error { return sentinel }},
		{variant{"an error in step 3", "", inDecide}, func() error { return sentinel }},
		{variant{"an error in the projector", "", inProject}, func() error { return sentinel }},
	}
	for _, selected := range []struct {
		name, key, selector, root string
	}{
		{"an exact contribution", "independent-cell", "declaration/Cell", "CellOfPerson"},
		{"a contribution of a form", "independent-word", "expression/word", "Word"},
	} {
		t.Run(selected.name, func(t *testing.T) {
			fixture := gateFixture(t)
			subject := indSelector(t, fixture, selected.selector).Key()
			clean := indComposition(t, fixture, indOptions{})
			if _, err := fixture.PlanContributions(context.Background(), clean, selected.root); err != nil {
				t.Fatalf("the composition does not plan %s without a refusal: %v", selected.root, err)
			}
			for _, v := range variants {
				composition := indComposition(t, fixture, indOptions{})
				contribution := &composition.Contributions[indFind(t, composition, selected.key)]
				called := 0
				v.apply(contribution, func() error {
					called++
					return v.result()
				})
				plan, err := fixture.PlanContributions(context.Background(), composition, selected.root)
				if err == nil || plan.Representation().Valid() {
					t.Errorf("%s: the planner went on and returned a plan (error: %v)", v.name, err)
					continue
				}
				if called != 1 {
					t.Errorf("%s: the step ran %d times", v.name, called)
				}
				if v.class == "" {
					var failure *g.RuleFailure
					if !errors.As(err, &failure) || failure.Err != sentinel || failure.Contribution != contribution.Reference {
						t.Errorf("%s: the planner returned %v and not the rule failure of %s with its cause", v.name, err, contribution.Reference)
					}
					if matched := gateRefusalClasses(err); len(matched) != 0 {
						t.Errorf("%s: the failure matches the refusal values %v", v.name, matched)
					}
					continue
				}
				got := gateExpectRefusal(t, v.name, err, v.class)
				if got.Owner != gatePublisher || got.Reason != reason {
					t.Errorf("%s: the refusal has the owner %q and the reason %q", v.name, got.Owner, got.Reason)
				}
				if got.Contribution != contribution.Reference || got.Subject != subject || got.Target != goTarget || got.Profile != f.Profile {
					t.Errorf("%s: the refusal names the contribution %s, the subject %q, the target %q and the profile %q", v.name, got.Contribution, got.Subject, got.Target, got.Profile)
				}
				if len(got.Path) == 0 || got.Path[len(got.Path)-1] != subject {
					t.Errorf("%s: the path of the refusal is %v and does not end at the subject", v.name, got.Path)
				}
			}
		})
	}
}

// gateKept is what a contribution of the next test kept of its capabilities,
// and what it saw of them.
type gateKept struct {
	name     indNameStep
	decide   indDecideStep
	answers  []g.Answer
	use      goProjection
	declared goProjection
	// usedRan and declaredRan say that the two projectors ran, so that the two
	// kept projections are the planner's and not zero values.
	usedRan, declaredRan bool
	// open are the faults of a capability while its own step ran, and late what
	// a kept capability still gave while a later step of the same selection ran.
	open, late []string
}

const gateKeptHelper = "kept"

func gateKeptOwned(name string) []g.OwnedDescription[gs.Language] {
	return []g.OwnedDescription[gs.Language]{{Name: gateKeptHelper, Declaration: g.TargetDeclaration[gs.Language]{Name: name, Form: gomodel.Form(gomodel.FormAlias), Body: gomodel.Primitive(gomodel.PrimitiveString)}}}
}

// closed lists what a kept step of a rule and a kept projection still give.
// It is empty when each gives nothing.
func (k *gateKept) closed(steps, projections bool) []string {
	var given []string
	isScope := func(err error) bool { return errors.Is(err, g.ErrScope) }
	if steps {
		if reader, err := k.name.Reader(); !isScope(err) || !reflect.ValueOf(&reader).Elem().IsZero() {
			given = append(given, "the kept step 1 gives a reader")
		}
		if reader, err := k.decide.Reader(); !isScope(err) || !reflect.ValueOf(&reader).Elem().IsZero() {
			given = append(given, "the kept step 3 gives a reader")
		}
		if answers := k.decide.Answers(0); len(answers) != 0 {
			given = append(given, "the kept step 3 lists its answers")
		}
		if len(k.answers) == 1 {
			if outcome, content, err := k.decide.Present(0, k.answers[0]); !isScope(err) || outcome != 0 || content != "" {
				given = append(given, "the kept step 3 presents an answer")
			}
		}
		if outcome, _, established := k.decide.Evidence(0, gateNilFree); outcome != g.Unknown || established {
			given = append(given, "the kept step 3 gives evidence")
		}
	}
	if projections {
		if _, err := k.use.Owned(gateKeptHelper); !isScope(err) {
			given = append(given, "the kept projection of a use gives an owned declaration")
		}
		if _, err := k.declared.Owned(gateKeptHelper); !isScope(err) {
			given = append(given, "the kept projection of a declaration gives an owned declaration")
		}
	}
	return given
}

// The closing of a step and of a projection (docs/contribution-contract.md,
// "What a rule does: the protocol", and rows 4 and 28), gate line 1.
//
// Shows: a rule that keeps its two steps and an answer, and two projectors,
// of a use and of a declaration, that keep their projections, get nothing
// from them afterwards. While its own function runs, each capability works:
// step 1 gives the reader, step 3 gives the reader and presents the answer
// that was issued, and each projection gives the declaration that its
// contribution owns. Once step 1 has returned it gives no reader, which the
// rule sees in step 3. Once step 3 has returned it gives no reader, lists no
// answer, presents none and gives no evidence, which the projector sees. After
// planning neither step gives anything, and neither projection gives an owned
// declaration; each call fails as a scope error. The calls after planning
// leave the plan as it was: its serialized form and the generated text are
// those of before the calls.
// Does not show: that a rule cannot keep a value that it took from a
// capability while that was open. A reader value or a use that a rule copied
// is the source's to bound ("the kernel cannot revoke a value of a type it
// does not know"). It does not show the views of a projection that are plain
// values: the payload, the dependencies, the stated name and the origins stay
// readable, and they are copies. The subject of a kept step stays readable
// too: its uses and declarations are handles, and this test does not show
// what the planner does with a handle that a rule presents in a later run.
// The zero value of a projection, which no planner made, is closed in the
// same way.
// Red under: X-name-step-stays-open, X-decide-step-stays-open,
// X-reader-ignores-closed, X-answers-ignore-closed, X-present-ignores-closed,
// X-projection-stays-open, X-declaration-projection-stays-open,
// X-owned-ignores-closed, X-zero-projection-is-open.
func TestGateKeptStepAndProjectionGiveNoCapability(t *testing.T) {
	fixture := gateFixture(t)
	composition := indComposition(t, fixture, indOptions{})
	kept := &gateKept{}
	note := func(list *[]string, what string, err error) {
		if err != nil {
			what += ": " + err.Error()
		}
		*list = append(*list, what)
	}
	probe := &composition.Contributions[indFind(t, composition, "independent-probe")]
	name, decide, project := probe.Rule.Name, probe.Rule.Decide, probe.Project
	probe.Owned = gateKeptOwned("KeptByProbe")
	probe.Rule.Name = func(ctx context.Context, step indNameStep) ([]indRequest, error) {
		kept.name = step
		if reader, err := step.Reader(); err != nil || reader == nil {
			note(&kept.open, "step 1 gives no reader while it runs", err)
		}
		return name(ctx, step)
	}
	probe.Rule.Decide = func(ctx context.Context, step indDecideStep) (g.Decision, error) {
		kept.decide, kept.answers = step, step.Answers(0)
		if reader, err := step.Reader(); err != nil || reader == nil {
			note(&kept.open, "step 3 gives no reader while it runs", err)
		}
		if len(kept.answers) != 1 {
			note(&kept.open, "step 3 lists no answer while it runs", nil)
		} else if outcome, _, err := step.Present(0, kept.answers[0]); err != nil || outcome != g.Proved {
			note(&kept.open, "step 3 does not present its answer while it runs", err)
		}
		// Step 1 has returned.
		if reader, err := kept.name.Reader(); !errors.Is(err, g.ErrScope) || reader != nil {
			note(&kept.late, "the kept step 1 gives a reader while step 3 runs", nil)
		}
		return decide(ctx, step)
	}
	probe.Project = func(ctx context.Context, p goProjection) (goProjected, error) {
		kept.use, kept.usedRan = p, true
		if _, err := p.Owned(gateKeptHelper); err != nil {
			note(&kept.open, "the projection of a use gives no owned declaration while its projector runs", err)
		}
		// The rule has returned.
		kept.late = append(kept.late, kept.closed(true, false)...)
		return project(ctx, p)
	}
	policy := &composition.Contributions[indFind(t, composition, "independent-declarations")]
	declare := policy.Declare
	policy.Owned = gateKeptOwned("KeptByPolicy")
	policy.Declare = func(ctx context.Context, p goProjection) (g.DeclarationDescription[gs.Language], error) {
		kept.declared, kept.declaredRan = p, true
		if _, err := p.Owned(gateKeptHelper); err != nil {
			note(&kept.open, "the projection of a declaration gives no owned declaration while its projector runs", err)
		}
		return declare(ctx, p)
	}
	plan, err := fixture.PlanContributions(context.Background(), composition, "ProbeWord")
	if err != nil {
		t.Fatal(err)
	}
	if kept.name.Subject.Kind != g.UseSubject || kept.decide.Subject.Kind != g.UseSubject || len(kept.answers) != 1 {
		t.Fatalf("the fixture kept no step of the probe's rule: %+v", kept)
	}
	if !kept.usedRan || !kept.declaredRan {
		t.Fatalf("a projector of the fixture did not run (of a use: %v, of a declaration: %v)", kept.usedRan, kept.declaredRan)
	}
	// A projection that no planner made is closed too: it gives nothing, and
	// reading it fails nothing.
	var none goProjection
	if _, err := none.Owned(gateKeptHelper); !errors.Is(err, g.ErrScope) {
		t.Errorf("the zero projection gives an owned declaration: %v", err)
	}
	if none.Payload() != nil || none.Dependencies() != 0 || none.Binders() != 0 || none.Stated() != "" || none.Dependency(0).Stated != "" {
		t.Errorf("the zero projection has content")
	}
	for _, fault := range kept.open {
		t.Errorf("a capability does not work during its own step: %s", fault)
	}
	for _, given := range kept.late {
		t.Errorf("during a later step of the same selection: %s", given)
	}
	before, err := plan.Representation().Serialize()
	if err != nil {
		t.Fatal(err)
	}
	for _, given := range kept.closed(true, true) {
		t.Errorf("after planning: %s", given)
	}
	after, err := plan.Representation().Serialize()
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Errorf("the calls on the kept capabilities changed the serialized plan")
	}
	// The plan is the one that was made: it runs, and the probe took the direct
	// form on the evidence that its step presented.
	if native := gateParseGo(t, gateText[gs.Language](t, golangText(), plan)); native.aliases["ProbeWord"] != "[]Word" {
		t.Errorf("the plan generates ProbeWord as %q", native.aliases["ProbeWord"])
	}
}

// Rows 1, 15 and 17 (O4, O20), gate line 2: Contribution.Validate at
// registration.
//
// Shows: the registry refuses a contribution that lacks a part: either part
// of its identity, either part of its meaning, its representation revision,
// its profile, its target, either part of its schema reference, the deciding
// step of its rule, its domain, the projector of its domain (none, both, the
// one of a declaration for a form, the one of a use for the declaration
// policy), and a name of an owned declaration, or that owns two declarations
// under one name. Each is refused as an invalid description, by Validate and
// by the registry. The one exception is the target: Validate calls a missing
// target invalid, and the registry, which compares the stated target with the
// one it registers for before it validates, calls it a scope error. The
// complete contribution registers, and so does the complete declaration
// policy. A refused contribution leaves the registry as it was: the complete
// one registers afterwards under the same identity.
// Does not show: that a complete contribution is right. Validate checks that
// the parts are stated and not what they say. It does not show the planner's
// own validation of a contribution that reaches it without a registry.
// Red under: X-registry-skips-validate, X-validate-ignores-target,
// X-validate-ignores-schema, X-validate-ignores-decide,
// X-validate-ignores-domain, X-validate-weak-projector,
// X-validate-ignores-repeated-owned, X-validate-ignores-unnamed-owned,
// X-reference-valid-without-revision, B-reference-valid-without-meaning.
func TestGateRegistrationRefusesAnIncompleteContribution(t *testing.T) {
	from, _ := gateSource(t, gatePagination)
	to := nativeChecker(t, golangText()).Language()
	complete := gateNilFreeText(t, from, "golang/1", "1")
	complete.Owned = gateKeptOwned("Helper")
	policy := complete
	policy.Reference = gateReference("declarations", "golang/1")
	policy.Domain, policy.Project = g.ForDeclarationPolicy(), nil
	policy.Declare = func(context.Context, goProjection) (g.DeclarationDescription[gs.Language], error) {
		return g.DeclarationDescription[gs.Language]{}, nil
	}
	register := func(contribution goContribution) (*c.Registry, error) {
		registry := &c.Registry{}
		return registry, c.RegisterContribution[h.Language, gs.Language, source.Reader](registry, from, to, contribution)
	}
	for name, contribution := range map[string]goContribution{"the complete contribution": complete, "the complete declaration policy": policy} {
		if err := contribution.Validate(); err != nil {
			t.Fatalf("%s is not valid: %v", name, err)
		}
		if _, err := register(contribution); err != nil {
			t.Fatalf("%s does not register: %v", name, err)
		}
	}
	declare := policy.Declare
	cases := []struct {
		name string
		base goContribution
		// registry is the value that the registry refuses with, where it is not
		// the invalid description that Validate reports.
		registry error
		lacks    func(*goContribution)
	}{
		{"no publisher", complete, nil, func(k *goContribution) { k.Reference.Identity.Publisher = "" }},
		{"no local key", complete, nil, func(k *goContribution) { k.Reference.Identity.Key = "" }},
		{"no meaning identity", complete, nil, func(k *goContribution) { k.Reference.Meaning.Identity = "" }},
		{"no meaning revision", complete, nil, func(k *goContribution) { k.Reference.Meaning.Revision = "" }},
		{"no representation revision", complete, nil, func(k *goContribution) { k.Reference.Representation = "" }},
		{"no profile", complete, nil, func(k *goContribution) { k.Profile = "" }},
		{"no target", complete, g.ErrScope, func(k *goContribution) { k.Target = "" }},
		{"no schema identity", complete, nil, func(k *goContribution) { k.Schema.Identity = "" }},
		{"no schema version", complete, nil, func(k *goContribution) { k.Schema.Version = "" }},
		{"no deciding step", complete, nil, func(k *goContribution) { k.Rule.Decide = nil }},
		{"no domain", complete, nil, func(k *goContribution) { k.Domain = g.Domain{} }},
		{"no projector", complete, nil, func(k *goContribution) { k.Project = nil }},
		{"both projectors", complete, nil, func(k *goContribution) { k.Declare = declare }},
		{"the projector of a declaration for a form", complete, nil, func(k *goContribution) { k.Project, k.Declare = nil, declare }},
		{"the projector of a use for the declaration policy", policy, nil, func(k *goContribution) { k.Project, k.Declare = complete.Project, nil }},
		{"no projector for the declaration policy", policy, nil, func(k *goContribution) { k.Declare = nil }},
		{"an owned declaration without a name", complete, nil, func(k *goContribution) { k.Owned = gateKeptOwned("Helper"); k.Owned[0].Name = "" }},
		{"two owned declarations under one name", complete, nil, func(k *goContribution) { k.Owned = append(gateKeptOwned("Helper"), gateKeptOwned("Other")...) }},
	}
	for _, incomplete := range cases {
		contribution := incomplete.base
		incomplete.lacks(&contribution)
		if err := contribution.Validate(); !errors.Is(err, g.ErrInvalid) {
			t.Errorf("%s: Validate does not call the contribution invalid: %v", incomplete.name, err)
		}
		want := g.ErrInvalid
		if incomplete.registry != nil {
			want = incomplete.registry
		}
		registry, err := register(contribution)
		if !errors.Is(err, want) {
			t.Errorf("%s: the registry does not refuse the contribution as %q: %v", incomplete.name, want, err)
		}
		// The refused contribution is not in the registry: the complete one
		// registers under the same identity. An identity that was registered
		// would fail as registered twice.
		if err := c.RegisterContribution[h.Language, gs.Language, source.Reader](registry, from, to, incomplete.base); err != nil {
			t.Errorf("%s: after the refusal the complete contribution does not register: %v", incomplete.name, err)
		}
	}
}
