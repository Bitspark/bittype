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
	codec "github.com/bitspark/ontos/codec/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

func bounded(body value, name, bound string) value {
	return d("type", data.EncodeList([]value{d("parameter", txt(name), atom(bound))}), body)
}

func unifiedModule(forwardBound string) value {
	return module("unified.pagination", mapping(),
		entry("Page", bounded(page(param(0)), "T", "type")),
		entry("DataPage", bounded(local("Page", param(0)), "D", "data")),
		entry("Forward", bounded(local("DataPage", param(0)), "Forwarded", forwardBound)),
		entry("People", definition(local("Page", person()))),
		entry("Live", definition(local("Page", h("wire", person())))))
}

func TestUnifiedPaginationBoundsSurviveASTAndSemanticExport(t *testing.T) {
	root := unifiedModule("data")
	tree, err := resolve.Resolve(root, nil, "unified")
	if err != nil {
		t.Fatal(err)
	}
	if err := check.Check(tree); err != nil {
		t.Fatal(err)
	}
	program, err := resolve.Model(tree)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"People", "Live"} {
		if value, err := program.Specialize(decl.DeclarationID{Module: program.Root(), Name: name}, nil); err != nil || value == nil {
			t.Fatal("one frontend rejected an admitted Page application", err)
		}
	}
	if got, err := program.DataDescriptor(decl.DeclarationID{Module: program.Root(), Name: "People"}, nil); err != nil || !core.Equal(got, page(person())) {
		t.Fatal("data page changed through the AST projection", err)
	}
	if got, err := program.DataDescriptor(decl.DeclarationID{Module: program.Root(), Name: "Live"}, nil); err == nil || got != nil {
		t.Fatal("live page was admitted as a data descriptor")
	}
	model := requireModel(t, run(t, tree, nil))
	for _, wanted := range []struct{ name, bound string }{{"T", "type"}, {"D", "data"}, {"Forwarded", "data"}} {
		site := ast.NodeRef{Tree: tree.ID, Node: find(tree, s.KindTypeParameter, wanted.name)}
		termID, ok := model.TermAt(semantics.Occurrence{Site: site, Role: resolve.BoundRole()})
		if !ok {
			t.Fatalf("missing bound for %s", wanted.name)
		}
		fact := term(t, model, termID)
		if fact.Kind != semantics.ProviderExpression || fact.Operator != (semantics.Key{Owner: "horos", Name: "bound/" + wanted.bound}) || string(fact.Data) != string(codec.Encode(atom(wanted.bound))) {
			t.Fatal("parameter bound lost its provider-owned meaning")
		}
	}
	_, alias := entity(t, model, "unified.pagination", "DataPage")
	pageID, pageEntity := entity(t, model, "unified.pagination", "Page")
	application := term(t, model, alias.Body)
	if application.Kind != semantics.Application || application.Target != pageID || pageEntity.Parameters[0].Sort != (semantics.Key{Owner: "horos", Name: "type"}) {
		t.Fatal("constrained alias became a second Page or a different parameter sort")
	}
}

func TestUnboundedForwardingReachesSemanticRejection(t *testing.T) {
	tree, err := resolve.Resolve(unifiedModule("data"), nil, "changed-bound")
	if err != nil {
		t.Fatal(err)
	}
	stripAnnotations(tree)
	parameter := find(tree, s.KindTypeParameter, "Forwarded")
	bound := tree.Nodes[parameter-1].Children[0].Node
	setAttribute(&tree.Nodes[bound-1], s.AttributeName, "type")
	setAttribute(&tree.Nodes[find(tree, s.KindModule, "unified.pagination")-1], s.AttributePin, pin(t, unifiedModule("type")).String())
	if err := check.Check(tree); err != nil {
		t.Fatal("syntax checker performed bound semantics", err)
	}
	for _, result := range []semantic.Result{runResolution(t, tree, nil), run(t, tree, nil)} {
		requireFailure(t, result, semantic.Complete, semantic.Rejected, "bound-violation", "source")
		if result.Diagnostics()[0].At == nil {
			t.Fatal("bound failure lost its source location")
		}
	}
	setAttribute(&tree.Nodes[bound-1], s.AttributeName, "unknown")
	if err := check.Check(tree); err == nil {
		t.Fatal("unknown bound passed syntax admission")
	}
}
