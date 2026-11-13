package typescriptimports

import (
	"testing"

	"github.com/Bitspark/bittype/files/go"
)

func TestRelativeImportUsesBothFilePlacements(t *testing.T) {
	for _, example := range []struct {
		from, to files.Path
		want     string
	}{
		{"generated/pages/page.ts", "generated/model/person.ts", "../model/person.js"},
		{"generated/pages/page.ts", "generated/records/version2/person.ts", "../records/version2/person.js"},
		{"generated/pages/page.ts", "generated/pages/person.ts", "./person.js"},
		{"page.ts", "shared/person.ts", "./shared/person.js"},
		{"x.js/page.ts", "x.ts", "../x.js"},
	} {
		spec, err := Relative(example.from, example.to)
		if err != nil || spec != example.want {
			t.Fatalf("%s -> %s: %q, %v", example.from, example.to, spec, err)
		}
	}
	for _, p := range []files.Path{"../person.ts", "person.d.ts", "thing.tsx", "q?/person.ts", "q#x.ts", "q%2F.ts"} {
		if _, err := Relative("page.ts", p); err == nil {
			t.Fatalf("invalid ESM path admitted: %s", p)
		}
	}
	if _, err := Relative("page.ts", "page.ts"); err == nil {
		t.Fatal("self import admitted")
	}
}

func TestExcludedPathIsATypeScriptPathOutsideTheProfile(t *testing.T) {
	for _, p := range []files.Path{"person.d.ts", "person.mts", "person.cts", "person.tsx", "person.d.mts", "per#son.ts", "per%son.ts", "q?/person.ts"} {
		if SourcePath(p) || !ExcludedPath(p) {
			t.Fatalf("%s is a TypeScript path that the profile excludes", p)
		}
	}
	for _, p := range []files.Path{"person.ts", "person.txt", "PERSON.TS", "per#son.txt", "person", "../person.mts", "", "ts"} {
		if ExcludedPath(p) {
			t.Fatalf("%s is admitted or no TypeScript path", p)
		}
	}
}
