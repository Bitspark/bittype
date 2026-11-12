package integration_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	tsboolean "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/boolean/go"
	tsdeclarations "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/declarations/go"
	tslist "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/list/go"
	tsrecord "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/record/go"
	tstext "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/text/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
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
// against its own capture, including the operands of a view that a library
// assembles itself and a handle kept from an earlier run.
func TestHorosReaderRefusesForeignHandles(t *testing.T) {
	_, _, input := generationInput(t, paginationModule())
	var kept g.Use[h.Language]
	var refusals []error
	run := 0
	record := tsContribution(func(from, to representation.Language) (g.MappingDescription[h.Language, ts.Language, source.Reader], error) {
		pattern, err := patterns.Head(from, "record", "closed")
		return g.MappingDescription[h.Language, ts.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, ts.Language, source.Reader]{{ID: "record/foreign", Input: pattern, Expand: func(ctx context.Context, r g.Expansion[h.Language, ts.Language, source.Reader]) (g.Expression[ts.Language], error) {
			view, err := r.Use()
			if err != nil {
				return g.Expression[ts.Language]{}, err
			}
			s, err := r.Reader()
			if err != nil {
				return g.Expression[ts.Language]{}, err
			}
			fields, err := s.Fields(view)
			if err != nil || len(fields) == 0 {
				return g.Expression[ts.Language]{}, errors.Join(errors.New("own fields"), err)
			}
			if run == 1 {
				kept = fields[0].Use
			} else {
				doctored := view.Clone()
				for i := range doctored.Operands {
					doctored.Operands[i].Use = kept
				}
				_, foreign := s.Fields(doctored)
				refusals = append(refusals, foreign, s.Unsupported(kept, "kept from an earlier run"))
			}
			var members []g.Expression[ts.Language]
			for _, field := range fields {
				ty, err := r.Mapped(field.Use)
				if err != nil {
					return g.Expression[ts.Language]{}, err
				}
				members = append(members, model.Field(field.Name, field.Optional, ty))
			}
			return model.Record(members...), nil
		}}}}, err
	})
	plan := func() error {
		run++
		var registry c.Registry
		to := nativeChecker(t, typescriptText()).Language()
		q, err := genfixture.Request(&registry, input, to, "foreign-reader", target.Provider{}, tstext.Library(), tsboolean.Library(), record, tslist.Library(), tsdeclarations.Library())
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

// Independent libraries advertise data at registration. Expansion is planning only.
type tsContribution func(representation.Language, representation.Language) (g.MappingDescription[h.Language, ts.Language, source.Reader], error)

func (p tsContribution) Describe(from, to representation.Language) (g.MappingDescription[h.Language, ts.Language, source.Reader], error) {
	d, err := p(from, to)
	d.Profile = rules.RepresentationProfile
	return d, err
}
func selectedTS(p tsContribution, replacementID string) requestBuilder {
	return func(registry *c.Registry, input transform.Input, to representation.Language) (c.Request, error) {
		request, err := genfixture.Request(registry, input, to, "selected-ts", target.Provider{}, tstext.Library(), tsboolean.Library(), tsrecord.Library(), tsdeclarations.Library())
		if err != nil {
			return c.Request{}, err
		}
		selected, err := c.RegisterMapping(registry, input.Source.Language(), to, p)
		if err != nil {
			return c.Request{}, err
		}
		base, err := c.RegisterMapping(registry, input.Source.Language(), to, tslist.Library())
		if err != nil {
			return c.Request{}, err
		}
		if replacementID != "" {
			key, _ := source.Builtin(input.Source.Language(), "list-of", "")
			err = registry.Replace(c.Replacement{Selector: key, Expected: c.BindingRef{Provider: base, ID: "list/typescript/1"}, With: c.BindingRef{Provider: selected, ID: replacementID}})
			if err != nil {
				return c.Request{}, err
			}
		}
		return request, nil
	}
}
func listContribution(id, module, name string) tsContribution {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, ts.Language, source.Reader], error) {
		slot := g.TypeSlot("E")
		pattern, err := patterns.Element(from, slot)
		return g.MappingDescription[h.Language, ts.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, ts.Language, source.Reader]{{ID: id, Input: pattern, Output: model.External(module, name, g.Mapped[ts.Language](slot))}}}, err
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
func TestSelectedMappingFailureNeverFallsBack(t *testing.T) {
	_, _, input := generationInput(t, paginationModule())
	sentinel := errors.New("selected mapping failed")
	called := 0
	provider := tsContribution(func(from, to representation.Language) (g.MappingDescription[h.Language, ts.Language, source.Reader], error) {
		pattern, err := patterns.Head(from, "list-of", "")
		return g.MappingDescription[h.Language, ts.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, ts.Language, source.Reader]{{ID: "consumer/failure/1", Input: pattern, Expand: func(context.Context, g.Expansion[h.Language, ts.Language, source.Reader]) (g.Expression[ts.Language], error) {
			called++
			return g.Expression[ts.Language]{}, sentinel
		}}}}, err
	})
	out, err := runGeneration(t, input, nativeChecker(t, typescriptText()), selectedTS(provider, "consumer/failure/1"))
	if !errors.Is(err, sentinel) || out.Snapshot().Valid() || called != 1 {
		t.Fatal(err, called)
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
func cellMapping(from representation.Language, pin string, calls *int) tsContribution {
	return func(_, to representation.Language) (g.MappingDescription[h.Language, ts.Language, source.Reader], error) {
		key, err := source.Named(from, pin, "Cell")
		if err != nil {
			return g.MappingDescription[h.Language, ts.Language, source.Reader]{}, err
		}
		slot := g.TypeSlot("V")
		pattern := g.Pattern[h.Language, source.Reader]{Selector: key, Slots: []g.Slot{slot}, Capture: func(ctx context.Context, src source.Reader, use g.UseView[h.Language]) ([]g.Capture[h.Language], error) {
			(*calls)++
			// Select the actual for Cell's own binder T, never a position.
			for _, a := range use.Actuals {
				if p, err := src.Parameter(a.Binder); err != nil {
					return nil, err
				} else if p.Name == "T" {
					return []g.Capture[h.Language]{{Slot: slot, Use: a.Use}}, nil
				}
			}
			return nil, g.ErrScope
		}}
		return g.MappingDescription[h.Language, ts.Language, source.Reader]{Bindings: []g.BindingDescription[h.Language, ts.Language, source.Reader]{{ID: "cells/typescript/1", Input: pattern, Output: model.External("@bitspark/cell", "Cell", g.Mapped[ts.Language](slot))}}}, nil
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
