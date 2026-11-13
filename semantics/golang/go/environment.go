package golang

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	// refused and no output root is assumed. ImportConfiguration refuses a
	// root in a directory that UnsupportedDirectory names with ErrUnsupported,
	// and every other refused root with semantic.ErrInvalidRequest.
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
// equal environments produce identical settings. An output root in a directory
// that UnsupportedDirectory names is refused with ErrUnsupported, every other
// invalid environment with semantic.ErrInvalidRequest, and one that is too
// large with semantic.ErrLimit.
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
	if errors.Is(err, ErrUnsupported) {
		// ImportConfiguration writes no setting for an output root it does not
		// support, so these bytes are no setting.
		return ImportEnvironment{}, false, semantic.ErrInvalidRequest
	}
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
// its Files entry states it. The compiler bridge derives each captured file's
// package once, by this function, and groups the files by the result.
//
// It reports false, and no path, when the environment does not state a valid
// module path and output root, when filePath names no file, when a directory
// of the file is one that UnsupportedDirectory names, or when the result is not
// a valid import path. The 512 bytes of that bound hold the module path, the
// output root and the directory together. The file's own name is not judged:
// the bridge refuses test files and ignored file names itself. Only Module and
// OutputRoot are read.
func CapturedImportPath(environment ImportEnvironment, filePath string) (string, bool) {
	importPath, err := outputRootImportPath(environment.Module, environment.OutputRoot)
	if err != nil {
		return "", false
	}
	// A Files path has no empty, "." or ".." element, so its directory is the
	// text before the last slash. Nothing is cleaned: an element that is not an
	// import-path element is refused below.
	slash := strings.LastIndex(filePath, "/")
	if slash >= 0 {
		directory := filePath[:slash]
		if slices.ContainsFunc(strings.Split(directory, "/"), UnsupportedDirectory) {
			return "", false
		}
		importPath += "/" + directory
	}
	if filePath[slash+1:] == "" || !ValidImportPath(importPath) {
		return "", false
	}
	return importPath, true
}

// UnsupportedDirectory reports whether this profile leaves a directory of this
// name outside the packages it covers: testdata, vendor, and a name that starts
// with "_" or ".". The go tool leaves testdata and a leading underscore or dot
// out of the "./..." pattern and still builds such a package when it is
// imported explicitly, and it refuses an import through a vendor directory.
//
// This is the one declaration of the rule. It applies to every directory below
// the module root, wherever the directory enters: an element of the output
// root is refused by ImportConfiguration, the directory of a captured file
// gives no identity in CapturedImportPath, and the compiler bridge reports such
// a file as unsupported coverage. All three give the answer ErrUnsupported.
func UnsupportedDirectory(element string) bool {
	return element == "testdata" || element == "vendor" || strings.HasPrefix(element, "_") || strings.HasPrefix(element, ".")
}

// windowsDeviceNames are the reserved file names that Go's rule for a path
// element refuses: badWindowsNames in golang.org/x/mod/module, module.go:503-528
// of the pinned go1.27.0 toolchain (x/mod v0.36.1-0.20260813213634-8569e2639ca1).
var windowsDeviceNames = []string{
	"CON", "PRN", "AUX", "NUL",
	"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
	"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
}

// importPathElement is the rule for one element of an import path. It is Go's
// rule, as checkElem states it for an import path in golang.org/x/mod/module
// (module.go:417-482 of the pinned go1.27.0 toolchain): not empty, not ending
// in a dot, ASCII letters, digits and "-._~", and, for the text before the
// first dot, no Windows device name in any case and no tilde followed only by
// digits to its end. It is narrower than Go in two places and wider in none: a
// leading dot and a plus sign are refused.
//
// The elements of an output root follow it. ValidImportPath does not yet: it
// admits the device names and the tilde form for module, supplied and captured
// paths, which is Bitspark/bittype#287.
func importPathElement(element string) bool {
	if element == "" || strings.HasPrefix(element, ".") || strings.HasSuffix(element, ".") {
		return false
	}
	for _, r := range element {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || strings.ContainsRune("-._~", r)) {
			return false
		}
	}
	short, _, _ := strings.Cut(element, ".")
	if slices.ContainsFunc(windowsDeviceNames, func(device string) bool { return strings.EqualFold(device, short) }) {
		return false
	}
	if tilde := strings.LastIndexByte(short, '~'); tilde >= 0 && tilde < len(short)-1 && strings.Trim(short[tilde+1:], "0123456789") == "" {
		return false
	}
	return true
}

// outputRootImportPath is the import path of the output root: the module path
// for ModuleRoot, and otherwise the module path joined with the stated
// directory. One directory has one spelling. A path with an empty, "." or ".."
// element is not clean and is an invalid request: the empty string, an
// absolute path and a trailing slash among them. A clean path is then judged
// as the directory of a captured file is: an element that UnsupportedDirectory
// names gives ErrUnsupported, and after that an element that is no import-path
// element is an invalid request.
func outputRootImportPath(module, outputRoot string) (string, error) {
	if !ValidImportPath(module) {
		return "", semantic.ErrInvalidRequest
	}
	if outputRoot == ModuleRoot {
		return module, nil
	}
	elements := strings.Split(outputRoot, "/")
	if slices.Contains(elements, "") || slices.Contains(elements, ".") || slices.Contains(elements, "..") {
		return "", semantic.ErrInvalidRequest
	}
	if slices.ContainsFunc(elements, UnsupportedDirectory) {
		return "", fmt.Errorf("%w: the output root %q has a testdata or vendor directory, or one that starts with an underscore or a dot", ErrUnsupported, outputRoot)
	}
	if slices.ContainsFunc(elements, func(element string) bool { return !importPathElement(element) }) {
		return "", semantic.ErrInvalidRequest
	}
	// Every element is an import-path element, so only the bound on the length
	// of the whole path can still refuse. A root's validity therefore depends on
	// the length of the module path.
	root := module + "/" + outputRoot
	if !ValidImportPath(root) {
		return "", semantic.ErrInvalidRequest
	}
	return root, nil
}

func normalizeEnvironment(e ImportEnvironment) (ImportEnvironment, error) {
	root, err := outputRootImportPath(e.Module, e.OutputRoot)
	if err != nil {
		return ImportEnvironment{}, err
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
