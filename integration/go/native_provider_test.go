package integration_test

import (
	checkcpp "github.com/Bitspark/bittype/check/syntax/cpp/provider/go"
	checkgolang "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	checkhoros "github.com/Bitspark/bittype/check/syntax/horos/provider/go"
	checkrust "github.com/Bitspark/bittype/check/syntax/rust/provider/go"
	checktypescript "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	formatcpp "github.com/Bitspark/bittype/format/syntax/cpp/provider/go"
	formatgolang "github.com/Bitspark/bittype/format/syntax/golang/provider/go"
	formathoros "github.com/Bitspark/bittype/format/syntax/horos/provider/go"
	formatrust "github.com/Bitspark/bittype/format/syntax/rust/provider/go"
	formattypescript "github.com/Bitspark/bittype/format/syntax/typescript/provider/go"
	tf "github.com/Bitspark/bittype/integration/go/testdata/textfixture"
	parsecpp "github.com/Bitspark/bittype/parse/syntax/cpp/provider/go"
	parsegolang "github.com/Bitspark/bittype/parse/syntax/golang/provider/go"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	parserust "github.com/Bitspark/bittype/parse/syntax/rust/provider/go"
	parsetypescript "github.com/Bitspark/bittype/parse/syntax/typescript/provider/go"
	cpp "github.com/Bitspark/bittype/syntax/cpp/go"
	golang "github.com/Bitspark/bittype/syntax/golang/go"
	horos "github.com/Bitspark/bittype/syntax/horos/go"
	rust "github.com/Bitspark/bittype/syntax/rust/go"
	typescript "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

// Native conformance fixtures exercise four capabilities of one implementation.
// Engine registration and production contracts do not require this combination.
type nativeProvider[L any] interface {
	parse.Provider[parse.Input, cst.CST[L]]
	check.Provider[L]
	format.Provider[ast.View[L], cst.CST[L]]
	print.Provider[cst.CST[L], string]
}

func typescriptText() tf.Fixture[typescript.Language] {
	return tf.New(parsetypescript.Provider{}, checktypescript.Provider{}, formattypescript.Provider{})
}

func golangText() tf.Fixture[golang.Language] {
	return tf.New(parsegolang.Provider{}, checkgolang.Provider{}, formatgolang.Provider{})
}

func rustText() tf.Fixture[rust.Language] {
	return tf.New(parserust.Provider{}, checkrust.Provider{}, formatrust.Provider{})
}

func cppText() tf.Fixture[cpp.Language] {
	return tf.New(parsecpp.Provider{}, checkcpp.Provider{}, formatcpp.Provider{})
}

func horosText() tf.Fixture[horos.Language] {
	return tf.New(parsehoros.Provider{}, checkhoros.Provider{}, formathoros.Provider{})
}
