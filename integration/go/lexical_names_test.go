package integration_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/language/go"
	cpplexical "github.com/Bitspark/bittype/model/syntax/cpp/lexical/go"
	"github.com/Bitspark/bittype/model/syntax/go"
	golexical "github.com/Bitspark/bittype/model/syntax/golang/lexical/go"
	horoslexical "github.com/Bitspark/bittype/model/syntax/horos/lexical/go"
	rustlexical "github.com/Bitspark/bittype/model/syntax/rust/lexical/go"
	typescriptlexical "github.com/Bitspark/bittype/model/syntax/typescript/lexical/go"
	cpp "github.com/Bitspark/bittype/syntax/cpp/go"
	golang "github.com/Bitspark/bittype/syntax/golang/go"
	horos "github.com/Bitspark/bittype/syntax/horos/go"
	rust "github.com/Bitspark/bittype/syntax/rust/go"
	typescript "github.com/Bitspark/bittype/syntax/typescript/go"
)

// lexicalNames are the predicate names that one language's lexical part
// declares, with the identifier of the language that owns them.
type lexicalNames struct {
	owner language.ID
	names []syntax.Predicate
}

// declaredLexicalNames holds, per directory below model/syntax, every name of
// that language's lexical part. A lexical part cannot say its language's
// identifier itself: it imports no syntax vocabulary. This composition can.
func declaredLexicalNames() map[string]lexicalNames {
	return map[string]lexicalNames{
		"cpp":    {cpp.Language{}.ID(), []syntax.Predicate{cpplexical.Identifier(), cpplexical.FundamentalType(), cpplexical.HeaderName()}},
		"golang": {golang.Language{}.ID(), []syntax.Predicate{golexical.Identifier(), golexical.ImportName()}},
		"horos":  {horos.Language{}.ID(), []syntax.Predicate{horoslexical.Identifier(), horoslexical.ModuleName(), horoslexical.Hexadecimal(), horoslexical.Pin()}},
		"rust": {rust.Language{}.ID(), []syntax.Predicate{
			rustlexical.Identifier(), rustlexical.IdentifierOrUnderscore(), rustlexical.LifetimeParameter(), rustlexical.Lifetime(), rustlexical.PathSegment(),
		}},
		"typescript": {typescript.Language{}.ID(), []syntax.Predicate{typescriptlexical.Identifier(), typescriptlexical.TypeName()}},
	}
}

// checkLexicalNames is the one guard for the names of any lexical part. A name
// is an identity that tables hold and checkers bind, so within one part the
// names are pairwise distinct, each is complete, and the owner of each is the
// language's identifier. count is the number of name functions that the part's
// source declares, so a name that this test does not list fails too.
func checkLexicalNames(t *testing.T, owner language.ID, count int, names []syntax.Predicate) {
	t.Helper()
	if len(names) != count {
		t.Errorf("the part declares %d names and the test lists %d", count, len(names))
	}
	seen := map[syntax.Predicate]bool{}
	for _, name := range names {
		if name.Owner != string(owner) {
			t.Errorf("the owner of %s/%s is not the language's identifier %s", name.Owner, name.Name, owner)
		}
		if name.Name == "" || strings.TrimSpace(name.Name) != name.Name {
			t.Errorf("%s/%q is no complete name", name.Owner, name.Name)
		}
		if seen[name] {
			t.Errorf("%s/%s is declared twice: two functions return one name", name.Owner, name.Name)
		}
		seen[name] = true
	}
}

// nameFunctions counts the functions of a lexical part that declare a name:
// those without input that return a syntax.Predicate.
func nameFunctions(t *testing.T, directory string) int {
	t.Helper()
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(directory, entry.Name()), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, declaration := range file.Decls {
			function, ok := declaration.(*ast.FuncDecl)
			if !ok || function.Recv != nil || function.Type.Params.NumFields() != 0 || function.Type.Results.NumFields() != 1 {
				continue
			}
			if result, ok := function.Type.Results.List[0].Type.(*ast.SelectorExpr); ok && result.Sel.Name == "Predicate" {
				count++
			}
		}
	}
	return count
}

// Every lexical part of the module is guarded: the parts are discovered from
// the directories, so a new language's part without a row fails.
func TestLexicalNamesAreDistinctAndOwnedByTheirLanguage(t *testing.T) {
	declared := declaredLexicalNames()
	directories, err := filepath.Glob(filepath.Join("../..", "model", "syntax", "*", "lexical", "go"))
	if err != nil || len(directories) == 0 {
		t.Fatal("no lexical part found", err)
	}
	for _, directory := range directories {
		scope := filepath.Base(filepath.Dir(filepath.Dir(directory)))
		part, ok := declared[scope]
		if !ok {
			t.Errorf("the lexical part of %s has no row in this test", scope)
			continue
		}
		delete(declared, scope)
		t.Run(scope, func(t *testing.T) { checkLexicalNames(t, part.owner, nameFunctions(t, directory), part.names) })
	}
	for scope := range declared {
		t.Errorf("a row names the lexical part of %s, which does not exist", scope)
	}
}
