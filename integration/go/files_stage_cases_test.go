package integration_test

import (
	"reflect"
	"testing"

	formathoros "github.com/Bitspark/bittype/format/files/horos/provider/go"
	formatts "github.com/Bitspark/bittype/format/files/typescript/provider/go"
	parsehoros "github.com/Bitspark/bittype/parse/files/horos/provider/go"
)

// This test suite's selected providers and their independently runnable tests.
// The coverage check discovers stage packages and consumes this composition.
func filesStageFixtures() map[string]func(*testing.T) {
	return map[string]func(*testing.T){
		reflect.TypeFor[parsehoros.Provider]().PkgPath():  TestFilesPaginationAndReceiptPreservingRoundTrips,
		reflect.TypeFor[formathoros.Provider]().PkgPath(): TestFilesPaginationAndReceiptPreservingRoundTrips,
		reflect.TypeFor[formatts.Provider]().PkgPath():    TestTypeScriptTypeImportPlacement,
	}
}
