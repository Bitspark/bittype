package integration_test

import (
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	tsbooleancodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/boolean/go"
	tsdeclarationscodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/declarations/go"
	tslistcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/list/go"
	tsoptionalcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/optional/go"
	tsrecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/record/go"
	tstextcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/codec/text/go"
	tsbooleanmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/boolean/go"
	tsdeclarationsmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/declarations/go"
	tslistmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/list/go"
	tsoptionalmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/optional/go"
	tsrecordmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/record/go"
	tstextmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/text/go"
)

// typescriptBundleNames lists the builtin horos-to-typescript contributions
// in the order a test registers them. Registration order is no preference.
var typescriptBundleNames = []string{"text", "boolean", "list", "record", "declarations", "optional"}

// typescriptBundleTable names every builtin horos-to-typescript contribution
// with the codec of its payload schema, by the directory of its package
// under transform/syntax/horos-to-typescript/mapping. It is the one place in
// this package that lists the whole builtin set.
func typescriptBundleTable() map[string]genfixture.Bundle[s.Language] {
	return map[string]genfixture.Bundle[s.Language]{
		"text":         {Contribution: tstextmapping.Contribution, Codec: tstextcodec.Codec},
		"boolean":      {Contribution: tsbooleanmapping.Contribution, Codec: tsbooleancodec.Codec},
		"list":         {Contribution: tslistmapping.Contribution, Codec: tslistcodec.Codec},
		"record":       {Contribution: tsrecordmapping.Contribution, Codec: tsrecordcodec.Codec},
		"declarations": {Contribution: tsdeclarationsmapping.Contribution, Codec: tsdeclarationscodec.Codec},
		"optional":     {Contribution: tsoptionalmapping.Contribution, Codec: tsoptionalcodec.Codec},
	}
}

// typescriptBundles is the whole builtin set, as a consumer's composition
// registers it. The record contribution delegates the presence of an optional
// field to the optional contribution through the planner: neither names the
// other, and no consumer creates a key for them.
func typescriptBundles() []genfixture.Bundle[s.Language] {
	table := typescriptBundleTable()
	var bundles []genfixture.Bundle[s.Language]
	for _, name := range typescriptBundleNames {
		bundles = append(bundles, table[name])
	}
	return bundles
}
