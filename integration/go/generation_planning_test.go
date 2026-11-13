package integration_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	checkts "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/engine/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	src "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	target "github.com/Bitspark/bittype/transform/syntax/typescript/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

type countedTarget struct {
	target.Provider
	builds *atomic.Int32
}

func (p countedTarget) Build(ctx context.Context, in g.BuildInput[ts.Language]) (g.Assembly, error) {
	p.builds.Add(1)
	return p.Provider.Build(ctx, in)
}
func planningEngine(t *testing.T, root core.Value, deps ...core.Value) (*engine.Engine, c.Request, *atomic.Int32) {
	t.Helper()
	e, _, input := generationInput(t, root, deps...)
	to := mustLanguage[ts.Language](t, ts.Language{})
	if err := e.RegisterLanguage(to); err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[ts.Language](to, checkts.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterChecker(checker); err != nil {
		t.Fatal(err)
	}
	builds := &atomic.Int32{}
	request, err := genfixture.Compose(e.Generation(), input, to, "planned-output", countedTarget{builds: builds}, typescriptBundles()...)
	if err != nil {
		t.Fatal(err)
	}
	return e, request, builds
}

// builtinListIdentity is the identity of the builtin horos-to-typescript list
// contribution, read from the contribution itself.
func builtinListIdentity(t *testing.T, from representation.Language) g.ContributionIdentity {
	t.Helper()
	list, err := typescriptBundleTable()["list"].Contribution(from)
	if err != nil {
		t.Fatal(err)
	}
	return list.Reference.Identity
}

func TestGenerationPlanPageCellSnapshot(t *testing.T) {
	library := genModule("cells", genMap(), genEntry("Cell", genType(genTuple("horos/wire", genParam(0)), "T")))
	pin, err := decl.Pin(library)
	if err != nil {
		t.Fatal(err)
	}
	imported := genTuple("horos/decl/import", data.MustEncodeText("cells"), core.NewAtom(pin[:]))
	root := genModule("consumer", genMap(genEntry("cells", imported)),
		genEntry("Page", genType(genPage(), "T")), genEntry("Person", genType(genPerson())),
		genEntry("Cell", genType(genScalar("utf8-text"))),
		genEntry("Cells", genType(genLocal("Page", genImported("cells", "Cell", genLocal("Person"))))))
	e, q, builds := planningEngine(t, root, library)
	from := q.Input.Source.Language()
	calls := 0
	cell, err := registerTS(e.Generation(), from, q.To, cellMapping(from, pin.String(), &calls))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := e.PlanGeneration(context.Background(), q)
	if err != nil {
		t.Fatal(err)
	}
	if builds.Load() != 0 || calls != 1 {
		t.Fatal("planning constructed ASTs or missed Cell", builds.Load(), calls)
	}
	found := false
	for _, step := range plan.Steps() {
		// A step names the contribution that was selected for it.
		if strings.Contains(step.Expression, "@bitspark/cell Cell") && step.Binding == cell.String() {
			found = true
		}
	}
	if !found {
		t.Fatal("Cell destination absent from inspectable plan", plan.Steps())
	}
	steps := plan.Steps()
	steps[0].Binding = "changed"
	if reflect.DeepEqual(steps, plan.Steps()) {
		t.Fatal("mutable plan inspection")
	}
	first, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	second, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if first.Snapshot().Same(second.Snapshot()) || calls != 1 || builds.Load() != 2 {
		t.Fatal("run reused capture or reentered mapping", calls, builds.Load())
	}
	if !reflect.DeepEqual(outputTree[ts.Language](t, first), outputTree[ts.Language](t, second)) {
		t.Fatal("runs changed intended AST")
	}
	// Another exact contribution for Cell makes future planning ambiguous, but cannot alter this plan.
	if _, err = registerTS(e.Generation(), from, q.To, cellContribution("cells/typescript/2", from, pin.String(), new(int))); err != nil {
		t.Fatal(err)
	}
	if p, err := e.PlanGeneration(context.Background(), q); !errors.Is(err, g.ErrAmbiguous) || p.Valid() {
		t.Fatal("ambiguous Cell accepted", err)
	}
	if _, err = plan.Run(context.Background()); err != nil || calls != 1 {
		t.Fatal("registry mutation affected captured plan", err)
	}
	tree := placedTS(t, first, n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations})
	checkImportedUse(t, tree, "Cell_1", "@bitspark/cell", "Cell")
}

func TestGenerationPlanReversedLibraryArguments(t *testing.T) {
	library := genModule("pairs", genMap(), genEntry("Pair", genType(genTuple("horos/wire", genParam(0)), "L", "R")))
	pin, err := decl.Pin(library)
	if err != nil {
		t.Fatal(err)
	}
	imported := genTuple("horos/decl/import", data.MustEncodeText("pairs"), core.NewAtom(pin[:]))
	root := genModule("consumer", genMap(genEntry("pairs", imported)), genEntry("Page", genType(genPage(), "T")), genEntry("Person", genType(genPerson())),
		genEntry("Pairs", genType(genLocal("Page", genImported("pairs", "Pair", genLocal("Person"), genScalar("bool"))))))
	e, q, _ := planningEngine(t, root, library)
	contribution := tsConstructor(func(from representation.Language) (tsContribution, error) {
		key, err := src.Named(from, pin.String(), "Pair")
		if err != nil {
			return tsContribution{}, err
		}
		return tsContribution{
			Reference: testReference("TestGenerationPlanReversedLibraryArguments", "reverse-pair"), Target: tsTarget, Profile: rules.RepresentationProfile,
			Domain: g.ForExact(key), Schema: f.EmptySchema,
			Rule: horosRule{Decide: decideNothing, Name: func(_ context.Context, step horosNameStep) ([]horosRequest, error) {
				reader, err := step.Reader()
				if err != nil {
					return nil, err
				}
				// Actuals are selected by the real binder, never by position.
				var left, right g.Use[h.Language]
				for _, a := range step.Subject.Use.Actuals {
					p, err := reader.Parameter(a.Binder)
					if err != nil {
						return nil, err
					}
					switch p.Name {
					case "L":
						left = a.Use
					case "R":
						right = a.Use
					}
				}
				return []horosRequest{g.Reference(left), g.Reference(right)}, nil
			}},
			// The target type takes what was selected for R first and for L second.
			Project: func(_ context.Context, p tsProjection) (tsProjected, error) {
				return tsProjected{Type: model.External("@example/pairs", "Pair", p.Dependency(1).Type, p.Dependency(0).Type)}, nil
			},
		}, nil
	})
	if _, err = registerTS(e.Generation(), q.Input.Source.Language(), q.To, contribution); err != nil {
		t.Fatal(err)
	}
	plan, err := e.PlanGeneration(context.Background(), q)
	if err != nil {
		t.Fatal(err)
	}
	out, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	tree := placedTS(t, out, n.PackageTypes{Module: "@example/pairs", Declaration: "export interface Pair<R,L> { right:R; left:L }"})
	text := formatted(t, typescriptText(), tree)
	if !strings.Contains(text, "Page<Pair<boolean, Person>>") {
		t.Fatal(text)
	}
	checkImportedUse(t, tree, "Pair", "@example/pairs", "Pair")
}

func TestGenerationPlanPreflightAndSourceEvidence(t *testing.T) {
	_, _, foreignInput := generationInput(t, paginationModule())
	cases := map[string]func(*c.Request){
		"empty roots":      func(q *c.Request) { q.Roots = nil },
		"foreign root":     func(q *c.Request) { q.Roots = []ast.NodeRef{{Tree: "foreign", Node: 1}} },
		"duplicate roots":  func(q *c.Request) { q.Roots = []ast.NodeRef{q.Roots[0], q.Roots[0]} },
		"output collision": func(q *c.Request) { q.Output = q.Input.Source.Entry().Tree },
		"missing evidence": func(q *c.Request) { q.Input.Evidence = nil },
		"foreign evidence": func(q *c.Request) { q.Input.Evidence = foreignInput.Evidence },
		"no profile":       func(q *c.Request) { q.Profile = "" },
		"unknown profile":  func(q *c.Request) { q.Profile = "unknown" },
		"budget":           func(q *c.Request) { q.Limits.Operations = 1 },
		"foreign profile":  func(q *c.Request) { q.To, _ = representation.NewLanguage[ts.Language](ts.Language{}, "foreign") },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			e, q, builds := planningEngine(t, paginationModule())
			change(&q)
			p, err := e.PlanGeneration(context.Background(), q)
			if err == nil || p.Valid() || builds.Load() != 0 {
				t.Fatal(err, p.Valid(), builds.Load())
			}
		})
	}
}
func TestGenerationPlanConcurrentRunsAndCancellation(t *testing.T) {
	e, q, builds := planningEngine(t, paginationModule())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if p, err := e.PlanGeneration(ctx, q); !errors.Is(err, context.Canceled) || p.Valid() {
		t.Fatal(err)
	}
	plan, err := e.PlanGeneration(context.Background(), q)
	if err != nil {
		t.Fatal(err)
	}
	if out, err := plan.Run(ctx); !errors.Is(err, context.Canceled) || out.Snapshot().Valid() || builds.Load() != 0 {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, err := plan.Run(context.Background())
			if err != nil || !out.Snapshot().Valid() {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if builds.Load() != 6 {
		t.Fatal(builds.Load())
	}
}

// TestContributionRegistrationFreezesItsValue: the registry keeps the
// contribution as it was registered. What the caller does afterwards with the
// list of owned declarations it handed over, or with a declaration in it, does
// not reach the registered value. A contribution enters under the target it states and once under its identity.
// A second contribution for the list domain is ambiguous at a demanded use
// until the consumer states which of the two answers.
func TestContributionRegistrationFreezesItsValue(t *testing.T) {
	e, q, _ := planningEngine(t, paginationModule())
	from := q.Input.Source.Language()
	const helper = "items"
	owned := []g.OwnedDescription[ts.Language]{{Name: helper, Declaration: g.TargetDeclaration[ts.Language]{
		Name: "Items", Form: model.Form(model.FormAlias), Parameters: []g.TargetParameter{{Name: "T"}},
		Body: model.Collection(g.Parameter[ts.Language](0)),
	}}}
	custom := func(key, target string) tsConstructor {
		return func(from representation.Language) (tsContribution, error) {
			selector, err := src.Builtin(from, "list-of", "")
			return tsContribution{
				Reference: testReference("TestContributionRegistrationFreezesItsValue", key), Target: target, Profile: rules.RepresentationProfile,
				Domain: g.ForForm(selector), Schema: f.EmptySchema,
				Rule:  horosRule{Name: nameElement, Decide: decideNothing},
				Owned: owned,
				Project: func(_ context.Context, p tsProjection) (tsProjected, error) {
					items, err := p.Owned(helper, p.Dependency(0).Type)
					return tsProjected{Type: items}, err
				},
			}, err
		}
	}
	reference, err := registerTS(e.Generation(), from, q.To, custom("custom-list", tsTarget))
	if err != nil {
		t.Fatal(err)
	}
	// The caller writes into the parameters of the declaration it handed over,
	// and then its list names another helper without a declaration. The
	// registered value keeps the declaration it was registered with: the
	// registry copied the list and each declaration in it.
	owned[0].Declaration.Parameters[0].Name = "Changed"
	owned[0] = g.OwnedDescription[ts.Language]{Name: "another"}
	// A contribution that states another target than the one it is registered
	// for is out of scope, and an identity enters once.
	if _, err = registerTS(e.Generation(), from, q.To, custom("foreign-target", "another-target")); !errors.Is(err, g.ErrScope) {
		t.Fatal("a contribution for another target was registered", err)
	}
	if _, err = registerTS(e.Generation(), from, q.To, custom("custom-list", tsTarget)); !errors.Is(err, g.ErrInvalid) {
		t.Fatal("one identity was registered twice", err)
	}
	// The two legitimate list offers are ambiguous at a demanded use.
	if _, err = e.PlanGeneration(context.Background(), q); !errors.Is(err, g.ErrAmbiguous) {
		t.Fatal(err)
	}
	// The consumer selects the registered one. It plans with the helper it
	// was registered with.
	if err = e.Generation().ReplaceContribution(builtinListIdentity(t, from), reference.Identity); err != nil {
		t.Fatal(err)
	}
	plan, err := e.PlanGeneration(context.Background(), q)
	if err != nil {
		t.Fatal("the registered contribution did not keep its value:", err)
	}
	out, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if text := formatted(t, typescriptText(), outputTree[ts.Language](t, out)); !strings.Contains(text, "type Items<T> = T[]") || !strings.Contains(text, "items: Items<T>") {
		t.Fatal(text)
	}
}

// TestPlanningProjectionScopeAndNoFallback: the consumer selects its own
// contribution for the list in place of the builtin one. What that
// contribution does outside its scope fails the planning, and the replaced
// builtin contribution never answers in its place.
//
//   - stale: a use that the rule kept from one planning is no child that the
//     kernel minted in the next one. The first plan still runs.
//   - swallowed-scope-error: the projector asks for a declaration that its
//     contribution does not own, drops the error and returns a valid type.
//   - foreign-native-form: the projector returns a form of another issuer.
//   - missing-helper: the projector asks for a declaration that its
//     contribution does not own and returns the error.
//   - unnamed-dependency: the projector reads a dependency that its rule did
//     not name. It receives no type, and what it builds from that is no
//     target intent.
//   - shared-helper: the rule names a helper by a key. No contribution owns a
//     helper for another one yet, so the kernel refuses the request.
func TestPlanningProjectionScopeAndNoFallback(t *testing.T) {
	const key = "projection"
	for _, mode := range []string{"stale", "swallowed-scope-error", "foreign-native-form", "missing-helper", "unnamed-dependency", "shared-helper"} {
		t.Run(mode, func(t *testing.T) {
			_, _, input := generationInput(t, paginationModule())
			var kept g.Use[h.Language]
			planning := 0
			contribution := tsConstructor(func(from representation.Language) (tsContribution, error) {
				selector, err := src.Builtin(from, "list-of", "")
				return tsContribution{
					Reference: testReference("TestPlanningProjectionScopeAndNoFallback", key), Target: tsTarget, Profile: rules.RepresentationProfile,
					Domain: g.ForForm(selector), Schema: f.EmptySchema,
					Rule: horosRule{Decide: decideNothing, Name: func(ctx context.Context, step horosNameStep) ([]horosRequest, error) {
						reader, err := step.Reader()
						if err != nil {
							return nil, err
						}
						element, err := reader.Element(step.Subject.Use)
						if err != nil {
							return nil, err
						}
						switch {
						case mode == "unnamed-dependency":
							return nil, nil
						case mode == "shared-helper":
							return []horosRequest{g.Reference(element), g.SharedHelper[h.Language]("absent")}, nil
						case mode == "stale" && planning == 1:
							kept = element
						case mode == "stale":
							element = kept
						}
						return []horosRequest{g.Reference(element)}, nil
					}},
					Project: func(_ context.Context, p tsProjection) (tsProjected, error) {
						switch mode {
						case "swallowed-scope-error":
							_, _ = p.Owned("absent")
							return tsProjected{Type: model.Primitive("string")}, nil
						case "foreign-native-form":
							return tsProjected{Type: g.Apply(g.NewFormIssuer[ts.Language]().Form("array"), nil)}, nil
						case "missing-helper":
							helper, err := p.Owned("absent")
							return tsProjected{Type: helper}, err
						}
						return tsProjected{Type: model.Collection(p.Dependency(0).Type)}, nil
					},
				}, err
			})
			plan := func() (g.GenerationPlan, error) {
				planning++
				var registry c.Registry
				q, err := selectedTS(contribution, key)(&registry, input, nativeChecker(t, typescriptText()).Language())
				if err != nil {
					t.Fatal(err)
				}
				return registry.PlanGeneration(context.Background(), q, nativeChecker(t, typescriptText()))
			}
			first, err := plan()
			switch mode {
			case "stale":
				if err != nil || !kept.Valid() {
					t.Fatal(err)
				}
				second, err := plan()
				if !errors.Is(err, g.ErrScope) || !strings.Contains(err.Error(), "a request must name a child") || second.Valid() {
					t.Fatal("a use kept from an earlier planning was planned", err)
				}
				if _, err = first.Run(context.Background()); err != nil {
					t.Fatal(err)
				}
			case "shared-helper":
				var refusal *g.Refusal
				if !errors.Is(err, g.ErrUnsupported) || !errors.As(err, &refusal) || refusal.Owner != g.KernelOwner || refusal.Reason != g.ReasonSharedHelper || first.Valid() {
					t.Fatal("a request for a shared helper produced a plan", err)
				}
			default:
				if !errors.Is(err, g.ErrScope) || first.Valid() {
					t.Fatal("invalid projection produced a plan", mode, err)
				}
			}
		})
	}
}

// Keep the public source capability explicit in this separate consumer package.
var _ g.SourceProvider[h.Language, src.Reader] = src.Provider{}

// TestPlanFreezesProjectedNativeRepresentation: the rule decides and the
// projector projects while the plan is made. A run builds what was selected
// then. It enters neither the rule nor the projector again, so what they would
// answer later does not change the output.
func TestPlanFreezesProjectedNativeRepresentation(t *testing.T) {
	_, _, input := generationInput(t, paginationModule())
	const key = "chosen"
	readonly := false
	calls := 0
	contribution := tsConstructor(func(from representation.Language) (tsContribution, error) {
		selector, err := src.Builtin(from, "list-of", "")
		return tsContribution{
			Reference: testReference("TestPlanFreezesProjectedNativeRepresentation", key), Target: tsTarget, Profile: rules.RepresentationProfile,
			Domain: g.ForForm(selector), Schema: f.EmptySchema,
			Rule: horosRule{Name: nameElement, Decide: func(ctx context.Context, step horosDecideStep) (g.Decision, error) {
				calls++
				return decideNothing(ctx, step)
			}},
			Project: func(_ context.Context, p tsProjection) (tsProjected, error) {
				calls++
				value := p.Dependency(0).Type
				if readonly {
					return tsProjected{Type: model.External("", "ReadonlyArray", value)}, nil
				}
				return tsProjected{Type: model.Collection(value)}, nil
			},
		}, err
	})
	var registry c.Registry
	q, err := selectedTS(contribution, key)(&registry, input, nativeChecker(t, typescriptText()).Language())
	if err != nil {
		t.Fatal(err)
	}
	plan, err := registry.PlanGeneration(context.Background(), q, nativeChecker(t, typescriptText()))
	if err != nil {
		t.Fatal(err)
	}
	plannedCalls := calls
	readonly = true
	out, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	text := formatted(t, typescriptText(), outputTree[ts.Language](t, out))
	if plannedCalls == 0 || strings.Contains(text, "ReadonlyArray") || !strings.Contains(text, "items: T[]") || calls != plannedCalls {
		t.Fatal("run changed selected representation", text, calls, plannedCalls)
	}
}

// A native projection may omit a source argument, but cannot grant admission.
// The unmapped Wire actual is valid source and deliberately has no target rule.
//
// The contribution contract has one declaration policy for a profile and a
// target, and no policy for one named declaration. So the consumer states its
// own policy in place of the builtin one: it decides as the builtin policy
// does, and it gives Ignore no native parameter and binds Ignore's binder to
// string.
func TestGenerationPlanOmittedNativeArgument(t *testing.T) {
	root := genModule("consumer", genMap(),
		genEntry("Ignore", genType(genScalar("utf8-text"), "T")),
		genEntry("Ignored", genType(genLocal("Ignore", genTuple("horos/wire", genScalar("utf8-text"))))))
	e, q, _ := planningEngine(t, root)
	from := q.Input.Source.Language()
	// The builtin policy and the codec of its schema are registered already.
	builtin, err := typescriptBundleTable()["declarations"].Contribution(from)
	if err != nil {
		t.Fatal(err)
	}
	policy := builtin
	policy.Reference = testReference("TestGenerationPlanOmittedNativeArgument", "erase-ignored-parameter")
	policy.Declare = func(ctx context.Context, p tsProjection) (g.DeclarationDescription[ts.Language], error) {
		description, err := builtin.Declare(ctx, p)
		if err == nil && description.Name == "Ignore" {
			// The description of Ignore as this test has always stated it: no
			// target parameter, its binder bound to string, and no origin.
			description.Generic = g.GenericConvention[ts.Language]{Bindings: []g.Expression[ts.Language]{model.Primitive("string")}}
			description.Origin = ast.NodeRef{}
		}
		return description, err
	}
	if err = c.RegisterContribution(e.Generation(), from, q.To, policy); err != nil {
		t.Fatal(err)
	}
	if err = e.Generation().ReplaceContribution(builtin.Reference.Identity, policy.Reference.Identity); err != nil {
		t.Fatal(err)
	}
	plan, err := e.PlanGeneration(context.Background(), q)
	if err != nil {
		t.Fatal(err)
	}
	out, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	text := formatted(t, typescriptText(), outputTree[ts.Language](t, out))
	if !strings.Contains(text, "type Ignore = string") || !strings.Contains(text, "type Ignored = Ignore") {
		t.Fatal(text)
	}
	// The same descriptions never authorize missing source evidence.
	q.Input.Evidence = nil
	if p, err := e.PlanGeneration(context.Background(), q); err == nil || p.Valid() {
		t.Fatal("omitted argument bypassed source admission", err)
	}
}
