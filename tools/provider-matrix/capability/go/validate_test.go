package capability_test

import (
	"fmt"
	"testing"

	c "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
)

type source struct{}

var root = c.Symbol{Package: "example.test/model", Name: "Forms"}
var at = c.Span{File: "capture/one", Start: 0, End: 10}

func (source) File(ref c.FileRef) (c.File, error) {
	if ref != at.File {
		return c.File{}, fmt.Errorf("foreign file")
	}
	return c.File{Ref: ref, Path: "forms.go", Text: "package model\n"}, nil
}
func (source) Declaration(symbol c.Symbol) (c.Declaration, error) {
	if symbol != root {
		return c.Declaration{}, fmt.Errorf("foreign declaration")
	}
	return c.Declaration{Symbol: root, Span: at}, nil
}
func (source) Facts(c.Span) (c.Facts, error) { return c.Facts{}, nil }

func known[T any](value T) c.Field[T] { return c.Field[T]{State: c.Known, Value: value, At: at} }

func example() c.Inventory {
	return c.Inventory{Extent: c.Decoded, Features: []c.Feature{{
		Key: "item", At: at,
		Roles: known([]c.Role{{Name: "value", At: at, Min: known(0), Max: known(c.UpperBound{Unbounded: true}), Kinds: known([]string{})}}),
		Attributes: known([]c.Attribute{{Name: "name", At: at, Required: known(false), Values: known(c.ValueSet{Unrestricted: true})}}),
		References: known([]string{}),
	}}}
}

func TestKnowledgePreservesKnownEmptyUnboundedAndOpaqueFields(t *testing.T) {
	request := c.DefinitionRequest{Root: root, Source: source{}}
	inventory := example()
	if err := c.Validate(request, inventory); err != nil {
		t.Fatal(err)
	}
	inventory.Extent = c.PartlyDecoded
	inventory.Features[0].Attributes.Value[0].Values = c.Field[c.ValueSet]{State: c.Opaque, At: at, Reason: "custom predicate"}
	if err := c.Validate(request, inventory); err != nil {
		t.Fatal(err)
	}
	if inventory.Features[0].Attributes.Value[0].Required.State != c.Known {
		t.Fatal("opaque value predicate erased known presence")
	}
	inventory.Features[0].Roles.Value[0].Kinds = c.Field[[]string]{State: c.Partial, Value: []string{"item"}, At: at, Reason: "unknown spread"}
	if err := c.Validate(request, inventory); err != nil {
		t.Fatal(err)
	}
}

func TestInvalidRecordsFailInsteadOfBecomingUnknown(t *testing.T) {
	cases := map[string]func(*c.Inventory){
		"missing knowledge": func(i *c.Inventory) { i.Features[0].References.State = "" },
		"unknown value": func(i *c.Inventory) { i.Features[0].References = c.Field[[]string]{State: c.Unknown, Value: []string{"x"}, At: at, Reason: "unknown"} },
		"opaque role": func(i *c.Inventory) { i.Features[0].Roles = c.Field[[]c.Role]{State: c.Opaque, At: at, Reason: "opaque"} },
		"partial scalar": func(i *c.Inventory) { i.Features[0].Roles.Value[0].Min.State = c.Partial; i.Features[0].Roles.Value[0].Min.Reason = "partial" },
		"unbounded with count": func(i *c.Inventory) { i.Features[0].Roles.Value[0].Max.Value.Value = 2 },
		"contradictory bounds": func(i *c.Inventory) { i.Features[0].Roles.Value[0].Min.Value = 3; i.Features[0].Roles.Value[0].Max.Value = c.UpperBound{Value: 2} },
		"unrestricted enum": func(i *c.Inventory) { i.Features[0].Attributes.Value[0].Values.Value.Values = []string{"x"} },
		"duplicate feature": func(i *c.Inventory) { i.Features = append(i.Features, i.Features[0]) },
		"duplicate role": func(i *c.Inventory) { i.Features[0].Roles.Value = append(i.Features[0].Roles.Value, i.Features[0].Roles.Value[0]) },
		"duplicate enum": func(i *c.Inventory) { i.Features[0].Attributes.Value[0].Values.Value = c.ValueSet{Values: []string{"x", "x"}} },
		"foreign span": func(i *c.Inventory) { i.Features[0].At.File = "other-capture/one" },
		"outside file": func(i *c.Inventory) { i.Features[0].At.End = 1000 },
		"decoded with limitation": func(i *c.Inventory) { i.Limitations = []c.Limitation{{At: at, Message: "unknown"}} },
		"partial without limitation": func(i *c.Inventory) { i.Extent = c.PartlyDecoded },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			inventory := example()
			mutate(&inventory)
			if err := c.Validate(c.DefinitionRequest{Root: root, Source: source{}}, inventory); err == nil {
				t.Fatal("invalid extractor result accepted")
			}
		})
	}
}

func TestLimitsAndRootMembership(t *testing.T) {
	for _, request := range []c.DefinitionRequest{
		{Root: c.Symbol{Package: "other", Name: "Forms"}, Source: source{}},
		{Root: root, Source: source{}, Limits: c.Limits{MaxRecords: 1}},
		{Root: root, Source: source{}, Limits: c.Limits{MaxDepth: -1}},
	} {
		if err := c.Validate(request, example()); err == nil {
			t.Fatal("invalid request or exhausted limit accepted")
		}
	}
	for _, text := range []string{"example.test/p#Forms", "example.test/p#Provider.Check"} {
		symbol, err := c.ParseSymbol(text)
		if err != nil || symbol.String() != text {
			t.Fatalf("identity round trip: %s %v", symbol, err)
		}
	}
	for _, text := range []string{"Forms", "#Forms", "p#", "p#P.A.B", "p #Forms", "p#A#B"} {
		if _, err := c.ParseSymbol(text); err == nil {
			t.Fatalf("invalid identity accepted: %q", text)
		}
	}
}
