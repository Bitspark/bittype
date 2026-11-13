package integration_test

// PROTOTYPE (bittype#196), stage 7: a second reader in another process. Gate
// line 4, the fourth version axis of row 16, and the trigger of row 26a.

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/completeness/golang/backend/go"
	tsvalues "github.com/Bitspark/bittype/completeness/typescript/values/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func standinPlan(t *testing.T, input transform.Input, b bundle[ts.Language]) (g.ContributionPlan, []byte) {
	t.Helper()
	plan, err := protoPlanTS(t, input, b, standinConversion)
	if err != nil {
		t.Fatal(err)
	}
	serialized, err := plan.Plan.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	return plan, serialized
}

// Gate line 4, first test: "The TypeScript stand-in round-trips
// `Page<Person>`".
//
// Shows: a reader in another language and another process, given the
// serialized plan and the source snapshot and nothing of the planning
// process, converts Page<Person> there and back. It finds each factory by the
// module and export that the plan carries, gives it the payload decoded from
// the canonical bytes, and builds the open declaration Page<T> with the
// adapter argument of the application. The facet that names a factory crossed
// the process boundary through the transport descriptor of its kind.
// Does not show: that the service library's factory can be written against
// the plan (row 31a). The stand-in's adapters have no budget, no call context
// and no memo, and one seat wrote both sides.
func TestPrototypeStage7StandinRoundTripsPageOfPerson(t *testing.T) {
	if testing.Short() {
		t.Skip("runs node")
	}
	from, input := protoSource(t, protoPagination)
	plan, serialized := standinPlan(t, input, standinBundle(t, from))
	for _, c := range plan.Plan.Contributions() {
		if len(c.Facets) != 1 || c.Facets[0].MediaType != "application/json" || !strings.Contains(string(c.Facets[0].Value), `"module":"structural.mjs"`) {
			t.Errorf("contribution %s crossed with the facets %s", c.Reference, c.Facets)
		}
	}
	ada, bob := map[string]any{"name": "Ada", "email": "ada@example.org"}, map[string]any{"name": "Bob"}
	vectors := []standinVector{
		{Root: "PeoplePage", Value: map[string]any{"items": []any{ada, bob}, "next": "cursor-2"}},
		{Root: "PeoplePage", Value: map[string]any{"items": []any{}}},
		{Root: "NestedPage", Value: map[string]any{"items": []any{map[string]any{"items": []any{ada}, "next": "n"}}}},
		{Root: "Names", Value: []any{"a", "b"}},
		{Root: "PeoplePage", Value: map[string]any{"items": []any{map[string]any{"email": "x"}}}},
		{Root: "PeoplePage", Value: map[string]any{"items": "none"}},
	}
	out := runStandin(t, serialized, protoSnapshot(t, protoPagination), vectors)
	if out.Failure != nil {
		t.Fatalf("the stand-in failed: %+v", out.Failure)
	}
	for i, v := range vectors[:4] {
		if r := out.Results[i]; r.Error != nil || !reflect.DeepEqual(r.Value, v.Value) || !reflect.DeepEqual(r.Native, v.Value) {
			t.Errorf("vector %d of %s came back as %+v", i, v.Root, r)
		} else {
			t.Logf("GREEN the stand-in converts vector %d of %s there and back: %v", i, v.Root, r.Value)
		}
	}
	for i := 4; i < 6; i++ {
		if r := out.Results[i]; r.Error == nil || r.Error.Kind != "structural mismatch" {
			t.Errorf("vector %d outside the type came back as %+v", i, r)
		} else {
			t.Logf("GREEN the stand-in refuses vector %d: %s", i, r.Error.Message)
		}
	}
	if evidence := protoEvidence(); evidence != "" {
		protoWrite(evidence, "stage7-plan.json", serialized)
		protoWrite(evidence, "stage7-snapshot.json", protoSnapshot(t, protoPagination))
	}
}

// Gate line 4 and row 16: the stand-in "refuses a payload schema version it
// does not know", and each other axis that a reader can check.
//
// Shows: the stand-in refuses, each under its own name, an unknown envelope
// protocol, payload schema version, representation revision, runtime
// compatibility and source semantics, and reads the plan that differs in
// nothing else. The runtime compatibility is the fourth axis: it is stated in
// the facet that names the factory, and only the reader of that facet can
// refuse it.
// Does not show: a reader with more than one known value of an axis.
func TestPrototypeStage7StandinRefusesUnknownVersions(t *testing.T) {
	if testing.Short() {
		t.Skip("runs node")
	}
	const module = "module example.axes;\n\ntype Names = list<text>;\n"
	from, input := protoSource(t, module)
	snapshot := protoSnapshot(t, module)
	base := standinBundle(t, from)
	text := base.find(t, "horos-to-typescript/text")
	_, baseline := standinPlan(t, input, base)
	if out := runStandin(t, baseline, snapshot, nil); out.Failure != nil {
		t.Fatalf("the stand-in refuses the plan it knows: %+v", out.Failure)
	}
	edit := func(change func(map[string]any)) []byte {
		var data map[string]any
		if err := json.Unmarshal(baseline, &data); err != nil {
			t.Fatal(err)
		}
		change(data)
		out, err := json.Marshal(data)
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	// A text contribution with a schema of another version.
	otherSchema := text
	otherSchema.Schema.Version = "2"
	var textCodec g.PayloadCodec
	for _, codec := range base.codecs {
		if codec.Schema() == text.Schema {
			textCodec = codec
		}
	}
	_, schemaPlan := standinPlan(t, input, base.replace(otherSchema).withCodecs(versioned{textCodec, otherSchema.Schema}))
	// A text contribution of another representation revision.
	otherRevision := text
	otherRevision.Reference.Representation = "typescript/2"
	_, revisionPlan := standinPlan(t, input, base.replace(otherRevision))
	// A factory facet that requires another runtime.
	otherRuntime := protoMust[g.Facet](t)(tsvalues.ValueFactory.Issue(text.Reference.Identity.Publisher, text.Reference.Identity, tsvalues.Target, tsvalues.Factory{Module: "structural.mjs", Export: "text", Runtime: "bittype-standin/2"}))
	_, runtimePlan := standinPlan(t, input, base.replaceFacets(text.Reference.Identity, otherRuntime))
	for _, axis := range []struct {
		name, want string
		serialized []byte
	}{
		{"the envelope protocol", "envelope protocol", edit(func(d map[string]any) { d["envelope"] = "bittype.plan/9" })},
		{"the payload schema", "payload schema", schemaPlan},
		{"the representation semantics", "representation semantics", revisionPlan},
		{"generated-code and runtime compatibility", "runtime compatibility", runtimePlan},
		{"the source semantics", "source semantics", edit(func(d map[string]any) { d["source"].(map[string]any)["semantics"] = "horos/declarations-5" })},
	} {
		out := runStandin(t, axis.serialized, snapshot, nil)
		if out.Failure == nil || out.Failure.Kind != "version" || !strings.Contains(out.Failure.Message, axis.want) {
			t.Errorf("%s: an unknown value is not refused under its axis: %+v", axis.name, out.Failure)
		} else {
			t.Logf("GREEN the stand-in refuses an unknown value of %s: %s", axis.name, out.Failure.Message)
		}
	}
}

// Row 34: "Given the plan without the snapshot it fails explicitly, and a
// snapshot with another pin is a linking error".
//
// Shows: serializing the plan alone is not enough for a reader, and both
// readers say so: the stand-in and the Go backend fail explicitly without the
// snapshot. A snapshot of a module with the same name and the same used
// declarations under another pin is a linking error in both, not an
// incomplete plan and not a wrong conversion.
// Does not show: a general reader protocol for source semantics (row 34a).
func TestPrototypeStage7ReadersNeedTheMatchingSnapshot(t *testing.T) {
	if testing.Short() {
		t.Skip("runs node")
	}
	from, input := protoSource(t, protoPagination)
	plan, serialized := standinPlan(t, input, standinBundle(t, from))
	otherPin := protoSnapshot(t, protoPagination+"type Unrelated = text;\n")
	if out := runStandin(t, serialized, nil, nil); out.Failure == nil || out.Failure.Kind != "incomplete" || !strings.Contains(out.Failure.Message, "snapshot") {
		t.Errorf("the stand-in without the snapshot: %+v", out.Failure)
	} else {
		t.Logf("GREEN the stand-in without the snapshot fails explicitly: %s", out.Failure.Message)
	}
	if out := runStandin(t, serialized, otherPin, nil); out.Failure == nil || out.Failure.Kind != "linking" {
		t.Errorf("the stand-in with a snapshot of another pin: %+v", out.Failure)
	} else {
		t.Logf("GREEN the stand-in with a snapshot of another pin: a linking error: %s", out.Failure.Message)
	}
	if _, err := backend.Link(plan.Plan, nil); !errors.Is(err, g.ErrIncomplete) || !strings.Contains(err.Error(), "snapshot") {
		t.Errorf("the Go backend without the snapshot: %v", err)
	}
	if _, err := backend.Link(plan.Plan, otherPin); !errors.Is(err, g.ErrLinking) || errors.Is(err, g.ErrIncomplete) {
		t.Errorf("the Go backend with a snapshot of another pin: %v", err)
	} else {
		t.Logf("GREEN the Go backend with a snapshot of another pin: %v", err)
	}
	if _, err := backend.Link(plan.Plan, protoSnapshot(t, protoPagination)); err != nil {
		t.Errorf("the Go backend with the matching snapshot: %v", err)
	}
}

// Row 22, O20: "The stand-in derives a memo key from the serialized plan
// alone: two uses with equal nodes, applications and selected dependencies
// get equal keys, and a different selected dependency gives a different key".
//
// Shows: the stand-in derives equal keys for two uses of Page<Person> and
// another key when another contribution is selected for a dependency. The
// Go reader of stage 4 derives the same key from the same bytes: the key
// depends on the plan and on nothing of a process.
// Does not show: a memo. The token's form and the states of an entry are the
// service library's (bitspark-service#40).
func TestPrototypeStage7StandinDerivesMemoKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("runs node")
	}
	const module = "module example.memo;\n\ntype Page<T> = { items: list<T>; };\ntype Person = { name: text; };\ntype People = Page<Person>;\ntype PeopleAgain = Page<Person>;\n"
	from, input := protoSource(t, module)
	snapshot := protoSnapshot(t, module)
	base := standinBundle(t, from)
	_, one := standinPlan(t, input, base)
	// Another contribution for text: the same meaning and representation
	// revision under another local key, with its own facet.
	text := base.find(t, "horos-to-typescript/text")
	other := text
	other.Reference.Identity.Key = "horos-to-typescript/text-second"
	facet := protoMust[g.Facet](t)(tsvalues.Issue(tsvalues.ValueFactory, other.Reference, "text"))
	_, two := standinPlan(t, input, base.replace(other).replaceFacets(text.Reference.Identity, facet))
	first, second := runStandin(t, one, snapshot, nil), runStandin(t, two, snapshot, nil)
	if first.Failure != nil || second.Failure != nil {
		t.Fatalf("the stand-in failed: %+v %+v", first.Failure, second.Failure)
	}
	if first.Memo["People"] == "" || first.Memo["People"] != first.Memo["PeopleAgain"] || first.Memo["People"] == second.Memo["People"] {
		t.Errorf("memo keys %v and, with another dependency, %v", first.Memo, second.Memo)
	} else {
		t.Logf("GREEN the stand-in: equal keys %s for two uses of Page<Person>, %s with another selected dependency", first.Memo["People"], second.Memo["People"])
	}
	if fromGo := memoSource(t, one, "People"); fromGo != first.Memo["People"] {
		t.Errorf("the Go reader derives %s and the stand-in %s from the same plan", fromGo, first.Memo["People"])
	} else {
		t.Log("GREEN the Go reader and the stand-in derive the same key from the same serialized plan")
	}
}

var _ = bytes.Equal
var _ representation.Language
