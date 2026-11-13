// Package structtype is the vocabulary of the structtype contribution of the pair
// golang-to-typescript (docs/contribution-contract.md, O25). It holds the contribution's
// reference and the payload schema with the plain Go type of that schema. The
// mapping package and its codec import this package, and neither imports the
// other.
//
// The payload records the names of the members of the target record, in the
// order of the fields of the Go struct type.
//
// payload.horos is the schema source. payload.ontos is its canonical module,
// payload_schema.go its pin and payload_types.go the Go type: the pinned
// bootstrap artifact of row 38. The command tools/payload-schema/go makes all
// three. The schema is a horos type although horos is no part of this pair: a payload's schema
// is a horos type whatever the pair is (docs/representation-families.md).
package structtype

import (
	_ "embed"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	family "github.com/Bitspark/bittype/transform/syntax/golang/family/structtype/go"
	rules "github.com/Bitspark/bittype/transform/syntax/golang/mapping/go"
)

// Reference is the contribution reference: the identity, the family's meaning
// and this contribution's own representation revision.
var Reference = g.ContributionReference{
	Identity:       g.ContributionIdentity{Publisher: rules.Publisher, Key: "golang-to-typescript/structtype"},
	Meaning:        family.Meaning,
	Representation: "typescript/1",
}

// Schema is the reference of the payload schema.
var Schema = g.SchemaReference{Identity: "bittype.payload/golang-to-typescript/structtype", Version: "1"}

// SchemaType is the name of the payload type in the schema module.
const SchemaType = "Payload"

// SchemaSource is the schema source and SchemaModule its canonical module.
//
//go:embed payload.horos
var SchemaSource string

//go:embed payload.ontos
var SchemaModule []byte
