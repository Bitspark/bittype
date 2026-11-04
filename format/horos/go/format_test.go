package horosformat_test

import (
	"errors"
	"math/big"
	"slices"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	format "github.com/Bitspark/bittype/format/horos/go"
	provider "github.com/Bitspark/bittype/provider/horos/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

func atom(s string) core.Value { return core.NewAtom([]byte(s)) }
func tag(name string, args ...core.Value) core.Value {
	return core.NewTuple(append([]core.Value{atom(name)}, args...)...)
}
func h(name string, args ...core.Value) core.Value { return tag("horos/"+name, args...) }
func d(name string, args ...core.Value) core.Value { return tag("horos/decl/"+name, args...) }
func integer(n int64) core.Value                   { return data.EncodeInt(big.NewInt(n)) }
func entries(t *testing.T, pairs ...[2]core.Value) core.Value {
	t.Helper()
	value, err := data.EncodeMap(pairs)
	if err != nil {
		t.Fatal(err)
	}
	return value
}
func definition(body core.Value, parameters ...core.Value) core.Value {
	return d("type", data.EncodeList(parameters), atom("type"), body)
}
func apply(name string, arguments ...core.Value) core.Value {
	return d("apply", d("local", data.MustEncodeText(name)), data.EncodeList(arguments))
}
func field(name, presence string, typ core.Value) core.Value {
	return h("field", data.MustEncodeText(name), atom(presence), typ)
}

func TestCanonicalPaginationAndImportedLibraryRoundTrip(t *testing.T) {
	txt := data.MustEncodeText
	textType := h("scalar", atom("utf8-text"))
	slot := d("param", integer(0))
	pageBody := h("record", atom("closed"), field("items", "req", h("list-of", slot)), field("next", "opt", textType))
	page := d("module", integer(3), txt("example.pagination"), entries(t), entries(t,
		[2]core.Value{txt("Page"), definition(pageBody, d("parameter", txt("T"), d("proper", atom("type"))))},
		[2]core.Value{txt("DataPage"), definition(apply("Page", slot), d("parameter", txt("D"), d("proper", atom("data"))))},
	))
	pagePin, err := decl.Pin(page)
	if err != nil {
		t.Fatal(err)
	}
	cell := d("module", integer(3), txt("example.cell"), entries(t), entries(t,
		[2]core.Value{txt("Cell"), definition(h("wire", slot), d("parameter", txt("T"), d("proper", atom("type"))))},
	))
	cellPin, err := decl.Pin(cell)
	if err != nil {
		t.Fatal(err)
	}
	imported := func(alias, name string, args ...core.Value) core.Value {
		return d("apply", d("imported", txt(alias), txt(name)), data.EncodeList(args))
	}
	personBody := h("record", atom("closed"), field("name", "req", textType))
	person := d("local", txt("Person"))
	people := d("module", integer(3), txt("example.people"), entries(t,
		[2]core.Value{txt("pages"), d("import", txt("example.pagination"), core.NewAtom(pagePin[:]))},
		[2]core.Value{txt("cells"), d("import", txt("example.cell"), core.NewAtom(cellPin[:]))},
	), entries(t,
		[2]core.Value{txt("Person"), definition(personBody)},
		[2]core.Value{txt("PeoplePage"), definition(imported("pages", "Page", person))},
		[2]core.Value{txt("StoredPeople"), definition(imported("pages", "DataPage", person))},
		[2]core.Value{txt("EditablePeople"), definition(imported("pages", "Page", imported("cells", "Cell", person)))},
	))
	canonical, err := decl.Resolve(people, []core.Value{page, cell})
	if err != nil {
		t.Fatal(err)
	}
	projected, err := resolve.Project(canonical, "canonical")
	if err != nil {
		t.Fatal(err)
	}
	p := provider.Provider{}
	formatted, err := p.Format(projected)
	if err != nil {
		t.Fatal(err)
	}
	if projected.ID != formatted.Abstract.ID || projected.Root != formatted.Abstract.Root || len(projected.Nodes) != len(formatted.Abstract.Nodes) {
		t.Fatal("formatter did not retain the input AST snapshot")
	}
	for i, want := range projected.Nodes {
		got := formatted.Abstract.Nodes[i]
		if want.Kind != got.Kind || !slices.Equal(want.Attributes, got.Attributes) || !slices.Equal(want.Children, got.Children) || !slices.Equal(want.References, got.References) || !slices.Equal(want.Origin.Spans, got.Origin.Spans) || !slices.Equal(want.Origin.Nodes, got.Origin.Nodes) {
			t.Fatalf("formatter lost snapshot node %d", i+1)
		}
	}
	source, err := p.Print(formatted)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(source, "module example.people;") || !strings.Contains(source, "pages.Page<cells.Cell<Person>>") {
		t.Fatalf("entry or imported application lost:\n%s", source)
	}
	parsed, err := p.Parse(parse.Input{SourceID: "bundle.horos", TreeID: "source", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := resolve.Model(parsed)
	if err != nil {
		t.Fatal(err)
	}
	if canonical.Root() != recovered.Root() || len(canonical.Modules()) != len(recovered.Modules()) {
		t.Fatal("module closure identity changed")
	}
	for _, module := range canonical.Modules() {
		actual, ok := recovered.Module(module.Pin())
		if !ok || !core.Equal(module.Value(), actual.Value()) {
			t.Fatal("canonical module value changed")
		}
	}
	want := h("record", atom("closed"), field("items", "req", h("list-of", personBody)), field("next", "opt", textType))
	for _, name := range []string{"PeoplePage", "StoredPeople"} {
		value, err := recovered.DataDescriptor(decl.DeclarationID{Module: recovered.Root(), Name: name}, nil)
		if err != nil || !core.Equal(value, want) {
			t.Fatalf("%s specialization changed: %v", name, err)
		}
	}
	if _, err := recovered.DataDescriptor(decl.DeclarationID{Module: recovered.Root(), Name: "EditablePeople"}, nil); err == nil {
		t.Fatal("Cell capability was silently made data")
	}
	// Optional fields and closed record structure survive, not merely a list alias.
	if !strings.Contains(source, "next?: text;") || !strings.Contains(source, "DataPage<D: Data> = Page<D>") {
		t.Fatal("pagination syntax lost")
	}
}

func TestCanonicalConstructorInterfacesRoundTrip(t *testing.T) {
	txt := data.MustEncodeText
	proper := d("proper", atom("type"))
	dataOnly := d("proper", atom("data"))
	unary := d("constructor", data.EncodeList([]core.Value{proper}), atom("type"))
	higher := d("constructor", data.EncodeList([]core.Value{unary, dataOnly}), atom("data"))
	pair := d("constructor", data.EncodeList([]core.Value{proper, dataOnly}), atom("type"))
	unaryHigher := d("constructor", data.EncodeList([]core.Value{unary}), atom("type"))
	slot := func(index int64) core.Value { return d("param", integer(index)) }
	application := func(indices ...int64) core.Value {
		var arguments []core.Value
		for _, index := range indices {
			arguments = append(arguments, slot(index))
		}
		return d("apply", slot(0), data.EncodeList(arguments))
	}
	module := d("module", integer(3), txt("example.interfaces"), entries(t), entries(t,
		[2]core.Value{txt("Apply"), definition(application(1),
			d("parameter", txt("F"), unary), d("parameter", txt("T"), proper))},
		[2]core.Value{txt("Higher"), d("type", data.EncodeList([]core.Value{
			d("parameter", txt("H"), higher), d("parameter", txt("F"), unary), d("parameter", txt("D"), dataOnly),
		}), atom("data"), application(1, 2))},
		[2]core.Value{txt("Pair"), definition(application(1, 2),
			d("parameter", txt("F"), pair), d("parameter", txt("T"), proper), d("parameter", txt("D"), dataOnly))},
		[2]core.Value{txt("UnaryHigher"), definition(application(1),
			d("parameter", txt("H"), unaryHigher), d("parameter", txt("F"), unary))},
	))
	projected, err := resolve.Resolve(module, nil, "canonical-interfaces")
	if err != nil {
		t.Fatal(err)
	}
	p := provider.Provider{}
	formatted, err := p.Format(projected)
	if err != nil {
		t.Fatal(err)
	}
	source, err := p.Print(formatted)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"type Apply<F: Type -> Type, T> = F<T>;",
		"type Higher<H: (Type -> Type, Data) -> Data, F: Type -> Type, D: Data>: Data = H<F, D>;",
		"type Pair<F: (Type, Data) -> Type, T, D: Data> = F<T, D>;",
		"type UnaryHigher<H: (Type -> Type) -> Type, F: Type -> Type> = H<F>;",
	} {
		if !strings.Contains(source, want) {
			t.Fatalf("missing expected interface spelling %q:\n%s", want, source)
		}
	}
	parsed, err := p.Parse(parse.Input{SourceID: "interfaces.horos", TreeID: "interfaces", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := resolve.Model(parsed)
	if err != nil {
		t.Fatal(err)
	}
	actual, ok := recovered.Module(recovered.Root())
	if !ok || !core.Equal(actual.Value(), module) {
		t.Fatal("canonical interface trees, result claims or parameter callees changed")
	}
}

func TestFormatterRejectsUnspellableASTWithoutPartialOutput(t *testing.T) {
	p := provider.Provider{}
	parsed, err := p.Parse(parse.Input{SourceID: "small", TreeID: "small", Text: "module small; type Identity<T> = T;"})
	if err != nil {
		t.Fatal(err)
	}
	for _, change := range []func(*s.AST){
		func(tree *s.AST) { tree.Nodes[tree.Root-1].Kind = "unknown" },
		func(tree *s.AST) {
			for i := range tree.Nodes {
				for j := range tree.Nodes[i].Attributes {
					if tree.Nodes[i].Attributes[j].Name == s.AttributePin {
						tree.Nodes[i].Attributes[j].Text = strings.Repeat("0", 64)
					}
				}
			}
		},
	} {
		tree := parsed.Abstract
		tree.Nodes = append([]s.Node(nil), tree.Nodes...)
		for i := range tree.Nodes {
			tree.Nodes[i].Attributes = append([]ast.Attribute(nil), tree.Nodes[i].Attributes...)
		}
		change(&tree)
		out, err := format.Format(tree)
		if err == nil || len(out.Nodes) != 0 || len(out.Abstract.Nodes) != 0 {
			t.Fatal("unspellable AST returned partial output")
		}
	}
	var nilTree *s.AST
	for _, view := range []ast.View[s.Language]{nil, nilTree} {
		if _, err := format.Format(view); err == nil {
			t.Fatal("nil view admitted")
		}
	}
}

func TestFormatterBoundsLexicalExpansion(t *testing.T) {
	// One AST attribute and concrete leaf can expand into more than 200,000
	// lexical tokens. Counting CST leaves alone would emit unparseable source.
	module := d("module", integer(3), data.MustEncodeText(strings.Repeat("a.", 100000)+"a"), entries(t), entries(t))
	tree, err := resolve.Resolve(module, nil, "many-segments")
	if err != nil {
		t.Fatal(err)
	}
	out, err := format.Format(tree)
	var diagnostic *text.Diagnostic
	if !errors.As(err, &diagnostic) || diagnostic.Code != text.LimitExceeded || len(out.Nodes) != 0 || len(out.Abstract.Nodes) != 0 {
		t.Fatalf("lexical expansion must fail without partial output: %v", err)
	}
}
