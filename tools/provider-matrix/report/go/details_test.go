package report

import (
	"context"
	"strings"
	"testing"

	capability "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
)

func TestDelegateCasesAreRelatedRatherThanWrapperQualification(t *testing.T) {
	p := fixture(t)
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	target := capability.Symbol{Package: p.module + "/provider/ember/go", Name: "emit"}
	a := &annotations{operations: map[capability.Symbol][]authored{}, cases: []caseAssociation{{target: target, directive: directive{values: map[string]string{"feature": "printing", "case": "example", "expects": "shape-preserved"}}, test: testDeclaration{id: "test-id", name: "TestEmit", file: "emit_test.go", variant: "internal"}}}}
	if err := associateCases(m, a); err != nil {
		t.Fatal(err)
	}
	text, err := p.renderDetails(context.Background(), m, a, Config{}, "https://example.test/source", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"Related declaration; its cases and observations do not qualify the public wrapper", "shape-preserved", "No execution observations are included", "No structural definition attached"} {
		if !strings.Contains(text, expected) {
			t.Fatal("missing separation", expected)
		}
	}
	if strings.Contains(text, "percent") || strings.Contains(text, "**pass**") {
		t.Fatal("source-only report invented outcome")
	}
}

func TestAnchorsDoNotDependOnOrderingOrLocations(t *testing.T) {
	p := fixture(t)
	first, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	p.root = "another-root"
	second, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	a, b := operations(first), operations(second)
	if len(a) != len(b) {
		t.Fatal("changed operation count")
	}
	for i := range a {
		if a[i].id != b[i].id {
			t.Fatal("location affected canonical anchor")
		}
	}
}

func TestExtractionSnapshotOwnsNestedRecords(t *testing.T) {
	original := capability.Inventory{Features: []capability.Feature{{Roles: capability.Field[[]capability.Role]{Value: []capability.Role{{Kinds: capability.Field[[]string]{Value: []string{"original"}}}}}}}}
	copy := cloneInventory(original)
	original.Features[0].Roles.Value[0].Kinds.Value[0] = "changed"
	if copy.Features[0].Roles.Value[0].Kinds.Value[0] != "original" {
		t.Fatal("retained extractor slice mutated accepted record")
	}
	var output limitedText
	if _, err := output.Write([]byte(strings.Repeat("x", (64<<20)+1))); err == nil {
		t.Fatal("render output budget ignored")
	}
}
