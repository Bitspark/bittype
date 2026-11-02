package integration_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/engine/go"
	provider "github.com/Bitspark/bittype/provider/horos/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// A test-only capability catalogue qualifies every AST-only provider explicitly.
// It is not imported by any production registry or language package.
func astOnlyConformanceCases() []struct { directory string; qualify func(*testing.T) } {
	return []struct { directory string; qualify func(*testing.T) }{{"horos", qualifyHoros}}
}

func TestASTOnlyProviders(t *testing.T) {
	for _, fixture := range astOnlyConformanceCases() { t.Run(fixture.directory,fixture.qualify) }
}

func qualifyHoros(t *testing.T) {
	provider:=provider.Provider{}
	if provider.ID()!="horos" || provider.Name()!="Horos" { t.Fatal("unstable Horos identity") }
	if _,ok:=any(provider).(text.ParserProvider[h.Language]); ok { t.Fatal("AST-only provider advertises parsing") }
	if _,ok:=any(provider).(text.FormatterProvider[h.Language]); ok { t.Fatal("AST-only provider advertises formatting") }
	// Independent consumer: a generic list and its explicit text specialization.
	atom:=func(s string) core.Value { return core.NewAtom([]byte(s)) }
	tuple:=func(tag string,args ...core.Value) core.Value { return core.NewTuple(append([]core.Value{atom(tag)},args...)...) }
	txt:=data.MustEncodeText
	textType:=tuple("horos/scalar",atom("utf8-text"))
	parameter:=tuple("horos/decl/param",data.EncodeInt(big.NewInt(0)))
	page:=tuple("horos/decl/type",data.EncodeList([]core.Value{txt("T")}),tuple("horos/list-of",parameter))
	applied:=tuple("horos/decl/type",data.EncodeList(nil),tuple("horos/decl/apply",tuple("horos/decl/local",txt("Page")),data.EncodeList([]core.Value{textType})))
	declarations,err:=data.EncodeMap([][2]core.Value{{txt("Page"),page},{txt("TextPage"),applied}})
	if err!=nil { t.Fatal(err) }
	empty,_:=data.EncodeMap(nil)
	module:=tuple("horos/decl/module",data.EncodeInt(big.NewInt(1)),txt("consumer.pagination"),empty,declarations)
	tree,err:=resolve.Resolve(module,nil,"consumer")
	if err!=nil { t.Fatal(err) }
	definition,err:=representation.NewLanguage[h.Language](provider,h.Profile)
	if err!=nil { t.Fatal(err) }
	envelope,err:=representation.Pack[h.Language](definition,tree)
	if err!=nil { t.Fatal(err) }
	tree.Nodes[tree.Root-1].Kind="broken" // caller storage is not the envelope
	var registry engine.Engine
	if err:=registry.RegisterLanguage(definition); err!=nil { t.Fatal(err) }
	output,err:=registry.TransformTo(envelope,provider.ID())
	if err!=nil { t.Fatal(err) }
	checked,err:=representation.Unpack[h.Language](definition,output)
	if err!=nil { t.Fatal(err) }
	model,err:=resolve.Model(checked)
	if err!=nil { t.Fatal(err) }
	specialized,err:=model.Specialize(decl.DeclarationID{Module:model.Root(),Name:"TextPage"},nil)
	if err!=nil || !core.Equal(specialized,tuple("horos/list-of",textType)) { t.Fatalf("specialization through registered AST: %v",err) }
	checked.Nodes[checked.Root-1].Kind="mutated"
	if fresh,err:=representation.Unpack[h.Language](definition,output); err!=nil || fresh.Nodes[fresh.Root-1].Kind!=h.KindProgram { t.Fatal("unpacked AST aliases envelope") }
	if _,err:=registry.Parse(provider.ID(),text.Input{}); !errors.Is(err,engine.ErrMissingCapability) { t.Fatalf("AST-only parse: %v",err) }
	if _,err:=registry.Render(envelope); !errors.Is(err,engine.ErrMissingCapability) { t.Fatalf("AST-only render: %v",err) }
}
