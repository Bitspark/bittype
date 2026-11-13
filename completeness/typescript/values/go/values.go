// Package values is the vocabulary of the TypeScript stand-in product
// (prototype, bittype#196, rows 26a and 31, O8). It is a test instrument: only
// test compositions and the TypeScript facet packages import it.
//
// The product is "a reader in another process converts values". Its facet
// kinds name a run-time factory by module and export: "By default a plan
// entry names the family's run-time factory and gives it the family's payload
// as arguments." A value of these kinds crosses the process boundary, so the
// kinds have a transport descriptor: the trigger of row 26a.
package values

import (
	"encoding/json"
	"fmt"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
)

// Product is the identity of the product that a request names.
const Product = "bittype.completeness/typescript-conversion"

// Runtime is the generated-code and runtime compatibility that a factory
// reference requires of its reader: the fourth version axis (row 16). It is
// stated in the facet, and the reader of the facet refuses a value it does
// not know.
const Runtime = "bittype-standin/1"

// Target is the target that the facets of this vocabulary are for.
var Target = string(s.Language{}.ID())

// Factory is the value of a factory facet: the factory reference, module and
// export, and the runtime it requires.
type Factory struct {
	Module  string `json:"module"`
	Export  string `json:"export"`
	Runtime string `json:"runtime"`
}

// transport is the transport descriptor of the two kinds: a media type and
// the encoder of the kind's owner. The kernel calls the encoder when it
// writes a selected contribution into a plan and never reads the result.
var transport = &g.FacetTransport{MediaType: "application/json", Encode: func(value any) ([]byte, error) {
	return json.Marshal(value)
}}

func kind(identity string) g.FacetKind {
	k, err := g.NewFacetKind(g.FacetKindDescription{Identity: identity, Version: "1", Transport: transport, Validate: func(value any) error {
		if f, ok := value.(Factory); !ok || f.Module == "" || f.Export == "" || f.Runtime == "" {
			return fmt.Errorf("the value is no factory reference")
		}
		return nil
	}})
	if err != nil {
		panic(err)
	}
	return k
}

// ValueFactory names the factory of a value adapter. FieldStateFactory names
// the factory of a field-state protocol, for a contribution at a field
// position, which has no standalone value adapter (row 14).
var (
	ValueFactory      = kind("bittype.completeness/typescript/value-factory")
	FieldStateFactory = kind("bittype.completeness/typescript/field-state-factory")
)

// Issue makes the factory facet of one kind for a contribution.
func Issue(kind g.FacetKind, contribution g.ContributionReference, export string) (g.Facet, error) {
	return kind.Issue(contribution.Identity.Publisher, contribution.Identity, Target, Factory{Module: "structural.mjs", Export: export, Runtime: Runtime})
}

// Conversion is the product: what it requires of every selected contribution.
func Conversion() g.Product {
	return g.Product{Identity: Product, Requires: map[g.SubjectKind][]g.FacetKind{
		g.UseSubject:         {ValueFactory},
		g.DeclarationSubject: {ValueFactory},
		g.FieldSubject:       {FieldStateFactory},
	}}
}
