package report

import (
	"fmt"
	"net/url"
	"path"
	"regexp"
	"strings"
)

// Excerpts select existing owner-written prose, not inferred support from rule
// keys or expectation labels. The exact scope fragment and profile stay attached
// to the operation that declared them. A file-wide scope is not a wildcard.
type scopeExcerpt struct {
	fragment, part, text string
}

var excerptBegin = regexp.MustCompile(`^<!-- bittype:overview scope=(\S+) part=(features|exclusions) -->$`)
var excerptLink = regexp.MustCompile(`\[([^\[\]\n]+)\]\(([^()\s]+)\)`)
var excerptBlock = regexp.MustCompile(`^(#{1,6}\s|[-+*>|]\s|[0-9]+[.)]\s|[-=]{3,}$)`)
var excerptHTML = regexp.MustCompile(`<[/!?A-Za-z][^>]*>`)

// Excerpts admit ordinary prose, emphasis, single-backtick code and simple
// inline links. Handle code separately so a literal link example stays literal.
func proseParts(text string, visit func(string) (string, error)) (string, error) {
	parts := strings.Split(text, "`")
	if len(parts)%2 == 0 {
		return "", fmt.Errorf("unclosed inline code in overview excerpt")
	}
	for i, part := range parts {
		if i%2 != 0 {
			if part == "" || strings.Contains(part, "\n") {
				return "", fmt.Errorf("overview code spans require single backticks on one line")
			}
			continue
		}
		changed, err := visit(part)
		if err != nil {
			return "", err
		}
		parts[i] = changed
	}
	return strings.Join(parts, "`"), nil
}

func parseExcerpts(document string) ([]scopeExcerpt, error) {
	var result []scopeExcerpt
	var current *scopeExcerpt
	var lines []string
	var fence byte
	var fenceSize int
	anchors := headingAnchors(document)
	for _, line := range strings.Split(strings.ReplaceAll(document, "\r\n", "\n"), "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) >= 3 && (trimmed[0] == '`' || trimmed[0] == '~') {
			size := len(trimmed) - len(strings.TrimLeft(trimmed, string(trimmed[0])))
			if size >= 3 && (fence == 0 || fence == trimmed[0] && size >= fenceSize && strings.TrimSpace(trimmed[size:]) == "") {
				if current != nil {
					return nil, fmt.Errorf("overview excerpts contain prose, not fenced blocks")
				}
				if fence == 0 {
					fence, fenceSize = trimmed[0], size
				} else {
					fence = 0
				}
				continue
			}
		}
		if fence != 0 {
			continue
		}
		if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			if current != nil {
				return nil, fmt.Errorf("indented block inside overview excerpt")
			}
			continue
		}
		if trimmed == "<!-- bittype:overview:end -->" {
			if current == nil {
				return nil, fmt.Errorf("overview end without begin")
			}
			current.text = strings.TrimSpace(strings.Join(lines, "\n"))
			if current.text == "" || len(current.text) > 4096 {
				return nil, fmt.Errorf("overview excerpt must contain 1..4096 bytes of prose")
			}
			if _, err := proseParts(current.text, func(text string) (string, error) {
				withoutLinks := excerptLink.ReplaceAllString(text, "")
				if strings.ContainsAny(withoutLinks, "[]") || excerptHTML.MatchString(withoutLinks) {
					return "", fmt.Errorf("unsupported overview inline markup")
				}
				return text, nil
			}); err != nil {
				return nil, err
			}
			result = append(result, *current)
			current, lines = nil, nil
			continue
		}
		if strings.HasPrefix(trimmed, "<!-- bittype:overview") {
			match := excerptBegin.FindStringSubmatch(trimmed)
			if match == nil || current != nil {
				return nil, fmt.Errorf("malformed or nested overview excerpt")
			}
			fragment := match[1]
			if fragment == "-" {
				fragment = ""
			} else if !anchors[fragment] {
				return nil, fmt.Errorf("overview scope fragment does not exist: %s", fragment)
			}
			current = &scopeExcerpt{fragment: fragment, part: match[2]}
			continue
		}
		if current != nil {
			// Keep the transcluded surface deliberately small: paragraphs with
			// inline Markdown. Reference links, HTML and blocks belong in the scope
			// document, where the overview always supplies a link.
			if excerptBlock.MatchString(trimmed) || strings.Contains(line, "![") {
				return nil, fmt.Errorf("overview excerpt requires paragraphs with inline Markdown")
			}
			lines = append(lines, line)
		}
	}
	if current != nil {
		return nil, fmt.Errorf("unclosed overview excerpt")
	}
	return result, nil
}

// Rebase ordinary inline links to the owning document. This narrow prose format
// does not support reference-style links, images, raw HTML or nested link labels.
func excerptText(text, document, sourceURL string) (string, error) {
	return proseParts(text, func(text string) (string, error) {
		var failure error
		text = excerptLink.ReplaceAllStringFunc(text, func(link string) string {
			match := excerptLink.FindStringSubmatch(link)
			u, err := url.Parse(match[2])
			if err != nil {
				failure = err
				return ""
			}
			if u.IsAbs() {
				if u.Scheme != "https" && u.Scheme != "http" {
					failure = fmt.Errorf("unsupported excerpt link scheme")
				}
				return link
			}
			if u.Host != "" || u.RawQuery != "" || strings.HasPrefix(u.Path, "/") || strings.Contains(u.Path, "\\") {
				failure = fmt.Errorf("invalid excerpt link")
				return ""
			}
			name := document
			if u.Path != "" {
				name = path.Join(path.Dir(document), u.Path)
			}
			if name == ".." || strings.HasPrefix(name, "../") {
				failure = fmt.Errorf("excerpt link escapes selected module")
				return ""
			}
			target := sourceTarget("", sourceURL, name, 0)
			if u.Fragment != "" {
				target += "#" + url.PathEscape(u.Fragment)
			}
			return "[" + match[1] + "](" + target + ")"
		})
		return text, failure
	})
}
