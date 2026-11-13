package integration_test

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/engine/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/source/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// These qualify the generated-Files contract (#273) with a fixture language:
// the shared formatter runs the common admission before the provider, asks the
// language which evidence one request requires, and guards the publication.
// The fixture places the Go declarations of the independent generation fixture.
// Its placement type, its preparation step, its provider, its analysis and its
// guard are its own. None of them is the Go Files stage, which does not exist.

// fixturePlacement is the placement type of the fixture language. Beside the
// two tables it names a root file, as a language with a compilation root does.
type fixturePlacement struct {
	Root        files.ID
	Files       map[files.ID]files.Path
	Definitions map[g.DeclarationHandle]files.ID
}

func fixtureKey(name string) semantics.Key {
	return semantics.Key{Owner: "placement-fixture", Name: name}
}

// fixtureRoot is the one setting of the evidence requirement that follows from
// the placement.
func fixtureRoot(root files.ID) []semantic.Setting {
	return []semantic.Setting{{Key: fixtureKey("root"), Value: []byte(root)}}
}

func fixtureAnalysis() analysis.Contract {
	return analysis.Contract{Semantic: semantic.Contract{
		Profile:               "placement-fixture-v1",
		Obligations:           []semantic.Obligation{{Key: fixtureKey("validity"), Domain: fixtureKey("entry")}},
		SemanticsRequirements: []semantic.Requirement{{Obligation: fixtureKey("validity"), Verdict: semantic.Holds}},
	}}
}

// fixtureAnalyzer publishes a graph without a fact for whatever it is given.
type fixtureAnalyzer struct{ gs.Language }

func (fixtureAnalyzer) Analyze(_ context.Context, request analysis.Request[gs.Language]) (semantic.Result, error) {
	return request.Publish(semantic.Publication{Facts: &semantics.Data{}, Reports: []semantic.Report{
		{Obligation: fixtureKey("validity"), Status: semantic.Complete, Verdict: semantic.Holds},
	}})
}

// fixtureLanguage is the fixture's contribution to the generated-Files
// contract, with the switches that the tests turn.
type fixtureLanguage struct {
	gs.Language
	definition representation.Language
	checker    check.Checker
	analyzer   analysis.Analyzer
	issuer     format.GeneratedPlacementIssuer[fixturePlacement]

	// calls counts the calls of the provider.
	calls int
	// skipAdmission makes the step return without the common admission.
	skipAdmission bool
	// again is what a second call of Admit in one step answers.
	again error
	// kept is the request of the last step, for a call after the step returned.
	kept format.GeneratedFilesPreparation
	// evidence changes the requirement that the step returns.
	evidence func(*semantic.EvidenceRequirement)
	// withoutGuard makes the step return no guard.
	withoutGuard bool
	// reconstruct is the second run of the layout that the guard compares the
	// publication with. A test replaces it by one that returns the publication.
	reconstruct func(format.AdmittedPlacement, g.Output, source.Publication[gs.Language]) ([]ast.AST[gs.Language], error)
	// acceptAll makes the language's guard accept every publication.
	acceptAll bool
	// where is the file into which the provider puts a definition.
	where func(format.PlacedDeclaration) files.ID
	// paths changes the files that the provider publishes.
	paths func([]files.File) []files.File
	// evidenceRoot is the root under which the provider has its output analyzed,
	// when it differs from the root of the request's placement.
	evidenceRoot files.ID
}

func newFixtureLanguage(t *testing.T, output g.Output) *fixtureLanguage {
	t.Helper()
	definition := output.Snapshot().Language()
	checker, err := check.NewChecker[gs.Language](definition, golangText())
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[gs.Language](definition, fixtureAnalysis(), fixtureAnalyzer{})
	if err != nil {
		t.Fatal(err)
	}
	return &fixtureLanguage{definition: definition, checker: checker, analyzer: analyzer, issuer: format.NewGeneratedPlacementIssuer[fixturePlacement]()}
}

func (l *fixtureLanguage) formatter(t *testing.T) format.GeneratedFilesFormatter[gs.Language] {
	t.Helper()
	formatter, err := format.NewGeneratedFilesFormatter(l.definition, format.GeneratedFilesContract[gs.Language]{Print: golangText().Print, Prepare: l.prepare}, l)
	if err != nil {
		t.Fatal(err)
	}
	return formatter
}

func (l *fixtureLanguage) requirement(root files.ID) semantic.EvidenceRequirement {
	return semantic.EvidenceRequirement{Operation: analysis.Kind(), Profile: fixtureAnalysis().Semantic.Profile, Configuration: fixtureRoot(root), Reports: fixtureAnalysis().Semantic.SemanticsRequirements, Semantics: true}
}

// prepare is the fixture's preparation step.
func (l *fixtureLanguage) prepare(_ context.Context, preparation format.GeneratedFilesPreparation) (format.GeneratedFilesPrepared[gs.Language], error) {
	zero := format.GeneratedFilesPrepared[gs.Language]{}
	input := preparation.Input()
	placement, issued := l.issuer.Read(input.Placement)
	if !issued {
		return zero, source.ErrInvalid
	}
	prepared := format.GeneratedFilesPrepared[gs.Language]{Evidence: l.requirement(placement.Root)}
	if !l.skipAdmission {
		admitted, err := preparation.Admit(format.GeneratedPlacementView{Files: placement.Files, Definitions: placement.Definitions})
		if err != nil {
			return zero, err
		}
		prepared.Placement = admitted
		_, l.again = preparation.Admit(format.GeneratedPlacementView{Files: placement.Files, Definitions: placement.Definitions})
	}
	l.kept = preparation
	// The fixture's own rule: the root is a file of the placement.
	if _, declared := placement.Files[placement.Root]; !declared {
		return zero, source.ErrInvalid
	}
	if l.evidence != nil {
		l.evidence(&prepared.Evidence)
	}
	if !l.withoutGuard {
		admitted := prepared.Placement
		prepared.Check = func(_ context.Context, publication source.Publication[gs.Language]) error {
			return l.guard(admitted, input.Output, publication)
		}
	}
	return prepared, nil
}

// guard is the fixture's own guard. It compares the published value with the
// admitted placement: each admitted file holds exactly the definitions that the
// placement gives it. The comparison with a second run of the layout is the
// integrity check, and it is the part that a test replaces.
func (l *fixtureLanguage) guard(admitted format.AdmittedPlacement, output g.Output, publication source.Publication[gs.Language]) error {
	if l.acceptAll {
		return nil
	}
	reconstruct := l.reconstruct
	if reconstruct == nil {
		reconstruct = func(admitted format.AdmittedPlacement, output g.Output, _ source.Publication[gs.Language]) ([]ast.AST[gs.Language], error) {
			return fixtureLayout(admitted, output, func(d format.PlacedDeclaration) files.ID { return d.File })
		}
	}
	expected, err := reconstruct(admitted, output, publication)
	if err != nil {
		return err
	}
	if len(expected) != len(publication.Units) {
		return source.ErrInvalid
	}
	for i, unit := range publication.Units {
		if !reflect.DeepEqual(expected[i], unit.Syntax.Abstract) {
			return source.ErrInvalid
		}
	}
	// The published value against the admitted input.
	original, err := fixtureTree(output)
	if err != nil {
		return err
	}
	published := map[files.ID]ast.AST[gs.Language]{}
	for _, unit := range publication.Units {
		published[unit.File] = unit.Syntax.Abstract
	}
	next := map[files.ID]int{}
	for _, definition := range admitted.Definitions() {
		tree, exists := published[definition.File]
		if !exists {
			return source.ErrInvalid
		}
		declarations := tree.Nodes[tree.Root-1].Children
		position := next[definition.File]
		if position >= len(declarations) || !fixtureSameSubtree(original, definition.Declaration.Site.Node, tree, declarations[position].Node) {
			return source.ErrInvalid
		}
		next[definition.File]++
	}
	for file, tree := range published {
		if len(tree.Nodes[tree.Root-1].Children) != next[file] {
			return source.ErrInvalid
		}
	}
	return nil
}

func fixtureTree(output g.Output) (ast.AST[gs.Language], error) {
	entry, err := output.Snapshot().EntryTree()
	if err != nil {
		return ast.AST[gs.Language]{}, err
	}
	return representation.Unpack[gs.Language](output.Snapshot().Language(), entry)
}

// fixtureSameSubtree compares a declaration with its copy: kinds, attributes
// and children in order. A copy carries no reference.
func fixtureSameSubtree(a ast.AST[gs.Language], x ast.NodeID, b ast.AST[gs.Language], y ast.NodeID) bool {
	left, right := a.Nodes[x-1], b.Nodes[y-1]
	if left.Kind != right.Kind || !slices.Equal(left.Attributes, right.Attributes) || len(left.Children) != len(right.Children) {
		return false
	}
	for i := range left.Children {
		if left.Children[i].Role != right.Children[i].Role || !fixtureSameSubtree(a, left.Children[i].Node, b, right.Children[i].Node) {
			return false
		}
	}
	return true
}

// fixtureLayout builds one Go file per admitted file, in the order of the file
// identities, with the definitions that where gives it in the order of the
// inventory.
func fixtureLayout(admitted format.AdmittedPlacement, output g.Output, where func(format.PlacedDeclaration) files.ID) ([]ast.AST[gs.Language], error) {
	original, err := fixtureTree(output)
	if err != nil {
		return nil, err
	}
	var trees []ast.AST[gs.Language]
	for _, file := range admitted.Files() {
		tree := ast.AST[gs.Language]{ID: ast.TreeID("fixture/files/" + string(file.ID))}
		var copyNode func(ast.NodeID) ast.NodeID
		copyNode = func(id ast.NodeID) ast.NodeID {
			node := original.Nodes[id-1]
			children := node.Children
			node.Children, node.References, node.Origin = nil, nil, ast.Origin{}
			node.Attributes = slices.Clone(node.Attributes)
			for _, child := range children {
				node.Children = append(node.Children, ast.Child{Role: child.Role, Node: copyNode(child.Node)})
			}
			tree.Nodes = append(tree.Nodes, node)
			return ast.NodeID(len(tree.Nodes))
		}
		root := ast.Node[gs.Language]{Kind: gs.KindFile, Attributes: slices.Clone(original.Nodes[original.Root-1].Attributes)}
		for _, definition := range admitted.Definitions() {
			if where(definition) == file.ID {
				root.Children = append(root.Children, ast.Child{Role: gs.RoleDeclaration, Node: copyNode(definition.Declaration.Site.Node)})
			}
		}
		tree.Nodes = append(tree.Nodes, root)
		tree.Root = ast.NodeID(len(tree.Nodes))
		trees = append(trees, tree)
	}
	return trees, nil
}

// Format is the fixture's provider.
func (l *fixtureLanguage) Format(ctx context.Context, request format.GeneratedFilesRequest[gs.Language]) (source.Source[gs.Language], error) {
	l.calls++
	zero := source.Source[gs.Language]{}
	input := request.Input()
	placement, issued := l.issuer.Read(input.Placement)
	if !issued {
		return zero, source.ErrInvalid
	}
	where := l.where
	if where == nil {
		where = func(d format.PlacedDeclaration) files.ID { return d.File }
	}
	trees, err := fixtureLayout(request.Placement(), input.Output, where)
	if err != nil {
		return zero, err
	}
	admittedFiles := request.Placement().Files()
	if l.paths != nil {
		admittedFiles = l.paths(admittedFiles)
	}
	captured := semantic.Input{Files: &files.Context{}}
	var units []source.Unit[gs.Language]
	for i, file := range admittedFiles {
		tree := trees[i]
		syntax, err := golangText().Format(ctx, tree)
		if err != nil {
			return zero, err
		}
		text, err := golangText().Print(ctx, syntax)
		if err != nil {
			return zero, err
		}
		packed, err := representation.Pack(l.definition, syntax)
		if err != nil {
			return zero, err
		}
		if i == 0 {
			captured.Entry = ast.NodeRef{Tree: tree.ID, Node: tree.Root}
		}
		captured.Trees = append(captured.Trees, packed)
		captured.Files.Files.Entries = append(captured.Files.Files.Entries, files.File{ID: file.ID, Path: file.Path, Text: text})
		captured.Files.Links = append(captured.Files.Links, files.SyntaxLink{File: file.ID, Tree: tree.ID, Source: ast.SourceID(file.ID)})
		units = append(units, source.Unit[gs.Language]{File: file.ID, Path: file.Path, Source: ast.SourceID(file.ID), Syntax: syntax})
	}
	snapshot, err := semantic.Capture(ctx, l.checker, captured)
	if err != nil {
		return zero, err
	}
	root := placement.Root
	if l.evidenceRoot != "" {
		root = l.evidenceRoot
	}
	result, err := l.analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Options: semantic.Options{Configuration: fixtureRoot(root)}})
	if err != nil {
		return zero, err
	}
	return request.Publish(source.Publication[gs.Language]{Snapshot: snapshot, Analysis: result, Units: units})
}

// fixtureCase is a sealed output of the generation fixture with a placement of
// its five declarations into two files and a third file without a definition.
type fixtureCase struct {
	output   g.Output
	language *fixtureLanguage
	handles  []g.DeclarationHandle
}

func newFixtureCase(t *testing.T) *fixtureCase {
	t.Helper()
	output := inventoryOutput(t, independent(t), inventoryRoots...)
	c := &fixtureCase{output: output, language: newFixtureLanguage(t, output)}
	for _, record := range output.Inventory() {
		c.handles = append(c.handles, record.Handle)
	}
	if len(c.handles) != 5 {
		t.Fatal("the fixture output owns five declarations", len(c.handles))
	}
	return c
}

// placement puts the first two declarations into the file a and the others
// into the file b. The file c holds no definition.
func (c *fixtureCase) placement() fixturePlacement {
	p := fixturePlacement{Root: "a", Files: map[files.ID]files.Path{"a": "out/a.go", "b": "out/b.go", "c": "out/c.go"}, Definitions: map[g.DeclarationHandle]files.ID{}}
	for i, handle := range c.handles {
		p.Definitions[handle] = []files.ID{"a", "a", "b", "b", "b"}[i]
	}
	return p
}

func (c *fixtureCase) format(t *testing.T, placement fixturePlacement, limits format.FilesLimits) (source.Source[gs.Language], error) {
	t.Helper()
	return c.language.formatter(t).Format(context.Background(), format.GeneratedFilesInput{Output: c.output, Placement: c.language.issuer.Issue(placement), Limits: limits})
}

// onlyCategory requires that a refusal matches want and no other Files
// category. A nil want is a refusal without a category.
func onlyCategory(t *testing.T, name string, err, want error) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: accepted", name)
		return
	}
	for _, category := range filesCategories {
		if errors.Is(err, category.value) != (category.value == want) {
			t.Errorf("%s: errors.Is(err, %s) = %v for %q", name, category.name, errors.Is(err, category.value), err)
		}
	}
}

// The shared formatter runs the common admission before it calls the provider.
// A placement that the admission refuses never reaches the provider, and each
// refusal has its category.
func TestGeneratedFilesAdmissionRunsBeforeTheProvider(t *testing.T) {
	c := newFixtureCase(t)
	if out, err := c.format(t, c.placement(), format.FilesLimits{}); err != nil || !out.Valid() || c.language.calls != 1 {
		t.Fatal("the placement that every case changes must format, with one call of the provider", err, c.language.calls)
	}
	other := inventoryOutput(t, independent(t), inventoryRoots...)
	refused := []struct {
		name   string
		want   error
		also   error
		limits func(*format.FilesLimits)
		change func(*fixturePlacement)
	}{
		{name: "a declaration without an assignment", want: source.ErrInvalid, also: files.ErrInvalid, change: func(p *fixturePlacement) {
			delete(p.Definitions, c.handles[4])
		}},
		{name: "an assignment to an undeclared file", want: source.ErrInvalid, also: files.ErrInvalid, change: func(p *fixturePlacement) {
			p.Definitions[c.handles[4]] = "absent"
		}},
		{name: "a handle of another output", want: source.ErrInvalid, also: g.ErrScope, change: func(p *fixturePlacement) {
			p.Definitions[other.Inventory()[4].Handle] = "b"
		}},
		{name: "a handle of another output in place of this output's", want: source.ErrInvalid, also: g.ErrScope, change: func(p *fixturePlacement) {
			delete(p.Definitions, c.handles[4])
			p.Definitions[other.Inventory()[4].Handle] = "b"
		}},
		{name: "the zero handle", want: source.ErrInvalid, also: g.ErrScope, change: func(p *fixturePlacement) {
			p.Definitions[g.DeclarationHandle{}] = "b"
		}},
		{name: "a file set that the files contract refuses: one path for two files", want: source.ErrInvalid, also: files.ErrInvalid, change: func(p *fixturePlacement) {
			p.Files["c"] = p.Files["a"]
		}},
		{name: "a file set that the files contract refuses: an empty file identity", want: source.ErrInvalid, also: files.ErrInvalid, change: func(p *fixturePlacement) {
			p.Files[""] = "out/d.go"
		}},
		{name: "a file set that the files contract refuses: a path that leaves the root", want: source.ErrInvalid, also: files.ErrInvalid, change: func(p *fixturePlacement) {
			p.Files["c"] = "../c.go"
		}},
		{name: "no file", want: source.ErrInvalid, also: files.ErrInvalid, change: func(p *fixturePlacement) {
			clear(p.Files)
		}},
		{name: "more files than the file limit", want: source.ErrLimit, also: files.ErrLimit, limits: func(l *format.FilesLimits) { l.Source.Files.MaxFiles = 2 }},
		{name: "more files than the link limit", want: source.ErrLimit, also: files.ErrLimit, limits: func(l *format.FilesLimits) { l.Source.Files.MaxLinks = 2 }},
		{name: "file identities and paths over the byte limit", want: source.ErrLimit, also: files.ErrLimit, limits: func(l *format.FilesLimits) { l.Source.Files.MaxBytes = 8 }},
		{name: "the sealed output over the input node limit", want: source.ErrLimit, also: semantic.ErrLimit, limits: func(l *format.FilesLimits) { l.Semantic.MaxInputNodes = 1 }},
	}
	for _, r := range refused {
		placement, limits := c.placement(), format.DefaultFilesLimits()
		if r.change != nil {
			r.change(&placement)
		}
		if r.limits != nil {
			r.limits(&limits)
		}
		c.language.calls = 0
		out, err := c.format(t, placement, limits)
		if out.Valid() {
			t.Errorf("%s: a Source was returned", r.name)
		}
		onlyCategory(t, r.name, err, r.want)
		if r.also != nil && !errors.Is(err, r.also) {
			t.Errorf("%s: the refusal does not match %q: %v", r.name, r.also, err)
		}
		if c.language.calls != 0 {
			t.Errorf("%s: the provider was called %d times", r.name, c.language.calls)
		}
	}
	// The limits that refuse above admit a placement that meets them exactly.
	limits := format.DefaultFilesLimits()
	limits.Source.Files.MaxFiles, limits.Source.Files.MaxLinks = 3, 3
	if out, err := c.format(t, c.placement(), limits); err != nil || !out.Valid() {
		t.Fatal("three files under a file limit and a link limit of three", err)
	}
}

// The common admission accepts a file without a definition. It interprets no
// scope, so it cannot say that such a file is empty in a language.
func TestGeneratedFilesAdmissionAcceptsAFileWithoutADefinition(t *testing.T) {
	c := newFixtureCase(t)
	out, err := c.format(t, c.placement(), format.FilesLimits{})
	if err != nil || c.language.calls != 1 {
		t.Fatal(err, c.language.calls)
	}
	units := out.Units()
	if len(units) != 3 || units[2].File != "c" || len(units[2].Syntax.Abstract.Nodes) != 1 {
		t.Fatal("the file without a definition is published as a file without a declaration", len(units))
	}
}

// The admitted placement states the order that a formatter realizes: the files
// by their identities, the declarations by their sites in the sealed output.
// Two placements with the same assignments are the same placement.
func TestGeneratedFilesAdmittedPlacementStatesTheOrder(t *testing.T) {
	c := newFixtureCase(t)
	var texts [][]string
	for run := 0; run < 6; run++ {
		// The same assignments, written in another order.
		placement := fixturePlacement{Root: "a", Files: map[files.ID]files.Path{}, Definitions: map[g.DeclarationHandle]files.ID{}}
		base := c.placement()
		ids := []files.ID{"a", "b", "c"}
		handles := slices.Clone(c.handles)
		if run%2 == 1 {
			slices.Reverse(ids)
			slices.Reverse(handles)
		}
		for _, id := range ids {
			placement.Files[id] = base.Files[id]
		}
		for _, handle := range handles {
			placement.Definitions[handle] = base.Definitions[handle]
		}
		out, err := c.format(t, placement, format.FilesLimits{})
		if err != nil {
			t.Fatal(err)
		}
		var text []string
		for i, unit := range out.Units() {
			if unit.File != []files.ID{"a", "b", "c"}[i] {
				t.Fatal("the files are not in the order of their identities", i, unit.File)
			}
			printed, err := golangText().Print(context.Background(), unit.Syntax)
			if err != nil {
				t.Fatal(err)
			}
			text = append(text, printed)
		}
		texts = append(texts, text)
	}
	for _, text := range texts[1:] {
		if !slices.Equal(text, texts[0]) {
			t.Fatal("two placements with the same assignments gave different files")
		}
	}
	// The declarations of a file are in the order of the inventory.
	original, err := fixtureTree(c.output)
	if err != nil {
		t.Fatal(err)
	}
	out, err := c.format(t, c.placement(), format.FilesLimits{})
	if err != nil {
		t.Fatal(err)
	}
	inventory := c.output.Inventory()
	for i, expected := range [][]int{{0, 1}, {2, 3, 4}, {}} {
		tree := out.Units()[i].Syntax.Abstract
		declarations := tree.Nodes[tree.Root-1].Children
		if len(declarations) != len(expected) {
			t.Fatal("declarations of file", i, len(declarations))
		}
		for position, record := range expected {
			if !fixtureSameSubtree(original, inventory[record].Site.Node, tree, declarations[position].Node) {
				t.Fatal("the declarations of a file are not in the order of the inventory", i, position)
			}
		}
	}
}

// The step states the evidence requirement for one request. Here one setting of
// the requirement follows from the placement, as a compilation root does.
// Evidence that was established for one placement is refused for a request that
// carries another: the right files under the wrong root.
func TestGeneratedFilesEvidenceRequirementFollowsThePlacement(t *testing.T) {
	c := newFixtureCase(t)
	first, second := c.placement(), c.placement()
	second.Root = "b"
	// The provider has its output analyzed under the root of the second
	// placement, whatever the request carries.
	c.language.evidenceRoot = "b"
	out, err := c.format(t, first, format.FilesLimits{})
	if err == nil || out.Valid() {
		t.Fatal("evidence for the root b was accepted for a request with the root a")
	}
	if !errors.Is(err, semantic.ErrInvalidRequest) || c.language.calls != 1 {
		t.Fatal("the publisher refuses a receipt that does not match the requirement of the request", err, c.language.calls)
	}
	out, err = c.format(t, second, format.FilesLimits{})
	if err != nil || !out.Valid() {
		t.Fatal("evidence for the root b is the evidence that a request with the root b requires", err)
	}
	if got := out.Evidence().Configuration(); len(got) != 1 || string(got[0].Value) != "b" {
		t.Fatal("the receipt of the Source is not bound to the root", got)
	}
}

// A step whose requirement asks for no semantic graph, a step that returns no
// guard and a step that skips the common admission fail the request. No Source
// is returned and the provider is not called.
func TestGeneratedFilesStepWithoutEvidenceGuardOrAdmissionFails(t *testing.T) {
	for name, change := range map[string]func(*fixtureLanguage){
		"a requirement without semantic evidence": func(l *fixtureLanguage) {
			l.evidence = func(r *semantic.EvidenceRequirement) { r.Semantics = false }
		},
		"no guard":            func(l *fixtureLanguage) { l.withoutGuard = true },
		"no common admission": func(l *fixtureLanguage) { l.skipAdmission = true },
	} {
		c := newFixtureCase(t)
		change(c.language)
		out, err := c.format(t, c.placement(), format.FilesLimits{})
		if !errors.Is(err, format.ErrInvalidFormatter) || out.Valid() || c.language.calls != 0 {
			t.Errorf("%s: err=%v valid=%v calls=%d", name, err, out.Valid(), c.language.calls)
		}
		onlyCategory(t, name, err, nil)
	}
	// A requirement that the semantic contract refuses fails with that
	// contract's value.
	c := newFixtureCase(t)
	c.language.evidence = func(r *semantic.EvidenceRequirement) { r.Profile = "" }
	if out, err := c.format(t, c.placement(), format.FilesLimits{}); !errors.Is(err, semantic.ErrInvalidOperation) || out.Valid() || c.language.calls != 0 {
		t.Fatal("a malformed requirement", err, c.language.calls)
	}
	// The step has its view admitted once, while it runs. A second call and a
	// call after the step returned are refused, so no admitted placement exists
	// that the shared formatter did not receive from the step of its request.
	c = newFixtureCase(t)
	if out, err := c.format(t, c.placement(), format.FilesLimits{}); err != nil || !out.Valid() {
		t.Fatal(err)
	}
	if !errors.Is(c.language.again, format.ErrInvalidFormatter) {
		t.Fatal("a second admission in one step", c.language.again)
	}
	late := c.placement()
	view := format.GeneratedPlacementView{Files: late.Files, Definitions: late.Definitions}
	if admitted, err := c.language.kept.Admit(view); !errors.Is(err, format.ErrInvalidFormatter) || admitted.Valid() {
		t.Fatal("an admission after the step returned", err)
	}
	// Also when the step itself never asked for one.
	c = newFixtureCase(t)
	c.language.skipAdmission = true
	if out, err := c.format(t, c.placement(), format.FilesLimits{}); err == nil || out.Valid() {
		t.Fatal("a step without the common admission")
	}
	if admitted, err := c.language.kept.Admit(view); !errors.Is(err, format.ErrInvalidFormatter) || admitted.Valid() {
		t.Fatal("an admission after a step that asked for none", err)
	}
	// A placement of another issuer, and no placement, are refusals of the
	// step. They are invalid input and no failed type assertion.
	c = newFixtureCase(t)
	foreign := format.NewGeneratedPlacementIssuer[fixturePlacement]().Issue(c.placement())
	for name, placement := range map[string]format.GeneratedPlacement{"a placement of another issuer": foreign, "no placement": {}} {
		out, err := c.language.formatter(t).Format(context.Background(), format.GeneratedFilesInput{Output: c.output, Placement: placement})
		onlyCategory(t, name, err, source.ErrInvalid)
		if out.Valid() || c.language.calls != 0 {
			t.Errorf("%s: valid=%v calls=%d", name, out.Valid(), c.language.calls)
		}
	}
	if _, read := c.language.issuer.Read(foreign); read {
		t.Fatal("an issuer read the placement of another issuer")
	}
}

// A provider publishes a definition in another file than the admitted placement
// states. The guard is the fixture language's own, and the formatter is built
// from that language's step and the tampering provider. The publication is
// refused. It is also refused when the second run of the layout is replaced by
// one that returns the publication itself: what refuses is the comparison of
// the published value with the admitted placement.
func TestGeneratedFilesGuardComparesThePublicationWithTheAdmittedPlacement(t *testing.T) {
	publishedItself := func(_ format.AdmittedPlacement, _ g.Output, publication source.Publication[gs.Language]) ([]ast.AST[gs.Language], error) {
		var trees []ast.AST[gs.Language]
		for _, unit := range publication.Units {
			trees = append(trees, unit.Syntax.Abstract)
		}
		return trees, nil
	}
	for _, replaced := range []bool{false, true} {
		c := newFixtureCase(t)
		if replaced {
			c.language.reconstruct = publishedItself
		}
		if out, err := c.format(t, c.placement(), format.FilesLimits{}); err != nil || !out.Valid() {
			t.Fatal("the untampered provider must publish", replaced, err)
		}
		// The last declaration belongs to the file b. The provider puts it into c.
		moved := c.handles[4]
		c.language.where = func(d format.PlacedDeclaration) files.ID {
			if d.Declaration.Handle == moved {
				return "c"
			}
			return d.File
		}
		c.language.calls = 0
		out, err := c.format(t, c.placement(), format.FilesLimits{})
		if err == nil || out.Valid() {
			t.Fatal("a definition in another file than the placement states was published", replaced)
		}
		onlyCategory(t, "a definition in another file", err, source.ErrInvalid)
		if c.language.calls != 1 {
			t.Fatal("the provider must have been called", c.language.calls)
		}
	}
}

// After the provider returns, the shared side refuses a publication whose file
// identities and paths differ from the admitted files table, whatever the
// language's guard says.
func TestGeneratedFilesPublicationMatchesTheAdmittedFiles(t *testing.T) {
	for name, paths := range map[string]func([]files.File) []files.File{
		"a file at another path": func(admitted []files.File) []files.File {
			admitted[2].Path = "out/elsewhere.go"
			return admitted
		},
		"a file under another identity": func(admitted []files.File) []files.File {
			admitted[2].ID = "d"
			return admitted
		},
		"one file fewer": func(admitted []files.File) []files.File { return admitted[:2] },
	} {
		c := newFixtureCase(t)
		c.language.acceptAll = true
		if out, err := c.format(t, c.placement(), format.FilesLimits{}); err != nil || !out.Valid() {
			t.Fatal("the untampered provider must publish", name, err)
		}
		c.language.paths = paths
		c.language.calls = 0
		// One file fewer leaves a definition without a file: keep them in a and b.
		out, err := c.format(t, c.placement(), format.FilesLimits{})
		if err == nil || out.Valid() || c.language.calls != 1 {
			t.Errorf("%s: published, err=%v calls=%d", name, err, c.language.calls)
		}
		onlyCategory(t, name, err, source.ErrInvalid)
	}
}

// The engine carries a placement and does not read it. A formatter that is
// registered for the language of the output receives the value as issued.
func TestEngineCarriesTheGeneratedPlacement(t *testing.T) {
	c := newFixtureCase(t)
	kernel := engine.New()
	if err := kernel.RegisterLanguage(c.language.definition); err != nil {
		t.Fatal(err)
	}
	if err := kernel.RegisterGeneratedFilesFormatter(c.language.formatter(t)); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	boxed, err := kernel.FormatGeneratedFiles(ctx, format.GeneratedFilesInput{Output: c.output, Placement: c.language.issuer.Issue(c.placement())})
	if err != nil {
		t.Fatal(err)
	}
	out, err := source.Unpack[gs.Language](c.language.definition, boxed)
	if err != nil || len(out.Units()) != 3 {
		t.Fatal(err)
	}
	foreign := format.NewGeneratedPlacementIssuer[fixturePlacement]().Issue(c.placement())
	boxed, err = kernel.FormatGeneratedFiles(ctx, format.GeneratedFilesInput{Output: c.output, Placement: foreign})
	onlyCategory(t, "a placement of another issuer through the engine", err, source.ErrInvalid)
	if boxed.Valid() {
		t.Fatal("a Source for a placement of another issuer")
	}
}
