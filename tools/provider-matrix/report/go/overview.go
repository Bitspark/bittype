package report

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	capability "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
	observation "github.com/Bitspark/bittype/tools/provider-matrix/observation/go"
)

type overviewBinding struct {
	target capability.Symbol
	entry  authored
	op     *operation
}

type languageOverview struct {
	language language
	bindings []overviewBinding
	cases    []caseAssociation
	tests    map[string]bool
}

func caseKey(c caseAssociation) string {
	return stableKey(c.target.String(), c.marker.String(), c.test.id, c.directive.values["feature"], c.directive.values["case"], c.directive.values["expects"], fmt.Sprint(c.directive.profile.present), c.directive.profile.value)
}

func collectOverview(m *matrix, a *annotations) []languageOverview {
	var result []languageOverview
	for _, lang := range m.languages {
		view := languageOverview{language: lang, tests: map[string]bool{}}
		bindings, cases := map[string]bool{}, map[string]bool{}
		for _, op := range operations(m) {
			belongs := sameLanguage(op.language, lang)
			if op.from != nil {
				belongs = sameLanguage(*op.from, lang) || sameLanguage(*op.to, lang)
			}
			if !belongs {
				continue
			}
			for _, target := range operationTargets(op) {
				for _, entry := range a.operations[target] {
					key := bindingAnchor(op, target, entry)
					if !bindings[key] {
						view.bindings = append(view.bindings, overviewBinding{target: target, entry: entry, op: op})
						bindings[key] = true
					}
				}
				for _, c := range a.cases {
					if c.target != target || !markerMatches(op, c.marker) || cases[caseKey(c)] {
						continue
					}
					view.cases = append(view.cases, c)
					cases[caseKey(c)], view.tests[c.test.id] = true, true
				}
			}
		}
		slices.SortFunc(view.bindings, func(a, b overviewBinding) int {
			return strings.Compare(bindingAnchor(a.op, a.target, a.entry), bindingAnchor(b.op, b.target, b.entry))
		})
		slices.SortFunc(view.cases, func(a, b caseAssociation) int { return strings.Compare(caseKey(a), caseKey(b)) })
		result = append(result, view)
	}
	return result
}

func relativeReportLink(from, to string) string {
	link, _ := filepath.Rel(filepath.Dir(from), to)
	return filepath.ToSlash(link)
}

func (p *program) renderOverview(ctx context.Context, m *matrix, a *annotations, config Config, sourceURL string, definitions map[capability.Symbol]definitionReport) (string, error) {
	var b limitedText
	views := collectOverview(m, a)
	matrixLink := relativeReportLink(config.OverviewPath, config.MatrixPath)
	detailsLink := relativeReportLink(config.OverviewPath, config.DetailsPath)
	if config.Context != nil {
		matrixLink = sourceTarget(p.root, sourceURL, config.MatrixPath, 0)
		detailsLink = sourceTarget(p.root, sourceURL, config.DetailsPath, 0)
	}
	b.WriteString("# Language support at a glance\n\nGenerated from the same discovered providers, definitions, scope passages and test associations as the matrices. Expand a language's evidence for exact features and tests. This is bounded support, not a completeness score.\n\n")
	fmt.Fprintf(&b, "[Provider and transformer matrices](%s) · [Full capability details](%s) · [Test runs and observation artifacts](https://%s/actions/workflows/go.yml)\n\n", matrixLink, detailsLink, p.module)
	b.WriteString("Available stages are discovered implementations. **Supported features and exclusions are owner-authored scope**, not conclusions inferred from rule names. Syntax support does not imply semantic analysis or generation support. **Test counts are associated declarations**, not passing results or coverage percentages. Delegate evidence does not qualify wrappers.\n\n")
	if len(config.Executions) == 0 {
		b.WriteString("This committed view contains no execution results. Existing CI artifacts pair invocation reports with compact overviews; those retain source/build applicability and process failures. An absent observation does not mean a test never ran.\n\n")
	}
	fmt.Fprintf(&b, "Inspection build: %s, %s/%s; explicit build flags: %s.\n\n", escape(p.build.GOVERSION), escape(p.build.GOOS), escape(p.build.GOARCH), escape(strings.Join(config.BuildFlags, " ")))
	b.WriteString("Languages: ")
	for i, view := range views {
		if i > 0 {
			b.WriteString(" · ")
		}
		fmt.Fprintf(&b, "[%s](#language-%s)", escape(view.language.name), stableKey(view.language.id))
	}
	b.WriteString("\n\n## Direct transformer matrix\n\nRows are source languages; columns are target languages. Each populated cell links to its exact factory and capability evidence. `/` means no directly implemented transformer is discovered, including on the diagonal. Composed routes and configured applicability are separate.\n\n")
	transformers, err := p.renderTransformers(m, sourceURL, detailsLink, "")
	if err != nil {
		return "", err
	}
	b.WriteString(transformers)
	b.WriteString("\nProvider matrices below use discovered dimensions and provider contracts. `/` means no implementation is discovered for that pairing; it does not declare the pairing impossible. Each implementation retains its own source and evidence links.\n\n")
	grid, err := collectProviderGrid(m)
	if err != nil {
		return "", err
	}
	scopes, err := p.overviewScopes(views, sourceURL)
	if err != nil {
		return "", err
	}
	for _, view := range views {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		fmt.Fprintf(&b, "<a id=\"language-%s\"></a>\n\n## %s\n\n", stableKey(view.language.id), escape(view.language.name))
		identity, providers, err := p.renderProviderGrid(grid, view.language, sourceURL, detailsLink)
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&b, "Language ID: %s. Identity: %s.\n\n", escape(view.language.id), identity)
		b.WriteString(providers)
		b.WriteString("\n")
		var scoped []overviewScope
		for _, scope := range scopes {
			if slices.Contains(scope.languages, view.language.id) {
				scoped = append(scoped, scope)
			}
		}
		for _, scope := range scoped {
			if scope.hasExcerpt && len(scope.languages) == 1 {
				b.WriteString(scope.text)
			}
		}
		if len(scoped) == 0 {
			b.WriteString("No authored scope is attached to this language's discovered operations. Supported breadth and exclusions are unreported.\n\n")
		}
		fmt.Fprintf(&b, "<details>\n<summary>Structural inventory and exact scope associations</summary>\n\n")
		for _, scope := range scoped {
			if len(scope.languages) > 1 {
				fmt.Fprintf(&b, "- [Shared scope: %s, profile %s](#scope-%s).\n", escape(scope.entry.document), escape(profileLabel(scope.entry.directive.profile)), scope.key)
			} else if !scope.hasExcerpt {
				b.WriteString(scope.text)
			}
		}
		b.WriteString("\n")
		seenDefinition := map[string]bool{}
		for _, binding := range view.bindings {
			entry := binding.entry
			if entry.directive.name != "grammar" {
				continue
			}
			definition := definitions[entry.root]
			key := stableKey(entry.root.String(), definition.format, fmt.Sprint(entry.directive.profile.present), entry.directive.profile.value)
			if seenDefinition[key] {
				continue
			}
			seenDefinition[key] = true
			var keys []string
			for _, feature := range definition.inventory.Features {
				keys = append(keys, feature.Key)
			}
			slices.Sort(keys)
			fmt.Fprintf(&b, "- %s; format: %s; extraction: **%s**; profile: %s; **%d declared forms**.\n\n  %s\n\n", p.symbolLink(entry.root, sourceURL), escape(definition.format), escape(definition.state), escape(profileLabel(entry.directive.profile)), len(keys), escape(strings.Join(keys, ", ")))
		}
		if len(seenDefinition) == 0 {
			b.WriteString("No structural definition attached. No feature inventory is inferred from operation bodies.\n\n")
		}
		b.WriteString("These are structural forms, not a claim that every form works at every stage. Partial/unknown fields and procedural limits remain in the full details. Exact bindings (including related declarations, not wrapper qualification):\n\n")
		for _, binding := range view.bindings {
			fmt.Fprintf(&b, "- %s: [%s, profile %s](%s#binding-%s).\n", p.symbolLink(binding.target, sourceURL), escape(binding.entry.directive.name), escape(profileLabel(binding.entry.directive.profile)), detailsLink, bindingAnchor(binding.op, binding.target, binding.entry))
		}
		b.WriteString("\n</details>\n\n")
		fmt.Fprintf(&b, "<details>\n<summary>Test evidence: %d associated test declarations, %d authored cases</summary>\n\n", len(view.tests), len(view.cases))
		b.WriteString("Counts deduplicate declarations and case associations across wrappers. One test may declare several cases; a test observation never proves separate execution of each annotation. Expectations and profiles below are opaque owner labels.\n\n")
		b.WriteString("| Feature | Case | Expectation | Profile | Declared target | Test source |\n| --- | --- | --- | --- | --- | --- |\n")
		for _, c := range view.cases {
			pos := p.fset.PositionFor(c.test.pos, false)
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | [%s](%s) (%s) |\n", escape(c.directive.values["feature"]), escape(c.directive.values["case"]), escape(c.directive.values["expects"]), escape(profileLabel(c.directive.profile)), p.symbolLink(c.target, sourceURL), escape(c.test.name), sourceTarget(p.root, sourceURL, c.test.file, pos.Line), escape(c.test.variant))
		}
		if len(view.cases) == 0 {
			b.WriteString("| No associated cases in the selected sources | | | | | |\n")
		}
		b.WriteString("\n</details>\n\n")
	}
	var shared []overviewScope
	for _, scope := range scopes {
		if len(scope.languages) > 1 {
			shared = append(shared, scope)
		}
	}
	if len(shared) > 0 {
		b.WriteString("<details>\n<summary>Shared pipeline and transformer scope</summary>\n\nThese scopes are explicitly associated with more than one language; their exact per-operation bindings remain above.\n\n")
		for _, scope := range shared {
			fmt.Fprintf(&b, "<a id=\"scope-%s\"></a>\n\n%s", scope.key, scope.text)
		}
		b.WriteString("</details>\n\n")
	}
	if len(config.Executions) > 0 {
		renderOverviewObservations(&b, a, config, views)
	}
	if b.err != nil {
		return "", b.err
	}
	return strings.TrimRight(b.String(), "\n") + "\n", nil
}

func renderOverviewObservations(b *limitedText, a *annotations, config Config, views []languageOverview) {
	snapshot := observationSnapshot(a, config)
	b.WriteString("## Supplied invocation evidence\n\nCounts below are actual matched test occurrences, not case annotations. Failures, repeats and process exit remain visible. One test can be associated with multiple languages; rows must not be added together. Unmatched observations receive no language credit; the accompanying full invocation report retains all diagnostics.\n\n")
	for i, execution := range config.Executions {
		assessment := observation.Associate(snapshot, execution)
		exit := "not observed"
		if execution.Invocation.ExitCode != nil {
			exit = fmt.Sprint(*execution.Invocation.ExitCode)
		}
		fmt.Fprintf(b, "### Invocation %d: %s\n\nApplicability: **%s**; process exit: **%s**; capture complete: %t; protocol complete: %t.\n\n", i+1, escape(execution.Invocation.CommandID), assessment.Status, exit, execution.Invocation.CaptureComplete, execution.ProtocolComplete)
		for _, reason := range assessment.Reasons {
			fmt.Fprintf(b, "- %s\n", escape(reason))
		}
		columns := []observation.Outcome{observation.Passed, observation.Failed, observation.Skipped, observation.Incomplete}
		var additional []observation.Outcome
		for _, association := range assessment.Observations {
			if !slices.Contains(columns, association.Run.Outcome) {
				additional = append(additional, association.Run.Outcome)
			}
		}
		slices.Sort(additional)
		columns = append(columns, slices.Compact(additional)...)
		b.WriteString("\n`/` means no matching observation. Zero means observations matched, but none had that outcome.\n\n| Language |")
		for _, outcome := range columns {
			fmt.Fprintf(b, " %s |", escape(string(outcome)))
		}
		b.WriteString("\n| --- |")
		for range columns {
			b.WriteString(" --- |")
		}
		b.WriteString("\n")
		for _, view := range views {
			outcomes := map[observation.Outcome]int{}
			for _, association := range assessment.Observations {
				if association.Status == observation.Matched && view.tests[association.DeclarationID] {
					outcomes[association.Run.Outcome]++
				}
			}
			fmt.Fprintf(b, "| %s |", escape(view.language.name))
			for _, outcome := range columns {
				if len(outcomes) == 0 {
					b.WriteString(" / |")
				} else {
					fmt.Fprintf(b, " %d |", outcomes[outcome])
				}
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
}
