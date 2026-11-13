package integration_test

import (
	"context"
	"errors"
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
