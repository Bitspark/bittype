package report

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var markdownAnchor = regexp.MustCompile(`<a\s+(?:id|name)=["']([^"']+)["'][^>]*>`)
var markdownLink = regexp.MustCompile(`\[([^\]]+)\]\([^)]*\)`)
var markdownTag = regexp.MustCompile(`<[^>]+>`)

func contained(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel)
}

// scopeDocument admits only an existing Markdown file within the selected
// module. Lexical and resolved filesystem containment are both checked.
func scopeDocument(root, reference string) (string, error) {
	reference, _, err := scopeSource(root, reference)
	return reference, err
}

func scopeSource(root, reference string) (string, string, error) {
	name, fragment, _ := strings.Cut(reference, "#")
	if name == "" || strings.ContainsAny(name, "\\:") || strings.HasPrefix(name, "/") || !strings.HasSuffix(name, ".md") {
		return "", "", fmt.Errorf("scope must name a module-relative Markdown file: %q", reference)
	}
	path := filepath.Join(root, filepath.FromSlash(name))
	if !contained(root, path) {
		return "", "", fmt.Errorf("scope path escapes selected module")
	}
	realRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", "", err
	}
	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", "", fmt.Errorf("scope document: %w", err)
	}
	if !contained(realRoot, realPath) {
		return "", "", fmt.Errorf("scope symlink escapes selected module")
	}
	info, err := os.Stat(realPath)
	if err != nil || !info.Mode().IsRegular() || info.Size() > 2<<20 {
		return "", "", fmt.Errorf("scope document unavailable or exceeds 2 MiB: %s", name)
	}
	data, err := os.ReadFile(realPath)
	if err != nil {
		return "", "", err
	}
	if fragment != "" {
		decoded, err := url.PathUnescape(fragment)
		if err != nil || !headingAnchors(string(data))[decoded] {
			return "", "", fmt.Errorf("scope fragment does not exist: %s", reference)
		}
		fragment = decoded
	}
	rel, _ := filepath.Rel(root, path)
	result := filepath.ToSlash(rel)
	if fragment != "" {
		result += "#" + fragment
	}
	return result, string(data), nil
}

// Heading fragments follow GitHub's punctuation removal, whitespace and duplicate
// suffix rules for ordinary Markdown headings. Explicit anchors are preferred.
func headingAnchors(document string) map[string]bool {
	anchors := map[string]bool{}
	counts := map[string]int{}
	used := map[string]bool{}
	fenced := false
	lines := strings.Split(document, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			continue
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			fenced = !fenced
			continue
		}
		if fenced {
			continue
		}
		for _, match := range markdownAnchor.FindAllStringSubmatch(line, -1) {
			anchors[match[1]] = true
		}
		text := ""
		if strings.HasPrefix(line, "#") {
			j := 0
			for j < len(line) && line[j] == '#' {
				j++
			}
			if j <= 6 && j < len(line) && line[j] == ' ' {
				text = strings.TrimSpace(strings.TrimRight(line[j:], "#"))
			}
		} else if i+1 < len(lines) && line != "" {
			next := strings.TrimSpace(lines[i+1])
			if next != "" && (strings.Trim(next, "=") == "" || strings.Trim(next, "-") == "") {
				text = line
			}
		}
		if text == "" {
			continue
		}
		text = markdownLink.ReplaceAllString(text, "$1")
		text = markdownTag.ReplaceAllString(text, "")
		var slug strings.Builder
		for _, r := range strings.ToLower(text) {
			if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' || r == '_' {
				slug.WriteRune(r)
			} else if unicode.IsSpace(r) {
				slug.WriteByte('-')
			}
		}
		base := slug.String()
		key := base
		for used[key] {
			counts[base]++
			key = fmt.Sprintf("%s-%d", base, counts[base])
		}
		used[key] = true
		anchors[key] = true
	}
	return anchors
}
