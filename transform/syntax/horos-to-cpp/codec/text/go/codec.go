// Package text is the payload codec of the text contribution for
// cpp (docs/contribution-contract.md, row 28, O9, O37). It maps the
// plain Go type of the payload schema to an ontos value and back, by hand. The
// shared part of the codec role loads the pinned schema, validates by horos
// conformance and encodes the canonical bytes. This package imports the ontos
// value packages and never horos. A consumer's composition registers this
// codec beside the contribution.
package text

import (
	"errors"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/vocabulary/text/go"
	"github.com/Bitspark/bittype/transform/syntax/horos/payload/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// Codec returns the codec of the payload schema that the vocabulary pins.
func Codec() (g.PayloadCodec, error) {
	schema, err := payload.Load(vocabulary.Schema, vocabulary.SchemaModule, vocabulary.SchemaPin, vocabulary.SchemaType)
	if err != nil {
		return nil, err
	}
	return payload.New(schema, encode, decode), nil
}

func encode(vocabulary.Payload) (core.Value, error) { return data.EncodeMap(nil) }

func decode(v core.Value) (vocabulary.Payload, error) {
	entries, err := data.ReadMap(v)
	if err == nil && len(entries) != 0 {
		err = errors.New("the payload of this schema has no field")
	}
	return vocabulary.Payload{}, err
}
