package integration_test

import (
	"context"
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
	contract, err := analyze.ContractFor(settings)
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[ts.Language](def, contract, analyze.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	return tsfiles.Stages{Language: def, Syntax: text, Printer: text, Checker: checker, Resolver: resolver, Analyzer: analyzer, Configuration: settings}
}

func generatedCellPage(t *testing.T) ts.AST {
	t.Helper()
	library := genModule("cells", genMap(), genEntry("Cell", genType(genTuple("horos/wire", genParam(0)), "T")))
	pin, err := decl.Pin(library)
	if err != nil {
		t.Fatal(err)
	}
	imported := genTuple("horos/decl/import", data.MustEncodeText("cells"), core.NewAtom(pin[:]))
	root := genModule("consumer", genMap(genEntry("cells", imported)),
		genEntry("Page", genType(genPage(), "T")), genEntry("Person", genType(genPerson())),
		genEntry("Cells", genType(genLocal("Page", genImported("cells", "Cell", genLocal("Person"))))))
	e, from, input := generationInput(t, root, library)
	calls := 0
	return generated(t, e, input, typescriptText(), selectedTS(cellMapping(from, pin.String(), &calls), ""))
}

func captureGenerated(t *testing.T, stages tsfiles.Stages, tree ts.AST) format.FilesInput[ts.Language] {
	t.Helper()
	syntax, err := stages.Syntax.Format(context.Background(), tree)
	if err != nil {
		t.Fatal(err)
	}
	u := source.Unit[ts.Language]{File: "logical", Path: "logical.ts", Source: "logical", Syntax: syntax}
	input, err := stages.Capture(context.Background(), []source.Unit[ts.Language]{u}, u.File, semantic.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	return input
}

// bittype:case target=github.com/Bitspark/bittype/format/files/typescript/provider/go#Provider.Format feature=type-import-placement case=relocate-person-in-page expects=local-import-removed-package-import-retained
func TestTypeScriptTypeImportPlacement(t *testing.T) {
	ctx := context.Background()
	tree := generatedCellPage(t)
	stages := placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations})
	input := captureGenerated(t, stages, tree)
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
	sites := map[string]ast.NodeRef{}
	for _, child := range tree.Nodes[tree.Root-1].Children {
		if child.Role == ts.RoleDeclaration {
			sites[nodeAttribute(tree.Nodes[child.Node-1], ts.AttributeName)] = ast.NodeRef{Tree: tree.ID, Node: child.Node}
		}
	}
	page := tsfiles.Placement{File: "pages", Path: "generated/pages/page.ts", Declarations: []ast.NodeRef{sites["Page"], sites["Cells"]}}
	person := tsfiles.Placement{File: "model", Path: "generated/model/person.ts", Declarations: []ast.NodeRef{sites["Person"]}}
	for _, split := range []bool{true, false} {
		org := tsfiles.Organization{Entry: page.File, Units: []tsfiles.Placement{page}}
		if split {
			org.Units = append(org.Units, person)
		} else {
			org.Units[0].Declarations = append(org.Units[0].Declarations, sites["Person"])
		}
		setting, err := org.Setting()
		if err != nil {
			t.Fatal(err)
		}
		input.Configuration = []semantic.Setting{setting}
		out, err := formatter.Format(ctx, input)
		if err != nil {
			t.Fatalf("split=%v: %v", split, err)
		}
		if out.Snapshot().Same(input.Snapshot) || !out.Evidence().Snapshot().Same(out.Snapshot()) {
			t.Fatal("expected fresh TypeScript evidence")
		}
		emitted, err := printer.Print(ctx, out)
		if err != nil {
			t.Fatal(err)
		}
		if len(emitted.Entries) != len(org.Units) {
			t.Fatal(emitted)
		}
		unit := out.Units()[0]
		pages := emitted.Entries[0]
		if pages.Path != page.Path || !strings.Contains(pages.Text, `from "@bitspark/cell"`) || !strings.Contains(pages.Text, "Page<Cell<Person>>") {
			t.Fatal(pages)
		}
		checkImportedUse(t, unit.Syntax.Abstract, "Cell", "@bitspark/cell", "Cell")
		if split {
			if emitted.Entries[1].Path != person.Path || !strings.Contains(emitted.Entries[1].Text, "interface Person") || strings.Contains(pages.Text, "interface Person") || !strings.Contains(pages.Text, `from "../model/person.js"`) {
				t.Fatal(emitted)
			}
			checkImportedUse(t, unit.Syntax.Abstract, "Person", "../model/person.js", "Person")
			rels := out.Relations()
			if len(rels) != 1 || rels[0].From != page.File || rels[0].To != person.File || rels[0].At == nil || rels[0].At.Tree != unit.Syntax.Abstract.ID || unit.Syntax.Abstract.Nodes[rels[0].At.Node-1].Kind != ts.KindImport {
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
		writeText(t, dir, "use.ts", "import type { Cells } from './generated/pages/page.js';\nconst cells: Cells = {items:[{value:{name:'Ada', active:true}}]};\nvoid cells;\n")
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
		stale := input
		stale.Snapshot = out.Snapshot()
		if bad, err := formatter.Format(ctx, stale); err == nil || len(bad.Units()) != 0 {
			t.Fatal("stale receipt accepted")
		}
		if split {
			// Move the already imported Person from the actual two-file Source.
			// The new placement names sites from that Source's own receipt.
			for _, u := range out.Units() {
				tree := u.Syntax.Abstract
				for _, child := range tree.Nodes[tree.Root-1].Children {
					if child.Role == ts.RoleDeclaration {
						sites[nodeAttribute(tree.Nodes[child.Node-1], ts.AttributeName)] = ast.NodeRef{Tree: tree.ID, Node: child.Node}
					}
				}
			}
			page.Declarations = []ast.NodeRef{sites["Page"], sites["Cells"]}
			person.Declarations = []ast.NodeRef{sites["Person"]}
			input = format.FilesInput[ts.Language]{Abstract: out.SemanticsView(), Snapshot: out.Snapshot(), Evidence: out.Evidence()}
		}
	}
	for _, units := range [][]tsfiles.Placement{{page}, {page, person, person}, {{File: page.File, Path: page.Path, Declarations: []ast.NodeRef{{Tree: "foreign", Node: 1}}}}} {
		setting, _ := (tsfiles.Organization{Entry: page.File, Units: units}).Setting()
		input.Configuration = []semantic.Setting{setting}
		if bad, err := formatter.Format(ctx, input); err == nil || len(bad.Units()) != 0 {
			t.Fatal("invalid placement accepted")
		}
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
	for _, stages := range []tsfiles.Stages{placementStages(t), placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: "export interface Other {}"}), placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: "export type Cell = import('./ambient').Missing;"})} {
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
	organization := tsfiles.Organization{Entry: "public"}
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
	if out, err := formatter.Format(ctx, in); err == nil || len(out.Units()) != 0 {
		t.Fatal("script presented as an ESM module")
	}
}
