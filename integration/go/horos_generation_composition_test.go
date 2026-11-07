package integration_test

import (
	checktypescript "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	formattypescript "github.com/Bitspark/bittype/format/syntax/typescript/provider/go"
	textfixture "github.com/Bitspark/bittype/integration/go/testdata/textfixture"
	parsetypescript "github.com/Bitspark/bittype/parse/syntax/typescript/provider/go"

	"context"
	"errors"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	mappings "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	target "github.com/Bitspark/bittype/transform/syntax/typescript/target/go"

	h "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func externalTS(b target.Builder[h.Language], name string, argument g.Fragment[ts.Language]) (g.Fragment[ts.Language], error) {
	node, err := b.Node(m.Node(ts.KindTypeReference, m.Attributes("name", name), ast.NodeRef{}))
	if err != nil {
		return node, err
	}
	return b.Node(m.Node(ts.KindTypeApplication, nil, ast.NodeRef{}), m.Child[ts.Language]{Role: ts.RoleConstructor, Fragment: node}, m.Child[ts.Language]{Role: ts.RoleArgument, Fragment: argument})
}
func readonlyList(src source.Reader, b target.Builder[h.Language]) g.Mapping[h.Language, ts.Language] {
	return rules.Direct("consumer/readonly-list/1", src, b, func(_ context.Context, r g.UseRequest[h.Language, ts.Language]) (g.Fragment[ts.Language], error) {
		operand, err := src.Operand(r.Want().Use().Use, 1)
		if err != nil {
			return g.Fragment[ts.Language]{}, err
		}
		child, err := r.Child(operand)
		if err != nil {
			return g.Fragment[ts.Language]{}, err
		}
		return externalTS(b, "ReadonlyArray", child.Payload())
	})
}

// This independent consumer receives only its selected source/target APIs.
type tsContribution func(source.Reader, target.Builder[h.Language]) (g.Contribution[h.Language, ts.Language], error)

func (p tsContribution) Open(_ context.Context, r g.MappingRequest[source.Reader, target.Builder[h.Language]]) (g.Contribution[h.Language, ts.Language], error) {
	src, b, err := r.Opened()
	if err != nil {
		return g.Contribution[h.Language, ts.Language]{}, err
	}
	return p(src, b)
}
func selectedTS(input transform.Input, to representation.Language, p tsContribution, replacementID string) (transform.ASTTransformer, error) {
	var registry c.Registry
	request, err := genfixture.Request(&registry, input, to, "selected-ts", target.Provider[h.Language]{}, mappings.Text(rules.RepresentationProfile), mappings.Bool(rules.RepresentationProfile), mappings.Record(rules.RepresentationProfile), mappings.Binder(rules.RepresentationProfile), mappings.Declarations(rules.RepresentationProfile))
	if err != nil {
		return nil, err
	}
	selected, err := c.RegisterMapping[h.Language, ts.Language, source.Reader, target.Builder[h.Language]](&registry, input.Source.Language(), to, p)
	if err != nil {
		return nil, err
	}
	base, err := c.RegisterMapping[h.Language, ts.Language, source.Reader, target.Builder[h.Language]](&registry, input.Source.Language(), to, mappings.List(rules.RepresentationProfile))
	if err != nil {
		return nil, err
	}
	request.Selection.Mappings = append(request.Selection.Mappings, base)
	request.Selection.Mappings = append(request.Selection.Mappings, selected)
	if replacementID != "" {
		key, _ := source.Builtin(input.Source.Language(), "list-of", "")
		request.Selection.Replace = []c.Replacement{{Selector: key, Expected: c.BindingRef{Provider: base, ID: "list/typescript/1"}, With: c.BindingRef{Provider: selected, ID: replacementID}}}
	}
	return registry.Transformer(request)
}
func TestListOverrideComposesThroughGenericPage(t *testing.T) {
	e, from, input := generationInput(t, paginationModule())
	provider := tsContribution(func(src source.Reader, b target.Builder[h.Language]) (g.Contribution[h.Language, ts.Language], error) {
		key, _ := source.Builtin(from, "list-of", "")
		return g.Contribution[h.Language, ts.Language]{Bindings: []g.RepresentationBinding[h.Language, ts.Language]{{Selector: key, Mapping: readonlyList(src, b)}}}, nil
	})
	tree := generated(t, e, from, input, textfixture.New(parsetypescript.Provider{}, checktypescript.Provider{}, formattypescript.Provider{}), func(input transform.Input, f, to representation.Language) (transform.ASTTransformer, error) {
		return selectedTS(input, to, provider, "consumer/readonly-list/1")
	})
	output := formatted(t, textfixture.New(parsetypescript.Provider{}, checktypescript.Provider{}, formattypescript.Provider{}), tree)
	if !strings.Contains(output, "items: ReadonlyArray<T>") || !strings.Contains(output, "Page<Page<Person>>") || strings.Contains(output, "T[]") {
		t.Fatal(output)
	}
	checkGenerationLinks(t, tree)
	compileTSComposition(t, output, "const page: NestedPage = {items: [{items: [{name: 'Ada', active: true}]}]};\nvoid page;\n")
}
func TestSelectedMappingFailureNeverFallsBack(t *testing.T) {
	_, from, input := generationInput(t, paginationModule())
	sentinel := errors.New("selected mapping failed")
	called := 0
	provider := tsContribution(func(src source.Reader, b target.Builder[h.Language]) (g.Contribution[h.Language, ts.Language], error) {
		key, _ := source.Builtin(from, "list-of", "")
		mapping := rules.Direct("consumer/failure/1", src, b, func(context.Context, g.UseRequest[h.Language, ts.Language]) (g.Fragment[ts.Language], error) {
			called++
			return g.Fragment[ts.Language]{}, sentinel
		})
		return g.Contribution[h.Language, ts.Language]{Bindings: []g.RepresentationBinding[h.Language, ts.Language]{{Selector: key, Mapping: mapping}}}, nil
	})
	mapper, err := selectedTS(input, mustLanguage[ts.Language](t, ts.Language{}), provider, "consumer/failure/1")
	if err != nil {
		t.Fatal(err)
	}
	out, err := mapper.Transform(context.Background(), input)
	if !errors.Is(err, sentinel) || out.LanguageID() != "" || called != 1 {
		t.Fatalf("mapping retried or partial output: %v (%d calls)", err, called)
	}
}
func TestLibraryOwnedCellAvoidsBodyExpansion(t *testing.T) {
	library := genModule("cells", genMap(), genEntry("Cell", genType(genTuple("horos/wire", genParam(0)), "T")))
	pin, err := decl.Pin(library)
	if err != nil {
		t.Fatal(err)
	}
	imported := genTuple("horos/decl/import", data.MustEncodeText("cells"), core.NewAtom(pin[:]))
	root := genModule("consumer", genMap(genEntry("cells", imported)), genEntry("Page", genType(genPage(), "T")), genEntry("Cells", genType(genLocal("Page", genImported("cells", "Cell", genScalar("utf8-text"))))))
	e, from, input := generationInput(t, root, library)
	plain, err := tsMapper(input, from, mustLanguage[ts.Language](t, ts.Language{}))
	if err != nil {
		t.Fatal(err)
	}
	if out, err := plain.Transform(context.Background(), input); err == nil || out.LanguageID() != "" {
		t.Fatal("missing Cell mapping accepted")
	}
	calls := 0
	provider := tsContribution(func(src source.Reader, b target.Builder[h.Language]) (g.Contribution[h.Language, ts.Language], error) {
		key, err := source.Named(from, pin.String(), "Cell")
		if err != nil {
			return g.Contribution[h.Language, ts.Language]{}, err
		}
		mapping := rules.Direct("cells/typescript/1", src, b, func(_ context.Context, r g.UseRequest[h.Language, ts.Language]) (g.Fragment[ts.Language], error) {
			calls++
			use := r.Want().Use()
			if len(use.Actuals) != 1 {
				return g.Fragment[ts.Language]{}, g.ErrInvalid
			}
			child, err := r.Child(use.Actuals[0].Use)
			if err != nil {
				return g.Fragment[ts.Language]{}, err
			}
			return externalTS(b, "LibraryCell", child.Payload())
		})
		return g.Contribution[h.Language, ts.Language]{Bindings: []g.RepresentationBinding[h.Language, ts.Language]{{Selector: key, Mapping: mapping}}}, nil
	})
	tree := generated(t, e, from, input, textfixture.New(parsetypescript.Provider{}, checktypescript.Provider{}, formattypescript.Provider{}), func(input transform.Input, f, to representation.Language) (transform.ASTTransformer, error) {
		return selectedTS(input, to, provider, "")
	})
	output := formatted(t, textfixture.New(parsetypescript.Provider{}, checktypescript.Provider{}, formattypescript.Provider{}), tree)
	if calls != 1 || !strings.Contains(output, "Page<LibraryCell<string>>") || strings.Contains(output, "_Cell") {
		t.Fatal(output, calls)
	}
	for _, node := range tree.Nodes {
		if nodeAttribute(node, ts.AttributeName) == "LibraryCell" && len(node.References) != 0 {
			t.Fatal("invented local target")
		}
	}
	compileTSComposition(t, output, "interface LibraryCell<T> { value: T; }\nconst cells: Cells = {items:[{value:'text'}]};\nvoid cells;\n")
}
func compileTSComposition(t *testing.T, source, consumer string) {
	t.Helper()
	if os.Getenv("BITTYPE_NATIVE") != "1" {
		return
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "types.ts"), []byte(source+"\n"+consumer), 0600); err != nil {
		t.Fatal(err)
	}
	command(t, dir, "tsc", "--strict", "--noEmit", "--exactOptionalPropertyTypes", "types.ts")
}
