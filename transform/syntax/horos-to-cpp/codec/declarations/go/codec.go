// Package declarations is the payload codec of the declarations contribution for
// cpp (docs/contribution-contract.md, row 28, O9, O37). It maps the
// plain Go type of the payload schema to an ontos value and back, by hand. The
// shared part of the codec role loads the pinned schema, validates by horos
// conformance and encodes the canonical bytes. This package imports the ontos
// value packages and never horos. A consumer's composition registers this
// codec beside the contribution.
package declarations

import (
	"fmt"
	"slices"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/vocabulary/declarations/go"
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

var fields = []string{"name", "nominal", "parameters"}

func encode(p vocabulary.Payload) (core.Value, error) {
	name, err := data.EncodeText(p.Name)
	if err != nil {
		return nil, err
	}
	parameters, err := texts(p.Parameters)
	if err != nil {
		return nil, err
	}
	return record(fields, name, data.EncodeBool(p.Nominal), parameters)
}

func decode(v core.Value) (vocabulary.Payload, error) {
	var p vocabulary.Payload
	values, err := readRecord(v, fields)
	if err != nil {
		return p, err
	}
	if p.Name, err = data.ReadText(values[0]); err != nil {
		return p, err
	}
	if p.Nominal, err = data.ReadBool(values[1]); err != nil {
		return p, err
	}
	p.Parameters, err = readTexts(values[2])
	return p, err
}

// texts and readTexts map a list of text.
func texts(items []string) (core.Value, error) {
	values := make([]core.Value, len(items))
	for i, item := range items {
		value, err := data.EncodeText(item)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return data.EncodeList(values), nil
}

func readTexts(v core.Value) ([]string, error) {
	values, err := data.ReadList(v)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(values))
	for i, value := range values {
		if out[i], err = data.ReadText(value); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// record and readRecord map a record value with the given field names.
func record(names []string, values ...core.Value) (core.Value, error) {
	entries := make([][2]core.Value, len(names))
	for i, name := range names {
		key, err := data.EncodeText(name)
		if err != nil {
			return nil, err
		}
		entries[i] = [2]core.Value{key, values[i]}
	}
	return data.EncodeMap(entries)
}

func readRecord(v core.Value, names []string) ([]core.Value, error) {
	entries, err := data.ReadMap(v)
	if err != nil {
		return nil, err
	}
	if len(entries) != len(names) {
		return nil, fmt.Errorf("the payload has %d fields and its schema %d", len(entries), len(names))
	}
	out := make([]core.Value, len(names))
	for _, entry := range entries {
		key, err := data.ReadText(entry[0])
		if err != nil {
			return nil, err
		}
		i := slices.Index(names, key)
		if i < 0 || out[i] != nil {
			return nil, fmt.Errorf("the payload has a field %q that its schema does not have once", key)
		}
		out[i] = entry[1]
	}
	return out, nil
}
