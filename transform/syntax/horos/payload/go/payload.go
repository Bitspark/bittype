// Package payload is the shared part of the payload codec role
// (docs/contribution-contract.md, row 28, O9, O18, O37, O41).
//
// A payload is an ontos value, and its schema is a horos type, whatever the
// source language of a generation is. This package is what every payload
// codec shares: it loads a pinned schema module, resolves the ground
// descriptor of the payload type, validates a value by horos conformance,
// encodes the canonical ontos bytes and decodes them. It knows no family.
//
// The codec of one schema is a package of its contribution. It maps the
// generated plain Go type of the schema to an ontos value and back, by hand,
// and imports the ontos value packages for that (O37). It reaches horos only
// through this package. So a rule's value is frozen by its own codec, and
// every codec validates and encodes through one mechanism.
//
// This package is no kernel package: the generation contract and the registry
// import neither horos nor ontos, and they do not import this package. It
// belongs to horos as the schema language of payloads, and it is no part of
// the horos source reader: a generation from another source language uses the
// same codecs.
package payload

import (
	"fmt"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	horos "github.com/Bitspark/horos/core/go"
	"github.com/Bitspark/horos/core/go/declarations"
	codec "github.com/bitspark/ontos/codec/go"
	core "github.com/bitspark/ontos/core/go"
)

// Schema is one pinned payload schema: its reference and the ground
// descriptor of the payload type.
type Schema struct {
	reference  g.SchemaReference
	descriptor core.Value
}

func (s Schema) Reference() g.SchemaReference { return s.reference }

// Pin returns the pin of a canonical schema module. A vocabulary package
// commits it beside the module, and Load compares the two.
func Pin(module []byte) (string, error) {
	value, err := codec.Decode(module)
	if err != nil {
		return "", fmt.Errorf("%w: the canonical module does not decode: %v", g.ErrInvalid, err)
	}
	digest, err := declarations.Pin(value)
	if err != nil {
		return "", fmt.Errorf("%w: the canonical module has no pin: %v", g.ErrInvalid, err)
	}
	return digest.String(), nil
}

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

type codecOf[P any] struct {
	schema Schema
	encode func(P) (core.Value, error)
	decode func(core.Value) (P, error)
}

// New makes the codec of one schema from a contribution's two mappings
// between its generated plain Go type P and an ontos value. What this package
// adds is the schema language's part: validation by horos conformance and the
// canonical bytes. The adaptation to the codec interface is the contract's
// (generation.TypedCodec).
func New[P any](schema Schema, encode func(P) (core.Value, error), decode func(core.Value) (P, error)) g.PayloadCodec {
	c := codecOf[P]{schema, encode, decode}
	return g.TypedCodec(schema.reference, c.freeze, c.view)
}

// freeze encodes the value a rule returned, validates the encoding against the
// schema and returns the canonical bytes. The bytes are the frozen value: they
// share nothing with what the rule holds.
func (c codecOf[P]) freeze(typed P) ([]byte, error) {
	encoded, err := c.encode(typed)
	if err != nil {
		return nil, err
	}
	if !horos.Conforms(encoded, c.schema.descriptor) {
		return nil, fmt.Errorf("the value does not conform to schema %s/%s", c.schema.reference.Identity, c.schema.reference.Version)
	}
	return codec.Encode(encoded), nil
}

// view decodes frozen bytes into a fresh typed value. It validates again: the
// bytes may come from a serialized plan.
func (c codecOf[P]) view(frozen []byte) (P, error) {
	var zero P
	value, err := codec.Decode(frozen)
	if err != nil {
		return zero, err
	}
	if !horos.Conforms(value, c.schema.descriptor) {
		return zero, fmt.Errorf("the frozen value does not conform to schema %s/%s", c.schema.reference.Identity, c.schema.reference.Version)
	}
	return c.decode(value)
}
