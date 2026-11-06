package report

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type profile struct {
	present bool
	value   string
}

type directive struct {
	name     string
	argument string
	values   map[string]string
	profile  profile
}

// rawDirectives reads declaration documentation, not rendered comment text.
// Fenced and indented examples are not annotations. Existing dimension comments
// retain their separate type-only attachment and block-comment rules.
func rawDirectives(group *ast.CommentGroup) ([]directive, error) {
	if group == nil {
		return nil, nil
	}
	var result []directive
	fenced := false
	for _, comment := range group.List {
		if !strings.HasPrefix(comment.Text, "//") {
			continue
		}
		line := strings.TrimPrefix(comment.Text, "//")
		if strings.HasPrefix(line, " ") {
			line = line[1:]
		}
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			fenced = !fenced
			continue
		}
		if fenced || !strings.HasPrefix(line, "bittype:") {
			continue
		}
		if len(line) > 4096 || len(result) >= 256 {
			return nil, fmt.Errorf("directive size/count limit exceeded")
		}
		words, err := directiveWords(line)
		if err != nil {
			return nil, err
		}
		d := directive{name: strings.TrimPrefix(words[0], "bittype:"), values: map[string]string{}}
		if d.name == "dimension" || d.name == "shared" {
			continue
		}
		var allowed map[string]bool
		switch d.name {
		case "grammar", "scope":
			if len(words) < 2 || strings.Contains(words[1], "=") {
				return nil, fmt.Errorf("%s requires one reference", d.name)
			}
			d.argument = words[1]
			words = append(words[:1], words[2:]...)
			allowed = map[string]bool{"profile": true}
		case "case":
			allowed = map[string]bool{"target": true, "feature": true, "case": true, "expects": true, "profile": true, "language": true}
		default:
			return nil, fmt.Errorf("unknown bittype directive %q", d.name)
		}
		for _, word := range words[1:] {
			key, value, ok := strings.Cut(word, "=")
			if !ok || !allowed[key] {
				return nil, fmt.Errorf("unknown %s argument %q", d.name, word)
			}
			if _, exists := d.values[key]; exists {
				return nil, fmt.Errorf("repeated %s argument %q", d.name, key)
			}
			if value == "" || !utf8.ValidString(value) {
				return nil, fmt.Errorf("empty or invalid %s value", key)
			}
			d.values[key] = value
		}
		if d.name == "case" {
			for _, key := range []string{"target", "feature", "case", "expects"} {
				if d.values[key] == "" {
					return nil, fmt.Errorf("case requires %s", key)
				}
			}
		}
		if value, ok := d.values["profile"]; ok {
			d.profile = profile{true, value}
		}
		result = append(result, d)
	}
	return result, nil
}

func directiveWords(line string) ([]string, error) {
	var result []string
	for i := 0; i < len(line); {
		for i < len(line) && (line[i] == ' ' || line[i] == '\t') {
			i++
		}
		if i == len(line) {
			break
		}
		var value strings.Builder
		for i < len(line) && line[i] != ' ' && line[i] != '\t' {
			if line[i] != '"' {
				value.WriteByte(line[i])
				i++
				continue
			}
			start := i
			i++
			for i < len(line) && line[i] != '"' {
				if line[i] == '\\' {
					i++
				}
				i++
			}
			if i >= len(line) {
				return nil, fmt.Errorf("unterminated quoted directive value")
			}
			i++
			unquoted, err := strconv.Unquote(line[start:i])
			if err != nil {
				return nil, fmt.Errorf("quoted directive value: %w", err)
			}
			value.WriteString(unquoted)
			if i < len(line) && line[i] != ' ' && line[i] != '\t' {
				return nil, fmt.Errorf("text after quoted directive value")
			}
		}
		word := value.String()
		if !utf8.ValidString(word) || strings.ContainsFunc(word, func(r rune) bool { return unicode.IsControl(r) && r != '\t' }) {
			return nil, fmt.Errorf("invalid control/UTF-8 in directive")
		}
		result = append(result, word)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("empty directive")
	}
	return result, nil
}
