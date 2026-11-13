package integration_test

import (
	shared "github.com/Bitspark/bittype/integration/go/testdata/composition"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	gobooleancodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/boolean/go"
	godeclarationscodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/declarations/go"
	golistcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/list/go"
	gooptionalcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/optional/go"
	gorecordcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/record/go"
	gotextcodec "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/codec/text/go"
	gobooleanmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/boolean/go"
	godeclarationsmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/declarations/go"
	golistmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/list/go"
	gooptionalmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/optional/go"
	gorecordmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/record/go"
	gotextmapping "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/text/go"
)

// golangBundleTable names every builtin horos-to-golang contribution
// with the codec of its payload schema, by the directory of its package
// under transform/syntax/horos-to-golang/mapping. It is the one place in
// this package that lists the whole builtin set.
func golangBundleTable() map[string]genfixture.Bundle[s.Language] {
	return map[string]genfixture.Bundle[s.Language]{
		"text":         {Contribution: gotextmapping.Contribution, Codec: gotextcodec.Codec},
		"boolean":      {Contribution: gobooleanmapping.Contribution, Codec: gobooleancodec.Codec},
		"list":         {Contribution: golistmapping.Contribution, Codec: golistcodec.Codec},
		"record":       {Contribution: gorecordmapping.Contribution, Codec: gorecordcodec.Codec},
		"declarations": {Contribution: godeclarationsmapping.Contribution, Codec: godeclarationscodec.Codec},
		"optional":     {Contribution: gooptionalmapping.Contribution, Codec: gooptionalcodec.Codec},
	}
}

// golangBundles is the whole builtin set, as a consumer's composition
// registers it.
func golangBundles() []genfixture.Bundle[s.Language] {
	return shared.Select(golangBundleTable(), builtinBundleNames...)
}
