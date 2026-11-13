// Package list is the payload codec of the list contribution for
// typescript (docs/contribution-contract.md, row 28, O9, O37). It holds the two
// conversions between the plain Go type of the payload schema and an ontos
// value. The shared part of the codec role loads the pinned schema, validates
// by horos conformance and encodes the canonical bytes, and it has the value
// forms that these conversions use. This package never imports horos. A
// consumer's composition registers this codec beside the contribution.
package list

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/list/go"
	"github.com/Bitspark/bittype/transform/syntax/horos/payload/go"
	core "github.com/bitspark/ontos/core/go"
)

// Codec returns the codec of the payload schema that the vocabulary pins.
func Codec() (g.PayloadCodec, error) {
	return payload.Codec(vocabulary.Schema, vocabulary.SchemaModule, vocabulary.SchemaPin, vocabulary.SchemaType, encode, decode)
}

func encode(vocabulary.Payload) (core.Value, error) { return payload.Record(nil) }

func decode(v core.Value) (vocabulary.Payload, error) {
	_, err := payload.ReadRecord(v, nil)
	return vocabulary.Payload{}, err
}
