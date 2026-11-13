// Package boolean is the builtin horos-to-typescript contribution for the horos
// scalar bool (docs/contribution-contract.md). A consumer's composition
// registers it on its own, with the codec of its payload schema.
package boolean

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/vocabulary/boolean/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	shared "github.com/Bitspark/bittype/transform/syntax/shared/contribution/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[h.Language, s.Language, source.Reader]

// target is the target that the contribution states (O4).
var target = string(s.Language{}.ID())

// Contribution returns the boolean contribution: the horos scalar bool
// becomes the TypeScript model primitive bool. Nothing varies between two
// selections, so the payload is empty. The native type is an invariant that
// the representation revision fixes (row 5).
func Contribution(from representation.Language) (g.Contribution[h.Language, s.Language, source.Reader], error) {
	selector, err := source.Builtin(from, "scalar", "bool")
	if err != nil {
		return contribution{}, err
	}
	return shared.Constant(contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.RepresentationProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
	}, vocabulary.Payload{}, model.Primitive(model.PrimitiveBool)), nil
}
