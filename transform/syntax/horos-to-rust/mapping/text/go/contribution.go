// Package text is the builtin horos-to-rust contribution for the horos
// scalar utf8-text (docs/contribution-contract.md). A consumer's composition
// registers it on its own, with the codec of its payload schema.
package text

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	vocabulary "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/vocabulary/text/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/rust/model/go"
	shared "github.com/Bitspark/bittype/transform/syntax/shared/contribution/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contribution = g.Contribution[h.Language, s.Language, source.Reader]

// target is the target that the contribution states (O4).
var target = string(s.Language{}.ID())

// Contribution returns the text contribution: the horos scalar utf8-text
// becomes the Rust model primitive string. Nothing varies between two
// selections, so the payload is empty. The native type is an invariant that
// the representation revision fixes (row 5).
func Contribution(from representation.Language) (g.Contribution[h.Language, s.Language, source.Reader], error) {
	selector, err := source.Builtin(from, "scalar", "utf8-text")
	if err != nil {
		return contribution{}, err
	}
	return shared.Constant(contribution{
		Reference: vocabulary.Reference, Target: target, Profile: rules.RepresentationProfile,
		Domain: g.ForForm(selector), Schema: vocabulary.Schema,
	}, vocabulary.Payload{}, model.Primitive(model.PrimitiveString)), nil
}
