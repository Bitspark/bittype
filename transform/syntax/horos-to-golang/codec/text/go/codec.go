// Package text is the payload codec of the text contribution for golang
// (prototype, bittype#196, row 28, O9, O18). It maps the generated plain Go
// type of the payload schema to an ontos value and back, by hand. The shared
// part of the codec role validates, encodes and freezes. A consumer's
// composition registers this codec beside the contribution.
package text

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/text/go"
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

func encode(vocabulary.Payload) (payload.Value, error) { return payload.Record() }

func decode(v payload.Value) (vocabulary.Payload, error) {
	_, err := payload.ReadRecord(v)
	return vocabulary.Payload{}, err
}
