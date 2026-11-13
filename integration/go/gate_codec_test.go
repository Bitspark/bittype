package integration_test

import (
	"errors"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
)

// gateTypedPayload and gateOtherPayload are two payload types of two schema
// owners that know nothing of each other and of no schema language.
type gateTypedPayload struct{ Text string }
type gateOtherPayload struct{ Count int }

// Row 28 (O9, O47), guarantee 6's association of a value with its schema.
//
// Shows: the adaptation of two typed functions to the codec interface is one
// mechanism of the contract and knows no schema language. Instantiated with
// two payload types, each codec states its own schema reference, freezes a
// value of its type and a pointer to one to the same bytes, gives back a value
// of its type, and refuses a value of the other type, a nil pointer and a
// value of no payload type before its freeze function runs. A codec without a
// function is an invalid description.
// Does not show: validation. The functions here encode without a schema; the
// schema language's part is the shared part of the codec role
// (TestGuarantee6ExactPayloadPassedValidation).
// Red under: O47-typed-codec-ignores-pointer, O47-typed-codec-admits-any.
func TestGateTypedCodecAdaptsWithoutASchemaLanguage(t *testing.T) {
	text := g.SchemaReference{Identity: gatePublisher + "/typed/text", Version: "1"}
	count := g.SchemaReference{Identity: gatePublisher + "/typed/count", Version: "1"}
	frozenText, frozenCount := 0, 0
	textCodec := g.TypedCodec(text,
		func(p gateTypedPayload) ([]byte, error) { frozenText++; return []byte(p.Text), nil },
		func(b []byte) (gateTypedPayload, error) { return gateTypedPayload{Text: string(b)}, nil })
	countCodec := g.TypedCodec(count,
		func(p gateOtherPayload) ([]byte, error) { frozenCount++; return []byte{byte(p.Count)}, nil },
		func(b []byte) (gateOtherPayload, error) {
			if len(b) != 1 {
				return gateOtherPayload{}, errors.New("no count")
			}
			return gateOtherPayload{Count: int(b[0])}, nil
		})
	if textCodec.Schema() != text || countCodec.Schema() != count {
		t.Fatalf("the codecs state %v and %v", textCodec.Schema(), countCodec.Schema())
	}
	byValue, err := textCodec.Freeze(gateTypedPayload{Text: "a"})
	if err != nil || string(byValue) != "a" {
		t.Fatalf("a value of the payload type froze to %q: %v", byValue, err)
	}
	byPointer, err := textCodec.Freeze(&gateTypedPayload{Text: "a"})
	if err != nil || string(byPointer) != "a" {
		t.Errorf("a pointer to a value of the payload type froze to %q: %v", byPointer, err)
	}
	if view, err := textCodec.View(byValue); err != nil || view != (gateTypedPayload{Text: "a"}) {
		t.Errorf("the view is %#v: %v", view, err)
	}
	if frozen, err := countCodec.Freeze(gateOtherPayload{Count: 7}); err != nil || len(frozen) != 1 || frozen[0] != 7 {
		t.Errorf("the second codec froze %v: %v", frozen, err)
	} else if view, err := countCodec.View(frozen); err != nil || view != (gateOtherPayload{Count: 7}) {
		t.Errorf("the second codec's view is %#v: %v", view, err)
	}
	before := frozenText
	var none *gateTypedPayload
	for name, value := range map[string]any{"a value of another payload type": gateOtherPayload{Count: 1}, "a nil pointer": none, "a string": "a", "nothing": nil} {
		if frozen, err := textCodec.Freeze(value); err == nil || frozen != nil {
			t.Errorf("%s froze to %q", name, frozen)
		}
	}
	if frozenText != before {
		t.Errorf("the freeze function ran %d times for values that are not of its type", frozenText-before)
	}
	if frozenCount != 1 {
		t.Errorf("the second codec's freeze function ran %d times", frozenCount)
	}
	empty := g.TypedCodec[gateTypedPayload](text, nil, nil)
	if _, err := empty.Freeze(gateTypedPayload{}); !errors.Is(err, g.ErrInvalid) {
		t.Errorf("a codec without a freeze function: %v", err)
	}
	if _, err := empty.View(nil); !errors.Is(err, g.ErrInvalid) {
		t.Errorf("a codec without a view function: %v", err)
	}
}
