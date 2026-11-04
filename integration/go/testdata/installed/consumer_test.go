package consumer_test

import (
	"math/big"
	"testing"

	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	wireResolve "github.com/Bitspark/bittype/resolve/horoswire/go"
	"github.com/Bitspark/horos/core/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	wire "github.com/Bitspark/horos/core/go/wire"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

func tag(name string, values ...core.Value) core.Value {
	return core.NewTuple(append([]core.Value{core.NewAtom([]byte(name))}, values...)...)
}

// This fixture is copied outside the checkout and imports the exact public
// Bittype revision, Horos release and Ontos release without replace directives.
func TestPublicPaginationAndLift(t *testing.T) {
	text := tag("horos/scalar", core.NewAtom([]byte("utf8-text")))
	parameter := tag("horos/decl/param", data.EncodeInt(big.NewInt(0)))
	page := tag("horos/decl/type", data.EncodeList([]core.Value{data.MustEncodeText("T")}), tag("horos/list-of", parameter))
	members, err := data.EncodeMap([][2]core.Value{{data.MustEncodeText("Page"), page}})
	if err != nil {
		t.Fatal(err)
	}
	empty, err := data.EncodeMap(nil)
	if err != nil {
		t.Fatal(err)
	}
	module := tag("horos/decl/module", data.EncodeInt(big.NewInt(1)), data.MustEncodeText("installed.pagination"), empty, members)
	tree, err := resolve.Resolve(module, nil, "installed-ground")
	if err != nil {
		t.Fatal(err)
	}
	model, err := resolve.Model(tree)
	if err != nil {
		t.Fatal(err)
	}
	id := decl.DeclarationID{Module: model.Root(), Name: "Page"}
	specialized, err := model.Specialize(id, []core.Value{text})
	if err != nil || !core.Equal(specialized, tag("horos/list-of", text)) {
		t.Fatalf("ground generic binding: %v", err)
	}
	if !horos.Conforms(data.EncodeList([]core.Value{data.MustEncodeText("Ada")}), specialized) {
		t.Fatal("public Ontos values did not conform to the specialized Horos descriptor")
	}
	wireArgument := tag("horos/wire/wire", tag("horos/wire/data", text))
	if _, err := model.Specialize(id, []core.Value{wireArgument}); err == nil {
		t.Fatal("ground generic accepted an optional capability type")
	}
	lifted, err := wire.Lift(model, []decl.DeclarationID{id}, wire.DefaultLiftOptions())
	if err != nil {
		t.Fatal(err)
	}
	root := lifted.Roots()[0]
	optional, err := lifted.Program(root.Declaration.Module, wire.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	optionalTree, err := wireResolve.Project(optional, "installed-optional")
	if err != nil {
		t.Fatal(err)
	}
	optional, err = wireResolve.Model(optionalTree)
	if err != nil {
		t.Fatal(err)
	}
	specialized, err = optional.Specialize(root.Declaration, []core.Value{wireArgument})
	if err != nil || !core.Equal(specialized, tag("horos/wire/list-of", wireArgument)) {
		t.Fatalf("lifted generic binding: %v", err)
	}
	if root.Declaration.Module == id.Module {
		t.Fatal("explicit lifting reused the ground module identity")
	}
}
