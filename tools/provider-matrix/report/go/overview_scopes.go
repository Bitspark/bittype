package report

import (
	"fmt"
	"net/url"
	"slices"
	"strings"
)

type overviewScope struct {
	key        string
	entry      authored
	languages  []string
	hasExcerpt bool
	text       string
}

func (p *program) overviewScopes(views []languageOverview, sourceURL string) ([]overviewScope, error) {
	var scopes []overviewScope
	indices := map[string]int{}
	excerpts := map[string][]scopeExcerpt{}
	for _, view := range views {
		for _, binding := range view.bindings {
			entry := binding.entry
			if entry.directive.name != "scope" {
				continue
			}
			key := stableKey(entry.document, fmt.Sprint(entry.directive.profile.present), entry.directive.profile.value)
			if index, ok := indices[key]; ok {
				if !slices.Contains(scopes[index].languages, view.language.id) {
					scopes[index].languages = append(scopes[index].languages, view.language.id)
				}
				continue
			}
			name, fragment, _ := strings.Cut(entry.document, "#")
			passages, ok := excerpts[name]
			if !ok {
				_, document, err := scopeSource(p.root, entry.document)
				if err != nil {
					return nil, err
				}
				passages, err = parseExcerpts(document)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", name, err)
				}
				excerpts[name] = passages
			}
			link := sourceTarget(p.root, sourceURL, name, 0)
			if fragment != "" {
				link += "#" + url.PathEscape(fragment)
			}
			var b limitedText
			fmt.Fprintf(&b, "**[Scope: %s](%s)** — profile: %s.\n\n", escape(entry.document), link, escape(profileLabel(entry.directive.profile)))
			found := map[string]bool{}
			for _, passage := range passages {
				if passage.fragment != fragment {
					continue
				}
				text, err := excerptText(passage.text, name, sourceURL)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", name, err)
				}
				label := "Supported features"
				if passage.part == "exclusions" {
					label = "Exclusions / limits"
				}
				fmt.Fprintf(&b, "**%s:** %s\n\n", label, text)
				found[passage.part] = true
			}
			if !found["features"] {
				b.WriteString("No compact feature passage supplied for this exact scope; consult its source and structural inventory.\n\n")
			}
			if !found["exclusions"] {
				b.WriteString("No compact exclusion passage supplied for this exact scope; this does not mean no exclusions.\n\n")
			}
			if b.err != nil {
				return nil, b.err
			}
			indices[key] = len(scopes)
			scopes = append(scopes, overviewScope{key: key, entry: entry, languages: []string{view.language.id}, hasExcerpt: len(found) > 0, text: b.String()})
		}
	}
	slices.SortFunc(scopes, func(a, b overviewScope) int {
		if a.entry.document != b.entry.document {
			return strings.Compare(a.entry.document, b.entry.document)
		}
		return strings.Compare(a.key, b.key)
	})
	return scopes, nil
}
