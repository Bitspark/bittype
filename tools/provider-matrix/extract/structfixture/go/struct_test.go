// This external-package test is an independent consumer of the public host.
// Its struct-tag vocabulary belongs entirely to this fixture, not the runner,
// the grammar-table extractor or the renderer.
package structfixture_test

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	c "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
	"github.com/Bitspark/bittype/tools/provider-matrix/report/go"
)

const module = "example.test/struct-format"

type structExtractor struct {
	root    c.Symbol
	seen    []c.Inventory
	foreign bool
}

func (s *structExtractor) Selectors() []c.Selector {
	root := s.root
	return []c.Selector{{Definition: &root}}
}

func fact[T any](at c.Span, value T) c.Field[T] {
	return c.Field[T]{State: c.Known, Value: value, At: at}
}

type capturedNode struct {
	node  ast.Node
	file  c.File
	token *token.File
}

func readNode(source c.SourceView, span c.Span) (capturedNode, error) {
	if err := c.ValidateSpan(source, span); err != nil {
		return capturedNode{}, err
	}
	file, err := source.File(span.File)
	if err != nil {
		return capturedNode{}, err
	}
	set := token.NewFileSet()
	parsed, err := parser.ParseFile(set, file.Path, file.Text, parser.SkipObjectResolution)
	if err != nil {
		return capturedNode{}, err
	}
	tokenFile := set.File(parsed.Pos())
	var found ast.Node
	ast.Inspect(parsed, func(node ast.Node) bool {
		if node != nil && tokenFile.Offset(node.Pos()) == span.Start && tokenFile.Offset(node.End()) == span.End {
			found = node
		}
		return true
	})
	if found == nil {
		return capturedNode{}, fmt.Errorf("captured node unavailable")
	}
	return capturedNode{node: found, file: file, token: tokenFile}, nil
}

func (n capturedNode) span(node ast.Node) c.Span {
	return c.Span{File: n.file.Ref, Start: n.token.Offset(node.Pos()), End: n.token.Offset(node.End())}
}

func structure(source c.SourceView, span c.Span) (capturedNode, *ast.StructType, error) {
	node, err := readNode(source, span)
	if err != nil {
		return capturedNode{}, nil, err
	}
	if value, ok := node.node.(*ast.StructType); ok {
		return node, value, nil
	}
	// Named origins are supplied by the host's compiler facts. This follows a
	// true alias and imported declaration without interpreting import spellings.
	facts, err := source.Facts(span)
	if err != nil {
		return capturedNode{}, nil, err
	}
	if facts.Type.Origin == nil {
		return node, nil, nil
	}
	declaration, err := source.Declaration(*facts.Type.Origin)
	if err != nil {
		return capturedNode{}, nil, err
	}
	node, err = readNode(source, declaration.Type)
	if err != nil {
		return capturedNode{}, nil, err
	}
	value, _ := node.node.(*ast.StructType)
	return node, value, nil
}

func tags(field *ast.Field) (reflect.StructTag, error) {
	if field.Tag == nil {
		return "", fmt.Errorf("missing fixture tags")
	}
	text, err := strconv.Unquote(field.Tag.Value)
	return reflect.StructTag(text), err
}

func (s *structExtractor) Extract(ctx context.Context, request c.DefinitionRequest) (c.Inventory, error) {
	if request.Root != s.root {
		return c.Inventory{}, fmt.Errorf("unexpected selected root")
	}
	declaration, err := request.Source.Declaration(request.Root)
	if err != nil {
		return c.Inventory{}, err
	}
	root, forms, err := structure(request.Source, declaration.Type)
	if err != nil {
		return c.Inventory{}, err
	}
	if forms == nil {
		return c.Inventory{}, fmt.Errorf("fixture root is not a struct")
	}
	result := c.Inventory{Extent: c.Decoded}
	for _, form := range forms.Fields.List {
		if err := ctx.Err(); err != nil {
			return c.Inventory{}, err
		}
		meta, err := tags(form)
		if err != nil {
			return c.Inventory{}, err
		}
		key := meta.Get("form")
		if key == "" {
			return c.Inventory{}, fmt.Errorf("missing form key")
		}
		at := root.span(form.Type)
		feature := c.Feature{Key: key, At: at, Roles: fact[[]c.Role](at, nil), Attributes: fact[[]c.Attribute](at, nil), References: fact[[]string](at, nil)}
		node, fields, err := structure(request.Source, at)
		if err != nil {
			return c.Inventory{}, err
		}
		if fields == nil {
			feature.Roles = c.Field[[]c.Role]{State: c.Unknown, At: at, Reason: "Fixture form representation is not a struct"}
			result.Extent = c.PartlyDecoded
		} else {
			for _, field := range fields.Fields.List {
				meta, err := tags(field)
				if err != nil {
					return c.Inventory{}, err
				}
				minimum, err := strconv.Atoi(meta.Get("min"))
				if err != nil {
					return c.Inventory{}, err
				}
				maximum := c.UpperBound{}
				if meta.Get("max") == "*" {
					maximum.Unbounded = true
				} else {
					maximum.Value, err = strconv.Atoi(meta.Get("max"))
					if err != nil {
						return c.Inventory{}, err
					}
				}
				fieldAt := node.span(field.Type)
				role := c.Role{Name: meta.Get("role"), At: fieldAt, Min: fact(fieldAt, minimum), Max: fact(fieldAt, maximum), Kinds: fact(fieldAt, []string{meta.Get("kind")})}
				feature.Roles.Value = append(feature.Roles.Value, role)
			}
		}
		result.Features = append(result.Features, feature)
	}
	if s.foreign {
		result.Features[0].At.File = "unissued/source"
	}
	s.seen = append(s.seen, result)
	return result, nil
}

func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	sources := map[string]string{
		"go.mod": "module " + module + "\n\ngo 1.27.0\n",
		"contracts/language/go/language.go": `package language
type ID string
// bittype:shared
type Provider interface {ID()ID;Name()string}
`,
		"contracts/check/go/check.go": `package check
import "example.test/struct-format/contracts/language/go"
type Provider[I any] interface {language.Provider;Check(I)error}
`,
		"contracts/transform/go/transform.go": `package transform
type ASTTransformer interface {FromID()string;ToID()string;Transform()}
`,
		"trees/ast/go/ast.go": `package ast
// bittype:dimension Syntax
type AST[L any]struct{}
`,
		"syntax/ember/go/language.go": `package ember
import "example.test/struct-format/contracts/language/go"
type Language struct{}
func(Language)ID()language.ID{return "ember"}
func(Language)Name()string{return "Ember"}
`,
		"definition/fields/fields.go": "package fields\ntype Node struct{}\ntype Record struct{ Value Node `role:\"imported\" min:\"0\" max:\"*\" kind:\"node\"` }\n",
		"check/syntax/ember/go/check.go": `package check
import (
 s "example.test/struct-format/syntax/ember/go"
 a "example.test/struct-format/trees/ast/go"
 f "example.test/struct-format/definition/fields"
)
type Provider struct{s.Language}
// bittype:scope docs/scope.md#types
// bittype:grammar Forms
func(Provider)Check(a.AST[s.Language])error{return nil}
type Node struct{}
type Alias = Conditional
` + "type Conditional struct{ Check Node `role:\"check\" min:\"1\" max:\"1\" kind:\"node\"`; True Node `role:\"true\" min:\"1\" max:\"1\" kind:\"node\"` }\n" +
			"type Forms struct{ Conditional Conditional `form:\"conditional-type\"`; Alias Alias `form:\"aliased-type\"`; Imported f.Record `form:\"imported-type\"`; Unsupported func() `form:\"unsupported-form\"` }\n",
		"docs/providers.md": "# Providers\n<!-- generated:providers:begin -->\n<!-- generated:providers:end -->\n<!-- generated:transformers:begin -->\n<!-- generated:transformers:end -->\n",
		"docs/scope.md":     "# Fixture\n## Types\n<!-- bittype:overview scope=types part=features -->\nThe fixture supplies its own struct-tag interpretation.\n<!-- bittype:overview:end -->\n<!-- bittype:overview scope=types part=exclusions -->\nNon-struct forms remain unknown.\n<!-- bittype:overview:end -->\n",
	}
	for path, text := range sources {
		filename := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filename, []byte(text), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestIndependentStructFormatUsesPublicRunner(t *testing.T) {
	root := fixture(t)
	extractor := &structExtractor{root: c.Symbol{Package: module + "/check/syntax/ember/go", Name: "Forms"}}
	config := report.Config{Root: root, SourceURL: "https://example.test/source", Extractors: []report.Binding{{Format: "example.test/struct-tags/v1", Extractor: extractor}}}
	documents, err := report.Generate(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	if len(extractor.seen) != 1 {
		t.Fatalf("expected one definition invocation, got %d", len(extractor.seen))
	}
	result := extractor.seen[0]
	if result.Extent != c.PartlyDecoded || len(result.Features) != 4 {
		t.Fatalf("result: %#v", result)
	}
	if len(result.Features[0].Roles.Value) != 2 || len(result.Features[1].Roles.Value) != 2 || result.Features[2].Roles.Value[0].Name != "imported" || !result.Features[2].Roles.Value[0].Max.Value.Unbounded || result.Features[3].Roles.State != c.Unknown {
		t.Fatalf("independent format facts: %#v", result)
	}
	for _, value := range []string{"conditional-type", "aliased-type", "imported-type", "unsupported-form", "Fixture form representation is not a struct"} {
		if !strings.Contains(documents.Capabilities, value) {
			t.Fatalf("missing %q in report:\n%s", value, documents.Capabilities)
		}
	}
	if !strings.Contains(documents.Matrix, "provider-capabilities.md#") {
		t.Fatal("matrix does not link capability report")
	}
	if strings.Count(documents.Matrix, "[Ember](language-support.md#language-") != 3 {
		t.Fatal("both matrices must link to the exact language overview")
	}
	for _, value := range []string{"Ember", "example.test/struct-tags/v1", "4 declared forms", "conditional-type", "Non-struct forms remain unknown", "partial", "0 associated test declarations"} {
		if !strings.Contains(documents.Overview, value) {
			t.Fatalf("missing %q in compact overview", value)
		}
	}
	if err := report.Write(documents); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(filepath.Join(root, "docs/language-support.md")); err != nil || string(data) != documents.Overview {
		t.Fatal("overview not published", err)
	}
	// Malformed scope passages fail before the caller can publish any document.
	scopePath := filepath.Join(root, "docs/scope.md")
	original, err := os.ReadFile(scopePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(scopePath, append(original, []byte("\n<!-- bittype:overview scope=types part=features -->\nunclosed\n")...), 0600); err != nil {
		t.Fatal(err)
	}
	if partial, err := report.Generate(context.Background(), config); err == nil || partial.Matrix != "" || partial.Capabilities != "" || partial.Overview != "" {
		t.Fatal("malformed excerpt produced publishable output", err)
	}
	if err := os.WriteFile(scopePath, original, 0600); err != nil {
		t.Fatal(err)
	}
	// The same public host rejects ambiguous owners and a result pointing outside
	// its source capture; an extractor cannot grant itself association authority.
	config.Extractors = append(config.Extractors, report.Binding{Format: "example.test/competing-format", Extractor: &structExtractor{root: extractor.root}})
	if _, err := report.Generate(context.Background(), config); err == nil {
		t.Fatal("ambiguous extractor owners accepted")
	}
	config.Extractors = []report.Binding{{Format: "example.test/struct-tags/v1", Extractor: &structExtractor{root: extractor.root, foreign: true}}}
	if _, err := report.Generate(context.Background(), config); err == nil {
		t.Fatal("extractor result with foreign source accepted")
	}
}
