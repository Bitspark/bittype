package integration_test

// The pinned bootstrap artifact of every payload schema
// (docs/contribution-contract.md, row 38, O9): "the schema source, its
// canonical module, the committed generated types, the codec, and a test that
// regenerates the types and round-trips every value". This file is written by
// a script of the migration and lists every vocabulary package with its codec.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/tools/payload-schema/artifact/go"
	cppbooleancodec "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/codec/boolean/go"
	cppdeclarationscodec "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/codec/declarations/go"
	cpplistcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/codec/list/go"
	cppoptionalcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/codec/optional/go"
	cpprecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/codec/record/go"
	cpptextcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/codec/text/go"
	cppbooleanvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/vocabulary/boolean/go"
	cppdeclarationsvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/vocabulary/declarations/go"
	cpplistvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/vocabulary/list/go"
	cppoptionalvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/vocabulary/optional/go"
	cpprecordvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/vocabulary/record/go"
	cpptextvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/vocabulary/text/go"
	gobooleancodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/boolean/go"
	godeclarationscodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/declarations/go"
	golistcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/list/go"
	gooptionalcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/optional/go"
	gorecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/record/go"
	gotextcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/text/go"
	gobooleanvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/boolean/go"
	godeclarationsvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/declarations/go"
	golistvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/list/go"
	gooptionalvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/optional/go"
	gorecordvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/record/go"
	gotextvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/text/go"
	rustbooleancodec "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/codec/boolean/go"
	rustdeclarationscodec "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/codec/declarations/go"
	rustlistcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/codec/list/go"
	rustoptionalcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/codec/optional/go"
	rustrecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/codec/record/go"
	rusttextcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/codec/text/go"
	rustbooleanvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/vocabulary/boolean/go"
	rustdeclarationsvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/vocabulary/declarations/go"
	rustlistvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/vocabulary/list/go"
	rustoptionalvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/vocabulary/optional/go"
	rustrecordvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/vocabulary/record/go"
	rusttextvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/vocabulary/text/go"
	tsbooleancodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/boolean/go"
	tsdeclarationscodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/declarations/go"
	tslistcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/list/go"
	tsoptionalcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/optional/go"
	tsrecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/record/go"
	tstextcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/text/go"
	tsbooleanvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/boolean/go"
	tsdeclarationsvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/declarations/go"
	tslistvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/list/go"
	tsoptionalvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/optional/go"
	tsrecordvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/record/go"
	tstextvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/text/go"
	booleanfamily "github.com/Bitspark/bittype/transform/syntax/horos/family/boolean/go"
	declarationsfamily "github.com/Bitspark/bittype/transform/syntax/horos/family/declarations/go"
	listfamily "github.com/Bitspark/bittype/transform/syntax/horos/family/list/go"
	optionalfamily "github.com/Bitspark/bittype/transform/syntax/horos/family/optional/go"
	recordfamily "github.com/Bitspark/bittype/transform/syntax/horos/family/record/go"
	textfamily "github.com/Bitspark/bittype/transform/syntax/horos/family/text/go"
	"github.com/Bitspark/bittype/transform/syntax/horos/payload/go"
	"github.com/Bitspark/bittype/transform/syntax/shared/publisher/go"
)

// payloadCase is one vocabulary package with the codec of its schema and the
// values that the round trip covers.
type payloadCase struct {
	native, name string
	reference    g.ContributionReference
	schema       g.SchemaReference
	source       string
	module       []byte
	pin          string
	typeName     string
	codec        func() (g.PayloadCodec, error)
	values       []any
}

func payloadCases() []payloadCase {
	records := func(make func([]string) any) []any {
		return []any{make(nil), make([]string{"Items"}), make([]string{"Name", "Größe", ""})}
	}
	return []payloadCase{
		{"golang", "text", gotextvocabulary.Reference, gotextvocabulary.Schema, gotextvocabulary.SchemaSource, gotextvocabulary.SchemaModule, gotextvocabulary.SchemaPin, gotextvocabulary.SchemaType, gotextcodec.Codec, []any{gotextvocabulary.Payload{}}},
		{"golang", "boolean", gobooleanvocabulary.Reference, gobooleanvocabulary.Schema, gobooleanvocabulary.SchemaSource, gobooleanvocabulary.SchemaModule, gobooleanvocabulary.SchemaPin, gobooleanvocabulary.SchemaType, gobooleancodec.Codec, []any{gobooleanvocabulary.Payload{}}},
		{"golang", "list", golistvocabulary.Reference, golistvocabulary.Schema, golistvocabulary.SchemaSource, golistvocabulary.SchemaModule, golistvocabulary.SchemaPin, golistvocabulary.SchemaType, golistcodec.Codec, []any{golistvocabulary.Payload{}}},
		{"golang", "record", gorecordvocabulary.Reference, gorecordvocabulary.Schema, gorecordvocabulary.SchemaSource, gorecordvocabulary.SchemaModule, gorecordvocabulary.SchemaPin, gorecordvocabulary.SchemaType, gorecordcodec.Codec, records(func(m []string) any { return gorecordvocabulary.Payload{Members: m} })},
		{"golang", "optional", gooptionalvocabulary.Reference, gooptionalvocabulary.Schema, gooptionalvocabulary.SchemaSource, gooptionalvocabulary.SchemaModule, gooptionalvocabulary.SchemaPin, gooptionalvocabulary.SchemaType, gooptionalcodec.Codec, []any{gooptionalvocabulary.Payload{}}},
		{"golang", "declarations", godeclarationsvocabulary.Reference, godeclarationsvocabulary.Schema, godeclarationsvocabulary.SchemaSource, godeclarationsvocabulary.SchemaModule, godeclarationsvocabulary.SchemaPin, godeclarationsvocabulary.SchemaType, godeclarationscodec.Codec, []any{godeclarationsvocabulary.Payload{Name: "Page", Nominal: true, Parameters: []string{"T"}}, godeclarationsvocabulary.Payload{Name: "Names"}}},
		{"typescript", "text", tstextvocabulary.Reference, tstextvocabulary.Schema, tstextvocabulary.SchemaSource, tstextvocabulary.SchemaModule, tstextvocabulary.SchemaPin, tstextvocabulary.SchemaType, tstextcodec.Codec, []any{tstextvocabulary.Payload{}}},
		{"typescript", "boolean", tsbooleanvocabulary.Reference, tsbooleanvocabulary.Schema, tsbooleanvocabulary.SchemaSource, tsbooleanvocabulary.SchemaModule, tsbooleanvocabulary.SchemaPin, tsbooleanvocabulary.SchemaType, tsbooleancodec.Codec, []any{tsbooleanvocabulary.Payload{}}},
		{"typescript", "list", tslistvocabulary.Reference, tslistvocabulary.Schema, tslistvocabulary.SchemaSource, tslistvocabulary.SchemaModule, tslistvocabulary.SchemaPin, tslistvocabulary.SchemaType, tslistcodec.Codec, []any{tslistvocabulary.Payload{}}},
		{"typescript", "record", tsrecordvocabulary.Reference, tsrecordvocabulary.Schema, tsrecordvocabulary.SchemaSource, tsrecordvocabulary.SchemaModule, tsrecordvocabulary.SchemaPin, tsrecordvocabulary.SchemaType, tsrecordcodec.Codec, records(func(m []string) any { return tsrecordvocabulary.Payload{Members: m} })},
		{"typescript", "optional", tsoptionalvocabulary.Reference, tsoptionalvocabulary.Schema, tsoptionalvocabulary.SchemaSource, tsoptionalvocabulary.SchemaModule, tsoptionalvocabulary.SchemaPin, tsoptionalvocabulary.SchemaType, tsoptionalcodec.Codec, []any{tsoptionalvocabulary.Payload{}}},
		{"typescript", "declarations", tsdeclarationsvocabulary.Reference, tsdeclarationsvocabulary.Schema, tsdeclarationsvocabulary.SchemaSource, tsdeclarationsvocabulary.SchemaModule, tsdeclarationsvocabulary.SchemaPin, tsdeclarationsvocabulary.SchemaType, tsdeclarationscodec.Codec, []any{tsdeclarationsvocabulary.Payload{Name: "Page", Nominal: true, Parameters: []string{"T"}}, tsdeclarationsvocabulary.Payload{Name: "Names"}}},
		{"rust", "text", rusttextvocabulary.Reference, rusttextvocabulary.Schema, rusttextvocabulary.SchemaSource, rusttextvocabulary.SchemaModule, rusttextvocabulary.SchemaPin, rusttextvocabulary.SchemaType, rusttextcodec.Codec, []any{rusttextvocabulary.Payload{}}},
		{"rust", "boolean", rustbooleanvocabulary.Reference, rustbooleanvocabulary.Schema, rustbooleanvocabulary.SchemaSource, rustbooleanvocabulary.SchemaModule, rustbooleanvocabulary.SchemaPin, rustbooleanvocabulary.SchemaType, rustbooleancodec.Codec, []any{rustbooleanvocabulary.Payload{}}},
		{"rust", "list", rustlistvocabulary.Reference, rustlistvocabulary.Schema, rustlistvocabulary.SchemaSource, rustlistvocabulary.SchemaModule, rustlistvocabulary.SchemaPin, rustlistvocabulary.SchemaType, rustlistcodec.Codec, []any{rustlistvocabulary.Payload{}}},
		{"rust", "record", rustrecordvocabulary.Reference, rustrecordvocabulary.Schema, rustrecordvocabulary.SchemaSource, rustrecordvocabulary.SchemaModule, rustrecordvocabulary.SchemaPin, rustrecordvocabulary.SchemaType, rustrecordcodec.Codec, records(func(m []string) any { return rustrecordvocabulary.Payload{Members: m} })},
		{"rust", "optional", rustoptionalvocabulary.Reference, rustoptionalvocabulary.Schema, rustoptionalvocabulary.SchemaSource, rustoptionalvocabulary.SchemaModule, rustoptionalvocabulary.SchemaPin, rustoptionalvocabulary.SchemaType, rustoptionalcodec.Codec, []any{rustoptionalvocabulary.Payload{}}},
		{"rust", "declarations", rustdeclarationsvocabulary.Reference, rustdeclarationsvocabulary.Schema, rustdeclarationsvocabulary.SchemaSource, rustdeclarationsvocabulary.SchemaModule, rustdeclarationsvocabulary.SchemaPin, rustdeclarationsvocabulary.SchemaType, rustdeclarationscodec.Codec, []any{rustdeclarationsvocabulary.Payload{Name: "Page", Nominal: true, Parameters: []string{"T"}}, rustdeclarationsvocabulary.Payload{Name: "Names"}}},
		{"cpp", "text", cpptextvocabulary.Reference, cpptextvocabulary.Schema, cpptextvocabulary.SchemaSource, cpptextvocabulary.SchemaModule, cpptextvocabulary.SchemaPin, cpptextvocabulary.SchemaType, cpptextcodec.Codec, []any{cpptextvocabulary.Payload{}}},
		{"cpp", "boolean", cppbooleanvocabulary.Reference, cppbooleanvocabulary.Schema, cppbooleanvocabulary.SchemaSource, cppbooleanvocabulary.SchemaModule, cppbooleanvocabulary.SchemaPin, cppbooleanvocabulary.SchemaType, cppbooleancodec.Codec, []any{cppbooleanvocabulary.Payload{}}},
		{"cpp", "list", cpplistvocabulary.Reference, cpplistvocabulary.Schema, cpplistvocabulary.SchemaSource, cpplistvocabulary.SchemaModule, cpplistvocabulary.SchemaPin, cpplistvocabulary.SchemaType, cpplistcodec.Codec, []any{cpplistvocabulary.Payload{}}},
		{"cpp", "record", cpprecordvocabulary.Reference, cpprecordvocabulary.Schema, cpprecordvocabulary.SchemaSource, cpprecordvocabulary.SchemaModule, cpprecordvocabulary.SchemaPin, cpprecordvocabulary.SchemaType, cpprecordcodec.Codec, records(func(m []string) any { return cpprecordvocabulary.Payload{Members: m} })},
		{"cpp", "optional", cppoptionalvocabulary.Reference, cppoptionalvocabulary.Schema, cppoptionalvocabulary.SchemaSource, cppoptionalvocabulary.SchemaModule, cppoptionalvocabulary.SchemaPin, cppoptionalvocabulary.SchemaType, cppoptionalcodec.Codec, []any{cppoptionalvocabulary.Payload{}}},
		{"cpp", "declarations", cppdeclarationsvocabulary.Reference, cppdeclarationsvocabulary.Schema, cppdeclarationsvocabulary.SchemaSource, cppdeclarationsvocabulary.SchemaModule, cppdeclarationsvocabulary.SchemaPin, cppdeclarationsvocabulary.SchemaType, cppdeclarationscodec.Codec, []any{cppdeclarationsvocabulary.Payload{Name: "Page", Nominal: true, Parameters: []string{"T"}}, cppdeclarationsvocabulary.Payload{Name: "Names"}}},
	}
}

// familyMeanings is the meaning of each family, from the family's own
// vocabulary package.
func familyMeanings() map[string]g.Meaning {
	return map[string]g.Meaning{
		"text":         textfamily.Meaning,
		"boolean":      booleanfamily.Meaning,
		"list":         listfamily.Meaning,
		"record":       recordfamily.Meaning,
		"optional":     optionalfamily.Meaning,
		"declarations": declarationsfamily.Meaning,
	}
}

// Row 38 and O9, the module, the pin and the codec.
//
// Shows, for each of the 24 vocabulary packages: the committed canonical
// module and the committed pin are what the schema source lowers to; the
// codec states the schema reference of its vocabulary; every listed value
// freezes, decodes to an equal value of the same Go type and freezes to the
// same bytes a second time.
// Does not show: that the committed Go type is what bittype generates from
// the schema. TestPayloadArtifactsOfEveryPairRegenerate shows that, for every
// vocabulary package that the repository holds, found by path. It does not show a field
// that the existing families cannot type: no payload of the delivered
// contributions has one.
func TestPayloadBootstrapArtifacts(t *testing.T) {
	cases := payloadCases()
	if len(cases) != 24 {
		t.Fatalf("%d vocabulary packages, expected 24", len(cases))
	}
	for _, c := range cases {
		t.Run(c.native+"/"+c.name, func(t *testing.T) {
			module, pin, err := artifact.Module(context.Background(), c.source)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(module, c.module) {
				t.Errorf("the committed module is not what payload.horos lowers to")
			}
			if pin != c.pin {
				t.Errorf("the committed pin %s is not the pin %s of what payload.horos lowers to", c.pin, pin)
			}
			payloadCodec, err := c.codec()
			if err != nil {
				t.Fatal(err)
			}
			if payloadCodec.Schema() != c.schema {
				t.Errorf("the codec is of schema %v and the vocabulary states %v", payloadCodec.Schema(), c.schema)
			}
			for _, value := range c.values {
				frozen, err := payloadCodec.Freeze(value)
				if err != nil {
					t.Errorf("%#v does not freeze: %v", value, err)
					continue
				}
				view, err := payloadCodec.View(frozen)
				if err != nil {
					t.Errorf("%#v does not decode: %v", value, err)
					continue
				}
				if fmt.Sprintf("%#v", view) != fmt.Sprintf("%#v", normalizePayload(value)) || reflect.TypeOf(view) != reflect.TypeOf(value) {
					t.Errorf("%#v came back as %#v", value, view)
				}
				again, err := payloadCodec.Freeze(view)
				if err != nil || !bytes.Equal(again, frozen) {
					t.Errorf("%#v does not freeze to the same bytes a second time", value)
				}
			}
		})
	}
}

// normalizePayload gives a value the form a decoded value has: an absent list
// is an empty list.
func normalizePayload(value any) any {
	v := reflect.New(reflect.TypeOf(value)).Elem()
	v.Set(reflect.ValueOf(value))
	for i := 0; i < v.NumField(); i++ {
		if f := v.Field(i); f.Kind() == reflect.Slice && f.IsNil() {
			f.Set(reflect.MakeSlice(f.Type(), 0, 0))
		}
	}
	return v.Interface()
}

// Row 38: a schema and its committed artifact belong together.
//
// Shows: the shared part of the codec role refuses a canonical module under
// another pin, a module that does not decode, a type the module does not
// have and a missing schema reference, each as an invalid description.
// Does not show: who may change a pin. Nothing but the committed file binds it.
func TestPayloadSchemaIsPinned(t *testing.T) {
	c := payloadCases()[3]
	if c.name != "record" {
		t.Fatal("the case list changed")
	}
	if _, err := payload.Load(c.schema, c.module, c.pin, c.typeName); err != nil {
		t.Fatal("the committed artifact does not load:", err)
	}
	other := payloadCases()[5]
	for name, load := range map[string]func() (payload.Schema, error){
		"another pin":    func() (payload.Schema, error) { return payload.Load(c.schema, c.module, other.pin, c.typeName) },
		"another module": func() (payload.Schema, error) { return payload.Load(c.schema, other.module, c.pin, c.typeName) },
		"no module":      func() (payload.Schema, error) { return payload.Load(c.schema, []byte("no module"), c.pin, c.typeName) },
		"no such type":   func() (payload.Schema, error) { return payload.Load(c.schema, c.module, c.pin, "Absent") },
		"no reference":   func() (payload.Schema, error) { return payload.Load(g.SchemaReference{}, c.module, c.pin, c.typeName) },
		"no schema version": func() (payload.Schema, error) {
			return payload.Load(g.SchemaReference{Identity: c.schema.Identity}, c.module, c.pin, c.typeName)
		},
	} {
		if _, err := load(); !errors.Is(err, g.ErrInvalid) {
			t.Errorf("%s: loaded, or failed with another value: %v", name, err)
		}
	}
}

// Row 17, O4, O19, O27: one meaning, a representation revision per target.
//
// Shows: for each of the six names, the contributions of the four targets
// carry the meaning of the family's own vocabulary package, which is one Go
// value and one spelling; their identities differ, each names bittype as its
// publisher and has a local key of its own; each states a representation
// revision; each has a schema reference of its own; and a reference that lacks
// either part of its third axis is no valid reference.
// Does not show: who owns a family's meaning where two publishers contribute
// to it. The page leaves that unspecified (O19, O27).
func TestOneMeaningAndARepresentationRevisionPerTarget(t *testing.T) {
	meanings := familyMeanings()
	identities := map[g.ContributionIdentity]bool{}
	schemas := map[g.SchemaReference]bool{}
	seenMeanings := map[g.Meaning]string{}
	for _, c := range payloadCases() {
		meaning, ok := meanings[c.name]
		if !ok || c.reference.Meaning != meaning {
			t.Errorf("%s/%s carries the meaning %v, the family's vocabulary states %v", c.native, c.name, c.reference.Meaning, meaning)
		}
		if other, dup := seenMeanings[c.reference.Meaning]; dup && other != c.name {
			t.Errorf("the families %s and %s share the meaning %v", other, c.name, c.reference.Meaning)
		}
		seenMeanings[c.reference.Meaning] = c.name
		if !c.reference.Valid() || c.reference.Identity.Publisher != publisher.Name {
			t.Errorf("%s/%s has the reference %v", c.native, c.name, c.reference)
		}
		if identities[c.reference.Identity] {
			t.Errorf("the local key %q is used twice", c.reference.Identity.Key)
		}
		identities[c.reference.Identity] = true
		if c.reference.Representation != c.native+"/1" {
			t.Errorf("%s/%s states the representation revision %q", c.native, c.name, c.reference.Representation)
		}
		if !c.schema.Valid() || schemas[c.schema] {
			t.Errorf("%s/%s has no schema reference of its own: %v", c.native, c.name, c.schema)
		}
		schemas[c.schema] = true
		lacksMeaning, lacksRevision, lacksRepresentation := c.reference, c.reference, c.reference
		lacksMeaning.Meaning.Identity, lacksRevision.Meaning.Revision, lacksRepresentation.Representation = "", "", ""
		if lacksMeaning.Valid() || lacksRevision.Valid() || lacksRepresentation.Valid() {
			t.Errorf("a reference of %s/%s that lacks a part of its third axis is valid", c.native, c.name)
		}
	}
	if len(identities) != 24 || len(seenMeanings) != 6 {
		t.Errorf("%d identities and %d meanings, expected 24 and 6", len(identities), len(seenMeanings))
	}
}
