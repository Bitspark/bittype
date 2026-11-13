package integration_test

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	tsfiles "github.com/Bitspark/bittype/format/files/typescript/provider/go"
	printFiles "github.com/Bitspark/bittype/print/files/shared/provider/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	"github.com/Bitspark/bittype/source/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// These qualify the TypeScript placement of generated declarations (#273): two
// keyed tables of handles, no entry and no order.

// Run one plan twice. The two outputs have the same sites, and a placement that
// was built from the first is refused with the second as invalid input: a
// handle carries its output.
func TestTypeScriptPlacementBelongsToOneRun(t *testing.T) {
	ctx := context.Background()
	root := genModule("visibility", genMap(), genEntry("Hidden", genType(genPerson())), genEntry("Visible", genType(genLocal("Hidden"))))
	e, _, input := generationInput(t, root)
	to, err := representation.NewLanguage[ts.Language](typescriptText(), "native-types-v1")
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterLanguage(to); err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[ts.Language](to, typescriptText())
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterChecker(checker); err != nil {
		t.Fatal(err)
	}
	request, err := tsRequest(e.Generation(), input, to)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := e.PlanGeneration(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	first, err := plan.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
	second, err := plan.Run(ctx)
	if err != nil {
		t.Fatal(err)
	}
	sites := func(output g.Output) []ast.NodeRef {
		var out []ast.NodeRef
		for _, record := range output.Inventory() {
			out = append(out, record.Site)
		}
		return out
	}
	if len(sites(first)) != 2 || !reflect.DeepEqual(sites(first), sites(second)) {
		t.Fatal("the two runs must own the same sites, or an address would tell them apart", sites(first), sites(second))
	}
	formatter, err := tsfiles.NewGenerated(placementStages(t))
	if err != nil {
		t.Fatal(err)
	}
	placement := tsOneFile(first, "types", "types.ts").Issue()
	if out, err := formatter.Format(ctx, format.GeneratedFilesInput{Output: first, Placement: placement}); err != nil || !out.Valid() {
		t.Fatal("the placement must format with the output it was built from", err)
	}
	out, err := formatter.Format(ctx, format.GeneratedFilesInput{Output: second, Placement: placement})
	if out.Valid() {
		t.Fatal("a placement of the first run was accepted with the second")
	}
	onlyCategory(t, "a placement of the first run with the second", err, source.ErrInvalid)
	if !errors.Is(err, g.ErrScope) {
		t.Fatal("the refusal no longer carries the value with which the output refuses the handle", err)
	}
}

// threeFiles places the four declarations of the Cell page into three files.
// forward decides the order in which the two tables are written.
func threeFiles(handles map[string]g.DeclarationHandle, forward bool) tsfiles.GeneratedPlacement {
	paths := []struct {
		file files.ID
		path files.Path
	}{{"a-model", "generated/model/person.ts"}, {"b-cell", "generated/cell.ts"}, {"c-pages", "generated/pages/page.ts"}}
	declarations := []struct {
		name string
		file files.ID
	}{{"Person", "a-model"}, {"Cell", "b-cell"}, {"Page", "c-pages"}, {"Cells", "c-pages"}}
	if !forward {
		slices.Reverse(paths)
		slices.Reverse(declarations)
	}
	placement := tsfiles.GeneratedPlacement{Files: map[files.ID]files.Path{}, Definitions: map[g.DeclarationHandle]files.ID{}}
	for _, file := range paths {
		placement.Files[file.file] = file.path
	}
	for _, declaration := range declarations {
		placement.Definitions[handles[declaration.name]] = declaration.file
	}
	return placement
}

// Two placements that state the same assignments give byte-identical files and
// the same unit order, whatever order they were built in. The order is the
// formatter's: the files by their identities, and the declarations of one file
// by their sites in the sealed output.
func TestTypeScriptPlacementStatesNoOrder(t *testing.T) {
	ctx := context.Background()
	output := generatedCellPage(t)
	handles := tsDeclarations(t, output)
	stages := placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations})
	formatter, err := tsfiles.NewGenerated(stages)
	if err != nil {
		t.Fatal(err)
	}
	printerProvider, err := printFiles.New[ts.Language](ts.Language{}, stages.Printer)
	if err != nil {
		t.Fatal(err)
	}
	printer, err := print.NewFilesPrinter[ts.Language](stages.Language, printerProvider)
	if err != nil {
		t.Fatal(err)
	}
	var emitted []files.Files
	var last source.Source[ts.Language]
	for run := 0; run < 6; run++ {
		out, err := formatter.Format(ctx, format.GeneratedFilesInput{Output: output, Placement: threeFiles(handles, run%2 == 0).Issue()})
		if err != nil {
			t.Fatal(run, err)
		}
		for i, unit := range out.Units() {
			if unit.File != []files.ID{"a-model", "b-cell", "c-pages"}[i] {
				t.Fatal("the files are not in the order of their identities", run, i, unit.File)
			}
		}
		printed, err := printer.Print(ctx, out)
		if err != nil {
			t.Fatal(err)
		}
		emitted = append(emitted, printed)
		last = out
	}
	for run, printed := range emitted[1:] {
		if !reflect.DeepEqual(printed, emitted[0]) {
			t.Fatal("two placements with the same assignments gave different files", run+1)
		}
	}
	// The declarations of the file c-pages follow the sites of the output.
	var want []string
	entry, err := output.Snapshot().EntryTree()
	if err != nil {
		t.Fatal(err)
	}
	original, err := representation.Unpack[ts.Language](output.Snapshot().Language(), entry)
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range output.Inventory() {
		if name := nodeAttribute(original.Nodes[record.Site.Node-1], ts.AttributeName); name == "Page" || name == "Cells" {
			want = append(want, name)
		}
	}
	var got []string
	pages := last.Units()[2].Syntax.Abstract
	for _, child := range pages.Nodes[pages.Root-1].Children {
		if child.Role == ts.RoleDeclaration {
			got = append(got, nodeAttribute(pages.Nodes[child.Node-1], ts.AttributeName))
		}
	}
	if len(want) != 2 || !slices.Equal(got, want) {
		t.Fatal("the declarations of a file are not in the order of their sites in the output", got, want)
	}
}

// A generated-Files request names no entry. It succeeds, and formatting it
// twice gives the same snapshot entry: the root of the file of least identity.
// The entry is an address and names no root of a TypeScript program.
func TestTypeScriptGeneratedFilesNeedNoEntry(t *testing.T) {
	ctx := context.Background()
	output := generatedCellPage(t)
	handles := tsDeclarations(t, output)
	formatter, err := tsfiles.NewGenerated(placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations}))
	if err != nil {
		t.Fatal(err)
	}
	var entries []ast.NodeRef
	for run := 0; run < 2; run++ {
		placement := threeFiles(handles, run == 0)
		issued := placement.Issue()
		// The issued value owns its tables: what the consumer changes afterwards
		// is no part of the request.
		delete(placement.Definitions, handles["Person"])
		placement.Files["later"] = "later.ts"
		out, err := formatter.Format(ctx, format.GeneratedFilesInput{Output: output, Placement: issued})
		if err != nil {
			t.Fatal(run, err)
		}
		least := out.Units()[0]
		if least.File != "a-model" || out.Entry() != (ast.NodeRef{Tree: least.Syntax.Abstract.ID, Node: least.Syntax.Abstract.Root}) || out.Snapshot().Entry() != out.Entry() {
			t.Fatal("the snapshot entry is not the root of the file of least identity", run, out.Entry())
		}
		entries = append(entries, out.Entry())
	}
	if entries[0] != entries[1] {
		t.Fatal("formatting one placement twice gave two snapshot entries", entries)
	}
}

// A placement of existing source names no entry either. The lists of its
// setting state an order, and the entry of the fresh snapshot does not follow
// it: it is the file of least identity.
func TestTypeScriptExistingPlacementNeedsNoEntry(t *testing.T) {
	ctx := context.Background()
	stages := placementStages(t)
	formatter, err := tsfiles.New(stages)
	if err != nil {
		t.Fatal(err)
	}
	unit := importedUnit(t, "export type Person = string; export type People = [Person, Person];")
	in, err := stages.Capture(ctx, []source.Unit[ts.Language]{unit}, unit.File, semantic.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	tree := unit.Syntax.Abstract
	site := func(name string) ast.NodeRef {
		for _, child := range tree.Nodes[tree.Root-1].Children {
			if nodeAttribute(tree.Nodes[child.Node-1], ts.AttributeName) == name {
				return ast.NodeRef{Tree: tree.ID, Node: child.Node}
			}
		}
		t.Fatal("no declaration", name)
		return ast.NodeRef{}
	}
	// The file pages is the first of the list and not the least identity.
	setting, err := (tsfiles.Organization{Units: []tsfiles.Placement{
		{File: "pages", Path: "pages.ts", Declarations: []ast.NodeRef{site("People")}},
		{File: "model", Path: "model.ts", Declarations: []ast.NodeRef{site("Person")}},
	}}).Setting()
	if err != nil {
		t.Fatal(err)
	}
	in.Configuration = []semantic.Setting{setting}
	out, err := formatter.Format(ctx, in)
	if err != nil {
		t.Fatal(err)
	}
	units := out.Units()
	if len(units) != 2 || units[0].File != "pages" || units[1].File != "model" {
		t.Fatal("the units of existing source follow the lists of the setting", len(units))
	}
	model := units[1].Syntax.Abstract
	if out.Entry() != (ast.NodeRef{Tree: model.ID, Node: model.Root}) {
		t.Fatal("the snapshot entry is not the file of least identity", out.Entry())
	}
}

// replayingProvider publishes what another request published.
type replayingProvider struct {
	ts.Language
	from  source.Source[ts.Language]
	calls *int
}

func (p replayingProvider) Format(_ context.Context, request format.GeneratedFilesRequest[ts.Language]) (source.Source[ts.Language], error) {
	*p.calls++
	return request.Publish(source.Publication[ts.Language]{Snapshot: p.from.Snapshot(), Analysis: p.from.Evidence(), Units: p.from.Units(), Groups: p.from.Groups(), Relations: p.from.Relations()})
}

// A test builds a generated-Files formatter from TypeScript's own preparation
// step and another provider. The provider publishes what the TypeScript
// formatter published for another placement: well formed, analyzed and bound to
// its receipt. TypeScript's guard refuses a definition in another file than the
// admitted placement states, and the shared contract refuses a file at another
// path.
func TestTypeScriptGuardReceivesAPublicationOfAnotherProvider(t *testing.T) {
	ctx := context.Background()
	output := generatedCellPage(t)
	handles := tsDeclarations(t, output)
	stages := placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations})
	genuine, err := tsfiles.NewGenerated(stages)
	if err != nil {
		t.Fatal(err)
	}
	contract, err := stages.GeneratedContract()
	if err != nil {
		t.Fatal(err)
	}
	paths := func() map[files.ID]files.Path {
		return map[files.ID]files.Path{"pages": "generated/pages/page.ts", "model": "generated/model/person.ts"}
	}
	admitted := tsPlaced(t, handles, paths(), map[string]files.ID{"Page": "pages", "Cells": "pages", "Cell": "pages", "Person": "model"})
	// The same files, with Cell in the file model in place of Person.
	otherDefinitions := tsPlaced(t, handles, paths(), map[string]files.ID{"Page": "pages", "Cells": "pages", "Cell": "model", "Person": "pages"})
	// The same definitions, with the file model at another path.
	otherPath := tsPlaced(t, handles, paths(), map[string]files.ID{"Page": "pages", "Cells": "pages", "Cell": "pages", "Person": "model"})
	otherPath.Files["model"] = "generated/model/people.ts"
	replay := func(published tsfiles.GeneratedPlacement) (source.Source[ts.Language], int, error) {
		from, err := genuine.Format(ctx, format.GeneratedFilesInput{Output: output, Placement: published.Issue()})
		if err != nil {
			t.Fatal("the TypeScript formatter must publish the placement that is replayed", err)
		}
		calls := 0
		formatter, err := format.NewGeneratedFilesFormatter(stages.Language, contract, replayingProvider{from: from, calls: &calls})
		if err != nil {
			t.Fatal(err)
		}
		out, err := formatter.Format(ctx, format.GeneratedFilesInput{Output: output, Placement: admitted.Issue()})
		return out, calls, err
	}
	// What the formatter published for the admitted placement itself is accepted,
	// so a refusal below is a refusal of the difference.
	if out, calls, err := replay(admitted); err != nil || !out.Valid() || calls != 1 {
		t.Fatal("the publication of the admitted placement", err, calls)
	}
	for name, published := range map[string]tsfiles.GeneratedPlacement{
		"a definition in another file than the admitted placement states": otherDefinitions,
		"a file at another path than the admitted placement states":       otherPath,
	} {
		out, calls, err := replay(published)
		if out.Valid() || calls != 1 {
			t.Errorf("%s: published=%v calls=%d", name, out.Valid(), calls)
		}
		onlyCategory(t, name, err, source.ErrInvalid)
	}
}

// countedProvider counts its calls and publishes nothing.
type countedProvider struct {
	ts.Language
	calls *int
}

func (p countedProvider) Format(context.Context, format.GeneratedFilesRequest[ts.Language]) (source.Source[ts.Language], error) {
	*p.calls++
	return source.Source[ts.Language]{}, errReached
}

// The TypeScript step refuses what TypeScript asks of a placement, before the
// provider is called: the formatter is TypeScript's own contract with a
// provider that counts its calls.
func TestTypeScriptGeneratedStepRefusesBeforeTheProvider(t *testing.T) {
	ctx := context.Background()
	stages := placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations})
	contract, err := stages.GeneratedContract()
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	formatter, err := format.NewGeneratedFilesFormatter(stages.Language, contract, countedProvider{calls: &calls})
	if err != nil {
		t.Fatal(err)
	}
	output := generatedCellPage(t)
	// A placement that TypeScript admits reaches the provider once.
	if _, err := formatter.Format(ctx, format.GeneratedFilesInput{Output: output, Placement: tsOneFile(output, "types", "types.ts").Issue()}); !errors.Is(err, errReached) || calls != 1 {
		t.Fatal("an admitted placement must reach the provider", err, calls)
	}
	empty := tsOneFile(output, "types", "types.ts")
	empty.Files["empty"] = "empty.ts"
	renamed := tamperedCellPage(t, tamperedTree(func(tree *ts.AST) {
		for _, child := range tree.Nodes[tree.Root-1].Children {
			node := &tree.Nodes[child.Node-1]
			for i := range node.Attributes {
				if node.Attributes[i].Name == ts.AttributeName && node.Attributes[i].Text == "Person" {
					node.Attributes[i].Text = "Page"
				}
			}
		}
	}), "")
	for _, refused := range []struct {
		name      string
		output    g.Output
		placement tsfiles.GeneratedPlacement
		want      error
	}{
		{"a file without a definition", output, empty, source.ErrUnsupported},
		{"a path that the profile excludes", output, tsOneFile(output, "types", "types.mts"), source.ErrUnsupported},
		{"a path that is no TypeScript path", output, tsOneFile(output, "types", "types.txt"), source.ErrInvalid},
		{"two declarations of one name in one file", renamed, tsOneFile(renamed, "types", "types.ts"), source.ErrInvalid},
	} {
		calls = 0
		out, err := formatter.Format(ctx, format.GeneratedFilesInput{Output: refused.output, Placement: refused.placement.Issue()})
		if out.Valid() || calls != 0 {
			t.Errorf("%s: published=%v, calls of the provider=%d", refused.name, out.Valid(), calls)
		}
		onlyCategory(t, refused.name, err, refused.want)
	}
}

// republishingProvider takes the units that the TypeScript formatter published,
// changes them, has them analyzed again and publishes the result.
type republishingProvider struct {
	ts.Language
	stages tsfiles.Stages
	from   source.Source[ts.Language]
	entry  files.ID
	change func([]source.Unit[ts.Language]) error
}

func (p republishingProvider) Format(ctx context.Context, request format.GeneratedFilesRequest[ts.Language]) (source.Source[ts.Language], error) {
	units := p.from.Units()
	if p.change != nil {
		if err := p.change(units); err != nil {
			return source.Source[ts.Language]{}, err
		}
	}
	fresh, err := p.stages.Capture(ctx, units, p.entry, request.Input().Limits.Semantic)
	if err != nil {
		return source.Source[ts.Language]{}, err
	}
	return request.Publish(source.Publication[ts.Language]{Snapshot: fresh.Snapshot, Analysis: fresh.Evidence, Units: units, Groups: p.from.Groups(), Relations: p.from.Relations()})
}

// The TypeScript guard holds a publication to the two things that the
// formatter decides for a placement that states neither: the entry of the
// fresh snapshot and the order of the declarations of a file. The publications
// here have the admitted files, the same declarations, the same imports and the
// compiler's bindings and exports.
func TestTypeScriptGuardComparesTheEntryAndTheOrder(t *testing.T) {
	ctx := context.Background()
	output := generatedCellPage(t)
	handles := tsDeclarations(t, output)
	stages := placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations})
	genuine, err := tsfiles.NewGenerated(stages)
	if err != nil {
		t.Fatal(err)
	}
	contract, err := stages.GeneratedContract()
	if err != nil {
		t.Fatal(err)
	}
	placement := tsPlaced(t, handles, map[files.ID]files.Path{"pages": "generated/pages/page.ts", "model": "generated/model/person.ts"}, map[string]files.ID{"Page": "pages", "Cells": "pages", "Cell": "pages", "Person": "model"})
	from, err := genuine.Format(ctx, format.GeneratedFilesInput{Output: output, Placement: placement.Issue()})
	if err != nil {
		t.Fatal(err)
	}
	republish := func(entry files.ID, change func([]source.Unit[ts.Language]) error) (source.Source[ts.Language], error) {
		formatter, err := format.NewGeneratedFilesFormatter(stages.Language, contract, republishingProvider{stages: stages, from: from, entry: entry, change: change})
		if err != nil {
			t.Fatal(err)
		}
		return formatter.Format(ctx, format.GeneratedFilesInput{Output: output, Placement: placement.Issue()})
	}
	// The same units under the entry that the rule chooses are accepted, so a
	// refusal below is a refusal of the one difference.
	if out, err := republish("model", nil); err != nil || !out.Valid() {
		t.Fatal("the republished units under the entry of least identity", err)
	}
	out, err := republish("pages", nil)
	if out.Valid() {
		t.Error("a publication whose snapshot entry is not the file of least identity was accepted")
	}
	onlyCategory(t, "another snapshot entry", err, source.ErrInvalid)
	// The last two declarations of the file pages change places. Every node
	// keeps its identity, so only the order differs.
	swapped := func(units []source.Unit[ts.Language]) error {
		for i := range units {
			if units[i].File != "pages" {
				continue
			}
			tree := units[i].Syntax.Abstract
			children := tree.Nodes[tree.Root-1].Children
			last := len(children) - 1
			if last < 1 || children[last].Role != ts.RoleDeclaration || children[last-1].Role != ts.RoleDeclaration {
				return errors.New("the file pages has no two declarations at its end")
			}
			children[last], children[last-1] = children[last-1], children[last]
			syntax, err := stages.Syntax.Format(ctx, tree)
			if err != nil {
				return err
			}
			units[i].Syntax = syntax
			return nil
		}
		return errors.New("no file pages")
	}
	out, err = republish("model", swapped)
	if out.Valid() {
		t.Error("a publication with the declarations of a file in another order was accepted")
	}
	onlyCategory(t, "another order of the declarations", err, source.ErrInvalid)
}

// errReached is what reachedProvider answers.
var errReached = errors.New("the provider was reached")

// reachedProvider reports that the request passed the preparation step.
type reachedProvider struct{ ts.Language }

func (reachedProvider) Format(context.Context, format.GeneratedFilesRequest[ts.Language]) (source.Source[ts.Language], error) {
	return source.Source[ts.Language]{}, errReached
}

// The common admission accepts a file without a definition and two
// declarations of one name in one file: it interprets no scope and reads no
// name. A language whose step asks nothing beyond the admission reaches its
// provider with both. The TypeScript step refuses the first as unsupported
// coverage and the second as invalid input.
func TestCommonAdmissionLeavesScopeRulesToTheLanguage(t *testing.T) {
	ctx := context.Background()
	stages := placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations})
	typescript, err := tsfiles.NewGenerated(stages)
	if err != nil {
		t.Fatal(err)
	}
	// A language that states its placement as the view itself and has no rule
	// of its own.
	views := format.NewGeneratedPlacementIssuer[format.GeneratedPlacementView]()
	plain, err := format.NewGeneratedFilesFormatter(stages.Language, format.GeneratedFilesContract[ts.Language]{
		Print: stages.Printer.Print,
		Prepare: func(_ context.Context, preparation format.GeneratedFilesPreparation) (format.GeneratedFilesPrepared[ts.Language], error) {
			view, issued := views.Read(preparation.Input().Placement)
			if !issued {
				return format.GeneratedFilesPrepared[ts.Language]{}, source.ErrInvalid
			}
			admitted, err := preparation.Admit(view)
			if err != nil {
				return format.GeneratedFilesPrepared[ts.Language]{}, err
			}
			return format.GeneratedFilesPrepared[ts.Language]{Placement: admitted, Evidence: n.Evidence(), Check: func(context.Context, source.Publication[ts.Language]) error { return nil }}, nil
		},
	}, reachedProvider{})
	if err != nil {
		t.Fatal(err)
	}
	both := func(name string, output g.Output, placement tsfiles.GeneratedPlacement, want error) {
		t.Helper()
		view := format.GeneratedPlacementView{Files: placement.Files, Definitions: placement.Definitions}
		if out, err := plain.Format(ctx, format.GeneratedFilesInput{Output: output, Placement: views.Issue(view)}); !errors.Is(err, errReached) || out.Valid() {
			t.Errorf("%s: the common admission refused: %v", name, err)
		}
		out, err := typescript.Format(ctx, format.GeneratedFilesInput{Output: output, Placement: placement.Issue()})
		if out.Valid() {
			t.Errorf("%s: TypeScript published", name)
		}
		onlyCategory(t, name, err, want)
	}
	output := generatedCellPage(t)
	empty := tsOneFile(output, "types", "types.ts")
	empty.Files["empty"] = "empty.ts"
	both("a file without a definition", output, empty, source.ErrUnsupported)
	// Generation refuses two declarations of one name until the flat-name rule
	// leaves it, so the output is built by a target that renames one.
	renamed := tamperedCellPage(t, tamperedTree(func(tree *ts.AST) {
		for _, child := range tree.Nodes[tree.Root-1].Children {
			node := &tree.Nodes[child.Node-1]
			for i := range node.Attributes {
				if node.Attributes[i].Name == ts.AttributeName && node.Attributes[i].Text == "Person" {
					node.Attributes[i].Text = "Page"
				}
			}
		}
	}), "")
	both("two declarations of one name in one file", renamed, tsOneFile(renamed, "types", "types.ts"), source.ErrInvalid)
}
