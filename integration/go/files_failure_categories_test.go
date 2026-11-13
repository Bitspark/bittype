package integration_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	checkhoros "github.com/Bitspark/bittype/check/syntax/horos/provider/go"
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
	horosDecoded("malformed JSON", `[{"ID":}]`)
	add("truncated JSON", typescript, source.ErrInvalid, tsRaw(valid[:len(valid)-1]))
	horosDecoded("truncated JSON", horosValid[:len(horosValid)-1])
	add("a JSON value of another type", typescript, source.ErrInvalid, tsRaw(`[]`))
	horosDecoded("a JSON value of another type", `{}`)
	add("an empty setting value", typescript, source.ErrInvalid, tsRaw(``))
	horosDecoded("an empty setting value", ``)
	add("the JSON value null", typescript, source.ErrInvalid, tsRaw(`null`))
	add("the JSON value null", horosName, source.ErrInvalid, horosRaw(`null`), files.ErrInvalid)
	add("a second JSON value", typescript, source.ErrInvalid, tsRaw(valid+` {}`))
	add("a second JSON value", horosName, source.ErrInvalid, horosRaw(horosValid+` []`), semantic.ErrInvalidRequest)

	// No field is accepted and ignored.
	add("unknown field", typescript, source.ErrInvalid, tsRaw(`{"Later":1,`+valid[1:]))
	add("unknown field of a file", typescript, source.ErrInvalid, tsRaw(strings.Replace(valid, `{"File":`, `{"Later":1,"File":`, 1)))
	add("unknown field of an item", typescript, source.ErrInvalid, tsRaw(strings.Replace(valid, `{"Tree":`, `{"Later":1,"Tree":`, 1)))
	horosDecoded("unknown field", `[{"Later":1,`+horosValid[2:])

	// A member that is named twice is refused. The decoder would keep the last
	// value, so the first would be accepted and ignored. A member name in another
	// case names the same member, and one such name alone stays admitted.
	for name, run := range map[string]func(*testing.T) (bool, error){
		"TypeScript, a member name in another case": tsRaw(strings.Replace(strings.Replace(valid, `"Entry"`, `"entry"`, 1), `"Units"`, `"UNITS"`, 1)),
		"TypeScript, trailing white space":          tsRaw(valid + " \n\t"),
		"horos, a member name in another case":      horosRaw(strings.Replace(horosValid, `"ID"`, `"id"`, 1)),
		"horos, trailing white space":               horosRaw(horosValid + " \n"),
	} {
		if published, err := run(t); err != nil || !published {
			t.Fatal("an admitted placement is refused:", name, err)
		}
	}
	add("a repeated member", typescript, source.ErrInvalid, tsRaw(`{"Entry":"absent",`+valid[1:]))
	add("a repeated member whose first value is another placement", typescript, source.ErrInvalid, tsRaw(`{"Units":[{"File":"x","Path":"../x","Declarations":[]}],`+valid[1:]))
	add("a repeated member of a file", typescript, source.ErrInvalid, tsRaw(strings.Replace(valid, `{"File":`, `{"File":"other","File":`, 1)))
	add("a repeated member of an item", typescript, source.ErrInvalid, tsRaw(strings.Replace(valid, `{"Tree":`, `{"Tree":"other","Tree":`, 1)))
	add("a member repeated in another case", typescript, source.ErrInvalid, tsRaw(`{"entry":"absent",`+valid[1:]))
	add("a repeated member", horosName, source.ErrInvalid, horosRaw(`[{"ID":"../other",`+horosValid[2:]))
	add("a member repeated in another case", horosName, source.ErrInvalid, horosRaw(`[{"id":"../other",`+horosValid[2:]))

	// The inventory of files.
	add("too many files", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Source.Files.MaxFiles = 1 }))
	add("too many files", horosName, source.ErrLimit, horosLimited(func(l *format.FilesLimits) { l.Source.Files.MaxFiles = 1 }))
	add("file identities and paths exceed the byte limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) { l.Source.Files.MaxBytes = 8 }), files.ErrLimit)
	add("file identities and paths exceed the byte limit", horosName, source.ErrLimit, horosLimited(func(l *format.FilesLimits) { l.Source.Files.MaxBytes = 8 }), files.ErrLimit)
	add("the fresh files exceed the Files byte limit", typescript, source.ErrLimit, tsLimited(func(l *format.FilesLimits) {
		// The identities and paths fit. The text of the fresh files does not.
		l.Source.Files.MaxBytes = 0
		for _, unit := range tsPlacement().Units {
			l.Source.Files.MaxBytes += len(unit.File) + len(unit.Path)
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
		o := tsPlacement()
		o.Units[0].Path = files.Path("generated/" + strings.Repeat("d", 6000) + "/page.ts")
		setting := semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: tsBytes(o)}
		if published, err := tsFormat(generated, output, format.FilesLimits{}, setting); err != nil || !published {
			t.Fatal("the long path must format under the default limits", err)
		}
		limits := format.DefaultFilesLimits()
		limits.Semantic.MaxInputBytes = 5000
		if published, err := tsFormat(generated, output, limits, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: []byte(valid)}); err != nil || !published {
			t.Fatal("the short path must format under the limit", err)
		}
		return tsFormat(generated, output, limits, setting)
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
	// A path that is no TypeScript path is invalid input. A path that TypeScript
	// itself reads and that this profile excludes is unsupported coverage.
	tsPath := func(path string) func(*testing.T) (bool, error) {
		return tsPlaced(func(o *tsfiles.Organization) { o.Units[1].Path = files.Path(path) })
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
	add("a file without a declaration", typescript, source.ErrUnsupported, tsPlaced(func(o *tsfiles.Organization) {
		o.Units = append(o.Units, tsfiles.Placement{File: "empty", Path: "generated/empty.ts"})
	}))
	formatted := 0
	add("a file without a declaration, before any fresh file is formatted", typescript, source.ErrUnsupported, func(t *testing.T) (bool, error) {
		stages := placementStages(t, cell)
		stages.Syntax = countingSyntax{stages.Syntax, &formatted}
		counting, err := tsfiles.NewGenerated(stages)
		if err != nil {
			t.Fatal(err)
		}
		if published, err := tsFormat(counting, output, format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: []byte(valid)}); err != nil || !published || formatted != 2 {
			t.Fatal("the base placement must format its two files", err, formatted)
		}
		formatted = 0
		o := tsPlacement()
		o.Units = append(o.Units, tsfiles.Placement{File: "empty", Path: "generated/empty.ts"})
		return tsFormat(counting, output, format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: tsBytes(o)})
	})
	cases[len(cases)-1].check = func(t *testing.T, _ error) {
		if formatted != 0 {
			t.Errorf("the placement was refused after %d fresh files were formatted", formatted)
		}
	}
	add("a module of imports only in its original placement", typescriptExisting, source.ErrUnsupported, tsExisting(nil, tsUnit("a", "export type A = string;"), tsUnit("user", "import type { A } from \"./a.js\";")))
	add("a file without a module", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) {
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
	// has the defect that a later slice could cover.
	const excluded = "generated/elsewhere/file.mts"
	coverage := map[string]func(o *tsfiles.Organization, unit int){
		"an excluded path": func(o *tsfiles.Organization, unit int) { o.Units[unit].Path = excluded },
	}
	validity := []struct {
		name   string
		change func(o *tsfiles.Organization, unit int)
		also   []error
	}{
		{"an unknown declaration", func(o *tsfiles.Organization, unit int) {
			o.Units[unit].Declarations = append(o.Units[unit].Declarations, ast.NodeRef{Tree: "foreign", Node: 1})
		}, nil},
		{"the path ..", func(o *tsfiles.Organization, unit int) { o.Units[unit].Path = ".." }, nil},
		{"a path that is no TypeScript path", func(o *tsfiles.Organization, unit int) { o.Units[unit].Path = "generated/model/person.txt" }, nil},
		{"a repeated declaration", func(o *tsfiles.Organization, unit int) {
			o.Units[unit].Declarations = append(o.Units[unit].Declarations, o.Units[1-unit].Declarations[0])
		}, nil},
	}
	for name, cover := range coverage {
		for _, defect := range validity {
			for _, order := range []struct {
				name            string
				covered, broken int
			}{{"in the first file, then", 0, 1}, {"in the second file, after", 1, 0}} {
				add(name+" "+order.name+" "+defect.name, typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) {
					cover(o, order.covered)
					defect.change(o, order.broken)
				}), defect.also...)
			}
		}
	}
	add("an excluded path and an unknown declaration in one file", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) {
		o.Units[0].Path = excluded
		o.Units[0].Declarations = append(o.Units[0].Declarations, ast.NodeRef{Tree: "foreign", Node: 1})
	}))
	for _, unit := range []int{0, 1} {
		place := []string{"the first", "the second"}[unit]
		add("an excluded path of "+place+" file and a missing entry", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) {
			o.Units[unit].Path = excluded
			o.Entry = "absent"
		}))
		add("an excluded path of "+place+" file and a missing item", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) {
			o.Units[unit].Path = excluded
			o.Units[0].Declarations = o.Units[0].Declarations[:2]
		}))
		add("an excluded path of "+place+" file and one file identity for two files", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) {
			o.Units[unit].Path = excluded
			o.Units[1].File = o.Units[0].File
		}), files.ErrInvalid)
	}
	empty := tsfiles.Placement{File: "empty", Path: "generated/empty.ts"}
	for _, order := range []string{"first", "last"} {
		with := func(o *tsfiles.Organization) {
			if order == "first" {
				o.Units = append([]tsfiles.Placement{empty}, o.Units...)
			} else {
				o.Units = append(o.Units, empty)
			}
		}
		add("a file without a declaration "+order+" and an empty file identity", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) {
			o.Units[1].File = ""
			with(o)
		}), files.ErrInvalid)
		add("a file without a declaration "+order+" and a repeated declaration", typescript, source.ErrInvalid, tsPlaced(func(o *tsfiles.Organization) {
			o.Units[1].Declarations = append(o.Units[1].Declarations, o.Units[0].Declarations[0])
			with(o)
		}))
	}
	add("a file without a module and one file identity for two files", horosName, source.ErrInvalid, horosPlaced(func(p *[]horosPlacement) {
		*p = append(*p, horosPlacement{ID: (*p)[0].ID, Path: "generated/empty.horos"})
	}))
	// A limit is counted before the files are read.
	add("an excluded path and too many files", typescript, source.ErrLimit, func(*testing.T) (bool, error) {
		o := tsPlacement()
		o.Units[0].Path = excluded
		limits := format.DefaultFilesLimits()
		limits.Source.Files.MaxFiles = 1
		return tsFormat(generated, output, limits, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: tsBytes(o)})
	})
	// The order of two checks is no part of the contract where one of them runs
	// before the placement is read, or after it is admitted. These rows pin what
	// is found first today: the category describes that defect and says nothing
	// about the rest of the input.
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
	add("a script and malformed JSON", typescriptExisting, source.ErrUnsupported, scripted(tsfiles.OrganizationSetting(), `{"Entry":`))
	add("a script and another setting key", typescriptExisting, source.ErrUnsupported, scripted(wrongKey, `{}`))
	add("a script and an unknown field", typescriptExisting, source.ErrUnsupported, scripted(tsfiles.OrganizationSetting(), `{"Later":1}`))
	add("an excluded path and a use of an unexposed declaration of another file", typescriptExisting, source.ErrUnsupported, tsExisting(&tsfiles.Organization{Entry: "public", Units: []tsfiles.Placement{
		{File: "public", Path: "public.mts", Declarations: []ast.NodeRef{tsSite(hidden, "Visible")}},
		{File: "private", Path: "private.ts", Declarations: []ast.NodeRef{tsSite(hidden, "Hidden")}},
	}}, hidden))
	add("an excluded path and two declarations of one name in one file", typescriptExisting, source.ErrUnsupported, tsExisting(&tsfiles.Organization{Entry: "both", Units: []tsfiles.Placement{
		{File: "both", Path: "both.mts", Declarations: []ast.NodeRef{tsSite(first, "Same"), tsSite(second, "Same")}},
	}}, first, second))

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
		return tsFormat(formatter, output, format.FilesLimits{}, semantic.Setting{Key: tsfiles.OrganizationSetting(), Value: []byte(valid)})
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
