package integration_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	target "github.com/Bitspark/bittype/transform/syntax/typescript/target/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// TestHorosReaderRefusesForeignHandles: the horos reader resolves every handle
// against its own capture, including the operands of a view that a rule
// assembles itself and a handle kept from an earlier run. The reader is the
// one that a rule receives in its step. The record contribution here is the
// builtin one with a first step that probes the reader before it names the
// builtin's requests.
func TestHorosReaderRefusesForeignHandles(t *testing.T) {
	_, _, input := generationInput(t, paginationModule())
	var kept g.Use[h.Language]
	var refusals []error
	run := 0
	record := typescriptBundleTable()["record"]
	builtin := record.Contribution
	record.Contribution = func(from representation.Language) (tsContribution, error) {
		contribution, err := builtin(from)
		if err != nil {
			return contribution, err
		}
		name := contribution.Rule.Name
		contribution.Rule.Name = func(ctx context.Context, step horosNameStep) ([]horosRequest, error) {
			reader, err := step.Reader()
			if err != nil {
				return nil, err
			}
			view := step.Subject.Use
			fields, err := reader.Fields(view)
			if err != nil || len(fields) == 0 {
				return nil, errors.Join(errors.New("own fields"), err)
			}
			if run == 1 {
				kept = fields[0].Use
			} else {
				doctored := view.Clone()
				for i := range doctored.Operands {
					doctored.Operands[i].Use = kept
				}
				_, foreign := reader.Fields(doctored)
				refusals = append(refusals, foreign, reader.Unsupported(kept, "kept from an earlier run"))
			}
			return name(ctx, step)
		}
		return contribution, nil
	}
	plan := func() error {
		run++
		var registry c.Registry
		to := nativeChecker(t, typescriptText()).Language()
		table := typescriptBundleTable()
		q, err := genfixture.Compose(&registry, input, to, "foreign-reader", target.Provider{}, table["text"], table["boolean"], record, table["list"], table["declarations"], table["optional"])
		if err != nil {
			t.Fatal(err)
		}
		_, err = registry.PlanGeneration(context.Background(), q, nativeChecker(t, typescriptText()))
		return err
	}
	if err := plan(); err != nil || !kept.Valid() {
		t.Fatal("first run:", err)
	}
	if err := plan(); err != nil || len(refusals) == 0 {
		t.Fatal("second run:", err, len(refusals))
	}
	for _, err := range refusals {
		if !errors.Is(err, g.ErrScope) {
			t.Fatal("foreign handle read:", err)
		}
	}
}

// A test writes a contribution for TypeScript as the builtin mapping packages
// export theirs: a constructor that takes the source language. Each such
// contribution records nothing, so it states the schema of the empty payload.
type tsContribution = g.Contribution[h.Language, ts.Language, source.Reader]
type tsConstructor func(from representation.Language) (tsContribution, error)
type horosRule = g.Rule[h.Language, source.Reader]
type horosNameStep = g.NameStep[h.Language, source.Reader]
type horosDecideStep = g.DecideStep[h.Language, source.Reader]
type horosRequest = g.Request[h.Language]
type tsProjection = g.Projection[ts.Language]
type tsProjected = g.Projected[ts.Language]

// tsTarget is the target that a contribution written for TypeScript states.
var tsTarget = string(ts.Language{}.ID())

// testReference is the reference of a contribution that a test writes. The
// publisher names the test or the fixture, and key is local to it.
func testReference(publisher, key string) g.ContributionReference {
	return g.ContributionReference{
		Identity:       g.ContributionIdentity{Publisher: "example.test/" + publisher, Key: key},
		Meaning:        g.Meaning{Identity: "example.test/" + publisher + "/" + key, Revision: "1"},
		Representation: "1",
	}
}

// decideNothing is the second step of a rule that records nothing.
func decideNothing(context.Context, horosDecideStep) (g.Decision, error) {
	return g.Decision{Payload: f.Empty{}}, nil
}

// nameElement is the first step of a rule for a horos list: it names the
// element, which the source reads from the kernel's view.
func nameElement(_ context.Context, step horosNameStep) ([]horosRequest, error) {
	reader, err := step.Reader()
	if err != nil {
		return nil, err
	}
	element, err := reader.Element(step.Subject.Use)
	if err != nil {
		return nil, err
	}
	return []horosRequest{g.Reference(element)}, nil
}

// emptyPayloadCodec is the codec of the empty payload schema, made once:
// lowering a schema source is not free, and a codec holds no state of a plan.
var emptyPayloadCodec = sync.OnceValues(f.EmptyCodec)

// registerEmpty registers a contribution that states the schema of the empty
// payload, and the codec of that schema. A registry holds one codec for each
// schema and answers ErrAmbiguous to a second one, so that answer means that
// an earlier contribution of the same registry brought the codec.
func registerEmpty[T any](registry *c.Registry, from, to representation.Language, contribution g.Contribution[h.Language, T, source.Reader]) error {
	if contribution.Schema != f.EmptySchema {
		return g.ErrInvalid
	}
	if err := c.RegisterContribution(registry, from, to, contribution); err != nil {
		return err
	}
	codec, err := emptyPayloadCodec()
	if err != nil {
		return err
	}
	if err := registry.RegisterCodec(codec); err != nil && !errors.Is(err, g.ErrAmbiguous) {
		return err
	}
	return nil
}

// registerTS builds and registers a contribution that a test wrote for
// TypeScript, and returns its reference.
func registerTS(registry *c.Registry, from, to representation.Language, p tsConstructor) (g.ContributionReference, error) {
	contribution, err := p(from)
	if err != nil {
		return g.ContributionReference{}, err
	}
	return contribution.Reference, registerEmpty(registry, from, to, contribution)
}

// selectedTS composes the builtin horos-to-typescript set with one
// contribution of the test. With a replacing key the consumer states the
// explicit composition: the contribution of that local key answers the list
// domain and the builtin list contribution is not selected. The key is the
// local key of p. Without one, p stands beside the builtin set.
func selectedTS(p tsConstructor, replacingKey string) requestBuilder {
	return func(registry *c.Registry, input transform.Input, to representation.Language) (c.Request, error) {
		request, err := genfixture.Compose(registry, input, to, "selected-ts", target.Provider{}, typescriptBundles()...)
		if err != nil {
			return c.Request{}, err
		}
		selected, err := registerTS(registry, input.Source.Language(), to, p)
		if err != nil {
			return c.Request{}, err
		}
		if replacingKey != "" {
			list, err := typescriptBundleTable()["list"].Contribution(input.Source.Language())
			if err != nil {
				return c.Request{}, err
			}
			err = registry.ReplaceContribution(list.Reference.Identity, g.ContributionIdentity{Publisher: selected.Identity.Publisher, Key: replacingKey})
			if err != nil {
				return c.Request{}, err
			}
		}
		return request, nil
	}
}

// listContribution is a consumer's contribution for the horos list: the list
// becomes the type name of module, applied to what was selected for the
// element. key is its local key.
func listContribution(key, module, name string) tsConstructor {
	return func(from representation.Language) (tsContribution, error) {
		selector, err := source.Builtin(from, "list-of", "")
		return tsContribution{
			Reference: testReference("consumer", key), Target: tsTarget, Profile: rules.RepresentationProfile,
			Domain: g.ForForm(selector), Schema: f.EmptySchema,
			Rule: horosRule{Name: nameElement, Decide: decideNothing},
			Project: func(_ context.Context, p tsProjection) (tsProjected, error) {
				return tsProjected{Type: model.External(module, name, p.Dependency(0).Type)}, nil
			},
		}, err
	}
}
func TestListOverrideComposesThroughGenericPage(t *testing.T) {
	e, _, input := generationInput(t, paginationModule())
	tree := generated(t, e, input, typescriptText(), selectedTS(listContribution("consumer/readonly-list/1", "", "ReadonlyArray"), "consumer/readonly-list/1"))
	output := formatted(t, typescriptText(), tree)
	if !strings.Contains(output, "items: ReadonlyArray<T>") || !strings.Contains(output, "Page<Page<Person>>") || strings.Contains(output, "T[]") {
		t.Fatal(output)
	}
	checkGenerationLinks(t, tree)
	compileTSComposition(t, output, "const page: NestedPage = {items: [{items: [{name: 'Ada', active: true}]}]};\nvoid page;\n")
}

// TestSelectedMappingFailureNeverFallsBack: the consumer selects its own
// contribution for the list in place of the builtin one, and that contribution
// fails, once in its rule and once in its projector. The generation fails with
// that failure and yields no output. The replaced builtin contribution does
// not answer in its place, and the failing step runs once. An error that is no
// refusal reaches the caller as a rule failure of the selected contribution
// with the error as its cause, and it matches no refusal value.
func TestSelectedMappingFailureNeverFallsBack(t *testing.T) {
	_, _, input := generationInput(t, paginationModule())
	sentinel := errors.New("selected mapping failed")
	const key = "consumer/failure/1"
	for _, step := range []string{"rule", "projector"} {
		called := 0
		fail := func() error {
			called++
			return sentinel
		}
		provider := tsConstructor(func(from representation.Language) (tsContribution, error) {
			selector, err := source.Builtin(from, "list-of", "")
			contribution := tsContribution{
				Reference: testReference("consumer", key), Target: tsTarget, Profile: rules.RepresentationProfile,
				Domain: g.ForForm(selector), Schema: f.EmptySchema,
				Rule: horosRule{Decide: decideNothing},
				Project: func(context.Context, tsProjection) (tsProjected, error) {
					return tsProjected{}, fail()
				},
			}
			if step == "rule" {
				contribution.Rule.Decide = func(context.Context, horosDecideStep) (g.Decision, error) {
					return g.Decision{}, fail()
				}
			}
			return contribution, err
		})
		out, err := runGeneration(t, input, nativeChecker(t, typescriptText()), selectedTS(provider, key))
		var failure *g.RuleFailure
		if !errors.As(err, &failure) || failure.Err != sentinel || failure.Contribution != testReference("consumer", key) || out.Snapshot().Valid() || called != 1 {
			t.Fatal(step, err, called)
		}
		for _, refusal := range []error{g.ErrUnsupported, g.ErrRefusedByDesign, g.ErrRequirementUnsatisfied, g.ErrRequirementUnresolved, g.ErrMissingFacet, g.ErrAmbiguous, g.ErrMissingContribution, g.ErrLimit} {
			if errors.Is(err, refusal) {
				t.Errorf("%s: the failure also matches the refusal value %v", step, refusal)
			}
		}
	}
}
func TestGeneratedNamedImportsTrackOnlyEmittedFragments(t *testing.T) {
	root := genModule("consumer", genMap(), genEntry("Box", genType(genScalar("bool"))), genEntry("Batch", genType(genTuple("horos/list-of", genParam(0)), "T")), genEntry("Names", genType(genTuple("horos/list-of", genScalar("utf8-text")))))
	e, _, input := generationInput(t, root)
	out := generatedOutput(t, e, input, typescriptText(), selectedTS(listContribution("consumer/collection/1", "@example/collections", "Box"), "consumer/collection/1"))
	if len(out.ExternalUses()) != 2 {
		t.Fatal(out.ExternalUses())
	}
	tree := placedTS(t, out, n.PackageTypes{Module: "@example/collections", Declaration: "export interface Box<T> { value:T }"})
	output := formatted(t, typescriptText(), tree)
	checkImportedUse(t, tree, "Box_1", "@example/collections", "Box")
	if strings.Count(output, "import type") != 1 || !strings.Contains(output, "Box_1<T>") || !strings.Contains(output, "Box_1<string>") {
		t.Fatal(output)
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
	if out, err := runGeneration(t, input, nativeChecker(t, typescriptText()), tsRequest); err == nil || out.Snapshot().Valid() {
		t.Fatal("missing Cell mapping accepted")
	}
	calls := 0
	tree := placedTS(t, generatedOutput(t, e, input, typescriptText(), selectedTS(cellMapping(from, pin.String(), &calls), "")), n.PackageTypes{Module: "@bitspark/cell", Declaration: cellDeclarations})
	output := formatted(t, typescriptText(), tree)
	if calls != 1 || !strings.Contains(output, "Page<Cell<string>>") || strings.Contains(output, "_Cell") {
		t.Fatal(output, calls)
	}
	checkImportedUse(t, tree, "Cell", "@bitspark/cell", "Cell")
	compileTSComposition(t, output, "const cells: Cells = {items:[{value:'text'}]};\nvoid cells;\n")
}

// cellMapping is the exact contribution of the library that declares Cell: it
// answers for the uses of that one declaration, so Cell's body is never
// opened. calls counts how often its rule names its request, once for each use.
func cellMapping(from representation.Language, pin string, calls *int) tsConstructor {
	return cellContribution("cells/typescript/1", from, pin, calls)
}

// cellContribution is cellMapping under the local key key.
func cellContribution(key string, from representation.Language, pin string, calls *int) tsConstructor {
	return func(representation.Language) (tsContribution, error) {
		selector, err := source.Named(from, pin, "Cell")
		if err != nil {
			return tsContribution{}, err
		}
		return tsContribution{
			Reference: testReference("cells", key), Target: tsTarget, Profile: rules.RepresentationProfile,
			Domain: g.ForExact(selector), Schema: f.EmptySchema,
			Rule: horosRule{Decide: decideNothing, Name: func(_ context.Context, step horosNameStep) ([]horosRequest, error) {
				reader, err := step.Reader()
				if err != nil {
					return nil, err
				}
				(*calls)++
				// Select the actual for Cell's own binder T, never a position.
				for _, a := range step.Subject.Use.Actuals {
					if p, err := reader.Parameter(a.Binder); err != nil {
						return nil, err
					} else if p.Name == "T" {
						return []horosRequest{g.Reference(a.Use)}, nil
					}
				}
				return nil, g.ErrScope
			}},
			Project: func(_ context.Context, p tsProjection) (tsProjected, error) {
				return tsProjected{Type: model.External("@bitspark/cell", "Cell", p.Dependency(0).Type)}, nil
			},
		}, nil
	}
}
func compileTSComposition(t *testing.T, source, consumer string) {
	t.Helper()
	if os.Getenv("BITTYPE_NATIVE") != "1" {
		return
	}
	dir := t.TempDir()
	writeCellPackage(t, dir)
	if err := os.WriteFile(filepath.Join(dir, "types.ts"), []byte(source+"\n"+consumer), 0600); err != nil {
		t.Fatal(err)
	}
	command(t, dir, "tsc", "--strict", "--noEmit", "--module", "NodeNext", "--exactOptionalPropertyTypes", "types.ts")
}
