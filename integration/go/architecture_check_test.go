package integration_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"strconv"
	"strings"
)

// architecturePackage is what the check reads about one package. The real test
// fills it from `go list` and the package's non-test files; the check's own
// tests fill it by hand.
type architecturePackage struct {
	Standard bool
	// Production is true for a package that has non-test Go files.
	Production bool
	// Imports are import paths.
	Imports []string
	// Uses lists, per imported module package, the exported symbols selected on
	// it in the package's non-test files.
	Uses map[string][]string
	// DotImports are the import paths that a non-test file imports as ".".
	DotImports []string
}

// providerException admits one import of a stage provider package by a
// production package that is not a composition. An entry with Issue 0 is
// accepted; any other entry is an open finding that the named issue removes.
type providerException struct {
	Importer string
	Provider string
	Symbols  []string
	Reason   string
	Issue    int
}

type architectureSpec struct {
	Module string
	// Rules is the per-package allow-list of module imports.
	Rules      map[string][]string
	Exceptions []providerException
	// ExternalModels are the non-module packages the production graph may hold;
	// ExternalImports names the module packages that may import each directly.
	ExternalModels  []string
	ExternalImports map[string][]string
}

// isStageProvider reports whether a module-relative path has the shape
// <stage>/<dimension>/<language>/provider/go. The shared printers are kernel
// packages, not a language's contribution.
func isStageProvider(relative string) bool {
	s := strings.Split(relative, "/")
	return len(s) == 5 && s[3] == "provider" && s[4] == "go" && s[2] != "shared"
}

// isComposition reports whether a package chooses and registers contributions.
// A composition may import any stage provider and use any of its symbols.
func isComposition(relative string) bool {
	return isFreeComposition(relative) || strings.HasPrefix(relative, "tools/")
}

// The test compositions are not classified at all.
func isFreeComposition(relative string) bool {
	return relative == "integration/go" || relative == "verification/integration/go"
}

// roleRule is the import rule of every package of one path shape. A package
// with a role needs no entry in the per-package list.
type roleRule struct {
	allowed []string
	// external are the non-module packages every package of the role may import.
	external []string
	// forbidden returns why a dependency is never permitted for the role, or "".
	// It also binds a package of this shape that has a per-package entry.
	forbidden func(dependency string) string
}

// isSemanticHelper reports whether a path has the shape internal/<language>sem.
func isSemanticHelper(relative string) bool {
	s := strings.Split(relative, "/")
	return len(s) == 2 && s[0] == "internal" && strings.HasSuffix(s[1], "sem") && s[1] != "sem"
}

func isMappingContribution(relative string) bool {
	s := strings.Split(relative, "/")
	return len(s) == 6 && s[0] == "transform" && s[1] == "syntax" && strings.HasPrefix(s[2], "horos-to-") && s[2] != "horos-to-" && s[3] == "mapping" && s[5] == "go"
}

func isVerificationLaw(relative string) bool {
	s := strings.Split(relative, "/")
	return len(s) == 6 && strings.HasPrefix(relative, "verification/files/horos-to-go/mappings/") && s[5] == "go"
}

// packageRole derives a package's rule from its layout path. A role's
// allow-list is not checked for use: it states what the role may import, and
// one package of the role need not import all of it.
func packageRole(relative string) (roleRule, bool) {
	s := strings.Split(relative, "/")
	switch {
	case len(s) == 3 && s[0] == "semantics" && s[2] == "go":
		// semantics/<language>/go: the vocabulary of one language. It may use
		// its language's lexical helper, never a provider or internal/<language>sem.
		language := s[1]
		return roleRule{
			allowed: []string{
				"contracts/analysis/go", "contracts/resolution/go", "contracts/semantic/go", "semantics/go",
				"syntax/" + language + "/go", "trees/ast/go", "trees/representation/go", "internal/" + language,
			},
			forbidden: func(dependency string) string {
				if isStageProvider(dependency) {
					return "a vocabulary never imports a stage provider"
				}
				if isSemanticHelper(dependency) {
					return "a vocabulary never imports a semantic helper"
				}
				return ""
			},
		}, true
	case len(s) == 5 && s[0] == "analyze" && s[1] == "files" && s[3] == "queries" && s[4] == "go":
		// analyze/files/<language>/queries/go: a native query service. It may
		// import what the analysis provider it is split from could, without the
		// resolver.
		language := s[2]
		resolver := "resolve/files/" + language + "/provider/go"
		return roleRule{
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
	case isMappingContribution(relative):
		// transform/syntax/horos-to-<native>/mapping/<name>/go: one mapping. It
		// may import what the pair's mapping package it is split from could.
		native := strings.TrimPrefix(s[2], "horos-to-")
		return roleRule{
			allowed: []string{
				"contracts/generation/go", "syntax/horos/go", "syntax/" + native + "/go", "transform/syntax/horos/source/go",
				"transform/syntax/horos/mapping/go", "transform/syntax/" + native + "/model/go", "trees/ast/go", "trees/representation/go",
				// The common horos mapping packages: source patterns, the
				// declaration recipe and naming. They are not mapping contributions.
				"transform/syntax/horos/mapping/pattern/go", "transform/syntax/horos/mapping/declarations/go", "transform/syntax/horos/mapping/naming/go",
			},
			forbidden: func(dependency string) string {
				if dependency != relative && isMappingContribution(dependency) {
					return "a mapping contribution never imports another mapping contribution"
				}
				return ""
			},
		}, true
	case isVerificationLaw(relative):
		// verification/files/horos-to-go/mappings/<name>/go: one verification
		// law. It may import what the mappings package it is split from could,
		// the external models included.
		return roleRule{
			allowed:  []string{"verification/files/horos-to-go/go", "semantics/golang/go", "semantics/go"},
			external: []string{"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go"},
			forbidden: func(dependency string) string {
				if dependency != relative && isVerificationLaw(dependency) {
					return "a verification law never imports another law"
				}
				return ""
			},
		}, true
	}
	return roleRule{}, false
}

// checkArchitecture returns one message per violated rule, in a stable order.
//
// An entry for a package that the set does not hold is ignored: a rule, a
// provider exception by its importer, an external allowance. A package can
// then be deleted without an edit here, and "unused" is said only of a package
// that exists.
func checkArchitecture(spec architectureSpec, packages map[string]architecturePackage) []string {
	var failures []string
	fail := func(format string, arguments ...any) {
		failures = append(failures, fmt.Sprintf(format, arguments...))
	}
	prefix := spec.Module + "/"

	exceptions := map[[2]string]providerException{}
	for _, e := range spec.Exceptions {
		key := [2]string{e.Importer, e.Provider}
		name := e.Importer + " -> " + e.Provider
		if _, ok := exceptions[key]; ok {
			fail("provider exception %s is listed twice", name)
			continue
		}
		exceptions[key] = e
		if !isStageProvider(e.Provider) {
			fail("provider exception %s does not name a stage provider package", name)
		}
		if isComposition(e.Importer) {
			fail("provider exception %s is not needed: %s is a composition", name, e.Importer)
		}
		if len(e.Symbols) == 0 {
			fail("provider exception %s names no symbol", name)
		}
		if strings.TrimSpace(e.Reason) == "" {
			fail("provider exception %s gives no reason", name)
		}
	}

	paths := make([]string, 0, len(packages))
	for path := range packages {
		paths = append(paths, path)
	}
	slices.Sort(paths)
	for _, path := range paths {
		p := packages[path]
		if p.Standard || !p.Production {
			continue
		}
		if !strings.HasPrefix(path, prefix) {
			if !slices.Contains(spec.ExternalModels, path) {
				fail("unclassified external production dependency: %s", path)
			}
			continue
		}
		relative := strings.TrimPrefix(path, prefix)
		if isFreeComposition(relative) {
			continue
		}
		for _, dependency := range sorted(p.DotImports) {
			if strings.HasPrefix(dependency, prefix) {
				fail("%s dot-imports %s, which hides the symbols it uses", relative, strings.TrimPrefix(dependency, prefix))
			}
		}
		role, hasRole := packageRole(relative)
		allowed, listed := spec.Rules[relative]
		if !listed {
			if !hasRole {
				fail("unclassified production package: %s", relative)
				continue
			}
			allowed = role.allowed
		}
		for _, dependency := range sorted(p.Imports) {
			if dep, ok := packages[dependency]; ok && dep.Standard {
				continue
			}
			if slices.Contains(spec.ExternalImports[relative], dependency) || slices.Contains(role.external, dependency) {
				continue
			}
			if !strings.HasPrefix(dependency, prefix) {
				fail("%s must not import %s", relative, dependency)
				continue
			}
			target := strings.TrimPrefix(dependency, prefix)
			if hasRole {
				if reason := role.forbidden(target); reason != "" {
					fail("%s must not import %s: %s", relative, target, reason)
					continue
				}
			}
			if !slices.Contains(allowed, target) {
				fail("%s must not import %s", relative, dependency)
			}
			if !isStageProvider(target) || isComposition(relative) {
				continue
			}
			used := sorted(p.Uses[dependency])
			e, ok := exceptions[[2]string{relative, target}]
			if !ok {
				fail("%s imports stage provider %s (uses %s) without a provider exception", relative, target, symbolList(used))
				continue
			}
			for _, symbol := range used {
				if !slices.Contains(e.Symbols, symbol) {
					fail("%s uses %s from %s, but its provider exception names only %s", relative, symbol, target, symbolList(sorted(e.Symbols)))
				}
			}
		}
		if listed {
			for _, allowance := range allowed {
				if !slices.Contains(p.Imports, prefix+allowance) {
					fail("unused allowance: %s may import %s and does not", relative, allowance)
				}
			}
		}
	}

	for i, e := range spec.Exceptions {
		if slices.IndexFunc(spec.Exceptions, func(o providerException) bool { return o.Importer == e.Importer && o.Provider == e.Provider }) != i {
			continue
		}
		name := e.Importer + " -> " + e.Provider
		p, ok := packages[prefix+e.Importer]
		if !ok || p.Standard || !p.Production {
			continue
		}
		if !slices.Contains(p.Imports, prefix+e.Provider) {
			fail("unused provider exception: %s does not import %s", e.Importer, e.Provider)
			continue
		}
		for _, symbol := range e.Symbols {
			if !slices.Contains(p.Uses[prefix+e.Provider], symbol) {
				fail("unused symbol in provider exception %s: %s", name, symbol)
			}
		}
	}
	return failures
}

func sorted(values []string) []string {
	result := slices.Clone(values)
	slices.Sort(result)
	return slices.Compact(result)
}

func symbolList(symbols []string) string {
	if len(symbols) == 0 {
		return "no symbol"
	}
	return strings.Join(symbols, ", ")
}

// fileSymbolUse reads one parsed file and adds to uses the symbols it selects
// on each imported package for which name returns a package name, and returns
// the paths it dot-imports. An import without an alias is known by the imported
// package's own name, which its path does not give.
//
// Use is syntactic. The parser's own scopes tell a local declaration that
// shadows an import name from the import; nothing is type-checked.
func fileSymbolUse(file *ast.File, name func(path string) (string, bool), uses map[string][]string) (dots []string) {
	local := map[string]string{}
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			continue
		}
		packageName, ok := name(path)
		if !ok {
			continue
		}
		if spec.Name != nil {
			packageName = spec.Name.Name
		}
		switch packageName {
		case "_":
		case ".":
			dots = append(dots, path)
		default:
			local[packageName] = path
		}
	}
	ast.Inspect(file, func(node ast.Node) bool {
		selector, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if id, ok := selector.X.(*ast.Ident); ok && id.Obj == nil {
			if path, ok := local[id.Name]; ok && !slices.Contains(uses[path], selector.Sel.Name) {
				uses[path] = append(uses[path], selector.Sel.Name)
			}
		}
		return true
	})
	return dots
}

// parseSymbolUse applies fileSymbolUse to source files on disk.
func parseSymbolUse(files []string, name func(path string) (string, bool)) (uses map[string][]string, dots []string, err error) {
	uses = map[string][]string{}
	set := token.NewFileSet()
	for _, path := range files {
		file, err := parser.ParseFile(set, path, nil, 0)
		if err != nil {
			return nil, nil, err
		}
		dots = append(dots, fileSymbolUse(file, name, uses)...)
	}
	return uses, dots, nil
}
