package lift_test

import (
	"errors"
	"math/big"
	"reflect"
	"testing"

	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/engine/go"
	hp "github.com/Bitspark/bittype/provider/horos/go"
	wp "github.com/Bitspark/bittype/provider/horoswire/go"
	hr "github.com/Bitspark/bittype/resolve/horos/go"
	wr "github.com/Bitspark/bittype/resolve/horoswire/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	w "github.com/Bitspark/bittype/syntax/horoswire/go"
	lift "github.com/Bitspark/bittype/transforms/horos-to-horoswire/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	ground "github.com/Bitspark/horos/core/go/declarations"
	wire "github.com/Bitspark/horos/core/go/wire"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

func tag(prefix, name string, args ...core.Value) core.Value {
	return core.NewTuple(append([]core.Value{core.NewAtom([]byte(prefix + name))}, args...)...)
}
func g(name string, args ...core.Value) core.Value     { return tag("horos/", name, args...) }
func d(name string, args ...core.Value) core.Value     { return tag("horos/decl/", name, args...) }
func mixed(name string, args ...core.Value) core.Value { return tag("horos/wire/", name, args...) }
func num(i int64) core.Value                           { return data.EncodeInt(big.NewInt(i)) }
func declaration(body core.Value, names ...string) core.Value {
	var params []core.Value
	for _, name := range names {
		params = append(params, data.MustEncodeText(name))
	}
	return d("type", data.EncodeList(params), body)
}
func module(t *testing.T, name string, members map[string]core.Value) core.Value {
	t.Helper()
	var entries [][2]core.Value
	for name, value := range members {
		entries = append(entries, [2]core.Value{data.MustEncodeText(name), value})
	}
	table, err := data.EncodeMap(entries)
	if err != nil {
		t.Fatal(err)
	}
	empty, _ := data.EncodeMap(nil)
	return d("module", num(1), data.MustEncodeText(name), empty, table)
}
func fixture(t *testing.T) (h.AST, lift.Options) {
	t.Helper()
	value := module(t, "consumer.pagination", map[string]core.Value{
		"Page":  declaration(g("list-of", d("param", num(0))), "T"),
		"Alias": declaration(d("apply", d("local", data.MustEncodeText("Page")), data.EncodeList([]core.Value{g("scalar", core.NewAtom([]byte("int")))}))),
	})
	tree, err := hr.Resolve(value, nil, "authored")
	if err != nil {
		t.Fatal(err)
	}
	pin, _ := ground.Pin(value)
	for i := range tree.Nodes {
		tree.Nodes[i].Origin.Spans = []ast.Span{{Source: "source", Start: uint64(i), End: uint64(i + 1)}}
	}
	return tree, lift.Options{Roots: []ground.DeclarationID{{Module: pin, Name: "Alias"}}, Entry: pin, TreeID: "derived", Lift: wire.DefaultLiftOptions()}
}

func TestExplicitLiftRetainsSlotsApplicationsAndSourceOrigins(t *testing.T) {
	input, options := fixture(t)
	before := input.Nodes[0].Origin.Spans[0]
	result, err := lift.Derive(input, options)
	if err != nil {
		t.Fatal(err)
	}
	if result.Tree.ID == input.ID || result.Derivation.Roots()[0].Declaration == options.Roots[0] {
		t.Fatal("derived identity collapsed")
	}
	program, err := wr.Model(result.Tree)
	if err != nil {
		t.Fatal(err)
	}
	foundParameter, foundApplication := false, false
	for _, node := range result.Tree.Nodes {
		if node.Kind == w.KindParameter {
			foundParameter = true
		}
		if node.Kind == w.KindApplication {
			foundApplication = true
		}
		for _, origin := range node.Origin.Nodes {
			if origin.Tree != input.ID || origin.Node == 0 {
				t.Fatal("lost direct source address")
			}
		}
	}
	if !foundParameter || !foundApplication {
		t.Fatal("lift specialized away the template")
	}
	for _, mapping := range result.Derivation.Declarations() {
		if mapping.Source.Name == "Page" {
			argument := mixed("wire", mixed("data", g("scalar", core.NewAtom([]byte("int")))))
			actual, err := program.Specialize(mapping.Target.Declaration, []core.Value{argument})
			if err != nil || !core.Equal(actual, mixed("list-of", argument)) {
				t.Fatalf("derived generic: %v", err)
			}
		}
	}
	if len(result.Tree.Nodes[result.Tree.Root-1].Origin.Spans) != 1 {
		t.Fatal("source span not preserved")
	}
	result.Tree.Nodes[result.Tree.Root-1].Origin.Spans[0].Start = 999
	if input.Nodes[0].Origin.Spans[0] != before {
		t.Fatal("output origin aliases caller storage")
	}
}

func TestTransformerIsExplicitAndOwnsConfiguration(t *testing.T) {
	input, options := fixture(t)
	from, err := representation.NewLanguage[h.Language](hp.Provider{}, h.Profile)
	if err != nil {
		t.Fatal(err)
	}
	to, err := representation.NewLanguage[w.Language](wp.Provider{}, w.Profile)
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := representation.Pack[h.Language](from, input)
	if err != nil {
		t.Fatal(err)
	}
	var registry engine.Engine
	if err := registry.RegisterLanguage(from); err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterLanguage(to); err != nil {
		t.Fatal(err)
	}
	if _, err := registry.TransformTo(envelope, to.ID()); err == nil {
		t.Fatal("provider silently registered a lift")
	}
	options.Lift.ModuleNames = map[ground.Digest]string{options.Entry: "selected.output"}
	transformer, err := lift.NewTransformer(from, to, options)
	if err != nil {
		t.Fatal(err)
	}
	options.Roots[0].Name = "Missing"
	options.Lift.ModuleNames[options.Entry] = "mutated.output"
	if err := registry.RegisterTransformer(transformer); err != nil {
		t.Fatal(err)
	}
	result, err := registry.TransformTo(envelope, to.ID())
	if err != nil {
		t.Fatal(err)
	}
	tree, err := representation.Unpack[w.Language](to, result)
	if err != nil {
		t.Fatal(err)
	}
	model, err := wr.Model(tree)
	if err != nil {
		t.Fatal(err)
	}
	entry, _ := model.Module(model.Root())
	if entry.Name() != "selected.output" {
		t.Fatal("caller mutation changed captured selection")
	}
	wrong, _ := representation.NewLanguage[w.Language](wp.Provider{}, "different-profile")
	if _, err := lift.NewTransformer(from, wrong, options); !errors.Is(err, lift.ErrOptions) {
		t.Fatalf("profile mismatch: %v", err)
	}
}

func TestUnsupportedLiftReturnsNoASTOrFallback(t *testing.T) {
	value := module(t, "consumer.unsupported", map[string]core.Value{
		"Other":  declaration(g("not", d("param", num(0))), "T"),
		"Closed": declaration(d("apply", d("local", data.MustEncodeText("Other")), data.EncodeList([]core.Value{g("scalar", core.NewAtom([]byte("int")))}))),
	})
	input, err := hr.Resolve(value, nil, "unsupported")
	if err != nil {
		t.Fatal(err)
	}
	pin, _ := ground.Pin(value)
	options := lift.Options{Roots: []ground.DeclarationID{{Module: pin, Name: "Closed"}}, Entry: pin, TreeID: "result", Lift: wire.DefaultLiftOptions()}
	result, err := lift.Derive(input, options)
	var detail *wire.Error
	if !errors.As(err, &detail) || detail.Code != wire.UnsupportedLift || !reflect.DeepEqual(result, lift.Result{}) {
		t.Fatalf("unsupported application fell back or published output: %+v %v", result, err)
	}
}

func TestSelectionAndBudgetRefusals(t *testing.T) {
	input, options := fixture(t)
	for _, change := range []func(*lift.Options){
		func(o *lift.Options) { o.TreeID = input.ID },
		func(o *lift.Options) { o.Entry = ground.Digest{} },
		func(o *lift.Options) { o.Lift.Limits.MaxNodes = 1 },
	} {
		selected := options
		change(&selected)
		result, err := lift.Derive(input, selected)
		if err == nil || !reflect.DeepEqual(result, lift.Result{}) {
			t.Fatal("refusal returned usable output")
		}
	}
}

func TestIndependentSelectedRootsAreNeverDropped(t *testing.T) {
	other := module(t, "consumer.other", map[string]core.Value{"Other": declaration(g("tuple"))})
	otherPin, _ := ground.Pin(other)
	imports, _ := data.EncodeMap([][2]core.Value{{data.MustEncodeText("other"), d("import", data.MustEncodeText("consumer.other"), core.NewAtom(otherPin[:]))}})
	definitions, _ := data.EncodeMap([][2]core.Value{{data.MustEncodeText("Root"), declaration(g("tuple"))}})
	root := d("module", num(1), data.MustEncodeText("consumer.roots"), imports, definitions)
	rootPin, _ := ground.Pin(root)
	input, err := hr.Resolve(root, []core.Value{other}, "roots")
	if err != nil {
		t.Fatal(err)
	}
	options := lift.Options{Roots: []ground.DeclarationID{{Module: rootPin, Name: "Root"}, {Module: otherPin, Name: "Other"}}, Entry: rootPin, TreeID: "result", Lift: wire.DefaultLiftOptions()}
	result, err := lift.Derive(input, options)
	if !errors.Is(err, lift.ErrOptions) || !reflect.DeepEqual(result, lift.Result{}) {
		t.Fatalf("independent root was silently dropped: %v", err)
	}
}

func TestSelectionRefusalPrecedesCollectionCopies(t *testing.T) {
	_, options := fixture(t)
	from, _ := representation.NewLanguage[h.Language](hp.Provider{}, h.Profile)
	to, _ := representation.NewLanguage[w.Language](wp.Provider{}, w.Profile)
	options.Roots = make([]ground.DeclarationID, 4096)
	options.Lift.Limits.MaxNodes = 1
	allocations := func(count int) float64 {
		selected := options
		selected.Roots = options.Roots[:count]
		return testing.AllocsPerRun(20, func() {
			transformer, err := lift.NewTransformer(from, to, selected)
			var detail *text.Diagnostic
			if transformer != nil || !errors.As(err, &detail) || detail.Code != text.LimitExceeded {
				t.Fatalf("invalid selection captured: %v", err)
			}
		})
	}
	if narrow, wide := allocations(2), allocations(4096); wide > narrow {
		t.Fatalf("rejected selection was copied: %v -> %v", narrow, wide)
	}
}
