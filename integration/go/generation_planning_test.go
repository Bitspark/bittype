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
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	mappings "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/go"
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
	request, err := genfixture.Request(e.Generation(), input, to, "planned-output", countedTarget{builds: builds}, mappings.Libraries(rules.RepresentationProfile)...)
	if err != nil {
		t.Fatal(err)
	}
	return e, request, builds
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
	calls := 0
	_, err = c.RegisterMapping(e.Generation(), q.Input.Source.Language(), q.To, cellMapping(q.Input.Source.Language(), pin.String(), &calls))
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
		if strings.Contains(step.Expression, "@bitspark/cell Cell") {
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
	// Adding another exact offer makes future planning ambiguous, but cannot alter this plan.
	_, err = c.RegisterMapping(e.Generation(), q.Input.Source.Language(), q.To, cellMapping(q.Input.Source.Language(), pin.String(), new(int)))
	if err != nil {
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
	contribution := tsContribution(func(from, to representation.Language) (g.MappingDescription[h.Language, ts.Language, src.Reader], error) {
		key, err := src.Named(from, pin.String(), "Pair")
		if err != nil {
			return g.MappingDescription[h.Language, ts.Language, src.Reader]{}, err
		}
		left, right := g.TypeSlot("L"), g.TypeSlot("R")
		pattern := g.Pattern[h.Language, src.Reader]{Selector: key, Slots: []g.Slot{left, right}, Capture: func(ctx context.Context, s src.Reader, u g.UseData[h.Language]) ([]g.Capture[h.Language], error) {
			if len(u.Actuals) != 2 {
				return nil, g.ErrScope
			}
			return []g.Capture[h.Language]{{Slot: left, Use: u.Actuals[0].Use}, {Slot: right, Use: u.Actuals[1].Use}}, nil
		}}
		return g.MappingDescription[h.Language, ts.Language, src.Reader]{Bindings: []g.BindingDescription[h.Language, ts.Language, src.Reader]{{
			ID: "reverse-pair", Input: pattern, Output: model.External("@example/pairs", "Pair", g.Mapped[ts.Language](right), g.Mapped[ts.Language](left)),
		}}}, nil
	})
	if _, err = c.RegisterMapping(e.Generation(), q.Input.Source.Language(), q.To, contribution); err != nil {
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
func TestMappingRegistrationFreezesDescriptionAndSlots(t *testing.T) {
	e, q, _ := planningEngine(t, paginationModule())
	from := q.Input.Source.Language()
	slot := g.TypeSlot("E")
	pattern, err := rules.Element(from, slot)
	if err != nil {
		t.Fatal(err)
	}
	description := g.MappingDescription[h.Language, ts.Language, src.Reader]{Bindings: []g.BindingDescription[h.Language, ts.Language, src.Reader]{{ID: "custom-list", Input: pattern, Output: model.Collection(g.Mapped[ts.Language](slot))}}}
	lib := tsContribution(func(representation.Language, representation.Language) (g.MappingDescription[h.Language, ts.Language, src.Reader], error) {
		return description, nil
	})
	_, err = c.RegisterMapping(e.Generation(), from, q.To, lib)
	if err != nil {
		t.Fatal(err)
	}
	// The copy now contains a foreign same-named slot; the registered value remains valid.
	description.Bindings[0].Output = model.Collection(g.Mapped[ts.Language](g.TypeSlot("E")))
	if _, err = c.RegisterMapping(e.Generation(), from, q.To, lib); !errors.Is(err, g.ErrScope) {
		t.Fatal("foreign slot accepted", err)
	}
	// The two legitimate list offers are ambiguous at a demanded use.
	if _, err = e.PlanGeneration(context.Background(), q); !errors.Is(err, g.ErrAmbiguous) {
		t.Fatal(err)
	}
}
func TestPlanningExpansionScopeAndNoFallback(t *testing.T) {
	for _, mode := range []string{"stale", "swallowed-scope-error", "foreign-native-form", "missing-helper"} {
		t.Run(mode, func(t *testing.T) {
			_, _, input := generationInput(t, paginationModule())
			var saved g.Expansion[h.Language, ts.Language, src.Reader]
			library := tsContribution(func(from, to representation.Language) (g.MappingDescription[h.Language, ts.Language, src.Reader], error) {
				pattern, err := rules.Pattern(from, "list-of", "")
				return g.MappingDescription[h.Language, ts.Language, src.Reader]{Bindings: []g.BindingDescription[h.Language, ts.Language, src.Reader]{{ID: "expansion", Input: pattern, Expand: func(ctx context.Context, r g.Expansion[h.Language, ts.Language, src.Reader]) (g.Expression[ts.Language], error) {
					saved = r
					u, err := r.Use()
					if err != nil {
						return g.Expression[ts.Language]{}, err
					}
					switch mode {
					case "swallowed-scope-error":
						_, _ = r.Mapped(u.Use)
						return model.Primitive("string"), nil
					case "foreign-native-form":
						return g.Native(g.NewFormIssuer[ts.Language]().Form("array"), nil), nil
					case "missing-helper":
						return g.Helper[ts.Language](g.NewHelperKey("absent")), nil
					default:
						s, err := r.Source()
						if err != nil {
							return g.Expression[ts.Language]{}, err
						}
						operand, err := s.Operand(u.Use, 1)
						if err != nil {
							return g.Expression[ts.Language]{}, err
						}
						v, err := r.Mapped(operand)
						return model.Collection(v), err
					}
				}}}}, err
			})
			var registry c.Registry
			q, err := selectedTS(library, "expansion")(&registry, input, nativeChecker(t, typescriptText()).Language())
			if err != nil {
				t.Fatal(err)
			}
			plan, err := registry.PlanGeneration(context.Background(), q, nativeChecker(t, typescriptText()))
			if mode == "stale" {
				if err != nil {
					t.Fatal(err)
				}
				if _, err = saved.Source(); !errors.Is(err, g.ErrScope) {
					t.Fatal("retained expansion is live", err)
				}
				if _, err = plan.Run(context.Background()); err != nil {
					t.Fatal(err)
				}
			} else if err == nil || plan.Valid() {
				t.Fatal("invalid expansion produced a plan", mode, err)
			}
		})
	}
}

// Keep the public source capability explicit in this separate consumer package.
var _ g.SourceProvider[h.Language, src.Reader] = src.Provider{}

func TestPlanFreezesExpandedNativeRepresentation(t *testing.T) {
	_, _, input := generationInput(t, paginationModule())
	readonly := false
	calls := 0
	library := tsContribution(func(from, to representation.Language) (g.MappingDescription[h.Language, ts.Language, src.Reader], error) {
		pattern, err := rules.Pattern(from, "list-of", "")
		return g.MappingDescription[h.Language, ts.Language, src.Reader]{Bindings: []g.BindingDescription[h.Language, ts.Language, src.Reader]{{ID: "chosen", Input: pattern, Expand: func(ctx context.Context, r g.Expansion[h.Language, ts.Language, src.Reader]) (g.Expression[ts.Language], error) {
			calls++
			source, err := r.Source()
			if err != nil {
				return g.Expression[ts.Language]{}, err
			}
			use, err := r.Use()
			if err != nil {
				return g.Expression[ts.Language]{}, err
			}
			element, err := source.Operand(use.Use, 1)
			if err != nil {
				return g.Expression[ts.Language]{}, err
			}
			value, err := r.Mapped(element)
			if err != nil {
				return g.Expression[ts.Language]{}, err
			}
			if readonly {
				return model.External("", "ReadonlyArray", value), nil
			}
			return model.Collection(value), nil
		}}}}, err
	})
	var registry c.Registry
	q, err := selectedTS(library, "chosen")(&registry, input, nativeChecker(t, typescriptText()).Language())
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
	if strings.Contains(text, "ReadonlyArray") || !strings.Contains(text, "items: T[]") || calls != plannedCalls {
		t.Fatal("run changed selected representation", text, calls, plannedCalls)
	}
}

// A native projection may omit a source argument, but cannot grant admission.
// The unmapped Wire actual is valid source and deliberately has no target rule.
func TestGenerationPlanOmittedNativeArgument(t *testing.T) {
	root := genModule("consumer", genMap(),
		genEntry("Ignore", genType(genScalar("utf8-text"), "T")),
		genEntry("Ignored", genType(genLocal("Ignore", genTuple("horos/wire", genScalar("utf8-text"))))))
	pin, err := decl.Pin(root)
	if err != nil {
		t.Fatal(err)
	}
	e, q, _ := planningEngine(t, root)
	contribution := tsContribution(func(from, to representation.Language) (g.MappingDescription[h.Language, ts.Language, src.Reader], error) {
		key, err := src.Named(from, pin.String(), "Ignore")
		if err != nil {
			return g.MappingDescription[h.Language, ts.Language, src.Reader]{}, err
		}
		return g.MappingDescription[h.Language, ts.Language, src.Reader]{Declarations: []g.DeclarationRecipe[h.Language, ts.Language, src.Reader]{{
			ID: "erase-ignored-parameter", Selector: key,
			Expand: func(ctx context.Context, r g.DeclarationRequest[h.Language, ts.Language, src.Reader]) (g.DeclarationDescription[ts.Language], error) {
				open, err := r.Open()
				if err != nil {
					return g.DeclarationDescription[ts.Language]{}, err
				}
				body, err := r.Mapped(open.Body)
				return g.DeclarationDescription[ts.Language]{Name: "Ignore", Form: model.Form("alias"), Body: body, Generic: g.GenericConvention[ts.Language]{Bindings: []g.Expression[ts.Language]{model.Primitive("string")}}}, err
			},
		}}}, nil
	})
	if _, err = c.RegisterMapping(e.Generation(), q.Input.Source.Language(), q.To, contribution); err != nil {
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
