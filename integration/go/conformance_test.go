package integration_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/internal/tree"
	cppprovider "github.com/Bitspark/bittype/provider/cpp/go"
	goprovider "github.com/Bitspark/bittype/provider/golang/go"
	rustprovider "github.com/Bitspark/bittype/provider/rust/go"
	tsprovider "github.com/Bitspark/bittype/provider/typescript/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

type nativeConformanceCase struct {
	directory string
	conformance func(*testing.T)
	roundTrips func(*testing.T)
	fuzz func(*testing.T, string)
}

// This is qualification data, never a production registration/default list.
// The directory-discovery test requires coverage for every implemented language.
func nativeConformanceCases() []nativeConformanceCase {
	return []nativeConformanceCase{
		nativeCase(tsprovider.Provider{}, "typescript", "ts", "type Page<T> = T[];"),
		nativeCase(goprovider.Provider{}, "golang", "go", "package p; type Page[T any] []T"),
		nativeCase(rustprovider.Provider{}, "rust", "rs", "type Page<T> = Vec<T>;"),
		nativeCase(cppprovider.Provider{}, "cpp", "hpp", "template<typename T> struct Page { T value; };"),
	}
}

func nativeCase[L any](provider text.Provider[L], directory, extension, source string) nativeConformanceCase {
	return nativeConformanceCase{
		directory: directory,
		conformance: func(t *testing.T) { providerConformance(t, provider, source) },
		roundTrips: func(t *testing.T) { exercise(t, provider, directory, extension) },
		fuzz: func(t *testing.T, input string) { fuzzOne(t, provider, input) },
	}
}

func TestNativeProviderConformance(t *testing.T) {
	for _, fixture := range nativeConformanceCases() { t.Run(fixture.directory, fixture.conformance) }
}

func providerConformance[L any](t *testing.T, provider text.Provider[L], source string) {
	t.Helper()
	id, name := provider.ID(), provider.Name()
	if id == "" || name == "" { t.Fatal("empty provider identity") }
	input := text.Input{SourceID: "conformance", TreeID: "conformance-tree", Text: source}
	parsed, err := provider.Parse(input)
	if err != nil { t.Fatal(err) }
	for _, err := range []error{provider.Check(parsed), provider.Check(parsed.ASTView())} {
		if err != nil { t.Fatalf("AST and CST views must satisfy the same checker: %v", err) }
	}
	before := tree.Clone(parsed.Abstract)
	formatted, err := provider.Format(parsed)
	if err != nil { t.Fatal(err) }
	if !reflect.DeepEqual(parsed.Abstract, before) || !reflect.DeepEqual(formatted.Abstract, before) {
		t.Fatal("formatter changed input or lost AST metadata")
	}
	printed, err := provider.Print(formatted)
	if err != nil { t.Fatal(err) }
	parsed.Abstract.Nodes[parsed.Abstract.Root-1].Kind = "changed-by-caller"
	if !reflect.DeepEqual(formatted.Abstract, before) { t.Fatal("formatter borrowed mutable input storage") }
	if again, err := provider.Print(formatted); err != nil || again != printed { t.Fatal("output changed with caller storage", err) }

	input.Text = source + " @"
	partial, err := provider.Parse(input)
	var diagnostic *text.Diagnostic
	if err == nil || !errors.As(err, &diagnostic) || !reflect.DeepEqual(partial, cst.CST[L]{}) {
		t.Fatalf("whole-input rejection must return a located error and no partial result: %v", err)
	}
	invalid := tree.Clone(before)
	invalid.Root = 0
	if err := provider.Check(invalid); err == nil { t.Fatal("checker accepted an invalid carrier") }
	partial, err = provider.Format(invalid)
	if err == nil || !reflect.DeepEqual(partial, cst.CST[L]{}) { t.Fatal("formatter returned partial output on invalid input") }
	formatted.Root = 0
	if output, err := provider.Print(formatted); err == nil || output != "" { t.Fatal("printer returned partial text on invalid input") }
	if provider.ID() != id || provider.Name() != name { t.Fatal("provider identity changed during use") }
}
