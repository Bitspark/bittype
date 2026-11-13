// Package optional is the vocabulary of the optional contribution for rust
// (docs/contribution-contract.md, O25). It holds the contribution's reference
// and the payload schema with the plain Go type of that schema. The mapping
// package, its codec and its facet packages import this package, and none of
// them imports another.
//
// The payload is empty: the operand is a link of the frozen record, and the
// member's name is stated in the request that selected this contribution (O3).
//
// payload.horos is the schema source. payload.ontos is its canonical module,
// payload_schema.go its pin and payload_types.go the Go type: the pinned
// bootstrap artifact of row 38. The command tools/payload-schema/go makes the
// module and the pin, and the regeneration test of the integration package
// fails when a committed file differs from what the source gives.
package optional

import (
	_ "embed"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	family "github.com/Bitspark/bittype/transform/syntax/horos/family/optional/go"
	"github.com/Bitspark/bittype/transform/syntax/shared/publisher/go"
)

// Reference is the contribution reference: the identity, the family's meaning,
// which is the same in every target, and this contribution's own
// representation revision, which fixes the invariants of this target's
// representation.
var Reference = g.ContributionReference{
	Identity:       g.ContributionIdentity{Publisher: publisher.Name, Key: "horos-to-rust/optional"},
	Meaning:        family.Meaning,
	Representation: "rust/1",
}

// Schema is the reference of the payload schema. The schema belongs to this
// contribution's contract for this target (O21).
var Schema = g.SchemaReference{Identity: "bittype.payload/horos-to-rust/optional", Version: "1"}

// SchemaType is the name of the payload type in the schema module.
const SchemaType = "Payload"

// SchemaSource is the schema source and SchemaModule its canonical module.
//
//go:embed payload.horos
var SchemaSource string

//go:embed payload.ontos
var SchemaModule []byte
