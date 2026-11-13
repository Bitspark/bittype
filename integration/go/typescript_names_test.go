package integration_test

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/Bitspark/bittype/contracts/format/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	tsfiles "github.com/Bitspark/bittype/format/files/typescript/provider/go"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	"github.com/Bitspark/bittype/source/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// These qualify what TypeScript does with names once generation states no rule
// about names across declarations (#273): the TypeScript Files formatter
// decides them per file.

const (
	// Both declarations are spelled Page.
	sameNames = "module collision; type page = text; type Page = bool;"
	// The parameter of Page is spelled like the declaration T, and the list
	// uses the parameter.
	parameterNamed = "module collision; type T = text; type Page<T> = list<T>;"
	// The declaration t and the parameter of Box are both spelled T, and Box
	// uses both.
	capturedNames = "module collision; type t = text; type Box<T> = { first: t; second: T; };"
	// Both declarations are spelled Page, and the record uses the alias.
	importedSameName = "module collision; type page = text; type Page = { inner: page; };"
)

// namesOutput generates TypeScript from one horos module.
func namesOutput(t *testing.T, text string) g.Output {
	t.Helper()
	parsed, err := parsehoros.Parse(context.Background(), parse.Input{SourceID: "names.horos", TreeID: "canonical", Text: text})
	if err != nil {
		t.Fatal(err)
	}
	_, _, input := generationSource(t, parsed.Abstract)
	output, err := runGeneration(t, input, nativeChecker(t, typescriptText()), tsRequest)
	if err != nil {
		t.Fatal("generation refused the names: ", err)
	}
	return output
}

// namesTree is the one tree of an output, read without the helpers that
// record it.
func namesTree(t *testing.T, output g.Output) ts.AST {
	t.Helper()
	entry, err := output.Snapshot().EntryTree()
	if err != nil {
		t.Fatal(err)
	}
	tree, err := representation.Unpack[ts.Language](output.Snapshot().Language(), entry)
	if err != nil {
		t.Fatal(err)
	}
	return tree
}

// typeUse is one name use below a node: how it is spelled and what it is bound to.
type typeUse struct {
	name   string
	target ast.NodeRef
}

func typeUses(tree ts.AST, from ast.NodeID) []typeUse {
	var out []typeUse
	node := tree.Nodes[from-1]
	if node.Kind == ts.KindTypeReference {
		use := typeUse{name: nodeAttribute(node, ts.AttributeName)}
		if len(node.References) == 1 {
			use.target = node.References[0].Target
		}
		out = append(out, use)
	}
	for _, child := range node.Children {
		out = append(out, typeUses(tree, child.Node)...)
	}
	return out
}

// declarationOf returns the record of the one declaration of a kind.
func declarationOf(t *testing.T, output g.Output, tree ts.AST, kind ts.Kind) g.OwnedDeclaration {
	t.Helper()
	var found []g.OwnedDeclaration
	for _, record := range output.Inventory() {
		if tree.Nodes[record.Site.Node-1].Kind == kind {
			found = append(found, record)
		}
	}
	if len(found) != 1 {
		t.Fatalf("the output has %d declarations of the kind %s", len(found), kind)
	}
	return found[0]
}

// namesFormat formats an output under a placement and returns each file's text
// by its identity.
func namesFormat(t *testing.T, output g.Output, placement tsfiles.GeneratedPlacement) (source.Source[ts.Language], map[files.ID]string, error) {
	t.Helper()
	stages := placementStages(t)
	formatter, err := tsfiles.NewGenerated(stages)
	if err != nil {
		t.Fatal(err)
	}
	out, err := formatter.Format(context.Background(), format.GeneratedFilesInput{Output: output, Placement: placement.Issue()})
	if err != nil {
		return out, nil, err
	}
	texts := map[files.ID]string{}
	for _, unit := range out.Units() {
		text, err := stages.Printer.Print(context.Background(), unit.Syntax)
		if err != nil {
			t.Fatal(err)
		}
		texts[unit.File] = text
	}
	return out, texts, nil
}

// unitTree returns the tree of one published file.
func unitTree(t *testing.T, out source.Source[ts.Language], file files.ID) ts.AST {
	t.Helper()
	for _, unit := range out.Units() {
		if unit.File == file {
			return unit.Syntax.Abstract
		}
	}
	t.Fatal("no file", file)
	return ts.AST{}
}

// Two declarations of one name generate for TypeScript. Placed in one file
// they are refused as invalid input: no slice admits that placement. Placed in
// two files they are published, each under its name.
func TestTypeScriptPlacesTwoDeclarationsOfOneNameInTwoFiles(t *testing.T) {
	output := namesOutput(t, sameNames)
	tree := namesTree(t, output)
	inventory := output.Inventory()
	if len(inventory) != 2 || inventory[0].Site == inventory[1].Site {
		t.Fatal("the output owns two declarations at two sites", len(inventory))
	}
	for _, record := range inventory {
		if name := nodeAttribute(tree.Nodes[record.Site.Node-1], ts.AttributeName); name != "Page" {
			t.Fatal("both declarations are spelled Page", name)
		}
	}
	out, _, err := namesFormat(t, output, tsOneFile(output, "types", "types.ts"))
	if out.Valid() {
		t.Fatal("two declarations of one name were published in one file")
	}
	onlyCategory(t, "two declarations of one name in one file", err, source.ErrInvalid)
	apart := tsfiles.GeneratedPlacement{
		Files:       map[files.ID]files.Path{"text": "text.ts", "flag": "flag.ts"},
		Definitions: map[g.DeclarationHandle]files.ID{inventory[0].Handle: "text", inventory[1].Handle: "flag"},
	}
	out, texts, err := namesFormat(t, output, apart)
	if err != nil || !out.Valid() {
		t.Fatal("two declarations of one name in two files were refused: ", err)
	}
	if texts["text"] != "export type Page = string;\n" || texts["flag"] != "export type Page = boolean;\n" {
		t.Fatalf("the two files: %q", texts)
	}
}

// A type parameter that is named like a declaration and captures no use of it
// generates and is published in one file. The element of the list is linked to
// the parameter, in the output and in the published file.
func TestTypeScriptParameterNamedLikeADeclaration(t *testing.T) {
	output := namesOutput(t, parameterNamed)
	tree := namesTree(t, output)
	var generic g.OwnedDeclaration
	for _, record := range output.Inventory() {
		if len(typeUses(tree, record.Site.Node)) == 1 {
			generic = record
		}
	}
	if !generic.Handle.Valid() {
		t.Fatal("no declaration with one name use")
	}
	parameter := func(tree ts.AST, declaration ast.NodeID) ast.NodeRef {
		for _, child := range tree.Nodes[declaration-1].Children {
			if child.Role == ts.RoleTypeParameter {
				return ast.NodeRef{Tree: tree.ID, Node: child.Node}
			}
		}
		t.Fatal("the declaration has no type parameter")
		return ast.NodeRef{}
	}
	if use := typeUses(tree, generic.Site.Node)[0]; use.name != "T" || use.target != parameter(tree, generic.Site.Node) {
		t.Fatal("the element of the list is not linked to the parameter in the output", use)
	}
	out, texts, err := namesFormat(t, output, tsOneFile(output, "types", "types.ts"))
	if err != nil || !out.Valid() {
		t.Fatal("a parameter that captures nothing was refused: ", err)
	}
	if texts["types"] != "export type Page<T> = T[];\n\nexport type T = string;\n" {
		t.Fatalf("the file: %q", texts["types"])
	}
	placed := unitTree(t, out, "types")
	for _, child := range placed.Nodes[placed.Root-1].Children {
		uses := typeUses(placed, child.Node)
		if len(uses) == 1 && (uses[0].name != "T" || uses[0].target != parameter(placed, child.Node)) {
			t.Fatal("the element of the list is not linked to the parameter in the published file", uses[0])
		}
	}
}

// The example of point 6 of #273. A declaration and a type parameter are both
// spelled T, and the record uses both. Placed in one file the parameter would
// capture the use of the declaration: TypeScript refuses the placement as
// unsupported coverage, and no Source is returned. Placed in two files it is
// published, the import is aliased, and the use binds to the declaration.
func TestTypeScriptParameterCapturesADeclarationOfItsFile(t *testing.T) {
	output := namesOutput(t, capturedNames)
	tree := namesTree(t, output)
	alias, record := declarationOf(t, output, tree, ts.KindTypeAlias), declarationOf(t, output, tree, ts.KindInterface)
	var parameter ast.NodeRef
	for _, child := range tree.Nodes[record.Site.Node-1].Children {
		if child.Role == ts.RoleTypeParameter {
			parameter = ast.NodeRef{Tree: tree.ID, Node: child.Node}
		}
	}
	// Generation recorded the two destinations. Both uses are spelled T.
	uses := typeUses(tree, record.Site.Node)
	if len(uses) != 2 || uses[0] != (typeUse{"T", alias.Site}) || uses[1] != (typeUse{"T", parameter}) {
		t.Fatal("the record uses the declaration T and then its parameter T", uses)
	}
	out, _, err := namesFormat(t, output, tsOneFile(output, "types", "types.ts"))
	if out.Valid() {
		t.Fatal("a use that the parameter of its declaration captures was published")
	}
	onlyCategory(t, "a parameter that captures a declaration of its file", err, source.ErrUnsupported)
	if errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal("the refusal is the bridge's and not the formatter's: ", err)
	}
	apart := tsfiles.GeneratedPlacement{
		Files:       map[files.ID]files.Path{"box": "box.ts", "t": "t.ts"},
		Definitions: map[g.DeclarationHandle]files.ID{record.Handle: "box", alias.Handle: "t"},
	}
	out, texts, err := namesFormat(t, output, apart)
	if err != nil || !out.Valid() {
		t.Fatal("the two files were refused: ", err)
	}
	const box = "import type { T as T_1 } from \"./t.js\";\n\nexport interface Box<T> {\n  first: T_1;\n  second: T;\n}\n"
	if texts["box"] != box || texts["t"] != "export type T = string;\n" {
		t.Fatalf("the two files: %q", texts)
	}
	// The use of the declaration is bound to the aliased import, and the use of
	// the parameter to the parameter. The guard compared both with the compiler.
	placed := unitTree(t, out, "box")
	checkImportedUse(t, placed, "T_1", "./t.js", "T")
	for _, child := range placed.Nodes[placed.Root-1].Children {
		if placed.Nodes[child.Node-1].Kind != ts.KindInterface {
			continue
		}
		var own ast.NodeRef
		for _, grand := range placed.Nodes[child.Node-1].Children {
			if grand.Role == ts.RoleTypeParameter {
				own = ast.NodeRef{Tree: placed.ID, Node: grand.Node}
			}
		}
		uses := typeUses(placed, child.Node)
		if len(uses) != 2 || uses[1] != (typeUse{"T", own}) {
			t.Fatal("the use of the parameter is not bound to the parameter", uses)
		}
	}
}

// Two declarations of one name of which one uses the other. In one file they
// are invalid input. In two files the file that uses the other declaration
// imports it under an alias, because its own declaration has the name.
func TestTypeScriptAliasesAnImportOfADeclarationOfTheSameName(t *testing.T) {
	output := namesOutput(t, importedSameName)
	tree := namesTree(t, output)
	alias, record := declarationOf(t, output, tree, ts.KindTypeAlias), declarationOf(t, output, tree, ts.KindInterface)
	if uses := typeUses(tree, record.Site.Node); len(uses) != 1 || uses[0] != (typeUse{"Page", alias.Site}) {
		t.Fatal("the record uses the alias of its own name", uses)
	}
	out, _, err := namesFormat(t, output, tsOneFile(output, "types", "types.ts"))
	if out.Valid() {
		t.Fatal("two declarations of one name were published in one file")
	}
	onlyCategory(t, "two declarations of one name in one file", err, source.ErrInvalid)
	apart := tsfiles.GeneratedPlacement{
		Files:       map[files.ID]files.Path{"record": "record.ts", "alias": "alias.ts"},
		Definitions: map[g.DeclarationHandle]files.ID{record.Handle: "record", alias.Handle: "alias"},
	}
	out, texts, err := namesFormat(t, output, apart)
	if err != nil || !out.Valid() {
		t.Fatal("the two files were refused: ", err)
	}
	const want = "import type { Page as Page_1 } from \"./alias.js\";\n\nexport interface Page {\n  inner: Page_1;\n}\n"
	if texts["record"] != want || texts["alias"] != "export type Page = string;\n" {
		t.Fatalf("the two files: %q", texts)
	}
	checkImportedUse(t, unitTree(t, out, "record"), "Page_1", "./alias.js", "Page")
}

// A type parameter takes its name in the file also when no use is written with
// it. The existing source imports the declaration T under another name beside
// a record whose parameter T is unused. Formatting it again decides the alias
// anew: the import cannot be called T.
func TestTypeScriptAliasAvoidsATypeParameterWithoutAUse(t *testing.T) {
	ctx := context.Background()
	stages := placementStages(t)
	formatter, err := tsfiles.New(stages)
	if err != nil {
		t.Fatal(err)
	}
	unit := func(id, text string) source.Unit[ts.Language] {
		parsed, err := typescriptText().Parse(ctx, parse.Input{SourceID: ast.SourceID(id), TreeID: ast.TreeID(id), Text: text})
		if err != nil {
			t.Fatal(err)
		}
		return source.Unit[ts.Language]{File: files.ID(id), Path: files.Path(id + ".ts"), Source: ast.SourceID(id), Syntax: parsed}
	}
	units := []source.Unit[ts.Language]{unit("a", "export type T = string;"), unit("user", "import type { T as U } from \"./a.js\"; export interface Box<T> { own: U }")}
	in, err := stages.Capture(ctx, units, "a", semantic.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := formatter.Format(ctx, in)
	if err != nil {
		t.Fatal(err)
	}
	text, err := stages.Printer.Print(ctx, out.Units()[1].Syntax)
	if err != nil {
		t.Fatal(err)
	}
	if out.Units()[1].File != "user" || text != "import type { T as T_1 } from \"./a.js\";\n\nexport interface Box<T> {\n  own: T_1;\n}\n" {
		t.Fatalf("the file user: %q", text)
	}
	checkImportedUse(t, out.Units()[1].Syntax.Abstract, "T_1", "./a.js", "T")
}

// namesOutputWith generates TypeScript from one horos module with a request of
// the caller, for a mapping that no builtin library has.
func namesOutputWith(t *testing.T, text string, build requestBuilder) g.Output {
	t.Helper()
	parsed, err := parsehoros.Parse(context.Background(), parse.Input{SourceID: "names.horos", TreeID: "canonical", Text: text})
	if err != nil {
		t.Fatal(err)
	}
	_, _, input := generationSource(t, parsed.Abstract)
	output, err := runGeneration(t, input, nativeChecker(t, typescriptText()), build)
	if err != nil {
		t.Fatal("generation refused the names: ", err)
	}
	return output
}

// namesExisting captures existing TypeScript, one file per text under the
// path of its key, and formats it. A placement gives each new file the
// declarations that it holds, each named as file.Name. Without a placement
// the files stay as they are. It returns the text of each published file.
// The probes of the independent review of #304 are the source of this helper
// and of the tests below that use it.
func namesExisting(t *testing.T, texts map[string]string, placement map[string][]string) (map[files.ID]string, error) {
	t.Helper()
	ctx := context.Background()
	stages := placementStages(t)
	formatter, err := tsfiles.New(stages)
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for id := range texts {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	var units []source.Unit[ts.Language]
	sites := map[string]ast.NodeRef{}
	for _, id := range ids {
		parsed, err := typescriptText().Parse(ctx, parse.Input{SourceID: ast.SourceID(id), TreeID: ast.TreeID(id), Text: texts[id]})
		if err != nil {
			t.Fatal(id, err)
		}
		tree := parsed.Abstract
		for _, child := range tree.Nodes[tree.Root-1].Children {
			if child.Role == ts.RoleDeclaration {
				sites[id+"."+nodeAttribute(tree.Nodes[child.Node-1], ts.AttributeName)] = ast.NodeRef{Tree: tree.ID, Node: child.Node}
			}
		}
		units = append(units, source.Unit[ts.Language]{File: files.ID(id), Path: files.Path(id + ".ts"), Source: ast.SourceID(id), Syntax: parsed})
	}
	in, err := stages.Capture(ctx, units, files.ID(ids[0]), semantic.Limits{})
	if err != nil {
		t.Fatal("the existing source must be admitted before it is placed: ", err)
	}
	if placement != nil {
		var organization tsfiles.Organization
		var placed []string
		for file := range placement {
			placed = append(placed, file)
		}
		slices.Sort(placed)
		for _, file := range placed {
			unit := tsfiles.Placement{File: files.ID(file), Path: files.Path(file + ".ts")}
			for _, declaration := range placement[file] {
				site, known := sites[declaration]
				if !known {
					t.Fatal("no declaration", declaration)
				}
				unit.Declarations = append(unit.Declarations, site)
			}
			organization.Units = append(organization.Units, unit)
		}
		setting, err := organization.Setting()
		if err != nil {
			t.Fatal(err)
		}
		in.Configuration = []semantic.Setting{setting}
	}
	out, err := formatter.Format(ctx, in)
	if err != nil {
		return nil, err
	}
	published := map[files.ID]string{}
	for _, unit := range out.Units() {
		text, err := stages.Printer.Print(ctx, unit.Syntax)
		if err != nil {
			t.Fatal(err)
		}
		published[unit.File] = text
	}
	return published, nil
}

// Two imports of one export name in one file get two spellings: an alias that
// was chosen takes its name. Here a record T uses three declarations T of
// three other files.
func TestTypeScriptAliasesOfOneExportNameDiffer(t *testing.T) {
	published, err := namesExisting(t, map[string]string{
		"a":    "export type T = string;",
		"b":    "export type T = boolean;",
		"c":    "export type T = string[];",
		"user": "import type { T as A } from \"./a.js\"; import type { T as B } from \"./b.js\"; import type { T as C } from \"./c.js\"; export interface T { a: A; b: B; c: C }",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	const want = "import type { T as T_1 } from \"./a.js\";\nimport type { T as T_2 } from \"./b.js\";\nimport type { T as T_3 } from \"./c.js\";\n\nexport interface T {\n  a: T_1;\n  b: T_2;\n  c: T_3;\n}\n"
	if published["user"] != want {
		t.Fatalf("the file user: %q", published["user"])
	}
}

// A name that another use of the file is written with is taken. The use of
// the global Array is written Array, so the import of a declaration Array
// gets an alias, although the file declares no Array and has no parameter of
// that name.
func TestTypeScriptAliasAvoidsANameThatAnotherUseIsWrittenWith(t *testing.T) {
	published, err := namesExisting(t, map[string]string{
		"a": "export type Array = string;",
		"b": "import type { Array as Mine } from \"./a.js\"; export interface Box { items: Array<string>; mine: Mine }",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	const want = "import type { Array as Array_1 } from \"./a.js\";\n\nexport interface Box {\n  items: Array<string>;\n  mine: Array_1;\n}\n"
	if published["b"] != want {
		t.Fatalf("the file b: %q", published["b"])
	}
}

// The imports of a file are decided in the order of their module specifier
// and then of their export name. A file imports T and an export that is
// literally named T_1, beside a type parameter T. With the module of T first,
// T becomes T_1 and the export T_1 has to give way. With the module of T_1
// first, T_1 keeps its name and T takes the next suffix. An order by export
// name would give the first answer in both.
func TestTypeScriptAliasesAreDecidedByModuleSpecifierThenExportName(t *testing.T) {
	texts := map[string]string{
		"a":    "export type T = string;",
		"b":    "export type T_1 = boolean;",
		"user": "import type { T as A } from \"./a.js\"; import type { T_1 as B } from \"./b.js\"; export interface Box<T> { a: A; b: B; own: T }",
	}
	published, err := namesExisting(t, texts, nil)
	if err != nil {
		t.Fatal(err)
	}
	const first = "import type { T as T_1 } from \"./a.js\";\nimport type { T_1 as T_1_1 } from \"./b.js\";\n\nexport interface Box<T> {\n  a: T_1;\n  b: T_1_1;\n  own: T;\n}\n"
	if published["user"] != first {
		t.Fatalf("the module of T first: %q", published["user"])
	}
	published, err = namesExisting(t, texts, map[string][]string{"0": {"b.T_1"}, "a": {"a.T"}, "user": {"user.Box"}})
	if err != nil {
		t.Fatal(err)
	}
	const second = "import type { T_1 } from \"./0.js\";\nimport type { T as T_2 } from \"./a.js\";\n\nexport interface Box<T> {\n  a: T_2;\n  b: T_1;\n  own: T;\n}\n"
	if published["user"] != second {
		t.Fatalf("the module of T_1 first: %q", published["user"])
	}
}

// The capture is decided by the type parameters of the declaration that holds
// the use, and not by those of another declaration of the file. Box has a
// parameter T and uses only that. User uses the declaration T. In one file
// nothing is captured, and the placement is published.
func TestTypeScriptCaptureReadsTheDeclarationThatHoldsTheUse(t *testing.T) {
	published, err := namesExisting(t, map[string]string{
		"a": "export type T = string;",
		"b": "import type { T as U } from \"./a.js\"; export interface Box<T> { own: T } export interface User { u: U }",
	}, map[string][]string{"all": {"a.T", "b.Box", "b.User"}})
	if err != nil {
		t.Fatal("a use outside the declaration with the parameter was refused: ", err)
	}
	const want = "export type T = string;\n\nexport interface Box<T> {\n  own: T;\n}\n\nexport interface User {\n  u: T;\n}\n"
	if published["all"] != want {
		t.Fatalf("the file all: %q", published["all"])
	}
}

// A type parameter of an alias captures as one of an interface does.
func TestTypeScriptParameterOfAnAliasCaptures(t *testing.T) {
	texts := map[string]string{
		"a": "export type T = string;",
		"b": "import type { T as U } from \"./a.js\"; export type Pair<T> = [U, T];",
	}
	published, err := namesExisting(t, texts, nil)
	if err != nil {
		t.Fatal(err)
	}
	if published["b"] != "import type { T as T_1 } from \"./a.js\";\n\nexport type Pair<T> = [T_1, T];\n" {
		t.Fatalf("the file b: %q", published["b"])
	}
	published, err = namesExisting(t, texts, map[string][]string{"all": {"a.T", "b.Pair"}})
	if published != nil {
		t.Fatal("a use that the parameter of an alias captures was published")
	}
	onlyCategory(t, "a parameter of an alias that captures a declaration of its file", err, source.ErrUnsupported)
}

// A type parameter that is named like its own declaration generates for
// TypeScript and is published: inside the declaration the name is the
// parameter, and nothing uses the declaration there.
func TestTypeScriptParameterNamedLikeItsOwnDeclaration(t *testing.T) {
	output := namesOutput(t, "module collision; type box<box> = list<box>;")
	out, texts, err := namesFormat(t, output, tsOneFile(output, "types", "types.ts"))
	if err != nil || !out.Valid() {
		t.Fatal("a parameter named like its own declaration was refused: ", err)
	}
	if texts["types"] != "export type Box<Box> = Box[];\n" {
		t.Fatalf("the file: %q", texts["types"])
	}
}

// A use of a global type is written by the global's name. A declaration of
// its file with that name, or a type parameter of the declaration that holds
// the use, would capture it, and TypeScript refuses the placement as
// unsupported coverage. With the declaration in another file the placement is
// published, and where the file also uses the declaration its import gets an
// alias. Where the parameter and the use are in one declaration no placement
// helps. No builtin mapping writes a global type: the list here is a
// consumer's contribution that maps to the global Array.
func TestTypeScriptGlobalTypeCapturedByANameOfItsFile(t *testing.T) {
	arrays := selectedTS(listContribution("names/array/1", "", "Array"), "names/array/1")
	output := namesOutputWith(t, "module collision; type array = text; type box = list<array>;", arrays)
	tree := namesTree(t, output)
	inventory := output.Inventory()
	if len(inventory) != 2 || nodeAttribute(tree.Nodes[inventory[0].Site.Node-1], ts.AttributeName) != "Array" || len(output.ExternalUses()) != 1 {
		t.Fatal("the output owns a declaration Array and a declaration that applies the global Array")
	}
	out, _, err := namesFormat(t, output, tsOneFile(output, "types", "types.ts"))
	if out.Valid() {
		t.Fatal("a use of the global Array beside a declaration Array was published")
	}
	onlyCategory(t, "a declaration named like a global type that its file uses", err, source.ErrUnsupported)
	apart := tsfiles.GeneratedPlacement{
		Files:       map[files.ID]files.Path{"a": "a.ts", "b": "b.ts"},
		Definitions: map[g.DeclarationHandle]files.ID{inventory[0].Handle: "a", inventory[1].Handle: "b"},
	}
	out, texts, err := namesFormat(t, output, apart)
	if err != nil || !out.Valid() {
		t.Fatal("the two files were refused: ", err)
	}
	if texts["a"] != "export type Array = string;\n" || texts["b"] != "import type { Array as Array_1 } from \"./a.js\";\n\nexport type Box = Array<Array_1>;\n" {
		t.Fatalf("the two files: %q", texts)
	}
	output = namesOutputWith(t, "module collision; type box<array> = list<array>;", arrays)
	out, _, err = namesFormat(t, output, tsOneFile(output, "types", "types.ts"))
	if out.Valid() {
		t.Fatal("a use of the global Array under a parameter Array was published")
	}
	onlyCategory(t, "a type parameter named like a global type that its declaration uses", err, source.ErrUnsupported)
}
