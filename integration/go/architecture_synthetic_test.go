package integration_test

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// These tests feed checkArchitecture hand-built package sets. Each failing set
// breaks one rule and must produce that rule's message and no other.

const syntheticModule = "example.test/m"

// syntheticSet is a package set with its rules. Paths are module-relative.
type syntheticSet struct {
	spec     architectureSpec
	packages map[string]architecturePackage
}

func newSyntheticSet() *syntheticSet {
	return &syntheticSet{
		spec: architectureSpec{
			Module: syntheticModule, Foundations: map[string][]string{}, Rules: map[string][]string{},
			Extensions: map[string][]string{}, Owners: map[string]string{}, ExternalImports: map[string][]string{},
		},
		packages: map[string]architecturePackage{},
	}
}

// add declares a production package and what it imports. It gets no entry in
// any table: it has the role of its path, or none.
func (s *syntheticSet) add(path string, imports ...string) {
	p := architecturePackage{Production: true, Uses: map[string][]string{}}
	for _, dependency := range imports {
		p.Imports = append(p.Imports, syntheticModule+"/"+dependency)
	}
	s.packages[syntheticModule+"/"+path] = p
}

// listed declares a production package whose hand entry allows exactly what it
// imports.
func (s *syntheticSet) listed(path string, imports ...string) {
	s.add(path, imports...)
	s.spec.Rules[path] = slices.Clone(imports)
}

// leaves declares packages without module imports, each with an empty hand entry.
func (s *syntheticSet) leaves(paths ...string) {
	for _, path := range paths {
		s.listed(path)
	}
}

func (s *syntheticSet) uses(path, dependency string, symbols ...string) {
	s.packages[syntheticModule+"/"+path].Uses[syntheticModule+"/"+dependency] = symbols
}

func (s *syntheticSet) dotImports(path, dependency string) {
	p := s.packages[syntheticModule+"/"+path]
	p.DotImports = append(p.DotImports, syntheticModule+"/"+dependency)
	s.packages[syntheticModule+"/"+path] = p
}

func (s *syntheticSet) except(importer, provider string, symbols ...string) {
	s.spec.Exceptions = append(s.spec.Exceptions, providerException{Importer: importer, Provider: provider, Symbols: symbols, Reason: "synthetic"})
}

func (s *syntheticSet) expect(t *testing.T, want ...string) {
	t.Helper()
	got := checkArchitecture(s.spec, s.packages)
	if !slices.Equal(got, want) {
		t.Errorf("failures:\n got: %q\nwant: %q", got, want)
	}
}

const (
	syntheticChecker = "check/syntax/x/provider/go"
	syntheticParser  = "parse/syntax/x/provider/go"
)

// providerSet is a conforming set: a parser that calls its own language's
// Check under a provider exception, and a tool that composes both providers.
func providerSet() *syntheticSet {
	s := newSyntheticSet()
	s.leaves("trees/ast/go")
	s.listed(syntheticChecker, "trees/ast/go")
	s.listed(syntheticParser, "trees/ast/go", syntheticChecker)
	s.uses(syntheticParser, syntheticChecker, "Check")
	s.except(syntheticParser, syntheticChecker, "Check")
	return s
}

func TestArchitectureAcceptsConformingProviderSet(t *testing.T) {
	providerSet().expect(t)
}

func TestArchitectureProviderImportNeedsException(t *testing.T) {
	s := providerSet()
	s.listed("format/syntax/x/provider/go", syntheticChecker)
	s.uses("format/syntax/x/provider/go", syntheticChecker, "Check")
	s.expect(t, "format/syntax/x/provider/go imports stage provider check/syntax/x/provider/go (uses Check) without a provider exception")
}

func TestArchitectureProviderSymbolOutsideException(t *testing.T) {
	s := providerSet()
	s.uses(syntheticParser, syntheticChecker, "Check", "Provider")
	s.expect(t, "parse/syntax/x/provider/go uses Provider from check/syntax/x/provider/go, but its provider exception names only Check")
}

func TestArchitectureUnusedProviderException(t *testing.T) {
	s := providerSet()
	s.except(syntheticChecker, syntheticParser, "Parse")
	s.expect(t, "unused provider exception: check/syntax/x/provider/go does not import parse/syntax/x/provider/go")
}

func TestArchitectureUnusedProviderExceptionSymbol(t *testing.T) {
	s := providerSet()
	s.spec.Exceptions[0].Symbols = []string{"Check", "Provider"}
	s.expect(t, "unused symbol in provider exception parse/syntax/x/provider/go -> check/syntax/x/provider/go: Provider")
}

func TestArchitectureUnusedAllowance(t *testing.T) {
	s := providerSet()
	s.spec.Rules[syntheticParser] = append(s.spec.Rules[syntheticParser], "trees/cst/go")
	s.expect(t, "unused allowance: parse/syntax/x/provider/go may import trees/cst/go and does not")
}

// A deleted package leaves its entries behind until a cleanup removes them.
// They do not fail, and they are not unused: that is said of a package that exists.
func TestArchitectureEntriesOfAbsentPackagesPass(t *testing.T) {
	s := providerSet()
	s.spec.Rules["legacy/go"] = []string{"trees/ast/go", syntheticChecker}
	s.spec.ExternalImports["legacy/go"] = []string{"example.org/model"}
	s.except("legacy/go", syntheticChecker, "Check")
	s.spec.Foundations["contracts/gone/go"] = []string{"trees/ast/go"}
	s.spec.Extensions["analyze/files/x/provider/go"] = []string{"trees/ast/go"}
	s.spec.Owners["internal/gone"] = "x"
	// A directory that keeps only test files is not a production package either.
	s.packages[syntheticModule+"/format/syntax/x/provider/go"] = architecturePackage{Imports: []string{syntheticModule + "/" + syntheticChecker}}
	s.spec.Rules["format/syntax/x/provider/go"] = []string{"trees/ast/go", syntheticChecker}
	s.except("format/syntax/x/provider/go", syntheticChecker, "Check")
	s.expect(t)

	// The same entries are checked once the package exists.
	s.add("legacy/go", "trees/ast/go")
	s.expect(t,
		"unused external allowance: legacy/go may import example.org/model and does not",
		"unused allowance: legacy/go may import check/syntax/x/provider/go and does not",
		"unused provider exception: legacy/go does not import check/syntax/x/provider/go",
	)
}

func TestArchitectureDotImport(t *testing.T) {
	s := providerSet()
	s.dotImports(syntheticParser, "trees/ast/go")
	s.expect(t, "parse/syntax/x/provider/go dot-imports trees/ast/go, which hides the symbols it uses")
}

func TestArchitectureCompositionsStayFree(t *testing.T) {
	s := providerSet()
	// A tool keeps its per-package rule and needs no provider exception.
	s.listed("tools/x/go", syntheticParser, syntheticChecker)
	s.uses("tools/x/go", syntheticParser, "Parse", "Provider")
	s.uses("tools/x/go", syntheticChecker, "Check", "Provider")
	// The test compositions are not classified at all.
	for _, composition := range []string{"integration/go", "verification/integration/go"} {
		s.add(composition, syntheticParser, syntheticChecker, "unlisted/go")
		s.uses(composition, syntheticParser, "Parse")
		s.dotImports(composition, syntheticChecker)
	}
	s.expect(t)

	s.except("tools/x/go", syntheticParser, "Parse", "Provider")
	s.expect(t, "provider exception tools/x/go -> parse/syntax/x/provider/go is not needed: tools/x/go is a composition")
}

func TestArchitectureUnclassifiedPackage(t *testing.T) {
	s := providerSet()
	s.add("misc/go")
	// One segment more or less than a role's shape gives no role, and neither
	// does a component the layout does not have.
	s.add("semantics/x/extra/go")
	s.add("syntax/x/extra/go")
	s.add("transform/syntax/x/other/go")
	s.add("transform/syntax/x/mapping/part/extra/go")
	s.add("transform/syntax/shared/extra/go")
	s.add("verification/files/x/go")
	s.add("tools/x/go")
	s.expect(t,
		"unclassified production package: misc/go",
		"unclassified production package: semantics/x/extra/go",
		"unclassified production package: syntax/x/extra/go",
		"unclassified production package: tools/x/go",
		"unclassified production package: transform/syntax/shared/extra/go",
		"unclassified production package: transform/syntax/x/mapping/part/extra/go",
		"unclassified production package: transform/syntax/x/other/go",
		"unclassified production package: verification/files/x/go",
	)
}

func TestArchitectureForbiddenImport(t *testing.T) {
	s := providerSet()
	s.leaves("trees/cst/go")
	s.packages[syntheticModule+"/"+syntheticChecker] = architecturePackage{Production: true, Imports: []string{syntheticModule + "/trees/ast/go", syntheticModule + "/trees/cst/go"}}
	s.expect(t, "check/syntax/x/provider/go must not import example.test/m/trees/cst/go")
}

func TestArchitectureExternalDependencies(t *testing.T) {
	s := providerSet()
	s.spec.ExternalModels = []string{"example.org/model"}
	s.spec.ExternalImports[syntheticChecker] = []string{"example.org/model"}
	s.packages["example.org/model"] = architecturePackage{Production: true}
	s.packages["strings"] = architecturePackage{Standard: true, Production: true}
	for _, path := range []string{syntheticChecker, syntheticParser} {
		p := s.packages[syntheticModule+"/"+path]
		p.Imports = append(p.Imports, "strings", "example.org/model")
		s.packages[syntheticModule+"/"+path] = p
	}
	s.packages["example.org/other"] = architecturePackage{Production: true}
	s.expect(t,
		"unclassified external production dependency: example.org/other",
		"parse/syntax/x/provider/go must not import example.org/model",
	)

	// An external allowance that its package does not use fails.
	s = providerSet()
	s.spec.ExternalModels = []string{"example.org/model"}
	s.packages["example.org/model"] = architecturePackage{Production: true}
	s.spec.ExternalImports[syntheticChecker] = []string{"example.org/model"}
	s.expect(t, "unused external allowance: check/syntax/x/provider/go may import example.org/model and does not")
}

// An import of a module package that the set does not hold fails: `go list`
// leaves out a package whose every file the build constraints exclude, so
// nothing would check it.
func TestArchitectureImportOfUnlistedPackage(t *testing.T) {
	s := providerSet()
	s.imports(syntheticParser, "internal/elsewhere")
	s.spec.Rules[syntheticParser] = append(s.spec.Rules[syntheticParser], "internal/elsewhere")
	s.expect(t, "parse/syntax/x/provider/go imports internal/elsewhere, which the package set does not hold: no file of it is built here, so nothing checks it")

	// A role does not admit it either, also where its path is one the role names.
	s = layoutSet()
	delete(s.packages, syntheticModule+"/transform/syntax/x/model/go")
	s.expect(t,
		"transform/syntax/x/target/go imports transform/syntax/x/model/go, which the package set does not hold: no file of it is built here, so nothing checks it",
		"transform/syntax/y-to-x/mapping/record/go imports transform/syntax/x/model/go, which the package set does not hold: no file of it is built here, so nothing checks it",
		"transform/syntax/y-to-x/mapping/text/go imports transform/syntax/x/model/go, which the package set does not hold: no file of it is built here, so nothing checks it",
	)
}

func TestArchitectureMalformedProviderException(t *testing.T) {
	s := providerSet()
	s.spec.Exceptions[0].Reason = " "
	s.expect(t, "provider exception parse/syntax/x/provider/go -> check/syntax/x/provider/go gives no reason")

	s = providerSet()
	s.except(syntheticParser, syntheticChecker, "Check")
	s.expect(t, "provider exception parse/syntax/x/provider/go -> check/syntax/x/provider/go is listed twice")

	// A blank import cannot be admitted: an entry names at least one symbol,
	// and every named symbol must be used.
	s = providerSet()
	s.uses(syntheticParser, syntheticChecker)
	s.spec.Exceptions[0].Symbols = nil
	s.expect(t, "provider exception parse/syntax/x/provider/go -> check/syntax/x/provider/go names no symbol")

	// The shared printers are kernel packages, and an entry is not a second
	// allow-list for other imports.
	s = providerSet()
	s.except(syntheticParser, "trees/ast/go", "Node")
	s.expect(t,
		"provider exception parse/syntax/x/provider/go -> trees/ast/go does not name a stage provider or query service package",
		"unused symbol in provider exception parse/syntax/x/provider/go -> trees/ast/go: Node",
	)
}

// The tests below classify packages by role. Nothing in a layout set has a
// hand entry: every package is in the foundation table or takes the rule of
// its path, as a new package of a known shape does.

// syntheticFoundations stands for the closed foundation and kernel table.
var syntheticFoundations = []string{
	"contracts/analysis/go", "contracts/check/go", "contracts/format/go", "contracts/generation/go", "contracts/language/go",
	"contracts/parse/go", "contracts/print/go", "contracts/resolution/go", "contracts/semantic/go", "contracts/text/go",
	"contracts/transform/go", "files/go", "internal/concrete", "internal/grammar", "internal/lex", "internal/parseutil",
	"internal/tree", "semantics/go", "source/go", "transform/syntax/shared/contribution/go", "transform/syntax/shared/publisher/go", "transform/syntax/shared/target/go", "trees/ast/go", "trees/representation/go",
}

// languageLayout is every package shape of one language, each importing all
// that its role permits. "<l>" stands for the language.
var languageLayout = map[string][]string{
	"syntax/<l>/go":    {"contracts/language/go", "trees/ast/go", "trees/cst/go"},
	"semantics/<l>/go": {"contracts/analysis/go", "contracts/resolution/go", "contracts/semantic/go", "internal/<l>", "semantics/go", "syntax/<l>/go", "trees/ast/go", "trees/representation/go"},
	"internal/<l>":     {"internal/grammar"},
	"internal/<l>sem":  {"contracts/parse/go", "contracts/semantic/go", "files/go", "parse/syntax/<l>/provider/go", "semantics/go", "semantics/<l>/go", "syntax/<l>/go", "trees/ast/go"},

	"check/syntax/<l>/provider/go": {"contracts/check/go", "contracts/text/go", "internal/<l>", "internal/grammar", "internal/tree", "syntax/<l>/go", "trees/ast/go"},
	"parse/syntax/<l>/provider/go": {
		"check/syntax/<l>/provider/go", "contracts/parse/go", "contracts/text/go", "internal/<l>", "internal/concrete", "internal/lex",
		"internal/parseutil", "internal/tree", "syntax/<l>/go", "trees/ast/go", "trees/cst/go",
	},
	"format/syntax/<l>/provider/go": {"check/syntax/<l>/provider/go", "contracts/format/go", "internal/<l>", "internal/concrete", "internal/tree", "syntax/<l>/go", "trees/ast/go", "trees/cst/go"},

	"resolve/files/<l>/provider/go": {"contracts/resolution/go", "contracts/semantic/go", "internal/<l>sem", "semantics/go", "semantics/<l>/go", "syntax/<l>/go", "trees/representation/go"},
	"analyze/files/<l>/provider/go": {"contracts/analysis/go", "contracts/semantic/go", "semantics/go", "semantics/<l>/go", "syntax/<l>/go"},
	"analyze/files/<l>/queries/go":  {"contracts/semantic/go", "internal/<l>sem", "semantics/<l>/go", "syntax/<l>/go", "trees/representation/go"},
	"parse/files/<l>/provider/go": {
		"contracts/analysis/go", "contracts/check/go", "contracts/parse/go", "contracts/print/go", "contracts/resolution/go", "contracts/semantic/go",
		"semantics/go", "source/go", "syntax/<l>/go", "trees/ast/go", "trees/cst/go", "trees/representation/go",
	},
	"format/files/<l>/provider/go": {
		"contracts/analysis/go", "contracts/check/go", "contracts/format/go", "contracts/print/go", "contracts/resolution/go", "contracts/semantic/go",
		"semantics/go", "semantics/<l>/go", "source/go", "syntax/<l>/go", "trees/ast/go", "trees/cst/go", "trees/representation/go",
	},

	"transform/syntax/<l>/model/go":  {"contracts/generation/go", "syntax/<l>/go"},
	"transform/syntax/<l>/symbol/go": {"contracts/generation/go", "internal/<l>", "semantics/go", "syntax/<l>/go", "trees/ast/go", "trees/representation/go"},
	"transform/syntax/<l>/source/go": {
		"contracts/analysis/go", "contracts/generation/go", "contracts/semantic/go", "contracts/text/go", "contracts/transform/go",
		"semantics/go", "semantics/<l>/go", "syntax/<l>/go", "trees/ast/go", "trees/representation/go",
	},
	"transform/syntax/<l>/target/go": {"contracts/generation/go", "internal/<l>", "syntax/<l>/go", "transform/syntax/<l>/model/go", "transform/syntax/shared/target/go", "trees/ast/go"},

	"transform/syntax/<l>/mapping/go":         {},
	"transform/syntax/<l>/mapping/pattern/go": {"contracts/generation/go", "syntax/<l>/go", "transform/syntax/<l>/source/go", "transform/syntax/shared/contribution/go"},
	"transform/syntax/<l>/mapping/naming/go":  {"contracts/generation/go", "syntax/<l>/go", "transform/syntax/<l>/source/go"},
}

// pairImports is what one mapping contribution of <s>-to-<n> may import.
var pairImports = []string{
	"contracts/generation/go", "syntax/<n>/go", "syntax/<s>/go", "transform/syntax/<n>/model/go", "transform/syntax/shared/contribution/go",
	"transform/syntax/<s>/mapping/declarations/go", "transform/syntax/<s>/mapping/go", "transform/syntax/<s>/mapping/naming/go",
	"transform/syntax/<s>/mapping/pattern/go", "transform/syntax/<s>/source/go", "trees/representation/go",
}

func replaced(values []string, replacer *strings.Replacer) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = replacer.Replace(value)
	}
	return result
}

// foundation declares a package of the foundation table with its exact imports.
func (s *syntheticSet) foundation(path string, imports ...string) {
	s.add(path, imports...)
	s.spec.Foundations[path] = slices.Clone(imports)
}

// imports adds imports to a declared package and touches no table.
func (s *syntheticSet) imports(path string, dependencies ...string) *syntheticSet {
	p, ok := s.packages[syntheticModule+"/"+path]
	if !ok {
		panic("undeclared synthetic package: " + path)
	}
	for _, dependency := range dependencies {
		p.Imports = append(p.Imports, syntheticModule+"/"+dependency)
	}
	s.packages[syntheticModule+"/"+path] = p
	return s
}

// language declares every package of one language. The only lines it adds to a
// table are the provider exceptions for the calls into its own checker and
// parser, and the extension of the one common mapping part that builds on others.
func (s *syntheticSet) language(name string) {
	replacer := strings.NewReplacer("<l>", name)
	for path, imports := range languageLayout {
		s.add(replacer.Replace(path), replaced(imports, replacer)...)
	}
	checker, parser := "check/syntax/"+name+"/provider/go", "parse/syntax/"+name+"/provider/go"
	for _, importer := range []string{parser, "format/syntax/" + name + "/provider/go"} {
		s.uses(importer, checker, "Check")
		s.except(importer, checker, "Check")
	}
	s.uses("internal/"+name+"sem", parser, "Parse")
	s.except("internal/"+name+"sem", parser, "Parse")

	adapter, naming := "transform/syntax/"+name+"/mapping/go", "transform/syntax/"+name+"/mapping/naming/go"
	declarations := "transform/syntax/" + name + "/mapping/declarations/go"
	s.add(declarations, append(replaced(languageLayout["transform/syntax/<l>/mapping/naming/go"], replacer), adapter, naming)...)
	s.spec.Extensions[declarations] = []string{adapter, naming}
}

// The external models that the codec role brings with it.
var codecModels = []string{"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"}

// pair declares two mapping contributions of one ordered pair, each with the
// vocabulary of its family, its own vocabulary and its codec.
func (s *syntheticSet) pair(source, native string) {
	replacer := strings.NewReplacer("<s>", source, "<n>", native)
	base := "transform/syntax/" + source + "-to-" + native
	for _, name := range []string{"text", "record"} {
		family := "transform/syntax/" + source + "/family/" + name + "/go"
		vocabulary := base + "/vocabulary/" + name + "/go"
		s.add(family, "contracts/generation/go")
		s.add(vocabulary, "contracts/generation/go", "transform/syntax/shared/publisher/go", family)
		s.add(base+"/codec/"+name+"/go", "contracts/generation/go", payloadShared, vocabulary)
		s.importsExternal(base+"/codec/"+name+"/go", codecModels...)
		s.add(base+"/mapping/"+name+"/go", append(replaced(pairImports, replacer), vocabulary)...)
	}
}

// layoutSet is a conforming module: the foundation table, two complete
// languages x and y, and the mapping contributions of both ordered pairs.
func layoutSet() *syntheticSet {
	s := newSyntheticSet()
	for _, path := range syntheticFoundations {
		s.foundation(path)
	}
	s.foundation("trees/cst/go", "trees/ast/go")
	s.foundation("engine/go", "contracts/check/go", "contracts/parse/go", "trees/representation/go")
	s.language("x")
	s.language("y")
	s.listed(payloadShared, "contracts/generation/go")
	s.spec.ExternalModels = append(s.spec.ExternalModels, codecModels...)
	s.pair("x", "y")
	s.pair("y", "x")
	return s
}

const (
	onlyFoundations      = "a foundation or kernel package imports only foundation and kernel packages"
	noContribution       = "a vocabulary never imports a contribution"
	noSemanticHelper     = "a vocabulary never imports a semantic helper"
	noStage              = "a lexical helper depends on no stage"
	noSourceOrMapping    = "a native target reads no source and calls no mapping"
	noOtherMapping       = "a mapping contribution never imports another mapping contribution"
	noTreeConstructor    = "a mapping contribution never imports a native tree constructor"
	queryWithoutResolver = "a query service does not import its resolver"
	mappingOwnVocabulary = "a mapping contribution imports no codec and no vocabulary but its own"
	mappingNoMeaning     = "a mapping contribution reaches its family's meaning through its own vocabulary and no payload codec"
	vocabularyAlone      = "a vocabulary package imports no mapping, no codec and no vocabulary but its own"
	codecAlone           = "a codec package imports no mapping, no codec and no vocabulary but its own"
	familyByVocabulary   = "only the vocabulary of a contribution imports the vocabulary of its family"
	languageY            = "it belongs to another language, y"
)

// refused is the failure of an import that a rule prohibits for a reason.
func refused(importer, imported, reason string) string {
	return importer + " must not import " + imported + ": " + reason
}

// unlisted is the failure of an import that the importer's rule does not name.
func unlisted(importer, imported string) string {
	return importer + " must not import " + syntheticModule + "/" + imported
}

// unadmitted is the second failure of an unlisted import of a stage provider
// or a query service.
func unadmitted(importer, provider string) string {
	kind, _ := selectedKind(provider)
	return importer + " imports " + kind + " " + provider + " (uses no symbol) without a provider exception"
}

type roleBreak struct {
	importer, imported string
	want               []string
}

func prohibited(importer, imported, reason string) roleBreak {
	return roleBreak{importer, imported, []string{refused(importer, imported, reason)}}
}

func notAllowed(importer, imported string) roleBreak {
	want := []string{unlisted(importer, imported)}
	if _, selected := selectedKind(imported); selected {
		want = append(want, unadmitted(importer, imported))
	}
	return roleBreak{importer, imported, want}
}

// expectBreaks adds each import to a fresh conforming layout and expects
// exactly the failures of the rule it breaks.
func expectBreaks(t *testing.T, breaks ...roleBreak) {
	t.Helper()
	for _, b := range breaks {
		s := layoutSet()
		s.imports(b.importer, b.imported)
		if got := checkArchitecture(s.spec, s.packages); !slices.Equal(got, b.want) {
			t.Errorf("%s importing %s:\n got: %q\nwant: %q", b.importer, b.imported, got, b.want)
		}
	}
}

func TestArchitectureAcceptsConformingLayout(t *testing.T) {
	layoutSet().expect(t)
}

// A third language, a pair into it, one more mapping of an existing pair and
// one more common mapping part are classified without a line in the foundation
// table, the hand entries or the owners. The language brings only what
// language() adds for every language.
func TestArchitectureNewPackagesOfKnownShapesNeedNoEntry(t *testing.T) {
	s := layoutSet()
	rules, foundations, extensions := len(s.spec.Rules), len(s.spec.Foundations), len(s.spec.Extensions)
	s.language("z")
	s.pair("x", "z")
	s.add("transform/syntax/x/mapping/defaults/go", "contracts/generation/go", "syntax/x/go")
	s.add("transform/syntax/x-to-y/mapping/optional/go",
		"contracts/generation/go", "transform/syntax/x/mapping/defaults/go", "transform/syntax/x/mapping/pattern/go", "transform/syntax/y/model/go")
	s.expect(t)
	if len(s.spec.Rules) != rules || len(s.spec.Owners) != 0 || len(s.spec.Foundations) != foundations || len(s.spec.Extensions) != extensions+1 {
		t.Errorf("tables grew: %d hand entries, %d owners, %d foundations, %d extensions", len(s.spec.Rules), len(s.spec.Owners), len(s.spec.Foundations), len(s.spec.Extensions))
	}

	// The same packages are policed by their roles.
	s.imports("parse/syntax/z/provider/go", "syntax/x/go")
	s.imports("transform/syntax/x-to-z/mapping/text/go", "transform/syntax/z/target/go")
	s.imports("transform/syntax/x/mapping/defaults/go", "transform/syntax/x/mapping/pattern/go")
	s.expect(t,
		refused("parse/syntax/z/provider/go", "syntax/x/go", "it belongs to another language, x"),
		refused("transform/syntax/x-to-z/mapping/text/go", "transform/syntax/z/target/go", noTreeConstructor),
		unlisted("transform/syntax/x/mapping/defaults/go", "transform/syntax/x/mapping/pattern/go"),
	)
}

// Rule 1. A foundation or kernel package imports what its entry lists, and
// never a package outside the table, also when its entry lists one.
func TestArchitectureFoundationTable(t *testing.T) {
	expectBreaks(t,
		prohibited("engine/go", "syntax/x/go", onlyFoundations),
		prohibited("engine/go", "parse/syntax/x/provider/go", onlyFoundations),
		prohibited("engine/go", "transform/syntax/x-to-y/mapping/text/go", onlyFoundations),
		prohibited("contracts/generation/go", "transform/syntax/x/model/go", onlyFoundations),
		prohibited("internal/grammar", "internal/x", onlyFoundations),
		prohibited("internal/lex", "internal/xsem", onlyFoundations),
		// The layering inside the table is each entry's own list.
		notAllowed("trees/ast/go", "trees/cst/go"),
		notAllowed("contracts/text/go", "contracts/parse/go"),
		notAllowed("engine/go", "internal/lex"),
	)

	s := layoutSet().imports("engine/go", "syntax/x/go")
	s.spec.Foundations["engine/go"] = append(s.spec.Foundations["engine/go"], "syntax/x/go")
	s.expect(t, refused("engine/go", "syntax/x/go", onlyFoundations))

	s = layoutSet()
	s.spec.Foundations["engine/go"] = append(s.spec.Foundations["engine/go"], "contracts/print/go")
	s.expect(t, "unused allowance: engine/go may import contracts/print/go and does not")

	s = layoutSet()
	s.spec.Rules["engine/go"] = slices.Clone(s.spec.Foundations["engine/go"])
	s.expect(t, "engine/go has a foundation entry and a hand entry")

	// A package that looks like a foundation and is not in the table is not one.
	s = layoutSet()
	s.add("contracts/extra/go")
	s.expect(t, "unclassified production package: contracts/extra/go")
}

// Rule 2. A vocabulary imports foundations and its own language's vocabulary,
// never a contribution and never another language.
func TestArchitectureVocabularyRoles(t *testing.T) {
	expectBreaks(t,
		prohibited("syntax/x/go", "check/syntax/x/provider/go", noContribution),
		prohibited("syntax/x/go", "syntax/y/go", languageY),
		notAllowed("syntax/x/go", "contracts/check/go"),
		notAllowed("syntax/x/go", "semantics/x/go"),
		notAllowed("syntax/x/go", "internal/x"),

		prohibited("semantics/x/go", "resolve/files/x/provider/go", noContribution),
		prohibited("semantics/x/go", "analyze/files/x/queries/go", noContribution),
		prohibited("semantics/x/go", "transform/syntax/x/source/go", noContribution),
		prohibited("semantics/x/go", "internal/xsem", noSemanticHelper),
		prohibited("semantics/x/go", "internal/y", languageY),
		prohibited("semantics/x/go", "semantics/y/go", languageY),
		notAllowed("semantics/x/go", "contracts/parse/go"),

		prohibited("transform/syntax/x/model/go", "transform/syntax/x/target/go", noContribution),
		prohibited("transform/syntax/x/model/go", "transform/syntax/x-to-y/mapping/text/go", noContribution),
		prohibited("transform/syntax/x/model/go", "transform/syntax/y/model/go", languageY),
		notAllowed("transform/syntax/x/model/go", "trees/ast/go"),
		notAllowed("transform/syntax/x/model/go", "semantics/x/go"),

		prohibited("transform/syntax/x/symbol/go", "format/files/x/provider/go", noContribution),
		prohibited("transform/syntax/x/symbol/go", "syntax/y/go", languageY),
		notAllowed("transform/syntax/x/symbol/go", "transform/syntax/x/model/go"),
	)
}

// The roles of the contribution contract (docs/contribution-contract.md, O25,
// O27, O37): the vocabulary of a family, the vocabulary of a contribution and
// its payload codec. No contribution imports a contribution, a mapping
// imports of its bundle only its own vocabulary, and only a codec reaches the
// shared part of the codec role.
func TestArchitectureContributionBundleRoles(t *testing.T) {
	const (
		family     = "transform/syntax/x/family/text/go"
		vocabulary = "transform/syntax/x-to-y/vocabulary/text/go"
		codec      = "transform/syntax/x-to-y/codec/text/go"
		mapping    = "transform/syntax/x-to-y/mapping/text/go"
	)
	expectBreaks(t,
		// A mapping imports no codec, no other vocabulary, no family vocabulary
		// and not the shared part of the codec role.
		prohibited(mapping, codec, mappingOwnVocabulary),
		prohibited(mapping, "transform/syntax/x-to-y/vocabulary/record/go", mappingOwnVocabulary),
		prohibited(mapping, family, mappingNoMeaning),
		prohibited(mapping, payloadShared, mappingNoMeaning),

		// A vocabulary of a contribution imports its own family and nothing
		// else of a bundle. It is a vocabulary: it never imports a contribution.
		prohibited(vocabulary, mapping, vocabularyAlone),
		prohibited(vocabulary, codec, vocabularyAlone),
		prohibited(vocabulary, "transform/syntax/x-to-y/vocabulary/record/go", vocabularyAlone),
		prohibited(vocabulary, "transform/syntax/x/family/record/go", familyByVocabulary),
		notAllowed(vocabulary, payloadShared),
		notAllowed(vocabulary, "transform/syntax/y/model/go"),
		prohibited(vocabulary, "transform/syntax/x/source/go", noContribution),

		// A codec imports its own vocabulary and the shared part, and neither
		// the mapping, nor another codec, nor a family vocabulary.
		prohibited(codec, mapping, codecAlone),
		prohibited(codec, "transform/syntax/x-to-y/codec/record/go", codecAlone),
		prohibited(codec, "transform/syntax/x-to-y/vocabulary/record/go", codecAlone),
		prohibited(codec, family, familyByVocabulary),
		notAllowed(codec, "transform/syntax/x/source/go"),
		notAllowed(codec, "transform/syntax/y/model/go"),

		// The vocabulary of a family imports the contract and nothing else.
		prohibited(family, mapping, noContribution),
		notAllowed(family, "transform/syntax/x/mapping/go"),
		notAllowed(family, "syntax/x/go"),

		// No kernel package and no native target reaches the codec role.
		prohibited("engine/go", payloadShared, onlyFoundations),
		notAllowed("transform/syntax/y/target/go", payloadShared),
	)

	// A codec imports the ontos value packages by its role, and never horos.
	layoutSet().importsExternal(codec, "github.com/Bitspark/horos/core/go").expect(t, codec+" must not import github.com/Bitspark/horos/core/go")
	layoutSet().importsExternal(mapping, codecModels[0]).expect(t, mapping+" must not import "+codecModels[0])
	layoutSet().importsExternal(vocabulary, codecModels[0]).expect(t, vocabulary+" must not import "+codecModels[0])
}

// Rule 3 for the Syntax stages. Each stage has its own list; parser and
// formatter import their own checker under a provider exception, never each
// other, and no stage imports another language.
func TestArchitectureSyntaxStageProviderRoles(t *testing.T) {
	const checker, parser, formatter = "check/syntax/x/provider/go", "parse/syntax/x/provider/go", "format/syntax/x/provider/go"
	expectBreaks(t,
		notAllowed(parser, formatter),
		notAllowed(formatter, parser),
		notAllowed(checker, parser),
		notAllowed(checker, "trees/cst/go"),
		notAllowed(checker, "internal/lex"),
		notAllowed(checker, "internal/concrete"),
		notAllowed(formatter, "internal/lex"),
		notAllowed(formatter, "contracts/text/go"),
		notAllowed(formatter, "contracts/parse/go"),
		notAllowed(parser, "contracts/semantic/go"),
		notAllowed(parser, "semantics/x/go"),
		notAllowed(parser, "internal/xsem"),
		notAllowed(parser, "transform/syntax/x/model/go"),
		prohibited(parser, "syntax/y/go", languageY),
		prohibited(parser, "internal/y", languageY),
		prohibited(parser, "check/syntax/y/provider/go", languageY),
		prohibited(checker, "syntax/y/go", languageY),
		prohibited(formatter, "transform/syntax/y/model/go", languageY),
	)

	// The role admits the import of the own checker, and the provider exception
	// still names the symbols.
	s := layoutSet()
	s.uses(parser, checker, "Check", "Provider")
	s.expect(t, "parse/syntax/x/provider/go uses Provider from check/syntax/x/provider/go, but its provider exception names only Check")
}

// Rule 3 for the Files stages.
func TestArchitectureFilesStageProviderRoles(t *testing.T) {
	const resolver, analyzer = "resolve/files/x/provider/go", "analyze/files/x/provider/go"
	const parser, formatter, queries = "parse/files/x/provider/go", "format/files/x/provider/go", "analyze/files/x/queries/go"
	expectBreaks(t,
		notAllowed(analyzer, resolver),
		notAllowed(analyzer, "internal/xsem"),
		notAllowed(analyzer, "contracts/resolution/go"),
		notAllowed(analyzer, queries),
		notAllowed(resolver, analyzer),
		notAllowed(resolver, "contracts/analysis/go"),
		notAllowed(resolver, "parse/syntax/x/provider/go"),
		prohibited(resolver, "internal/ysem", languageY),
		prohibited(resolver, "semantics/y/go", languageY),
		prohibited(analyzer, "syntax/y/go", languageY),
		// A Files parser or formatter composes its collaborators through contracts.
		notAllowed(parser, "parse/syntax/x/provider/go"),
		notAllowed(parser, resolver),
		notAllowed(parser, "contracts/format/go"),
		notAllowed(parser, "internal/xsem"),
		notAllowed(parser, "semantics/x/go"),
		notAllowed(formatter, "contracts/parse/go"),
		notAllowed(formatter, "files/go"),
		prohibited(formatter, "transform/syntax/y/symbol/go", languageY),

		prohibited(queries, resolver, queryWithoutResolver),
		notAllowed(queries, analyzer),
		notAllowed(queries, "trees/ast/go"),
		notAllowed(queries, "contracts/analysis/go"),
		notAllowed(queries, "semantics/go"),
		prohibited(queries, "internal/ysem", languageY),
	)

	// A stage the layout does not have gives no role.
	s := layoutSet()
	s.add("check/files/x/provider/go", "contracts/check/go")
	s.add("lint/syntax/x/provider/go", "trees/ast/go")
	s.expect(t,
		"unclassified production package: check/files/x/provider/go",
		"unclassified production package: lint/syntax/x/provider/go",
	)
}

// Rule 3 for a source reader and a native target.
func TestArchitectureSourceAndTargetRoles(t *testing.T) {
	const source, target = "transform/syntax/x/source/go", "transform/syntax/x/target/go"
	expectBreaks(t,
		prohibited(target, source, noSourceOrMapping),
		prohibited(target, "transform/syntax/y/source/go", noSourceOrMapping),
		prohibited(target, "transform/syntax/y-to-x/mapping/text/go", noSourceOrMapping),
		prohibited(target, "transform/syntax/y/mapping/pattern/go", noSourceOrMapping),
		prohibited(target, "transform/syntax/y/target/go", languageY),
		prohibited(target, "transform/syntax/y/model/go", languageY),
		notAllowed(target, "check/syntax/x/provider/go"),
		notAllowed(target, "semantics/x/go"),
		notAllowed(target, "trees/representation/go"),

		prohibited(source, "syntax/y/go", languageY),
		prohibited(source, "transform/syntax/y/model/go", languageY),
		notAllowed(source, target),
		notAllowed(source, "transform/syntax/x/model/go"),
		notAllowed(source, "analyze/files/x/provider/go"),
		notAllowed(source, "transform/syntax/x/mapping/pattern/go"),
	)
}

// The common mapping parts of one source language import no other part and no
// native language. The one part that builds on two others says so in an
// extension, and without it the same imports fail.
func TestArchitectureCommonMappingParts(t *testing.T) {
	const (
		adapter      = "transform/syntax/x/mapping/go"
		pattern      = "transform/syntax/x/mapping/pattern/go"
		naming       = "transform/syntax/x/mapping/naming/go"
		declarations = "transform/syntax/x/mapping/declarations/go"
	)
	expectBreaks(t,
		notAllowed(pattern, naming),
		notAllowed(pattern, adapter),
		notAllowed(pattern, declarations),
		notAllowed(naming, pattern),
		notAllowed(naming, adapter),
		notAllowed(adapter, pattern),
		notAllowed(adapter, declarations),
		notAllowed(declarations, pattern),
		prohibited(adapter, "transform/syntax/y/model/go", languageY),
		prohibited(pattern, "syntax/y/go", languageY),
		prohibited(declarations, "transform/syntax/x-to-y/mapping/text/go", languageY),
		notAllowed(pattern, "transform/syntax/x/model/go"),
		notAllowed(naming, "transform/syntax/x/target/go"),
		notAllowed(adapter, "semantics/x/go"),
	)

	s := layoutSet()
	delete(s.spec.Extensions, declarations)
	s.expect(t, unlisted(declarations, adapter), unlisted(declarations, naming))
}

// Rule 3 for a mapping contribution: its own pair, its source language's
// common parts, never another mapping and never a tree constructor.
func TestArchitectureMappingContributionRole(t *testing.T) {
	const text = "transform/syntax/x-to-y/mapping/text/go"
	expectBreaks(t,
		prohibited(text, "transform/syntax/x-to-y/mapping/record/go", noOtherMapping),
		prohibited(text, "transform/syntax/y-to-x/mapping/text/go", noOtherMapping),
		prohibited(text, "transform/syntax/y/target/go", noTreeConstructor),
		prohibited(text, "transform/syntax/x/target/go", noTreeConstructor),
		// Both languages of the pair are its scope, and each side has its part.
		notAllowed(text, "transform/syntax/y/mapping/pattern/go"),
		notAllowed(text, "transform/syntax/y/source/go"),
		notAllowed(text, "transform/syntax/x/model/go"),
		notAllowed(text, "semantics/x/go"),
		notAllowed(text, "internal/y"),
		notAllowed(text, "parse/syntax/x/provider/go"),
		notAllowed(text, "transform/syntax/shared/target/go"),
		notAllowed(text, "trees/ast/go"),
	)

	s := layoutSet()
	s.language("z")
	s.imports(text, "transform/syntax/z/model/go")
	s.expect(t, refused(text, "transform/syntax/z/model/go", "it belongs to another language, z"))

	// One segment more or less is not a mapping contribution.
	s = layoutSet()
	s.add("transform/syntax/x-to-y/mapping/go")
	s.add("transform/syntax/x-to-y/mapping/text/extra/go")
	s.add("transform/syntax/x-to-/mapping/text/go")
	s.expect(t,
		"unclassified production package: transform/syntax/x-to-/mapping/text/go",
		"unclassified production package: transform/syntax/x-to-y/mapping/go",
		"unclassified production package: transform/syntax/x-to-y/mapping/text/extra/go",
	)
}

// The private parts of one language: a lexical helper knows spelling only, and
// a semantic helper reparses with its own parser under a provider exception.
// Both exist only for a language that exists.
func TestArchitectureLanguageHelperRoles(t *testing.T) {
	expectBreaks(t,
		prohibited("internal/x", "contracts/parse/go", noStage),
		prohibited("internal/x", "contracts/text/go", noStage),
		prohibited("internal/x", "check/syntax/x/provider/go", noStage),
		prohibited("internal/x", "internal/y", languageY),
		notAllowed("internal/x", "trees/ast/go"),
		notAllowed("internal/x", "syntax/x/go"),
		notAllowed("internal/x", "internal/lex"),
		notAllowed("internal/x", "internal/xsem"),

		prohibited("internal/xsem", "parse/syntax/y/provider/go", languageY),
		prohibited("internal/xsem", "internal/ysem", languageY),
		prohibited("internal/xsem", "semantics/y/go", languageY),
		notAllowed("internal/xsem", "resolve/files/x/provider/go"),
		notAllowed("internal/xsem", "check/syntax/x/provider/go"),
		notAllowed("internal/xsem", "contracts/resolution/go"),
		notAllowed("internal/xsem", "internal/x"),
	)

	s := layoutSet()
	s.uses("internal/xsem", "parse/syntax/x/provider/go", "Parse", "Provider")
	s.expect(t, "internal/xsem uses Provider from parse/syntax/x/provider/go, but its provider exception names only Parse")

	// A helper of no language has no role: a generic helper is an entry in the
	// foundation table, and that is an edit.
	s = layoutSet()
	s.add("internal/common", "internal/grammar")
	s.add("internal/wsem", "semantics/go")
	s.add("internal/x/deeper")
	s.expect(t,
		"unclassified production package: internal/common (as a language's helper it needs syntax/common/go, which does not exist)",
		"unclassified production package: internal/wsem (as a language's helper it needs syntax/w/go, which does not exist)",
		"unclassified production package: internal/x/deeper",
	)
	s.foundation("internal/common", "internal/grammar")
	delete(s.packages, syntheticModule+"/internal/wsem")
	delete(s.packages, syntheticModule+"/internal/x/deeper")
	s.expect(t)
}

// drop removes one import from a declared package.
func (s *syntheticSet) drop(path, dependency string) *syntheticSet {
	p := s.packages[syntheticModule+"/"+path]
	p.Imports = slices.DeleteFunc(slices.Clone(p.Imports), func(imported string) bool { return imported == syntheticModule+"/"+dependency })
	s.packages[syntheticModule+"/"+path] = p
	return s
}

// A role names its language's private part by the path: internal/<language>.
// It grants that import only where the imported package is that language's
// helper by its own role. A package cannot give itself the private part of
// another owner by choosing a name.
func TestArchitectureRoleGrantsOnlyALanguagesHelper(t *testing.T) {
	// internal/xmodel is a private part of x with a hand entry. A vocabulary
	// named after it would be allowed "internal/" + its own name.
	s := layoutSet()
	s.listed("internal/xmodel", "semantics/x/go")
	s.add("semantics/xmodel/go", "semantics/go", "internal/xmodel")
	s.expect(t, "semantics/xmodel/go must not import internal/xmodel: its role grants internal/<language> to a language's helper only, and internal/xmodel is not one")
	// With an owner the private part is refused as another language's already.
	s.spec.Owners["internal/xmodel"] = "x"
	s.expect(t, refused("semantics/xmodel/go", "internal/xmodel", "it belongs to another language, x"))

	// The same for every role that names a helper by the path. The importer
	// brings other failures of its own; the one about the helper is among them.
	for _, c := range []struct{ importer, helper, pattern string }{
		{"check/syntax/xmodel/provider/go", "internal/xmodel", "internal/<language>"},
		{"parse/syntax/xmodel/provider/go", "internal/xmodel", "internal/<language>"},
		{"format/syntax/xmodel/provider/go", "internal/xmodel", "internal/<language>"},
		{"transform/syntax/xmodel/symbol/go", "internal/xmodel", "internal/<language>"},
		{"transform/syntax/xmodel/target/go", "internal/xmodel", "internal/<language>"},
		{"resolve/files/xmo/provider/go", "internal/xmosem", "internal/<language>sem"},
		{"analyze/files/xmo/queries/go", "internal/xmosem", "internal/<language>sem"},
	} {
		s = layoutSet()
		s.listed(c.helper, "semantics/x/go")
		s.add(c.importer, c.helper)
		got := checkArchitecture(s.spec, s.packages)
		want := c.importer + " must not import " + c.helper + ": its role grants " + c.pattern + " to a language's helper only, and " + c.helper + " is not one"
		if !slices.Equal(got, []string{want}) {
			t.Errorf("%s importing %s:\n got: %q\nwant: %q", c.importer, c.helper, got, want)
		}
	}

	// Declaring the language too makes the private part that language's
	// lexical helper, and the helper's prohibitions then bind it.
	s = layoutSet()
	s.listed("internal/xmodel", "contracts/semantic/go")
	s.add("syntax/xmodel/go")
	s.add("semantics/xmodel/go", "semantics/go", "internal/xmodel")
	s.expect(t, refused("internal/xmodel", "contracts/semantic/go", noStage))

	// A foundation helper of the same name as a language is not that
	// language's helper either.
	s = layoutSet()
	s.add("syntax/grammar/go")
	s.add("semantics/grammar/go", "semantics/go", "internal/grammar")
	s.expect(t, "semantics/grammar/go must not import internal/grammar: its role grants internal/<language> to a language's helper only, and internal/grammar is not one")
}

// A role is never wider than what its packages use together: each import it
// permits must be imported by at least one package that the role classifies.
func TestArchitectureUnusedRoleAllowance(t *testing.T) {
	const pattern = "trees/cst/go"
	s := layoutSet().drop("format/syntax/x/provider/go", pattern)
	s.expect(t)
	s.drop("format/syntax/y/provider/go", pattern)
	s.expect(t, "unused role allowance: no package of format/syntax/<language>/provider/go imports trees/cst/go")

	// The allowance is one across languages: a language's own package counts.
	s = layoutSet().drop("check/syntax/x/provider/go", "internal/x")
	s.expect(t)
	s.drop("check/syntax/y/provider/go", "internal/y")
	s.expect(t, "unused role allowance: no package of check/syntax/<language>/provider/go imports internal/<language>")

	// A package with a hand entry is checked by its entry and does not count.
	s = layoutSet().drop("check/syntax/x/provider/go", "internal/x")
	s.spec.Rules["check/syntax/y/provider/go"] = []string{"contracts/check/go", "contracts/text/go", "internal/y", "internal/grammar", "internal/tree", "syntax/y/go", "trees/ast/go"}
	s.expect(t, "unused role allowance: no package of check/syntax/<language>/provider/go imports internal/<language>")

	// A shape that no package has is not said to be unused.
	s = layoutSet()
	for _, language := range []string{"x", "y"} {
		delete(s.packages, syntheticModule+"/analyze/files/"+language+"/queries/go")
	}
	s.expect(t)

	// What an extension adds and what a shape admits is not part of the list.
	s = layoutSet()
	for _, mapping := range []string{"x-to-y/mapping/text", "x-to-y/mapping/record", "y-to-x/mapping/text", "y-to-x/mapping/record"} {
		for _, language := range []string{"x", "y"} {
			s.drop("transform/syntax/"+mapping+"/go", "transform/syntax/"+language+"/mapping/pattern/go")
		}
	}
	s.expect(t)
}

// An extension widens one package's role by named imports. It must be used,
// must not repeat the role and lifts no prohibition.
func TestArchitectureRoleExtensions(t *testing.T) {
	const parser = "parse/syntax/x/provider/go"
	s := layoutSet().imports(parser, "contracts/semantic/go")
	s.expect(t, unlisted(parser, "contracts/semantic/go"))
	s.spec.Extensions[parser] = []string{"contracts/semantic/go"}
	s.expect(t)

	s = layoutSet()
	s.spec.Extensions[parser] = []string{"contracts/semantic/go"}
	s.expect(t, "unused extension: parse/syntax/x/provider/go may import contracts/semantic/go and does not")

	s = layoutSet()
	s.spec.Extensions[parser] = []string{"trees/ast/go"}
	s.expect(t, "extension of parse/syntax/x/provider/go repeats its role: trees/ast/go")

	s = layoutSet()
	s.spec.Extensions["transform/syntax/x-to-y/mapping/text/go"] = []string{"transform/syntax/x/mapping/pattern/go"}
	s.expect(t, "extension of transform/syntax/x-to-y/mapping/text/go repeats its role: transform/syntax/x/mapping/pattern/go")

	for _, lifted := range []roleBreak{
		prohibited(parser, "syntax/y/go", languageY),
		prohibited("syntax/x/go", "check/syntax/x/provider/go", noContribution),
		prohibited("internal/x", "contracts/parse/go", noStage),
		prohibited("transform/syntax/x/target/go", "transform/syntax/x/source/go", noSourceOrMapping),
		prohibited("transform/syntax/x-to-y/mapping/text/go", "transform/syntax/y/target/go", noTreeConstructor),
	} {
		s = layoutSet().imports(lifted.importer, lifted.imported)
		s.spec.Extensions[lifted.importer] = []string{lifted.imported}
		s.expect(t, lifted.want...)
	}

	// A foundation or a package with a hand entry has no role to extend.
	s = layoutSet()
	s.spec.Extensions["engine/go"] = []string{"contracts/print/go"}
	s.expect(t, "extension of engine/go extends no role: the package has its own entry")
}

// A hand entry replaces the role's allow-list and is checked for use. It lifts
// no prohibition of the role.
func TestArchitectureHandEntryOverRole(t *testing.T) {
	const analyzer = "analyze/files/x/provider/go"
	byRole := replaced(languageLayout["analyze/files/<l>/provider/go"], strings.NewReplacer("<l>", "x"))

	s := layoutSet().imports(analyzer, "contracts/resolution/go")
	s.expect(t, unlisted(analyzer, "contracts/resolution/go"))
	s.spec.Rules[analyzer] = append(slices.Clone(byRole), "contracts/resolution/go")
	s.expect(t)
	s.spec.Rules[analyzer] = append(slices.Clone(byRole), "contracts/resolution/go", "trees/ast/go")
	s.expect(t, "unused allowance: analyze/files/x/provider/go may import trees/ast/go and does not")

	s = layoutSet().imports(analyzer, "syntax/y/go")
	s.spec.Rules[analyzer] = append(slices.Clone(byRole), "syntax/y/go")
	s.expect(t, refused(analyzer, "syntax/y/go", languageY))
}

// A package with a hand entry can be bound to one language. It then imports
// nothing of another language, and no package of another language imports it,
// whatever an entry or an extension lists.
func TestArchitectureLanguageOwners(t *testing.T) {
	const part = "internal/xmodel"
	owned := func() *syntheticSet {
		s := layoutSet()
		s.listed(part, "semantics/x/go", "syntax/x/go")
		s.spec.Owners[part] = "x"
		return s
	}
	extended := func(s *syntheticSet, importer string) *syntheticSet {
		s.imports(importer, part)
		s.spec.Extensions[importer] = []string{part}
		return s
	}
	// No role allows the part: its importers name it.
	owned().imports("parse/syntax/x/provider/go", part).expect(t, unlisted("parse/syntax/x/provider/go", part))
	extended(owned(), "parse/syntax/x/provider/go").expect(t)
	extended(owned(), "parse/syntax/y/provider/go").expect(t, refused("parse/syntax/y/provider/go", part, "it belongs to another language, x"))

	s := owned()
	s.listed(part, "semantics/x/go", "syntax/y/go")
	s.expect(t, refused(part, "syntax/y/go", languageY))

	// Without the owner neither import is refused.
	s = extended(owned(), "parse/syntax/y/provider/go")
	delete(s.spec.Owners, part)
	s.expect(t)

	s = layoutSet()
	s.spec.Owners["syntax/x/go"] = "x"
	s.spec.Owners["engine/go"] = "x"
	s.expect(t,
		"owner of engine/go is not needed: its path or the foundation table places it",
		"owner of syntax/x/go is not needed: its path or the foundation table places it",
	)
}

// tools/provider-matrix/extract/<format>/go sees the extraction contract only.
func TestArchitectureExtractorRole(t *testing.T) {
	const extractor = "tools/provider-matrix/extract/rules/go"
	set := func(extra ...string) *syntheticSet {
		s := layoutSet()
		s.leaves("tools/provider-matrix/capability/go")
		s.listed("tools/provider-matrix/report/go", "tools/provider-matrix/capability/go")
		s.add("tools/provider-matrix/extract/other/go")
		s.add(extractor, append([]string{"tools/provider-matrix/capability/go"}, extra...)...)
		return s
	}
	set().expect(t)
	set("tools/provider-matrix/report/go").expect(t, unlisted(extractor, "tools/provider-matrix/report/go"))
	set("tools/provider-matrix/extract/other/go").expect(t, unlisted(extractor, "tools/provider-matrix/extract/other/go"))
	set("trees/ast/go").expect(t, unlisted(extractor, "trees/ast/go"))
	// An extractor is a contribution: no vocabulary imports it.
	set().imports("syntax/x/go", extractor).expect(t, refused("syntax/x/go", extractor, noContribution))
}

// The sets below keep hand entries beside a role. No real package has both
// today; the check still has to hold for one that does.

var vocabularyImports = []string{
	"contracts/analysis/go", "contracts/resolution/go", "contracts/semantic/go", "semantics/go",
	"syntax/x/go", "trees/ast/go", "trees/representation/go", "internal/x",
}

func vocabularySet(extra ...string) *syntheticSet {
	s := newSyntheticSet()
	s.leaves(vocabularyImports...)
	s.leaves("internal/xsem", "internal/ysem", "internal/y", "syntax/y/go")
	s.listed(syntheticParser)
	s.add("semantics/x/go", append(slices.Clone(vocabularyImports), extra...)...)
	return s
}

// A hand entry replaces the role's allow-list, never its prohibitions, and a
// provider exception does not either.
func TestArchitectureVocabularyEntryCannotLiftProhibition(t *testing.T) {
	s := vocabularySet(syntheticParser, "internal/xsem")
	s.spec.Rules["semantics/x/go"] = append(slices.Clone(vocabularyImports), syntheticParser, "internal/xsem")
	s.uses("semantics/x/go", syntheticParser, "Parse")
	s.except("semantics/x/go", syntheticParser, "Parse")
	s.expect(t,
		"semantics/x/go must not import internal/xsem: a vocabulary never imports a semantic helper",
		"semantics/x/go must not import parse/syntax/x/provider/go: a vocabulary never imports a contribution",
	)
}

var queryImports = []string{"contracts/semantic/go", "internal/xsem", "semantics/x/go", "syntax/x/go", "trees/representation/go"}

func querySet(extra ...string) *syntheticSet {
	s := newSyntheticSet()
	for _, path := range queryImports {
		if path != "semantics/x/go" {
			s.leaves(path)
		}
	}
	s.listed("semantics/x/go")
	s.leaves("resolve/files/x/provider/go", "analyze/files/x/provider/go", "trees/ast/go")
	s.add("analyze/files/x/queries/go", append(slices.Clone(queryImports), extra...)...)
	return s
}

func TestArchitectureQueryServiceRole(t *testing.T) {
	querySet().expect(t)
	querySet("resolve/files/x/provider/go").expect(t, "analyze/files/x/queries/go must not import resolve/files/x/provider/go: a query service does not import its resolver")
	querySet("analyze/files/x/provider/go").expect(t,
		"analyze/files/x/queries/go must not import example.test/m/analyze/files/x/provider/go",
		"analyze/files/x/queries/go imports stage provider analyze/files/x/provider/go (uses no symbol) without a provider exception",
	)
	querySet("trees/ast/go").expect(t, "analyze/files/x/queries/go must not import example.test/m/trees/ast/go")
}

// A query service is selected by a composition, like a stage provider. Another
// importer needs a provider exception that names the symbols it uses.
func TestArchitectureQueryServiceImportNeedsException(t *testing.T) {
	const queries, pair = "analyze/files/x/queries/go", "verification/files/pair/go"
	set := func() *syntheticSet {
		s := querySet()
		s.listed(pair, queries)
		s.uses(pair, queries, "Admission")
		return s
	}
	set().expect(t, pair+" imports query service "+queries+" (uses Admission) without a provider exception")
	s := set()
	s.except(pair, queries, "Admission")
	s.expect(t)
	s = set()
	s.except(pair, queries, "Admission")
	s.uses(pair, queries, "Admission", "NewQueries")
	s.expect(t, pair+" uses NewQueries from "+queries+", but its provider exception names only Admission")
	// A composition selects the query service and needs no exception.
	s = querySet()
	s.listed("tools/x/go", queries)
	s.uses("tools/x/go", queries, "NewQueries")
	s.expect(t)
}

var lawImports = []string{"verification/files/horos-to-go/go", "semantics/golang/go", "semantics/go"}

// The first two are what a law may import; the third is another classified model.
var lawModels = []string{"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"}

func lawSet(extra ...string) *syntheticSet {
	s := newSyntheticSet()
	s.leaves("verification/files/horos-to-go/go", "semantics/go", "verification/files/horos-to-go/mappings/go")
	s.listed("semantics/golang/go")
	s.add("verification/files/horos-to-go/mappings/list/go", lawImports...)
	s.add("verification/files/horos-to-go/mappings/text/go", append(slices.Clone(lawImports), extra...)...)
	s.spec.ExternalModels = slices.Clone(lawModels)
	for _, model := range lawModels {
		s.packages[model] = architecturePackage{Production: true}
	}
	// One law uses the external models of the role, so they are not unused.
	return s.importsExternal("verification/files/horos-to-go/mappings/list/go", lawModels[0], lawModels[1])
}

func (s *syntheticSet) importsExternal(path string, models ...string) *syntheticSet {
	p := s.packages[syntheticModule+"/"+path]
	p.Imports = append(p.Imports, models...)
	s.packages[syntheticModule+"/"+path] = p
	return s
}

func TestArchitectureVerificationLawRole(t *testing.T) {
	lawSet().expect(t)
	lawSet("verification/files/horos-to-go/mappings/list/go").expect(t, "verification/files/horos-to-go/mappings/text/go must not import verification/files/horos-to-go/mappings/list/go: a verification law never imports another law")
	lawSet("verification/files/horos-to-go/mappings/go").expect(t, "verification/files/horos-to-go/mappings/text/go must not import example.test/m/verification/files/horos-to-go/mappings/go")
}

// A law's external model imports come with its role: the package has no line
// in the external allowances. A vocabulary's do not. Like the role's module
// imports, each must be imported by at least one law.
func TestArchitectureVerificationLawExternalModels(t *testing.T) {
	const law = "verification/files/horos-to-go/mappings/text/go"
	unused := lawSet()
	unused.packages[syntheticModule+"/verification/files/horos-to-go/mappings/list/go"] = unused.packages[syntheticModule+"/"+law]
	unused.expect(t,
		"unused role allowance: no package of verification/files/horos-to-go/mappings/<name>/go imports github.com/Bitspark/horos/core/go/declarations",
		"unused role allowance: no package of verification/files/horos-to-go/mappings/<name>/go imports github.com/bitspark/ontos/core/go",
	)
	lawSet().importsExternal(law, lawModels[0], lawModels[1]).expect(t)
	lawSet().importsExternal(law, lawModels...).expect(t, "verification/files/horos-to-go/mappings/text/go must not import github.com/bitspark/ontos/data/go")

	s := vocabularySet()
	s.spec.ExternalModels = slices.Clone(lawModels)
	s.packages[lawModels[0]] = architecturePackage{Production: true}
	s.importsExternal("semantics/x/go", lawModels[0])
	s.expect(t, "semantics/x/go must not import github.com/Bitspark/horos/core/go/declarations")
	s.spec.ExternalImports["semantics/x/go"] = []string{lawModels[0]}
	s.expect(t)
}

// A forbidden import in a file that the build constraints of the run exclude
// is production code on another build, and the check reads it: behind a build
// tag, in a file of another operating system, and in a cgo file when cgo is
// off. `go list` reports none of these among a package's imports.
//
// The fixture is a module written at test time; no nested module is committed.
func TestArchitectureReadsFilesExcludedByBuildConstraints(t *testing.T) {
	const (
		checker = "check/syntax/x/provider/go"
		parser  = "parse/syntax/x/provider/go"
	)
	files := map[string]string{
		"go.mod":                      "module example.test/m\n\ngo 1.21\n",
		"trees/ast/go/ast.go":         "package ast\n\ntype Node struct{}\n",
		"contracts/parse/go/parse.go": "package parse\n\ntype Input struct{}\n",
		"syntax/x/go/x.go":            "package x\n\nimport \"example.test/m/trees/ast/go\"\n\ntype Tree = ast.Node\n",
		"syntax/y/go/y.go":            "package y\n\nconst Name = \"y\"\n",

		checker + "/check.go": "package check\n\nimport \"example.test/m/syntax/x/go\"\n\nfunc Check(x.Tree) error { return nil }\n",
		// A build tag: another language's vocabulary.
		checker + "/tagged.go": "//go:build sometag\n\npackage check\n\nimport _ \"example.test/m/syntax/y/go\"\n",

		parser + "/parse.go": "package parse\n\nimport \"example.test/m/syntax/x/go\"\n\nfunc Parse() x.Tree { return x.Tree{} }\n",
		// Another operating system: a stage provider, used without an exception.
		// Its standard-library import is of no concern.
		parser + "/admit_plan9.go": "package parse\n\nimport (\n\t\"os/signal\"\n\n\t\"example.test/m/check/syntax/x/provider/go\"\n)\n\nvar _ = signal.Notify\n\nfunc admit() error { return check.Check(Parse()) }\n",

		"internal/x/x.go": "package x\n\nconst Quote = '\"'\n",
		// cgo: a stage contract in a lexical helper.
		"internal/x/bridge.go": "package x\n\n// #include <stddef.h>\nimport \"C\"\n\nimport _ \"example.test/m/contracts/parse/go\"\n\nvar _ C.size_t\n",
		// An excluded test file stays out, as an included one does.
		"internal/x/tagged_test.go": "//go:build sometag\n\npackage x\n\nimport _ \"example.test/m/syntax/y/go\"\n",
	}
	root := t.TempDir()
	for name, content := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	spec := architectureSpec{
		Module:      syntheticModule,
		Foundations: map[string][]string{"trees/ast/go": {}, "contracts/parse/go": {}},
		// Hand entries keep the fixture small. The roles of these paths still
		// bind them. The parser's entry names its checker, which only the
		// excluded file imports.
		Rules: map[string][]string{
			"syntax/x/go": {"trees/ast/go"}, "syntax/y/go": {},
			checker: {"syntax/x/go"}, parser: {"syntax/x/go", checker}, "internal/x": {},
		},
	}
	want := []string{
		refused(checker, "syntax/y/go", languageY),
		"internal/x must not import C",
		refused("internal/x", "contracts/parse/go", noStage),
		"parse/syntax/x/provider/go imports stage provider check/syntax/x/provider/go (uses Check) without a provider exception",
	}
	for _, cgo := range []string{"0", "1"} {
		listed := listPackages(t, root, "CGO_ENABLED="+cgo, "GOWORK=off", "GOFLAGS=")
		// What `go list` reports alone does not show the seeded imports.
		for importer, hidden := range map[string]string{checker: "syntax/y/go", parser: checker, "internal/x": "contracts/parse/go"} {
			reported := slices.Contains(listed[syntheticModule+"/"+importer].Imports, syntheticModule+"/"+hidden)
			if reported != (importer == "internal/x" && cgo == "1") {
				t.Errorf("CGO_ENABLED=%s: go list reports the import of %s by %s: %v", cgo, hidden, importer, reported)
			}
		}
		if got := checkArchitecture(spec, architecturePackages(t, syntheticModule, listed)); !slices.Equal(got, want) {
			t.Errorf("CGO_ENABLED=%s:\n got: %q\nwant: %q", cgo, got, want)
		}
	}
}

// Symbol use is read per file: an import without an alias goes by the imported
// package's own name, and a local declaration of that name is not the import.
func TestArchitectureSymbolUse(t *testing.T) {
	const source = `package parse

import (
	"strings"

	"example.test/m/check/syntax/x/provider/go"
	renamed "example.test/m/parse/syntax/y/provider/go"
	. "example.test/m/trees/ast/go"
	_ "example.test/m/trees/cst/go"
)

type local struct{ Hidden int }

func Parse(text string) error {
	_ = strings.TrimSpace(text)
	var _ renamed.Provider
	if err := check.Check(nil); err != nil {
		return err
	}
	return shadowed(local{})
}

func shadowed(check local) error {
	_ = check.Hidden
	renamed := check
	_ = renamed.Hidden
	return nil
}
`
	file, err := parser.ParseFile(token.NewFileSet(), "parse.go", source, 0)
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]string{
		"example.test/m/check/syntax/x/provider/go": "check",
		"example.test/m/parse/syntax/y/provider/go": "parse",
		"example.test/m/trees/ast/go":               "ast",
		"example.test/m/trees/cst/go":               "cst",
	}
	uses := map[string][]string{}
	dots := fileSymbolUse(file, func(path string) (string, bool) { name, ok := names[path]; return name, ok }, uses)
	if !slices.Equal(dots, []string{"example.test/m/trees/ast/go"}) {
		t.Errorf("dot imports: %q", dots)
	}
	want := map[string][]string{
		"example.test/m/check/syntax/x/provider/go": {"Check"},
		"example.test/m/parse/syntax/y/provider/go": {"Provider"},
	}
	if len(uses) != len(want) {
		t.Errorf("uses: %q", uses)
	}
	for path, symbols := range want {
		if !slices.Equal(uses[path], symbols) {
			t.Errorf("uses of %s: got %q, want %q", path, uses[path], symbols)
		}
	}
}
