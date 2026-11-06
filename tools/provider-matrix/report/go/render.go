// Command provider-matrix derives documentation from implemented Go providers.
package report

import (
	
	
	"fmt"
	"html"
	"net/url"
	
	"path/filepath"
	"slices"
	"strings"
)

func (p *program) links(locations []location, sourceURL string, detailPaths ...string) (string, error) {
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
		if loc.op != nil {
			detailPath := "provider-capabilities.md"
			if len(detailPaths)>0 { detailPath=detailPaths[0] }
			link += fmt.Sprintf(" ([details](%s#cap-%s))",detailPath,loc.op.id)
		}
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
	return strings.NewReplacer("|", "&#124;", "\n", " ", "\r", " ", "[", "&#91;", "]", "&#93;", "`", "&#96;", "*", "&#42;", "_", "&#95;", "\\", "&#92;").Replace(html.EscapeString(text))
}

func (p *program) render(m *matrix, sourceURL string, detailPaths ...string) (string, string, error) {
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
			links, err := p.links(row.cells[lang.id], sourceURL, detailPaths...)
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
			links, err := p.links(m.edges[from.id][to.id], sourceURL, detailPaths...)
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

