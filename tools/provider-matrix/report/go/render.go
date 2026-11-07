// Command provider-matrix derives documentation from implemented Go providers.
package report

import (
	"fmt"
	"go/token"
	"go/types"
	"html"
	"net/url"

	"path/filepath"
	"slices"
	"strings"
)

// Label the exposed provider and link to its actual public declaration. Shared
// callees remain separate source relationships in the details; they never stand
// in for the public method's label or location. Labels decide no identity.
func (p *program) publicLocation(op *operation) (string, token.Pos, error) {
	if op == nil {
		return "", token.NoPos, fmt.Errorf("report contribution lacks operation identity")
	}
	qualify := func(pkg *types.Package) string { return strings.TrimPrefix(pkg.Path(), p.module+"/") }
	if op.public != nil {
		var label string
		if op.receiver != nil {
			label = types.TypeString(op.receiver, qualify)
			if _, pointer := op.receiver.(*types.Pointer); pointer {
				label = "(" + label + ")"
			}
		} else {
			label = qualify(op.public.Pkg())
		}
		return label + "." + op.public.Name(), op.public.Origin().Pos(), nil
	}
	if op.kind == "identity" && op.language.marker != nil {
		return types.TypeString(op.language.marker, qualify), op.language.marker.Obj().Pos(), nil
	}
	return "", token.NoPos, fmt.Errorf("report operation has no public declaration")
}

func (p *program) links(locations []location, sourceURL string, detailPaths ...string) (string, error) {
	unique := map[string]bool{}
	var links []string
	for _, loc := range locations {
		label, at, err := p.publicLocation(loc.op)
		if err != nil {
			return "", err
		}
		pos := p.fset.PositionFor(at, false)
		if !pos.IsValid() {
			return "", fmt.Errorf("missing source position for %s", label)
		}
		rel, err := filepath.Rel(p.root, pos.Filename)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return "", fmt.Errorf("implementation is outside the repository: %s", pos.Filename)
		}
		if loc.op.kind == "mapping provider" {
			label = "mapping " + label
		}
		parts := strings.Split(filepath.ToSlash(rel), "/")
		for i := range parts {
			parts[i] = url.PathEscape(parts[i])
		}
		target := fmt.Sprintf("%s/%s#L%d", strings.TrimSuffix(sourceURL, "/"), strings.Join(parts, "/"), pos.Line)
		link := fmt.Sprintf("[%s](%s)", escape(label), target)
		detailPath := "provider-capabilities.md"
		if len(detailPaths) > 0 {
			detailPath = detailPaths[0]
		}
		link += fmt.Sprintf(" ([details](%s#cap-%s))", detailPath, loc.op.id)
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
	var providers strings.Builder
	providers.WriteString("| Level | Provider |")
	for _, lang := range m.languages {
		fmt.Fprintf(&providers, " %s |", matrixLanguageLabel(lang, detailPaths...))
	}
	providers.WriteString("\n| --- | --- |")
	for range m.languages {
		providers.WriteString(" :---: |")
	}
	providers.WriteString("\n")
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
	transformers, err := p.renderTransformers(m, sourceURL, detailPaths...)
	if err != nil {
		return "", "", err
	}
	return providers.String(), transformers, nil
}

func matrixLanguageLabel(lang language, detailPaths ...string) string {
	if len(detailPaths) > 1 {
		return fmt.Sprintf("[%s](%s#language-%s)", escape(lang.name), detailPaths[1], stableKey(lang.id))
	}
	return escape(lang.name)
}

// Both report views use the same direct-edge matrix. Never add inferred routes
// or identity edges while rendering it.
func (p *program) renderTransformers(m *matrix, sourceURL string, detailPaths ...string) (string, error) {
	var b strings.Builder
	b.WriteString("| Source → target |")
	for _, lang := range m.languages {
		fmt.Fprintf(&b, " %s |", matrixLanguageLabel(lang, detailPaths...))
	}
	b.WriteString("\n| --- |")
	for range m.languages {
		b.WriteString(" :---: |")
	}
	b.WriteString("\n")
	for _, from := range m.languages {
		fmt.Fprintf(&b, "| %s |", matrixLanguageLabel(from, detailPaths...))
		for _, to := range m.languages {
			links, err := p.links(m.edges[from.id][to.id], sourceURL, detailPaths...)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, " %s |", links)
		}
		b.WriteString("\n")
	}
	return b.String(), nil
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
