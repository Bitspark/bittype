package integration_test

import "strings"

// The kinds of role a layout path gives. Foundation and kernel packages have
// no kind: they are the closed table foundationRules, not a path shape.
const (
	// A vocabulary holds names two parties agree on. It imports no contribution.
	roleVocabulary = "vocabulary"
	// A contribution is one thing a consumer registers or selects on its own.
	roleContribution = "contribution"
	// A private part serves the packages of one language and is never selected.
	rolePrivatePart = "private part"
	// A common part is shared by the mapping contributions of one source language.
	roleCommonPart = "common part"
)

// roleRule is the import rule of every package of one path shape. A package
// with a role needs no entry in any table.
type roleRule struct {
	// shape is the path pattern that gives the role.
	shape string
	kind  string
	// scope holds the languages the path names. A package with a scope never
	// imports a package whose scope holds another language.
	scope []string
	// needs is a package without which the path gives no role. A private helper
	// is named after its language, so a helper of no language is unclassified.
	needs string
	// patterns are the module imports the role permits, written once for every
	// language: "<language>" stands for the language of the path. allowed are
	// the same imports for this package, index by index.
	//
	// One package of the role need not import all of them. Each must be
	// imported by at least one package of the role, so a role is never wider
	// than what its packages use together.
	//
	// A pattern that starts with "internal/<" is a language's private part. The
	// role grants it only where the imported package is that language's helper
	// by its own role. A package cannot give itself a private part of another
	// owner by choosing its name.
	patterns []string
	allowed  []string
	// admits reports further permitted imports that a shape states and a list
	// cannot, or is nil. They are not checked for use.
	admits func(dependency string) bool
	// external are the non-module packages every package of the role may
	// import. They are checked for use like the patterns.
	external []string
	// forbidden returns why a dependency is never permitted for the role, or "".
	// It may be nil. It also binds a package of this shape that has a hand entry
	// or an extension: those replace or widen the allow-list, never a prohibition.
	forbidden func(dependency string) string
}

// role states the rule of one shape for the languages a path names. Each
// name and value in languages replaces "<name>" in the patterns.
func role(shape, kind string, languages []string, patterns ...string) roleRule {
	rule := roleRule{shape: shape, kind: kind, patterns: patterns}
	var pairs []string
	for i := 0; i+1 < len(languages); i += 2 {
		pairs = append(pairs, "<"+languages[i]+">", languages[i+1])
		rule.scope = append(rule.scope, languages[i+1])
	}
	replacer := strings.NewReplacer(pairs...)
	for _, pattern := range patterns {
		rule.allowed = append(rule.allowed, replacer.Replace(pattern))
	}
	return rule
}

// isLanguageName reports whether a scope segment names one language. "shared"
// is kernel machinery and "<source>-to-<target>" is an ordered pair.
func isLanguageName(name string) bool {
	return name != "" && name != "shared" && !strings.Contains(name, "-to-")
}

// isStageProvider reports whether a module-relative path has the shape
// <stage>/<dimension>/<language>/provider/go. The shared printers are kernel
// packages, not a language's contribution.
func isStageProvider(relative string) bool {
	s := strings.Split(relative, "/")
	return len(s) == 5 && s[3] == "provider" && s[4] == "go" && s[2] != "shared"
}

// isQueryService reports whether a path has the shape
// analyze/files/<language>/queries/go.
func isQueryService(relative string) bool {
	s := strings.Split(relative, "/")
	return len(s) == 5 && s[0] == "analyze" && s[1] == "files" && s[3] == "queries" && s[4] == "go" && isLanguageName(s[2])
}

// selectedKind names the kind of a package that a consumer selects on its own:
// a stage provider or a query service. Only a composition imports one freely.
// Any other importer needs a provider exception that names the symbols it uses.
func selectedKind(relative string) (string, bool) {
	switch {
	case isStageProvider(relative):
		return "stage provider", true
	case isQueryService(relative):
		return "query service", true
	}
	return "", false
}

// isSemanticHelper reports whether a path has the shape internal/<language>sem.
func isSemanticHelper(relative string) bool {
	s := strings.Split(relative, "/")
	return len(s) == 2 && s[0] == "internal" && strings.HasSuffix(s[1], "sem") && s[1] != "sem"
}

// transformPath splits transform/syntax/<scope>/<component...>/go.
func transformPath(relative string) (scope string, component []string, ok bool) {
	s := strings.Split(relative, "/")
	if len(s) < 5 || s[0] != "transform" || s[1] != "syntax" || s[len(s)-1] != "go" {
		return "", nil, false
	}
	return s[2], s[3 : len(s)-1], true
}

// mappingPair reads transform/syntax/<source>-to-<target>/mapping/<name>/go.
// The two names are the roles of the pair. Which pairs exist follows from
// what the module holds: the test reads no language's name.
func mappingPair(relative string) (source, native string, ok bool) {
	scope, component, ok := transformPath(relative)
	if !ok || len(component) != 2 || component[0] != "mapping" {
		return "", "", false
	}
	source, native, ok = strings.Cut(scope, "-to-")
	return source, native, ok && isLanguageName(source) && isLanguageName(native)
}

func isMappingContribution(relative string) bool {
	_, _, ok := mappingPair(relative)
	return ok
}

// The contribution contract gives an ordered pair two package roles beside the
// mapping (docs/contribution-contract.md, O25, row 28):
//
//	transform/syntax/<source>-to-<target>/vocabulary/<name>/go  the vocabulary of a contribution
//	transform/syntax/<source>-to-<target>/codec/<name>/go       the payload codec of its schema
//
// pairPart reads such a path: the pair, the part and the contribution's name.
func pairPart(relative string) (source, native, part, name string, ok bool) {
	scope, component, ok := transformPath(relative)
	if !ok || len(component) != 2 {
		return "", "", "", "", false
	}
	source, native, ok = strings.Cut(scope, "-to-")
	return source, native, component[0], component[1], ok && isLanguageName(source) && isLanguageName(native)
}

// isPairPart reports whether a path is the given part of some contribution.
func isPairPart(relative, part string) bool {
	_, _, p, _, ok := pairPart(relative)
	return ok && p == part
}

// ownVocabulary is the vocabulary package of the contribution that a mapping
// or a codec belongs to: same pair, same name.
func ownVocabulary(relative string) string {
	source, native, _, name, _ := pairPart(relative)
	return "transform/syntax/" + source + "-to-" + native + "/vocabulary/" + name + "/go"
}

// familyVocabulary reads transform/syntax/<source>/family/<name>/go: the
// vocabulary of a family, without a target (O27).
func familyVocabulary(relative string) (source, name string, ok bool) {
	scope, component, ok := transformPath(relative)
	if !ok || !isLanguageName(scope) || len(component) != 2 || component[0] != "family" {
		return "", "", false
	}
	return scope, component[1], true
}

// payloadShared is the shared part of the payload codec role. It has a hand
// entry: a payload schema is a horos type whatever the source language of a
// generation is, so there is one such package.
const payloadShared = "transform/syntax/horos/payload/go"

// bundleRole classifies the vocabulary and the codec of a contribution.
// Neither imports the mapping, and the codec imports no vocabulary but its
// own: "Only a composition imports a contribution" (AGENTS.md), and O25.
func bundleRole(relative string) (roleRule, bool) {
	source, native, part, name, ok := pairPart(relative)
	if !ok {
		return roleRule{}, false
	}
	languages := []string{"source", source, "target", native}
	var rule roleRule
	switch part {
	case "vocabulary":
		// The reference of one contribution and the plain Go type of its
		// payload schema. It imports the contract for the identity types, the
		// vocabulary of its family for the meaning (O27) and the profile
		// identities of its source's mapping side for the publisher. It
		// imports nothing of horos or ontos.
		rule = role("transform/syntax/<source>-to-<target>/vocabulary/<name>/go", roleVocabulary, languages,
			"contracts/generation/go", "transform/syntax/<source>/mapping/go",
		)
		family := "transform/syntax/" + source + "/family/" + name + "/go"
		rule.admits = func(dependency string) bool { return dependency == family }
	case "codec":
		// The codec of one payload schema: two functions between the Go type
		// and an ontos value. It imports the ontos value packages (O37) and
		// reaches horos only through the shared part of the codec role.
		rule = role("transform/syntax/<source>-to-<target>/codec/<name>/go", roleContribution, languages,
			"contracts/generation/go", payloadShared,
		)
		rule.external = []string{"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"}
		rule.admits = func(dependency string) bool { return dependency == ownVocabulary(relative) }
	default:
		return roleRule{}, false
	}
	rule.forbidden = func(dependency string) string {
		if dependency == relative || dependency == ownVocabulary(relative) && part != "vocabulary" {
			return ""
		}
		for _, other := range []string{"mapping", "vocabulary", "codec"} {
			if isPairPart(dependency, other) {
				return "a " + part + " package imports no mapping, no codec and no vocabulary but its own"
			}
		}
		if _, family, ok := familyVocabulary(dependency); ok && (part != "vocabulary" || family != name) {
			return "only the vocabulary of a contribution imports the vocabulary of its family"
		}
		return ""
	}
	return rule, true
}

// commonMappingSource reads transform/syntax/<source>/mapping/go and
// transform/syntax/<source>/mapping/<part>/go.
func commonMappingSource(relative string) (source string, ok bool) {
	scope, component, ok := transformPath(relative)
	return scope, ok && isLanguageName(scope) && len(component) <= 2 && component[0] == "mapping"
}

// isTransformComponent reports whether a path is transform/syntax/<language>/<component>/go.
func isTransformComponent(relative, name string) bool {
	scope, component, ok := transformPath(relative)
	return ok && isLanguageName(scope) && len(component) == 1 && component[0] == name
}

func isVerificationLaw(relative string) bool {
	s := strings.Split(relative, "/")
	return len(s) == 6 && strings.HasPrefix(relative, "verification/files/horos-to-go/mappings/") && s[5] == "go"
}

// packageRole derives a package's rule from its layout path. The language, the
// stage and the pair come from the path, so a new package of a known shape has
// its rule without an edit here.
func packageRole(relative string) (roleRule, bool) {
	s := strings.Split(relative, "/")
	switch {
	case len(s) == 2 && s[0] == "internal":
		return helperRole(s[1])
	case s[len(s)-1] != "go":
		return roleRule{}, false
	case len(s) == 3 && s[0] == "syntax" && isLanguageName(s[1]):
		// The syntax kinds and identity of one language.
		return role("syntax/<language>/go", roleVocabulary, []string{"language", s[1]},
			"contracts/language/go", "trees/ast/go", "trees/cst/go",
		), true
	case len(s) == 3 && s[0] == "semantics" && isLanguageName(s[1]):
		// The semantic names and contract identities of one language. It may use
		// its language's lexical helper, never internal/<language>sem, which
		// runs the compiler behind the names.
		rule := role("semantics/<language>/go", roleVocabulary, []string{"language", s[1]},
			"contracts/analysis/go", "contracts/resolution/go", "contracts/semantic/go", "semantics/go",
			"syntax/<language>/go", "trees/ast/go", "trees/representation/go", "internal/<language>",
		)
		rule.forbidden = func(dependency string) string {
			if isSemanticHelper(dependency) {
				return "a vocabulary never imports a semantic helper"
			}
			return ""
		}
		return rule, true
	case len(s) == 5 && s[3] == "provider" && isLanguageName(s[2]):
		return stageProviderRole(s[0], s[1], s[2])
	case isQueryService(relative):
		// A native query service over the facts of its language's compiler
		// bridge. It is not an analysis, and it does not import its resolver.
		rule := role("analyze/files/<language>/queries/go", roleContribution, []string{"language", s[2]},
			"contracts/semantic/go", "internal/<language>sem", "semantics/<language>/go", "syntax/<language>/go", "trees/representation/go",
		)
		resolver := "resolve/files/" + s[2] + "/provider/go"
		rule.forbidden = func(dependency string) string {
			if dependency == resolver {
				return "a query service does not import its resolver"
			}
			return ""
		}
		return rule, true
	case s[0] == "transform":
		return transformRole(relative)
	case isVerificationLaw(relative):
		// One verification law of the horos-to-go pair. The pair is spelled out
		// because it is the only one: its directory says "go" where the language
		// is "golang", so no path gives the scope.
		rule := role("verification/files/horos-to-go/mappings/<name>/go", roleContribution, nil,
			"verification/files/horos-to-go/go", "semantics/golang/go", "semantics/go",
		)
		rule.external = []string{"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go"}
		rule.forbidden = func(dependency string) string {
			if dependency != relative && isVerificationLaw(dependency) {
				return "a verification law never imports another law"
			}
			return ""
		}
		return rule, true
	case len(s) == 5 && strings.HasPrefix(relative, "tools/provider-matrix/extract/"):
		// One capability extractor. It sees the public extraction contract and
		// nothing else of the tool.
		return role("tools/provider-matrix/extract/<format>/go", roleContribution, nil, "tools/provider-matrix/capability/go"), true
	}
	return roleRule{}, false
}

// helperRole classifies internal/<language> and internal/<language>sem. Both
// are bound to their language by name: a role grants them only to that
// language's packages, and they have their role only when the language exists.
// The generic helpers are in the foundation table and the one-of-a-kind
// private parts have hand entries.
func helperRole(name string) (roleRule, bool) {
	if language, ok := strings.CutSuffix(name, "sem"); ok && isLanguageName(language) {
		// The bridge to one language's compiler. It may reparse captured source
		// with its own language's parser, under a provider exception like every
		// import of a stage provider.
		rule := role("internal/<language>sem", rolePrivatePart, []string{"language", language},
			"contracts/parse/go", "contracts/semantic/go", "files/go", "parse/syntax/<language>/provider/go",
			"semantics/go", "semantics/<language>/go", "syntax/<language>/go", "trees/ast/go",
		)
		rule.needs = "syntax/" + language + "/go"
		return rule, true
	}
	if !isLanguageName(name) {
		return roleRule{}, false
	}
	// The lexical helper of one language. It knows spelling, not structure: no
	// carrier, no stage contract, no provider.
	rule := role("internal/<language>", rolePrivatePart, []string{"language", name}, "internal/grammar")
	rule.needs = "syntax/" + name + "/go"
	rule.forbidden = func(dependency string) string {
		if strings.HasPrefix(dependency, "contracts/") || isStageProvider(dependency) {
			return "a lexical helper depends on no stage"
		}
		return ""
	}
	return rule, true
}

// stageProviderPatterns holds, per stage and dimension the layout has, what a
// provider of one language may import. Each stage has its own list. A checker
// gets no concrete tree, a formatter no lexer, and no Syntax stage a semantic
// contract. Parser and formatter import their own checker and never each
// other. A Files parser or formatter composes the collaborators it is given
// through their contracts and publishes Source.
var stageProviderPatterns = map[string][]string{
	"check/syntax": {
		"contracts/check/go", "contracts/text/go", "internal/grammar", "internal/tree", "internal/<language>", "syntax/<language>/go", "trees/ast/go",
	},
	"parse/syntax": {
		"contracts/parse/go", "contracts/text/go", "check/syntax/<language>/provider/go", "internal/concrete", "internal/lex",
		"internal/parseutil", "internal/tree", "internal/<language>", "syntax/<language>/go", "trees/ast/go", "trees/cst/go",
	},
	"format/syntax": {
		"contracts/format/go", "check/syntax/<language>/provider/go", "internal/concrete", "internal/tree", "internal/<language>",
		"syntax/<language>/go", "trees/ast/go", "trees/cst/go",
	},
	"resolve/files": {
		"contracts/resolution/go", "contracts/semantic/go", "internal/<language>sem", "semantics/go", "semantics/<language>/go",
		"syntax/<language>/go", "trees/representation/go",
	},
	"analyze/files": {
		"contracts/analysis/go", "contracts/semantic/go", "semantics/go", "semantics/<language>/go", "syntax/<language>/go",
	},
	"parse/files": {
		"contracts/analysis/go", "contracts/check/go", "contracts/parse/go", "contracts/print/go", "contracts/resolution/go", "contracts/semantic/go",
		"semantics/go", "source/go", "syntax/<language>/go", "trees/ast/go", "trees/cst/go", "trees/representation/go",
	},
	"format/files": {
		"contracts/analysis/go", "contracts/check/go", "contracts/format/go", "contracts/print/go", "contracts/resolution/go", "contracts/semantic/go",
		"semantics/go", "semantics/<language>/go", "source/go", "syntax/<language>/go", "trees/ast/go", "trees/cst/go", "trees/representation/go",
	},
}

// stageProviderRole classifies <stage>/<dimension>/<language>/provider/go.
// A stage or dimension the layout does not have gives no role: it is a new
// kernel contract, and that is an edit.
func stageProviderRole(stage, dimension, language string) (roleRule, bool) {
	patterns, ok := stageProviderPatterns[stage+"/"+dimension]
	if !ok {
		return roleRule{}, false
	}
	return role(stage+"/"+dimension+"/<language>/provider/go", roleContribution, []string{"language", language}, patterns...), true
}

// transformRole classifies the generation packages of one language and the
// mapping contributions of one ordered pair. transform/syntax/shared is kernel.
func transformRole(relative string) (roleRule, bool) {
	if source, native, ok := mappingPair(relative); ok {
		// One mapping. It describes native forms for source meanings. It may
		// build on the common mapping parts of its source language, which are
		// not contributions. It never imports another mapping, and never the
		// target that constructs trees from the forms it describes.
		rule := role("transform/syntax/<source>-to-<target>/mapping/<name>/go", roleContribution, []string{"source", source, "target", native},
			"contracts/generation/go", "syntax/<source>/go", "syntax/<target>/go",
			"transform/syntax/<source>/source/go", "transform/syntax/<target>/model/go", "trees/representation/go",
		)
		rule.admits = func(dependency string) bool {
			// A mapping package gains one import, of its own vocabulary (O25).
			if dependency == ownVocabulary(relative) {
				return true
			}
			common, ok := commonMappingSource(dependency)
			return ok && common == source
		}
		rule.forbidden = func(dependency string) string {
			if dependency != relative && isMappingContribution(dependency) {
				return "a mapping contribution never imports another mapping contribution"
			}
			if isTransformComponent(dependency, "target") {
				return "a mapping contribution never imports a native tree constructor"
			}
			// A mapping imports nothing of its bundle but its own vocabulary,
			// and it imports neither horos nor ontos: the codec does (O25, O37).
			if dependency != ownVocabulary(relative) && (isPairPart(dependency, "vocabulary") || isPairPart(dependency, "codec")) {
				return "a mapping contribution imports no codec and no vocabulary but its own"
			}
			if _, _, ok := familyVocabulary(dependency); ok || dependency == payloadShared {
				return "a mapping contribution reaches its family's meaning through its own vocabulary and no payload codec"
			}
			return ""
		}
		return rule, true
	}
	if rule, ok := bundleRole(relative); ok {
		return rule, true
	}
	if source, _, ok := familyVocabulary(relative); ok {
		// The vocabulary of a family, without a target: the one home of the
		// family's meaning identity and its revision (O27).
		return role("transform/syntax/<language>/family/<name>/go", roleVocabulary, []string{"language", source}, "contracts/generation/go"), true
	}
	scope, component, ok := transformPath(relative)
	if !ok || !isLanguageName(scope) {
		return roleRule{}, false
	}
	language := []string{"language", scope}
	switch strings.Join(component, "/") {
	case "model":
		// The intended native forms. No tree is allocated here and no source
		// language is read.
		return role("transform/syntax/<language>/model/go", roleVocabulary, language, "contracts/generation/go", "syntax/<language>/go"), true
	case "symbol":
		// The public type symbols that generation and Files formatting of one
		// language agree on.
		return role("transform/syntax/<language>/symbol/go", roleVocabulary, language,
			"contracts/generation/go", "internal/<language>", "semantics/go", "syntax/<language>/go", "trees/ast/go", "trees/representation/go",
		), true
	case "source":
		// The reader of one source language's analysed meaning.
		return role("transform/syntax/<language>/source/go", roleContribution, language,
			"contracts/analysis/go", "contracts/generation/go", "contracts/semantic/go", "contracts/text/go", "contracts/transform/go",
			"semantics/go", "semantics/<language>/go", "syntax/<language>/go", "trees/ast/go", "trees/representation/go",
		), true
	case "target":
		// Tree construction from the planned native forms. It reads no source
		// and calls no mapping.
		rule := role("transform/syntax/<language>/target/go", roleContribution, language,
			"contracts/generation/go", "internal/<language>", "syntax/<language>/go",
			"transform/syntax/<language>/model/go", "transform/syntax/shared/target/go", "trees/ast/go",
		)
		rule.forbidden = func(dependency string) string {
			_, common := commonMappingSource(dependency)
			if common || isMappingContribution(dependency) || isTransformComponent(dependency, "source") {
				return "a native target reads no source and calls no mapping"
			}
			return ""
		}
		return rule, true
	}
	if len(component) == 1 && component[0] == "mapping" {
		// The identities of one source language's mapping side: its profiles
		// and the publisher of their contributions. It imports nothing, so that
		// the vocabulary of a contribution, which imports it, reaches no reader.
		return role("transform/syntax/<language>/mapping/go", roleVocabulary, language), true
	}
	if _, ok := commonMappingSource(relative); ok {
		// What the mapping contributions of one source language share. A part
		// knows the source language only (its scope refuses every package of a
		// target language) and imports no other part. The one part that builds
		// on another says so in roleExtensions.
		return role("transform/syntax/<language>/mapping/<part>/go", roleCommonPart, language,
			"contracts/generation/go", "syntax/<language>/go", "transform/syntax/<language>/source/go",
			"transform/syntax/shared/contribution/go",
		), true
	}
	return roleRule{}, false
}
