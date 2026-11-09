package golang

import (
	"bytes"
	"encoding/json"
	"path"
	"reflect"
	"slices"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/Bitspark/bittype/contracts/semantic/go"
)

// Import environments bound the imported-packages extension of the first Go
// profile. A captured file's package import path is the module path joined with
// the file's directory. Every other imported package is supplied here as Go
// source; nothing is discovered from a module cache, GOPATH or the standard
// library. Package text participates in receipt identity and is never an output.
type ImportEnvironment struct {
	Module   string
	Packages []Package
}

// Package is the supplied source of one dependency package.
type Package struct {
	ImportPath string
	Files      []PackageFile
}

// PackageFile is one source file of a supplied package. Name is a base name ending in ".go".
type PackageFile struct {
	Name string
	Text string
}

const importsKey = "imports-1"

// ImportConfiguration admits imports of captured packages under module and of
// the explicitly supplied packages. The result is normalized, so equal
// environments produce identical settings.
func ImportConfiguration(environment ImportEnvironment) ([]semantic.Setting, error) {
	normalized, err := normalizeEnvironment(environment)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(normalized)
	if err != nil {
		return nil, err
	}
	settings, err := semantic.CaptureSettings(append(Configuration(), semantic.Setting{Key: Key(importsKey), Value: data}), semantic.DefaultLimits())
	if err != nil {
		return nil, err
	}
	// One analysis invocation charges its own configuration and its
	// resolution's against one budget, so two copies must fit with room for
	// the operation, profile, report and requested keys.
	if 2*settingsBytes(settings)+4096 > semantic.DefaultLimits().MaxConfigurationBytes || 2*len(settings) > semantic.DefaultLimits().MaxConfigurationEntries {
		return nil, semantic.ErrLimit
	}
	return settings, nil
}

// Imports decodes the receipt-bound environment. The base Configuration has no
// import environment: one captured package with path PackagePath and no imports.
func Imports(settings []semantic.Setting) (environment ImportEnvironment, enabled bool, err error) {
	base := Configuration()
	if reflect.DeepEqual(settings, base) {
		return ImportEnvironment{}, false, nil
	}
	if len(settings) != len(base)+1 || !reflect.DeepEqual(settings[:len(base)], base) {
		return ImportEnvironment{}, false, semantic.ErrInvalidRequest
	}
	last := settings[len(base)]
	if last.Key != Key(importsKey) || len(last.Value) > semantic.DefaultLimits().MaxConfigurationBytes {
		return ImportEnvironment{}, false, semantic.ErrInvalidRequest
	}
	if err := json.Unmarshal(last.Value, &environment); err != nil {
		return ImportEnvironment{}, false, semantic.ErrInvalidRequest
	}
	want, err := ImportConfiguration(environment)
	if err != nil {
		return ImportEnvironment{}, false, err
	}
	if !bytes.Equal(want[len(base)].Value, last.Value) {
		return ImportEnvironment{}, false, semantic.ErrInvalidRequest
	}
	return environment, true, nil
}

// EvidenceFor is the analysis evidence requirement for exactly this configuration.
func EvidenceFor(settings []semantic.Setting) (semantic.EvidenceRequirement, error) {
	if !MatchesConfiguration(settings) {
		return semantic.EvidenceRequirement{}, semantic.ErrInvalidRequest
	}
	contract := Evidence()
	contract.Configuration = settings
	return contract.Clone(), nil
}

// CapturedImportPath is the import path of a captured file at filePath under module.
func CapturedImportPath(module, filePath string) string {
	dir := path.Dir(filePath)
	if dir == "." {
		return module
	}
	return module + "/" + dir
}

func normalizeEnvironment(e ImportEnvironment) (ImportEnvironment, error) {
	if !ValidImportPath(e.Module) {
		return ImportEnvironment{}, semantic.ErrInvalidRequest
	}
	if len(e.Packages) > 64 {
		return ImportEnvironment{}, semantic.ErrLimit
	}
	out := ImportEnvironment{Module: e.Module, Packages: make([]Package, 0, len(e.Packages))}
	size := len(e.Module)
	for _, p := range e.Packages {
		lower, module := strings.ToLower(p.ImportPath), strings.ToLower(e.Module)
		if !ValidImportPath(p.ImportPath) || lower == module || strings.HasPrefix(lower, module+"/") || slices.Contains(strings.Split(p.ImportPath, "/"), "vendor") {
			// A supplied package must not claim the captured module's paths.
			return ImportEnvironment{}, semantic.ErrInvalidRequest
		}
		if len(p.Files) == 0 || len(p.Files) > 16 {
			return ImportEnvironment{}, semantic.ErrInvalidRequest
		}
		files := slices.Clone(p.Files)
		sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })
		// Go refuses file names in one package that differ only in case.
		names := map[string]bool{}
		for _, f := range files {
			size += len(p.ImportPath) + len(f.Name) + len(f.Text)
			if size > 768<<10 {
				return ImportEnvironment{}, semantic.ErrLimit
			}
			lower := strings.ToLower(f.Name)
			if !validFileName(f.Name) || !utf8.ValidString(f.Text) || names[lower] {
				return ImportEnvironment{}, semantic.ErrInvalidRequest
			}
			names[lower] = true
		}
		out.Packages = append(out.Packages, Package{ImportPath: p.ImportPath, Files: files})
	}
	sort.Slice(out.Packages, func(i, j int) bool { return out.Packages[i].ImportPath < out.Packages[j].ImportPath })
	// Go refuses import paths that differ only in case.
	folded := map[string]bool{}
	for _, p := range out.Packages {
		lower := strings.ToLower(p.ImportPath)
		if folded[lower] {
			return ImportEnvironment{}, semantic.ErrInvalidRequest
		}
		folded[lower] = true
	}
	return out, nil
}

// ValidImportPath accepts the conservative subset of Go module-style import
// paths: slash-separated elements of ASCII letters, digits and "-._~", no empty,
// "." or ".." elements, and a first element containing a dot. The last rule is
// Go's module-path rule; it also keeps supplied packages from naming standard
// library paths, which this profile does not provide.
func ValidImportPath(p string) bool {
	if p == "" || len(p) > 512 {
		return false
	}
	elements := strings.Split(p, "/")
	if !strings.Contains(elements[0], ".") {
		return false
	}
	for _, e := range elements {
		if e == "" || e == "." || e == ".." || strings.HasPrefix(e, ".") || strings.HasSuffix(e, ".") {
			return false
		}
		for _, r := range e {
			if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("-._~", r)) {
				return false
			}
		}
	}
	return true
}

func settingsBytes(settings []semantic.Setting) int {
	n := 0
	for _, s := range settings {
		n += len(s.Key.Owner) + len(s.Key.Name) + len(s.Value)
	}
	return n
}

// validFileName accepts package source files only: no test files and no names
// the go tool ignores.
func validFileName(name string) bool {
	return strings.HasSuffix(name, ".go") && len(name) > 3 && len(name) <= 255 && !strings.ContainsAny(name, "/\\\x00") &&
		!strings.HasPrefix(name, ".") && !strings.HasPrefix(name, "_") && !strings.HasSuffix(name, "_test.go")
}
