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
	// Analysis contracts and model admission carry these settings in an evidence
	// requirement with its own budget; refuse an environment they could not carry.
	requirement := Evidence()
	requirement.Configuration = settings
	if err := semantic.ValidateEvidenceRequirements([]semantic.EvidenceRequirement{requirement}); err != nil {
		return nil, err
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
		if !ValidImportPath(p.ImportPath) || p.ImportPath == e.Module || strings.HasPrefix(p.ImportPath, e.Module+"/") {
			// A supplied package must not claim the captured module's paths.
			return ImportEnvironment{}, semantic.ErrInvalidRequest
		}
		if len(p.Files) == 0 || len(p.Files) > 16 {
			return ImportEnvironment{}, semantic.ErrInvalidRequest
		}
		files := slices.Clone(p.Files)
		sort.Slice(files, func(i, j int) bool { return files[i].Name < files[j].Name })
		for i, f := range files {
			size += len(p.ImportPath) + len(f.Name) + len(f.Text)
			if size > 768<<10 {
				return ImportEnvironment{}, semantic.ErrLimit
			}
			if !validFileName(f.Name) || !utf8.ValidString(f.Text) || i > 0 && f.Name == files[i-1].Name {
				return ImportEnvironment{}, semantic.ErrInvalidRequest
			}
		}
		out.Packages = append(out.Packages, Package{ImportPath: p.ImportPath, Files: files})
	}
	sort.Slice(out.Packages, func(i, j int) bool { return out.Packages[i].ImportPath < out.Packages[j].ImportPath })
	for i := 1; i < len(out.Packages); i++ {
		if out.Packages[i].ImportPath == out.Packages[i-1].ImportPath {
			return ImportEnvironment{}, semantic.ErrInvalidRequest
		}
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

func validFileName(name string) bool {
	return strings.HasSuffix(name, ".go") && len(name) > 3 && len(name) <= 255 && !strings.ContainsAny(name, "/\\\x00") && !strings.HasPrefix(name, ".")
}
