package integration_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/engine/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Each implementation deliberately supplies only one capability.
type stageIdentity struct{ id language.ID }

func (p stageIdentity) ID() language.ID { return p.id }
func (p stageIdentity) Name() string { return "Independent stage" }

type parserOnly[L any] struct {
	stageIdentity
	run func(parse.Input) (cst.CST[L], error)
}
func (p parserOnly[L]) Parse(input parse.Input) (cst.CST[L], error) { return p.run(input) }

type formatterOnly[L any] struct {
	stageIdentity
	run func(ast.View[L]) (cst.CST[L], error)
}
func (p formatterOnly[L]) Format(input ast.View[L]) (cst.CST[L], error) { return p.run(input) }

type printerOnly[L any] struct {
	stageIdentity
	run func(cst.CST[L]) (string, error)
}
func (p printerOnly[L]) Print(input cst.CST[L]) (string, error) { return p.run(input) }

func TestEngineIndependentStages(t *testing.T) {
	identity := stageIdentity{"split"}
	d, err := representation.NewLanguage[amber](identity, "v1")
	if err != nil { t.Fatal(err) }
	e := engine.New()
	register(t, e, d)
	leaf := "page"
	parser, err := parse.NewParser[amber](d, parserOnly[amber]{identity, func(parse.Input) (cst.CST[amber], error) { return exampleConcrete[amber](&leaf), nil }})
	if err != nil { t.Fatal(err) }
	if err := e.RegisterParser(parser); err != nil { t.Fatal(err) }
	parsed, err := e.Parse(d.ID(), parse.Input{})
	if err != nil { t.Fatal("parser-only registration", err) }
	if err := e.Check(parsed.ASTView()); !errors.Is(err, engine.ErrMissingCapability) { t.Fatal(err) }
	if _, err := e.Format(parsed.ASTView()); !errors.Is(err, engine.ErrMissingCapability) { t.Fatal(err) }
	if _, err := e.Print(parsed); !errors.Is(err, engine.ErrMissingCapability) { t.Fatal(err) }
	if _, err := e.TransformTo(parsed.ASTView(), d.ID()); !errors.Is(err, engine.ErrMissingCapability) { t.Fatal("zero-hop admission skipped", err) }
	if _, err := e.Capture(context.Background(), semantic.Input{Trees: []representation.AST{parsed.ASTView()}}); !errors.Is(err, engine.ErrMissingCapability) { t.Fatal(err) }
	formats := 0
	formatter, err := format.NewFormatter[amber](d, formatterOnly[amber]{identity, func(ast.View[amber]) (cst.CST[amber], error) { formats++; return exampleConcrete[amber](&leaf), nil }})
	if err != nil { t.Fatal(err) }
	if err := e.RegisterFormatter(formatter); err != nil { t.Fatal(err) }
	if _, err := e.Render(parsed.ASTView()); !errors.Is(err, engine.ErrMissingCapability) || formats != 0 { t.Fatal("render must preflight its printer", err) }
	printer, err := print.NewPrinter[amber](d, printerOnly[amber]{identity, func(c cst.CST[amber]) (string, error) { return *c.Nodes[0].Text, nil }})
	if err != nil { t.Fatal(err) }
	if err := e.RegisterPrinter(printer); err != nil { t.Fatal(err) }
	if output, err := e.Render(parsed.ASTView()); err != nil || output != "page" || formats != 1 { t.Fatal("render without a checker", output, err) }
	checker, err := check.NewChecker[amber](d, testChecker[amber]{id: d.ID()})
	if err != nil { t.Fatal(err) }
	if err := e.RegisterChecker(checker); err != nil { t.Fatal(err) }
	if err := e.Check(parsed.ASTView()); err != nil { t.Fatal(err) }
	for _, duplicate := range []struct { got, want error }{
		{e.RegisterParser(parser), engine.ErrDuplicateParser},
		{e.RegisterChecker(checker), engine.ErrDuplicateChecker},
		{e.RegisterFormatter(formatter), engine.ErrDuplicateFormatter},
		{e.RegisterPrinter(printer), engine.ErrDuplicatePrinter},
	} {
		if !errors.Is(duplicate.got, duplicate.want) { t.Fatal(duplicate) }
	}
	if _, err := e.Render(parsed.ASTView()); err != nil { t.Fatal("failed duplicate damaged registration", err) }

	// A separate engine can render with neither a parser nor a checker.
	renderOnly := engine.New()
	register(t, renderOnly, d)
	if err := renderOnly.RegisterPrinter(printer); err != nil { t.Fatal(err) }
	if output, err := renderOnly.Print(parsed); err != nil || output != "page" { t.Fatal(output, err) }
	if err := renderOnly.RegisterFormatter(formatter); err != nil { t.Fatal(err) }
	if _, err := renderOnly.Render(parsed.ASTView()); err != nil { t.Fatal(err) }
	if _, err := renderOnly.Parse(d.ID(), parse.Input{}); !errors.Is(err, engine.ErrMissingCapability) { t.Fatal(err) }
	if err := renderOnly.Check(parsed.ASTView()); !errors.Is(err, engine.ErrMissingCapability) { t.Fatal(err) }

	checkOnly := engine.New()
	register(t, checkOnly, d)
	if err := checkOnly.RegisterChecker(checker); err != nil { t.Fatal(err) }
	if err := checkOnly.Check(parsed.ASTView()); err != nil { t.Fatal(err) }
	if _, err := checkOnly.Render(parsed.ASTView()); !errors.Is(err, engine.ErrMissingCapability) { t.Fatal(err) }
}

// Bind every stage against a supplied profile, using independent implementations.
func stageRegistrations[L any](t *testing.T, d representation.Language) map[string]func(*engine.Engine) error {
	t.Helper()
	identity := stageIdentity{d.ID()}
	parser, err := parse.NewParser[L](d, parserOnly[L]{stageIdentity: identity})
	if err != nil { t.Fatal(err) }
	checker, err := check.NewChecker[L](d, testChecker[L]{id: d.ID()})
	if err != nil { t.Fatal(err) }
	formatter, err := format.NewFormatter[L](d, formatterOnly[L]{stageIdentity: identity})
	if err != nil { t.Fatal(err) }
	printer, err := print.NewPrinter[L](d, printerOnly[L]{stageIdentity: identity})
	if err != nil { t.Fatal(err) }
	return map[string]func(*engine.Engine) error{
		"parse": func(e *engine.Engine) error { return e.RegisterParser(parser) },
		"check": func(e *engine.Engine) error { return e.RegisterChecker(checker) },
		"format": func(e *engine.Engine) error { return e.RegisterFormatter(formatter) },
		"print": func(e *engine.Engine) error { return e.RegisterPrinter(printer) },
	}
}

func TestStageRegistrationIdentityAndProfiles(t *testing.T) {
	d := definition[amber](t, "stages")
	otherProfile, err := representation.NewLanguage[amber](stageIdentity{d.ID()}, "other")
	if err != nil { t.Fatal(err) }
	otherMarker := definition[blue](t, d.ID())
	for name, install := range stageRegistrations[amber](t, d) {
		t.Run(name, func(t *testing.T) {
			empty := engine.New()
			if err := install(empty); !errors.Is(err, engine.ErrUnknownLanguage) { t.Fatal(err) }
			for _, wrong := range []representation.Language{otherProfile, otherMarker} {
				e := engine.New()
				register(t, e, wrong)
				if err := install(e); !errors.Is(err, representation.ErrLanguageMismatch) { t.Fatal(err) }
			}
			register(t, empty, d)
			if err := install(empty); err != nil { t.Fatal("failed registration polluted registry", err) }
		})
	}
	// Actual Go markers are checked at adapter construction, not only by ID.
	if _, err := parse.NewParser[blue](d, parserOnly[blue]{stageIdentity: stageIdentity{d.ID()}}); !errors.Is(err, representation.ErrTypeMismatch) { t.Fatal(err) }
	if _, err := check.NewChecker[blue](d, testChecker[blue]{id: d.ID()}); !errors.Is(err, representation.ErrTypeMismatch) { t.Fatal(err) }
	if _, err := format.NewFormatter[blue](d, formatterOnly[blue]{stageIdentity: stageIdentity{d.ID()}}); !errors.Is(err, representation.ErrTypeMismatch) { t.Fatal(err) }
	if _, err := print.NewPrinter[blue](d, printerOnly[blue]{stageIdentity: stageIdentity{d.ID()}}); !errors.Is(err, representation.ErrTypeMismatch) { t.Fatal(err) }
	if _, err := parse.NewParser[amber](d, parserOnly[amber]{stageIdentity: stageIdentity{"wrong"}}); !errors.Is(err, parse.ErrInvalidParser) { t.Fatal(err) }
	if _, err := check.NewChecker[amber](d, testChecker[amber]{id: "wrong"}); !errors.Is(err, check.ErrInvalidChecker) { t.Fatal(err) }
	if _, err := format.NewFormatter[amber](d, formatterOnly[amber]{stageIdentity: stageIdentity{"wrong"}}); !errors.Is(err, format.ErrInvalidFormatter) { t.Fatal(err) }
	if _, err := print.NewPrinter[amber](d, printerOnly[amber]{stageIdentity: stageIdentity{"wrong"}}); !errors.Is(err, print.ErrInvalidPrinter) { t.Fatal(err) }
	var nilParser *parserOnly[amber]
	var nilChecker *testChecker[amber]
	var nilFormatter *formatterOnly[amber]
	var nilPrinter *printerOnly[amber]
	if _, err := parse.NewParser[amber](d, nilParser); !errors.Is(err, parse.ErrInvalidParser) { t.Fatal(err) }
	if _, err := check.NewChecker[amber](d, nilChecker); !errors.Is(err, check.ErrInvalidChecker) { t.Fatal(err) }
	if _, err := format.NewFormatter[amber](d, nilFormatter); !errors.Is(err, format.ErrInvalidFormatter) { t.Fatal(err) }
	if _, err := print.NewPrinter[amber](d, nilPrinter); !errors.Is(err, print.ErrInvalidPrinter) { t.Fatal(err) }
	e := engine.New()
	for _, err := range []error{e.RegisterParser(parse.Parser{}), e.RegisterChecker(check.Checker{}), e.RegisterFormatter(format.Formatter{}), e.RegisterPrinter(print.Printer{})} {
		if err == nil { t.Fatal("zero stage accepted") }
	}
}

func TestIndependentStageInputSnapshots(t *testing.T) {
	d := definition[amber](t, "owned")
	input := packed[amber](t, d)
	leaf := "page"
	formatter, err := format.NewFormatter[amber](d, formatterOnly[amber]{stageIdentity{d.ID()}, func(v ast.View[amber]) (cst.CST[amber], error) {
		v.ASTView().Nodes[1].Attributes[0].Text = "formatter mutation"
		return exampleConcrete[amber](&leaf), nil
	}})
	if err != nil { t.Fatal(err) }
	concrete, err := formatter.Format(input)
	if err != nil { t.Fatal(err) }
	printer, err := print.NewPrinter[amber](d, printerOnly[amber]{stageIdentity{d.ID()}, func(c cst.CST[amber]) (string, error) {
		*c.Nodes[0].Text = "printer mutation"
		c.Abstract.Nodes[1].Attributes[0].Text = "printer mutation"
		return "printed", nil
	}})
	if err != nil { t.Fatal(err) }
	if _, err := printer.Print(concrete); err != nil { t.Fatal(err) }
	abstract, _ := representation.Unpack[amber](d, input)
	retained, _ := representation.UnpackCST[amber](d, concrete)
	if abstract.Nodes[1].Attributes[0].Text != "Page" || retained.Abstract.Nodes[1].Attributes[0].Text != "Page" || *retained.Nodes[0].Text != "page" { t.Fatal("stage mutated envelope storage") }
}

func TestEnginePreflightsEveryRouteChecker(t *testing.T) {
	a, b, c := definition[amber](t, "a"), definition[blue](t, "b"), definition[coral](t, "c")
	e := engine.New()
	register(t, e, a, b, c)
	checks, edges := 0, 0
	first, err := check.NewChecker[amber](a, testChecker[amber]{id: a.ID(), visit: func(ast.AST[amber]) error { checks++; return nil }})
	if err != nil { t.Fatal(err) }
	last, err := check.NewChecker[coral](c, testChecker[coral]{id: c.ID()})
	if err != nil { t.Fatal(err) }
	if err := e.RegisterChecker(first); err != nil { t.Fatal(err) }
	if err := e.RegisterChecker(last); err != nil { t.Fatal(err) }
	wrong := exampleTree[blue]()
	wrong.Nodes[0].Kind = "not-a-document"
	invalid, err := representation.Pack(b, wrong)
	if err != nil { t.Fatal(err) }
	connect(t, e, &directTransformer{from: a.ID(), to: b.ID(), run: func(representation.AST) (representation.AST, error) { edges++; return invalid, nil }}, bind[blue, coral](t, b, c))
	input := packed[amber](t, a)
	plan, err := e.Plan(a.ID(), c.ID())
	if err != nil { t.Fatal("planning should require connectivity only", err) }
	if _, err := e.Execute(input, plan); !errors.Is(err, engine.ErrMissingCapability) || checks != 0 || edges != 0 { t.Fatal("missing intermediate checker discovered after callbacks", err) }
	middle, err := check.NewChecker[blue](b, testChecker[blue]{id: b.ID()})
	if err != nil { t.Fatal(err) }
	if err := e.RegisterChecker(middle); err != nil { t.Fatal(err) }
	var step *engine.StepError
	if output, err := e.Execute(input, plan); !errors.As(err, &step) || step.Index != 0 || output.LanguageID() != "" || checks != 1 || edges != 1 { t.Fatal("intermediate syntax was not checked", err) }
	if _, err := e.TransformTo(invalid, b.ID()); err == nil { t.Fatal("zero-hop syntax was not checked") }
}

func TestEngineConcurrentIndependentStages(t *testing.T) {
	d := definition[amber](t, "concurrent-stages")
	e := engine.New()
	register(t, e, d)
	input := packed[amber](t, d)
	leaf := "page"
	raw := exampleConcrete[amber](&leaf)
	concrete, err := representation.PackCST(d, raw)
	if err != nil { t.Fatal(err) }
	extra := definition[echo](t, "callback-registration")
	var once sync.Once
	reenter := func() { _ = e.Languages(); once.Do(func() { if err := e.RegisterLanguage(extra); err != nil { t.Error(err) } }) }
	identity := stageIdentity{d.ID()}
	parser, err := parse.NewParser[amber](d, parserOnly[amber]{identity, func(parse.Input) (cst.CST[amber], error) { reenter(); return raw, nil }})
	if err != nil { t.Fatal(err) }
	checker, err := check.NewChecker[amber](d, testChecker[amber]{id: d.ID(), visit: func(ast.AST[amber]) error { reenter(); return nil }})
	if err != nil { t.Fatal(err) }
	formatter, err := format.NewFormatter[amber](d, formatterOnly[amber]{identity, func(ast.View[amber]) (cst.CST[amber], error) { reenter(); return raw, nil }})
	if err != nil { t.Fatal(err) }
	printer, err := print.NewPrinter[amber](d, printerOnly[amber]{identity, func(cst.CST[amber]) (string, error) { reenter(); return leaf, nil }})
	if err != nil { t.Fatal(err) }
	var group sync.WaitGroup
	for _, run := range []func() error{
		func() error { if err := e.RegisterParser(parser); err != nil { return err }; _, err := e.Parse(d.ID(), parse.Input{}); return err },
		func() error { if err := e.RegisterChecker(checker); err != nil { return err }; return e.Check(input) },
		func() error { if err := e.RegisterFormatter(formatter); err != nil { return err }; _, err := e.Format(input); return err },
		func() error { if err := e.RegisterPrinter(printer); err != nil { return err }; _, err := e.Print(concrete); return err },
	} {
		group.Add(1)
		go func() { defer group.Done(); if err := run(); err != nil { t.Error(err) } }()
	}
	done := make(chan struct{})
	go func() { group.Wait(); close(done) }()
	select { case <-done: case <-time.After(5*time.Second): t.Fatal("stage callback held a registry lock") }
	if output, err := e.Render(input); err != nil || output != "page" { t.Fatal("concurrent registrations lost a stage", err) }
}
