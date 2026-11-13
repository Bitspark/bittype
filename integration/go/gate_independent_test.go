package integration_test

// The fixture contributions of the gate for the independent fixture source
// (bittype#196). None is a contribution of bittype's structural profile and
// none reads horos. They are written against the public packages only: the
// contract, the fixture's reader, the Go model and the shared part of the
// codec role through the fixture's Codec.

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	gomodel "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

type indRule = g.Rule[f.SourceLanguage, f.Reader]
type indNameStep = g.NameStep[f.SourceLanguage, f.Reader]
type indDecideStep = g.DecideStep[f.SourceLanguage, f.Reader]
type indRequest = g.Request[f.SourceLanguage]

// indPolicyPayload is what the fixture's declaration policy records. The name
// of a target parameter states the source binder it stands for: "T1" is the
// parameter of binder 1. Parameters are in the order of the target declaration.
type indPolicyPayload struct {
	Name       string
	Nominal    bool
	Erased     bool
	Parameters []string
}

var indPolicySchema = g.SchemaReference{Identity: gatePublisher + "/payload/independent-policy", Version: "1"}

var indPolicyFields = []string{"name", "nominal", "erased", "parameters"}

// indPolicyCodec is the codec of indPolicySchema, made once.
var indPolicyCodec = sync.OnceValues(func() (g.PayloadCodec, error) {
	return f.Codec(indPolicySchema, "module fixtures.gate_independent_policy_payload;\n\ntype Payload = { name: text; nominal: bool; erased: bool; parameters: list<text>; };\n",
		func(p indPolicyPayload) (core.Value, error) {
			name, err := data.EncodeText(p.Name)
			if err != nil {
				return nil, err
			}
			parameters := make([]core.Value, len(p.Parameters))
			for i, parameter := range p.Parameters {
				if parameters[i], err = data.EncodeText(parameter); err != nil {
					return nil, err
				}
			}
			values := []core.Value{name, data.EncodeBool(p.Nominal), data.EncodeBool(p.Erased), data.EncodeList(parameters)}
			entries := make([][2]core.Value, len(indPolicyFields))
			for i, field := range indPolicyFields {
				key, err := data.EncodeText(field)
				if err != nil {
					return nil, err
				}
				entries[i] = [2]core.Value{key, values[i]}
			}
			return data.EncodeMap(entries)
		},
		func(v core.Value) (indPolicyPayload, error) {
			var p indPolicyPayload
			entries, err := data.ReadMap(v)
			if err != nil {
				return p, err
			}
			fields := map[string]core.Value{}
			for _, entry := range entries {
				key, err := data.ReadText(entry[0])
				if err != nil {
					return p, err
				}
				fields[key] = entry[1]
			}
			if len(fields) != len(indPolicyFields) {
				return p, errors.New("the policy payload has other fields than its schema")
			}
			if p.Name, err = data.ReadText(fields["name"]); err != nil {
				return p, err
			}
			if p.Nominal, err = data.ReadBool(fields["nominal"]); err != nil {
				return p, err
			}
			if p.Erased, err = data.ReadBool(fields["erased"]); err != nil {
				return p, err
			}
			parameters, err := data.ReadList(fields["parameters"])
			if err != nil {
				return p, err
			}
			p.Parameters = make([]string, len(parameters))
			for i, parameter := range parameters {
				if p.Parameters[i], err = data.ReadText(parameter); err != nil {
					return p, err
				}
			}
			return p, nil
		})
})

// indDirectCodec is the codec of gateDirectSchema, made once.
var indDirectCodec = sync.OnceValues(func() (g.PayloadCodec, error) {
	key, err := data.EncodeText("direct")
	if err != nil {
		return nil, err
	}
	return f.Codec(gateDirectSchema, "module fixtures.gate_direct_payload;\n\ntype Payload = { direct: bool; };\n",
		func(p gateDirect) (core.Value, error) {
			return data.EncodeMap([][2]core.Value{{key, data.EncodeBool(p.Direct)}})
		},
		func(v core.Value) (gateDirect, error) {
			entries, err := data.ReadMap(v)
			if err != nil || len(entries) != 1 {
				return gateDirect{}, errors.Join(errors.New("the payload has one field"), err)
			}
			direct, err := data.ReadBool(entries[0][1])
			return gateDirect{direct}, err
		})
})

// indOptions configures the fixture composition for the independent source.
type indOptions struct {
	// probeRequires makes the probe family declare its fact as a requirement
	// of a binder.
	probeRequires bool
	// personFacts are the facts the person family declares beside the record
	// form. wordFacts replace the fact that the word family declares, which is
	// "nil-free, proved".
	personFacts []g.Fact
	wordFacts   []g.Fact
	// cellConversionOnly makes the exact Cell contribution omit its parameter
	// in the target and name its argument for conversion only.
	cellConversionOnly bool
	// reversed names the declarations whose target parameters are in the
	// reverse order of their source binders. This is a decision of the policy's
	// own configuration: the source does not determine it.
	reversed []string
}

func indSelector(t *testing.T, fixture *f.Fixture, key string) g.Selector {
	t.Helper()
	selector, err := g.NewSelector(fixture.Source(), key)
	if err != nil {
		t.Fatal(err)
	}
	return selector
}

// indEmpty is a fixture contribution that records nothing.
func indEmpty(key string, domain g.Domain, name func(context.Context, indNameStep) ([]indRequest, error), project g.Projector[gs.Language], facts ...g.Fact) f.Contribution {
	return f.Contribution{
		Reference: gateReference(key, "golang/1"), Target: goTarget, Profile: f.Profile, Domain: domain, Schema: f.EmptySchema, Facts: facts,
		Rule:    indRule{Name: name, Decide: func(context.Context, indDecideStep) (g.Decision, error) { return g.Decision{Payload: f.Empty{}}, nil }},
		Project: project,
	}
}

// indBinder reads the source binder that a target parameter of the fixture
// policy stands for from the parameter's name.
func indBinder(parameter string) (int, error) {
	var binder int
	if _, err := fmt.Sscanf(parameter, "T%d", &binder); err != nil {
		return 0, fmt.Errorf("the parameter %q names no binder: %w", parameter, err)
	}
	return binder, nil
}

// indComposition is the fixture composition for the independent source: a
// declaration policy, the forms "person-body", "page-body", "pair-body",
// "word", "probe", "carry" and "outcome-body", a mock family for the role of
// an optional nullable field, and an exact contribution for Cell.
func indComposition(t *testing.T, fixture *f.Fixture, o indOptions) f.Composition {
	t.Helper()
	emptyCodec, policyCodec, directCodec := gateEmptyCodec(t), gateMust[g.PayloadCodec](t)(indPolicyCodec()), gateMust[g.PayloadCodec](t)(indDirectCodec())
	// The phantom declaration: its target side omits its source parameter.
	erased := func(name string) bool { return name == "Ignore" }
	policy := f.Contribution{
		Reference: gateReference("independent-declarations", "golang/1"), Target: goTarget, Profile: f.Profile, Domain: g.ForDeclarationPolicy(), Schema: indPolicySchema,
		Rule: indRule{
			Name: func(_ context.Context, step indNameStep) ([]indRequest, error) {
				reader, err := step.Reader()
				if err != nil {
					return nil, err
				}
				name, err := reader.Name(step.Subject.Declaration.Declaration)
				if err != nil || erased(name) {
					return nil, err
				}
				return []indRequest{g.Evidence(step.Subject.Declaration.Body, g.Demand{Kind: gomodel.EvidenceRecordForm})}, nil
			},
			Decide: func(_ context.Context, step indDecideStep) (g.Decision, error) {
				reader, err := step.Reader()
				if err != nil {
					return g.Decision{}, err
				}
				name, err := reader.Name(step.Subject.Declaration.Declaration)
				if err != nil {
					return g.Decision{}, err
				}
				p := indPolicyPayload{Name: name, Erased: erased(name)}
				if !p.Erased {
					outcome, _, _ := step.Evidence(0, gomodel.EvidenceRecordForm)
					p.Nominal = outcome == g.Proved
					for i := range step.Subject.Declaration.Formals {
						p.Parameters = append(p.Parameters, fmt.Sprintf("T%d", i))
					}
					if slices.Contains(o.reversed, name) {
						slices.Reverse(p.Parameters)
					}
				}
				// Where the policy decides an alias it states that the declaration
				// is transparent, and states no fact of its own (O35).
				return g.Decision{Payload: p, Transparent: !p.Nominal && !p.Erased}, nil
			},
		},
		Declare: func(_ context.Context, p goProjection) (g.DeclarationDescription[gs.Language], error) {
			zero := g.DeclarationDescription[gs.Language]{}
			view, ok := p.Payload().(indPolicyPayload)
			if !ok {
				return zero, g.ErrInvalid
			}
			if view.Erased {
				d := g.DeclarationDescription[gs.Language]{Name: view.Name, Form: gomodel.Form(gomodel.FormAlias), Body: gomodel.Primitive(gomodel.PrimitiveString)}
				for i := 0; i < p.Binders(); i++ {
					d.Generic.Bindings = append(d.Generic.Bindings, gomodel.Primitive(gomodel.PrimitiveString))
				}
				return d, nil
			}
			form := gomodel.Form(gomodel.FormAlias)
			if view.Nominal {
				form = gomodel.Form(gomodel.FormDefined)
			}
			d := g.DeclarationDescription[gs.Language]{Name: view.Name, Form: form, Origin: p.Origin(), Body: p.Dependency(0).Type}
			d.Generic.Bindings = make([]g.Expression[gs.Language], len(view.Parameters))
			for position, name := range view.Parameters {
				binder, err := indBinder(name)
				if err != nil || binder < 0 || binder >= len(view.Parameters) {
					return zero, g.ErrInvalid
				}
				d.Generic.Parameters = append(d.Generic.Parameters, g.ParameterDescription{Name: name})
				d.Generic.Bindings[binder] = g.Parameter[gs.Language](position)
				d.Generic.Arguments = append(d.Generic.Arguments, g.Argument[gs.Language](binder))
			}
			return d, nil
		},
	}
	operand := func(_ context.Context, step indNameStep) ([]indRequest, error) {
		child, ok := step.Subject.Use.Operand(0)
		if !ok {
			return nil, g.ErrScope
		}
		return []indRequest{g.Reference(child)}, nil
	}
	person := indEmpty("independent-person", g.ForForm(indSelector(t, fixture, "expression/person-body")), nil,
		func(context.Context, goProjection) (goProjected, error) {
			return goProjected{Type: gomodel.Record(gomodel.Field("Name", false, gomodel.Primitive(gomodel.PrimitiveString))), Members: []g.Member{{Name: "Name"}}}, nil
		}, append([]g.Fact{gomodel.EvidenceRecordForm.Proved("")}, o.personFacts...)...)
	page := indEmpty("independent-page", g.ForForm(indSelector(t, fixture, "expression/page-body")), operand,
		func(_ context.Context, p goProjection) (goProjected, error) {
			return goProjected{Type: gomodel.Record(gomodel.Field("Items", false, gomodel.Collection(p.Dependency(0).Type))), Members: []g.Member{{Name: "Items"}}}, nil
		}, gomodel.EvidenceRecordForm.Proved(""))
	pair := indEmpty("independent-pair", g.ForForm(indSelector(t, fixture, "expression/pair-body")),
		func(_ context.Context, step indNameStep) ([]indRequest, error) {
			first, ok := step.Subject.Use.Operand(0)
			second, ok2 := step.Subject.Use.Operand(1)
			if !ok || !ok2 {
				return nil, g.ErrScope
			}
			return []indRequest{g.Reference(first), g.Reference(second)}, nil
		},
		func(_ context.Context, p goProjection) (goProjected, error) {
			return goProjected{Type: gomodel.Record(gomodel.Field("First", false, p.Dependency(0).Type), gomodel.Field("Second", false, p.Dependency(1).Type)), Members: []g.Member{{Name: "First"}, {Name: "Second"}}}, nil
		}, gomodel.EvidenceRecordForm.Proved(""))
	wordFacts := o.wordFacts
	if wordFacts == nil {
		wordFacts = []g.Fact{gateNilFree.Proved("")}
	}
	word := indEmpty("independent-word", g.ForForm(indSelector(t, fixture, "expression/word")), nil,
		func(context.Context, goProjection) (goProjected, error) {
			return goProjected{Type: gomodel.Primitive(gomodel.PrimitiveString)}, nil
		}, wordFacts...)
	probe := f.Contribution{
		Reference: gateReference("independent-probe", "golang/1"), Target: goTarget, Profile: f.Profile, Domain: g.ForForm(indSelector(t, fixture, "expression/probe")), Schema: gateDirectSchema,
		Rule: indRule{
			Name: func(_ context.Context, step indNameStep) ([]indRequest, error) {
				child, ok := step.Subject.Use.Operand(0)
				if !ok {
					return nil, g.ErrScope
				}
				return []indRequest{g.Evidence(child, g.Demand{Kind: gateNilFree, Require: o.probeRequires})}, nil
			},
			Decide: func(_ context.Context, step indDecideStep) (g.Decision, error) {
				outcome, _, established := step.Evidence(0, gateNilFree)
				return g.Decision{Payload: gateDirect{Direct: outcome == g.Proved || outcome == g.Deferred && established}}, nil
			},
		},
		Project: func(_ context.Context, p goProjection) (goProjected, error) {
			view, ok := p.Payload().(gateDirect)
			if !ok {
				return goProjected{}, g.ErrInvalid
			}
			if view.Direct {
				return goProjected{Type: gomodel.Collection(p.Dependency(0).Type)}, nil
			}
			return goProjected{Type: gomodel.Collection(gomodel.Collection(p.Dependency(0).Type))}, nil
		},
	}
	// A family that names its operand for conversion only: the target type
	// does not refer to it.
	carry := indEmpty("independent-carry", g.ForForm(indSelector(t, fixture, "expression/carry")),
		func(_ context.Context, step indNameStep) ([]indRequest, error) {
			child, ok := step.Subject.Use.Operand(0)
			if !ok {
				return nil, g.ErrScope
			}
			return []indRequest{g.Conversion(child)}, nil
		},
		func(context.Context, goProjection) (goProjected, error) {
			return goProjected{Type: gomodel.Primitive(gomodel.PrimitiveString)}, nil
		})
	// The mock record of vertical case 2: it delegates every field to the role
	// the source issues, and lays out nothing itself.
	outcome := indEmpty("independent-outcome", g.ForForm(indSelector(t, fixture, "expression/outcome-body")),
		func(_ context.Context, step indNameStep) ([]indRequest, error) {
			reader, err := step.Reader()
			if err != nil {
				return nil, err
			}
			roles, err := reader.Roles(step.Subject.Use)
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
			return goProjected{Type: gomodel.Record(members...), Members: table}, nil
		}, gomodel.EvidenceRecordForm.Proved(""))
	// The mock family. It stands for "optional nullable fields", which no built
	// family covers: it represents the field as a slice of the operand.
	mock := indEmpty("mock-optional-nullable", g.ForRole(f.OptionalNullableRole),
		func(_ context.Context, step indNameStep) ([]indRequest, error) {
			return []indRequest{g.Reference(step.Subject.Field.Operand)}, nil
		},
		func(_ context.Context, p goProjection) (goProjected, error) {
			return goProjected{Fragment: []g.Expression[gs.Language]{gomodel.Field(p.Stated(), true, gomodel.Collection(p.Dependency(0).Type))}, Members: []g.Member{{Name: p.Stated()}}}, nil
		})
	// An exact generic. In its first form the target type refers to its
	// argument. In its second it omits the parameter and names the argument
	// for conversion only.
	cell := indEmpty("independent-cell", g.ForExact(indSelector(t, fixture, "declaration/Cell")),
		func(_ context.Context, step indNameStep) ([]indRequest, error) {
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
				return goProjected{Type: gomodel.Primitive(gomodel.PrimitiveString)}, nil
			}
			return goProjected{Type: gomodel.Collection(p.Dependency(0).Type)}, nil
		})
	return f.Composition{Contributions: []f.Contribution{policy, person, page, pair, word, probe, carry, outcome, mock, cell}, Codecs: []g.PayloadCodec{emptyCodec, policyCodec, directCodec}}
}

// indFind returns the index of the fixture contribution with the local key key.
func indFind(t *testing.T, composition f.Composition, key string) int {
	t.Helper()
	i := slices.IndexFunc(composition.Contributions, func(c f.Contribution) bool { return c.Reference.Identity.Key == key })
	if i < 0 {
		t.Fatalf("the fixture composition has no contribution %s", key)
	}
	return i
}

func gateFixture(t *testing.T) *f.Fixture {
	t.Helper()
	fixture, err := f.NewGate(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return fixture
}

// indPlan plans the named declarations of the gate program with the fixture
// composition, for types alone or, with bindings, for types and bindings.
func indPlan(t *testing.T, o indOptions, bindings bool, roots ...string) (g.GenerationPlan, error) {
	t.Helper()
	fixture := gateFixture(t)
	composition := indComposition(t, fixture, o)
	if bindings {
		composition.Products, composition.Requested = []g.Product{gateBindings}, []string{gateBindings.Identity}
	}
	return fixture.PlanContributions(context.Background(), composition, roots...)
}
