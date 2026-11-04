package analyze_test

import (
	"testing"

	check "github.com/Bitspark/bittype/check/horos/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

func constrained(body value, name string, contract value) value {
	return d("type", data.EncodeList([]value{d("parameter", txt(name), contract)}), d("top", d("sort-proper")), body)
}
func unifiedModule(forward value) value {
	top := d("top", d("sort-proper"))
	dataKind := d("kind-imported", txt("kinds"), txt("Data"))
	return module("unified.pagination", mapping(entry("kinds", decl.PreludeImport())),
		entry("Page", constrained(page(param(0)), "T", top)),
		entry("DataPage", constrained(local("Page", param(0)), "D", dataKind)),
		entry("Forward", constrained(local("DataPage", param(0)), "Forwarded", forward)),
		entry("People", definition(local("Page", person()))),
		entry("Live", definition(local("Page", h("wire", person())))))
}

func TestUnifiedPaginationContractsSurviveASTAndSemanticExport(t *testing.T) {
	tree, err := resolve.Resolve(unifiedModule(d("kind-imported", txt("kinds"), txt("Data"))), nil, "unified")
	if err != nil { t.Fatal(err) }
	if err := check.Check(tree); err != nil { t.Fatal(err) }
	program, err := resolve.Model(tree)
	if err != nil { t.Fatal(err) }
	for _, name := range []string{"People", "Live"} {
		value, err := program.Value(decl.DeclarationID{Module: program.Root(), Name: name})
		if err != nil { t.Fatal(err) }
		if got, err := program.Descriptor(value); err != nil || got == nil { t.Fatal("one frontend rejected an admitted Page application", err) }
		got, err := program.GroundDescriptor(value)
		if name == "People" && (err != nil || !core.Equal(got, page(person()))) { t.Fatal("ground page changed through projection", err) }
		if name == "Live" && (err == nil || got != nil) { t.Fatal("capability page admitted as a ground descriptor") }
	}
	model := requireModel(t, run(t, tree, nil))
	for _, wanted := range []string{"T", "D", "Forwarded"} {
		var id semantics.EntityID
		for i, entity := range model.Data().Entities { if entity.Kind.Name == "type-parameter" && entity.Name == wanted { id = semantics.EntityID(i+1) } }
		facts := model.Facts(semantics.Subject{Entity: id}, resolve.DeclaredKindRelation())
		if len(facts) != 1 { t.Fatalf("missing declared kind for %s", wanted) }
		fact := term(t, model, facts[0])
		if wanted == "T" {
			if fact.Kind != semantics.ProviderExpression || fact.Operator.Name != "kind/top" { t.Fatal("top contract lost") }
		} else {
			target, ok := model.Entity(fact.Target)
			if fact.Kind != semantics.Reference || !ok || target.Kind.Name != "kind-definition" || target.Name != "Data" { t.Fatal("ordinary imported contract identity lost") }
		}
	}
	_, alias := entity(t, model, "unified.pagination", "DataPage")
	pageID, pageEntity := entity(t, model, "unified.pagination", "Page")
	application := term(t, model, alias.Body)
	if application.Kind != semantics.Application || term(t, model, application.Callee).Target != pageID || pageEntity.Parameters[0].Sort != (semantics.Key{Owner: "horos", Name: "proper"}) { t.Fatal("constrained alias changed Page identity or parameter sort") }
}

func TestUnconstrainedForwardingReachesSemanticRejection(t *testing.T) {
	tree, err := resolve.Resolve(unifiedModule(d("kind-imported", txt("kinds"), txt("Data"))), nil, "changed-contract")
	if err != nil { t.Fatal(err) }
	stripAnnotations(tree)
	parameter := find(tree, s.KindTypeParameter, "Forwarded")
	contract := tree.Nodes[parameter-1].Children[0].Node
	tree.Nodes = append(tree.Nodes, s.Node{Kind: s.KindProperSort})
	tree.Nodes[contract-1] = s.Node{Kind: s.KindTop, Children: []ast.Child{{Role: s.RoleSort, Node: ast.NodeID(len(tree.Nodes))}}}
	setAttribute(&tree.Nodes[find(tree, s.KindModule, "unified.pagination")-1], s.AttributePin, pin(t, unifiedModule(d("top", d("sort-proper")))).String())
	if err := check.Check(tree); err != nil { t.Fatal("syntax checker performed contract semantics", err) }
	for _, result := range []semantic.Result{runResolution(t, tree, nil), run(t, tree, nil)} {
		if result.HasSemantics() { t.Fatal("unproved source contract published semantics") }
		report, _ := result.Report(resolve.Validity())
		if report.Verdict == semantic.Holds || len(result.Diagnostics()) == 0 || result.Diagnostics()[0].At == nil { t.Fatal("contract failure lost its verdict or source location") }
	}
	tree.Nodes[contract-1].Kind = "unknown-kind"
	if err := check.Check(tree); err == nil { t.Fatal("unknown kind passed syntax admission") }
}
