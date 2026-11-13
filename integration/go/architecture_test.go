package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

var nativeStages = []string{"syntax", "parse/syntax", "check/syntax", "format/syntax"}

type productionPackage struct {
	ImportPath string
	Name       string
	Dir        string
	Standard   bool
	Imports    []string
	GoFiles    []string
	CgoFiles   []string
	// IgnoredGoFiles are the Go files that the build constraints of this run
	// exclude: a build tag, another operating system or architecture, and a
	// cgo file when cgo is off. `go list` leaves their imports out of Imports.
	IgnoredGoFiles []string
}

func productionPackages(t *testing.T) map[string]productionPackage {
	t.Helper()
	return listPackages(t, "../..")
}

// listPackages reads the packages of the module in a directory and their
// dependencies. Each environment entry overrides the test's own.
func listPackages(t *testing.T, directory string, environment ...string) map[string]productionPackage {
	t.Helper()
	command := exec.Command("go", "list", "-deps", "-json", "./...")
	command.Dir = directory
	command.Env = append(os.Environ(), environment...)
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	result := map[string]productionPackage{}
	decoder := json.NewDecoder(bytes.NewReader(output))
	for {
		var p productionPackage
		if err := decoder.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		result[p.ImportPath] = p
	}
	return result
}

// Discover the stage matrix from source directories, not a production language list.
// The union catches an orphan parser as well as an incomplete syntax implementation.
func stageLanguages(t *testing.T) []string {
	t.Helper()
	names := map[string]bool{}
	for _, stage := range append(slices.Clone(nativeStages), "resolve/files", "analyze/files") {
		entries, err := os.ReadDir(filepath.Join("../..", stage))
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if !entry.IsDir() || entry.Name() == "shared" {
				continue
			}
			files, err := filepath.Glob(filepath.Join("../..", strings.TrimPrefix(stagePath(stage, entry.Name()), modulePath+"/"), "*.go"))
			if err != nil {
				t.Fatal(err)
			}
			if len(files) > 0 {
				names[entry.Name()] = true
			}
		}
	}
	var result []string
	for name := range names {
		result = append(result, name)
	}
	slices.Sort(result)
	return result
}

func stagePath(stage, name string) string {
	if stage == "syntax" {
		return modulePath + "/syntax/" + name + "/go"
	}
	return modulePath + "/" + stage + "/" + name + "/provider/go"
}

func TestLanguageStageMatrix(t *testing.T) {
	packages := productionPackages(t)
	languages := stageLanguages(t)
	if len(languages) == 0 {
		t.Fatal("no language implementations")
	}
	qualified := map[string]bool{}
	profiles := map[string][]string{}
	for _, fixture := range nativeConformanceCases() {
		if qualified[fixture.directory] {
			t.Fatalf("duplicate conformance fixture: %s", fixture.directory)
		}
		qualified[fixture.directory] = true
		profiles[fixture.directory] = nativeStages
	}
	semanticQualified := map[string]bool{}
	for _, fixture := range semanticConformanceCases() {
		if semanticQualified[fixture.directory] || fixture.qualify == nil {
			t.Fatalf("duplicate or unqualified semantic/text fixture: %s", fixture.directory)
		}
		qualified[fixture.directory] = true
		semanticQualified[fixture.directory] = true
		profiles[fixture.directory] = []string{"syntax", "parse/syntax", "resolve/files", "check/syntax", "analyze/files", "format/syntax"}
	}
	for _, name := range languages {
		for _, stage := range profiles[name] {
			path := stagePath(stage, name)
			if p, ok := packages[path]; !ok || len(p.GoFiles) == 0 {
				t.Errorf("missing implemented stage %s", path)
			}
		}
		for _, stage := range append(slices.Clone(nativeStages), "resolve/files", "analyze/files") {
			if _, ok := packages[stagePath(stage, name)]; ok && !slices.Contains(profiles[name], stage) {
				t.Errorf("%s has an unqualified capability stage: %s", name, stage)
			}
		}
		if !qualified[name] {
			t.Errorf("%s is missing a tested native-text or semantic/text capability profile", name)
		}
		delete(qualified, name)
	}
	for name := range qualified {
		t.Errorf("conformance fixture has no language implementation: %s", name)
	}
}

// foundationRules is the closed set of foundation and kernel packages: the
// carriers, the contracts, the engine, the shared transform, check and print
// machinery, and the private helpers that serve no single language. Each entry
// is the exact list of module imports of its package, and an unused allowance
// fails.
//
// These stay explicit, and are not a role, for two reasons. Their layering
// among each other is a fact per package: the AST imports nothing, the CST
// only the AST, text diagnostics only AST addresses, representation never
// Source. One rule for "a foundation" would lose every such fact. And the set
// is closed: a new foundation or kernel package is rare and is meant to be an
// edit here.
//
// What the set shares is checked as a rule on top of the entries: a package in
// this table imports no package outside it, whatever its entry lists. So no
// foundation or kernel package can import a vocabulary, a contribution or a
// language's private part.
var foundationRules = map[string][]string{
	// Carriers.
	"trees/ast/go":            {},
	"trees/cst/go":            {"trees/ast/go"},
	"trees/representation/go": {"check/syntax/shared/ast/go", "check/syntax/shared/cst/go", "contracts/language/go", "internal/tree", "semantics/go", "trees/ast/go", "trees/cst/go"},
	"files/go":                {"trees/ast/go"},
	"source/go":               {"contracts/semantic/go", "files/go", "semantics/go", "trees/ast/go", "trees/cst/go", "trees/representation/go"},
	"semantics/go":            {"trees/ast/go"},
	// Contracts.
	"contracts/language/go":     {},
	"contracts/pipeline/go":     {},
	"contracts/text/go":         {"trees/ast/go"},
	"contracts/check/go":        {"contracts/language/go", "trees/ast/go", "trees/representation/go"},
	"contracts/parse/go":        {"contracts/language/go", "contracts/semantic/go", "files/go", "source/go", "trees/ast/go", "trees/cst/go", "trees/representation/go"},
	"contracts/format/go":       {"contracts/generation/go", "contracts/language/go", "contracts/semantic/go", "semantics/go", "source/go", "trees/ast/go", "trees/cst/go", "trees/representation/go"},
	"contracts/print/go":        {"contracts/language/go", "files/go", "source/go", "trees/cst/go", "trees/representation/go"},
	"contracts/transform/go":    {"contracts/language/go", "contracts/semantic/go", "internal/tree", "trees/ast/go", "trees/representation/go"},
	"contracts/semantic/go":     {"contracts/check/go", "contracts/language/go", "files/go", "internal/tree", "semantics/go", "trees/ast/go", "trees/representation/go"},
	"contracts/resolution/go":   {"contracts/language/go", "contracts/semantic/go", "semantics/go", "trees/representation/go"},
	"contracts/analysis/go":     {"contracts/language/go", "contracts/resolution/go", "contracts/semantic/go", "semantics/go", "trees/representation/go"},
	"contracts/generation/go":   {"contracts/check/go", "contracts/semantic/go", "contracts/transform/go", "semantics/go", "trees/ast/go", "trees/representation/go"},
	"contracts/verification/go": {"contracts/semantic/go", "semantics/go", "trees/ast/go", "trees/representation/go"},
	// Kernel.
	"engine/go": {
		"contracts/analysis/go", "contracts/check/go", "contracts/format/go", "contracts/generation/go", "contracts/language/go",
		"contracts/parse/go", "contracts/print/go", "contracts/resolution/go", "contracts/semantic/go", "contracts/transform/go",
		"files/go", "source/go", "transform/syntax/shared/generation/go", "trees/representation/go",
	},
	"transform/syntax/shared/generation/go": {"contracts/check/go", "contracts/generation/go", "contracts/semantic/go", "contracts/text/go", "contracts/transform/go", "trees/ast/go", "trees/representation/go"},
	"transform/syntax/shared/target/go":     {"contracts/generation/go", "trees/ast/go", "trees/representation/go"},
	"check/syntax/shared/ast/go":            {"contracts/text/go", "trees/ast/go"},
	"check/syntax/shared/cst/go":            {"check/syntax/shared/ast/go", "contracts/text/go", "trees/cst/go"},
	"print/syntax/shared/provider/go":       {"check/syntax/shared/cst/go", "contracts/language/go", "contracts/print/go", "trees/cst/go", "trees/representation/go"},
	"print/files/shared/provider/go":        {"contracts/language/go", "contracts/print/go", "files/go", "source/go", "trees/cst/go", "trees/representation/go"},
	// Private helpers of no single language. A shared helper cannot conceal a
	// language: it is bound by the same rule as the rest of the table.
	"internal/tree":      {"trees/ast/go"},
	"internal/grammar":   {"check/syntax/shared/ast/go", "contracts/text/go", "trees/ast/go"},
	"internal/lex":       {"contracts/parse/go", "contracts/text/go", "trees/ast/go"},
	"internal/parseutil": {"check/syntax/shared/ast/go", "contracts/parse/go", "contracts/text/go", "internal/lex", "trees/ast/go"},
	"internal/concrete":  {"check/syntax/shared/cst/go", "contracts/parse/go", "contracts/text/go", "internal/lex", "internal/tree", "trees/ast/go", "trees/cst/go"},
}

// oneOfAKindRules holds the packages whose path gives no role because there is
// one of them and no second is planned: no pattern would be honest. Each entry
// is the exact list of module imports, and an unused allowance fails. A new
// package of such a kind is an edit here, or a new role once there are two.
var oneOfAKindRules = map[string][]string{
	// The horos frontend's private parts. Their names give no role: a role reads
	// internal/<language> and internal/<language>sem, and these are neither.
	// languageOwners binds each to horos. No role allows them, so an importer
	// names them in roleExtensions or in its own entry here.
	"internal/declsyntax": {"check/syntax/shared/ast/go", "contracts/text/go", "internal/grammar", "syntax/horos/go", "trees/ast/go"},
	"internal/horosmodel": {"check/syntax/shared/ast/go", "contracts/semantic/go", "contracts/text/go", "semantics/go", "semantics/horos/go", "syntax/horos/go", "trees/ast/go"},
	// A private composition of the horos Files stages through their contracts,
	// not a substitute language registry.
	"internal/horosfiles": {
		"contracts/analysis/go", "contracts/check/go", "contracts/format/go", "contracts/parse/go", "contracts/print/go", "contracts/resolution/go",
		"contracts/semantic/go", "files/go", "internal/tree", "semantics/go", "semantics/horos/go", "source/go", "syntax/horos/go",
		"trees/ast/go", "trees/cst/go", "trees/representation/go",
	},
	// Verification projections are explicit concrete compositions, not registered
	// language stages. The agnostic association never imports their models.
	"verification/files/source/horos/go": {"contracts/resolution/go", "contracts/semantic/go", "internal/horosmodel", "semantics/horos/go", "syntax/horos/go", "trees/ast/go", "trees/representation/go"},
	// The TypeScript import reader: a private part beside internal/typescript
	// and internal/typescriptsem, with a name that gives neither role.
	"internal/typescriptimports": {"files/go", "internal/typescript", "syntax/typescript/go", "trees/ast/go"},
	// The horos-to-go verification pair and the logos adapter. The pair takes
	// the Go query service's admission type, never the Go analysis provider.
	"verification/files/horos-to-go/go": {"contracts/semantic/go", "contracts/verification/go", "semantics/go", "trees/ast/go", "verification/files/source/horos/go", "semantics/golang/go", "analyze/files/golang/queries/go", "verification/logos/go"},
	"verification/logos/go":             {"contracts/verification/go", "semantics/go"},
	// The reporting tool: a command, its host and the two contracts the host is
	// written against. The host never imports a concrete extractor; extractors
	// have a role.
	"tools/provider-matrix/go":             {"tools/provider-matrix/extract/grammar/go", "tools/provider-matrix/observation/go", "tools/provider-matrix/report/go"},
	"tools/provider-matrix/report/go":      {"tools/provider-matrix/capability/go", "tools/provider-matrix/observation/go"},
	"tools/provider-matrix/capability/go":  {},
	"tools/provider-matrix/observation/go": {},
	// A command over the horos parser and private module lowering.
	"tools/horos-module/go": {"contracts/parse/go", "internal/horosmodel", "internal/lex", "parse/syntax/horos/provider/go", "trees/ast/go"},
}

// languageOwners names the one language that a package with a hand entry
// belongs to, where its path does not say so in a shape a role reads. The
// package then imports nothing of another language, and no package of another
// language imports it, whatever an entry or an extension lists.
var languageOwners = map[string]string{
	"internal/declsyntax":                "horos",
	"internal/horosmodel":                "horos",
	"internal/horosfiles":                "horos",
	"internal/typescriptimports":         "typescript",
	"verification/files/source/horos/go": "horos",
}

// roleExtensions names, per package, the imports that one language needs and
// its role does not give. Languages differ, and the difference is stated here
// once per package; the role stays as narrow as the other languages need. An
// extension must be used, must not repeat the role, and cannot lift a
// prohibition of the role.
var roleExtensions = map[string][]string{
	// The horos checker admits declarations with the private declaration framing.
	"check/syntax/horos/provider/go": {"internal/declsyntax"},
	// The horos text stages read and write source pins and literals with the
	// private model projection. The parser also keeps to the carrier's node
	// bound; the formatter keeps to the lexer's bounds and reports with the
	// shared diagnostics.
	"parse/syntax/horos/provider/go":  {"check/syntax/shared/ast/go", "internal/horosmodel"},
	"format/syntax/horos/provider/go": {"contracts/text/go", "internal/horosmodel", "internal/lex"},
	// The horos resolver has no compiler bridge. It admits the tree with its
	// own checker (a provider exception), reads the private model projection
	// and names the tree it projects into.
	"resolve/files/horos/provider/go": {"check/syntax/horos/provider/go", "internal/horosmodel", "trees/ast/go"},
	// Both horos Files stages are the same private composition.
	"parse/files/horos/provider/go":  {"internal/horosfiles"},
	"format/files/horos/provider/go": {"contracts/parse/go", "internal/horosfiles"},
	// The TypeScript formatter asks the shared identifier grammar whether a
	// name needs quoting.
	"format/syntax/typescript/provider/go": {"internal/grammar"},
	// TypeScript Files formatting places generated declarations in the supplied
	// Files and writes named type imports for their public symbols.
	"format/files/typescript/provider/go": {"contracts/generation/go", "files/go", "internal/typescriptimports", "transform/syntax/typescript/symbol/go"},
	// The TypeScript target publishes those symbols with its output.
	"transform/syntax/typescript/target/go": {"transform/syntax/typescript/symbol/go", "trees/representation/go"},
	// The Go target publishes the symbol of each predeclared type it uses with
	// its output.
	"transform/syntax/golang/target/go": {"transform/syntax/golang/symbol/go", "trees/representation/go"},
	// The Rust target publishes the symbol of each item of the standard
	// library and of each primitive type it uses with its output.
	"transform/syntax/rust/target/go": {"transform/syntax/rust/symbol/go", "trees/representation/go"},
	// Of the common horos mapping parts, only the declaration recipe builds on
	// another: on naming. It states no profile and imports no profile identity.
	// Patterns and naming import no part.
	"transform/syntax/horos/mapping/declarations/go": {"transform/syntax/horos/mapping/naming/go"},
}

// Each production package is checked, including private helpers. This prevents
// a forbidden dependency from being hidden behind an otherwise allowed helper.
//
// A package gets its rule from the role its path gives (packageRole, with
// roleExtensions), from foundationRules or from a hand entry. A package with
// none of these fails. An import of a stage provider package also needs an
// entry in providerExceptions that names the symbols used. The rules
// themselves are in checkArchitecture.
func TestProductionArchitecture(t *testing.T) {
	packages := productionPackages(t)
	spec := architectureSpec{
		Module: modulePath, Foundations: foundationRules, Rules: oneOfAKindRules, Extensions: roleExtensions, Owners: languageOwners,
		Exceptions: providerExceptions, ExternalModels: modelDependencies, ExternalImports: modelExternalImports,
	}
	for _, failure := range checkArchitecture(spec, architecturePackages(t, modulePath, packages)) {
		t.Error(failure)
	}
}

// architecturePackages adds to the `go list` data what the check needs from
// the source of each production package of the module: the symbols it selects
// on its module imports, and the imports and symbols of the files that the
// build constraints of this run exclude.
//
// `go list` reports for one operating system, one architecture, one set of
// tags and one cgo setting. A file outside that build is production code on
// another, so the check reads it too. Otherwise a forbidden import would pass
// behind a build tag, in a _linux.go file on Windows or in a cgo file when cgo
// is off. Test files stay out, as they do in a build.
func architecturePackages(t *testing.T, module string, packages map[string]productionPackage) map[string]architecturePackage {
	t.Helper()
	name := func(path string) (string, bool) {
		p, ok := packages[path]
		return p.Name, ok && strings.HasPrefix(path, module+"/")
	}
	result := map[string]architecturePackage{}
	for path, p := range packages {
		files := append(slices.Clone(p.GoFiles), p.CgoFiles...)
		var excluded []string
		for _, file := range p.IgnoredGoFiles {
			if !strings.HasSuffix(file, "_test.go") {
				excluded = append(excluded, filepath.Join(p.Dir, file))
			}
		}
		data := architecturePackage{Standard: p.Standard, Production: len(files)+len(excluded) > 0, Imports: slices.Clone(p.Imports)}
		if !p.Standard && strings.HasPrefix(path, module+"/") {
			for i, file := range files {
				files[i] = filepath.Join(p.Dir, file)
			}
			var err error
			if data.Uses, data.DotImports, err = parseSymbolUse(append(files, excluded...), name); err != nil {
				t.Fatal(err)
			}
			hidden, err := parseImports(excluded)
			if err != nil {
				t.Fatal(err)
			}
			for _, imported := range hidden {
				// The pseudo-package "C" is kept: with cgo on, `go list` reports
				// it as an import too, so a cgo file fails the same either way.
				listed, known := packages[imported]
				standard := known && listed.Standard || !known && imported != "C" && isStandardImportPath(imported)
				if !standard && !slices.Contains(data.Imports, imported) {
					data.Imports = append(data.Imports, imported)
				}
			}
		}
		result[path] = data
	}
	return result
}

// Only owning frontends/compositions consume external semantic models. A global exemption for
// external modules would let model dependencies leak into generic infrastructure.
var modelDependencies = []string{
	"github.com/Bitspark/horos/core/go",
	"github.com/Bitspark/horos/core/go/declarations",
	"github.com/Bitspark/horos/core/go/internal/declmodel",
	"github.com/Bitspark/horos/core/go/internal/kindmodel",
	"github.com/bitspark/ontos/core/go",
	"github.com/bitspark/ontos/data/go",
	"github.com/bitspark/ontos/codec/go",
	"github.com/Bitspark/logos-contract/contract/go/typed",
}

var modelExternalImports = map[string][]string{
	"verification/logos/go":              {"github.com/Bitspark/logos-contract/contract/go/typed"},
	"tools/horos-module/go":              {"github.com/bitspark/ontos/codec/go"},
	"transform/syntax/horos/source/go":   {"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go", "github.com/bitspark/ontos/codec/go"},
	"parse/syntax/horos/provider/go":     {"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"},
	"format/syntax/horos/provider/go":    {"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"},
	"internal/horosmodel":                {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go", "github.com/bitspark/ontos/codec/go"},
	"resolve/files/horos/provider/go":    {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go"},
	"semantics/horos/go":                 {"github.com/Bitspark/horos/core/go/declarations"},
	"verification/files/source/horos/go": {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"},
	"verification/files/horos-to-go/go":  {"github.com/Bitspark/horos/core/go/declarations", "github.com/Bitspark/logos-contract/contract/go/typed"},
}

func TestFilesDimensionStageCoverage(t *testing.T) {
	fixtures := filesStageFixtures()
	for _, stage := range []string{"parse", "format"} {
		entries, err := os.ReadDir(filepath.Join("../..", stage, "files"))
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			pkg := stagePath(stage+"/files", entry.Name())
			if fixtures[pkg] == nil {
				t.Errorf("Files stage %s lacks a conformance fixture", pkg)
			}
			delete(fixtures, pkg)
		}
	}
	for pkg := range fixtures {
		t.Errorf("Files fixture names absent stage %s", pkg)
	}
}
