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
//
// The allow-list is not checked for use: it states what the role may import,
// and one package of the role need not import all of it.
type roleRule struct {
	kind string
	// scope holds the languages the path names. A package with a scope never
	// imports a package whose scope holds another language.
	scope []string
	// needs is a package without which the path gives no role. A private helper
	// is named after its language, so a helper of no language is unclassified.
	needs   string
	allowed []string
	// admits reports further permitted imports that a shape states and a list
	// cannot, or is nil.
	admits func(dependency string) bool
	// external are the non-module packages every package of the role may import.
	external []string
	// forbidden returns why a dependency is never permitted for the role, or "".
	// It may be nil. It also binds a package of this shape that has a hand entry
	// or an extension: those replace or widen the allow-list, never a prohibition.
	forbidden func(dependency string) string
}

// isLanguageName reports whether a scope segment names one language. "shared"
// is kernel machinery and "<source>-to-<native>" is an ordered pair.
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

// mappingPair reads transform/syntax/<source>-to-<native>/mapping/<name>/go.
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
		// syntax/<language>/go: the syntax kinds and identity of one language.
		return roleRule{
			kind:    roleVocabulary,
			scope:   []string{s[1]},
			allowed: []string{"contracts/language/go", "trees/ast/go", "trees/cst/go"},
		}, true
	case len(s) == 3 && s[0] == "semantics" && isLanguageName(s[1]):
		// semantics/<language>/go: the semantic names and contract identities of
		// one language. It may use its language's lexical helper, never
		// internal/<language>sem, which runs the compiler behind the names.
		language := s[1]
		return roleRule{
			kind:  roleVocabulary,
			scope: []string{language},
			allowed: []string{
				"contracts/analysis/go", "contracts/resolution/go", "contracts/semantic/go", "semantics/go",
				"syntax/" + language + "/go", "trees/ast/go", "trees/representation/go", "internal/" + language,
			},
			forbidden: func(dependency string) string {
				if isSemanticHelper(dependency) {
					return "a vocabulary never imports a semantic helper"
				}
				return ""
			},
		}, true
	case len(s) == 5 && s[3] == "provider" && isLanguageName(s[2]):
		return stageProviderRole(s[0], s[1], s[2])
	case len(s) == 5 && s[0] == "analyze" && s[1] == "files" && s[3] == "queries" && isLanguageName(s[2]):
		// analyze/files/<language>/queries/go: a native query service. It may
		// import what the analysis provider it is split from could, without the
		// resolver.
		language := s[2]
		resolver := "resolve/files/" + language + "/provider/go"
		return roleRule{
			kind:  roleContribution,
			scope: []string{language},
			allowed: []string{
				"contracts/analysis/go", "contracts/resolution/go", "contracts/semantic/go", "internal/" + language + "sem",
				"semantics/go", "semantics/" + language + "/go", "syntax/" + language + "/go", "trees/representation/go",
			},
			forbidden: func(dependency string) string {
				if dependency == resolver {
					return "a query service does not import its resolver"
				}
				return ""
			},
		}, true
	case s[0] == "transform":
		return transformRole(relative)
	case isVerificationLaw(relative):
		// verification/files/horos-to-go/mappings/<name>/go: one verification
		// law. The pair is spelled out because it is the only one: its directory
		// says "go" where the language is "golang", so no path gives the scope.
		return roleRule{
			kind:     roleContribution,
			allowed:  []string{"verification/files/horos-to-go/go", "semantics/golang/go", "semantics/go"},
			external: []string{"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go"},
			forbidden: func(dependency string) string {
				if dependency != relative && isVerificationLaw(dependency) {
					return "a verification law never imports another law"
				}
				return ""
			},
		}, true
	case len(s) == 5 && strings.HasPrefix(relative, "tools/provider-matrix/extract/"):
		// tools/provider-matrix/extract/<format>/go: one capability extractor. It
		// sees the public extraction contract and nothing else of the tool.
		return roleRule{kind: roleContribution, allowed: []string{"tools/provider-matrix/capability/go"}}, true
	}
	return roleRule{}, false
}

// helperRole classifies internal/<language> and internal/<language>sem. Both
// are bound to their language by name: only that language's packages are
// allowed them, and the helper needs the language to exist. The generic
// helpers and the one-of-a-kind private parts have explicit entries instead.
func helperRole(name string) (roleRule, bool) {
	if language, ok := strings.CutSuffix(name, "sem"); ok && isLanguageName(language) {
		// internal/<language>sem: the bridge to one language's compiler. It may
		// reparse captured source with its own language's parser, under a
		// provider exception like every import of a stage provider.
		return roleRule{
			kind:  rolePrivatePart,
			scope: []string{language},
			needs: "syntax/" + language + "/go",
			allowed: []string{
				"contracts/parse/go", "contracts/semantic/go", "files/go", "parse/syntax/" + language + "/provider/go",
				"semantics/go", "semantics/" + language + "/go", "syntax/" + language + "/go", "trees/ast/go",
			},
		}, true
	}
	if !isLanguageName(name) {
		return roleRule{}, false
	}
	// internal/<language>: the lexical helper of one language. It knows
	// spelling, not structure: no carrier, no stage contract, no provider.
	return roleRule{
		kind:    rolePrivatePart,
		scope:   []string{name},
		needs:   "syntax/" + name + "/go",
		allowed: []string{"internal/grammar"},
		forbidden: func(dependency string) string {
			if strings.HasPrefix(dependency, "contracts/") || isStageProvider(dependency) {
				return "a lexical helper depends on no stage"
			}
			return ""
		},
	}, true
}

// stageProviderRole classifies <stage>/<dimension>/<language>/provider/go for
// the stages the layout has. Another stage gives no role: a new stage is a new
// kernel contract, and that is an edit.
//
// Each stage has its own list. A checker gets no concrete tree, a formatter no
// lexer, and no Syntax stage a semantic contract. Parser and formatter import
// their own checker and never each other.
func stageProviderRole(stage, dimension, language string) (roleRule, bool) {
	syntax := "syntax/" + language + "/go"
	semantics := "semantics/" + language + "/go"
	checker := "check/syntax/" + language + "/provider/go"
	lexical := "internal/" + language
	// A Files parser or formatter composes the collaborators it is given through
	// their contracts and publishes Source.
	files := []string{
		"contracts/analysis/go", "contracts/check/go", "contracts/print/go", "contracts/resolution/go", "contracts/semantic/go",
		"semantics/go", semantics, "source/go", syntax, "trees/ast/go", "trees/cst/go", "trees/representation/go",
	}
	var allowed []string
	switch stage + "/" + dimension {
	case "check/syntax":
		allowed = []string{"contracts/check/go", "contracts/text/go", "internal/grammar", "internal/tree", lexical, syntax, "trees/ast/go"}
	case "parse/syntax":
		allowed = []string{
			"contracts/parse/go", "contracts/text/go", checker, "internal/concrete", "internal/lex", "internal/parseutil",
			"internal/tree", lexical, syntax, "trees/ast/go", "trees/cst/go",
		}
	case "format/syntax":
		allowed = []string{"contracts/format/go", checker, "internal/concrete", "internal/tree", lexical, syntax, "trees/ast/go", "trees/cst/go"}
	case "resolve/files":
		allowed = []string{
			"contracts/resolution/go", "contracts/semantic/go", "internal/" + language + "sem",
			"semantics/go", semantics, syntax, "trees/representation/go",
		}
	case "analyze/files":
		allowed = []string{"contracts/analysis/go", "contracts/semantic/go", "semantics/go", semantics, syntax}
	case "parse/files":
		allowed = append(files, "contracts/parse/go")
	case "format/files":
		allowed = append(files, "contracts/format/go")
	default:
		return roleRule{}, false
	}
	return roleRule{kind: roleContribution, scope: []string{language}, allowed: allowed}, true
}

// transformRole classifies the generation packages of one language and the
// mapping contributions of one ordered pair. transform/syntax/shared is kernel.
func transformRole(relative string) (roleRule, bool) {
	if source, native, ok := mappingPair(relative); ok {
		// transform/syntax/<source>-to-<native>/mapping/<name>/go: one mapping. It
		// describes native forms for source meanings. It may build on the common
		// mapping parts of its source language, which are not contributions. It
		// never imports another mapping, and never the target that constructs
		// trees from the forms it describes.
		return roleRule{
			kind:  roleContribution,
			scope: []string{source, native},
			allowed: []string{
				"contracts/generation/go", "syntax/" + source + "/go", "syntax/" + native + "/go",
				"transform/syntax/" + source + "/source/go", "transform/syntax/" + native + "/model/go",
				"trees/ast/go", "trees/representation/go",
			},
			admits: func(dependency string) bool {
				common, ok := commonMappingSource(dependency)
				return ok && common == source
			},
			forbidden: func(dependency string) string {
				if dependency != relative && isMappingContribution(dependency) {
					return "a mapping contribution never imports another mapping contribution"
				}
				if isTransformComponent(dependency, "target") {
					return "a mapping contribution never imports a native tree constructor"
				}
				return ""
			},
		}, true
	}
	scope, component, ok := transformPath(relative)
	if !ok || !isLanguageName(scope) {
		return roleRule{}, false
	}
	syntax := "syntax/" + scope + "/go"
	switch strings.Join(component, "/") {
	case "model":
		// transform/syntax/<native>/model/go: the intended native forms. No tree
		// is allocated here and no source language is read.
		return roleRule{kind: roleVocabulary, scope: []string{scope}, allowed: []string{"contracts/generation/go", syntax}}, true
	case "symbol":
		// transform/syntax/<native>/symbol/go: the public type symbols that
		// generation and Files formatting of one language agree on.
		return roleRule{
			kind:  roleVocabulary,
			scope: []string{scope},
			allowed: []string{
				"contracts/generation/go", "internal/" + scope, "semantics/go", syntax, "trees/ast/go", "trees/representation/go",
			},
		}, true
	case "source":
		// transform/syntax/<language>/source/go: the reader of one source
		// language's analysed meaning.
		return roleRule{
			kind:  roleContribution,
			scope: []string{scope},
			allowed: []string{
				"contracts/analysis/go", "contracts/generation/go", "contracts/semantic/go", "contracts/text/go", "contracts/transform/go",
				"semantics/go", "semantics/" + scope + "/go", syntax, "trees/ast/go", "trees/representation/go",
			},
		}, true
	case "target":
		// transform/syntax/<native>/target/go: tree construction from the
		// planned native forms. It reads no source and calls no mapping.
		return roleRule{
			kind:  roleContribution,
			scope: []string{scope},
			allowed: []string{
				"contracts/generation/go", "internal/" + scope, syntax,
				"transform/syntax/" + scope + "/model/go", "transform/syntax/shared/target/go", "trees/ast/go",
			},
			forbidden: func(dependency string) string {
				_, common := commonMappingSource(dependency)
				if common || isMappingContribution(dependency) || isTransformComponent(dependency, "source") {
					return "a native target reads no source and calls no mapping"
				}
				return ""
			},
		}, true
	}
	if _, ok := commonMappingSource(relative); ok {
		// transform/syntax/<source>/mapping/go and .../mapping/<part>/go: what
		// the mapping contributions of one source language share. A part knows
		// the source language only (its scope refuses every native package) and
		// imports no other part. The one part that builds on two others says so
		// in roleExtensions.
		return roleRule{
			kind:    roleCommonPart,
			scope:   []string{scope},
			allowed: []string{"contracts/generation/go", syntax, "transform/syntax/" + scope + "/source/go", "trees/representation/go"},
		}, true
	}
	return roleRule{}, false
}
