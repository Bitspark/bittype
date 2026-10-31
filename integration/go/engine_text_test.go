package integration_test

import (
	"errors"
	"github.com/Bitspark/bittype/trees/representation/go"
	"testing"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/engine/go"
	cppprovider "github.com/Bitspark/bittype/provider/cpp/go"
	goprovider "github.com/Bitspark/bittype/provider/golang/go"
	rustprovider "github.com/Bitspark/bittype/provider/rust/go"
	tsprovider "github.com/Bitspark/bittype/provider/typescript/go"
	cpp "github.com/Bitspark/bittype/syntax/cpp/go"
	golang "github.com/Bitspark/bittype/syntax/golang/go"
	rust "github.com/Bitspark/bittype/syntax/rust/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

func nativeRegistration[L any](t *testing.T, e *engine.Engine, p text.Provider[L], source string) engine.Provider {
	t.Helper()
	registered, err := engine.NewProvider[L](p, "native-types-v1")
	if err != nil {
		t.Fatal(err)
	}
	if registered.ID() != p.ID() || registered.Name() != p.Name() {
		t.Fatal("provider identity changed")
	}
	if err := e.RegisterProvider(registered); err != nil {
		t.Fatal(err)
	}
	input := text.Input{SourceID: "engine-native", TreeID: "engine-native", Text: source}
	parsed, err := e.Parse(p.ID(), input)
	if err != nil {
		t.Fatal(err)
	}
	original, err := e.Print(parsed)
	if err != nil || original != source {
		t.Fatalf("retention: %v", err)
	}
	output, err := e.Render(parsed.ASTView())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.Parse(text.Input{SourceID: "rendered", TreeID: "rendered", Text: output}); err != nil {
		t.Fatal(err)
	}
	if err := e.Check(parsed.ASTView()); err != nil {
		t.Fatal(err)
	}
	plain, err := representation.UnpackCST[L](registered.Language(), parsed)
	if err != nil {
		t.Fatal(err)
	}
	for i := range plain.Nodes {
		if plain.Nodes[i].Text != nil {
			*plain.Nodes[i].Text = "changed"
			break
		}
	}
	plain.Abstract.Nodes[plain.Abstract.Root-1].Kind = "changed"
	after, err := e.Print(parsed)
	if err != nil || after != source {
		t.Fatal("concrete storage escaped")
	}
	return registered
}

func TestEngineUsesFourExistingTextProviders(t *testing.T) {
	e := engine.New()
	if len(e.Languages()) != 0 {
		t.Fatal("hidden language registration")
	}
	p := nativeRegistration[ts.Language](t, e, tsprovider.Provider{}, "// retain\r\ntype Page<T> = T[];\r\n")
	nativeRegistration[golang.Language](t, e, goprovider.Provider{}, "package p\n\ntype Page[T any] []T\n")
	nativeRegistration[rust.Language](t, e, rustprovider.Provider{}, "type Page<T> = Vec<T>;\n")
	nativeRegistration[cpp.Language](t, e, cppprovider.Provider{}, "template<typename T> struct Page { T value; };\n")
	if err := e.RegisterLanguage(p.Language()); !errors.Is(err, engine.ErrDuplicateLanguage) {
		t.Fatal(err)
	}
	if _, err := e.Plan("typescript", "rust"); !errors.Is(err, engine.ErrNoPath) {
		t.Fatal("invented native transformer", err)
	}
	if _, err := e.Parse("missing", text.Input{}); !errors.Is(err, engine.ErrUnknownLanguage) {
		t.Fatal(err)
	}
}

type fakeText[L any] struct {
	testChecker[L]
	parse  func(text.Input) (cst.CST[L], error)
	format func(ast.View[L]) (cst.CST[L], error)
	print  func(cst.CST[L]) (string, error)
}

func (p fakeText[L]) Parse(input text.Input) (cst.CST[L], error)   { return p.parse(input) }
func (p fakeText[L]) Format(input ast.View[L]) (cst.CST[L], error) { return p.format(input) }
func (p fakeText[L]) Print(input cst.CST[L]) (string, error)       { return p.print(input) }
func exampleConcrete[L any](leaf *string) cst.CST[L] {
	return cst.CST[L]{Abstract: exampleTree[L](), Root: 1, Nodes: []cst.Node[L]{{Kind: "leaf", Text: leaf, Syntax: []ast.NodeID{1}, Origin: ast.Origin{Spans: []ast.Span{{Source: "text", Start: 0, End: 4}}, Nodes: []ast.NodeRef{{Tree: "earlier", Node: 1}}}}}}
}
func TestEngineConcreteSnapshotsAndPartialErrors(t *testing.T) {
	leaf := "page"
	raw := exampleConcrete[amber](&leaf)
	failure := errors.New("provider failed")
	var parseErr, formatErr, printErr error
	p := fakeText[amber]{
		testChecker: testChecker[amber]{id: "external-language"},
		parse:       func(text.Input) (cst.CST[amber], error) { return raw, parseErr },
		format:      func(ast.View[amber]) (cst.CST[amber], error) { return raw, formatErr },
		print:       func(c cst.CST[amber]) (string, error) { return *c.Nodes[0].Text, printErr },
	}
	adapted, err := engine.NewProvider[amber](p, "v1")
	if err != nil {
		t.Fatal(err)
	}
	e := engine.New()
	if err := e.RegisterProvider(adapted); err != nil {
		t.Fatal(err)
	}
	parsed, err := e.Parse(adapted.ID(), text.Input{})
	if err != nil {
		t.Fatal(err)
	}
	leaf = "mutated"
	raw.Nodes[0].Origin.Spans[0].Source = "mutated"
	raw.Nodes[0].Origin.Nodes[0].Tree = "mutated"
	raw.Nodes[0].Syntax[0] = 2
	output, err := e.Print(parsed)
	if err != nil || output != "page" {
		t.Fatal("provider retained concrete output storage")
	}
	unpacked, err := representation.UnpackCST[amber](adapted.Language(), parsed)
	if err != nil || unpacked.Nodes[0].Origin.Spans[0].Source != "text" || unpacked.Nodes[0].Origin.Nodes[0].Tree != "earlier" || unpacked.Nodes[0].Syntax[0] != 1 {
		t.Fatal("concrete metadata aliases provider")
	}
	parseErr = failure
	if c, err := e.Parse(adapted.ID(), text.Input{}); !errors.Is(err, failure) || c.LanguageID() != "" {
		t.Fatal("partial parse output", err)
	}
	formatErr = failure
	if c, err := e.Format(parsed.ASTView()); !errors.Is(err, failure) || c.LanguageID() != "" {
		t.Fatal("partial format output", err)
	}
	printErr = failure
	if output, err := e.Print(parsed); !errors.Is(err, failure) || output != "" {
		t.Fatal("partial print output", err)
	}
	parseErr = nil
	raw.Root = 99
	if c, err := e.Parse(adapted.ID(), text.Input{}); err == nil || c.LanguageID() != "" {
		t.Fatal("invalid concrete parse accepted")
	}
	if _, err := representation.UnpackCST[blue](adapted.Language(), parsed); !errors.Is(err, representation.ErrTypeMismatch) {
		t.Fatal(err)
	}
}

func TestEngineASTOnlyCapabilitiesAndProfiles(t *testing.T) {
	a := definition[amber](t, "amber")
	e := engine.New()
	register(t, e, a)
	if _, err := e.Parse("amber", text.Input{}); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal(err)
	}
	if _, err := e.Render(packed[amber](t, a)); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal(err)
	}
	if err := e.RegisterProvider(engine.Provider{}); err == nil {
		t.Fatal("zero provider accepted")
	}
	var nilProvider *fakeText[amber]
	if _, err := engine.NewProvider[amber](nilProvider, "v1"); !errors.Is(err, representation.ErrInvalidLanguage) {
		t.Fatal(err)
	}
	p, err := engine.NewProvider[ts.Language](tsprovider.Provider{}, "v1")
	if err != nil {
		t.Fatal(err)
	}
	other, err := engine.NewProvider[ts.Language](tsprovider.Provider{}, "v2")
	if err != nil {
		t.Fatal(err)
	}
	if err := e.RegisterProvider(p); err != nil {
		t.Fatal(err)
	}
	typedCST, err := (tsprovider.Provider{}).Parse(text.Input{SourceID: "profile", TreeID: "profile", Text: "type X=string;"})
	if err != nil {
		t.Fatal(err)
	}
	foreign, err := representation.PackCST(other.Language(), typedCST)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.Print(foreign); !errors.Is(err, representation.ErrLanguageMismatch) {
		t.Fatal(err)
	}
	if _, err := e.Format(foreign.ASTView()); !errors.Is(err, representation.ErrLanguageMismatch) {
		t.Fatal(err)
	}
	if _, err := e.Print(representation.CST{}); !errors.Is(err, engine.ErrUnknownLanguage) {
		t.Fatal(err)
	}
}

var _ language.LanguageProvider = engine.Provider{}
