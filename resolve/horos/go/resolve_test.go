package resolve_test

import (
	"errors"
	"math/big"
	"reflect"
	"slices"
	"strings"
	"testing"

	structure "github.com/Bitspark/bittype/check/ast/go"
	check "github.com/Bitspark/bittype/check/horos/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

type value = core.Value

func atom(v string) value   { return core.NewAtom([]byte(v)) }
func txt(v string) value    { return data.MustEncodeText(v) }
func integer(v int64) value { return data.EncodeInt(big.NewInt(v)) }
func form(prefix, tag string, values ...value) value {
	return core.NewTuple(append([]value{atom(prefix + tag)}, values...)...)
}
func h(tag string, values ...value) value { return form("horos/", tag, values...) }
func d(tag string, values ...value) value { return form("horos/decl/", tag, values...) }
func mapping(entries ...[2]value) value {
	v, err := data.EncodeMap(entries)
	if err != nil {
		panic(err)
	}
	return v
}
func entry(name string, v value) [2]value { return [2]value{txt(name), v} }
func definition(body value, names ...string) value {
	var parameters []value
	for _, name := range names {
		parameters = append(parameters, d("parameter", txt(name), d("top", d("sort-proper"))))
	}
	return d("type", data.EncodeList(parameters), d("top", d("sort-proper")), body)
}
func module(name string, imports value, declarations ...[2]value) value {
	return d("module", integer(4), txt(name), imports, mapping(declarations...))
}
func param(i int64) value { return d("param", integer(i)) }
func local(name string, args ...value) value {
	if len(args) == 0 { return d("local", txt(name)) }
	return d("apply", d("local", txt(name)), data.EncodeList(args))
}
func page(element value) value {
	return h("record", atom("closed"),
		h("field", txt("items"), atom("req"), h("list-of", element)),
		h("field", txt("next"), atom("opt"), h("scalar", atom("utf8-text"))))
}
func person() value {
	return h("record", atom("closed"), h("field", txt("name"), atom("req"), h("scalar", atom("utf8-text"))))
}

func pagination(t *testing.T) (value, value) {
	t.Helper()
	pages := module("example.pagination", mapping(),
		entry("Page", definition(page(param(0)), "T")),
		entry("Person", definition(h("scalar", atom("int"))))) // same name in callee scope
	pin, err := decl.Pin(pages)
	if err != nil {
		t.Fatal(err)
	}
	imported := d("import", txt("example.pagination"), core.NewAtom(pin[:]))
	application := d("apply", d("imported", txt("pagination"), txt("Page")), data.EncodeList([]value{local("Person")}))
	people := module("example.people", mapping(entry("pagination", imported)),
		entry("Person", definition(person())), entry("PeoplePage", definition(application)))
	return people, pages
}

func fixture(t *testing.T) s.AST {
	t.Helper()
	root, dependency := pagination(t)
	tree, err := resolve.Resolve(root, []value{dependency}, "pagination")
	if err != nil {
		t.Fatal(err)
	}
	return tree
}
func attr(node s.Node, name string) string {
	for _, a := range node.Attributes {
		if a.Name == name {
			return a.Text
		}
	}
	return ""
}
func setAttr(node *s.Node, name, value string) {
	for i := range node.Attributes {
		if node.Attributes[i].Name == name {
			node.Attributes[i].Text = value
			return
		}
	}
	node.Attributes = append(node.Attributes, ast.Attribute{Name: name, Text: value})
}
func find(tree s.AST, kind s.Kind, name string) ast.NodeID {
	for i, node := range tree.Nodes {
		if node.Kind == kind && (name == "" || attr(node, s.AttributeName) == name) {
			return ast.NodeID(i + 1)
		}
	}
	panic("missing fixture node")
}
func target(node s.Node, role string) ast.NodeID {
	for _, ref := range node.References {
		if ref.Role == role {
			return ref.Target.Node
		}
	}
	return 0
}

func materialize(program *decl.Program, id decl.DeclarationID) (core.Value, error) {
	value, err := program.Value(id)
	if err != nil { return nil, err }
	return program.Descriptor(value)
}

func TestPaginationBindingsAndSpecialization(t *testing.T) {
	tree := fixture(t)
	if err := check.Check(tree); err != nil {
		t.Fatal(err)
	}
	pageID := find(tree, s.KindTypeDefinition, "Page")
	parameter := find(tree, s.KindTypeParameter, "T")
	use := find(tree, s.KindParameter, "")
	if target(tree.Nodes[use-1], s.RoleTarget) != parameter {
		t.Fatal("slot lost its binder")
	}
	application := find(tree, s.KindReference, "Page")
	if target(tree.Nodes[application-1], s.RoleTarget) != pageID || target(tree.Nodes[application-1], s.RoleImport) == 0 {
		t.Fatal("imported application lost declaration/import identity")
	}
	model, err := resolve.Model(tree)
	if err != nil {
		t.Fatal(err)
	}
	root, _ := pagination(t)
	stored, _ := model.Module(model.Root())
	if !core.Equal(stored.Value(), root) {
		t.Fatal("canonical value changed during AST round trip")
	}
	result, err := materialize(model, decl.DeclarationID{Module: model.Root(), Name: "PeoplePage"})
	if err != nil || !core.Equal(result, page(person())) {
		t.Fatalf("caller binding or specialization changed: %v", err)
	}
	if tree.Nodes[use-1].Kind != s.KindParameter {
		t.Fatal("projection monomorphized the open definition")
	}
}

func TestNestedGenericBinding(t *testing.T) {
	root := module("example.nested", mapping(),
		entry("Page", definition(page(param(0)), "T")),
		entry("Nested", definition(local("Page", local("Page", param(0))), "T")),
		entry("Applied", definition(local("Nested", person()))))
	tree, err := resolve.Resolve(root, nil, "nested")
	if err != nil {
		t.Fatal(err)
	}
	model, err := resolve.Model(tree)
	if err != nil {
		t.Fatal(err)
	}
	result, err := materialize(model, decl.DeclarationID{Module: model.Root(), Name: "Applied"})
	if err != nil || !core.Equal(result, page(page(person()))) {
		t.Fatalf("nested binding: %v", err)
	}
	binders := map[ast.NodeID]bool{}
	for _, node := range tree.Nodes {
		if node.Kind == s.KindParameter {
			binders[target(node, s.RoleTarget)] = true
		}
	}
	if len(binders) != 2 || binders[0] {
		t.Fatal("same-spelled T binders were collapsed")
	}
}

func TestExactLiteralIsOpaque(t *testing.T) {
	literal := h("exact", core.NewTuple(param(99), local("Missing")))
	root := module("example.literal", mapping(), entry("Literal", definition(literal)))
	tree, err := resolve.Resolve(root, nil, "literal")
	if err != nil {
		t.Fatal(err)
	}
	for _, node := range tree.Nodes {
		if node.Kind == s.KindParameter || node.Kind == s.KindApplication {
			t.Fatal("literal became declaration syntax")
		}
	}
	model, err := resolve.Model(tree)
	if err != nil {
		t.Fatal(err)
	}
	result, err := materialize(model, decl.DeclarationID{Module: model.Root(), Name: "Literal"})
	if err != nil || !core.Equal(result, literal) {
		t.Fatalf("literal changed: %v", err)
	}
}

func TestModelAcceptsAbsentBindingAnnotations(t *testing.T) {
	tree := fixture(t)
	for i := range tree.Nodes {
		if tree.Nodes[i].Kind != s.KindProgram {
			tree.Nodes[i].References = nil
		}
	}
	before := clone(tree)
	model, err := resolve.Model(tree)
	if err != nil || model == nil {
		t.Fatalf("unannotated syntax did not resolve: %v", err)
	}
	if !reflect.DeepEqual(tree, before) {
		t.Fatal("resolution wrote annotations into borrowed syntax")
	}
}

type changingView struct {
	first, second s.AST
	calls         int
}

func (v *changingView) ASTView() s.AST {
	v.calls++
	if v.calls == 1 {
		return v.first
	}
	return v.second
}

func TestModelCapturesAnExternalViewOnce(t *testing.T) {
	valid := fixture(t)
	view := &changingView{first: valid, second: s.AST{}}
	if model, err := resolve.Model(view); err != nil || model == nil || view.calls != 1 {
		t.Fatalf("changed view reached reconstruction: calls=%d, error=%v", view.calls, err)
	}
	view = &changingView{first: s.AST{}, second: valid}
	if model, err := resolve.Model(view); err == nil || model != nil || view.calls != 1 {
		t.Fatalf("invalid first capture was replaced: calls=%d, error=%v", view.calls, err)
	}
}

func TestInvalidASTsFailWithoutRepair(t *testing.T) {
	cases := map[string]func(s.AST){
		"wrong parameter binding": func(a s.AST) {
			i := find(a, s.KindParameter, "")
			a.Nodes[i-1].References[0].Target.Node = find(a, s.KindTypeDefinition, "Page")
		},
		"wrong application binding": func(a s.AST) {
			i := find(a, s.KindReference, "Page")
			a.Nodes[i-1].References[0].Target.Node = find(a, s.KindTypeDefinition, "PeoplePage")
		},
		"external binding": func(a s.AST) { i := find(a, s.KindParameter, ""); a.Nodes[i-1].References[0].Target.Tree = "external" },
		"wrong module pin": func(a s.AST) {
			i := find(a, s.KindModule, "")
			setAttr(&a.Nodes[i-1], s.AttributePin, strings.Repeat("0", 64))
		},
		"wrong import pin":   func(a s.AST) { i := find(a, s.KindImport, ""); setAttr(&a.Nodes[i-1], s.AttributePin, "ff") },
		"wrong index":        func(a s.AST) { i := find(a, s.KindParameter, ""); setAttr(&a.Nodes[i-1], s.AttributeIndex, "1") },
		"noncanonical index": func(a s.AST) { i := find(a, s.KindParameter, ""); setAttr(&a.Nodes[i-1], s.AttributeIndex, "00") },
		"extra attribute": func(a s.AST) {
			a.Nodes[0].Attributes = append(a.Nodes[0].Attributes, ast.Attribute{Name: "extra", Text: "yes"})
		},
		"unknown kind":     func(a s.AST) { a.Nodes[0].Kind = "unknown" },
		"wrong child role": func(a s.AST) { a.Nodes[0].Children[0].Role = "unknown" },
		"invalid hex":      func(a s.AST) { i := find(a, s.KindAtom, ""); setAttr(&a.Nodes[i-1], s.AttributeHex, "xy") },
		"invalid UTF-8": func(a s.AST) {
			i := find(a, s.KindModule, "")
			setAttr(&a.Nodes[i-1], s.AttributeName, string([]byte{255}))
		},
		"literal replacing slot": func(a s.AST) {
			i := find(a, s.KindParameter, "")
			a.Nodes[i-1] = s.Node{Kind: s.KindAtom, Attributes: []ast.Attribute{{Name: s.AttributeHex, Text: "00"}}}
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			tree := fixture(t)
			mutate(tree)
			before := clone(tree)
			if model, err := resolve.Model(tree); err == nil || model != nil {
				t.Fatal("malformed AST accepted or partial model returned")
			}
			if !reflect.DeepEqual(before, tree) {
				t.Fatal("checker mutated its input")
			}
		})
	}
}

func clone(tree s.AST) s.AST {
	result := tree
	result.Nodes = slices.Clone(tree.Nodes)
	for i, node := range tree.Nodes {
		result.Nodes[i].Attributes = slices.Clone(node.Attributes)
		result.Nodes[i].Children = slices.Clone(node.Children)
		result.Nodes[i].References = slices.Clone(node.References)
		result.Nodes[i].Origin.Spans = slices.Clone(node.Origin.Spans)
		result.Nodes[i].Origin.Nodes = slices.Clone(node.Origin.Nodes)
	}
	return result
}

func TestArenaAddressesAndOriginsAreNotPins(t *testing.T) {
	tree := fixture(t)
	// Reverse the entire arena while preserving containment and reference edges.
	renumber := func(id ast.NodeID) ast.NodeID { return ast.NodeID(len(tree.Nodes)) + 1 - id }
	result := clone(tree)
	result.Root = renumber(tree.Root)
	for i, node := range tree.Nodes {
		n := clone(s.AST{Nodes: []s.Node{node}}).Nodes[0]
		for j := range n.Children {
			n.Children[j].Node = renumber(n.Children[j].Node)
		}
		for j := range n.References {
			n.References[j].Target.Node = renumber(n.References[j].Target.Node)
		}
		n.Origin.Spans = []ast.Span{{Source: "source", Start: 1, End: 2}}
		result.Nodes[len(tree.Nodes)-1-i] = n
	}
	before := clone(result)
	if err := check.Check(result); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, before) {
		t.Fatal("checker rewrote addresses or provenance")
	}
}

func TestAdmissionAndProjectionLimits(t *testing.T) {
	var typedNil *s.AST
	for _, view := range []ast.View[s.Language]{nil, typedNil, s.AST{}} {
		if model, err := resolve.Model(view); err == nil || model != nil {
			t.Fatal("invalid view accepted")
		}
	}
	for _, model := range []*decl.Program{nil, {}} {
		if tree, err := resolve.Project(model, "empty"); err == nil || len(tree.Nodes) != 0 {
			t.Fatal("invalid model accepted")
		}
	}
	root, _ := pagination(t)
	if tree, err := resolve.Resolve(root, nil, "missing"); err == nil || len(tree.Nodes) != 0 {
		t.Fatal("missing dependency accepted")
	}
	bigLiteral := module("example.large", mapping(), entry("Large", definition(h("exact", core.NewAtom(make([]byte, 9<<20))))))
	program, err := decl.Resolve(bigLiteral, nil)
	if err != nil {
		t.Fatal(err)
	}
	tree, err := resolve.Project(program, "large")
	var diagnostic *text.Diagnostic
	if len(tree.Nodes) != 0 || !errors.As(err, &diagnostic) || diagnostic.Code != text.LimitExceeded {
		t.Fatalf("hex growth must fail boundedly: %v", err)
	}
	tree = fixture(t)
	tree.Nodes[0].Attributes = []ast.Attribute{{Name: "large", Text: strings.Repeat("x", structure.MaxTextBytes)}}
	if err := check.Check(tree); !errors.As(err, &diagnostic) || diagnostic.Code != text.LimitExceeded {
		t.Fatalf("oversize AST: %v", err)
	}
}

func TestProjectionCountsReferenceTextAndDepth(t *testing.T) {
	empty, err := decl.Resolve(module("example.empty", mapping()), nil)
	if err != nil {
		t.Fatal(err)
	}
	// The ID fits once, but its repeated occurrence on the root reference does
	// not. Even an otherwise empty program must return no partial AST.
	tree, err := resolve.Project(empty, ast.TreeID(strings.Repeat("x", structure.MaxTextBytes/2)))
	var diagnostic *text.Diagnostic
	if len(tree.Nodes) != 0 || !errors.As(err, &diagnostic) || diagnostic.Code != text.LimitExceeded {
		t.Fatalf("reference text admission: %v", err)
	}
	var literal value = atom("leaf")
	for i := 0; i < structure.MaxDepth+1; i++ {
		literal = core.NewTuple(literal)
	}
	limits := decl.DefaultLimits()
	limits.MaxDepth = 512
	deep, err := decl.ResolveWithLimits(module("example.deep", mapping(), entry("Deep", definition(h("exact", literal)))), nil, limits)
	if err != nil {
		t.Fatal(err)
	}
	tree, err = resolve.Project(deep, "deep")
	var modelError *decl.Error
	if len(tree.Nodes) != 0 || !errors.As(err, &modelError) || modelError.Code != decl.Exhausted {
		t.Fatalf("canonical depth admission: %v", err)
	}
}

func TestProjectionAndCheckerShareCanonicalLimits(t *testing.T) {
	// Canonical module wrappers and AST wrappers have different depths. Near
	// this seam an AST can fit even when its canonical module exceeds defaults.
	for wrappers := 249; wrappers <= 253; wrappers++ {
		var literal value = atom("leaf")
		for i := 0; i < wrappers; i++ {
			literal = core.NewTuple(literal)
		}
		root := module("example.seam", mapping(), entry("Seam", definition(h("exact", literal))))
		limits := decl.DefaultLimits()
		limits.MaxDepth = 512
		custom, err := decl.ResolveWithLimits(root, nil, limits)
		if err != nil {
			t.Fatal(err)
		}
		_, admission := decl.Resolve(root, nil)
		tree, err := resolve.Project(custom, "seam")
		if admission != nil {
			if err == nil || len(tree.Nodes) != 0 {
				t.Fatalf("%d wrappers: Project admitted what Model cannot check", wrappers)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%d wrappers: unexpected projection failure: %v", wrappers, err)
		}
		if err := check.Check(tree); err != nil {
			t.Fatalf("%d wrappers: projected AST fails its own profile: %v", wrappers, err)
		}
		if _, err := resolve.Model(tree); err != nil {
			t.Fatalf("%d wrappers: projected AST has no checked model: %v", wrappers, err)
		}
	}
}
