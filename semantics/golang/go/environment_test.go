package golang

import (
	"testing"

	"github.com/Bitspark/bittype/semantics/go"
)

// The element rule is Go's rule for an element of an import path, as
// checkElem states it in golang.org/x/mod/module, narrowed in two places: no
// leading dot and no plus sign. Each refused row has an admitted neighbour.
func TestImportPathElementFollowsGoRule(t *testing.T) {
	for _, element := range []string{
		"generated", "Generated", "a", "v2", "a-b", "a_b", "_a", "a.b", "a..b", "a~b", "-a",
		// A tilde refuses only with digits after it, at the end of the text before the first dot.
		"gen~x", "gen~", "~", "~x", "gen~1x", "gen~1~x", "x.gen~1", "a.b~2",
		// A device name refuses only as the whole text before the first dot.
		"console", "con1", "acon", "x.con", "x.nul.y", "null", "auxiliary", "prnt", "com", "com0", "com10", "lpt", "lpt0", "lpt10", "comx",
		"testdata", "vendor",
	} {
		if !importPathElement(element) {
			t.Errorf("%q is refused", element)
		}
	}
	for _, element := range []string{
		"", ".", "..", "...", ".a", "a.", "a b", " a", "a/b", "a\\b", "a:b", "a+b", "a@b", "é", "a\x00",
		"gen~1", "~1", "gen~12", "a~b~3", "gen~1.x", "gen~1.x.y", "~0.a",
		"con", "CON", "Con", "cOn", "prn", "PRN", "aux", "AUX", "nul", "NUL", "Nul",
		"com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9", "COM1", "Com9",
		"lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9", "LPT1", "Lpt9",
		"con.x", "NUL.txt", "aux.a.b", "com1.go", "lpt9.x",
	} {
		if importPathElement(element) {
			t.Errorf("%q is admitted", element)
		}
	}
}

// The directory rule has one declaration. A name that only resembles a
// refused one is covered.
func TestUnsupportedDirectory(t *testing.T) {
	for _, element := range []string{"testdata", "vendor", "_", "_gen", "__", "...", ".gen", ".git"} {
		if !UnsupportedDirectory(element) {
			t.Errorf("%q is covered", element)
		}
	}
	for _, element := range []string{"", "generated", "Testdata", "TESTDATA", "testdata2", "atestdata", "Vendor", "vendors", "a_b", "gen_", "a.b", "gen.", "internal", "main"} {
		if UnsupportedDirectory(element) {
			t.Errorf("%q is not covered", element)
		}
	}
}

// The names that the bridge writes and the name readers compare with.
func TestPublishedPackageNameForm(t *testing.T) {
	for name, tc := range map[string][2]semantics.Key{
		"relation": {PackageNameRelation(), {Owner: "go", Name: "package-name"}},
		"operator": {PackageNameOperator(), {Owner: "go", Name: "package-name"}},
		"encoding": {TextEncoding(), {Owner: "go", Name: "utf8"}},
	} {
		if tc[0] != tc[1] {
			t.Error(name, "changed:", tc[0])
		}
	}
}
