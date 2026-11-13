package integration_test

// Gate line 1 (bittype#196), native projection: "The native type is computed
// from the frozen record, by a projector that cannot read the source". The
// second layer of row 39 is in gate_entry_test.go.

import (
	"context"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
	"slices"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	gorecordvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/record/go"
	naming "github.com/Bitspark/bittype/transform/syntax/horos/mapping/naming/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// gatePagination is the pagination module in horos text, with a generic of two
// binders beside it, so that a tie to a binder can be told from a position.
const gatePagination = `module example.pagination;

type Page<T> = { items: list<T>; next?: text; };
type Person = { name: text; active: bool; email?: text; };
type PeoplePage = Page<Person>;
type NestedPage = Page<Page<Person>>;
type Names = list<text>;
type Pair<A, B> = { first: A; second: B; };
type Both = Pair<text, bool>;
`

// gateSourceReach lists what a capability type can reach of the source: a
// type that is, or that mentions as a type argument, a kernel handle of the
// source, a selector, a role, a step of a rule or a source reader. It follows
// the exported fields, the element types and the method sets of every type it
// meets. A value of an interface type is not followed: what it holds is not
// part of the type.
func gateSourceReach(capability reflect.Type) []string {
	forbidden := []string{
		"generation.Use[", "generation.UseView[", "generation.Declaration[", "generation.DeclarationView[", "generation.Binder[",
		"generation.Selector", "generation.Captured[", "generation.NameStep[", "generation.DecideStep[", "generation.Subject[",
		"generation.FieldPosition[", "generation.Role[", "generation.RoleIssuer[", "generation.Interpretation[", "generation.Source[",
		"generation.SourceProvider[", "generation.Request[", "generation.Answer", "source.Reader", "generationfixture.Reader",
	}
	seen := map[reflect.Type]bool{}
	var reach []string
	var visit func(t reflect.Type, path string)
	visit = func(t reflect.Type, path string) {
		if seen[t] {
			return
		}
		seen[t] = true
		name := t.String()
		for _, f := range forbidden {
			if strings.Contains(name, f) {
				reach = append(reach, path+" reaches "+name)
				return
			}
		}
		switch t.Kind() {
		case reflect.Pointer, reflect.Slice, reflect.Array, reflect.Chan:
			visit(t.Elem(), path)
		case reflect.Map:
			visit(t.Key(), path)
			visit(t.Elem(), path)
		case reflect.Func:
			for i := 0; i < t.NumIn(); i++ {
				visit(t.In(i), path)
			}
			for i := 0; i < t.NumOut(); i++ {
				visit(t.Out(i), path)
			}
		case reflect.Struct:
			for i := 0; i < t.NumField(); i++ {
				if field := t.Field(i); field.IsExported() {
					visit(field.Type, path+"."+field.Name)
				}
			}
		}
		for _, receiver := range []reflect.Type{t, reflect.PointerTo(t)} {
			for i := 0; i < receiver.NumMethod(); i++ {
				m := receiver.Method(i)
				first := 1
				if receiver.Kind() == reflect.Interface {
					first = 0
				}
				for j := first; j < m.Type.NumIn(); j++ {
					visit(m.Type.In(j), path+"."+m.Name+"()")
				}
				for j := 0; j < m.Type.NumOut(); j++ {
					visit(m.Type.Out(j), path+"."+m.Name+"()")
				}
			}
		}
	}
	visit(capability, capability.Name())
	return reach
}

// Row 4, gate line 1.
//
// Shows: the capability that a projector is given, Projection, has no method
// and no reachable exported field whose type mentions a use, a declaration, a
// binder, a selector, a role, a step of a rule or a source reader, for the Go
// and the TypeScript target. The same reflection finds the reader and the
// subject in the two steps of a rule, so the check can see what it looks for.
// A projector that calls a reader, a use, a declaration, a subject or a
// selector on its projection does not type-check, and one that reads the
// payload does.
// Does not show: that a contribution cannot hand what its rule read to its
// projector through state of its own package. The page says so ("A limit").
// It does not show what a value behind an interface type holds: Payload
// returns the decoded view as any, and its type is the contribution's own.
// Red under: L1-projection-reaches-selector.
func TestGateProjectorCannotReadTheSource(t *testing.T) {
	for name, capability := range map[string]reflect.Type{
		"golang":     reflect.TypeFor[g.Projection[gs.Language]](),
		"typescript": reflect.TypeFor[g.Projection[ts.Language]](),
	} {
		if reach := gateSourceReach(capability); len(reach) != 0 {
			t.Errorf("%s: a projector can reach the source: %v", name, reach)
		}
	}
	// The control: the same check sees the reader and the subject of a rule's steps.
	for name, capability := range map[string]reflect.Type{
		"step 1 of a rule": reflect.TypeFor[g.NameStep[h.Language, source.Reader]](),
		"step 3 of a rule": reflect.TypeFor[g.DecideStep[h.Language, source.Reader]](),
	} {
		if reach := gateSourceReach(capability); len(reach) == 0 {
			t.Errorf("the reflection does not see the source in %s", name)
		}
	}
	compile := func(body string) error {
		text := `package probe

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
)

var _ g.Projector[s.Language] = func(ctx context.Context, p g.Projection[s.Language]) (g.Projected[s.Language], error) {
` + body + `
	return g.Projected[s.Language]{}, nil
}
`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "probe.go", text, parser.AllErrors)
		if err != nil {
			return err
		}
		_, err = (&types.Config{Importer: importer.ForCompiler(fset, "source", nil)}).Check("probe", fset, []*ast.File{file}, nil)
		return err
	}
	if err := compile("_ = p.Payload()\n_ = p.Dependency(0).Type\n_ = p.Stated()"); err != nil {
		t.Fatalf("a projector that reads the frozen record does not type-check: %v", err)
	}
	for _, attempt := range []string{"_, _ = p.Reader()", "_ = p.Use()", "_ = p.Declaration()", "_ = p.Subject", "_ = p.Selector()", "_ = p.Dependency(0).Use"} {
		if err := compile(attempt); err == nil {
			t.Errorf("a projector with %q type-checks", attempt)
		}
	}
}

// gateRetaining wraps the Go record contribution so that it keeps the value
// that its rule returned and mutates it when its projector starts, which is
// after the rule returned and before the planner writes the record.
func gateRetaining(base goContribution) goContribution {
	var retained [][]string
	decide, project := base.Rule.Decide, base.Project
	base.Rule.Decide = func(ctx context.Context, step horosDecideStep) (g.Decision, error) {
		d, err := decide(ctx, step)
		if payload, ok := d.Payload.(gorecordvocabulary.Payload); ok {
			retained = append(retained, payload.Members)
		}
		return d, err
	}
	base.Project = func(ctx context.Context, p goProjection) (goProjected, error) {
		for _, members := range retained {
			for i := range members {
				if !strings.HasPrefix(members[i], "Mutated") {
					members[i] = "Mutated" + members[i]
				}
			}
		}
		return project(ctx, p)
	}
	return base
}

func gateFirstLineWith(text, needle string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, needle) {
			return strings.TrimSpace(line)
		}
	}
	return ""
}

// Row 28 (O9, O20), gate line 1 and guarantee 5.
//
// Shows: a contribution that keeps the value its rule returned and writes
// into it after the rule returned changes neither the generated Go text nor
// the serialized plan: both are what the same composition gives without the
// mutation. So the projector's view and the published payload come from the
// value that the codec froze when the rule returned.
// Does not show: immutability against a codec that keeps the value itself.
// The codec is trusted to return bytes (O37), and the kernel copies them. It
// does not show the same for a declaration's projector: the mutated value is
// the record contribution's.
// Red under: G5-projector-gets-returned-value, G5-record-encodes-late.
func TestGuarantee5FrozenValueSurvivesMutation(t *testing.T) {
	from, input := gateSource(t, gatePagination)
	clean, err := gatePlanGo(t, input, gateGolang(t, from), gateRequest{})
	if err != nil {
		t.Fatal(err)
	}
	wantText := gateText[gs.Language](t, golangText(), clean)
	wantPlan, err := clean.Representation().Serialize()
	if err != nil {
		t.Fatal(err)
	}
	b := gateGolang(t, from)
	record := b.find(t, "horos-to-golang/record")
	plan, err := gatePlanGo(t, input, b.replace(gateRetaining(record)), gateRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if text := gateText[gs.Language](t, golangText(), plan); text != wantText {
		t.Errorf("the native type changed with the mutation: it now holds %q", gateFirstLineWith(text, "Mutated"))
	}
	serialized, err := plan.Representation().Serialize()
	if err != nil {
		t.Fatal(err)
	}
	if string(serialized) != string(wantPlan) {
		t.Errorf("the serialized plan changed with the mutation: it now holds %q", gateFirstLineWith(string(serialized), `"Mutated`))
	}
	if !strings.Contains(wantText, "Items []T") || strings.Contains(wantText, "Mutated") {
		t.Fatalf("the unmutated text is not the expected one:\n%s", wantText)
	}
}

// gateConfigured wraps the Go record contribution as a fixture whose member
// names come from its own configuration and not from the source.
func gateConfigured(base goContribution, rename map[string]string) goContribution {
	base.Reference.Identity = g.ContributionIdentity{Publisher: gatePublisher, Key: "configured-record"}
	decide := base.Rule.Decide
	base.Rule.Decide = func(ctx context.Context, step horosDecideStep) (g.Decision, error) {
		d, err := decide(ctx, step)
		if payload, ok := d.Payload.(gorecordvocabulary.Payload); ok {
			renamed := slices.Clone(payload.Members)
			for i, name := range renamed {
				if to, ok := rename[name]; ok {
					renamed[i] = to
				}
			}
			d.Payload = gorecordvocabulary.Payload{Members: renamed}
		}
		return d, err
	}
	return base
}

// Rows 6 and 39 (the third layer), gate line 1.
//
// Shows: with a fixture whose member names come from its own configuration,
// the plan entry of each record (the member list of its node, and the member
// names of its frozen payload with the name stated in the request of the
// delegated field) and the generated Go struct have the same members in the
// same order. A pass of this test that rebuilds the entry of Page's body from
// the source, with the naming rule that the builtin contribution uses, gives
// other members than the native type has. So a second origin of the decision
// would be visible here, where it is invisible on the builtin contributions.
// Does not show: that no such pass can be written. The structure prevents a
// sanctioned one (TestGateProjectorCannotReadTheSource). This test detects an
// incorrect one.
// Red under: L1-entry-from-source.
func TestGateDecisionTheSourceDoesNotDetermine(t *testing.T) {
	from, input := gateSource(t, gatePagination)
	b := gateGolang(t, from)
	fixture := gateConfigured(b.find(t, "horos-to-golang/record"), map[string]string{"Items": "Elements", "Name": "FullName"})
	plan, err := gatePlanGo(t, input, b.replace(fixture), gateRequest{})
	if err != nil {
		t.Fatal(err)
	}
	native := gateParseGo(t, gateText[gs.Language](t, golangText(), plan))
	nativeNames := func(declaration string) []string {
		var names []string
		for _, field := range native.fields[declaration] {
			names = append(names, field.name)
		}
		return names
	}
	codec := gateMust[g.PayloadCodec](t)(golangBundleTable()["record"].Codec())
	// The members of an entry as the frozen record states them: for each link
	// in order, the next name of the payload for a field the record lays out,
	// and the name stated in the request for a delegated field.
	recorded := func(node g.PlanNode) []string {
		frozen, err := hex.DecodeString(node.Payload)
		if err != nil {
			t.Fatal(err)
		}
		view, err := codec.View(frozen)
		if err != nil {
			t.Fatal(err)
		}
		own := view.(gorecordvocabulary.Payload).Members
		var members []string
		for _, link := range node.Links {
			if link.Kind == "field" {
				members = append(members, link.Stated)
			} else if len(own) > 0 {
				members, own = append(members, own[0]), own[1:]
			}
		}
		return append(members, own...)
	}
	for declaration, want := range map[string][]string{"Page": {"Elements", "Next"}, "Person": {"FullName", "Active", "Email"}, "Pair": {"First", "Second"}} {
		node := gateBody(t, plan.Representation(), declaration)
		if node.Contribution == nil || *node.Contribution != fixture.Reference {
			t.Fatalf("the body of %s is not selected with the fixture: %+v", declaration, node.Contribution)
		}
		if got := nativeNames(declaration); !slices.Equal(got, want) {
			t.Errorf("%s: the native type has the members %v, the configuration gives %v", declaration, got, want)
		}
		if !slices.Equal(node.Members, nativeNames(declaration)) {
			t.Errorf("%s: the entry has the members %v and the native type has %v", declaration, node.Members, nativeNames(declaration))
		}
		if got := recorded(node); !slices.Equal(got, nativeNames(declaration)) {
			t.Errorf("%s: the frozen record gives the members %v and the native type has %v", declaration, got, nativeNames(declaration))
		}
	}
	// The pass that rebuilds the entry of Page's body from the source: the
	// field names of the module, through the naming rule of the builtin
	// contribution. It is a pass of the test and no part of the kernel.
	var fromSource []string
	for _, field := range []string{"items", "next"} {
		name, err := naming.Exported(field)
		if err != nil {
			t.Fatal(err)
		}
		fromSource = append(fromSource, name)
	}
	if slices.Equal(fromSource, nativeNames("Page")) {
		t.Errorf("an entry rebuilt from the source equals the native type %v: the fixture's decision is determined by the source", fromSource)
	}
	if got := fmt.Sprint(fromSource); got != "[Items Next]" {
		t.Errorf("the pass over the source gives %s", got)
	}
}
