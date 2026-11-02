package integration_test

import (
    "testing"
    "github.com/Bitspark/bittype/contracts/analysis/go"
    "github.com/Bitspark/bittype/contracts/resolution/go"
    "github.com/Bitspark/bittype/contracts/text/go"
    cpp "github.com/Bitspark/bittype/provider/cpp/go"
    golang "github.com/Bitspark/bittype/provider/golang/go"
    horos "github.com/Bitspark/bittype/provider/horos/go"
    rust "github.com/Bitspark/bittype/provider/rust/go"
    typescript "github.com/Bitspark/bittype/provider/typescript/go"
)

func TestDocumentedProviderMatrix(t *testing.T) {
    providerCapabilities(t,horos.Provider{},false,true)
    providerCapabilities(t,typescript.Provider{},true,false)
    providerCapabilities(t,golang.Provider{},true,false)
    providerCapabilities(t,rust.Provider{},true,false)
    providerCapabilities(t,cpp.Provider{},true,false)
}

func providerCapabilities[L any](t *testing.T, p text.CheckerProvider[L], nativeText, semantic bool) {
    t.Helper()
    t.Run(string(p.ID()),func(t *testing.T){
        if p.ID()=="" || p.Name()=="" {t.Fatal("missing language identity")}
        _,parser:=any(p).(text.ParserProvider[L])
        _,formatter:=any(p).(text.FormatterProvider[L])
        _,printer:=any(p).(text.PrinterProvider[L])
        _,combined:=any(p).(text.Provider[L])
        _,resolver:=any(p).(resolution.ResolverProvider[L])
        _,analyzer:=any(p).(analysis.Provider[L])
        if parser!=nativeText || formatter!=nativeText || printer!=nativeText || combined!=nativeText || resolver!=semantic || analyzer!=semantic {t.Fatalf("capability matrix changed: parser=%v formatter=%v printer=%v text=%v resolver=%v analyzer=%v",parser,formatter,printer,combined,resolver,analyzer)}
    })
}
