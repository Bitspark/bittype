package integration_test

// PROTOTYPE (bittype#196): the pinned bootstrap artifact of a payload schema
// (row 38, O9): "the schema source, its canonical module, the committed
// generated types, the codec, and a test that regenerates the types and
// round-trips every value".

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/internal/horosmodel"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
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
	tsdeclarationscodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/declarations/go"
	tslistcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/list/go"
	tsoptionalcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/optional/go"
	tsrecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/record/go"
	tstextcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/text/go"
	tsdeclarationsvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/declarations/go"
	tslistvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/list/go"
	tsoptionalvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/optional/go"
	tsrecordvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/record/go"
	tstextvocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/text/go"
	"github.com/Bitspark/horos/core/go/declarations"
	codec "github.com/bitspark/ontos/codec/go"
)

// vocabularyCase is one vocabulary package with the codec of its schema and
// the values that the round trip covers. empty says that the payload is empty.
type vocabularyCase struct {
	native, name string
	reference    g.ContributionReference
	schema       g.SchemaReference
	source       string
	module       []byte
	pin          string
	codec        func() (g.PayloadCodec, error)
	values       []any
	empty        bool
}

func (c vocabularyCase) directory() string {
	return filepath.Join("..", "..", "transform", "syntax", "horos-to-"+c.native, "vocabulary", c.name, "go")
}

func vocabularyCases() []vocabularyCase {
	records := func(make func([]string) any) []any {
		return []any{make(nil), make([]string{"Items"}), make([]string{"Name", "Größe", ""})}
	}
	return []vocabularyCase{
		{"golang", "text", gotextvocabulary.Reference, gotextvocabulary.Schema, gotextvocabulary.SchemaSource, gotextvocabulary.SchemaModule, gotextvocabulary.SchemaPin, gotextcodec.Codec, []any{gotextvocabulary.Payload{}}, true},
		{"golang", "boolean", gobooleanvocabulary.Reference, gobooleanvocabulary.Schema, gobooleanvocabulary.SchemaSource, gobooleanvocabulary.SchemaModule, gobooleanvocabulary.SchemaPin, gobooleancodec.Codec, []any{gobooleanvocabulary.Payload{}}, true},
		{"golang", "list", golistvocabulary.Reference, golistvocabulary.Schema, golistvocabulary.SchemaSource, golistvocabulary.SchemaModule, golistvocabulary.SchemaPin, golistcodec.Codec, []any{golistvocabulary.Payload{}}, true},
		{"golang", "record", gorecordvocabulary.Reference, gorecordvocabulary.Schema, gorecordvocabulary.SchemaSource, gorecordvocabulary.SchemaModule, gorecordvocabulary.SchemaPin, gorecordcodec.Codec, records(func(m []string) any { return gorecordvocabulary.Payload{Members: m} }), false},
		{"golang", "declarations", godeclarationsvocabulary.Reference, godeclarationsvocabulary.Schema, godeclarationsvocabulary.SchemaSource, godeclarationsvocabulary.SchemaModule, godeclarationsvocabulary.SchemaPin, godeclarationscodec.Codec, []any{godeclarationsvocabulary.Payload{Name: "Page", Nominal: true, Parameters: []string{"T"}}, godeclarationsvocabulary.Payload{Name: "Names"}}, false},
		{"golang", "optional", gooptionalvocabulary.Reference, gooptionalvocabulary.Schema, gooptionalvocabulary.SchemaSource, gooptionalvocabulary.SchemaModule, gooptionalvocabulary.SchemaPin, gooptionalcodec.Codec, []any{gooptionalvocabulary.Payload{}}, true},
		{"typescript", "text", tstextvocabulary.Reference, tstextvocabulary.Schema, tstextvocabulary.SchemaSource, tstextvocabulary.SchemaModule, tstextvocabulary.SchemaPin, tstextcodec.Codec, []any{tstextvocabulary.Payload{}}, true},
		{"typescript", "list", tslistvocabulary.Reference, tslistvocabulary.Schema, tslistvocabulary.SchemaSource, tslistvocabulary.SchemaModule, tslistvocabulary.SchemaPin, tslistcodec.Codec, []any{tslistvocabulary.Payload{}}, true},
		{"typescript", "record", tsrecordvocabulary.Reference, tsrecordvocabulary.Schema, tsrecordvocabulary.SchemaSource, tsrecordvocabulary.SchemaModule, tsrecordvocabulary.SchemaPin, tsrecordcodec.Codec, records(func(m []string) any { return tsrecordvocabulary.Payload{Members: m} }), false},
		{"typescript", "declarations", tsdeclarationsvocabulary.Reference, tsdeclarationsvocabulary.Schema, tsdeclarationsvocabulary.SchemaSource, tsdeclarationsvocabulary.SchemaModule, tsdeclarationsvocabulary.SchemaPin, tsdeclarationscodec.Codec, []any{tsdeclarationsvocabulary.Payload{Name: "Page", Nominal: true, Parameters: []string{"T"}}, tsdeclarationsvocabulary.Payload{Name: "Names"}}, false},
		{"typescript", "optional", tsoptionalvocabulary.Reference, tsoptionalvocabulary.Schema, tsoptionalvocabulary.SchemaSource, tsoptionalvocabulary.SchemaModule, tsoptionalvocabulary.SchemaPin, tsoptionalcodec.Codec, []any{tsoptionalvocabulary.Payload{}}, true},
	}
}

// protoSchemaModule lowers a schema source to its canonical module and pin.
func protoSchemaModule(t *testing.T, source string) ([]byte, string) {
	t.Helper()
	parsed, err := parsehoros.Parse(context.Background(), parse.Input{SourceID: "payload.horos", TreeID: "canonical", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	modules, err := horosmodel.SourceModules(parsed.Abstract)
	if err != nil || len(modules) != 1 {
		t.Fatal("the schema source is not one module", err)
	}
	pin, err := declarations.Pin(modules[0].Value)
	if err != nil {
		t.Fatal(err)
	}
	return codec.Encode(modules[0].Value), pin.String()
}

// protoGeneratedTypes generates the Go types of a schema with the Go
// contributions of the prototype, as package name.
func protoGeneratedTypes(t *testing.T, source, name string) string {
	t.Helper()
	from, input := protoSource(t, source)
	plan, err := protoPlanGo(t, input, golangBundle(t, from), protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	text := strings.Replace(protoText[gs.Language](t, golangText(), plan), "package generated", "// Code generated by bittype from payload.horos. DO NOT EDIT.\n\npackage "+name, 1)
	formatted, err := format.Source([]byte(text))
	if err != nil {
		t.Fatal(err, text)
	}
	return string(formatted)
}

// Row 38 and O9.
//
// Shows, for each of the eleven vocabulary packages the prototype has: the
// committed canonical module and pin are what the schema source lowers to, the
// committed Go type is what bittype generates from the schema with the
// prototype's own Go contributions, and every listed value round-trips through
// the hand-written codec.
// Does not show: TypeScript types of a payload. The prototype generates none:
// the stand-in reads a payload as data. It does not show a schema that the
// existing families cannot type; no payload of the moved contributions needs one.
//
// With BITTYPE_PROTOTYPE_UPDATE=1 it writes the artifact instead of comparing:
// run it twice, once for the modules and pins and once for the types.
func TestPrototypePayloadBootstrap(t *testing.T) {
	update := os.Getenv("BITTYPE_PROTOTYPE_UPDATE") == "1"
	cases := vocabularyCases()
	loaded := true
	for _, c := range cases {
		if _, err := c.codec(); err != nil {
			loaded = false
		}
	}
	for _, c := range cases {
		t.Run(c.native+"/"+c.name, func(t *testing.T) {
			module, pin := protoSchemaModule(t, c.source)
			pinFile := "// Code generated by the regeneration test of the prototype. DO NOT EDIT.\n\npackage " + c.name + "\n\n// SchemaPin is the pin of the canonical schema module.\nconst SchemaPin = \"" + pin + "\"\n"
			if update {
				if err := os.WriteFile(filepath.Join(c.directory(), "payload.ontos"), module, 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(c.directory(), "payload_schema.go"), []byte(pinFile), 0o644); err != nil {
					t.Fatal(err)
				}
			} else if !bytes.Equal(module, c.module) || pin != c.pin {
				t.Errorf("the committed module or pin is not what payload.horos lowers to: pin %s, committed %s", pin, c.pin)
			}
			if !loaded {
				t.Skip("the pins were written in this run; run again for the types")
			}
			types := protoGeneratedTypes(t, c.source, c.name)
			committed, err := os.ReadFile(filepath.Join(c.directory(), "payload_types.go"))
			if err != nil {
				t.Fatal(err)
			}
			if update {
				if err := os.WriteFile(filepath.Join(c.directory(), "payload_types.go"), []byte(types), 0o644); err != nil {
					t.Fatal(err)
				}
			} else if types != strings.ReplaceAll(string(committed), "\r\n", "\n") {
				t.Errorf("the committed types differ from what the schema generates:\n%s\ncommitted:\n%s", types, committed)
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
				if fmt.Sprintf("%#v", view) != fmt.Sprintf("%#v", normalize(value)) || reflect.TypeOf(view) != reflect.TypeOf(value) {
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

// normalize gives a value the form a decoded value has: an absent list is an
// empty list.
func normalize(value any) any {
	v := reflect.New(reflect.TypeOf(value)).Elem()
	v.Set(reflect.ValueOf(value))
	for i := 0; i < v.NumField(); i++ {
		if f := v.Field(i); f.Kind() == reflect.Slice && f.IsNil() {
			f.Set(reflect.MakeSlice(f.Type(), 0, 0))
		}
	}
	return v.Interface()
}
