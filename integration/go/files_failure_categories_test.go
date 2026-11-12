package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

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
	filesparse "github.com/Bitspark/bittype/parse/files/horos/provider/go"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	"github.com/Bitspark/bittype/semantics/go"
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

// commentingPrinter is a TypeScript Syntax printer that writes a long comment
// after each unit, so that the fresh text exceeds a byte limit that the input met.
type commentingPrinter struct {
	print.Provider[cst.CST[ts.Language], string]
}

func (p commentingPrinter) Print(ctx context.Context, unit cst.CST[ts.Language]) (string, error) {
	text, err := p.Provider.Print(ctx, unit)
	return text + "\n// " + strings.Repeat("x", 1<<16) + "\n", err
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

	// TypeScript, generated declarations.
	output := generatedCellPage(t)
	tree := outputTree[ts.Language](t, output)
	cell := n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations}
	generated, err := tsfiles.NewGenerated(placementStages(t, cell))
	if err != nil {
		t.Fatal(err)
	}
	sites := map[string]ast.NodeRef{}
	for _, child := range tree.Nodes[tree.Root-1].Children {
		sites[nodeAttribute(tree.Nodes[child.Node-1], ts.AttributeName)] = ast.NodeRef{Tree: tree.ID, Node: child.Node}
	}
	tsPlacement := func() tsfiles.Organization {
		return tsfiles.Organization{Entry: "pages", Units: []tsfiles.Placement{
			{File: "pages", Path: "generated/pages/page.ts", Declarations: []ast.NodeRef{sites["Page"], sites["Cells"], sites["Cell"]}},
			{File: "model", Path: "generated/model/person.ts", Declarations: []ast.NodeRef{sites["Person"]}},
		}}
	}
	tsBytes := func(o tsfiles.Organization) []byte {
		setting, err := o.Setting()
		if err != nil {
			t.Fatal(err)
		}
		return setting.Value
	}
	tsFormat := func(f format.GeneratedFilesFormatter[ts.Language], out g.Output, limits format.FilesLimits, settings ...semantic.Setting) (bool, error) {
		result, err := f.Format(ctx, format.GeneratedFilesInput{Output: out, Configuration: settings, Limits: limits})
		return result.Valid(), err
	}
	tsRaw := func(value string) func(*testing.T) (bool, error) {
		return func(*testing.T) (bool, error) {
			return tsFormat(generated, output, format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: []byte(value)})
		}
	}
	tsPlaced := func(change func(*tsfiles.Organization)) func(*testing.T) (bool, error) {
		return func(*testing.T) (bool, error) {
			o := tsPlacement()
			change(&o)
			return tsFormat(generated, output, format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: tsBytes(o)})
		}
	}
	tsLimited := func(change func(*format.FilesLimits)) func(*testing.T) (bool, error) {
		return func(*testing.T) (bool, error) {
			limits := format.DefaultFilesLimits()
			change(&limits)
			return tsFormat(generated, output, limits, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: tsBytes(tsPlacement())})
		}
	}
	if published, err := tsPlaced(func(*tsfiles.Organization) {})(t); err != nil || !published {
		t.Fatal("the TypeScript placement that every case changes must format", err)
	}

	// TypeScript, existing source.
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
			tree := outputTree[ts.Language](t, out)
			placement := tsfiles.Placement{File: "types", Path: "types.ts"}
			for _, child := range tree.Nodes[tree.Root-1].Children {
				if child.Role == ts.RoleDeclaration {
					placement.Declarations = append(placement.Declarations, ast.NodeRef{Tree: tree.ID, Node: child.Node})
				}
			}
			return tsFormat(generated, out, format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: tsBytes(tsfiles.Organization{Entry: "types", Units: []tsfiles.Placement{placement}})})
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
	valid := string(tsBytes(tsPlacement()))
	horosValid := string(horosBytes(horosPlacements()))

	// The setting.
	add("wrong setting key", typescript, source.ErrInvalid, func(*testing.T) (bool, error) {
		return tsFormat(generated, output, format.FilesLimits{}, semantic.Setting{Key: wrongKey, Value: []byte(valid)})
	})
	add("wrong setting key", horosName, source.ErrInvalid, func(*testing.T) (bool, error) {
		return horosFormat(horosFormatter, parsed, format.FilesLimits{}, semantic.Setting{Key: wrongKey, Value: []byte(horosValid)})
	}, semantic.ErrInvalidRequest)
	add("two settings", typescript, source.ErrInvalid, func(*testing.T) (bool, error) {
		return tsFormat(generated, output, format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: []byte(valid)}, semantic.Setting{Key: wrongKey, Value: []byte(valid)})
	})
	add("two settings", horosName, source.ErrInvalid, func(*testing.T) (bool, error) {
		return horosFormat(horosFormatter, parsed, format.FilesLimits{}, semantic.Setting{Key: filesformat.OrganizationSetting(), Value: []byte(horosValid)}, semantic.Setting{Key: wrongKey, Value: []byte(horosValid)})
	}, semantic.ErrInvalidRequest)
	add("malformed JSON", typescript, source.ErrInvalid, tsRaw(`{"Entry":"pages","Units":[{"File":}]}`))
	add("malformed JSON", horosName, source.ErrInvalid, horosRaw(`[{"ID":}]`))
	cases[len(cases)-1].check = func(t *testing.T, err error) {
		var syntax *json.SyntaxError
		if !errors.As(err, &syntax) {
			t.Errorf("the decoder's error is no longer part of the refusal: %q", err)
		}
	}
	add("truncated JSON", typescript, source.ErrInvalid, tsRaw(valid[:len(valid)-1]))
	add("truncated JSON", horosName, source.ErrInvalid, horosRaw(horosValid[:len(horosValid)-1]))
	add("a JSON value of another type", typescript, source.ErrInvalid, tsRaw(`[]`))
	add("a JSON value of another type", horosName, source.ErrInvalid, horosRaw(`{}`))
	add("a second JSON value", typescript, source.ErrInvalid, tsRaw(valid+` {}`))
	add("a second JSON value", horosName, source.ErrInvalid, horosRaw(horosValid+` []`), semantic.ErrInvalidRequest)

	// No field is accepted and ignored.
	add("unknown field", typescript, source.ErrInvalid, tsRaw(`{"Later":1,`+valid[1:]))
	add("unknown field of a file", typescript, source.ErrInvalid, tsRaw(strings.Replace(valid, `{"File":`, `{"Later":1,"File":`, 1)))
	add("unknown field of an item", typescript, source.ErrInvalid, tsRaw(strings.Replace(valid, `{"Tree":`, `{"Later":1,"Tree":`, 1)))
	add("unknown field", horosName, source.ErrInvalid, horosRaw(`[{"Later":1,`+horosValid[2:]))

	// The inventory of files.
	add("too many files", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Source.Files.MaxFiles = 1 }))
	add("too many files", horosName, source.ErrLimit, horosLimited(func(l *format.FilesLimits) { l.Source.Files.MaxFiles = 1 }))
	add("file identities and paths exceed the byte limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Source.Files.MaxBytes = 8 }), files.ErrLimit)
	add("file identities and paths exceed the byte limit", horosName, source.ErrLimit, horosLimited(func(l *format.FilesLimits) { l.Source.Files.MaxBytes = 8 }), files.ErrLimit)
	add("one file identity for two files", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) { o.Units[1].File = o.Units[0].File }), files.ErrInvalid)
	add("one file identity for two files", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].ID = (*p)[0].ID }), files.ErrInvalid)
	add("empty file identity", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) { o.Units[1].File = "" }), files.ErrInvalid)
	add("empty file identity", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].ID = "" }), files.ErrInvalid)
	add("one path for two files", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) { o.Units[1].Path = o.Units[0].Path }), files.ErrInvalid)
	add("one path for two files", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].Path = (*p)[0].Path }), files.ErrInvalid)
	add("a path that is also a directory of another", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) { o.Units[1].Path = o.Units[0].Path + "/person.ts" }), files.ErrInvalid)
	add("a path that is also a directory of another", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].Path = (*p)[0].Path + "/page.horos" }), files.ErrInvalid)
	add("the path ..", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) { o.Units[1].Path = ".." }))
	add("the path ..", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].Path = ".." }), files.ErrInvalid)
	add("an absolute path", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) { o.Units[1].Path = "/generated/person.ts" }))
	add("an absolute path", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].Path = "/generated/page.horos" }), files.ErrInvalid)
	add("a valid path that is not a TypeScript path", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) { o.Units[1].Path = "generated/model/person.txt" }))
	add("a declaration file path", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) { o.Units[1].Path = "generated/model/person.d.ts" }))

	// The items of the placement.
	add("a file without an item", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) {
		o.Units = append(o.Units, tsfiles.Placement{File: "empty", Path: "generated/empty.ts"})
	}))
	add("a file without an item", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) {
		*p = append(*p, horosPlacement{ID: "empty", Path: "generated/empty.horos"})
	}))
	add("no file", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) { o.Units = nil }))
	add("no file", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { *p = nil }))
	add("unknown item", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) {
		o.Units[1].Declarations = append(o.Units[1].Declarations, ast.NodeRef{Tree: "foreign", Node: 1})
	}))
	add("unknown item", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].Modules = append((*p)[1].Modules, "0000") }))
	add("repeated item", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) {
		o.Units[1].Declarations = append(o.Units[1].Declarations, sites["Page"])
	}))
	add("repeated item", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { (*p)[1].Modules = append((*p)[1].Modules, peoplePin) }))
	add("missing item", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) { o.Units = o.Units[:1] }))
	add("missing item", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) { *p = (*p)[1:] }))
	add("missing entry", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) { o.Entry = "absent" }))

	// What the placement asks of the declarations.
	// A placement that separates an unexposed declaration from its user is
	// invalid input: no later profile can admit it for the same output, because
	// Files never exports a declaration without authorization.
	add("a declaration uses an unexposed one of another file", typescriptExisting, source.ErrInvalid, tsExisting(&tsfiles.Organization{Entry: "public", Units: []tsfiles.Placement{
		{File: "public", Path: "public.ts", Declarations: []ast.NodeRef{tsSite(hidden, "Visible")}},
		{File: "private", Path: "private.ts", Declarations: []ast.NodeRef{tsSite(hidden, "Hidden")}},
	}}, hidden))
	add("a declaration uses an unexposed one of another file", typescript, source.ErrInvalid, func(t *testing.T) (bool, error) {
		out := tamperedCellPage(t, func(g.BuildInput[ts.Language], *g.Assembly) error { return nil }, "Cells")
		tree := outputTree[ts.Language](t, out)
		used := tsfiles.Placement{File: "used", Path: "used.ts"}
		user := tsfiles.Placement{File: "user", Path: "user.ts"}
		for _, child := range tree.Nodes[tree.Root-1].Children {
			site := ast.NodeRef{Tree: tree.ID, Node: child.Node}
			if nodeAttribute(tree.Nodes[child.Node-1], ts.AttributeName) == "Cells" {
				user.Declarations = append(user.Declarations, site)
			} else {
				used.Declarations = append(used.Declarations, site)
			}
		}
		if exposed := out.Exposed(); len(exposed) != 1 || len(user.Declarations) != 1 || exposed[0] != user.Declarations[0] || len(used.Declarations) == 0 {
			t.Fatal("expected Cells as the one exposed declaration", exposed)
		}
		together := tsBytes(tsfiles.Organization{Entry: "user", Units: []tsfiles.Placement{{File: "user", Path: "user.ts", Declarations: append(append([]ast.NodeRef{}, user.Declarations...), used.Declarations...)}}})
		if published, err := tsFormat(generated, out, format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: together}); err != nil || !published {
			t.Fatal("the same declarations in one file must format", err)
		}
		return tsFormat(generated, out, format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: tsBytes(tsfiles.Organization{Entry: "user", Units: []tsfiles.Placement{user, used}})})
	})
	add("two declarations of one name in one file", typescriptExisting, source.ErrInvalid, tsExisting(&tsfiles.Organization{Entry: "both", Units: []tsfiles.Placement{
		{File: "both", Path: "both.ts", Declarations: []ast.NodeRef{tsSite(first, "Same"), tsSite(second, "Same")}},
	}}, first, second))

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
		assembly.ExternalUses[0].At.Site = assembly.Exposed[0]
		return nil
	}))
	add("a name use without a destination", typescript, source.ErrInvalid, tampered(func(_ g.BuildInput[ts.Language], assembly *g.Assembly) error {
		assembly.ExternalUses = nil
		return nil
	}))
	add("an exposed site that is no top-level declaration", typescript, source.ErrInvalid, tampered(func(_ g.BuildInput[ts.Language], assembly *g.Assembly) error {
		assembly.Exposed[0] = assembly.ExternalUses[0].At.Site
		return nil
	}))
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

	// Well formed, and outside this profile.
	add("a file that is a script", typescriptExisting, source.ErrUnsupported, tsExisting(nil, tsUnit("input", "type Global = string;")))
	add("a placed file that is a script", typescriptExisting, source.ErrUnsupported, tsExisting(&tsfiles.Organization{Entry: "public", Units: []tsfiles.Placement{
		{File: "public", Path: "public.ts", Declarations: []ast.NodeRef{tsSite(unused, "Visible")}},
		{File: "private", Path: "private.ts", Declarations: []ast.NodeRef{tsSite(unused, "Unused")}},
	}}, unused))

	// The fresh output.
	add("the fresh output is not admitted", typescript, semantic.ErrEvidenceUnavailable, func(*testing.T) (bool, error) {
		return tsFormat(unavailable, output, format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: []byte(valid)})
	})
	add("the fresh output is not admitted", horosName, semantic.ErrEvidenceUnavailable, func(*testing.T) (bool, error) {
		return horosFormat(detached, parsed, format.FilesLimits{}, semantic.Setting{Key: filesformat.OrganizationSetting(), Value: []byte(horosValid)})
	})
	add("the fresh output exceeds the input node limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) {
		l.Semantic.MaxInputNodes = smallest(output.Snapshot(), inputNodes)
	}))
	add("the fresh output exceeds the input node limit", horosName, source.ErrLimit, func(*testing.T) (bool, error) {
		limits := format.DefaultFilesLimits()
		limits.Semantic.MaxInputNodes = smallest(bundle.Snapshot(), inputNodes)
		return horosFormat(horosFormatter, bundle, limits, semantic.Setting{Key: filesformat.OrganizationSetting(), Value: []byte(horosValid)})
	}, semantic.ErrLimit)
	add("more fresh units than the input limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Semantic.MaxInputs = 1 }))
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
		return tsFormat(commenting, output, limits, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: []byte(valid)})
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
		// nodes of its declarations.
		var count func(ast.NodeID) int
		count = func(id ast.NodeID) int {
			nodes := 1
			for _, child := range tree.Nodes[id-1].Children {
				nodes += count(child.Node)
			}
			return nodes
		}
		l.Source.MaxNodes = 1
		for _, site := range tsPlacement().Units[0].Declarations {
			l.Source.MaxNodes += count(site.Node)
		}
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
			if matched != 1 {
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
