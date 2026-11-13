// Package declarations is the payload codec of the declarations contribution for
// cpp (docs/contribution-contract.md, row 28, O9, O37). It holds the two
// conversions between the plain Go type of the payload schema and an ontos
// value. The shared part of the codec role loads the pinned schema, validates
// by horos conformance and encodes the canonical bytes, and it has the value
// forms that these conversions use. This package never imports horos. A
// consumer's composition registers this codec beside the contribution.
package declarations

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/vocabulary/declarations/go"
	"github.com/Bitspark/bittype/transform/syntax/horos/payload/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// Codec returns the codec of the payload schema that the vocabulary pins.
func Codec() (g.PayloadCodec, error) {
	return payload.Codec(vocabulary.Schema, vocabulary.SchemaModule, vocabulary.SchemaPin, vocabulary.SchemaType, encode, decode)
}

var fields = []string{"name", "nominal", "parameters"}

func encode(p vocabulary.Payload) (core.Value, error) {
	name, err := data.EncodeText(p.Name)
	if err != nil {
		return nil, err
	}
	parameters, err := payload.List(p.Parameters, data.EncodeText)
	if err != nil {
		return nil, err
	}
	return payload.Record(fields, name, data.EncodeBool(p.Nominal), parameters)
}

func decode(v core.Value) (vocabulary.Payload, error) {
	var p vocabulary.Payload
	values, err := payload.ReadRecord(v, fields)
	if err != nil {
		return p, err
	}
	if p.Name, err = data.ReadText(values[0]); err != nil {
		return p, err
	}
	if p.Nominal, err = data.ReadBool(values[1]); err != nil {
		return p, err
	}
	p.Parameters, err = payload.ReadList(values[2], data.ReadText)
	return p, err
}
