package integration_test

// PROTOTYPE (bittype#196), stage 1: the staged rule and the frozen record.
// Gate line 1, native projection: "The native type is computed from the frozen
// record, by a projector that cannot read the source".

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
	gorecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/record/go"
	gorecordvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/record/go"
	naming "github.com/Bitspark/bittype/transform/syntax/horos/mapping/naming/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// protoPagination is the pagination fixture without booleans, which the
// TypeScript side of the prototype does not move.
const protoPagination = `module example.pagination;

type Page<T> = { items: list<T>; next?: text; };
type Person = { name: text; email?: text; };
type PeoplePage = Page<Person>;
type NestedPage = Page<Page<Person>>;
type Names = list<text>;
`

// mutant runs a check with one property of the kernel switched off.
func mutant(name string, check func() error) error {
	g.PrototypeMutants[name] = true
	defer delete(g.PrototypeMutants, name)
	return check()
}

// redGreen runs a check as it is, where it must pass, and under each named
// mutant, where it must fail. It logs both, so that a run shows the red and
// the green of one property side by side.
func redGreen(t *testing.T, property string, check func() error, mutants ...string) {
	t.Helper()
	if err := check(); err != nil {
		t.Errorf("%s: %v", property, err)
	} else {
		t.Logf("GREEN %s", property)
	}
	for _, name := range mutants {
		if err := mutant(name, check); err == nil {
			t.Errorf("%s: the mutant %q is not detected", property, name)
		} else {
			t.Logf("RED   %s under mutant %q: %v", property, name, err)
		}
	}
}

// The protocol end to end. The moved contributions of Go and TypeScript run
// through name, resolve, decide, freeze and project, and generate what the
// mapping contract generates today for the same module.
//
// Shows: the staged rule, the codec and the registered projector produce the
// native declarations of today's pagination fixtures unchanged, for Go and
// TypeScript, with no contribution for generics (the two checks that are green
// today and stay green, for the two targets the prototype moved).
// Does not show: Rust and C++, which stay on the mapping contract.
func TestPrototypeStage1NativeDeclarationsDoNotChange(t *testing.T) {
	_, _, full := generationInput(t, paginationModule())
	goChecker := nativeChecker(t, golangText())
	before, err := runGeneration(t, full, goChecker, goRequest)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := protoPlanGo(t, full, golangBundle(t, full.Source.Language()), protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	want := formatted(t, golangText(), outputTree[gs.Language](t, before))
	got := strings.Replace(protoText[gs.Language](t, golangText(), plan), "package generated", "package pagination", 1)
	if got != want {
		t.Errorf("Go declarations changed.\nmapping contract:\n%s\ncontribution contract:\n%s", want, got)
	}
	from, reduced := protoSource(t, protoPagination)
	tsChecker := nativeChecker(t, typescriptText())
	beforeTS, err := runGeneration(t, reduced, tsChecker, tsRequest)
	if err != nil {
		t.Fatal(err)
	}
	planTS, err := protoPlanTS(t, reduced, typescriptBundle(t, from), protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if want, got := formatted(t, typescriptText(), outputTree[ts.Language](t, beforeTS)), protoText[ts.Language](t, typescriptText(), planTS); got != want {
		t.Errorf("TypeScript declarations changed.\nmapping contract:\n%s\ncontribution contract:\n%s", want, got)
	}
	t.Logf("Go, contribution contract:\n%s", got)
}

// sourceReach lists what a callback's capability type can reach of the
// source: a method whose signature mentions a kernel handle of the source, a
// selector, or the source reader.
func sourceReach(capability reflect.Type) []string {
	var reach []string
	mentions := func(t reflect.Type) bool {
		name := t.String()
		for _, forbidden := range []string{"generation.Use[", "generation.UseView[", "generation.Declaration[", "generation.DeclarationView[", "generation.Binder[", "generation.Selector", "source.Reader", "generation.Captured["} {
			if strings.Contains(name, forbidden) {
				return true
			}
		}
		return false
	}
	for i := 0; i < capability.NumMethod(); i++ {
		m := capability.Method(i)
		for j := 1; j < m.Type.NumIn(); j++ {
			if mentions(m.Type.In(j)) {
				reach = append(reach, m.Name)
			}
		}
		for j := 0; j < m.Type.NumOut(); j++ {
			if mentions(m.Type.Out(j)) {
				reach = append(reach, m.Name)
			}
		}
	}
	return reach
}

// Row 4: "A projector that tries to read the source does not compile or fails
// the contract test".
//
// Shows: the capability a projector is given has no method that mentions a
// source handle, a selector or the reader, and Go refuses to compile a
// projector that calls one. The same contract test is red on the callback
// capability of the mapping contract, which is why the line is red on main.
// Does not show: that a contribution cannot pass what its rule read to its
// projector through state of its own package. The page says so ("A limit").
func TestPrototypeStage1ProjectorCannotReadTheSource(t *testing.T) {
	if reach := sourceReach(reflect.TypeFor[g.Projection[gs.Language]]()); len(reach) != 0 {
		t.Errorf("a projector can reach the source through %v", reach)
	} else {
		t.Log("GREEN the projector's capability reaches no source handle")
	}
	if reach := sourceReach(reflect.TypeFor[g.Expansion[h.Language, gs.Language, source.Reader]]()); len(reach) == 0 {
		t.Error("the contract test does not see the reader of the mapping contract's callback")
	} else {
		t.Logf("RED   on the branch point: the one callback of the mapping contract reaches the source through %v", reach)
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
	if err := compile("_ = p.Payload()"); err != nil {
		t.Fatalf("a projector that reads the frozen record does not compile: %v", err)
	}
	for _, attempt := range []string{"_, _ = p.Reader()", "_, _ = p.Use()", "_, _ = p.Declaration()"} {
		if err := compile(attempt); err == nil {
			t.Errorf("a projector that calls %q compiles", attempt)
		} else {
			t.Logf("GREEN a projector with %q does not compile: %v", attempt, err)
		}
	}
}

// retaining wraps the Go record contribution so that it keeps what its rule
// returned and mutates it when its projector starts: after the rule returned.
func retaining(t *testing.T, base golangContribution) golangContribution {
	var retained [][]string
	decide, project := base.Rule.Decide, base.Project
	base.Rule.Decide = func(ctx context.Context, step horosDecide) (g.Decision, error) {
		d, err := decide(ctx, step)
		if payload, ok := d.Payload.(gorecordvocabulary.Payload); ok {
			retained = append(retained, payload.Members)
		}
		return d, err
	}
	base.Project = func(ctx context.Context, p g.Projection[gs.Language]) (g.Projected[gs.Language], error) {
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

// Row 28 and guarantee 5: "Mutating the value a rule returned, after it
// returned, changes neither the native type nor the serialized plan."
//
// Shows: the projector's view and the published payload are decoded from the
// value the codec froze when the rule returned, not from what the rule holds.
// Red under a kernel that hands the projector the returned value, and under a
// kernel that encodes the value only when the plan is serialized.
// Does not show: immutability against a codec that keeps the value itself.
// The codec is trusted to return bytes, and the kernel copies them.
func TestPrototypeStage1FrozenValueSurvivesMutation(t *testing.T) {
	from, input := protoSource(t, protoPagination)
	clean, err := protoPlanGo(t, input, golangBundle(t, from), protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	wantText := protoText[gs.Language](t, golangText(), clean)
	wantPlan, err := clean.Plan.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	check := func() error {
		b := golangBundle(t, from)
		plan, err := protoPlanGo(t, input, b.replace(retaining(t, b.find(t, "horos-to-golang/record"))), protoRequest{})
		if err != nil {
			return err
		}
		if text := protoText[gs.Language](t, golangText(), plan); text != wantText {
			return fmt.Errorf("the native type changed with the mutation: it now holds %q", firstLineWith(text, "Mutated"))
		}
		serialized, err := plan.Plan.Serialize()
		if err != nil {
			return err
		}
		if string(serialized) != string(wantPlan) {
			return fmt.Errorf("the serialized plan changed with the mutation")
		}
		return nil
	}
	redGreen(t, "a mutation after the rule returned changes neither the native type nor the serialized plan", check, "projector-gets-returned-value", "freeze-at-serialization")
}

func firstLineWith(text, needle string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, needle) {
			return strings.TrimSpace(line)
		}
	}
	return ""
}

// configured wraps the Go record contribution as a fixture whose member names
// come from its own configuration and not from the source.
func configured(base golangContribution, rename map[string]string) golangContribution {
	base.Reference.Identity = g.ContributionIdentity{Publisher: "example.org/fixtures", Key: "configured-record"}
	decide := base.Rule.Decide
	base.Rule.Decide = func(ctx context.Context, step horosDecide) (g.Decision, error) {
		d, err := decide(ctx, step)
		if payload, ok := d.Payload.(gorecordvocabulary.Payload); ok {
			for i, name := range payload.Members {
				if to, ok := rename[name]; ok {
					payload.Members[i] = to
				}
			}
			d.Payload = payload
		}
		return d, err
	}
	return base
}

// goFields parses generated Go and returns the field names of each struct.
func goFields(t *testing.T, text string) map[string][]string {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "generated.go", text, parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string][]string{}
	ast.Inspect(file, func(n ast.Node) bool {
		spec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if s, ok := spec.Type.(*ast.StructType); ok {
			out[spec.Name.Name] = []string{}
			for _, f := range s.Fields.List {
				for _, name := range f.Names {
					out[spec.Name.Name] = append(out[spec.Name.Name], name.Name)
				}
			}
		}
		return true
	})
	return out
}

// Rows 6 and 39: "A fixture whose decision the source does not determine is
// green on the planner and red on a mutant that derives the entry from the
// source".
//
// Shows: the plan entry and the native type of one selection come from one
// frozen record, in a case where a second pass over the source would compute
// something else. The mutant here is a pass of the test that rebuilds the
// entry from the source; the check detects it.
// Does not show: that no such pass can be written. The structure prevents a
// sanctioned one (the projector test above); this test detects an incorrect one.
func TestPrototypeStage1DecisionTheSourceDoesNotDetermine(t *testing.T) {
	from, input := protoSource(t, protoPagination)
	b := golangBundle(t, from)
	fixture := configured(b.find(t, "horos-to-golang/record"), map[string]string{"Items": "Elements", "Name": "FullName"})
	plan, err := protoPlanGo(t, input, b.replace(fixture), protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	native := goFields(t, protoText[gs.Language](t, golangText(), plan))
	codec := protoMust[g.PayloadCodec](t)(gorecordcodec.Codec())
	// The members of the entry of Page's body, as the plan records them: the
	// payload names for the fields the record lays out, and the name stated in
	// the request for the delegated field.
	entry := func(fromSource bool) ([]string, error) {
		for _, node := range plan.Plan.Nodes() {
			if node.Kind != "use" || node.Contribution == nil || *node.Contribution != fixture.Reference || len(node.Links) != 2 || node.Links[1].Stated != "Next" {
				continue
			}
			frozen, err := hex.DecodeString(node.Payload)
			if err != nil {
				return nil, err
			}
			view, err := codec.View(frozen)
			if err != nil {
				return nil, err
			}
			members := append([]string(nil), view.(gorecordvocabulary.Payload).Members...)
			if fromSource {
				// The mutant: a separate pass derives the entry from the source.
				name, err := naming.Exported("items")
				if err != nil {
					return nil, err
				}
				members = []string{name}
			}
			return append(members, node.Links[1].Stated), nil
		}
		return nil, fmt.Errorf("the plan has no entry for the body of Page")
	}
	check := func(fromSource bool) error {
		members, err := entry(fromSource)
		if err != nil {
			return err
		}
		if !slices.Equal(members, native["Page"]) {
			return fmt.Errorf("the entry has the members %v and the native type has %v", members, native["Page"])
		}
		return nil
	}
	if err := check(false); err != nil {
		t.Error(err)
	} else {
		t.Logf("GREEN the planner's entry and the native type agree: %v", native["Page"])
	}
	if err := check(true); err == nil {
		t.Error("an entry derived from the source is not detected")
	} else {
		t.Logf("RED   on a mutant that derives the entry from the source: %v", err)
	}
}
