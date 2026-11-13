package integration_test

import (
	shared "github.com/Bitspark/bittype/integration/go/testdata/composition"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	cppbooleancodec "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/codec/boolean/go"
	cppdeclarationscodec "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/codec/declarations/go"
	cpplistcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/codec/list/go"
	cppoptionalcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/codec/optional/go"
	cpprecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/codec/record/go"
	cpptextcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/codec/text/go"
	cppbooleanmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/boolean/go"
	cppdeclarationsmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/declarations/go"
	cpplistmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/list/go"
	cppoptionalmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/optional/go"
	cpprecordmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/record/go"
	cpptextmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/text/go"
)

// cppBundleTable names every builtin horos-to-cpp contribution
// with the codec of its payload schema, by the directory of its package
// under transform/syntax/horos-to-cpp/mapping. It is the one place in
// this package that lists the whole builtin set.
func cppBundleTable() map[string]genfixture.Bundle[s.Language] {
	return map[string]genfixture.Bundle[s.Language]{
		"text":         {Contribution: cpptextmapping.Contribution, Codec: cpptextcodec.Codec},
		"boolean":      {Contribution: cppbooleanmapping.Contribution, Codec: cppbooleancodec.Codec},
		"list":         {Contribution: cpplistmapping.Contribution, Codec: cpplistcodec.Codec},
		"record":       {Contribution: cpprecordmapping.Contribution, Codec: cpprecordcodec.Codec},
		"declarations": {Contribution: cppdeclarationsmapping.Contribution, Codec: cppdeclarationscodec.Codec},
		"optional":     {Contribution: cppoptionalmapping.Contribution, Codec: cppoptionalcodec.Codec},
	}
}

// cppBundles is the whole builtin set, as a consumer's composition
// registers it.
func cppBundles() []genfixture.Bundle[s.Language] {
	return shared.Select(cppBundleTable(), builtinBundleNames...)
}
