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
	// Foundations is the closed table of foundation and kernel packages with the
	// exact module imports of each. Such a package imports nothing outside the table.
	Foundations map[string][]string
	// Rules holds the hand entries: the exact module imports of a package whose
	// path gives no role, or of one whose entry replaces its role's allow-list.
	Rules map[string][]string
	// Extensions names, per package, imports that its role does not give.
	Extensions map[string][]string
	// Owners names the one language that a package with a hand entry belongs
	// to. Such a package is bound like a package whose path names the language.
	Owners     map[string]string
	Exceptions []providerException
	// ExternalModels are the non-module packages the production graph may hold;
	// ExternalImports names the module packages that may import each directly.
	ExternalModels  []string
	ExternalImports map[string][]string
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

// checkArchitecture returns one message per violated rule, in a stable order.
//
// Every production package of the module is classified by exactly one of:
//
//   - its entry in Foundations. It imports what the entry lists, and nothing
//     outside the table, whatever the entry says;
//   - its hand entry in Rules. It imports what the entry lists;
//   - the role its path gives (packageRole), with the package's Extensions.
//
// A package with none of the three fails. Every allowance in an entry or an
// extension must be used; a role's allow-list need not be. The prohibitions of
// a role bind every package of its shape, also under a hand entry or an
// extension: its own (roleRule.forbidden) and "a vocabulary never imports a
// contribution". A package that belongs to a language, by its path or by
// Owners, imports no package that belongs to another. An import of a stage
// provider by a package that is not a composition also needs a provider
// exception that names the symbols used.
//
// An entry for a package that the set does not hold is ignored: a rule, an
// extension, an owner, a provider exception by its importer, an external
// allowance. A package can then be deleted without an edit here, and "unused"
// is said only of a package that exists.
func checkArchitecture(spec architectureSpec, packages map[string]architecturePackage) []string {
	var failures []string
	fail := func(format string, arguments ...any) {
		failures = append(failures, fmt.Sprintf(format, arguments...))
	}
	prefix := spec.Module + "/"

	// roleOf is the role a path gives in this package set: none for a
	// foundation, and none for a helper whose language does not exist.
	roleOf := func(relative string) (roleRule, bool) {
		if _, ok := spec.Foundations[relative]; ok {
			return roleRule{}, false
		}
		role, ok := packageRole(relative)
		if ok && role.needs != "" {
			if owner, exists := packages[prefix+role.needs]; !exists || !owner.Production {
				return roleRule{}, false
			}
		}
		return role, ok
	}
	// languagesOf are the languages a package belongs to: those its path names,
	// or the owner of a package whose path names none.
	languagesOf := func(relative string) []string {
		if role, ok := roleOf(relative); ok {
			return role.scope
		}
		if owner, ok := spec.Owners[relative]; ok {
			return []string{owner}
		}
		return nil
	}

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

		foundation, isFoundation := spec.Foundations[relative]
		hand, listed := spec.Rules[relative]
		role, hasRole := roleOf(relative)
		extension, extended := spec.Extensions[relative]
		var allowed []string
		switch {
		case isFoundation:
			if listed {
				fail("%s has a foundation entry and a hand entry", relative)
			}
			allowed = foundation
		case listed:
			allowed = hand
		case hasRole:
			allowed = append(slices.Clone(role.allowed), extension...)
		default:
			if shape, ok := packageRole(relative); ok {
				fail("unclassified production package: %s (as a language's helper it needs %s, which does not exist)", relative, shape.needs)
			} else {
				fail("unclassified production package: %s", relative)
			}
			continue
		}
		byRole := !isFoundation && !listed
		if extended && !byRole {
			fail("extension of %s extends no role: the package has its own entry", relative)
		}
		if _, owned := spec.Owners[relative]; owned && (isFoundation || hasRole) {
			fail("owner of %s is not needed: its path or the foundation table places it", relative)
		}

		// prohibition returns why the package may never import a module package,
		// whatever its allow-list says, or "".
		prohibition := func(target string) string {
			if isFoundation {
				if _, ok := spec.Foundations[target]; !ok {
					return "a foundation or kernel package imports only foundation and kernel packages"
				}
				return ""
			}
			if hasRole {
				if role.forbidden != nil {
					if reason := role.forbidden(target); reason != "" {
						return reason
					}
				}
				if imported, ok := roleOf(target); ok && role.kind == roleVocabulary && imported.kind == roleContribution {
					return "a vocabulary never imports a contribution"
				}
			}
			if own := languagesOf(relative); len(own) > 0 {
				for _, language := range languagesOf(target) {
					if !slices.Contains(own, language) {
						return "it belongs to another language, " + language
					}
				}
			}
			return ""
		}

		for _, dependency := range sorted(p.Imports) {
			if dep, ok := packages[dependency]; ok && dep.Standard {
				continue
			}
			if slices.Contains(spec.ExternalImports[relative], dependency) || hasRole && slices.Contains(role.external, dependency) {
				continue
			}
			if !strings.HasPrefix(dependency, prefix) {
				fail("%s must not import %s", relative, dependency)
				continue
			}
			target := strings.TrimPrefix(dependency, prefix)
			if reason := prohibition(target); reason != "" {
				fail("%s must not import %s: %s", relative, target, reason)
				continue
			}
			if !slices.Contains(allowed, target) && !(byRole && role.admits != nil && role.admits(target)) {
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
		if !byRole {
			for _, allowance := range allowed {
				if !slices.Contains(p.Imports, prefix+allowance) {
					fail("unused allowance: %s may import %s and does not", relative, allowance)
				}
			}
			continue
		}
		for _, addition := range extension {
			switch {
			case slices.Contains(role.allowed, addition) || role.admits != nil && role.admits(addition):
				fail("extension of %s repeats its role: %s", relative, addition)
			case !slices.Contains(p.Imports, prefix+addition):
				fail("unused extension: %s may import %s and does not", relative, addition)
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
