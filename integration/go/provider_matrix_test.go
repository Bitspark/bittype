package integration_test

import (
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	cpp "github.com/Bitspark/bittype/provider/cpp/go"
	golang "github.com/Bitspark/bittype/provider/golang/go"
	horos "github.com/Bitspark/bittype/provider/horos/go"
	horoswire "github.com/Bitspark/bittype/provider/horoswire/go"
	rust "github.com/Bitspark/bittype/provider/rust/go"
	typescript "github.com/Bitspark/bittype/provider/typescript/go"
	"testing"
)

func TestDocumentedProviderMatrix(t *testing.T) {
	providerCapabilities(t, horos.Provider{}, false, true)
	providerCapabilities(t, horoswire.Provider{}, false, true)
	providerCapabilities(t, typescript.Provider{}, true, false)
	providerCapabilities(t, golang.Provider{}, true, false)
	providerCapabilities(t, rust.Provider{}, true, false)
	providerCapabilities(t, cpp.Provider{}, true, false)
}

func providerCapabilities[L any](t *testing.T, p check.Provider[L], nativeText, semantic bool) {
	t.Helper()
	t.Run(string(p.ID()), func(t *testing.T) {
		if p.ID() == "" || p.Name() == "" {
			t.Fatal("missing language identity")
		}
		_, parser := any(p).(parse.Provider[L])
		_, formatter := any(p).(format.Provider[L])
		_, printer := any(p).(print.Provider[L])
		_, combined := any(p).(text.Provider[L])
		_, resolver := any(p).(resolution.Provider[L])
		_, analyzer := any(p).(analysis.Provider[L])
		if parser != nativeText || formatter != nativeText || printer != nativeText || combined != nativeText || resolver != semantic || analyzer != semantic {
			t.Fatalf("capability matrix changed: parser=%v formatter=%v printer=%v text=%v resolver=%v analyzer=%v", parser, formatter, printer, combined, resolver, analyzer)
		}
	})
}
