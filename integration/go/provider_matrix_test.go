package integration_test

import (
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	cpp "github.com/Bitspark/bittype/provider/cpp/go"
	golang "github.com/Bitspark/bittype/provider/golang/go"
	horos "github.com/Bitspark/bittype/provider/horos/go"
	rust "github.com/Bitspark/bittype/provider/rust/go"
	typescript "github.com/Bitspark/bittype/provider/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"testing"
)

func TestDocumentedProviderMatrix(t *testing.T) {
	providerCapabilities(t, horos.Provider{}, true, true)
	providerCapabilities(t, typescript.Provider{}, true, false)
	providerCapabilities(t, golang.Provider{}, true, true)
	providerCapabilities(t, rust.Provider{}, true, false)
	providerCapabilities(t, cpp.Provider{}, true, false)
}

func providerCapabilities[L any](t *testing.T, p check.Provider[L], textStages, semantic bool) {
	t.Helper()
	t.Run(string(p.ID()), func(t *testing.T) {
		if p.ID() == "" || p.Name() == "" {
			t.Fatal("missing language identity")
		}
		_, parser := any(p).(parse.Provider[parse.Input, cst.CST[L]])
		_, formatter := any(p).(format.Provider[ast.View[L], cst.CST[L]])
		_, printer := any(p).(print.Provider[cst.CST[L], string])
		_, resolver := any(p).(resolution.Provider[L])
		_, analyzer := any(p).(analysis.Provider[L])
		if parser != textStages || formatter != textStages || printer != textStages || resolver != semantic || analyzer != semantic {
			t.Fatalf("capability matrix changed: parser=%v formatter=%v printer=%v resolver=%v analyzer=%v", parser, formatter, printer, resolver, analyzer)
		}
	})
}
