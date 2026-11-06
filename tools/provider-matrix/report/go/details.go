package report

import (
	"context"
	"encoding/base64"
	"fmt"
	"go/types"
	"io"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	capability "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
)

func stableKey(parts ...string) string {
	var tuple strings.Builder
	for _, part := range parts {
		fmt.Fprintf(&tuple, "%d:%s", len(part), part)
	}
	return base64.RawURLEncoding.EncodeToString([]byte(tuple.String()))
}

func operations(m *matrix) []*operation {
	byID := map[string]*operation{}
	add := func(locations []location) {
		for _, loc := range locations {
			if loc.op != nil {
				byID[loc.op.id] = loc.op
			}
		}
	}
	for _, row := range m.rows {
		for _, cell := range row.cells {
			add(cell)
		}
	}
	for _, edge := range m.edges {
		for _, cell := range edge {
			add(cell)
		}
	}
	result := make([]*operation, 0, len(byID))
	for _, op := range byID {
		result = append(result, op)
	}
	slices.SortFunc(result, func(a, b *operation) int { return strings.Compare(a.id, b.id) })
	return result
}

func operationTargets(op *operation) []capability.Symbol {
	set := map[capability.Symbol]bool{}
	add := func(fn *types.Func) {
		if fn != nil {
			set[canonicalSymbol(fn)] = true
		}
	}
	add(op.public)
	add(op.entry)
	for _, call := range op.chain {
		add(call.caller)
		add(call.callee)
	}
	if op.kind == "identity" {
		add(method(op.language.marker, "ID"))
		add(method(op.language.marker, "Name"))
	}
	result := make([]capability.Symbol, 0, len(set))
	for symbol := range set {
		result = append(result, symbol)
	}
	slices.SortFunc(result, func(a, b capability.Symbol) int { return strings.Compare(a.String(), b.String()) })
	return result
}

func markerMatches(op *operation, marker capability.Symbol) bool {
	if marker.Package == "" {
		return true
	}
	if op.language.marker != nil && canonicalSymbol(op.language.marker.Obj()) == marker {
		return true
	}
	return op.from != nil && canonicalSymbol(op.from.marker.Obj()) == marker || op.to != nil && canonicalSymbol(op.to.marker.Obj()) == marker
}

func targetContext(op *operation, target capability.Symbol) string {
	parts := []string{op.kind, target.String(), fullType(op.contract)}
	if op.language.marker != nil {
		parts = append(parts, canonicalSymbol(op.language.marker.Obj()).String())
	}
	if op.from != nil {
		parts = append(parts, canonicalSymbol(op.from.marker.Obj()).String(), canonicalSymbol(op.to.marker.Obj()).String())
	}
	return stableKey(parts...)
}

func associateCases(m *matrix, a *annotations) error {
	for _, c := range a.cases {
		contexts := map[string]bool{}
		for _, op := range operations(m) {
			if markerMatches(op, c.marker) && slices.Contains(operationTargets(op), c.target) {
				contexts[targetContext(op, c.target)] = true
			}
		}
		if len(contexts) == 0 {
			return fmt.Errorf("%s: case target %s is not a discovered operation for the selected language", c.test.file, c.target)
		}
		if len(contexts) > 1 {
			return fmt.Errorf("%s: ambiguous case target %s; supply its exact language marker or a more specific operation", c.test.file, c.target)
		}
	}
	return nil
}

type definitionReport struct {
	root          capability.Symbol
	format, state string
	inventory     capability.Inventory
}

type queryView struct {
	base      capability.SourceView
	ctx       context.Context
	mu        sync.Mutex
	remaining int
	closed    bool
}

func (v *queryView) admit() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.closed {
		return fmt.Errorf("definition source query after extraction returned")
	}
	if err := v.ctx.Err(); err != nil {
		return err
	}
	v.remaining--
	if v.remaining < 0 {
		return fmt.Errorf("definition source query limit exceeded")
	}
	return nil
}
func (v *queryView) File(ref capability.FileRef) (capability.File, error) {
	if err := v.admit(); err != nil {
		return capability.File{}, err
	}
	return v.base.File(ref)
}
func (v *queryView) Declaration(symbol capability.Symbol) (capability.Declaration, error) {
	if err := v.admit(); err != nil {
		return capability.Declaration{}, err
	}
	return v.base.Declaration(symbol)
}
func (v *queryView) Facts(span capability.Span) (capability.Facts, error) {
	if err := v.admit(); err != nil {
		return capability.Facts{}, err
	}
	return v.base.Facts(span)
}

func extractDefinition(ctx context.Context, source capability.SourceView, root capability.Symbol, bindings []Binding, limits capability.Limits) (definitionReport, error) {
	result := definitionReport{root: root, state: "no matching extractor"}
	decl, err := source.Declaration(root)
	if err != nil {
		result.state = "source unavailable"
		return result, nil
	}
	span := decl.Value
	if span.File == "" {
		span = decl.Type
	}
	var origin *capability.Symbol
	if span.File != "" {
		facts, err := source.Facts(span)
		if err != nil {
			return result, err
		}
		origin = facts.Type.Origin
	}
	var selected *Binding
	formats := map[string]bool{}
	for i := range bindings {
		binding := &bindings[i]
		if binding.Format == "" || binding.Extractor == nil || formats[binding.Format] {
			return result, fmt.Errorf("extractors require unique nonempty format identities and implementations")
		}
		formats[binding.Format] = true
		matches := false
		selectors := binding.Extractor.Selectors()
		if len(selectors) == 0 || len(selectors) > 256 {
			return result, fmt.Errorf("extractor %s has invalid selector count", binding.Format)
		}
		for _, selector := range selectors {
			if (selector.Definition == nil) == (selector.NamedType == nil) {
				return result, fmt.Errorf("extractor %s selector requires exactly one identity", binding.Format)
			}
			matches = matches || selector.Definition != nil && *selector.Definition == root || selector.NamedType != nil && origin != nil && *selector.NamedType == *origin
		}
		if matches {
			if selected != nil {
				return result, fmt.Errorf("ambiguous extractors for %s: %s and %s", root, selected.Format, binding.Format)
			}
			selected = binding
		}
	}
	if selected == nil {
		return result, nil
	}
	limits, err = capability.EffectiveLimits(limits)
	if err != nil {
		return result, err
	}
	view := &queryView{base: source, ctx: ctx, remaining: limits.MaxExpansions}
	request := capability.DefinitionRequest{Root: root, Source: view, Limits: limits}
	inventory, err := selected.Extractor.Extract(ctx, request)
	view.mu.Lock()
	view.closed = true
	view.mu.Unlock()
	if err != nil {
		return result, fmt.Errorf("extract %s with %s: %w", root, selected.Format, err)
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	// Validation uses a separately bounded view after the extractor's lease ends.
	validation := &queryView{base: source, ctx: ctx, remaining: limits.MaxExpansions}
	request.Source = validation
	if err := capability.Validate(request, inventory); err != nil {
		return result, fmt.Errorf("invalid inventory from %s: %w", selected.Format, err)
	}
	result.format, result.state, result.inventory = selected.Format, string(inventory.Extent), cloneInventory(inventory)
	return result, nil
}

func cloneInventory(inventory capability.Inventory) capability.Inventory {
	inventory.Features = slices.Clone(inventory.Features)
	inventory.Limitations = slices.Clone(inventory.Limitations)
	for i := range inventory.Features {
		f := &inventory.Features[i]
		f.Roles.Value = slices.Clone(f.Roles.Value)
		f.Attributes.Value = slices.Clone(f.Attributes.Value)
		f.References.Value = slices.Clone(f.References.Value)
		for j := range f.Roles.Value {
			f.Roles.Value[j].Kinds.Value = slices.Clone(f.Roles.Value[j].Kinds.Value)
		}
		for j := range f.Attributes.Value {
			f.Attributes.Value[j].Values.Value.Values = slices.Clone(f.Attributes.Value[j].Values.Value.Values)
		}
	}
	return inventory
}

type limitedText struct {
	strings.Builder
	err error
}

func (b *limitedText) Write(p []byte) (int, error) {
	if b.err != nil {
		return 0, b.err
	}
	if b.Len()+len(p) > 64<<20 {
		b.err = fmt.Errorf("rendered capability report exceeds 64 MiB")
		return 0, b.err
	}
	return b.Builder.Write(p)
}
func (b *limitedText) WriteString(s string) (int, error) { return b.Write([]byte(s)) }

func profileLabel(p profile) string {
	if !p.present {
		return "unspecified (authored profile absent)"
	}
	return p.value
}

func bindingAnchor(op *operation, target capability.Symbol, entry authored) string {
	return stableKey(op.id, target.String(), entry.directive.name, entry.root.String(), entry.document, fmt.Sprint(entry.directive.profile.present), entry.directive.profile.value)
}

func sourceTarget(root, sourceURL, path string, line int) string {
	if filepath.IsAbs(path) {
		path, _ = filepath.Rel(root, path)
	}
	parts := strings.Split(filepath.ToSlash(path), "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	target := strings.TrimSuffix(sourceURL, "/") + "/" + strings.Join(parts, "/")
	if line > 0 {
		target += fmt.Sprintf("#L%d", line)
	}
	return target
}

func (p *program) symbolLink(symbol capability.Symbol, sourceURL string) string {
	pkg := p.packages[symbol.Package]
	if pkg == nil {
		return escape(symbol.String())
	}
	var obj types.Object
	if symbol.Receiver == "" {
		obj = pkg.types.Scope().Lookup(symbol.Name)
	} else if typ := pkg.types.Scope().Lookup(symbol.Receiver); typ != nil {
		obj = method(types.NewPointer(typ.Type()), symbol.Name)
	}
	if obj == nil {
		return escape(symbol.String())
	}
	pos := p.fset.PositionFor(obj.Pos(), false)
	return fmt.Sprintf("[%s](%s)", escape(symbol.String()), sourceTarget(p.root, sourceURL, pos.Filename, pos.Line))
}

func (p *program) renderDetails(ctx context.Context, m *matrix, a *annotations, config Config, sourceURL string) (string, error) {
	var b limitedText
	b.WriteString("# Provider capability details\n\nGenerated from selected production declarations and explicit test associations.\nThese reports describe bounded implementation scope, not language completeness.\n\n")
	fmt.Fprintf(&b, "Inspection build: %s, %s/%s; Go flags %s; explicit build flags %s.\n\n", escape(p.build.GOVERSION), escape(p.build.GOOS), escape(p.build.GOARCH), escape(p.build.GOFLAGS), escape(strings.Join(config.BuildFlags, " ")))
	b.WriteString("Definitions describe selected initializers; mutation, operation bodies and procedural semantics are not inferred. Profiles and expectations are authored assertions. A test function's observation is not separate execution of every associated case, nor evidence that a wrapper or every branch ran.\n\n")
	contractLink, matrixLink := "capability-reports.md", "providers.md"
	if config.Context != nil {
		contractLink = sourceTarget(p.root, sourceURL, "docs/capability-reports.md", 0)
		matrixLink = sourceTarget(p.root, sourceURL, config.MatrixPath, 0)
	}
	fmt.Fprintf(&b, "[Report contract](%s) · [Provider matrices](%s) · [Existing CI invocations](https://%s/actions/workflows/go.yml)\n\n", contractLink, matrixLink, p.module)
	if len(config.Executions) == 0 {
		b.WriteString("No execution observations are included in this committed inventory. Per-invocation CI artifacts contain observations with their actual source, build selection and process result. Missing observations do not mean a test never ran.\n\n")
	}
	source := newSourceView(p)
	definitions := map[capability.Symbol]definitionReport{}
	for _, op := range operations(m) {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		fmt.Fprintf(&b, "<a id=\"cap-%s\"></a>\n\n## %s\n\n", op.id, escape(operationTitle(op)))
		fmt.Fprintf(&b, "Contract: %s.\n\n", escape(fullType(op.contract)))
		if op.public != nil {
			fmt.Fprintf(&b, "Public operation: %s.", p.symbolLink(canonicalSymbol(op.public), sourceURL))
			if op.receiver != nil {
				fmt.Fprintf(&b, " Receiver: %s.", escape(fullType(op.receiver)))
			}
			b.WriteString("\n\n")
		}
		if op.from != nil {
			fmt.Fprintf(&b, "Direct endpoints: %s → %s. Typed calls in the factory body establish this discovered edge; runtime reachability, returned-value flow, configured mappings and composed routes are not inferred.\n\n", escape(fullType(op.from.marker)), escape(fullType(op.to.marker)))
		} else {
			fmt.Fprintf(&b, "Language marker: %s.\n\n", escape(fullType(op.language.marker)))
		}
		if len(op.alternatives) > 0 {
			var alternatives []string
			for _, set := range op.alternatives {
				var requirements []string
				for _, dep := range set {
					requirements = append(requirements, fullType(dep))
				}
				slices.Sort(requirements)
				if len(requirements) == 0 {
					requirements = []string{"no provider collaborators"}
				}
				alternatives = append(alternatives, strings.Join(requirements, " + "))
			}
			slices.Sort(alternatives)
			alternatives = slices.Compact(alternatives)
			fmt.Fprintf(&b, "Construction alternatives: %s.\n\n", escape(strings.Join(alternatives, " OR ")))
		}
		if len(op.chain) > 0 {
			b.WriteString("Static call relationships (no evidence transfer):\n\n")
			for _, call := range op.chain {
				pos := p.fset.PositionFor(call.pos, false)
				var args []string
				for _, arg := range call.arguments {
					args = append(args, fullType(arg))
				}
				fmt.Fprintf(&b, "- %s → %s ([call](%s)); receiver %s; type arguments %s.\n", p.symbolLink(canonicalSymbol(call.caller), sourceURL), p.symbolLink(canonicalSymbol(call.callee), sourceURL), sourceTarget(p.root, sourceURL, pos.Filename, pos.Line), escape(fullType(call.receiver)), escape(strings.Join(args, ", ")))
			}
			fmt.Fprintf(&b, "\nTraversal boundary: %s.\n\n", escape(op.stop))
		}
		hasScope, hasDefinition, hasCases := false, false, false
		for _, target := range operationTargets(op) {
			entries := a.operations[target]
			var cases []caseAssociation
			for _, c := range a.cases {
				if c.target == target && markerMatches(op, c.marker) {
					cases = append(cases, c)
				}
			}
			if len(entries) == 0 && len(cases) == 0 {
				continue
			}
			fmt.Fprintf(&b, "### Declared target %s\n\n", p.symbolLink(target, sourceURL))
			if canonicalSymbol(op.public) != target && op.kind != "identity" {
				b.WriteString("Related declaration; its cases and observations do not qualify the public wrapper.\n\n")
			}
			for _, entry := range entries {
				fmt.Fprintf(&b, "<a id=\"binding-%s\"></a>\n\n", bindingAnchor(op, target, entry))
				if entry.directive.name == "scope" {
					hasScope = true
					name, frag, _ := strings.Cut(entry.document, "#")
					link := sourceTarget(p.root, sourceURL, name, 0)
					if frag != "" {
						link += "#" + url.PathEscape(frag)
					}
					fmt.Fprintf(&b, "Authored scope: [%s](%s); profile: %s.\n\n", escape(entry.document), link, escape(profileLabel(entry.directive.profile)))
					continue
				}
				hasDefinition = true
				definition, ok := definitions[entry.root]
				if !ok {
					var err error
					definition, err = extractDefinition(ctx, source, entry.root, config.Extractors, config.Limits)
					if err != nil {
						return "", err
					}
					definitions[entry.root] = definition
				}
				fmt.Fprintf(&b, "Definition: %s; format: %s; extraction: **%s**; profile: %s.\n\n", p.symbolLink(entry.root, sourceURL), escape(definition.format), escape(definition.state), escape(profileLabel(entry.directive.profile)))
				renderInventory(&b, p, source, sourceURL, definition, cases, entry.directive.profile, entries)
			}
			if len(cases) > 0 {
				hasCases = true
				b.WriteString("| Feature | Case | Authored expectation | Authored profile | Test declaration |\n| --- | --- | --- | --- | --- |\n")
				for _, c := range cases {
					pos := p.fset.PositionFor(c.test.pos, false)
					observed := ""
					if len(config.Executions) > 0 {
						observed = " [observations](#test-" + c.test.id + ")"
					}
					fmt.Fprintf(&b, "| %s | %s | %s | %s | [%s](%s) (%s)%s |\n", escape(c.directive.values["feature"]), escape(c.directive.values["case"]), escape(c.directive.values["expects"]), escape(profileLabel(c.directive.profile)), escape(c.test.name), sourceTarget(p.root, sourceURL, c.test.file, pos.Line), escape(c.test.variant), observed)
				}
				b.WriteString("\n")
			}
		}
		if !hasScope {
			b.WriteString("No authored scope link attached to the discovered operation or its static related declarations.\n\n")
		}
		if !hasDefinition {
			b.WriteString("No structural definition attached. Operation source and authored scope remain the boundary; no grammar is inferred from AST kinds, switches or tests.\n\n")
		}
		if !hasCases {
			b.WriteString("No associated cases found in the selected test sources. This does not establish absence of tests.\n\n")
		}
		if b.err != nil {
			return "", b.err
		}
	}
	if len(config.Executions) > 0 {
		renderObservations(&b, p, a, config, sourceURL)
	}
	if b.err != nil {
		return "", b.err
	}
	return b.String(), nil
}

func operationTitle(op *operation) string {
	if op.kind == "identity" {
		return op.language.name + " language identity"
	}
	if op.from != nil {
		return op.from.name + " → " + op.to.name + " / " + canonicalSymbol(op.public).String()
	}
	return op.language.name + " / " + canonicalSymbol(op.public).String()
}

func spanLink(p *program, source capability.SourceView, sourceURL string, span capability.Span) string {
	file, err := source.File(span.File)
	if err != nil {
		return "source unavailable"
	}
	line := 1 + strings.Count(file.Text[:span.Start], "\n")
	return "[source](" + sourceTarget(p.root, sourceURL, file.Path, line) + ")"
}

func knowledge[T any](field capability.Field[T], known func(T) string) string {
	if field.State == capability.Known {
		return known(field.Value)
	}
	if field.State == capability.Partial {
		return "partial: " + known(field.Value) + " (" + field.Reason + ")"
	}
	return string(field.State) + ": " + field.Reason
}

func renderInventory(b io.Writer, p *program, source capability.SourceView, sourceURL string, d definitionReport, cases []caseAssociation, profile profile, entries []authored) {
	for _, limitation := range d.inventory.Limitations {
		fmt.Fprintf(b, "- Extraction limit: %s (%s).\n", escape(limitation.Message), spanLink(p, source, sourceURL, limitation.At))
	}
	if len(d.inventory.Limitations) > 0 {
		fmt.Fprintln(b)
	}
	for _, feature := range d.inventory.Features {
		fmt.Fprintf(b, "#### %s (%s)\n\n", escape(feature.Key), spanLink(p, source, sourceURL, feature.At))
		fmt.Fprintf(b, "Feature identity: definition %s / format %s / key %s.\n\n", escape(d.root.String()), escape(d.format), escape(feature.Key))
		fmt.Fprintf(b, "Roles: %s. Attributes: %s. References: %s.\n\n", escape(knowledge(feature.Roles, func(v []capability.Role) string { return fmt.Sprintf("%d declared", len(v)) })), escape(knowledge(feature.Attributes, func(v []capability.Attribute) string { return fmt.Sprintf("%d declared", len(v)) })), escape(knowledge(feature.References, func(v []string) string {
			if len(v) == 0 {
				return "none"
			}
			return strings.Join(v, ", ")
		})))
		if len(feature.Roles.Value) > 0 {
			fmt.Fprintln(b, "| Role | Minimum | Maximum | Allowed kinds |\n| --- | --- | --- | --- |")
			for _, role := range feature.Roles.Value {
				fmt.Fprintf(b, "| %s | %s | %s | %s |\n", escape(role.Name), escape(knowledge(role.Min, func(v int) string { return fmt.Sprint(v) })), escape(knowledge(role.Max, func(v capability.UpperBound) string {
					if v.Unbounded {
						return "unbounded"
					}
					return fmt.Sprint(v.Value)
				})), escape(knowledge(role.Kinds, func(v []string) string { return strings.Join(v, ", ") })))
			}
			fmt.Fprintln(b)
		}
		if len(feature.Attributes.Value) > 0 {
			fmt.Fprintln(b, "| Attribute | Required | Values or predicate |\n| --- | --- | --- |")
			for _, attr := range feature.Attributes.Value {
				fmt.Fprintf(b, "| %s | %s | %s |\n", escape(attr.Name), escape(knowledge(attr.Required, func(v bool) string { return fmt.Sprint(v) })), escape(knowledge(attr.Values, func(v capability.ValueSet) string {
					if v.Unrestricted {
						return "unrestricted"
					}
					if len(v.Values) == 0 {
						return "empty set"
					}
					return strings.Join(v.Values, ", ")
				})))
			}
			fmt.Fprintln(b)
		}
		var matches []string
		for _, c := range cases {
			if c.directive.profile == profile && c.directive.values["feature"] == feature.Key {
				matches = append(matches, c.directive.values["case"]+" / "+c.test.name)
			}
		}
		roots := map[capability.Symbol]bool{}
		for _, entry := range entries {
			if entry.directive.name == "grammar" && entry.directive.profile == profile {
				roots[entry.root] = true
			}
		}
		if len(matches) == 0 {
			fmt.Fprint(b, "No feature-key case association at this exact authored profile.\n\n")
		} else if len(roots) > 1 {
			fmt.Fprint(b, "Matching operation-scoped case labels exist; multiple roots prevent assigning them to this inventory item.\n\n")
		} else {
			fmt.Fprintf(b, "Associated labels at this exact authored profile: %s. Function-level associations do not establish all scenarios.\n\n", escape(strings.Join(matches, "; ")))
		}
	}
}
