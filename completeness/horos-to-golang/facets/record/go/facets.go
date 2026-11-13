// Package record is the facet package of the record contribution for Go
// (prototype, bittype#196, O8, O16): its lowering to Go conversion code.
// It is test code of bittype's own bundle. It imports the contribution's
// vocabulary and the Go value vocabulary, and neither the mapping package nor
// the codec. A consumer's composition registers it apart from the mapping.
package record

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Bitspark/bittype/completeness/golang/values/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/vocabulary/record/go"
)

// Facets returns the facets of the contribution for the Go conversion product.
func Facets() ([]g.Facet, error) {
	f, err := values.Issue(values.ValueLowering, vocabulary.Reference, values.Lowering{Lower: Lower})
	return []g.Facet{f}, err
}

// Lower lowers a record to an adapter of its declared type. It constructs the
// native value with a typed construction and reads it by member access, and
// every member enters through a member reference of the plan: the lowering
// spells no member. For a field whose presence the record delegates it
// invokes the field-state protocol of whatever was selected for the field
// position: it tests the abstract record for the field, hands the state to
// Construct, and on the way back asks State whether to write the field.
func Lower(ctx values.Context) (values.Code, error) {
	var holes []values.Hole
	hole := func(h values.Hole) string {
		holes = append(holes, h)
		return fmt.Sprintf("«%d»", len(holes)-1)
	}
	self := hole(values.SelfType())
	var names []string
	var locals, decode, literal, encode strings.Builder
	for i := 0; i < ctx.Links(); i++ {
		canonical, err := ctx.Canonical(i)
		if err != nil {
			return values.Code{}, err
		}
		reference, err := ctx.Member(i)
		if err != nil {
			return values.Code{}, err
		}
		member, name := hole(values.Reference(reference, g.MemberReference)), strconv.Quote(canonical)
		names = append(names, name)
		fmt.Fprintf(&locals, "\t\tf%d := %s\n", i, hole(values.Dependency(i)))
		if ctx.Fragment(i) {
			fmt.Fprintf(&decode, "\t\t\t\tv%d, present%d := m[%s]\n\t\t\t\tx%d, err%d := f%d.Construct(present%d, v%d)\n\t\t\t\tif err%d != nil {\n\t\t\t\t\treturn out, err%d\n\t\t\t\t}\n", i, i, name, i, i, i, i, i, i, i)
			fmt.Fprintf(&encode, "\t\t\t\tif present, v, err := f%d.State(n.%s); err != nil {\n\t\t\t\t\treturn nil, err\n\t\t\t\t} else if present {\n\t\t\t\t\tm[%s] = v\n\t\t\t\t}\n", i, member, name)
		} else {
			fmt.Fprintf(&decode, "\t\t\t\tv%d, present%d := m[%s]\n\t\t\t\tif !present%d {\n\t\t\t\t\treturn out, fmt.Errorf(\"structural mismatch: the required field %%q is absent\", %s)\n\t\t\t\t}\n\t\t\t\tx%d, err%d := f%d.Decode(v%d)\n\t\t\t\tif err%d != nil {\n\t\t\t\t\treturn out, err%d\n\t\t\t\t}\n", i, i, name, i, name, i, i, i, i, i, i)
			fmt.Fprintf(&encode, "\t\t\t\tif v, err := f%d.Encode(n.%s); err != nil {\n\t\t\t\t\treturn nil, err\n\t\t\t\t} else {\n\t\t\t\t\tm[%s] = v\n\t\t\t\t}\n", i, member, name)
		}
		fmt.Fprintf(&literal, "%s: x%d, ", member, i)
	}
	text := "func() Adapter[" + self + "] {\n" + locals.String() +
		"\t\treturn Adapter[" + self + "]{\n" +
		"\t\t\tDecode: func(v Value) (" + self + ", error) {\n\t\t\t\tvar out " + self + "\n\t\t\t\tm, err := Fields(v, " + strings.Join(names, ", ") + ")\n\t\t\t\tif err != nil {\n\t\t\t\t\treturn out, err\n\t\t\t\t}\n" +
		decode.String() + "\t\t\t\treturn " + self + "{" + literal.String() + "}, nil\n\t\t\t},\n" +
		"\t\t\tEncode: func(n " + self + ") (Value, error) {\n\t\t\t\tm := map[string]Value{}\n" + encode.String() + "\t\t\t\treturn m, nil\n\t\t\t},\n\t\t}\n\t}()"
	return values.Code{Text: text, Holes: holes}, nil
}
