// Package declarations is the payload codec of the declarations contribution for typescript
// (prototype, bittype#196, row 28, O9, O18). It maps the generated plain Go
// type of the payload schema to an ontos value and back, by hand. The shared
// part of the codec role validates, encodes and freezes. A consumer's
// composition registers this codec beside the contribution.
package declarations

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/declarations/go"
	"github.com/Bitspark/bittype/transform/syntax/shared/payload/go"
)

// Codec returns the codec of the payload schema that the vocabulary pins.
func Codec() (g.PayloadCodec, error) {
	schema, err := payload.Load(vocabulary.Schema, vocabulary.SchemaModule, vocabulary.SchemaPin, vocabulary.SchemaType)
	if err != nil {
		return nil, err
	}
	return payload.New(schema, encode, decode), nil
}

func encode(p vocabulary.Payload) (payload.Value, error) {
	name, err := payload.Text(p.Name)
	if err != nil {
		return nil, err
	}
	parameters, err := payload.Texts(p.Parameters)
	if err != nil {
		return nil, err
	}
	return payload.Record(payload.Field{Name: "name", Value: name}, payload.Field{Name: "nominal", Value: payload.Bool(p.Nominal)}, payload.Field{Name: "parameters", Value: parameters})
}

func decode(v payload.Value) (vocabulary.Payload, error) {
	fields, err := payload.ReadRecord(v)
	if err != nil {
		return vocabulary.Payload{}, err
	}
	var p vocabulary.Payload
	if p.Name, err = payload.ReadText(fields["name"]); err != nil {
		return p, err
	}
	if p.Nominal, err = payload.ReadBool(fields["nominal"]); err != nil {
		return p, err
	}
	p.Parameters, err = payload.ReadTexts(fields["parameters"])
	return p, err
}
