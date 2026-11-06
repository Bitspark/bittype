package report

import (
	"strings"
	"testing"
)

func TestExcerptsSelectExactDeclaredScopes(t *testing.T) {
	text := "# A\n<!-- bittype:overview scope=a part=features -->\n`Form<T>` and [details](other.md#part).\n<!-- bittype:overview:end -->\n## B\n<!-- bittype:overview scope=b part=exclusions -->\nNo bodies.\n<!-- bittype:overview:end -->\n~~~text\n<!-- bittype:overview scope=missing part=features -->\n~~~\n"
	got, err := parseExcerpts(text)
	if err != nil || len(got) != 2 || got[0].fragment != "a" || got[1].part != "exclusions" {
		t.Fatal(got, err)
	}
	rendered, err := excerptText(got[0].text, "docs/scope.md", "https://example.test/source")
	if err != nil || !strings.Contains(rendered, "[details](https://example.test/source/docs/other.md#part)") {
		t.Fatal(rendered, err)
	}
	literal := "`[literal](other.md)` and [real](other.md)."
	rebased, err := excerptText(literal, "docs/scope.md", "https://example.test/source")
	if err != nil || !strings.Contains(rebased, "`[literal](other.md)`") || !strings.Contains(rebased, "[real](https://example.test/source/docs/other.md)") {
		t.Fatal("inline code was interpreted", rebased, err)
	}
	for _, link := range []string{"../../outside.md", "//elsewhere.test/page", "javascript:alert", "/absolute.md"} {
		if _, err := excerptText("[link]("+link+")", "docs/scope.md", "https://example.test/source"); err == nil {
			t.Fatalf("accepted %s", link)
		}
	}
}

func TestMalformedExcerptsFailInsteadOfPublishingPartialScope(t *testing.T) {
	begin := "<!-- bittype:overview scope=a part=features -->\n"
	end := "\n<!-- bittype:overview:end -->"
	for name, body := range map[string]string{
		"unclosed": begin + "text", "orphan": end, "nested": begin + begin + "text" + end,
		"unknown scope": strings.Replace(begin, "scope=a", "scope=absent", 1) + "text" + end,
		"unknown part":  strings.Replace(begin, "features", "inferred", 1) + "text" + end,
		"empty":         begin + end, "oversized": begin + strings.Repeat("a", 4097) + end,
		"fenced block":   begin + "```\nx\n```" + end,
		"indented block": begin + "    code" + end,
		"reference link": begin + "[name][ref]" + end,
		"html":           begin + "<img src=x>" + end,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := parseExcerpts("# A\n" + body); err == nil {
				t.Fatal("accepted malformed excerpt")
			}
		})
	}
}
