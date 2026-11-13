package generation

// PROTOTYPE (bittype#196): the planner of the contribution contract. It
// drives each selected contribution through the five steps of the protocol:
// the rule names, the planner resolves, the rule decides, the codec freezes,
// the registered projector projects. A rule never invokes the planner.

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// ContributionInput is one generation request under the contribution
// contract. Products are the products the request names besides native types
// (O5): with none, the request asks for types alone.
type ContributionInput[F, T, R any] struct {
	From          representation.Language
	Input         transform.Input
	Source        SourceProvider[F, R]
	Roots         []ast.NodeRef
	Profile       string
	Contributions []Contribution[F, T, R]
	Codecs        []PayloadCodec
	Facets        []Facet
	Products      []Product
	// Configuration is the target configuration that affects selection, what is
	// admissible natively or whether generated code is compatible (row 35).
	Configuration []Setting
	Execution     Execution[T]
}

// ContributionPlan is the result of planning: the closed native program as a
// runnable generation plan, and the selected representation plan.
type ContributionPlan struct {
	Generation GenerationPlan
	Plan       RepresentationPlan
}

// MemberVocabulary is a target's statement of its member vocabulary (row 13):
// which native expressions are members of a record. A target without it
// admits no member fragment.
type MemberVocabulary[T any] interface {
	CheckMember(context.Context, Expression[T]) error
}

type declState[F, T any] struct {
	id           uint32
	symbol, node string
	formals      []Binder[F]
	convention   GenericConvention[T]
	requirements []PlanRequirement
	facts        []Fact
	reference    ContributionReference
	start        int // the number of open request edges when the declaration began
}

func (s *declState[F, T]) require(binder int, kind EvidenceKind) {
	r := PlanRequirement{binder, kind}
	if !slices.Contains(s.requirements, r) {
		s.requirements = append(s.requirements, r)
	}
}

type cframe[F, T any] struct {
	declaration semantics.EntityID
	state       *declState[F, T]
	binders     map[semantics.EntityID]int
}

// resolved is what the planner resolved for one request: the dependency.
type resolved[T any] struct {
	node      string
	expr      Expression[T]
	fragment  []Expression[T]
	members   []string
	facts     []Fact
	reference ContributionReference
	binder    int // the frame's binder the use is bound to, or -1
	origin    ast.NodeRef
}

type symbolInfo struct{ symbol, node string }

type cplanner[F, T, R any] struct {
	*planner[F, T, R]
	in       ContributionInput[F, T, R]
	target   string
	roles    RoleIssuer[F]
	byDomain map[string][]*Contribution[F, T, R]
	specs    []*Contribution[F, T, R]
	codecs   map[SchemaReference]PayloadCodec
	bindings bool
	decls    map[semantics.EntityID]*declState[F, T]
	owned    map[string]uint32
	symbols  map[uint32]symbolInfo
	selected map[ContributionReference]map[SubjectKind]bool
	b        *planBuilder
	edges    []bool // one per open request: true for a demand for evidence
	path     []string
}

func PlanContributions[F, T, R any](ctx context.Context, in ContributionInput[F, T, R]) (ContributionPlan, error) {
	fail := func(e error) (ContributionPlan, error) { return ContributionPlan{}, e }
	if ctx == nil || representation.IsNil(in.Source) || representation.CheckMarker[F](in.From) != nil || !in.Execution.Limits.Valid() || representation.IsNil(in.Execution.Target) || representation.CheckMarker[T](in.Execution.Language) != nil || in.Execution.Output == "" || in.Execution.Checker.Validate() != nil || !in.Execution.Language.Compatible(in.Execution.Checker.Language()) || len(in.Roots) == 0 || in.Profile == "" {
		return fail(ErrInvalid)
	}
	if err := ctx.Err(); err != nil {
		return fail(err)
	}
	if len(in.Roots) > in.Execution.Limits.References || len(in.Contributions) > in.Execution.Limits.Plans {
		return fail(ErrLimit)
	}
	base := &planner[F, T, R]{ctx: ctx, language: in.From, execution: in.Execution, bindings: map[string][]selectedBinding[F, T, R]{}, recipes: map[string][]selectedRecipe[F, T, R]{}, helpers: map[HelperKey]HelperRecipe[T]{}, declarations: map[semantics.EntityID]uint32{}, helperIDs: map[HelperKey]uint32{}, formals: map[uint32][]Binder[F]{}, conventions: map[uint32]GenericConvention[T]{}, active: map[uint32]bool{}}
	p := &cplanner[F, T, R]{planner: base, in: in, target: string(in.Execution.Language.ID()), byDomain: map[string][]*Contribution[F, T, R]{}, codecs: map[SchemaReference]PayloadCodec{}, bindings: len(in.Products) > 0, decls: map[semantics.EntityID]*declState[F, T]{}, owned: map[string]uint32{}, symbols: map[uint32]symbolInfo{}, selected: map[ContributionReference]map[SubjectKind]bool{}, b: &planBuilder{index: map[string]int{}, lazy: map[string]lazyPayload{}}}
	for _, codec := range in.Codecs {
		if representation.IsNil(codec) || !codec.Schema().Valid() {
			return fail(ErrInvalid)
		}
		if _, dup := p.codecs[codec.Schema()]; dup {
			return fail(fmt.Errorf("%w: two codecs for payload schema %s/%s", ErrAmbiguous, codec.Schema().Identity, codec.Schema().Version))
		}
		p.codecs[codec.Schema()] = codec
	}
	identities := map[ContributionIdentity]bool{}
	for i := range in.Contributions {
		c := &in.Contributions[i]
		if err := c.validate(); err != nil {
			return fail(err)
		}
		if c.Profile != in.Profile || c.Target != p.target {
			return fail(fmt.Errorf("%w: contribution %s states profile %q and target %q", ErrScope, c.Reference, c.Profile, c.Target))
		}
		if identities[c.Reference.Identity] {
			return fail(fmt.Errorf("%w: contribution identity %s/%s twice", ErrInvalid, c.Reference.Identity.Publisher, c.Reference.Identity.Key))
		}
		identities[c.Reference.Identity] = true
		for _, s := range append([]Selector{c.Domain.selector}, c.Domain.actuals...) {
			if s.Valid() && !in.From.Compatible(s.Language()) {
				return fail(ErrInvalid)
			}
		}
		// The codec of a payload is the one registered under the schema that the
		// contribution states. That is how the kernel knows which codec validates.
		if _, ok := p.codecs[c.Schema]; !ok {
			return fail(fmt.Errorf("%w: no codec is registered for payload schema %s/%s of contribution %s", ErrIncomplete, c.Schema.Identity, c.Schema.Version, c.Reference))
		}
		if c.Domain.kind == SpecializationDomain {
			p.specs = append(p.specs, c)
		} else {
			p.byDomain[c.Domain.Key()] = append(p.byDomain[c.Domain.Key()], c)
		}
	}
	for _, f := range in.Facets {
		// A facet names a registered contribution of the same publisher (O8).
		if !f.Valid() || !identities[f.contribution] || f.publisher != f.contribution.Publisher {
			return fail(fmt.Errorf("%w: a facet of kind %s names no registered contribution of its publisher", ErrScope, f.Kind().Identity))
		}
	}
	captured, err := admitInput[F](in.From, in.Source.Contract(), in.Input)
	if err != nil {
		return fail(err)
	}
	p.walk = newWalk(captured)
	source, err := in.Source.Open(ctx, captured)
	if err != nil {
		return fail(err)
	}
	if representation.IsNil(source.Interpretation) {
		return fail(ErrInvalid)
	}
	p.source, p.reader, p.roles = source.Interpretation, source.Reader, source.Roles
	p.b.data = planData{Envelope: EnvelopeVersion, Profile: in.Profile, Products: []string{TypesProduct}, Target: PlanTarget{p.target, slices.Clone(in.Configuration)}, Source: PlanSource{string(in.From.ID()), source.Semantics}}
	for _, product := range in.Products {
		if product.Identity == "" || product.Identity == TypesProduct || slices.Contains(p.b.data.Products, product.Identity) {
			return fail(ErrInvalid)
		}
		p.b.data.Products = append(p.b.data.Products, product.Identity)
	}
	seen := map[ast.NodeRef]bool{}
	for _, site := range in.Roots {
		if seen[site] {
			return fail(ErrScope)
		}
		seen[site] = true
		root, err := p.root(site)
		if err != nil {
			return fail(err)
		}
		st, err := p.declare(root, 0)
		if err != nil {
			return fail(err)
		}
		if p.program.declarations[st.id-1].Role == RequestedDeclaration {
			return fail(ErrScope)
		}
		p.program.exposed = append(p.program.exposed, st.id)
		p.program.declarations[st.id-1].Role = RequestedDeclaration
		p.b.data.Roots = append(p.b.data.Roots, st.symbol)
	}
	if err := p.products(); err != nil {
		return fail(err)
	}
	if err := in.Execution.Target.CheckProgram(ctx, p.program); err != nil {
		return fail(err)
	}
	if err := ctx.Err(); err != nil {
		return fail(err)
	}
	program := p.program
	exec := in.Execution
	owned := make([]planned, len(program.declarations))
	roles := map[DeclarationRole]string{RequestedDeclaration: "requested", DependencyDeclaration: "dependency", HelperDeclaration: "owned"}
	for i, d := range program.declarations {
		info := p.symbols[d.ID]
		owned[i] = planned{role: d.Role, provenance: d.Provenance, values: d.Description.Values, symbol: info.symbol}
		symbol := PlanSymbol{ID: info.symbol, Role: roles[d.Role], Intent: d.Description.Name, Node: info.node}
		for _, site := range d.Provenance {
			symbol.Provenance = append(symbol.Provenance, planSite(site))
		}
		p.b.data.Symbols = append(p.b.data.Symbols, symbol)
	}
	plan := RepresentationPlan{owner: fresh(), data: p.b.data, index: p.b.index, lazy: p.b.lazy}
	return ContributionPlan{Plan: plan, Generation: GenerationPlan{steps: slices.Clone(p.steps), run: func(ctx context.Context) (Output, error) {
		built, err := exec.Target.Build(ctx, BuildInput[T]{program, exec.Language, exec.Output, exec.Limits})
		if err != nil {
			return Output{}, err
		}
		if err = ctx.Err(); err != nil {
			return Output{}, err
		}
		out, err := seal[T](ctx, exec.Checker, built, owned, program.exposed, exec.Limits)
		if err != nil {
			return Output{}, err
		}
		if out.Snapshot().Entry().Tree != exec.Output {
			return Output{}, ErrScope
		}
		return out, nil
	}}}, nil
}

// products is guarantees 1 and 2: every facet that a requested product
// requires belongs to the selected contribution and matches its kind. It also
// writes the contributions that were selected into the plan, with the facets
// that cross a process boundary.
func (p *cplanner[F, T, R]) products() error {
	references := make([]ContributionReference, 0, len(p.selected))
	for r := range p.selected {
		references = append(references, r)
	}
	sort.Slice(references, func(i, j int) bool { return references[i].String() < references[j].String() })
	for _, reference := range references {
		var c *Contribution[F, T, R]
		for i := range p.in.Contributions {
			if p.in.Contributions[i].Reference == reference {
				c = &p.in.Contributions[i]
			}
		}
		entry := PlanContribution{Reference: reference, Target: c.Target, Profile: c.Profile, Domain: c.Domain.Key(), Schema: c.Schema, Facts: planFacts(c.Facts)}
		for _, product := range p.in.Products {
			for _, kind := range []SubjectKind{UseSubject, DeclarationSubject, FieldSubject} {
				if !p.selected[reference][kind] {
					continue
				}
				for _, required := range product.Requires[kind] {
					facet, err := matchFacet(p.in.Facets, required, reference.Identity, p.target)
					if errors.Is(err, ErrMissingFacet) {
						return &Refusal{Class: RefusalMissingFacet, Owner: KernelOwner, Reason: ReasonFacet + ": " + required.Reference().Identity + " for product " + product.Identity, Contribution: reference, Target: p.target, Profile: p.in.Profile}
					}
					if err != nil {
						return err
					}
					if t := required.description.Transport; t != nil && !slices.ContainsFunc(entry.Facets, func(f PlanFacet) bool { return f.Kind == required.Reference() }) {
						encoded, err := t.Encode(facet.value)
						if err != nil {
							return fmt.Errorf("%w: transport of facet kind %s: %v", ErrInvalid, required.Reference().Identity, err)
						}
						entry.Facets = append(entry.Facets, PlanFacet{required.Reference(), t.MediaType, encoded})
					}
				}
			}
		}
		p.b.data.Contributions = append(p.b.data.Contributions, entry)
	}
	return nil
}

func (p *cplanner[F, T, R]) site(o *occurrence[F]) ast.NodeRef {
	if site, ok := p.walk.termSites[o.term]; ok && !p.walk.shared[o.term] {
		return site
	}
	if sites := p.walk.declared(o.declaration); len(sites) > 0 {
		return sites[0]
	}
	return ast.NodeRef{}
}

func (p *cplanner[F, T, R]) refusal(class RefusalClass, owner, reason, subject string, site ast.NodeRef, c *Contribution[F, T, R]) *Refusal {
	r := &Refusal{Class: class, Owner: owner, Reason: reason, Subject: subject, Site: site, Target: p.target, Profile: p.in.Profile, Path: slices.Clone(p.path)}
	if c != nil {
		r.Contribution = c.Reference
	}
	return r
}

// failure classifies what a rule or a projector returned. A refusal keeps its
// class and gains its context. Anything else is an implementation failure.
func (p *cplanner[F, T, R]) failure(c *Contribution[F, T, R], err error, subject string, site ast.NodeRef) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	if r, ok := err.(*Refusal); ok && r.Class >= RefusalUnsupported && r.Class <= RefusalRequirementUnresolved {
		out := *r
		out.Subject, out.Site, out.Contribution, out.Target, out.Profile, out.Path = subject, site, c.Reference, p.target, p.in.Profile, slices.Clone(p.path)
		return &out
	}
	return &RuleFailure{c.Reference, err}
}

// choose selects by the domain. None is "missing contribution" and more than
// one is "ambiguous contribution": two values (O26). Registration order is
// never a preference.
func (p *cplanner[F, T, R]) choose(domain, subject string, site ast.NodeRef) (*Contribution[F, T, R], error) {
	choices := p.byDomain[domain]
	switch len(choices) {
	case 0:
		return nil, p.refusal(RefusalMissingContribution, KernelOwner, ReasonNoContribution+" "+domain, subject, site, nil)
	case 1:
		return choices[0], nil
	}
	return nil, p.refusal(RefusalAmbiguousContribution, KernelOwner, ReasonSeveralContributions+" "+domain, subject, site, nil)
}

func (p *cplanner[F, T, R]) note(c *Contribution[F, T, R], kind SubjectKind) {
	if p.selected[c.Reference] == nil {
		p.selected[c.Reference] = map[SubjectKind]bool{}
	}
	p.selected[c.Reference][kind] = true
}

// useKey is the source identity of a use: the selector of its form or callee.
func (p *cplanner[F, T, R]) useKey(view UseView[F]) (string, error) {
	u := view.Use
	switch view.Kind {
	case BoundUse:
		return "binder", nil
	case DeclarationUse:
		return p.selector(view.Callee, u.o.term, u.o.declaration)
	case ExpressionUse:
		selector, err := p.source.Expression(p.ctx, u)
		if err != nil {
			return "", err
		}
		return p.qualified(selector, u.o.term, u.o.declaration, "expression selector")
	}
	return "", ErrInvalid
}

// use plans one occurrence and records its use entry.
func (p *cplanner[F, T, R]) use(u Use[F], f cframe[F, T], depth int) (resolved[T], error) {
	var zero resolved[T]
	if err := p.tick(depth); err != nil {
		return zero, err
	}
	if !p.walk.owns(u) || u.o.declaration != f.declaration {
		return zero, ErrScope
	}
	if err := p.admit(u.o, depth); err != nil {
		return zero, err
	}
	view, err := p.walk.structure(u.o)
	if err != nil {
		return zero, err
	}
	site := p.site(u.o)
	key, err := p.useKey(view)
	if err != nil {
		return zero, err
	}
	var out resolved[T]
	switch view.Kind {
	case BoundUse:
		i, ok := f.binders[view.Bound.entity]
		if !ok {
			return zero, ErrScope
		}
		node := p.b.node(PlanNode{Kind: "binder", Binder: &PlanBinder{f.state.symbol, i}})
		out = resolved[T]{node: node, expr: Expression[T]{kind: BinderExpression, declaration: f.state.id, index: i}, binder: i}
	case DeclarationUse:
		if spec, name, err := p.specialized(key, view); err != nil {
			return zero, err
		} else if spec != nil {
			r := p.refusal(RefusalUnsupported, KernelOwner, ReasonSpecialization, key, site, spec)
			r.Application = name
			return zero, r
		}
		if len(p.byDomain["exact:"+key]) > 0 {
			c, err := p.choose("exact:"+key, key, site)
			if err != nil {
				return zero, err
			}
			out, err = p.represent(c, Subject[F]{Kind: UseSubject, Use: view.Clone()}, key, site, f, depth)
		} else {
			out, err = p.application(view, key, site, f, depth)
		}
	case ExpressionUse:
		c, cerr := p.choose("form:"+key, key, site)
		if cerr != nil {
			return zero, cerr
		}
		out, err = p.represent(c, Subject[F]{Kind: UseSubject, Use: view.Clone()}, key, site, f, depth)
	default:
		return zero, ErrInvalid
	}
	if err != nil {
		return zero, err
	}
	out.origin = site
	p.b.use(out.node, site, f.state.symbol)
	return out, nil
}

// specialized finds a contribution whose domain is this application.
func (p *cplanner[F, T, R]) specialized(key string, view UseView[F]) (*Contribution[F, T, R], string, error) {
	if len(p.specs) == 0 {
		return nil, "", nil
	}
	actuals := make([]string, len(view.Actuals))
	for i, a := range view.Actuals {
		av, err := p.walk.structure(a.Use.o)
		if err != nil {
			return nil, "", err
		}
		if actuals[i], err = p.useKey(av); err != nil {
			return nil, "", err
		}
	}
	for _, c := range p.specs {
		if c.Domain.selector.Key() != key || len(c.Domain.actuals) != len(actuals) {
			continue
		}
		match := true
		for i, a := range c.Domain.actuals {
			match = match && a.Key() == actuals[i]
		}
		if match {
			return c, key + "<" + strings.Join(actuals, ", ") + ">", nil
		}
	}
	return nil, "", nil
}

type ran[T any] struct {
	links    []PlanLink
	deps     []resolved[T]
	views    []DependencyView[T]
	decision Decision
	frozen   []byte
	view     any
}

func occurrencePath[F any](o *occurrence[F]) string {
	if o == nil {
		return ""
	}
	if o.parent == nil {
		return "body"
	}
	for position, child := range o.parent.children {
		if child == o {
			return fmt.Sprintf("%s.%d", occurrencePath(o.parent), position)
		}
	}
	return occurrencePath(o.parent) + ".?"
}

// run drives steps 1 to 4 for one selection: name, resolve, decide, freeze.
func (p *cplanner[F, T, R]) run(c *Contribution[F, T, R], subject Subject[F], key string, site ast.NodeRef, f cframe[F, T], depth int) (ran[T], error) {
	var out ran[T]
	var requests []Request[F]
	var err error
	if c.Rule.Name != nil {
		if requests, err = c.Rule.Name(p.ctx, NameStep[F, R]{subject, p.reader}); err != nil {
			return out, p.failure(c, err, key, site)
		}
	}
	if len(requests) > p.execution.Limits.References {
		return out, ErrLimit
	}
	// The subject as an answer is bound to it: its place in its declaration.
	// The declaration itself is the generic environment, bound apart.
	subjectKey := "declaration"
	switch subject.Kind {
	case UseSubject:
		subjectKey = occurrencePath(subject.Use.Use.o)
	case FieldSubject:
		subjectKey = occurrencePath(subject.Field.Operand.o) + "@" + subject.Field.Role.key
	}
	state := &answerState{}
	for i, r := range requests {
		if r.kind != helperRequest && !p.walk.owns(r.use) {
			return out, ErrScope
		}
		link := PlanLink{Stated: r.stated}
		var dep resolved[T]
		switch r.kind {
		case childRequest:
			switch subject.Kind {
			case UseSubject:
				if !childOf(r.use, subject.Use.Use) {
					return out, errNotChild
				}
				for _, operand := range subject.Use.Operands {
					if operand.Use == r.use {
						link.Position = slices.Clone(operand.Locator)
					}
				}
				for n, actual := range subject.Use.Actuals {
					if actual.Use == r.use {
						link.Position = []uint32{uint32(n)}
					}
				}
			case DeclarationSubject:
				if r.use != subject.Declaration.Body {
					return out, errNotBody
				}
			case FieldSubject:
				if r.use != subject.Field.Operand {
					return out, errNotChild
				}
			}
			link.Kind, link.Need = "child", "reference"
			if len(r.demands) > 0 {
				link.Need = "evidence"
			}
			if r.conversion {
				link.Need = "conversion"
				if !p.bindings {
					// Types alone: the dependency is not demanded. The link
					// records that the rule named it and that nothing was planned.
					link.Need = "conversion, not demanded"
					state.issued, state.expected = append(state.issued, nil), append(state.expected, answerBinding{})
					out.links, out.deps, out.views = append(out.links, link), append(out.deps, resolved[T]{binder: -1}), append(out.views, DependencyView[T]{})
					continue
				}
			}
			p.edges = append(p.edges, len(r.demands) > 0)
			dep, err = p.use(r.use, f, depth+1)
			p.edges = p.edges[:len(p.edges)-1]
		case fieldRequest:
			// The kernel keeps its scope check (row 12): the operand is a child
			// of the parent use, and the role was issued for that pair.
			if subject.Kind != UseSubject || !PrototypeMutants["field-ignores-parent"] && !childOf(r.use, subject.Use.Use) {
				return out, errNotChild
			}
			if !PrototypeMutants["field-ignores-parent"] && (r.role.owner == nil || r.role.owner != p.roles.owner || r.role.parent != subject.Use.Use.o || r.role.operand != r.use.o) {
				return out, fmt.Errorf("%w: the role selector was not issued for this parent and operand", ErrScope)
			}
			for _, operand := range subject.Use.Operands {
				if operand.Use == r.use {
					link.Position = slices.Clone(operand.Locator)
				}
			}
			fieldKey := "role/" + r.role.key
			fc, cerr := p.choose("role:"+r.role.key, fieldKey, p.site(r.use.o))
			if cerr != nil {
				return out, cerr
			}
			link.Kind, link.Need, link.Role = "field", "reference", r.role.key
			p.edges = append(p.edges, false)
			dep, err = p.represent(fc, Subject[F]{Kind: FieldSubject, Field: FieldPosition[F]{subject.Use.Use, r.use, r.role, r.stated}}, fieldKey, p.site(r.use.o), f, depth+1)
			p.edges = p.edges[:len(p.edges)-1]
		case helperRequest:
			return out, p.refusal(RefusalUnsupported, KernelOwner, ReasonSharedHelper, key, site, c)
		default:
			return out, ErrInvalid
		}
		if err != nil {
			return out, err
		}
		link.Node = dep.node
		expected := answerBinding{asker: c.Reference, subject: subjectKey, dependency: i, reference: dep.reference, environment: f.state.symbol}
		var issued []Answer
		for _, d := range r.demands {
			a := Answer{kind: d.Kind, outcome: Unknown, bound: expected}
			if dep.binder >= 0 {
				// Evidence about an operand bound to a binder comes from the
				// binder's requirements and never from one application (row 7).
				a.outcome = Deferred
				if d.Require {
					a.established = true
					f.state.require(dep.binder, d.Kind)
				}
			} else {
				for _, fact := range dep.facts {
					if fact.Kind == d.Kind {
						a.outcome, a.content = fact.Outcome, fact.Content
						break
					}
				}
			}
			issued = append(issued, a)
			answer := PlanAnswer{Kind: d.Kind, Outcome: a.outcome.String(), Content: a.content, Required: a.established}
			if dep.reference.Valid() {
				reference := dep.reference
				answer.Reference = &reference
			}
			link.Answers = append(link.Answers, answer)
		}
		state.issued, state.expected = append(state.issued, issued), append(state.expected, expected)
		out.links, out.deps = append(out.links, link), append(out.deps, dep)
		out.views = append(out.views, DependencyView[T]{Type: dep.expr, Fragment: slices.Clone(dep.fragment), Stated: r.stated, Origin: dep.origin})
	}
	if out.decision, err = c.Rule.Decide(p.ctx, DecideStep[F, R]{subject, p.reader, state}); err != nil {
		return out, p.failure(c, err, key, site)
	}
	if state.failure != nil {
		return out, state.failure
	}
	for _, fact := range out.decision.Facts {
		if fact.Kind.Identity == "" || fact.Outcome != Proved && fact.Outcome != Refuted && fact.Outcome != Exhausted {
			return out, &RuleFailure{c.Reference, fmt.Errorf("%w: a stated fact needs a kind and one of proved, refuted, exhausted", ErrInvalid)}
		}
	}
	// Step 4. The codec of the stated schema encodes, validates and freezes
	// the value at the moment the rule returns, for every request.
	codec := p.codecs[c.Schema]
	frozen, err := codec.Freeze(out.decision.Payload)
	if err != nil {
		return out, fmt.Errorf("%w: contribution %s: %v", ErrPayload, c.Reference, err)
	}
	out.frozen = slices.Clone(frozen)
	if len(out.frozen) > p.execution.Limits.OutputBytes-p.bytes {
		return out, ErrLimit
	}
	p.bytes += len(out.frozen)
	// Step 5 receives the typed view decoded from the frozen value.
	if out.view, err = codec.View(slices.Clone(out.frozen)); err != nil {
		return out, fmt.Errorf("%w: contribution %s: %v", ErrPayload, c.Reference, err)
	}
	if PrototypeMutants["projector-gets-returned-value"] {
		out.view = out.decision.Payload
	}
	return out, nil
}

func (p *cplanner[F, T, R]) record(kind, key string, c *Contribution[F, T, R], run ran[T]) PlanNode {
	reference, schema := c.Reference, c.Schema
	return PlanNode{Kind: kind, Source: key, Contribution: &reference, Schema: &schema, Payload: hex.EncodeToString(run.frozen), Receipt: receipt(schema, run.frozen), Links: run.links, Facts: planFacts(run.decision.Facts)}
}

func (p *cplanner[F, T, R]) store(node PlanNode, c *Contribution[F, T, R], run ran[T]) string {
	id := p.b.node(node)
	if PrototypeMutants["freeze-at-serialization"] {
		p.b.lazy[id] = lazyPayload{p.codecs[c.Schema], run.decision.Payload}
	}
	return id
}

// represent plans one use or one field position with its selected
// contribution: steps 1 to 4, then the registered projector.
func (p *cplanner[F, T, R]) represent(c *Contribution[F, T, R], subject Subject[F], key string, site ast.NodeRef, f cframe[F, T], depth int) (resolved[T], error) {
	var zero resolved[T]
	p.path = append(p.path, key)
	defer func() { p.path = p.path[:len(p.path)-1] }()
	run, err := p.run(c, subject, key, site, f, depth)
	if err != nil {
		return zero, err
	}
	projected, err := c.Project(p.ctx, Projection[T]{&projectionState[T]{payload: run.view, dependencies: run.views, stated: subject.Field.Stated, origin: site, owned: p.ownedOf(c, depth)}})
	if err != nil {
		return zero, p.failure(c, err, key, site)
	}
	out := resolved[T]{binder: -1, reference: c.Reference, facts: append(slices.Clone(run.decision.Facts), c.Facts...)}
	node := p.record("use", key, c, run)
	if subject.Kind == FieldSubject {
		node.Kind = "field"
		if len(projected.Fragment) == 0 || projected.Type.kind != 0 {
			return zero, &RuleFailure{c.Reference, fmt.Errorf("%w: a field position projects a member fragment", ErrInvalid)}
		}
		vocabulary, ok := any(p.execution.Target).(MemberVocabulary[T])
		if !ok {
			return zero, p.refusal(RefusalUnsupported, KernelOwner, ReasonNoMemberVocabulary, key, site, c)
		}
		for _, member := range projected.Fragment {
			if err := p.check(member, depth+1); err != nil {
				return zero, err
			}
			if err := vocabulary.CheckMember(p.ctx, member); err != nil {
				return zero, p.refusal(RefusalUnsupported, "member vocabulary of "+p.target, ReasonMemberFragment+": "+err.Error(), key, site, c)
			}
		}
		out.fragment = slices.Clone(projected.Fragment)
		for _, m := range projected.Members {
			out.members = append(out.members, m.Name)
		}
	} else {
		if projected.Type.kind == 0 || len(projected.Fragment) != 0 {
			return zero, &RuleFailure{c.Reference, fmt.Errorf("%w: a use projects a native type", ErrInvalid)}
		}
		if err := p.check(projected.Type, depth+1); err != nil {
			return zero, err
		}
		out.expr = projected.Type
		for _, m := range projected.Members {
			if m.Fragment > 0 && m.Fragment <= len(run.deps) {
				node.Members = append(node.Members, run.deps[m.Fragment-1].members...)
			} else {
				node.Members = append(node.Members, m.Name)
			}
		}
	}
	out.node = p.store(node, c, run)
	p.note(c, subject.Kind)
	p.steps = append(p.steps, Step{Binding: c.Reference.String(), Source: key, Form: out.expr.form.Name(), Expression: describe(out.expr)})
	return out, nil
}

// check validates native intent that a projector returned.
func (p *cplanner[F, T, R]) check(e Expression[T], depth int) error {
	if err := p.tick(depth); err != nil {
		return err
	}
	switch e.kind {
	case NativeExpression, DeclarationExpression:
		for _, a := range e.args {
			if err := p.check(a, depth+1); err != nil {
				return err
			}
		}
		if e.kind == NativeExpression {
			return p.execution.Target.CheckExpression(p.ctx, e)
		}
		return nil
	case BinderExpression:
		return nil
	}
	return ErrScope
}

func (p *cplanner[F, T, R]) account(e Expression[T]) error {
	size := len(e.form.name)
	for _, value := range e.data {
		size += len(value)
	}
	if size > p.execution.Limits.OutputBytes-p.bytes {
		return ErrLimit
	}
	p.bytes += size
	for _, a := range e.args {
		if err := p.account(a); err != nil {
			return err
		}
	}
	return nil
}

// ownedOf gives a projector the declarations that its contribution owns. The
// planner issues the symbol and shares the declaration by the owned key: the
// contribution reference and the local name. The spelling is not the key.
func (p *cplanner[F, T, R]) ownedOf(c *Contribution[F, T, R], depth int) func(string, []Expression[T]) (Expression[T], error) {
	return func(name string, args []Expression[T]) (Expression[T], error) {
		var zero Expression[T]
		index := slices.IndexFunc(c.Owned, func(o OwnedDescription[T]) bool { return o.Name == name })
		if index < 0 {
			return zero, ErrScope
		}
		desc := c.Owned[index].Declaration.Clone()
		key := c.Reference.String() + "/" + name
		if PrototypeMutants["owned-shared-by-spelling"] {
			key = desc.Name
		}
		id := p.owned[key]
		if id == 0 {
			var err error
			if id, err = p.reserve(HelperDeclaration, nil); err != nil {
				return zero, err
			}
			p.owned[key] = id
			if len(desc.Parameters) > p.execution.Limits.Formals {
				return zero, ErrLimit
			}
			if err = p.carry(desc.Values); err != nil {
				return zero, err
			}
			if desc.Body, err = p.template(desc.Body, id, len(desc.Parameters), nil, depth+1); err != nil {
				return zero, err
			}
			if err = p.execution.Target.CheckDeclaration(p.ctx, desc.Clone()); err != nil {
				return zero, err
			}
			p.program.declarations[id-1].Description = desc
			p.active[id] = false
			reference := c.Reference
			symbol := "owned:" + key
			p.symbols[id] = symbolInfo{symbol, p.b.node(PlanNode{Kind: "owned", Contribution: &reference, Symbol: symbol})}
			p.steps = append(p.steps, Step{Binding: "owned/" + key, Declaration: desc.Name, Form: desc.Form.Name(), Expression: describe(desc.Body)})
		}
		if len(args) != len(p.program.declarations[id-1].Description.Parameters) {
			return zero, ErrScope
		}
		return Expression[T]{kind: DeclarationExpression, declaration: id, args: slices.Clone(args)}, nil
	}
}

func neededArguments[T any](convention GenericConvention[T]) map[int]bool {
	needed := map[int]bool{}
	var visit func(Expression[T])
	visit = func(e Expression[T]) {
		if e.kind == ArgumentExpression {
			needed[e.index] = true
		}
		for _, a := range e.args {
			visit(a)
		}
	}
	for _, arg := range convention.Arguments {
		visit(arg)
	}
	return needed
}

// application plans a use of a declaration that the declaration policy
// represents. The open declaration is primary (row 20): it is planned once,
// and the application links its binders to the arguments of this use.
func (p *cplanner[F, T, R]) application(view UseView[F], key string, site ast.NodeRef, f cframe[F, T], depth int) (resolved[T], error) {
	var zero resolved[T]
	p.path = append(p.path, key)
	defer func() { p.path = p.path[:len(p.path)-1] }()
	st, err := p.declare(view.Callee, depth+1)
	if err != nil {
		return zero, err
	}
	if len(view.Actuals) != len(st.formals) {
		return zero, ErrScope
	}
	needed := neededArguments(st.convention)
	required := map[int]bool{}
	for _, r := range st.requirements {
		required[r.Binder] = true
	}
	actuals := make([]Expression[T], len(st.formals))
	plannedActuals := make([]*resolved[T], len(st.formals))
	app := &PlanApplication{Declaration: st.node}
	for i, binder := range st.formals {
		a, ok := view.Actual(binder)
		if !ok {
			return zero, ErrScope
		}
		// For types alone only demanded native actuals need mappings. When
		// bindings are asked for, an argument that the native convention omits
		// is still represented, because its adapter is needed (row 20, O5).
		adapters := p.bindings && !PrototypeMutants["application-drops-omitted-argument"]
		if !needed[i] && !adapters && !required[i] {
			continue
		}
		p.edges = append(p.edges, false)
		d, err := p.use(a, f, depth+1)
		p.edges = p.edges[:len(p.edges)-1]
		if err != nil {
			return zero, err
		}
		plannedActuals[i], actuals[i] = &d, d.expr
		if needed[i] {
			app.Native = append(app.Native, PlanArgument{i, d.node})
		}
		if adapters {
			app.Adapters = append(app.Adapters, PlanArgument{i, d.node})
		}
	}
	args := make([]Expression[T], len(st.convention.Arguments))
	for i, arg := range st.convention.Arguments {
		if args[i], err = p.template(arg, 0, 0, actuals, depth+1); err != nil {
			return zero, err
		}
	}
	// Each application discharges the requirements of the binders (row 8).
	for _, r := range st.requirements {
		d := plannedActuals[r.Binder]
		if d.binder >= 0 {
			// The actual is itself bound to a binder: the requirement moves to
			// that binder, and the applications of this declaration discharge it.
			f.state.require(d.binder, r.Kind)
			app.Discharges = append(app.Discharges, PlanDischarge{Binder: r.Binder, Kind: r.Kind, Outcome: Deferred.String()})
			continue
		}
		outcome := Unknown
		for _, fact := range d.facts {
			if fact.Kind == r.Kind {
				outcome = fact.Outcome
				break
			}
		}
		class := RefusalRequirementUnresolved
		switch outcome {
		case Proved:
			reference := d.reference
			app.Discharges = append(app.Discharges, PlanDischarge{Binder: r.Binder, Kind: r.Kind, Outcome: Proved.String(), Reference: &reference})
			continue
		case Refuted:
			class = RefusalRequirementUnsatisfied
		}
		refusal := p.refusal(class, KernelOwner, fmt.Sprintf("%s: binder %d requires %s/%s and the evidence is %s", ReasonRequirement, r.Binder, r.Kind.Identity, r.Kind.Version, outcome), key, site, nil)
		refusal.Requirement, refusal.Contribution = r.Kind, d.reference
		return zero, refusal
	}
	node := p.b.node(PlanNode{Kind: "application", Source: key, Symbol: st.symbol, Application: app})
	return resolved[T]{node: node, expr: Expression[T]{kind: DeclarationExpression, declaration: st.id, args: args}, facts: st.facts, reference: st.reference, binder: -1}, nil
}

// substitute replaces the binder placeholders of one declaration by the native
// meaning that the declaration's convention gives its binders.
func substitute[T any](e Expression[T], scope uint32, bindings []Expression[T]) Expression[T] {
	if e.kind == BinderExpression {
		if e.declaration == scope && e.index >= 0 && e.index < len(bindings) {
			return bindings[e.index]
		}
		return e
	}
	if len(e.args) > 0 {
		args := make([]Expression[T], len(e.args))
		for i, a := range e.args {
			args[i] = substitute(a, scope, bindings)
		}
		e.args = args
	}
	return e
}

// declare plans the open declaration of a source declaration, once.
func (p *cplanner[F, T, R]) declare(d Declaration[F], depth int) (*declState[F, T], error) {
	if err := p.tick(depth); err != nil {
		return nil, err
	}
	if st := p.decls[d.entity]; st != nil {
		if p.active[st.id] {
			// Two refusals that can be told apart (row 9). The cycle is an
			// evidence cycle when every open request on it demands evidence:
			// then no decision on it can finish. With one reference on it, it
			// is a recursive reference, which placeholders would admit (#164).
			open := p.edges[st.start:]
			reason := ReasonEvidenceCycle
			if len(open) == 0 || slices.Contains(open, false) {
				reason = ReasonRecursiveReference
			}
			return nil, p.refusal(RefusalUnsupported, KernelOwner, reason, strings.TrimPrefix(st.symbol, "source:"), ast.NodeRef{}, nil)
		}
		return st, nil
	}
	key, err := p.selector(d, 0, d.entity)
	if err != nil {
		return nil, err
	}
	var site ast.NodeRef
	if sites := p.walk.declared(d.entity); len(sites) > 0 {
		site = sites[0]
	}
	if len(p.byDomain["exact:"+key]) > 0 {
		// An exact mapping never silently becomes a structural publication.
		return nil, p.refusal(RefusalUnsupported, KernelOwner, ReasonExactWithoutPolicy, key, site, p.byDomain["exact:"+key][0])
	}
	c, err := p.choose("policy", key, site)
	if err != nil {
		return nil, err
	}
	view, err := p.walk.declaration(d)
	if err != nil {
		return nil, err
	}
	if len(view.Formals) > p.execution.Limits.Formals {
		return nil, ErrLimit
	}
	j, err := p.source.JudgeDeclaration(p.ctx, d)
	if err == nil {
		err = verdict(j)
	}
	if err == nil {
		err = p.admit(view.Body.o, depth+1)
	}
	if err != nil {
		return nil, err
	}
	id, err := p.reserve(DependencyDeclaration, p.walk.declared(d.entity))
	if err != nil {
		return nil, err
	}
	st := &declState[F, T]{id: id, symbol: "source:" + key, formals: view.Formals, start: len(p.edges), reference: c.Reference}
	p.decls[d.entity] = st
	f := cframe[F, T]{declaration: d.entity, state: st, binders: map[semantics.EntityID]int{}}
	var origins []ast.NodeRef
	for i, binder := range view.Formals {
		if _, dup := f.binders[binder.entity]; dup {
			return nil, ErrScope
		}
		f.binders[binder.entity] = i
		var origin ast.NodeRef
		if sites := p.walk.declared(binder.entity); len(sites) > 0 {
			origin = sites[0]
		}
		origins = append(origins, origin)
	}
	p.path = append(p.path, key)
	defer func() { p.path = p.path[:len(p.path)-1] }()
	run, err := p.run(c, Subject[F]{Kind: DeclarationSubject, Declaration: view.Clone()}, key, site, f, depth)
	if err != nil {
		return nil, err
	}
	desc, err := c.Declare(p.ctx, Projection[T]{&projectionState[T]{payload: run.view, dependencies: run.views, origin: site, binders: origins, owned: p.ownedOf(c, depth)}})
	if err != nil {
		return nil, p.failure(c, err, key, site)
	}
	desc = desc.Clone()
	if len(desc.Generic.Bindings) != len(view.Formals) || len(desc.Generic.Arguments) != len(desc.Generic.Parameters) {
		return nil, ErrScope
	}
	if err = p.carry(desc.Values); err != nil {
		return nil, err
	}
	if len(desc.Generic.Parameters) > p.execution.Limits.Formals {
		return nil, ErrLimit
	}
	bindings := make([]Expression[T], len(view.Formals))
	for i := range view.Formals {
		if bindings[i], err = p.template(desc.Generic.Bindings[i], id, len(desc.Generic.Parameters), nil, depth+1); err != nil {
			return nil, err
		}
	}
	for _, arg := range desc.Generic.Arguments {
		if err = p.convention(arg, len(view.Formals), depth+1); err != nil {
			return nil, err
		}
	}
	if err = p.check(desc.Body, depth+1); err != nil {
		return nil, err
	}
	if err = p.account(desc.Body); err != nil {
		return nil, err
	}
	desc.Body = substitute(desc.Body, id, bindings)
	if err = checkParameterUses(id, desc); err != nil {
		return nil, err
	}
	if err = p.execution.Target.CheckDeclaration(p.ctx, nativeDeclaration(desc)); err != nil {
		return nil, err
	}
	p.program.declarations[id-1].Description = nativeDeclaration(desc)
	st.convention = desc.Generic
	st.facts = append(slices.Clone(run.decision.Facts), c.Facts...)
	node := p.record("declaration", key, c, run)
	node.Symbol, node.Binders, node.Requirements = st.symbol, len(view.Formals), slices.Clone(st.requirements)
	st.node = p.store(node, c, run)
	p.symbols[id] = symbolInfo{st.symbol, st.node}
	p.active[id] = false
	p.note(c, DeclarationSubject)
	p.steps = append(p.steps, Step{Binding: c.Reference.String(), Source: key, Declaration: desc.Name, Form: desc.Form.Name(), Expression: describe(desc.Body)})
	return st, nil
}
