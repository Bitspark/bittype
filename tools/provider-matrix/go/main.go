// Command provider-matrix derives documentation from implemented Go providers.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func (p *program) links(locations []location, sourceURL string) (string, error) {
	unique := map[string]bool{}
	var links []string
	for _, loc := range locations {
		pos := p.fset.PositionFor(loc.pos, false)
		if !pos.IsValid() {
			return "", fmt.Errorf("missing source position for %s", loc.label)
		}
		rel, err := filepath.Rel(p.root, pos.Filename)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return "", fmt.Errorf("implementation is outside the repository: %s", pos.Filename)
		}
		parts := strings.Split(filepath.ToSlash(rel), "/")
		for i := range parts {
			parts[i] = url.PathEscape(parts[i])
		}
		target := fmt.Sprintf("%s/%s#L%d", strings.TrimSuffix(sourceURL, "/"), strings.Join(parts, "/"), pos.Line)
		link := fmt.Sprintf("[%s](%s)", escape(loc.label), target)
		if !unique[link] {
			unique[link] = true
			links = append(links, link)
		}
	}
	slices.Sort(links)
	if len(links) == 0 {
		return "/", nil
	}
	return strings.Join(links, "<br>"), nil
}

func escape(text string) string {
	return strings.NewReplacer("|", "&#124;", "\n", " ", "\r", " ", "[", "&#91;", "]", "&#93;").Replace(text)
}

func (p *program) render(m *matrix, sourceURL string) (string, string, error) {
	var providers, transformers strings.Builder
	providers.WriteString("| Level | Provider |")
	transformers.WriteString("| Source → target |")
	for _, lang := range m.languages {
		fmt.Fprintf(&providers, " %s |", escape(lang.name))
		fmt.Fprintf(&transformers, " %s |", escape(lang.name))
	}
	providers.WriteString("\n| --- | --- |")
	transformers.WriteString("\n| --- |")
	for range m.languages {
		providers.WriteString(" :---: |")
		transformers.WriteString(" :---: |")
	}
	providers.WriteString("\n")
	transformers.WriteString("\n")
	var rows []*row
	for _, row := range m.rows {
		rows = append(rows, row)
	}
	slices.SortFunc(rows, func(a, b *row) int {
		if a.level != b.level {
			return strings.Compare(a.level, b.level)
		}
		return strings.Compare(a.contract, b.contract)
	})
	for _, row := range rows {
		fmt.Fprintf(&providers, "| %s | `%s` |", escape(row.level), strings.ReplaceAll(row.contract, "|", "&#124;"))
		for _, lang := range m.languages {
			links, err := p.links(row.cells[lang.id], sourceURL)
			if err != nil {
				return "", "", err
			}
			fmt.Fprintf(&providers, " %s |", links)
		}
		providers.WriteString("\n")
	}
	for _, from := range m.languages {
		fmt.Fprintf(&transformers, "| %s |", escape(from.name))
		for _, to := range m.languages {
			links, err := p.links(m.edges[from.id][to.id], sourceURL)
			if err != nil {
				return "", "", err
			}
			fmt.Fprintf(&transformers, " %s |", links)
		}
		transformers.WriteString("\n")
	}
	return providers.String(), transformers.String(), nil
}

func replaceRegion(document, name, table string) (string, error) {
	begin, end := "<!-- generated:"+name+":begin -->", "<!-- generated:"+name+":end -->"
	if strings.Count(document, begin) != 1 || strings.Count(document, end) != 1 {
		return "", fmt.Errorf("expected one %s generated region", name)
	}
	start, stop := strings.Index(document, begin)+len(begin), strings.Index(document, end)
	if start > stop {
		return "", fmt.Errorf("reversed %s generated markers", name)
	}
	return document[:start] + "\n\n" + strings.TrimSpace(table) + "\n\n" + document[stop:], nil
}

func run() error {
	rootFlag := flag.String("root", ".", "repository root")
	outputFlag := flag.String("output", "docs/providers.md", "Markdown document relative to the repository root")
	urlFlag := flag.String("source-url", "", "source URL prefix (defaults to the module's GitHub main branch)")
	flag.Parse()
	root, err := filepath.Abs(*rootFlag)
	if err != nil {
		return err
	}
	p, err := load(root)
	if err != nil {
		return err
	}
	m, err := p.discover()
	if err != nil {
		return err
	}
	sourceURL := *urlFlag
	if sourceURL == "" {
		if !strings.HasPrefix(p.module, "github.com/") {
			return fmt.Errorf("supply -source-url for a non-GitHub module")
		}
		sourceURL = "https://" + p.module + "/blob/main"
	}
	providers, transformers, err := p.render(m, sourceURL)
	if err != nil {
		return err
	}
	path := filepath.Join(root, *outputFlag)
	original, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	document := strings.ReplaceAll(string(original), "\r\n", "\n")
	document, err = replaceRegion(document, "providers", providers)
	if err != nil {
		return err
	}
	document, err = replaceRegion(document, "transformers", transformers)
	if err != nil {
		return err
	}
	// Preserve the document's existing newline convention and all authored prose.
	if bytes.Contains(original, []byte("\r\n")) {
		document = strings.ReplaceAll(document, "\n", "\r\n")
	}
	if document == string(original) {
		fmt.Println("Provider matrices are unchanged.")
		return nil
	}
	if err := os.WriteFile(path, []byte(document), 0644); err != nil {
		return err
	}
	fmt.Println("Updated", *outputFlag)
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
