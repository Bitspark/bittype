package integration_test

import (
	checkcpp "github.com/Bitspark/bittype/check/syntax/cpp/provider/go"
	checkgolang "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	checkrust "github.com/Bitspark/bittype/check/syntax/rust/provider/go"
	checktypescript "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	formatcpp "github.com/Bitspark/bittype/format/syntax/cpp/provider/go"
	formatgolang "github.com/Bitspark/bittype/format/syntax/golang/provider/go"
	formatrust "github.com/Bitspark/bittype/format/syntax/rust/provider/go"
	formattypescript "github.com/Bitspark/bittype/format/syntax/typescript/provider/go"
	textfixture "github.com/Bitspark/bittype/integration/go/testdata/textfixture"
	parsecpp "github.com/Bitspark/bittype/parse/syntax/cpp/provider/go"
	parsegolang "github.com/Bitspark/bittype/parse/syntax/golang/provider/go"
	parserust "github.com/Bitspark/bittype/parse/syntax/rust/provider/go"
	parsetypescript "github.com/Bitspark/bittype/parse/syntax/typescript/provider/go"

	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/internal/tree"

	"github.com/Bitspark/bittype/trees/cst/go"
)

type nativeConformanceCase struct {
	directory   string
	conformance func(*testing.T)
	roundTrips  func(*testing.T)
	fuzz        func(*testing.T, string)
}

// This is qualification data, never a production registration/default list.
// The directory-discovery test requires coverage for every implemented language.
func nativeConformanceCases() []nativeConformanceCase {
	return []nativeConformanceCase{
		nativeCase(textfixture.New(parsetypescript.Provider{}, checktypescript.Provider{}, formattypescript.Provider{}), "typescript", "ts", "type Page<T> = T[];"),
		nativeCase(textfixture.New(parsegolang.Provider{}, checkgolang.Provider{}, formatgolang.Provider{}), "golang", "go", "package p; type Page[T any] []T"),
		nativeCase(textfixture.New(parserust.Provider{}, checkrust.Provider{}, formatrust.Provider{}), "rust", "rs", "type Page<T> = Vec<T>;"),
		nativeCase(textfixture.New(parsecpp.Provider{}, checkcpp.Provider{}, formatcpp.Provider{}), "cpp", "hpp", "template<typename T> struct Page { T value; };"),
	}
}

func nativeCase[L any](provider nativeProvider[L], directory, extension, source string) nativeConformanceCase {
	return nativeConformanceCase{
		directory:   directory,
		conformance: func(t *testing.T) { providerConformance(t, provider, source) },
		roundTrips:  func(t *testing.T) { exercise(t, provider, directory, extension) },
		fuzz:        func(t *testing.T, input string) { fuzzOne(t, provider, input) },
	}
}

// bittype:case target=github.com/Bitspark/bittype/parse/syntax/typescript/provider/go#Provider.Parse feature=provider-contract case=native-fixture-contract expects=contract-preserved
// bittype:case target=github.com/Bitspark/bittype/check/syntax/typescript/provider/go#Provider.Check feature=provider-contract case=native-fixture-contract expects=contract-preserved
// bittype:case target=github.com/Bitspark/bittype/format/syntax/typescript/provider/go#Provider.Format feature=provider-contract case=native-fixture-contract expects=contract-preserved
// bittype:case target=github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print feature=provider-contract case=native-fixture-contract expects=contract-preserved language=github.com/Bitspark/bittype/syntax/typescript/go#Language
// bittype:case target=github.com/Bitspark/bittype/parse/syntax/golang/provider/go#Provider.Parse feature=provider-contract case=native-fixture-contract expects=contract-preserved
// bittype:case target=github.com/Bitspark/bittype/check/syntax/golang/provider/go#Provider.Check feature=provider-contract case=native-fixture-contract expects=contract-preserved
// bittype:case target=github.com/Bitspark/bittype/format/syntax/golang/provider/go#Provider.Format feature=provider-contract case=native-fixture-contract expects=contract-preserved
// bittype:case target=github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print feature=provider-contract case=native-fixture-contract expects=contract-preserved language=github.com/Bitspark/bittype/syntax/golang/go#Language
// bittype:case target=github.com/Bitspark/bittype/parse/syntax/rust/provider/go#Provider.Parse feature=provider-contract case=native-fixture-contract expects=contract-preserved
// bittype:case target=github.com/Bitspark/bittype/check/syntax/rust/provider/go#Provider.Check feature=provider-contract case=native-fixture-contract expects=contract-preserved
// bittype:case target=github.com/Bitspark/bittype/format/syntax/rust/provider/go#Provider.Format feature=provider-contract case=native-fixture-contract expects=contract-preserved
// bittype:case target=github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print feature=provider-contract case=native-fixture-contract expects=contract-preserved language=github.com/Bitspark/bittype/syntax/rust/go#Language
// bittype:case target=github.com/Bitspark/bittype/parse/syntax/cpp/provider/go#Provider.Parse feature=provider-contract case=native-fixture-contract expects=contract-preserved
// bittype:case target=github.com/Bitspark/bittype/check/syntax/cpp/provider/go#Provider.Check feature=provider-contract case=native-fixture-contract expects=contract-preserved
// bittype:case target=github.com/Bitspark/bittype/format/syntax/cpp/provider/go#Provider.Format feature=provider-contract case=native-fixture-contract expects=contract-preserved
// bittype:case target=github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print feature=provider-contract case=native-fixture-contract expects=contract-preserved language=github.com/Bitspark/bittype/syntax/cpp/go#Language
func TestNativeProviderConformance(t *testing.T) {
	for _, fixture := range nativeConformanceCases() {
		t.Run(fixture.directory, fixture.conformance)
	}
}

func providerConformance[L any](t *testing.T, provider nativeProvider[L], source string) {
	t.Helper()
	id, name := provider.ID(), provider.Name()
	if id == "" || name == "" {
		t.Fatal("empty provider identity")
	}
	input := parse.Input{SourceID: "conformance", TreeID: "conformance-tree", Text: source}
	parsed, err := provider.Parse(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	for _, err := range []error{provider.Check(parsed), provider.Check(parsed.ASTView())} {
		if err != nil {
			t.Fatalf("AST and CST views must satisfy the same checker: %v", err)
		}
	}
	before := tree.Clone(parsed.Abstract)
	formatted, err := provider.Format(context.Background(), parsed)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(parsed.Abstract, before) || !reflect.DeepEqual(formatted.Abstract, before) {
		t.Fatal("formatter changed input or lost AST metadata")
	}
	printed, err := provider.Print(context.Background(), formatted)
	if err != nil {
		t.Fatal(err)
	}
	parsed.Abstract.Nodes[parsed.Abstract.Root-1].Kind = "changed-by-caller"
	if !reflect.DeepEqual(formatted.Abstract, before) {
		t.Fatal("formatter borrowed mutable input storage")
	}
	if again, err := provider.Print(context.Background(), formatted); err != nil || again != printed {
		t.Fatal("output changed with caller storage", err)
	}

	input.Text = source + " @"
	partial, err := provider.Parse(context.Background(), input)
	var diagnostic *text.Diagnostic
	if err == nil || !errors.As(err, &diagnostic) || !reflect.DeepEqual(partial, cst.CST[L]{}) {
		t.Fatalf("whole-input rejection must return a located error and no partial result: %v", err)
	}
	if span := diagnostic.Span; span == nil || span.Source != input.SourceID || span.Start > span.End || span.End > uint64(len(input.Text)) {
		t.Fatalf("parse rejection has no valid source location: %#v", diagnostic)
	}
	invalid := tree.Clone(before)
	invalid.Root = 0
	if err := provider.Check(invalid); err == nil {
		t.Fatal("checker accepted an invalid carrier")
	}
	partial, err = provider.Format(context.Background(), invalid)
	if err == nil || !reflect.DeepEqual(partial, cst.CST[L]{}) {
		t.Fatal("formatter returned partial output on invalid input")
	}
	formatted.Root = 0
	if output, err := provider.Print(context.Background(), formatted); err == nil || output != "" {
		t.Fatal("printer returned partial text on invalid input")
	}
	if provider.ID() != id || provider.Name() != name {
		t.Fatal("provider identity changed during use")
	}
}
