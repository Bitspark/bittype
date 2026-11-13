package integration_test

// PROTOTYPE (bittype#196), stage 5: facets, the Go handlers of the lowering
// facets, and executed Go lowering. Guarantees 1, 2, 3 and 6, gate line 5,
// and the first test of gate line 3.

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/completeness/golang/backend/go"
	"github.com/Bitspark/bittype/completeness/golang/values/go"
	gobooleanfacets "github.com/Bitspark/bittype/completeness/horos-to-golang/facets/boolean/go"
	godeclarationsfacets "github.com/Bitspark/bittype/completeness/horos-to-golang/facets/declarations/go"
	golistfacets "github.com/Bitspark/bittype/completeness/horos-to-golang/facets/list/go"
	gooptionalfacets "github.com/Bitspark/bittype/completeness/horos-to-golang/facets/optional/go"
	gorecordfacets "github.com/Bitspark/bittype/completeness/horos-to-golang/facets/record/go"
	gotextfacets "github.com/Bitspark/bittype/completeness/horos-to-golang/facets/text/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	naming "github.com/Bitspark/bittype/transform/syntax/horos/mapping/naming/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/transform/syntax/shared/payload/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var conversion = protoRequest{products: []string{values.Product}}

// goFacets are the facet packages of bittype's own bundles for the Go
// conversion product. The test composition registers each apart from its
// mapping (O16).
func goFacets(t *testing.T) []g.Facet {
	t.Helper()
	var out []g.Facet
	for _, facets := range []func() ([]g.Facet, error){gotextfacets.Facets, gobooleanfacets.Facets, golistfacets.Facets, gorecordfacets.Facets, godeclarationsfacets.Facets, gooptionalfacets.Facets} {
		out = append(out, protoMust[[]g.Facet](t)(facets())...)
	}
	return out
}

func goCodecs(t *testing.T) []g.PayloadCodec {
	t.Helper()
	var out []g.PayloadCodec
	for _, c := range vocabularyCases() {
		if c.native == "golang" {
			out = append(out, protoMust[g.PayloadCodec](t)(c.codec()))
		}
	}
	return out
}

// conversionBundle is the Go bundle with its facet packages and the product.
func conversionBundle(t *testing.T, from representation.Language) bundle[gs.Language] {
	t.Helper()
	b := golangBundle(t, from)
	b.facets, b.products = goFacets(t), []g.Product{values.Conversion()}
	return b
}

// replaceFacets puts other facets in the place of those of one contribution.
func (b bundle[T]) replaceFacets(of g.ContributionIdentity, with ...g.Facet) bundle[T] {
	out := b
	out.facets = nil
	for _, f := range b.facets {
		if f.Contribution() != of {
			out.facets = append(out.facets, f)
		}
	}
	out.facets = append(out.facets, with...)
	return out
}

// stringly is a record lowering that takes each member by a string: the name
// the naming rule gives the source field. It is the lowering that row 29 rules out.
func stringly(ctx values.Context) (values.Code, error) {
	code, err := gorecordfacets.Lower(ctx)
	if err != nil {
		return code, err
	}
	for i, hole := range code.Holes {
		if hole.Want != g.MemberReference {
			continue
		}
		link := (i - 1) / 2
		canonical, err := ctx.Canonical(link)
		if err != nil {
			return code, err
		}
		spelled, err := naming.Exported(canonical)
		if err != nil {
			return code, err
		}
		code.Text = strings.ReplaceAll(code.Text, fmt.Sprintf("«%d»", i), spelled)
	}
	return code, nil
}

var configuredReference = g.ContributionIdentity{Publisher: "example.org/fixtures", Key: "configured-record"}

// lowerModes are further facets that a lowering process registers, by mode.
// A test file that brings a contribution of its own adds its mode here, so
// that it changes no other file.
var lowerModes = map[string]func(*testing.T) ([]g.Facet, []g.PayloadCodec){}

// lowerCodecs are the codecs a lowering process registers, by mode. A reader
// reads a payload only through the codec of its schema, so it needs the codec
// of every selected contribution.
func lowerCodecs(t *testing.T, mode string) []g.PayloadCodec {
	codecs := goCodecs(t)
	if more, ok := lowerModes[mode]; ok {
		_, added := more(t)
		codecs = append(codecs, added...)
	}
	return codecs
}

// lowerFacets are the facets a lowering process registers, by mode.
func lowerFacets(t *testing.T, mode string) []g.Facet {
	facets := goFacets(t)
	if more, ok := lowerModes[mode]; ok {
		added, _ := more(t)
		return append(facets, added...)
	}
	lower := gorecordfacets.Lower
	if mode == "configured-stringly" {
		lower = stringly
	}
	if strings.HasPrefix(mode, "configured") {
		// The fixture record of stage 1, whose member names come from its
		// configuration. Its facet package is the fixture's own.
		f, err := values.ValueLowering.Issue(configuredReference.Publisher, configuredReference, values.Target, values.Lowering{Lower: lower})
		if err != nil {
			t.Fatal(err)
		}
		facets = append(facets, f)
	}
	return facets
}

// TestPrototypeHelperLower is no test. It is the process that lowers from the
// serialized plan: started by the tests below with a directory, it reads the
// plan and the snapshot from files and writes the lowered Go. It never held
// an object of the planner.
func TestPrototypeHelperLower(t *testing.T) {
	dir := os.Getenv("BITTYPE_PROTOTYPE_LOWER_DIR")
	if dir == "" {
		t.Skip("the lowering process of the prototype's tests")
	}
	plan, err := os.ReadFile(filepath.Join(dir, "plan.json"))
	if err != nil {
		t.Fatal(err)
	}
	snapshot, _ := os.ReadFile(filepath.Join(dir, "snapshot.json"))
	lowered, err := backend.Lower(backend.Input{Plan: plan, Snapshot: snapshot, Facets: lowerFacets(t, os.Getenv("BITTYPE_PROTOTYPE_LOWER_MODE")), Codecs: lowerCodecs(t, os.Getenv("BITTYPE_PROTOTYPE_LOWER_MODE")), Semantics: []string{source.Semantics}})
	if err != nil {
		if werr := os.WriteFile(filepath.Join(dir, "lower.err"), []byte(err.Error()), 0o644); werr != nil {
			t.Fatal(werr)
		}
		return
	}
	if err := os.WriteFile(filepath.Join(dir, "lowered.go"), []byte(lowered), 0o644); err != nil {
		t.Fatal(err)
	}
}

// conversionProgram plans a module with the conversion product, serializes
// the plan, exports the snapshot, has a fresh process lower the conversion
// from those two files, and writes the Go program: the generated types, the
// lowered conversion and the runtime. It returns the program's directory.
func conversionProgram(t *testing.T, module string, b bundle[gs.Language], mode string, snapshot []byte) (string, error) {
	t.Helper()
	_, input := protoSource(t, module)
	plan, err := protoPlanWith[gs.Language](t, input, nativeChecker(t, golangText()), ggtarget.Provider{Package: "main"}, b, conversion)
	if err != nil {
		return "", err
	}
	dir, err := os.MkdirTemp(os.Getenv("GOTMPDIR"), "bittype-prototype-conversion-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	serialized, err := plan.Plan.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	write := func(name string, content []byte) {
		if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("plan.json", serialized)
	if snapshot != nil {
		write("snapshot.json", snapshot)
	}
	write("types.go", []byte(protoText[gs.Language](t, golangText(), plan)))
	write("runtime.go", []byte(backend.Runtime))
	write("go.mod", []byte("module completeness\n\ngo 1.27.0\n"))
	child := exec.Command(os.Args[0], "-test.run=^TestPrototypeHelperLower$", "-test.count=1")
	child.Env = append(os.Environ(), "BITTYPE_PROTOTYPE_LOWER_DIR="+dir, "BITTYPE_PROTOTYPE_LOWER_MODE="+mode)
	if out, err := child.CombinedOutput(); err != nil {
		t.Fatalf("the lowering process failed: %v\n%s", err, out)
	}
	if failure, err := os.ReadFile(filepath.Join(dir, "lower.err")); err == nil {
		return dir, errors.New(string(failure))
	}
	return dir, nil
}

type conversionVector struct {
	Root  string `json:"root"`
	Value any    `json:"value"`
	Error string `json:"error,omitempty"`
}

// runConversion compiles and runs a conversion program over vectors.
func runConversion(t *testing.T, dir string, vectors []conversionVector) ([]conversionVector, error) {
	t.Helper()
	input, err := json.Marshal(vectors)
	if err != nil {
		t.Fatal(err)
	}
	command := exec.Command("go", "run", ".")
	command.Dir, command.Stdin = dir, bytes.NewReader(input)
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil {
		return nil, fmt.Errorf("%v: %s", err, stderr.String())
	}
	var results []conversionVector
	if err := json.Unmarshal(stdout.Bytes(), &results); err != nil {
		t.Fatal(err, stdout.String())
	}
	return results, nil
}

// Gate line 5, and the first test of gate line 3: "Construction of
// `Page[Person]`, access to a member, an owned helper and a generic adapter
// argument compile and run", and "The completeness test runs from the
// serialized plan in a process that never held the planner's objects".
//
// Shows: a second process, given the serialized plan, the source snapshot and
// the facets of the selected contributions, lowers Go conversion code for
// Page[Person] that compiles and converts values there and back. The lowered
// code constructs Page[B0] with a typed construction, reads its members,
// uses the helper the optional contribution owns, and passes the adapter of
// the binder as a generic adapter argument. A vector outside the type is
// refused by the lowered code.
// Does not show: a typed expression vocabulary. The Go vocabulary of the
// prototype is Go text with typed holes (finding F14). It does not show the
// conversion laws: one seat wrote the facets and the vectors.
func TestPrototypeStage5LoweredGoCompilesAndRuns(t *testing.T) {
	if testing.Short() {
		t.Skip("compiles and runs Go")
	}
	from, _ := protoSource(t, protoPagination)
	dir, err := conversionProgram(t, protoPagination, conversionBundle(t, from), "", protoSnapshot(t, protoPagination))
	if err != nil {
		t.Fatal(err)
	}
	lowered, err := os.ReadFile(filepath.Join(dir, "lowered.go"))
	if err != nil {
		t.Fatal(err)
	}
	for what, text := range map[string]string{
		"a typed construction of the declared type": "return Page[B0]{Items: x0, Next: x1, }",
		"access to a member":                        "n.Items",
		"an owned helper":                           "FieldState[BittypeOptional[T]]",
		"a generic adapter argument":                "func adaptPage[B0 any](b0 Adapter[B0]) Adapter[Page[B0]]",
		"the application Page[Person]":              "adaptPage(adaptPerson())",
	} {
		if !strings.Contains(string(lowered), text) {
			t.Errorf("the lowered code has no %s (%q):\n%s", what, text, lowered)
		}
	}
	ada := map[string]any{"name": "Ada", "email": "ada@example.org"}
	bob := map[string]any{"name": "Bob"}
	vectors := []conversionVector{
		{Root: "PeoplePage", Value: map[string]any{"items": []any{ada, bob}, "next": "cursor-2"}},
		{Root: "PeoplePage", Value: map[string]any{"items": []any{}}},
		{Root: "NestedPage", Value: map[string]any{"items": []any{map[string]any{"items": []any{ada}}}}},
		{Root: "Names", Value: []any{"a", "b"}},
		{Root: "Person", Value: bob},
		// Outside the type: a required field is absent, a field is unknown, a
		// value has another structure.
		{Root: "Person", Value: map[string]any{"email": "x"}, Error: "the required field \"name\" is absent"},
		{Root: "Person", Value: map[string]any{"name": "Ada", "age": "36"}, Error: "the record has no field \"age\""},
		{Root: "Names", Value: []any{"a", true}, Error: "text is expected"},
	}
	results, err := runConversion(t, dir, vectors)
	if err != nil {
		t.Fatalf("the lowered code does not compile or run: %v\n%s", err, lowered)
	}
	for i, v := range vectors {
		r := results[i]
		switch {
		case v.Error != "" && !strings.Contains(r.Error, v.Error):
			t.Errorf("vector %d of %s: got %q, want the refusal %q", i, v.Root, r.Error, v.Error)
		case v.Error == "" && (r.Error != "" || !reflect.DeepEqual(r.Value, v.Value)):
			t.Errorf("vector %d of %s: %v came back as %v (%s)", i, v.Root, v.Value, r.Value, r.Error)
		default:
			t.Logf("GREEN vector %d of %s: %v", i, v.Root, r)
		}
	}
	if evidence := os.Getenv("BITTYPE_PROTOTYPE_EVIDENCE"); evidence != "" {
		os.WriteFile(filepath.Join(evidence, "stage5-lowered.go.txt"), lowered, 0o644)
		plan, _ := os.ReadFile(filepath.Join(dir, "plan.json"))
		os.WriteFile(filepath.Join(evidence, "stage5-plan.json"), plan, 0o644)
		os.WriteFile(filepath.Join(evidence, "stage5-snapshot.json"), protoSnapshot(t, protoPagination), 0o644)
	}
}

// Row 29: "A lowering that takes a member by a string the naming policy
// changed fails to compile, and with a member reference it compiles".
//
// Shows: with a record contribution whose member names come from its
// configuration, the lowering that refers to members compiles and runs, and
// the same lowering with each member spelled from the source's field name
// does not compile.
// Does not show: a naming policy of the naming design. The change of names is
// the fixture's configuration.
func TestPrototypeStage5MemberByReferenceNotByString(t *testing.T) {
	if testing.Short() {
		t.Skip("compiles and runs Go")
	}
	const module = "module example.members;\n\ntype Person = { name: text; email?: text; };\n"
	from, _ := protoSource(t, module)
	b := conversionBundle(t, from)
	fixture := configured(b.find(t, "horos-to-golang/record"), map[string]string{"Name": "FullName"})
	// The fixture takes the place of the record contribution, so that
	// contribution's facet package is not registered either.
	b = b.replaceFacets(b.find(t, "horos-to-golang/record").Reference.Identity)
	facet := func(lower func(values.Context) (values.Code, error)) g.Facet {
		return protoMust[g.Facet](t)(values.ValueLowering.Issue(configuredReference.Publisher, configuredReference, values.Target, values.Lowering{Lower: lower}))
	}
	vectors := []conversionVector{{Root: "Person", Value: map[string]any{"name": "Ada"}}}
	byReference, err := conversionProgram(t, module, b.replace(fixture).replaceFacets(configuredReference, facet(gorecordfacets.Lower)), "configured", protoSnapshot(t, module))
	if err != nil {
		t.Fatal(err)
	}
	if results, err := runConversion(t, byReference, vectors); err != nil || results[0].Error != "" || !reflect.DeepEqual(results[0].Value, vectors[0].Value) {
		t.Errorf("the lowering by member reference does not compile and run: %v %v", err, results)
	} else {
		t.Log("GREEN with a member reference the lowering compiles and runs: the member is FullName")
	}
	byString, err := conversionProgram(t, module, b.replace(fixture).replaceFacets(configuredReference, facet(stringly)), "configured-stringly", protoSnapshot(t, module))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := runConversion(t, byString, vectors); err == nil || !strings.Contains(err.Error(), "Name") {
		t.Errorf("the lowering by string compiles: %v", err)
	} else {
		t.Logf("RED   a lowering that takes the member by a string does not compile: %s", firstLineWith(err.Error(), "Name"))
	}
}

// Guarantee 1: "Remove a facet that a bindings product requires. The request
// for bindings fails as "missing facet", and a request for types succeeds."
// Guarantee 2: "A facet from another issuer under the same display name, for
// another target, or of another contract version is rejected."
//
// Shows: a product states the facet kinds it requires, a request names its
// products, and the planner checks the facets of each selected contribution
// against the kind the product requires: its issuer, its target and its
// contract version. A contribution at a field position is asked for another
// kind than a use, so a contribution without a value adapter is expressible.
// Each of the three comparisons is red under a kernel without it.
// Does not show: what a facet does. The kernel hands the value to the owner
// of its kind and reads nothing of it.
func TestPrototypeStage5FacetsBelongAndMatch(t *testing.T) {
	from, input := protoSource(t, protoPagination)
	b := conversionBundle(t, from)
	list := b.find(t, "horos-to-golang/list").Reference.Identity
	if _, err := protoPlanGo(t, input, b, conversion); err != nil {
		t.Fatalf("the whole bundle does not plan for bindings: %v", err)
	}
	_, err := protoPlanGo(t, input, b.replaceFacets(list), conversion)
	refusal := expectRefusal(t, "guarantee 1: a request for bindings without the list contribution's facet", err, "missing facet")
	if refusal.Contribution.Identity != list {
		t.Errorf("the refusal names %s", refusal.Contribution)
	}
	if _, err := protoPlanGo(t, input, b.replaceFacets(list), protoRequest{}); err != nil {
		t.Errorf("guarantee 1: a request for types fails without a facet nobody asked for: %v", err)
	} else {
		t.Log("GREEN guarantee 1: a request for types alone succeeds without the facet")
	}
	// The optional contribution has a field-state facet and no value facet.
	optional := b.find(t, "horos-to-golang/optional").Reference.Identity
	for _, f := range b.facets {
		if f.Contribution() == optional && f.Kind() != values.FieldStateLowering.Reference() {
			t.Errorf("the optional contribution has a facet of kind %v", f.Kind())
		}
	}
	lowering := values.Lowering{Lower: func(values.Context) (values.Code, error) { return values.Code{}, nil }}
	sameName := func(version string) g.FacetKind {
		return protoMust[g.FacetKind](t)(g.NewFacetKind(g.FacetKindDescription{Identity: values.ValueLowering.Reference().Identity, Version: version}))
	}
	mismatched := func(facet g.Facet, err error) func() error {
		return func() error {
			if err != nil {
				return err
			}
			_, err := protoPlanGo(t, input, b.replaceFacets(list, facet), conversion)
			if !errors.Is(err, g.ErrFacet) {
				return fmt.Errorf("the facet was not rejected: %v", err)
			}
			return nil
		}
	}
	redGreen(t, "guarantee 2: a facet from another issuer under the same display name is rejected",
		mismatched(sameName("1").Issue(list.Publisher, list, values.Target, lowering)), "facet-ignores-issuer")
	redGreen(t, "guarantee 2: a facet for another target is rejected",
		mismatched(values.ValueLowering.Issue(list.Publisher, list, "typescript", lowering)), "facet-ignores-target")
	// The kind of the other contract version has the issuer of the required
	// kind in the mutant only when both comparisons are off; the version alone
	// is exercised with a kind that differs in its version and is compared by
	// a kernel that ignores the issuer.
	g.PrototypeMutants["facet-ignores-issuer"] = true
	redGreen(t, "guarantee 2: a facet of another contract version is rejected",
		mismatched(sameName("2").Issue(list.Publisher, list, values.Target, lowering)), "facet-ignores-version")
	delete(g.PrototypeMutants, "facet-ignores-issuer")
	// A facet is admitted only for a registered contribution of its publisher.
	stranger := protoMust[g.Facet](t)(values.ValueLowering.Issue("example.org/a-stranger", list, values.Target, lowering))
	if _, err := protoPlanGo(t, input, b.replaceFacets(list, stranger), conversion); !errors.Is(err, g.ErrScope) {
		t.Errorf("a facet of another publisher than its contribution was admitted: %v", err)
	} else {
		t.Logf("GREEN a facet whose publisher is not its contribution's is not admitted: %v", err)
	}
}

// Guarantee 3: "A reference of another plan, a member where a type is
// expected, and a binder outside its scope are each rejected."
//
// Shows: a reference is issued by one plan with a category and, for a binder,
// a scope, and the plan rejects each of the three misuses. Each is red under
// a kernel without that comparison. The Go backend checks every reference of
// a lowering this way before it resolves it.
// Does not show: references to an external symbol or to a requirement and
// its evidence. The prototype issues the categories type, binder, adapter,
// declaration, owned helper and member.
func TestPrototypeStage5ReferencesAreChecked(t *testing.T) {
	from, input := protoSource(t, protoPagination)
	planned := func() g.RepresentationPlan {
		plan, err := protoPlanGo(t, input, conversionBundle(t, from), conversion)
		if err != nil {
			t.Fatal(err)
		}
		return plan.Plan
	}
	one, other := planned(), planned()
	page, person := symbolNode(t, one, "Page"), symbolNode(t, one, "Person")
	record := bodyNode(t, one, "Page")
	member := protoMust[g.PlanReference](t)(one.MemberOf(record.ID, 0))
	binder := protoMust[g.PlanReference](t)(one.AdapterOf(page.ID, 0))
	if name, err := one.Member(member); err != nil || name != "Items" {
		t.Fatalf("the member reference resolves to %q: %v", name, err)
	}
	rejected := func(err error) error {
		if !errors.Is(err, g.ErrScope) {
			return fmt.Errorf("not rejected: %v", err)
		}
		return nil
	}
	redGreen(t, "guarantee 3: a reference of another plan is rejected", func() error { return rejected(other.Check(member, g.MemberReference, "")) }, "reference-ignores-plan")
	redGreen(t, "guarantee 3: a member where a type is expected is rejected", func() error { return rejected(one.Check(member, g.TypeReference, "")) }, "reference-ignores-category")
	redGreen(t, "guarantee 3: a binder outside its scope is rejected", func() error { return rejected(one.Check(binder, g.AdapterReference, person.Symbol)) }, "reference-ignores-scope")
	if err := one.Check(binder, g.AdapterReference, page.Symbol); err != nil {
		t.Errorf("a binder in its scope is rejected: %v", err)
	}
}

type modedPayload struct{ Mode string }

// moded is a fixture family whose schema has a constraint that a Go type
// cannot express: its one field is one of two atoms. It represents a list.
func moded(t *testing.T, from representation.Language, mode string, projected *bool) (golangContribution, g.PayloadCodec) {
	t.Helper()
	schema, codec := fixtureCodec(t, "moded_list", "module fixtures.moded_list_payload;\n\ntype Payload = { mode: atom-exact(#\"direct\") | atom-exact(#\"boxed\"); };\n",
		func(p modedPayload) (payload.Value, error) {
			return payload.Record(payload.Field{Name: "mode", Value: payload.Atom(p.Mode)})
		},
		func(v payload.Value) (modedPayload, error) {
			fields, err := payload.ReadRecord(v)
			if err != nil {
				return modedPayload{}, err
			}
			mode, err := payload.ReadAtom(fields["mode"])
			return modedPayload{mode}, err
		})
	selector, err := source.Builtin(from, "list-of", "")
	if err != nil {
		t.Fatal(err)
	}
	return golangContribution{
		Reference: fixtureReference("moded-list", "golang/1"), Target: golangTarget, Profile: rules.RepresentationProfile, Domain: g.ForForm(selector), Schema: schema,
		Rule: horosRule{
			Name: func(_ context.Context, step horosName) ([]horosRequest, error) {
				element, err := step.Reader.Element(step.Subject.Use)
				return []horosRequest{g.Reference(element)}, err
			},
			Decide: func(context.Context, horosDecide) (g.Decision, error) {
				return g.Decision{Payload: modedPayload{Mode: mode}}, nil
			},
		},
		Project: func(_ context.Context, p g.Projection[gs.Language]) (g.Projected[gs.Language], error) {
			*projected = true
			return g.Projected[gs.Language]{Type: model.Collection(p.Dependency(0).Type)}, nil
		},
	}, codec
}

// Guarantee 6: "A fixture family whose schema has a constraint that a Go type
// cannot express, and a rule that returns a value violating it: the value is
// refused before any backend runs. A receipt issued for another snapshot
// fails the same way."
//
// Shows: the value a rule returns is validated against its schema when the
// rule returns, before the projector and before any reader: a value that the
// Go type admits and the schema does not is refused, and the projector never
// runs. Red under a codec that does not validate. A receipt is the kernel's
// record over one frozen value and one schema; presented for another frozen
// value or another schema it fails, and a parsed plan with a payload that its
// receipt does not cover is refused.
// Does not show: that the kernel can tell a codec that validates from one
// that does not. The codec is found by the schema reference that the
// contribution states and is trusted to be that schema's codec.
func TestPrototypeStage5ExactPayloadPassedValidation(t *testing.T) {
	from, input := protoSource(t, "module example.moded;\n\ntype Names = list<text>;\n")
	plan := func(mode string) (bool, g.ContributionPlan, error) {
		projected := false
		contribution, codec := moded(t, from, mode, &projected)
		p, err := protoPlanGo(t, input, golangBundle(t, from).replace(contribution).withCodecs(codec), protoRequest{})
		return projected, p, err
	}
	projected, valid, err := plan("boxed")
	if err != nil || !projected {
		t.Fatalf("a value inside the schema does not plan: %v", err)
	}
	redGreen(t, "guarantee 6: a value that violates a constraint of its schema is refused before any backend runs", func() error {
		projected, _, err := plan("sideways")
		if !errors.Is(err, g.ErrPayload) || projected {
			return fmt.Errorf("planning went on (projector ran: %v): %v", projected, err)
		}
		return nil
	}, "codec-skips-validation")
	// A receipt issued for another snapshot.
	list, text := bodyNode(t, valid.Plan, "Names"), g.PlanNode{}
	text, _ = valid.Plan.Node(list.Links[0].Node)
	frozenList, _ := hex.DecodeString(list.Payload)
	frozenText, _ := hex.DecodeString(text.Payload)
	if err := g.CheckReceipt(list.Receipt, *list.Schema, frozenList); err != nil {
		t.Errorf("a receipt fails for its own frozen value: %v", err)
	}
	for name, err := range map[string]error{
		"another frozen value": g.CheckReceipt(list.Receipt, *list.Schema, frozenText),
		"another schema":       g.CheckReceipt(list.Receipt, *text.Schema, frozenList),
	} {
		if !errors.Is(err, g.ErrPayload) {
			t.Errorf("a receipt issued for %s is accepted: %v", name, err)
		} else {
			t.Logf("GREEN guarantee 6: a receipt presented for %s fails: %v", name, err)
		}
	}
	serialized, err := valid.Plan.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	tampered := bytes.Replace(serialized, []byte(`"payload": "`+list.Payload+`"`), []byte(`"payload": "`+text.Payload+`"`), 1)
	if _, err := g.ParsePlan(tampered); err == nil {
		t.Error("a plan whose payload is not the one its receipt covers is read")
	}
}
