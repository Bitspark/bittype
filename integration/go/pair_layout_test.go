package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/tools/payload-schema/artifact/go"
)

// The checks of this file find their cases by path (O46, the finding F15 of
// the two other pairs): a new pair's packages are under them without an edit
// of a list. A hand-written list is a case catalogue that a new pair must
// edit, or silently escapes.

// pairDirectories returns the directories of the ordered pairs below root:
// transform/syntax/<source>-to-<target>.
func pairDirectories(root string) ([]string, error) {
	found, err := filepath.Glob(filepath.Join(root, "transform", "syntax", "*-to-*"))
	if err != nil {
		return nil, err
	}
	var pairs []string
	for _, path := range found {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			pairs = append(pairs, path)
		}
	}
	slices.Sort(pairs)
	return pairs, nil
}

// subdirectories returns the sorted names of the directories in a directory,
// or none where it does not exist.
func subdirectories(path string) []string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	slices.Sort(names)
	return names
}

// vocabularyDirectories returns every vocabulary package of every pair below
// root: transform/syntax/<source>-to-<target>/vocabulary/<name>/go.
func vocabularyDirectories(root string) ([]string, error) {
	pairs, err := pairDirectories(root)
	if err != nil {
		return nil, err
	}
	var found []string
	for _, pair := range pairs {
		for _, name := range subdirectories(filepath.Join(pair, "vocabulary")) {
			found = append(found, filepath.Join(pair, "vocabulary", name, "go"))
		}
	}
	return found, nil
}

// checkPayloadArtifact checks the pinned bootstrap artifact in one vocabulary
// directory against what its schema source gives: the canonical module, the
// file that states its pin and the generated Go types. It reads files only,
// so it needs no import of the package.
func checkPayloadArtifact(directory string) error {
	name := filepath.Base(filepath.Dir(directory))
	read := func(file string, text bool) ([]byte, error) {
		content, err := os.ReadFile(filepath.Join(directory, file))
		if err == nil && text {
			content = bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
		}
		return content, err
	}
	source, err := read(artifact.SourceFile, true)
	if err != nil {
		return fmt.Errorf("%s has no schema source: %w", directory, err)
	}
	module, pin, err := artifact.Module(context.Background(), string(source))
	if err != nil {
		return fmt.Errorf("%s: %w", directory, err)
	}
	types, err := artifact.Types(context.Background(), string(source), name)
	if err != nil {
		return fmt.Errorf("%s: %w", directory, err)
	}
	for _, want := range []struct {
		file    string
		content []byte
		text    bool
	}{{artifact.ModuleFile, module, false}, {artifact.PinFile, []byte(artifact.Pin(name, pin)), true}, {artifact.TypesFile, []byte(types), true}} {
		committed, err := read(want.file, want.text)
		if err != nil {
			return fmt.Errorf("%s: %w", directory, err)
		}
		if !bytes.Equal(committed, want.content) {
			return fmt.Errorf("%s: the committed %s is not what %s gives", directory, want.file, artifact.SourceFile)
		}
	}
	return nil
}

// checkPairLayout checks that every contribution of one pair has its three
// packages: the names below mapping, vocabulary and codec are the same, each
// mapping package exports Contribution and each codec package exports Codec.
func checkPairLayout(pair string) error {
	mappings := subdirectories(filepath.Join(pair, "mapping"))
	if len(mappings) == 0 {
		return fmt.Errorf("%s has no mapping package", pair)
	}
	for part, export := range map[string]string{"vocabulary": "", "codec": "func Codec(", "mapping": "func Contribution("} {
		if names := subdirectories(filepath.Join(pair, part)); !slices.Equal(names, mappings) {
			return fmt.Errorf("%s has the mapping packages %v and the %s packages %v", pair, mappings, part, names)
		}
		for _, name := range mappings {
			if export == "" {
				continue
			}
			files, _ := filepath.Glob(filepath.Join(pair, part, name, "go", "*.go"))
			exported := false
			for _, file := range files {
				content, err := os.ReadFile(file)
				exported = exported || err == nil && !strings.HasSuffix(file, "_test.go") && bytes.Contains(content, []byte("\n"+export))
			}
			if !exported {
				return fmt.Errorf("%s/%s/%s exports no %s", pair, part, name, strings.TrimSuffix(strings.TrimPrefix(export, "func "), "("))
			}
		}
	}
	return nil
}

// Row 38 and O46 (F15).
//
// Shows, for every vocabulary package of every pair that the repository
// holds, found by path: the committed canonical module, the file that states
// its pin and the Go types are, byte for byte, what bittype gives for the
// schema source. Every package of the hand-written list of the codec test is
// among the found ones, so the list names no package that the walk misses.
// Does not show: the round trip of values through a codec, which needs the
// codec's package and typed values and stays with the list
// (TestPayloadBootstrapArtifacts) and, for a pair that is not listed, with
// that pair's own tests.
// Red under: B-schema-source-changed, B-committed-type-changed,
// B-committed-module-changed, B-committed-pin-changed.
func TestPayloadArtifactsOfEveryPairRegenerate(t *testing.T) {
	root := filepath.Join("..", "..")
	found, err := vocabularyDirectories(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(found) < len(payloadCases()) {
		t.Fatalf("the walk finds %d vocabulary packages and the list has %d", len(found), len(payloadCases()))
	}
	for _, c := range payloadCases() {
		if want := filepath.Join(root, "transform", "syntax", "horos-to-"+c.native, "vocabulary", c.name, "go"); !slices.Contains(found, want) {
			t.Errorf("the walk does not find the listed package %s", want)
		}
	}
	for _, directory := range found {
		if err := checkPayloadArtifact(directory); err != nil {
			t.Error(err)
		}
	}
	t.Logf("checked %d vocabulary packages", len(found))
}

// O46 (F15): the walk against the list.
//
// Shows, in a temporary tree that holds one pair that no list of this package
// names: the walk finds its vocabulary and its packages; a complete artifact
// and a complete layout pass; a changed pin, a changed type and a mapping
// without its codec are each reported. The hand-written list of the codec
// test does not name the pair, which is what a list cannot do without an edit.
// Does not show: anything about the two real further pairs. They are on
// another branch.
// Red under: F15-walk-skips-other-pairs.
func TestPairChecksFindAPairThatNoListNames(t *testing.T) {
	root := t.TempDir()
	pair := filepath.Join(root, "transform", "syntax", "elsewhere-to-nowhere")
	vocabulary := filepath.Join(pair, "vocabulary", "thing", "go")
	write := func(path, content string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	const source = "module elsewhere.thing;\n\ntype Payload = { label: text; };\n"
	module, pin, err := artifact.Module(context.Background(), source)
	if err != nil {
		t.Fatal(err)
	}
	types, err := artifact.Types(context.Background(), source, "thing")
	if err != nil {
		t.Fatal(err)
	}
	write(filepath.Join(vocabulary, artifact.SourceFile), source)
	write(filepath.Join(vocabulary, artifact.ModuleFile), string(module))
	write(filepath.Join(vocabulary, artifact.PinFile), artifact.Pin("thing", pin))
	write(filepath.Join(vocabulary, artifact.TypesFile), types)
	write(filepath.Join(pair, "mapping", "thing", "go", "contribution.go"), "package thing\n\nfunc Contribution() {}\n")
	write(filepath.Join(pair, "codec", "thing", "go", "codec.go"), "package thing\n\nfunc Codec() {}\n")

	found, err := vocabularyDirectories(root)
	if err != nil || len(found) != 1 || found[0] != vocabulary {
		t.Fatalf("the walk finds %v: %v", found, err)
	}
	pairs, err := pairDirectories(root)
	if err != nil || len(pairs) != 1 || pairs[0] != pair {
		t.Fatalf("the walk finds the pairs %v: %v", pairs, err)
	}
	for _, c := range payloadCases() {
		if c.native == "nowhere" {
			t.Fatal("the hand-written list names the temporary pair")
		}
	}
	if err := checkPayloadArtifact(vocabulary); err != nil {
		t.Errorf("a complete artifact is reported: %v", err)
	}
	if err := checkPairLayout(pair); err != nil {
		t.Errorf("a complete pair is reported: %v", err)
	}
	// A changed pin.
	write(filepath.Join(vocabulary, artifact.PinFile), artifact.Pin("thing", strings.Repeat("0", len(pin))))
	if err := checkPayloadArtifact(vocabulary); err == nil || !strings.Contains(err.Error(), artifact.PinFile) {
		t.Errorf("a changed pin is not reported: %v", err)
	}
	write(filepath.Join(vocabulary, artifact.PinFile), artifact.Pin("thing", pin))
	// A changed type.
	write(filepath.Join(vocabulary, artifact.TypesFile), strings.Replace(types, "Label", "Other", 1))
	if err := checkPayloadArtifact(vocabulary); err == nil || !strings.Contains(err.Error(), artifact.TypesFile) {
		t.Errorf("a changed type is not reported: %v", err)
	}
	write(filepath.Join(vocabulary, artifact.TypesFile), types)
	// A mapping without its codec.
	if err := os.RemoveAll(filepath.Join(pair, "codec")); err != nil {
		t.Fatal(err)
	}
	if err := checkPairLayout(pair); err == nil || !strings.Contains(err.Error(), "codec") {
		t.Errorf("a mapping without its codec is not reported: %v", err)
	}
}

// O46 (F15), with LAYOUT.md's path form of a component with several
// contributions.
//
// Shows, for every ordered pair that the repository holds, found by path:
// each contribution has its three packages under the same name, its mapping
// package exports Contribution and its codec package exports Codec.
// Does not show: what a contribution states. For the pairs from horos,
// TestBuiltinMappingsStateTheStructuralProfile builds each contribution and
// reads its profile, target and schema, which needs the package's types.
// Red: its check is shown to report on the temporary pair of
// TestPairChecksFindAPairThatNoListNames. No mutant of the repository's own
// tree is one replacement of text: removing a package is not.
func TestEveryPairHasThreePackagesPerContribution(t *testing.T) {
	pairs, err := pairDirectories(filepath.Join("..", ".."))
	if err != nil || len(pairs) < 4 {
		t.Fatalf("the walk finds the pairs %v: %v", pairs, err)
	}
	for _, pair := range pairs {
		if err := checkPairLayout(pair); err != nil {
			t.Error(err)
		}
	}
	t.Logf("checked %d pairs", len(pairs))
}
