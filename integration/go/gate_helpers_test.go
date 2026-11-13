package integration_test

// The helpers of the kernel-level tests of the merge gate (bittype#196,
// docs/contribution-contract.md, "The gate"). A gate test composes what a
// consumer registers, plans through the public registry and reads the selected
// representation plan and the generated native text. No helper here reads a
// switch of the kernel: a test is shown red by a mutant that a script makes on
// the worktree and removes again.

import (
	"context"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	shared "github.com/Bitspark/bittype/integration/go/testdata/composition"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	cs "github.com/Bitspark/bittype/syntax/cpp/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	rs "github.com/Bitspark/bittype/syntax/rust/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	cgtarget "github.com/Bitspark/bittype/transform/syntax/cpp/target/go"
	gomodel "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	rgtarget "github.com/Bitspark/bittype/transform/syntax/rust/target/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	tgtarget "github.com/Bitspark/bittype/transform/syntax/typescript/target/go"
	astree "github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type goContribution = g.Contribution[h.Language, gs.Language, source.Reader]
type goProjection = g.Projection[gs.Language]
type goProjected = g.Projected[gs.Language]

// goTarget is the target that a contribution written for Go states.
var goTarget = string(gs.Language{}.ID())

// gatePublisher is the publisher of the fixture contributions of the gate.
const gatePublisher = "example.test/gate"

// gateReference is the reference of a fixture contribution of the gate.
func gateReference(key, revision string) g.ContributionReference {
	return g.ContributionReference{
		Identity:       g.ContributionIdentity{Publisher: gatePublisher, Key: key},
		Meaning:        g.Meaning{Identity: gatePublisher + "/" + key, Revision: "1"},
		Representation: revision,
	}
}

// gateBundle is what a consumer's composition registers for one target under
// the horos source: contributions, the codecs of their payload schemas, facets
// and products.
type gateBundle[T any] struct {
	contributions []g.Contribution[h.Language, T, source.Reader]
	codecs        []g.PayloadCodec
	facets        []g.Facet
	products      []g.Product
}

// replace puts a contribution in the place of every contribution with the
// same domain: the composition of a consumer that selects it instead.
func (b gateBundle[T]) replace(with g.Contribution[h.Language, T, source.Reader]) gateBundle[T] {
	out := b
	out.contributions = nil
	for _, old := range b.contributions {
		if old.Domain.Key() != with.Domain.Key() {
			out.contributions = append(out.contributions, old)
		}
	}
	out.contributions = append(out.contributions, with)
	return out
}

// with adds contributions beside the ones the bundle holds.
func (b gateBundle[T]) with(more ...g.Contribution[h.Language, T, source.Reader]) gateBundle[T] {
	out := b
	out.contributions = append(slices.Clone(b.contributions), more...)
	return out
}

// withCodecs adds the codecs of further payload schemas. A registry holds one
// codec for each schema reference, so a codec of a schema that the bundle has
// already is left out.
func (b gateBundle[T]) withCodecs(more ...g.PayloadCodec) gateBundle[T] {
	out := b
	out.codecs = slices.Clone(b.codecs)
	for _, codec := range more {
		if !slices.ContainsFunc(out.codecs, func(old g.PayloadCodec) bool { return old.Schema() == codec.Schema() }) {
			out.codecs = append(out.codecs, codec)
		}
	}
	return out
}

func (b gateBundle[T]) withFacets(more ...g.Facet) gateBundle[T] {
	out := b
	out.facets = append(slices.Clone(b.facets), more...)
	return out
}

// without leaves out the contribution with the local key key.
func (b gateBundle[T]) without(key string) gateBundle[T] {
	out := b
	out.contributions = nil
	for _, old := range b.contributions {
		if old.Reference.Identity.Key != key {
			out.contributions = append(out.contributions, old)
		}
	}
	return out
}

func (b gateBundle[T]) find(t *testing.T, key string) g.Contribution[h.Language, T, source.Reader] {
	t.Helper()
	for _, old := range b.contributions {
		if old.Reference.Identity.Key == key {
			return old
		}
	}
	t.Fatalf("no contribution %s in the bundle", key)
	panic("unreachable")
}

// gateBuiltin is the builtin set of one target as a bundle.
func gateBuiltin[T any](t *testing.T, from representation.Language, table map[string]genfixture.Bundle[T], names []string) gateBundle[T] {
	t.Helper()
	set, err := shared.Build(from, shared.Select(table, names...)...)
	if err != nil {
		t.Fatal(err)
	}
	return gateBundle[T]{contributions: set.Contributions, codecs: set.Codecs}
}

func gateGolang(t *testing.T, from representation.Language) gateBundle[gs.Language] {
	t.Helper()
	return gateBuiltin(t, from, golangBundleTable(), builtinBundleNames)
}

func gateTypescript(t *testing.T, from representation.Language) gateBundle[ts.Language] {
	t.Helper()
	return gateBuiltin(t, from, typescriptBundleTable(), builtinBundleNames)
}

func gateRust(t *testing.T, from representation.Language) gateBundle[rs.Language] {
	t.Helper()
	return gateBuiltin(t, from, rustBundleTable(), builtinBundleNames)
}

func gateCpp(t *testing.T, from representation.Language) gateBundle[cs.Language] {
	t.Helper()
	return gateBuiltin(t, from, cppBundleTable(), builtinBundleNames)
}

// gateSource parses horos text and analyses it into a generation input.
func gateSource(t *testing.T, text string) (representation.Language, transform.Input) {
	t.Helper()
	parsed, err := parsehoros.Parse(context.Background(), parse.Input{SourceID: "gate.horos", TreeID: "canonical", Text: text})
	if err != nil {
		t.Fatal(err)
	}
	_, from, input := generationSource(t, parsed.Abstract)
	return from, input
}

// gateRequest is what a gate test varies of a request.
type gateRequest struct {
	products      []string
	configuration []g.Setting
	roots         []string // the names of the root declarations; every type definition of the entry module when empty
	limits        *g.Limits
}

// gatePlan registers a bundle in a fresh registry, as a consumer does, and
// plans the request.
func gatePlan[T any](t *testing.T, input transform.Input, checker check.Checker, target g.Target[T], b gateBundle[T], q gateRequest) (g.GenerationPlan, error) {
	t.Helper()
	var registry c.Registry
	from, to := input.Source.Language(), checker.Language()
	if err := c.RegisterSource[h.Language, source.Reader](&registry, from, source.Provider{}); err != nil {
		t.Fatal(err)
	}
	if err := c.RegisterTarget(&registry, to, target); err != nil {
		t.Fatal(err)
	}
	set := shared.Set[h.Language, T, source.Reader]{Contributions: b.contributions, Codecs: b.codecs, Facets: b.facets, Products: b.products}
	if err := set.Register(&registry, from, to); err != nil {
		return g.GenerationPlan{}, err
	}
	roots, err := genfixture.Roots(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(q.roots) > 0 {
		roots = gateRoots(t, input, q.roots)
	}
	limits := g.DefaultLimits()
	if q.limits != nil {
		limits = *q.limits
	}
	return registry.PlanGeneration(context.Background(), c.Request{Input: input, To: to, Output: "generated", Roots: roots, Profile: rules.RepresentationProfile, Limits: limits, Products: q.products, Configuration: q.configuration}, checker)
}

func gatePlanGo(t *testing.T, input transform.Input, b gateBundle[gs.Language], q gateRequest) (g.GenerationPlan, error) {
	t.Helper()
	return gatePlan[gs.Language](t, input, nativeChecker(t, golangText()), ggtarget.Provider{Package: "generated"}, b, q)
}

func gatePlanTS(t *testing.T, input transform.Input, b gateBundle[ts.Language], q gateRequest) (g.GenerationPlan, error) {
	t.Helper()
	return gatePlan[ts.Language](t, input, nativeChecker(t, typescriptText()), tgtarget.Provider{}, b, q)
}

func gatePlanRust(t *testing.T, input transform.Input, b gateBundle[rs.Language], q gateRequest) (g.GenerationPlan, error) {
	t.Helper()
	return gatePlan[rs.Language](t, input, nativeChecker(t, rustText()), rgtarget.Provider{}, b, q)
}

func gatePlanCpp(t *testing.T, input transform.Input, b gateBundle[cs.Language], q gateRequest) (g.GenerationPlan, error) {
	t.Helper()
	return gatePlan[cs.Language](t, input, nativeChecker(t, cppText()), cgtarget.Provider{Namespace: "generated"}, b, q)
}

// gateText runs a plan and prints the generated native text.
func gateText[L any](t *testing.T, p nativeProvider[L], plan g.GenerationPlan) string {
	t.Helper()
	out, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return formatted(t, p, outputTree[L](t, out))
}

// gateRoots finds the declaring sites of the named type definitions.
func gateRoots(t *testing.T, input transform.Input, names []string) []astree.NodeRef {
	t.Helper()
	var out []astree.NodeRef
	for _, envelope := range input.Source.Inputs() {
		tree, err := representation.Unpack[h.Language](input.Source.Language(), envelope)
		if err != nil {
			t.Fatal(err)
		}
		for i, node := range tree.Nodes {
			if node.Kind == h.KindTypeDefinition && slices.Contains(names, nodeAttribute(node, h.AttributeName)) {
				out = append(out, astree.NodeRef{Tree: tree.ID, Node: astree.NodeID(i + 1)})
			}
		}
	}
	if len(out) != len(names) {
		t.Fatalf("found %d of the roots %v", len(out), names)
	}
	return out
}

// gateMust fails the test on an error and returns the value otherwise.
func gateMust[V any](t *testing.T) func(V, error) V {
	return func(v V, err error) V {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
}

// gateSymbol finds the planned declaration with the given source name.
func gateSymbol(plan g.RepresentationPlan, name string) (g.PlanSymbol, bool) {
	for _, s := range plan.Symbols() {
		if strings.HasSuffix(s.ID, "/"+name) {
			return s, true
		}
	}
	return g.PlanSymbol{}, false
}

// gateDeclaration is the node of the declaration with the given source name.
func gateDeclaration(t *testing.T, plan g.RepresentationPlan, name string) g.PlanNode {
	t.Helper()
	s, ok := gateSymbol(plan, name)
	if !ok {
		t.Fatalf("the plan has no symbol %s", name)
	}
	node, _ := plan.Node(s.Node)
	return node
}

// gateBody is the node selected for the body of a declaration.
func gateBody(t *testing.T, plan g.RepresentationPlan, name string) g.PlanNode {
	t.Helper()
	declaration := gateDeclaration(t, plan, name)
	if len(declaration.Links) != 1 {
		t.Fatalf("declaration %s has %d links", name, len(declaration.Links))
	}
	node, _ := plan.Node(declaration.Links[0].Node)
	return node
}

// gateUsesOf returns the use entries that refer to a node.
func gateUsesOf(plan g.RepresentationPlan, node string) []g.PlanUse {
	var out []g.PlanUse
	for _, u := range plan.Uses() {
		if u.Node == node {
			out = append(out, u)
		}
	}
	return out
}

// gateRefusalValues are the eight values of the closed set of refusal classes
// (row 10, O17, O26).
var gateRefusalValues = map[string]error{
	"unsupported": g.ErrUnsupported, "refused by design": g.ErrRefusedByDesign, "requirement unsatisfied": g.ErrRequirementUnsatisfied,
	"requirement unresolved": g.ErrRequirementUnresolved, "missing facet": g.ErrMissingFacet, "ambiguous contribution": g.ErrAmbiguous,
	"missing contribution": g.ErrMissingContribution, "limit": g.ErrLimit,
}

// gateRefusalClasses names every refusal value that an error matches.
func gateRefusalClasses(err error) []string {
	var matched []string
	for name, value := range gateRefusalValues {
		if errors.Is(err, value) {
			matched = append(matched, name)
		}
	}
	sort.Strings(matched)
	return matched
}

// gateExpectRefusal requires that err is a refusal that matches exactly the
// one value named class, and returns it.
func gateExpectRefusal(t *testing.T, name string, err error, class string) *g.Refusal {
	t.Helper()
	var refusal *g.Refusal
	if matched := gateRefusalClasses(err); !errors.As(err, &refusal) || len(matched) != 1 || matched[0] != class {
		t.Errorf("%s: got %v, which matches %v; want a refusal that matches exactly %q", name, err, matched, class)
		return &g.Refusal{}
	}
	return refusal
}

// gateGoTypes parses generated Go and returns, for each struct, its fields in
// order with the type of each as Go spells it, the type parameters of each
// declaration, and for every other declaration the type it is declared as.
type gateGoTypes struct {
	fields     map[string][]gateGoField
	parameters map[string][]string
	aliases    map[string]string
}

type gateGoField struct{ name, typ string }

func (g gateGoTypes) field(declaration, name string) string {
	for _, field := range g.fields[declaration] {
		if field.name == name {
			return field.typ
		}
	}
	return ""
}

func gateParseGo(t *testing.T, text string) gateGoTypes {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "generated.go", text, parser.AllErrors)
	if err != nil {
		t.Fatal(err, text)
	}
	out := gateGoTypes{map[string][]gateGoField{}, map[string][]string{}, map[string]string{}}
	ast.Inspect(file, func(n ast.Node) bool {
		spec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if spec.TypeParams != nil {
			for _, list := range spec.TypeParams.List {
				for _, name := range list.Names {
					out.parameters[spec.Name.Name] = append(out.parameters[spec.Name.Name], name.Name)
				}
			}
		}
		if s, ok := spec.Type.(*ast.StructType); ok {
			out.fields[spec.Name.Name] = []gateGoField{}
			for _, list := range s.Fields.List {
				for _, name := range list.Names {
					out.fields[spec.Name.Name] = append(out.fields[spec.Name.Name], gateGoField{name.Name, types.ExprString(list.Type)})
				}
			}
		} else {
			out.aliases[spec.Name.Name] = types.ExprString(spec.Type)
		}
		return true
	})
	return out
}

// ---------------------------------------------------------------------------
// Fixture contributions for the horos source. Each is written against the
// public packages only and has the gate's publisher.

// gateNilFree is an evidence kind of a fixture vocabulary of the Go target:
// "a Go value of the selected type is never nil". The kernel knows no kind.
var gateNilFree = g.EvidenceKind{Identity: gatePublisher + "/golang/nil-free", Version: "1"}

// gateEmptyCodec is the codec of the fixture's empty payload schema, made once.
func gateEmptyCodec(t *testing.T) g.PayloadCodec {
	t.Helper()
	codec, err := emptyPayloadCodec()
	if err != nil {
		t.Fatal(err)
	}
	return codec
}

// gateNilFreeText is a fixture text contribution that declares at registration
// a fact that its revision fixes: a Go string is never nil (O23). Its payload
// is empty.
func gateNilFreeText(t *testing.T, from representation.Language, representationRevision, meaningRevision string) goContribution {
	t.Helper()
	selector, err := source.Builtin(from, "scalar", "utf8-text")
	if err != nil {
		t.Fatal(err)
	}
	reference := gateReference("nil-free-text", representationRevision)
	reference.Meaning.Revision = meaningRevision
	return goContribution{
		Reference: reference, Target: goTarget, Profile: rules.RepresentationProfile, Domain: g.ForForm(selector), Schema: f.EmptySchema,
		Facts: []g.Fact{gateNilFree.Proved("")},
		Rule:  horosRule{Decide: decideNothing},
		Project: func(context.Context, goProjection) (goProjected, error) {
			return goProjected{Type: gomodel.Primitive(gomodel.PrimitiveString)}, nil
		},
	}
}

// gateDirect is the payload of the fixture families whose one decision is
// whether the direct form is taken.
type gateDirect struct{ Direct bool }

var gateDirectSchema = g.SchemaReference{Identity: gatePublisher + "/payload/direct", Version: "1"}

// gateDirectCodec is the codec of gateDirectSchema.
func gateDirectCodec(t *testing.T) g.PayloadCodec {
	t.Helper()
	return gateMust[g.PayloadCodec](t)(indDirectCodec())
}

// gateEvidenceListOptions configures the fixture family gateEvidenceList.
type gateEvidenceListOptions struct {
	// require declares the fact as a requirement of a binder (the second
	// course of a generic family). Without it the rule makes a choice for want
	// of proof whenever the operand is bound to a binder.
	require bool
	// states are the facts that the rule states beside its decision (O23).
	states []g.Fact
}

// gateEvidenceList is a fixture family whose decision depends on evidence
// about its operand (row 3). It represents a horos list. Where the selected
// element is proved nil-free it takes the direct form []E. Where nothing is
// proved it takes the form [][]E, which needs no proof.
func gateEvidenceList(t *testing.T, from representation.Language, o gateEvidenceListOptions) goContribution {
	t.Helper()
	selector, err := source.Builtin(from, "list-of", "")
	if err != nil {
		t.Fatal(err)
	}
	return goContribution{
		Reference: gateReference("evidence-list", "golang/1"), Target: goTarget, Profile: rules.RepresentationProfile, Domain: g.ForForm(selector), Schema: gateDirectSchema,
		Rule: horosRule{
			Name: func(_ context.Context, step horosNameStep) ([]horosRequest, error) {
				reader, err := step.Reader()
				if err != nil {
					return nil, err
				}
				element, err := reader.Element(step.Subject.Use)
				return []horosRequest{g.Evidence(element, g.Demand{Kind: gateNilFree, Require: o.require})}, err
			},
			Decide: func(_ context.Context, step horosDecideStep) (g.Decision, error) {
				outcome, _, established := step.Evidence(0, gateNilFree)
				return g.Decision{Payload: gateDirect{Direct: outcome == g.Proved || outcome == g.Deferred && established}, Facts: o.states}, nil
			},
		},
		Project: func(_ context.Context, p goProjection) (goProjected, error) {
			view, ok := p.Payload().(gateDirect)
			if !ok {
				return goProjected{}, g.ErrInvalid
			}
			element := p.Dependency(0).Type
			if !view.Direct {
				element = gomodel.Collection(element)
			}
			return goProjected{Type: gomodel.Collection(element)}, nil
		},
	}
}

// gateEvidenceBundle is the Go bundle with the fixture family in the place of
// the list contribution and the fixture text, which declares the fact, in the
// place of the text contribution.
func gateEvidenceBundle(t *testing.T, from representation.Language, o gateEvidenceListOptions) gateBundle[gs.Language] {
	t.Helper()
	return gateGolang(t, from).replace(gateEvidenceList(t, from, o)).replace(gateNilFreeText(t, from, "golang/1", "1")).withCodecs(gateDirectCodec(t), gateEmptyCodec(t))
}

// gateBindings is a product that asks for bindings and requires no facet. It
// is enough to make a request "types and bindings" (O5).
var gateBindings = g.Product{Identity: gatePublisher + "/bindings"}
