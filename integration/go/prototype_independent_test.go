package integration_test

// PROTOTYPE (bittype#196): contributions for the independent fixture source.
// None of them is a contribution of bittype's structural profile, and none
// reads horos. They are written against the public packages only.

import (
	"context"
	"fmt"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	"github.com/Bitspark/bittype/transform/syntax/shared/payload/go"
)

type indRule = g.Rule[f.SourceLanguage, f.Reader]
type indName = g.NameStep[f.SourceLanguage, f.Reader]
type indDecide = g.DecideStep[f.SourceLanguage, f.Reader]
type indRequest = g.Request[f.SourceLanguage]
type goProjection = g.Projection[gs.Language]
type goProjected = g.Projected[gs.Language]

type indPolicyPayload struct {
	Name       string
	Nominal    bool
	Erased     bool
	Parameters []string
}

type indProbePayload struct{ Direct bool }

// independentOptions configures the fixture composition for the independent
// source.
type independentOptions struct {
	// probeRequires makes the probe family declare its fact as a requirement
	// of a binder. personFacts are the facts the person family declares.
	probeRequires bool
	personFacts   []g.Fact
	// cellConversionOnly makes the exact Cell contribution omit its parameter
	// natively and name its argument for conversion only.
	cellConversionOnly bool
}

func indSelector(t *testing.T, fixture *f.Fixture, key string) g.Selector {
	t.Helper()
	selector, err := g.NewSelector(fixture.Source(), key)
	if err != nil {
		t.Fatal(err)
	}
	return selector
}

func indEmpty(key string, domain g.Domain, schema g.SchemaReference, name func(context.Context, indName) ([]indRequest, error), project g.Projector[gs.Language], facts ...g.Fact) f.Contribution {
	return f.Contribution{
		Reference: fixtureReference(key, "golang/1"), Target: golangTarget, Profile: f.Profile, Domain: domain, Schema: schema, Facts: facts,
		Rule:    indRule{Name: name, Decide: func(context.Context, indDecide) (g.Decision, error) { return g.Decision{Payload: struct{}{}}, nil }},
		Project: project,
	}
}

// independentComposition is the fixture composition for the independent
// source: a declaration policy, the forms "person-body", "page-body",
// "probe" and "outcome-body", a mock family for the role of an optional
// nullable field, and an exact contribution for Cell.
func independentComposition(t *testing.T, fixture *f.Fixture, o independentOptions) f.Composition {
	t.Helper()
	empty, emptyCodec := emptyFixtureCodec(t, "independent_empty")
	policySchema, policyCodec := fixtureCodec(t, "independent_policy", "module fixtures.independent_policy_payload;\n\ntype Payload = { name: text; nominal: bool; erased: bool; parameters: list<text>; };\n",
		func(p indPolicyPayload) (payload.Value, error) {
			name, err := payload.Text(p.Name)
			if err != nil {
				return nil, err
			}
			parameters, err := payload.Texts(p.Parameters)
			if err != nil {
				return nil, err
			}
			return payload.Record(payload.Field{Name: "name", Value: name}, payload.Field{Name: "nominal", Value: payload.Bool(p.Nominal)}, payload.Field{Name: "erased", Value: payload.Bool(p.Erased)}, payload.Field{Name: "parameters", Value: parameters})
		},
		func(v payload.Value) (indPolicyPayload, error) {
			fields, err := payload.ReadRecord(v)
			if err != nil {
				return indPolicyPayload{}, err
			}
			var p indPolicyPayload
			if p.Name, err = payload.ReadText(fields["name"]); err != nil {
				return p, err
			}
			if p.Nominal, err = payload.ReadBool(fields["nominal"]); err != nil {
				return p, err
			}
			if p.Erased, err = payload.ReadBool(fields["erased"]); err != nil {
				return p, err
			}
			p.Parameters, err = payload.ReadTexts(fields["parameters"])
			return p, err
		})
	probeSchema, probeCodec := fixtureCodec(t, "independent_probe", "module fixtures.independent_probe_payload;\n\ntype Payload = { direct: bool; };\n",
		func(p indProbePayload) (payload.Value, error) {
			return payload.Record(payload.Field{Name: "direct", Value: payload.Bool(p.Direct)})
		},
		func(v payload.Value) (indProbePayload, error) {
			fields, err := payload.ReadRecord(v)
			if err != nil {
				return indProbePayload{}, err
			}
			direct, err := payload.ReadBool(fields["direct"])
			return indProbePayload{direct}, err
		})
	// The phantom declaration: its native side omits its source parameter.
	erased := func(name string) bool { return name == "Ignore" }
	policy := f.Contribution{
		Reference: fixtureReference("independent-declarations", "golang/1"), Target: golangTarget, Profile: f.Profile, Domain: g.ForDeclarationPolicy(), Schema: policySchema,
		Rule: indRule{
			Name: func(_ context.Context, step indName) ([]indRequest, error) {
				name, err := step.Reader.Name(step.Subject.Declaration.Declaration)
				if err != nil || erased(name) {
					return nil, err
				}
				return []indRequest{g.Evidence(step.Subject.Declaration.Body, g.Demand{Kind: model.EvidenceRecordForm}, g.Demand{Kind: fixtureNilFree})}, nil
			},
			Decide: func(_ context.Context, step indDecide) (g.Decision, error) {
				name, err := step.Reader.Name(step.Subject.Declaration.Declaration)
				if err != nil {
					return g.Decision{}, err
				}
				p := indPolicyPayload{Name: name, Erased: erased(name)}
				if !p.Erased {
					outcome, _, _ := step.Evidence(0, model.EvidenceRecordForm)
					p.Nominal = outcome == g.Proved
					for i := range step.Subject.Declaration.Formals {
						p.Parameters = append(p.Parameters, fmt.Sprintf("T%d", i))
					}
				}
				// A use of a declaration is answered from what the declaration
				// contribution states. This policy states beside its decision
				// what was established about its body (O23).
				var facts []g.Fact
				if !p.Erased {
					switch outcome, content, _ := step.Evidence(0, fixtureNilFree); outcome {
					case g.Proved:
						facts = append(facts, fixtureNilFree.Proved(content))
					case g.Refuted:
						facts = append(facts, fixtureNilFree.Refuted(content))
					}
				}
				return g.Decision{Payload: p, Facts: facts}, nil
			},
		},
		Declare: func(_ context.Context, p goProjection) (g.DeclarationDescription[gs.Language], error) {
			view := p.Payload().(indPolicyPayload)
			if view.Erased {
				d := g.DeclarationDescription[gs.Language]{Name: view.Name, Form: model.Form(model.FormAlias), Body: model.Primitive(model.PrimitiveString)}
				for i := 0; i < p.Binders(); i++ {
					d.Generic.Bindings = append(d.Generic.Bindings, model.Primitive(model.PrimitiveString))
				}
				return d, nil
			}
			form := model.Form(model.FormAlias)
			if view.Nominal {
				form = model.Form(model.FormDefined)
			}
			d := g.DeclarationDescription[gs.Language]{Name: view.Name, Form: form, Origin: p.Origin(), Body: p.Dependency(0).Type}
			for i, name := range view.Parameters {
				d.Generic.Parameters = append(d.Generic.Parameters, g.ParameterDescription{Name: name})
				d.Generic.Bindings = append(d.Generic.Bindings, g.Parameter[gs.Language](i))
				d.Generic.Arguments = append(d.Generic.Arguments, g.Argument[gs.Language](i))
			}
			return d, nil
		},
	}
	operand := func(_ context.Context, step indName) ([]indRequest, error) {
		child, ok := step.Subject.Use.Operand(0)
		if !ok {
			return nil, g.ErrScope
		}
		return []indRequest{g.Reference(child)}, nil
	}
	person := indEmpty("independent-person", g.ForForm(indSelector(t, fixture, "expression/person-body")), empty, nil,
		func(context.Context, goProjection) (goProjected, error) {
			return goProjected{Type: model.Record(model.Field("Name", false, model.Primitive(model.PrimitiveString))), Members: []g.Member{{Name: "Name"}}}, nil
		}, append([]g.Fact{model.EvidenceRecordForm.Proved("")}, o.personFacts...)...)
	page := indEmpty("independent-page", g.ForForm(indSelector(t, fixture, "expression/page-body")), empty, operand,
		func(_ context.Context, p goProjection) (goProjected, error) {
			return goProjected{Type: model.Record(model.Field("Items", false, model.Collection(p.Dependency(0).Type))), Members: []g.Member{{Name: "Items"}}}, nil
		}, model.EvidenceRecordForm.Proved(""))
	probe := f.Contribution{
		Reference: fixtureReference("independent-probe", "golang/1"), Target: golangTarget, Profile: f.Profile, Domain: g.ForForm(indSelector(t, fixture, "expression/probe")), Schema: probeSchema,
		Rule: indRule{
			Name: func(_ context.Context, step indName) ([]indRequest, error) {
				child, ok := step.Subject.Use.Operand(0)
				if !ok {
					return nil, g.ErrScope
				}
				return []indRequest{g.Evidence(child, g.Demand{Kind: fixtureNilFree, Require: o.probeRequires})}, nil
			},
			Decide: func(_ context.Context, step indDecide) (g.Decision, error) {
				outcome, _, established := step.Evidence(0, fixtureNilFree)
				return g.Decision{Payload: indProbePayload{Direct: outcome == g.Proved || outcome == g.Deferred && established}}, nil
			},
		},
		Project: func(_ context.Context, p goProjection) (goProjected, error) {
			if p.Payload().(indProbePayload).Direct {
				return goProjected{Type: model.Collection(p.Dependency(0).Type)}, nil
			}
			return goProjected{Type: model.Collection(model.Collection(p.Dependency(0).Type))}, nil
		},
	}
	// The mock record of vertical case 2: it delegates every field to the
	// role the source issues, and lays out nothing itself.
	outcome := indEmpty("independent-outcome", g.ForForm(indSelector(t, fixture, "expression/outcome-body")), empty,
		func(_ context.Context, step indName) ([]indRequest, error) {
			roles, err := step.Reader.Roles(step.Subject.Use)
			if err != nil {
				return nil, err
			}
			var requests []indRequest
			for i, operand := range step.Subject.Use.Operands {
				requests = append(requests, g.Field(operand.Use, roles[i], fmt.Sprintf("Field%d", i)))
			}
			return requests, nil
		},
		func(_ context.Context, p goProjection) (goProjected, error) {
			var members []g.Expression[gs.Language]
			var table []g.Member
			for i := 0; i < p.Dependencies(); i++ {
				members = append(members, p.Dependency(i).Fragment...)
				table = append(table, g.FragmentMembers(i))
			}
			return goProjected{Type: model.Record(members...), Members: table}, nil
		}, model.EvidenceRecordForm.Proved(""))
	// The mock family. It stands for "optional nullable fields", which no
	// built family covers: it represents the field as a slice of the operand.
	mock := indEmpty("mock-optional-nullable", g.ForRole(f.OptionalNullableRole), empty,
		func(_ context.Context, step indName) ([]indRequest, error) {
			return []indRequest{g.Reference(step.Subject.Field.Operand)}, nil
		},
		func(_ context.Context, p goProjection) (goProjected, error) {
			return goProjected{Fragment: []g.Expression[gs.Language]{model.Field(p.Stated(), true, model.Collection(p.Dependency(0).Type))}, Members: []g.Member{{Name: p.Stated()}}}, nil
		})
	// An exact generic. In its first form the native type refers to its
	// argument. In its second it omits the parameter natively and names the
	// argument for conversion only.
	cell := indEmpty("independent-cell", g.ForExact(indSelector(t, fixture, "declaration/Cell")), empty,
		func(_ context.Context, step indName) ([]indRequest, error) {
			if len(step.Subject.Use.Actuals) != 1 {
				return nil, g.ErrScope
			}
			if o.cellConversionOnly {
				return []indRequest{g.Conversion(step.Subject.Use.Actuals[0].Use)}, nil
			}
			return []indRequest{g.Reference(step.Subject.Use.Actuals[0].Use)}, nil
		},
		func(_ context.Context, p goProjection) (goProjected, error) {
			if o.cellConversionOnly {
				return goProjected{Type: model.Primitive(model.PrimitiveString)}, nil
			}
			return goProjected{Type: model.Collection(p.Dependency(0).Type)}, nil
		})
	return f.Composition{Contributions: []f.Contribution{policy, person, page, probe, outcome, mock, cell}, Codecs: []g.PayloadCodec{emptyCodec, policyCodec, probeCodec}}
}

func protoFixture(t *testing.T) *f.Fixture {
	t.Helper()
	fixture, err := f.NewPrototype(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return fixture
}
