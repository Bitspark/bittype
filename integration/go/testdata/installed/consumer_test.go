package consumer_test

import (
	"errors"
	"math/big"
	"testing"

	resolve "github.com/Bitspark/bittype/resolve/files/horos/provider/go"
	"github.com/Bitspark/horos/core/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

func tag(name string, values ...core.Value) core.Value {
	return core.NewTuple(append([]core.Value{core.NewAtom([]byte(name))}, values...)...)
}

// Copied outside the checkout and installed with exact public dependencies,
// without replace directives or access to repository-private helpers.
func TestPublicUnifiedPagination(t *testing.T) {
	atom := func(s string) core.Value { return core.NewAtom([]byte(s)) }
	text := tag("horos/scalar", atom("utf8-text"))
	parameter := tag("horos/decl/param", data.EncodeInt(big.NewInt(0)))
	pageOf := func(element core.Value) core.Value {
		return tag("horos/record", atom("closed"),
			tag("horos/field", data.MustEncodeText("items"), atom("req"), tag("horos/list-of", element)),
			tag("horos/field", data.MustEncodeText("next"), atom("opt"), text))
	}
	pageBody := pageOf(parameter)
	top := tag("horos/decl/top", tag("horos/decl/sort-proper"))
	dataKind := tag("horos/decl/kind-imported", data.MustEncodeText("h"), data.MustEncodeText("Data"))
	page := tag("horos/decl/type", data.EncodeList([]core.Value{
		tag("horos/decl/parameter", data.MustEncodeText("T"), top),
	}), top, pageBody)
	dataPage := tag("horos/decl/type", data.EncodeList([]core.Value{
		tag("horos/decl/parameter", data.MustEncodeText("D"), dataKind),
	}), top, tag("horos/decl/apply", tag("horos/decl/local", data.MustEncodeText("Page")), data.EncodeList([]core.Value{parameter})))
	members, err := data.EncodeMap([][2]core.Value{
		{data.MustEncodeText("Page"), page}, {data.MustEncodeText("DataPage"), dataPage},
	})
	if err != nil {
		t.Fatal(err)
	}
	imports, err := data.EncodeMap([][2]core.Value{{data.MustEncodeText("h"), decl.PreludeImport()}})
	if err != nil {
		t.Fatal(err)
	}
	module := tag("horos/decl/module", data.EncodeInt(big.NewInt(4)), data.MustEncodeText("installed.pagination"), imports, members)
	tree, err := resolve.Resolve(module, nil, "installed-unified")
	if err != nil {
		t.Fatal(err)
	}
	model, err := resolve.Model(tree)
	if err != nil {
		t.Fatal(err)
	}
	id := decl.DeclarationID{Module: model.Root(), Name: "Page"}
	dataID := decl.DeclarationID{Module: model.Root(), Name: "DataPage"}
	apply := func(id decl.DeclarationID, argument core.Value) (decl.TypeValue, error) {
		callee, err := model.Value(id)
		if err != nil {
			return decl.TypeValue{}, err
		}
		value, err := model.Proper(argument)
		if err != nil {
			return decl.TypeValue{}, err
		}
		return model.Apply(callee, []decl.TypeValue{value})
	}
	wireArgument := tag("horos/wire", text)
	for _, argument := range []core.Value{text, wireArgument} {
		value, err := apply(id, argument)
		if err != nil {
			t.Fatal(err)
		}
		specialized, err := model.Descriptor(value)
		if err != nil || !core.Equal(specialized, pageOf(argument)) {
			t.Fatalf("same public generic binding: %v", err)
		}
	}
	value, err := apply(dataID, text)
	if err != nil {
		t.Fatal(err)
	}
	specialized, err := model.GroundDescriptor(value)
	if err != nil {
		t.Fatal(err)
	}
	instance, err := data.EncodeMap([][2]core.Value{{data.MustEncodeText("items"), data.EncodeList([]core.Value{data.MustEncodeText("Ada")})}})
	if err != nil {
		t.Fatal(err)
	}
	if !horos.Conforms(instance, specialized) {
		t.Fatal("public Ontos instance did not conform to the unchanged ground descriptor")
	}
	for _, argument := range []core.Value{wireArgument, tag("horos/list-of", wireArgument)} {
		_, err := apply(dataID, argument)
		var detail *decl.Error
		if !errors.As(err, &detail) || detail.Code != decl.BoundViolation {
			t.Fatalf("data-only alias accepted capability argument: %v", err)
		}
	}
	declaration, ok := model.Declaration(id)
	if !ok || !core.Equal(declaration.Expression().Value(), pageBody) || !core.Equal(declaration.Parameters()[0].Kind.Value(), top) {
		t.Fatal("specialization changed the authored generic")
	}
}
