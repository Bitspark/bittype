package integration_test

// PROTOTYPE (bittype#196), stage 9: gate line 7 in its small form. A
// materially different contribution is added after the eight stages, one
// registered contribution is replaced by it for the same domain, and nothing
// outside this file changes: not a shared package, not an existing
// contribution, not the composition file.
//
// The contribution: another representation of an optional record field for
// Go. Where bittype's optional contribution wraps the operand in a helper
// declaration with a presence flag, this one represents the field as a slice
// that holds no value or one. It owns no declaration, its native type is
// another, its invariant is another (a slice of at most one value), and its
// conversion refuses a native value that bittype's cannot even express. It is
// one of the three candidates of the advice (section 4.2): an optional
// fragment of another shape. It is written against the public packages and
// the Go value vocabulary of the test tree only, under its own publisher.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/completeness/golang/values/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var slicePresenceReference = g.ContributionReference{
	Identity:       g.ContributionIdentity{Publisher: "example.org/independent-family", Key: "slice-presence"},
	Meaning:        g.Meaning{Identity: "bittype.structural/optional-field", Revision: "1"},
	Representation: "golang-slice/1",
}

// slicePresence is the contribution, with the codec of its payload schema.
func slicePresence(t *testing.T) (golangContribution, g.PayloadCodec) {
	t.Helper()
	schema, codec := emptyFixtureCodec(t, "slice_presence")
	return golangContribution{
		Reference: slicePresenceReference, Target: golangTarget, Profile: rules.RepresentationProfile,
		Domain: g.ForRole(source.OptionalFieldRole), Schema: schema,
		Rule: horosRule{
			Name: func(_ context.Context, step horosName) ([]horosRequest, error) {
				return []horosRequest{g.Reference(step.Subject.Field.Operand)}, nil
			},
			Decide: func(context.Context, horosDecide) (g.Decision, error) { return g.Decision{Payload: struct{}{}}, nil },
		},
		Project: func(_ context.Context, p g.Projection[gs.Language]) (g.Projected[gs.Language], error) {
			operand := p.Dependency(0)
			return g.Projected[gs.Language]{Fragment: []g.Expression[gs.Language]{g.At(model.Field(p.Stated(), false, model.Collection(operand.Type)), operand.Origin)}, Members: []g.Member{{Name: p.Stated()}}}, nil
		},
	}, codec
}

// slicePresenceFacet is the contribution's facet for the Go conversion
// product: its field-state protocol.
func slicePresenceFacet(t *testing.T) g.Facet {
	t.Helper()
	return protoMust[g.Facet](t)(values.Issue(values.FieldStateLowering, slicePresenceReference, values.Lowering{
		Helpers: func(values.Context) ([]values.HelperCode, error) {
			return []values.HelperCode{{Name: "state", Code: values.Go(`func %s[T any](operand Adapter[T]) FieldState[[]T] {
	return FieldState[[]T]{
		Construct: func(present bool, value Value) ([]T, error) {
			if !present {
				return nil, nil
			}
			n, err := operand.Decode(value)
			return []T{n}, err
		},
		State: func(member []T) (bool, Value, error) {
			switch len(member) {
			case 0:
				return false, nil, nil
			case 1:
				v, err := operand.Encode(member[0])
				return true, v, err
			}
			return false, nil, fmt.Errorf("invalid native value: an optional field holds %%d values", len(member))
		},
	}
}`, values.Helper("state"))}}, nil
		},
		Lower: func(values.Context) (values.Code, error) {
			return values.Go("%s(%s)", values.Helper("state"), values.Dependency(0)), nil
		},
	}))
}

func init() {
	lowerModes["slice-presence"] = func(t *testing.T) ([]g.Facet, []g.PayloadCodec) {
		_, codec := slicePresence(t)
		return []g.Facet{slicePresenceFacet(t)}, []g.PayloadCodec{codec}
	}
}

// replacedBundle is the consumer's composition with the contribution in the
// place of bittype's optional contribution.
func replacedBundle(t *testing.T, b bundle[gs.Language]) bundle[gs.Language] {
	t.Helper()
	contribution, codec := slicePresence(t)
	replaced := b.find(t, "horos-to-golang/optional").Reference.Identity
	out := b.replace(contribution).withCodecs(codec)
	if len(b.facets) > 0 {
		out = out.replaceFacets(replaced, slicePresenceFacet(t))
	}
	return out
}

// dependsOn reports, for each node of a plan, whether it or a node it links
// to, at any depth, was selected with the given contribution.
func dependsOn(plan g.RepresentationPlan, identity g.ContributionIdentity) map[string]bool {
	nodes := map[string]g.PlanNode{}
	for _, n := range plan.Nodes() {
		nodes[n.ID] = n
	}
	memo := map[string]bool{}
	var visit func(string) bool
	visit = func(id string) bool {
		if v, ok := memo[id]; ok {
			return v
		}
		memo[id] = false
		n := nodes[id]
		depends := n.Contribution != nil && n.Contribution.Identity == identity
		for _, l := range n.Links {
			depends = depends || l.Node != "" && visit(l.Node)
		}
		if a := n.Application; a != nil {
			depends = depends || visit(a.Declaration)
			for _, argument := range append(append([]g.PlanArgument(nil), a.Native...), a.Adapters...) {
				depends = depends || visit(argument.Node)
			}
		}
		memo[id] = depends
		return depends
	}
	for id := range nodes {
		visit(id)
	}
	return memo
}

// Gate line 7, row 1, O20, O26: "One registered contribution is replaced by
// another for the same domain. The payloads of every other selection are
// byte-identical, and so are the frozen records of every selection that does
// not depend on the replaced one."
//
// Shows: a consumer replaces one contribution and keeps every other. No
// payload of another selection changes. Every node that does not depend on
// the replaced contribution is the same node. A node that depends on it is
// another node (O26), with the same payload. With both registered the domain
// is ambiguous until the consumer composes them explicitly.
// Does not show: a mutant. The line is red on main because there is no
// payload, no node and no replacement of a role there.
func TestPrototypeStage9ReplacedContributionLeavesTheOthers(t *testing.T) {
	from, input := protoSource(t, protoPagination)
	b := golangBundle(t, from)
	replaced := b.find(t, "horos-to-golang/optional").Reference.Identity
	before, err := protoPlanGo(t, input, b, protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	after, err := protoPlanGo(t, input, replacedBundle(t, b), protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	type selection struct {
		kind, source, contribution, payload string
	}
	count := func(plan g.RepresentationPlan, skip g.ContributionIdentity) map[selection]int {
		out := map[selection]int{}
		for _, n := range plan.Nodes() {
			if n.Contribution != nil && n.Contribution.Identity != skip && n.Kind != "owned" {
				out[selection{n.Kind, n.Source, n.Contribution.String(), n.Payload}]++
			}
		}
		return out
	}
	if a, b := count(before.Plan, replaced), count(after.Plan, slicePresenceReference.Identity); !reflect.DeepEqual(a, b) {
		t.Errorf("the payloads of the other selections changed:\n%v\n%v", a, b)
	} else {
		t.Logf("GREEN the payloads of the %d other selections are byte-identical", len(a))
	}
	depends := dependsOn(before.Plan, replaced)
	same, other := 0, 0
	for _, n := range before.Plan.Nodes() {
		_, kept := after.Plan.Node(n.ID)
		switch {
		case !depends[n.ID] && !kept && n.Kind != "owned":
			t.Errorf("the frozen record of %s %s changed, and it does not depend on the replaced contribution", n.Kind, n.Source)
		case !depends[n.ID]:
			same++
		case kept:
			t.Errorf("the node of %s %s depends on the replaced contribution and is unchanged", n.Kind, n.Source)
		default:
			other++
		}
	}
	t.Logf("GREEN %d frozen records that do not depend on the replaced contribution are identical; %d that depend on it are other nodes (O26)", same, other)
	// Both registered: ambiguous, until the consumer composes them explicitly.
	contribution, codec := slicePresence(t)
	_, err = protoPlanGo(t, input, b.with(contribution).withCodecs(codec), protoRequest{})
	expectRefusal(t, "two contributions for the role, uncomposed", err, "ambiguous contribution")
	var registry c.Registry
	checker := nativeChecker(t, golangText())
	if err := c.RegisterSource[h.Language, source.Reader](&registry, from, source.Provider{}); err != nil {
		t.Fatal(err)
	}
	if err := c.RegisterTarget[gs.Language](&registry, checker.Language(), ggtarget.Provider{Package: "generated"}); err != nil {
		t.Fatal(err)
	}
	for _, registered := range append(append([]golangContribution(nil), b.contributions...), contribution) {
		if err := c.RegisterContribution(&registry, from, checker.Language(), registered); err != nil {
			t.Fatal(err)
		}
	}
	for _, registered := range append(append([]g.PayloadCodec(nil), b.codecs...), codec) {
		if err := registry.RegisterCodec(registered); err != nil {
			t.Fatal(err)
		}
	}
	if err := registry.ReplaceContribution(replaced, slicePresenceReference.Identity); err != nil {
		t.Fatal(err)
	}
	composed, err := registry.PlanContributions(context.Background(), c.ContributionRequest{Request: c.Request{Input: input, To: checker.Language(), Output: "generated", Roots: protoRoots(t, input, []string{"PeoplePage"}), Profile: rules.RepresentationProfile, Limits: g.DefaultLimits()}}, checker)
	if err != nil {
		t.Fatalf("the explicit composition does not plan: %v", err)
	}
	if _, selected := composed.Plan.Contribution(slicePresenceReference); !selected {
		t.Error("the explicit composition did not select the replacement")
	} else {
		t.Log("GREEN with both registered, the consumer's explicit replacement selects the new contribution")
	}
}

// Gate line 7, row 40: "The contribution registers, and `Page<Person>`
// regenerates with it. The diff of that test touches no shared package".
//
// Shows: with the shared packages fixed, the contribution registers, the
// native declarations of Page<Person> regenerate with its representation,
// and the Go conversion lowered by the unchanged backend, with the unchanged
// facet of the record, compiles and converts through the new contribution's
// field-state protocol. The commit that adds this file adds nothing else.
// Does not show: authorship by another owner; one seat wrote it. It does not
// show the same for TypeScript.
func TestPrototypeStage9PageOfPersonRegeneratesAndConverts(t *testing.T) {
	if testing.Short() {
		t.Skip("compiles and runs Go")
	}
	from, input := protoSource(t, protoPagination)
	b := replacedBundle(t, conversionBundle(t, from))
	plan, err := protoPlanGo(t, input, b, conversion)
	if err != nil {
		t.Fatal(err)
	}
	text := protoText[gs.Language](t, golangText(), plan)
	fields, _, _ := goFieldTypes(t, text)
	if fields["Page"]["Next"] != "[]string" || fields["Person"]["Email"] != "[]string" || strings.Contains(text, "BittypeOptional") {
		t.Fatalf("Page<Person> did not regenerate with the contribution:\n%s", text)
	}
	t.Logf("GREEN Page<Person> regenerates with the contribution:\n%s", text)
	dir, err := conversionProgram(t, protoPagination, b, "slice-presence", protoSnapshot(t, protoPagination))
	if err != nil {
		t.Fatal(err)
	}
	vectors := []conversionVector{
		{Root: "PeoplePage", Value: map[string]any{"items": []any{map[string]any{"name": "Ada", "email": "ada@example.org"}, map[string]any{"name": "Bob"}}, "next": "cursor-2"}},
		{Root: "PeoplePage", Value: map[string]any{"items": []any{}}},
	}
	results, err := runConversion(t, dir, vectors)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range vectors {
		if results[i].Error != "" || !reflect.DeepEqual(results[i].Value, v.Value) {
			t.Errorf("vector %d came back as %+v", i, results[i])
		}
	}
	t.Log("GREEN the conversion lowered with the unchanged backend and record facet converts through the new field-state protocol")
}

const outOfTreeContribution = `// Package outoftree authors a contribution in a Go module outside bittype.
package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/transform/syntax/shared/payload/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

//go:embed payload.ontos
var schemaModule []byte

type Payload struct{}

var reference = g.ContributionReference{
	Identity:       g.ContributionIdentity{Publisher: "example.org/outoftree", Key: "text"},
	Meaning:        g.Meaning{Identity: "bittype.structural/text", Revision: "1"},
	Representation: "golang/outoftree-1",
}
var schema = g.SchemaReference{Identity: "example.org/outoftree/payload/text", Version: "1"}

func main() {
	from, err := representation.NewLanguage[h.Language](h.Language{}, h.Profile)
	check(err)
	to, err := representation.NewLanguage[s.Language](s.Language{}, "native-types-v1")
	check(err)
	selector, err := source.Builtin(from, "scalar", "utf8-text")
	check(err)
	loaded, err := payload.Load(schema, schemaModule, schemaPin, "Payload")
	check(err)
	codec := payload.New(loaded, func(Payload) (payload.Value, error) { return payload.Record() }, func(payload.Value) (Payload, error) { return Payload{}, nil })
	contribution := g.Contribution[h.Language, s.Language, source.Reader]{
		Reference: reference, Target: string(s.Language{}.ID()), Profile: rules.RepresentationProfile, Domain: g.ForForm(selector), Schema: schema,
		Rule: g.Rule[h.Language, source.Reader]{Decide: func(context.Context, g.DecideStep[h.Language, source.Reader]) (g.Decision, error) {
			return g.Decision{Payload: Payload{}}, nil
		}},
		Project: func(context.Context, g.Projection[s.Language]) (g.Projected[s.Language], error) {
			return g.Projected[s.Language]{Type: model.Primitive(model.PrimitiveString)}, nil
		},
	}
	var registry c.Registry
	check(c.RegisterContribution(&registry, from, to, contribution))
	check(registry.RegisterCodec(codec))
	frozen, err := codec.Freeze(Payload{})
	check(err)
	fmt.Printf("registered %s with a codec that freezes %d bytes\n", reference, len(frozen))
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
`

// Gate line 7, row 37: "A Go module outside the repository authors and
// registers a contribution with public packages only".
//
// Shows: a Go module outside this repository, which reaches bittype through
// its module path alone, authors a contribution that states native types,
// loads its payload schema, builds its codec and registers both. It imports
// no internal package and nothing of the test tree. A second file that needs
// an internal package does not compile.
// Does not show: a generation with that contribution; the module registers
// and does not plan. It does not show how the author gets the canonical
// schema module, its pin and the generated payload type: this test writes
// them, with a test helper that uses an internal package (finding F25).
func TestPrototypeStage9OutOfTreeModule(t *testing.T) {
	if testing.Short() {
		t.Skip("builds a Go module")
	}
	dir, err := os.MkdirTemp(os.Getenv("GOTMPDIR"), "bittype-prototype-outoftree-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	module, pin := protoSchemaModule(t, "module outoftree.text_payload;\n\ntype Payload = { };\n")
	sums, err := os.ReadFile(filepath.Join(root, "go.sum"))
	if err != nil {
		t.Fatal(err)
	}
	files := map[string][]byte{
		"go.mod":        []byte("module example.org/outoftree\n\ngo 1.27.0\n\nrequire (\n\tgithub.com/Bitspark/bittype v0.0.0\n\tgithub.com/Bitspark/horos/core/go v0.10.2\n\tgithub.com/Bitspark/logos-contract/contract/go v0.26.0\n\tgithub.com/bitspark/ontos v0.16.0\n)\n\nreplace github.com/Bitspark/bittype => " + filepath.ToSlash(root) + "\n"),
		"go.sum":        sums,
		"main.go":       []byte(outOfTreeContribution),
		"pin.go":        []byte("package main\n\nconst schemaPin = \"" + pin + "\"\n"),
		"payload.ontos": module,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	run := func(args ...string) (string, error) {
		command := exec.Command("go", args...)
		command.Dir = dir
		command.Env = append(os.Environ(), "GOFLAGS=-mod=mod", "GOPROXY=off", "GOSUMDB=off")
		var out bytes.Buffer
		command.Stdout, command.Stderr = &out, &out
		err := command.Run()
		return out.String(), err
	}
	out, err := run("run", ".")
	if err != nil || !strings.Contains(out, "registered example.org/outoftree/text") {
		t.Fatalf("the module outside the repository does not author and register: %v\n%s", err, out)
	}
	t.Logf("GREEN %s", strings.TrimSpace(out))
	imports, err := run("list", "-f", `{{join .Imports "\n"}}`, ".")
	if err != nil {
		t.Fatal(err, imports)
	}
	for _, imported := range strings.Fields(imports) {
		for _, forbidden := range []string{"/internal/", "/integration/", "/completeness/"} {
			if strings.HasPrefix(imported, "github.com/Bitspark/bittype/") && strings.Contains(imported+"/", forbidden) {
				t.Errorf("the module imports %s", imported)
			}
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "internal.go"), []byte("package main\n\nimport _ \"github.com/Bitspark/bittype/internal/horosmodel\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, err := run("build", "."); err == nil || !strings.Contains(out, "internal") {
		t.Errorf("a module that needs an internal package builds: %v\n%s", err, out)
	} else {
		t.Logf("RED   a module that needs an internal symbol does not build: %s", firstLineWith(out, "internal"))
	}
}

// Gate line 7, row 33: "One fixture for each item that is in the extension
// domain, and one recorded refusal for each that is out".
//
// This test runs the one refusal that no earlier stage has, and lists where
// each other item is shown.
//
//	A further family                  in   fixtures evidenceList (stage 2), moded (stage 5); refused "missing contribution" (stage 3)
//	Further products and facets       in   the fixture product and facet kind (stage 8, case 3); refused "missing facet" (stages 5, 8)
//	A contextual representation       in   slicePresence (this stage); refused: a fragment outside the member vocabulary (stage 6)
//	Cross-field layout coordination   out  refused below
//	Application-specific specialization  in, refused until #166 (stage 3)
//	Authoring routes                  one  the module outside the repository (this stage)
//
// Shows: a contribution at a field position cannot reach another field. Its
// subject is one operand, and a request for a sibling's operand is refused.
// So a layout that several fields share cannot be answered field by field.
// Does not show: a record policy that coordinates fields. That is the
// replacement of the record contribution, which the domain names as the way.
func TestPrototypeStage9ExtensionDomainRefusals(t *testing.T) {
	from, input := protoSource(t, "module example.crossfield;\n\ntype Pair = { first?: text; second?: text; };\n")
	b := golangBundle(t, from)
	var sibling *g.Use[h.Language]
	contribution, codec := slicePresence(t)
	contribution.Rule.Name = func(_ context.Context, step horosName) ([]horosRequest, error) {
		if sibling == nil {
			operand := step.Subject.Field.Operand
			sibling = &operand
			return []horosRequest{g.Reference(operand)}, nil
		}
		// A layout shared across fields would need the other field too.
		return []horosRequest{g.Reference(step.Subject.Field.Operand), g.Reference(*sibling)}, nil
	}
	_, err := protoPlanGo(t, input, b.replace(contribution).withCodecs(codec), protoRequest{})
	if !errors.Is(err, g.ErrScope) {
		t.Errorf("a contribution at a field position reached another field: %v", err)
	} else {
		t.Logf("GREEN cross-field layout coordination is out: a request for a sibling's operand is refused: %v", err)
	}
}

var _ = fmt.Sprint
var _ representation.Language
