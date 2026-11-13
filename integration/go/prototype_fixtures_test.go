package integration_test

// PROTOTYPE (bittype#196): fixture contributions of the test composition.
// They are written against the public packages only: the contract, the horos
// source reader, the Go model and the shared part of the codec role. Each has
// its own publisher, its own schema and its own codec.

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/transform/syntax/shared/payload/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

const fixturePublisher = "example.org/fixtures"

// fixtureNilFree is an evidence kind of a fixture vocabulary: "a native Go
// value of the selected type is never nil". It concerns the native side of
// the target, so a vocabulary of the target owns it. The kernel knows no kind.
var fixtureNilFree = g.EvidenceKind{Identity: "example.org/fixtures/golang/nil-free", Version: "1"}

var golangTarget = string(gs.Language{}.ID())

func fixtureReference(key, representation string) g.ContributionReference {
	return g.ContributionReference{Identity: g.ContributionIdentity{Publisher: fixturePublisher, Key: key}, Meaning: g.Meaning{Identity: fixturePublisher + "/" + key, Revision: "1"}, Representation: representation}
}

// fixtureCodec makes the schema and the codec of a fixture family from its
// schema source. A fixture pins its schema when the test starts.
func fixtureCodec[P any](t *testing.T, key, schemaSource string, encode func(P) (payload.Value, error), decode func(payload.Value) (P, error)) (g.SchemaReference, g.PayloadCodec) {
	t.Helper()
	reference := g.SchemaReference{Identity: "example.org/fixtures/payload/" + key, Version: "1"}
	module, pin := protoSchemaModule(t, schemaSource)
	schema, err := payload.Load(reference, module, pin, "Payload")
	if err != nil {
		t.Fatal(err)
	}
	return reference, payload.New(schema, encode, decode)
}

func emptyFixtureCodec(t *testing.T, key string) (g.SchemaReference, g.PayloadCodec) {
	t.Helper()
	return fixtureCodec(t, key, "module fixtures."+key+"_payload;\n\ntype Payload = { };\n",
		func(struct{}) (payload.Value, error) { return payload.Record() },
		func(payload.Value) (struct{}, error) { return struct{}{}, nil })
}

// nilFreeText is a fixture text contribution that declares at registration a
// fact its revision fixes: a Go string is never nil (O23).
func nilFreeText(t *testing.T, from representation.Language, representationRevision, meaningRevision string) (golangContribution, g.PayloadCodec) {
	t.Helper()
	schema, codec := emptyFixtureCodec(t, "nil_free_text")
	selector, err := source.Builtin(from, "scalar", "utf8-text")
	if err != nil {
		t.Fatal(err)
	}
	reference := fixtureReference("nil-free-text", representationRevision)
	reference.Meaning.Revision = meaningRevision
	return golangContribution{
		Reference: reference, Target: golangTarget, Profile: rules.RepresentationProfile, Domain: g.ForForm(selector), Schema: schema,
		Facts: []g.Fact{fixtureNilFree.Proved("")},
		Rule:  horosRule{Decide: func(context.Context, horosDecide) (g.Decision, error) { return g.Decision{Payload: struct{}{}}, nil }},
		Project: func(context.Context, g.Projection[gs.Language]) (g.Projected[gs.Language], error) {
			return g.Projected[gs.Language]{Type: model.Primitive(model.PrimitiveString)}, nil
		},
	}, codec
}

type evidenceListPayload struct{ Direct bool }

// evidenceListOptions configures the fixture family below.
type evidenceListOptions struct {
	// require declares the fact as a requirement of a binder (the second
	// course of a generic family). Without it the rule makes a choice for want
	// of proof whenever the operand is bound to a binder.
	require bool
	// refuseUnproved refuses where the fact is not established, instead of
	// taking the representation that needs no proof.
	refuseUnproved bool
	// states are the facts the rule states beside its decision (O23).
	states []g.Fact
	// stash and replay are the entry point of the test of guarantee 4.
	decide func(step horosDecide) error
}

// evidenceList is a fixture family whose decision depends on evidence about
// its operand (row 3). It represents a horos list. Where the selected element
// is proved nil-free it takes the direct form []E. Where nothing is proved it
// takes the form [][]E, which needs no proof.
func evidenceList(t *testing.T, from representation.Language, o evidenceListOptions) (golangContribution, g.PayloadCodec) {
	t.Helper()
	schema, codec := fixtureCodec(t, "evidence_list", "module fixtures.evidence_list_payload;\n\ntype Payload = { direct: bool; };\n",
		func(p evidenceListPayload) (payload.Value, error) {
			return payload.Record(payload.Field{Name: "direct", Value: payload.Bool(p.Direct)})
		},
		func(v payload.Value) (evidenceListPayload, error) {
			fields, err := payload.ReadRecord(v)
			if err != nil {
				return evidenceListPayload{}, err
			}
			direct, err := payload.ReadBool(fields["direct"])
			return evidenceListPayload{direct}, err
		})
	selector, err := source.Builtin(from, "list-of", "")
	if err != nil {
		t.Fatal(err)
	}
	return golangContribution{
		Reference: fixtureReference("evidence-list", "golang/1"), Target: golangTarget, Profile: rules.RepresentationProfile, Domain: g.ForForm(selector), Schema: schema,
		Rule: horosRule{
			Name: func(_ context.Context, step horosName) ([]horosRequest, error) {
				element, err := step.Reader.Element(step.Subject.Use)
				return []horosRequest{g.Evidence(element, g.Demand{Kind: fixtureNilFree, Require: o.require})}, err
			},
			Decide: func(_ context.Context, step horosDecide) (g.Decision, error) {
				if o.decide != nil {
					if err := o.decide(step); err != nil {
						return g.Decision{}, err
					}
				}
				outcome, _, established := step.Evidence(0, fixtureNilFree)
				direct := outcome == g.Proved || outcome == g.Deferred && established
				if !direct && o.refuseUnproved {
					class := g.RefusalRequirementUnresolved
					if outcome == g.Refuted {
						class = g.RefusalRequirementUnsatisfied
					}
					return g.Decision{}, g.Refuse(class, fixturePublisher, "the element is not proved nil-free: "+outcome.String())
				}
				return g.Decision{Payload: evidenceListPayload{Direct: direct}, Facts: o.states}, nil
			},
		},
		Project: func(_ context.Context, p g.Projection[gs.Language]) (g.Projected[gs.Language], error) {
			view, ok := p.Payload().(evidenceListPayload)
			if !ok {
				return g.Projected[gs.Language]{}, fmt.Errorf("the view is a %T", p.Payload())
			}
			element := p.Dependency(0).Type
			if !view.Direct {
				element = model.Collection(element)
			}
			return g.Projected[gs.Language]{Type: model.Collection(element)}, nil
		},
	}, codec
}

// goFieldTypes parses generated Go and returns, for each struct, the type of
// each field as Go spells it, and the type parameters of each declaration.
func goFieldTypes(t *testing.T, text string) (fields map[string]map[string]string, parameters map[string][]string, aliases map[string]string) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "generated.go", text, parser.AllErrors)
	if err != nil {
		t.Fatal(err, text)
	}
	fields, parameters, aliases = map[string]map[string]string{}, map[string][]string{}, map[string]string{}
	ast.Inspect(file, func(n ast.Node) bool {
		spec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if spec.TypeParams != nil {
			for _, f := range spec.TypeParams.List {
				for _, name := range f.Names {
					parameters[spec.Name.Name] = append(parameters[spec.Name.Name], name.Name)
				}
			}
		}
		if s, ok := spec.Type.(*ast.StructType); ok {
			fields[spec.Name.Name] = map[string]string{}
			for _, f := range s.Fields.List {
				for _, name := range f.Names {
					fields[spec.Name.Name][name.Name] = types.ExprString(f.Type)
				}
			}
		} else {
			aliases[spec.Name.Name] = types.ExprString(spec.Type)
		}
		return true
	})
	return fields, parameters, aliases
}

var _ = h.Profile
