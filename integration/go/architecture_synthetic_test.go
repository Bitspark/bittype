package integration_test

import (
	"go/parser"
	"go/token"
	"slices"
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
		spec:     architectureSpec{Module: syntheticModule, Rules: map[string][]string{}, ExternalImports: map[string][]string{}},
		packages: map[string]architecturePackage{},
	}
}

// add declares a production package and what it imports. It gets no rule.
func (s *syntheticSet) add(path string, imports ...string) {
	p := architecturePackage{Production: true, Uses: map[string][]string{}}
	for _, dependency := range imports {
		p.Imports = append(p.Imports, syntheticModule+"/"+dependency)
	}
	s.packages[syntheticModule+"/"+path] = p
}

// listed declares a production package whose per-package rule allows exactly
// what it imports.
func (s *syntheticSet) listed(path string, imports ...string) {
	s.add(path, imports...)
	s.spec.Rules[path] = slices.Clone(imports)
}

// leaves declares packages without module imports, each with an empty rule.
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
	s.spec.Rules["transform/syntax/horos-to-x/mapping/go"] = []string{"trees/ast/go", syntheticChecker}
	s.spec.ExternalImports["transform/syntax/horos-to-x/mapping/go"] = []string{"example.org/model"}
	s.except("transform/syntax/horos-to-x/mapping/go", syntheticChecker, "Check")
	// A directory that keeps only test files is not a production package either.
	s.packages[syntheticModule+"/format/syntax/x/provider/go"] = architecturePackage{Imports: []string{syntheticModule + "/" + syntheticChecker}}
	s.spec.Rules["format/syntax/x/provider/go"] = []string{"trees/ast/go", syntheticChecker}
	s.except("format/syntax/x/provider/go", syntheticChecker, "Check")
	s.expect(t)

	// The same entries are checked once the package exists.
	s.add("transform/syntax/horos-to-x/mapping/go", "trees/ast/go")
	s.expect(t,
		"unused allowance: transform/syntax/horos-to-x/mapping/go may import check/syntax/x/provider/go and does not",
		"unused provider exception: transform/syntax/horos-to-x/mapping/go does not import check/syntax/x/provider/go",
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
	// One segment more or less than a role's shape gives no role.
	s.add("semantics/x/extra/go")
	s.add("transform/syntax/horos-to-x/mapping/text/extra/go")
	s.expect(t,
		"unclassified production package: misc/go",
		"unclassified production package: semantics/x/extra/go",
		"unclassified production package: transform/syntax/horos-to-x/mapping/text/extra/go",
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
		"provider exception parse/syntax/x/provider/go -> trees/ast/go does not name a stage provider package",
		"unused symbol in provider exception parse/syntax/x/provider/go -> trees/ast/go: Node",
	)
}

// The role rules below classify package shapes that the tree does not hold
// yet. None of these packages has an entry in the per-package rules.

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

func TestArchitectureVocabularyRole(t *testing.T) {
	vocabularySet().expect(t)
	vocabularySet(syntheticParser).expect(t, "semantics/x/go must not import parse/syntax/x/provider/go: a vocabulary never imports a stage provider")
	vocabularySet("internal/xsem").expect(t, "semantics/x/go must not import internal/xsem: a vocabulary never imports a semantic helper")
	vocabularySet("internal/ysem").expect(t, "semantics/x/go must not import internal/ysem: a vocabulary never imports a semantic helper")
	vocabularySet("internal/y").expect(t, "semantics/x/go must not import example.test/m/internal/y")
	vocabularySet("syntax/y/go").expect(t, "semantics/x/go must not import example.test/m/syntax/y/go")
}

// A per-package entry replaces the role's allow-list, never its prohibitions.
func TestArchitectureVocabularyEntryCannotLiftProhibition(t *testing.T) {
	s := vocabularySet(syntheticParser, "internal/xsem")
	s.spec.Rules["semantics/x/go"] = append(slices.Clone(vocabularyImports), syntheticParser, "internal/xsem")
	s.uses("semantics/x/go", syntheticParser, "Parse")
	s.except("semantics/x/go", syntheticParser, "Parse")
	s.expect(t,
		"semantics/x/go must not import internal/xsem: a vocabulary never imports a semantic helper",
		"semantics/x/go must not import parse/syntax/x/provider/go: a vocabulary never imports a stage provider",
	)
}

var queryImports = []string{
	"contracts/analysis/go", "contracts/resolution/go", "contracts/semantic/go", "internal/xsem",
	"semantics/go", "semantics/x/go", "syntax/x/go", "trees/representation/go",
}

func querySet(extra ...string) *syntheticSet {
	s := newSyntheticSet()
	for _, path := range queryImports {
		if path != "semantics/x/go" {
			s.leaves(path)
		}
	}
	s.add("semantics/x/go")
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

var mappingImports = []string{
	"contracts/generation/go", "syntax/horos/go", "syntax/x/go", "transform/syntax/horos/source/go",
	"transform/syntax/horos/mapping/go", "transform/syntax/x/model/go", "trees/ast/go", "trees/representation/go",
}

func mappingSet(extra ...string) *syntheticSet {
	s := newSyntheticSet()
	s.leaves(mappingImports...)
	s.leaves("transform/syntax/y/model/go", "transform/syntax/horos-to-x/mapping/go")
	s.add("transform/syntax/horos-to-x/mapping/record/go", mappingImports...)
	s.add("transform/syntax/horos-to-y/mapping/text/go", "contracts/generation/go")
	s.add("transform/syntax/horos-to-x/mapping/text/go", append(slices.Clone(mappingImports), extra...)...)
	return s
}

func TestArchitectureMappingContributionRole(t *testing.T) {
	mappingSet().expect(t)
	mappingSet("transform/syntax/horos-to-x/mapping/record/go").expect(t, "transform/syntax/horos-to-x/mapping/text/go must not import transform/syntax/horos-to-x/mapping/record/go: a mapping contribution never imports another mapping contribution")
	mappingSet("transform/syntax/horos-to-y/mapping/text/go").expect(t, "transform/syntax/horos-to-x/mapping/text/go must not import transform/syntax/horos-to-y/mapping/text/go: a mapping contribution never imports another mapping contribution")
	mappingSet("transform/syntax/y/model/go").expect(t, "transform/syntax/horos-to-x/mapping/text/go must not import example.test/m/transform/syntax/y/model/go")
	mappingSet("transform/syntax/horos-to-x/mapping/go").expect(t, "transform/syntax/horos-to-x/mapping/text/go must not import example.test/m/transform/syntax/horos-to-x/mapping/go")
}

const (
	commonAdapter      = "transform/syntax/horos/mapping/go"
	commonPattern      = "transform/syntax/horos/mapping/pattern/go"
	commonDeclarations = "transform/syntax/horos/mapping/declarations/go"
	commonNaming       = "transform/syntax/horos/mapping/naming/go"
)

var commonMappingImports = []string{"contracts/generation/go", "syntax/horos/go", "transform/syntax/horos/source/go", "trees/representation/go"}

// commonMappingSet is the divided common horos mapping package under the
// per-package entries of the real rules, and one mapping contribution that
// imports all four parts. patternExtra and namingExtra are imports that the
// entries of pattern and naming do not allow.
func commonMappingSet(patternExtra, namingExtra []string) *syntheticSet {
	s := newSyntheticSet()
	s.leaves(commonMappingImports...)
	s.leaves("syntax/x/go", "transform/syntax/x/model/go", "trees/ast/go")
	s.listed(commonAdapter, commonMappingImports...)
	s.add(commonPattern, append(slices.Clone(commonMappingImports), patternExtra...)...)
	s.spec.Rules[commonPattern] = slices.Clone(commonMappingImports)
	s.add(commonNaming, namingExtra...)
	s.spec.Rules[commonNaming] = []string{}
	s.listed(commonDeclarations, append(slices.Clone(commonMappingImports), commonAdapter, commonNaming)...)
	s.add("transform/syntax/horos-to-x/mapping/declarations/go", "syntax/x/go", "transform/syntax/x/model/go", commonAdapter, commonDeclarations)
	s.add("transform/syntax/horos-to-x/mapping/record/go", append(slices.Clone(mappingImports), commonPattern, commonNaming)...)
	return s
}

// A mapping contribution may import each part of the common horos mapping
// package; they are not mapping contributions. Among the parts, the entries of
// pattern and naming allow neither each other nor the adapter. An entry states
// that as a missing allowance, so the failure carries no reason of its own.
func TestArchitectureCommonMappingPackages(t *testing.T) {
	commonMappingSet(nil, nil).expect(t)
	commonMappingSet([]string{commonNaming}, nil).expect(t, "transform/syntax/horos/mapping/pattern/go must not import example.test/m/transform/syntax/horos/mapping/naming/go")
	commonMappingSet([]string{commonAdapter}, nil).expect(t, "transform/syntax/horos/mapping/pattern/go must not import example.test/m/transform/syntax/horos/mapping/go")
	commonMappingSet([]string{commonDeclarations}, nil).expect(t, "transform/syntax/horos/mapping/pattern/go must not import example.test/m/transform/syntax/horos/mapping/declarations/go")
	commonMappingSet(nil, []string{commonPattern}).expect(t, "transform/syntax/horos/mapping/naming/go must not import example.test/m/transform/syntax/horos/mapping/pattern/go")
	commonMappingSet(nil, []string{commonAdapter}).expect(t, "transform/syntax/horos/mapping/naming/go must not import example.test/m/transform/syntax/horos/mapping/go")
}

var lawImports = []string{"verification/files/horos-to-go/go", "semantics/golang/go", "semantics/go"}

// The first two are what a law may import; the third is another classified model.
var lawModels = []string{"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"}

func lawSet(extra ...string) *syntheticSet {
	s := newSyntheticSet()
	s.leaves("verification/files/horos-to-go/go", "semantics/go", "verification/files/horos-to-go/mappings/go")
	s.add("semantics/golang/go")
	s.add("verification/files/horos-to-go/mappings/list/go", lawImports...)
	s.add("verification/files/horos-to-go/mappings/text/go", append(slices.Clone(lawImports), extra...)...)
	s.spec.ExternalModels = slices.Clone(lawModels)
	for _, model := range lawModels {
		s.packages[model] = architecturePackage{Production: true}
	}
	return s
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
// in the external allowances. A vocabulary's do not.
func TestArchitectureVerificationLawExternalModels(t *testing.T) {
	const law = "verification/files/horos-to-go/mappings/text/go"
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
