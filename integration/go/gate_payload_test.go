package integration_test

// Guarantee 6 of the contribution contract (bittype#196): "The exact payload
// passed validation". The fixture family here has a payload schema with a
// constraint that a Go type cannot express.

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"strings"
	"sync"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	gomodel "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	ontoscodec "github.com/bitspark/ontos/codec/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// gateModed is the payload of the fixture family of guarantee 6: how a list
// is represented. Its Go type is a string. Its schema admits two atoms.
type gateModed struct{ Mode string }

var gateModedSchema = g.SchemaReference{Identity: gatePublisher + "/payload/moded-list", Version: "1"}

// gateModedSource is the schema source: the one field is one of two exact
// atoms, which the Go type string does not express.
const gateModedSource = "module fixtures.gate_moded_list_payload;\n\ntype Payload = { mode: atom-exact(#\"direct\") | atom-exact(#\"boxed\"); };\n"

// gateModedValue is the ontos value of a payload with the given mode, whether
// or not the schema admits the mode.
func gateModedValue(mode string) (core.Value, error) {
	key, err := data.EncodeText("mode")
	if err != nil {
		return nil, err
	}
	return data.EncodeMap([][2]core.Value{{key, core.NewAtom([]byte(mode))}})
}

// gateModedCodec is the codec of gateModedSchema, made once through the
// public route: the schema source is lowered, and the shared part of the
// codec role validates and encodes. The two conversions validate nothing.
var gateModedCodec = sync.OnceValues(func() (g.PayloadCodec, error) {
	return f.Codec(gateModedSchema, gateModedSource,
		func(p gateModed) (core.Value, error) { return gateModedValue(p.Mode) },
		func(v core.Value) (gateModed, error) {
			entries, err := data.ReadMap(v)
			if err != nil || len(entries) != 1 {
				return gateModed{}, errors.Join(errors.New("the payload has one field"), err)
			}
			mode, ok := entries[0][1].(core.Atom)
			if !ok {
				return gateModed{}, errors.New("the mode is no atom")
			}
			return gateModed{string(mode.Bytes())}, nil
		})
})

// gateModedList is the fixture family: it represents a horos list as a slice
// of the selected element, or, in the mode "boxed", as a slice of slices. Its
// rule returns whatever value the test gives it. projected counts the runs of
// its projector.
func gateModedList(t *testing.T, from representation.Language, value any, projected *int) goContribution {
	t.Helper()
	selector, err := source.Builtin(from, "list-of", "")
	if err != nil {
		t.Fatal(err)
	}
	return goContribution{
		Reference: gateReference("moded-list", "golang/1"), Target: goTarget, Profile: rules.RepresentationProfile, Domain: g.ForForm(selector), Schema: gateModedSchema,
		Rule: horosRule{
			Name: nameElement,
			Decide: func(context.Context, horosDecideStep) (g.Decision, error) {
				return g.Decision{Payload: value}, nil
			},
		},
		Project: func(_ context.Context, p goProjection) (goProjected, error) {
			*projected++
			view, ok := p.Payload().(gateModed)
			if !ok {
				return goProjected{}, g.ErrInvalid
			}
			element := p.Dependency(0).Type
			if view.Mode == "boxed" {
				element = gomodel.Collection(element)
			}
			return goProjected{Type: gomodel.Collection(element)}, nil
		},
	}
}

// Row 38 (O9, O20, O37) and guarantee 6, gate line 3 ("A payload outside its
// schema is refused before a backend runs").
//
// Shows, with a fixture family whose schema admits one of two atoms where its
// Go type is a string:
//   - a value inside the schema plans, the projector runs once and reads the
//     decoded view, and the two admitted values give two native types;
//   - a value that the Go type admits and the schema does not is refused when
//     the rule returns, as a payload outside its schema, with the contribution
//     named: no plan results, so nothing reaches a backend, and the projector
//     does not run. The same holds for a value of another Go type and for no
//     value. The refusal is no refusal of a subject: it matches none of the
//     eight values;
//   - the codec that the shared part of the codec role makes refuses such a
//     value when it freezes, and refuses frozen bytes outside the schema when
//     it decodes them;
//   - a contribution whose schema has no registered codec is not planned;
//   - a receipt is the kernel's record over one frozen value and one schema:
//     it is accepted for its own value and schema, and it fails for another
//     frozen value of the same schema, for the frozen value of another
//     schema, for another schema and when it is missing;
//   - a serialized plan whose receipt was issued for another frozen value, or
//     whose receipt is missing, is not read, and neither is one whose payload
//     was exchanged.
//
// Does not show: that the kernel can tell a codec that validates from one
// that does not. The kernel finds the codec by the schema reference that the
// contribution states and trusts it as it trusts a target (O37). What is the
// kernel's is the association of value, schema reference and receipt. It does
// not show who may issue a receipt: the receipt is a digest that anyone who
// holds the value and the schema reference can compute, and the page says so
// ("The receipt is a digest over the schema reference and the frozen bytes").
// Red under: G6-codec-skips-validation (the planner goes on and the projector
// runs), G6-freeze-skips-validation, G6-view-skips-validation,
// G6-receipt-ignores-value, G6-receipt-ignores-schema,
// G6-receipt-ignores-frozen, G6-parse-skips-receipt,
// G6-planner-issues-no-receipt.
func TestGuarantee6ExactPayloadPassedValidation(t *testing.T) {
	from, input := gateSource(t, "module example.moded;\n\ntype Names = list<text>;\n")
	codec := gateMust[g.PayloadCodec](t)(gateModedCodec())
	reference := gateReference("moded-list", "golang/1")
	plan := func(value any) (g.GenerationPlan, int, error) {
		projected := 0
		p, err := gatePlanGo(t, input, gateGolang(t, from).replace(gateModedList(t, from, value, &projected)).withCodecs(codec), gateRequest{})
		return p, projected, err
	}
	// Inside the schema.
	plans := map[string]g.GenerationPlan{}
	for mode, native := range map[string]string{"direct": "[]string", "boxed": "[][]string"} {
		p, projected, err := plan(gateModed{mode})
		if err != nil || projected != 1 {
			t.Fatalf("the mode %q, which the schema admits, does not plan (projector runs: %d): %v", mode, projected, err)
		}
		if got := gateParseGo(t, gateText[gs.Language](t, golangText(), p)).aliases["Names"]; got != native {
			t.Errorf("the mode %q generates Names as %s, expected %s", mode, got, native)
		}
		plans[mode] = p
	}
	// Outside the schema.
	for name, value := range map[string]any{
		"a mode that the Go type admits and the schema does not": gateModed{"sideways"},
		"an empty mode":           gateModed{},
		"a value of another type": gateDirect{Direct: true},
		"no value":                nil,
	} {
		p, projected, err := plan(value)
		if !errors.Is(err, g.ErrPayload) || p.Representation().Valid() {
			t.Errorf("%s: the planner does not refuse the value as a payload outside its schema: %v", name, err)
			continue
		}
		if projected != 0 {
			t.Errorf("%s: the projector ran %d times on a value that was refused", name, projected)
		}
		if !strings.Contains(err.Error(), reference.String()) {
			t.Errorf("%s: the refusal does not name the contribution: %v", name, err)
		}
		var failure *g.RuleFailure
		if matched := gateRefusalClasses(err); len(matched) != 0 || errors.As(err, &failure) {
			t.Errorf("%s: the refusal of the payload matches the refusal values %v, or is a rule failure", name, matched)
		}
	}
	// The codec itself, which the kernel trusts.
	if _, err := codec.Freeze(gateModed{"sideways"}); err == nil {
		t.Errorf("the codec freezes a value outside its schema")
	}
	outside, err := gateModedValue("sideways")
	if err != nil {
		t.Fatal(err)
	}
	if view, err := codec.View(ontoscodec.Encode(outside)); err == nil {
		t.Errorf("the codec decodes frozen bytes outside its schema as %+v", view)
	}
	inside, err := gateModedValue("boxed")
	if err != nil {
		t.Fatal(err)
	}
	if view, err := codec.View(ontoscodec.Encode(inside)); err != nil || view != (gateModed{"boxed"}) {
		t.Errorf("the codec does not decode frozen bytes inside its schema: %+v, %v", view, err)
	}
	// A schema without a registered codec.
	unused := 0
	if _, err := gatePlanGo(t, input, gateGolang(t, from).replace(gateModedList(t, from, gateModed{"boxed"}, &unused)), gateRequest{}); !errors.Is(err, g.ErrIncomplete) || unused != 0 {
		t.Errorf("a contribution whose schema has no registered codec is planned, or fails with another value: %v", err)
	}
	// Receipts.
	boxed, direct := gateBody(t, plans["boxed"].Representation(), "Names"), gateBody(t, plans["direct"].Representation(), "Names")
	if len(boxed.Links) != 1 {
		t.Fatalf("the list has %d links", len(boxed.Links))
	}
	text, _ := plans["boxed"].Representation().Node(boxed.Links[0].Node)
	if boxed.Schema == nil || text.Schema == nil || *boxed.Schema != gateModedSchema || *text.Schema == gateModedSchema || boxed.Payload == direct.Payload {
		t.Fatalf("the nodes are not the expected ones: %+v, %+v, %+v", boxed, direct, text)
	}
	frozen := func(n g.PlanNode) []byte {
		b, err := hex.DecodeString(n.Payload)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	if view, err := codec.View(frozen(boxed)); err != nil || view != (gateModed{"boxed"}) {
		t.Errorf("the frozen record of the list decodes to %+v: %v", view, err)
	}
	if err := g.CheckReceipt(boxed.Receipt, *boxed.Schema, frozen(boxed)); err != nil {
		t.Errorf("a receipt fails for its own frozen value and schema: %v", err)
	}
	if boxed.Receipt == direct.Receipt {
		t.Errorf("two frozen values of one schema have one receipt")
	}
	otherVersion := *boxed.Schema
	otherVersion.Version = "2"
	for name, err := range map[string]error{
		"another frozen value of the same schema":   g.CheckReceipt(boxed.Receipt, *boxed.Schema, frozen(direct)),
		"the frozen value of another schema":        g.CheckReceipt(boxed.Receipt, *boxed.Schema, frozen(text)),
		"another schema":                            g.CheckReceipt(boxed.Receipt, *text.Schema, frozen(boxed)),
		"another version of the schema":             g.CheckReceipt(boxed.Receipt, otherVersion, frozen(boxed)),
		"the receipt of another value, for its own": g.CheckReceipt(direct.Receipt, *boxed.Schema, frozen(boxed)),
		"no receipt": g.CheckReceipt("", *boxed.Schema, frozen(boxed)),
	} {
		if !errors.Is(err, g.ErrPayload) {
			t.Errorf("a receipt is accepted for %s: %v", name, err)
		}
	}
	// A serialized plan.
	serialized, err := plans["boxed"].Representation().Serialize()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.ParsePlan(serialized); err != nil {
		t.Fatalf("the serialized plan is not read: %v", err)
	}
	tamper := func(field, from, to string) []byte {
		old, replacement := []byte(`"`+field+`": "`+from+`"`), []byte(`"`+field+`": "`+to+`"`)
		if bytes.Count(serialized, old) != 1 {
			t.Fatalf("the serialized plan holds the %s of the list %d times", field, bytes.Count(serialized, old))
		}
		return bytes.Replace(serialized, old, replacement, 1)
	}
	if _, err := g.ParsePlan(tamper("receipt", boxed.Receipt, direct.Receipt)); !errors.Is(err, g.ErrPayload) {
		t.Errorf("a plan whose receipt was issued for another frozen value is read, or fails with another value: %v", err)
	}
	if _, err := g.ParsePlan(tamper("receipt", boxed.Receipt, "")); !errors.Is(err, g.ErrPayload) {
		t.Errorf("a plan without the receipt of a payload is read, or fails with another value: %v", err)
	}
	if _, err := g.ParsePlan(tamper("payload", boxed.Payload, direct.Payload)); err == nil {
		t.Errorf("a plan whose payload was exchanged is read")
	}
}
