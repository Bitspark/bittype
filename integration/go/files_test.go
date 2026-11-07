package integration_test

import (
	analyze "github.com/Bitspark/bittype/analyze/files/horos/provider/go"
	checkhoros "github.com/Bitspark/bittype/check/syntax/horos/provider/go"
	checkerapi "github.com/Bitspark/bittype/contracts/check/go"
	printapi "github.com/Bitspark/bittype/contracts/print/go"
	checkfiles "github.com/Bitspark/bittype/contracts/resolution/go"
	formathoros "github.com/Bitspark/bittype/format/syntax/horos/provider/go"
	textfixture "github.com/Bitspark/bittype/integration/go/testdata/textfixture"
	filesparse "github.com/Bitspark/bittype/parse/files/horos/provider/go"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	filesprint "github.com/Bitspark/bittype/print/files/shared/provider/go"
	syntaxprint "github.com/Bitspark/bittype/print/syntax/shared/provider/go"
	resolve "github.com/Bitspark/bittype/resolve/files/horos/provider/go"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/engine/go"
	"github.com/Bitspark/bittype/files/go"
	filesformat "github.com/Bitspark/bittype/format/files/horos/provider/go"

	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/source/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"reflect"
	"testing"
)

func fileStages(t *testing.T, def representation.Language) (checkerapi.Checker, checkfiles.Resolver, analysis.Analyzer) {
	t.Helper()
	checker, err := checkerapi.NewChecker[s.Language](def, checkhoros.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	resolver, err := checkfiles.NewResolver[s.Language](def, resolve.Contract(), resolve.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[s.Language](def, analyze.Contract(), analyze.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	return checker, resolver, analyzer
}
func filePagination(t *testing.T) (parse.FilesInput, string, string) {
	t.Helper()
	ctx := context.Background()
	p := horosText()
	page := "// shared pagination\r\nmodule example.pagination;\r\ntype Page<T> = { items: list<T>; next?: text; };\r\n"
	pinOf := func(text string) string {
		concrete, err := p.Parse(ctx, parse.Input{SourceID: "fixture", TreeID: "fixture", Text: text})
		if err != nil {
			t.Fatal(err)
		}
		for _, node := range concrete.Abstract.Nodes {
			if node.Kind == s.KindModule {
				for _, a := range node.Attributes {
					if a.Name == s.AttributePin {
						return a.Text
					}
				}
			}
		}
		t.Fatal("missing module pin")
		return ""
	}
	pagePin := pinOf(page)
	people := fmt.Sprintf("module example.people;\nimport pages from example.pagination pin \"sha256:%s\";\ntype Person = { name: text; };\ntype People = pages.Page<Person>;\n", pagePin)
	return parse.FilesInput{Entry: "people", Files: files.Files{Entries: []files.File{{ID: "page", Path: "lib/page.horos", Text: page}, {ID: "people", Path: "people.horos", Text: people}}}}, pagePin, pinOf(people)
}

// bittype:case target=github.com/Bitspark/bittype/parse/files/horos/provider/go#Provider.Parse feature=multi-file-pagination case=receipt-preserving-round-trips expects=module-context-preserved
// bittype:case target=github.com/Bitspark/bittype/format/files/horos/provider/go#Provider.Format feature=multi-file-pagination case=receipt-preserving-round-trips expects=fresh-output-and-receipt
// bittype:case target=github.com/Bitspark/bittype/print/files/shared/provider/go#Provider.Print feature=multi-file-pagination case=receipt-preserving-round-trips expects=file-content-and-placement-preserved language=github.com/Bitspark/bittype/syntax/horos/go#Language
func TestFilesPaginationAndReceiptPreservingRoundTrips(t *testing.T) {
	ctx := context.Background()
	input, pagePin, peoplePin := filePagination(t)
	definition, err := representation.NewLanguage[s.Language](s.Language{}, s.Profile)
	if err != nil {
		t.Fatal(err)
	}
	c, r, a := fileStages(t, definition)
	syntaxPrinter := textfixture.Printing[s.Language](s.Language{})
	parser, err := filesparse.New(definition, parsehoros.Provider{}, syntaxPrinter, c, r, a, semantic.Options{}, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	formatter, err := filesformat.New(definition, parsehoros.Provider{}, formathoros.Provider{}, syntaxPrinter, c, r, a, semantic.Options{}, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	filesPrinter, err := filesprint.New[s.Language](s.Language{}, syntaxPrinter)
	if err != nil {
		t.Fatal(err)
	}
	printer, err := printapi.NewFilesPrinter[s.Language](definition, filesPrinter)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parser.Parse(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.Units()) != 2 || len(parsed.Snapshot().Inputs()) != 2 || len(parsed.Groups()) != 2 || len(parsed.Relations()) != 1 {
		t.Fatal("lost multi-file context")
	}
	printed, err := printer.Print(ctx, parsed)
	if err != nil || !reflect.DeepEqual(printed, input.Files) {
		t.Fatal("parse/print changed input files", err)
	}
	graph := parsed.SemanticsView()
	if !graph.Valid() || !parsed.Evidence().Snapshot().Same(parsed.Snapshot()) {
		t.Fatal("discarded receipt")
	}
	var pageEntity, peopleEntity semantics.EntityID
	for i, e := range graph.Data().Entities {
		if e.Name == "Page" {
			pageEntity = semantics.EntityID(i + 1)
		}
		if e.Name == "People" {
			peopleEntity = semantics.EntityID(i + 1)
		}
	}
	people, _ := graph.Entity(peopleEntity)
	application, _ := graph.Term(people.Body)
	callee, _ := graph.Term(application.Callee)
	if callee.Target != pageEntity {
		t.Fatal("cross-file generic target changed")
	}
	if bindings, ok := graph.DirectParameterBindings(people.Body); !ok || len(bindings) != 1 {
		t.Fatal("lost generic substitution binding")
	}
	copy := parsed.Units()
	copy[0].Syntax.Abstract.Nodes[0].Kind = "mutated"
	if parsed.Units()[0].Syntax.Abstract.Nodes[0].Kind == "mutated" {
		t.Fatal("Source exposes mutable storage")
	}
	placement, _ := json.Marshal([]map[string]any{{"ID": "bundle", "Path": "generated/types.horos", "Modules": []string{peoplePin, pagePin}}})
	formatted, err := formatter.Format(ctx, format.FilesInput[s.Language]{Abstract: graph, Snapshot: parsed.Snapshot(), Evidence: parsed.Evidence(), Configuration: []semantic.Setting{{Key: filesformat.OrganizationSetting(), Value: placement}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(formatted.Units()) != 1 || formatted.Snapshot().Same(parsed.Snapshot()) || !formatted.Evidence().Snapshot().Same(formatted.Snapshot()) {
		t.Fatal("format did not establish fresh output")
	}
	if _, err := formatter.Format(ctx, format.FilesInput[s.Language]{Abstract: formatted.SemanticsView(), Snapshot: formatted.Snapshot(), Evidence: formatted.Evidence()}); err != nil {
		t.Fatal("format result lost its reusable receipt", err)
	}
	if _, err := formatter.Format(ctx, format.FilesInput[s.Language]{Abstract: graph, Snapshot: formatted.Snapshot(), Evidence: parsed.Evidence()}); err == nil {
		t.Fatal("stale receipt accepted")
	}
	filesOut, err := printer.Print(ctx, formatted)
	if err != nil {
		t.Fatal(err)
	}
	if filesOut.Entries[0].Path != "generated/types.horos" {
		t.Fatal("printer selected another path")
	}
	if _, err := parser.Parse(ctx, parse.FilesInput{Files: filesOut, Entry: "bundle"}); err != nil {
		t.Fatal("formatted files cannot be parsed", err)
	}
}

type replayFiles struct {
	s.Language
	saved source.Source[s.Language]
}

func (r replayFiles) Parse(context.Context, parse.FilesRequest[s.Language]) (source.Source[s.Language], error) {
	return r.saved, nil
}
func TestFilesCapabilitiesAreIndependentAndRejectStalePublications(t *testing.T) {
	ctx := context.Background()
	input, _, _ := filePagination(t)
	definition, err := representation.NewLanguage[s.Language](s.Language{}, s.Profile)
	if err != nil {
		t.Fatal(err)
	}
	c, r, a := fileStages(t, definition)
	syntaxPrinter := textfixture.Printing[s.Language](s.Language{})
	parser, err := filesparse.New(definition, parsehoros.Provider{}, syntaxPrinter, c, r, a, semantic.Options{}, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parser.Parse(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	e := engine.New()
	registerNativeStages(t, e, parser.Language(), horosText())
	if _, err := e.ParseFiles(ctx, parser.ID(), input); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal("Files fell back to Syntax", err)
	}
	if err := e.RegisterFilesParser(parser); err != nil {
		t.Fatal(err)
	}
	out, err := e.ParseFiles(ctx, parser.ID(), input)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.PrintFiles(ctx, out); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal("parser installed printer")
	}
	if err := e.RegisterFilesParser(parser); !errors.Is(err, engine.ErrDuplicateFilesParser) {
		t.Fatal(err)
	}
	result := parsed.Evidence()
	contract := parse.FilesContract[s.Language]{Evidence: semantic.EvidenceRequirement{Operation: analysis.Kind(), Profile: result.Contract().Profile, Reports: result.Contract().SemanticsRequirements, Semantics: true}, Print: syntaxprint.Print[s.Language], Check: func(context.Context, parse.FilesInput, source.Publication[s.Language]) error { return nil }}
	replay, err := parse.NewFilesParser(parser.Language(), contract, replayFiles{saved: parsed})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := replay.Parse(ctx, input); !errors.Is(err, source.ErrStale) || got.Valid() {
		t.Fatal("old Source accepted", err)
	}
	input.Files.Entries = input.Files.Entries[1:]
	if got, err := parser.Parse(ctx, input); err == nil || got.Valid() {
		t.Fatal("missing dependency published Source")
	}
}

func TestMultiUnitRouteRetainsEveryOutput(t *testing.T) {
	ctx := context.Background()
	e := engine.New()
	a, b, c := definition[amber](t, "a"), definition[blue](t, "b"), definition[coral](t, "c")
	registerChecked[amber](t, e, a)
	registerChecked[blue](t, e, b)
	registerChecked[coral](t, e, c)
	first, err := transform.NewTransformer[amber, blue](a, b, transform.Contract{}, func(context.Context, transform.Request[amber]) (ast.Set[blue], error) {
		one, two := exampleTree[blue](), exampleTree[blue]()
		one.ID = "one"
		two.ID = "two"
		return ast.Set[blue]{Entry: ast.NodeRef{Tree: "two", Node: two.Root}, Trees: []ast.AST[blue]{one, two}}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	second, err := transform.NewTransformer[blue, coral](b, c, transform.Contract{}, func(_ context.Context, r transform.Request[blue]) (ast.Set[coral], error) {
		if len(r.Inputs()) != 2 || r.Source().ID != "two" {
			t.Fatal("route discarded units or changed entry")
		}
		result := ast.Set[coral]{Entry: r.Entry()}
		for _, envelope := range r.Inputs() {
			typed, err := representation.Unpack[blue](b, envelope)
			if err != nil {
				return ast.Set[coral]{}, err
			}
			out, err := mapExample[blue, coral](typed)
			if err != nil {
				return ast.Set[coral]{}, err
			}
			result.Trees = append(result.Trees, out)
		}
		return result, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	connect(t, e, first, second)
	output, err := e.TransformTo(ctx, captureForTransform[amber](t, a, packed[amber](t, a)), c.ID())
	if err != nil || len(output.Trees) != 2 || output.Entry.Tree != "two" {
		t.Fatal("incomplete multi-unit route", err)
	}
}
