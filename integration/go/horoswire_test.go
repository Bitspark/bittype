package integration_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/horoswire/go"
	check "github.com/Bitspark/bittype/check/horoswire/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/engine/go"
	groundProvider "github.com/Bitspark/bittype/provider/horos/go"
	provider "github.com/Bitspark/bittype/provider/horoswire/go"
	groundResolve "github.com/Bitspark/bittype/resolve/horos/go"
	resolve "github.com/Bitspark/bittype/resolve/horoswire/go"
	"github.com/Bitspark/bittype/semantics/go"
	groundSyntax "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/horoswire/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/wire"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

func wireTag(tag string, args ...core.Value) core.Value {
	return core.NewTuple(append([]core.Value{core.NewAtom([]byte("horos/wire/" + tag))}, args...)...)
}
func wireParam(i int64) core.Value { return wireTag("param", data.EncodeInt(big.NewInt(i))) }
func wireDecl(body core.Value, names ...string) core.Value {
	parameters := make([]core.Value, len(names))
	for i, name := range names {
		parameters[i] = data.MustEncodeText(name)
	}
	return wireTag("type", data.EncodeList(parameters), body)
}
func wireApply(name string, args ...core.Value) core.Value {
	return wireTag("apply", wireTag("local", data.MustEncodeText(name)), data.EncodeList(args))
}
func wireModule(t *testing.T, name string, members map[string]core.Value) core.Value {
	t.Helper()
	var pairs [][2]core.Value
	for key, value := range members {
		pairs = append(pairs, [2]core.Value{data.MustEncodeText(key), value})
	}
	table, err := data.EncodeMap(pairs)
	if err != nil {
		t.Fatal(err)
	}
	empty, _ := data.EncodeMap(nil)
	return wireTag("module", data.EncodeInt(big.NewInt(1)), data.MustEncodeText(name), empty, table)
}
func wireText() core.Value {
	return wireTag("data", core.NewTuple(core.NewAtom([]byte("horos/scalar")), core.NewAtom([]byte("utf8-text"))))
}
func wireFixture(t *testing.T) core.Value {
	return wireModule(t, "example.optional", map[string]core.Value{
		"Page":     wireDecl(wireTag("list-of", wireParam(0)), "T"),
		"SendPage": wireDecl(wireApply("Page", wireTag("wire", wireText()))),
		// All names below are authored ordinary declarations, including an unfamiliar one.
		"Function": wireDecl(wireTag("wire", wireTag("tuple", wireParam(0), wireTag("wire", wireParam(1)))), "A", "R"),
		"Stream":   wireDecl(wireTag("wire", wireTag("wire", wireParam(0))), "T"),
		"Cell":     wireDecl(wireTag("tuple", wireApply("Function", wireTag("tuple"), wireParam(0)), wireApply("Function", wireParam(0), wireTag("tuple"))), "T"),
		"Mailbox":  wireDecl(wireTag("wire", wireApply("Page", wireParam(0))), "T"),
	})
}

func wireEngine(t *testing.T, tree s.AST) (*engine.Engine, semantic.Snapshot) {
	t.Helper()
	p := provider.Provider{}
	definition, err := representation.NewLanguage[s.Language](p, s.Profile)
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := representation.Pack[s.Language](definition, tree)
	if err != nil {
		t.Fatal(err)
	}
	registry := &engine.Engine{}
	if err := registry.RegisterLanguage(definition); err != nil {
		t.Fatal(err)
	}
	resolver, err := resolution.NewResolver[s.Language](definition, resolve.Contract(), p)
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterResolver(resolver); err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[s.Language](definition, analyze.Contract(), p)
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterAnalyzer(analyzer); err != nil {
		t.Fatal(err)
	}
	snapshot, err := registry.Capture(context.Background(), semantic.Input{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: []representation.AST{envelope}})
	if err != nil {
		t.Fatal(err)
	}
	return registry, snapshot
}

func qualifyHorosWire(t *testing.T) {
	p := provider.Provider{}
	if p.ID() != "horoswire" || p.Name() != "Horos/Wire" {
		t.Fatal("unstable optional identity")
	}
	if _, ok := any(p).(parse.Provider[s.Language]); ok {
		t.Fatal("optional provider claims text parsing")
	}
	if _, ok := any(p).(format.Provider[s.Language]); ok {
		t.Fatal("optional provider claims text formatting")
	}
	module := wireFixture(t)
	tree, err := resolve.Resolve(module, nil, "optional")
	if err != nil {
		t.Fatal(err)
	}
	model, err := resolve.Model(tree)
	if err != nil {
		t.Fatal(err)
	}
	pin, _ := decl.Pin(module)
	result, err := model.Specialize(decl.DeclarationID{Module: pin, Name: "SendPage"}, nil)
	if err != nil || !core.Equal(result, wireTag("list-of", wireTag("wire", wireText()))) {
		t.Fatalf("wire pagination: %v", err)
	}
	if _, err := groundResolve.Resolve(module, nil, "wrong"); err == nil {
		t.Fatal("ground resolver accepted optional module")
	}
	registry, snapshot := wireEngine(t, tree)
	if _, err := registry.Analyze(context.Background(), analysis.Input{Snapshot: snapshot}); err == nil {
		t.Fatal("analysis silently resolved")
	}
	resolved, err := registry.Resolve(context.Background(), snapshot, semantic.Options{})
	if err != nil || !resolved.HasModel() {
		t.Fatalf("resolution: %v", err)
	}
	analyzed, err := registry.Analyze(context.Background(), analysis.Input{Snapshot: snapshot, Resolution: &resolved})
	if err != nil || !analyzed.HasModel() {
		t.Fatalf("analysis: %v", err)
	}
	first, _ := resolved.Model()
	second, _ := analyzed.Model()
	if first != second {
		t.Fatal("analysis rebuilt or replaced the resolution model")
	}
	for _, k := range []semantics.Key{resolve.Validity(), resolve.ReferenceResolution()} {
		report, ok := resolved.Report(k)
		if !ok || report.Status != semantic.Complete || report.Verdict != semantic.Holds {
			t.Fatalf("unproved %v: %+v", k, report)
		}
	}
	if _, ok := resolved.Report(groundResolve.Validity()); ok {
		t.Fatal("optional evidence impersonates ground evidence")
	}
	facts := first.Data()
	foundWire, foundData := false, false
	for _, term := range facts.Terms {
		if term.Kind != semantics.ProviderExpression {
			continue
		}
		if term.Operator.Owner != "horos/wire" {
			t.Fatalf("wrong semantic owner: %+v", term.Operator)
		}
		if term.Operator.Name == "descriptor/wire" {
			foundWire = true
			if len(term.Operands) != 1 {
				t.Fatal("Wire message operand missing")
			}
		}
		if term.Operator.Name == "descriptor/data" {
			foundData = true
			if len(term.Operands) != 0 {
				t.Fatal("Data payload exposed upper type positions")
			}
		}
	}
	if !foundWire || !foundData {
		t.Fatal("optional terms missing")
	}
	// Go marker, language ID and compatibility profile are independent boundaries.
	groundDefinition, err := representation.NewLanguage[groundSyntax.Language](groundProvider.Provider{}, groundSyntax.Profile)
	if err != nil {
		t.Fatal(err)
	}
	optionalDefinition, _ := representation.NewLanguage[s.Language](p, s.Profile)
	packed, _ := representation.Pack[s.Language](optionalDefinition, tree)
	if _, err := representation.Unpack[groundSyntax.Language](groundDefinition, packed); err == nil {
		t.Fatal("optional AST unpacked as ground")
	}
}

func TestHorosWireUnresolvedSyntaxAndLocatedRejection(t *testing.T) {
	module := wireModule(t, "example.unresolved", map[string]core.Value{"Use": wireDecl(wireApply("Use"))})
	// Construct the syntax directly: source admission must not require resolution.
	pin, _ := decl.Pin(module)
	tree := s.AST{ID: "unresolved", Root: 1, Nodes: []s.Node{
		{Kind: s.KindProgram, Children: []ast.Child{{Role: s.RoleModule, Node: 2}}, References: []ast.Reference{{Role: s.RoleRoot, Target: ast.NodeRef{Tree: "unresolved", Node: 2}}}},
		{Kind: s.KindModule, Attributes: []ast.Attribute{{Name: s.AttributeName, Text: "example.unresolved"}, {Name: s.AttributePin, Text: pin.String()}}, Children: []ast.Child{{Role: s.RoleDeclaration, Node: 3}}},
		{Kind: s.KindDeclaration, Attributes: []ast.Attribute{{Name: s.AttributeName, Text: "Use"}}, Children: []ast.Child{{Role: s.RoleBody, Node: 4}}},
		{Kind: s.KindApplication, Attributes: []ast.Attribute{{Name: s.AttributeScope, Text: "local"}, {Name: s.AttributeName, Text: "Use"}}},
	}}
	if err := check.Check(tree); err != nil {
		t.Fatal(err)
	}
	registry, snapshot := wireEngine(t, tree)
	result, err := registry.Resolve(context.Background(), snapshot, semantic.Options{})
	if err != nil || result.HasModel() {
		t.Fatalf("invalid cycle published model: %v", err)
	}
	report, _ := result.Report(resolve.Validity())
	if report.Status != semantic.Complete || report.Verdict != semantic.Rejected {
		t.Fatalf("wrong rejection: %+v", report)
	}
	diagnostics := result.Diagnostics()
	if len(diagnostics) == 0 || diagnostics[0].At == nil || diagnostics[0].At.Node != 4 {
		t.Fatalf("unlocated failure: %+v", diagnostics)
	}
}

func TestHorosWireRefusalsAndLiteralOpacity(t *testing.T) {
	// A slot-shaped exact value is data, never a binder occurrence.
	literal := core.NewTuple(core.NewAtom([]byte("horos/exact")), wireParam(50))
	module := wireModule(t, "example.literal", map[string]core.Value{"Literal": wireDecl(wireTag("data", literal))})
	tree, err := resolve.Resolve(module, nil, "literal")
	if err != nil {
		t.Fatal(err)
	}
	for _, node := range tree.Nodes {
		if node.Kind == s.KindParameter {
			t.Fatal("literal promoted to binder")
		}
	}
	_, err = resolve.Model(tree)
	if err != nil {
		t.Fatal(err)
	}
	for _, body := range []core.Value{wireTag("data", wireParam(0)), wireTag("union", wireText()), wireTag("wire"), wireTag("cell", wireText())} {
		got, err := resolve.Resolve(wireModule(t, "example.invalid", map[string]core.Value{"Bad": wireDecl(body, "T")}), nil, "bad")
		if err == nil || len(got.Nodes) != 0 {
			t.Fatal("invalid constructor published AST")
		}
	}
	registry, snapshot := wireEngine(t, tree)
	result, err := registry.Resolve(context.Background(), snapshot, semantic.Options{Configuration: []semantic.Setting{{Key: resolve.LimitsSetting(), Value: []byte("[1,1,1,1]")}}})
	if err != nil || result.HasModel() {
		t.Fatalf("bounded resolution: %v", err)
	}
	report, _ := result.Report(resolve.Validity())
	if report.Status != semantic.Incomplete || report.Verdict != semantic.Undetermined {
		t.Fatalf("exhaustion became rejection: %+v", report)
	}
	_, err = registry.Resolve(context.Background(), snapshot, semantic.Options{Configuration: []semantic.Setting{{Key: groundResolve.LimitsSetting(), Value: []byte("[1,1,1,1]")}}})
	if !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatalf("foreign configuration accepted: %v", err)
	}
}
