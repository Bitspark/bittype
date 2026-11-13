// Package payload is the shared part of the payload codec role (prototype,
// bittype#196, row 28, O9, O18).
//
// The division the prototype fixes. The shared part loads a pinned payload
// schema, validates an encoded value against it by horos conformance, encodes
// the canonical bytes and decodes them. It knows no family. The codec of one
// schema, a package of its family, maps the family's generated plain Go type
// to an ontos value and back, and nothing else. So a rule's value is frozen by
// its family's codec, and every family's codec validates and encodes through
// one mechanism.
//
// A payload schema is a horos type whatever the source language of a
// generation is, so this package is no part of the horos source reader. It is
// the one shared generation package that imports horos and ontos; the
// contract package imports neither.
package payload

import (
	"fmt"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	horos "github.com/Bitspark/horos/core/go"
	"github.com/Bitspark/horos/core/go/declarations"
	codec "github.com/bitspark/ontos/codec/go"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// Value is an ontos value. A family's codec builds and reads one through the
// functions of this package and imports neither library itself.
type Value = core.Value

// Schema is one pinned payload schema: its reference and the ground
// descriptor of the payload type.
type Schema struct {
	reference  g.SchemaReference
	descriptor core.Value
}

func (s Schema) Reference() g.SchemaReference { return s.reference }

// Load reads the canonical module of a payload schema, checks it against its
// pin and resolves the ground descriptor of the type named typeName. A module
// with another pin is refused: the committed schema and the committed types
// then no longer belong together.
func Load(reference g.SchemaReference, module []byte, pin, typeName string) (Schema, error) {
	fail := func(format string, args ...any) (Schema, error) {
		return Schema{}, fmt.Errorf("%w: payload schema %s/%s: %s", g.ErrInvalid, reference.Identity, reference.Version, fmt.Sprintf(format, args...))
	}
	if !reference.Valid() {
		return fail("no reference")
	}
	value, err := codec.Decode(module)
	if err != nil {
		return fail("the canonical module does not decode: %v", err)
	}
	digest, err := declarations.Pin(value)
	if err != nil {
		return fail("the canonical module has no pin: %v", err)
	}
	if digest.String() != pin {
		return fail("the canonical module has pin %s and the schema is pinned to %s", digest, pin)
	}
	program, err := declarations.Resolve(value, nil)
	if err != nil {
		return fail("the module does not resolve: %v", err)
	}
	typed, err := program.Value(declarations.DeclarationID{Module: program.Root(), Name: typeName})
	if err != nil {
		return fail("no type %s: %v", typeName, err)
	}
	descriptor, err := program.GroundDescriptor(typed)
	if err != nil {
		return fail("type %s has no ground descriptor: %v", typeName, err)
	}
	if !horos.WF(descriptor) {
		return fail("the descriptor of %s is not well formed", typeName)
	}
	return Schema{reference, descriptor}, nil
}

// MustLoad is Load for a codec package whose schema is committed beside it.
func MustLoad(reference g.SchemaReference, module []byte, pin, typeName string) Schema {
	s, err := Load(reference, module, pin, typeName)
	if err != nil {
		panic(err)
	}
	return s
}

type codecOf[P any] struct {
	schema Schema
	encode func(P) (Value, error)
	decode func(Value) (P, error)
}

// New makes the codec of one schema from a family's two mappings between its
// generated plain Go type P and an ontos value.
func New[P any](schema Schema, encode func(P) (Value, error), decode func(Value) (P, error)) g.PayloadCodec {
	return codecOf[P]{schema, encode, decode}
}

func (c codecOf[P]) Schema() g.SchemaReference { return c.schema.reference }

// Freeze encodes the value a rule returned, validates the encoding against the
// schema and returns the canonical bytes. The bytes are the frozen value: they
// share nothing with what the rule holds.
func (c codecOf[P]) Freeze(value any) ([]byte, error) {
	typed, ok := value.(P)
	if !ok {
		pointer, isPointer := value.(*P)
		if !isPointer || pointer == nil {
			return nil, fmt.Errorf("the rule returned a %T and schema %s/%s has the type %T", value, c.schema.reference.Identity, c.schema.reference.Version, typed)
		}
		typed = *pointer
	}
	encoded, err := c.encode(typed)
	if err != nil {
		return nil, err
	}
	if !horos.Conforms(encoded, c.schema.descriptor) {
		return nil, fmt.Errorf("the value does not conform to schema %s/%s", c.schema.reference.Identity, c.schema.reference.Version)
	}
	return codec.Encode(encoded), nil
}

// View decodes frozen bytes into a fresh typed value.
func (c codecOf[P]) View(frozen []byte) (any, error) {
	value, err := codec.Decode(frozen)
	if err != nil {
		return nil, err
	}
	if !horos.Conforms(value, c.schema.descriptor) {
		return nil, fmt.Errorf("the frozen value does not conform to schema %s/%s", c.schema.reference.Identity, c.schema.reference.Version)
	}
	return c.decode(value)
}

// The value forms a hand-written codec needs. They are the canonical horos
// values of the families that exist: text, booleans, lists and closed records.

func Text(s string) (Value, error)      { return data.EncodeText(s) }
func ReadText(v Value) (string, error)  { return data.ReadText(v) }
func Bool(b bool) Value                 { return data.EncodeBool(b) }
func ReadBool(v Value) (bool, error)    { return data.ReadBool(v) }
func List(items []Value) Value          { return data.EncodeList(items) }
func ReadList(v Value) ([]Value, error) { return data.ReadList(v) }

// Field is one field of a record value.
type Field struct {
	Name  string
	Value Value
}

func Record(fields ...Field) (Value, error) {
	entries := make([][2]core.Value, len(fields))
	for i, f := range fields {
		name, err := data.EncodeText(f.Name)
		if err != nil {
			return nil, err
		}
		entries[i] = [2]core.Value{name, f.Value}
	}
	return data.EncodeMap(entries)
}

// ReadRecord returns the fields of a record value by name.
func ReadRecord(v Value) (map[string]Value, error) {
	entries, err := data.ReadMap(v)
	if err != nil {
		return nil, err
	}
	out := map[string]Value{}
	for _, e := range entries {
		name, err := data.ReadText(e[0])
		if err != nil {
			return nil, err
		}
		out[name] = e[1]
	}
	return out, nil
}

// Texts and ReadTexts map a list of text.
func Texts(items []string) (Value, error) {
	values := make([]Value, len(items))
	for i, s := range items {
		v, err := data.EncodeText(s)
		if err != nil {
			return nil, err
		}
		values[i] = v
	}
	return data.EncodeList(values), nil
}
func ReadTexts(v Value) ([]string, error) {
	values, err := data.ReadList(v)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(values))
	for i, item := range values {
		if out[i], err = data.ReadText(item); err != nil {
			return nil, err
		}
	}
	return out, nil
}
