package integration_test

import (
	analyzegolang "github.com/Bitspark/bittype/analyze/files/golang/provider/go"
	analyzehoros "github.com/Bitspark/bittype/analyze/files/horos/provider/go"
	analyzerust "github.com/Bitspark/bittype/analyze/files/rust/provider/go"
	analyzetypescript "github.com/Bitspark/bittype/analyze/files/typescript/provider/go"
	checkcpp "github.com/Bitspark/bittype/check/syntax/cpp/provider/go"
	checkgolang "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	checkhoros "github.com/Bitspark/bittype/check/syntax/horos/provider/go"
	checkrust "github.com/Bitspark/bittype/check/syntax/rust/provider/go"
	checktypescript "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	formatcpp "github.com/Bitspark/bittype/format/syntax/cpp/provider/go"
	formatgolang "github.com/Bitspark/bittype/format/syntax/golang/provider/go"
	formathoros "github.com/Bitspark/bittype/format/syntax/horos/provider/go"
	formatrust "github.com/Bitspark/bittype/format/syntax/rust/provider/go"
	formattypescript "github.com/Bitspark/bittype/format/syntax/typescript/provider/go"
	textfixture "github.com/Bitspark/bittype/integration/go/testdata/textfixture"
	parsecpp "github.com/Bitspark/bittype/parse/syntax/cpp/provider/go"
	parsegolang "github.com/Bitspark/bittype/parse/syntax/golang/provider/go"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	parserust "github.com/Bitspark/bittype/parse/syntax/rust/provider/go"
	parsetypescript "github.com/Bitspark/bittype/parse/syntax/typescript/provider/go"
	resolvegolang "github.com/Bitspark/bittype/resolve/files/golang/provider/go"
	resolvehoros "github.com/Bitspark/bittype/resolve/files/horos/provider/go"
	resolverust "github.com/Bitspark/bittype/resolve/files/rust/provider/go"
	resolvetypescript "github.com/Bitspark/bittype/resolve/files/typescript/provider/go"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"

	"github.com/Bitspark/bittype/contracts/resolution/go"

	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"testing"
)

// Each role is a separately supplied implementation with the same language ID.
func TestDocumentedProviderMatrix(t *testing.T) {
	providerCapabilities(t, parsehoros.Provider{}, checkhoros.Provider{}, formathoros.Provider{}, resolvehoros.Provider{}, analyzehoros.Provider{})
	providerCapabilities(t, parsetypescript.Provider{}, checktypescript.Provider{}, formattypescript.Provider{}, resolvetypescript.Provider{}, analyzetypescript.Provider{})
	providerCapabilities(t, parsegolang.Provider{}, checkgolang.Provider{}, formatgolang.Provider{}, resolvegolang.Provider{}, analyzegolang.Provider{})
	providerCapabilities(t, parserust.Provider{}, checkrust.Provider{}, formatrust.Provider{}, resolverust.Provider{}, analyzerust.Provider{})
	providerCapabilities(t, parsecpp.Provider{}, checkcpp.Provider{}, formatcpp.Provider{}, nil, nil)
}
func providerCapabilities[L any](t *testing.T, p parse.Provider[parse.Input, cst.CST[L]], c check.Provider[L], f format.Provider[ast.View[L], cst.CST[L]], r resolution.Provider[L], a analysis.Provider[L]) {
	t.Helper()
	printer := textfixture.Printing[L](p)
	if p.ID() == "" || c.ID() != p.ID() || f.ID() != p.ID() || printer.ID() != p.ID() {
		t.Fatal("syntax provider identity mismatch", p.ID())
	}
	if (r == nil) != (a == nil) || r != nil && (r.ID() != p.ID() || a.ID() != p.ID()) {
		t.Fatal("semantic provider identity mismatch", p.ID())
	}
}
