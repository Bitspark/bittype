package report

import (
	"fmt"
	"go/types"
	"slices"
	"strings"
)

type providerColumn struct {
	contract *types.Named
	label    string
	rows     []*row
}

type providerGrid struct {
	dimensions []string
	columns    []providerColumn
	identities map[string][]location
}

// Pivot discovered contributions without discarding their dimensions or Go
// identities. The same generic provider declaration spans dimensions; its exact
// instantiations and implementations remain separate inside each cell.
func collectProviderGrid(m *matrix) (providerGrid, error) {
	grid := providerGrid{identities: map[string][]location{}}
	for _, row := range m.rows {
		cells := map[string][]location{}
		for lang, locations := range row.cells {
			for _, loc := range locations {
				if loc.op == nil {
					return grid, fmt.Errorf("provider matrix contribution lacks operation identity")
				}
				if loc.op.kind == "identity" {
					grid.identities[lang] = append(grid.identities[lang], loc)
				} else {
					cells[lang] = append(cells[lang], loc)
				}
			}
		}
		if len(cells) == 0 {
			continue
		}
		contract, ok := row.identity.(*types.Named)
		if !ok {
			return grid, fmt.Errorf("provider matrix contract is not a named Go declaration: %s", row.contract)
		}
		origin := contract.Origin()
		column := slices.IndexFunc(grid.columns, func(c providerColumn) bool { return sameNamed(c.contract, origin) })
		if column < 0 {
			column = len(grid.columns)
			grid.columns = append(grid.columns, providerColumn{contract: origin, label: origin.Obj().Pkg().Name() + "." + origin.Obj().Name()})
		}
		copy := *row
		copy.cells = cells
		grid.columns[column].rows = append(grid.columns[column].rows, &copy)
		grid.dimensions = append(grid.dimensions, row.level)
	}
	slices.Sort(grid.dimensions)
	grid.dimensions = slices.Compact(grid.dimensions)
	// Labels never decide column identity. Qualify colliding display names rather
	// than merging unrelated declarations from packages with the same name.
	labels := map[string]int{}
	for _, column := range grid.columns {
		labels[column.label]++
	}
	for i := range grid.columns {
		column := &grid.columns[i]
		if labels[column.label] > 1 {
			column.label = column.contract.Obj().Pkg().Path() + "." + column.contract.Obj().Name()
		}
	}
	slices.SortFunc(grid.columns, func(a, b providerColumn) int { return strings.Compare(a.label, b.label) })
	return grid, nil
}

func (p *program) renderProviderGrid(grid providerGrid, lang language, sourceURL, detailsLink string) (string, string, error) {
	identity, err := p.links(grid.identities[lang.id], sourceURL, detailsLink)
	if err != nil {
		return "", "", err
	}
	var b strings.Builder
	if len(grid.columns) == 0 {
		return identity, "No provider contributions discovered.\n", nil
	}
	b.WriteString("| Dimension |")
	for _, column := range grid.columns {
		fmt.Fprintf(&b, " %s |", escape(column.label))
	}
	b.WriteString("\n| --- |")
	for range grid.columns {
		b.WriteString(" --- |")
	}
	b.WriteString("\n")
	for _, dimension := range grid.dimensions {
		fmt.Fprintf(&b, "| %s |", escape(dimension))
		for _, column := range grid.columns {
			var locations []location
			for _, row := range column.rows {
				if row.level == dimension {
					locations = append(locations, row.cells[lang.id]...)
				}
			}
			links, err := p.links(locations, sourceURL, detailsLink)
			if err != nil {
				return "", "", err
			}
			fmt.Fprintf(&b, " %s |", links)
		}
		b.WriteString("\n")
	}
	return identity, b.String(), nil
}
