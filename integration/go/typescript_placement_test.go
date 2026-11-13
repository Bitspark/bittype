package integration_test

import (
	"context"
	"errors"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/engine/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	compose "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/files/typescript/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	tsfiles "github.com/Bitspark/bittype/format/files/typescript/provider/go"
	printFiles "github.com/Bitspark/bittype/print/files/shared/provider/go"
	resolve "github.com/Bitspark/bittype/resolve/files/typescript/provider/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	"github.com/Bitspark/bittype/source/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

const cellDeclarations = "export interface Cell<T> { value: T }\n"

func writeText(t *testing.T, dir, path, text string) {
	t.Helper()
	file := filepath.Join(dir, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(file), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
}

func writeCellPackage(t *testing.T, dir string) {
	t.Helper()
	writeText(t, dir, "package.json", `{"type":"module"}`)
	writeText(t, dir, "node_modules/@bitspark/cell/package.json", `{"name":"@bitspark/cell","type":"module","exports":{".":{"types":"./api.d.ts"}}}`)
	writeText(t, dir, "node_modules/@bitspark/cell/api.d.ts", cellDeclarations)
}

func checkImportedUse(t *testing.T, tree ts.AST, name, module, exported string) {
	t.Helper()
	var binding ast.NodeID
	for _, child := range tree.Nodes[tree.Root-1].Children {
		node := tree.Nodes[child.Node-1]
		if node.Kind != ts.KindImport || nodeAttribute(node, ts.AttributeModule) != module {
			continue
		}
		for _, b := range node.Children {
			n := tree.Nodes[b.Node-1]
			if n.Kind == ts.KindImportBinding && nodeAttribute(n, ts.AttributeName) == exported {
				if binding != 0 {
					t.Fatal("duplicate named import", module, exported)
				}
				binding = b.Node
			}
		}
	}
	if binding == 0 {
		t.Fatal("missing import", module, exported)
	}
	count := 0
	for _, node := range tree.Nodes {
		if node.Kind != ts.KindTypeReference || nodeAttribute(node, ts.AttributeName) != name {
			continue
		}
		want := ast.NodeRef{Tree: tree.ID, Node: binding}
		if len(node.References) != 1 || node.References[0].Role != ts.RoleTarget || node.References[0].Target != want {
			t.Fatal("unbound import use", name, node.References)
		}
		count++
	}
	if count == 0 {
		t.Fatal("unused import", module, exported)
	}
}

func placementStages(t *testing.T, packages ...n.PackageTypes) tsfiles.Stages {
	t.Helper()
	text := typescriptText()
	def, err := representation.NewLanguage[ts.Language](text, "native-types-v1")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[ts.Language](def, text)
	if err != nil {
		t.Fatal(err)
	}
	settings, err := n.ImportConfiguration(packages)
	if err != nil {
		t.Fatal(err)
	}
	resolver, err := resolve.New(def, tsCompiler(t))
	if err != nil {
		t.Fatal(err)
	}
	contract, err := n.AnalysisContractFor(settings)
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[ts.Language](def, contract, analyze.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	return tsfiles.Stages{Language: def, Syntax: text, Printer: text, Checker: checker, Resolver: resolver, Analyzer: analyzer, Configuration: settings}
}

func generatedCellPage(t *testing.T) g.Output {
	t.Helper()
	library := genModule("cells", genMap(), genEntry("Cell", genType(genTuple("horos/wire", genParam(0)), "T")))
	pin, err := decl.Pin(library)
	if err != nil {
		t.Fatal(err)
	}
	imported := genTuple("horos/decl/import", data.MustEncodeText("cells"), core.NewAtom(pin[:]))
	root := genModule("consumer", genMap(genEntry("cells", imported)),
		genEntry("Cell", genType(genScalar("utf8-text"))), genEntry("Page", genType(genPage(), "T")), genEntry("Person", genType(genPerson())),
		genEntry("Cells", genType(genLocal("Page", genImported("cells", "Cell", genLocal("Person"))))))
	e, from, input := generationInput(t, root, library)
	calls := 0
	out := generatedOutput(t, e, input, typescriptText(), selectedTS(cellMapping(from, pin.String(), &calls), ""))
	if calls != 1 {
		t.Fatalf("Cell mapping invoked %d times", calls)
	}
	return out
}

// tsDeclarations names the handles of a generated TypeScript output by the names
// of the declarations at their sites. The output has one tree.
func tsDeclarations(t *testing.T, output g.Output) map[string]g.DeclarationHandle {
	t.Helper()
	entry, err := output.Snapshot().EntryTree()
	if err != nil {
		t.Fatal(err)
	}
	tree, err := representation.Unpack[ts.Language](output.Snapshot().Language(), entry)
	if err != nil {
		t.Fatal(err)
	}
	handles := map[string]g.DeclarationHandle{}
	for _, record := range output.Inventory() {
		if record.Site.Tree != tree.ID {
			t.Fatal("a declaration outside the entry tree", record.Site)
		}
		name := nodeAttribute(tree.Nodes[record.Site.Node-1], ts.AttributeName)
		if _, twice := handles[name]; twice || name == "" {
			t.Fatal("the declarations of the output are not told apart by their names", name)
		}
		handles[name] = record.Handle
	}
	return handles
}

// tsPlaced builds a placement from the paths of its files and the file of each
// declaration, which is named as tsDeclarations names it.
func tsPlaced(t *testing.T, handles map[string]g.DeclarationHandle, paths map[files.ID]files.Path, declarations map[string]files.ID) tsfiles.GeneratedPlacement {
	t.Helper()
	placement := tsfiles.GeneratedPlacement{Files: paths, Definitions: map[g.DeclarationHandle]files.ID{}}
	for name, file := range declarations {
		handle, exists := handles[name]
		if !exists {
			t.Fatal("no declaration", name)
		}
		placement.Definitions[handle] = file
	}
	return placement
}

// tsOneFile places every declaration of an output in one file.
func tsOneFile(output g.Output, file files.ID, path files.Path) tsfiles.GeneratedPlacement {
	placement := tsfiles.GeneratedPlacement{Files: map[files.ID]files.Path{file: path}, Definitions: map[g.DeclarationHandle]files.ID{}}
	for _, record := range output.Inventory() {
		placement.Definitions[record.Handle] = file
	}
	return placement
}

// Place a complete generated result in one module for native compiler tests.
func placedTS(t *testing.T, output g.Output, packages ...n.PackageTypes) ts.AST {
	t.Helper()
	formatter, err := tsfiles.NewGenerated(placementStages(t, packages...))
	if err != nil {
		t.Fatal(err)
	}
	result, err := formatter.Format(context.Background(), format.GeneratedFilesInput{Output: output, Placement: tsOneFile(output, "types", "types.ts").Issue()})
	if err != nil {
		t.Fatal(err)
	}
	return result.Units()[0].Syntax.Abstract
}

// bittype:case target=github.com/Bitspark/bittype/format/files/typescript/provider/go#GeneratedProvider.Format feature=generated-placement case=relocate-person-and-shadow-cell expects=local-import-removed-external-cell-aliased
func TestTypeScriptTypeImportPlacement(t *testing.T) {
	ctx := context.Background()
	output := generatedCellPage(t)
	tree := outputTree[ts.Language](t, output)
	stages := placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations})
	input := format.GeneratedFilesInput{Output: output}
	for _, node := range tree.Nodes {
		if node.Kind == ts.KindImport || nodeAttribute(node, ts.AttributeExported) != "" {
			t.Fatalf("file decision in generation: %+v", node)
		}
	}
	uses := output.ExternalUses()
	if len(uses) != 1 || uses[0].At.Site.Tree != tree.ID || nodeAttribute(tree.Nodes[uses[0].At.Site.Node-1], ts.AttributeName) != "Cell" {
		t.Fatal("expected one external Cell", uses)
	}
	uses[0] = g.ExternalUse{}
	exposures := output.Exposed()
	if len(exposures) != 4 {
		t.Fatal("expected four published roots", exposures)
	}
	exposures[0] = ast.NodeRef{}
	generatedFormatter, err := tsfiles.NewGenerated(stages)
	if err != nil {
		t.Fatal(err)
	}
	kernel := engine.New()
	if err = kernel.RegisterLanguage(stages.Language); err != nil {
		t.Fatal(err)
	}
	if err = kernel.RegisterGeneratedFilesFormatter(generatedFormatter); err != nil {
		t.Fatal(err)
	}
	formatter, err := tsfiles.New(stages)
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
	handles := tsDeclarations(t, output)
	const pagesPath, personPath = "generated/pages/page.ts", "generated/model/person.ts"
	for _, split := range []bool{true, false} {
		placement := tsPlaced(t, handles, map[files.ID]files.Path{"pages": pagesPath}, map[string]files.ID{"Page": "pages", "Cells": "pages", "Cell": "pages", "Person": "pages"})
		if split {
			placement.Files["model"] = personPath
			placement.Definitions[handles["Person"]] = "model"
		}
		input.Placement = placement.Issue()
		boxed, err := kernel.FormatGeneratedFiles(ctx, input)
		if err != nil {
			t.Fatalf("split=%v: %v", split, err)
		}
		out, err := source.Unpack[ts.Language](stages.Language, boxed)
		if err != nil {
			t.Fatalf("split=%v: %v", split, err)
		}
		if out.Snapshot().Same(output.Snapshot()) || !out.Evidence().Snapshot().Same(out.Snapshot()) {
			t.Fatal("expected fresh TypeScript evidence")
		}
		emitted, err := printer.Print(ctx, out)
		if err != nil {
			t.Fatal(err)
		}
		if len(emitted.Entries) != len(placement.Files) {
			t.Fatal(emitted)
		}
		// The files come in the order of their identities, so pages is the last.
		units := out.Units()
		unit := units[len(units)-1]
		pages := emitted.Entries[len(emitted.Entries)-1]
		if unit.File != "pages" || pages.Path != pagesPath || !strings.Contains(pages.Text, `from "@bitspark/cell"`) || !strings.Contains(pages.Text, "Page<Cell_1<Person>>") {
			t.Fatal(pages)
		}
		checkImportedUse(t, unit.Syntax.Abstract, "Cell_1", "@bitspark/cell", "Cell")
		if split {
			if emitted.Entries[0].Path != personPath || !strings.Contains(emitted.Entries[0].Text, "interface Person") || strings.Contains(pages.Text, "interface Person") || !strings.Contains(pages.Text, `from "../model/person.js"`) {
				t.Fatal(emitted)
			}
			checkImportedUse(t, unit.Syntax.Abstract, "Person", "../model/person.js", "Person")
			rels := out.Relations()
			if len(rels) != 1 || rels[0].From != "pages" || rels[0].To != "model" || rels[0].At == nil || rels[0].At.Tree != unit.Syntax.Abstract.ID || unit.Syntax.Abstract.Nodes[rels[0].At.Node-1].Kind != ts.KindImport {
				t.Fatal(rels)
			}
		} else {
			if !strings.Contains(pages.Text, "interface Person") || strings.Contains(pages.Text, "person.js") || len(out.Relations()) != 0 {
				t.Fatal(pages.Text, out.Relations())
			}
			imports := 0
			personRef := false
			for _, node := range unit.Syntax.Abstract.Nodes {
				if node.Kind == ts.KindImport {
					imports++
				}
				if node.Kind == ts.KindTypeReference && nodeAttribute(node, ts.AttributeName) == "Person" && len(node.References) == 1 {
					target := node.References[0].Target
					personRef = target.Tree == unit.Syntax.Abstract.ID && unit.Syntax.Abstract.Nodes[target.Node-1].Kind == ts.KindInterface
				}
			}
			if imports != 1 || !personRef {
				t.Fatal("Person must refer directly to its declaration", imports, personRef)
			}
		}
		dir := t.TempDir()
		writeCellPackage(t, dir)
		for _, file := range emitted.Entries {
			writeText(t, dir, string(file.Path), file.Text)
		}
		writeText(t, dir, "use.ts", "import type { Cells, Cell } from './generated/pages/page.js';\nconst cells: Cells = {items:[{value:{name:'Ada', active:true}}]};\nvoid cells;\n")
		compiler := tsCompiler(t)
		command(t, dir, compiler.NodePath, filepath.Join(filepath.Dir(compiler.TypeScriptPath), "tsc.js"), "--module", "NodeNext", "--strict", "--exactOptionalPropertyTypes", "--noEmit", "use.ts")

		again, err := formatter.Format(ctx, format.FilesInput[ts.Language]{Abstract: out.SemanticsView(), Snapshot: out.Snapshot(), Evidence: out.Evidence()})
		if err != nil {
			t.Fatal("reuse placed source", err)
		}
		reprinted, err := printer.Print(ctx, again)
		if err != nil || !reflect.DeepEqual(reprinted, emitted) {
			t.Fatal("placement did not round-trip", err)
		}
		// Neither placement names an entry. The entry of each fresh snapshot is
		// the file of least identity: model when the placement has it.
		for _, formatted := range []source.Source[ts.Language]{out, again} {
			least := formatted.Units()[0].Syntax.Abstract
			if formatted.Entry() != (ast.NodeRef{Tree: least.ID, Node: least.Root}) {
				t.Fatal("the snapshot entry is not the file of least identity", split, formatted.Entry())
			}
		}
		stale := format.FilesInput[ts.Language]{Abstract: out.SemanticsView(), Snapshot: output.Snapshot(), Evidence: out.Evidence()}
		if bad, err := formatter.Format(ctx, stale); err == nil || len(bad.Units()) != 0 {
			t.Fatal("stale receipt accepted")
		}

	}
	// A same-named local declaration cannot satisfy a missing package export.
	unavailable, err := tsfiles.NewGenerated(placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: "export interface Other {}"}))
	if err != nil {
		t.Fatal(err)
	}
	if bad, err := unavailable.Format(ctx, input); err == nil || bad.Valid() {
		t.Fatal("local Cell captured an unavailable export")
	}
	other := generatedCellPage(t)
	incomplete := tsOneFile(output, "pages", pagesPath)
	delete(incomplete.Definitions, handles["Person"])
	undeclared := tsOneFile(output, "pages", pagesPath)
	undeclared.Definitions[handles["Person"]] = "model"
	for name, placement := range map[string]tsfiles.GeneratedPlacement{
		"a declaration without a file":  incomplete,
		"a file that is not declared":   undeclared,
		"the handles of another output": tsOneFile(other, "pages", pagesPath),
	} {
		input.Placement = placement.Issue()
		if bad, err := generatedFormatter.Format(ctx, input); err == nil || len(bad.Units()) != 0 {
			t.Fatal("invalid placement accepted:", name)
		}
	}
}

func TestGeneratedFilesPreservePublicationIntent(t *testing.T) {
	ctx := context.Background()
	root := genModule("visibility", genMap(), genEntry("Hidden", genType(genPerson())), genEntry("Visible", genType(genLocal("Hidden"))))
	e, _, input := generationInput(t, root)
	output := generatedOutput(t, e, input, typescriptText(), func(reg *compose.Registry, input transform.Input, to representation.Language) (compose.Request, error) {
		request, err := tsRequest(reg, input, to)
		if err != nil {
			return request, err
		}
		treeEnvelope, err := input.Source.EntryTree()
		if err != nil {
			return request, err
		}
		tree, err := representation.Unpack[h.Language](input.Source.Language(), treeEnvelope)
		if err != nil {
			return request, err
		}
		for _, site := range request.Roots {
			if nodeAttribute(tree.Nodes[site.Node-1], h.AttributeName) == "Visible" {
				request.Roots = []ast.NodeRef{site}
				return request, nil
			}
		}
		t.Fatal("Visible root missing")
		return request, nil
	})
	tree := outputTree[ts.Language](t, output)
	sites := map[string]ast.NodeRef{}
	for _, c := range tree.Nodes[tree.Root-1].Children {
		sites[nodeAttribute(tree.Nodes[c.Node-1], ts.AttributeName)] = ast.NodeRef{Tree: tree.ID, Node: c.Node}
	}
	if got := output.Exposed(); len(got) != 1 || got[0] != sites["Visible"] {
		t.Fatal("wrong public declarations", got)
	}
	placed := placedTS(t, output)
	for _, c := range placed.Nodes[placed.Root-1].Children {
		n := placed.Nodes[c.Node-1]
		name := nodeAttribute(n, ts.AttributeName)
		if (name == "Visible") != (nodeAttribute(n, ts.AttributeExported) == "true") {
			t.Fatal("public exposure changed", name)
		}
	}
	formatter, err := tsfiles.NewGenerated(placementStages(t))
	if err != nil {
		t.Fatal(err)
	}
	separate := tsPlaced(t, tsDeclarations(t, output), map[files.ID]files.Path{"visible": "visible.ts", "hidden": "hidden.ts"}, map[string]files.ID{"Visible": "visible", "Hidden": "hidden"})
	if bad, err := formatter.Format(ctx, format.GeneratedFilesInput{Output: output, Placement: separate.Issue()}); err == nil || bad.Valid() {
		t.Fatal("private dependency exported to satisfy placement")
	}
}

func importedUnit(t *testing.T, text string) source.Unit[ts.Language] {
	t.Helper()
	tree, err := typescriptText().Parse(context.Background(), parse.Input{SourceID: "input", TreeID: "input", Text: text})
	if err != nil {
		t.Fatal(err)
	}
	return source.Unit[ts.Language]{File: "input", Path: "input.ts", Source: "input", Syntax: tree}
}

// bittype:case target=github.com/Bitspark/bittype/format/files/typescript/provider/go#Provider.Format feature=type-import-placement case=same-export-name-and-lexical-binder expects=distinct-aliases-and-deduplicated-imports
func TestTypeScriptTypeImportAliases(t *testing.T) {
	ctx := context.Background()
	stages := placementStages(t, n.PackageTypes{Module: "@example/a", Declaration: "export interface Element<T> { a: T }"}, n.PackageTypes{Module: "@example/b", Declaration: "export interface Element<T> { b: T }"})
	u := importedUnit(t, `import type { Element as X } from "@example/a";
import type { Element as Y } from "@example/b";
export type Pair<Element> = [X<string>, Y<Element>, X<string>];`)
	in, err := stages.Capture(ctx, []source.Unit[ts.Language]{u}, u.File, semantic.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	formatter, err := tsfiles.New(stages)
	if err != nil {
		t.Fatal(err)
	}
	out, err := formatter.Format(ctx, in)
	if err != nil {
		t.Fatal(err)
	}
	tree := out.Units()[0].Syntax.Abstract
	checkImportedUse(t, tree, "Element_1", "@example/a", "Element")
	checkImportedUse(t, tree, "Element_2", "@example/b", "Element")
	text, err := stages.Printer.Print(ctx, out.Units()[0].Syntax)
	if err != nil || !strings.Contains(text, "[Element_1<string>, Element_2<Element>, Element_1<string>]") || strings.Count(text, "import type") != 2 {
		t.Fatal(err, text)
	}
	if len(out.Relations()) != 0 {
		t.Fatal("package imports invented local files", out.Relations())
	}
}

func TestTypeScriptImportEnvironmentIsReceiptBound(t *testing.T) {
	ctx := context.Background()
	u := importedUnit(t, `import type { Cell } from "@bitspark/cell"; export type Use = Cell<string>;`)
	check := func(stages tsfiles.Stages) (format.FilesInput[ts.Language], error) {
		return stages.Capture(ctx, []source.Unit[ts.Language]{u}, u.File, semantic.Limits{})
	}
	for _, stages := range []tsfiles.Stages{
		placementStages(t),
		placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: "export interface Other {}"}),
		placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: "export type Cell = import('./ambient').Missing;"}),
		placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: "// @ts-nocheck\nexport type Cell<T> = Missing<T>;"}),
		placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: "// @ts-ignore\nexport type Cell<T> = Missing<T>;"}),
	} {
		if in, err := check(stages); err == nil || in.Abstract.Valid() {
			t.Fatal("unresolved dependency or export admitted")
		}
	}
	stages := placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations})
	in, err := check(stages)
	if err != nil {
		t.Fatal(err)
	}
	changed := placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: "export interface Cell<T> { updated: T }"})
	formatter, err := tsfiles.New(changed)
	if err != nil {
		t.Fatal(err)
	}
	if out, err := formatter.Format(ctx, in); err == nil || len(out.Units()) != 0 {
		t.Fatal("accepted receipt checked against different package text")
	}
	resolved, err := stages.Resolver.Resolve(ctx, in.Snapshot, semantic.Options{Configuration: stages.Configuration})
	if err != nil {
		t.Fatal(err)
	}
	if receipt, err := stages.Analyzer.Analyze(ctx, analysis.Input{Snapshot: in.Snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: changed.Configuration}}); err == nil || receipt.HasSemantics() {
		t.Fatal("resolution was relabelled with different package text")
	}
}

func TestTypeScriptPlacementPreservesModuleScope(t *testing.T) {
	ctx := context.Background()
	stages := placementStages(t)
	formatter, err := tsfiles.New(stages)
	if err != nil {
		t.Fatal(err)
	}
	u := importedUnit(t, "type Hidden = string; export type Visible = Hidden;")
	in, err := stages.Capture(ctx, []source.Unit[ts.Language]{u}, u.File, semantic.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	organization := tsfiles.Organization{}
	for _, child := range u.Syntax.Abstract.Nodes[u.Syntax.Abstract.Root-1].Children {
		site := ast.NodeRef{Tree: u.Syntax.Abstract.ID, Node: child.Node}
		placement := tsfiles.Placement{File: "public", Path: "public.ts", Declarations: []ast.NodeRef{site}}
		if nodeAttribute(u.Syntax.Abstract.Nodes[child.Node-1], ts.AttributeName) == "Hidden" {
			placement.File = "private"
			placement.Path = "private.ts"
		}
		organization.Units = append(organization.Units, placement)
	}
	setting, _ := organization.Setting()
	in.Configuration = []semantic.Setting{setting}
	if out, err := formatter.Format(ctx, in); err == nil || len(out.Units()) != 0 {
		t.Fatal("private declaration became public through placement")
	}
	script := importedUnit(t, "type Global = string;")
	in, err = stages.Capture(ctx, []source.Unit[ts.Language]{script}, script.File, semantic.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := formatter.Format(ctx, in)
	if err == nil || len(out.Units()) != 0 {
		t.Fatal("script presented as an ESM module")
	}
	// A script is well formed and outside this profile, so the refusal is
	// unsupported coverage and no invalid input.
	if !errors.Is(err, source.ErrUnsupported) || errors.Is(err, source.ErrInvalid) {
		t.Fatal("a script is unsupported coverage", err)
	}
}

// The compiler helper is JavaScript and spells the graph names itself. Three of
// the names that Go reads through the vocabulary are read by no other test: the
// builtin kind of a binding's target, the parameter kind of a declaration opened
// by id and the native-type relation of a type opened by id.
func TestTypeScriptHelperNamesMatchVocabulary(t *testing.T) {
	ctx := context.Background()
	stages := placementStages(t)
	u := importedUnit(t, "export type Items<T> = Array<T>;")
	in, err := stages.Capture(ctx, []source.Unit[ts.Language]{u}, u.File, semantic.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	formatter, err := tsfiles.New(stages)
	if err != nil {
		t.Fatal(err)
	}
	if out, err := formatter.Format(ctx, in); err != nil || len(out.Units()) != 1 {
		t.Fatal("binding to a builtin constructor", err)
	}
	model, err := n.Open(in.Evidence)
	if err != nil {
		t.Fatal(err)
	}
	items, err := model.Lookup(u.Syntax.Abstract.ID, "Items")
	if err != nil || len(items.Parameters()) != 1 {
		t.Fatal("generic alias", err)
	}
	if opened, err := model.Declaration(items.Parameters()[0].ID()); err != nil || opened.Entity().Name != "T" {
		t.Fatal("parameter declaration by id", err)
	}
	if opened, err := model.Type(items.Type().ID()); err != nil || !opened.Same(items.Type()) {
		t.Fatal("native type by id", err)
	}
}
