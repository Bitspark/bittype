package generationfixture

// The independent fixture source under the contribution contract
// (bittype#196). The source is not horos: it has its own language, analysis,
// interpretation and reader. The gate uses it where horos cannot state a
// case: recursion, which the horos analysis does not cover, and vertical case
// 2 at the kernel's level with a mock family.

import (
	"context"

	checkgo "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	shared "github.com/Bitspark/bittype/integration/go/testdata/composition"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	target "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Semantics names the fixture's source interpretation: its version axis.
const Semantics = "generation-fixture/1"

// Profile is the profile the fixture's contributions state.
const Profile = "independent/1"

// OptionalNullableRole is the key of the role that this source issues for
// each operand of an "outcome-body" expression: the fields of the mock record
// of vertical case 2. Whether a position has the role is a fact of the source.
const OptionalNullableRole = "generation-fixture/optional-nullable-field"

func (r reader) Roles(view g.UseView[SourceLanguage]) ([]g.Role[SourceLanguage], error) {
	id, ok := view.Use.In(r.captured)
	if !ok {
		return nil, g.ErrScope
	}
	t, _ := r.graph.Term(id)
	roles := make([]g.Role[SourceLanguage], len(view.Operands))
	if t.Operator.Name == "outcome-body" {
		for i, operand := range view.Operands {
			roles[i] = r.roles.Issue(OptionalNullableRole, view.Use, operand.Use)
		}
	}
	return roles, nil
}

// gateProgram is the standard program and the declarations that the cases of
// the gate need.
func gateProgram() *graph {
	p := program()
	// An evidence cycle: a body whose rule asks about an operand that is a
	// use of the declaration itself.
	p.declare("Loop")
	p.body("Loop", p.expression("probe", p.ref(p.names["Loop"])))
	// A generic that requires evidence about its argument (vertical case 5).
	p.declare("Probed", "T")
	p.body("Probed", p.expression("probe", p.ref(p.parameter("Probed", 0))))
	p.declare("ProbedPerson")
	p.body("ProbedPerson", p.apply("Probed", p.ref(p.names["Person"])))
	p.declare("ProbedPage")
	p.body("ProbedPage", p.apply("Probed", p.apply("Page", p.ref(p.names["Person"]))))
	// Vertical case 2 at the kernel's level: Page<Outcome>, where the fields
	// of Outcome have a role that a mock family answers.
	p.declare("Outcome")
	p.body("Outcome", p.expression("outcome-body", p.expression("person-body"), p.ref(p.names["Person"])))
	p.declare("Outcomes")
	p.body("Outcomes", p.apply("Page", p.ref(p.names["Outcome"])))
	// A second use of Page<Person>, and two applications of the phantom
	// declaration Ignore with different arguments.
	p.declare("PeopleAgain")
	p.body("PeopleAgain", p.apply("Page", p.ref(p.names["Person"])))
	p.declare("KeptPage")
	p.body("KeptPage", p.apply("Ignore", p.apply("Page", p.ref(p.names["Person"]))))
	// An exact generic whose native side omits its parameter.
	p.declare("CellOfPerson")
	p.body("CellOfPerson", p.apply("Cell", p.ref(p.names["Person"])))
	// A second generic whose body asks about its own binder: another generic
	// environment than Probed's (guarantee 4).
	p.declare("ProbedToo", "U")
	p.body("ProbedToo", p.expression("probe", p.ref(p.parameter("ProbedToo", 0))))
	// Evidence through a declaration (O35). Word's body is no record, so a
	// declaration contribution decides an alias for it. Wrap's body is its binder.
	p.declare("Word")
	p.body("Word", p.expression("word"))
	p.declare("ProbeWord")
	p.body("ProbeWord", p.expression("probe", p.ref(p.names["Word"])))
	p.declare("ProbePerson")
	p.body("ProbePerson", p.expression("probe", p.ref(p.names["Person"])))
	p.declare("Wrap", "T")
	p.body("Wrap", p.ref(p.parameter("Wrap", 0)))
	p.declare("ProbeWrapped")
	p.body("ProbeWrapped", p.expression("probe", p.apply("Wrap", p.ref(p.names["Word"]))))
	p.declare("ProbedWord")
	p.body("ProbedWord", p.apply("Probed", p.ref(p.names["Word"])))
	// A cycle whose one request names its operand for conversion only (O40).
	p.declare("Spin")
	p.body("Spin", p.expression("carry", p.ref(p.names["Spin"])))
	// A generic with two binders and one application of it, for the ties of an
	// application's arguments to binders (row 20, row 39).
	p.declare("Pair", "A", "B")
	p.body("Pair", p.expression("pair-body", p.ref(p.parameter("Pair", 0)), p.ref(p.parameter("Pair", 1))))
	p.declare("Words")
	p.body("Words", p.apply("Pair", p.ref(p.names["Word"]), p.ref(p.names["Person"])))
	// One refuted term that two bodies share. The source judges it and states
	// no site: the kernel reports each occurrence at its own declaration.
	p.declare("FaultA")
	p.declare("FaultB")
	fault := p.expression("person-body")
	p.refute(fault, false)
	p.body("FaultA", fault)
	p.body("FaultB", fault)
	// A form of a family that no composition of the gate builds (row 33).
	p.declare("Unbuilt")
	p.body("Unbuilt", p.expression("unbuilt-family"))
	return p
}

// NewGate is the fixture over the program of the gate's cases.
func NewGate(ctx context.Context) (*Fixture, error) { return newFixture(ctx, gateProgram()) }

// Source is the fixture's source language.
func (f *Fixture) Source() representation.Language { return f.captured.language }

// Contribution is a contribution for the fixture source and the Go target.
type Contribution = g.Contribution[SourceLanguage, s.Language, Reader]

// Composition is what a consumer registers for a request.
type Composition struct {
	Contributions []Contribution
	Codecs        []g.PayloadCodec
	Facets        []g.Facet
	Products      []g.Product
	// Requested are the products the request names besides native types.
	Requested []string
	// Target replaces the Go target, for a test that needs another one.
	Target g.Target[s.Language]
}

// PlanContributions plans the named declarations through the public registry
// under the contribution contract.
func (f *Fixture) PlanContributions(ctx context.Context, composition Composition, names ...string) (g.GenerationPlan, error) {
	roots, err := f.roots(names)
	if err != nil {
		return g.GenerationPlan{}, err
	}
	var registry c.Registry
	from := f.captured.language
	if err := c.RegisterSource[SourceLanguage, Reader](&registry, from, sourceProvider{}); err != nil {
		return g.GenerationPlan{}, err
	}
	built := composition.Target
	if built == nil {
		built = target.Provider{Package: "generated"}
	}
	if err := c.RegisterTarget(&registry, f.target, built); err != nil {
		return g.GenerationPlan{}, err
	}
	if err := composition.register(&registry, from, f.target); err != nil {
		return g.GenerationPlan{}, err
	}
	checker, err := check.NewChecker[s.Language](f.target, checkgo.Provider{})
	if err != nil {
		return g.GenerationPlan{}, err
	}
	return registry.PlanGeneration(ctx, c.Request{Input: f.Input(), To: f.target, Roots: roots, Output: "independent", Profile: Profile, Limits: g.DefaultLimits(), Products: composition.Requested}, checker)
}

// register registers what a composition holds, as a consumer does.
func (composition Composition) register(registry *c.Registry, from, to representation.Language) error {
	return shared.Set[SourceLanguage, s.Language, Reader]{Contributions: composition.Contributions, Codecs: composition.Codecs, Facets: composition.Facets, Products: composition.Products}.Register(registry, from, to)
}
