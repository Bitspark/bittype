package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/files/horos/provider/go"
	checkhoros "github.com/Bitspark/bittype/check/syntax/horos/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/files/go"
	filesformat "github.com/Bitspark/bittype/format/files/horos/provider/go"
	tsfiles "github.com/Bitspark/bittype/format/files/typescript/provider/go"
	formathoros "github.com/Bitspark/bittype/format/syntax/horos/provider/go"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	textfixture "github.com/Bitspark/bittype/integration/go/testdata/textfixture"
	"github.com/Bitspark/bittype/internal/horosfiles"
	filesparse "github.com/Bitspark/bittype/parse/files/horos/provider/go"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/semantics/horos/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	"github.com/Bitspark/bittype/source/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	tsboolean "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/boolean/go"
	tsdeclarations "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/declarations/go"
	tslist "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/list/go"
	tsrecord "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/record/go"
	tstext "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/text/go"
	compose "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	tstarget "github.com/Bitspark/bittype/transform/syntax/typescript/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// filesCategories are the four categories of a Files refusal. A refusal matches
// exactly one of them through errors.Is.
var filesCategories = []struct {
	name  string
	value error
}{
	{"invalid input", source.ErrInvalid},
	{"unsupported coverage", source.ErrUnsupported},
	{"missing evidence", semantic.ErrEvidenceUnavailable},
	{"exceeded limit", source.ErrLimit},
}

// filesRefusal is one refusal of one formatter. The same defect has one row per
// formatter that can receive it, and every row of a defect names one category.
// A row without a category is a refusal that a formatter passes on from a
// collaborator: it matches none of the four.
type filesRefusal struct {
	defect    string
	formatter string
	want      error
	// also lists values of other packages that the refusal still matches.
	also  []error
	check func(*testing.T, error)
	run   func(*testing.T) (bool, error)
}

// tamperedTarget builds the TypeScript assembly and then changes it, so that a
// sealed output carries what the TypeScript target never produces.
type tamperedTarget struct {
	tstarget.Provider
	change func(g.BuildInput[ts.Language], *g.Assembly) error
}

func (p tamperedTarget) Build(ctx context.Context, in g.BuildInput[ts.Language]) (g.Assembly, error) {
	assembly, err := p.Provider.Build(ctx, in)
	if err != nil {
		return g.Assembly{}, err
	}
	if err = p.change(in, &assembly); err != nil {
		return g.Assembly{}, err
	}
	return assembly, nil
}

// tamperedCellPage generates Cells = Page<Cell<Person>> with one external use,
// the Cell of a package, through a target that changes the assembly. A name in
// only requests that one declaration, so that the others are not exposed.
func tamperedCellPage(t *testing.T, change func(g.BuildInput[ts.Language], *g.Assembly) error, only string) g.Output {
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
	return generatedOutput(t, e, input, typescriptText(), func(registry *compose.Registry, input transform.Input, to representation.Language) (compose.Request, error) {
		request, err := genfixture.Request(registry, input, to, "tampered-ts", tamperedTarget{change: change}, tstext.Library(), tsboolean.Library(), tsrecord.Library(), tsdeclarations.Library(), tslist.Library())
		if err != nil {
			return compose.Request{}, err
		}
		if _, err = compose.RegisterMapping(registry, input.Source.Language(), to, cellMapping(from, pin.String(), &calls)); err != nil {
			return compose.Request{}, err
		}
		if only == "" {
			return request, nil
		}
		envelope, err := input.Source.EntryTree()
		if err != nil {
			return compose.Request{}, err
		}
		tree, err := representation.Unpack[s.Language](input.Source.Language(), envelope)
		if err != nil {
			return compose.Request{}, err
		}
		for _, site := range request.Roots {
			if nodeAttribute(tree.Nodes[site.Node-1], s.AttributeName) == only {
				request.Roots = []ast.NodeRef{site}
				return request, nil
			}
		}
		t.Fatal("no requested declaration", only)
		return request, nil
	})
}

// tamperedTree changes the one tree of an assembly.
func tamperedTree(change func(*ts.AST)) func(g.BuildInput[ts.Language], *g.Assembly) error {
	return func(in g.BuildInput[ts.Language], assembly *g.Assembly) error {
		tree, err := representation.Unpack[ts.Language](in.Language, assembly.Syntax.Trees[0])
		if err != nil {
			return err
		}
		change(&tree)
		assembly.Syntax.Trees[0], err = representation.Pack(in.Language, tree)
		return err
	}
}

// substitutedSyntax is a Syntax formatter for the horos Files formatter that
// parses a text of its own in place of a unit, so that the fresh output differs
// from what the horos Syntax formatter writes.
type substitutedSyntax struct {
	formathoros.Provider
	text func(ast.AST[s.Language]) (string, bool)
}

func (p substitutedSyntax) Format(ctx context.Context, view ast.View[s.Language]) (cst.CST[s.Language], error) {
	tree := view.ASTView()
	if text, ok := p.text(tree); ok {
		return parsehoros.Provider{}.Parse(ctx, parse.Input{SourceID: ast.SourceID(tree.ID), TreeID: tree.ID, Text: text})
	}
	return p.Provider.Format(ctx, view)
}

// errFreshTree is what freshRefusal answers. It is a value of no package that a
// Files category is joined to.
var errFreshTree = errors.New("fixture checker refuses a fresh tree")

// freshRefusal is a Syntax checker that refuses the trees a Files formatter
// constructs, which it knows by the prefix of their identities. The capture of
// the fresh output returns its error.
type freshRefusal[L any] struct {
	check.Provider[L]
	prefix string
}

func (p freshRefusal[L]) Check(view ast.View[L]) error {
	if strings.HasPrefix(string(view.ASTView().ID), p.prefix) {
		return errFreshTree
	}
	return p.Provider.Check(view)
}

// countingSyntax is a TypeScript Syntax formatter that counts its calls.
type countingSyntax struct {
	format.Provider[ast.View[ts.Language], cst.CST[ts.Language]]
	calls *int
}

func (p countingSyntax) Format(ctx context.Context, view ast.View[ts.Language]) (cst.CST[ts.Language], error) {
	*p.calls++
	return p.Provider.Format(ctx, view)
}

// commentingPrinter is a TypeScript Syntax printer that writes a long comment
// after each unit, so that the fresh text exceeds a byte limit that the input met.
type commentingPrinter struct {
	print.Provider[cst.CST[ts.Language], string]
}

func (p commentingPrinter) Print(ctx context.Context, unit cst.CST[ts.Language]) (string, error) {
	text, err := p.Provider.Print(ctx, unit)
	return text + "\n// " + strings.Repeat("x", 1<<16) + "\n", err
}

// errorTreeHas reports whether an error or one that it wraps or joins matches.
func errorTreeHas(err error, match func(error) bool) bool {
	if err == nil {
		return false
	}
	if match(err) {
		return true
	}
	switch wrapped := err.(type) {
	case interface{ Unwrap() error }:
		return errorTreeHas(wrapped.Unwrap(), match)
	case interface{ Unwrap() []error }:
		for _, each := range wrapped.Unwrap() {
			if errorTreeHas(each, match) {
				return true
			}
		}
	}
	return false
}

type horosPlacement struct {
	ID      string
	Path    string
	Modules []string
}

// TestFilesFailureCategories covers every refusal that the TypeScript Files
// formatters and the horos Files formatter produce themselves. Each refusal
// matches exactly one of the four categories, both formatters answer the same
// defect with the same category, and a value of another package still matches.
func TestFilesFailureCategories(t *testing.T) {
	ctx := context.Background()
	wrongKey := semantics.Key{Owner: "fixture", Name: "another-setting"}

	// TypeScript, generated declarations. The placement is a typed value built
	// from the handles of the output.
	output := generatedCellPage(t)
	tree := outputTree[ts.Language](t, output)
	cell := n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations}
	generated, err := tsfiles.NewGenerated(placementStages(t, cell))
	if err != nil {
		t.Fatal(err)
	}
	handles := tsDeclarations(t, output)
	// A second run of the same generation. Its handles belong to another output.
	otherOutput := generatedCellPage(t)
	otherHandles := tsDeclarations(t, otherOutput)
	const pagesPath, modelPath = "generated/pages/page.ts", "generated/model/person.ts"
	tsPlacement := func() tsfiles.GeneratedPlacement {
		return tsPlaced(t, handles, map[files.ID]files.Path{"pages": pagesPath, "model": modelPath}, map[string]files.ID{"Page": "pages", "Cells": "pages", "Cell": "pages", "Person": "model"})
	}
	tsFormat := func(f format.GeneratedFilesFormatter[ts.Language], out g.Output, limits format.FilesLimits, placement format.GeneratedPlacement) (bool, error) {
		result, err := f.Format(ctx, format.GeneratedFilesInput{Output: out, Placement: placement, Limits: limits})
		return result.Valid(), err
	}
	tsChanged := func(change func(*tsfiles.GeneratedPlacement)) func(*testing.T) (bool, error) {
		return func(*testing.T) (bool, error) {
			p := tsPlacement()
			change(&p)
			return tsFormat(generated, output, format.FilesLimits{}, p.Issue())
		}
	}
	tsLimited := func(change func(*format.FilesLimits)) func(*testing.T) (bool, error) {
		return func(*testing.T) (bool, error) {
			limits := format.DefaultFilesLimits()
			change(&limits)
			return tsFormat(generated, output, limits, tsPlacement().Issue())
		}
	}
	if published, err := tsChanged(func(*tsfiles.GeneratedPlacement) {})(t); err != nil || !published {
		t.Fatal("the TypeScript placement that every case changes must format", err)
	}

	// TypeScript, existing source. Its placement travels as a setting.
	existingStages := placementStages(t)
	plain, err := tsfiles.New(existingStages)
	if err != nil {
		t.Fatal(err)
	}
	tsUnit := func(id, text string) source.Unit[ts.Language] {
		parsed, err := typescriptText().Parse(ctx, parse.Input{SourceID: ast.SourceID(id), TreeID: ast.TreeID(id), Text: text})
		if err != nil {
			t.Fatal(err)
		}
		return source.Unit[ts.Language]{File: files.ID(id), Path: files.Path(id + ".ts"), Source: ast.SourceID(id), Syntax: parsed}
	}
	tsSite := func(unit source.Unit[ts.Language], name string) ast.NodeRef {
		tree := unit.Syntax.Abstract
		for _, child := range tree.Nodes[tree.Root-1].Children {
			if nodeAttribute(tree.Nodes[child.Node-1], ts.AttributeName) == name {
				return ast.NodeRef{Tree: tree.ID, Node: child.Node}
			}
		}
		t.Fatal("no declaration", name)
		return ast.NodeRef{}
	}
	tsBytes := func(o tsfiles.Organization) []byte {
		setting, err := o.Setting()
		if err != nil {
			t.Fatal(err)
		}
		return setting.Value
	}
	tsExisting := func(organization *tsfiles.Organization, units ...source.Unit[ts.Language]) func(*testing.T) (bool, error) {
		return func(t *testing.T) (bool, error) {
			in, err := existingStages.Capture(ctx, units, units[0].File, semantic.Limits{})
			if err != nil {
				t.Fatal("the existing source must be admitted before it is placed", err)
			}
			if organization != nil {
				in.Configuration = []semantic.Setting{{Key: tsfiles.OrganizationSetting(), Value: tsBytes(*organization)}}
			}
			out, err := plain.Format(ctx, in)
			return out.Valid(), err
		}
	}
	// One captured module. The rows about the form of a setting place its three
	// declarations into two files.
	library := tsUnit("input", "export type Person = string; export type Page<T> = [T, T]; export type People = Page<Person>;")
	captured, err := existingStages.Capture(ctx, []source.Unit[ts.Language]{library}, library.File, semantic.Limits{})
	if err != nil {
		t.Fatal(err)
	}
	existingPlacement := func() tsfiles.Organization {
		return tsfiles.Organization{Units: []tsfiles.Placement{
			{File: "pages", Path: pagesPath, Declarations: []ast.NodeRef{tsSite(library, "Page"), tsSite(library, "People")}},
			{File: "model", Path: modelPath, Declarations: []ast.NodeRef{tsSite(library, "Person")}},
		}}
	}
	existingFormat := func(limits format.FilesLimits, settings ...semantic.Setting) (bool, error) {
		in := captured
		in.Configuration, in.Limits = settings, limits
		out, err := plain.Format(ctx, in)
		return out.Valid(), err
	}
	existingRaw := func(value string) func(*testing.T) (bool, error) {
		return func(*testing.T) (bool, error) {
			return existingFormat(format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: []byte(value)})
		}
	}
	existingChanged := func(change func(*tsfiles.Organization)) func(*testing.T) (bool, error) {
		return func(*testing.T) (bool, error) {
			o := existingPlacement()
			change(&o)
			return existingFormat(format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: tsBytes(o)})
		}
	}
	existingLimited := func(change func(*format.FilesLimits)) func(*testing.T) (bool, error) {
		return func(*testing.T) (bool, error) {
			limits := format.DefaultFilesLimits()
			change(&limits)
			return existingFormat(limits, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: tsBytes(existingPlacement())})
		}
	}
	valid := string(tsBytes(existingPlacement()))
	if published, err := existingRaw(valid)(t); err != nil || !published {
		t.Fatal("the placement of existing TypeScript source that every case changes must format", err)
	}

	// horos.
	input, pagePin, peoplePin := filePagination(t)
	definition, err := representation.NewLanguage[s.Language](s.Language{}, s.Profile)
	if err != nil {
		t.Fatal(err)
	}
	checker, resolver, analyzer := fileStages(t, definition)
	syntaxPrinter := textfixture.Printing[s.Language](s.Language{})
	parser, err := filesparse.New(definition, parsehoros.Provider{}, syntaxPrinter, checker, resolver, analyzer, semantic.Options{}, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	horosFormatter, err := filesformat.New(definition, parsehoros.Provider{}, formathoros.Provider{}, syntaxPrinter, checker, resolver, analyzer, semantic.Options{}, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parser.Parse(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	horosFormat := func(f format.FilesFormatter[s.Language], from source.Source[s.Language], limits format.FilesLimits, settings ...semantic.Setting) (bool, error) {
		out, err := f.Format(ctx, format.FilesInput[s.Language]{Abstract: from.SemanticsView(), Snapshot: from.Snapshot(), Evidence: from.Evidence(), Configuration: settings, Limits: limits})
		return out.Valid(), err
	}
	horosPlacements := func() []horosPlacement {
		return []horosPlacement{{"people", "generated/people.horos", []string{peoplePin}}, {"page", "generated/lib/page.horos", []string{pagePin}}}
	}
	horosBytes := func(placements []horosPlacement) []byte {
		value, err := json.Marshal(placements)
		if err != nil {
			t.Fatal(err)
		}
		return value
	}
	horosRaw := func(value string) func(*testing.T) (bool, error) {
		return func(*testing.T) (bool, error) {
			return horosFormat(horosFormatter, parsed, format.FilesLimits{}, semantic.Setting{Key: filesformat.OrganizationSetting(), Value: []byte(value)})
		}
	}
	horosPlaced := func(change func(*[]horosPlacement)) func(*testing.T) (bool, error) {
		return func(*testing.T) (bool, error) {
			placements := horosPlacements()
			change(&placements)
			return horosFormat(horosFormatter, parsed, format.FilesLimits{}, semantic.Setting{Key: filesformat.OrganizationSetting(), Value: horosBytes(placements)})
		}
	}
	horosLimited := func(change func(*format.FilesLimits)) func(*testing.T) (bool, error) {
		return func(*testing.T) (bool, error) {
			limits := format.DefaultFilesLimits()
			change(&limits)
			return horosFormat(horosFormatter, parsed, limits, semantic.Setting{Key: filesformat.OrganizationSetting(), Value: horosBytes(horosPlacements())})
		}
	}
	if published, err := horosPlaced(func(*[]horosPlacement) {})(t); err != nil || !published {
		t.Fatal("the horos placement that every case changes must format", err)
	}

	// The fresh output of each formatter is analyzed again. These formatters make
	// that analysis fail: a package that lacks the export a declaration uses, and
	// a unit whose import has no supplied target.
	unavailable, err := tsfiles.NewGenerated(placementStages(t, n.PackageTypes{Module: "@bitspark/cell", Declaration: "export interface Other {}"}))
	if err != nil {
		t.Fatal(err)
	}
	detached, err := filesformat.New(definition, parsehoros.Provider{}, substitutedSyntax{text: func(tree ast.AST[s.Language]) (string, bool) {
		modules := tree.Nodes[tree.Root-1].Children
		return "module example.unrelated;\ntype Unrelated = text;\n", len(modules) == 1 && nodeAttribute(tree.Nodes[modules[0].Node-1], s.AttributePin) == pagePin
	}}, syntaxPrinter, checker, resolver, analyzer, semantic.Options{}, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	// A sealed output that the TypeScript target never produces, with every
	// declaration placed in one file.
	tampered := func(change func(g.BuildInput[ts.Language], *g.Assembly) error) func(*testing.T) (bool, error) {
		return func(t *testing.T) (bool, error) {
			out := tamperedCellPage(t, change, "")
			return tsFormat(generated, out, format.FilesLimits{}, tsOneFile(out, "types", "types.ts").Issue())
		}
	}
	if published, err := tampered(func(g.BuildInput[ts.Language], *g.Assembly) error { return nil })(t); err != nil || !published {
		t.Fatal("the generated output that every case changes must format", err)
	}
	otherIssuer := g.NewSymbolIssuer()

	// The smallest limit under which the shared contract admits an input. A
	// formatter that receives it refuses its own larger output.
	smallest := func(snapshot semantic.Snapshot, set func(*semantic.Limits, int)) int {
		low, high := 1, 1<<20
		for low < high {
			limits, middle := semantic.DefaultLimits(), (low+high)/2
			set(&limits, middle)
			if snapshot.CheckLimits(limits) == nil {
				high = middle
			} else {
				low = middle + 1
			}
		}
		return low
	}
	inputNodes := func(l *semantic.Limits, limit int) { l.MaxInputNodes = limit }
	inputBytes := func(l *semantic.Limits, limit int) { l.MaxInputBytes = limit }

	// One horos file that holds both modules. Placing each module in its own file
	// adds a program node, so the fresh output exceeds a node limit that the
	// input met.
	bundle, err := parser.Parse(ctx, parse.FilesInput{Entry: "bundle", Files: files.Files{Entries: []files.File{{ID: "bundle", Path: "bundle.horos", Text: input.Files.Entries[1].Text + input.Files.Entries[0].Text}}}})
	if err != nil {
		t.Fatal(err)
	}
	// A Syntax formatter that writes a long comment before the same modules. The
	// fresh trees equal what the input met, and the fresh text exceeds the byte
	// limit.
	commented, err := filesformat.New(definition, parsehoros.Provider{}, substitutedSyntax{text: func(ast.AST[s.Language]) (string, bool) {
		return "// " + strings.Repeat("x", 1<<16) + "\n" + input.Files.Entries[1].Text + input.Files.Entries[0].Text, true
	}}, syntaxPrinter, checker, resolver, analyzer, semantic.Options{}, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}

	hidden := tsUnit("input", "type Hidden = string; export type Visible = Hidden;")
	unused := tsUnit("input", "type Unused = string; export type Visible = string;")
	first, second := tsUnit("first", "export type Same = string;"), tsUnit("second", "export type Same = number;")

	const typescript, typescriptExisting, horosName = "typescript generated", "typescript existing", "horos"
	var cases []filesRefusal
	add := func(defect, formatter string, want error, run func(*testing.T) (bool, error), also ...error) {
		cases = append(cases, filesRefusal{defect: defect, formatter: formatter, want: want, also: also, run: run})
	}
	horosValid := string(horosBytes(horosPlacements()))
	// A refusal that the horos decoder decides is invalid input and still carries
	// the decoder's own error, whatever its type.
	horosDecoded := func(defect, value string) {
		add(defect, horosName, source.ErrInvalid, horosRaw(value))
		var placements []horosfiles.Placement
		decoder := json.NewDecoder(strings.NewReader(value))
		decoder.DisallowUnknownFields()
		want := decoder.Decode(&placements)
		if want == nil {
			t.Fatal("the decoder accepts:", defect)
		}
		cases[len(cases)-1].check = func(t *testing.T, err error) {
			if !errorTreeHas(err, func(each error) bool {
				return reflect.TypeOf(each) == reflect.TypeOf(want) && each.Error() == want.Error()
			}) {
				t.Errorf("the decoder's error %q (%T) is no longer part of the refusal: %q", want, want, err)
			}
		}
	}

	// The setting. A generated placement is a typed value and no setting, so
	// these are refusals of the two formatters that read a placement from bytes.
	add("wrong setting key", typescriptExisting, source.ErrInvalid, func(*testing.T) (bool, error) {
		return existingFormat(format.FilesLimits{}, semantic.Setting{Key: wrongKey, Value: []byte(valid)})
	})
	add("wrong setting key", horosName, source.ErrInvalid, func(*testing.T) (bool, error) {
		return horosFormat(horosFormatter, parsed, format.FilesLimits{}, semantic.Setting{Key: wrongKey, Value: []byte(horosValid)})
	}, semantic.ErrInvalidRequest)
	add("two settings", typescriptExisting, source.ErrInvalid, func(*testing.T) (bool, error) {
		return existingFormat(format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: []byte(valid)}, semantic.Setting{Key: wrongKey, Value: []byte(valid)})
	})
	add("two settings", horosName, source.ErrInvalid, func(*testing.T) (bool, error) {
		return horosFormat(horosFormatter, parsed, format.FilesLimits{}, semantic.Setting{Key: filesformat.OrganizationSetting(), Value: []byte(horosValid)}, semantic.Setting{Key: wrongKey, Value: []byte(horosValid)})
	}, semantic.ErrInvalidRequest)
	add("malformed JSON", typescriptExisting, source.ErrInvalid, existingRaw(`{"Units":[{"File":}]}`))
	horosDecoded("malformed JSON", `[{"ID":}]`)
	add("truncated JSON", typescriptExisting, source.ErrInvalid, existingRaw(valid[:len(valid)-1]))
	horosDecoded("truncated JSON", horosValid[:len(horosValid)-1])
	add("a JSON value of another type", typescriptExisting, source.ErrInvalid, existingRaw(`[]`))
	horosDecoded("a JSON value of another type", `{}`)
	add("an empty setting value", typescriptExisting, source.ErrInvalid, existingRaw(``))
	horosDecoded("an empty setting value", ``)
	add("the JSON value null", typescriptExisting, source.ErrInvalid, existingRaw(`null`), files.ErrInvalid)
	add("the JSON value null", horosName, source.ErrInvalid, horosRaw(`null`), files.ErrInvalid)
	add("a second JSON value", typescriptExisting, source.ErrInvalid, existingRaw(valid+` {}`))
	add("a second JSON value", horosName, source.ErrInvalid, horosRaw(horosValid+` []`), semantic.ErrInvalidRequest)

	// The request of a generated placement. The value is issued by the language
	// of the formatter, and the engine and the shared contract carry it unread.
	add("no placement", typescript, source.ErrInvalid, func(*testing.T) (bool, error) {
		return tsFormat(generated, output, format.FilesLimits{}, format.GeneratedPlacement{})
	})
	add("a placement that another language issued", typescript, source.ErrInvalid, func(*testing.T) (bool, error) {
		return tsFormat(generated, output, format.FilesLimits{}, format.NewGeneratedPlacementIssuer[tsfiles.GeneratedPlacement]().Issue(tsPlacement()))
	})
	// The step refuses both before anything is admitted: the refusal is its own
	// and carries no value of the common admission.
	for i := len(cases) - 2; i < len(cases); i++ {
		cases[i].check = func(t *testing.T, err error) {
			if errors.Is(err, files.ErrInvalid) {
				t.Errorf("the refusal comes from the common admission: %q", err)
			}
		}
	}

	// No field is accepted and ignored.
	add("unknown field", typescriptExisting, source.ErrInvalid, existingRaw(`{"Later":1,`+valid[1:]))
	add("unknown field of a file", typescriptExisting, source.ErrInvalid, existingRaw(strings.Replace(valid, `{"File":`, `{"Later":1,"File":`, 1)))
	add("unknown field of an item", typescriptExisting, source.ErrInvalid, existingRaw(strings.Replace(valid, `{"Tree":`, `{"Later":1,"Tree":`, 1)))
	// The placement of existing source named an entry until #273. No TypeScript
	// reader has one, so the vocabulary without it is published under a new key.
	// A setting under the earlier key is another setting, whatever it holds: with
	// the entry that it required, without one, and with the defects that the rows
	// "missing entry" held for it. Under the new key the member is an unknown
	// field.
	earlierKey := n.Key("declaration-placement-1")
	earlier := func(value string) func(*testing.T) (bool, error) {
		return func(*testing.T) (bool, error) {
			return existingFormat(format.FilesLimits{}, semantic.Setting{Key: earlierKey, Value: []byte(value)})
		}
	}
	withEntry := func(entry string, change func(*tsfiles.Organization)) string {
		o := existingPlacement()
		change(&o)
		return `{"Entry":"` + entry + `",` + string(tsBytes(o))[1:]
	}
	add("the earlier setting key with its entry", typescriptExisting, source.ErrInvalid, earlier(withEntry("pages", func(*tsfiles.Organization) {})))
	add("the earlier setting key without an entry", typescriptExisting, source.ErrInvalid, earlier(valid))
	add("missing entry under the earlier setting key", typescriptExisting, source.ErrInvalid, earlier(withEntry("absent", func(*tsfiles.Organization) {})))
	add("an excluded path of the first file and a missing entry under the earlier setting key", typescriptExisting, source.ErrInvalid, earlier(withEntry("absent", func(o *tsfiles.Organization) { o.Units[0].Path = "generated/elsewhere/file.mts" })))
	add("an excluded path of the second file and a missing entry under the earlier setting key", typescriptExisting, source.ErrInvalid, earlier(withEntry("absent", func(o *tsfiles.Organization) { o.Units[1].Path = "generated/elsewhere/file.mts" })))
	add("the entry member of the earlier vocabulary", typescriptExisting, source.ErrInvalid, existingRaw(`{"Entry":"pages",`+valid[1:]))
	horosDecoded("unknown field", `[{"Later":1,`+horosValid[2:])

	// A member that is named twice is refused. The decoder would keep the last
	// value, so the first would be accepted and ignored. A member name in another
	// case names the same member, and one such name alone stays admitted.
	for name, run := range map[string]func(*testing.T) (bool, error){
		"TypeScript, a member name in another case": existingRaw(strings.Replace(strings.Replace(valid, `"File"`, `"file"`, 1), `"Units"`, `"UNITS"`, 1)),
		"TypeScript, trailing white space":          existingRaw(valid + " \n\t"),
		"horos, a member name in another case":      horosRaw(strings.Replace(horosValid, `"ID"`, `"id"`, 1)),
		"horos, trailing white space":               horosRaw(horosValid + " \n"),
	} {
		if published, err := run(t); err != nil || !published {
			t.Fatal("an admitted placement is refused:", name, err)
		}
	}
	add("a repeated member", typescriptExisting, source.ErrInvalid, existingRaw(`{"Units":[],`+valid[1:]))
	add("a repeated member whose first value is another placement", typescriptExisting, source.ErrInvalid, existingRaw(`{"Units":[{"File":"x","Path":"../x","Declarations":[]}],`+valid[1:]))
	add("a repeated member of a file", typescriptExisting, source.ErrInvalid, existingRaw(strings.Replace(valid, `{"File":`, `{"File":"other","File":`, 1)))
	add("a repeated member of an item", typescriptExisting, source.ErrInvalid, existingRaw(strings.Replace(valid, `{"Tree":`, `{"Tree":"other","Tree":`, 1)))
	add("a member repeated in another case", typescriptExisting, source.ErrInvalid, existingRaw(`{"units":[],`+valid[1:]))
	add("a repeated member", horosName, source.ErrInvalid, horosRaw(`[{"ID":"../other",`+horosValid[2:]))
	add("a member repeated in another case", horosName, source.ErrInvalid, horosRaw(`[{"id":"../other",`+horosValid[2:]))

	// The inventory of files.
	add("too many files", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Source.Files.MaxFiles = 1 }), files.ErrLimit)
	add("too many files", typescriptExisting, source.ErrLimit, existingLimited(func(l *format.FilesLimits) { l.Source.Files.MaxFiles = 1 }))
	add("too many files", horosName, source.ErrLimit, horosLimited(func(l *format.FilesLimits) { l.Source.Files.MaxFiles = 1 }))
	// The number of files is counted first: before a handle is read, and before
	// the lists of a setting are.
	add("a handle of another output with too many files", typescript, source.ErrLimit, func(*testing.T) (bool, error) {
		p := tsPlacement()
		p.Definitions[otherHandles["Person"]] = "model"
		limits := format.DefaultFilesLimits()
		limits.Source.Files.MaxFiles = 1
		return tsFormat(generated, output, limits, p.Issue())
	}, files.ErrLimit)
	add("too many files and one file identity for two files", typescriptExisting, source.ErrLimit, func(*testing.T) (bool, error) {
		o := existingPlacement()
		o.Units[1].File = o.Units[0].File
		limits := format.DefaultFilesLimits()
		limits.Source.Files.MaxFiles = 1
		return existingFormat(limits, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: tsBytes(o)})
	})
	// Each file of a placement is published with one syntax link. The shared
	// admission applies the link limit to the number of files. The horos
	// formatter is not on that admission and does not apply the limit.
	add("more files than the link limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Source.Files.MaxLinks = 1 }), files.ErrLimit)
	add("more files than the link limit", typescriptExisting, source.ErrLimit, existingLimited(func(l *format.FilesLimits) { l.Source.Files.MaxLinks = 1 }), files.ErrLimit)
	add("file identities and paths exceed the byte limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Source.Files.MaxBytes = 8 }), files.ErrLimit)
	add("file identities and paths exceed the byte limit", horosName, source.ErrLimit, horosLimited(func(l *format.FilesLimits) { l.Source.Files.MaxBytes = 8 }), files.ErrLimit)
	add("the fresh files exceed the Files byte limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) {
		// The identities and paths fit. The text of the fresh files does not.
		l.Source.Files.MaxBytes = 0
		for file, path := range tsPlacement().Files {
			l.Source.Files.MaxBytes += len(file) + len(path)
		}
	}), files.ErrLimit)
	add("the fresh files exceed the Files byte limit", typescriptExisting, source.ErrLimit, func(t *testing.T) (bool, error) {
		unit := tsUnit("a", "export type A = string;")
		in, err := existingStages.Capture(ctx, []source.Unit[ts.Language]{unit}, unit.File, semantic.Limits{})
		if err != nil {
			t.Fatal(err)
		}
		if out, err := plain.Format(ctx, in); err != nil || !out.Valid() {
			t.Fatal("the module must format under the default limits", err)
		}
		in.Limits = format.DefaultFilesLimits()
		in.Limits.Source.Files.MaxBytes = len(unit.File) + len(unit.Path) + 1
		out, err := plain.Format(ctx, in)
		return out.Valid(), err
	}, files.ErrLimit)
	add("the fresh files exceed the Files byte limit", horosName, source.ErrLimit, horosLimited(func(l *format.FilesLimits) {
		l.Source.Files.MaxBytes = 0
		for _, placement := range horosPlacements() {
			l.Source.Files.MaxBytes += len(placement.ID) + len(placement.Path)
		}
	}), files.ErrLimit)
	add("the fresh paths exceed the input byte limit", typescript, source.ErrLimit, func(t *testing.T) (bool, error) {
		// The long path belongs to the file that imports, so no import names it and
		// the fresh text stays below the limit.
		long := tsPlacement()
		long.Files["pages"] = files.Path("generated/" + strings.Repeat("d", 6000) + "/page.ts")
		if published, err := tsFormat(generated, output, format.FilesLimits{}, long.Issue()); err != nil || !published {
			t.Fatal("the long path must format under the default limits", err)
		}
		limits := format.DefaultFilesLimits()
		limits.Semantic.MaxInputBytes = 5000
		if published, err := tsFormat(generated, output, limits, tsPlacement().Issue()); err != nil || !published {
			t.Fatal("the short path must format under the limit", err)
		}
		return tsFormat(generated, output, limits, long.Issue())
	}, files.ErrLimit)
	add("the fresh paths exceed the input byte limit", horosName, source.ErrLimit, func(t *testing.T) (bool, error) {
		placements := horosPlacements()
		placements[0].Path = "generated/" + strings.Repeat("d", 6000) + "/people.horos"
		setting := semantic.Setting{Key: filesformat.OrganizationSetting(), Value: horosBytes(placements)}
		if published, err := horosFormat(horosFormatter, parsed, format.FilesLimits{}, setting); err != nil || !published {
			t.Fatal("the long path must format under the default limits", err)
		}
		limits := format.DefaultFilesLimits()
		limits.Semantic.MaxInputBytes = 5000
		if published, err := horosFormat(horosFormatter, parsed, limits, semantic.Setting{Key: filesformat.OrganizationSetting(), Value: []byte(horosValid)}); err != nil || !published {
			t.Fatal("the short path must format under the limit", err)
		}
		return horosFormat(horosFormatter, parsed, limits, setting)
	}, files.ErrLimit)
	add("two captured files with one path", typescriptExisting, source.ErrInvalid, func(*testing.T) (bool, error) {
		first, second := tsUnit("first", "export type A = string;"), tsUnit("second", "export type B = string;")
		second.Path = first.Path
		in, err := existingStages.Capture(ctx, []source.Unit[ts.Language]{first, second}, first.File, semantic.Limits{})
		return in.Abstract.Valid(), err
	}, files.ErrInvalid)
	// A keyed table cannot name a file identity twice. A list can.
	add("one file identity for two files", typescriptExisting, source.ErrInvalid, existingChanged(func(o *tsfiles.Organization) { o.Units[1].File = o.Units[0].File }), files.ErrInvalid)
	add("one file identity for two files", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].ID = (*p)[0].ID }), files.ErrInvalid)
	// The file model under the empty identity, with its declaration.
	emptyIdentity := func(p *tsfiles.GeneratedPlacement) {
		p.Files[""] = p.Files["model"]
		delete(p.Files, "model")
		p.Definitions[handles["Person"]] = ""
	}
	add("empty file identity", typescript, source.ErrInvalid, tsChanged(emptyIdentity), files.ErrInvalid)
	add("empty file identity", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].ID = "" }), files.ErrInvalid)
	add("one path for two files", typescript, source.ErrInvalid, tsChanged(func(p *tsfiles.GeneratedPlacement) { p.Files["model"] = p.Files["pages"] }), files.ErrInvalid)
	add("one path for two files", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].Path = (*p)[0].Path }), files.ErrInvalid)
	add("a path that is also a directory of another", typescript, source.ErrInvalid, tsChanged(func(p *tsfiles.GeneratedPlacement) { p.Files["model"] = p.Files["pages"] + "/person.ts" }), files.ErrInvalid)
	add("a path that is also a directory of another", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].Path = (*p)[0].Path + "/page.horos" }), files.ErrInvalid)
	// The shared admission refuses a path that is no path of the files contract
	// before TypeScript reads it, so the refusal carries the files value.
	add("the path ..", typescript, source.ErrInvalid, tsChanged(func(p *tsfiles.GeneratedPlacement) { p.Files["model"] = ".." }), files.ErrInvalid)
	add("the path ..", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].Path = ".." }), files.ErrInvalid)
	add("an absolute path", typescript, source.ErrInvalid, tsChanged(func(p *tsfiles.GeneratedPlacement) { p.Files["model"] = "/generated/person.ts" }), files.ErrInvalid)
	add("an absolute path", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].Path = "/generated/page.horos" }), files.ErrInvalid)
	// A path that is no TypeScript path is invalid input. A path that TypeScript
	// itself reads and that this profile excludes is unsupported coverage.
	tsPath := func(path string) func(*testing.T) (bool, error) {
		return tsChanged(func(p *tsfiles.GeneratedPlacement) { p.Files["model"] = files.Path(path) })
	}
	add("a valid path that is not a TypeScript path", typescript, source.ErrInvalid, tsPath("generated/model/person.txt"))
	add("a path with the extension in upper case", typescript, source.ErrInvalid, tsPath("generated/model/PERSON.TS"))
	add("a path that is no TypeScript path and has an excluded character", typescript, source.ErrInvalid, tsPath("generated/model/per#son.txt"))
	add("a declaration file path", typescript, source.ErrUnsupported, tsPath("generated/model/person.d.ts"))
	add("a path with the extension mts", typescript, source.ErrUnsupported, tsPath("generated/model/person.mts"))
	add("a path with the extension cts", typescript, source.ErrUnsupported, tsPath("generated/model/person.cts"))
	add("a path with the extension tsx", typescript, source.ErrUnsupported, tsPath("generated/model/person.tsx"))
	add("a path with a number sign", typescript, source.ErrUnsupported, tsPath("generated/model/per#son.ts"))
	add("a path with a percent sign", typescript, source.ErrUnsupported, tsPath("generated/model/per%son.ts"))
	add("a path with a question mark", typescript, source.ErrUnsupported, tsPath("generated/model/per?son.ts"))

	// The items of the placement.
	add("a file without a declaration", typescript, source.ErrUnsupported, tsChanged(func(p *tsfiles.GeneratedPlacement) { p.Files["empty"] = "generated/empty.ts" }))
	formatted := 0
	add("a file without a declaration, before any fresh file is formatted", typescript, source.ErrUnsupported, func(t *testing.T) (bool, error) {
		stages := placementStages(t, cell)
		stages.Syntax = countingSyntax{stages.Syntax, &formatted}
		counting, err := tsfiles.NewGenerated(stages)
		if err != nil {
			t.Fatal(err)
		}
		if published, err := tsFormat(counting, output, format.FilesLimits{}, tsPlacement().Issue()); err != nil || !published || formatted != 2 {
			t.Fatal("the base placement must format its two files", err, formatted)
		}
		formatted = 0
		p := tsPlacement()
		p.Files["empty"] = "generated/empty.ts"
		return tsFormat(counting, output, format.FilesLimits{}, p.Issue())
	})
	cases[len(cases)-1].check = func(t *testing.T, _ error) {
		if formatted != 0 {
			t.Errorf("the placement was refused after %d fresh files were formatted", formatted)
		}
	}
	add("a file without a declaration", typescriptExisting, source.ErrUnsupported, existingChanged(func(o *tsfiles.Organization) {
		o.Units = append(o.Units, tsfiles.Placement{File: "empty", Path: "generated/empty.ts"})
	}))
	add("a module of imports only in its original placement", typescriptExisting, source.ErrUnsupported, tsExisting(nil, tsUnit("a", "export type A = string;"), tsUnit("user", "import type { A } from \"./a.js\";")))
	add("a file without a module", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) {
		*p = append(*p, horosPlacement{ID: "empty", Path: "generated/empty.horos"})
	}))
	add("no file", typescript, source.ErrInvalid, tsChanged(func(p *tsfiles.GeneratedPlacement) {
		clear(p.Files)
		clear(p.Definitions)
	}), files.ErrInvalid)
	add("no file", typescriptExisting, source.ErrInvalid, existingChanged(func(o *tsfiles.Organization) { o.Units = nil }), files.ErrInvalid)
	add("no file", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { *p = nil }))
	// An item that the placed set does not hold. For a generated placement that
	// is the handle of another output, also of another run of the same plan, and
	// the refusal carries the value with which the output refuses the handle.
	add("unknown item", typescript, source.ErrInvalid, tsChanged(func(p *tsfiles.GeneratedPlacement) {
		p.Definitions[otherHandles["Person"]] = "model"
	}), g.ErrScope)
	add("a placement that was built for another run", typescript, source.ErrInvalid, func(*testing.T) (bool, error) {
		p := tsPlaced(t, otherHandles, map[files.ID]files.Path{"pages": pagesPath, "model": modelPath}, map[string]files.ID{"Page": "pages", "Cells": "pages", "Cell": "pages", "Person": "model"})
		if published, err := tsFormat(generated, otherOutput, format.FilesLimits{}, p.Issue()); err != nil || !published {
			t.Fatal("the placement must format with the output it was built for", err)
		}
		return tsFormat(generated, output, format.FilesLimits{}, p.Issue())
	}, g.ErrScope)
	add("unknown item", typescriptExisting, source.ErrInvalid, existingChanged(func(o *tsfiles.Organization) {
		o.Units[1].Declarations = append(o.Units[1].Declarations, ast.NodeRef{Tree: "foreign", Node: 1})
	}), files.ErrInvalid)
	add("unknown item", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].Modules = append((*p)[1].Modules, "0000") }))
	// A keyed table cannot place an item twice. A list can.
	add("repeated item", typescriptExisting, source.ErrInvalid, existingChanged(func(o *tsfiles.Organization) {
		o.Units[1].Declarations = append(o.Units[1].Declarations, tsSite(library, "Page"))
	}))
	add("repeated item", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].Modules = append((*p)[1].Modules, peoplePin) }))
	add("missing item", typescript, source.ErrInvalid, tsChanged(func(p *tsfiles.GeneratedPlacement) {
		delete(p.Files, "model")
		delete(p.Definitions, handles["Person"])
	}), files.ErrInvalid)
	add("missing item", typescriptExisting, source.ErrInvalid, existingChanged(func(o *tsfiles.Organization) { o.Units = o.Units[:1] }), files.ErrInvalid)
	add("missing item", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { *p = (*p)[1:] }))
	add("an item in a file that the placement does not declare", typescript, source.ErrInvalid, tsChanged(func(p *tsfiles.GeneratedPlacement) {
		delete(p.Files, "model")
	}), files.ErrInvalid)

	// What the placement asks of the declarations.
	// A placement that separates an unexposed declaration from its user is
	// invalid input: no later profile can admit it for the same output, because
	// Files never exports a declaration without authorization.
	add("a declaration uses an unexposed one of another file", typescriptExisting, source.ErrInvalid, tsExisting(&tsfiles.Organization{Units: []tsfiles.Placement{
		{File: "public", Path: "public.ts", Declarations: []ast.NodeRef{tsSite(hidden, "Visible")}},
		{File: "private", Path: "private.ts", Declarations: []ast.NodeRef{tsSite(hidden, "Hidden")}},
	}}, hidden))
	// separated generates Cells as the one exposed declaration and places it in
	// another file than the declarations that it uses.
	separated := func(t *testing.T, path files.Path) (bool, error) {
		out := tamperedCellPage(t, func(g.BuildInput[ts.Language], *g.Assembly) error { return nil }, "Cells")
		named := tsDeclarations(t, out)
		if exposed := out.Exposed(); len(exposed) != 1 || len(named) < 2 {
			t.Fatal("expected Cells as the one exposed declaration", exposed)
		} else if record, err := out.Resolve(named["Cells"]); err != nil || record.Site != exposed[0] {
			t.Fatal("expected Cells as the one exposed declaration", exposed, err)
		}
		if published, err := tsFormat(generated, out, format.FilesLimits{}, tsOneFile(out, "user", "user.ts").Issue()); err != nil || !published {
			t.Fatal("the same declarations in one file must format", err)
		}
		apart := tsOneFile(out, "used", "used.ts")
		apart.Files["user"] = path
		apart.Definitions[named["Cells"]] = "user"
		return tsFormat(generated, out, format.FilesLimits{}, apart.Issue())
	}
	add("a declaration uses an unexposed one of another file", typescript, source.ErrInvalid, func(t *testing.T) (bool, error) {
		return separated(t, "user.ts")
	})
	add("two declarations of one name in one file", typescriptExisting, source.ErrInvalid, tsExisting(&tsfiles.Organization{Units: []tsfiles.Placement{
		{File: "both", Path: "both.ts", Declarations: []ast.NodeRef{tsSite(first, "Same"), tsSite(second, "Same")}},
	}}, first, second))
	// Generation refuses two declarations of one name until the flat-name rule
	// leaves it, so the output is built by a target that renames one.
	renamed := tamperedTree(func(tree *ts.AST) {
		for _, child := range tree.Nodes[tree.Root-1].Children {
			node := &tree.Nodes[child.Node-1]
			for i := range node.Attributes {
				if node.Attributes[i].Name == ts.AttributeName && node.Attributes[i].Text == "Person" {
					node.Attributes[i].Text = "Page"
				}
			}
		}
	})
	add("two declarations of one name in one file", typescript, source.ErrInvalid, tampered(renamed))
	// A type parameter that would capture a use of a declaration of its file.
	// A later slice may rename the parameter where that is permitted, so the
	// placement is unsupported coverage. Another file for the declaration is
	// admitted today: the import gets an alias.
	capturing := func(t *testing.T, path files.Path) (bool, error) {
		out := namesOutput(t, capturedNames)
		return tsFormat(generated, out, format.FilesLimits{}, tsOneFile(out, "types", path).Issue())
	}
	add("a type parameter that captures a declaration of its file", typescript, source.ErrUnsupported, func(t *testing.T) (bool, error) {
		return capturing(t, "types.ts")
	})
	add("a type parameter that captures a declaration of its file, in a file whose path is no TypeScript path", typescript, source.ErrInvalid, func(t *testing.T) (bool, error) {
		return capturing(t, "types.txt")
	})
	capturedExisting := func(path string) func(*testing.T) (bool, error) {
		declared := tsUnit("a", "export type T = string;")
		user := tsUnit("user", "import type { T as U } from \"./a.js\"; export interface Box<T> { own: U; item: T }")
		return tsExisting(&tsfiles.Organization{Units: []tsfiles.Placement{
			{File: "both", Path: files.Path(path), Declarations: []ast.NodeRef{tsSite(declared, "T"), tsSite(user, "Box")}},
		}}, declared, user)
	}
	add("a type parameter that captures a declaration of its file", typescriptExisting, source.ErrUnsupported, capturedExisting("both.ts"))
	// Two captures whose text the compiler refuses. Until 13 November 2026 the
	// compiler's answer came first and was missing evidence (TS2315, TS2313).
	// TypeScript Files judges a capture itself, before a tree is built.
	capturedIn := func(declaration, user string) func(*testing.T) (bool, error) {
		declared := tsUnit("a", declaration)
		using := tsUnit("user", user)
		return tsExisting(&tsfiles.Organization{Units: []tsfiles.Placement{
			{File: "both", Path: "both.ts", Declarations: []ast.NodeRef{tsSite(declared, "T"), tsSite(using, "Box")}},
		}}, declared, using)
	}
	add("a captured use that is applied to arguments", typescriptExisting, source.ErrUnsupported, capturedIn("export type T<V> = V[];", "import type { T as U } from \"./a.js\"; export interface Box<T> { own: U<T> }"))
	add("a capture in the constraint of the parameter", typescriptExisting, source.ErrUnsupported, capturedIn("export type T = string;", "import type { T as U } from \"./a.js\"; export interface Box<T extends U> { v: T }"))
	// A use of a global type is written by the global's name. A declaration of
	// the file, or a type parameter of the declaration that holds the use, that
	// has that name would capture it. No builtin mapping writes a global type:
	// a consumer's contribution does, here a list as the global Array.
	arrays := selectedTS(listContribution("names/array/1", "", "Array"), "names/array/1")
	add("a declaration named Array and a use of the global Array in one file", typescript, source.ErrUnsupported, func(t *testing.T) (bool, error) {
		out := namesOutputWith(t, "module collision; type array = text; type box = list<text>;", arrays)
		return tsFormat(generated, out, format.FilesLimits{}, tsOneFile(out, "types", "types.ts").Issue())
	})
	add("a type parameter named Array in the declaration that applies the global Array, which no placement helps", typescript, source.ErrUnsupported, func(t *testing.T) (bool, error) {
		out := namesOutputWith(t, "module collision; type box<array> = list<array>;", arrays)
		return tsFormat(generated, out, format.FilesLimits{}, tsOneFile(out, "types", "types.ts").Issue())
	})
	add("a declaration named Array and a use of the global Array in one file", typescriptExisting, source.ErrUnsupported, func(t *testing.T) (bool, error) {
		declared := tsUnit("a", "export type Array = string;")
		using := tsUnit("user", "export interface Box { items: Array<string> }")
		return tsExisting(&tsfiles.Organization{Units: []tsfiles.Placement{
			{File: "both", Path: "both.ts", Declarations: []ast.NodeRef{tsSite(declared, "Array"), tsSite(using, "Box")}},
		}}, declared, using)(t)
	})
	// A member that is named like a type parameter of its interface. The
	// compiler keeps the member and the parameter in one symbol, so the use of
	// the parameter has two declarations in the graph. That is valid TypeScript
	// which this profile does not tell apart: unsupported coverage, and no
	// defect of the input. Until 13 November 2026 it was refused as invalid.
	add("a member named like a type parameter of its declaration", typescriptExisting, source.ErrUnsupported, tsExisting(nil, tsUnit("input", "export interface Box<T> { T: T }")))
	add("a member named like a type parameter of its declaration", typescript, source.ErrUnsupported, func(t *testing.T) (bool, error) {
		out := namesOutput(t, "module collision; type Box<T> = { T: T; };")
		return tsFormat(generated, out, format.FilesLimits{}, tsOneFile(out, "types", "types.ts").Issue())
	})
	// The reader refuses a member named like a type parameter while it reads the
	// graph, which is before the two refusals of relocation. This row pins what
	// is found first, as the row of a script does: the member, and not the use
	// of an unexposed declaration of another file. The category describes the
	// defect that was found.
	add("a member named like a type parameter of its declaration and a use of an unexposed declaration of another file", typescriptExisting, source.ErrUnsupported, func(t *testing.T) (bool, error) {
		member := tsUnit("member", "export interface Box<T> { T: T }")
		return tsExisting(&tsfiles.Organization{Units: []tsfiles.Placement{
			{File: "member", Path: "member.ts", Declarations: []ast.NodeRef{tsSite(member, "Box")}},
			{File: "public", Path: "public.ts", Declarations: []ast.NodeRef{tsSite(hidden, "Visible")}},
			{File: "private", Path: "private.ts", Declarations: []ast.NodeRef{tsSite(hidden, "Hidden")}},
		}}, member, hidden)(t)
	})
	// The capture of a global type is a refusal of coverage too, and it is
	// reported after everything invalid, wherever its file stands among the
	// files: before the file with the use of an unexposed declaration, or after.
	for _, order := range []string{"first", "last"} {
		add("a global type captured in the "+order+" file and a use of an unexposed declaration of another file", typescriptExisting, source.ErrInvalid, func(t *testing.T) (bool, error) {
			declared := tsUnit("a", "export type Array = string;")
			using := tsUnit("user", "export interface Box { items: Array<string> }")
			capture := tsfiles.Placement{File: "both", Path: "both.ts", Declarations: []ast.NodeRef{tsSite(declared, "Array"), tsSite(using, "Box")}}
			apart := []tsfiles.Placement{
				{File: "public", Path: "public.ts", Declarations: []ast.NodeRef{tsSite(hidden, "Visible")}},
				{File: "private", Path: "private.ts", Declarations: []ast.NodeRef{tsSite(hidden, "Hidden")}},
			}
			units := append([]tsfiles.Placement{capture}, apart...)
			if order == "last" {
				units = append(apart, capture)
			}
			return tsExisting(&tsfiles.Organization{Units: units}, declared, using, hidden)(t)
		})
	}
	// The capture is found in the same pass as the refusals of relocation and is
	// reported after them: a placement that also separates an unexposed
	// declaration from its user is invalid input.
	add("a type parameter that captures a declaration of its file and a use of an unexposed declaration of another file", typescriptExisting, source.ErrInvalid, func(t *testing.T) (bool, error) {
		declared := tsUnit("a", "export type T = string;")
		user := tsUnit("user", "import type { T as U } from \"./a.js\"; export interface Box<T> { own: U; item: T }")
		return tsExisting(&tsfiles.Organization{Units: []tsfiles.Placement{
			{File: "both", Path: "both.ts", Declarations: []ast.NodeRef{tsSite(declared, "T"), tsSite(user, "Box")}},
			{File: "public", Path: "public.ts", Declarations: []ast.NodeRef{tsSite(hidden, "Visible")}},
			{File: "private", Path: "private.ts", Declarations: []ast.NodeRef{tsSite(hidden, "Hidden")}},
		}}, declared, user, hidden)(t)
	})
	add("a type parameter that captures a declaration of its file, in a file whose path is no TypeScript path", typescriptExisting, source.ErrInvalid, capturedExisting("both.txt"))

	// What the sealed generated output carries.
	add("an external use of another issuer", typescript, source.ErrInvalid, tampered(func(in g.BuildInput[ts.Language], assembly *g.Assembly) error {
		symbol, err := otherIssuer.Issue(in.Language, "@bitspark/cell\x00Cell")
		assembly.ExternalUses[0].Symbol = symbol
		return err
	}), g.ErrScope)
	add("an external use in another role", typescript, source.ErrInvalid, tampered(func(_ g.BuildInput[ts.Language], assembly *g.Assembly) error {
		assembly.ExternalUses[0].At.Role = wrongKey
		return nil
	}), g.ErrScope)
	add("an external use at a site that does not name it", typescript, source.ErrInvalid, tampered(func(_ g.BuildInput[ts.Language], assembly *g.Assembly) error {
		assembly.ExternalUses[0].At.Site = assembly.Declarations[0]
		return nil
	}))
	add("a name use without a destination", typescript, source.ErrInvalid, tampered(func(_ g.BuildInput[ts.Language], assembly *g.Assembly) error {
		assembly.ExternalUses = nil
		return nil
	}))
	add("an exposed site that is no top-level declaration", typescript, source.ErrInvalid, tampered(func(in g.BuildInput[ts.Language], assembly *g.Assembly) error {
		// The declaration that holds the external use is listed at that use, a
		// node below it. No listed site contains another, so the output seals.
		tree, err := representation.Unpack[ts.Language](in.Language, assembly.Syntax.Trees[0])
		if err != nil {
			return err
		}
		use := assembly.ExternalUses[0].At.Site
		var holds func(ast.NodeID) bool
		holds = func(id ast.NodeID) bool {
			if id == use.Node {
				return true
			}
			for _, child := range tree.Nodes[id-1].Children {
				if holds(child.Node) {
					return true
				}
			}
			return false
		}
		for i, site := range assembly.Declarations {
			if holds(site.Node) {
				assembly.Declarations[i] = use
				return nil
			}
		}
		return errors.New("no declaration holds the external use")
	}))
	add("a top-level declaration that the output does not own", typescript, source.ErrInvalid, tampered(tamperedTree(func(tree *ts.AST) {
		// A declaration that no record of the inventory names.
		tree.Nodes = append(tree.Nodes, ts.Node{Kind: ts.KindKeywordType, Attributes: []ast.Attribute{{Name: ts.AttributeName, Text: "string"}}})
		body := ast.NodeID(len(tree.Nodes))
		tree.Nodes = append(tree.Nodes, ts.Node{Kind: ts.KindTypeAlias, Attributes: []ast.Attribute{{Name: ts.AttributeName, Text: "Unowned"}}, Children: []ast.Child{{Role: ts.RoleType, Node: body}}})
		root := &tree.Nodes[tree.Root-1]
		root.Children = append(root.Children, ast.Child{Role: ts.RoleDeclaration, Node: ast.NodeID(len(tree.Nodes))})
	})))
	add("an owned declaration that is not exposed and is no top-level declaration", typescript, source.ErrInvalid, func(t *testing.T) (bool, error) {
		// Only Cells is requested. Another declaration is listed at its first
		// child, a node below it, and an unowned declaration keeps the count of
		// the top-level declarations equal to the count of the inventory.
		out := tamperedCellPage(t, func(in g.BuildInput[ts.Language], assembly *g.Assembly) error {
			tree, err := representation.Unpack[ts.Language](in.Language, assembly.Syntax.Trees[0])
			if err != nil {
				return err
			}
			for i, site := range assembly.Declarations {
				node := tree.Nodes[site.Node-1]
				if nodeAttribute(node, ts.AttributeName) != "Cells" && len(node.Children) > 0 {
					assembly.Declarations[i] = ast.NodeRef{Tree: site.Tree, Node: node.Children[0].Node}
					return nil
				}
			}
			return errors.New("no declaration with a child beside Cells")
		}, "Cells")
		return tsFormat(generated, out, format.FilesLimits{}, tsOneFile(out, "types", "types.ts").Issue())
	})
	// The TypeScript provider is exported, and only the shared formatter builds
	// its request. It refuses a request that carries no admitted placement.
	add("a request that the shared formatter did not build", typescript, source.ErrInvalid, func(*testing.T) (bool, error) {
		out, err := (tsfiles.GeneratedProvider{}).Format(ctx, format.GeneratedFilesRequest[ts.Language]{})
		return out.Valid(), err
	})
	add("a generated declaration that is already exported", typescript, source.ErrInvalid, tampered(tamperedTree(func(tree *ts.AST) {
		node := &tree.Nodes[tree.Nodes[tree.Root-1].Children[0].Node-1]
		node.Attributes = append(node.Attributes, ast.Attribute{Name: ts.AttributeExported, Text: "true"})
	})))
	add("a generated tree with an import", typescript, source.ErrInvalid, tampered(tamperedTree(func(tree *ts.AST) {
		tree.Nodes = append(tree.Nodes, ts.Node{Kind: ts.KindImportBinding, Attributes: []ast.Attribute{{Name: ts.AttributeName, Text: "Other"}}})
		binding := ast.NodeID(len(tree.Nodes))
		tree.Nodes = append(tree.Nodes, ts.Node{Kind: ts.KindImport, Attributes: []ast.Attribute{{Name: ts.AttributeModule, Text: "@example/other"}}, Children: []ast.Child{{Role: ts.RoleBinding, Node: binding}}})
		root := &tree.Nodes[tree.Root-1]
		root.Children = append([]ast.Child{{Role: ts.RoleImport, Node: ast.NodeID(len(tree.Nodes))}}, root.Children...)
	})))

	// Several defects in one placement. Validity is checked before coverage: a
	// placement that no slice can admit is invalid input, whichever of its files
	// has the defect that a later slice could cover. For a generated placement
	// the shared admission decides the common part of validity, and it runs
	// before TypeScript refuses what it does not cover.
	const excluded = "generated/elsewhere/file.mts"
	for _, defect := range []struct {
		name   string
		change func(p *tsfiles.GeneratedPlacement, file files.ID)
		also   []error
	}{
		{"a declaration of another output", func(p *tsfiles.GeneratedPlacement, file files.ID) { p.Definitions[otherHandles["Person"]] = file }, []error{g.ErrScope}},
		{"the path ..", func(p *tsfiles.GeneratedPlacement, file files.ID) { p.Files[file] = ".." }, []error{files.ErrInvalid}},
		{"a path that is no TypeScript path", func(p *tsfiles.GeneratedPlacement, file files.ID) { p.Files[file] = "generated/model/person.txt" }, nil},
	} {
		for _, order := range [][2]files.ID{{"model", "pages"}, {"pages", "model"}} {
			add("an excluded path of the file "+string(order[0])+" and "+defect.name+" in the file "+string(order[1]), typescript, source.ErrInvalid, tsChanged(func(p *tsfiles.GeneratedPlacement) {
				p.Files[order[0]] = excluded
				defect.change(p, order[1])
			}), defect.also...)
		}
	}
	add("an excluded path and a declaration of another output in one file", typescript, source.ErrInvalid, tsChanged(func(p *tsfiles.GeneratedPlacement) {
		p.Files["pages"] = excluded
		p.Definitions[otherHandles["Person"]] = "pages"
	}), g.ErrScope)
	for _, file := range []files.ID{"model", "pages"} {
		add("an excluded path of the file "+string(file)+" and a missing item", typescript, source.ErrInvalid, tsChanged(func(p *tsfiles.GeneratedPlacement) {
			p.Files[file] = excluded
			delete(p.Definitions, handles["Cell"])
		}), files.ErrInvalid)
	}
	// The file without a declaration comes before or after the others in the
	// order of the identities.
	for _, empty := range []files.ID{"a-empty", "z-empty"} {
		add("a file without a declaration under the identity "+string(empty)+" and an empty file identity", typescript, source.ErrInvalid, tsChanged(func(p *tsfiles.GeneratedPlacement) {
			emptyIdentity(p)
			p.Files[empty] = "generated/empty.ts"
		}), files.ErrInvalid)
		add("a file without a declaration under the identity "+string(empty)+" and two declarations of one name in another file", typescript, source.ErrInvalid, func(t *testing.T) (bool, error) {
			out := tamperedCellPage(t, renamed, "")
			p := tsOneFile(out, "types", "types.ts")
			p.Files[empty] = "generated/empty.ts"
			return tsFormat(generated, out, format.FilesLimits{}, p.Issue())
		})
	}
	// The relocation refusals are validity. They are decided before a tree is
	// constructed and before a refusal of coverage.
	add("an excluded path and a use of an unexposed declaration of another file", typescript, source.ErrInvalid, func(t *testing.T) (bool, error) {
		return separated(t, "user.mts")
	})
	add("an excluded path and two declarations of one name in one file", typescript, source.ErrInvalid, func(t *testing.T) (bool, error) {
		out := tamperedCellPage(t, renamed, "")
		return tsFormat(generated, out, format.FilesLimits{}, tsOneFile(out, "types", "types.mts").Issue())
	})
	// The same order for a placement of existing source, whose lists also state
	// what a keyed table cannot.
	for _, defect := range []struct {
		name   string
		change func(o *tsfiles.Organization, unit int)
		also   []error
	}{
		{"an unknown declaration", func(o *tsfiles.Organization, unit int) {
			o.Units[unit].Declarations = append(o.Units[unit].Declarations, ast.NodeRef{Tree: "foreign", Node: 1})
		}, []error{files.ErrInvalid}},
		{"the path ..", func(o *tsfiles.Organization, unit int) { o.Units[unit].Path = ".." }, []error{files.ErrInvalid}},
		{"a path that is no TypeScript path", func(o *tsfiles.Organization, unit int) { o.Units[unit].Path = "generated/model/person.txt" }, nil},
		{"a repeated declaration", func(o *tsfiles.Organization, unit int) {
			o.Units[unit].Declarations = append(o.Units[unit].Declarations, o.Units[1-unit].Declarations[0])
		}, nil},
	} {
		for _, order := range []struct {
			name            string
			covered, broken int
		}{{"in the first file, then", 0, 1}, {"in the second file, after", 1, 0}} {
			add("an excluded path "+order.name+" "+defect.name, typescriptExisting, source.ErrInvalid, existingChanged(func(o *tsfiles.Organization) {
				o.Units[order.covered].Path = excluded
				defect.change(o, order.broken)
			}), defect.also...)
		}
	}
	for _, unit := range []int{0, 1} {
		place := []string{"the first", "the second"}[unit]
		add("an excluded path of "+place+" file and one file identity for two files", typescriptExisting, source.ErrInvalid, existingChanged(func(o *tsfiles.Organization) {
			o.Units[unit].Path = excluded
			o.Units[1].File = o.Units[0].File
		}), files.ErrInvalid)
	}
	empty := tsfiles.Placement{File: "empty", Path: "generated/empty.ts"}
	for _, order := range []string{"first", "last"} {
		add("a file without a declaration "+order+" and a repeated declaration", typescriptExisting, source.ErrInvalid, existingChanged(func(o *tsfiles.Organization) {
			o.Units[1].Declarations = append(o.Units[1].Declarations, o.Units[0].Declarations[0])
			if order == "first" {
				o.Units = append([]tsfiles.Placement{empty}, o.Units...)
			} else {
				o.Units = append(o.Units, empty)
			}
		}))
	}
	add("a file without a module and one file identity for two files", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) {
		*p = append(*p, horosPlacement{ID: (*p)[0].ID, Path: "generated/empty.horos"})
	}))
	// A limit is counted before the files are read.
	add("an excluded path and too many files", typescript, source.ErrLimit, func(*testing.T) (bool, error) {
		p := tsPlacement()
		p.Files["pages"] = excluded
		limits := format.DefaultFilesLimits()
		limits.Source.Files.MaxFiles = 1
		return tsFormat(generated, output, limits, p.Issue())
	}, files.ErrLimit)
	// The form of a setting is decided before the captured program is read, so a
	// script whose setting is malformed is invalid input.
	scripted := func(key semantics.Key, value string) func(*testing.T) (bool, error) {
		return func(t *testing.T) (bool, error) {
			script := tsUnit("input", "type Global = string;")
			in, err := existingStages.Capture(ctx, []source.Unit[ts.Language]{script}, script.File, semantic.Limits{})
			if err != nil {
				t.Fatal(err)
			}
			in.Configuration = []semantic.Setting{{Key: key, Value: []byte(value)}}
			out, err := plain.Format(ctx, in)
			return out.Valid(), err
		}
	}
	add("a script and malformed JSON", typescriptExisting, source.ErrInvalid, scripted(tsfiles.OrganizationSetting(), `{"Units":`))
	add("a script and another setting key", typescriptExisting, source.ErrInvalid, scripted(wrongKey, `{}`))
	add("a script and an unknown field", typescriptExisting, source.ErrInvalid, scripted(tsfiles.OrganizationSetting(), `{"Later":1}`))
	add("an excluded path and a use of an unexposed declaration of another file", typescriptExisting, source.ErrInvalid, tsExisting(&tsfiles.Organization{Units: []tsfiles.Placement{
		{File: "public", Path: "public.mts", Declarations: []ast.NodeRef{tsSite(hidden, "Visible")}},
		{File: "private", Path: "private.ts", Declarations: []ast.NodeRef{tsSite(hidden, "Hidden")}},
	}}, hidden))
	add("an excluded path and two declarations of one name in one file", typescriptExisting, source.ErrInvalid, tsExisting(&tsfiles.Organization{Units: []tsfiles.Placement{
		{File: "both", Path: "both.mts", Declarations: []ast.NodeRef{tsSite(first, "Same"), tsSite(second, "Same")}},
	}}, first, second))
	// A placement is judged against the top-level declarations of the captured
	// program before a script is refused: its items, its files and its paths need
	// nothing that only a module has.
	script := tsUnit("input", "type A = string; type B = A;")
	scriptPlaced := func(units ...tsfiles.Placement) func(*testing.T) (bool, error) {
		return func(t *testing.T) (bool, error) {
			in, err := existingStages.Capture(ctx, []source.Unit[ts.Language]{script}, script.File, semantic.Limits{})
			if err != nil {
				t.Fatal(err)
			}
			in.Configuration = []semantic.Setting{{Key: tsfiles.OrganizationSetting(), Value: tsBytes(tsfiles.Organization{Units: units})}}
			out, err := plain.Format(ctx, in)
			return out.Valid(), err
		}
	}
	scriptA, scriptB := tsSite(script, "A"), tsSite(script, "B")
	add("a script and a placement with an unknown declaration", typescriptExisting, source.ErrInvalid, scriptPlaced(tsfiles.Placement{File: "a", Path: "a.ts", Declarations: []ast.NodeRef{scriptA, scriptB, {Tree: "foreign", Node: 1}}}), files.ErrInvalid)
	add("a script and a placement with a repeated declaration", typescriptExisting, source.ErrInvalid, scriptPlaced(tsfiles.Placement{File: "a", Path: "a.ts", Declarations: []ast.NodeRef{scriptA, scriptB, scriptA}}))
	add("a script and a placement with a missing item", typescriptExisting, source.ErrInvalid, scriptPlaced(tsfiles.Placement{File: "a", Path: "a.ts", Declarations: []ast.NodeRef{scriptA}}), files.ErrInvalid)
	add("a script and a path that is no TypeScript path", typescriptExisting, source.ErrInvalid, scriptPlaced(tsfiles.Placement{File: "a", Path: "a.txt", Declarations: []ast.NodeRef{scriptA, scriptB}}))
	// What needs the bindings of the program is judged after the program is read,
	// and a script is refused while it is read. This row pins what is found
	// first: the script, and not the use of an unexposed declaration of another
	// file. The category describes the defect that was found.
	add("a script and a use of an unexposed declaration of another file", typescriptExisting, source.ErrUnsupported, scriptPlaced(
		tsfiles.Placement{File: "a", Path: "a.ts", Declarations: []ast.NodeRef{scriptA}},
		tsfiles.Placement{File: "b", Path: "b.ts", Declarations: []ast.NodeRef{scriptB}}))

	// Well formed, and outside this profile.
	add("a file that is a script", typescriptExisting, source.ErrUnsupported, tsExisting(nil, tsUnit("input", "type Global = string;")))
	add("a placed file that is a script", typescriptExisting, source.ErrUnsupported, tsExisting(&tsfiles.Organization{Units: []tsfiles.Placement{
		{File: "public", Path: "public.ts", Declarations: []ast.NodeRef{tsSite(unused, "Visible")}},
		{File: "private", Path: "private.ts", Declarations: []ast.NodeRef{tsSite(unused, "Unused")}},
	}}, unused))

	// The fresh output.
	add("the fresh output is not admitted", typescript, semantic.ErrEvidenceUnavailable, func(*testing.T) (bool, error) {
		return tsFormat(unavailable, output, format.FilesLimits{}, tsPlacement().Issue())
	})
	add("the fresh output is not admitted", horosName, semantic.ErrEvidenceUnavailable, func(*testing.T) (bool, error) {
		return horosFormat(detached, parsed, format.FilesLimits{}, semantic.Setting{Key: filesformat.OrganizationSetting(), Value: []byte(horosValid)})
	})
	// The shared contract counts the sealed output against the semantic limits
	// of the request before the placement is read, and its refusal has the limit
	// category.
	add("the sealed output exceeds the input node limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) {
		l.Semantic.MaxInputNodes = smallest(output.Snapshot(), inputNodes) - 1
	}), semantic.ErrLimit)
	add("the fresh output exceeds the input node limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) {
		l.Semantic.MaxInputNodes = smallest(output.Snapshot(), inputNodes)
	}))
	add("the fresh output exceeds the input node limit", horosName, source.ErrLimit, func(*testing.T) (bool, error) {
		limits := format.DefaultFilesLimits()
		limits.Semantic.MaxInputNodes = smallest(bundle.Snapshot(), inputNodes)
		return horosFormat(horosFormatter, bundle, limits, semantic.Setting{Key: filesformat.OrganizationSetting(), Value: []byte(horosValid)})
	}, semantic.ErrLimit)
	add("more fresh units than the input limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Semantic.MaxInputs = 1 }))
	add("more fresh units than the input limit", horosName, source.ErrLimit, func(*testing.T) (bool, error) {
		limits := format.DefaultFilesLimits()
		limits.Semantic.MaxInputs = 1
		return horosFormat(horosFormatter, bundle, limits, semantic.Setting{Key: filesformat.OrganizationSetting(), Value: []byte(horosValid)})
	}, semantic.ErrLimit)
	add("the capture of the fresh output exceeds the input byte limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) {
		l.Semantic.MaxInputBytes = smallest(output.Snapshot(), inputBytes)
	}), semantic.ErrLimit)

	// An error of the capture that is no value of the files package and no
	// semantic limit gains no category: the formatter passes it on as it is.
	// Only the capture of the fresh output is joined. A limit that the resolution
	// of the fresh output exceeds is passed on as the resolver returns it.
	add("the resolution of the fresh output exceeds the graph limit", typescript, nil, tsLimited(func(l *format.FilesLimits) {
		l.Semantic.Semantics.MaxRecords = 1
	}), semantic.ErrLimit, semantics.ErrLimit)
	add("the resolver refuses its configuration under the configuration byte limit", horosName, nil, func(t *testing.T) (bool, error) {
		// The request carries no setting. The resolver has one of its own, which
		// it refuses when it resolves the fresh output under the request's limit.
		settings := []semantic.Setting{{Key: horos.LimitsSetting(), Value: []byte("[256,100000,16777216,256]")}}
		configured, err := analysis.NewAnalyzer[s.Language](definition, horos.AnalysisContractWithResolutionConfiguration(settings), analyze.Provider{})
		if err != nil {
			t.Fatal(err)
		}
		options := semantic.Options{Configuration: settings}
		configuredParser, err := filesparse.New(definition, parsehoros.Provider{}, syntaxPrinter, checker, resolver, configured, options, semantic.Options{})
		if err != nil {
			t.Fatal(err)
		}
		formatter, err := filesformat.New(definition, parsehoros.Provider{}, formathoros.Provider{}, syntaxPrinter, checker, resolver, configured, options, semantic.Options{})
		if err != nil {
			t.Fatal(err)
		}
		from, err := configuredParser.Parse(ctx, input)
		if err != nil {
			t.Fatal(err)
		}
		if published, err := horosFormat(formatter, from, format.FilesLimits{}); err != nil || !published {
			t.Fatal("the input must format under the default limits", err)
		}
		limits := format.DefaultFilesLimits()
		limits.Semantic.MaxConfigurationBytes = 1
		return horosFormat(formatter, from, limits)
	}, semantic.ErrLimit)
	add("the capture refuses a fresh tree with an error of its checker", typescript, nil, func(t *testing.T) (bool, error) {
		stages := placementStages(t, cell)
		refusing, err := check.NewChecker[ts.Language](stages.Language, freshRefusal[ts.Language]{typescriptText(), "typescript/files/"})
		if err != nil {
			t.Fatal(err)
		}
		stages.Checker = refusing
		formatter, err := tsfiles.NewGenerated(stages)
		if err != nil {
			t.Fatal(err)
		}
		return tsFormat(formatter, output, format.FilesLimits{}, tsPlacement().Issue())
	}, errFreshTree)
	add("the capture refuses a fresh tree with an error of its checker", horosName, nil, func(t *testing.T) (bool, error) {
		refusing, err := check.NewChecker[s.Language](definition, freshRefusal[s.Language]{checkhoros.Provider{}, "horos/files/"})
		if err != nil {
			t.Fatal(err)
		}
		formatter, err := filesformat.New(definition, parsehoros.Provider{}, formathoros.Provider{}, syntaxPrinter, refusing, resolver, analyzer, semantic.Options{}, semantic.Options{})
		if err != nil {
			t.Fatal(err)
		}
		return horosFormat(formatter, parsed, format.FilesLimits{}, semantic.Setting{Key: filesformat.OrganizationSetting(), Value: []byte(horosValid)})
	}, errFreshTree)
	add("the fresh text exceeds the input byte limit", typescript, source.ErrLimit, func(*testing.T) (bool, error) {
		stages := placementStages(t, cell)
		stages.Printer = commentingPrinter{stages.Printer}
		commenting, err := tsfiles.NewGenerated(stages)
		if err != nil {
			t.Fatal(err)
		}
		limits := format.DefaultFilesLimits()
		limits.Semantic.MaxInputBytes = smallest(output.Snapshot(), inputBytes)
		if limits.Semantic.MaxInputBytes >= 1<<16 {
			t.Fatal("the comment must exceed the limit that admits the input", limits.Semantic.MaxInputBytes)
		}
		return tsFormat(commenting, output, limits, tsPlacement().Issue())
	})
	add("the fresh text exceeds the input byte limit", horosName, source.ErrLimit, func(*testing.T) (bool, error) {
		limits := format.DefaultFilesLimits()
		limits.Semantic.MaxInputBytes = smallest(bundle.Snapshot(), inputBytes)
		if limits.Semantic.MaxInputBytes >= 1<<16 {
			t.Fatal("the comment must exceed the limit that admits the input", limits.Semantic.MaxInputBytes)
		}
		return horosFormat(commented, bundle, limits)
	}, semantic.ErrLimit)

	// The limits of the placed trees.
	add("the placed declarations exceed the node limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Source.MaxNodes = 2 }))
	add("the placed declarations exceed the node limit", horosName, source.ErrLimit, horosLimited(func(l *format.FilesLimits) { l.Source.MaxNodes = 2 }))
	add("the placed declarations exceed the depth limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Source.MaxDepth = 1 }))
	add("the placed declarations exceed the depth limit", horosName, source.ErrLimit, horosLimited(func(l *format.FilesLimits) { l.Source.MaxDepth = 1 }))
	add("the placed declarations exceed the work limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Source.MaxWork = 1 }))
	add("the placed declarations exceed the work limit", horosName, source.ErrLimit, horosLimited(func(l *format.FilesLimits) { l.Source.MaxWork = 1 }))
	add("the root of a further file exceeds the node limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) {
		// The first file takes every node that the limit allows: its root and the
		// nodes of its declarations. The first file is the one of least identity,
		// model, which holds Person.
		var count func(ast.NodeID) int
		count = func(id ast.NodeID) int {
			nodes := 1
			for _, child := range tree.Nodes[id-1].Children {
				nodes += count(child.Node)
			}
			return nodes
		}
		person, err := output.Resolve(handles["Person"])
		if err != nil {
			t.Fatal(err)
		}
		l.Source.MaxNodes = 1 + count(person.Site.Node)
	}))
	add("the constructed imports exceed the node limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Source.MaxNodes = len(tree.Nodes) + 1 }))
	add("the constructed imports exceed the byte limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Source.MaxBytes = 1 }))

	for _, c := range cases {
		t.Run(c.formatter+"/"+c.defect, func(t *testing.T) {
			published, err := c.run(t)
			if err == nil || published {
				t.Fatalf("accepted: published=%v err=%v", published, err)
			}
			matched := 0
			for _, category := range filesCategories {
				is := errors.Is(err, category.value)
				if is {
					matched++
				}
				if is != (category.value == c.want) {
					t.Errorf("errors.Is(err, %s) = %v for %q", category.name, is, err)
				}
			}
			expected := 1
			if c.want == nil {
				expected = 0
			}
			if matched != expected {
				t.Errorf("the refusal matches %d categories: %q", matched, err)
			}
			for _, original := range c.also {
				if !errors.Is(err, original) {
					t.Errorf("the refusal no longer matches %q: %q", original, err)
				}
			}
			if c.check != nil {
				c.check(t, err)
			}
		})
	}
}

// TestHorosEstablishJoinsTheFilesValue calls the step that the horos Files
// parser and formatter share. Each formatter call hands it an inventory that
// the formatter has checked, so only a direct call reaches the capture's
// refusal of an invalid inventory.
func TestHorosEstablishJoinsTheFilesValue(t *testing.T) {
	ctx := context.Background()
	input, _, _ := filePagination(t)
	definition, err := representation.NewLanguage[s.Language](s.Language{}, s.Profile)
	if err != nil {
		t.Fatal(err)
	}
	checker, resolver, analyzer := fileStages(t, definition)
	syntaxPrinter := textfixture.Printing[s.Language](s.Language{})
	parser, err := filesparse.New(definition, parsehoros.Provider{}, syntaxPrinter, checker, resolver, analyzer, semantic.Options{}, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parser.Parse(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	stages, err := (horosfiles.Stages{Language: definition, Parser: parsehoros.Provider{}, Formatter: formathoros.Provider{}, Printer: syntaxPrinter, Checker: checker, Resolver: resolver, Analyzer: analyzer}).Capture()
	if err != nil {
		t.Fatal(err)
	}
	units := parsed.Units()
	if _, err := stages.Establish(ctx, units, "people", semantic.DefaultLimits()); err != nil {
		t.Fatal("the parsed units must be established", err)
	}
	units[1].Path = units[0].Path
	publication, err := stages.Establish(ctx, units, "people", semantic.DefaultLimits())
	if err == nil || publication.Snapshot.Valid() {
		t.Fatal("two units with one path were established")
	}
	for _, category := range filesCategories {
		if errors.Is(err, category.value) != (category.value == source.ErrInvalid) {
			t.Errorf("errors.Is(err, %s) = %v for %q", category.name, errors.Is(err, category.value), err)
		}
	}
	if !errors.Is(err, files.ErrInvalid) {
		t.Errorf("the refusal no longer matches %q: %q", files.ErrInvalid, err)
	}
}
