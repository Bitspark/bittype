package integration_test

import (
	"context"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	cs "github.com/Bitspark/bittype/syntax/cpp/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	rs "github.com/Bitspark/bittype/syntax/rust/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	cppmodel "github.com/Bitspark/bittype/transform/syntax/cpp/model/go"
	golangmodel "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	rustmodel "github.com/Bitspark/bittype/transform/syntax/rust/model/go"
	typescriptmodel "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
)

// gateFragment is a fixture contribution for the role of an optional record
// field. It answers the field position with one expression that fragment
// makes from the stated member name and the operand's type. It stands in the
// place of the target's builtin optional contribution.
func gateFragment[T any](key, target string, fragment func(name string, operand g.Expression[T]) g.Expression[T]) g.Contribution[h.Language, T, source.Reader] {
	return g.Contribution[h.Language, T, source.Reader]{
		Reference: gateReference(key, target+"/1"), Target: target, Profile: rules.RepresentationProfile,
		Domain: g.ForRole(source.OptionalFieldRole), Schema: f.EmptySchema,
		Rule: g.Rule[h.Language, source.Reader]{
			Name: func(_ context.Context, step g.NameStep[h.Language, source.Reader]) ([]g.Request[h.Language], error) {
				return []g.Request[h.Language]{g.Reference(step.Subject.Field.Operand)}, nil
			},
			Decide: func(context.Context, g.DecideStep[h.Language, source.Reader]) (g.Decision, error) {
				return g.Decision{Payload: f.Empty{}}, nil
			},
		},
		Project: func(_ context.Context, p g.Projection[T]) (g.Projected[T], error) {
			return g.Projected[T]{Fragment: []g.Expression[T]{fragment(p.Stated(), p.Dependency(0).Type)}, Members: []g.Member{{Name: p.Stated()}}}, nil
		},
	}
}

// gateMemberVocabulary runs the two cases of one target: a fragment of the
// target's field form is admitted, and a fragment of another form of the same
// target, which the target admits as an expression, is refused as unsupported
// with the target's member vocabulary as owner.
func gateMemberVocabulary[T any](t *testing.T, target string, bundle gateBundle[T], plan func(gateBundle[T]) (g.GenerationPlan, error), field, other func(name string, operand g.Expression[T]) g.Expression[T]) {
	t.Helper()
	bundle = bundle.withCodecs(gateEmptyCodec(t))
	if p, err := plan(bundle.replace(gateFragment(target+"-member-field", target, field))); err != nil || !p.Valid() {
		t.Errorf("a fragment of the field form of %s is not admitted: %v", target, err)
	}
	p, err := plan(bundle.replace(gateFragment(target+"-member-other", target, other)))
	if p.Valid() {
		t.Errorf("a fragment that is no field of %s was planned", target)
	}
	refusal := gateExpectRefusal(t, "a fragment that is no field of "+target, err, "unsupported")
	if want := "member vocabulary of " + target; refusal.Owner != want {
		t.Errorf("the refusal's owner is %q, want %q", refusal.Owner, want)
	}
}

// Row 14 (and row 13, O40), gate line 2.
//
// Shows, for each of the four targets, with the pagination module, whose Page
// has one optional field: the target states its member vocabulary. A
// contribution at a field position that answers with one field form of the
// target is admitted. One that answers with a collection, which the same
// target admits as a type, is refused as unsupported, the owner of the
// refusal is the member vocabulary of that target, and the refusal matches no
// other value. The rule "a member is one expression of the target's member
// form" is one function of the shared target package, and each target selects
// it with its own form.
// Does not show: a target without a member vocabulary (it refuses every field
// position; that case runs on the independent fixture source), or a fragment
// of more than one member.
// Red under: S1-member-golang, S1-member-typescript, S1-member-rust,
// S1-member-cpp (each target's own statement), S1-member-shared-form (the
// shared rule).
func TestGateMemberFragmentIsTheTargetsFieldForm(t *testing.T) {
	_, from, input := generationInput(t, paginationModule())
	t.Run("golang", func(t *testing.T) {
		gateMemberVocabulary(t, string(gs.Language{}.ID()), gateGolang(t, from),
			func(b gateBundle[gs.Language]) (g.GenerationPlan, error) {
				return gatePlanGo(t, input, b, gateRequest{})
			},
			func(name string, operand g.Expression[gs.Language]) g.Expression[gs.Language] {
				return golangmodel.Field(name, true, operand)
			},
			func(_ string, operand g.Expression[gs.Language]) g.Expression[gs.Language] {
				return golangmodel.Collection(operand)
			})
	})
	t.Run("typescript", func(t *testing.T) {
		gateMemberVocabulary(t, string(ts.Language{}.ID()), gateTypescript(t, from),
			func(b gateBundle[ts.Language]) (g.GenerationPlan, error) {
				return gatePlanTS(t, input, b, gateRequest{})
			},
			func(name string, operand g.Expression[ts.Language]) g.Expression[ts.Language] {
				return typescriptmodel.Field(name, true, operand)
			},
			func(_ string, operand g.Expression[ts.Language]) g.Expression[ts.Language] {
				return typescriptmodel.Collection(operand)
			})
	})
	t.Run("rust", func(t *testing.T) {
		gateMemberVocabulary(t, string(rs.Language{}.ID()), gateRust(t, from),
			func(b gateBundle[rs.Language]) (g.GenerationPlan, error) {
				return gatePlanRust(t, input, b, gateRequest{})
			},
			func(name string, operand g.Expression[rs.Language]) g.Expression[rs.Language] {
				return rustmodel.Field(name, true, operand)
			},
			func(_ string, operand g.Expression[rs.Language]) g.Expression[rs.Language] {
				return rustmodel.Collection(operand)
			})
	})
	t.Run("cpp", func(t *testing.T) {
		gateMemberVocabulary(t, string(cs.Language{}.ID()), gateCpp(t, from),
			func(b gateBundle[cs.Language]) (g.GenerationPlan, error) {
				return gatePlanCpp(t, input, b, gateRequest{})
			},
			func(name string, operand g.Expression[cs.Language]) g.Expression[cs.Language] {
				return cppmodel.Field(name, true, operand)
			},
			func(_ string, operand g.Expression[cs.Language]) g.Expression[cs.Language] {
				return cppmodel.Collection(operand)
			})
	})
}
