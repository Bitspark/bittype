package integration_test

import (
	"context"
	"github.com/Bitspark/bittype/contracts/parse/go"

	"strings"
	"testing"
)

func fuzzOne[L any](t *testing.T, p nativeProvider[L], source string) {
	t.Helper()
	parsed, err := p.Parse(context.Background(), parse.Input{SourceID: "fuzz", TreeID: "fuzz-tree", Text: source})
	if err != nil {
		if len(parsed.Nodes) > 0 {
			t.Fatal("partial concrete result on error")
		}
		return
	}
	printed, err := p.Print(context.Background(), parsed)
	if err != nil || printed != source {
		t.Fatalf("retention: %v", err)
	}
	if err := p.Check(parsed); err != nil {
		t.Fatal(err)
	}
}
func FuzzNativeText(f *testing.F) {
	f.Add(uint8(0), "type X = string;")
	f.Add(uint8(0), "type X<T> = T extends [infer U extends string] ? U : never;")
	f.Add(uint8(0), "type X<T> = { -readonly [K in keyof T as K extends string ? K : never]-?: T[K] };")
	f.Add(uint8(0), "type X = A extends B ? C extends D ? E : F : G;")
	f.Add(uint8(1), "package p; type X []string")
	f.Add(uint8(2), "type X<'a> = &'a str;")
	f.Add(uint8(2), "type X<T> = <<T as Source>::Item as Iterator>::Item;")
	f.Add(uint8(2), "trait X<'a, T>: Source<'a, T, Item<'a> = &'a T> {}")
	f.Add(uint8(2), "trait X: Iterator<(Item) = u8> {}")
	f.Add(uint8(2), "trait X { type Item<'a, T: ?Sized>: 'a where T: 'a; }")
	f.Add(uint8(2), "impl<T> X for Vec<T> { type Item<'a> = &'a T where Self: 'a; }")
	f.Add(uint8(3), "struct X { int (*p)[3]; };")
	f.Add(uint8(3), "using X = int(?);")
	cases := nativeConformanceCases()
	f.Fuzz(func(t *testing.T, which uint8, source string) {
		if len(source) > 65536 {
			t.Skip()
		}
		cases[int(which)%len(cases)].fuzz(t, source)
	})
}
func TestDeepSyntaxAndNativeNames(t *testing.T) {
	reject(t, typescriptText(), []string{"//  @ts-ignore\ntype X = string;", "//\t@ts-ignore\ntype X = string;", "/* @ts-ignore */ type X = string;"})
	reject(t, cppText(), []string{"// continued \\\nusing X = int;\n", "// continued \\\r\nusing X = int;\r\n", "// continued \\ \t\nusing X = int;\n"})
	reject(t, typescriptText(), []string{"type string = number;", "interface any {}", "type X<number> = number;"})
	reject(t, rustText(), []string{`type X = extern "\u{_43}" fn();`})
	reject(t, rustText(), []string{"enum X { _ }", "trait X { fn _(); }", "trait X { type _; }", "type X<_> = u8;", "type X<'static> = &'static u8;", "type X<'_> = &'_ u8;", "type X = &'fn u8;"})
	reject(t, rustText(), []string{"type X = (impl);", "type\u00a0X = u8;", "type X = extern \"\\a\" fn();"})
	reject(t, cppText(), []string{"template<typename T,> struct X {};", "using X = Box<int,>;", "struct X { int f(int,); };"})
	reject(t, rustText(), []string{"type X = " + strings.Repeat("(", 130) + "u8" + strings.Repeat(")", 130) + ";", "struct fn {}", "trait if {}", "type X = " + strings.Repeat("&", 300) + "u8;"})
	reject(t, typescriptText(), []string{"type if = string;", "interface X { 1e3: string }", "interface X { 0x10: string }", "interface X { 9007199254740993: string }", "type X = string" + strings.Repeat("[]", 400) + ";"})
	reject(t, cppText(), []string{"struct if {};", "using if = int;", "using X = int(?);", "using X = int(auto);", "#\ninclude \"x\"\n", "#include <x\n>\n", "using X = int(" + strings.Repeat("(", 150) + "*" + strings.Repeat(")", 150) + ");"})
	reject(t, golangText(), []string{"package p; type X " + strings.Repeat("[]", 300) + "int"})
}
