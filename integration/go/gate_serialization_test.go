package integration_test

// Gate line 3, the kernel's part (bittype#196, docs/contribution-contract.md,
// "Identities", row 16, and "The gate", line 3): each version axis that the
// plan states is exercised alone, and the serialized plan, parsed back, is
// the plan that was written. The fourth axis, the runtime compatibility, is a
// test of the reader of the facet that carries it and is not here.

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// gateAxes are the values of the four version axes that a plan states, as a
// reader finds them.
type gateAxes struct {
	envelope  string
	semantics string
	schemas   []g.SchemaReference
	meanings  []string // identity and revision of the meaning of each selected contribution
	revisions []string // identity and representation revision of each selected contribution
}

// gateAxesOf reads the axes of a serialized plan without the planner: the
// envelope from the bytes, the rest from the parsed plan.
func gateAxesOf(t *testing.T, serialized []byte) gateAxes {
	t.Helper()
	var envelope struct {
		Envelope string `json:"envelope"`
	}
	if err := json.Unmarshal(serialized, &envelope); err != nil {
		t.Fatal(err)
	}
	plan, err := g.ParsePlan(serialized)
	if err != nil {
		t.Fatal(err)
	}
	axes := gateAxes{envelope: envelope.Envelope, semantics: plan.Source().Semantics}
	for _, c := range plan.Contributions() {
		axes.schemas = append(axes.schemas, c.Schema)
		axes.meanings = append(axes.meanings, c.Reference.Identity.Publisher+"/"+c.Reference.Identity.Key+" "+c.Reference.Meaning.Identity+"@"+c.Reference.Meaning.Revision)
		axes.revisions = append(axes.revisions, c.Reference.Identity.Publisher+"/"+c.Reference.Identity.Key+"#"+c.Reference.Representation)
	}
	return axes
}

// gateKnowledge is what a reader knows that can read exactly this plan.
func gateKnowledge(plan g.RepresentationPlan) g.ReaderKnowledge {
	k := g.ReaderKnowledge{SourceSemantics: []string{plan.Source().Semantics}}
	for _, c := range plan.Contributions() {
		k.Schemas, k.Contributions = append(k.Schemas, c.Schema), append(k.Contributions, c.Reference)
	}
	return k
}

// gateReplaceOnce replaces the one occurrence of a text in a serialized plan.
func gateReplaceOnce(t *testing.T, serialized []byte, old, replacement string) []byte {
	t.Helper()
	if bytes.Count(serialized, []byte(old)) != 1 {
		t.Fatalf("the serialized plan holds %q %d times", old, bytes.Count(serialized, []byte(old)))
	}
	return bytes.Replace(serialized, []byte(old), []byte(replacement), 1)
}

// Row 16 (O19, O24, O39), gate line 3.
//
// Shows: the plan states four of the five version axes, each where it
// belongs: the envelope protocol and the source semantics on the plan, the
// payload schema and the contribution reference, with the revision of the
// family's meaning and the representation revision, on each selected
// contribution. A reader that knows exactly the values of one plan reads it.
// For each axis, a plan in which one required value of that axis is another
// is refused as an unknown required version, under the name of that axis and
// with the value, while a comparison of the two plans shows the other axes
// unchanged. The envelope and the source semantics are changed in the
// serialized bytes. The payload schema, the meaning revision and the
// representation revision are changed by planning with a contribution that
// states the other value, because a node's identity covers them. The third
// axis is exercised for each of its two parts.
// Does not show: the fourth axis, generated-code and runtime compatibility.
// It is a test of the reader of the facet that carries it, because the kernel
// carries that facet unread (O39). A reader in another process: the readers
// here are functions of this process that hold only the serialized bytes.
// Red under: A-envelope-unchecked, A-plan-states-no-source-semantics,
// A-source-semantics-unchecked, A-schema-unchecked,
// A-schema-version-unchecked, A-reference-unchecked,
// A-meaning-revision-named-as-representation,
// A-representation-revision-named-as-meaning.
func TestGateEachVersionAxisAlone(t *testing.T) {
	from, input := gateSource(t, "module example.axes;\n\ntype Names = list<text>;\n")
	schemaTwo := g.SchemaReference{Identity: f.EmptySchema.Identity, Version: "2"}
	codecTwo, err := f.Codec(schemaTwo, "module fixtures.empty_payload;\n\ntype Payload = { };\n",
		func(f.Empty) (core.Value, error) { return data.EncodeMap(nil) },
		func(core.Value) (f.Empty, error) { return f.Empty{}, nil })
	if err != nil {
		t.Fatal(err)
	}
	serialize := func(change func(*goContribution)) []byte {
		text := gateNilFreeText(t, from, "golang/1", "1")
		if change != nil {
			change(&text)
		}
		generation, err := gatePlanGo(t, input, gateGolang(t, from).replace(text).withCodecs(gateEmptyCodec(t), codecTwo), gateRequest{})
		if err != nil {
			t.Fatal(err)
		}
		return gateMust[[]byte](t)(generation.Representation().Serialize())
	}
	base := serialize(nil)
	baseline := gateAxesOf(t, base)
	basePlan := gateMust[g.RepresentationPlan](t)(g.ParsePlan(base))
	knowledge := gateKnowledge(basePlan)
	if baseline.envelope != g.EnvelopeVersion || baseline.semantics != source.Semantics || baseline.semantics == "" || len(baseline.schemas) == 0 {
		t.Fatalf("the plan states the axes %+v", baseline)
	}
	if err := basePlan.CheckKnown(knowledge); err != nil {
		t.Fatalf("a reader that knows the plan's own values does not read it: %v", err)
	}
	// read is a reader with the baseline's knowledge.
	read := func(serialized []byte) error {
		plan, err := g.ParsePlan(serialized)
		if err != nil {
			return err
		}
		return plan.CheckKnown(knowledge)
	}
	axisNames := []string{"envelope protocol", "payload schema", "meaning revision", "representation revision", "source semantics"}
	for _, c := range []struct {
		axis, value string
		serialized  []byte
		// changed names the one field of gateAxes that differs from the baseline.
		changed string
	}{
		{"envelope protocol", "bittype.plan/another", gateReplaceOnce(t, base, `"envelope": "`+g.EnvelopeVersion+`"`, `"envelope": "bittype.plan/another"`), "envelope"},
		{"source semantics", "another-semantics/9", gateReplaceOnce(t, base, `"semantics": "`+source.Semantics+`"`, `"semantics": "another-semantics/9"`), "semantics"},
		{"payload schema", `"2"`, serialize(func(c *goContribution) { c.Schema = schemaTwo }), "schemas"},
		{"meaning revision", `"2"`, serialize(func(c *goContribution) { c.Reference.Meaning.Revision = "2" }), "meanings"},
		{"representation revision", `"golang/2"`, serialize(func(c *goContribution) { c.Reference.Representation = "golang/2" }), "revisions"},
	} {
		err := read(c.serialized)
		if !errors.Is(err, g.ErrVersion) {
			t.Errorf("another %s: got %v, want an unknown required version", c.axis, err)
			continue
		}
		for _, name := range axisNames {
			if strings.Contains(err.Error(), name) != (name == c.axis) {
				t.Errorf("another %s is refused as %q, which names %q or does not name its own axis", c.axis, err, name)
			}
		}
		if !strings.Contains(err.Error(), c.value) {
			t.Errorf("another %s is refused without its value %s: %v", c.axis, c.value, err)
		}
		// The other axes are unchanged. The envelope is read from the bytes,
		// because a plan with another envelope is not parsed.
		var other gateAxes
		if c.changed == "envelope" {
			other = baseline
			var envelope struct {
				Envelope string `json:"envelope"`
			}
			if err := json.Unmarshal(c.serialized, &envelope); err != nil {
				t.Fatal(err)
			}
			other.envelope = envelope.Envelope
			if restored := gateReplaceOnce(t, c.serialized, `"envelope": "bittype.plan/another"`, `"envelope": "`+g.EnvelopeVersion+`"`); !bytes.Equal(restored, base) {
				t.Errorf("the plan with another envelope differs from the baseline in more than its envelope")
			}
		} else {
			other = gateAxesOf(t, c.serialized)
		}
		differs := map[string]bool{
			"envelope": other.envelope != baseline.envelope, "semantics": other.semantics != baseline.semantics, "schemas": !slices.Equal(other.schemas, baseline.schemas),
			"meanings": !slices.Equal(other.meanings, baseline.meanings), "revisions": !slices.Equal(other.revisions, baseline.revisions),
		}
		for field, differ := range differs {
			if differ != (field == c.changed) {
				t.Errorf("another %s: the axis %s differs from the baseline: %v, want only %s to differ", c.axis, field, differ, c.changed)
			}
		}
	}
}

// gatePlanJSON is a serialized plan as a reader without the contract's types
// sees it.
type gatePlanJSON map[string]any

func gateDecodePlan(t *testing.T, serialized []byte) gatePlanJSON {
	t.Helper()
	var plan gatePlanJSON
	if err := json.Unmarshal(serialized, &plan); err != nil {
		t.Fatal(err)
	}
	return plan
}

func (p gatePlanJSON) encode(t *testing.T) []byte {
	t.Helper()
	encoded, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func (p gatePlanJSON) list(name string) []any { list, _ := p[name].([]any); return list }

// gateCompactFacets returns the entries of the selected contributions with
// the value of each facet as compact JSON, so that two entries compare by the
// JSON value of a facet and not by its white space (finding F-B4).
func gateCompactFacets(t *testing.T, entries []g.PlanContribution) []g.PlanContribution {
	t.Helper()
	out := slices.Clone(entries)
	for i := range out {
		out[i].Facets = slices.Clone(out[i].Facets)
		for j, facet := range out[i].Facets {
			var compact bytes.Buffer
			if err := json.Compact(&compact, facet.Value); err != nil {
				t.Fatal(err)
			}
			out[i].Facets[j].Value = compact.Bytes()
		}
	}
	return out
}

// gateEncodedFacet is the value of the fixture facet kind that crosses into
// the serialized plan.
type gateEncodedFacet struct {
	Module string `json:"module"`
}

// Rows 18, 19 and 30 (O36, O41), gate line 3: "The plan, read from its
// serialized form in a fresh process, is the plan that was written."
//
// Shows: the same request planned twice, each in a fresh registry with
// freshly made contributions and a freshly made facet kind, gives the same
// serialized bytes and the same digest. So the plan holds no identity of the
// planning process: an issuer of a facet kind is one, and the two plans have
// two. The parsed plan serializes to the same bytes again, and its nodes, use
// entries, symbols, roots, contributions, products, target, source and
// profile equal the written plan's. A facet whose kind states a transport is
// in the plan as the kind's identity and version, the media type and the
// encoded bytes. A reader refuses bytes that are no plan, a node whose
// identity is not the digest of its content, a node that the plan holds
// twice, a link, an application, a use entry and a symbol that name a node
// which the plan does not hold.
// Does not show: that the bytes of a facet cross unchanged. They are compared
// as JSON values here (finding F-B4, TestGateFacetBytesCrossUnread). A fresh
// process. The reader here is ParsePlan in this
// process, on the bytes alone. The completeness test from the serialized plan
// in a process that never held the planner's objects is another agent's
// (line 3, V3). The layout and the encoding, which are provisional (#195).
// Red under: S-plan-holds-a-process-identity, S-parse-skips-node-identity,
// S-parse-admits-node-twice, S-parse-skips-links,
// S-parse-skips-application-declaration, S-parse-skips-application-arguments,
// S-parse-skips-use-entries, S-parse-skips-symbols,
// S-facet-transport-not-stored.
func TestGateSerializedPlanIsThePlanThatWasWritten(t *testing.T) {
	from, input := gateSource(t, gateFacetModule)
	plan := func() g.RepresentationPlan {
		kind := gateMust[g.FacetKind](t)(g.NewFacetKind(g.FacetKindDescription{Identity: gatePublisher + "/facet/factory", Version: "1",
			Transport: &g.FacetTransport{MediaType: "application/json", Encode: func(value any) ([]byte, error) { return json.Marshal(value) }}}))
		b := gateGolang(t, from)
		b.products = []g.Product{{Identity: gatePublisher + "/conversion", Requires: map[g.SubjectKind][]g.FacetKind{g.UseSubject: {kind}}}}
		for _, key := range []string{"horos-to-golang/text", "horos-to-golang/list", "horos-to-golang/record"} {
			identity := b.find(t, key).Reference.Identity
			b = b.withFacets(gateMust[g.Facet](t)(kind.Issue(identity.Publisher, identity, goTarget, gateEncodedFacet{Module: key})))
		}
		generation, err := gatePlanGo(t, input, b, gateRequest{products: []string{gatePublisher + "/conversion"}, configuration: []g.Setting{{Name: "example.test/gate/setting", Value: "on"}}})
		if err != nil {
			t.Fatal(err)
		}
		return generation.Representation()
	}
	one, two := plan(), plan()
	first, second := gateMust[[]byte](t)(one.Serialize()), gateMust[[]byte](t)(two.Serialize())
	if !bytes.Equal(first, second) {
		t.Errorf("the same request planned in two fresh registries gives two serialized plans")
	}
	if a, b := gateMust[string](t)(one.Digest()), gateMust[string](t)(two.Digest()); a != b || a == "" {
		t.Errorf("the two plans have the digests %q and %q", a, b)
	}
	parsed, err := g.ParsePlan(first)
	if err != nil {
		t.Fatal(err)
	}
	if again := gateMust[[]byte](t)(parsed.Serialize()); !bytes.Equal(again, first) {
		t.Errorf("the parsed plan does not serialize to the bytes it was parsed from")
	}
	for name, pair := range map[string][2]any{
		"nodes": {parsed.Nodes(), one.Nodes()}, "use entries": {parsed.Uses(), one.Uses()}, "symbols": {parsed.Symbols(), one.Symbols()}, "roots": {parsed.Roots(), one.Roots()},
		"contributions": {gateCompactFacets(t, parsed.Contributions()), gateCompactFacets(t, one.Contributions())}, "products": {parsed.Products(), one.Products()}, "target": {parsed.Target(), one.Target()},
		"source": {parsed.Source(), one.Source()}, "profile": {parsed.Profile(), one.Profile()},
	} {
		if !reflect.DeepEqual(pair[0], pair[1]) {
			t.Errorf("the %s of the parsed plan are not those of the written plan", name)
		}
	}
	if len(one.Nodes()) == 0 || len(one.Uses()) == 0 || len(one.Symbols()) == 0 || len(one.Roots()) == 0 || len(one.Target().Configuration) != 1 {
		t.Fatalf("the plan is empty in a part: %d nodes, %d use entries, %d symbols, %d roots", len(one.Nodes()), len(one.Uses()), len(one.Symbols()), len(one.Roots()))
	}
	for _, node := range parsed.Nodes() {
		if again, ok := parsed.Node(node.ID); !ok || !reflect.DeepEqual(again, node) {
			t.Errorf("the parsed plan does not find its node %s by its identity", node.ID)
		}
	}
	// The facet that crossed: the kind's identity and version, the media type
	// and the bytes that the kind's owner encoded.
	list := gateGolang(t, from).find(t, "horos-to-golang/list").Reference
	entry, _ := parsed.Contribution(list)
	if len(entry.Facets) != 1 || entry.Facets[0].Kind != (g.FacetKindReference{Identity: gatePublisher + "/facet/factory", Version: "1"}) || entry.Facets[0].MediaType != "application/json" {
		t.Fatalf("the list contribution's entry holds the facets %+v", entry.Facets)
	}
	var decoded gateEncodedFacet
	if err := json.Unmarshal(entry.Facets[0].Value, &decoded); err != nil || decoded.Module != "horos-to-golang/list" {
		t.Errorf("the facet of the list crossed as %s: %v", entry.Facets[0].Value, err)
	}
	// What a reader refuses.
	if _, err := g.ParsePlan([]byte("no plan")); !errors.Is(err, g.ErrInvalid) {
		t.Errorf("bytes that are no plan are read: %v", err)
	}
	people := gateBody(t, one, "People")
	page := gateDeclaration(t, one, "Page")
	if people.Application == nil || len(people.Application.Target) != 1 || people.Application.Declaration != page.ID {
		t.Fatalf("People is not the application Page<Person>: %+v", people)
	}
	argument := people.Application.Target[0].Node
	record := gateBody(t, one, "Page")
	if len(record.Links) != 1 {
		t.Fatalf("the record of Page has the links %+v", record.Links)
	}
	withoutNode := func(id string) func(gatePlanJSON) {
		return func(p gatePlanJSON) {
			nodes := p.list("nodes")
			p["nodes"] = slices.DeleteFunc(nodes, func(n any) bool { return n.(map[string]any)["id"] == id })
			if len(p.list("nodes")) != len(nodes)-1 {
				t.Fatalf("the plan does not hold the node %s once", id)
			}
		}
	}
	// keep leaves only the nodes, so that a case meets the check it is about
	// and no later one.
	keepNodes := func(change func(gatePlanJSON)) func(gatePlanJSON) {
		return func(p gatePlanJSON) {
			change(p)
			p["uses"], p["symbols"] = []any{}, []any{}
		}
	}
	for _, c := range []struct {
		name   string
		change func(gatePlanJSON)
		want   error
		text   string
	}{
		{"a node whose content was changed", func(p gatePlanJSON) {
			p.list("nodes")[0].(map[string]any)["source"] = "changed"
		}, g.ErrInvalid, "is not identified by its content"},
		{"a node that the plan holds twice", func(p gatePlanJSON) {
			p["nodes"] = append(p.list("nodes"), p.list("nodes")[0])
		}, g.ErrInvalid, "twice"},
		{"a link to a node that the plan does not hold", keepNodes(withoutNode(record.Links[0].Node)), g.ErrIncomplete, "links to a node the plan does not hold"},
		{"an application without its declaration", keepNodes(withoutNode(page.ID)), g.ErrIncomplete, "without its declaration"},
		{"an application without an argument", keepNodes(withoutNode(argument)), g.ErrIncomplete, "without an argument"},
		{"a use entry without its node", func(p gatePlanJSON) {
			p["uses"] = append(p.list("uses"), map[string]any{"id": 9999, "node": "no such node", "site": map[string]any{"tree": "canonical", "node": 1}})
		}, g.ErrIncomplete, "use entry 9999 without its node"},
		{"a symbol without its node", func(p gatePlanJSON) {
			p.list("symbols")[0].(map[string]any)["node"] = "no such node"
		}, g.ErrIncomplete, "without its node"},
	} {
		tampered := gateDecodePlan(t, first)
		c.change(tampered)
		_, err := g.ParsePlan(tampered.encode(t))
		if !errors.Is(err, c.want) || !strings.Contains(err.Error(), c.text) {
			t.Errorf("%s: got %v, want %v with %q", c.name, err, c.want, c.text)
		}
	}
	// The control of the tampering: the plan decoded and encoded again without
	// a change is read.
	if _, err := g.ParsePlan(gateDecodePlan(t, first).encode(t)); err != nil {
		t.Errorf("the plan encoded again without a change is not read: %v", err)
	}
}

// Row 26a (O36, O41), gate line 3: the planner stores, for a facet whose kind
// states a transport, "the kind's identity and contract version, the media
// type and the bytes, and reads none of it".
//
// Shows: the bytes that a kind's encoder returned are in the parsed plan byte
// for byte, for an encoder that returns JSON with its own white space, and a
// plan with a facet whose encoder returns bytes that are no JSON serializes
// and parses, with those bytes.
// Does not show: a reader of the facet.
// Red under: no mutant. The test is skipped as PENDING: it is red against the
// kernel as built (finding F-B4), so no mutant can show it red.
func TestGateFacetBytesCrossUnread(t *testing.T) {
	t.Skip("PENDING F-B4: the plan stores the transport bytes of a facet as a raw JSON value; the JSON encoder rewrites their white space, and bytes that are no JSON fail the serialization of the plan")
	from, input := gateSource(t, "module example.transport;\n\ntype Names = list<text>;\n")
	for _, c := range []struct {
		name, media string
		encoded     []byte
	}{
		{"JSON with white space of its own", "application/json", []byte(`{ "module" :  "m" }`)},
		{"bytes that are no JSON", "application/octet-stream", []byte{0x00, 0xff, 'b', 'i', 'n'}},
	} {
		kind := gateMust[g.FacetKind](t)(g.NewFacetKind(g.FacetKindDescription{Identity: gatePublisher + "/facet/opaque", Version: "1",
			Transport: &g.FacetTransport{MediaType: c.media, Encode: func(any) ([]byte, error) { return slices.Clone(c.encoded), nil }}}))
		b := gateGolang(t, from)
		b.products = []g.Product{{Identity: gatePublisher + "/opaque", Requires: map[g.SubjectKind][]g.FacetKind{g.UseSubject: {kind}}}}
		for _, key := range []string{"horos-to-golang/text", "horos-to-golang/list"} {
			identity := b.find(t, key).Reference.Identity
			b = b.withFacets(gateMust[g.Facet](t)(kind.Issue(identity.Publisher, identity, goTarget, "a value")))
		}
		generation, err := gatePlanGo(t, input, b, gateRequest{products: []string{gatePublisher + "/opaque"}})
		if err != nil {
			t.Errorf("%s: the request does not plan: %v", c.name, err)
			continue
		}
		serialized, err := generation.Representation().Serialize()
		if err != nil {
			t.Errorf("%s: the plan does not serialize: %v", c.name, err)
			continue
		}
		parsed, err := g.ParsePlan(serialized)
		if err != nil {
			t.Errorf("%s: the serialized plan is not read: %v", c.name, err)
			continue
		}
		entry, _ := parsed.Contribution(b.find(t, "horos-to-golang/list").Reference)
		if len(entry.Facets) != 1 || entry.Facets[0].MediaType != c.media || !bytes.Equal(entry.Facets[0].Value, c.encoded) {
			t.Errorf("%s: the facets crossed as %+v, want one of the media type %s with the bytes %q", c.name, entry.Facets, c.media, c.encoded)
		}
	}
}
