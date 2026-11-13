package integration_test

// PROTOTYPE (bittype#196, row 34): the source snapshot, exported from the Go
// side. A reader's second input is the source semantics, and no serialized
// form of it exists. The snapshot is the module pins and, for each
// declaration, the source's description of its type with the position of
// each child, under the source identity that the plan uses.

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Bitspark/bittype/completeness/golang/backend/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/internal/horosmodel"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

func snapshotAtom(v core.Value) string {
	if a, ok := v.(core.Atom); ok {
		return string(a.Bytes())
	}
	return ""
}

// snapshotType renders one canonical horos type expression.
func snapshotType(v core.Value, pin string) (backend.Descriptor, error) {
	tuple, ok := v.(core.Tuple)
	if !ok || tuple.Len() == 0 {
		return backend.Descriptor{}, fmt.Errorf("no type expression")
	}
	local := func(v core.Value) (string, error) {
		t, ok := v.(core.Tuple)
		if !ok || t.Len() != 2 || snapshotAtom(t.At(0)) != "horos/decl/local" {
			return "", fmt.Errorf("the snapshot renders references to local declarations only")
		}
		name, err := data.ReadText(t.At(1))
		return "declaration/" + pin + "/" + name, err
	}
	switch tag := snapshotAtom(tuple.At(0)); tag {
	case "horos/scalar":
		switch snapshotAtom(tuple.At(1)) {
		case "utf8-text":
			return backend.Descriptor{Kind: "text"}, nil
		case "bool":
			return backend.Descriptor{Kind: "bool"}, nil
		}
	case "horos/list-of":
		element, err := snapshotType(tuple.At(1), pin)
		return backend.Descriptor{Kind: "list", Element: &backend.Child{Position: []uint32{1}, Type: element}}, err
	case "horos/record":
		d := backend.Descriptor{Kind: "record"}
		for i := 2; i < tuple.Len(); i++ {
			field, ok := tuple.At(i).(core.Tuple)
			if !ok || field.Len() != 4 {
				return d, fmt.Errorf("no field")
			}
			name, err := data.ReadText(field.At(1))
			if err != nil {
				return d, err
			}
			fieldType, err := snapshotType(field.At(3), pin)
			if err != nil {
				return d, err
			}
			d.Fields = append(d.Fields, backend.Field{Name: name, Optional: snapshotAtom(field.At(2)) == "opt", Position: []uint32{uint32(i), 3}, Type: fieldType})
		}
		return d, nil
	case "horos/decl/local":
		identity, err := local(tuple)
		return backend.Descriptor{Kind: "reference", Declaration: identity}, err
	case "horos/decl/apply":
		identity, err := local(tuple.At(1))
		if err != nil {
			return backend.Descriptor{}, err
		}
		arguments, err := data.ReadList(tuple.At(2))
		d := backend.Descriptor{Kind: "reference", Declaration: identity}
		for _, a := range arguments {
			argument, err := snapshotType(a, pin)
			if err != nil {
				return d, err
			}
			d.Arguments = append(d.Arguments, argument)
		}
		return d, err
	case "horos/decl/param":
		index, err := data.ReadInt(tuple.At(1))
		if err != nil {
			return backend.Descriptor{}, err
		}
		return backend.Descriptor{Kind: "parameter", Index: int(index.Int64())}, nil
	}
	return backend.Descriptor{}, fmt.Errorf("the snapshot does not render %s", snapshotAtom(tuple.At(0)))
}

// protoSnapshot exports the snapshot of one horos module.
func protoSnapshot(t *testing.T, text string) []byte {
	t.Helper()
	parsed, err := parsehoros.Parse(context.Background(), parse.Input{SourceID: "prototype.horos", TreeID: "canonical", Text: text})
	if err != nil {
		t.Fatal(err)
	}
	modules, err := horosmodel.SourceModules(parsed.Abstract)
	if err != nil || len(modules) != 1 {
		t.Fatal("the source is not one module", err)
	}
	digest, err := declarations.Pin(modules[0].Value)
	if err != nil {
		t.Fatal(err)
	}
	module := modules[0].Value.(core.Tuple)
	name, err := data.ReadText(module.At(2))
	if err != nil {
		t.Fatal(err)
	}
	entries, err := data.ReadMap(module.At(4))
	if err != nil {
		t.Fatal(err)
	}
	snapshot := backend.Snapshot{Semantics: source.Semantics, Modules: []backend.Module{{Name: name, Pin: digest.String()}}, Declarations: map[string]backend.Descriptor{}}
	for _, entry := range entries {
		declared, err := data.ReadText(entry[0])
		if err != nil {
			t.Fatal(err)
		}
		declaration, ok := entry[1].(core.Tuple)
		if !ok || snapshotAtom(declaration.At(0)) != "horos/decl/type" {
			continue
		}
		descriptor, err := snapshotType(declaration.At(3), digest.String())
		if err != nil {
			t.Fatal(declared, err)
		}
		snapshot.Declarations["declaration/"+digest.String()+"/"+declared] = descriptor
	}
	out, err := json.MarshalIndent(snapshot, "", " ")
	if err != nil {
		t.Fatal(err)
	}
	return out
}
