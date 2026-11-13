package generationfixture

// PROTOTYPE (bittype#196): the independent fixture source under the
// contribution contract. The source is not horos: it has its own language,
// analysis, interpretation and reader. The prototype uses it where horos
// cannot state a case: recursion, which the horos analysis does not cover,
// and vertical case 2 at the kernel's level with a mock family.

import (
	"context"

	checkgo "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
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

// prototypeProgram is the standard program and the declarations the
// prototype's cases need.
func prototypeProgram() *graph {
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
	return p
}

// NewPrototype is the fixture over the prototype's program.
func NewPrototype(ctx context.Context) (*Fixture, error) { return newFixture(ctx, prototypeProgram()) }

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
func (f *Fixture) PlanContributions(ctx context.Context, composition Composition, names ...string) (g.ContributionPlan, error) {
	roots, err := f.roots(names)
	if err != nil {
		return g.ContributionPlan{}, err
	}
	var registry c.Registry
	from := f.captured.language
	if err := c.RegisterSource[SourceLanguage, Reader](&registry, from, sourceProvider{}); err != nil {
		return g.ContributionPlan{}, err
	}
	built := composition.Target
	if built == nil {
		built = target.Provider{Package: "generated"}
	}
	if err := c.RegisterTarget(&registry, f.target, built); err != nil {
		return g.ContributionPlan{}, err
	}
	for _, contribution := range composition.Contributions {
		if err := c.RegisterContribution(&registry, from, f.target, contribution); err != nil {
			return g.ContributionPlan{}, err
		}
	}
	for _, codec := range composition.Codecs {
		if err := registry.RegisterCodec(codec); err != nil {
			return g.ContributionPlan{}, err
		}
	}
	for _, facet := range composition.Facets {
		if err := registry.RegisterFacet(facet); err != nil {
			return g.ContributionPlan{}, err
		}
	}
	for _, product := range composition.Products {
		if err := registry.RegisterProduct(product); err != nil {
			return g.ContributionPlan{}, err
		}
	}
	checker, err := check.NewChecker[s.Language](f.target, checkgo.Provider{})
	if err != nil {
		return g.ContributionPlan{}, err
	}
	return registry.PlanContributions(ctx, c.ContributionRequest{Request: c.Request{Input: f.Input(), To: f.target, Roots: roots, Output: "independent", Profile: Profile, Limits: g.DefaultLimits()}, Products: composition.Requested}, checker)
}
