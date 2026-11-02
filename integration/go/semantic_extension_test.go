package integration_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/Bitspark/bittype/contracts/analysis/go"
 "github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/engine/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// This independent, test-only language is not a native compiler adapter. Its
// fixed fixture grammar tests structural extensibility without asserting Go,
// Rust, TypeScript, C++, or Horos semantic support.
type extensionLanguage struct{}

func extensionKey(name string) semantics.Key {
	return semantics.Key{Owner: "third-party.example/types", Name: name}
}

func extensionContract() analysis.Contract {
	return analysis.Contract{Semantic: semantic.Contract{
		Profile: "foreign-declarations-v1",
		Obligations: []semantic.Obligation{
			{Key: extensionKey("validity"), Domain: extensionKey("entry")},
			{Key: extensionKey("bindings"), Domain: extensionKey("entry")},
			{Key: extensionKey("applications"), Domain: extensionKey("entry")},
			{Key: extensionKey("dependencies"), Domain: extensionKey("entry")},
		},
		ModelRequirements: []semantic.Requirement{
			{Obligation: extensionKey("validity"), Verdict: semantic.Holds},
			{Obligation: extensionKey("bindings"), Verdict: semantic.Undetermined},
			{Obligation: extensionKey("applications"), Verdict: semantic.Undetermined},
			{Obligation: extensionKey("dependencies"), Verdict: semantic.Undetermined},
		},
	}}
}

func extensionTree() ast.AST[extensionLanguage] {
	tree := ast.AST[extensionLanguage]{ID: "foreign-input", Root: 1, Nodes: []ast.Node[extensionLanguage]{
		{Kind: "document"},
		{Kind: "generic", Attributes: attr("name", "Select")},
		{Kind: "parameter", Attributes: attr("name", "T", "sort", "lifetime")},
		{Kind: "parameter", Attributes: attr("name", "T", "sort", "type")},
		{Kind: "parameter", Attributes: attr("name", "N", "sort", "constant")},
		{Kind: "alias", Attributes: attr("name", "Same")},
		{Kind: "defined", Attributes: attr("name", "Same")},
		{Kind: "recursive", Attributes: attr("name", "Node")},
		{Kind: "array", Attributes: attr("name", "Three", "length", "3")},
		{Kind: "array", Attributes: attr("name", "Four", "length", "4")},
		{Kind: "application", Attributes: attr("name", "Select")},
		{Kind: "use", Attributes: attr("name", "Word")},
		{Kind: "qualified-uses", Attributes: attr("name", "T")},
		{Kind: "generic", Attributes: attr("name", "Other")},
		{Kind: "parameter", Attributes: attr("name", "T", "sort", "type")},
		{Kind: "defined", Attributes: attr("name", "Cell")},
		{Kind: "defined", Attributes: attr("name", "Function")},
		{Kind: "defined", Attributes: attr("name", "Stream")},
	}}
	for i := 1; i < len(tree.Nodes); i++ {
		tree.Nodes[0].Children = append(tree.Nodes[0].Children, ast.Child{Role: "item", Node: ast.NodeID(i + 1)})
	}
	// Historical provenance is retained but is not a current analysis input.
	tree.Nodes[1].Origin.Nodes = []ast.NodeRef{{Tree: "historical-input", Node: 77}}
	return tree
}

func extensionAttribute(tree ast.AST[extensionLanguage], node ast.NodeID, name string) string {
	for _, attribute := range tree.Nodes[node-1].Attributes {
		if attribute.Name == name {
			return attribute.Text
		}
	}
	return ""
}

func extensionOccurrence(tree ast.TreeID, node ast.NodeID, role string, ordinal uint32) semantics.Occurrence {
	return semantics.Occurrence{Site: ast.NodeRef{Tree: tree, Node: node}, Role: extensionKey(role), Ordinal: ordinal}
}

func extensionFacts(tree ast.AST[extensionLanguage]) semantics.Data {
	entity := func(node ast.NodeID, kind string, body semantics.TermID) semantics.Entity {
		return semantics.Entity{
			Kind: extensionKey(kind), Name: extensionAttribute(tree, node, "name"), Body: body,
			Declarations: []ast.NodeRef{{Tree: tree.ID, Node: node}},
		}
	}
	providerTerm := func(operator, encoding string, data []byte, operands ...semantics.Operand) semantics.Term {
		return semantics.Term{Kind: semantics.ProviderExpression, Operator: extensionKey(operator), Encoding: extensionKey(encoding), Data: data, Operands: operands}
	}
	entities := []semantics.Entity{
		entity(2, "defined", 1), entity(3, "parameter", 0), entity(4, "parameter", 0), entity(5, "parameter", 0),
		entity(6, "alias", 6), entity(7, "defined", 7), entity(8, "recursive", 8),
		{Kind: extensionKey("builtin"), Name: "Word"},
		entity(9, "defined", 10), entity(10, "defined", 11), entity(14, "defined", 12), entity(15, "parameter", 0),
		entity(16, "defined", 7), entity(17, "defined", 7), entity(18, "defined", 7),
	}
	for i, node := range []ast.NodeID{3, 4, 5} {
		parameter := semantics.EntityID(i + 2)
		entities[0].Parameters = append(entities[0].Parameters, semantics.Parameter{Entity: parameter, Sort: extensionKey(extensionAttribute(tree, node, "sort"))})
		entities[parameter-1].Owner = 1
	}
	entities[10].Parameters = []semantics.Parameter{{Entity: 12, Sort: extensionKey(extensionAttribute(tree, 15, "sort"))}}
	entities[11].Owner = 11
	terms := []semantics.Term{
		providerTerm("tuple", "opaque-v1", []byte{255, 0}, semantics.Operand{Locator: []uint32{0}, Term: 2}, semantics.Operand{Locator: []uint32{1}, Term: 3}, semantics.Operand{Locator: []uint32{2}, Term: 4}),
		{Kind: semantics.Reference, Target: 2},
		{Kind: semantics.Reference, Target: 3},
		{Kind: semantics.Reference, Target: 4},
		providerTerm("natural", "ascii-natural-v1", []byte("3")),
		{Kind: semantics.Application, Target: 1, Arguments: []semantics.TermID{13, 7, 5}},
		{Kind: semantics.Reference, Target: 8},
		providerTerm("mutable-pointer", "opaque-v1", nil, semantics.Operand{Locator: []uint32{0}, Term: 9}),
		{Kind: semantics.Reference, Target: 7},
		providerTerm("array", "element-and-length-v1", []byte(extensionAttribute(tree, 9, "length")), semantics.Operand{Locator: []uint32{0}, Term: 7}),
		providerTerm("array", "element-and-length-v1", []byte(extensionAttribute(tree, 10, "length")), semantics.Operand{Locator: []uint32{0}, Term: 7}),
		{Kind: semantics.Reference, Target: 12},
		providerTerm("lifetime-static", "opaque-v1", []byte("static")),
		{Kind: semantics.Application, Target: 8},
		{Kind: semantics.Reference, Target: 1},
	}
	use := extensionOccurrence(tree.ID, 11, "use", 0)
	return semantics.Data{
		Entities: entities, Terms: terms,
		Bindings: []semantics.Binding{
			{At: use, Target: 1},
			{At: extensionOccurrence(tree.ID, 12, "use", 0), Target: 8},
			{At: extensionOccurrence(tree.ID, 13, "use", 0), Target: 2},
			{At: extensionOccurrence(tree.ID, 13, "use", 1), Target: 3},
			{At: extensionOccurrence(tree.ID, 14, "body-use", 0), Target: 12},
		},
		Expressions: []semantics.Expression{
			{At: extensionOccurrence(tree.ID, 11, "expression", 0), Term: 6},
			{At: extensionOccurrence(tree.ID, 8, "expression", 0), Term: 8},
			{At: extensionOccurrence(tree.ID, 9, "expression", 0), Term: 10},
			{At: extensionOccurrence(tree.ID, 10, "expression", 0), Term: 11},
		},
		Dependencies: []semantics.Dependency{
			{From: 5, To: 1, Relation: extensionKey("direct"), At: &use},
			{From: 5, To: 8, Relation: extensionKey("direct")},
			{From: 7, To: 7, Relation: extensionKey("direct")},
			{From: 9, To: 8, Relation: extensionKey("direct")},
			{From: 10, To: 8, Relation: extensionKey("direct")},
		},
	}
}

type extensionProvider struct {
	badSite    *ast.NodeRef
	incomplete bool
}

func (extensionProvider) ID() language.ID { return "foreign-semantic-fixture" }
func (extensionProvider) Name() string    { return "Foreign semantic fixture" }
func (extensionProvider) Check(view ast.View[extensionLanguage]) error {
	tree := view.ASTView()
	if tree.Root != 1 || len(tree.Nodes) != 18 || tree.Nodes[0].Kind != "document" {
		return errors.New("outside the independent fixture grammar")
	}
	return nil
}

func (p extensionProvider) Analyze(ctx context.Context, request analysis.Request[extensionLanguage]) (semantic.Result, error) {
	if err := ctx.Err(); err != nil {
		return semantic.Result{}, err
	}
	inputs := request.Inputs()
	if len(inputs) != 1 {
		return semantic.Result{}, errors.New("fixture expects one input")
	}
	facts := extensionFacts(inputs[0])
	publication := semantic.Publication{Facts: &facts}
	known := make(map[semantics.Key]bool)
	for _, requirement := range extensionContract().ModelRequirements {
		known[requirement.Obligation] = true
		report := semantic.Report{Obligation: requirement.Obligation, Status: semantic.Complete, Verdict: requirement.Verdict}
		if p.incomplete {
			report.Status, report.Verdict = semantic.Incomplete, semantic.Undetermined
			publication.Facts = nil
		}
		publication.Reports = append(publication.Reports, report)
	}
	for _, key := range request.Requested() {
		if !known[key] {
			publication.Reports = append(publication.Reports, semantic.Report{Obligation: key, Status: semantic.Unsupported})
		}
	}
	if p.badSite != nil {
		facts.Entities[0].Declarations[0] = *p.badSite
	}
	return request.Publish(publication)
}

func extensionKernel(t *testing.T, provider extensionProvider, tree ast.AST[extensionLanguage]) (*engine.Engine, analysisCase) {
	t.Helper()
	syntax, err := representation.NewLanguage[extensionLanguage](provider, "fixture-syntax-v1")
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := representation.Pack(syntax, tree)
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[extensionLanguage](syntax, extensionContract(), provider)
	if err != nil {
		t.Fatal(err)
	}
	kernel := engine.New()
	if err := kernel.RegisterLanguage(syntax); err != nil {
		t.Fatal(err)
	}
	if err := kernel.RegisterAnalyzer(analyzer); err != nil {
		t.Fatal(err)
	}
	return kernel, analysisCase{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: []representation.AST{envelope}}
}

type semanticQueryResult struct {
	Target       semantics.EntityID
	Bindings     []semantics.ArgumentBinding
	Dependencies []semantics.Dependency
}

// This consumer depends only on shared result/model contracts. It never opens a
// language AST, decodes an opaque payload, or invokes a concrete provider.
func consumeSemanticFacts(result semantic.Result, reference, expression semantics.Occurrence, owner semantics.EntityID, relation semantics.Key, required []semantic.Requirement) (semanticQueryResult, error) {
	for _, requirement := range required {
		report, found := result.Report(requirement.Obligation)
		if !found || report.Status != semantic.Complete || report.Verdict != requirement.Verdict {
			return semanticQueryResult{}, fmt.Errorf("required facts were not established: %v", requirement.Obligation)
		}
	}
	model, ok := result.Model()
	if !ok {
		return semanticQueryResult{}, errors.New("required model is absent")
	}
	target, ok := model.Resolve(reference)
	if !ok {
		return semanticQueryResult{}, errors.New("reference was not bound")
	}
	term, ok := model.TermAt(expression)
	if !ok {
		return semanticQueryResult{}, errors.New("expression meaning is absent")
	}
	bindings, ok := model.ParameterBindings(term)
	if !ok {
		return semanticQueryResult{}, errors.New("expression is not an application")
	}
	return semanticQueryResult{Target: target, Bindings: bindings, Dependencies: model.DependenciesFrom(owner, relation)}, nil
}

func TestSemanticForeignProviderAndGenericConsumer(t *testing.T) {
	tree := extensionTree()
	kernel, input := extensionKernel(t, extensionProvider{}, tree)
	result, err := analyzeCase(kernel, context.Background(), input)
	if err != nil || !result.Valid() || !result.HasModel() {
		t.Fatalf("foreign analysis: %v", err)
	}
	queried, err := consumeSemanticFacts(result,
		extensionOccurrence(tree.ID, 11, "use", 0), extensionOccurrence(tree.ID, 11, "expression", 0),
		5, extensionKey("direct"), extensionContract().ModelRequirements)
	if err != nil {
		t.Fatal(err)
	}
	wantBindings := []semantics.ArgumentBinding{{Parameter: 2, Argument: 13}, {Parameter: 3, Argument: 7}, {Parameter: 4, Argument: 5}}
	if queried.Target != 1 || !reflect.DeepEqual(queried.Bindings, wantBindings) || len(queried.Dependencies) != 2 || queried.Dependencies[0].To != 1 || queried.Dependencies[1].To != 8 {
		t.Fatalf("four agnostic queries lost facts: %#v", queried)
	}
	model, _ := result.Model()
	alias, _ := model.Entity(5)
	defined, _ := model.Entity(6)
	if alias.Name != defined.Name || alias.Kind == defined.Kind {
		t.Fatal("same display names collapsed alias and defined entities")
	}
	for i, sort := range []string{"lifetime", "type", "constant"} {
		generic, _ := model.Entity(1)
		if generic.Parameters[i].Sort != extensionKey(sort) {
			t.Fatalf("parameter sort %s was narrowed", sort)
		}
	}
	firstT, _ := model.Resolve(extensionOccurrence(tree.ID, 13, "use", 1))
	otherT, _ := model.Resolve(extensionOccurrence(tree.ID, 14, "body-use", 0))
	if firstT != 3 || otherT != 12 || firstT == otherT {
		t.Fatal("same-spelled binders were merged")
	}
	builtin, _ := model.Entity(8)
	if len(builtin.Declarations) != 0 {
		t.Fatal("builtin acquired a fabricated source declaration")
	}
	three, _ := model.Term(10)
	four, _ := model.Term(11)
	if three.Operator != four.Operator || string(three.Data) != "3" || string(four.Data) != "4" {
		t.Fatal("literal array lengths collapsed")
	}
	if edges := model.DependenciesFrom(7, extensionKey("direct")); len(edges) != 1 || edges[0].To != 7 {
		t.Fatal("entity recursion was rejected or erased")
	}
	if bindings, ok := model.ParameterBindings(14); !ok || len(bindings) != 0 {
		t.Fatal("zero-argument application lost its form")
	}
	if _, ok := model.ParameterBindings(15); ok {
		t.Fatal("unapplied generic reference acquired a binding")
	}
	for i, name := range []string{"Cell", "Function", "Stream"} {
		entity, ok := model.Entity(semantics.EntityID(13 + i))
		if !ok || entity.Name != name || entity.Kind != extensionKey("defined") {
			t.Fatalf("library name %s acquired special treatment", name)
		}
	}
}

func TestSemanticConsumerRequiresCoverage(t *testing.T) {
	kernel, input := extensionKernel(t, extensionProvider{incomplete: true}, extensionTree())
	result, err := analyzeCase(kernel, context.Background(), input)
	if err != nil || !result.Valid() || result.HasModel() {
		t.Fatalf("incomplete result boundary: %v", err)
	}
	_, err = consumeSemanticFacts(result,
		extensionOccurrence("foreign-input", 11, "use", 0), extensionOccurrence("foreign-input", 11, "expression", 0),
		5, extensionKey("direct"), extensionContract().ModelRequirements)
	if err == nil {
		t.Fatal("nil Go error was mistaken for established semantic coverage")
	}
}

func TestSemanticForeignProviderCannotPublishHistoricalOrZeroSites(t *testing.T) {
	for _, site := range []ast.NodeRef{{Tree: "historical-input", Node: 77}, {Tree: "foreign-input", Node: 0}} {
		t.Run(fmt.Sprintf("%s-%d", site.Tree, site.Node), func(t *testing.T) {
			kernel, input := extensionKernel(t, extensionProvider{badSite: &site}, extensionTree())
			result, err := analyzeCase(kernel, context.Background(), input)
			if !errors.Is(err, semantic.ErrInvalidResult) || result.Valid() || result.HasModel() {
				t.Fatalf("non-current site became semantic provenance: %#v, %v", site, err)
			}
		})
	}
}
