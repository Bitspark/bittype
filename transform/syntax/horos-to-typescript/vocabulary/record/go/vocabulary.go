// Package record is the vocabulary of the record contribution for typescript
// (prototype, bittype#196, O25). It holds the contribution's identity and the
// payload schema with the plain Go type that bittype generates from it. The
// mapping package, its codec and its facet packages import this package, and
// none of them imports another.
//
// The payload records the names of the members that the record lays out itself, in order. A delegated field's name is stated in its request (O3) and is not restated here.
//
// payload.horos is the schema source. payload.ontos is its canonical module,
// payload_schema.go its pin and payload_types.go the generated type: the
// pinned bootstrap artifact of row 38. The regeneration test of the
// integration package fails when one of them differs from what the source
// generates.
package record

import (
	_ "embed"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/transform/syntax/shared/publisher/go"
)

// Reference is the contribution reference: the identity, the family's meaning,
// which is the same in every target, and this contribution's own
// representation revision.
var Reference = g.ContributionReference{
	Identity:       g.ContributionIdentity{Publisher: publisher.Name, Key: "horos-to-typescript/record"},
	Meaning:        g.Meaning{Identity: "bittype.structural/closed-record", Revision: "1"},
	Representation: "typescript/1",
}

// Schema is the reference of the payload schema. The schema belongs to this
// contribution's contract for this target (O21).
var Schema = g.SchemaReference{Identity: "bittype.payload/horos-to-typescript/record", Version: "1"}

// SchemaType is the name of the payload type in the schema module.
const SchemaType = "Payload"

// SchemaSource is the schema source and SchemaModule its canonical module.
//
//go:embed payload.horos
var SchemaSource string

//go:embed payload.ontos
var SchemaModule []byte
