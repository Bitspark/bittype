package integration_test

import (
	cppprovider "github.com/Bitspark/bittype/provider/cpp/go"
	goprovider "github.com/Bitspark/bittype/provider/golang/go"
	rustprovider "github.com/Bitspark/bittype/provider/rust/go"
	tsprovider "github.com/Bitspark/bittype/provider/typescript/go"
	"github.com/Bitspark/bittype/text/go"
	"strings"
	"testing"
)

func fuzzOne[L any](t *testing.T, p text.Provider[L], source string) {
	t.Helper()
	parsed, err := p.Parse(text.Input{SourceID: "fuzz", TreeID: "fuzz-tree", Text: source})
	if err != nil {
		if len(parsed.Nodes) > 0 {
			t.Fatal("partial concrete result on error")
		}
		return
	}
	printed, err := p.Print(parsed)
	if err != nil || printed != source {
		t.Fatalf("retention: %v", err)
	}
	if err := p.Check(parsed); err != nil {
		t.Fatal(err)
	}
}
func FuzzNativeText(f *testing.F) {
	f.Add(uint8(0), "type X = string;")
	f.Add(uint8(1), "package p; type X []string")
	f.Add(uint8(2), "type X<'a> = &'a str;")
	f.Add(uint8(3), "struct X { int (*p)[3]; };")
	f.Add(uint8(3), "using X = int(?);")
	f.Fuzz(func(t *testing.T, which uint8, source string) {
		if len(source) > 65536 {
			t.Skip()
		}
		switch which % 4 {
		case 0:
			fuzzOne(t, tsprovider.Provider{}, source)
		case 1:
			fuzzOne(t, goprovider.Provider{}, source)
		case 2:
			fuzzOne(t, rustprovider.Provider{}, source)
		case 3:
			fuzzOne(t, cppprovider.Provider{}, source)
		}
	})
}
func TestDeepSyntaxAndNativeNames(t *testing.T) {
	reject(t, tsprovider.Provider{}, []string{"//  @ts-ignore\ntype X = string;", "//\t@ts-ignore\ntype X = string;", "/* @ts-ignore */ type X = string;"})
	reject(t, cppprovider.Provider{}, []string{"// continued \\\nusing X = int;\n", "// continued \\\r\nusing X = int;\r\n", "// continued \\ \t\nusing X = int;\n"})
	reject(t, tsprovider.Provider{}, []string{"type string = number;", "interface any {}", "type X<number> = number;"})
	reject(t, rustprovider.Provider{}, []string{`type X = extern "\u{_43}" fn();`})
	reject(t, rustprovider.Provider{}, []string{"enum X { _ }", "trait X { fn _(); }", "trait X { type _; }", "type X<_> = u8;", "type X<'static> = &'static u8;", "type X<'_> = &'_ u8;", "type X = &'fn u8;"})
	reject(t, rustprovider.Provider{}, []string{"type X = (impl);", "type\u00a0X = u8;", "type X = extern \"\\a\" fn();"})
	reject(t, cppprovider.Provider{}, []string{"template<typename T,> struct X {};", "using X = Box<int,>;", "struct X { int f(int,); };"})
	reject(t, rustprovider.Provider{}, []string{"type X = " + strings.Repeat("(", 130) + "u8" + strings.Repeat(")", 130) + ";", "struct fn {}", "trait if {}", "type X = " + strings.Repeat("&", 300) + "u8;"})
	reject(t, tsprovider.Provider{}, []string{"type if = string;", "interface X { 1e3: string }", "interface X { 0x10: string }", "interface X { 9007199254740993: string }", "type X = string" + strings.Repeat("[]", 400) + ";"})
	reject(t, cppprovider.Provider{}, []string{"struct if {};", "using if = int;", "using X = int(?);", "using X = int(auto);", "#\ninclude \"x\"\n", "#include <x\n>\n", "using X = int(" + strings.Repeat("(", 150) + "*" + strings.Repeat(")", 150) + ");"})
	reject(t, goprovider.Provider{}, []string{"package p; type X " + strings.Repeat("[]", 300) + "int"})
}
