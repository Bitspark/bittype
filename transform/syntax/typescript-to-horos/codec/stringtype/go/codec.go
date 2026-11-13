// Package stringtype is the payload codec of the stringtype contribution of the
// pair typescript-to-horos (docs/contribution-contract.md, row 28, O9, O37). It maps the
// plain Go type of the payload schema to an ontos value and back, by hand. The
// shared part of the codec role loads the pinned schema, validates by horos
// conformance and encodes the canonical bytes. This package imports the ontos
// value packages and never horos. A consumer's composition registers this
// codec beside the contribution.
//
// COPY (AC-11, finding AG-2 of the two other pairs, bittype#196 O43): the codec of a schema without fields.
// These are the operations, invariants and failure rules of the codec packages
// of the horos pairs, with only the schema's field names different. They have
// no agnostic implementation: the shared part of the codec role exports no
// value forms, so each codec package repeats them by hand. This package is one
// more repetition. That the pair is green does not show that the mechanism
// has an owner.
package stringtype

import (
	"errors"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/transform/syntax/horos/payload/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/typescript-to-horos/vocabulary/stringtype/go"
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
