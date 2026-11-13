// Package objecttype is the vocabulary of the objecttype contribution of the pair
// typescript-to-horos (docs/contribution-contract.md, O25). It holds the contribution's
// reference and the payload schema with the plain Go type of that schema. The
// mapping package and its codec import this package, and neither imports the
// other.
//
// The payload records the names of the fields of the horos record, in the
// order of the properties of the TypeScript object type.
//
// payload.horos is the schema source. payload.ontos is its canonical module,
// payload_schema.go its pin and payload_types.go the Go type: the pinned
// bootstrap artifact of row 38. The command tools/payload-schema/go makes all
// three. The schema is a horos type although horos is the target of this pair and not its source: a payload's schema
// is a horos type whatever the pair is (docs/representation-families.md).
package objecttype

import (
	_ "embed"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	family "github.com/Bitspark/bittype/transform/syntax/typescript/family/objecttype/go"
	rules "github.com/Bitspark/bittype/transform/syntax/typescript/mapping/go"
)

// Reference is the contribution reference: the identity, the family's meaning
// and this contribution's own representation revision.
var Reference = g.ContributionReference{
	Identity:       g.ContributionIdentity{Publisher: rules.Publisher, Key: "typescript-to-horos/objecttype"},
	Meaning:        family.Meaning,
	Representation: "horos/1",
}

// Schema is the reference of the payload schema.
var Schema = g.SchemaReference{Identity: "bittype.payload/typescript-to-horos/objecttype", Version: "1"}

// SchemaType is the name of the payload type in the schema module.
const SchemaType = "Payload"

// SchemaSource is the schema source and SchemaModule its canonical module.
//
//go:embed payload.horos
var SchemaSource string

//go:embed payload.ontos
var SchemaModule []byte
