package integration_test

import (
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	rustbooleancodec "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/codec/boolean/go"
	rustdeclarationscodec "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/codec/declarations/go"
	rustlistcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/codec/list/go"
	rustoptionalcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/codec/optional/go"
	rustrecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/codec/record/go"
	rusttextcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/codec/text/go"
	rustbooleanmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/boolean/go"
	rustdeclarationsmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/declarations/go"
	rustlistmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/list/go"
	rustoptionalmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/optional/go"
	rustrecordmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/record/go"
	rusttextmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/text/go"
)

// rustBundleNames lists the builtin horos-to-rust contributions
// in the order a test registers them. Registration order is no preference.
var rustBundleNames = []string{"text", "boolean", "list", "record", "declarations", "optional"}

// rustBundleTable names every builtin horos-to-rust contribution
// with the codec of its payload schema, by the directory of its package
// under transform/syntax/horos-to-rust/mapping. It is the one place in
// this package that lists the whole builtin set.
func rustBundleTable() map[string]genfixture.Bundle[s.Language] {
	return map[string]genfixture.Bundle[s.Language]{
		"text":         {Contribution: rusttextmapping.Contribution, Codec: rusttextcodec.Codec},
		"boolean":      {Contribution: rustbooleanmapping.Contribution, Codec: rustbooleancodec.Codec},
		"list":         {Contribution: rustlistmapping.Contribution, Codec: rustlistcodec.Codec},
		"record":       {Contribution: rustrecordmapping.Contribution, Codec: rustrecordcodec.Codec},
		"declarations": {Contribution: rustdeclarationsmapping.Contribution, Codec: rustdeclarationscodec.Codec},
		"optional":     {Contribution: rustoptionalmapping.Contribution, Codec: rustoptionalcodec.Codec},
	}
}

// rustBundles is the whole builtin set, as a consumer's composition
// registers it. The record contribution delegates the presence of an optional
// field to the optional contribution through the planner: neither names the
// other, and no consumer creates a key for them.
func rustBundles() []genfixture.Bundle[s.Language] {
	table := rustBundleTable()
	var bundles []genfixture.Bundle[s.Language]
	for _, name := range rustBundleNames {
		bundles = append(bundles, table[name])
	}
	return bundles
}
