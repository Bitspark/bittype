package integration_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/engine/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func analysisKey(name string) semantics.Key {
	return semantics.Key{Owner: "independent-analysis", Name: name}
}

func analysisContract() analysis.Contract {
	return analysis.Contract{Semantic: semantic.Contract{
		Profile: "independent-facts-v1",
		Obligations: []semantic.Obligation{
			{Key: analysisKey("validity"), Domain: analysisKey("entry-program")},
			{Key: analysisKey("bindings"), Domain: analysisKey("entry-closure")},
		},
		ModelRequirements: []semantic.Requirement{
			{Obligation: analysisKey("validity"), Verdict: semantic.Holds},
			{Obligation: analysisKey("bindings"), Verdict: semantic.Undetermined},
		},
	}}
}

func analysisReports() []semantic.Report {
	return []semantic.Report{
		{Obligation: analysisKey("validity"), Status: semantic.Complete, Verdict: semantic.Holds},
		{Obligation: analysisKey("bindings"), Status: semantic.Complete},
	}
}

func analysisFacts() semantics.Data {
	at := semantics.Occurrence{Site: ast.NodeRef{Tree: "input", Node: 3}, Role: analysisKey("type-use")}
	return semantics.Data{
		Entities: []semantics.Entity{
			{Kind: analysisKey("declaration"), Name: "Page", Declarations: []ast.NodeRef{{Tree: "input", Node: 2}}, Parameters: []semantics.Parameter{{Entity: 2, Sort: analysisKey("type")}}, Body: 2},
			{Kind: analysisKey("parameter"), Name: "T", Owner: 1, Declarations: []ast.NodeRef{{Tree: "input", Node: 3}}},
		},
		Terms: []semantics.Term{
			{Kind: semantics.Reference, Target: 2},
			{Kind: semantics.ProviderExpression, Operator: analysisKey("array"), Encoding: analysisKey("array-v1"), Data: []byte{3}, Operands: []semantics.Operand{{Locator: []uint32{1}, Term: 1}}},
		},
		Bindings:     []semantics.Binding{{At: at, Target: 2}},
		Expressions:  []semantics.Expression{{At: at, Term: 1}},
		Dependencies: []semantics.Dependency{{From: 1, To: 1, Relation: analysisKey("recursive-reference"), At: &at}},
	}
}

type analysisChecker[L any] struct {
	id  language.ID
	run func(context.Context, analysis.Request[L]) (semantic.Result, error)
}

func (p analysisChecker[L]) ID() language.ID { return p.id }
func (p analysisChecker[L]) Name() string    { return "independent " + string(p.id) }
func (p analysisChecker[L]) Analyze(ctx context.Context, request analysis.Request[L]) (semantic.Result, error) {
	if p.run != nil {
		return p.run(ctx, request)
	}
	data := analysisFacts()
	return request.Publish(semantic.Publication{Facts: &data, Reports: analysisReports()})
}

func analysisFixture[L any](t *testing.T, definition representation.Language) analysisCase {
	t.Helper()
	tree := exampleTree[L]()
	tree.Nodes[1].References = nil // the old fixture's external library is not supplied
	envelope, err := representation.Pack(definition, tree)
	if err != nil {
		t.Fatal(err)
	}
	return analysisCase{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: []representation.AST{envelope}}
}

func analysisAnalyzer[L any](t *testing.T, definition representation.Language, run func(context.Context, analysis.Request[L]) (semantic.Result, error)) analysis.Analyzer {
	t.Helper()
	analyzer, err := analysis.NewAnalyzer[L](definition, analysisContract(), analysisChecker[L]{id: definition.ID(), run: run})
	if err != nil {
		t.Fatal(err)
	}
	return analyzer
}

func analysisEngine(t *testing.T, run func(context.Context, analysis.Request[amber]) (semantic.Result, error)) (*engine.Engine, analysisCase) {
	t.Helper()
	definition := definition[amber](t, "analysis-language")
	kernel := engine.New()
	registerChecked[amber](t, kernel, definition)
	if err := kernel.RegisterAnalyzer(analysisAnalyzer[amber](t, definition, run)); err != nil {
		t.Fatal(err)
	}
	return kernel, analysisFixture[amber](t, definition)
}

func TestAnalysisRegistrationBoundaries(t *testing.T) {
	a := definition[amber](t, "analysis-language")
	if _, err := analysis.NewAnalyzer[blue](a, analysisContract(), analysisChecker[blue]{id: a.ID()}); !errors.Is(err, representation.ErrTypeMismatch) {
		t.Fatal(err)
	}
	var nilProvider *analysisChecker[amber]
	if _, err := analysis.NewAnalyzer[amber](a, analysisContract(), nilProvider); !errors.Is(err, semantic.ErrInvalidOperation) {
		t.Fatal(err)
	}
	if _, err := analysis.NewAnalyzer[amber](a, analysisContract(), analysisChecker[amber]{id: "wrong"}); !errors.Is(err, semantic.ErrInvalidOperation) {
		t.Fatal(err)
	}
	for _, alter := range []func(*analysis.Contract){
		func(c *analysis.Contract) { c.Semantic.Profile = " " },
		func(c *analysis.Contract) {
			c.Semantic.Obligations = append(c.Semantic.Obligations, c.Semantic.Obligations[0])
		},
		func(c *analysis.Contract) { c.Semantic.ModelRequirements[0].Obligation = analysisKey("absent") },
	} {
		contract := analysisContract()
		alter(&contract)
		if _, err := analysis.NewAnalyzer[amber](a, contract, analysisChecker[amber]{id: a.ID()}); !errors.Is(err, semantic.ErrInvalidOperation) {
			t.Fatal(err)
		}
	}
	kernel := engine.New()
	good := analysisAnalyzer[amber](t, a, nil)
	if err := kernel.RegisterAnalyzer(good); !errors.Is(err, engine.ErrUnknownLanguage) {
		t.Fatal(err)
	}
	registerChecked[amber](t, kernel, a)
	if err := kernel.RegisterAnalyzer(analysis.Analyzer{}); !errors.Is(err, semantic.ErrInvalidOperation) {
		t.Fatal(err)
	}
	wrongMarker := definition[blue](t, a.ID())
	if err := kernel.RegisterAnalyzer(analysisAnalyzer[blue](t, wrongMarker, nil)); !errors.Is(err, representation.ErrLanguageMismatch) {
		t.Fatalf("same ID cast marker: %v", err)
	}
	wrongProfile, err := representation.NewLanguage[amber](testChecker[amber]{id: a.ID()}, "another-syntax-profile")
	if err != nil {
		t.Fatal(err)
	}
	if err := kernel.RegisterAnalyzer(analysisAnalyzer[amber](t, wrongProfile, nil)); !errors.Is(err, representation.ErrLanguageMismatch) {
		t.Fatalf("syntax profile ignored: %v", err)
	}
	if err := kernel.RegisterAnalyzer(good); err != nil {
		t.Fatal(err)
	}
	if err := kernel.RegisterAnalyzer(good); !errors.Is(err, engine.ErrDuplicateAnalyzer) {
		t.Fatal(err)
	}
	for name, input := range map[string]analysisCase{
		"marker":  analysisFixture[blue](t, wrongMarker),
		"profile": analysisFixture[amber](t, wrongProfile),
	} {
		t.Run(name, func(t *testing.T) {
			result, err := analyzeCase(kernel, context.Background(), input)
			if err == nil || result.Valid() {
				t.Fatalf("wrong input boundary admitted: %v", err)
			}
		})
	}
}

func TestAnalysisMissingCapability(t *testing.T) {
	kernel := engine.New()
	a := definition[amber](t, "unfamiliar")
	input := analysisFixture[amber](t, a)
	if _, err := analyzeCase(kernel, context.Background(), input); !errors.Is(err, engine.ErrUnknownLanguage) {
		t.Fatal(err)
	}
	registerChecked[amber](t, kernel, a)
	if _, err := analyzeCase(kernel, context.Background(), input); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal(err)
	}
	if _, err := analyzeCase(kernel, nil, input); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal(err)
	}
	if _, err := analyzeCase(kernel, context.Background(), analysisCase{}); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal(err)
	}
}

func TestAnalysisRejectsForeignRequestResult(t *testing.T) {
	changes := map[string]func(*testing.T, *analysisCase){
		"identical": func(*testing.T, *analysisCase) {},
		"syntax": func(t *testing.T, input *analysisCase) {
			d := definition[amber](t, "analysis-language")
			tree, err := representation.Unpack[amber](d, input.Trees[0])
			if err != nil {
				t.Fatal(err)
			}
			tree.Nodes[1].Attributes[0].Text = "Changed"
			input.Trees[0], err = representation.Pack(d, tree)
			if err != nil {
				t.Fatal(err)
			}
		},
		"configuration": func(_ *testing.T, input *analysisCase) {
			input.Configuration = []semantic.Setting{{Key: analysisKey("mode"), Value: []byte("changed")}}
		},
		"limits": func(_ *testing.T, input *analysisCase) {
			input.Limits = semantic.DefaultLimits()
			input.Limits.MaxDiagnostics--
		},
		"dependency": func(t *testing.T, input *analysisCase) {
			d := definition[amber](t, "analysis-language")
			tree := exampleTree[amber]()
			tree.ID, tree.Nodes[1].References = "dependency", nil
			dependency, err := representation.Pack(d, tree)
			if err != nil {
				t.Fatal(err)
			}
			input.Trees = append(input.Trees, dependency)
		},
	}
	for name, change := range changes {
		t.Run(name, func(t *testing.T) {
			var retained semantic.Result
			kernel, input := analysisEngine(t, func(_ context.Context, request analysis.Request[amber]) (semantic.Result, error) {
				if retained.Valid() {
					return retained, nil
				}
				data := analysisFacts()
				var err error
				retained, err = request.Publish(semantic.Publication{Facts: &data, Reports: analysisReports()})
				return retained, err
			})
			if _, err := analyzeCase(kernel, context.Background(), input); err != nil {
				t.Fatal(err)
			}
			change(t, &input)
			result, err := analyzeCase(kernel, context.Background(), input)
			if !errors.Is(err, semantic.ErrStaleResult) || result.Valid() {
				t.Fatalf("stale result accepted: %v", err)
			}
		})
	}
}

func TestAnalysisOwnershipAndPublication(t *testing.T) {
	want := analysisFacts()
	kernel, input := analysisEngine(t, func(_ context.Context, request analysis.Request[amber]) (semantic.Result, error) {
		copies := request.Inputs()
		copies[0].Nodes[1].Attributes[0].Text = "mutated input"
		copies[0].Nodes[0].Children[0].Node = 99
		copies[0].Nodes[1].Origin.Nodes[0].Tree = "changed history"
		if got := request.Inputs()[0]; got.Nodes[1].Attributes[0].Text != "Page" || got.Nodes[0].Children[0].Node != 2 || got.Nodes[1].Origin.Nodes[0].Tree != "earlier" {
			return semantic.Result{}, errors.New("request inputs alias accessor")
		}
		configuration := request.Configuration()
		configuration[0].Value[0] = 'X'
		requested := request.Requested()
		requested[0] = analysisKey("mutated")
		contract := request.Contract()
		contract.Obligations[0].Domain = analysisKey("mutated")
		if string(request.Configuration()[0].Value) != "original" || request.Requested()[0] != analysisKey("validity") || request.Contract().Obligations[0].Domain != analysisKey("entry-program") {
			return semantic.Result{}, errors.New("captured metadata aliases accessor")
		}
		data, reports := analysisFacts(), analysisReports()
		site := ast.NodeRef{Tree: "input", Node: 2}
		diagnostics := []semantic.Diagnostic{{Code: analysisKey("note"), Subject: analysisKey("entry-program"), Message: "original", At: &site}}
		result, err := request.Publish(semantic.Publication{Facts: &data, Reports: reports, Diagnostics: diagnostics})
		data.Entities[0].Parameters[0].Sort = analysisKey("mutated")
		data.Terms[1].Data[0] = 4
		data.Terms[1].Operands[0].Locator[0] = 7
		data.Dependencies[0].At.Site.Node = 99
		reports[0].Verdict = semantic.Rejected
		diagnostics[0].Message = "mutated diagnostic"
		site.Node = 99
		return result, err
	})
	input.Configuration = []semantic.Setting{{Key: analysisKey("mode"), Value: []byte("original")}}
	result, err := analyzeCase(kernel, context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	model, present := result.Model()
	if !present || !reflect.DeepEqual(model.Data(), want) {
		t.Fatal("publication aliases provider graph")
	}
	if result.Reports()[0].Verdict != semantic.Holds || result.Diagnostics()[0].Message != "original" || result.Diagnostics()[0].At.Node != 2 {
		t.Fatal("publication aliases provider metadata")
	}
	input.Configuration[0].Value[0] = 'Y'
	result.Configuration()[0].Value[0] = 'Z'
	result.Reports()[0].Verdict = semantic.Rejected
	result.Diagnostics()[0].At.Node = 99
	result.Contract().Obligations[0].Domain = analysisKey("changed")
	extracted := model.Data()
	extracted.Terms[1].Data[0] = 9
	extracted.Entities[0].Declarations[0].Node = 99
	if !reflect.DeepEqual(model.Data(), want) || string(result.Configuration()[0].Value) != "original" || result.Reports()[0].Verdict != semantic.Holds || result.Diagnostics()[0].At.Node != 2 || result.Contract().Obligations[0].Domain != analysisKey("entry-program") {
		t.Fatal("result accessor exposes private storage")
	}
}

func TestAnalysisFailureDiscardsOutput(t *testing.T) {
	sentinel := errors.New("provider infrastructure failure")
	calls := 0
	kernel, input := analysisEngine(t, func(_ context.Context, request analysis.Request[amber]) (semantic.Result, error) {
		calls++
		data := analysisFacts()
		result, err := request.Publish(semantic.Publication{Facts: &data, Reports: analysisReports()})
		if err != nil {
			return semantic.Result{}, err
		}
		return result, sentinel
	})
	result, err := analyzeCase(kernel, context.Background(), input)
	if !errors.Is(err, sentinel) || result.Valid() || calls != 1 {
		t.Fatalf("output escaped failure or callback retried: %v", err)
	}
	kernel, input = analysisEngine(t, func(_ context.Context, request analysis.Request[amber]) (semantic.Result, error) {
		reports := analysisReports()
		reports[0].Verdict = semantic.Rejected
		reports[1].Status = semantic.Incomplete
		return request.Publish(semantic.Publication{Reports: reports, Diagnostics: []semantic.Diagnostic{{Code: analysisKey("invalid-declaration"), Subject: analysisKey("entry-program"), Message: "definite counterexample"}}})
	})
	result, err = analyzeCase(kernel, context.Background(), input)
	if err != nil || !result.Valid() || result.HasModel() {
		t.Fatalf("semantic rejection lost: %v", err)
	}
	validity, found := result.Report(analysisKey("validity"))
	bindings, _ := result.Report(analysisKey("bindings"))
	if !found || validity.Status != semantic.Complete || validity.Verdict != semantic.Rejected || bindings.Status != semantic.Incomplete {
		t.Fatal("completion became validity or complete enumeration")
	}
}

func TestAnalysisRejectsInvalidPublications(t *testing.T) {
	cases := map[string]func(*semantic.Publication){
		"missing-report":            func(p *semantic.Publication) { p.Reports = p.Reports[:1] },
		"duplicate-report":          func(p *semantic.Publication) { p.Reports[1] = p.Reports[0] },
		"unreported":                func(p *semantic.Publication) { p.Reports[1].Status = semantic.Unreported },
		"unknown-status":            func(p *semantic.Publication) { p.Reports[1].Status = 99 },
		"unsupported-holds":         func(p *semantic.Publication) { p.Reports[0].Status = semantic.Unsupported },
		"missing-model-requirement": func(p *semantic.Publication) { p.Reports[0].Verdict = semantic.Rejected },
		"historical-only-site": func(p *semantic.Publication) {
			p.Facts.Entities[0].Declarations[0] = ast.NodeRef{Tree: "earlier", Node: 4}
		},
		"dangling-entity": func(p *semantic.Publication) { p.Facts.Bindings[0].Target = 99 },
		"zero-node":       func(p *semantic.Publication) { p.Facts.Entities[0].Declarations[0].Node = 0 },
		"diagnostic-site": func(p *semantic.Publication) {
			p.Diagnostics = []semantic.Diagnostic{{Code: analysisKey("note"), Subject: analysisKey("entry"), At: &ast.NodeRef{Tree: "input", Node: 99}}}
		},
	}
	for name, alter := range cases {
		t.Run(name, func(t *testing.T) {
			kernel, input := analysisEngine(t, func(_ context.Context, request analysis.Request[amber]) (semantic.Result, error) {
				data := analysisFacts()
				publication := semantic.Publication{Facts: &data, Reports: analysisReports()}
				alter(&publication)
				return request.Publish(publication)
			})
			result, err := analyzeCase(kernel, context.Background(), input)
			if !errors.Is(err, semantic.ErrInvalidResult) || result.Valid() {
				t.Fatalf("invalid result accepted: %v", err)
			}
		})
	}
}

func TestAnalysisCaptureAndOutputLimits(t *testing.T) {
	cases := map[string]func(*analysisCase){
		"duplicate-tree": func(input *analysisCase) { input.Trees = append(input.Trees, input.Trees[0]) },
		"zero-entry":     func(input *analysisCase) { input.Entry.Node = 0 },
		"unknown-entry":  func(input *analysisCase) { input.Entry.Tree = "absent" },
		"duplicate-request": func(input *analysisCase) {
			input.Requested = []semantics.Key{analysisKey("validity"), analysisKey("validity")}
		},
		"duplicate-configuration": func(input *analysisCase) {
			input.Configuration = []semantic.Setting{{Key: analysisKey("mode")}, {Key: analysisKey("mode")}}
		},
		"partial-limits": func(input *analysisCase) { input.Limits.MaxInputNodes = 1 },
	}
	for name, alter := range cases {
		t.Run(name, func(t *testing.T) {
			calls := 0
			kernel, input := analysisEngine(t, func(context.Context, analysis.Request[amber]) (semantic.Result, error) {
				calls++
				return semantic.Result{}, nil
			})
			alter(&input)
			result, err := analyzeCase(kernel, context.Background(), input)
			if !errors.Is(err, semantic.ErrInvalidRequest) || result.Valid() || calls != 0 {
				t.Fatalf("invalid input reached provider: %v", err)
			}
		})
	}
	for name, alter := range map[string]func(*semantic.Limits){
		"nodes":          func(l *semantic.Limits) { l.MaxInputNodes = 2 },
		"entries":        func(l *semantic.Limits) { l.MaxInputEntries = 1 },
		"input-bytes":    func(l *semantic.Limits) { l.MaxInputBytes = 1 },
		"metadata-bytes": func(l *semantic.Limits) { l.MaxConfigurationBytes = 1 },
		"model-bytes":    func(l *semantic.Limits) { l.Model.MaxBytes = 1 },
	} {
		t.Run(name, func(t *testing.T) {
			kernel, input := analysisEngine(t, nil)
			input.Limits = semantic.DefaultLimits()
			alter(&input.Limits)
			result, err := analyzeCase(kernel, context.Background(), input)
			if !errors.Is(err, semantic.ErrLimit) || result.Valid() {
				t.Fatalf("resource limit accepted or became a semantic rejection: %v", err)
			}
		})
	}
	var zero analysis.Request[amber]
	if result, err := zero.Publish(semantic.Publication{}); !errors.Is(err, semantic.ErrInvalidRequest) || result.Valid() {
		t.Fatal(err)
	}
}

func TestAnalysisBudgetsPrecedeCheckerAndCopy(t *testing.T) {
	checks := 0
	provider := testChecker[amber]{id: "analysis-language", visit: func(ast.AST[amber]) error { checks++; return nil }}
	d, err := representation.NewLanguage[amber](provider, "test-v1")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[amber](d, provider)
	if err != nil {
		t.Fatal(err)
	}
	input := analysisFixture[amber](t, d)
	checks = 0
	info, err := input.Trees[0].Inspect()
	if err != nil || info.ID != "input" || info.Root != 1 || info.Nodes != 3 || info.Entries != 6 || info.TextBytes != 79 || checks != 0 {
		t.Fatalf("incorrect or active structural inspection: %+v, %v", info, err)
	}
	info.Nodes = 99
	if fresh, err := input.Trees[0].Inspect(); err != nil || fresh.Nodes != 3 {
		t.Fatal("structural metadata aliases envelope")
	}
	if _, err := (representation.AST{}).Inspect(); !errors.Is(err, representation.ErrInvalidAST) {
		t.Fatal(err)
	}
	analyzer := analysisAnalyzer[amber](t, d, nil)
	for _, alter := range []func(*semantic.Limits){
		func(l *semantic.Limits) { l.MaxInputNodes = 2 },
		func(l *semantic.Limits) { l.MaxInputEntries = 5 },
		func(l *semantic.Limits) { l.MaxInputBytes = 78 },
	} {
		input.Limits = semantic.DefaultLimits()
		alter(&input.Limits)
		if result, err := runAnalyzerCase(analyzer, checker, context.Background(), input); !errors.Is(err, semantic.ErrLimit) || result.Valid() || checks != 0 {
			t.Fatalf("budget checked after callback/copy boundary: %v (%d callbacks)", err, checks)
		}
	}
	input.Limits = semantic.DefaultLimits()
	if result, err := runAnalyzerCase(analyzer, checker, context.Background(), input); err != nil || !result.HasModel() || checks != 1 {
		t.Fatalf("valid request did not execute checker: %v", err)
	}
}

func TestAnalysisUnknownObligationsAndAdditionalCoverage(t *testing.T) {
	unknown := analysisKey("future-relation")
	kernel, input := analysisEngine(t, func(_ context.Context, request analysis.Request[amber]) (semantic.Result, error) {
		data := analysisFacts()
		reports := append(analysisReports(), semantic.Report{Obligation: unknown, Status: semantic.Unsupported})
		return request.Publish(semantic.Publication{Facts: &data, Reports: reports})
	})
	input.Requested = []semantics.Key{unknown}
	result, err := analyzeCase(kernel, context.Background(), input)
	if err != nil || !result.HasModel() {
		t.Fatalf("supported graph with unsupported extra obligation: %v", err)
	}
	report, ok := result.Report(unknown)
	if !ok || report.Status != semantic.Unsupported || report.Verdict != semantic.Undetermined {
		t.Fatal("unsupported obligation became a verdict")
	}
	if _, err := analyzeCase(kernel, context.Background(), analysisCase{Entry: input.Entry, Trees: input.Trees}); !errors.Is(err, semantic.ErrInvalidResult) {
		t.Fatalf("unrequested unknown report was accepted: %v", err)
	}
}

func TestAnalysisCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	kernel, input := analysisEngine(t, func(context.Context, analysis.Request[amber]) (semantic.Result, error) {
		calls++
		return semantic.Result{}, nil
	})
	if result, err := analyzeCase(kernel, ctx, input); !errors.Is(err, context.Canceled) || result.Valid() || calls != 0 {
		t.Fatal(err)
	}
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	kernel, input = analysisEngine(t, func(_ context.Context, request analysis.Request[amber]) (semantic.Result, error) {
		data := analysisFacts()
		result, err := request.Publish(semantic.Publication{Facts: &data, Reports: analysisReports()})
		cancel()
		return result, err
	})
	if result, err := analyzeCase(kernel, ctx, input); !errors.Is(err, context.Canceled) || result.Valid() {
		t.Fatalf("cancelled output escaped: %v", err)
	}
}

func TestEngineConcurrentAnalysisAndReentry(t *testing.T) {
	kernel := engine.New()
	d := definition[amber](t, "analysis-language")
	registerChecked[amber](t, kernel, d)
	input := analysisFixture[amber](t, d)
	analyzer := analysisAnalyzer[amber](t, d, func(_ context.Context, request analysis.Request[amber]) (semantic.Result, error) {
		_ = kernel.Languages()
		if _, ok := kernel.Language(d.ID()); !ok {
			return semantic.Result{}, errors.New("reentrant registry read lost language")
		}
		data := analysisFacts()
		return request.Publish(semantic.Publication{Facts: &data, Reports: analysisReports()})
	})
	if err := kernel.RegisterAnalyzer(analyzer); err != nil {
		t.Fatal(err)
	}
	var group sync.WaitGroup
	failures := make(chan error, 16)
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			extra, err := representation.NewLanguage[echo](testChecker[echo]{id: language.ID(fmt.Sprintf("analysis-extra-%d", i))}, "test-v1")
			if err == nil {
				err = kernel.RegisterLanguage(extra)
			}
			if err != nil {
				failures <- err
				return
			}
			for range 3 {
				result, err := analyzeCase(kernel, context.Background(), input)
				if err != nil {
					failures <- err
					return
				}
				if !result.HasModel() {
					failures <- errors.New("concurrent result lost model")
					return
				}
			}
		}(i)
	}
	done := make(chan struct{})
	go func() { group.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("analysis callbacks held registry lock or failed to complete")
	}
	close(failures)
	for err := range failures {
		t.Error(err)
	}
	if len(kernel.Languages()) != 9 {
		t.Fatal("concurrent registration lost entries")
	}
}

func TestAnalysisCheckerCannotReplacePreflightedInputs(t *testing.T) {
	var caller analysisCase
	var replacement representation.AST
	armed := false
	provider := testChecker[amber]{id: "analysis-language", visit: func(value ast.AST[amber]) error {
		if armed && value.ID == "input" {
			caller.Trees[1] = replacement
		}
		return nil
	}}
	d, err := representation.NewLanguage[amber](provider, "test-v1")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[amber](d, provider)
	if err != nil {
		t.Fatal(err)
	}
	caller = analysisFixture[amber](t, d)
	dependency := exampleTree[amber]()
	dependency.ID, dependency.Nodes[1].References = "dependency", nil
	second, err := representation.Pack(d, dependency)
	if err != nil {
		t.Fatal(err)
	}
	caller.Trees = append(caller.Trees, second)
	oversized := ast.AST[amber]{ID: "replacement", Root: 1, Nodes: []ast.Node[amber]{{Kind: "document"}}}
	for range 8 {
		oversized.Nodes = append(oversized.Nodes, ast.Node[amber]{Kind: "declaration"})
		oversized.Nodes[0].Children = append(oversized.Nodes[0].Children, ast.Child{Role: "declaration", Node: ast.NodeID(len(oversized.Nodes))})
	}
	replacement, err = representation.Pack(d, oversized)
	if err != nil {
		t.Fatal(err)
	}
	caller.Limits = semantic.DefaultLimits()
	caller.Limits.MaxInputNodes = 6 // exactly the two original three-node trees
	analyzer := analysisAnalyzer[amber](t, d, func(_ context.Context, request analysis.Request[amber]) (semantic.Result, error) {
		inputs := request.Inputs()
		if len(inputs) != 2 || inputs[1].ID != "dependency" || len(inputs[1].Nodes) != 3 {
			return semantic.Result{}, errors.New("checker replaced an input after catalogue/budget preflight")
		}
		data := analysisFacts()
		return request.Publish(semantic.Publication{Facts: &data, Reports: analysisReports()})
	})
	armed = true
	result, err := runAnalyzerCase(analyzer, checker, context.Background(), caller)
	if err != nil || !result.HasModel() {
		t.Fatalf("captured input collection changed: %v", err)
	}
	changed, err := caller.Trees[1].Inspect()
	if err != nil || changed.ID != "replacement" {
		t.Fatal("checker mutation control did not execute")
	}
}
