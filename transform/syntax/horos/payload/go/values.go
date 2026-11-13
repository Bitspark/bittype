package payload

import (
	"fmt"
	"slices"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// Codec makes the codec of one pinned payload schema: it loads the schema and
// joins it with a contribution's two conversions between its plain Go payload
// type P and an ontos value. Every codec package is this one call with its own
// vocabulary and its own two functions.
func Codec[P any](reference g.SchemaReference, module []byte, pin, typeName string, encode func(P) (core.Value, error), decode func(core.Value) (P, error)) (g.PayloadCodec, error) {
	schema, err := Load(reference, module, pin, typeName)
	if err != nil {
		return nil, err
	}
	return New(schema, encode, decode), nil
}

// The framings below are those of the schema language's own composite kinds,
// a closed record and a list, generic in their parts (O47). They are here
// once, so that no codec repeats them. A scalar comes from the ontos value
// packages directly, and a value form of one family never enters this
// package: a codec builds it with the ontos value packages itself (O37).

// Record is the value of a closed record with the given field names, in the
// order of names. A payload without fields is Record(nil).
func Record(names []string, values ...core.Value) (core.Value, error) {
	if len(names) != len(values) {
		return nil, fmt.Errorf("a record of %d fields was given %d values", len(names), len(values))
	}
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

// ReadRecord returns the values of a record's fields in the order of names. A
// record with another set of fields is refused.
func ReadRecord(v core.Value, names []string) ([]core.Value, error) {
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

// List is the value of a list: the schema language's list kind, generic in
// its element. The caller gives the conversion of one element, for a text
// element the text conversion of the ontos value package.
func List[E any](items []E, element func(E) (core.Value, error)) (core.Value, error) {
	values := make([]core.Value, len(items))
	for i, item := range items {
		value, err := element(item)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}
	return data.EncodeList(values), nil
}

// ReadList reads a list with the caller's conversion of one element.
func ReadList[E any](v core.Value, element func(core.Value) (E, error)) ([]E, error) {
	values, err := data.ReadList(v)
	if err != nil {
		return nil, err
	}
	out := make([]E, len(values))
	for i, value := range values {
		if out[i], err = element(value); err != nil {
			return nil, err
		}
	}
	return out, nil
}
