package analyze_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"reflect"
	"slices"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/horos/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	provider "github.com/Bitspark/bittype/provider/horos/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	codec "github.com/bitspark/ontos/codec/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

type value = core.Value

func atom(v string) value { return core.NewAtom([]byte(v)) }
func txt(v string) value { return data.MustEncodeText(v) }
func integer(v int64) value { return data.EncodeInt(big.NewInt(v)) }
func form(prefix, tag string, values ...value) value {
	return core.NewTuple(append([]value{atom(prefix + tag)}, values...)...)
}
func h(tag string, values ...value) value { return form("horos/", tag, values...) }
func d(tag string, values ...value) value { return form("horos/decl/", tag, values...) }
func mapping(entries ...[2]value) value {
	v, err := data.EncodeMap(entries)
	if err != nil { panic(err) }
	return v
}
func entry(name string, v value) [2]value { return [2]value{txt(name), v} }
func definition(body value, names ...string) value {
	var parameters []value
	for _, name := range names { parameters = append(parameters, txt(name)) }
	return d("type", data.EncodeList(parameters), body)
}
func module(name string, imports value, declarations ...[2]value) value {
	return d("module", integer(1), txt(name), imports, mapping(declarations...))
}
func param(index int64) value { return d("param", integer(index)) }
func local(name string, args ...value) value {
	return d("apply", d("local", txt(name)), data.EncodeList(args))
}
func imported(alias, name string, args ...value) value {
	return d("apply", d("imported", txt(alias), txt(name)), data.EncodeList(args))
}
func page(element value) value {
	return h("record", atom("closed"),
		h("field", txt("items"), atom("req"), h("list-of", element)),
		h("field", txt("next"), atom("opt"), h("scalar", atom("utf8-text"))))
}
func person() value {
	return h("record", atom("closed"), h("field", txt("name"), atom("req"), h("scalar", atom("utf8-text"))))
}
func pin(t *testing.T, value value) decl.Digest {
	t.Helper()
	p, err := decl.Pin(value)
	if err != nil { t.Fatal(err) }
	return p
}

// The options change canonical source as well as its handwritten AST mutations
// below, so diagnostics concern real unresolved source, never a stale module pin.
type fixtureOptions struct { importedName, expectedModule, alias string; parameter int64 }

func pagination(t *testing.T, options fixtureOptions) (value, value) {
	t.Helper()
	if options.importedName == "" { options.importedName = "Page" }
	if options.expectedModule == "" { options.expectedModule = "example.pagination" }
	if options.alias == "" { options.alias = "pagination" }
	pages := module("example.pagination", mapping(),
		entry("Page", definition(page(param(0)), "T")),
		entry("Person", definition(h("scalar", atom("int")))))
	pagesPin := pin(t, pages)
	dependency := d("import", txt(options.expectedModule), core.NewAtom(pagesPin[:]))
	people := module("example.people", mapping(entry("pagination", dependency)),
		entry("Person", definition(person())),
		entry("Outer", definition(imported(options.alias, options.importedName, param(options.parameter)), "T")),
		entry("PeoplePage", definition(local("Outer", local("Person")))))
	return people, pages
}

func fixture(t *testing.T) s.AST {
	t.Helper()
	root, dependency := pagination(t, fixtureOptions{})
	tree, err := resolve.Resolve(root, []value{dependency}, "pagination-current")
	if err != nil { t.Fatal(err) }
	return tree
}
func attribute(node s.Node, name string) string {
	for _, a := range node.Attributes { if a.Name == name { return a.Text } }
	return ""
}
func setAttribute(node *s.Node, name, value string) {
	for i := range node.Attributes {
		if node.Attributes[i].Name == name { node.Attributes[i].Text = value; return }
	}
	node.Attributes = append(node.Attributes, ast.Attribute{Name: name, Text: value})
}
func find(tree s.AST, kind s.Kind, name string) ast.NodeID {
	for i, node := range tree.Nodes {
		if node.Kind == kind && (name == "" || attribute(node, s.AttributeName) == name) { return ast.NodeID(i + 1) }
	}
	panic("missing fixture node")
}
func body(tree s.AST, declaration ast.NodeID) ast.NodeID {
	for _, child := range tree.Nodes[declaration-1].Children { if child.Role == s.RoleBody { return child.Node } }
	panic("missing declaration body")
}
func stripAnnotations(tree s.AST) {
	for i := range tree.Nodes { if tree.Nodes[i].Kind != s.KindProgram { tree.Nodes[i].References = nil } }
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

func analyzer(t *testing.T, tree s.AST) (analysis.Analyzer, analysis.Input) {
	t.Helper()
	language, err := representation.NewLanguage[s.Language](provider.Provider{}, s.Profile)
	if err != nil { t.Fatal(err) }
	envelope, err := representation.Pack(language, tree)
	if err != nil { t.Fatalf("syntax must reach analysis: %v", err) }
	analyzer, err := analysis.NewAnalyzer[s.Language](language, analyze.Contract(), provider.Provider{})
	if err != nil { t.Fatal(err) }
	return analyzer, analysis.Input{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: []representation.AST{envelope}}
}
func run(t *testing.T, tree s.AST, modify func(*analysis.Input)) analysis.Result {
	t.Helper()
	analyzer, input := analyzer(t, tree)
	if modify != nil { modify(&input) }
	result, err := analyzer.Analyze(context.Background(), input)
	if err != nil { t.Fatal(err) }
	if !result.Valid() { t.Fatal("analysis returned no request-bound result") }
	return result
}
func requireReport(t *testing.T, result analysis.Result, key semantics.Key, status analysis.Completion, verdict analysis.Verdict) {
	t.Helper()
	report, found := result.Report(key)
	if !found || report.Status != status || report.Verdict != verdict { t.Fatalf("report %v: %+v, found=%v", key, report, found) }
}
func requireModel(t *testing.T, result analysis.Result) semantics.Model {
	t.Helper()
	model, ok := result.Model()
	if !ok { t.Fatalf("missing semantic model: %+v", result.Diagnostics()) }
	for _, obligation := range analyze.Contract().Obligations {
		verdict := analysis.Undetermined
		if obligation.Key == analyze.Validity() { verdict = analysis.Holds }
		requireReport(t, result, obligation.Key, analysis.Complete, verdict)
	}
	return model
}
func entity(t *testing.T, model semantics.Model, module, name string) (semantics.EntityID, semantics.Entity) {
	t.Helper()
	for i, e := range model.Data().Entities {
		if e.Kind.Name != "declaration" || e.Name != name { continue }
		owner, ok := model.Entity(e.Owner)
		if ok && owner.Kind.Name == "module" && owner.Name == module { return semantics.EntityID(i + 1), e }
	}
	t.Fatalf("missing declaration %s.%s", module, name)
	return 0, semantics.Entity{}
}
func term(t *testing.T, model semantics.Model, id semantics.TermID) semantics.Term {
	t.Helper()
	term, found := model.Term(id)
	if !found { t.Fatalf("missing term %d", id) }
	return term
}

func TestPaginationBindingsPreserveCallerScopeAndDistinctParameters(t *testing.T) {
	tree := fixture(t)
	stripAnnotations(tree)
	before := clone(tree)
	result := run(t, tree, nil)
	model := requireModel(t, result)
	if !reflect.DeepEqual(tree, before) { t.Fatal("analysis rewrote syntax") }
	pageID, page := entity(t, model, "example.pagination", "Page")
	outerID, outer := entity(t, model, "example.people", "Outer")
	personID, _ := entity(t, model, "example.people", "Person")
	calleePersonID, _ := entity(t, model, "example.pagination", "Person")
	peopleID, people := entity(t, model, "example.people", "PeoplePage")
	if personID == calleePersonID || len(page.Parameters) != 1 || len(outer.Parameters) != 1 || page.Parameters[0].Entity == outer.Parameters[0].Entity { t.Fatal("module or generic binder identities collapsed") }
	outerBody := term(t, model, outer.Body)
	if outerBody.Kind != semantics.Application || outerBody.Target != pageID { t.Fatal("Outer does not apply imported Page") }
	bindings, ok := model.ParameterBindings(outer.Body)
	if !ok || len(bindings) != 1 || bindings[0].Parameter != page.Parameters[0].Entity { t.Fatal("Page parameter binding missing") }
	argument := term(t, model, bindings[0].Argument)
	if argument.Kind != semantics.Reference || argument.Target != outer.Parameters[0].Entity { t.Fatal("Page argument captured the callee's T") }
	peopleBody := term(t, model, people.Body)
	if peopleBody.Kind != semantics.Application || peopleBody.Target != outerID { t.Fatal("PeoplePage lost its Outer application") }
	personApplication := term(t, model, peopleBody.Arguments[0])
	if personApplication.Kind != semantics.Application || personApplication.Target != personID || len(personApplication.Arguments) != 0 { t.Fatal("zero-argument Person application lost caller identity") }
	if bindings, ok := model.ParameterBindings(peopleBody.Arguments[0]); !ok || len(bindings) != 0 { t.Fatal("zero-argument application must retain application status") }
	deps := model.DependenciesFrom(peopleID, analyze.ApplicationDependency())
	targets := map[semantics.EntityID]bool{}
	for _, dep := range deps { targets[dep.To] = true; if dep.At == nil { t.Fatal("dependency has no occurrence") } }
	if len(deps) != 2 || !targets[outerID] || !targets[personID] { t.Fatalf("wrong dependency facts: %+v", deps) }
	for _, binding := range model.Data().Bindings {
		if resolved, ok := model.Resolve(binding.At); !ok || resolved != binding.Target { t.Fatal("published occurrence binding is not queryable") }
		at := binding.At
		at.Role = analyze.ExpressionRole()
		if _, ok := model.TermAt(at); !ok { t.Fatal("binding lost its occurrence-to-term association") }
	}
}

func TestDescriptorPayloadsAndOperandPathsRetainCanonicalMeaning(t *testing.T) {
	bodies := map[string]value{
		"Int": h("scalar", atom("int")), "Text": h("scalar", atom("utf8-text")),
		"ExactA": h("exact", atom("a")), "ExactB": h("exact", atom("b")),
		"Required": h("record", atom("closed"), h("field", txt("value"), atom("req"), h("scalar", atom("int")))),
		"Optional": h("record", atom("closed"), h("field", txt("value"), atom("opt"), h("scalar", atom("int")))),
		"Open": h("record", atom("open"), h("field", txt("value"), atom("req"), h("scalar", atom("int")))),
		"Opaque": h("exact", core.NewTuple(param(99), local("Missing"))),
	}
	var definitions [][2]value
	for name, body := range bodies { definitions = append(definitions, entry(name, definition(body))) }
	tree, err := resolve.Resolve(module("example.payloads", mapping(), definitions...), nil, "payloads")
	if err != nil { t.Fatal(err) }
	model := requireModel(t, run(t, tree, nil))
	payloads := make(map[string][]byte)
	for name, expected := range bodies {
		_, declaration := entity(t, model, "example.payloads", name)
		body := term(t, model, declaration.Body)
		if body.Kind != semantics.ProviderExpression || body.Encoding != (semantics.Key{Owner: "ontos", Name: "codec-v1"}) || !bytes.Equal(body.Data, codec.Encode(expected)) { t.Fatalf("%s lost canonical payload", name) }
		decoded, err := codec.Decode(body.Data)
		if err != nil || !core.Equal(decoded, expected) { t.Fatalf("%s payload cannot round-trip: %v", name, err) }
		payloads[name] = body.Data
		if name == "Required" || name == "Optional" || name == "Open" {
			if len(body.Operands) != 1 || !slices.Equal(body.Operands[0].Locator, []uint32{2, 3}) { t.Fatalf("%s lost its descriptor-relative operand path: %+v", name, body.Operands) }
		} else if len(body.Operands) != 0 { t.Fatalf("literal/scalar %s unexpectedly contains semantic operands", name) }
	}
	for _, pair := range [][2]string{{"Int", "Text"}, {"ExactA", "ExactB"}, {"Required", "Optional"}, {"Required", "Open"}} {
		if bytes.Equal(payloads[pair[0]], payloads[pair[1]]) { t.Fatalf("%v became indistinguishable", pair) }
	}
	if len(model.Data().Bindings) != 0 || len(model.Data().Dependencies) != 0 { t.Fatal("declaration-looking literal tuples were interpreted as types") }
}

func TestCurrentSyntaxSitesSurviveArenaReorderingAndHistoricalOrigins(t *testing.T) {
	tree := fixture(t)
	stripAnnotations(tree)
	old := clone(tree)
	renumber := func(id ast.NodeID) ast.NodeID { return ast.NodeID(len(tree.Nodes)) + 1 - id }
	tree.Root = renumber(old.Root)
	for i, node := range old.Nodes {
		for j := range node.Children { node.Children[j].Node = renumber(node.Children[j].Node) }
		for j := range node.References { node.References[j].Target.Node = renumber(node.References[j].Target.Node) }
		node.Origin.Nodes = []ast.NodeRef{{Tree: "historical", Node: 999}}
		tree.Nodes[len(tree.Nodes)-1-i] = node
	}
	model := requireModel(t, run(t, tree, nil))
	for _, e := range model.Data().Entities {
		for _, site := range e.Declarations {
			if site.Tree != tree.ID || site.Node == 0 || int(site.Node) > len(tree.Nodes) { t.Fatal("entity used historical origin as current address") }
			if attribute(tree.Nodes[site.Node-1], s.AttributeName) != e.Name { t.Fatal("entity provenance points at another current declaration") }
		}
	}
	for _, expression := range model.Data().Expressions {
		if expression.At.Site.Tree != tree.ID { t.Fatal("expression did not refer to captured syntax") }
	}
}

func requireFailure(t *testing.T, result analysis.Result, status analysis.Completion, verdict analysis.Verdict, code, subject string) {
	t.Helper()
	if result.HasModel() { t.Fatal("failed or incomplete analysis published a model") }
	requireReport(t, result, analyze.Validity(), status, verdict)
	for _, key := range []semantics.Key{analyze.Indexing(), analyze.Bindings(), analyze.Applications(), analyze.Dependencies()} { requireReport(t, result, key, analysis.Incomplete, analysis.Undetermined) }
	diagnostics := result.Diagnostics()
	if len(diagnostics) != 1 || diagnostics[0].Code != (semantics.Key{Owner: "horos", Name: code}) || diagnostics[0].Subject != (semantics.Key{Owner: "horos", Name: subject}) { t.Fatalf("wrong diagnostic: %+v", diagnostics) }
}

func TestUnresolvedSourceAndContextHaveDistinctReports(t *testing.T) {
	cases := []struct { name string; options fixtureOptions; mutate func(s.AST); code, subject string }{
		{"unknown declaration", fixtureOptions{importedName: "Missing"}, func(tree s.AST) { setAttribute(&tree.Nodes[find(tree, s.KindApplication, "Page")-1], s.AttributeName, "Missing") }, "unknown-declaration", "source"},
		{"unknown import", fixtureOptions{alias: "missing"}, func(tree s.AST) { setAttribute(&tree.Nodes[find(tree, s.KindApplication, "Page")-1], s.AttributeAlias, "missing") }, "unknown-import", "source"},
		{"out of range parameter", fixtureOptions{parameter: 9}, func(tree s.AST) {
			outer := find(tree, s.KindDeclaration, "Outer")
			application := body(tree, outer)
			setAttribute(&tree.Nodes[tree.Nodes[application-1].Children[0].Node-1], s.AttributeIndex, "9")
		}, "invalid", "source"},
		{"wrong expected module", fixtureOptions{expectedModule: "example.wrong"}, func(tree s.AST) { setAttribute(&tree.Nodes[find(tree, s.KindImport, "")-1], s.AttributeName, "example.wrong") }, "import-mismatch", "context"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			tree := fixture(t)
			stripAnnotations(tree)
			test.mutate(tree)
			root, _ := pagination(t, test.options)
			setAttribute(&tree.Nodes[find(tree, s.KindModule, "example.people")-1], s.AttributePin, pin(t, root).String())
			result := run(t, tree, nil)
			requireFailure(t, result, analysis.Complete, analysis.Rejected, test.code, test.subject)
			if result.Diagnostics()[0].At == nil { t.Fatal("located failure lost its current syntax site") }
		})
	}
}

// Compact keeps only the containment subtree after a test edit. It is a test
// fixture operation, not a parser or a production canonical-value walker.
func compact(tree s.AST) s.AST {
	result := s.AST{ID: tree.ID}
	ids := map[ast.NodeID]ast.NodeID{}
	var copyNode func(ast.NodeID) ast.NodeID
	copyNode = func(id ast.NodeID) ast.NodeID {
		node := clone(s.AST{Nodes: []s.Node{tree.Nodes[id-1]}}).Nodes[0]
		result.Nodes = append(result.Nodes, node)
		newID := ast.NodeID(len(result.Nodes))
		ids[id] = newID
		for i, child := range node.Children { result.Nodes[newID-1].Children[i].Node = copyNode(child.Node) }
		return newID
	}
	result.Root = copyNode(tree.Root)
	for i := range result.Nodes {
		for j, reference := range result.Nodes[i].References {
			if reference.Target.Tree == tree.ID { result.Nodes[i].References[j].Target.Node = ids[reference.Target.Node] }
		}
	}
	return result
}

func TestMissingPinnedDependencyIsIncomplete(t *testing.T) {
	tree := fixture(t)
	stripAnnotations(tree)
	root := find(tree, s.KindModule, "example.people")
	tree.Nodes[tree.Root-1].Children = []ast.Child{{Role: s.RoleModule, Node: root}}
	result := run(t, compact(tree), nil)
	requireFailure(t, result, analysis.Incomplete, analysis.Undetermined, "missing-dependency", "context")
	if result.Diagnostics()[0].At == nil { t.Fatal("missing dependency diagnostic lost its import site") }
}

func TestPresentAnnotationsAndTypePositionsAreCheckedWithoutRepair(t *testing.T) {
	t.Run("incorrect present target", func(t *testing.T) {
		tree := fixture(t)
		application := find(tree, s.KindApplication, "Page")
		for i := range tree.Nodes[application-1].References {
			if tree.Nodes[application-1].References[i].Role == s.RoleTarget { tree.Nodes[application-1].References[i].Target.Node = find(tree, s.KindDeclaration, "Outer") }
		}
		before := clone(tree)
		result := run(t, tree, nil)
		requireFailure(t, result, analysis.Complete, analysis.Rejected, "invalid-annotation", "source")
		if !reflect.DeepEqual(tree, before) { t.Fatal("incorrect annotation was repaired") }
	})
	t.Run("literal represented as a type", func(t *testing.T) {
		root := module("example.literal", mapping(), entry("Literal", definition(h("exact", param(0)), "T")))
		tree, err := resolve.Resolve(root, nil, "literal-position")
		if err != nil { t.Fatal(err) }
		literal := tree.Nodes[body(tree, find(tree, s.KindDeclaration, "Literal"))-1].Children[0].Node
		tree.Nodes[literal-1] = s.Node{Kind: s.KindParameter, Attributes: []ast.Attribute{{Name: s.AttributeIndex, Text: "0"}}}
		result := run(t, compact(tree), nil)
		requireFailure(t, result, analysis.Complete, analysis.Rejected, "invalid-annotation", "source")
	})
}

func TestResolverAndExportBudgetsAreIndependent(t *testing.T) {
	tree := fixture(t)
	result := run(t, tree, func(input *analysis.Input) {
		input.Configuration = []analysis.Setting{{Key: analyze.LimitsSetting(), Value: []byte("[256,1,16777216,256]")}}
	})
	requireFailure(t, result, analysis.Incomplete, analysis.Undetermined, "exhausted", "request")
	result = run(t, tree, func(input *analysis.Input) { input.Limits = analysis.DefaultLimits(); input.Limits.Model.MaxRecords = 1 })
	requireFailure(t, result, analysis.Complete, analysis.Holds, "export-exhausted", "request")
	if result.Limits().Model.MaxRecords != 1 { t.Fatal("result lost its export budget") }
}

func TestInvalidConfigurationEntryAndCancellationReturnNoResult(t *testing.T) {
	if result, err := analyze.Analyze(nil, analysis.Request[s.Language]{}); !errors.Is(err, analysis.ErrInvalidRequest) || result.Valid() { t.Fatalf("nil context: %v", err) }
	if result, err := analyze.Analyze(context.Background(), analysis.Request[s.Language]{}); !errors.Is(err, analysis.ErrInvalidRequest) || result.Valid() { t.Fatalf("zero request: %v", err) }
	analyzer, input := analyzer(t, fixture(t))
	for _, settings := range [][]analysis.Setting{
		{{Key: analyze.LimitsSetting(), Value: []byte("[0,100000,16777216,256]")}},
		{{Key: analyze.LimitsSetting(), Value: []byte("[256,1]")}},
		{{Key: analyze.LimitsSetting(), Value: []byte("[256,100000,16777216,256] true")}},
		{{Key: semantics.Key{Owner: "other", Name: "configuration"}, Value: []byte("[]")}},
	} {
		request := input
		request.Configuration = settings
		result, err := analyzer.Analyze(context.Background(), request)
		if !errors.Is(err, analysis.ErrInvalidRequest) || result.Valid() { t.Fatalf("invalid configuration returned a publication: %v", err) }
	}
	wrongEntry := input
	wrongEntry.Entry.Node = find(fixture(t), s.KindDeclaration, "Outer")
	if result, err := analyzer.Analyze(context.Background(), wrongEntry); !errors.Is(err, analysis.ErrInvalidRequest) || result.Valid() { t.Fatalf("wrong entry: %v", err) }
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if result, err := analyzer.Analyze(ctx, input); !errors.Is(err, context.Canceled) || result.Valid() { t.Fatalf("cancelled analysis published: %v", err) }
}

func TestUnknownObligationIsUnsupportedAndDoesNotInventFacts(t *testing.T) {
	unknown := semantics.Key{Owner: "consumer", Name: "layout-equivalence"}
	result := run(t, fixture(t), func(input *analysis.Input) { input.Requested = []semantics.Key{unknown} })
	requireModel(t, result)
	requireReport(t, result, unknown, analysis.Unsupported, analysis.Undetermined)
}

func TestUnusedPoolAnnotationsDoNotEstablishFacts(t *testing.T) {
	tree := fixture(t)
	extra, err := resolve.Resolve(module("example.unused", mapping(), entry("Unused", definition(local("Target"))), entry("Target", definition(h("scalar", atom("int"))))), nil, "unused")
	if err != nil { t.Fatal(err) }
	// Move the extra module arena into this Program, preserving ordinary syntax
	// and intentionally supplying a false target annotation outside the domain.
	offset := ast.NodeID(len(tree.Nodes))
	for _, node := range extra.Nodes {
		for i := range node.Children { node.Children[i].Node += offset }
		for i := range node.References { node.References[i].Target = ast.NodeRef{Tree: tree.ID, Node: tree.Root} }
		tree.Nodes = append(tree.Nodes, node)
	}
	for _, child := range extra.Nodes[extra.Root-1].Children {
		tree.Nodes[tree.Root-1].Children = append(tree.Nodes[tree.Root-1].Children, ast.Child{Role: s.RoleModule, Node: child.Node})
	}
	// compact discards the unattached extra Program wrapper.
	tree = compact(tree)
	result := run(t, tree, nil)
	model := requireModel(t, result)
	for _, e := range model.Data().Entities { if e.Name == "Unused" || e.Name == "Target" || e.Name == "example.unused" { t.Fatal("unused context leaked into the claimed closure") } }
}

func TestUnknownDescriptorConstructorReachesAuthoritativeAnalysis(t *testing.T) {
	root := module("example.future", mapping(), entry("Future", definition(h("scalar", atom("int")))))
	tree, err := resolve.Resolve(root, nil, "unknown-constructor")
	if err != nil { t.Fatal(err) }
	id := body(tree, find(tree, s.KindDeclaration, "Future"))
	setAttribute(&tree.Nodes[id-1], s.AttributeConstructor, "future-constructor")
	changed := module("example.future", mapping(), entry("Future", definition(h("future-constructor", atom("int")))))
	setAttribute(&tree.Nodes[find(tree, s.KindModule, "example.future")-1], s.AttributePin, pin(t, changed).String())
	result := run(t, tree, nil)
	requireFailure(t, result, analysis.Complete, analysis.Rejected, "invalid", "source")
	if result.Diagnostics()[0].At == nil || result.Diagnostics()[0].At.Node != id { t.Fatal("unknown constructor did not retain its direct syntax location") }
}

func TestSourcePayloadMutationUsesNewPinAndNeverHistoricalOrigin(t *testing.T) {
	root := module("example.scalar", mapping(), entry("Value", definition(h("scalar", atom("int")))))
	tree, err := resolve.Resolve(root, nil, "edited-source")
	if err != nil { t.Fatal(err) }
	id := body(tree, find(tree, s.KindDeclaration, "Value"))
	leaf := tree.Nodes[id-1].Children[0].Node
	setAttribute(&tree.Nodes[leaf-1], s.AttributeHex, hex.EncodeToString([]byte("utf8-text")))
	changed := module("example.scalar", mapping(), entry("Value", definition(h("scalar", atom("utf8-text")))))
	setAttribute(&tree.Nodes[find(tree, s.KindModule, "example.scalar")-1], s.AttributePin, pin(t, changed).String())
	tree.Nodes[id-1].Origin.Nodes = []ast.NodeRef{{Tree: "old-value", Node: 1}}
	model := requireModel(t, run(t, tree, nil))
	_, declaration := entity(t, model, "example.scalar", "Value")
	if !bytes.Equal(term(t, model, declaration.Body).Data, codec.Encode(h("scalar", atom("utf8-text")))) { t.Fatal("analysis reused stale semantic payload") }
}
