package integration_test

import (
	"errors"
	"github.com/Bitspark/bittype/ast/go"
	cstcheck "github.com/Bitspark/bittype/check/cst/go"
	"github.com/Bitspark/bittype/cst/go"
	"github.com/Bitspark/bittype/language/go"
	cppprovider "github.com/Bitspark/bittype/provider/cpp/go"
	goprovider "github.com/Bitspark/bittype/provider/golang/go"
	rustprovider "github.com/Bitspark/bittype/provider/rust/go"
	tsprovider "github.com/Bitspark/bittype/provider/typescript/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/text/go"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type tsMarker = ts.Language

type syntaxShape struct {
	Kind       string
	Attributes []ast.Attribute
	Roles      []string
	Children   []syntaxShape
}

func shape[L any](t ast.AST[L], id ast.NodeID) syntaxShape {
	n := t.Nodes[id-1]
	out := syntaxShape{Kind: string(n.Kind)}
	for _, a := range n.Attributes {
		if a.Text == "false" || (a.Name == "direction" && a.Text == "both") || (a.Name == "visibility" && a.Text == "private") || (a.Name == "parameter-key" && a.Text == "typename") || (a.Name == "ref-qualifier" && a.Text == "none") {
			continue
		}
		out.Attributes = append(out.Attributes, a)
	}
	sort.Slice(out.Attributes, func(i, j int) bool { return out.Attributes[i].Name < out.Attributes[j].Name })
	for _, c := range n.Children {
		out.Roles = append(out.Roles, c.Role)
		out.Children = append(out.Children, shape(t, c.Node))
	}
	return out
}
func exercise[L any](t *testing.T, p text.Provider[L], directory, extension string) {
	t.Helper()
	paths, err := filepath.Glob("../../testdata/text/" + directory + "/valid/*." + extension)
	if err != nil || len(paths) == 0 {
		t.Fatalf("fixtures: %v", err)
	}
	for _, path := range paths {
		t.Run(filepath.Base(path), func(t *testing.T) {
			bytes, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			source := string(bytes)
			input := text.Input{SourceID: "fixture", TreeID: "parsed", Text: source}
			parsed, err := p.Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			printed, err := p.Print(parsed)
			if err != nil || printed != source {
				t.Fatalf("source retention failed: %v\n%q", err, printed)
			}
			if err := p.Check(parsed); err != nil {
				t.Fatal(err)
			}
			concrete, err := p.Format(parsed)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(parsed.Abstract, concrete.Abstract) {
				t.Fatal("formatter changed AST view or metadata")
			}
			output, err := p.Print(concrete)
			if err != nil {
				t.Fatal(err)
			}
			reparsed, err := p.Parse(text.Input{SourceID: "formatted", TreeID: "reparsed", Text: output})
			if err != nil {
				t.Fatalf("%v\n%s", err, output)
			}
			if !reflect.DeepEqual(shape(parsed.Abstract, parsed.Abstract.Root), shape(reparsed.Abstract, reparsed.Abstract.Root)) {
				t.Fatalf("AST syntax changed\nBEFORE %#v\nAFTER %#v\n%s", shape(parsed.Abstract, parsed.Abstract.Root), shape(reparsed.Abstract, reparsed.Abstract.Root), output)
			}
			second, err := p.Format(reparsed)
			if err != nil {
				t.Fatal(err)
			}
			stable, err := p.Print(second)
			if err != nil || stable != output {
				t.Fatalf("formatting unstable: %v\n%s\n%s", err, output, stable)
			}
			t.Log(output)
		})
	}
}
func TestNativeTextRoundTrips(t *testing.T) {
	t.Run("typescript", func(t *testing.T) { exercise(t, tsprovider.Provider{}, "typescript", "ts") })
	t.Run("go", func(t *testing.T) { exercise(t, goprovider.Provider{}, "golang", "go") })
	t.Run("rust", func(t *testing.T) { exercise(t, rustprovider.Provider{}, "rust", "rs") })
	t.Run("cpp", func(t *testing.T) { exercise(t, cppprovider.Provider{}, "cpp", "hpp") })
}
func reject[L any](t *testing.T, p text.Provider[L], sources []string) {
	t.Helper()
	for _, source := range sources {
		t.Run(source, func(t *testing.T) {
			c, err := p.Parse(text.Input{SourceID: "bad", TreeID: "bad-tree", Text: source})
			var d *text.Diagnostic
			if err == nil || !errors.As(err, &d) || len(c.Nodes) != 0 {
				t.Fatalf("expected located error and no partial CST, got %v", err)
			}
		})
	}
}
func TestUnsupportedAndMalformedSources(t *testing.T) {
	reject(t, tsprovider.Provider{}, []string{"const x = 1;", "type X = T extends U ? A : B;", "type X = { [P in keyof T]: T[P] };", "type X = string; garbage", "type X = (x: string;", "type X = [a?:string,b:number];", "type X = typeof foo;", "type X = 1 + 2;", "import type {} from 'x';", "/// <reference path='x'/>\ntype X=string;"})
	reject(t, goprovider.Provider{}, []string{"package p; var X int", "package p; type X [1+2]int", "package p; type X [N]int", "package p; type X struct {", "//go:build ignore\npackage p; type X int"})
	reject(t, rustprovider.Provider{}, []string{"struct Unit;", "pub fn f() {}", "type X = impl Clone;", "type X<const N:usize> = [u8;N];", "type X = u8; trailing", "struct X { value: }", "type X = for<'a> fn(&'a str);", "type X = 'x';", "use std::{string,vec};"})
	reject(t, cppprovider.Provider{}, []string{"#define X int\nusing A = X;", "using X = decltype(1);", "struct X { int a = 1; };", "struct X { int a[N]; };", "template<int N> struct X {};", "using X = int; trailing", "struct X { int (X::*f)() const; };", "struct X { X(); };", "struct X { virtual int f() = 1; };"})
}
func TestConcreteSourceTrivia(t *testing.T) {
	cases := []struct {
		name, source string
		run          func(text.Input) (string, error)
	}{
		{"typescript", "\ufeff// comment\r\ntype Café = 'é';\r\n\t", func(i text.Input) (string, error) {
			p := tsprovider.Provider{}
			c, e := p.Parse(i)
			if e != nil {
				return "", e
			}
			return p.Print(c)
		}},
		{"go", "\ufeffpackage p\r\n// comment\r\ntype Café string\r\n", func(i text.Input) (string, error) {
			p := goprovider.Provider{}
			c, e := p.Parse(i)
			if e != nil {
				return "", e
			}
			return p.Print(c)
		}},
		{"rust", "/* one /* two */ three */\r\ntype Ref<'a> = &'a str; // end\r\n", func(i text.Input) (string, error) {
			p := rustprovider.Provider{}
			c, e := p.Parse(i)
			if e != nil {
				return "", e
			}
			return p.Print(c)
		}},
		{"cpp", "// comment\r\nusing X = int;\r\n\t", func(i text.Input) (string, error) {
			p := cppprovider.Provider{}
			c, e := p.Parse(i)
			if e != nil {
				return "", e
			}
			return p.Print(c)
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.run(text.Input{SourceID: "trivia", TreeID: "trivia-tree", Text: tc.source})
			if err != nil || got != tc.source {
				t.Fatalf("%v: %q", err, got)
			}
		})
	}
}
func TestProviderIdentity(t *testing.T) {
	for _, p := range []language.LanguageProvider{tsprovider.Provider{}, goprovider.Provider{}, rustprovider.Provider{}, cppprovider.Provider{}} {
		if p.ID() == "" || p.Name() == "" {
			t.Fatal("missing provider identity")
		}
	}
}
func TestMalformedConcreteTree(t *testing.T) {
	p := tsprovider.Provider{}
	parsed, err := p.Parse(text.Input{SourceID: "good", TreeID: "tree", Text: "type X = string;"})
	if err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*cst.CST[tsMarker]){
		func(c *cst.CST[tsMarker]) { c.Root = 0 },
		func(c *cst.CST[tsMarker]) { c.Nodes[c.Root-1].Children = append(c.Nodes[c.Root-1].Children, c.Root) },
		func(c *cst.CST[tsMarker]) { c.Nodes[c.Root-1].Syntax = []ast.NodeID{99999} },
		func(c *cst.CST[tsMarker]) { c.Nodes = append(c.Nodes, cst.Node[tsMarker]{Kind: "orphan"}) },
	} {
		copy := parsed
		copy.Nodes = append([]cst.Node[tsMarker](nil), parsed.Nodes...)
		mutate(&copy)
		if err := cstcheck.Check(copy); err == nil {
			t.Fatal("malformed concrete tree accepted")
		}
		if output, err := p.Print(copy); err == nil || output != "" {
			t.Fatal("printer returned partial output")
		}
	}
}
func TestParserResourceLimit(t *testing.T) {
	p := tsprovider.Provider{}
	_, err := p.Parse(text.Input{SourceID: "large", TreeID: "large-tree", Text: strings.Repeat(" ", (4<<20)+1)})
	var d *text.Diagnostic
	if !errors.As(err, &d) || d.Code != text.LimitExceeded {
		t.Fatalf("expected input limit: %v", err)
	}
}
