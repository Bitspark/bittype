package golang

import (
	"bytes"
	"encoding/json"
	"reflect"
	"slices"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/Bitspark/bittype/contracts/semantic/go"
)

// Import environments bound the imported-packages extension of the first Go
// profile. Module is the module path and nothing else. OutputRoot states where
// the output root sits in that module; the output root is the directory that
// the paths of captured files are relative to. A captured file's package import
// path is the module path, the output root and the file's directory, joined:
// CapturedImportPath. Every other imported package is supplied here as Go
// source; nothing is discovered from a module cache, GOPATH or the standard
// library. Package text participates in receipt identity and is never an output.
type ImportEnvironment struct {
	Module string
	// OutputRoot is the directory of the output root relative to the module
	// root: slash-separated import-path elements such as "generated" or
	// "internal/gen", or ModuleRoot. It must be stated. The empty string is
	// refused and no output root is assumed.
	OutputRoot string
	Packages   []Package
}

// ModuleRoot is the OutputRoot of an environment whose output root is the
// module root. It is a value a caller writes, so it cannot be confused with a
// field that was left empty.
const ModuleRoot = "."

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

// importsKey names the setting of an environment that states its output root.
// Environments under the earlier key "imports-1" stated none; a setting under
// that key is refused and never read with an assumed output root.
const importsKey = "imports-2"

// ImportConfiguration admits imports of captured packages under the output
// root and of the explicitly supplied packages. The result is normalized, so
// equal environments produce identical settings.
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

// CapturedImportPath derives the identity of the package that a captured file
// belongs to: the import path of the environment's output root, joined with the
// file's directory. filePath is the file's path relative to the output root, as
// its Files entry states it. The compiler bridge groups captured files into
// packages by this function, so a caller with the same environment derives the
// identity the bridge used.
//
// It reports false, and no path, when the environment does not state a valid
// module path and output root, when filePath names no file, or when the result
// is not a valid import path. Only Module and OutputRoot are read.
func CapturedImportPath(environment ImportEnvironment, filePath string) (string, bool) {
	importPath, ok := outputRootImportPath(environment.Module, environment.OutputRoot)
	if !ok {
		return "", false
	}
	// A Files path has no empty, "." or ".." element, so its directory is the
	// text before the last slash. Nothing is cleaned: an element that is not an
	// import-path element is refused below.
	slash := strings.LastIndex(filePath, "/")
	if slash >= 0 {
		importPath += "/" + filePath[:slash]
	}
	if filePath[slash+1:] == "" || !ValidImportPath(importPath) {
		return "", false
	}
	return importPath, true
}

// outputRootImportPath is the import path of the output root: the module path
// for ModuleRoot, and otherwise the module path joined with the stated
// directory. Every other spelling is refused, so one directory has one
// spelling: the empty string, an absolute path, a backslash, a trailing slash
// and a "." or ".." element all fail as elements of an import path. A
// directory that the go tool keeps out of a package tree is refused as it is
// for a captured file: testdata, vendor and a leading underscore. A leading dot
// is no import-path element.
func outputRootImportPath(module, outputRoot string) (string, bool) {
	if !ValidImportPath(module) {
		return "", false
	}
	if outputRoot == ModuleRoot {
		return module, true
	}
	root := module + "/" + outputRoot
	if !ValidImportPath(root) {
		return "", false
	}
	for _, element := range strings.Split(outputRoot, "/") {
		if element == "testdata" || element == "vendor" || strings.HasPrefix(element, "_") {
			return "", false
		}
	}
	return root, true
}

func normalizeEnvironment(e ImportEnvironment) (ImportEnvironment, error) {
	root, ok := outputRootImportPath(e.Module, e.OutputRoot)
	if !ok {
		return ImportEnvironment{}, semantic.ErrInvalidRequest
	}
	if len(e.Packages) > 64 {
		return ImportEnvironment{}, semantic.ErrLimit
	}
	out := ImportEnvironment{Module: e.Module, OutputRoot: e.OutputRoot, Packages: make([]Package, 0, len(e.Packages))}
	size := len(e.Module)
	for _, p := range e.Packages {
		lower, captured := strings.ToLower(p.ImportPath), strings.ToLower(root)
		if !ValidImportPath(p.ImportPath) || lower == captured || strings.HasPrefix(lower, captured+"/") || slices.Contains(strings.Split(p.ImportPath, "/"), "vendor") {
			// A supplied package must not claim a path at or below the output
			// root, where the captured packages are. A package of the module
			// outside the output root is hand-written code and may be supplied.
			// A vendor element would make Go resolve its path differently.
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
