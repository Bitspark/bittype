package report

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"slices"
	"strings"
)

type language struct {
	id, name string
	marker   *types.Named
}
type location struct {
	op *operation
}

// An operation keeps compiler identities until presentation. A delegate edge is
// a source relationship only; it does not establish wrapper equivalence or carry
// a delegate's qualification to the public method.
type operation struct {
	id            string
	kind          string
	public, entry *types.Func
	chain         []callLink
	stop          string
	contract      *types.Named
	receiver      types.Type
	language      language
	from, to      *language
	alternatives  [][]*types.Named
}

type callLink struct {
	caller, callee *types.Func
	pos            token.Pos
	receiver       types.Type
	arguments      []types.Type
}

func typeLabel(t types.Type) string {
	if t == nil {
		return ""
	}
	return types.TypeString(t, func(pkg *types.Package) string { return pkg.Path() })
}

func operationID(kind string, fn *types.Func, contract *types.Named, receiver types.Type, markers ...language) string {
	parts := []string{kind, canonicalSymbol(fn).String(), typeLabel(contract), typeLabel(receiver)}
	for _, marker := range markers {
		parts = append(parts, canonicalSymbol(marker.marker.Obj()).String(), marker.id)
	}
	return stableKey(parts...)
}

type row struct {
	level, contract string
	identity        types.Type
	cells           map[string][]location
}
type matrix struct {
	languages []language
	rows      map[string]*row
	edges     map[string]map[string][]location
}
type contribution struct {
	language language
	contract *types.Named
	op       *operation
	// Each alternative is a public constructor's provider requirements.
	// Concrete implementations use one empty alternative.
	alternatives [][]*types.Named
}

func (p *program) constantString(fn *types.Func) (string, bool) {
	if fn == nil {
		return "", false
	}
	decl, ok := p.functions[fn.Origin()]
	if !ok || decl.node.Body == nil || len(decl.node.Body.List) != 1 {
		return "", false
	}
	ret, ok := decl.node.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return "", false
	}
	value := decl.pkg.info.Types[ret.Results[0]].Value
	if value == nil || value.Kind() != constant.String {
		return "", false
	}
	return constant.StringVal(value), true
}

func (p *program) languages(identity *types.Named) ([]language, error) {
	var result []language
	ids := map[string]bool{}
	for path, pkg := range p.packages {
		if !strings.HasPrefix(path, p.module+"/syntax/") {
			continue
		}
		for _, name := range pkg.types.Scope().Names() {
			obj, ok := pkg.types.Scope().Lookup(name).(*types.TypeName)
			if !ok || !obj.Exported() || obj.IsAlias() {
				continue
			}
			named, ok := obj.Type().(*types.Named)
			if !ok || named.TypeParams().Len() != 0 || !types.Implements(named, identity.Underlying().(*types.Interface)) {
				continue
			}
			id, idOK := p.constantString(method(named, "ID"))
			label, labelOK := p.constantString(method(named, "Name"))
			if !idOK || !labelOK || id == "" || label == "" {
				return nil, fmt.Errorf("%s: language ID and Name must return constant strings for source discovery", obj)
			}
			if ids[id] {
				return nil, fmt.Errorf("duplicate language ID %q", id)
			}
			ids[id] = true
			result = append(result, language{id, label, named})
		}
	}
	slices.SortFunc(result, func(a, b language) int { return strings.Compare(a.id, b.id) })
	if len(result) == 0 {
		return nil, fmt.Errorf("no language markers discovered")
	}
	return result, nil
}

// Infer the contract's generic arguments from the implementation's method set.
// Match Go types, including instantiated carriers, rather than spellings/import aliases.
func match(contract *types.Named, implementation types.Type) *types.Named {
	iface, ok := contract.Underlying().(*types.Interface)
	if !ok {
		return nil
	}
	bindings := map[*types.TypeParam]types.Type{}
	for i := 0; i < iface.NumMethods(); i++ {
		expected := iface.Method(i)
		actual := method(implementation, expected.Name())
		if actual == nil || !unify(expected.Type(), actual.Type(), bindings) {
			return nil
		}
	}
	args := make([]types.Type, contract.TypeParams().Len())
	for i := range args {
		args[i] = bindings[contract.TypeParams().At(i)]
		if args[i] == nil {
			return nil
		}
	}
	var instantiated types.Type = contract
	if len(args) > 0 {
		var err error
		instantiated, err = types.Instantiate(nil, contract, args, true)
		if err != nil {
			return nil
		}
	}
	if !types.Implements(implementation, instantiated.Underlying().(*types.Interface)) {
		return nil
	}
	return instantiated.(*types.Named)
}

func unify(pattern, actual types.Type, bindings map[*types.TypeParam]types.Type) bool {
	pattern, actual = types.Unalias(pattern), types.Unalias(actual)
	if parameter, ok := pattern.(*types.TypeParam); ok {
		if previous := bindings[parameter]; previous != nil {
			return types.Identical(previous, actual)
		}
		bindings[parameter] = actual
		return true
	}
	switch a := pattern.(type) {
	case *types.Named:
		b, ok := actual.(*types.Named)
		if !ok || a.Origin() != b.Origin() || a.TypeArgs().Len() != b.TypeArgs().Len() {
			return false
		}
		for i := 0; i < a.TypeArgs().Len(); i++ {
			if !unify(a.TypeArgs().At(i), b.TypeArgs().At(i), bindings) {
				return false
			}
		}
		return true
	case *types.Signature:
		b, ok := actual.(*types.Signature)
		return ok && a.Variadic() == b.Variadic() && unify(a.Params(), b.Params(), bindings) && unify(a.Results(), b.Results(), bindings)
	case *types.Tuple:
		b, ok := actual.(*types.Tuple)
		if !ok || a.Len() != b.Len() {
			return false
		}
		for i := 0; i < a.Len(); i++ {
			if !unify(a.At(i).Type(), b.At(i).Type(), bindings) {
				return false
			}
		}
		return true
	case *types.Pointer:
		b, ok := actual.(*types.Pointer)
		return ok && unify(a.Elem(), b.Elem(), bindings)
	case *types.Slice:
		b, ok := actual.(*types.Slice)
		return ok && unify(a.Elem(), b.Elem(), bindings)
	case *types.Array:
		b, ok := actual.(*types.Array)
		return ok && a.Len() == b.Len() && unify(a.Elem(), b.Elem(), bindings)
	case *types.Map:
		b, ok := actual.(*types.Map)
		return ok && unify(a.Key(), b.Key(), bindings) && unify(a.Elem(), b.Elem(), bindings)
	case *types.Chan:
		b, ok := actual.(*types.Chan)
		return ok && a.Dir() == b.Dir() && unify(a.Elem(), b.Elem(), bindings)
	default:
		return types.Identical(pattern, actual)
	}
}

func namedType(t types.Type) *types.Named {
	t = types.Unalias(t)
	if pointer, ok := t.(*types.Pointer); ok {
		t = types.Unalias(pointer.Elem())
	}
	named, _ := t.(*types.Named)
	return named
}

// Dependencies are read from real generic constructor signatures. Operation input
// carriers are deliberately not dependencies: printing Source does not require parsing it.
func (p *program) constructors(named *types.Named, lang language, contracts map[*types.Named]bool) ([][]*types.Named, error) {
	var alternatives [][]*types.Named
	for fn := range p.functions {
		if fn.Pkg() != named.Obj().Pkg() || !fn.Exported() {
			continue
		}
		signature := fn.Type().(*types.Signature)
		if signature.Recv() != nil {
			continue
		}
		returns := false
		for i := 0; i < signature.Results().Len(); i++ {
			if result := namedType(signature.Results().At(i).Type()); result != nil && result.Origin() == named {
				returns = true
			}
		}
		if !returns {
			continue
		}
		if signature.TypeParams().Len() != 1 {
			return nil, fmt.Errorf("%s: generic provider factory needs one discoverable language parameter", fn)
		}
		inst, err := types.Instantiate(nil, signature, []types.Type{lang.marker}, true)
		if err != nil {
			continue
		}
		var requirements []*types.Named
		params := inst.(*types.Signature).Params()
		for i := 0; i < params.Len(); i++ {
			if dep := namedType(params.At(i).Type()); dep != nil && contracts[dep.Origin()] {
				requirements = append(requirements, dep)
			}
		}
		alternatives = append(alternatives, requirements)
	}
	if len(alternatives) == 0 {
		return nil, fmt.Errorf("%s: no discoverable public generic provider constructor for %s", named, lang.id)
	}
	return alternatives, nil
}

func contractName(contract *types.Named, lang language) string {
	name := types.TypeString(contract, func(p *types.Package) string { return p.Name() })
	marker := types.TypeString(lang.marker, func(p *types.Package) string { return p.Name() })
	return strings.ReplaceAll(name, marker, "L")
}

func (p *program) discover() (*matrix, error) {
	identity, err := p.named("contracts/language/go", "Provider")
	if err != nil {
		return nil, err
	}
	languages, err := p.languages(identity)
	if err != nil {
		return nil, err
	}
	contracts := map[*types.Named]bool{}
	for path, pkg := range p.packages {
		if !strings.HasPrefix(path, p.module+"/contracts/") {
			continue
		}
		obj := pkg.types.Scope().Lookup("Provider")
		if obj == nil {
			continue
		}
		named, ok := obj.Type().(*types.Named)
		if !ok {
			continue
		}
		iface, ok := named.Underlying().(*types.Interface)
		if ok && types.Implements(iface, identity.Underlying().(*types.Interface)) {
			contracts[named] = true
		}
	}
	m := &matrix{languages: languages, rows: map[string]*row{}, edges: map[string]map[string][]location{}}
	identityRow := &row{level: p.level(identity), contract: "language.Provider", identity: identity, cells: map[string][]location{}}
	m.rows[identityRow.contract] = identityRow
	parameter := rowParameter()
	selectRow := func(c contribution) (*row, error) {
		key := contractName(c.contract, c.language)
		normalized, err := rowType(c.contract, c.language.marker, parameter)
		if err != nil {
			return nil, fmt.Errorf("matrix row %s: %w", key, err)
		}
		if existing := m.rows[key]; existing != nil {
			if !types.Identical(existing.identity, normalized) {
				return nil, fmt.Errorf("ambiguous matrix contract label %q: distinct Go identities share this display spelling", key)
			}
			return existing, nil
		}
		result := &row{level: p.level(c.contract), contract: key, identity: normalized, cells: map[string][]location{}}
		m.rows[key] = result
		return result, nil
	}
	available := map[string]*typeSet{}
	for _, lang := range languages {
		op := &operation{kind: "identity", contract: identity, receiver: lang.marker, language: lang}
		op.id = operationID(op.kind, nil, identity, lang.marker, lang)
		identityRow.cells[lang.id] = []location{{op: op}}
		available[lang.id] = &typeSet{identity}
	}
	var pending []contribution
	for path, pkg := range p.packages {
		if !p.implementation(path) || strings.HasPrefix(path, p.module+"/syntax/") {
			continue
		}
		for _, name := range pkg.types.Scope().Names() {
			obj, ok := pkg.types.Scope().Lookup(name).(*types.TypeName)
			if !ok || !obj.Exported() || obj.IsAlias() {
				continue
			}
			named, ok := obj.Type().(*types.Named)
			if !ok {
				continue
			}
			if _, ok := named.Underlying().(*types.Interface); ok {
				continue
			}
			for _, lang := range languages {
				var concrete types.Type = named
				generic := named.TypeParams().Len() > 0
				if generic {
					// This repository's reusable providers have one language marker. Do
					// not silently misdocument a provider with additional open parameters.
					if named.TypeParams().Len() != 1 {
						if method(named, "ID") != nil && method(named, "Name") != nil {
							return nil, fmt.Errorf("%s: multiple open provider parameters require explicit discovery support", named)
						}
						continue
					}
					concrete, err = types.Instantiate(nil, named, []types.Type{lang.marker}, true)
					if err != nil {
						continue
					}
				} else {
					id, ok := p.constantString(method(concrete, "ID"))
					if !ok || id != lang.id {
						continue
					}
				}
				for _, candidate := range []types.Type{concrete, types.NewPointer(concrete)} {
					for contract := range contracts {
						if contract == identity {
							continue
						}
						bound := match(contract, candidate)
						if bound == nil {
							continue
						}
						alternatives := [][]*types.Named{{}}
						if generic {
							alternatives, err = p.constructors(named, lang, contracts)
							if err != nil {
								return nil, err
							}
						}
						iface := bound.Underlying().(*types.Interface)
						for i := 0; i < iface.NumExplicitMethods(); i++ {
							public := method(candidate, iface.ExplicitMethod(i).Name())
							fn, chain, stop := p.entryChain(public)
							if fn != nil {
								receiver := public.Type().(*types.Signature).Recv().Type()
								op := &operation{kind: "provider", public: public, entry: fn, chain: chain, stop: stop, contract: bound, receiver: receiver, language: lang, alternatives: alternatives}
								op.id = operationID(op.kind, public, bound, receiver, lang)
								pending = append(pending, contribution{language: lang, contract: bound, op: op, alternatives: alternatives})
							}
						}
					}
				}
			}
		}
	}
	// Least fixed point: a cycle of generic collaborators cannot create availability.
	for {
		progress := false
		remaining := pending[:0]
		for _, c := range pending {
			ready := false
			for _, deps := range c.alternatives {
				all := true
				for _, dep := range deps {
					all = all && available[c.language.id].contains(dep)
				}
				ready = ready || all
			}
			if !ready {
				remaining = append(remaining, c)
				continue
			}
			row, err := selectRow(c)
			if err != nil {
				return nil, err
			}
			duplicate := false
			for _, loc := range row.cells[c.language.id] {
				duplicate = duplicate || sameOperation(loc.op, c.op)
			}
			if !duplicate {
				row.cells[c.language.id] = append(row.cells[c.language.id], location{op: c.op})
			}
			available[c.language.id].add(c.contract)
			progress = true
		}
		pending = remaining
		if !progress {
			break
		}
	}
	// Keep discovered but unsatisfied generic pairings visible as unavailable rows.
	for _, c := range pending {
		if _, err := selectRow(c); err != nil {
			return nil, err
		}
	}
	if err := p.mappingProviders(m); err != nil {
		return nil, err
	}
	if err := p.transformers(m); err != nil {
		return nil, err
	}
	if err := validateOperationAnchors(m); err != nil {
		return nil, err
	}
	return m, nil
}

// mappingProviders reports the endpoints in actual implementations of the public
// generation contract. A supplied mapping is available for selection. A mapping
// alone makes no claim about a configured route, source evidence, or roots.
func (p *program) mappingProviders(m *matrix) error {
	if p.packages[p.module+"/contracts/generation/go"] == nil {
		return nil
	}
	contract, err := p.named("contracts/generation/go", "MappingProvider")
	if err != nil {
		return err
	}
	marker := func(t types.Type) *language {
		for i := range m.languages {
			if types.Identical(t, m.languages[i].marker) {
				return &m.languages[i]
			}
		}
		return nil
	}
	for path, pkg := range p.packages {
		if !p.implementation(path) {
			continue
		}
		for _, name := range pkg.types.Scope().Names() {
			obj, ok := pkg.types.Scope().Lookup(name).(*types.TypeName)
			if !ok || !obj.Exported() || obj.IsAlias() {
				continue
			}
			named, ok := obj.Type().(*types.Named)
			if !ok || named.TypeParams().Len() != 0 {
				continue
			}
			if _, ok := named.Underlying().(*types.Interface); ok {
				continue
			}
			candidate := types.NewPointer(named)
			bound := match(contract, candidate)
			if bound == nil {
				continue
			}
			from, to := marker(bound.TypeArgs().At(0)), marker(bound.TypeArgs().At(1))
			if from == nil || to == nil {
				continue
			}
			public := method(candidate, "Open")
			entry, chain, stop := p.entryChain(public)
			op := &operation{kind: "mapping provider", public: public, entry: entry, chain: chain, stop: stop, contract: bound, receiver: public.Type().(*types.Signature).Recv().Type(), from: from, to: to}
			op.id = operationID(op.kind, public, bound, op.receiver, *from, *to)
			if m.edges[from.id] == nil {
				m.edges[from.id] = map[string][]location{}
			}
			m.edges[from.id][to.id] = append(m.edges[from.id][to.id], location{op})
		}
	}
	return nil
}

func (p *program) transformers(m *matrix) error {
	contract, err := p.named("contracts/transform/go", "ASTTransformer")
	if err != nil {
		return err
	}
	iface := contract.Underlying().(*types.Interface)
	lookup := func(t types.Type) string {
		for _, lang := range m.languages {
			if types.Identical(t, lang.marker) {
				return lang.id
			}
		}
		return ""
	}
	add := func(from, to string, loc location) {
		if m.edges[from] == nil {
			m.edges[from] = map[string][]location{}
		}
		m.edges[from][to] = append(m.edges[from][to], loc)
	}
	returnsTransformer := func(fn *types.Func) bool {
		signature := fn.Type().(*types.Signature)
		for i := 0; i < signature.Results().Len(); i++ {
			if types.Implements(signature.Results().At(i).Type(), iface) {
				return true
			}
		}
		return false
	}
	acceptsTransformer := func(fn *types.Func) bool {
		signature := fn.Type().(*types.Signature)
		for i := 0; i < signature.Params().Len(); i++ {
			if types.Implements(signature.Params().At(i).Type(), iface) {
				return true
			}
		}
		return false
	}
	// Follow typed factory delegation to the existing transformer boundary.
	// Generic helpers are not concrete edges. Bind their actual type arguments
	// at each call, including reordered parameters, without executing factories.
	type endpoint struct {
		from, to string
		chain    []callLink
	}
	var endpoints func(*types.Func, map[*types.TypeParam]types.Type, map[*types.Func]bool) ([]endpoint, error)
	endpoints = func(fn *types.Func, bindings map[*types.TypeParam]types.Type, active map[*types.Func]bool) ([]endpoint, error) {
		fn = fn.Origin()
		if active[fn] {
			return nil, fmt.Errorf("%s: cyclic transformer factory delegation", fn)
		}
		decl, ok := p.functions[fn]
		if !ok || decl.node.Body == nil {
			return nil, fmt.Errorf("%s: transformer factory source unavailable", fn)
		}
		active[fn] = true
		defer delete(active, fn)
		var result []endpoint
		var failure error
		ast.Inspect(decl.node.Body, func(n ast.Node) bool {
			if failure != nil {
				return false
			}
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			callee := p.called(decl.pkg, call.Fun)
			if callee == nil || !returnsTransformer(callee) {
				return true
			}
			var ident *ast.Ident
			ast.Inspect(call.Fun, func(n ast.Node) bool {
				if id, ok := n.(*ast.Ident); ok && decl.pkg.info.Uses[id] == callee {
					ident = id
				}
				return true
			})
			instance := decl.pkg.info.Instances[ident]
			parameters := callee.Origin().Type().(*types.Signature).TypeParams()
			if parameters.Len() != instance.TypeArgs.Len() {
				failure = fmt.Errorf("%s: transformer arguments are not statically instantiated", fn)
				return false
			}
			actuals := make([]types.Type, parameters.Len())
			nested := map[*types.TypeParam]types.Type{}
			for i := range actuals {
				actual := types.Unalias(instance.TypeArgs.At(i))
				if parameter, ok := actual.(*types.TypeParam); ok {
					actual = bindings[parameter]
				}
				if actual == nil {
					failure = fmt.Errorf("%s: open transformer endpoint", fn)
					return false
				}
				actuals[i] = actual
				nested[parameters.At(i)] = actual
			}
			link := callLink{caller: fn, callee: callee.Origin(), pos: call.Pos(), receiver: callReceiver(decl.pkg, call.Fun), arguments: actuals}
			if callee.Pkg().Path() == p.module+"/contracts/transform/go" && callee.Name() == "NewTransformer" {
				if len(actuals) != 2 {
					failure = fmt.Errorf("%s: transformer requires two endpoint arguments", fn)
					return false
				}
				from, to := lookup(actuals[0]), lookup(actuals[1])
				if from == "" || to == "" {
					failure = fmt.Errorf("%s: transformer uses an undiscovered language marker", fn)
					return false
				}
				result = append(result, endpoint{from: from, to: to, chain: []callLink{link}})
				return false
			}
			child, err := endpoints(callee, nested, active)
			if err != nil {
				failure = err
				return false
			}
			for _, edge := range child {
				edge.chain = append([]callLink{link}, edge.chain...)
				result = append(result, edge)
			}
			// Arguments can contain their own statically typed factories, as in
			// Decorate(NewTransformer[A, B](...)). The decorator's supplied input
			// does not establish an endpoint or prove returned-value flow.
			return true
		})
		if failure != nil {
			return nil, failure
		}
		if len(result) == 0 && !acceptsTransformer(fn) {
			return nil, fmt.Errorf("%s: cannot derive transformer endpoints from the typed factory", fn)
		}
		return result, nil
	}
	for fn := range p.functions {
		if !p.implementation(fn.Pkg().Path()) || !fn.Exported() {
			continue
		}
		signature := fn.Type().(*types.Signature)
		if signature.Recv() != nil || signature.TypeParams().Len() != 0 || !returnsTransformer(fn) {
			continue
		}
		edges, err := endpoints(fn, nil, map[*types.Func]bool{})
		if err != nil {
			return err
		}
		var seen []*operation
		for _, edge := range edges {
			var from, to *language
			for i := range m.languages {
				if m.languages[i].id == edge.from {
					from = &m.languages[i]
				}
				if m.languages[i].id == edge.to {
					to = &m.languages[i]
				}
			}
			op := &operation{kind: "transformer", public: fn, entry: edge.chain[len(edge.chain)-1].callee, chain: edge.chain, stop: "typed transformer boundary", contract: contract, from: from, to: to}
			op.id = operationID(op.kind, fn, contract, nil, *from, *to)
			previousIndex := slices.IndexFunc(seen, func(previous *operation) bool { return sameOperation(previous, op) })
			if previousIndex < 0 {
				add(edge.from, edge.to, location{op: op})
				seen = append(seen, op)
			} else {
				previous := seen[previousIndex]
				// Several typed call paths can establish the same candidate edge.
				// Retain their source relationships without treating any as returned
				// value flow or discarding a related declaration's own associations.
				for _, link := range op.chain {
					duplicate := false
					for _, existing := range previous.chain {
						if sameCallLink(link, existing) {
							duplicate = true
							break
						}
					}
					if !duplicate {
						previous.chain = append(previous.chain, link)
					}
				}
			}
		}
	}
	return nil
}
